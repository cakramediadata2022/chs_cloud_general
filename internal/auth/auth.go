package auth

import (
	"chs/internal/config"
	"chs/internal/db_var"
	"chs/internal/general"
	"chs/internal/global_query"
	"chs/internal/global_var"
	"chs/internal/master_data"
	"chs/internal/utils/cache"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetUserAccess() {

}

func GenerateJWT(UserCode string, isApp bool) (string, string, error) {
	Token := jwt.New(jwt.SigningMethodHS256)

	Claims := Token.Claims.(jwt.MapClaims)
	Claims["refresh"] = false
	Claims["user"] = UserCode
	Claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	if isApp {
		Claims["exp"] = time.Now().Add((time.Hour * 24) * 3600).Unix()
	}

	TokenString, Err := Token.SignedString(global_var.SigningKey)

	if Err != nil {
		TokenString = ""
	} else {
		RefreshToken := jwt.New(jwt.SigningMethodHS256)

		Claims := RefreshToken.Claims.(jwt.MapClaims)
		Claims["refresh"] = true
		Claims["user"] = UserCode
		Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
		if isApp {
			Claims["exp"] = time.Now().Add((time.Hour * 24) * 3700).Unix()
		}

		RefreshTokenString, Err := RefreshToken.SignedString(global_var.SigningKey)
		if Err != nil {
			TokenString = ""
			RefreshTokenString = ""
		}
		return TokenString, RefreshTokenString, Err
	}
	return "", "", Err
}

type dataUserStruct struct {
	db_var.User
	UserAccessLevelCode int
}

func GenerateAppTokenP(c *gin.Context) {
	AppCode := c.Param("AppCode")
	NewToken, RefreshToken, Err := GenerateJWT(strings.ToUpper(AppCode), true)
	if Err != nil {
		master_data.SendResponse(global_var.ResponseCode.NotAuthorized, "", nil, c)
		return
	}

	master_data.SendResponse(global_var.ResponseCode.Successfully, "", gin.H{
		"token":         NewToken,
		"refresh_token": RefreshToken,
	}, c)
}

func Login(c *gin.Context) {
	var user dataUserStruct
	type DataInputStruct struct {
		Code     string `json:"Code" binding:"required"`
		Password string `json:"Password" binding:"required"`
		Shift    string `json:"Shift"`
	}
	var DataInput DataInputStruct
	err := c.Bind(&DataInput)
	if err != nil {
		master_data.SendResponse(global_var.ResponseCode.InvalidDataFormat, err.Error(), nil, c)
		return
	}

	if LoginRateLimit(c, DataInput.Code, false) {
		master_data.SendResponse(global_var.ResponseCode.OtherResult, "Too many failed attempts", nil, c)
		return
	}
	// Get Program Configuration
	val, exist := c.Get("pConfig")
	if !exist {
		master_data.SendResponse(global_var.ResponseCode.DatabaseError, "pConfig", nil, c)
		return
	}
	pConfig := val.(*config.CompanyDataConfiguration)
	DB := pConfig.DB

	AuditDate := global_query.GetAuditDate(c, DB, true)
	Result := DB.Table(db_var.TableName.User).Select("user.*", "user_group_access.user_access_level_code").
		Joins("INNER JOIN user_group_access ON user.user_group_access_code=user_group_access.code").
		Where("user.code = ? AND user.password = ?", strings.ToUpper(DataInput.Code), general.GetMD5Hash(DataInput.Password)).Find(&user)

	if Result.RowsAffected > 0 {
		NewToken, RefreshToken, Err := GenerateJWT(strings.ToUpper(DataInput.Code), false)
		if Err != nil {
			master_data.SendResponse(global_var.ResponseCode.NotAuthorized, "", nil, c)
			return
		} else {
			UserInfo, err := SetUserInfo(c, DB, pConfig.CompanyCode, user, DataInput.Shift)
			if err != nil {
				master_data.SendResponse(global_var.ResponseCode.DatabaseError, "", nil, c)
				return
			}
			UserAccess, err := GetUserFormAccess(DB, user.UserGroupAccessCode)
			if err != nil {
				master_data.SendResponse(global_var.ResponseCode.DatabaseError, "", nil, c)
				return
			}
			key := fmt.Sprintf("loginRateLimit:%s", DataInput.Code)
			cache.DataCache.Del(c, c.GetString("UnitCode"), key)
			global_query.InsertLogUser(DB, "", global_var.LogUserAction.Login, AuditDate, c.ClientIP(), c.Request.Host, "", DataInput.Code, "", "", "Login Success", "SYSTEM")
			master_data.SendResponse(global_var.ResponseCode.Successfully, "", gin.H{
				"UserInfo":     UserInfo,
				"UserAccess":   UserAccess,
				"NewToken":     NewToken,
				"RefreshToken": RefreshToken}, c)
		}
	} else {
		LoginRateLimit(c, DataInput.Code, true)
		global_query.InsertLogUser(DB, "", global_var.LogUserAction.LoginDenied, AuditDate, c.ClientIP(), c.Request.Host, "", DataInput.Code, "", "", "Login Failed", "SYSTEM")
		master_data.SendResponse(global_var.ResponseCode.NotAuthorized, "", nil, c)
	}
}

var ctx = context.Background()

func LoginRateLimit(c *gin.Context, UserID string, IsFailed bool) bool {
	key := fmt.Sprintf("loginRateLimit:%s", UserID)
	fmt.Println(UserID)
	times, err := cache.DataCache.Get(ctx, c.GetString("UnitCode"), key)
	if err != nil {
		if IsFailed {
			cache.DataCache.Set(ctx, c.GetString("UnitCode"), key, 1, 1*time.Minute)
		}
		return false
	}

	timeInt := general.StrToInt(string(times))

	if IsFailed {
		timeInt++
		cache.DataCache.Set(ctx, c.GetString("UnitCode"), key, timeInt, 1*time.Minute)
	}
	return timeInt > 3

}

func SetUserInfo(c *gin.Context, DB *gorm.DB, CompanyCode string, UserData dataUserStruct, Shift string) (global_var.TUserInfo, error) {
	IsUseSessionID := false
	clientIP := c.ClientIP()
	if clientIP == "" {
		clientIP = c.Request.Header.Get("X-Forwarded-For")
	}
	if clientIP == "" {
		clientIP = strings.Split(c.Request.RemoteAddr, ":")[0]
	}
	if net.ParseIP(clientIP) == nil {
		clientIP = "unknown"
	}
	var ShiftID uint64
	err := DB.Table(db_var.TableName.LogShift).Select("id").Where("created_by=?", UserData.Code).Where("is_open=1").Limit(1).Scan(&ShiftID).Error
	if ShiftID <= 0 {
		ShiftID, err = global_query.InsertLogShift(DB, UserData.Code, Shift, time.Now(), global_query.GetAuditDate(c, DB, false), 0, "", clientIP, "", "")
	}
	UserInfo := global_var.TUserInfo{
		ID:           UserData.Code,
		GroupCode:    UserData.UserGroupAccessCode,
		Name:         UserData.Name,
		AccessLevel:  UserData.UserAccessLevelCode,
		CompanyCode:  CompanyCode,
		LogShiftID:   ShiftID,
		WorkingShift: Shift,
		SessionID:    GenerateSessionID(),
	}
	ID := UserInfo.ID
	if IsUseSessionID {
		ID = UserInfo.SessionID
	}
	cache.DataCache.Set(c, CompanyCode, "USER_INFO_"+ID, UserInfo, 0)
	// cookie := &http.Cookie{
	// 	Name:     "session_id",
	// 	Value:    global_var.UserInfo.SessionID,
	// 	MaxAge:   3600, // cookie will expire in 1 hour
	// 	Path:     "/",
	// 	HttpOnly: true,                     // cookie can only be accessed through HTTP(S)
	// 	Secure:   true,                     // cookie can only be transmitted over HTTPS
	// 	SameSite: http.SameSiteDefaultMode, // cookie will only be sent in first-party context
	// }
	// http.SetCookie(c.Writer, cookie)
	c.SetSameSite(http.SameSiteDefaultMode)
	c.SetCookie("_app_cakrasoft_session_id", UserInfo.SessionID, 2*(3600*24), "/", global_var.Config.Server.Domain, global_var.Config.Server.SSL, false)
	return UserInfo, err
}

func LoadUserInfo(c *gin.Context, Auth string, CompanyCode string) (UserInfo global_var.TUserInfo, err error) {
	IsUseSessionID := false
	ID := strings.TrimSpace(Auth)
	if IsUseSessionID {
		// Get session ID from cookie
		ID, err = c.Cookie("_app_cakrasoft_session_id")
		if err != nil {
			fmt.Println(err.Error())
			return global_var.TUserInfo{}, err
		}
	}
	// Lookup user ID from Redis using session ID
	b, err := cache.DataCache.Get(c, CompanyCode, "USER_INFO_"+ID)
	if err != nil {
		fmt.Println(err.Error())
		return global_var.TUserInfo{}, err
	}

	if b != nil {
		err = json.Unmarshal(b, &UserInfo)
		if err != nil {
			fmt.Println(err.Error())
			return global_var.TUserInfo{}, err
		}
	}
	return UserInfo, nil
}

func LogoutAllUser(c *gin.Context, DB *gorm.DB, CompanyCode string) (err error) {
	var Code []string
	UserID := c.GetString("ValidUserCode")
	if err := DB.Table(db_var.TableName.User).Select("code").Scan(&Code).Error; err != nil {
		return err
	}
	if err := DB.Table(db_var.TableName.LogShift).Where("is_open=1").Updates(&map[string]interface{}{
		"is_open":    "0",
		"updated_by": UserID,
	}).Error; err != nil {
		return err
	}

	for _, code := range Code {
		go cache.DataCache.Del(c, CompanyCode, "USER_INFO_"+code)
	}
	return nil
}

func GetRefreshToken(c *gin.Context) {
	Token := c.GetHeader("Token")
	if Token != "" {
		TokenCheck, Err := jwt.Parse(Token, func(TokenCheck *jwt.Token) (interface{}, error) {
			if _, ok := TokenCheck.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Check Token Error")
			}
			return global_var.SigningKey, nil
		})

		if Err == nil {
			if claims, ok := TokenCheck.Claims.(jwt.MapClaims); ok && TokenCheck.Valid {
				if claims["refresh"].(bool) && claims["user"] != "" {
					Token, RefreshToken, Err := GenerateJWT(claims["user"].(string), false)
					if Err != nil {
						master_data.SendResponse(global_var.ResponseCode.ErrorCreateToken, "", nil, c)
						return
					} else {
						master_data.SendResponse(global_var.ResponseCode.Successfully, "", gin.H{
							"NewToken":     Token,
							"RefreshToken": RefreshToken}, c)
						return
					}
				}
			}
		}
	}
	master_data.SendResponse(global_var.ResponseCode.NotAuthorized, "", nil, c)
}

func GenerateSessionID() string {
	id := uuid.New()
	return id.String()
}

func GetUserFormAccess(DB *gorm.DB, UserGroupAccessCode string) (interface{}, error) {
	type DataOutputStruct struct {
		Code            string                    `json:"code"`
		Id              uint64                    `json:"id"`
		AccessLevelCode int                       `json:"access_level_code"`
		General         db_var.General_user_group `json:"general" gorm:"embedded"`
		Hotel           db_var.User_group         `json:"hotel" gorm:"embedded"`
		PointOfSales    db_var.Pos_user_group     `json:"point_of_sales" gorm:"embedded"`
		Banquet         db_var.Ban_user_group     `json:"banquet" gorm:"embedded"`
		Accounting      db_var.Acc_user_group     `json:"accounting" gorm:"embedded"`
		Asset           db_var.Ast_user_group     `json:"asset" gorm:"embedded"`
		Report          db_var.Report_user_group  `json:"report" gorm:"embedded"`
		Tools           db_var.Tools_user_group   `json:"tools" gorm:"embedded"`
	}
	var DataOutput DataOutputStruct
	err := DB.Table(db_var.TableName.UserGroupAccess).Select(
		"user_group_access.code",
		"user_group_access.id",
		"user_group_access.user_access_level_code as access_level_code",
		"general_user_group.*",
		"user_group.*",
		"pos_user_group.*",
		"ban_user_group.*",
		"acc_user_group.*",
		"ast_user_group.*",
		"report_user_group.*",
		"tools_user_group.*").
		Joins("LEFT JOIN general_user_group ON user_group_access.general_user_group_id=general_user_group.id").
		Joins("LEFT JOIN user_group ON user_group_access.user_group_id=user_group.id").
		Joins("LEFT JOIN pos_user_group ON user_group_access.pos_user_group_id=pos_user_group.id").
		Joins("LEFT JOIN acc_user_group ON user_group_access.acc_user_group_id=acc_user_group.id").
		Joins("LEFT JOIN ast_user_group ON user_group_access.ast_user_group_id=ast_user_group.id").
		Joins("LEFT JOIN ban_user_group ON user_group_access.ban_user_group_id=ban_user_group.id").
		Joins("LEFT JOIN report_user_group ON user_group_access.report_user_group_id=report_user_group.id").
		Joins("LEFT JOIN tools_user_group ON user_group_access.tools_user_group_id=tools_user_group.id").
		Where("user_group_access.code=?", UserGroupAccessCode).
		Take(&DataOutput).Error

	return DataOutput, err
}

func GetUserFormAccessP(c *gin.Context) {
	// Get Program Configuration
	val, exist := c.Get("pConfig")
	if !exist {
		master_data.SendResponse(global_var.ResponseCode.DatabaseError, "pConfig", nil, c)
		return
	}
	pConfig := val.(*config.CompanyDataConfiguration)
	DB := pConfig.DB

	UserAccess, err := GetUserFormAccess(DB, global_var.UserInfo.GroupCode)
	if err != nil {
		master_data.SendResponse(global_var.ResponseCode.DatabaseError, "", nil, c)
		return
	}

	master_data.SendResponse(global_var.ResponseCode.Successfully, "", UserAccess, c)
}

func GetUserAccessString(DB *gorm.DB, UserID, Password, AccessType string, SystemCode string, IsInternalAccess bool) (string, error) {
	FieldAccess := "user.password"
	if SystemCode == global_var.SystemCode.Hotel {
		if AccessType == global_var.UserAccessType.Form {
			FieldAccess = "user_group.access_form"
		} else if AccessType == global_var.UserAccessType.Special {
			FieldAccess = "user_group.access_special"
		} else if AccessType == global_var.UserAccessType.Keylock {
			FieldAccess = "user_group.access_keylock"
		} else if AccessType == global_var.UserAccessType.Reservation {
			FieldAccess = "user_group.access_reservation"
		} else if AccessType == global_var.UserAccessType.Deposit {
			FieldAccess = "user_group.access_deposit"
		} else if AccessType == global_var.UserAccessType.InHouse {
			FieldAccess = "user_group.access_in_house"
		} else if AccessType == global_var.UserAccessType.WalkIn {
			FieldAccess = "user_group.access_walk_in"
		} else if AccessType == global_var.UserAccessType.Folio {
			FieldAccess = "user_group.access_folio"
		} else if AccessType == global_var.UserAccessType.FolioHistory {
			FieldAccess = "user_group.access_folio_history"
		} else if AccessType == global_var.UserAccessType.FloorPlan {
			FieldAccess = "user_group.access_floor_plan"
		} else if AccessType == global_var.UserAccessType.Company {
			FieldAccess = "user_group.access_company"
		} else if AccessType == global_var.UserAccessType.Invoice {
			FieldAccess = "user_group.access_invoice"
		} else if AccessType == global_var.UserAccessType.MemberVoucherGift {
			FieldAccess = "user_group.access_member_voucher_gift"
		} else if AccessType == global_var.UserAccessType.PreviewReport {
			FieldAccess = "user_group.access_preview_report"
		} else if AccessType == global_var.UserAccessType.PaymentByAPAR {
			FieldAccess = "user_group.access_payment_by_ap_ar"
		}
	} else if SystemCode == global_var.SystemCode.Pos {
		if AccessType == global_var.UserAccessType.Special {
			FieldAccess = "pos_user_group.access_special"
		}
	} else if SystemCode == global_var.SystemCode.Accounting {
		if AccessType == global_var.UserAccessType.Special {
			FieldAccess = "acc_user_group.access_special"
		}
	} else if SystemCode == global_var.SystemCode.Report {
		if AccessType == global_var.ReportAccessType.FrontDesk {
			FieldAccess = "report_user_group.access_fo_report"
		}
		if AccessType == global_var.ReportAccessType.Pos {
			FieldAccess = "report_user_group.access_pos_report"
		}
		if AccessType == global_var.ReportAccessType.Banquet {
			FieldAccess = "report_user_group.access_ban_report"
		}
		if AccessType == global_var.ReportAccessType.Accounting {
			FieldAccess = "report_user_group.access_acc_report"
		}
		if AccessType == global_var.ReportAccessType.Asset {
			FieldAccess = "report_user_group.access_ast_report"
		}
		if AccessType == global_var.ReportAccessType.Form {
			FieldAccess = "report_user_group.access_form"
		}
		if AccessType == global_var.ReportAccessType.Preview {
			FieldAccess = "report_user_group.access_preview_report"
		}
	} else {
		return "", errors.New("Invalid SystemCode")
	}

	type UserStruct struct {
		Data      string
		GroupCode string
	}
	var User UserStruct
	Query := DB.Table(db_var.TableName.User).Select("user_group_access.code AS GroupCode, IFNULL("+FieldAccess+",'') AS Data").
		Joins("LEFT JOIN user_group_access ON (user.user_group_access_code = user_group_access.code)").
		Joins("LEFT JOIN user_group ON (user_group_access.user_group_id = user_group.id)").
		Joins("LEFT JOIN pos_user_group ON (user_group_access.pos_user_group_id = pos_user_group.id)").
		Joins("LEFT JOIN report_user_group ON (user_group_access.report_user_group_id = report_user_group.id)").
		// Joins("INNER JOIN ban_user_group ON (user_group_access.ban_user_group_id = ban_user_group.id)").
		Joins("LEFT JOIN acc_user_group ON (user_group_access.acc_user_group_id = acc_user_group.id)").
		Joins("LEFT JOIN ast_user_group ON (user_group_access.ast_user_group_id = ast_user_group.id)").
		// Joins("INNER JOIN pyr_user_group ON (user_group_access.pyr_user_group_id = pyr_user_group.id)").
		// Joins("INNER JOIN cor_user_group ON (user_group_access.cor_user_group_id = cor_user_group.id)").
		Where("user.code=?", UserID)
	if !IsInternalAccess {
		Query.Where("user.password=?", Password)
	}
	Query.Limit(1).
		Scan(&User)

	Result2 := "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
	var err error
	fmt.Println("data", User.Data)
	if User.Data != "" {
		Result2, err = general.OpensslDecrypt(User.Data, User.GroupCode)
		fmt.Println("User", User)
		fmt.Println("Result2", Result2)
		if err != nil {
			Result2 := ""
			return Result2, err
		}
	}
	return Result2, err
}

func CanUserAccess(DB *gorm.DB, UserID, Password, UserAccessTypeCode string, AccessMode byte, IsPasswordEncrypted bool, SystemCode string) bool {
	Result := false
	PasswordEncrypted := Password
	if !IsPasswordEncrypted {
		PasswordEncrypted = general.GetMD5Hash(Password)
	}
	AccessString, _ := GetUserAccessString(DB, UserID, PasswordEncrypted, UserAccessTypeCode, SystemCode, false)
	if AccessString != "" {
		if len(AccessString) > int(AccessMode) {
			Result = string(AccessString[AccessMode]) == "1"
		}
	}
	return Result
}

func CanUserAccessP(c *gin.Context) {
	type DataInputStruct struct {
		UserID              string `json:"user_id" binding:"required"`
		Password            string `json:"password" binding:"required"`
		UserAccessTypeCode  string `json:"user_access_type_code" binding:"required"`
		AccessMode          []byte `json:"access_mode"`
		IsPasswordEncrypted bool   `json:"is_encrypted"`
		SystemCode          string `json:"system_code" binding:"required"` //Module
	}

	var DataInput DataInputStruct
	err := c.BindJSON(&DataInput)
	if err != nil {
		master_data.SendResponse(global_var.ResponseCode.InvalidDataFormat, "", nil, c)
	} else {
		// Get Program Configuration
		val, exist := c.Get("pConfig")
		if !exist {
			master_data.SendResponse(global_var.ResponseCode.DatabaseError, "pConfig", nil, c)
			return
		}
		pConfig := val.(*config.CompanyDataConfiguration)
		DB := pConfig.DB
		IsAuthorized := false
		for _, access := range DataInput.AccessMode {
			canAccess := CanUserAccess(DB, DataInput.UserID, DataInput.Password, DataInput.UserAccessTypeCode, access, DataInput.IsPasswordEncrypted, DataInput.SystemCode)
			IsAuthorized = canAccess
			if !canAccess {
				break
			}
		}
		// IsAuthorized := CanUserAccess(DB, DataInput.UserID, DataInput.Password, DataInput.UserAccessTypeCode, DataInput.AccessMode, DataInput.IsPasswordEncrypted, DataInput.SystemCode)
		master_data.SendResponse(global_var.ResponseCode.Successfully, "", IsAuthorized, c)
	}
}

func ChangeUserPasswordP(c *gin.Context) {
	type DataInputStruct struct {
		OldPassword     string `json:"old_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required"`
		ConfirmPassword string `json:"confirm_password" binding:"required"`
	}
	var DataInput DataInputStruct
	// fmt.Println("11")
	err := c.BindJSON(&DataInput)
	if err != nil {
		fmt.Println(err.Error())
		errMsg := general.GenerateValidateErrorMsg(c, err)
		master_data.SendResponse(global_var.ResponseCode.InvalidDataFormat, errMsg, nil, c)
		return
	}
	// fmt.Println("13")
	// Get Program Configuration
	val, exist := c.Get("pConfig")
	if !exist {
		master_data.SendResponse(global_var.ResponseCode.DatabaseError, "pConfig", nil, c)
		return
	}
	pConfig := val.(*config.CompanyDataConfiguration)
	DB := pConfig.DB
	// fmt.Println("1")
	if DataInput.NewPassword != DataInput.ConfirmPassword {
		master_data.SendResponse(global_var.ResponseCode.OtherResult, "Password doesn't match", nil, c)
		return
	}
	UserID := c.GetString("ValidUserCode")

	// fmt.Println("2")
	var Code string
	DB.Table(db_var.TableName.User).Select("user.code").
		Where("user.code = ? AND user.password = ?", strings.ToUpper(UserID), general.GetMD5Hash(DataInput.OldPassword)).
		Take(&Code)

	if Code == "" {
		master_data.SendResponse(global_var.ResponseCode.OtherResult, "Invalid Password", nil, c)
		return
	}

	// fmt.Println("3")
	if err := DB.Table(db_var.TableName.User).Where("code=?", UserID).Updates(map[string]interface{}{
		"password":   general.GetMD5Hash(DataInput.ConfirmPassword),
		"updated_by": UserID,
	}).Error; err != nil {
		master_data.SendResponse(global_var.ResponseCode.DatabaseError, "1001", nil, c)
		return
	}
	// fmt.Println("5")
	master_data.SendResponse(global_var.ResponseCode.Successfully, "", nil, c)
}
