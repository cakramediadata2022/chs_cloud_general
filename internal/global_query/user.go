package global_query

import (
	"chs_cloud_general/internal/config"
	"chs_cloud_general/internal/db_var"
	"chs_cloud_general/internal/general"
	General "chs_cloud_general/internal/general"
	"chs_cloud_general/internal/global_var"
	"chs_cloud_general/internal/master_data"
	"chs_cloud_general/internal/utils/cache"
	"chs_cloud_general/pkg/utils"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"gorm.io/gorm"
)

func loadUserInfo(c *gin.Context, Auth string, CompanyCode string) (UserInfo global_var.TUserInfo, err error) {
	IsUseSessionID := false
	ID := strings.TrimSpace(Auth)
	if IsUseSessionID {
		// Get session ID from cookie
		ID, err = c.Cookie("_app_cakrasoft_session_id")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
	// Lookup user ID from Redis using session ID
	b, err := cache.DataCache.Get(c, CompanyCode, "USER_INFO_"+ID)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if b != nil {
		err = json.Unmarshal(b, &UserInfo)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
	return UserInfo, nil
}

func InsertUserGroupAccessP(c *gin.Context) {
	ctx, span := global_var.Tracer.Start(c, "InsertUserGroupAccessP")
	defer span.End()

	type HotelStruct struct {
		AccessForm              string  `json:"access_form" binding:"required"`
		AccessSpecial           string  `json:"access_special" binding:"required"`
		AccessKeylock           string  `json:"access_keylock" binding:"required"`
		AccessReservation       string  `json:"access_reservation" binding:"required"`
		AccessDeposit           string  `json:"access_deposit" binding:"required"`
		AccessInHouse           string  `json:"access_in_house" binding:"required"`
		AccessWalkIn            string  `json:"access_walk_in" binding:"required"`
		AccessFolio             string  `json:"access_folio" binding:"required"`
		AccessFolioHistory      string  `json:"access_folio_history" binding:"required"`
		AccessFloorPlan         string  `json:"access_floor_plan" binding:"required"`
		AccessMemberVoucherGift string  `json:"access_member_voucher_gift" binding:"required"`
		SaMaxDiscountPercent    int     `json:"sa_max_discount_percent"`
		SaMaxDiscountAmount     float64 `json:"sa_max_discount_amount" `
	}

	type AssetStruct struct {
		AccessForm              string `json:"access_form" binding:"required"`
		AccessInventoryReceive  string `json:"access_inventory_receive" binding:"required"`
		AccessFixedAssetReceive string `json:"access_fixed_asset_receive" binding:"required"`
		AccessSpecial           string `json:"access_special" binding:"required"`
	}

	type PointOfSalesStruct struct {
		AccessForm                string `json:"access_form" binding:"required"`
		AccessSpecial             string `json:"access_special" binding:"required"`
		AccessTransactionTerminal string `json:"access_transaction_terminal" binding:"required"`
		AccessTableView           string `json:"access_table_view" binding:"required"`
		AccessReservation         string `json:"access_reservation"`
	}

	type BanquetStruct struct {
		AccessForm         string `json:"access_form" binding:"required"`
		AccessSpecial      string `json:"access_special" binding:"required"`
		AccessReservation  string `json:"access_reservation"`
		AccessDeposit      string `json:"access_deposit"`
		AccessInHouse      string `json:"access_in_house" binding:"required"`
		AccessFolio        string `json:"access_folio"`
		AccessFolioHistory string `json:"access_folio_history"`
	}

	type AccountingStruct struct {
		AccessForm        string `json:"access_form" binding:"required"`
		AccessSpecial     string `json:"access_special" binding:"required"`
		AccessInvoice     string `json:"access_invoice" binding:"required"`
		PrintInvoiceCount int    `json:"print_invoice_count"`
	}

	type GeneralStruct struct {
		AccessModule string `json:"access_module" binding:"required"`
		IsActive     uint8  `json:"is_active"`
	}

	type ReportStruct struct {
		AccessForm          string `json:"access_form"`
		AccessFoReport      string `json:"access_fo_report"`
		AccessPosReport     string `json:"access_pos_report"`
		AccessBanReport     string `json:"access_ban_report"`
		AccessAccReport     string `json:"access_acc_report"`
		AccessAstReport     string `json:"access_ast_report"`
		AccessPyrReport     string `json:"access_pyr_report"`
		AccessCorReport     string `json:"access_cor_report"`
		AccessPreviewReport string `json:"access_preview_report"`
	}

	type ToolsStruct struct {
		AccessForm          string `json:"access_form" binding:"required"`
		AccessConfiguration string `json:"access_configuration" binding:"required"`
		AccessCompany       string `json:"access_company" binding:"required"`
	}

	type DataInputStruct struct {
		AccessLevelCode int                `json:"access_level_code"`
		Code            string             `json:"code"`
		General         GeneralStruct      `json:"general"`
		Hotel           HotelStruct        `json:"hotel"`
		PointOfSales    PointOfSalesStruct `json:"point_of_sales"`
		// Banquet         BanquetStruct      `json:"banquet"`
		Accounting AccountingStruct `json:"accounting"`
		Asset      AssetStruct      `json:"asset"`
		Report     ReportStruct     `json:"report"`
		Tools      ToolsStruct      `json:"tools"`
	}

	var DataInput DataInputStruct
	err := c.BindJSON(&DataInput)
	if err != nil {
		fmt.Println(err.Error())
		errMsg := General.GenerateValidateErrorMsg(c, err)
		master_data.SendResponse(global_var.ResponseCode.InvalidDataFormat, errMsg, nil, c)
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
	UserID := c.GetString("ValidUserCode")

	// Get User Info
	UserInfo, err := loadUserInfo(c, UserID, pConfig.CompanyCode)
	if err != nil {
		utils.LogResponseError(c, ctx, global_var.AppLogger, errors.Wrap(err, "UpdateUserGroupAccessP.loadUserInfo"))
		master_data.SendResponse(global_var.ResponseCode.InternalServerError, "Failed load User", nil, c)
		return
	}
	if DataInput.AccessLevelCode <= UserInfo.AccessLevel {
		master_data.SendResponse(global_var.ResponseCode.OtherResult, "Your group access level is same or lower than group level", nil, c)
		return
	}

	err = DB.Transaction(func(tx *gorm.DB) error {
		GeneralAccessID := make(chan uint64)
		HotelAccessID := make(chan uint64)
		POSAccessID := make(chan uint64)
		// BanAccessID := make(chan uint64)
		AccAccessID := make(chan uint64)
		AstAccessID := make(chan uint64)
		ReportAccessID := make(chan uint64)
		ToolsAccessID := make(chan uint64)
		//General
		go InsertGeneralUserGroup(tx, GeneralAccessID, DataInput.General.AccessModule, DataInput.General.IsActive, UserID)

		//hotel
		go InsertUserGroup(tx, HotelAccessID, DataInput.Hotel.AccessForm, DataInput.Hotel.AccessSpecial, DataInput.Hotel.AccessKeylock, DataInput.Hotel.AccessReservation,
			DataInput.Hotel.AccessDeposit, DataInput.Hotel.AccessInHouse, DataInput.Hotel.AccessWalkIn, DataInput.Hotel.AccessFolio, DataInput.Hotel.AccessFolioHistory, DataInput.Hotel.AccessFloorPlan,
			DataInput.Hotel.AccessMemberVoucherGift, DataInput.Hotel.SaMaxDiscountPercent, DataInput.Hotel.SaMaxDiscountAmount, UserID)

		//pointOFSales
		go InsertPosUserGroup(tx, POSAccessID, DataInput.PointOfSales.AccessForm, DataInput.PointOfSales.AccessSpecial, DataInput.PointOfSales.AccessTransactionTerminal, DataInput.PointOfSales.AccessTableView, DataInput.PointOfSales.AccessReservation, UserID)

		//Banquet
		// go InsertBanUserGroup(tx, BanAccessID, DataInput.Banquet.AccessForm, DataInput.Banquet.AccessSpecial, DataInput.Banquet.AccessReservation, DataInput.Banquet.AccessDeposit, DataInput.Banquet.AccessInHouse,
		// 	DataInput.Banquet.AccessFolio, DataInput.Banquet.AccessFolioHistory, UserID)

		//Accounting
		go InsertAccUserGroup(tx, AccAccessID, DataInput.Accounting.AccessForm, DataInput.Accounting.AccessSpecial, DataInput.Accounting.AccessInvoice, DataInput.Accounting.PrintInvoiceCount, UserID)

		//Asset
		go InsertAstUserGroup(tx, AstAccessID, DataInput.Asset.AccessForm, DataInput.Asset.AccessInventoryReceive, DataInput.Asset.AccessFixedAssetReceive, DataInput.Asset.AccessSpecial, UserID)

		//Report
		go InsertReportUserGroup(tx, ReportAccessID, DataInput.Report.AccessForm, DataInput.Report.AccessFoReport, DataInput.Report.AccessPosReport, DataInput.Report.AccessBanReport, DataInput.Report.AccessAccReport, DataInput.Report.AccessAstReport, DataInput.Report.AccessPyrReport, DataInput.Report.AccessCorReport, DataInput.Report.AccessPreviewReport, UserID)

		//Tools
		go InsertToolsUserGroup(tx, ToolsAccessID, DataInput.Tools.AccessForm, DataInput.Tools.AccessConfiguration, DataInput.Tools.AccessCompany, UserID)

		// <-GeneralAccessID
		// <-HotelAccessID
		// <-POSAccessID
		// <-AccAccessID
		// <-AstAccessID
		GeneralAccessIDx, HotelAccessIDx, POSAccessIDx, AccAccessIDx, AstAccessIDx, ReportAccessIDx, ToolsAccessIDx := <-GeneralAccessID, <-HotelAccessID, <-POSAccessID, <-AccAccessID, <-AstAccessID, <-ReportAccessID, <-ToolsAccessID
		if err := InsertUserGroupAccess(tx, DataInput.Code, GeneralAccessIDx, HotelAccessIDx, POSAccessIDx, 0, AccAccessIDx, AstAccessIDx, 0, 0, ReportAccessIDx, ToolsAccessIDx, DataInput.AccessLevelCode, 1, UserID); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		master_data.SendResponse(global_var.ResponseCode.DatabaseError, "", nil, c)
		return
	}
	master_data.SendResponse(global_var.ResponseCode.Successfully, "", nil, c)
}

func GetUserAccessLevelListP(c *gin.Context) {
	ctx, span := global_var.Tracer.Start(c, "GetUserAccessLevelListP")
	defer span.End()

	var DataOutput []map[string]interface{}
	// Get Program Configuration
	val, exist := c.Get("pConfig")
	if !exist {
		master_data.SendResponse(global_var.ResponseCode.DatabaseError, "pConfig", nil, c)
		return
	}
	pConfig := val.(*config.CompanyDataConfiguration)
	DB := pConfig.DB
	UserID := c.GetString("ValidUserCode")
	// Get User Info
	UserInfo, err := loadUserInfo(c, UserID, pConfig.CompanyCode)
	if err != nil {
		utils.LogResponseError(c, ctx, global_var.AppLogger, errors.Wrap(err, "GetUserAccessLevelListP.loadUserInfo"))
		master_data.SendResponse(global_var.ResponseCode.InternalServerError, "Failed load User", nil, c)
		return
	}
	err = DB.Table(db_var.TableName.ConstUserAccessLevel).Where("code>?", UserInfo.AccessLevel).Scan(&DataOutput).Error
	if err != nil {
		master_data.SendResponse(global_var.ResponseCode.DatabaseError, "", nil, c)
		return
	}
	master_data.SendResponse(global_var.ResponseCode.Successfully, "", DataOutput, c)
}

func UpdateUserGroupAccessP(c *gin.Context) {
	ctx, span := global_var.Tracer.Start(c, "UpdateUserGroupAccessP")
	defer span.End()

	type HotelStruct struct {
		AccessForm              string  `json:"access_form" binding:"required"`
		AccessSpecial           string  `json:"access_special" binding:"required"`
		AccessKeylock           string  `json:"access_keylock" binding:"required"`
		AccessReservation       string  `json:"access_reservation" binding:"required"`
		AccessDeposit           string  `json:"access_deposit" binding:"required"`
		AccessInHouse           string  `json:"access_in_house" binding:"required"`
		AccessWalkIn            string  `json:"access_walk_in" binding:"required"`
		AccessFolio             string  `json:"access_folio" binding:"required"`
		AccessFolioHistory      string  `json:"access_folio_history" binding:"required"`
		AccessFloorPlan         string  `json:"access_floor_plan" binding:"required"`
		AccessMemberVoucherGift string  `json:"access_member_voucher_gift" binding:"required"`
		SaMaxDiscountPercent    int     `json:"sa_max_discount_percent"`
		SaMaxDiscountAmount     float64 `json:"sa_max_discount_amount" `
	}

	type AssetStruct struct {
		AccessForm              string `json:"access_form" binding:"required"`
		AccessInventoryReceive  string `json:"access_inventory_receive" binding:"required"`
		AccessFixedAssetReceive string `json:"access_fixed_asset_receive" binding:"required"`
		AccessSpecial           string `json:"access_special" binding:"required"`
	}

	type PointOfSalesStruct struct {
		AccessForm                string `json:"access_form" binding:"required"`
		AccessSpecial             string `json:"access_special" binding:"required"`
		AccessTransactionTerminal string `json:"access_transaction_terminal" binding:"required"`
		AccessTableView           string `json:"access_table_view" binding:"required"`
		AccessReservation         string `json:"access_reservation"`
	}

	type BanquetStruct struct {
		AccessForm         string `json:"access_form" binding:"required"`
		AccessSpecial      string `json:"access_special" binding:"required"`
		AccessReservation  string `json:"access_reservation"`
		AccessDeposit      string `json:"access_deposit"`
		AccessInHouse      string `json:"access_in_house" binding:"required"`
		AccessFolio        string `json:"access_folio"`
		AccessFolioHistory string `json:"access_folio_history"`
	}

	type AccountingStruct struct {
		AccessForm        string `json:"access_form" binding:"required"`
		AccessSpecial     string `json:"access_special" binding:"required"`
		AccessInvoice     string `json:"access_invoice" binding:"required"`
		PrintInvoiceCount int    `json:"print_invoice_count"`
	}

	type GeneralStruct struct {
		AccessModule string `json:"access_module" binding:"required"`
		IsActive     uint8  `json:"is_active"`
	}

	type ReportStruct struct {
		AccessForm          string `json:"access_form"`
		AccessFoReport      string `json:"access_fo_report"`
		AccessPosReport     string `json:"access_pos_report"`
		AccessBanReport     string `json:"access_ban_report"`
		AccessAccReport     string `json:"access_acc_report"`
		AccessAstReport     string `json:"access_ast_report"`
		AccessPyrReport     string `json:"access_pyr_report"`
		AccessCorReport     string `json:"access_cor_report"`
		AccessPreviewReport string `json:"access_preview_report"`
	}

	type ToolsStruct struct {
		AccessForm          string `json:"access_form" binding:"required"`
		AccessConfiguration string `json:"access_configuration" binding:"required"`
		AccessCompany       string `json:"access_company" binding:"required"`
	}

	type DataInputStruct struct {
		AccessLevelCode int                `json:"access_level_code"`
		Id              uint64             `json:"id"`
		General         GeneralStruct      `json:"general"`
		Hotel           HotelStruct        `json:"hotel"`
		PointOfSales    PointOfSalesStruct `json:"point_of_sales"`
		// Banquet         BanquetStruct      `json:"banquet"`
		Accounting AccountingStruct `json:"accounting"`
		Asset      AssetStruct      `json:"asset"`
		Report     ReportStruct     `json:"report"`
		Tools      ToolsStruct      `json:"tools"`
	}

	var DataInput DataInputStruct
	err := c.BindJSON(&DataInput)
	if err != nil {
		fmt.Println(err.Error())
		errMsg := General.GenerateValidateErrorMsg(c, err)
		master_data.SendResponse(global_var.ResponseCode.InvalidDataFormat, errMsg, nil, c)
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
	UserID := c.GetString("ValidUserCode")

	// Get User Info
	UserInfo, err := loadUserInfo(c, UserID, pConfig.CompanyCode)
	if err != nil {
		utils.LogResponseError(c, ctx, global_var.AppLogger, errors.Wrap(err, "UpdateUserGroupAccessP.loadUserInfo"))
		master_data.SendResponse(global_var.ResponseCode.InternalServerError, "Failed load User", nil, c)
		return
	}

	var UserGroupAccess db_var.User_group_access
	if err = DB.Table(db_var.TableName.UserGroupAccess).Where("id=?", DataInput.Id).Take(&UserGroupAccess).Error; err != nil {
		master_data.SendResponse(global_var.ResponseCode.DataNotFound, nil, nil, c)
		return
	}

	if UserGroupAccess.UserAccessLevelCode <= UserInfo.AccessLevel {
		master_data.SendResponse(global_var.ResponseCode.OtherResult, "Your group access level is same or lower than group level", nil, c)
		return
	}

	err = DB.Transaction(func(tx *gorm.DB) error {
		//General
		go UpdateGeneralUserGroup(tx, UserGroupAccess.GeneralUserGroupId, DataInput.General.AccessModule, DataInput.General.IsActive, UserID)

		//hotel
		go UpdateUserGroup(tx, UserGroupAccess.UserGroupId, DataInput.Hotel.AccessForm, DataInput.Hotel.AccessSpecial, DataInput.Hotel.AccessKeylock, DataInput.Hotel.AccessReservation,
			DataInput.Hotel.AccessDeposit, DataInput.Hotel.AccessInHouse, DataInput.Hotel.AccessWalkIn, DataInput.Hotel.AccessFolio, DataInput.Hotel.AccessFolioHistory, DataInput.Hotel.AccessFloorPlan,
			DataInput.Hotel.AccessMemberVoucherGift, DataInput.Hotel.SaMaxDiscountPercent, DataInput.Hotel.SaMaxDiscountAmount, UserID)

		//pointOFSales
		go UpdatePosUserGroup(tx, UserGroupAccess.PosUserGroupId, DataInput.PointOfSales.AccessForm, DataInput.PointOfSales.AccessSpecial, DataInput.PointOfSales.AccessTransactionTerminal, DataInput.PointOfSales.AccessTableView, DataInput.PointOfSales.AccessReservation, UserID)

		//Banquet
		// go InsertBanUserGroup(tx, BanAccessID, DataInput.Banquet.AccessForm, DataInput.Banquet.AccessSpecial, DataInput.Banquet.AccessReservation, DataInput.Banquet.AccessDeposit, DataInput.Banquet.AccessInHouse,
		// 	DataInput.Banquet.AccessFolio, DataInput.Banquet.AccessFolioHistory, UserID)

		//Accounting
		go UpdateAccUserGroup(tx, UserGroupAccess.AccUserGroupId, DataInput.Accounting.AccessForm, DataInput.Accounting.AccessSpecial, DataInput.Accounting.AccessInvoice, DataInput.Accounting.PrintInvoiceCount, UserID)

		//Asset
		go UpdateAstUserGroup(tx, UserGroupAccess.AstUserGroupId, DataInput.Asset.AccessForm, DataInput.Asset.AccessInventoryReceive, DataInput.Asset.AccessFixedAssetReceive, DataInput.Asset.AccessSpecial, UserID)

		//Report
		go UpdateReportUserGroup(tx, UserGroupAccess.ReportUserGroupId, DataInput.Report.AccessForm, DataInput.Report.AccessFoReport, DataInput.Report.AccessPosReport, DataInput.Report.AccessBanReport, DataInput.Report.AccessAccReport, DataInput.Report.AccessAstReport, DataInput.Report.AccessPyrReport, DataInput.Report.AccessCorReport, DataInput.Report.AccessPreviewReport, UserID)

		//Tools
		go UpdateToolsUserGroup(tx, UserGroupAccess.ToolsUserGroupId, DataInput.Tools.AccessForm, DataInput.Tools.AccessConfiguration, DataInput.Tools.AccessCompany, UserID)

		if err := tx.Table(db_var.TableName.UserGroupAccess).Where("id=?", DataInput.Id).Update("updated_by", UserID).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		master_data.SendResponse(global_var.ResponseCode.DatabaseError, "", nil, c)
		return
	}
	master_data.SendResponse(global_var.ResponseCode.Successfully, "", nil, c)
}

func GetAccessReportListP(c *gin.Context) {
	var DataOutput []map[string]interface{}

	// Get Program Configuration
	val, exist := c.Get("pConfig")
	if !exist {
		master_data.SendResponse(global_var.ResponseCode.DatabaseError, "pConfig", nil, c)
		return
	}
	pConfig := val.(*config.CompanyDataConfiguration)
	DB := pConfig.DB

	DB.Table(db_var.TableName.Report).Select("report.system_code, report.name as report, report1.name as parent").
		Joins("LEFT OUTER JOIN report report1 ON (report.parent_id = report1.code)").
		Where("report.code>1000").
		Order("report.id_sort").
		Scan(&DataOutput)

	DataOutput2 := make(map[string][]map[string]interface{}, 0)
	if len(DataOutput) > 0 {
		for _, data := range DataOutput {
			if data["system_code"].(string) == global_var.SystemCode.Hotel {
				DataOutput2["hotel"] = append(DataOutput2["hotel"], data)
				continue
			}
			if data["system_code"].(string) == global_var.SystemCode.Pos {
				DataOutput2["point_of_sales"] = append(DataOutput2["point_of_sales"], data)
				continue
			}
			if data["system_code"].(string) == global_var.SystemCode.Banquet {
				DataOutput2["banquet"] = append(DataOutput2["banquet"], data)
				continue
			}
			if data["system_code"].(string) == global_var.SystemCode.Accounting {
				DataOutput2["accounting"] = append(DataOutput2["accounting"], data)
				continue
			}
			if data["system_code"].(string) == global_var.SystemCode.Asset {
				DataOutput2["asset"] = append(DataOutput2["asset"], data)
				continue
			}
		}
	}

	master_data.SendResponse(global_var.ResponseCode.Successfully, "", DataOutput2, c)
}

func GetUserGroupAccessListP(c *gin.Context) {
	IsActive := c.Query("IsActive")

	// Get Program Configuration
	val, exist := c.Get("pConfig")
	if !exist {
		master_data.SendResponse(global_var.ResponseCode.DatabaseError, "pConfig", nil, c)
		return
	}
	pConfig := val.(*config.CompanyDataConfiguration)
	DB := pConfig.DB

	val2, exist := c.Get("UserInfo")
	if !exist {
		master_data.SendResponse(global_var.ResponseCode.InternalServerError, "Failed get User Info", nil, c)
		return
	}

	UserInfo := val2.(global_var.TUserInfo)
	Query := DB.Table(db_var.TableName.UserGroupAccess).Select("user_group_access.*", "const_user_access_level.name as AccessLevel").
		Joins("LEFT JOIN const_user_access_level ON user_group_access.user_access_level_code=const_user_access_level.code").
		Where("user_access_level_code>?", UserInfo.AccessLevel)
	if IsActive == "1" {
		Query.Where("is_active=1")
	} else if IsActive == "0" {
		Query.Where("is_active=0")
	}
	var DataOutput []map[string]interface{}
	Query.Scan(&DataOutput)

	if Query.Error != nil {
		master_data.SendResponse(global_var.ResponseCode.DatabaseError, "", nil, c)
		return
	}
	master_data.SendResponse(global_var.ResponseCode.Successfully, "", DataOutput, c)

}

func GetUserGroupAccessP(c *gin.Context) {
	Id := c.Param("Id")

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
	// Get Program Configuration
	val, exist := c.Get("pConfig")
	if !exist {
		master_data.SendResponse(global_var.ResponseCode.DatabaseError, "pConfig", nil, c)
		return
	}
	pConfig := val.(*config.CompanyDataConfiguration)
	DB := pConfig.DB

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
		Where("user_group_access.id=?", Id).
		Take(&DataOutput).Error

	if err != nil {
		master_data.SendResponse(global_var.ResponseCode.DataNotFound, "", nil, c)
		return
	}
	master_data.SendResponse(global_var.ResponseCode.Successfully, "", DataOutput, c)

}

func GetUserListP(c *gin.Context) {
	IsActive := c.Query("IsActive")
	// Get Program Configuration
	val, exist := c.Get("pConfig")
	if !exist {
		master_data.SendResponse(global_var.ResponseCode.DatabaseError, "pConfig", nil, c)
		return
	}
	pConfig := val.(*config.CompanyDataConfiguration)
	DB := pConfig.DB

	val2, exist := c.Get("UserInfo")
	if !exist {
		master_data.SendResponse(global_var.ResponseCode.InternalServerError, "Failed get User Info", nil, c)
		return
	}

	UserInfo := val2.(global_var.TUserInfo)
	Query := DB.Table(db_var.TableName.User).Select("user.*").
		Joins("LEFT JOIN user_group_access ON user.user_group_access_code=user_group_access.code").
		Joins("LEFT JOIN const_user_access_level ON user_group_access.user_access_level_code=const_user_access_level.code").
		Where("user_access_level_code>?", UserInfo.AccessLevel)
	if IsActive == "1" {
		Query.Where("user.is_active=1")
	} else if IsActive == "0" {
		Query.Where("user.is_active=0")
	}
	var DataOutput []map[string]interface{}
	Query.Scan(&DataOutput)

	if Query.Error != nil {
		master_data.SendResponse(global_var.ResponseCode.DatabaseError, "", nil, c)
		return
	}
	master_data.SendResponse(global_var.ResponseCode.Successfully, "", DataOutput, c)
}

func GetUserP(c *gin.Context) {
	UserId := c.Param("Id")
	// Get Program Configuration
	val, exist := c.Get("pConfig")
	if !exist {
		master_data.SendResponse(global_var.ResponseCode.DatabaseError, "pConfig", nil, c)
		return
	}
	pConfig := val.(*config.CompanyDataConfiguration)
	DB := pConfig.DB

	var DataOutput map[string]interface{}
	if err := DB.Table(db_var.TableName.User).Where("id=?", UserId).Take(&DataOutput).Error; err != nil {
		master_data.SendResponse(global_var.ResponseCode.DataNotFound, "", nil, c)
		return
	}
	master_data.SendResponse(global_var.ResponseCode.Successfully, "", DataOutput, c)
}

func InsertUserP(c *gin.Context) {
	type DataInputStruct struct {
		UserId    string `json:"user_id" binding:"required"`
		FullName  string `json:"full_name" binding:"required"`
		Password  string `json:"password" binding:"required"`
		GroupCode string `json:"group_code" binding:"required"`
	}
	var DataInput DataInputStruct
	err := c.BindJSON(&DataInput)
	if err != nil {
		fmt.Println(err.Error())
		errMsg := General.GenerateValidateErrorMsg(c, err)
		master_data.SendResponse(global_var.ResponseCode.InvalidDataFormat, errMsg, nil, c)
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
	UserID := c.GetString("ValidUserCode")

	var TotalUser int64
	if err := DB.Table(db_var.TableName.User).Count(&TotalUser).Error; err != nil {
		master_data.SendResponse(global_var.ResponseCode.DatabaseError, "", nil, c)
		return
	}

	if TotalUser >= pConfig.MaxUser {
		master_data.SendResponse(global_var.ResponseCode.OtherResult, "Cannot add more user!", nil, c)
		return
	}

	if err := InsertUser(DB, DataInput.UserId, DataInput.FullName, general.GetMD5Hash(DataInput.Password), DataInput.GroupCode, UserID); err != nil {
		master_data.SendResponse(global_var.ResponseCode.DatabaseError, "", nil, c)
		return
	}

	master_data.SendResponse(global_var.ResponseCode.Successfully, "", nil, c)
}

func UpdateUserP(c *gin.Context) {
	type DataInputStruct struct {
		Id              uint64 `json:"id" binding:"required"`
		PasswordChanged bool   `json:"password_changed"`
		FullName        string `json:"full_name" binding:"required"`
		Password        string `json:"password" binding:"required"`
		GroupCode       string `json:"group_code" binding:"required"`
	}
	var DataInput DataInputStruct
	err := c.BindJSON(&DataInput)
	if err != nil {
		fmt.Println(err.Error())
		errMsg := General.GenerateValidateErrorMsg(c, err)
		master_data.SendResponse(global_var.ResponseCode.InvalidDataFormat, errMsg, nil, c)
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

	UserID := c.GetString("ValidUserCode")

	if err := UpdateUser(DB, DataInput.Id, DataInput.FullName, general.GetMD5Hash(DataInput.Password), DataInput.GroupCode, DataInput.PasswordChanged, UserID); err != nil {
		master_data.SendResponse(global_var.ResponseCode.DatabaseError, "", nil, c)
		return
	}

	master_data.SendResponse(global_var.ResponseCode.Successfully, "", nil, c)
}

func DeleteUserP(c *gin.Context) {
	UserID := c.Param("Code")
	// Get Program Configuration
	val, exist := c.Get("pConfig")
	if !exist {
		master_data.SendResponse(global_var.ResponseCode.DatabaseError, "pConfig", nil, c)
		return
	}
	pConfig := val.(*config.CompanyDataConfiguration)
	DB := pConfig.DB

	err := DB.Transaction(func(tx *gorm.DB) error {
		if err := DB.Table(db_var.TableName.User).Where("code=?", UserID).Update("updated_by", UserID).Error; err != nil {
			return err
		}
		if err := DB.Table(db_var.TableName.User).Where("code=?", UserID).Delete(&UserID).Error; err != nil {
			return err
		}

		if err := cache.DataCache.Del(c, pConfig.CompanyCode, "USER_INFO_"+UserID); err != nil {
			return err

		}
		return nil
	})

	if err != nil {
		master_data.SendResponse(global_var.ResponseCode.DatabaseError, "", nil, c)
		return
	}
	master_data.SendResponse(global_var.ResponseCode.Successfully, "", nil, c)
}

func ActivateDeactivateUserP(c *gin.Context) {
	type DataInputStruct struct {
		UserId   uint64 `json:"user_id" binding:"required"`
		IsActive string `json:"is_active" binding:"enum=1_0"`
	}

	var DataInput DataInputStruct
	err := c.BindJSON(&DataInput)
	if err != nil {
		fmt.Println(err.Error())
		errMsg := General.GenerateValidateErrorMsg(c, err)
		master_data.SendResponse(global_var.ResponseCode.InvalidDataFormat, errMsg, nil, c)
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
	UserID := c.GetString("ValidUserCode")

	fmt.Println("aa", UserID)
	err = DB.Table(db_var.TableName.User).Where("id=?", DataInput.UserId).Updates(map[string]interface{}{
		"is_active":  DataInput.IsActive,
		"updated_by": UserID}).Error

	if err != nil {
		master_data.SendResponse(global_var.ResponseCode.DatabaseError, "", nil, c)
		return
	}
	master_data.SendResponse(global_var.ResponseCode.Successfully, "", nil, c)
}

func DeleteUserGroupAccessP(c *gin.Context) {
	ctx, span := global_var.Tracer.Start(c, "DeleteUserGroupAccessP")
	defer span.End()

	GroupId := c.Param("Id")
	// Get Program Configuration
	val, exist := c.Get("pConfig")
	if !exist {
		master_data.SendResponse(global_var.ResponseCode.DatabaseError, "pConfig", nil, c)
		return
	}
	pConfig := val.(*config.CompanyDataConfiguration)
	DB := pConfig.DB
	UserID := c.GetString("ValidUserCode")
	// Get User Info
	UserInfo, err := loadUserInfo(c, UserID, pConfig.CompanyCode)
	if err != nil {
		utils.LogResponseError(c, ctx, global_var.AppLogger, errors.Wrap(err, "UpdateUserGroupAccessP.loadUserInfo"))
		master_data.SendResponse(global_var.ResponseCode.InternalServerError, "Failed load User", nil, c)
		return
	}
	var UserGroupAccess db_var.User_group_access
	if err := DB.Table(db_var.TableName.UserGroupAccess).Where("id=?", GroupId).Take(&UserGroupAccess).Error; err != nil {
		master_data.SendResponse(global_var.ResponseCode.DataNotFound, nil, nil, c)
		return
	}

	if UserGroupAccess.UserAccessLevelCode <= UserInfo.AccessLevel {
		master_data.SendResponse(global_var.ResponseCode.OtherResult, "Your group access level is same or lower than group level", nil, c)
		return
	}
	err = DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(db_var.TableName.UserGroupAccess).Where("id=?", GroupId).Update("updated_by", UserID).Error; err != nil {
			return err
		}
		if err := tx.Table(db_var.TableName.UserGroup).Where("id=?", UserGroupAccess.UserGroupId).Delete(&UserGroupAccess.UserGroupId).Error; err != nil {
			return err
		}
		if err := tx.Table(db_var.TableName.PosUserGroup).Where("id=?", UserGroupAccess.PosUserGroupId).Delete(&UserGroupAccess.PosUserGroupId).Error; err != nil {
			return err
		}
		if err := tx.Table(db_var.TableName.AstUserGroup).Where("id=?", UserGroupAccess.AstUserGroupId).Delete(&UserGroupAccess.AstUserGroupId).Error; err != nil {
			return err
		}
		if err := tx.Table(db_var.TableName.AccUserGroup).Where("id=?", UserGroupAccess.AccUserGroupId).Delete(&UserGroupAccess.AccUserGroupId).Error; err != nil {
			return err
		}
		if err := tx.Table(db_var.TableName.BanUserGroup).Where("id=?", UserGroupAccess.BanUserGroupId).Delete(&UserGroupAccess.BanUserGroupId).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		master_data.SendResponse(global_var.ResponseCode.DatabaseError, "", nil, c)
		return
	}
	master_data.SendResponse(global_var.ResponseCode.Successfully, "", nil, c)
}

func GetUserInventorySubDepartmentAccessListP(c *gin.Context) {
	ctx, span := global_var.Tracer.Start(c, "GetUserInventorySubDepartmentAccessListP")
	defer span.End()

	UserCode := c.Param("UserCode")

	// Get Program Configuration
	val, exist := c.Get("pConfig")
	if !exist {
		master_data.SendResponse(global_var.ResponseCode.DatabaseError, "pConfig", nil, c)
		return
	}
	pConfig := val.(*config.CompanyDataConfiguration)
	DB := pConfig.DB
	var DataOutput []map[string]interface{}
	if err := DB.Table(db_var.TableName.AstUserSubDepartment).
		Select("ast_user_sub_department.*", "cfg_init_sub_department.name AS sub_department").
		Joins("LEFT JOIN cfg_init_sub_department ON ast_user_sub_department.sub_department_code=cfg_init_sub_department.code").
		Where("user_code", UserCode).Scan(&DataOutput).Error; err != nil {
		utils.LogResponseError(c, ctx, global_var.AppLogger, errors.Wrap(err, "GetUserInventoryAccessListP.Query"))
		master_data.SendResponse(global_var.ResponseCode.DatabaseError, "", nil, c)
		return
	}

	if DataOutput == nil {
		DataOutput = make([]map[string]interface{}, 0)
	}

	master_data.SendResponse(global_var.ResponseCode.Successfully, "", DataOutput, c)
}

func GetUserInventorySubDepartmentAccessP(c *gin.Context) {
	ctx, span := global_var.Tracer.Start(c, "GetUserInventorySubDepartmentAccessP")
	defer span.End()

	ID := c.Param("id")

	// Get Program Configuration
	val, exist := c.Get("pConfig")
	if !exist {
		master_data.SendResponse(global_var.ResponseCode.DatabaseError, "pConfig", nil, c)
		return
	}
	pConfig := val.(*config.CompanyDataConfiguration)
	DB := pConfig.DB
	var DataOutput map[string]interface{}
	if err := DB.Table(db_var.TableName.AstUserSubDepartment).
		Where("id", ID).Limit(1).Scan(&DataOutput).Error; err != nil {
		utils.LogResponseError(c, ctx, global_var.AppLogger, errors.Wrap(err, "GetUserInventoryAccessP.Query"))
		master_data.SendResponse(global_var.ResponseCode.DatabaseError, "", nil, c)
		return
	}

	if DataOutput == nil {
		master_data.SendResponse(global_var.ResponseCode.DataNotFound, "", nil, c)
		return
	}

	master_data.SendResponse(global_var.ResponseCode.Successfully, "Success", DataOutput, c)
}

func InsertUserInventorySubDepartmentP(c *gin.Context) {
	ctx, span := global_var.Tracer.Start(c, "InsertUserInventorySubDepartmentP")
	defer span.End()

	type DataInputStruct struct {
		UserCode              string `json:"user_code" binding:"required"`
		SubDepartmentCode     string `json:"sub_department_code" binding:"required"`
		IsCanInvPrApprove1    uint8  `json:"is_can_inv_pr_approve1"`
		IsCanInvPrApprove12   uint8  `json:"is_can_inv_pr_approve12"`
		IsCanInvPrApprove2    uint8  `json:"is_can_inv_pr_approve2"`
		IsCanInvPrApprove3    uint8  `json:"is_can_inv_pr_approve3"`
		IsCanInvPrAssignPrice uint8  `json:"is_can_inv_pr_assign_price"`
		IsCanInvSrApprove1    uint8  `json:"is_can_inv_sr_approve1"`
		IsCanInvSrApprove2    uint8  `json:"is_can_inv_sr_approve2"`
		CreatedBy             string
		UpdatedBy             string
	}
	var DataInput DataInputStruct
	err := c.BindJSON(&DataInput)
	if err != nil {
		//fmt.Println(err.Error())
		errMsg := General.GenerateValidateErrorMsg(c, err)
		master_data.SendResponse(global_var.ResponseCode.InvalidDataFormat, errMsg, nil, c)
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

	var UserCode string
	if err := DB.Table(db_var.TableName.AstUserSubDepartment).Select("user_code").
		Where("user_code", DataInput.UserCode).
		Where("sub_department_code", DataInput.SubDepartmentCode).
		Limit(1).Scan(&UserCode).Error; err != nil {
		utils.LogResponseError(c, ctx, global_var.AppLogger, errors.Wrap(err, "InsertUserInventorySubDepartmentP.CheckUserCode"))
		master_data.SendResponse(global_var.ResponseCode.DatabaseError, "", nil, c)
		return
	}

	if UserCode != "" {
		master_data.SendResponse(global_var.ResponseCode.DuplicateEntry, "Sub department already exist!", nil, c)
		return
	}
	DataInput.CreatedBy = c.GetString("ValidUserCode")
	DataInput.UpdatedBy = DataInput.CreatedBy
	if err := DB.Table(db_var.TableName.AstUserSubDepartment).Create(&DataInput).Error; err != nil {
		utils.LogResponseError(c, ctx, global_var.AppLogger, errors.Wrap(err, "InsertUserInventorySubDepartmentP.Create"))
		master_data.SendResponse(global_var.ResponseCode.DatabaseError, "", nil, c)
		return
	}
	master_data.SendResponse(global_var.ResponseCode.Successfully, "Success", nil, c)
}

func UpdateUserInventorySubDepartmentP(c *gin.Context) {
	ctx, span := global_var.Tracer.Start(c, "UpdateUserInventorySubDepartmentP")
	defer span.End()

	type DataInputStruct struct {
		UserCode              string `json:"user_code" binding:"required"`
		SubDepartmentCode     string `json:"sub_department_code" binding:"required"`
		IsCanInvPrApprove1    *uint8 `json:"is_can_inv_pr_approve1"`
		IsCanInvPrApprove12   *uint8 `json:"is_can_inv_pr_approve12"`
		IsCanInvPrApprove2    *uint8 `json:"is_can_inv_pr_approve2"`
		IsCanInvPrApprove3    *uint8 `json:"is_can_inv_pr_approve3"`
		IsCanInvPrAssignPrice *uint8 `json:"is_can_inv_pr_assign_price"`
		IsCanInvSrApprove1    *uint8 `json:"is_can_inv_sr_approve1"`
		IsCanInvSrApprove2    *uint8 `json:"is_can_inv_sr_approve2"`
		Id                    uint64 `json:"id"  binding:"required"`
		UpdatedBy             string
	}
	var DataInput DataInputStruct
	err := c.BindJSON(&DataInput)
	if err != nil {
		//fmt.Println(err.Error())
		errMsg := General.GenerateValidateErrorMsg(c, err)
		master_data.SendResponse(global_var.ResponseCode.InvalidDataFormat, errMsg, nil, c)
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

	var UserCode string
	if err := DB.Table(db_var.TableName.AstUserSubDepartment).Select("user_code").
		Where("user_code", DataInput.UserCode).
		Where("sub_department_code", DataInput.SubDepartmentCode).
		Where("id <> ?", DataInput.Id).
		Limit(1).Scan(&UserCode).Error; err != nil {
		utils.LogResponseError(c, ctx, global_var.AppLogger, errors.Wrap(err, "UpdateUserInventorySubDepartmentP.CheckUserCode"))
		master_data.SendResponse(global_var.ResponseCode.DatabaseError, "", nil, c)
		return
	}

	if UserCode != "" {
		master_data.SendResponse(global_var.ResponseCode.DuplicateEntry, "Sub department is duplicated!", nil, c)
		return
	}
	DataInput.UpdatedBy = c.GetString("ValidUserCode")
	if err := DB.Table(db_var.TableName.AstUserSubDepartment).Where("id", DataInput.Id).Updates(&DataInput).Error; err != nil {
		utils.LogResponseError(c, ctx, global_var.AppLogger, errors.Wrap(err, "UpdateUserInventorySubDepartmentP.Updates"))
		master_data.SendResponse(global_var.ResponseCode.DatabaseError, "", nil, c)
		return
	}
	master_data.SendResponse(global_var.ResponseCode.Successfully, "Success", nil, c)
}

func DeleteUserInventorySubDepartmentP(c *gin.Context) {
	ctx, span := global_var.Tracer.Start(c, "DeleteUserInventorySubDepartmentP")
	defer span.End()

	ID := c.Param("id")
	// Get Program Configuration
	val, exist := c.Get("pConfig")
	if !exist {
		master_data.SendResponse(global_var.ResponseCode.DatabaseError, "pConfig", nil, c)
		return
	}
	pConfig := val.(*config.CompanyDataConfiguration)
	DB := pConfig.DB

	if err := DB.Table(db_var.TableName.AstUserSubDepartment).Where("id", ID).Delete(&ID).Error; err != nil {
		utils.LogResponseError(c, ctx, global_var.AppLogger, errors.Wrap(err, "DeleteUserInventorySubDepartmentP.Delete"))
		master_data.SendResponse(global_var.ResponseCode.DatabaseError, "", nil, c)
		return
	}
	master_data.SendResponse(global_var.ResponseCode.Successfully, "Success", nil, c)
}
