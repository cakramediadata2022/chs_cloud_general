package global_query

import (
	"chs_cloud_general/internal/config"
	"chs_cloud_general/internal/db_var"
	DBVar "chs_cloud_general/internal/db_var"
	General "chs_cloud_general/internal/general"
	"chs_cloud_general/internal/global_var"
	GlobalVar "chs_cloud_general/internal/global_var"
	MasterData "chs_cloud_general/internal/master_data"
	"chs_cloud_general/internal/utils/cache"
	"chs_cloud_general/internal/utils/websocket"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RoomNumberListStruct struct {
	RoomNumber string
}

type RoomRateListStruct struct {
	Code, Name, WeekdayRate1, WeekendRate1 string
}

func SaveTextToFile(S, DestinationFolder string) {
	f, err := os.Create(DestinationFolder)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err2 := f.WriteString(S)

	if err2 != nil {
		log.Fatal(err2)
	}
}

func IsAuthorized(Token string) string {
	if Token != "" {
		TokenCheck, Err := jwt.Parse(Token, func(TokenCheck *jwt.Token) (interface{}, error) {
			if _, ok := TokenCheck.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Check Token Error")
			}
			return GlobalVar.SigningKey, nil
		})

		if Err == nil {
			if claims, ok := TokenCheck.Claims.(jwt.MapClaims); ok && TokenCheck.Valid {
				if !claims["refresh"].(bool) {
					return strings.ToUpper(claims["user"].(string))
				}
			}
		}
	}
	return ""
}

// validate

func GetTimeLocation() *time.Location {
	localTime, _ := time.LoadLocation("Asia/Makassar")
	return localTime
}

func GetAuditDate(c *gin.Context, DB *gorm.DB, Reload bool) time.Time {
	val, exist := c.Get("pConfig")
	if !exist {
	}
	pConfig := val.(*config.CompanyDataConfiguration)
	CompanyCode := pConfig.CompanyCode
	DateX, err := cache.DataCache.GetString(c, CompanyCode, "AUDIT_DATE")
	if err != nil || Reload {
		Date := MasterData.GetFieldTimeQuery(DB,
			"SELECT audit_date FROM audit_log "+
				"ORDER BY id DESC "+
				"LIMIT 1")

		cache.DataCache.Set(c, CompanyCode, "AUDIT_DATE", Date, 6*time.Hour)
		return General.DateOf(Date)
	}

	return General.StrZToDate(DateX)

}

func GetAuditDateP(c *gin.Context) {
	// Get Program Configuration
	val, exist := c.Get("pConfig")
	if !exist {
		MasterData.SendResponse(GlobalVar.ResponseCode.DatabaseError, "pConfig", nil, c)
		return
	}
	pConfig := val.(*config.CompanyDataConfiguration)
	DB := pConfig.DB
	AuditDate := GetAuditDate(c, DB, false)
	MasterData.SendResponse(GlobalVar.ResponseCode.Successfully, "", AuditDate, c)
}

func GetServerDateTime(c *gin.Context, DB *gorm.DB) time.Time {
	return MasterData.GetFieldTimeQuery(DB, "SELECT NOW() AS DateServer;")
}

func GetServerDateTimeP(c *gin.Context) {
	// Get Program Configuration
	val, exist := c.Get("pConfig")
	if !exist {
		MasterData.SendResponse(GlobalVar.ResponseCode.DatabaseError, "pConfig", nil, c)
		return
	}
	pConfig := val.(*config.CompanyDataConfiguration)
	DB := pConfig.DB
	MasterData.SendResponse(GlobalVar.ResponseCode.Successfully, "", GetServerDateTime(c, DB), c)
}

func GetServerDate(c *gin.Context, DB *gorm.DB) time.Time {
	return MasterData.GetFieldTimeQuery(DB, "SELECT DATE(NOW()) AS DateServer;")
}

func GetServerID(DB *gorm.DB) int {
	var ID int
	DB.Raw("SELECT @@server_id").Scan(&ID)
	return ID
}

func GetServerDateP(c *gin.Context) {
	// Get Program Configuration
	val, exist := c.Get("pConfig")
	if !exist {
		MasterData.SendResponse(global_var.ResponseCode.DatabaseError, "pConfig", nil, c)
		return
	}
	pConfig := val.(*config.CompanyDataConfiguration)
	DB := pConfig.DB
	MasterData.SendResponse(GlobalVar.ResponseCode.Successfully, "", GetServerDate(c, DB), c)
}

func IsDiscountLimit(c *gin.Context, DB *gorm.DB, UserID string, RateOriginal, RateOverride float64) bool {
	var DiscountPercent float64
	DiscountAmount := RateOriginal - RateOverride
	if RateOriginal > 0 {
		DiscountPercent = (DiscountAmount) / RateOriginal * 100
	}
	var Code string
	if err := DB.Table(DBVar.TableName.User).Select("IFNULL(user_group_access.code,'') AS Code").
		Joins("INNER JOIN user_group_access ON (user.user_group_access_code = user_group_access.code)").
		Joins("INNER JOIN user_group ON (user_group_access.user_group_id = user_group.id)").
		Where("user.code=?", UserID).
		Where("user_group.sa_max_discount_percent>=?", DiscountPercent).
		Where("user_group.sa_max_discount_amount>=?", DiscountAmount).
		Scan(&Code).Error; err != nil {
		return true
	}

	return Code == ""

}

func GetReservationStatusCode(DB *gorm.DB, ReservationNumber uint64) string {
	Status := ""
	DB.Table(DBVar.TableName.Reservation).Select("status_code").Where("number = ?", ReservationNumber).Take(&Status)

	return Status
}

func GetAvailableRoomCountByType(DB *gorm.DB, ArrivalDate, DepartureDate time.Time, RoomTypeCode, BedTypeCode string, ReservationNumber, FolioNumber, RoomUnavailableID, RoomAllotmentID uint64, ReadyOnly, AllotmentOnly bool) (int64, error) {
	var Result int64 = 0
	ArrivalDateStr := ArrivalDate.Format("2006-01-02")
	DepartureDateStr := DepartureDate.Format("2006-01-02")

	FieldCountRoomTotal := ""
	FieldCountRoom := ""
	FieldCountReservation := ""
	FieldCountFolio := ""
	FieldCountUnavailable := ""
	FieldCountAllotment := ""
	CountDay := General.DaysBetween(ArrivalDate, DepartureDate)
	if CountDay <= 0 {
		FieldCountRoomTotal = " SUM(IFNULL(A.RoomCount, 0)) AS AvailableRoomCount1 "
		FieldCountRoom = " COUNT(number) AS RoomCount "
		FieldCountReservation = " -COUNT(guest_detail.room_number) AS RoomCount "
		FieldCountFolio = " -COUNT(DISTINCT guest_detail.room_number, NULL) AS RoomCount "
		FieldCountUnavailable = " -COUNT(DISTINCT cfg_init_room.number) AS RoomCount "
		if AllotmentOnly {
			FieldCountAllotment = " COUNT(DISTINCT cfg_init_room.number) AS RoomCount "
		} else {
			FieldCountAllotment = " -COUNT(DISTINCT cfg_init_room.number) AS RoomCount "
		}
	} else {
		for Count := 1; Count <= CountDay; Count++ {
			CurrentDateStr := General.FormatDate1(General.IncDay(ArrivalDate, Count-1))
			CountString := strconv.Itoa(Count)
			FieldCountRoomTotal = FieldCountRoomTotal + " SUM(IFNULL(A.RoomCount" + CountString + ", 0)) AS AvailableRoomCount" + CountString
			FieldCountRoom = FieldCountRoom + " COUNT(number) AS RoomCount" + CountString
			FieldCountReservation = FieldCountReservation + " -COUNT(IF(DATE(guest_detail.arrival)<='" + CurrentDateStr + "' AND DATE(guest_detail.departure)>'" + CurrentDateStr + "', guest_detail.room_number, NULL)) AS RoomCount" + CountString
			FieldCountFolio = FieldCountFolio + " -COUNT(DISTINCT IF(DATE(guest_detail.arrival)<='" + CurrentDateStr + "' AND DATE(guest_detail.departure)>'" + CurrentDateStr + "', guest_detail.room_number, NULL)) AS RoomCount" + CountString
			FieldCountUnavailable = FieldCountUnavailable + " -COUNT(DISTINCT IF(DATE(room_unavailable.start_date)<='" + CurrentDateStr + "' AND DATE(room_unavailable.end_date)>='" + CurrentDateStr + "', cfg_init_room.number, NULL)) AS RoomCount" + CountString
			if AllotmentOnly {
				FieldCountAllotment = FieldCountAllotment + " COUNT(DISTINCT IF(DATE(room_allotment.from_date)<='" + CurrentDateStr + "' AND DATE(room_allotment.to_date)>='" + CurrentDateStr + "', cfg_init_room.number, NULL)) AS RoomCount" + CountString
			} else {
				FieldCountAllotment = FieldCountAllotment + " -COUNT(DISTINCT IF(DATE(room_allotment.from_date)<='" + CurrentDateStr + "' AND DATE(room_allotment.to_date)>='" + CurrentDateStr + "', cfg_init_room.number, NULL)) AS RoomCount" + CountString
			}

			if Count == CountDay {
				FieldCountRoomTotal = FieldCountRoomTotal + " "
				FieldCountRoom = FieldCountRoom + " "
				FieldCountReservation = FieldCountReservation + " "
				FieldCountFolio = FieldCountFolio + " "
				FieldCountUnavailable = FieldCountUnavailable + " "
				FieldCountAllotment = FieldCountAllotment + " "
			} else {
				FieldCountRoomTotal = FieldCountRoomTotal + ","
				FieldCountRoom = FieldCountRoom + ","
				FieldCountReservation = FieldCountReservation + ","
				FieldCountFolio = FieldCountFolio + ","
				FieldCountUnavailable = FieldCountUnavailable + ","
				FieldCountAllotment = FieldCountAllotment + ","
			}
		}
	}

	QueryCondition1 := ""
	QueryCondition2A := ""
	QueryCondition2B := ""
	if ReservationNumber != 0 {
		QueryCondition1 = " AND reservation.number<>'" + strconv.FormatUint(ReservationNumber, 10) + "'"
		QueryCondition2B = " AND folio.reservation_number<>'" + strconv.FormatUint(ReservationNumber, 10) + "'"
	}

	if FolioNumber != 0 {
		QueryCondition2A = " AND folio.number<>'" + strconv.FormatUint(FolioNumber, 10) + "'"
	}

	QueryCondition3 := ""
	if RoomUnavailableID != 0 {
		QueryCondition3 = " AND room_unavailable.id<>'" + strconv.FormatUint(RoomUnavailableID, 10) + "'"
	}

	QueryCondition4 := ""
	QueryCondition4x := ""
	if BedTypeCode != "" {
		QueryCondition4 = " AND cfg_init_room.bed_type_code='" + BedTypeCode + "'"
		QueryCondition4x = " AND guest_detail.bed_type_code='" + BedTypeCode + "'"
	}

	//Allotment Only
	QueryCondition1 = QueryCondition1 + " AND reservation.is_from_allotment='" + General.BoolToUint8String(AllotmentOnly) + "'"
	QueryNot := ""
	QueryRoom := ""
	QueryRoomUnavailable := ""
	QueryRoomAllotment := ""
	//TODO optimize query, not using prepare statement
	if AllotmentOnly {
		QueryCondition2A = QueryCondition2A + " AND folio.is_from_allotment='" + General.BoolToUint8String(AllotmentOnly) + "'"
		QueryRoomAllotment =
			" WHERE DATE(room_allotment.from_date)<='" + ArrivalDateStr + "'" +
				" AND ADDDATE(DATE(room_allotment.to_date), INTERVAL 1 DAY)>='" + DepartureDateStr + "'"
	} else {
		QueryNot = "NOT "
		QueryRoom =
			"SELECT" +
				" 'A' AS Code," +
				FieldCountRoom +
				"FROM" +
				" cfg_init_room" +
				" WHERE room_type_code='" + RoomTypeCode + "'" +
				QueryCondition4 +
				")UNION("
		QueryRoomUnavailable =
			"SELECT" +
				" 'D' AS Code," +
				FieldCountUnavailable +
				"FROM" +
				" room_unavailable" +
				" LEFT OUTER JOIN cfg_init_room ON (room_unavailable.room_number = cfg_init_room.number)" +
				" WHERE cfg_init_room.room_type_code='" + RoomTypeCode + "'" +
				QueryCondition4 +
				" AND DATE(room_unavailable.start_date)<'" + DepartureDateStr + "'" +
				" AND DATE(room_unavailable.end_date)>='" + ArrivalDateStr + "'" +
				QueryCondition3 +
				")UNION("
		QueryRoomAllotment =
			" WHERE DATE(room_allotment.from_date)<'" + DepartureDateStr + "'" +
				" AND DATE(room_allotment.to_date)>='" + ArrivalDateStr + "'"
	}

	QueryCondition5 := ""
	if RoomAllotmentID != 0 {
		QueryCondition5 = " AND room_allotment.id<>'" + strconv.FormatUint(RoomAllotmentID, 10) + "'"
	}

	var AvailableRoomCountArray []map[string]interface{}
	if err := DB.Raw(
		"SELECT" + FieldCountRoomTotal + "FROM (" +
			"(" +
			QueryRoom +
			"SELECT" +
			" 'B' AS Code," +
			FieldCountReservation +
			"FROM" +
			" reservation" +
			" LEFT OUTER JOIN guest_detail ON (reservation.guest_detail_id = guest_detail.id)" +
			" LEFT OUTER JOIN cfg_init_room ON (guest_detail.room_number = cfg_init_room.number)" +
			" WHERE reservation.status_code='" + GlobalVar.ReservationStatus.New + "'" +
			" AND guest_detail.room_type_code='" + RoomTypeCode + "'" +
			QueryCondition4x +
			" AND DATE(guest_detail.arrival)<'" + DepartureDateStr + "'" +
			" AND DATE(guest_detail.departure)>'" + ArrivalDateStr + "'" +
			QueryCondition1 +
			")UNION(" +
			"SELECT" +
			" 'C' AS Code," +
			FieldCountFolio +
			"FROM" +
			" folio" +
			" LEFT OUTER JOIN guest_detail ON (folio.guest_detail_id = guest_detail.id)" +
			" LEFT OUTER JOIN cfg_init_room ON (guest_detail.room_number = cfg_init_room.number)" +
			" WHERE folio.status_code='" + GlobalVar.FolioStatus.Open + "'" +
			" AND folio.type_code='" + GlobalVar.FolioType.GuestFolio + "'" +
			" AND guest_detail.room_type_code='" + RoomTypeCode + "'" +
			QueryCondition4 +
			" AND DATE(guest_detail.arrival)<'" + DepartureDateStr + "'" +
			" AND DATE(guest_detail.departure)>'" + ArrivalDateStr + "'" +
			QueryCondition2A + QueryCondition2B +
			")UNION(" +
			QueryRoomUnavailable +
			"SELECT" +
			" 'E' AS Code," +
			FieldCountAllotment +
			"FROM" +
			" room_allotment" +
			" LEFT OUTER JOIN cfg_init_room ON (room_allotment.room_number = cfg_init_room.number)" +
			" WHERE cfg_init_room.room_type_code='" + RoomTypeCode + "'" +
			QueryCondition4 +
			" AND DATE(room_allotment.from_date)<'" + DepartureDateStr + "'" +
			" AND DATE(room_allotment.to_date)>='" + ArrivalDateStr + "'" +
			QueryCondition5 + ")) AS A").Find(&AvailableRoomCountArray).Error; err != nil {
		return 0, err
	}

	if len(AvailableRoomCountArray) <= 0 {
		return 0, nil
	}
	var AvailableRoomCount int64
	AvailableRoomCount = General.StrToInt64(AvailableRoomCountArray[0]["AvailableRoomCount1"].(string))
	if CountDay > 1 {
		for _, roomCount := range AvailableRoomCountArray {
			for count := range roomCount {
				var countX = General.StrToInt64(roomCount[count].(string))
				if AvailableRoomCount > countX {
					AvailableRoomCount = countX
				}
			}
		}
	}

	Result = AvailableRoomCount
	var AvailableRoomCountReady int64
	if ReadyOnly {
		if err := DB.Raw(
			"SELECT COUNT(number) AS AvailableRoomCountReady FROM cfg_init_room" +
				" WHERE room_type_code='" + RoomTypeCode + "'" +
				" AND status_code='" + GlobalVar.RoomStatus.Ready + "'" +
				QueryCondition4 + " " +
				"AND number NOT IN(SELECT" +
				" guest_detail.room_number " +
				"FROM" +
				" reservation" +
				" LEFT OUTER JOIN guest_detail ON (reservation.guest_detail_id = guest_detail.id)" +
				" WHERE reservation.status_code='" + GlobalVar.ReservationStatus.New + "'" +
				" AND guest_detail.room_type_code='" + RoomTypeCode + "'" +
				" AND guest_detail.room_number<>''" +
				" AND DATE(guest_detail.arrival)<'" + DepartureDateStr + "'" +
				" AND DATE(guest_detail.departure)>'" + ArrivalDateStr + "'" +
				QueryCondition1 + ") " +
				"AND number NOT IN(SELECT" +
				" guest_detail.room_number " +
				"FROM" +
				" folio" +
				" LEFT OUTER JOIN guest_detail ON (folio.guest_detail_id = guest_detail.id)" +
				" WHERE folio.status_code='" + GlobalVar.FolioStatus.Open + "'" +
				" AND folio.type_code='" + GlobalVar.FolioType.GuestFolio + "'" +
				" AND guest_detail.room_type_code='" + RoomTypeCode + "'" +
				" AND DATE(guest_detail.arrival)<'" + DepartureDateStr + "'" +
				" AND DATE(guest_detail.departure)>'" + ArrivalDateStr + "'" +
				QueryCondition2A + QueryCondition2B + ") " +
				"AND number NOT IN(SELECT" +
				" room_unavailable.room_number " +
				"FROM" +
				" room_unavailable" +
				" LEFT OUTER JOIN cfg_init_room ON (room_unavailable.room_number = cfg_init_room.number)" +
				" WHERE DATE(room_unavailable.start_date)<'" + DepartureDateStr + "'" +
				" AND DATE(room_unavailable.end_date)>='" + ArrivalDateStr + "'" +
				" AND cfg_init_room.room_type_code='" + RoomTypeCode + "'" +
				QueryCondition3 + ") " +
				"AND number " + QueryNot + "IN(SELECT" +
				" room_allotment.room_number " +
				"FROM" +
				" room_allotment" +
				" LEFT OUTER JOIN cfg_init_room ON (room_allotment.room_number = cfg_init_room.number)" +
				QueryRoomAllotment +
				" AND cfg_init_room.room_type_code='" + RoomTypeCode + "'" +
				QueryCondition5 + ")").Limit(1).Find(&AvailableRoomCountReady).Error; err != nil {
			return 0, err
		}

		if ReadyOnly {
			if AvailableRoomCount <= AvailableRoomCountReady {
				return AvailableRoomCount, nil
			}
			return AvailableRoomCountReady, nil
		}
	}

	if AllotmentOnly && (Result < 0) {
		Result = 0
	}

	return Result, nil
}

func GetAvailableRoomByType(DB *gorm.DB, Dataset *GlobalVar.TDataset, ArrivalDate, DepartureDate time.Time, RoomTypeCode, BedTypeCode string, ReservationNumber, FolioNumber, RoomUnavailableID, RoomAllotmentID uint64, ReadyOnly, AllotmentOnly bool) []RoomNumberListStruct {
	ArrivalDateStr := ArrivalDate.Format("2006-01-02")
	DepartureDateStr := DepartureDate.Format("2006-01-02")

	QueryConditionRoomType1 := ""
	QueryConditionRoomType2 := ""
	QueryConditionRoomType3 := ""
	if RoomTypeCode != "" {
		QueryConditionRoomType1 = "room_type_code='" + RoomTypeCode + "' AND "
		QueryConditionRoomType2 = " AND guest_detail.room_type_code='" + RoomTypeCode + "'"
		QueryConditionRoomType3 = " AND cfg_init_room.room_type_code='" + RoomTypeCode + "'"
	}

	QueryCondition1 := ""
	QueryCondition2A := ""
	QueryCondition2B := ""
	if ReservationNumber != 0 {
		QueryCondition1 = " AND reservation.number<>'" + strconv.FormatUint(ReservationNumber, 10) + "'"
		QueryCondition2B = " AND folio.reservation_number<>'" + strconv.FormatUint(ReservationNumber, 10) + "'"
	}

	if FolioNumber != 0 {
		QueryCondition2A = " AND folio.number<>'" + strconv.FormatUint(FolioNumber, 10) + "'"
	}

	QueryCondition3 := ""
	if RoomUnavailableID != 0 {
		QueryCondition3 = " AND room_unavailable.id<>'" + strconv.FormatUint(RoomUnavailableID, 10) + "'"
	}

	QueryCondition4 := ""
	if BedTypeCode != "" {
		QueryCondition4 = "bed_type_code='" + BedTypeCode + "' AND "
	}

	QueryCondition5 := ""
	if ReadyOnly {
		QueryCondition5 = "status_code='" + GlobalVar.RoomStatus.Ready + "' AND "
	}

	//Allotment Only
	QueryCondition1 = QueryCondition1 + " AND reservation.is_from_allotment=" + strconv.FormatBool(AllotmentOnly)
	QueryNot := ""
	QueryRoomAllotment := ""
	if AllotmentOnly {
		QueryCondition2A = QueryCondition2A + " AND folio.is_from_allotment=" + strconv.FormatBool(AllotmentOnly)
		QueryRoomAllotment =
			" WHERE DATE(room_allotment.from_date)<='" + ArrivalDateStr + "'" +
				" AND ADDDATE(DATE(room_allotment.to_date), INTERVAL 1 DAY)>='" + DepartureDateStr + "'"
	} else {
		QueryNot = "NOT "
		QueryRoomAllotment =
			" WHERE DATE(room_allotment.from_date)<'" + DepartureDateStr + "'" +
				" AND DATE(room_allotment.to_date)>='" + ArrivalDateStr + "'"
	}

	QueryCondition6 := ""
	if RoomAllotmentID != 0 {
		QueryCondition6 = " AND room_allotment.id<>'" + strconv.FormatUint(RoomAllotmentID, 10) + "'"
	}

	type DataStruct struct {
		Number string
		Name   string
	}

	var DataArray []DataStruct

	DB.Raw(
		"SELECT number, name FROM cfg_init_room" +
			" WHERE " + QueryConditionRoomType1 + QueryCondition4 + QueryCondition5 +
			"number NOT IN(SELECT" +
			" guest_detail.room_number " +
			"FROM" +
			" reservation" +
			" LEFT OUTER JOIN guest_detail ON (reservation.guest_detail_id = guest_detail.id)" +
			" WHERE reservation.status_code='" + GlobalVar.ReservationStatus.New + "'" +
			QueryConditionRoomType2 +
			" AND guest_detail.room_number<>''" +
			" AND DATE(guest_detail.arrival)<'" + DepartureDateStr + "'" +
			" AND DATE(guest_detail.departure)>'" + ArrivalDateStr + "'" +
			QueryCondition1 + ") " +
			"AND number NOT IN(SELECT" +
			" guest_detail.room_number " +
			"FROM" +
			" folio" +
			" LEFT OUTER JOIN guest_detail ON (folio.guest_detail_id = guest_detail.id)" +
			" WHERE folio.status_code='" + GlobalVar.FolioStatus.Open + "'" +
			" AND folio.type_code='" + GlobalVar.FolioType.GuestFolio + "'" +
			QueryConditionRoomType2 +
			" AND DATE(guest_detail.arrival)<'" + DepartureDateStr + "'" +
			" AND DATE(guest_detail.departure)>'" + ArrivalDateStr + "'" +
			QueryCondition2A + QueryCondition2B + ") " +
			"AND number NOT IN(SELECT" +
			" room_unavailable.room_number " +
			"FROM" +
			" room_unavailable" +
			" LEFT OUTER JOIN cfg_init_room ON (room_unavailable.room_number = cfg_init_room.number)" +
			" WHERE DATE(room_unavailable.start_date)<'" + DepartureDateStr + "'" +
			" AND DATE(room_unavailable.end_date)>='" + ArrivalDateStr + "'" +
			QueryConditionRoomType3 +
			QueryCondition3 + ") " +
			"AND number " + QueryNot + "IN(SELECT" +
			" room_allotment.room_number " +
			"FROM" +
			" room_allotment" +
			" LEFT OUTER JOIN cfg_init_room ON (room_allotment.room_number = cfg_init_room.number)" +
			QueryRoomAllotment +
			QueryConditionRoomType3 +
			QueryCondition6 + ") " +
			"ORDER BY id_sort, number;").Scan(&DataArray)

	var RoomNumberList []RoomNumberListStruct
	for _, Data := range DataArray {
		if Dataset.ProgramConfiguration.IsRoomByName {
			RoomNumberList = append(RoomNumberList, RoomNumberListStruct{Data.Name})
		} else {
			RoomNumberList = append(RoomNumberList, RoomNumberListStruct{Data.Number})
		}
	}
	return RoomNumberList
}

func GetRoomRate(DB *gorm.DB, Dataset *global_var.TDataset, RoomTypeCode, BusinessSourceCode, MarketCode, CompanyCode, ArrivalDate string) []map[string]interface{} {
	var DataArray []map[string]interface{}
	ArrivalDateStr := General.FormatDate1(General.StrZToDate(ArrivalDate))
	if RoomTypeCode == "" || ArrivalDate == "" {
		DB.Raw(
			"SELECT code, name, weekday_rate1, weekend_rate1 " +
				"FROM cfg_init_room_rate" +
				" WHERE room_type_code='';").Scan(&DataArray)
	} else {
		QueryConditionMarket := ""
		if Dataset.ProgramConfiguration.FilterRateByMarket {
			if MarketCode == "" {
				QueryConditionMarket = " AND IFNULL(cfg_init_room_rate.market_code, '')='' "
			} else {
				if Dataset.ProgramConfiguration.AlwaysShowPublishRate {
					fmt.Println("3")
					QueryConditionMarket = " AND (IFNULL(cfg_init_room_rate.market_code, '')='" + MarketCode + "' OR IFNULL(cfg_init_room_rate.market_code, '')='')"
				} else {
					fmt.Println("14")
					QueryConditionMarket = " AND IFNULL(cfg_init_room_rate.market_code, '')='" + MarketCode + "' "
				}
			}
		}

		QueryConditionCompany := ""
		if Dataset.ProgramConfiguration.FilterRateByCompany {
			if CompanyCode == "" {
				QueryConditionCompany = " AND cfg_init_room_rate.company_code=''"
			} else {
				if Dataset.ProgramConfiguration.AlwaysShowPublishRate {
					fmt.Println("1")
					QueryConditionCompany = " AND (cfg_init_room_rate.company_code='" + CompanyCode + "' OR cfg_init_room_rate.company_code='')"
				} else {
					fmt.Println("13")
					QueryConditionCompany = " AND cfg_init_room_rate.company_code='" + CompanyCode + "'"
				}
			}
		} else {
			if BusinessSourceCode == "" {
				QueryConditionCompany = " AND cfg_init_room_rate.company_code=''"
			} else {
				if Dataset.ProgramConfiguration.AlwaysShowPublishRate {
					fmt.Println("2")
					QueryConditionCompany = " OR (cfg_init_room_rate.company_code='" + BusinessSourceCode + "' OR cfg_init_room_rate.company_code='')"
				} else {
					fmt.Println("12")
					QueryConditionCompany = " OR cfg_init_room_rate.company_code='" + BusinessSourceCode + "'"
				}
			}
		}

		QueryConditionBusinessSource := ""
		if BusinessSourceCode == "" {
			QueryConditionBusinessSource = "IFNULL(cfg_init_room_rate_business_source.company_code,'')=''"
		} else {
			if Dataset.ProgramConfiguration.AlwaysShowPublishRate {
				fmt.Println("1")
				QueryConditionBusinessSource = " (cfg_init_room_rate_business_source.company_code='" + BusinessSourceCode + "' OR IFNULL(cfg_init_room_rate_business_source.company_code,'')='')"
			} else {
				fmt.Println("11")
				QueryConditionBusinessSource = " cfg_init_room_rate_business_source.company_code='" + BusinessSourceCode + "'"
			}
		}

		DB.Raw(
			"SELECT cfg_init_room_rate.code, cfg_init_room_rate.name, cfg_init_room_rate.weekday_rate1, cfg_init_room_rate.weekend_rate1 " +
				"FROM cfg_init_room_rate" +
				" LEFT OUTER JOIN cfg_init_room_rate_business_source ON (cfg_init_room_rate.code = cfg_init_room_rate_business_source.room_rate_code)" +
				" WHERE cfg_init_room_rate.room_type_code LIKE '%" + RoomTypeCode + "%' " +
				" AND cfg_init_room_rate.from_date<='" + ArrivalDateStr + "'" +
				" AND cfg_init_room_rate.to_date>='" + ArrivalDateStr + "'" +
				" AND cfg_init_room_rate.is_active=1" +
				" AND (" + QueryConditionBusinessSource + QueryConditionCompany + ")" +
				QueryConditionMarket + " " +
				"GROUP BY cfg_init_room_rate.code " +
				"ORDER BY cfg_init_room_rate.name;").Scan(&DataArray)
	}

	return DataArray
}

func GetRoomRateAmount(ctx context.Context, DB *gorm.DB, Dataset *GlobalVar.TDataset, RoomRateCode, PostingDateStr string, Adult, Child int, IsWeekend bool) float64 {
	ctx, span := global_var.Tracer.Start(ctx, "GetRoomRateAmount")
	defer span.End()

	var Rate float64 = 0
	var Pax int = 0
	type DataStruct struct {
		WeekdayRate1      float64
		WeekdayRate2      float64
		WeekdayRate3      float64
		WeekdayRate4      float64
		WeekendRate1      float64
		WeekendRate2      float64
		WeekendRate3      float64
		WeekendRate4      float64
		WeekdayRateChild1 float64
		WeekdayRateChild2 float64
		WeekdayRateChild3 float64
		WeekdayRateChild4 float64
		WeekendRateChild1 float64
		WeekendRateChild2 float64
		WeekendRateChild3 float64
		WeekendRateChild4 float64
		IncludeChild      bool
		ExtraPax          float64
		PerPax            bool
	}
	var Data DataStruct
	if RoomRateCode == "" {
		return 0
	}
	QueryResult := DB.WithContext(ctx).Raw(
		"SELECT" +
			" weekday_rate1," +
			" weekday_rate2," +
			" weekday_rate3," +
			" weekday_rate4," +
			" weekend_rate1," +
			" weekend_rate2," +
			" weekend_rate3," +
			" weekend_rate4," +
			" weekday_rate_child1," +
			" weekday_rate_child2," +
			" weekday_rate_child3," +
			" weekday_rate_child4," +
			" weekend_rate_child1," +
			" weekend_rate_child2," +
			" weekend_rate_child3," +
			" weekend_rate_child4," +
			" include_child," +
			" extra_pax," +
			" per_pax " +
			"FROM" +
			" cfg_init_room_rate" +
			" WHERE code='" + RoomRateCode + "'").Limit(1).Scan(&Data)

	if QueryResult.RowsAffected > 0 {
		Pax = Adult
		if Data.IncludeChild {
			Pax = Adult + Child
		}

		//If Null Date
		if PostingDateStr != "" {
			PostingDate, _ := time.Parse("2006-01-02", PostingDateStr)
			IsWeekend = General.IsWeekend(PostingDate, Dataset)
		}

		if IsWeekend {
			if Dataset.ProgramConfiguration.UseChildRate {
				if Adult == 1 {
					Rate = Data.WeekendRate1
				} else if Adult == 2 {
					Rate = Data.WeekendRate2
				} else if Adult == 3 {
					Rate = Data.WeekendRate3
				} else if Adult >= 4 {
					Rate = Data.WeekendRate4
				}
				//Child Rate
				if Child == 1 {
					Rate = Data.WeekendRateChild1 + Data.WeekendRateChild1
				} else if Child == 2 {
					Rate = Data.WeekendRateChild2 + Data.WeekendRateChild2
				} else if Child == 3 {
					Rate = Data.WeekendRateChild3 + Data.WeekendRateChild3
				} else if Child >= 4 {
					Rate = Data.WeekendRateChild4 + Data.WeekendRateChild4
				}
			} else {
				if Pax == 1 {
					Rate = Data.WeekendRate1
				} else if Pax == 2 {
					Rate = Data.WeekendRate2
				} else if Pax == 3 {
					Rate = Data.WeekendRate3
				} else if Pax >= 4 {
					Rate = Data.WeekendRate4
				}
			}
		} else {
			if Dataset.ProgramConfiguration.UseChildRate {
				if Adult == 1 {
					Rate = Data.WeekdayRate1
				} else if Adult == 2 {
					Rate = Data.WeekdayRate2
				} else if Adult == 3 {
					Rate = Data.WeekdayRate3
				} else if Adult >= 4 {
					Rate = Data.WeekdayRate4
				}

				//Child Rate
				if Child == 1 {
					Rate = Rate + Data.WeekdayRateChild1
				} else if Child == 2 {
					Rate = Rate + Data.WeekdayRateChild2
				} else if Child == 3 {
					Rate = Rate + Data.WeekdayRateChild3
				} else if Child >= 4 {
					Rate = Rate + Data.WeekdayRateChild4
				}
			} else {
				if Pax == 1 {
					Rate = Data.WeekdayRate1
				} else if Pax == 2 {
					Rate = Data.WeekdayRate2
				} else if Pax == 3 {
					Rate = Data.WeekdayRate3
				} else if Pax >= 4 {
					Rate = Data.WeekdayRate4
				}
			}

			//Extra Pax
			if Pax > 4 {
				if Data.PerPax {
					Rate = Rate + (float64(Pax-4) * Data.ExtraPax)
				} else {
					Rate = Rate + Data.ExtraPax
				}
			}
		}
	}
	return Rate
}

func CopyGuestProfileToContactPerson(GuestProfileData DBVar.Guest_profile) DBVar.Contact_person {
	var ContactPersonData DBVar.Contact_person
	ContactPersonData.TitleCode = &GuestProfileData.TitleCode
	ContactPersonData.FullName = &GuestProfileData.FullName
	ContactPersonData.Street = &GuestProfileData.Street
	ContactPersonData.CountryCode = &GuestProfileData.CountryCode
	ContactPersonData.StateCode = &GuestProfileData.StateCode
	ContactPersonData.CityCode = &GuestProfileData.CityCode
	ContactPersonData.City = &GuestProfileData.City
	ContactPersonData.NationalityCode = &GuestProfileData.NationalityCode
	ContactPersonData.PostalCode = &GuestProfileData.PostalCode
	ContactPersonData.Phone1 = &GuestProfileData.Phone1
	ContactPersonData.Phone2 = &GuestProfileData.Phone2
	ContactPersonData.Fax = &GuestProfileData.Fax
	ContactPersonData.Email = &GuestProfileData.Email
	ContactPersonData.Website = &GuestProfileData.Website
	ContactPersonData.CompanyCode = &GuestProfileData.CompanyCode
	ContactPersonData.GuestTypeCode = &GuestProfileData.GuestTypeCode
	ContactPersonData.IdCardCode = &GuestProfileData.IdCardCode
	ContactPersonData.IdCardNumber = &GuestProfileData.IdCardNumber
	ContactPersonData.IsMale = &GuestProfileData.IsMale
	ContactPersonData.BirthPlace = &GuestProfileData.BirthPlace
	ContactPersonData.BirthDate = &GuestProfileData.BirthDate
	ContactPersonData.TypeCode = GuestProfileData.TypeCode
	ContactPersonData.CustomField01 = &GuestProfileData.CustomField01
	ContactPersonData.CustomField02 = &GuestProfileData.CustomField02
	ContactPersonData.CustomField03 = &GuestProfileData.CustomField03
	ContactPersonData.CustomField04 = &GuestProfileData.CustomField04
	ContactPersonData.CustomField05 = &GuestProfileData.CustomField05
	ContactPersonData.CustomField06 = &GuestProfileData.CustomField06
	ContactPersonData.CustomField07 = &GuestProfileData.CustomField07
	ContactPersonData.CustomField08 = &GuestProfileData.CustomField08
	ContactPersonData.CustomField09 = &GuestProfileData.CustomField09
	ContactPersonData.CustomField10 = &GuestProfileData.CustomField10
	ContactPersonData.CustomField11 = &GuestProfileData.CustomField11
	ContactPersonData.CustomField12 = &GuestProfileData.CustomField12
	ContactPersonData.CustomLookupFieldCode01 = &GuestProfileData.CustomLookupFieldCode01
	ContactPersonData.CustomLookupFieldCode02 = &GuestProfileData.CustomLookupFieldCode02
	ContactPersonData.CustomLookupFieldCode03 = &GuestProfileData.CustomLookupFieldCode03
	ContactPersonData.CustomLookupFieldCode04 = &GuestProfileData.CustomLookupFieldCode04
	ContactPersonData.CustomLookupFieldCode05 = &GuestProfileData.CustomLookupFieldCode05
	ContactPersonData.CustomLookupFieldCode06 = &GuestProfileData.CustomLookupFieldCode06
	ContactPersonData.CustomLookupFieldCode07 = &GuestProfileData.CustomLookupFieldCode07
	ContactPersonData.CustomLookupFieldCode08 = &GuestProfileData.CustomLookupFieldCode08
	ContactPersonData.CustomLookupFieldCode09 = &GuestProfileData.CustomLookupFieldCode09
	ContactPersonData.CustomLookupFieldCode10 = &GuestProfileData.CustomLookupFieldCode10
	ContactPersonData.CustomLookupFieldCode11 = &GuestProfileData.CustomLookupFieldCode11
	ContactPersonData.CustomLookupFieldCode12 = &GuestProfileData.CustomLookupFieldCode12
	ContactPersonData.CreatedBy = GuestProfileData.CreatedBy
	ContactPersonData.UpdatedBy = GuestProfileData.UpdatedBy

	return ContactPersonData
}

func SaveGuestProfileContactPerson(Db *gorm.DB, ValidUserCode string, GuestProfileData DBVar.Guest_profile, ContactPersonId, GuestProfileId uint64) (uint64, uint64, error) {
	GuestProfileData.TypeCode = GlobalVar.CPType.Guest
	GuestProfileData.UpdatedBy = ValidUserCode
	GuestProfileData.CreatedBy = ValidUserCode
	ContactPersonData := CopyGuestProfileToContactPerson(GuestProfileData)
	if GuestProfileData.FullName != "" {
		// insert
		OmitX := "updated_by"
		OmitY := "updated_by"
		OmitSource := ""

		//update
		if GuestProfileId != 0 {
			GuestProfileData.Id = GuestProfileId
			ContactPersonData.Id = ContactPersonId
			OmitX = "created_by"
			OmitSource = "source"
		}
		if ContactPersonId != 0 {
			GuestProfileData.Id = GuestProfileId
			ContactPersonData.Id = ContactPersonId
			OmitY = "created_by"
		}

		// InsertGuestProfile(Db, GuestProfileData.TitleCode,GuestProfileData.FullName,GuestProfileData.Street,GuestProfileData.CountryCode,GuestProfileData.StateCode,GuestProfileData.CityCode,
		// 	GuestProfileData.City,GuestProfileData.NationalityCode,GuestProfileData.PostalCode,GuestProfileData.Phone1,GuestProfileData.Phone2,GuestProfileData.Fax,GuestProfileData.Email,GuestProfileData.Website,
		// 	GuestProfileData.CompanyCode,GuestProfileData.GuestTypeCode,GuestProfileData.IdCardCode,GuestProfileData.IdCardNumber,GuestProfileData.IsMale,GuestProfileData.BirthPlace,
		// 	GuestProfileData.BirthDate,GuestProfileData.TypeCode,GuestProfileData.CustomField01, GuestProfileData.CustomField02,GuestProfileData.CustomField03,GuestProfileData.CustomField04,
		// 	GuestProfileData.CustomField05,GuestProfileData.CustomField06,GuestProfileData.CustomField07,GuestProfileData.CustomField08,GuestProfileData.CustomField09,GuestProfileData.CustomField10,
		// 	GuestProfileData.CustomField11,GuestProfileData.CustomField12,GuestProfileData.CustomLookupFieldCode01,GuestProfileData.CustomLookupFieldCode02,GuestProfileData.CustomLookupFieldCode03,
		// 	GuestProfileData.CustomLookupFieldCode04,GuestProfileData.CustomLookupFieldCode05,GuestProfileData.CustomLookupFieldCode06,GuestProfileData.CustomLookupFieldCode07,GuestProfileData.CustomLookupFieldCode08,
		// 	GuestProfileData.CustomLookupFieldCode09,GuestProfileData.CustomLookupFieldCode10,GuestProfileData.CustomLookupFieldCode11,GuestProfileData.CustomLookupFieldCode12,1,0,GuestProfileData.CustomerCode,GuestProfileData.Source,ValidUserCode)

		err := Db.Table(DBVar.TableName.GuestProfile).Omit(OmitX, OmitSource).Save(&GuestProfileData).Error
		if err != nil {
			return 0, 0, err
		}
		err = Db.Table(DBVar.TableName.ContactPerson).Omit(OmitY).Save(&ContactPersonData).Error
		if err != nil {
			return 0, 0, err
		}
	}
	return GuestProfileData.Id, ContactPersonData.Id, nil
}

// func IsFolioAutoTransferAccountP(c *gin.Context) {
// 	type DataInputStruct struct {
// 		FolioNumber int64
// 		AccountCode string
// 	}

// 	var DataInput DataInputStruct
// 	err := c.BindJSON(&DataInput)
// 	if err != nil {
// 		MasterData.SendResponse(GlobalVar.ResponseCode.InvalidDataFormat, "", nil, c)
// 	} else {
// 		var FolioTransfer []string
// 		var IsFolioAutoTransferAccount bool
// 		DB.Table(DBVar.TableName.FolioRouting).Select("folio_transfer").Where("folio_number = ? AND account_code = ?", DataInput.FolioNumber, DataInput.AccountCode).Find(&FolioTransfer)
// 		IsFolioAutoTransferAccount = len(FolioTransfer) > 0
// 		MasterData.SendResponse(GlobalVar.ResponseCode.Successfully, "", IsFolioAutoTransferAccount, c)
// 	}
// }

func GetDefaultCurrencyCode(DB *gorm.DB) string {
	Result := ""
	var Code []string
	DB.Table(DBVar.TableName.CfgInitCurrency).Select("code").Where("is_default = ?", 1).Find(&Code)

	if len(Code) > 0 {
		Result = Code[0]
		return Result
	}

	return Result
}

func GetExchangeRateCurrency(DB *gorm.DB, CurrencyCode string) float64 {
	return MasterData.GetFieldFloat(DB, DBVar.TableName.CfgInitCurrency, "exchange_rate", "code", CurrencyCode, 1)
}

func GetGuestDepositCorrectionBreakDown(DB *gorm.DB) uint64 {
	var Result uint64 = 1
	var CorrectionBreakdown []uint64
	DB.Table(DBVar.TableName.GuestDeposit).Select("correction_breakdown").Order("correction_breakdown desc").Limit(1).Find(&CorrectionBreakdown)

	if len(CorrectionBreakdown) > 0 {
		Result = CorrectionBreakdown[0] + 1
		return Result
	}
	return Result
}

func GetAccountSubGroupCode(DB *gorm.DB, AccountCode string) string {
	return MasterData.GetFieldString(DB, DBVar.TableName.CfgInitAccount, "sub_group_code", "code", AccountCode, "")
}

func GetTotalDepositReservation(DB *gorm.DB, ReservationNumber uint64) float64 {
	return MasterData.GetFieldFloatQuery(DB,
		"SELECT"+
			" SUM(IF(type_code='C', amount, -amount)) AS TotalDeposit "+
			"FROM"+
			" guest_deposit"+
			" WHERE reservation_number=?"+
			" AND void='0'"+
			" AND system_code='"+GlobalVar.ConstProgramVariable.DefaultSystemCode+"' "+
			"GROUP BY reservation_number", 0, ReservationNumber)
}

func UpdateReservationRoomNumber(c *gin.Context, DB *gorm.DB, GuestDetailID uint64, RoomNumber string, UpdatedBy string) string {
	// Get Program Configuration
	val, exist := c.Get("pConfig")
	if !exist {
		//
	}
	pConfig := val.(*config.CompanyDataConfiguration)
	if pConfig.Dataset.ProgramConfiguration.IsRoomByName {
		RoomNumber = MasterData.GetFieldString(DB, DBVar.TableName.CfgInitRoom, "number", "name", RoomNumber, "")
	}
	// if MasterData.GetConfigurationBool(DB, GlobalVar.SystemCode.Hotel, GlobalVar.ConfigurationCategory.General, GlobalVar.ConfigurationName.IsRoomByName, false) {
	// 	RoomNumber = MasterData.GetFieldString(DBVar.TableName.CfgInitRoom, "number", "name", RoomNumber, "")
	// }

	DB.Table(DBVar.TableName.GuestDetail).Select("room_number", "updated_by").Where("id = ?", GuestDetailID).Updates(&map[string]interface{}{
		"room_number": RoomNumber,
		"updated_by":  UpdatedBy,
	})
	return RoomNumber
}

func AssignRoom(c *gin.Context, DB *gorm.DB, Dataset *GlobalVar.TDataset, ReservationNumber uint64, BedTypeCode string, ReadyOnly bool, UpdatedBy string) string {
	var RoomNumber string
	var DataOutput map[string]interface{}
	DB.Raw(
		"SELECT"+
			" guest_detail.id AS id,"+
			" DATE(guest_detail.arrival) AS DateArrival,"+
			" DATE(guest_detail.departure) AS DateDeparture,"+
			" reservation.is_lock AS is_lock,"+
			" reservation.is_from_allotment AS is_from_allotment,"+
			" guest_detail.room_type_code AS room_type_code "+
			"FROM"+
			" reservation"+
			" LEFT OUTER JOIN guest_detail ON (reservation.guest_detail_id = guest_detail.id)"+
			" WHERE reservation.number=?"+
			" AND reservation.status_code=?"+
			" AND guest_detail.room_number='' "+
			"ORDER BY guest_detail.room_number", ReservationNumber, GlobalVar.ReservationStatus.New).Find(&DataOutput)

	if len(DataOutput) > 0 {
		RoomList := GetAvailableRoomByType(DB, Dataset, (DataOutput["DateArrival"].(time.Time)), (DataOutput["DateDeparture"].(time.Time)), DataOutput["room_type_code"].(string), BedTypeCode, 0, 0, 0, 0, ReadyOnly, General.InterfaceToBool(DataOutput["is_from_allotment"]))
		if len(RoomList) > 0 {
			RoomNumber = RoomList[0].RoomNumber
			if RoomNumber != "" {
				UpdateReservationRoomNumber(c, DB, General.InterfaceToUint64(DataOutput["id"]), RoomNumber, UpdatedBy)
				return RoomNumber
			}
		}
	}
	return ""
}

func GetBasicTaxService(DB *gorm.DB, AccountCode, TaxAndServiceCodeManual string, Amount float64) (Basic, Tax, Service float64) {
	Basic = Amount
	var TaxServiceCode string
	var TaxPercent, ServicePercent, ServiceTaxPercent float64
	var IsTaxIncluded, IsServiceIncluded bool
	isDebug := false

	GroupCode := MasterData.GetAccountField(DB, "group_code", "cfg_init_account.code", AccountCode)
	if GroupCode == GlobalVar.GlobalAccountGroup.Charge {
		if TaxAndServiceCodeManual == "" {
			TaxServiceCode = MasterData.GetAccountField(DB, "tax_and_service_code", "cfg_init_account.code", AccountCode)
		} else {
			TaxServiceCode = TaxAndServiceCodeManual
		}
		if isDebug {
			fmt.Println("txCode", TaxServiceCode)
		}
		if TaxServiceCode != "" {
			var TaxAndServiceData DBVar.Cfg_init_tax_and_service
			DB.Table(DBVar.TableName.CfgInitTaxAndService).Where("code", TaxServiceCode).Limit(1).Scan(&TaxAndServiceData)

			TaxPercent = TaxAndServiceData.Tax                         //MasterData.GetFieldFloat(DB, DBVar.TableName.CfgInitTaxAndService, "tax", "code", TaxServiceCode, 0)
			ServicePercent = TaxAndServiceData.Service                 // MasterData.GetFieldFloat(DB, DBVar.TableName.CfgInitTaxAndService, "service", "code", TaxServiceCode, 0)
			ServiceTaxPercent = TaxAndServiceData.ServiceTax           //MasterData.GetFieldFloat(DB, DBVar.TableName.CfgInitTaxAndService, "service_tax", "code", TaxServiceCode, 0)
			IsTaxIncluded = TaxAndServiceData.IsTaxInclude > 0         //MasterData.GetFieldBool(DB, DBVar.TableName.CfgInitTaxAndService, "is_tax_include", "code", TaxServiceCode, false)
			IsServiceIncluded = TaxAndServiceData.IsServiceInclude > 0 //MasterData.GetFieldBool(DB, DBVar.TableName.CfgInitTaxAndService, "is_service_include", "code", TaxServiceCode, false)
			//Tax and Service Include
			if IsTaxIncluded && IsServiceIncluded {
				Tax = General.RoundToX3(Amount * (TaxPercent + (ServiceTaxPercent * ServicePercent / 100)) / (100 + TaxPercent + ServicePercent + (ServiceTaxPercent * ServicePercent / 100)))
				Service = General.RoundToX3(Amount * ServicePercent / (100 + TaxPercent + ServicePercent + (ServiceTaxPercent * ServicePercent / 100)))
				Basic = Amount - Tax - Service
				if isDebug {
					fmt.Println("1")
					fmt.Println(Basic)
					fmt.Println(Amount)
					fmt.Println(Tax)
					fmt.Println(Service)
				}
				//Tax and Service Exclude
			} else if !IsTaxIncluded && !IsServiceIncluded {
				Basic = Amount
				Tax = General.RoundToX3(Amount * (TaxPercent + (ServiceTaxPercent * ServicePercent / 100)) / 100)
				Service = General.RoundToX3(Amount * ServicePercent / 100)
				if isDebug {
					fmt.Println("2")
					fmt.Println(Basic)
					fmt.Println(Amount)
					fmt.Println(Tax)
					fmt.Println(Service)
				}
				//Tax Exclude and Service Include
			} else if !IsTaxIncluded && IsServiceIncluded {
				Service = General.RoundToX3(Amount / (100 + ServicePercent) * ServicePercent)
				Basic = Amount - Service
				Tax = General.RoundToX3(Basic * TaxPercent / 100)
				if ServiceTaxPercent > 0 {
					Tax = General.RoundToX3(Amount * TaxPercent / 100)
				}
				Tax = General.RoundToX3(Amount * (TaxPercent + (ServiceTaxPercent * ServicePercent / 100)) / 100)
				if isDebug {
					fmt.Println("3")
					fmt.Println(Basic)
					fmt.Println(Amount)
					fmt.Println(Tax)
					fmt.Println(Service)
				}
				//Tax Include and Service Exclude
			} else if IsTaxIncluded && !IsServiceIncluded {
				Tax = General.RoundToX3(Amount / (100 + TaxPercent + (ServiceTaxPercent * ServicePercent / 100)) * (TaxPercent + (ServiceTaxPercent * ServicePercent / 100)))
				Basic = Amount - Tax
				Service = General.RoundToX3((Amount - (Amount / (100 + TaxPercent + (ServiceTaxPercent * ServicePercent / 100)) * (TaxPercent + (ServiceTaxPercent * ServicePercent / 100)))) * ServicePercent / 100)

				if isDebug {
					fmt.Println("4")
					fmt.Println(Basic)
					fmt.Println(Amount)
					fmt.Println(Tax)
					fmt.Println(Service)
				}
			}
		}
	}
	return Basic, Tax, Service
}

func GetBasicTaxService2(DB *gorm.DB, AccountCode, TaxServiceCodeManual string, Amount float64) (Basic, Tax, Service float64) {
	var TaxServiceCode string
	var TaxPercent, ServicePercent, ServiceTaxPercent float64
	// var IsTaxIncluded, IsServiceIncluded bool

	type DataOutputAccountStruct struct {
		AccountGroupCode  string
		TaxAndServiceCode string
		Tax               float64
		Service           float64
		ServiceTax        float64
	}

	var DataOutputAccount DataOutputAccountStruct
	DB.Table(DBVar.TableName.CfgInitAccount).
		Select("cfg_init_account.tax_and_service_code",
			"cfg_init_account_sub_group.group_code as account_group_code",
			"cfg_init_tax_and_service.tax",
			"cfg_init_tax_and_service.service",
			"cfg_init_tax_and_service.service_tax").
		Joins("LEFT JOIN cfg_init_account_sub_group ON (cfg_init_account.sub_group_code=cfg_init_account_sub_group.code)").
		Joins("LEFT JOIN cfg_init_tax_and_service ON (cfg_init_account.tax_and_service_code=cfg_init_tax_and_service.code)").
		Where("cfg_init_account.code = ?", AccountCode).Take(&DataOutputAccount)

	if DataOutputAccount.AccountGroupCode == GlobalVar.GlobalAccountGroup.Charge {
		if TaxServiceCodeManual == "" {
			TaxServiceCode = DataOutputAccount.TaxAndServiceCode
		} else {
			TaxServiceCode = TaxServiceCodeManual
		}
		//fmt.Println(TaxServiceCode)

		if TaxServiceCode != "" {
			TaxPercent = DataOutputAccount.Tax
			ServicePercent = DataOutputAccount.Service
			ServiceTaxPercent = DataOutputAccount.ServiceTax
			// IsTaxIncluded = DataOutputAccount.TaxAndService.IsTaxInclude > 0
			// IsServiceIncluded = DataOutputAccount.TaxAndService.IsServiceInclude > 0

			Tax = General.RoundToX3(Amount * (TaxPercent + (ServiceTaxPercent * ServicePercent / 100)) / (100 + TaxPercent + ServicePercent + (ServiceTaxPercent * ServicePercent / 100)))
			Service = General.RoundToX3(Amount * ServicePercent / (100 + TaxPercent + ServicePercent + (ServiceTaxPercent * ServicePercent / 100)))
			Basic = Amount - Tax - Service
		}
	}
	return Basic, Tax, Service
}

func GetBasicTaxServiceForeign(Tax, Service, AmountForeign, ExchangeRate float64) (BasicForeign float64, TaxForeign float64, ServiceForeign float64) {
	// var TaxForeign, ServiceForeign, BasicForeign float64
	if Tax > 0 {
		TaxForeign = General.RoundToX3(Tax / ExchangeRate)
	} else {
		TaxForeign = 0
	}

	if Service > 0 {
		ServiceForeign = General.RoundToX3(Service / ExchangeRate)
	} else {
		ServiceForeign = AmountForeign
	}
	BasicForeign = AmountForeign

	return BasicForeign, TaxForeign, ServiceForeign
}

func GetSubFolioCorrectionBreakdown(c *gin.Context, DB *gorm.DB) uint64 {
	var CorrectionBreakdown uint64
	uuid := time.Now().UnixMilli()
	return uint64(uuid)
	AuditDate := GetAuditDate(c, DB, false)
	DB.Raw("SELECT" +
		" A.correction_breakdown " +
		"FROM " +
		"(SELECT correction_breakdown FROM sub_folio " +
		"WHERE sub_folio.audit_date BETWEEN  DATE(DATE_ADD('" + General.FormatDate1(AuditDate) + "', INTERVAL -1 DAY)) AND DATE('" + General.FormatDate1(AuditDate) + "') " +
		") AS A " +
		"ORDER BY A.correction_breakdown DESC " +
		"LIMIT 1;").Scan(&CorrectionBreakdown)

	if CorrectionBreakdown <= 0 {
		DB.Table(DBVar.TableName.SubFolio).Select("correction_breakdown").Order("correction_breakdown DESC").Limit(1).Scan(&CorrectionBreakdown)
	}

	return CorrectionBreakdown + 1
}

func GetSubFolioBreakdown1(c *gin.Context, DB *gorm.DB) uint64 {
	var Breakdown1 uint64
	uuid := time.Now().UnixMilli()
	return uint64(uuid)
	AuditDate := GetAuditDate(c, DB, false)
	DB.Raw("SELECT" +
		" A.breakdown1 " +
		"FROM " +
		"(SELECT breakdown1 FROM sub_folio " +
		"WHERE sub_folio.audit_date BETWEEN  DATE(DATE_ADD('" + General.FormatDate1(AuditDate) + "', INTERVAL -1 DAY)) AND DATE('" + General.FormatDate1(AuditDate) + "') " +
		") AS A " +
		"ORDER BY A.breakdown1 DESC " +
		"LIMIT 1;").Scan(&Breakdown1)

	if Breakdown1 <= 0 {
		DB.Table(DBVar.TableName.SubFolio).Select("breakdown1").Order("breakdown1 DESC").Limit(1).Scan(&Breakdown1)
	}

	return Breakdown1 + 1
}

func GetSubFolioBreakdown2(c *gin.Context, DB *gorm.DB, Breakdown1 uint64) (Breakdown2 int) {
	uuid := time.Now().UnixMilli()
	return int(uuid)
	AuditDate := GetAuditDate(c, DB, false)
	DB.Raw("SELECT" +
		" A.breakdown2 " +
		"FROM " +
		"(SELECT breakdown2 FROM sub_folio " +
		"WHERE breakdown1='" + strconv.FormatUint(Breakdown1, 10) + "' " +
		"AND sub_folio.audit_date BETWEEN  DATE(DATE_ADD('" + General.FormatDate1(AuditDate) + "', INTERVAL -1 DAY)) AND DATE('" + General.FormatDate1(AuditDate) + "') " +
		") AS A " +
		"ORDER BY A.breakdown2 DESC " +
		"LIMIT 1;").Scan(&Breakdown2)

	return Breakdown2 + 1
}

func GetRoomNumberFromConfigurationIsRoomName(c *gin.Context, DB *gorm.DB, RoomName string) string {
	Result := RoomName

	// Get Program Configuration
	val, exist := c.Get("pConfig")
	if !exist {
		//
	}
	pConfig := val.(*config.CompanyDataConfiguration)
	if General.StrToBool(pConfig.Dataset.Configuration[GlobalVar.ConfigurationCategory.General][GlobalVar.ConfigurationName.IsRoomByName].(string)) {
		Result = MasterData.GetFieldString(DB, DBVar.TableName.CfgInitRoom, "number", "name", RoomName, RoomName)
	}
	return Result
}

func GetRoomNameFromConfigurationIsRoomName(c *gin.Context, DB *gorm.DB, RoomNumber string) string {
	Result := RoomNumber

	// Get Program Configuration
	val, exist := c.Get("pConfig")
	if !exist {
		// MasterData.SendResponse(global_var.ResponseCode.DatabaseError, "pConfig", nil, c)
		// return
	}
	pConfig := val.(*config.CompanyDataConfiguration)

	IsRoomByName := General.StrToBool(pConfig.Dataset.Configuration[GlobalVar.ConfigurationCategory.General][GlobalVar.ConfigurationName.IsRoomByName].(string))
	if IsRoomByName {
		Result = MasterData.GetFieldString(DB, DBVar.TableName.CfgInitRoom, "name", "number", RoomNumber, RoomNumber)
	}
	return Result
}

func GetAutoRouting(DB *gorm.DB, BelongsTo uint64, AccountCode string) (uint64, string) {
	var FolioRouting DBVar.Folio_routing
	DB.Table(DBVar.TableName.FolioRouting).Select("folio_transfer, sub_folio_transfer").Where("folio_number=?", BelongsTo).Where("account_code=?", AccountCode).Limit(1).Scan(&FolioRouting)

	return FolioRouting.FolioTransfer, FolioRouting.SubFolioTransfer
}

func InsertSubFolioX(c *gin.Context, Dataset *GlobalVar.TDataset, IDCorrected uint64, DataInput DBVar.Sub_folio, DB *gorm.DB) (uint64, error) {
	var err error
	var Id uint64

	DataInput.DefaultCurrencyCode = GetDefaultCurrencyCode(DB)
	if DataInput.CurrencyCode == "" {
		DataInput.CurrencyCode = DataInput.DefaultCurrencyCode
	}
	DataInput.ExchangeRate = GetExchangeRateCurrency(DB, DataInput.CurrencyCode)
	DataInput.AmountForeign = DataInput.Amount

	if DataInput.CurrencyCode != DataInput.DefaultCurrencyCode {
		DataInput.Amount = General.RoundToX3(DataInput.Amount * DataInput.ExchangeRate)
	}
	IsRoomByName := General.StrToBool(Dataset.Configuration[GlobalVar.ConfigurationCategory.General][GlobalVar.ConfigurationName.IsRoomByName].(string))
	if IsRoomByName {
		DataInput.RoomNumber = GetRoomNumberFromConfigurationIsRoomName(c, DB, DataInput.RoomNumber)
	}
	if DataInput.AuditDate.IsZero() {
		AuditDate := GetAuditDate(c, DB, false)
		DataInput.AuditDate = AuditDate
		DataInput.AuditDateUnixx = int(AuditDate.Unix())
	}
	DataInput.Id = 0
	if err := DB.Table(DBVar.TableName.SubFolio).Omit("id").Create(&DataInput).Error; err != nil {
		return Id, err
	}
	Id = DataInput.Id
	// Insert Foreign Cash
	if (GetAccountSubGroupCode(DB, DataInput.AccountCode) == GlobalVar.GlobalAccountSubGroup.Payment || GetAccountSubGroupCode(DB, DataInput.AccountCode) == GlobalVar.GlobalAccountSubGroup.CreditDebitCard || GetAccountSubGroupCode(DB, DataInput.AccountCode) == GlobalVar.GlobalAccountSubGroup.BankTransfer) && DataInput.CurrencyCode != DataInput.DefaultCurrencyCode {
		RemarkForeignCash := "Payment for Folio: " + strconv.FormatUint(DataInput.FolioNumber, 10) + ", Room: " + DataInput.RoomNumber + ", Doc#: " + DataInput.DocumentNumber
		TypeCodeX := GlobalVar.TransactionType.Debit
		if DataInput.TypeCode == GlobalVar.TransactionType.Debit {
			TypeCodeX = GlobalVar.TransactionType.Credit
		}
		if General.Uint8ToBool(DataInput.IsCorrection) {
			RemarkForeignCash = "Payment Correction  for Folio: " + strconv.FormatUint(DataInput.FolioNumber, 10) + ", Room: " + DataInput.RoomNumber + ", Doc#: " + DataInput.DocumentNumber
		}

		var ForeignCash DBVar.Acc_foreign_cash
		ForeignCash.IdTransaction = DataInput.FolioNumber
		ForeignCash.IdCorrected = IDCorrected
		ForeignCash.IdChange = 0
		ForeignCash.IdTable = GlobalVar.ForeignCashTableID.SubFolio
		ForeignCash.Breakdown = 0
		ForeignCash.RefNumber = ""
		ForeignCash.Date = GlobalVar.ProgramVariable.AuditDate
		ForeignCash.TypeCode = TypeCodeX
		ForeignCash.Amount = DataInput.Amount
		ForeignCash.Stock = DataInput.AmountForeign
		ForeignCash.DefaultCurrencyCode = DataInput.DefaultCurrencyCode
		ForeignCash.AmountForeign = DataInput.AmountForeign
		ForeignCash.ExchangeRate = DataInput.ExchangeRate
		ForeignCash.CurrencyCode = DataInput.CurrencyCode
		ForeignCash.Remark = RemarkForeignCash
		ForeignCash.IsCorrection = DataInput.IsCorrection
		ForeignCash.CreatedBy = DataInput.CreatedBy
		if err := DB.Table(DBVar.TableName.AccForeignCash).Create(&ForeignCash).Error; err != nil {
			return Id, err
		}
	}
	return Id, err
}

func InsertSubFolio(c *gin.Context, DB *gorm.DB, Dataset *global_var.TDataset, IsCanAutoTransfer bool, AccountCodeTransfer, TaxServiceManual string, DataInput DBVar.Sub_folio) (Result string, SubFolioID uint64, Error error) {
	var SubFolioId uint64
	var err error

	v, e := c.Get("UserInfo")
	if !e {

		return "", 0, errors.New("UserInfo cannot be found")
	}
	UserInfo := v.(global_var.TUserInfo)
	DataInput.Shift = UserInfo.WorkingShift
	DataInput.LogShiftId = UserInfo.LogShiftID

	if DataInput.AuditDate.IsZero() {
		AuditDate := GetAuditDate(c, DB, false)
		DataInput.AuditDate = AuditDate
		DataInput.AuditDateUnixx = int(AuditDate.Unix())
	}
	if DataInput.CurrencyCode == "" {
		DataInput.CurrencyCode = GetDefaultCurrencyCode(DB)
		DataInput.ExchangeRate = GetExchangeRateCurrency(DB, DataInput.CurrencyCode)
	}
	if DataInput.ExchangeRate == 0 {
		DataInput.ExchangeRate = GetExchangeRateCurrency(DB, DataInput.CurrencyCode)
	}
	AllowZeroAmount := Dataset.Configuration[GlobalVar.ConfigurationCategory.Folio][GlobalVar.ConfigurationName.AllowZeroAmount].(string)
	if ((DataInput.Quantity*DataInput.Amount) > 0 || ((DataInput.Quantity*DataInput.Amount) <= 0) && AllowZeroAmount != "0") && DataInput.CurrencyCode != "" && DataInput.ExchangeRate > 0 {

		DataInput.DefaultCurrencyCode = GetDefaultCurrencyCode(DB)
		if DataInput.CurrencyCode != DataInput.DefaultCurrencyCode {
			DataInput.Amount = General.RoundToX3(DataInput.Amount * DataInput.ExchangeRate)
		}
		DataInput.BelongsTo = DataInput.FolioNumber
		//Cek auto routing
		if IsCanAutoTransfer {
			FolioNumber, SubFolio := GetAutoRouting(DB, DataInput.BelongsTo, AccountCodeTransfer)
			if FolioNumber != 0 && SubFolio != "" {
				DataInput.FolioNumber = FolioNumber
				DataInput.GroupCode = SubFolio
			}

		}
		var Tax, Service, TaxForeign, ServiceForeign float64
		DataInput.Amount, Tax, Service = GetBasicTaxService(DB, DataInput.AccountCode, TaxServiceManual, DataInput.Amount)
		DataInput.AmountForeign = General.RoundToX3(DataInput.Amount / DataInput.ExchangeRate)
		DataInput.AmountForeign, TaxForeign, ServiceForeign = GetBasicTaxServiceForeign(Tax, Service, DataInput.AmountForeign, DataInput.ExchangeRate)
		if DataInput.CorrectionBreakdown == 0 {
			DataInput.CorrectionBreakdown = GetSubFolioCorrectionBreakdown(c, DB)
		}
		if DataInput.Breakdown1 == 0 {
			DataInput.Breakdown1 = GetSubFolioBreakdown1(c, DB)
		}
		DataInput.Breakdown2 = GetSubFolioBreakdown2(c, DB, DataInput.Breakdown1)
		AccountSubGroup := GetAccountSubGroupCode(DB, DataInput.AccountCode)
		DataInput.CorrectionBy = ""
		DataInput.CorrectionReason = ""
		//Begin Transaction
		err = DB.Transaction(func(tx *gorm.DB) error {
			//Jika Direct Bill atau Compliment atau AP Commission
			if AccountSubGroup == GlobalVar.GlobalAccountSubGroup.AccountReceivable || AccountSubGroup == GlobalVar.GlobalAccountSubGroup.AccountPayable || AccountSubGroup == GlobalVar.GlobalAccountSubGroup.Compliment {
				SubFolioId, err = InsertSubFolioX(c, Dataset, 0, DataInput, tx)
				if err != nil {
					return err
				}
			} else {
				DataInput.DirectBillCode = ""
				SubFolioId, err = InsertSubFolioX(c, Dataset, 0, DataInput, tx)
				if err != nil {
					return err
				}
			}
			if err != nil {
				// return nil
			}

			if Tax > 0 {
				DataInput.AccountCode = Dataset.GlobalAccount.Tax
				DataInput.ProductCode = DataInput.AccountCode
				DataInput.Amount = Tax
				DataInput.AmountForeign = TaxForeign
				_, err = InsertSubFolioX(c, Dataset, 0, DataInput, tx)
				if err != nil {
					return err
				}
			}

			if Service > 0 {
				DataInput.AccountCode = Dataset.GlobalAccount.Service
				DataInput.ProductCode = DataInput.AccountCode
				DataInput.Amount = Service
				DataInput.AmountForeign = ServiceForeign
				_, err = InsertSubFolioX(c, Dataset, 0, DataInput, tx)
				if err != nil {
					return err
				}
			}
			return nil
		})

		if DataInput.FolioNumber == DataInput.BelongsTo {
			Result = ""
		} else {
			Result = strconv.FormatUint(DataInput.FolioNumber, 10) + "/" + DataInput.GroupCode
		}
	}
	return Result, SubFolioId, err
}

func UpdateGuestDepositTransferPairIDWithFolio(Db *gorm.DB, GuestDepositID, SubFolioID uint64, ValidUserCode string) error {
	if err := Db.Exec("CALL update_guest_deposit_transfer_pair_id(?, ?, ?)", GuestDepositID, 1, SubFolioID).Error; err != nil {
		return err
	}

	if err := Db.Exec("CALL update_sub_folio_transfer_pair_id(?, ?, ?)", SubFolioID, 1, GuestDepositID).Error; err != nil {
		return err
	}

	return nil
}

func UpdateGuestDepositTransferPairID(Db *gorm.DB, GuestDepositID1, GuestDepositID2 uint64, ValidUserCode string) error {
	if err := Db.Table(DBVar.TableName.GuestDeposit).Where("id=?", GuestDepositID1).Updates(map[string]interface{}{"is_pair_with_folio": false, "transfer_pair_id": GuestDepositID2, "updated_by": ValidUserCode}).Error; err != nil {
		return err
	}

	if err := Db.Table(DBVar.TableName.GuestDeposit).Where("id=?", GuestDepositID2).Updates(map[string]interface{}{"is_pair_with_folio": false, "transfer_pair_id": GuestDepositID1, "updated_by": ValidUserCode}).Error; err != nil {
		return err
	}

	return nil
}

func UpdateRoomStatus(tx *gorm.DB, ValidUserCode, RoomNumber, RoomStatusCode string) error {
	if err := tx.Table(DBVar.TableName.CfgInitRoom).Where("number = ?", RoomNumber).Updates(&map[string]interface{}{
		"status_code": RoomStatusCode,
		"updated_by":  ValidUserCode}).Error; err != nil {
		return err
	}
	return nil
}

func GetFolioStatus(DB *gorm.DB, Number uint64, IsFolio bool) string {
	var FolioStatus struct {
		StatusCode string
	}

	Condition := "reservation_number = ?"
	if IsFolio {
		Condition = "number = ?"
	}
	DB.Table(DBVar.TableName.Folio).Select("status_code").Where(Condition, Number).Scan(&FolioStatus)

	return FolioStatus.StatusCode
}

func GetRoomStatus(DB *gorm.DB, Number string) (StatusCode string, BlockStatusCode string) {
	var RoomStatus struct {
		StatusCode      string
		BlockStatusCode string
	}
	DB.Table(DBVar.TableName.CfgInitRoom).Where("number = ?", Number).Limit(1).Scan(&RoomStatus)

	return RoomStatus.StatusCode, RoomStatus.BlockStatusCode
}

func GetRoomTypeCode(DB *gorm.DB, Number string) string {
	var Room struct {
		RoomTypeCode string
	}
	DB.Table(DBVar.TableName.CfgInitRoom).Where("number = ?", Number).Scan(&Room)

	return Room.RoomTypeCode
}

func GetBedTypeCode(DB *gorm.DB, Number string) string {
	var BedTypeCode string
	DB.Table(DBVar.TableName.CfgInitRoom).Select("bed_type_code").Where("number = ?", Number).Limit(1).Scan(&BedTypeCode)

	return BedTypeCode
}

func GetReservationStatus(DB *gorm.DB, Number uint64) (StatusCode string) {
	DB.Table(DBVar.TableName.Reservation).Select("status_code").Where("number = ? ", Number).Limit(1).Scan(&StatusCode)
	return StatusCode
}

func IsRoomAvailable(DB *gorm.DB, Dataset *global_var.TDataset, RoomTypeCode string, BedTypeCode string, RoomNumber string, ArrivalDate, DepartureDate time.Time, ReservationNumber, FolioNumber, RoomUnavailableID, RoomAllotmentID uint64, ReadyOnly, AllotmentOnly bool) (bool, error) {
	if RoomNumber != "" {
		RoomTypeCode = GetRoomTypeCode(DB, RoomNumber)
		BedTypeCode = GetBedTypeCode(DB, RoomNumber)
	}
	if RoomTypeCode == "" && RoomNumber == "" {
		return false, nil
	}

	RoomList := GetAvailableRoomByType(DB, Dataset, ArrivalDate, DepartureDate, RoomTypeCode, BedTypeCode, ReservationNumber, FolioNumber, RoomUnavailableID, RoomAllotmentID, ReadyOnly, AllotmentOnly)
	IsAvailableInRoomList := false
	RoomAvailableCount, err := GetAvailableRoomCountByType(DB, ArrivalDate, DepartureDate, RoomTypeCode, BedTypeCode, ReservationNumber, FolioNumber, RoomUnavailableID, RoomAllotmentID, ReadyOnly, AllotmentOnly)
	if err != nil {
		return false, err
	}

	if RoomNumber != "" {
		for _, el := range RoomList {
			if el.RoomNumber == RoomNumber {
				IsAvailableInRoomList = true
				break
			}
		}
	} else {
		IsAvailableInRoomList = true
	}

	if FolioNumber == 0 && ReservationNumber == 0 {
		if IsAvailableInRoomList && RoomAvailableCount > 0 {
			return true, nil
		}
	} else {

		if IsAvailableInRoomList {
			return true, nil
		}
	}

	return false, nil
}

func GetJournalRefNumber(c *gin.Context, DB *gorm.DB, PrefixX string, PostingDate time.Time) string {
	Prefix := PrefixX + General.FormatDatePrefix(PostingDate) + "-"
	// Number, err := cache.DataCache.GetString(c, c.GetString("UnitCode"), global_var.CacheKey.LastRefNumber+"_"+PrefixX)
	// if err != nil {
	var DataOutput uint64
	DB.Raw(
		"SELECT" +
			" CAST(RIGHT(ref_number,LENGTH(ref_number)-" + strconv.FormatInt(int64(len(Prefix)), 10) + ") AS UNSIGNED) AS MaxRefNumber " +
			"FROM" +
			" acc_journal" +
			" WHERE LEFT(ref_number," + strconv.FormatInt(int64(len(Prefix)), 10) + ")='" + Prefix + "' " +
			"ORDER BY MaxRefNumber DESC " +
			"LIMIT 1;").Scan(&DataOutput)

	return Prefix + General.Uint64ToStr(DataOutput+1)

	// }
	// return fmt.Sprintf("%s%d", Prefix, General.StrToInt64(Number)+1)
}

func GetJournalRefNumberTemp(c *gin.Context, DB *gorm.DB, PrefixX string, PostingDate time.Time) string {
	Prefix := PrefixX + General.FormatDatePrefix(PostingDate) + "-"
	// Number, err := cache.DataCache.GetString(c, c.GetString("UnitCode"), global_var.CacheKey.LastRefNumber+"_"+PrefixX+"_TEMP")
	// if err != nil {
	var DataOutput uint64
	DB.Raw(
		"SELECT" +
			" CAST(RIGHT(ref_number,LENGTH(ref_number)-" + strconv.FormatInt(int64(len(Prefix)), 10) + ") AS UNSIGNED) AS MaxRefNumber " +
			"FROM" +
			" acc_journal" +
			" WHERE LEFT(ref_number," + strconv.FormatInt(int64(len(Prefix)), 10) + ")='" + Prefix + "' " +
			"ORDER BY MaxRefNumber DESC " +
			"LIMIT 1;").Scan(&DataOutput)

	return Prefix + General.Uint64ToStr(DataOutput+1)

	// }
	// return fmt.Sprintf("%s%d", Prefix, General.StrToInt64(Number)+1)
}

func GetReceiveNumber(c *gin.Context, DB *gorm.DB, PostingDate time.Time) string {
	var Prefix string
	Prefix = GlobalVar.ConstProgramVariable.ReceiveNumberPrefix + General.FormatDatePrefix(PostingDate) + "-"
	// Number, err := cache.DataCache.GetString(c, c.GetString("UnitCode"), global_var.CacheKey.LastReceiveNumber)
	// if err != nil {
	var DataOutput uint64
	DB.Raw("SELECT" +
		" CAST(RIGHT(number, LENGTH(number) - " + strconv.FormatInt(int64(len(Prefix)), 10) + ") AS UNSIGNED) AS MaxReceiveNumber " +
		"FROM" +
		" inv_receiving" +
		" WHERE LEFT(number," + strconv.FormatInt(int64(len(Prefix)), 10) + ")='" + Prefix + "' " +
		"ORDER BY MaxReceiveNumber DESC " +
		"LIMIT 1").Scan(&DataOutput)

	return Prefix + General.Uint64ToStr(DataOutput+1)

	// }
	// return fmt.Sprintf("%s%d", Prefix, General.StrToInt64(Number)+1)
}

func GetPaymentNumber(c *gin.Context, DB *gorm.DB, PostingDate time.Time) string {
	var Prefix string
	Prefix = GlobalVar.ConstProgramVariable.PaymentNumberPrefix + General.FormatDatePrefix(PostingDate) + "-"
	// Number, err := cache.DataCache.GetString(c, c.GetString("UnitCode"), global_var.CacheKey.LastPaymentNumber)
	// if err != nil {
	var DataOutput uint64
	DB.Raw("SELECT" +
		" CAST(RIGHT(ap_ar_number, LENGTH(ap_ar_number) - " + strconv.FormatInt(int64(len(Prefix)), 10) + ") AS UNSIGNED) AS MaxPaymentNumber " +
		"FROM" +
		" acc_journal" +
		" WHERE LEFT(ap_ar_number," + strconv.FormatInt(int64(len(Prefix)), 10) + ")='" + Prefix + "' " +
		"ORDER BY MaxPaymentNumber DESC " +
		"LIMIT 1").Scan(&DataOutput)

	return Prefix + General.Uint64ToStr(DataOutput+1)
	// }
	// return fmt.Sprintf("%s%d", Prefix, General.StrToInt64(Number)+1)
}

func GetJournalAccountTypeCode(DB *gorm.DB, JournalAccountCode string) string {
	var TypeCode string
	DB.Table(DBVar.TableName.CfgInitAccount).Select("type_code").Where("code=?", JournalAccountCode).Limit(1).Scan(&TypeCode)
	return TypeCode
}

func GetAPNumber(c *gin.Context, DB *gorm.DB, PostingDate time.Time) string {
	var Prefix string
	Prefix = GlobalVar.ConstProgramVariable.APNumberPrefix + General.FormatDatePrefix(PostingDate) + "-"

	// Number, err := cache.DataCache.GetString(c, c.GetString("UnitCode"), global_var.CacheKey.LastAPNumber)
	// if err != nil {
	var DataOutput uint64
	DB.Raw("SELECT" +
		" CAST(RIGHT(number, LENGTH(number) - " + strconv.FormatInt(int64(len(Prefix)), 10) + ") AS UNSIGNED) AS MaxAPNumber " +
		"FROM" +
		" acc_ap_ar" +
		" WHERE LEFT(number," + strconv.FormatInt(int64(len(Prefix)), 10) + ")='" + Prefix + "' " +
		"ORDER BY MaxAPNumber DESC " +
		"LIMIT 1").Scan(&DataOutput)

	return Prefix + General.Uint64ToStr(DataOutput+1)
	// }
	// return fmt.Sprintf("%s%d", Prefix, General.StrToInt64(Number)+1)

}

func GetARNumber(c *gin.Context, DB *gorm.DB, PostingDate time.Time) string {
	var Prefix string
	Prefix = GlobalVar.ConstProgramVariable.ARNumberPrefix + General.FormatDatePrefix(PostingDate) + "-"
	// Number, err := cache.DataCache.GetString(c, c.GetString("UnitCode"), global_var.CacheKey.LastARNumber)
	// if err != nil {
	var DataOutput uint64
	DB.Raw("SELECT" +
		" CAST(RIGHT(number, LENGTH(number) - " + strconv.FormatInt(int64(len(Prefix)), 10) + ") AS UNSIGNED) AS MaxARNumber " +
		"FROM" +
		" acc_ap_ar" +
		" WHERE LEFT(number," + strconv.FormatInt(int64(len(Prefix)), 10) + ")='" + Prefix + "' " +
		"ORDER BY MaxARNumber DESC " +
		"LIMIT 1").Scan(&DataOutput)

	return Prefix + General.Uint64ToStr(DataOutput+1)
	// }

	// return fmt.Sprintf("%s%d", Prefix, General.StrToInt64(Number)+1)
}

func GetCostingNumber(c *gin.Context, DB *gorm.DB, PostingDate time.Time) string {
	Prefix := GlobalVar.ConstProgramVariable.CostingNumberPrefix + General.FormatDatePrefix(PostingDate) + "-"
	// Number, err := cache.DataCache.GetString(c, c.GetString("UnitCode"), global_var.CacheKey.LastCostingNumber)
	// if err != nil {
	var DataOutput uint64
	DB.Raw("SELECT" +
		" CAST(RIGHT(number, LENGTH(number) - " + strconv.FormatInt(int64(len(Prefix)), 10) + ") AS UNSIGNED) AS MaxCostingNumber " +
		"FROM" +
		" inv_costing" +
		" WHERE LEFT(number," + strconv.FormatInt(int64(len(Prefix)), 10) + ")='" + Prefix + "' " +
		"ORDER BY MaxCostingNumber DESC " +
		"LIMIT 1").Scan(&DataOutput)

	return Prefix + General.Uint64ToStr(DataOutput+1)
	// }
	// return fmt.Sprintf("%s%d", Prefix, General.StrToInt64(Number)+1)
}

func GetDepreciationNumber(c *gin.Context, DB *gorm.DB, PostingDate time.Time) string {
	Prefix := GlobalVar.ConstProgramVariable.DepreciationNumberPrefix + General.FormatDatePrefix(PostingDate) + "-"
	// Number, err := cache.DataCache.GetString(c, c.GetString("UnitCode"), global_var.CacheKey.LastDepreciationNumber)
	// if err != nil {
	var DataOutput uint64
	DB.Raw("SELECT" +
		" CAST(RIGHT(number, LENGTH(number) - " + strconv.FormatInt(int64(len(Prefix)), 10) + ") AS UNSIGNED) AS MaxDepreciationNumber " +
		"FROM" +
		" fa_depreciation" +
		" WHERE LEFT(number," + strconv.FormatInt(int64(len(Prefix)), 10) + ")='" + Prefix + "' " +
		"ORDER BY MaxDepreciationNumber DESC " +
		"LIMIT 1").Scan(&DataOutput)

	return Prefix + General.Uint64ToStr(DataOutput+1)
	// }
	// return fmt.Sprintf("%s%d", Prefix, General.StrToInt64(Number)+1)
}

func GetSRNumber(c *gin.Context, DB *gorm.DB, PostingDate time.Time) string {
	Prefix := GlobalVar.ConstProgramVariable.SRNumberPrefix + General.FormatDatePrefix(PostingDate) + "-"
	// Number, err := cache.DataCache.GetString(c, c.GetString("UnitCode"), global_var.CacheKey.LastSRNumber)
	// if err != nil {
	var DataOutput uint64
	DB.Raw("SELECT" +
		" CAST(RIGHT(number, LENGTH(number) - " + strconv.FormatInt(int64(len(Prefix)), 10) + ") AS UNSIGNED) AS MaxDepreciationNumber " +
		"FROM" +
		" inv_store_requisition" +
		" WHERE LEFT(number," + strconv.FormatInt(int64(len(Prefix)), 10) + ")='" + Prefix + "' " +
		"ORDER BY MaxDepreciationNumber DESC " +
		"LIMIT 1").Scan(&DataOutput)

	return Prefix + General.Uint64ToStr(DataOutput+1)
	// }
	// return fmt.Sprintf("%s%d", Prefix, General.StrToInt64(Number)+1)
}
func GetProductionNumber(c *gin.Context, DB *gorm.DB, PostingDate time.Time) string {
	Prefix := GlobalVar.ConstProgramVariable.ProductionNumberPrefix + General.FormatDatePrefix(PostingDate) + "-"
	// Number, err := cache.DataCache.GetString(c, c.GetString("UnitCode"), global_var.CacheKey.LastProductionNumber)
	// if err != nil {
	var DataOutput uint64
	DB.Raw("SELECT" +
		" CAST(RIGHT(number, LENGTH(number) - " + strconv.FormatInt(int64(len(Prefix)), 10) + ") AS UNSIGNED) AS MaxDepreciationNumber " +
		"FROM" +
		" inv_production" +
		" WHERE LEFT(number," + strconv.FormatInt(int64(len(Prefix)), 10) + ")='" + Prefix + "' " +
		"ORDER BY MaxDepreciationNumber DESC " +
		"LIMIT 1").Scan(&DataOutput)

	return Prefix + General.Uint64ToStr(DataOutput+1)
	// }
	// return fmt.Sprintf("%s%d", Prefix, General.StrToInt64(Number)+1)
}

func GetReturnStockNumber(c *gin.Context, DB *gorm.DB, PostingDate time.Time) string {
	Prefix := GlobalVar.ConstProgramVariable.ReturnStockPrefix + General.FormatDatePrefix(PostingDate) + "-"
	// Number, err := cache.DataCache.GetString(c, c.GetString("UnitCode"), global_var.CacheKey.LastReturnStock)
	// if err != nil {
	var DataOutput uint64
	DB.Raw("SELECT" +
		" CAST(RIGHT(number, LENGTH(number) - " + strconv.FormatInt(int64(len(Prefix)), 10) + ") AS UNSIGNED) AS MaxDepreciationNumber " +
		"FROM" +
		" inv_return_stock" +
		" WHERE LEFT(number," + strconv.FormatInt(int64(len(Prefix)), 10) + ")='" + Prefix + "' " +
		"ORDER BY MaxDepreciationNumber DESC " +
		"LIMIT 1").Scan(&DataOutput)

	return Prefix + General.Uint64ToStr(DataOutput+1)
	// }
	// return fmt.Sprintf("%s%d", Prefix, General.StrToInt64(Number)+1)
}

func GetOpnameNumber(c *gin.Context, DB *gorm.DB, PostingDate time.Time) string {
	Prefix := GlobalVar.ConstProgramVariable.OpnameNumberPrefix + General.FormatDatePrefix(PostingDate) + "-"
	// Number, err := cache.DataCache.GetString(c, c.GetString("UnitCode"), global_var.CacheKey.LastOpnameNumber)
	// if err != nil {
	var DataOutput uint64
	DB.Raw("SELECT" +
		" CAST(RIGHT(number, LENGTH(number) - " + strconv.FormatInt(int64(len(Prefix)), 10) + ") AS UNSIGNED) AS MaxDepreciationNumber " +
		"FROM" +
		" inv_opname" +
		" WHERE LEFT(number," + strconv.FormatInt(int64(len(Prefix)), 10) + ")='" + Prefix + "' " +
		"ORDER BY MaxDepreciationNumber DESC " +
		"LIMIT 1").Scan(&DataOutput)

	return Prefix + General.Uint64ToStr(DataOutput+1)
	// }
	// return fmt.Sprintf("%s%d", Prefix, General.StrToInt64(Number)+1)
}

func GetFAPONumber(c *gin.Context, DB *gorm.DB, PostingDate time.Time) string {
	Prefix := GlobalVar.ConstProgramVariable.FAPONumberPrefix + General.FormatDatePrefix(PostingDate) + "-"
	// Number, err := cache.DataCache.GetString(c, c.GetString("UnitCode"), global_var.CacheKey.LastFAPONumber)
	// if err != nil {
	var DataOutput uint64
	DB.Raw("SELECT" +
		" CAST(RIGHT(number, LENGTH(number) - " + strconv.FormatInt(int64(len(Prefix)), 10) + ") AS UNSIGNED) AS MaxDepreciationNumber " +
		"FROM" +
		" fa_purchase_order" +
		" WHERE LEFT(number," + strconv.FormatInt(int64(len(Prefix)), 10) + ")='" + Prefix + "' " +
		"ORDER BY MaxDepreciationNumber DESC " +
		"LIMIT 1").Scan(&DataOutput)

	return Prefix + General.Uint64ToStr(DataOutput+1)
	// }
	// return fmt.Sprintf("%s%d", Prefix, General.StrToInt64(Number)+1)
}

func GetFAReceiveNumber(c *gin.Context, DB *gorm.DB, PostingDate time.Time) string {
	Prefix := GlobalVar.ConstProgramVariable.FAReceiveNumberPrefix + General.FormatDatePrefix(PostingDate) + "-"
	// Number, err := cache.DataCache.GetString(c, c.GetString("UnitCode"), global_var.CacheKey.LastFAReceiveNumber)
	// if err != nil {
	var DataOutput uint64
	DB.Raw("SELECT" +
		" CAST(RIGHT(number, LENGTH(number) - " + strconv.FormatInt(int64(len(Prefix)), 10) + ") AS UNSIGNED) AS MaxDepreciationNumber " +
		"FROM" +
		" fa_receive" +
		" WHERE LEFT(number," + strconv.FormatInt(int64(len(Prefix)), 10) + ")='" + Prefix + "' " +
		"ORDER BY MaxDepreciationNumber DESC " +
		"LIMIT 1").Scan(&DataOutput)

	return Prefix + General.Uint64ToStr(DataOutput+1)
	// }
	// return fmt.Sprintf("%s%d", Prefix, General.StrToInt64(Number)+1)
}

func GetFACode(c *gin.Context, DB *gorm.DB, ItemCode string, PostingDate time.Time) (SortNumber uint64, Code string) {
	Prefix := General.FormatDatePrefix(PostingDate) + ItemCode
	Result := Prefix + "1"

	var DataOutput uint64
	DB.Raw("SELECT" +
		" CAST(RIGHT(code, LENGTH(code) - " + strconv.FormatInt(int64(len(Prefix)), 10) + ") AS UNSIGNED) AS MaxDepreciationNumber " +
		"FROM" +
		" fa_list" +
		" WHERE LEFT(code," + strconv.FormatInt(int64(len(Prefix)), 10) + ")='" + Prefix + "' " +
		"ORDER BY MaxDepreciationNumber DESC " +
		"LIMIT 1").Scan(&DataOutput)
	SortNumber = DataOutput + 1
	Result = fmt.Sprintf("%s%d", Prefix, SortNumber)
	return SortNumber, Result
}

func GetFACondition(DB *gorm.DB, FACode string) string {
	ConditionCode := ""
	DB.Table(DBVar.TableName.FaList).Select("condition_code").Where("code=?", FACode).Limit(1).Scan(&ConditionCode)
	return ConditionCode
}

func GetFAJournalAccount(DB *gorm.DB, ItemCode string) string {
	var JournalAccountCode string
	DB.Table(DBVar.TableName.FaCfgInitItem).Select("fa_cfg_init_item_category.journal_account_code").
		Joins("LEFT OUTER JOIN fa_cfg_init_item_category ON (fa_cfg_init_item.category_code = fa_cfg_init_item_category.code)").
		Where("fa_cfg_init_item.code=?", ItemCode).
		Scan(&JournalAccountCode)
	return JournalAccountCode
}

func GetFAJournalAccountDepreciation(DB *gorm.DB, ItemCode string) string {
	var JournalAccountCodeDepreciation string
	DB.Table(DBVar.TableName.FaCfgInitItem).Select("fa_cfg_init_item_category.journal_account_code_depreciation").
		Joins("LEFT OUTER JOIN fa_cfg_init_item_category ON (fa_cfg_init_item.category_code = fa_cfg_init_item_category.code)").
		Where("fa_cfg_init_item.code=?", ItemCode).
		Scan(&JournalAccountCodeDepreciation)
	return JournalAccountCodeDepreciation
}

func GetFACOGSAccount(DB *gorm.DB, ItemCode string) string {
	var JournalAccountCodeCogs string
	DB.Table(DBVar.TableName.FaCfgInitItem).Select("fa_cfg_init_item_category.journal_account_code_cogs").
		Joins("LEFT OUTER JOIN fa_cfg_init_item_category ON (fa_cfg_init_item.category_code = fa_cfg_init_item_category.code)").
		Where("fa_cfg_init_item.code=?", ItemCode).
		Limit(1).
		Scan(&JournalAccountCodeCogs)
	return JournalAccountCodeCogs
}

func GetFAExpenseAccount(DB *gorm.DB, ItemCode string) string {
	var JournalAccountCodeExpense string
	DB.Table(DBVar.TableName.FaCfgInitItem).Select("fa_cfg_init_item_category.journal_account_code_expense").
		Joins("LEFT OUTER JOIN fa_cfg_init_item_category ON (fa_cfg_init_item.category_code = fa_cfg_init_item_category.code)").
		Where("fa_cfg_init_item.code=?", ItemCode).
		Limit(1).
		Scan(&JournalAccountCodeExpense)
	return JournalAccountCodeExpense
}

func GetFASellAccount(DB *gorm.DB, ItemCode string) string {
	var JournalAccountCodeSell string
	DB.Table(DBVar.TableName.FaCfgInitItem).Select("fa_cfg_init_item_category.journal_account_code_sell").
		Joins("LEFT OUTER JOIN fa_cfg_init_item_category ON (fa_cfg_init_item.category_code = fa_cfg_init_item_category.code)").
		Where("fa_cfg_init_item.code=?", ItemCode).
		Limit(1).
		Scan(&JournalAccountCodeSell)
	return JournalAccountCodeSell
}

func GetFASpoilJournalAccount(DB *gorm.DB, ItemCode string) string {
	var JournalAccountCodeSpoil string
	DB.Table(DBVar.TableName.FaCfgInitItem).Select("fa_cfg_init_item_category.journal_account_code_spoil").
		Joins("LEFT OUTER JOIN fa_cfg_init_item_category ON (fa_cfg_init_item.category_code = fa_cfg_init_item_category.code)").
		Where("fa_cfg_init_item.code=?", ItemCode).
		Limit(1).
		Scan(&JournalAccountCodeSpoil)
	return JournalAccountCodeSpoil
}

func GetStockTransferNumber(c *gin.Context, DB *gorm.DB, PostingDate time.Time) string {
	var Prefix string
	Prefix = GlobalVar.ConstProgramVariable.StockTransferNumberPrefix + General.FormatDatePrefix(PostingDate) + "-"
	// Number, err := cache.DataCache.GetString(c, c.GetString("UnitCode"), global_var.CacheKey.LastStockTransferNumber)
	// if err != nil {
	var DataOutput uint64
	DB.Raw("SELECT" +
		" CAST(RIGHT(number, LENGTH(number) - " + strconv.FormatInt(int64(len(Prefix)), 10) + ") AS UNSIGNED) AS MaxStockTransferNumber " +
		"FROM" +
		" inv_stock_transfer" +
		" WHERE LEFT(number," + strconv.FormatInt(int64(len(Prefix)), 10) + ")='" + Prefix + "' " +
		"ORDER BY MaxStockTransferNumber DESC " +
		"LIMIT 1").Scan(&DataOutput)

	return Prefix + General.Uint64ToStr(DataOutput+1)

	// }
	// return fmt.Sprintf("%s%d", Prefix, General.StrToInt64(Number)+1)
}

func IsYearClosed(DB *gorm.DB, Year int) bool {
	var DataOutput []map[string]interface{}
	DB.Raw(
		"SELECT id FROM acc_close_year"+
			" WHERE year=? "+
			"LIMIT 1", Year).Scan(&DataOutput)

	return (len(DataOutput) > 0)
}

func IsJournalClosed(DB *gorm.DB, Mode byte, PostingDateB4, PostingDate time.Time) bool {
	Month := PostingDate.Month()
	Year := PostingDate.Year()

	if Mode != 1 {
		IsMonthClosed, _ := IsMonthClosed(DB, int(Month), Year)
		return IsMonthClosed || IsYearClosed(DB, Year)
	} else {
		MonthB4 := PostingDateB4.Month()
		YearB4 := PostingDateB4.Year()

		IsMonthClosed1, _ := IsMonthClosed(DB, int(Month), Year)
		IsMonthClosed2, _ := IsMonthClosed(DB, int(MonthB4), YearB4)

		return IsMonthClosed1 || IsYearClosed(DB, Year) || IsMonthClosed2 || IsYearClosed(DB, YearB4)
	}
}

func GetAPRefundDepositOutStanding(DB *gorm.DB, Dataset *global_var.TDataset, SubFolioId uint64, RefNumber string) float64 {
	var DataOutput float64
	DB.Raw(
		"SELECT Amount FROM ("+
			"SELECT"+
			" sub_folio.id,"+
			" (SUM(IF(sub_folio.type_code=?, sub_folio.quantity*sub_folio.amount, -(sub_folio.quantity*sub_folio.amount))) - IFNULL(Payment.TotalPaid,0)) AS Amount "+
			"FROM"+
			" sub_folio"+
			" LEFT OUTER JOIN ("+
			"SELECT"+
			" sub_folio_id,"+
			" SUM(amount) AS TotalPaid "+
			"FROM"+
			" acc_ap_refund_deposit_payment_detail"+
			" WHERE sub_folio_id=?"+
			" AND ref_number<>? "+
			"GROUP BY sub_folio_id) AS Payment ON (sub_folio.id = Payment.sub_folio_id) "+
			" WHERE sub_folio.account_code=?"+
			" AND sub_folio.void='0' "+
			"GROUP BY sub_folio.correction_breakdown) AS APNoShow"+
			" WHERE id=?", GlobalVar.TransactionType.Debit, SubFolioId, RefNumber, Dataset.GlobalAccount.APRefundDeposit, SubFolioId).Scan(&DataOutput)

	return DataOutput
}

func GetARCityLedgerInvoiceOutStanding(DB *gorm.DB, InvoiceNumber string) (float64, error) {
	var Amount float64
	if err := DB.Table(DBVar.TableName.InvoiceItem).Select(
		"SUM(amount_charged - amount_paid) AS Amount",
	).Where("invoice_number=?", InvoiceNumber).
		Scan(&Amount).Error; err != nil {
		return 0, err
	}
	return Amount, nil
}
func IsFolioClosed(DB *gorm.DB, FolioNumber uint64) bool {
	var StatusCode string
	DB.Table(DBVar.TableName.Folio).Select("status_code").Where("number = ? AND status_code=?", FolioNumber, GlobalVar.FolioStatus.Open).Find(StatusCode)
	return StatusCode != GlobalVar.FolioStatus.Open
}

// func GetGlobalAccount(Account string) string {
// 	return MasterData.GetConfiguration(GlobalVar.SystemCode.Hotel, GlobalVar.ConfigurationCategory.GlobalAccount, Account, false).(string)
// }

// func GetGlobalSubDepartment(SubDepartment string) string {
// 	return MasterData.GetConfiguration(GlobalVar.SystemCode.Hotel, GlobalVar.ConfigurationCategory.GlobalSubDepartment, SubDepartment, false).(string)
// }

// func GetGlobalDepartment(Department string) string {
// 	return MasterData.GetConfiguration(GlobalVar.SystemCode.Hotel, GlobalVar.ConfigurationCategory.GlobalDepartment, Department, false).(string)
// }
// func GetGlobalJournalAccount(AccountName string) string {
// 	return MasterData.GetConfiguration(GlobalVar.SystemCode.Hotel, GlobalVar.ConfigurationCategory.GlobalJournalAccount, AccountName, false).(string)
// }

func IsFolioHaveBreakfast(DB *gorm.DB, Dataset *global_var.TDataset, FolioNumber uint64) bool {
	Number := 0
	DB.Table(DBVar.TableName.Folio).Select("folio.number").
		Joins("LEFT OUTER JOIN guest_detail ON (folio.guest_detail_id = guest_detail.id)").
		Joins("LEFT OUTER JOIN cfg_init_room_rate ON (guest_detail.room_rate_code = cfg_init_room_rate.code)").
		Joins("LEFT OUTER JOIN cfg_init_room_rate_breakdown ON (guest_detail.room_rate_code = cfg_init_room_rate_breakdown.room_rate_code)").
		Where("folio.number = ? AND (IFNULL(cfg_init_room_rate_breakdown.account_code, '') = ? OR cfg_init_room_rate.include_breakfast='1')", FolioNumber, Dataset.GlobalAccount.Breakfast).
		Find(&Number)

	return Number > 0
}

func IsReservationHaveBreakfast(DB *gorm.DB, Dataset *global_var.TDataset, ReservationNumber uint64) bool {
	Number := 0
	DB.Table(DBVar.TableName.Reservation).Select("reservation.number").
		Joins("LEFT OUTER JOIN guest_detail ON (reservation.guest_detail_id = guest_detail.id)").
		Joins("LEFT OUTER JOIN cfg_init_room_rate ON (guest_detail.room_rate_code = cfg_init_room_rate.code)").
		Joins("LEFT OUTER JOIN cfg_init_room_rate_breakdown ON (guest_detail.room_rate_code = cfg_init_room_rate_breakdown.room_rate_code)").
		Where("reservation.number = ? AND (IFNULL(cfg_init_room_rate_breakdown.account_code, '') = ? OR cfg_init_room_rate.include_breakfast='1')", ReservationNumber, Dataset.GlobalAccount.Breakfast).
		Find(&Number)

	return Number > 0
}

func GetRoomRateAmountXXXX(DB *gorm.DB, Dataset *global_var.TDataset, RoomRateCode string, PostingDate time.Time, Adult, Child uint, IsWeekendX bool) float64 {
	var Rate float64
	var Pax uint

	var DataOutput DBVar.Cfg_init_room_rate
	DB.Table(DBVar.TableName.CfgInitRoomRate).Select(" weekday_rate1,"+
		" weekday_rate2,"+
		" weekday_rate3,"+
		" weekday_rate4,"+
		" weekend_rate1,"+
		" weekend_rate2,"+
		" weekend_rate3,"+
		" weekend_rate4,"+
		" weekday_rate_child1,"+
		" weekday_rate_child2,"+
		" weekday_rate_child3,"+
		" weekday_rate_child4,"+
		" weekend_rate_child1,"+
		" weekend_rate_child2,"+
		" weekend_rate_child3,"+
		" weekend_rate_child4,"+
		" include_child,"+
		" extra_pax,"+
		" per_pax ").
		Where("code = ? ", RoomRateCode).Take(&DataOutput)

	Pax = Adult
	if DataOutput.IncludeChild > 0 {
		Pax = Adult + Child
	}

	if !PostingDate.IsZero() {
		IsWeekendX = General.IsWeekend(PostingDate, Dataset)
	}

	if IsWeekendX {
		if General.StrToBool(Dataset.Configuration[GlobalVar.ConfigurationCategory.General][GlobalVar.ConfigurationName.UseChildRate].(string)) {
			if Adult == 1 {
				Rate = DataOutput.WeekendRate1
			} else if Adult == 2 {
				Rate = DataOutput.WeekendRate2
			} else if Adult == 3 {
				Rate = DataOutput.WeekendRate3
			} else if Adult >= 4 {
				Rate = DataOutput.WeekendRate4
			}

			if Child == 1 {
				Rate = Rate + DataOutput.WeekendRateChild1
			} else if Child == 2 {
				Rate = Rate + DataOutput.WeekendRateChild2
			} else if Child == 3 {
				Rate = Rate + DataOutput.WeekendRateChild3
			} else if Child == 4 {
				Rate = Rate + DataOutput.WeekendRateChild4
			}
		} else {
			if Pax == 1 {
				Rate = DataOutput.WeekendRate1
			} else if Pax == 2 {
				Rate = DataOutput.WeekendRate2
			} else if Pax == 3 {
				Rate = DataOutput.WeekendRate3
			} else if Pax >= 4 {
				Rate = DataOutput.WeekendRate4
			}
		}
	} else {
		if General.StrToBool(Dataset.Configuration[GlobalVar.ConfigurationCategory.General][GlobalVar.ConfigurationName.UseChildRate].(string)) {
			if Adult == 1 {
				Rate = DataOutput.WeekdayRate1
			} else if Adult == 2 {
				Rate = DataOutput.WeekdayRate2
			} else if Adult == 3 {
				Rate = DataOutput.WeekdayRate3
			} else if Adult >= 4 {
				Rate = DataOutput.WeekdayRate4
			}

			if Child == 1 {
				Rate = Rate + DataOutput.WeekdayRateChild1
			} else if Child == 2 {
				Rate = Rate + DataOutput.WeekdayRateChild2
			} else if Child == 3 {
				Rate = Rate + DataOutput.WeekdayRateChild3
			} else if Child == 4 {
				Rate = Rate + DataOutput.WeekdayRateChild4
			}

		} else {
			if Pax == 1 {
				Rate = DataOutput.WeekdayRate1
			} else if Pax == 2 {
				Rate = DataOutput.WeekdayRate2
			} else if Pax == 3 {
				Rate = DataOutput.WeekdayRate3
			} else if Pax >= 4 {
				Rate = DataOutput.WeekdayRate4
			}
		}
	}

	if Pax > 4 {
		if DataOutput.PerPax > 0 {
			Rate = Rate + (float64(Pax-4) * DataOutput.ExtraPax)
		} else {
			Rate = Rate + DataOutput.ExtraPax
		}
	}

	return Rate
}

func IsCanPostCharge(c *gin.Context, DB *gorm.DB, ChargeFrequencyCode string, ArrivalDate time.Time) bool {
	var WeeklyCharged, MonthlyCharged bool
	AuditDate := GetAuditDate(c, DB, false)
	ArrivalDate = General.DateOf(ArrivalDate)
	WeeklyCharged = General.DaysBetween(ArrivalDate, AuditDate)%7 == 0
	MonthlyCharged = General.DaysBetween(ArrivalDate, AuditDate)%30 == 0

	return ((ChargeFrequencyCode == GlobalVar.ChargeFrequency.OnceOnly && ArrivalDate == AuditDate) ||
		(ChargeFrequencyCode == GlobalVar.ChargeFrequency.Daily) ||
		(ChargeFrequencyCode == GlobalVar.ChargeFrequency.Weekly && WeeklyCharged) ||
		(ChargeFrequencyCode == GlobalVar.ChargeFrequency.Monthly && MonthlyCharged))

}

func GetTotalBreakdownAmount(Quantity, BreakdownAmount, BreakdownAmountExtra float64, PerPax, IncludeChild, PerPaxExtra bool, MaxPax, Adult, Child int) float64 {
	var TotalPax int
	var Amount, AmountExtra float64

	if IncludeChild {
		TotalPax = Adult + Child
	} else {
		TotalPax = Adult
	}

	if PerPax {
		if TotalPax <= MaxPax {
			Amount = float64(TotalPax) * BreakdownAmount
		} else {
			Amount = float64(MaxPax) * BreakdownAmount
		}
	} else {
		Amount = float64(Quantity) * BreakdownAmount
	}

	if PerPaxExtra {
		if TotalPax > MaxPax {
			AmountExtra = (float64(TotalPax) - float64(MaxPax)*BreakdownAmountExtra)
		} else {
			AmountExtra = 0
		}
	} else {
		if TotalPax > MaxPax {
			AmountExtra = BreakdownAmount
		}
	}

	return Amount + AmountExtra
}

func GetCommission(c *gin.Context, DB *gorm.DB, CommissionTypeCode string, CommissionValue, RoomRateAmount, RoomRateBasicAmount float64, ArrivalDate time.Time) float64 {
	AuditDate := GetAuditDate(c, DB, false)
	Value := 0.00
	if CommissionTypeCode == GlobalVar.CommissionType.PercentFirstNightFullRate {
		if General.DateOf(ArrivalDate).Equal(General.DateOf(AuditDate)) {
			Value = General.RoundToX3(CommissionValue * RoomRateAmount / 100)
			return Value
		}
	}

	if CommissionTypeCode == GlobalVar.CommissionType.PercentPerNightFullRate {
		Value = General.RoundToX3(CommissionValue * RoomRateAmount / 100)
		return Value
	}

	if CommissionTypeCode == GlobalVar.CommissionType.PercentFirstNightNettRate {
		if General.DateOf(ArrivalDate).Equal(General.DateOf(AuditDate)) {
			Value = General.RoundToX3(CommissionValue * RoomRateBasicAmount / 100)
			return Value
		}
	}

	if CommissionTypeCode == GlobalVar.CommissionType.PercentPerNightNettRate {
		Value = General.RoundToX3(CommissionValue * RoomRateBasicAmount / 100)
		return Value
	}

	if CommissionTypeCode == GlobalVar.CommissionType.FixAmountPerNight {
		Value = General.RoundToX3(CommissionValue)
		return Value
	}

	return Value
}

func GetCommissionPackage(c *gin.Context, DB *gorm.DB, CommissionTypeCode string, CommissionValue, PackageAmount float64, ArrivalDate time.Time) (value float64) {
	AuditDate := GetAuditDate(c, DB, false)
	if CommissionTypeCode == GlobalVar.CommissionType.PercentFirstNightFullRate {
		if General.DateOf(ArrivalDate).Equal(General.DateOf(AuditDate)) {
			value = General.RoundToX3(CommissionValue * PackageAmount / 100)
		}
	} else if CommissionTypeCode == GlobalVar.CommissionType.PercentPerNightFullRate {
		value = General.RoundToX3(CommissionValue * PackageAmount / 100)
	} else if CommissionTypeCode == GlobalVar.CommissionType.FixAmountFirstNight {
		if General.DateOf(ArrivalDate).Equal(General.DateOf(AuditDate)) {
			value = CommissionValue
		}
	} else if CommissionTypeCode == GlobalVar.CommissionType.FixAmountPerNight {
		value = CommissionValue
	}

	if CommissionTypeCode == GlobalVar.CommissionType.FixAmountPerNight {
		value = General.RoundToX3(CommissionValue)
		return value
	}

	return value
}

func InsertSubFolio2(c *gin.Context, DB *gorm.DB, Dataset *global_var.TDataset, FolioNumber uint64, GroupCode, RoomNumber, SubDepartmentCode, AccountCode,
	AccountCodeTransfer, ProductCode, PackageCode, CurrencyCode, Remark, DocumentNumber, VoucherNumber,
	TypeCode, CardBankCode, CardTypeCode string, CorrectionBreakdown, BreakDown1 uint64, PostingType string,
	ExtraChargeID uint64, Basic, Tax, Service, ExchangeRate float64, AllowZeroAmount, IsCanAutoTransferred bool, UserID string) (Result string, SubFolioID uint64, err error) {
	var BasicForeign, TaxForeign, ServiceForeign float64
	var BelongsTo uint64
	var CurrencyCodeDefault string
	var BreakDown2 int

	if CurrencyCode == "" {
		CurrencyCode = GetDefaultCurrencyCode(DB)
		ExchangeRate = GetExchangeRateCurrency(DB, CurrencyCode)
	}

	//Insert Sub folio yg langsung ditentukan tax dan servicenya
	if ((Basic > 0) || ((Basic == 0) && AllowZeroAmount)) && (CurrencyCode != "") && (ExchangeRate > 0) {

		CurrencyCodeDefault = GetDefaultCurrencyCode(DB)
		BasicForeign = Basic
		TaxForeign = Tax
		ServiceForeign = Service
		if CurrencyCode != CurrencyCodeDefault {

			Basic = General.RoundToX3(BasicForeign * ExchangeRate)
			Tax = General.RoundToX3(TaxForeign * ExchangeRate)
			Service = General.RoundToX3(ServiceForeign * ExchangeRate)
		}

		BelongsTo = FolioNumber
		//Cek jika di routing otomatis
		if IsCanAutoTransferred {
			F, G := GetAutoRouting(DB, BelongsTo, AccountCodeTransfer)
			if F > 0 && G != "" {
				FolioNumber = F
				GroupCode = G
			}
		}

		if CorrectionBreakdown == 0 {
			CorrectionBreakdown = GetSubFolioCorrectionBreakdown(c, DB)
		}
		if BreakDown1 == 0 {
			BreakDown1 = GetSubFolioBreakdown1(c, DB)
		}
		BreakDown2 = GetSubFolioBreakdown2(c, DB, BreakDown1)

		var SubFolioData DBVar.Sub_folio
		SubFolioData.FolioNumber = FolioNumber
		SubFolioData.BelongsTo = BelongsTo
		SubFolioData.GroupCode = GroupCode
		SubFolioData.RoomNumber = RoomNumber
		SubFolioData.SubDepartmentCode = SubDepartmentCode
		SubFolioData.AccountCode = AccountCode
		SubFolioData.ProductCode = ProductCode
		SubFolioData.DefaultCurrencyCode = CurrencyCodeDefault
		SubFolioData.CurrencyCode = CurrencyCode
		SubFolioData.Remark = Remark
		SubFolioData.DocumentNumber = DocumentNumber
		SubFolioData.VoucherNumber = VoucherNumber
		SubFolioData.TypeCode = TypeCode
		SubFolioData.CardBankCode = CardBankCode
		SubFolioData.CardTypeCode = CardTypeCode
		SubFolioData.CorrectionBreakdown = CorrectionBreakdown
		SubFolioData.Breakdown1 = BreakDown1
		SubFolioData.Breakdown2 = BreakDown2
		SubFolioData.PostingType = PostingType
		SubFolioData.ExtraChargeId = ExtraChargeID
		SubFolioData.Amount = Basic
		SubFolioData.AmountForeign = BasicForeign
		SubFolioData.ExchangeRate = ExchangeRate
		SubFolioData.Quantity = 1
		SubFolioData.CreatedBy = UserID

		SubFolioID, err = InsertSubFolioX(c, Dataset, 0, SubFolioData, DB)
		if err != nil {
			return "", 0, err
		}

		if Tax > 0 {
			SubFolioData.AccountCode = Dataset.GlobalAccount.Tax
			SubFolioData.ProductCode = SubFolioData.AccountCode
			SubFolioData.Amount = Tax
			SubFolioData.AmountForeign = TaxForeign
			_, err = InsertSubFolioX(c, Dataset, 0, SubFolioData, DB)
			if err != nil {
				return "", 0, err
			}
		}
		if Service > 0 {
			SubFolioData.AccountCode = Dataset.GlobalAccount.Service
			SubFolioData.ProductCode = SubFolioData.AccountCode
			SubFolioData.Amount = Service
			SubFolioData.AmountForeign = ServiceForeign
			_, err = InsertSubFolioX(c, Dataset, 0, SubFolioData, DB)
			if err != nil {
				return "", 0, err
			}
		}

		if FolioNumber == BelongsTo {
			Result = ""
		} else {
			Result = strconv.FormatUint(FolioNumber, 10) + "/" + GroupCode
		}

	}
	return Result, SubFolioID, err
}

func InsertGuestInHouseBreakdown(DB *gorm.DB, PostingDate time.Time, FolioNumber uint64, OutletCode, ProductCode, SubDepartment, AccountCode, Remark, TaxNServiceCode,
	ChargeFrequencyCode string, Quantity, Amount, ExtraPax float64, MaxPax int, PerPax, IncludeChild, PerPaxExtra uint8, UserID string) error {
	var GuestInHouseBreakdown DBVar.Guest_in_house_breakdown
	GuestInHouseBreakdown.AuditDate = PostingDate
	GuestInHouseBreakdown.FolioNumber = FolioNumber
	GuestInHouseBreakdown.OutletCode = OutletCode
	GuestInHouseBreakdown.ProductCode = ProductCode
	GuestInHouseBreakdown.SubDepartmentCode = SubDepartment
	GuestInHouseBreakdown.AccountCode = AccountCode
	GuestInHouseBreakdown.Remark = Remark
	GuestInHouseBreakdown.TaxAndServiceCode = TaxNServiceCode
	GuestInHouseBreakdown.ChargeFrequencyCode = ChargeFrequencyCode
	GuestInHouseBreakdown.Quantity = Quantity
	GuestInHouseBreakdown.Amount = Amount
	GuestInHouseBreakdown.ExtraPax = ExtraPax
	GuestInHouseBreakdown.MaxPax = MaxPax
	GuestInHouseBreakdown.PerPax = PerPax
	GuestInHouseBreakdown.IncludeChild = IncludeChild
	GuestInHouseBreakdown.PerPaxExtra = PerPaxExtra
	GuestInHouseBreakdown.CreatedBy = UserID

	err := DB.Table(DBVar.TableName.GuestInHouseBreakdown).Create(&GuestInHouseBreakdown).Error
	if err != nil {
		return err
	}
	return nil
}

func IsInHousePosted(DB *gorm.DB, AuditDate time.Time, FolioNumber uint64) bool {
	Id := 0
	DB.Table(DBVar.TableName.GuestInHouse).Select("id").Where("audit_date = ? AND folio_number=?", AuditDate, strconv.FormatUint(FolioNumber, 10)).Limit(1).Scan(&Id)

	return Id > 0
}

// func InsertGuestInHouse(DB *gorm.DB, AuditDate time.Time, FolioNumber uint64, GroupCode, RoomTypeCode, BedTypeCode,
// 	RoomNumber, RoomRateCode, BusinessSourceCode, CommissionTypeCode, PaymentTypeCode, MarketCode, TitleCode,
// 	FullName, Street, City, CityCode, CountryCode, StateCode, PostalCode, Phone1, Phone2, Fax, Email, Website, CompanyCode, GuestTypeCode, MarketingCode,
// 	ComplimentHu, Notes string, Adult, Child int, Rate, RateOriginal, Discount, CommissionValue float64,
// 	DiscountPercent, IsAdditional, IsScheduledRate, IsBreakfast uint8, UserID string) error {

// 	var GuestInHouseData DBVar.Guest_in_house
// 	GuestInHouseData.FolioNumber = FolioNumber
// 	GuestInHouseData.AuditDate = AuditDate
// 	GuestInHouseData.GroupCode = GroupCode
// 	GuestInHouseData.RoomTypeCode = RoomTypeCode
// 	GuestInHouseData.BedTypeCode = BedTypeCode
// 	GuestInHouseData.RoomNumber = RoomNumber
// 	GuestInHouseData.RoomRateCode = RoomRateCode
// 	GuestInHouseData.BusinessSourceCode = BusinessSourceCode
// 	GuestInHouseData.CommissionTypeCode = CommissionTypeCode
// 	GuestInHouseData.PaymentTypeCode = PaymentTypeCode
// 	GuestInHouseData.MarketCode = MarketCode
// 	GuestInHouseData.TitleCode = TitleCode
// 	GuestInHouseData.FullName = FullName
// 	GuestInHouseData.Street = Street
// 	GuestInHouseData.City = City
// 	GuestInHouseData.CityCode = CityCode
// 	GuestInHouseData.CountryCode = CountryCode
// 	GuestInHouseData.StateCode = StateCode
// 	GuestInHouseData.PostalCode = PostalCode
// 	GuestInHouseData.Phone1 = Phone1
// 	GuestInHouseData.Phone2 = Phone2
// 	GuestInHouseData.Fax = Fax
// 	GuestInHouseData.Email = Email
// 	GuestInHouseData.Website = Website
// 	GuestInHouseData.CompanyCode = CompanyCode
// 	GuestInHouseData.GuestTypeCode = GuestTypeCode
// 	GuestInHouseData.MarketingCode = MarketingCode
// 	GuestInHouseData.ComplimentHu = ComplimentHu
// 	GuestInHouseData.Notes = Notes
// 	GuestInHouseData.Adult = Adult
// 	GuestInHouseData.Child = Child
// 	GuestInHouseData.Rate = Rate
// 	GuestInHouseData.RateOriginal = RateOriginal
// 	GuestInHouseData.Discount = Discount
// 	GuestInHouseData.CommissionValue = CommissionValue
// 	GuestInHouseData.DiscountPercent = DiscountPercent
// 	GuestInHouseData.IsBreakfast = IsBreakfast
// 	GuestInHouseData.IsScheduledRate = IsScheduledRate
// 	GuestInHouseData.IsAdditional = IsAdditional
// 	GuestInHouseData.CreatedBy = UserID

// 	err := DB.Table(DBVar.TableName.GuestInHouse).Create(&GuestInHouseData).Error

// 	if err != nil {
// 		return err
// 	}

// 	return nil

// }

func PostingRoomChargeManual(ctx context.Context, c *gin.Context, DB *gorm.DB, Dataset *global_var.TDataset, FolioNumber uint64, SubFolioGroupCode, SubDepartmentCode, CurrencyCode, Remark, DocumentNumber string, Amount, ExchangeRate float64, IsCanAutoTransferred bool, UserID string) error {
	ctx, span := global_var.Tracer.Start(ctx, "PostingRoomChargeManual")
	defer span.End()

	type DataOutputStruct struct {
		Folio         DBVar.Folio              `gorm:"embedded"`
		ContactPerson DBVar.Contact_person     `gorm:"embedded"`
		GuestDetail   DBVar.Guest_detail       `gorm:"embedded"`
		GuestGeneral  DBVar.Guest_general      `gorm:"embedded"`
		RoomRate      DBVar.Cfg_init_room_rate `gorm:"embedded"`
		DateArrival   time.Time
	}

	var DataOutput DataOutputStruct
	var RoomRateAmountOriginal, RoomRateAmount, RoomChargeB4Breakdown, RoomChargeAfterBreakdown, RoomChargeBasic, RoomChargeTax, RoomChargeService, TotalBreakdown, Commission, BreakdownAmount, BreakdownBasic, BreakdownTax, BreakdownService float64
	var RoomNumber string
	var BusinessSourceCode string
	var IsBreakfast uint8
	var CorrectionBreakdown, BreakDown1 uint64

	if err := DB.Table(DBVar.TableName.Folio).Select(
		" folio.group_code,"+
			" guest_detail.currency_code,"+
			" guest_detail.exchange_rate,"+
			" guest_general.purpose_of_code,"+
			" guest_general.sales_code,"+
			" folio.voucher_number,"+
			" guest_general.notes,"+
			" folio.number,"+
			" folio.compliment_hu,"+
			" DATE(guest_detail.arrival) AS DateArrival,"+
			" contact_person.title_code,"+
			" contact_person.full_name,"+
			" contact_person.street,"+
			" contact_person.city_code,"+
			" contact_person.city,"+
			" contact_person.nationality_code,"+
			" contact_person.country_code,"+
			" contact_person.state_code,"+
			" contact_person.postal_code,"+
			" contact_person.phone1,"+
			" contact_person.phone2,"+
			" contact_person.fax,"+
			" contact_person.email,"+
			" contact_person.website,"+
			" contact_person.company_code,"+
			" contact_person.guest_type_code,"+
			" contact_person.custom_field01,"+
			" contact_person.custom_field02,"+
			" contact_person.custom_lookup_field_code01,"+
			" contact_person.custom_lookup_field_code02,"+
			" guest_detail.adult,"+
			" guest_detail.child,"+
			" guest_detail.room_type_code,"+
			" guest_detail.bed_type_code,"+
			" guest_detail.room_number,"+
			" guest_detail.room_rate_code,"+
			" guest_detail.weekday_rate,"+
			" guest_detail.weekend_rate,"+
			" guest_detail.discount_percent,"+
			" guest_detail.discount,"+
			" guest_detail.business_source_code,"+
			" guest_detail.commission_type_code,"+
			" guest_detail.commission_value,"+
			" guest_detail.payment_type_code,"+
			" guest_detail.market_code,"+
			" guest_detail.booking_source_code,"+
			" cfg_init_room_rate.tax_and_service_code,"+
			" cfg_init_room_rate.charge_frequency_code ").
		Joins("LEFT OUTER JOIN contact_person ON (folio.contact_person_id1 = contact_person.id)").
		Joins("LEFT OUTER JOIN guest_general ON (folio.guest_general_id = guest_general.id)").
		Joins("LEFT OUTER JOIN guest_detail ON (folio.guest_detail_id = guest_detail.id)").
		Joins("LEFT OUTER JOIN cfg_init_room_rate ON (guest_detail.room_rate_code = cfg_init_room_rate.code)").
		Where("folio.number = ?", FolioNumber).Take(&DataOutput).Error; err != nil {
		return err
	}

	if DataOutput.Folio.Number > 0 {
		if IsFolioHaveBreakfast(DB, Dataset, DataOutput.Folio.Number) {
			IsBreakfast = 1
		}
		RoomNumber = *DataOutput.GuestDetail.RoomNumber
		// RoomRateTaxServiceCode = DataOutput.RoomRate.TaxAndServiceCode

		if ExchangeRate <= 0 {
			ExchangeRate = 1
		}
		AuditDate := GetAuditDate(c, DB, false)
		RoomChargeB4Breakdown = Amount
		RoomRateAmount = RoomChargeB4Breakdown
		RoomRateAmountOriginal = GetRoomRateAmount(ctx, DB, Dataset, DataOutput.GuestDetail.RoomRateCode, AuditDate.String(), DataOutput.GuestDetail.Adult, *DataOutput.GuestDetail.Child, false)

		//Cari Total Tax Service Room Charge
		RoomChargeBasic, RoomChargeTax, RoomChargeService = GetBasicTaxService(DB, Dataset.GlobalAccount.RoomCharge, DataOutput.RoomRate.TaxAndServiceCode, RoomChargeB4Breakdown)
		RoomChargeB4Breakdown = RoomChargeBasic + RoomChargeTax + RoomChargeService

		var DataOutputGuestBreakdown []DBVar.Guest_breakdown
		//Proses Query Breakdown
		DB.Table(DBVar.TableName.GuestBreakdown).Select(
			" guest_breakdown.outlet_code,"+
				" guest_breakdown.product_code,"+
				" guest_breakdown.sub_department_code,"+
				" guest_breakdown.account_code,"+
				" guest_breakdown.company_code,"+
				" guest_breakdown.quantity,"+
				" guest_breakdown.is_amount_percent,"+
				" guest_breakdown.amount,"+
				" guest_breakdown.per_pax,"+
				" guest_breakdown.include_child,"+
				" guest_breakdown.remark,"+
				" guest_breakdown.tax_and_service_code,"+
				" guest_breakdown.charge_frequency_code,"+
				" guest_breakdown.max_pax,"+
				" guest_breakdown.extra_pax,"+
				" guest_breakdown.per_pax_extra,"+
				" guest_breakdown.created_by,"+
				" guest_breakdown.id").
			Where("folio_number = ?", FolioNumber).Find(&DataOutputGuestBreakdown)

		//Calculate Total Breakdown
		TotalBreakdown = 0
		for _, guestBreakdown := range DataOutputGuestBreakdown {
			if IsCanPostCharge(c, DB, guestBreakdown.ChargeFrequencyCode, DataOutput.DateArrival) {
				if guestBreakdown.IsAmountPercent > 0 {
					BreakdownAmount = GetTotalBreakdownAmount(guestBreakdown.Quantity, Amount*guestBreakdown.Amount/100, guestBreakdown.ExtraPax,
						guestBreakdown.PerPax > 0, guestBreakdown.IncludeChild > 0, guestBreakdown.PerPaxExtra > 0, guestBreakdown.MaxPax, DataOutput.GuestDetail.Adult, *DataOutput.GuestDetail.Child) / ExchangeRate
				} else {
					BreakdownAmount = GetTotalBreakdownAmount(guestBreakdown.Quantity, guestBreakdown.Amount, guestBreakdown.ExtraPax,
						guestBreakdown.PerPax > 0, guestBreakdown.IncludeChild > 0, guestBreakdown.PerPaxExtra > 0, guestBreakdown.MaxPax, DataOutput.GuestDetail.Adult, *DataOutput.GuestDetail.Child) / ExchangeRate
				}

				BreakdownBasic, BreakdownTax, BreakdownService = GetBasicTaxService(DB, guestBreakdown.AccountCode, guestBreakdown.TaxAndServiceCode, BreakdownAmount)
				BreakdownAmount = BreakdownBasic + BreakdownTax + BreakdownService
				TotalBreakdown = TotalBreakdown + BreakdownAmount
			}
		}
		//Get Commission from Business Source
		BusinessSourceCode = *DataOutput.GuestDetail.BusinessSourceCode
		Commission = 0
		if BusinessSourceCode != "" {
			Commission = GetCommission(c, DB, *DataOutput.GuestDetail.CommissionTypeCode, *DataOutput.GuestDetail.CommissionValue, RoomChargeB4Breakdown, RoomChargeBasic, DataOutput.GuestDetail.Arrival) / ExchangeRate
		}

		//Room Charge - Total Breakdown - Total Commission
		RoomChargeAfterBreakdown = RoomChargeB4Breakdown - TotalBreakdown - Commission
		if RoomChargeAfterBreakdown > 0 {
			RoomChargeBasic, RoomChargeTax, RoomChargeService = GetBasicTaxService2(DB, Dataset.GlobalAccount.RoomCharge, DataOutput.RoomRate.TaxAndServiceCode, RoomChargeAfterBreakdown)
			RoomChargeAfterBreakdown = RoomChargeBasic + RoomChargeTax + RoomChargeService

			//Post Room Charge
			CorrectionBreakdown = GetSubFolioCorrectionBreakdown(c, DB)
			BreakDown1 = GetSubFolioBreakdown1(c, DB)
			if Remark == "" {
				Remark = "Additional Room Charge"
			} else {
				Remark = "Additional Room Charge - " + Remark
			}

			err := DB.Transaction(func(tx *gorm.DB) error {
				_, _, err := InsertSubFolio2(c, tx, Dataset, FolioNumber, SubFolioGroupCode, RoomNumber, Dataset.GlobalSubDepartment.FrontOffice,
					Dataset.GlobalAccount.RoomCharge, Dataset.GlobalAccount.RoomCharge, "", "",
					CurrencyCode, Remark, DocumentNumber, "", GlobalVar.TransactionType.Debit, "", "", CorrectionBreakdown,
					BreakDown1, GlobalVar.SubFolioPostingType.Room, 0, RoomChargeBasic, RoomChargeTax, RoomChargeService,
					ExchangeRate, false, IsCanAutoTransferred, UserID)

				if err != nil {
					return err
				}

				AuditDate := GetAuditDate(c, tx, false)
				if !IsInHousePosted(tx, AuditDate, FolioNumber) {
					err = InsertGuestInHouse(tx, AuditDate,
						FolioNumber,
						DataOutput.Folio.GroupCode,
						DataOutput.GuestDetail.RoomTypeCode,
						DataOutput.GuestDetail.BedTypeCode,
						*DataOutput.GuestDetail.RoomNumber,
						DataOutput.GuestDetail.RoomRateCode,
						*DataOutput.GuestDetail.BusinessSourceCode,
						*DataOutput.GuestDetail.CommissionTypeCode,
						DataOutput.GuestDetail.PaymentTypeCode,
						*DataOutput.GuestDetail.MarketCode,
						*DataOutput.ContactPerson.TitleCode,
						*DataOutput.ContactPerson.FullName,
						*DataOutput.ContactPerson.Street,
						*DataOutput.ContactPerson.City,
						*DataOutput.ContactPerson.CityCode,
						*DataOutput.ContactPerson.CountryCode,
						*DataOutput.ContactPerson.StateCode,
						*DataOutput.ContactPerson.PostalCode,
						*DataOutput.ContactPerson.Phone1,
						*DataOutput.ContactPerson.Phone2,
						*DataOutput.ContactPerson.Fax,
						*DataOutput.ContactPerson.Email,
						*DataOutput.ContactPerson.Website,
						*DataOutput.ContactPerson.CompanyCode,
						*DataOutput.ContactPerson.GuestTypeCode,
						*DataOutput.GuestGeneral.SalesCode,
						DataOutput.Folio.ComplimentHu,
						*DataOutput.GuestGeneral.Notes,
						DataOutput.GuestDetail.Adult,
						*DataOutput.GuestDetail.Child,
						RoomRateAmountOriginal,
						RoomRateAmount,
						*DataOutput.GuestDetail.Discount,
						*DataOutput.GuestDetail.CommissionValue,
						*DataOutput.GuestDetail.DiscountPercent,
						0,
						0,
						IsBreakfast,
						*DataOutput.GuestDetail.BookingSourceCode,
						*DataOutput.GuestGeneral.PurposeOfCode,
						*DataOutput.ContactPerson.CustomField01,
						*DataOutput.ContactPerson.CustomField02, 0, "",
						*DataOutput.ContactPerson.NationalityCode,
						UserID)
					if err != nil {
						return err
					}
				}

				for _, guestBreakdown := range DataOutputGuestBreakdown {
					if IsCanPostCharge(c, tx, guestBreakdown.ChargeFrequencyCode, DataOutput.DateArrival) {
						if guestBreakdown.IsAmountPercent > 0 {
							BreakdownAmount = GetTotalBreakdownAmount(guestBreakdown.Quantity, Amount*guestBreakdown.Amount/100, guestBreakdown.ExtraPax, guestBreakdown.PerPax > 0, guestBreakdown.IncludeChild > 0, guestBreakdown.PerPaxExtra > 0, guestBreakdown.MaxPax, DataOutput.GuestDetail.Adult, *DataOutput.GuestDetail.Child) / ExchangeRate
						} else {
							BreakdownAmount = GetTotalBreakdownAmount(guestBreakdown.Quantity, guestBreakdown.Amount, guestBreakdown.ExtraPax, guestBreakdown.PerPax > 0, guestBreakdown.IncludeChild > 0, guestBreakdown.PerPaxExtra > 0, guestBreakdown.MaxPax, DataOutput.GuestDetail.Adult, *DataOutput.GuestDetail.Child) / ExchangeRate
						}
						SubFolioData := DBVar.Sub_folio{}
						SubFolioData.FolioNumber = FolioNumber
						SubFolioData.GroupCode = SubFolioGroupCode
						SubFolioData.RoomNumber = RoomNumber
						SubFolioData.SubDepartmentCode = SubDepartmentCode
						SubFolioData.AccountCode = guestBreakdown.AccountCode
						SubFolioData.ProductCode = guestBreakdown.ProductCode
						SubFolioData.CurrencyCode = CurrencyCode
						SubFolioData.Remark = "Breakdown: " + guestBreakdown.Remark
						SubFolioData.TypeCode = GlobalVar.TransactionType.Debit
						SubFolioData.CorrectionBreakdown = CorrectionBreakdown
						SubFolioData.Breakdown1 = BreakDown1
						SubFolioData.DirectBillCode = guestBreakdown.CompanyCode
						SubFolioData.PostingType = GlobalVar.SubFolioPostingType.Room
						SubFolioData.ExtraChargeId = guestBreakdown.Id
						SubFolioData.Quantity = 1
						SubFolioData.Amount = BreakdownAmount
						SubFolioData.ExchangeRate = ExchangeRate
						SubFolioData.CreatedBy = UserID

						_, _, err = InsertSubFolio(c, tx, Dataset, IsCanAutoTransferred, Dataset.GlobalAccount.RoomCharge, guestBreakdown.TaxAndServiceCode, SubFolioData)
						if err != nil {
							return err
						}
						// Inset guest in house breakdown
						err = InsertGuestInHouseBreakdown(tx, AuditDate, FolioNumber, guestBreakdown.OutletCode, guestBreakdown.ProductCode, guestBreakdown.SubDepartmentCode,
							guestBreakdown.AccountCode, guestBreakdown.Remark, guestBreakdown.TaxAndServiceCode,
							guestBreakdown.ChargeFrequencyCode, guestBreakdown.Quantity, guestBreakdown.Amount,
							guestBreakdown.ExtraPax, guestBreakdown.MaxPax, guestBreakdown.PerPax,
							guestBreakdown.IncludeChild, guestBreakdown.PerPaxExtra, UserID)
						if err != nil {
							return err
						}
					}
				}

				//Posting Commission from Business Source
				if Commission > 0 {
					SubFolioData := DBVar.Sub_folio{}
					SubFolioData.FolioNumber = FolioNumber
					SubFolioData.GroupCode = SubFolioGroupCode
					SubFolioData.RoomNumber = RoomNumber
					SubFolioData.SubDepartmentCode = Dataset.GlobalSubDepartment.FrontOffice
					SubFolioData.AccountCode = Dataset.GlobalAccount.APCommission
					SubFolioData.ProductCode = Dataset.GlobalAccount.RoomCharge
					SubFolioData.CurrencyCode = CurrencyCode
					SubFolioData.Remark = "Breakdown Commission"
					SubFolioData.TypeCode = GlobalVar.TransactionType.Debit
					SubFolioData.CorrectionBreakdown = CorrectionBreakdown
					SubFolioData.Breakdown1 = BreakDown1
					SubFolioData.DirectBillCode = BusinessSourceCode
					SubFolioData.PostingType = GlobalVar.SubFolioPostingType.Room
					SubFolioData.Quantity = 1
					SubFolioData.Amount = Commission
					SubFolioData.ExchangeRate = ExchangeRate
					SubFolioData.CreatedBy = UserID
					_, _, err = InsertSubFolio(c, tx, Dataset, IsCanAutoTransferred, Dataset.GlobalAccount.RoomCharge, "", SubFolioData)
					if err != nil {
						return err
					}
				}

				return nil
			})
			return err
		}
	}
	return nil
}

func UpdateSubFolioGroupByFolioNumber(DB *gorm.DB, FolioNumber uint64, SubFolioGroupCode, UserID string) error {
	err := DB.Table(DBVar.TableName.SubFolio).
		Where("folio_number=? AND void='0'", FolioNumber).
		Update("group_code", SubFolioGroupCode).Update("updated_by", UserID).Error

	return err
}

func UpdateSubFolioTransferPairID(DB *gorm.DB, SubFolioID1, SubFolioID2 uint64, UserID string) (err error) {
	err = DB.Table(DBVar.TableName.SubFolio).
		Where("id=?", SubFolioID1).Updates(&map[string]interface{}{
		"is_pair_with_deposit": 0,
		"transfer_pair_id":     SubFolioID2,
		"updated_by":           UserID,
	}).Error

	err = DB.Table(DBVar.TableName.SubFolio).
		Where("id=?", SubFolioID2).Updates(&map[string]interface{}{
		"is_pair_with_deposit": 0,
		"transfer_pair_id":     SubFolioID1,
		"updated_by":           UserID,
	}).Error

	return err
}

func MoveSubFolioByFolioNumber(DB *gorm.DB, FolioNumberFrom, FolioNumberTo uint64, SubFolioGroupCode, UserID string) error {
	err := DB.Table(DBVar.TableName.SubFolio).
		Select("folio_number", "group_code", "updated_by").
		Where("folio_number = ? AND void='0'", FolioNumberFrom).
		Updates(map[string]interface{}{"folio_number": FolioNumberTo, "group_code": SubFolioGroupCode, "updated_by": UserID}).Error

	return err
}

func MoveSubFolioByBreakdown(DB *gorm.DB, FolioNumber uint64, SubFolioGroupCode string, CorrectionBreakdown uint64, UserID string) error {
	err := DB.Table(DBVar.TableName.SubFolio).
		Select("folio_number", "group_code", "updated_by").
		Where("correction_breakdown = ? AND void='0'", CorrectionBreakdown).
		Updates(map[string]interface{}{"folio_number": FolioNumber, "group_code": SubFolioGroupCode, "updated_by": UserID}).Error
	return err
}

func GetTransferBalanceRemark(DB *gorm.DB, FolioNumber uint64) string {
	type DataOutputStruct struct {
		RoomNumber string
		FullName   string
	}

	var DataOutput DataOutputStruct
	DB.Table(DBVar.TableName.Folio).Select(
		" guest_detail.room_number,"+" CONCAT(contact_person.title_code, contact_person.full_name) AS FullName ").
		Joins("LEFT OUTER JOIN guest_detail ON (folio.guest_detail_id = guest_detail.id)").
		Joins("LEFT OUTER JOIN contact_person ON (folio.contact_person_id1 = contact_person.id)").
		Where("folio.number = ?", FolioNumber).
		Take(&DataOutput)

	return DataOutput.RoomNumber + " / " + DataOutput.FullName
}

func GetFolioReturns(DB *gorm.DB, FolioNumber string) ([]map[string]interface{}, error) {
	var DataOutput []map[string]interface{}

	err := DB.Raw(
		"SELECT" +
			" number," +
			" FullName," +
			" room_number," +
			" name " +
			"FROM ((" +
			"SELECT" +
			" folio.number," +
			" CONCAT(contact_person.title_code, contact_person.full_name) AS FullName," +
			" guest_detail.room_number," +
			" const_folio_type.name," +
			" const_folio_type.id_sort " +
			"FROM" +
			" folio_routing" +
			" LEFT OUTER JOIN folio ON (folio_routing.folio_number = folio.number)" +
			" LEFT OUTER JOIN contact_person ON (folio.contact_person_id1 = contact_person.id)" +
			" LEFT OUTER JOIN guest_detail ON (folio.guest_detail_id = guest_detail.id)" +
			" LEFT OUTER JOIN const_folio_type ON (folio.type_code = const_folio_type.code)" +
			" WHERE folio_routing.folio_transfer='" + FolioNumber + "'" +
			" AND folio.status_code='" + GlobalVar.FolioStatus.Open + "' " +
			"GROUP BY folio_routing.folio_number" +
			")UNION(" +
			"SELECT" +
			" sub_folio.belongs_to AS number," +
			" CONCAT(contact_person.title_code, contact_person.full_name) AS FullName," +
			" guest_detail.room_number," +
			" const_folio_type.name," +
			" const_folio_type.id_sort " +
			"FROM" +
			" sub_folio" +
			" LEFT OUTER JOIN folio ON (sub_folio.belongs_to = folio.number)" +
			" LEFT OUTER JOIN contact_person ON (folio.contact_person_id1 = contact_person.id)" +
			" LEFT OUTER JOIN guest_detail ON (folio.guest_detail_id = guest_detail.id)" +
			" LEFT OUTER JOIN const_folio_type ON (folio.type_code = const_folio_type.code)" +
			" WHERE sub_folio.folio_number='" + FolioNumber + "'" +
			" AND sub_folio.folio_number<>sub_folio.belongs_to" +
			" AND folio.status_code='" + GlobalVar.FolioStatus.Open + "' " +
			"GROUP BY sub_folio.belongs_to)) AS FolioReturn " +
			"GROUP BY number " +
			"ORDER BY id_sort, room_number, number").Scan(&DataOutput).Error

	return DataOutput, err
}

func GetInvoiceNumber(DB *gorm.DB, IssuedDate time.Time) string {
	ServerID := GetServerID(DB)
	Prefix := GlobalVar.ConstProgramVariable.InvoiceNumberPrefix + General.FormatDatePrefix(IssuedDate) + strconv.Itoa(ServerID) + "-"
	var DataOutput uint64
	DB.Raw(
		"SELECT" +
			" CAST(RIGHT(number, LENGTH(number) - " + strconv.FormatInt(int64(len(Prefix)), 10) + ") AS UNSIGNED) AS MaxInvoiceNumber " +
			"FROM" +
			" invoice" +
			" WHERE LEFT(number," + strconv.FormatInt(int64(len(Prefix)), 10) + ")='" + Prefix + "' " +
			"ORDER BY MaxInvoiceNumber DESC " +
			"LIMIT 1;").Scan(&DataOutput)

	return Prefix + General.Uint64ToStr(DataOutput+1)
}
func GetReceiptNumber(DB *gorm.DB, IssuedDate time.Time) string {
	ServerID := GetServerID(DB)
	Prefix := General.FormatDatePrefix(IssuedDate) + strconv.Itoa(ServerID) + "-"
	var DataOutput uint64
	DB.Raw(
		"SELECT" +
			" CAST(RIGHT(number, LENGTH(number) - " + strconv.FormatInt(int64(len(Prefix)), 10) + ") AS UNSIGNED) AS MaxReceiptNumber " +
			"FROM" +
			" receipt" +
			" WHERE LEFT(number," + strconv.FormatInt(int64(len(Prefix)), 10) + ")='" + Prefix + "' " +
			"ORDER BY MaxReceiptNumber DESC " +
			"LIMIT 1;").Scan(&DataOutput)
	//fmt.Println(ServerID)
	//fmt.Println(Prefix)
	//fmt.Println(DataOutput)
	//fmt.Println(Prefix + General.Uint64ToStr(DataOutput+1))

	return Prefix + General.Uint64ToStr(DataOutput+1)
}

func GetPRNumber(DB *gorm.DB, IssuedDate time.Time) string {
	ServerID := GetServerID(DB)
	Prefix := GlobalVar.ConstProgramVariable.PRNumberPrefix + General.FormatDatePrefix(IssuedDate) + strconv.Itoa(ServerID) + "-"
	var DataOutput uint64
	DB.Raw(
		"SELECT" +
			" CAST(RIGHT(number, LENGTH(number) - " + strconv.FormatInt(int64(len(Prefix)), 10) + ") AS UNSIGNED) AS MaxPRNumber " +
			"FROM" +
			" inv_purchase_request" +
			" WHERE LEFT(number," + strconv.FormatInt(int64(len(Prefix)), 10) + ")='" + Prefix + "' " +
			"ORDER BY MaxPRNumber DESC " +
			"LIMIT 1;").Scan(&DataOutput)

	return Prefix + General.Uint64ToStr(DataOutput+1)
}

func GetPONumber(DB *gorm.DB, IssuedDate time.Time) string {
	ServerID := GetServerID(DB)
	Prefix := GlobalVar.ConstProgramVariable.PONumberPrefix + General.FormatDatePrefix(IssuedDate) + strconv.Itoa(ServerID) + "-"
	var DataOutput uint64
	DB.Raw(
		"SELECT" +
			" CAST(RIGHT(number, LENGTH(number) - " + strconv.FormatInt(int64(len(Prefix)), 10) + ") AS UNSIGNED) AS MaxPONumber " +
			"FROM" +
			" inv_purchase_order" +
			" WHERE LEFT(number," + strconv.FormatInt(int64(len(Prefix)), 10) + ")='" + Prefix + "' " +
			"ORDER BY MaxPONumber DESC " +
			"LIMIT 1;").Scan(&DataOutput)

	return Prefix + General.Uint64ToStr(DataOutput+1)
}

func GetJournalAccountCompanyAR(DB *gorm.DB, CompanyCode string) string {
	var DataOutput string
	DB.Table(DBVar.TableName.Company).Select("cfg_init_company_type.journal_account_code_ar").Joins("LEFT OUTER JOIN cfg_init_company_type ON (company.type_code = cfg_init_company_type.code)").
		Where("company.code = ?", CompanyCode).Scan(&DataOutput)

	return DataOutput

}

func GetJournalIDSort(DB *gorm.DB, PostingDate time.Time) int {
	var IdSort int
	DB.Table(DBVar.TableName.AccJournal).Select("id_sort").Where("date=?", PostingDate).Order("id_sort DESC").Limit(1).Scan(&IdSort)

	return IdSort + 1
}

func IsInvoiceHadPayment(DB *gorm.DB, InvoiceNumberX string) bool {
	var InvoiceNumber string
	DB.Table(DBVar.TableName.InvoicePayment).Select("invoice_number").Where("invoice_number=?", InvoiceNumber).Take(&InvoiceNumber)
	return InvoiceNumber != ""
}

func IsSubFolioIdHadInvoice(DB *gorm.DB, SubFolioId uint64) bool {
	var SubFolioIdX uint64
	DB.Table(DBVar.TableName.InvoiceItem).Select("sub_folio_id").Where("sub_folio_id=?", SubFolioId).Limit(1).Scan(&SubFolioIdX)

	return SubFolioIdX > 0
}

func GetFolioType(DB *gorm.DB, FolioNumber uint64) string {
	var Type string
	DB.Table(DBVar.TableName.Folio).Select("type_code").Where("number=?", FolioNumber).Limit(1).Scan(&Type)

	return Type
}

func GetFolioSystemCode(DB *gorm.DB, FolioNumber uint64) string {
	var Code string
	DB.Table(DBVar.TableName.Folio).Select("system_code").Where("number=?", FolioNumber).Limit(1).Scan(&Code)

	return Code
}
func IsThereCardActiveFolio(DB *gorm.DB, FolioNumber uint64) bool {
	return MasterData.GetFieldBool(DB, DBVar.TableName.LogKeylock, "id", "is_active = '1' AND folio_number", General.Uint64ToStr(FolioNumber), false)
}

func IsCanCreateInvoice(DB *gorm.DB, FolioNumber string) (DirectBillCode string, CanCreate bool) {
	DB.Table(DBVar.TableName.SubFolio).Select("sub_folio.direct_bill_code").
		Joins("LEFT OUTER JOIN cfg_init_account ON (sub_folio.account_code = cfg_init_account.code)").
		Where("sub_folio.folio_number=?", FolioNumber).
		Where("sub_folio.direct_bill_code<>''").
		Where("cfg_init_account.sub_group_code=?", GlobalVar.GlobalAccountSubGroup.AccountReceivable).
		Where("void='0'").
		Group("sub_folio.correction_breakdown").
		Limit(1).
		Scan(&DirectBillCode)

	return DirectBillCode, DirectBillCode != ""
}

func CheckFolioReceiveTransfer(DB *gorm.DB, FolioNumber uint64) []string {
	var FolioTransferMessage []string
	var FolioDetail []string

	DB.Table(DBVar.TableName.FolioRouting).Select(" DISTINCT CONCAT(folio.number, '/Room: ', guest_detail.room_number, '/', contact_person.title_code, contact_person.full_name) AS FolioDetail").
		Joins("LEFT OUTER JOIN folio ON (folio_routing.folio_number = folio.number)").
		Joins("LEFT OUTER JOIN contact_person ON (folio.contact_person_id1 = contact_person.id)").
		Joins("LEFT OUTER JOIN guest_detail ON (folio.guest_detail_id = guest_detail.id)").
		Where("folio_routing.folio_transfer=?", FolioNumber).
		Scan(&FolioDetail)

	for _, detail := range FolioDetail {
		FolioTransferMessage = append(FolioTransferMessage, detail)
	}

	return FolioTransferMessage
}

func GetInvoiceTotalPayment(DB *gorm.DB, InvoiceNumber string) (TotalPayment float64) {
	DB.Table(DBVar.TableName.InvoicePayment).Select("SUM(amount_foreign) AS TotalAmount").
		Where("invoice_number=?", InvoiceNumber).
		Group("invoice_number").Scan(&TotalPayment)

	return TotalPayment
}

func IsPrepaidExpensePosted(DB *gorm.DB, PrepaidID, PrepaidPostedID uint64, PostingDate time.Time) bool {
	var Id uint64
	Year, Month, _ := PostingDate.Date()
	Query := DB.Table(DBVar.TableName.AccPrepaidExpensePosted).Select("id").
		Where("prepaid_id=?", PrepaidID).
		Where("MONTH(posting_date)=?", Month).
		Where("YEAR(posting_date)=?", Year)

	if PrepaidPostedID > 0 {
		Query.Where("id<>?", PrepaidPostedID)
	}

	Query.Limit(1).Scan(&Id)
	return Id > 0
}

func IsDifferedIncomePosted(DB *gorm.DB, DifferedID, DifferedPostedID uint64, PostingDate time.Time) bool {
	var Id uint64
	Year, Month, _ := PostingDate.Date()
	Query := DB.Table(DBVar.TableName.AccDefferedIncomePosted).Select("id").
		Where("deffered_id=?", DifferedID).
		Where("MONTH(posting_date)=?", Month).
		Where("YEAR(posting_date)=?", Year)

	if DifferedPostedID > 0 {
		Query.Where("id<>?", DifferedPostedID)
	}

	Query.Limit(1).Scan(&Id)
	return Id > 0
}

func GetJournalAccountCode(DB *gorm.DB, AccountCode string) string {
	var JournalAccountCode string
	DB.Table(DBVar.TableName.CfgInitAccount).Select("journal_account_code").Where("code=?", AccountCode).Limit(1).Scan(&JournalAccountCode)

	return JournalAccountCode
}

func GetJournalAccountCurrency(DB *gorm.DB, CurrencyCode string) (JournalAccountCode string) {
	DB.Table(DBVar.TableName.CfgInitCurrency).Select("cfg_init_account.journal_account_code").
		Joins("LEFT OUTER JOIN cfg_init_account ON (cfg_init_currency.account_code = cfg_init_account.code)").
		Where("cfg_init_currency.code=?", CurrencyCode).
		Limit(1).
		Scan(&JournalAccountCode)

	return
}

func GetForeignCashBreakdown(DB *gorm.DB) (BreakdownId uint64) {
	DB.Table(DBVar.TableName.AccForeignCash).Select("breakdown").
		Order("breakdown DESC").
		Limit(1).
		Scan(&BreakdownId)
	BreakdownId++
	return
}

func GetJournalAccountCodeFromTransaction(DB *gorm.DB, IdTable, IdTransaction uint64) (AccountCode string) {
	MainTableName := DBVar.TableName.SubFolio
	if IdTable == 1 {
		MainTableName = DBVar.TableName.GuestDeposit
	}
	DB.Table(MainTableName).Select("cfg_init_account.journal_account_code").
		Joins("LEFT OUTER JOIN cfg_init_account ON ("+MainTableName+".account_code = cfg_init_account.code)").
		Where(MainTableName+".id=?", IdTransaction).
		Limit(1).
		Scan(&AccountCode)

	return
}

func GetJournalAccountListByGroup(DB *gorm.DB, SubDepartmentCode, GroupCode1, GroupCode2, GroupCode3 string) (DataOutput []map[string]interface{}) {
	Query := DB.Table(DBVar.TableName.CfgInitJournalAccount).Select(
		" cfg_init_journal_account.code,"+
			" cfg_init_journal_account.name").
		Joins("LEFT OUTER JOIN cfg_init_journal_account_sub_group ON (cfg_init_journal_account.sub_group_code = cfg_init_journal_account_sub_group.code)").
		Where("cfg_init_journal_account_sub_group.group_code=?", GroupCode1)

	if GroupCode2 != "" {
		Query.Or("cfg_init_journal_account_sub_group.group_code=?", GroupCode2)
	}
	if GroupCode3 != "" {
		Query.Or("cfg_init_journal_account_sub_group.group_code=?", GroupCode3)
	}

	Query.Where("sub_department_code LIKE ?", "%"+SubDepartmentCode+"%").Order("cfg_init_journal_account.code").
		Scan(&DataOutput)
	return
}

func GetPurchaseRequestStatus(DB *gorm.DB, Id uint64) (Status string) {
	DB.Table(DBVar.TableName.InvPurchaseRequest).Select("status_code").Where("id=?", Id).Limit(1).Scan(&Status)
	if Status == "" {
		Status = GlobalVar.PurchaseRequestStatus.NotApproved
	}
	return
}

func IsPurchaseRequestApproved1(DB *gorm.DB, Id uint64) (IsApprove bool) {
	var Number string
	DB.Table(DBVar.TableName.InvPurchaseRequest).Select("number").Where("id = ? AND is_user_approved1=1", Id).Limit(1).Scan(&Number)
	IsApprove = Number != ""
	return
}

func IsPurchaseRequestApproved12(DB *gorm.DB, Id uint64) (IsApprove bool) {
	var Number string
	DB.Table(DBVar.TableName.InvPurchaseRequest).Select("number").Where("id = ? AND is_user_approved12=1", Id).Limit(1).Scan(&Number)
	IsApprove = Number != ""
	return
}

func IsPurchaseRequestApproved2(DB *gorm.DB, Id uint64) (IsApprove bool) {
	var Number string
	DB.Table(DBVar.TableName.InvPurchaseRequest).Select("number").Where("id = ? AND is_user_approved2=1", Id).Limit(1).Scan(&Number)
	IsApprove = Number != ""
	return
}

func IsPurchaseRequestApproved3(DB *gorm.DB, Id uint64) (IsApprove bool) {
	var Number string
	DB.Table(DBVar.TableName.InvPurchaseRequest).Select("number").Where("id = ? AND is_user_approved3=1", Id).Limit(1).Scan(&Number)
	IsApprove = Number != ""
	return
}
func IsPurchaseRequestPriceApplied(DB *gorm.DB, Number string) (IsApplied bool) {
	var PRNumber string
	DB.Table(DBVar.TableName.InvPurchaseRequestDetail).Select("pr_number").Where("pr_number=? AND quantity_approved > 0", Number).Limit(1).Scan(&PRNumber)
	IsApplied = PRNumber != ""
	return
}

func CanUserApprovedInventoryPR1(DB *gorm.DB, UserID, Password, SubDepartmentCode string, IsPasswordEncrypted bool) bool {
	var UserCode string
	// if IsPasswordEncrypted {
	//   PasswordEncrypted = Password
	// } else {
	//   PasswordEncrypted = EncryptString(Password, VariableDLL.KeyOtherDLL)
	// }
	DB.Table(DBVar.TableName.AstUserSubDepartment).Select("user_code").Where("user_code=? AND sub_department_code=? AND is_can_inv_pr_approve1=1", UserID, SubDepartmentCode).Limit(1).Scan(&UserCode)
	// if IsUserPasswordValid(UserID, PasswordEncrypted) {
	//TODO please Review validation
	return UserCode != ""
}

func CanUserApprovedInventoryPR2(DB *gorm.DB, UserID, Password, SubDepartmentCode string, IsPasswordEncrypted bool) bool {
	var UserCode string
	// if IsPasswordEncrypted {
	//   PasswordEncrypted = Password
	// } else {
	//   PasswordEncrypted = EncryptString(Password, VariableDLL.KeyOtherDLL)
	// }
	DB.Table(DBVar.TableName.AstUserSubDepartment).Select("user_code").Where("user_code=? AND sub_department_code=? AND is_can_inv_pr_approve2=1", UserID, SubDepartmentCode).Limit(1).Scan(&UserCode)
	// if IsUserPasswordValid(UserID, PasswordEncrypted) {
	//TODO please Review validation
	return UserCode != ""
}

func CanUserApprovedInventoryPR3(DB *gorm.DB, UserID, Password, SubDepartmentCode string, IsPasswordEncrypted bool) bool {
	var UserCode string
	// if IsPasswordEncrypted {
	//   PasswordEncrypted = Password
	// } else {
	//   PasswordEncrypted = EncryptString(Password, VariableDLL.KeyOtherDLL)
	// }
	DB.Table(DBVar.TableName.AstUserSubDepartment).Select("user_code").Where("user_code=? AND sub_department_code=? AND is_can_inv_pr_approve3=1", UserID, SubDepartmentCode).Limit(1).Scan(&UserCode)
	// if IsUserPasswordValid(UserID, PasswordEncrypted) {
	//TODO please Review validation
	return UserCode != ""
}

func CanUserApprovedInventoryPR12(DB *gorm.DB, UserID, Password, SubDepartmentCode string, IsPasswordEncrypted bool) bool {
	var UserCode string
	// if IsPasswordEncrypted {
	//   PasswordEncrypted = Password
	// } else {
	//   PasswordEncrypted = EncryptString(Password, VariableDLL.KeyOtherDLL)
	// }
	DB.Table(DBVar.TableName.AstUserSubDepartment).Select("user_code").Where("user_code=? AND sub_department_code=? AND is_can_inv_pr_approve12=1", UserID, SubDepartmentCode).Limit(1).Scan(&UserCode)
	// if IsUserPasswordValid(UserID, PasswordEncrypted) {
	//TODO please Review validation
	return UserCode != ""
}

func CanUserApplyPrice(DB *gorm.DB, UserID, Password, SubDepartmentCode string, IsPasswordEncrypted bool) bool {
	var UserCode string
	// if IsPasswordEncrypted {
	//   PasswordEncrypted = Password
	// } else {
	//   PasswordEncrypted = EncryptString(Password, VariableDLL.KeyOtherDLL)
	// }
	DB.Table(DBVar.TableName.AstUserSubDepartment).Select("user_code").Where("user_code=? AND sub_department_code=? AND is_can_inv_pr_assign_price=1", UserID, SubDepartmentCode).Limit(1).Scan(&UserCode)
	// if IsUserPasswordValid(UserID, PasswordEncrypted) {
	//TODO please Review validation
	return UserCode != ""
}

func GetMarketListCheaperPrice(DB *gorm.DB, ItemCode, UOMCode string) map[string]interface{} {
	DataOutput2 := make(map[string]interface{})
	var DataOutput map[string]interface{}
	DB.Table(DBVar.TableName.InvCfgInitMarketList).Select(
		" inv_cfg_init_market_list.company_code,"+
			" IF(inv_cfg_init_market_list.uom_code=inv_cfg_init_item.uom_code,IFNULL(inv_cfg_init_market_list.price,0)*IFNULL(inv_cfg_init_item_uom1.quantity,1),"+"IF(IFNULL(inv_cfg_init_item_uom.quantity,0)=0,0,IFNULL(inv_cfg_init_market_list.price,0)/inv_cfg_init_item_uom.quantity*IFNULL(inv_cfg_init_item_uom1.quantity,0))) AS MarketPrice,"+
			" IFNULL(company.name,'') AS name").
		Joins("INNER JOIN company ON (inv_cfg_init_market_list.company_code = company.code)").
		Joins("LEFT OUTER JOIN inv_cfg_init_item_uom ON (inv_cfg_init_item_uom.item_code = ? AND inv_cfg_init_market_list.uom_code = inv_cfg_init_item_uom.uom_code)", ItemCode).
		Joins("LEFT OUTER JOIN inv_cfg_init_item_uom inv_cfg_init_item_uom1 ON (inv_cfg_init_item_uom1.item_code = ? AND inv_cfg_init_item_uom1.uom_code = ?)", ItemCode, UOMCode).
		Joins("LEFT OUTER JOIN inv_cfg_init_item ON (inv_cfg_init_market_list.item_code = inv_cfg_init_item.code)").
		Where("inv_cfg_init_market_list.item_code=?", ItemCode).Limit(1).Scan(&DataOutput)
	fmt.Println(DataOutput)
	if DataOutput == nil {
		return nil
	}

	DataOutput2["item_code"] = ItemCode
	DataOutput2["uom_code"] = UOMCode
	DataOutput2["company_code"] = ""
	DataOutput2["company_name"] = ""
	DataOutput2["price"] = 0
	if DataOutput != nil {
		DataOutput2["company_code"] = DataOutput["company_code"].(string)
		DataOutput2["company_name"] = DataOutput["name"].(string)
		DataOutput2["price"] = General.StrToFloat64(DataOutput["MarketPrice"].(string))
	}
	return DataOutput2
}

func GetReceiveLastPrice(DB *gorm.DB, ItemCode, UOMCode, StockDate string) map[string]interface{} {
	DataOutput2 := make(map[string]interface{})
	var DataOutput map[string]interface{}
	DB.Table(DBVar.TableName.InvReceivingDetail).Select(
		" inv_receiving_detail.receive_uom_code,"+
			" inv_receiving_detail.receive_price,"+
			" inv_receiving_detail.basic_uom_code,"+
			" inv_receiving_detail.basic_price ").
		Joins("LEFT OUTER JOIN inv_receiving ON (inv_receiving_detail.receive_number = inv_receiving.number)").
		Where("inv_receiving_detail.item_code=?", ItemCode).
		Where("inv_receiving.date=?", StockDate).
		Order("inv_receiving.date DESC").
		Limit(1).Scan(&DataOutput)

	DataOutput2["item_code"] = ItemCode
	DataOutput2["uom_code"] = UOMCode
	DataOutput2["price"] = 0
	fmt.Println("uom", UOMCode)
	fmt.Println("2", DataOutput2)
	fmt.Println("3", DataOutput)
	if DataOutput != nil {
		UOMCodeReceive := DataOutput["receive_uom_code"].(string)
		ReceivePrice := DataOutput["receive_price"].(string)
		BasicUOMCodeReceive := DataOutput["basic_uom_code"].(string)
		BasicPrice := DataOutput["basic_price"].(string)
		if UOMCode == UOMCodeReceive {
			DataOutput2["price"] = ReceivePrice
		} else if UOMCode == BasicUOMCodeReceive {
			DataOutput2["price"] = BasicPrice
		} else {
			DataOutput2["uom_code"] = UOMCode
			DataOutput2["price"] = ConvertPriceFromBasic(DB, ItemCode, UOMCode, General.StrToFloat64(BasicPrice))
		}
	}
	return DataOutput2
}

func ConvertPriceFromBasic(DB *gorm.DB, ItemCode, UOMCode string, BasicPrice float64) (Price float64) {
	Price = BasicPrice
	//Ambil dari Receive
	Quantity := MasterData.GetFieldFloatQuery(DB, "SELECT quantity FROM inv_cfg_init_item_uom"+
		" WHERE item_code=?"+
		" AND uom_code=?", 0, ItemCode, UOMCode)
	if Quantity > 0 {
		Price = BasicPrice * Quantity
	}
	return
}

func GetStockStoreUpdateStockTransfer(DB *gorm.DB, Dataset *global_var.TDataset, StoreCode, ItemCode, StockTransferNumber string, StockDate time.Time) (float64, error) {
	StockDateStr := General.FormatDate1(StockDate)
	StockInDate := GetStockInDate(DB, StoreCode, ItemCode, StockDate)
	StockInDateStr := General.FormatDate1(StockInDate)
	var Stock float64
	if Dataset.ProgramConfiguration.CostingMethod == GlobalVar.InventoryCostingMethod.Average {
		if err := DB.Raw(
			"SELECT SUM(IFNULL(Stock.Quantity, 0)) FROM ((" +
				"SELECT" +
				" SUM(IF(inv_receiving_detail.store_code='" + StoreCode + "', inv_receiving_detail.basic_quantity, 0)) AS Quantity " +
				"FROM" +
				" inv_receiving_detail" +
				" LEFT OUTER JOIN inv_receiving ON (inv_receiving_detail.receive_number = inv_receiving.number)" +
				" WHERE inv_receiving_detail.item_code='" + ItemCode + "' " +
				" AND inv_receiving.date<='" + General.FormatDate1(StockDate) + "' " +
				"GROUP BY inv_receiving_detail.item_code" +
				") UNION ALL (" +
				"SELECT" +
				" SUM(inv_stock_transfer_detail.quantity) AS Quantity " +
				"FROM" +
				" inv_stock_transfer_detail" +
				" LEFT OUTER JOIN inv_stock_transfer ON (inv_stock_transfer_detail.st_number = inv_stock_transfer.number)" +
				" WHERE inv_stock_transfer_detail.item_code='" + ItemCode + "' " +
				" AND inv_stock_transfer.date<='" + StockDateStr + "'" +
				" AND inv_stock_transfer_detail.to_store_code='" + StoreCode + "' " +
				" AND inv_stock_transfer_detail.st_number<>'" + StockTransferNumber + "' " +
				"GROUP BY inv_stock_transfer_detail.item_code" +
				") UNION ALL (" +
				"SELECT" +
				" -SUM(inv_stock_transfer_detail.quantity) AS Quantity " +
				"FROM" +
				" inv_stock_transfer_detail" +
				" LEFT OUTER JOIN inv_stock_transfer ON (inv_stock_transfer_detail.st_number = inv_stock_transfer.number)" +
				" WHERE inv_stock_transfer_detail.item_code='" + ItemCode + "' " +
				" AND inv_stock_transfer.date<='" + StockInDateStr + "'" +
				" AND inv_stock_transfer_detail.from_store_code='" + StoreCode + "'" +
				" AND inv_stock_transfer_detail.st_number<>'" + StockTransferNumber + "' " +
				"GROUP BY inv_stock_transfer_detail.item_code" +
				") UNION ALL (" +
				"SELECT" +
				" -SUM(inv_costing_detail.quantity) AS Quantity " +
				"FROM" +
				" inv_costing_detail" +
				" LEFT OUTER JOIN inv_costing ON (inv_costing_detail.costing_number = inv_costing.number)" +
				" WHERE inv_costing_detail.item_code='" + ItemCode + "'" +
				" AND inv_costing.date<='" + StockInDateStr + "'" +
				" AND inv_costing_detail.store_code='" + StoreCode + "' " +
				"GROUP BY inv_costing_detail.item_code" +
				") UNION ALL (" +
				"SELECT" +
				" -SUM(inv_costing_detail.quantity) AS Quantity " +
				"FROM" +
				" inv_costing_detail" +
				" LEFT OUTER JOIN inv_costing ON (inv_costing_detail.costing_number = inv_costing.number)" +
				" WHERE inv_costing_detail.item_code='" + ItemCode + "'" +
				" AND inv_costing.date >'" + StockInDateStr + "'" +
				" AND inv_costing_detail.store_code=" + StoreCode + " " +
				"GROUP BY inv_costing_detail.item_code)) AS Stock").Scan(&Stock).Error; err != nil {
			return 0, err
		}
	} else {
		if err := DB.Raw("SELECT" +
			" (SUM(IF(inv_receiving_detail.store_code='" + StoreCode + "', inv_receiving_detail.basic_quantity, 0)) + SUM(IFNULL(StockTransfer.Quantity, 0)) - SUM(IFNULL(Costing.Quantity, 0))) AS Quantity " +
			"FROM" +
			" inv_receiving_detail" +
			" LEFT OUTER JOIN inv_receiving ON (inv_receiving_detail.receive_number = inv_receiving.number)" +
			" LEFT OUTER JOIN (" +
			"SELECT" +
			" SUM(IF(inv_stock_transfer_detail.to_store_code='" + StoreCode + "' AND inv_stock_transfer.date<='" + StockDateStr + "', inv_stock_transfer_detail.quantity, 0) - IF(inv_stock_transfer_detail.from_store_code='" + StoreCode + "', inv_stock_transfer_detail.quantity, 0)) AS Quantity," +
			" inv_stock_transfer_detail.receive_id " +
			"FROM" +
			" inv_stock_transfer_detail" +
			" LEFT OUTER JOIN inv_stock_transfer ON (inv_stock_transfer_detail.st_number = inv_stock_transfer.number)" +
			" LEFT OUTER JOIN inv_receiving_detail ON (inv_stock_transfer_detail.receive_id = inv_receiving_detail.id)" +
			" LEFT OUTER JOIN inv_receiving ON (inv_receiving_detail.receive_number = inv_receiving.number)" +
			" WHERE inv_receiving_detail.quantity>0" +
			" AND inv_receiving.date<='" + StockDateStr + "'" +
			" AND inv_stock_transfer_detail.item_code='" + ItemCode + "'" +
			" AND (inv_stock_transfer_detail.from_store_code='" + StoreCode + "' OR (inv_stock_transfer_detail.to_store_code='" + StoreCode + "' AND inv_stock_transfer.date<='" + StockDateStr + "'))" +
			" AND inv_stock_transfer_detail.st_number<>'" + StockTransferNumber + "' " +
			"GROUP BY inv_stock_transfer_detail.receive_id) AS StockTransfer ON (inv_receiving_detail.id = StockTransfer.receive_id)" +
			" LEFT OUTER JOIN (" +
			"SELECT" +
			" SUM(inv_costing_detail.quantity) AS Quantity," +
			" inv_costing_detail.receive_id " +
			"FROM" +
			" inv_costing_detail" +
			" LEFT OUTER JOIN inv_receiving_detail ON (inv_costing_detail.receive_id = inv_receiving_detail.id)" +
			" LEFT OUTER JOIN inv_receiving ON (inv_receiving_detail.receive_number = inv_receiving.number)" +
			" WHERE inv_receiving_detail.quantity>0" +
			" AND inv_receiving.date<='" + StockDateStr + "'" +
			" AND inv_costing_detail.item_code='" + ItemCode + "'" +
			" AND inv_costing_detail.store_code='" + StoreCode + "' " +
			"GROUP BY inv_costing_detail.receive_id) AS Costing ON (inv_receiving_detail.id = Costing.receive_id)" +
			" WHERE inv_receiving_detail.quantity>0" +
			" AND inv_receiving.date<='" + StockDateStr + "'" +
			" AND inv_receiving_detail.item_code='" + ItemCode + "' " +
			"GROUP BY inv_receiving_detail.item_code").Scan(&Stock).Error; err != nil {
			return 0, err
		}
	}
	return Stock, nil
}

func GetStockStore(DB *gorm.DB, Dataset *global_var.TDataset, StoreCode, ItemCode string, StockDate time.Time) (Stock float64, err error) {
	var StockDateStr, StockInDateStr, LastItemClosedDateStr, IncDay1ItemClosedDateStr string
	var StockInDate, LastItemClosedDate time.Time
	// TODO Optimize query
	// Result = 0
	StockDateStr = General.FormatDate1(StockDate)
	CostingMethod := Dataset.ProgramConfiguration.CostingMethod
	if CostingMethod == GlobalVar.InventoryCostingMethod.Average {
		StockInDate = GetStockInDate(DB, StoreCode, ItemCode, StockDate)
		StockInDateStr = General.FormatDate1(StockInDate)

		LastItemClosedDate = GetLastItemStoreClosedDate(DB, ItemCode)
		LastItemClosedDateStr = General.FormatDate1(LastItemClosedDate)
		IncDay1ItemClosedDateStr = General.FormatDate1(LastItemClosedDate.AddDate(0, 0, 1))

		Query1 := DB.Table(DBVar.TableName.InvReceivingDetail).
			Select("SUM(IF(inv_receiving_detail.store_code=?, inv_receiving_detail.basic_quantity, 0)) AS Quantity").
			Joins("LEFT OUTER JOIN inv_receiving ON (inv_receiving_detail.receive_number = inv_receiving.number)").
			Where("inv_receiving_detail.item_code=?", ItemCode).
			Where("(inv_receiving.date BETWEEN ? AND ?)", IncDay1ItemClosedDateStr, StockDateStr).
			Group("inv_receiving_detail.item_code")
		Query2 := DB.Table(DBVar.TableName.InvStockTransferDetail).
			Select("SUM(inv_stock_transfer_detail.quantity) AS Quantity").
			Joins("LEFT OUTER JOIN inv_stock_transfer ON (inv_stock_transfer_detail.st_number = inv_stock_transfer.number)").
			Where("inv_stock_transfer_detail.item_code=?", ItemCode).
			Where("(inv_stock_transfer.date BETWEEN ? AND ?)", IncDay1ItemClosedDateStr, StockDateStr).
			Where("inv_stock_transfer_detail.to_store_code=?", StoreCode).
			Group("inv_stock_transfer_detail.item_code")
		Query3 := DB.Table(DBVar.TableName.InvStockTransferDetail).
			Select("-SUM(inv_stock_transfer_detail.quantity) AS Quantity").
			Joins("LEFT OUTER JOIN inv_stock_transfer ON (inv_stock_transfer_detail.st_number = inv_stock_transfer.number)").
			Where("inv_stock_transfer_detail.item_code=?", ItemCode).
			Where("(inv_stock_transfer.date BETWEEN ? AND ?)", IncDay1ItemClosedDateStr, StockDateStr).
			Where("inv_stock_transfer_detail.from_store_code=?", StoreCode).
			Group("inv_stock_transfer_detail.item_code")
		Query4 := DB.Table(DBVar.TableName.InvCostingDetail).
			Select("-SUM(inv_costing_detail.quantity) AS Quantity").
			Joins("LEFT OUTER JOIN inv_costing ON (inv_costing_detail.costing_number = inv_costing.number)").
			Where("inv_costing_detail.item_code=?", ItemCode).
			Where("inv_costing.date<=?", StockInDateStr).
			Where("inv_costing_detail.store_code=?", StoreCode).
			Group("inv_costing_detail.item_code")
		Query5 := DB.Table(DBVar.TableName.InvCloseSummaryStore).
			Select("inv_close_summary_store.quantity_all AS Quantity").
			Where("inv_close_summary_store.date=?", LastItemClosedDateStr).
			Where("inv_close_summary_store.store_code=?", StoreCode).
			Where("inv_close_summary_store.item_code=?", ItemCode)

		var DataOutput []map[string]interface{}
		QueryUnion := DB.Raw("(?) UNION ALL (?) UNION ALL (?) UNION ALL (?) UNION ALL (?)", Query1, Query2, Query3, Query4, Query5).Scan(&DataOutput)
		if err := DB.Table("(?) as Stock", QueryUnion).Select("SUM(IFNULL(Stock.Quantity, 0)").Scan(&Stock).Error; err != nil {
			return 0, err
		}
	} else {
		JoinQuery1 := DB.Raw(
			"SELECT"+
				" SUM(IF(inv_stock_transfer_detail.to_store_code=? AND inv_stock_transfer.date<=?, inv_stock_transfer_detail.quantity, 0) - IF(inv_stock_transfer_detail.from_store_code=?, inv_stock_transfer_detail.quantity, 0)) AS Quantity,"+
				" inv_stock_transfer_detail.receive_id "+
				"FROM"+
				" inv_stock_transfer_detail"+
				" LEFT OUTER JOIN inv_stock_transfer ON (inv_stock_transfer_detail.st_number = inv_stock_transfer.number)"+
				" LEFT OUTER JOIN inv_receiving_detail ON (inv_stock_transfer_detail.receive_id = inv_receiving_detail.id)"+
				" LEFT OUTER JOIN inv_receiving ON (inv_receiving_detail.receive_number = inv_receiving.number)"+
				" WHERE inv_receiving_detail.quantity>0"+
				" AND inv_receiving.date<=?"+
				" AND inv_stock_transfer_detail.item_code=?"+
				" AND (inv_stock_transfer_detail.from_store_code=? OR (inv_stock_transfer_detail.to_store_code=? AND inv_stock_transfer.date<=?)) "+
				" GROUP BY inv_stock_transfer_detail.receive_id", StoreCode, StockDateStr, StoreCode, StockDateStr, ItemCode, StoreCode, StoreCode, StockDateStr)
		JoinQuery2 := DB.Raw("SELECT"+
			" SUM(inv_costing_detail.quantity) AS Quantity,"+
			" inv_costing_detail.receive_id "+
			"FROM"+
			" inv_costing_detail"+
			" LEFT OUTER JOIN inv_receiving_detail ON (inv_costing_detail.receive_id = inv_receiving_detail.id)"+
			" LEFT OUTER JOIN inv_receiving ON (inv_receiving_detail.receive_number = inv_receiving.number)"+
			" WHERE inv_receiving_detail.quantity>0"+
			" AND inv_receiving.date<=?"+
			" AND inv_costing_detail.item_code=?"+
			" AND inv_costing_detail.store_code=?"+
			" GROUP BY inv_costing_detail.receive_id", StockDateStr, ItemCode, StoreCode)

		if err := DB.Debug().Table(DBVar.TableName.InvReceivingDetail).
			Select("(SUM(IF(inv_receiving_detail.store_code='"+StoreCode+"', inv_receiving_detail.basic_quantity, 0)) + SUM(IFNULL(StockTransfer.Quantity, 0)) - SUM(IFNULL(Costing.Quantity, 0))) AS Quantity").
			Joins("LEFT OUTER JOIN inv_receiving ON (inv_receiving_detail.receive_number = inv_receiving.number)").
			Joins("LEFT OUTER JOIN (?) AS StockTransfer ON (inv_receiving_detail. id = StockTransfer.receive_id)", JoinQuery1).
			Joins("LEFT OUTER JOIN (?) AS Costing ON (inv_receiving_detail. id = Costing.receive_id)", JoinQuery2).
			Where("inv_receiving_detail.quantity>0").
			Where("inv_receiving.date<=?", StockDateStr).
			Where("inv_receiving_detail.item_code=?", ItemCode).
			Group("inv_receiving_detail.item_code").
			Scan(&Stock).Error; err != nil {
			return 0, err
		}
	}
	return Stock, nil
}

func GetStockStoreUpdateCosting(DB *gorm.DB, Dataset *global_var.TDataset, StoreCode, ItemCode, CostingNumber string, StockDate time.Time) (Stock float64, err error) {
	var StockDateStr, StockInDateStr, LastItemClosedDateStr, IncDay1ItemClosedDateStr string
	var StockInDate, LastItemClosedDate time.Time
	// TODO Optimize query
	// Result = 0
	StockDateStr = General.FormatDate1(StockDate)
	CostingMethod := Dataset.Configuration[GlobalVar.ConfigurationCategory.Inventory][GlobalVar.ConfigurationName.CostingMethod].(string)
	if CostingMethod == GlobalVar.InventoryCostingMethod.Average {
		StockInDate = GetStockInDate(DB, StoreCode, ItemCode, StockDate)
		StockInDateStr = General.FormatDate1(StockInDate)

		LastItemClosedDate = GetLastItemStoreClosedDate(DB, ItemCode)
		LastItemClosedDateStr = General.FormatDate1(LastItemClosedDate)
		IncDay1ItemClosedDateStr = General.FormatDate1(LastItemClosedDate.AddDate(0, 0, 1))

		Query1 := DB.Table(DBVar.TableName.InvReceivingDetail).
			Select("SUM(IF(inv_receiving_detail.store_code=?, inv_receiving_detail.basic_quantity, 0)) AS Quantity").
			Joins("LEFT OUTER JOIN inv_receiving ON (inv_receiving_detail.receive_number = inv_receiving.number)").
			Where("inv_receiving_detail.item_code=?", ItemCode).
			Where("(inv_receiving.date BETWEEN ? AND ?)", IncDay1ItemClosedDateStr, StockDateStr).
			Group("inv_receiving_detail.item_code")
		Query2 := DB.Table(DBVar.TableName.InvStockTransferDetail).
			Select("SUM(inv_stock_transfer_detail.quantity) AS Quantity ").
			Joins("LEFT OUTER JOIN inv_stock_transfer ON (inv_stock_transfer_detail.st_number = inv_stock_transfer.number)").
			Where("inv_stock_transfer_detail.item_code=?", ItemCode).
			Where("(inv_stock_transfer.date BETWEEN ? AND ?)", IncDay1ItemClosedDateStr, StockDateStr).
			Where("inv_stock_transfer_detail.to_store_code=?").
			Group("inv_stock_transfer_detail.item_code")
		Query3 := DB.Table(DBVar.TableName.InvStockTransferDetail).
			Select("-SUM(inv_stock_transfer_detail.quantity) AS Quantity").
			Joins("LEFT OUTER JOIN inv_stock_transfer ON (inv_stock_transfer_detail.st_number = inv_stock_transfer.number)").
			Where("inv_stock_transfer_detail.item_code=?", ItemCode).
			Where("(inv_stock_transfer.date BETWEEN ? AND ?)", IncDay1ItemClosedDateStr, StockDateStr).
			Where("inv_stock_transfer_detail.from_store_code=?", StoreCode).
			Group("inv_stock_transfer_detail.item_code")
		Query4 := DB.Table(DBVar.TableName.InvCostingDetail).
			Select("-SUM(inv_costing_detail.quantity) AS Quantity").
			Joins("LEFT OUTER JOIN inv_costing ON (inv_costing_detail.costing_number = inv_costing.number)").
			Where("inv_costing_detail.item_code=?", ItemCode).
			Where("inv_costing.date<=?", StockInDateStr).
			Where("inv_costing_detail.store_code=?", StoreCode).
			Where("inv_costing_detail.costing_number<>?", CostingNumber).
			Group("inv_costing_detail.item_code")
		Query5 := DB.Table(DBVar.TableName.InvCloseSummaryStore).
			Select("inv_close_summary_store.quantity_all AS Quantity").
			Where("inv_close_summary_store.date=?", LastItemClosedDateStr).
			Where("inv_close_summary_store.store_code=?", StoreCode).
			Where("inv_close_summary_store.item_code=?", ItemCode)

		var DataOutput []map[string]interface{}
		QueryUnion := DB.Raw("(?) UNION ALL (?) UNION ALL (?) UNION ALL (?) UNION ALL (?)", Query1, Query2, Query3, Query4, Query5).Scan(&DataOutput)
		if err := DB.Table("(?) as Stock", QueryUnion).Select("SUM(IFNULL(Stock.Quantity, 0)").Scan(&Stock).Error; err != nil {
			return 0, err
		}
	} else {
		JoinQuery1 := DB.Raw(
			"SELECT"+
				" SUM(IF(inv_stock_transfer_detail.to_store_code=? AND inv_stock_transfer.date<=?, inv_stock_transfer_detail.quantity, 0) - IF(inv_stock_transfer_detail.from_store_code=?, inv_stock_transfer_detail.quantity, 0)) AS Quantity,"+
				" inv_stock_transfer_detail.receive_id "+
				"FROM"+
				" inv_stock_transfer_detail"+
				" LEFT OUTER JOIN inv_stock_transfer ON (inv_stock_transfer_detail.st_number = inv_stock_transfer.number)"+
				" LEFT OUTER JOIN inv_receiving_detail ON (inv_stock_transfer_detail.receive_id = inv_receiving_detail. id)"+
				" LEFT OUTER JOIN inv_receiving ON (inv_receiving_detail.receive_number = inv_receiving.number)"+
				" LEFT OUTER JOIN ("+
				"SELECT DISTINCT"+
				" inv_costing_detail.receive_id "+
				"FROM"+
				" inv_costing_detail"+
				" WHERE inv_costing_detail.costing_number=?"+
				" AND inv_costing_detail.item_code=? "+
				"GROUP BY inv_costing_detail.receive_id) AS CostingItem ON (inv_receiving_detail. id = CostingItem.receive_id)"+
				" WHERE (inv_receiving_detail.quantity>0 OR (inv_receiving_detail.quantity=0 AND IFNULL(CostingItem.receive_id, '')<>''))"+
				" AND inv_receiving.date<=?"+
				" AND inv_stock_transfer_detail.item_code=?"+
				" AND (inv_stock_transfer_detail.from_store_code=? OR (inv_stock_transfer_detail.to_store_code=? AND inv_stock_transfer.date<=?)) "+
				"GROUP BY inv_stock_transfer_detail.receive_id", StoreCode, StockDateStr, StoreCode, CostingNumber, ItemCode, StockDateStr, ItemCode, StoreCode, StoreCode, StockDateStr)
		JoinQuery2 := DB.Raw("SELECT"+
			" SUM(inv_costing_detail.quantity) AS Quantity,"+
			" inv_costing_detail.receive_id "+
			"FROM"+
			" inv_costing_detail"+
			" LEFT OUTER JOIN inv_receiving_detail ON (inv_costing_detail.receive_id = inv_receiving_detail. id)"+
			" LEFT OUTER JOIN inv_receiving ON (inv_receiving_detail.receive_number = inv_receiving.number)"+
			" LEFT OUTER JOIN ("+
			"SELECT DISTINCT"+
			" inv_costing_detail.receive_id "+
			"FROM"+
			" inv_costing_detail"+
			" WHERE inv_costing_detail.costing_number=?"+
			" AND inv_costing_detail.item_code=? "+
			"GROUP BY inv_costing_detail.receive_id) AS CostingItem ON (inv_receiving_detail. id = CostingItem.receive_id)"+
			" WHERE (inv_receiving_detail.quantity>0 OR (inv_receiving_detail.quantity=0 AND IFNULL(CostingItem.receive_id, '')<>''))"+
			" AND inv_receiving.date<=?"+
			" AND inv_costing_detail.item_code=?"+
			" AND inv_costing_detail.store_code=?"+
			" AND inv_costing_detail.costing_number<>? "+
			"GROUP BY inv_costing_detail.receive_id", CostingNumber, ItemCode, StockDateStr, ItemCode, StoreCode, CostingNumber)
		JoinQuery3 := DB.Raw("SELECT DISTINCT"+
			" inv_costing_detail.receive_id "+
			"FROM"+
			" inv_costing_detail"+
			" WHERE inv_costing_detail.costing_number=?"+
			" AND inv_costing_detail.item_code=? "+
			"GROUP BY inv_costing_detail.receive_id", CostingNumber, ItemCode)

		if err := DB.Table(DBVar.TableName.InvReceivingDetail).
			Select("(SUM(IF(inv_receiving_detail.store_code='"+StoreCode+"', inv_receiving_detail.basic_quantity, 0)) + SUM(IFNULL(StockTransfer.Quantity, 0)) - SUM(IFNULL(Costing.Quantity, 0))) AS Quantity").
			Joins("LEFT OUTER JOIN inv_receiving ON (inv_receiving_detail.receive_number = inv_receiving.number)").
			Joins("LEFT OUTER JOIN (?) AS StockTransfer ON (inv_receiving_detail. id = StockTransfer.receive_id)", JoinQuery1).
			Joins("LEFT OUTER JOIN (?) AS Costing ON (inv_receiving_detail. id = Costing.receive_id)", JoinQuery2).
			Joins("LEFT OUTER JOIN (?) AS CostingItem ON (inv_receiving_detail. id = CostingItem.receive_id)", JoinQuery3).
			Where("(inv_receiving_detail.quantity>0 OR (inv_receiving_detail.quantity=0 AND IFNULL(CostingItem.receive_id, '')<>''))").
			Where("inv_receiving_detail.item_code=?", ItemCode).
			Where("inv_receiving.date<=?", StockDateStr).
			Group("inv_receiving_detail.item_code").
			Scan(&Stock).Error; err != nil {
			return 0, err
		}
	}
	return Stock, nil
}

func GetLastItemStoreClosedDate(DB *gorm.DB, ItemCode string) (CloseDate time.Time) {
	DB.Table(DBVar.TableName.InvCloseSummaryStore).Select("date").Where("item_code=?", ItemCode).Limit(1).Scan(&CloseDate)
	return
}

func GetLastInventoryClosedDate(DB *gorm.DB) (CloseDate time.Time) {
	DB.Table(DBVar.TableName.InvCloseLog).Select("date").Limit(1).Scan(&CloseDate)
	return
}

func GetStockInDate(DB *gorm.DB, StoreCode, ItemCode string, StockDate time.Time) time.Time {
	var StockDateStr string
	var ReceiveDate, TransferInDate, CostingDate time.Time
	var ReceiveDateQuery, TransferInDateQuery, CostingDateQuery time.Time
	ReceiveDate = StockDate
	TransferInDate = StockDate
	CostingDate = StockDate
	StockDateStr = General.FormatDate1(StockDate)

	DB.Table(DBVar.TableName.InvReceivingDetail).Select("inv_receiving.`date`").
		Joins("LEFT OUTER JOIN inv_receiving ON (inv_receiving_detail.receive_number = inv_receiving.number)").
		Where("inv_receiving.date>?", StockDateStr).
		Where("inv_receiving_detail.item_code=?", ItemCode).
		Where("inv_receiving_detail.store_code=?", StoreCode).
		Order("inv_receiving.date").
		Limit(1).
		Scan(&ReceiveDateQuery)

	if ReceiveDateQuery.IsZero() {
		ReceiveDate = ReceiveDateQuery
	}

	DB.Table(DBVar.TableName.InvStockTransferDetail).Select("inv_stock_transfer.`date`").
		Joins("LEFT OUTER JOIN inv_stock_transfer ON (inv_stock_transfer_detail.st_number = inv_stock_transfer.number)").
		Where("inv_stock_transfer.date>?", StockDateStr).
		Where("inv_stock_transfer_detail.item_code=?", ItemCode).
		Where("inv_stock_transfer_detail.to_store_code=?", StoreCode).
		Order("inv_stock_transfer.date").
		Limit(1).
		Scan(&TransferInDateQuery)

	if TransferInDateQuery.IsZero() {
		TransferInDate = TransferInDateQuery
	}

	if ReceiveDate.Unix() > StockDate.Unix() {
		if TransferInDate.Unix() > StockDate.Unix() {
			if ReceiveDate.Unix() < TransferInDate.Unix() {
				return ReceiveDate
			} else {
				return TransferInDate
			}
		} else {
			return ReceiveDate
		}
	} else if TransferInDate.Unix() > StockDate.Unix() {
		return TransferInDate
	} else {

		DB.Table(DBVar.TableName.InvCostingDetail).Select("inv_costing.`date`").
			Joins("LEFT OUTER JOIN inv_costing ON (inv_costing_detail.costing_number = inv_costing.number)").
			Where("inv_costing.date>?", StockDateStr).
			Where("inv_costing_detail.item_code=?", ItemCode).
			Where("inv_costing_detail.store_code=?", StoreCode).
			Order("inv_costing.date").
			Limit(1).
			Scan(&CostingDateQuery)

		if CostingDateQuery.IsZero() {
			CostingDate = CostingDateQuery
		}

		if CostingDate.Unix() > StockDate.Unix() {
			return CostingDate
		} else {
			return StockDate
		}
	}
}

func GetInventoryJournalAccount(DB *gorm.DB, ItemCode string) (JournalAccountCode string) {
	DB.Table(DBVar.TableName.InvCfgInitItem).Select("inv_cfg_init_item_category.journal_account_code").
		Joins("LEFT OUTER JOIN inv_cfg_init_item_category ON (inv_cfg_init_item.category_code = inv_cfg_init_item_category.code)").
		Where("inv_cfg_init_item.code=?", ItemCode).
		Limit(1).
		Scan(&JournalAccountCode)

	return
}

func GetStoreID(DB *gorm.DB, Code string) (ID uint64) {
	DB.Table(DBVar.TableName.InvCfgInitStore).Select("id").
		Where("code=?", Code).
		Limit(1).
		Scan(&ID)

	return
}

func IsStockMinusStockTransfer(DB *gorm.DB, StoreCode, ItemCode, StockTransferNumber, CostingMethod string, StockDate time.Time) (IsMinus bool) {
	Query1 := DB.Table(DBVar.TableName.InvCostingDetail).Distinct("inv_costing.`date").
		Joins("LEFT OUTER JOIN inv_costing ON (inv_costing_detail.costing_number = inv_costing.number)").
		Where("inv_costing_detail.item_code=?", ItemCode).
		Where("inv_costing_detail.store_code=?", StoreCode).
		Where("inv_costing.date>=?", StockDate)
	Query2 := DB.Table(DBVar.TableName.InvStockTransferDetail).Distinct("inv_stock_transfer.`date").
		Joins("LEFT OUTER JOIN inv_stock_transfer ON (inv_stock_transfer_detail.st_number = inv_stock_transfer.number)").
		Where("inv_stock_transfer_detail.item_code=?", ItemCode).
		Where("inv_stock_transfer_detail.from_store_code=?", StoreCode).
		Where("inv_stock_transfer.`date` >= ?", StockDate).
		Where("inv_stock_transfer_detail.st_number<>?", StockTransferNumber)

	var DataOutput []map[string]interface{}
	DB.Table("(?) AS Stock", DB.Raw("(?) UNION ALL (?)", Query1, Query2)).Distinct("Stock.`date`").Order("Stock.`date`").Scan(&DataOutput)

	if len(DataOutput) > 0 {
		for _, data := range DataOutput {
			if GetStockStoreUpdateReceive(DB, StoreCode, ItemCode, StockTransferNumber, CostingMethod, data["date"].(time.Time)) < 0 {
				IsMinus = true
				break
			}
		}
	}
	return
}

func GetStockStoreUpdateReceive(DB *gorm.DB, StoreCode, ItemCode, ReceiveNumber, CostingMethod string, StockDate time.Time) (Stock float64) {
	var StockDateStr, StockInDateStr string
	var StockInDate time.Time
	// TODO Optimize query
	StockDateStr = General.FormatDate1(StockDate)
	if CostingMethod == GlobalVar.InventoryCostingMethod.Average {
		StockInDate = GetStockInDate(DB, StoreCode, ItemCode, StockDate)
		StockInDateStr = General.FormatDate1(StockInDate)

		Query1 := DB.Table(DBVar.TableName.InvReceivingDetail).
			Select("SUM(IF(inv_receiving_detail.store_code=?, inv_receiving_detail.basic_quantity, 0)) AS Quantity").
			Joins("LEFT OUTER JOIN inv_receiving ON (inv_receiving_detail.receive_number = inv_receiving.number)").
			Where("inv_receiving_detail.item_code=?", ItemCode).
			Where("(inv_receiving.date<=?", StockDateStr).
			Where("inv_receiving_detail.receive_number<>?", ReceiveNumber).
			Group("inv_receiving_detail.item_code")
		Query2 := DB.Table(DBVar.TableName.InvStockTransferDetail).
			Select("SUM(inv_stock_transfer_detail.quantity) AS Quantity").
			Joins("LEFT OUTER JOIN inv_stock_transfer ON (inv_stock_transfer_detail.st_number = inv_stock_transfer.number)").
			Where("inv_stock_transfer_detail.item_code=?", ItemCode).
			Where("inv_stock_transfer.date<=?", StockDateStr).
			Where("inv_stock_transfer_detail.to_store_code=?", StoreCode).
			Where("inv_stock_transfer.number<>?", ReceiveNumber).
			Group("inv_stock_transfer_detail.item_code")
		Query3 := DB.Table(DBVar.TableName.InvStockTransferDetail).
			Select("-SUM(inv_stock_transfer_detail.quantity) AS Quantity").
			Joins("LEFT OUTER JOIN inv_stock_transfer ON (inv_stock_transfer_detail.st_number = inv_stock_transfer.number)").
			Where("inv_stock_transfer_detail.item_code=?", ItemCode).
			Where("inv_stock_transfer.date<=?", StockDateStr).
			Where("inv_stock_transfer_detail.from_store_code=?", StoreCode).
			Group("inv_stock_transfer_detail.item_code")
		Query4 := DB.Table(DBVar.TableName.InvCostingDetail).
			Select("-SUM(inv_costing_detail.quantity) AS Quantity").
			Joins("LEFT OUTER JOIN inv_costing ON (inv_costing_detail.costing_number = inv_costing.number)").
			Where("inv_costing_detail.item_code=?", ItemCode).
			Where("inv_costing.date<=?", StockInDateStr).
			Where("inv_costing_detail.store_code=?", StoreCode).
			Group("inv_costing_detail.item_code")

		var DataOutput []map[string]interface{}
		QueryUnion := DB.Raw("(?) UNION ALL (?) UNION ALL (?) UNION ALL (?) ", Query1, Query2, Query3, Query4).Scan(&DataOutput)
		DB.Table("(?) as Stock", QueryUnion).Select("SUM(IFNULL(Stock.Quantity, 0)").Scan(&Stock)
	}
	return
}

func GetStockUsedOnDestination(DB *gorm.DB, StockTransferNumber, ItemCode string) (Stock float64) {
	if ItemCode != "" {
		type DataOutputStruct struct {
			ToStoreCode, ItemCode string
			Quantity              float64
			ReceiveID             string
			Date                  time.Time
		}
		var DataOutput []DataOutputStruct
		DetailTableName := DBVar.TableName.InvStockTransferDetail
		DB.Table(DetailTableName).Select(
			DetailTableName+".to_store_code,"+
				DetailTableName+".item_code,"+
				DetailTableName+".receive_id,"+
				"SUM("+DetailTableName+".quantity) AS Quantity,"+
				"inv_stock_transfer.`date` ").
			Joins("LEFT OUTER JOIN inv_stock_transfer ON ("+DetailTableName+".st_number = inv_stock_transfer.number)").
			Where(DetailTableName+".st_number=?", StockTransferNumber).
			Where(DetailTableName+".item_code=?", ItemCode).
			Group(DetailTableName + ".receive_id").Scan(&DataOutput)

		if len(DataOutput) > 0 {
			for _, data := range DataOutput {
				Stock = Stock + data.Quantity - GetStockStoreDestination(DB, data.ToStoreCode, data.ReceiveID, data.Date)
			}
		}
	}
	return
}

func GetStockStoreDestination(DB *gorm.DB, StoreCode string, ReceiveID string, StockDate time.Time) (Stock float64) {
	var StockDateStr string

	StockDateStr = General.FormatDate1(StockDate)
	DB.Table(DBVar.TableName.InvReceivingDetail).Select("(SUM(IF(inv_receiving_detail.store_code='"+StoreCode+"', inv_receiving_detail.basic_quantity, 0)) + SUM(IFNULL(StockTransfer.Quantity, 0)) - SUM(IFNULL(Costing.Quantity, 0))) AS Quantity").
		Joins("LEFT OUTER JOIN inv_receiving ON (inv_receiving_detail.receive_number = inv_receiving.number)").
		Joins("LEFT OUTER JOIN ("+
			"SELECT"+
			" SUM(IF(inv_stock_transfer_detail.to_store_code=? AND inv_stock_transfer.date<='"+StockDateStr+"', inv_stock_transfer_detail.quantity, 0) - IF(inv_stock_transfer_detail.from_store_code=?, inv_stock_transfer_detail.quantity, 0)) AS Quantity,"+
			" inv_stock_transfer_detail.receive_id "+
			"FROM"+
			" inv_stock_transfer_detail"+
			" LEFT OUTER JOIN inv_stock_transfer ON (inv_stock_transfer_detail.st_number = inv_stock_transfer.number)"+
			" WHERE inv_stock_transfer_detail.receive_id="+ReceiveID+
			" AND (inv_stock_transfer_detail.from_store_code=? OR (inv_stock_transfer_detail.to_store_code=? AND inv_stock_transfer.date<='"+StockDateStr+"')) "+
			" GROUP BY inv_stock_transfer_detail.receive_id) AS StockTransfer ON (inv_receiving_detail.id = StockTransfer.receive_id)", StoreCode, StoreCode, StoreCode, StoreCode).
		Joins("LEFT OUTER JOIN ("+
			"SELECT"+
			" SUM(inv_costing_detail.quantity) AS Quantity,"+
			" inv_costing_detail.receive_id "+
			"FROM"+
			" inv_costing_detail"+
			" WHERE inv_costing_detail.receive_id="+ReceiveID+
			" AND inv_costing_detail.store_code=?"+
			" GROUP BY inv_costing_detail.receive_id) AS Costing ON (inv_receiving_detail.id = Costing.receive_id)", StoreCode).
		Where("inv_receiving_detail.id=?", ReceiveID).
		Where("inv_receiving.date<=?", StockDateStr).
		Group("inv_receiving_detail.item_code").
		Scan(&Stock)

	return
}

func IsReservationLocked(DB *gorm.DB, ReservationNumber interface{}) (IsLocked bool) {
	IsLocked = MasterData.GetFieldString(DB, DBVar.TableName.Reservation, "is_lock", "number", ReservationNumber, "") == "1"
	return
}

func IsThereCardActiveReservation(DB *gorm.DB, ReservationNumber interface{}) (IsThereCardActive bool) {
	IsThereCardActive = MasterData.GetFieldString(DB, DBVar.TableName.LogKeylock, "id", "is_active='1' AND reservation_number", ReservationNumber, "") != ""
	return
}

func IsGroupAlreadyUsed(DB *gorm.DB, GroupCode string) bool {
	IsUsedInReservation := MasterData.GetFieldBool(DB, DBVar.TableName.Reservation, "number", "group_code", GroupCode, false)
	IsUsedInFolio := MasterData.GetFieldBool(DB, DBVar.TableName.Folio, "number", "group_code", GroupCode, false)

	return IsUsedInReservation || IsUsedInFolio

}

func IsRoomOccupiedNow(c *gin.Context, DB *gorm.DB, RoomNumber string) bool {
	if MasterData.GetConfigurationBool(DB, GlobalVar.SystemCode.General, GlobalVar.ConfigurationCategory.General, GlobalVar.ConfigurationName.IsRoomByName, false) {
		RoomNumber = GetRoomNumberFromConfigurationIsRoomName(c, DB, RoomNumber)
	}

	var Number uint64
	DB.Table(DBVar.TableName.Folio).Select("folio.number").
		Joins("LEFT OUTER JOIN guest_detail ON (folio.guest_detail_id = guest_detail.id)").
		Where("guest_detail.room_number=?", RoomNumber).
		Where("folio.status_code=?", GlobalVar.FolioStatus.Open).
		Limit(1).
		Scan(&Number)

	return Number > 0
}

func GetRoomBlockStatus(c *gin.Context, DB *gorm.DB, RoomNumber string) string {
	if MasterData.GetConfigurationBool(DB, GlobalVar.SystemCode.General, GlobalVar.ConfigurationCategory.General, GlobalVar.ConfigurationName.IsRoomByName, false) {
		RoomNumber = GetRoomNumberFromConfigurationIsRoomName(c, DB, RoomNumber)
	}
	var Status sql.NullString
	DB.Table(DBVar.TableName.CfgInitRoom).Select("block_status_code").
		Where("number=?", RoomNumber).
		Limit(1).
		Scan(&Status)
	return Status.String
}

func IsRoomBlockedNow(c *gin.Context, DB *gorm.DB, RoomNumber string) bool {
	if MasterData.GetConfigurationBool(DB, GlobalVar.SystemCode.General, GlobalVar.ConfigurationCategory.General, GlobalVar.ConfigurationName.IsRoomByName, false) {
		RoomNumber = GetRoomNumberFromConfigurationIsRoomName(c, DB, RoomNumber)
	}

	return GetRoomBlockStatus(c, DB, RoomNumber) != ""
}

func IsCheckIn(DB *gorm.DB, ReservationNumber uint64) bool {
	var Number uint64
	DB.Table(DBVar.TableName.Folio).Select("reservation_number").
		Where("reservation_number=?", ReservationNumber).
		Limit(1).
		Scan(&Number)

	return Number > 0
}

func IsRoomReady(c *gin.Context, DB *gorm.DB, RoomNumber string) bool {
	if MasterData.GetConfigurationBool(DB, GlobalVar.SystemCode.General, GlobalVar.ConfigurationCategory.General, GlobalVar.ConfigurationName.IsRoomByName, false) {
		RoomNumber = GetRoomNumberFromConfigurationIsRoomName(c, DB, RoomNumber)
	}
	var Status string
	DB.Table(DBVar.TableName.CfgInitRoom).Select("status_code").
		Where("number=?", RoomNumber).
		Where("status_code=?", GlobalVar.RoomStatus.Ready).
		Limit(1).
		Scan(&Status)
	return Status != ""
}

func GetMemberCodeFromGuestProfile(DB *gorm.DB, GuestProfileID uint64) (Code string) {
	DB.Table(DBVar.TableName.Member).Select("code").
		Where("guest_profile_id=?", GuestProfileID).
		Limit(1).
		Scan(&Code)
	return
}

func GetMemberIDFromGuestProfile(DB *gorm.DB, MemberCode string) (ID uint64) {
	DB.Table(DBVar.TableName.Member).Select("guest_profile_id").
		Where("code=?", MemberCode).
		Limit(1).
		Scan(&ID)
	return
}

func IsScheduledRate(ctx context.Context, DB *gorm.DB, FolioNumber uint64, ADate time.Time) bool {
	var ID uint64
	DB.WithContext(ctx).Table(DBVar.TableName.GuestScheduledRate).Select("id").
		Where("folio_number=?", FolioNumber).
		Where("from_date<=?", General.FormatDate1(ADate)).
		Where("to_date>=?", General.FormatDate1(ADate)).
		Limit(1).
		Scan(&ID)
	return ID > 0
}

func GetFolioComplimentHU(DB *gorm.DB, FolioNumber uint64) (ComplimentHU string) {
	DB.Table(DBVar.TableName.Folio).Select("compliment_hu").Where("number=?", FolioNumber).Take(&ComplimentHU)
	return
}

func GetScheduledRateComplimentHU(ctx context.Context, DB *gorm.DB, FolioNumber uint64, PostingDate time.Time) (ComplimentHU string) {
	ctx, span := global_var.Tracer.Start(ctx, "GetScheduledRateComplimentHU")
	defer span.End()

	DB.Table(DBVar.TableName.GuestScheduledRate).Select("compliment_hu").
		Where("folio_number=?", FolioNumber).
		Where("from_date<=?", General.FormatDate1(PostingDate)).
		Where("to_date>=?", General.FormatDate1(PostingDate)).
		Limit(1).
		Scan(&ComplimentHU)
	return
}

func PostCharges(ctx context.Context, c *gin.Context, DB *gorm.DB, Dataset *global_var.TDataset, FolioNumber uint64, SubFolioGroupCode string, AllowZeroAmount bool, PostingDate time.Time, UserID string) error {
	ctx, span := global_var.Tracer.Start(ctx, "PostCharges")
	defer span.End()

	if _, err := PostingRoomCharge(ctx, c, DB, Dataset, FolioNumber, SubFolioGroupCode, AllowZeroAmount, PostingDate, UserID); err != nil {
		return err
	}
	if _, err := PostingExtraCharge(ctx, c, DB, Dataset, FolioNumber, 0, SubFolioGroupCode, UserID); err != nil {
		return err
	}
	return nil
}

func IsVoucherComplimentStillActive(ctx context.Context, DB *gorm.DB, VoucherNumber string, ArrivalDate, AuditDate time.Time) bool {
	Result := General.IncDay(General.DateOf(ArrivalDate), GetVoucherNights(ctx, DB, VoucherNumber)).Unix() > AuditDate.Unix()
	return Result
}

func GetVoucherNights(ctx context.Context, DB *gorm.DB, VoucherNumber string) (Night int) {
	DB.WithContext(ctx).Table(DBVar.TableName.Voucher).Select("nights").Where("number=?", VoucherNumber).Limit(1).Scan(&Night)
	return
}

func GetVoucherType(ctx context.Context, DB *gorm.DB, VoucherNumber string) (Type string) {
	DB.WithContext(ctx).Table(DBVar.TableName.Voucher).Select("type_code").Where("number=?", VoucherNumber).Limit(1).Scan(&Type)
	return
}

func GetScheduledRoomRateCode(DB *gorm.DB, FolioNumber uint64, PostingDate time.Time) (RoomRateCode string) {
	DB.Table(DBVar.TableName.GuestScheduledRate).Select("room_rate_code").
		Where("folio_number=?", FolioNumber).
		Where("from_date<=?", General.FormatDate1(PostingDate)).
		Where("to_date>=?", General.FormatDate1(PostingDate)).
		Limit(1).Scan(&RoomRateCode)
	return
}

func GetRoomRateTaxAndServiceCode(DB *gorm.DB, RoomRateCode string) (Code string) {
	DB.Table(DBVar.TableName.CfgInitRoomRate).Select("tax_and_service_code").Where("code=?", RoomRateCode).Limit(1).Scan(&Code)
	return
}

func GetScheduledRate(DB *gorm.DB, FolioNumber uint64, PostingDate time.Time) (RoomRate float64) {
	DB.Debug().Table(DBVar.TableName.GuestScheduledRate).Select("rate").
		Where("folio_number=?", FolioNumber).
		Where("from_date<=?", General.FormatDate1(PostingDate)).
		Where("to_date>=?", General.FormatDate1(PostingDate)).
		Limit(1).Scan(&RoomRate)
	return
}

func GetVoucherPrice(DB *gorm.DB, VoucherNumber string) (Price float64) {
	DB.Table(DBVar.TableName.Voucher).Select("price").Where("number=?", VoucherNumber).Limit(1).Scan(&Price)
	return
}

func PostingRoomCharge(ctx context.Context, c *gin.Context, DB *gorm.DB, Dataset *global_var.TDataset, FolioNumber uint64, SubFolioGroupCode string, AllowZeroAmount bool, PostingDate time.Time, UserID string) (Result int, err error) {
	ctx, span := global_var.Tracer.Start(ctx, "PostingRoomCharge")
	defer span.End()

	var RoomRateAmountOriginal, RoomRateAmount, RoomChargeB4Breakdown, RoomChargeAfterBreakdown, RoomChargeBasic, RoomChargeTax, RoomChargeService,
		Discount, TotalBreakdown, Commission, BreakdownAmount, BreakdownBasic, BreakdownTax, BreakdownService, ExchangeRate float64
	var RoomNumber, RoomRateCode, RoomRateTaxServiceCode, BusinessSourceCode, CurrencyCode, ComplimentHU, VoucherTypeCode string
	var IsScheduledRateX, IsVoucherActiveX, IsBreakfastX bool
	var BreakDown1, CorrectionBreakdown uint64

	Result = 255
	//0: Succes
	//1: Breakdown Too Large
	//2: Zero Amount Not Allowed
	//3: No Room Charge for Today
	//4: No Room Charge for Compliment
	//255: Posting Room Charge Failed

	//Proses Posting Room Charge and Breakdown
	type FolioStruct struct {
		Folio         DBVar.Folio              `gorm:"embedded"`
		ContactPerson DBVar.Contact_person     `gorm:"embedded"`
		GuestGeneral  DBVar.Guest_general      `gorm:"embedded"`
		GuestDetail   DBVar.Guest_detail       `gorm:"embedded"`
		RoomRate      DBVar.Cfg_init_room_rate `gorm:"embedded"`
		DateArrival   time.Time
	}
	var DataOutput FolioStruct
	err = DB.Table(DBVar.TableName.Folio).Select(
		" folio.group_code,"+
			" guest_detail.currency_code,"+
			" guest_detail.exchange_rate,"+
			" guest_general.purpose_of_code,"+
			" guest_general.sales_code,"+
			" folio.voucher_number,"+
			" guest_general.notes,"+
			" folio.compliment_hu,"+
			" DATE(guest_detail.arrival) AS DateArrival,"+
			" contact_person.title_code,"+
			" contact_person.full_name,"+
			" contact_person.street,"+
			" contact_person.city_code,"+
			" contact_person.city,"+
			" contact_person.nationality_code,"+
			" contact_person.country_code,"+
			" contact_person.state_code,"+
			" contact_person.postal_code,"+
			" contact_person.phone1,"+
			" contact_person.phone2,"+
			" contact_person.fax,"+
			" contact_person.email,"+
			" contact_person.website,"+
			" contact_person.company_code,"+
			" contact_person.guest_type_code,"+
			" contact_person.custom_lookup_field_code01,"+
			" contact_person.custom_lookup_field_code02,"+
			" guest_detail.adult,"+
			" guest_detail.child,"+
			" guest_detail.room_type_code,"+
			" guest_detail.bed_type_code,"+
			" guest_detail.room_number,"+
			" guest_detail.room_rate_code,"+
			" guest_detail.weekday_rate,"+
			" guest_detail.weekend_rate,"+
			" guest_detail.discount_percent,"+
			" guest_detail.discount,"+
			" guest_detail.business_source_code,"+
			" guest_detail.commission_type_code,"+
			" guest_detail.commission_value,"+
			" guest_detail.payment_type_code,"+
			" guest_detail.market_code,"+
			" guest_detail.booking_source_code,"+
			" cfg_init_room_rate.tax_and_service_code,"+
			" cfg_init_room_rate.charge_frequency_code ").
		Joins("LEFT JOIN contact_person ON (folio.contact_person_id1 = contact_person.id)").
		Joins("LEFT JOIN guest_detail ON (folio.guest_detail_id = guest_detail.id)").
		Joins("LEFT JOIN guest_general ON (folio.guest_general_id = guest_general.id)").
		Joins("LEFT JOIN cfg_init_room_rate ON (guest_detail.room_rate_code = cfg_init_room_rate.code)").
		Where("folio.number=?", FolioNumber).Take(&DataOutput).Error
	if err != nil {
		return
	}

	var ScheduledRateData db_var.Guest_scheduled_rate
	if err := DB.Table(db_var.TableName.GuestScheduledRate).
		Where("folio_number=?", FolioNumber).
		Where("from_date<=?", General.FormatDate1(PostingDate)).
		Where("to_date>=?", General.FormatDate1(PostingDate)).
		Limit(1).
		Scan(&ScheduledRateData).Error; err != nil {
		return 0, err
	}

	// if err == nil {
	IsCanPostRoomCharge := IsCanPostCharge(c, DB, DataOutput.RoomRate.ChargeFrequencyCode, DataOutput.DateArrival)
	IsBreakfastX = IsFolioHaveBreakfast(DB, Dataset, FolioNumber)
	IsVoucherActiveX = IsVoucherComplimentStillActive(ctx, DB, DataOutput.Folio.VoucherNumber, DataOutput.DateArrival, PostingDate)
	VoucherTypeCode = GetVoucherType(ctx, DB, DataOutput.Folio.VoucherNumber)
	if IsVoucherActiveX && (VoucherTypeCode == GlobalVar.VoucherType.Compliment) {
		ComplimentHU = GlobalVar.RoomStatus.Compliment
	} else {
		ComplimentHU = DataOutput.Folio.ComplimentHu
		if ScheduledRateData.Id > 0 {
			ComplimentHU = ScheduledRateData.ComplimentHu
		}
	}

	SDFrontOffice := Dataset.GlobalSubDepartment.FrontOffice
	GARoomCharge := Dataset.GlobalAccount.RoomCharge
	if !(!AllowZeroAmount && (ComplimentHU == "H" || ComplimentHU == "C")) {
		if IsCanPostRoomCharge {
			RoomNumber = *DataOutput.GuestDetail.RoomNumber

			CurrencyCode = DataOutput.GuestDetail.CurrencyCode
			ExchangeRate = DataOutput.GuestDetail.ExchangeRate
			if ExchangeRate <= 0 {
				ExchangeRate = 1
			}

			IsScheduledRateX = ScheduledRateData.Id > 0
			if IsScheduledRateX && !IsVoucherActiveX {
				RoomRateCode = ScheduledRateData.RoomRateCode
				if RoomRateCode == "" {
					RoomRateCode = DataOutput.GuestDetail.RoomRateCode
					RoomRateTaxServiceCode = DataOutput.RoomRate.TaxAndServiceCode
				} else {
					RoomRateTaxServiceCode = GetRoomRateTaxAndServiceCode(DB, RoomRateCode)
				}

				RoomChargeB4Breakdown = *ScheduledRateData.Rate / ExchangeRate
				RoomRateAmount = RoomChargeB4Breakdown

				RoomRateAmountOriginal = GetRoomRateAmount(ctx, DB, Dataset, DataOutput.GuestDetail.RoomRateCode, General.FormatDate1(PostingDate), DataOutput.GuestDetail.Adult, *DataOutput.GuestDetail.Child, false)
				ComplimentHU = ScheduledRateData.ComplimentHu
				Discount = 0
			} else {
				RoomRateCode = DataOutput.GuestDetail.RoomRateCode
				RoomRateTaxServiceCode = DataOutput.RoomRate.TaxAndServiceCode
				if General.IsWeekend(PostingDate, Dataset) {
					RoomChargeB4Breakdown = *DataOutput.GuestDetail.WeekendRate
				} else {
					RoomChargeB4Breakdown = *DataOutput.GuestDetail.WeekdayRate
				}
				RoomRateAmount = RoomChargeB4Breakdown
				RoomRateAmountOriginal = GetRoomRateAmount(ctx, DB, Dataset, DataOutput.GuestDetail.RoomRateCode, General.FormatDate1(PostingDate), DataOutput.GuestDetail.Adult, *DataOutput.GuestDetail.Child, false)
				ComplimentHU = DataOutput.Folio.ComplimentHu

				if (IsVoucherActiveX) && (VoucherTypeCode != GlobalVar.VoucherType.Compliment) {
					Discount = GetVoucherPrice(DB, DataOutput.Folio.VoucherNumber)
				} else {
					if *DataOutput.GuestDetail.DiscountPercent > 0 {
						Discount = General.RoundTo(RoomChargeB4Breakdown * *DataOutput.GuestDetail.Discount / 100)
					} else {
						Discount = *DataOutput.GuestDetail.Discount
					}
				}

				if !Dataset.ProgramConfiguration.PostDiscount {
					RoomChargeB4Breakdown = RoomChargeB4Breakdown - Discount
					Discount = 0
				}
			}

			if (ComplimentHU == "H" || ComplimentHU == "C") && (RoomChargeB4Breakdown > 0) {
				if AllowZeroAmount {
					Result = 0
					//Post Room Charge
					CorrectionBreakdown = GetSubFolioCorrectionBreakdown(c, DB)
					BreakDown1 = GetSubFolioBreakdown1(c, DB)
					_, _, err = InsertSubFolio2(c, DB, Dataset, FolioNumber, SubFolioGroupCode, RoomNumber, SDFrontOffice, GARoomCharge, GARoomCharge, "", "", CurrencyCode, "Auto Room Charge", "", "", GlobalVar.TransactionType.Debit, "", "", CorrectionBreakdown, BreakDown1, GlobalVar.SubFolioPostingType.Room, 0, 0, 0, 0, ExchangeRate, AllowZeroAmount, true, UserID)
					if err != nil {
						return
					}
					if IsInHousePosted(DB, PostingDate, FolioNumber) {
						err = DeleteGuestInHouse(ctx, DB, PostingDate, FolioNumber)
						if err != nil {
							return
						}
					}
					err = InsertGuestInHouse(DB, PostingDate, FolioNumber, DataOutput.Folio.GroupCode, DataOutput.GuestDetail.RoomTypeCode, DataOutput.GuestDetail.BedTypeCode, *DataOutput.GuestDetail.RoomNumber, DataOutput.GuestDetail.RoomRateCode, *DataOutput.GuestDetail.BusinessSourceCode, *DataOutput.GuestDetail.CommissionTypeCode, DataOutput.GuestDetail.PaymentTypeCode, *DataOutput.GuestDetail.MarketCode,
						*DataOutput.ContactPerson.TitleCode, *DataOutput.ContactPerson.FullName, *DataOutput.ContactPerson.Street, *DataOutput.ContactPerson.City, *DataOutput.ContactPerson.CityCode, *DataOutput.ContactPerson.CountryCode, *DataOutput.ContactPerson.StateCode, *DataOutput.ContactPerson.PostalCode, *DataOutput.ContactPerson.Phone1, *DataOutput.ContactPerson.Phone2, *DataOutput.ContactPerson.Fax, *DataOutput.ContactPerson.Email,
						*DataOutput.ContactPerson.Website, *DataOutput.ContactPerson.CompanyCode, *DataOutput.ContactPerson.GuestTypeCode, *DataOutput.GuestGeneral.SalesCode, DataOutput.Folio.ComplimentHu, *DataOutput.GuestGeneral.Notes, DataOutput.GuestDetail.Adult, *DataOutput.GuestDetail.Child, 0, 0, 0, *DataOutput.GuestDetail.CommissionValue, *DataOutput.GuestDetail.DiscountPercent, 0, General.BoolToUint8(IsScheduledRateX), General.BoolToUint8(IsBreakfastX), *DataOutput.GuestDetail.BookingSourceCode,
						*DataOutput.GuestGeneral.PurposeOfCode, *DataOutput.ContactPerson.CustomLookupFieldCode01, *DataOutput.ContactPerson.CustomLookupFieldCode02, 0, "", *DataOutput.ContactPerson.NationalityCode, UserID)
					if err != nil {
						return
					}
				} else {
					Result = 4
				}
			} else {
				//Cari Total Tax Service Room Charge
				RoomChargeBasic, RoomChargeTax, RoomChargeService = GetBasicTaxService(DB, GARoomCharge, RoomRateTaxServiceCode, RoomChargeB4Breakdown)
				RoomChargeB4Breakdown = RoomChargeBasic + RoomChargeTax + RoomChargeService
				//Proses Query Breakdown
				Breakdown := []DBVar.Cfg_init_room_rate_breakdown{}
				if IsScheduledRateX {
					err = DB.Table(DBVar.TableName.CfgInitRoomRateBreakdown).Select(
						" cfg_init_room_rate_breakdown.outlet_code,"+
							" cfg_init_room_rate_breakdown.product_code,"+
							" cfg_init_room_rate_breakdown.sub_department_code,"+
							" cfg_init_room_rate_breakdown.account_code,"+
							" cfg_init_room_rate_breakdown.company_code,"+
							" cfg_init_room_rate_breakdown.quantity,"+
							" cfg_init_room_rate_breakdown.is_amount_percent,"+
							" cfg_init_room_rate_breakdown.amount,"+
							" cfg_init_room_rate_breakdown.per_pax,"+
							" cfg_init_room_rate_breakdown.include_child,"+
							" cfg_init_room_rate_breakdown.remark,"+
							" cfg_init_room_rate_breakdown.tax_and_service_code,"+
							" cfg_init_room_rate_breakdown.charge_frequency_code,"+
							" cfg_init_room_rate_breakdown.max_pax,"+
							" cfg_init_room_rate_breakdown.extra_pax,"+
							" cfg_init_room_rate_breakdown.per_pax_extra,"+
							" cfg_init_room_rate_breakdown.id,"+
							" cfg_init_room_rate_breakdown.created_by,"+
							" cfg_init_room_rate_breakdown.updated_by").
						Where("cfg_init_room_rate_breakdown.room_rate_code=?", RoomRateCode).
						Scan(&Breakdown).Error
					if err != nil {
						return
					}
				} else {
					err = DB.Table(DBVar.TableName.GuestBreakdown).Select(
						" guest_breakdown.outlet_code,"+
							" guest_breakdown.product_code,"+
							" guest_breakdown.sub_department_code,"+
							" guest_breakdown.account_code,"+
							" guest_breakdown.company_code,"+
							" guest_breakdown.quantity,"+
							" guest_breakdown.is_amount_percent,"+
							" guest_breakdown.amount,"+
							" guest_breakdown.per_pax,"+
							" guest_breakdown.include_child,"+
							" guest_breakdown.remark,"+
							" guest_breakdown.tax_and_service_code,"+
							" guest_breakdown.charge_frequency_code,"+
							" guest_breakdown.max_pax,"+
							" guest_breakdown.extra_pax,"+
							" guest_breakdown.per_pax_extra,"+
							" guest_breakdown.id,"+
							" guest_breakdown.created_by,"+
							" guest_breakdown.updated_by").
						Where("guest_breakdown.folio_number=?", FolioNumber).
						Scan(&Breakdown).Error
					if err != nil {
						return
					}
				}
				//Calculate Total Breakdown
				TotalBreakdown = 0
				for _, breakdown := range Breakdown {
					if IsCanPostCharge(c, DB, breakdown.ChargeFrequencyCode, DataOutput.DateArrival) {
						if breakdown.IsAmountPercent > 0 {
							BreakdownAmount = GetTotalBreakdownAmount(breakdown.Quantity, RoomChargeB4Breakdown*breakdown.Amount/100, breakdown.ExtraPax, breakdown.PerPax > 0, breakdown.IncludeChild > 0, breakdown.PerPaxExtra > 0, breakdown.MaxPax, DataOutput.GuestDetail.Adult, *DataOutput.GuestDetail.Child) / ExchangeRate
							//                  BreakdownAmount = GetTotalBreakdownAmount(MyQGuestBreakdownCalculatequantity.AsFloat, RoomRateAmount * MyQGuestBreakdownCalculateamount.AsFloat/100, MyQGuestBreakdownCalculateextra_pax.AsFloat, MyQGuestBreakdownCalculateper_pax.AsVariant, MyQGuestBreakdownCalculateinclude_child.AsVariant, MyQGuestBreakdownCalculateper_pax_extra.AsVariant, MyQGuestBreakdownCalculatemax_pax.AsInteger, MyQFolioCalculateadult.AsInteger, MyQFolioCalculatechild.AsInteger) / ExchangeRate
						} else {
							BreakdownAmount = GetTotalBreakdownAmount(breakdown.Quantity, breakdown.Amount, breakdown.ExtraPax, breakdown.PerPax > 0, breakdown.IncludeChild > 0, breakdown.PerPaxExtra > 0, breakdown.MaxPax, DataOutput.GuestDetail.Adult, *DataOutput.GuestDetail.Child) / ExchangeRate
						}
						BreakdownBasic, BreakdownTax, BreakdownService = GetBasicTaxService(DB, breakdown.AccountCode, breakdown.TaxAndServiceCode, BreakdownAmount)
						BreakdownAmount = BreakdownBasic + BreakdownTax + BreakdownService
						TotalBreakdown = TotalBreakdown + BreakdownAmount
					}
				}

				//Get Commission from Business Source
				BusinessSourceCode = *DataOutput.GuestDetail.BusinessSourceCode
				Commission = 0
				if BusinessSourceCode != "" {
					Commission = GetCommission(c, DB, *DataOutput.GuestDetail.CommissionTypeCode, *DataOutput.GuestDetail.CommissionValue, RoomChargeB4Breakdown, RoomChargeBasic, DataOutput.DateArrival) / ExchangeRate
				}
				//Room Charge - Total Breakdown - Total Commission
				RoomChargeAfterBreakdown = RoomChargeB4Breakdown - TotalBreakdown - Commission
				if (RoomChargeAfterBreakdown > 0) || (AllowZeroAmount && (RoomChargeAfterBreakdown == 0)) {
					Result = 0
					//Cari Basic dari Room Charge Bersih (yang sudah dikurangi Breakdown dan Total Commission)
					RoomChargeBasic, RoomChargeTax, RoomChargeService = GetBasicTaxService2(DB, GARoomCharge, RoomRateTaxServiceCode, RoomChargeAfterBreakdown)
					RoomChargeAfterBreakdown = RoomChargeBasic + RoomChargeTax + RoomChargeService
					//Post Room Charge
					CorrectionBreakdown = GetSubFolioCorrectionBreakdown(c, DB)
					BreakDown1 = GetSubFolioBreakdown1(c, DB)
					_, _, err = InsertSubFolio2(c, DB, Dataset, FolioNumber, SubFolioGroupCode, RoomNumber, SDFrontOffice, GARoomCharge, GARoomCharge, "", "", CurrencyCode, "Auto Room Charge", "", "", GlobalVar.TransactionType.Debit, "", "", CorrectionBreakdown, BreakDown1, GlobalVar.SubFolioPostingType.Room, 0, RoomChargeBasic, RoomChargeTax, RoomChargeService, ExchangeRate, AllowZeroAmount, true, UserID)
					if err != nil {
						return
					}
					if IsInHousePosted(DB, PostingDate, FolioNumber) {
						DeleteGuestInHouse(ctx, DB, PostingDate, FolioNumber)
					}
					err = InsertGuestInHouse(DB, PostingDate, FolioNumber, DataOutput.Folio.GroupCode, DataOutput.GuestDetail.RoomTypeCode, DataOutput.GuestDetail.BedTypeCode, *DataOutput.GuestDetail.RoomNumber, RoomRateCode, *DataOutput.GuestDetail.BusinessSourceCode, *DataOutput.GuestDetail.CommissionTypeCode, DataOutput.GuestDetail.PaymentTypeCode, *DataOutput.GuestDetail.MarketCode,
						*DataOutput.ContactPerson.TitleCode, *DataOutput.ContactPerson.FullName, *DataOutput.ContactPerson.Street, *DataOutput.ContactPerson.City, *DataOutput.ContactPerson.CityCode, *DataOutput.ContactPerson.CountryCode, *DataOutput.ContactPerson.StateCode, *DataOutput.ContactPerson.PostalCode, *DataOutput.ContactPerson.Phone1, *DataOutput.ContactPerson.Phone2, *DataOutput.ContactPerson.Fax, *DataOutput.ContactPerson.Email,
						*DataOutput.ContactPerson.Website, *DataOutput.ContactPerson.CompanyCode, *DataOutput.ContactPerson.GuestTypeCode, *DataOutput.GuestGeneral.SalesCode, ComplimentHU, *DataOutput.GuestGeneral.Notes, DataOutput.GuestDetail.Adult, *DataOutput.GuestDetail.Child, RoomRateAmountOriginal, RoomRateAmount, *DataOutput.GuestDetail.Discount, *DataOutput.GuestDetail.CommissionValue, *DataOutput.GuestDetail.DiscountPercent, 0, General.BoolToUint8(IsScheduledRateX), General.BoolToUint8(IsBreakfastX), *DataOutput.GuestDetail.BookingSourceCode,
						*DataOutput.GuestGeneral.PurposeOfCode, *DataOutput.ContactPerson.CustomLookupFieldCode01, *DataOutput.ContactPerson.CustomLookupFieldCode02, 0, "", *DataOutput.ContactPerson.NationalityCode, UserID)
					if err != nil {
						return
					}
					//Posting Breakdown
					for _, breakdown := range Breakdown {
						if IsCanPostCharge(c, DB, breakdown.ChargeFrequencyCode, DataOutput.DateArrival) {
							if breakdown.IsAmountPercent > 0 {
								BreakdownAmount = GetTotalBreakdownAmount(breakdown.Quantity, RoomChargeB4Breakdown*breakdown.Amount/100, breakdown.ExtraPax, breakdown.PerPax > 0, breakdown.IncludeChild > 0, breakdown.PerPaxExtra > 0, breakdown.MaxPax, DataOutput.GuestDetail.Adult, *DataOutput.GuestDetail.Child) / ExchangeRate
								//                  BreakdownAmount = GetTotalBreakdownAmount(MyQGuestBreakdownCalculatequantity.AsFloat, RoomRateAmount * MyQGuestBreakdownCalculateamount.AsFloat/100, MyQGuestBreakdownCalculateextra_pax.AsFloat, MyQGuestBreakdownCalculateper_pax.AsVariant, MyQGuestBreakdownCalculateinclude_child.AsVariant, MyQGuestBreakdownCalculateper_pax_extra.AsVariant, MyQGuestBreakdownCalculatemax_pax.AsInteger, MyQFolioCalculateadult.AsInteger, MyQFolioCalculatechild.AsInteger) / ExchangeRate
							} else {
								BreakdownAmount = GetTotalBreakdownAmount(breakdown.Quantity, breakdown.Amount, breakdown.ExtraPax, breakdown.PerPax > 0, breakdown.IncludeChild > 0, breakdown.PerPaxExtra > 0, breakdown.MaxPax, DataOutput.GuestDetail.Adult, *DataOutput.GuestDetail.Child) / ExchangeRate
							}
							_, _, err = InsertSubFolio(c, DB, Dataset, true, GARoomCharge, breakdown.TaxAndServiceCode, DBVar.Sub_folio{
								FolioNumber:         FolioNumber,
								RoomNumber:          RoomNumber,
								SubDepartmentCode:   breakdown.SubDepartmentCode,
								AccountCode:         breakdown.AccountCode,
								ProductCode:         breakdown.ProductCode,
								GroupCode:           SubFolioGroupCode,
								Remark:              "Breakdown: " + breakdown.Remark,
								TypeCode:            GlobalVar.TransactionType.Debit,
								PostingType:         GlobalVar.SubFolioPostingType.Room,
								Quantity:            1,
								CurrencyCode:        CurrencyCode,
								CorrectionBreakdown: CorrectionBreakdown,
								Breakdown1:          BreakDown1,
								DirectBillCode:      breakdown.CompanyCode,
								Amount:              BreakdownAmount,
								ExchangeRate:        ExchangeRate,
								ExtraChargeId:       breakdown.Id})
							if err != nil {
								return
							}
							err = InsertGuestInHouseBreakdown(DB, PostingDate, FolioNumber, breakdown.OutletCode, breakdown.ProductCode, breakdown.SubDepartmentCode, breakdown.AccountCode, breakdown.Remark, breakdown.TaxAndServiceCode, breakdown.ChargeFrequencyCode, breakdown.Quantity, breakdown.Amount, breakdown.ExtraPax, breakdown.MaxPax, breakdown.PerPax,
								breakdown.IncludeChild, breakdown.PerPaxExtra, UserID)
							if err != nil {
								return
							}
						}
					}

					//Posting Commission from Business Source
					if Commission > 0 {
						GAAPCommission := Dataset.GlobalAccount.APCommission
						_, _, err = InsertSubFolio(c, DB, Dataset, true, GARoomCharge, "", DBVar.Sub_folio{
							FolioNumber:         FolioNumber,
							RoomNumber:          RoomNumber,
							SubDepartmentCode:   SDFrontOffice,
							AccountCode:         GAAPCommission,
							GroupCode:           SubFolioGroupCode,
							Remark:              "Breakdown Commission",
							TypeCode:            GlobalVar.TransactionType.Debit,
							PostingType:         GlobalVar.SubFolioPostingType.Room,
							Quantity:            1,
							CurrencyCode:        CurrencyCode,
							CorrectionBreakdown: CorrectionBreakdown,
							Breakdown1:          BreakDown1,
							DirectBillCode:      BusinessSourceCode,
							Amount:              Commission,
							ExchangeRate:        ExchangeRate})
						if err != nil {
							return
						}
					}
				} else {
					Result = 1
				}
			}
		} else {
			Result = 3
		}
	} else {
		Result = 2
	}
	return
}

func PostingExtraCharge(ctx context.Context, c *gin.Context, DB *gorm.DB, Dataset *global_var.TDataset, FolioNumber, ExtraChargeID uint64, SubFolioGroupCode, UserID string) (Result int, err error) {
	ctx, span := global_var.Tracer.Start(ctx, "PostingExtraCharge")
	defer span.End()

	var AmountB4Breakdown, AmountAfterBreakdown, AmountBasic, AmountTax, AmountService, TotalBreakdown, BreakdownAmount, BreakdownBasic, BreakdownTax, BreakdownService, ExchangeRate float64
	var RoomNumber, AccountCodeMaster, CurrencyCode string
	var Adult, Child int
	var CorrectionBreakdown, BreakDown1 uint64
	var SubFolioGroupCodeX string
	Result = 255
	//0: Succes
	//1: Breakdown Too Large
	//2: No Extra Charge for Today
	//255: Posting Extra Charge Failed

	//Proses Posting Extra Charge
	type ExtraChargeStruct struct {
		ExtraCharge DBVar.Guest_extra_charge `gorm:"embedded"`
		GuestDetail DBVar.Guest_detail       `gorm:"embedded"`
		DateArrival time.Time
	}
	var DataOutput []ExtraChargeStruct
	Query := DB.WithContext(ctx).Table(DBVar.TableName.GuestExtraCharge).Select(
		" guest_extra_charge.*,"+
			" guest_detail.currency_code,"+
			" guest_detail.exchange_rate,"+
			" DATE(guest_detail.arrival) AS DateArrival,"+
			" guest_detail.adult,"+
			" guest_detail.child,"+
			" guest_detail.room_number").
		Joins("LEFT OUTER JOIN folio ON (guest_extra_charge.folio_number = folio.number)").
		Joins("LEFT OUTER JOIN guest_detail ON (folio.guest_detail_id = guest_detail.id)").
		Where("folio_number=?", FolioNumber)
	if ExtraChargeID > 0 {
		Query.Where("guest_extra_charge.id=?", ExtraChargeID)
	}
	Query.Scan(&DataOutput)

	for _, extraCharge := range DataOutput {
		if IsCanPostCharge(c, DB, extraCharge.ExtraCharge.ChargeFrequencyCode, extraCharge.DateArrival) {
			SubFolioGroupCodeX = extraCharge.ExtraCharge.GroupCode
			if SubFolioGroupCodeX != "" {
				SubFolioGroupCode = SubFolioGroupCodeX
			}

			RoomNumber = *extraCharge.GuestDetail.RoomNumber
			Adult = extraCharge.GuestDetail.Adult
			Child = *extraCharge.GuestDetail.Child
			AmountB4Breakdown = GetTotalBreakdownAmount(extraCharge.ExtraCharge.Quantity, extraCharge.ExtraCharge.Amount, *extraCharge.ExtraCharge.ExtraPax, *extraCharge.ExtraCharge.PerPax > 0, *extraCharge.ExtraCharge.IncludeChild > 0,
				*extraCharge.ExtraCharge.PerPaxExtra > 0, extraCharge.ExtraCharge.MaxPax, Adult, Child)

			//Cari Total Tax Service Extra Charge
			AmountBasic, AmountTax, AmountService = GetBasicTaxService(DB, extraCharge.ExtraCharge.AccountCode, *extraCharge.ExtraCharge.TaxAndServiceCode, AmountB4Breakdown)
			AmountB4Breakdown = AmountBasic + AmountTax + AmountService
			//Proses Query Breakdown
			var GuestExtraChargeBreakdown []DBVar.Guest_extra_charge_breakdown
			DB.WithContext(ctx).Table(DBVar.TableName.GuestExtraChargeBreakdown).Where("guest_extra_charge_id=?", extraCharge.ExtraCharge.Id).Scan(&GuestExtraChargeBreakdown)

			//Calculate Total Breakdown
			TotalBreakdown = 0
			for _, breakdown := range GuestExtraChargeBreakdown {
				if IsCanPostCharge(c, DB, breakdown.ChargeFrequencyCode, extraCharge.DateArrival) {
					if *breakdown.IsAmountPercent > 0 {
						BreakdownAmount = GetTotalBreakdownAmount(breakdown.Quantity, extraCharge.ExtraCharge.Amount*breakdown.Amount/100, *breakdown.ExtraPax, *breakdown.PerPax > 0, *breakdown.IncludeChild > 0, *breakdown.PerPaxExtra > 0, breakdown.MaxPax, Adult, Child)
					} else {
						BreakdownAmount = GetTotalBreakdownAmount(breakdown.Quantity, breakdown.Amount, *breakdown.ExtraPax, *breakdown.PerPax > 0, *breakdown.IncludeChild > 0, *breakdown.PerPaxExtra > 0, breakdown.MaxPax, Adult, Child)
					}
					BreakdownBasic, BreakdownTax, BreakdownService = GetBasicTaxService(DB, breakdown.AccountCode, *breakdown.TaxAndServiceCode, BreakdownAmount)
					BreakdownAmount = BreakdownBasic + BreakdownTax + BreakdownService
					TotalBreakdown = TotalBreakdown + BreakdownAmount
				}
			}

			//Amount - Total Breakdown
			AmountAfterBreakdown = AmountB4Breakdown - TotalBreakdown

			//Cari Basic dari Extra Charge Bersih (yang sudah dikurangi Breakdown)

			if AmountAfterBreakdown > 0 {
				Result = 0
				AmountBasic, AmountTax, AmountService = GetBasicTaxService2(DB, extraCharge.ExtraCharge.AccountCode, *extraCharge.ExtraCharge.TaxAndServiceCode, AmountAfterBreakdown)
				AmountAfterBreakdown = AmountBasic + AmountTax + AmountService
				//Post Extra Charge
				AccountCodeMaster = extraCharge.ExtraCharge.AccountCode
				CorrectionBreakdown = GetSubFolioCorrectionBreakdown(c, DB)
				BreakDown1 = GetSubFolioBreakdown1(c, DB)
				_, _, err = InsertSubFolio2(c, DB, Dataset, FolioNumber, SubFolioGroupCode, RoomNumber, extraCharge.ExtraCharge.SubDepartmentCode, extraCharge.ExtraCharge.AccountCode,
					AccountCodeMaster, *extraCharge.ExtraCharge.ProductCode, *extraCharge.ExtraCharge.PackageCode, "", "Extra Charge", "", "", GlobalVar.TransactionType.Debit,
					"", "", CorrectionBreakdown, BreakDown1, GlobalVar.SubFolioPostingType.ExtraCharge, extraCharge.ExtraCharge.Id, AmountBasic, AmountTax, AmountService, 0, false, true, UserID)
				if err != nil {
					return
				}
				//Posting Breakdown
				for _, breakdown := range GuestExtraChargeBreakdown {
					if IsCanPostCharge(c, DB, breakdown.ChargeFrequencyCode, extraCharge.DateArrival) {
						if *breakdown.IsAmountPercent > 0 {
							BreakdownAmount = GetTotalBreakdownAmount(breakdown.Quantity, extraCharge.ExtraCharge.Amount*breakdown.Amount/100, *breakdown.ExtraPax, *breakdown.PerPax > 0, *breakdown.IncludeChild > 0, *breakdown.PerPaxExtra > 0, breakdown.MaxPax, Adult, Child)
						} else {
							BreakdownAmount = GetTotalBreakdownAmount(breakdown.Quantity, breakdown.Amount, *breakdown.ExtraPax, *breakdown.PerPax > 0, *breakdown.IncludeChild > 0, *breakdown.PerPaxExtra > 0, breakdown.MaxPax, Adult, Child)
						}
						_, _, err = InsertSubFolio(c, DB, Dataset, true, AccountCodeMaster, *breakdown.TaxAndServiceCode, DBVar.Sub_folio{
							FolioNumber:         FolioNumber,
							RoomNumber:          RoomNumber,
							SubDepartmentCode:   breakdown.SubDepartmentCode,
							AccountCode:         breakdown.AccountCode,
							PackageCode:         *extraCharge.ExtraCharge.PackageCode,
							ProductCode:         *breakdown.ProductCode,
							GroupCode:           SubFolioGroupCode,
							Remark:              "Extra Charge Breakdown",
							TypeCode:            GlobalVar.TransactionType.Debit,
							PostingType:         GlobalVar.SubFolioPostingType.ExtraCharge,
							Quantity:            1,
							CurrencyCode:        CurrencyCode,
							CorrectionBreakdown: CorrectionBreakdown,
							Breakdown1:          BreakDown1,
							DirectBillCode:      *breakdown.CompanyCode,
							Amount:              BreakdownAmount,
							ExchangeRate:        ExchangeRate,
							ExtraChargeId:       breakdown.Id})
						if err != nil {
							return
						}
					}
				}
			} else {
				Result = 1
			}
		} else {
			Result = 2
		}
	}
	return
}

func GetGuestDepositBreakdownAmountForeign(DB *gorm.DB, CorrectionBreakdown int64) (Amount float64) {
	DB.Table(DBVar.TableName.GuestDeposit).Select("SUM(IF(type_code='"+GlobalVar.TransactionType.Debit+"', amount_foreign, -amount_foreign)) AS TotalAmount").
		Where("correction_breakdown=?", CorrectionBreakdown).
		Where("void=0").
		Group("correction_breakdown").
		Limit(1).
		Scan(&Amount)

	return
}

func GetSubFolioBreakdownAmountForeign(DB *gorm.DB, CorrectionBreakdown int64) (Amount float64) {
	DB.Table(DBVar.TableName.SubFolio).Select("SUM(IF(type_code='"+GlobalVar.TransactionType.Debit+"',(quantity * amount_foreign),-(quantity * amount_foreign))) AS TotalAmount ").
		Where("correction_breakdown=?", CorrectionBreakdown).
		Where("void=0").
		Group("correction_breakdown").
		Limit(1).
		Scan(&Amount)

	return
}

func GetSubFolioBreakdown1AmountForeign(DB *gorm.DB, Breakdown1 int64) (Amount float64) {
	DB.Table(DBVar.TableName.SubFolio).Select("SUM(IF(type_code='"+GlobalVar.TransactionType.Debit+"',(quantity * amount_foreign),-(quantity * amount_foreign))) AS TotalAmount").
		Where("breakdown1=?", Breakdown1).
		Where("void=0").
		Group("breakdown1").
		Limit(1).
		Scan(&Amount)

	return
}

func GetSubFolioBreakdownAmountForeign2(DB *gorm.DB, CorrectionBreakdown, Breakdown2 int64) (Amount float64) {
	DB.Table(DBVar.TableName.SubFolio).Select("SUM(IF(type_code='"+GlobalVar.TransactionType.Debit+"',(quantity * amount_foreign),-(quantity * amount_foreign))) AS TotalAmount").
		Where("correction_breakdown=?", CorrectionBreakdown).
		Where("breakdown2=?", Breakdown2).
		Where("void=0").
		Group("correction_breakdown,breakdown2").
		Limit(1).
		Scan(&Amount)

	return
}

func GetSubFolioBreakdown1AmountForeign2(DB *gorm.DB, Breakdown1, Breakdown2 int64) (Amount float64) {
	DB.Table(DBVar.TableName.SubFolio).Select("SUM(IF(type_code='"+GlobalVar.TransactionType.Debit+"',(quantity * amount_foreign),-(quantity * amount_foreign))) AS TotalAmount").
		Where("breakdown1=?", Breakdown1).
		Where("breakdown2=?", Breakdown2).
		Where("void=0").
		Group("breakdown1").
		Limit(1).
		Scan(&Amount)

	return
}

func GetCheckQuantityByCorrectionBreakdown(DB *gorm.DB, CorrectionBreakdown int64) (Amount float64) {
	DB.Table(DBVar.TableName.SubFolio).Select("IF(cfg_init_account_sub_group.group_code='"+GlobalVar.GlobalAccountGroup.Charge+"', IF(sub_folio.type_code='"+GlobalVar.TransactionType.Debit+"', sub_folio.quantity, -sub_folio.quantity), IF(sub_folio.type_code='"+GlobalVar.TransactionType.Credit+"', sub_folio.quantity, -sub_folio.quantity)) AS TotalQuantity ").
		Joins("LEFT OUTER JOIN cfg_init_account ON (sub_folio.account_code = cfg_init_account.code)").
		Joins("LEFT OUTER JOIN cfg_init_account_sub_group ON (cfg_init_account.sub_group_code = cfg_init_account_sub_group.code)").
		Where("sub_folio.correction_breakdown=?", CorrectionBreakdown).
		Where("sub_folio.void=0").
		Group("sub_folio.correction_breakdown").
		Limit(1).
		Scan(&Amount)

	return
}

func IsSubFolioFromPOS(DB *gorm.DB, SubFolioID, TransferPairID int64) bool {
	Result := 0
	DB.Table(DBVar.TableName.PosCheckTransaction).Select("id").
		Where("(sub_folio_id=? OR sub_folio_id=?) AND sub_folio_id<>0", SubFolioID, TransferPairID).
		Scan(&Result)

	return Result > 0
}

func GetAccountGroupCode(DB *gorm.DB, AccountCode string) string {
	Result := ""
	DB.Table(DBVar.TableName.CfgInitAccount).Select("group_code").
		Joins("JOIN cfg_init_account_sub_group ON (cfg_init_account.sub_group_code = cfg_init_account_sub_group.code)").
		Where("cfg_init_account.code=?", AccountCode).
		Scan(&Result)

	return Result
}

func IsCanVoidOrCorrect(DB *gorm.DB, Dataset *global_var.TDataset, Mode byte, TransactionID uint64, AccountCode string) (Result bool) {
	var AccountSubGroupCode string
	AccountSubGroupCode = GetAccountSubGroupCode(DB, AccountCode)
	GAAPRefundDeposit := Dataset.GlobalAccount.APRefundDeposit
	if AccountCode == GAAPRefundDeposit {
		Result = !MasterData.GetFieldBool(DB, DBVar.TableName.AccApRefundDepositPaymentDetail, "id", "sub_folio_id", TransactionID, false)
	} else if AccountSubGroupCode == GlobalVar.GlobalAccountSubGroup.AccountPayable {
		Result = !MasterData.GetFieldBool(DB, DBVar.TableName.AccApCommissionPaymentDetail, "id", "sub_folio_id", TransactionID, false)
	} else if (AccountSubGroupCode == GlobalVar.GlobalAccountSubGroup.CreditDebitCard) || (AccountSubGroupCode == GlobalVar.GlobalAccountSubGroup.BankTransfer) {
		if Mode == 0 {
			Result = !MasterData.GetFieldBool(DB, DBVar.TableName.AccCreditCardReconDetail, "id", "guest_deposit_id", TransactionID, false)
		} else {
			Result = !MasterData.GetFieldBool(DB, DBVar.TableName.AccCreditCardReconDetail, "id", "sub_folio_id", TransactionID, false)
		}
	} else if AccountSubGroupCode == GlobalVar.GlobalAccountSubGroup.AccountReceivable {
		Result = !MasterData.GetFieldBool(DB, DBVar.TableName.InvoiceItem, "id", "sub_folio_id", TransactionID, false)
	} else if AccountSubGroupCode == GlobalVar.GlobalAccountSubGroup.Payment {
		if Mode == 0 {
			Result = !MasterData.GetFieldBool(DB, DBVar.TableName.AccForeignCash, "id", "id_transaction", General.Uint64ToStr(TransactionID)+"' AND id_change=0 AND stock<>amount_foreign AND id_table='"+strconv.Itoa(GlobalVar.ForeignCashTableID.GuestDeposit), false)
		} else {
			Result = !MasterData.GetFieldBool(DB, DBVar.TableName.InvoiceItem, "id", "sub_folio_id", General.Uint64ToStr(TransactionID)+"' AND id_change=0 AND stock<>amount_foreign AND id_table='"+strconv.Itoa(GlobalVar.ForeignCashTableID.SubFolio), false)
		}
	} else {
		Result = true
	}
	return
}

func GetCountAutoPosting(ctx context.Context, DB *gorm.DB, FolioNumber uint64, AccountCode, PostingType string, PostingDate time.Time) (Count int) {
	Query := DB.WithContext(ctx).Table(DBVar.TableName.SubFolio).Select("account_code, is_correction").
		Where("belongs_to=?", FolioNumber).
		Where("audit_date=? AND void=0", General.FormatDate1(PostingDate))
	if PostingType != "" {
		Query.Where("posting_type", PostingType)
	}
	Query.Group("correction_breakdown")
	DB.Table("(?) AS SubFolioX", Query).Select("COUNT(account_code)").Where("account_code=? AND is_correction=0", AccountCode).Limit(1).Scan(&Count)
	return
}

func IsAlreadyAutoPosting(ctx context.Context, DB *gorm.DB, FolioNumber uint64, AccountCode, PostingType string, PostingDate time.Time) bool {
	return GetCountAutoPosting(ctx, DB, FolioNumber, AccountCode, PostingType, PostingDate) > 0
}

func GetBreakdownAutoPosting(DB *gorm.DB, FolioNumber uint64, AccountCode, PostingType string, PostingDate time.Time) (CorrectionBreakdown uint64) {
	Query := DB.Table(DBVar.TableName.SubFolio).Select("account_code, correction_breakdown").
		Where("belongs_to=?", FolioNumber).
		Where("audit_date_unixx=UNIX_TIMESTAMP(?) AND void=0", General.FormatDate2("2006-01-02 00:00:00", PostingDate)).
		Where("posting_type", PostingType).
		Group("correction_breakdown")
	DB.Table("(?) AS ExtraCharge", Query).Select("correction_breakdown").Where("account_code=?", AccountCode).Limit(1).Scan(&CorrectionBreakdown)

	return

}

func IsCanCancelCheckInFolio(DB *gorm.DB, Dataset *global_var.TDataset, FolioNumber uint64) bool {
	var DataOutput []DBVar.Sub_folio
	DB.Table(DBVar.TableName.SubFolio).Select("account_code, id").Where("folio_number=? AND void=0", FolioNumber).Scan(&DataOutput)
	Result := true
	for _, data := range DataOutput {
		if !IsCanVoidOrCorrect(DB, Dataset, 1, data.Id, data.AccountCode) {
			Result = false
			break
		}
	}

	return Result
}

func CheckFolioHaveTransactionFromOtherFolio(DB *gorm.DB, FolioNumber uint64) (string, error) {
	var FolioHaveTransactionMessage string
	var FolioDetail []string
	err := DB.Table(DBVar.TableName.SubFolio).Select("CONCAT(sub_folio.belongs_to, '/Room: ', guest_detail.room_number, '/', contact_person.title_code, contact_person.full_name) AS FolioDetail").
		Joins("LEFT OUTER JOIN folio ON (sub_folio.belongs_to = folio.number)").
		Joins("LEFT OUTER JOIN contact_person ON (folio.contact_person_id1 = contact_person.id)").
		Joins("LEFT OUTER JOIN guest_detail ON (folio.guest_detail_id = guest_detail.id)").
		Where("sub_folio.folio_number=? AND sub_folio.folio_number<>sub_folio.belongs_to", FolioNumber).
		Group("sub_folio.belongs_to").
		Scan(&FolioDetail).Error
	if err != nil {
		return "", err
	}

	if len(FolioDetail) > 0 {
		FolioHaveTransactionMessage = "This folio have transaction from another folio(s)\nPlease return the transaction to original folio:\n"
		for _, detail := range FolioDetail {
			FolioHaveTransactionMessage += "- " + detail + "\n"
		}
	}
	return FolioHaveTransactionMessage, nil
}

func GetInvoiceNumberFromFolio(DB *gorm.DB, FolioNumber uint64) (InvoiceNumber string) {
	DB.Table(DBVar.TableName.InvoiceItem).Select("invoice_number").Where("folio_number=?", FolioNumber).Limit(1).Scan(&InvoiceNumber)
	return
}

func GetCompanyDetailListP(c *gin.Context) {
	var DataOutputCompany []map[string]interface{}

	// Get Program Configuration
	val, exist := c.Get("pConfig")
	if !exist {
		MasterData.SendResponse(global_var.ResponseCode.DatabaseError, "pConfig", nil, c)
		return
	}
	pConfig := val.(*config.CompanyDataConfiguration)
	DB := pConfig.DB
	DB.Table(DBVar.TableName.Company).Select("company.*").
		Joins("LEFT JOIN cfg_init_company_type ON company.type_code=cfg_init_company_type.code").
		// Where("company.is_direct_bill = '1'").
		Order("company.name").
		Find(&DataOutputCompany)

	MasterData.SendResponse(GlobalVar.ResponseCode.Successfully, "", DataOutputCompany, c)
}

func GetARCompaniesP(c *gin.Context) {
	var DataOutputCompany []map[string]interface{}

	// Get Program Configuration
	val, exist := c.Get("pConfig")
	if !exist {
		MasterData.SendResponse(global_var.ResponseCode.DatabaseError, "pConfig", nil, c)
		return
	}
	pConfig := val.(*config.CompanyDataConfiguration)
	DB := pConfig.DB
	DB.Table(DBVar.TableName.Company).Select("company.*").
		Joins("LEFT JOIN cfg_init_company_type ON company.type_code=cfg_init_company_type.code").
		Where("company.is_direct_bill = '1'").
		Order("company.name").
		Find(&DataOutputCompany)

	MasterData.SendResponse(GlobalVar.ResponseCode.Successfully, "", DataOutputCompany, c)
}

func IsDayendCloseRunning(c *gin.Context) (bool, error) {
	Status, err := GetDayendCloseStatus(c)

	return Status.Status, err
}

type DayendCloseStatus struct {
	Status    bool   `json:"status"`
	UserID    string `json:"user_id"`
	IPAddress string `json:"ip_address"`
	SessionID string `json:"session_id"`
}

func GetDayendCloseStatus(c *gin.Context) (DayendCloseStatus, error) {
	// Get Program Configuration
	// val, exist := c.Get("pConfig")
	// if !exist {
	// 	return DayendCloseStatus{}, errors.New("pConfig not Found")
	// }
	// pConfig := val.(*config.CompanyDataConfiguration)

	var Status DayendCloseStatus
	// pConfig.MxDayendClose.Lock()
	// defer pConfig.MxDayendClose.Unlock()

	b, _ := cache.DataCache.Get(c, c.GetString("UnitCode"), "DayendCloseStatus")
	if b != nil {
		err := json.Unmarshal(b, &Status)
		if err != nil {
			return Status, errors.New("Error on unmarshal")
		}
	}
	return Status, nil
}

func GetAPAROutstanding(DB *gorm.DB, APARNumber string) (Outstanding float64, err error) {
	if err := DB.Table(DBVar.TableName.AccApAr).Select("SUM(IFNULL(amount,0) - IFNULL(amount_paid,0)) AS Outstanding").
		Where("number=?", APARNumber).Limit(1).Scan(&Outstanding).Error; err != nil {
		return 0, err
	}
	return Outstanding, nil
}

func GetAPCommissionOutStanding(DB *gorm.DB, Dataset *global_var.TDataset, SubFolioID uint64, RefNumber string) (float64, error) {
	var Amount float64
	GAAPRefundDeposit := Dataset.GlobalAccount.APRefundDeposit
	GACreditCardAdm := Dataset.GlobalAccount.CreditCardAdm
	if err := DB.Table("(?) AS APNoShow",
		DB.Table(DBVar.TableName.SubFolio).Select(
			" sub_folio.id,"+
				" (SUM(IF(sub_folio.type_code='"+GlobalVar.TransactionType.Debit+"', sub_folio.quantity*sub_folio.amount, -(sub_folio.quantity*sub_folio.amount))) - IFNULL(Payment.TotalPaid,0)) AS Amount ").
			Joins(" LEFT OUTER JOIN ("+
				"SELECT"+
				" sub_folio_id,"+
				" SUM(amount) AS TotalPaid "+
				"FROM"+
				"  acc_ap_commission_payment_detail"+
				" WHERE sub_folio_id="+General.Uint64ToStr(SubFolioID)+
				" AND ref_number<>'"+RefNumber+"' "+
				"GROUP BY sub_folio_id) AS Payment ON (sub_folio.id = Payment.sub_folio_id) ").
			Joins(" LEFT OUTER JOIN cfg_init_account ON (sub_folio.account_code = cfg_init_account.code)").
			Where("cfg_init_account.sub_group_code=?", GlobalVar.GlobalAccountSubGroup.AccountPayable).
			Where("sub_folio.account_code<>?", GAAPRefundDeposit).
			Where("sub_folio.account_code<>?", GACreditCardAdm).
			Where("sub_folio.void='0'").
			Group("sub_folio.correction_breakdown, sub_folio.direct_bill_code")).
		Select("Amount").
		Where("id=?", SubFolioID).
		Scan(&Amount).Error; err != nil {
		return 0, err
	}
	return Amount, nil
}

func GetTimeSegmentCode() string {
	now := time.Now()
	Hour := now.UTC().Hour()

	if Hour > 5 && Hour < 12 {
		return GlobalVar.TimeSegment.Breakfast
	} else if Hour > 12 && Hour < 17 {
		return GlobalVar.TimeSegment.Lunch
	}
	return GlobalVar.TimeSegment.Dinner
}

func GetFirstMarketCodePOS(DB *gorm.DB) (string, error) {
	code := ""
	if err := DB.Table(DBVar.TableName.PosCfgInitMarket).Select("code").Order("id_sort").Limit(1).Scan(&code).Error; err != nil {
		return "", err
	}
	return code, nil
}

func GetJournalAccountPayable(c *gin.Context, DB *gorm.DB) ([]map[string]interface{}, error) {
	var DataOutput []map[string]interface{}
	if err := DB.Table(DBVar.TableName.CfgInitJournalAccount).Select(
		"cfg_init_journal_account.code",
		"cfg_init_journal_account.name",
		"cfg_init_journal_account_sub_group.name AS SubGroupName").
		Joins("LEFT OUTER JOIN cfg_init_journal_account_sub_group ON (cfg_init_journal_account.sub_group_code = cfg_init_journal_account_sub_group.code)").
		Where("cfg_init_journal_account_sub_group.group_code=?", 2).
		Order("cfg_init_journal_account.code").Scan(&DataOutput).Error; err != nil {
		return DataOutput, err
	}

	return DataOutput, nil
}

func GetJournalAccountIncome(c *gin.Context, DB *gorm.DB) ([]map[string]interface{}, error) {
	var DataOutput []map[string]interface{}
	if err := DB.Table(DBVar.TableName.CfgInitJournalAccount).Select(
		"cfg_init_journal_account.code",
		"cfg_init_journal_account.name",
		"cfg_init_journal_account_sub_group.name AS SubGroupName").
		Joins("LEFT OUTER JOIN cfg_init_journal_account_sub_group ON (cfg_init_journal_account.sub_group_code = cfg_init_journal_account_sub_group.code)").
		Where("cfg_init_journal_account_sub_group.group_code=?", 4).
		Or("cfg_init_journal_account_sub_group.group_code=?", 8).
		Order("cfg_init_journal_account.code").Scan(&DataOutput).Error; err != nil {
		return DataOutput, err
	}

	return DataOutput, nil
}

func GetJournalAccountExpense(c *gin.Context, DB *gorm.DB) ([]map[string]interface{}, error) {
	var DataOutput []map[string]interface{}
	if err := DB.Table(DBVar.TableName.CfgInitJournalAccount).Select(
		"cfg_init_journal_account.code",
		"cfg_init_journal_account.name",
		"cfg_init_journal_account_sub_group.name AS SubGroupName").
		Joins("LEFT OUTER JOIN cfg_init_journal_account_sub_group ON (cfg_init_journal_account.sub_group_code = cfg_init_journal_account_sub_group.code)").
		Where("cfg_init_journal_account_sub_group.group_code=?", 6).
		Or("cfg_init_journal_account_sub_group.group_code=?", 7).
		Or("cfg_init_journal_account_sub_group.group_code=?", 9).
		Order("cfg_init_journal_account.code").Scan(&DataOutput).Error; err != nil {
		return DataOutput, err
	}

	return DataOutput, nil
}

func GetJournalAccountCosting(c *gin.Context, DB *gorm.DB) ([]map[string]interface{}, error) {
	var DataOutput []map[string]interface{}
	if err := DB.Table(DBVar.TableName.CfgInitJournalAccount).Select(
		"cfg_init_journal_account.code",
		"cfg_init_journal_account.name",
		"cfg_init_journal_account_sub_group.name AS SubGroupName").
		Joins("LEFT OUTER JOIN cfg_init_journal_account_sub_group ON (cfg_init_journal_account.sub_group_code = cfg_init_journal_account_sub_group.code)").
		Where("cfg_init_journal_account_sub_group.group_code=?", 6).
		Or("cfg_init_journal_account_sub_group.group_code=?", 7).
		Or("cfg_init_journal_account_sub_group.group_code=?", 9).
		Order("cfg_init_journal_account.code").Scan(&DataOutput).Error; err != nil {
		return DataOutput, err
	}

	return DataOutput, nil
}

func GetBankAccountReceive(c *gin.Context, DB *gorm.DB) ([]map[string]interface{}, error) {
	var DataOutput []map[string]interface{}
	if err := DB.Table(DBVar.TableName.AccCfgInitBankAccount).Select(
		"code", "name", "journal_account_code", "type_code", "bank_account_number").
		Where("for_receive=1").
		Order("journal_account_code").Scan(&DataOutput).Error; err != nil {
		return DataOutput, err
	}

	return DataOutput, nil
}

func GetBankAccountPayment(c *gin.Context, DB *gorm.DB) ([]map[string]interface{}, error) {
	var DataOutput []map[string]interface{}
	if err := DB.Table(DBVar.TableName.AccCfgInitBankAccount).Select(
		"code", "name", "journal_account_code", "type_code", "bank_account_number").
		Where("for_payment=1").
		Order("journal_account_code").Scan(&DataOutput).Error; err != nil {
		return DataOutput, err
	}

	return DataOutput, nil
}

func GetCompany(c *gin.Context, DB *gorm.DB) ([]map[string]interface{}, error) {
	var DataOutput []map[string]interface{}
	if err := DB.Table(DBVar.TableName.Company).Select(
		"code", "name").
		Order("name").Scan(&DataOutput).Error; err != nil {
		return DataOutput, err
	}

	return DataOutput, nil
}

func IsAPPaid(c *gin.Context, DB *gorm.DB, Number string) (IsPaid uint8, err error) {
	if err := DB.Table(DBVar.TableName.AccApAr).
		Select("is_paid").Where("number=?", Number).
		Limit(1).
		Scan(&IsPaid).Error; err != nil {
		return 0, err
	}

	return
}

func GetReceiveIDLastPrice(DB *gorm.DB, ReceiveID uint64, Quantity float64) float64 {
	type Struct struct {
		BasicQuantity, Quantity, TotalPrice float64
	}
	var DataOutput Struct
	DB.Table(DBVar.TableName.InvReceivingDetail).Select("basic_quantity, quantity, total_price").Where("id=?", ReceiveID).Limit(1).Scan(&DataOutput)
	if Quantity == DataOutput.Quantity {
		if DataOutput.BasicQuantity > DataOutput.Quantity {
			return General.RoundToX3(DataOutput.TotalPrice - ((DataOutput.TotalPrice / DataOutput.BasicQuantity) * (DataOutput.BasicQuantity - DataOutput.Quantity)))
		} else {
			return General.RoundToX3(DataOutput.TotalPrice)
		}
	} else {
		return General.RoundToX3((DataOutput.TotalPrice / DataOutput.BasicQuantity)) * Quantity
	}
}

func IsFAReceiveUsed(DB *gorm.DB, ReceiveNumber string) (bool, error) {
	var Number string
	if err := DB.Table(DBVar.TableName.FaRevaluation).Select("number").Joins("LEFT OUTER JOIN fa_list ON (fa_revaluation.fa_code = fa_list.code)").Where("fa_list.receive_number=?", ReceiveNumber).Limit(1).Scan(&Number).Error; err != nil {
		return false, err
	}
	var Code string
	if err := DB.Table(DBVar.TableName.FaList).Select("code").Where("receive_number=?", ReceiveNumber).Where("condition_code<>?", global_var.FAItemCondition.Good).Limit(1).Scan(&Code).Error; err != nil {
		return false, err
	}
	return Number != "" || Code != "", nil

}

func IsStockMinusReceive(DB *gorm.DB, Dataset *global_var.TDataset, StoreCode, ItemCode, ReceiveNumber string, StockDate time.Time) bool {
	var Date []time.Time
	DB.Raw(
		"SELECT DISTINCT Stock.`date` FROM (("+
			"SELECT DISTINCT"+
			" inv_costing.`date` "+
			"FROM"+
			" inv_costing_detail"+
			" LEFT OUTER JOIN inv_costing ON (inv_costing_detail.costing_number = inv_costing.number)"+
			" WHERE inv_costing_detail.item_code=? "+
			" AND inv_costing_detail.store_code=? "+
			" AND inv_costing.date>='"+General.FormatDate1(StockDate)+"'"+
			") UNION ALL ("+
			"SELECT DISTINCT"+
			" inv_stock_transfer.`date` "+
			"FROM"+
			" inv_stock_transfer_detail"+
			" LEFT OUTER JOIN inv_stock_transfer ON (inv_stock_transfer_detail.st_number = inv_stock_transfer.number)"+
			" WHERE inv_stock_transfer_detail.item_code=? "+
			" AND inv_stock_transfer_detail.from_store_code=? "+
			" AND inv_stock_transfer.`date` >= '"+General.FormatDate1(StockDate)+"')) AS Stock "+
			"ORDER BY Stock.`date`", ItemCode, StoreCode, ItemCode, StoreCode).Scan(&Date)

	for _, date := range Date {
		if GetStockStoreUpdateReceive(DB, StoreCode, ItemCode, ReceiveNumber, Dataset.ProgramConfiguration.CostingMethod, date) < 0 {
			return true
		}
	}
	return false
}

func IsReceiveIDUsedInStockTransfer(DB *gorm.DB, ReceiveIDx int64) bool {
	var ReceiveID uint64
	DB.Table(DBVar.TableName.InvStockTransferDetail).Select("receive_id").Where("receive_id IN (?)", ReceiveIDx).Limit(1).Scan(&ReceiveID)

	return ReceiveID > 0
}

func IsReceiveIDUsedInCosting(DB *gorm.DB, ReceiveIDx int64) bool {
	var ReceiveID uint64
	DB.Table(DBVar.TableName.InvCostingDetail).Select("receive_id").Where("receive_id IN (?)", ReceiveIDx).Limit(1).Scan(&ReceiveID)

	return ReceiveID > 0
}

func IsFAPurchaseOrderUsedInReceive(DB *gorm.DB, PONumberx string) bool {
	var PONumber uint64
	DB.Table(DBVar.TableName.FaReceive).Select("po_number").Where("po_number IN (?)", PONumberx).Limit(1).Scan(&PONumber)

	return PONumber > 0
}

func GetInventoryCOGSJournalAccount(DB *gorm.DB, ItemCode, SubDepartmentCode string) string {
	var JournalAccountCode string
	DB.Raw(
		"SELECT"+
			" inv_cfg_init_item_category_other_cogs.journal_account_code "+
			"FROM"+
			" inv_cfg_init_item"+
			" LEFT OUTER JOIN inv_cfg_init_item_category_other_cogs ON (inv_cfg_init_item.category_code = inv_cfg_init_item_category_other_cogs.category_code)"+
			" WHERE inv_cfg_init_item.code=? "+
			" AND inv_cfg_init_item_category_other_cogs.sub_department_code=? ", ItemCode, SubDepartmentCode).Limit(1).Scan(&JournalAccountCode)
	var JournalAccountCodeCOGS string
	if JournalAccountCode == "" {
		DB.Raw(
			"SELECT"+
				" inv_cfg_init_item_category.journal_account_code_cogs "+
				"FROM"+
				" inv_cfg_init_item"+
				" LEFT OUTER JOIN inv_cfg_init_item_category ON (inv_cfg_init_item.category_code = inv_cfg_init_item_category.code)"+
				" WHERE inv_cfg_init_item.code=? ", ItemCode).Limit(1).Scan(&JournalAccountCodeCOGS)
		return JournalAccountCodeCOGS
	}
	return JournalAccountCode
}

func GetInventoryExpenseJournalAccount(DB *gorm.DB, ItemCode, SubDepartmentCode string) string {
	var JournalAccountCode string
	DB.Raw(
		"SELECT"+
			" inv_cfg_init_item_category_other_expense.journal_account_code "+
			"FROM"+
			" inv_cfg_init_item"+
			" LEFT OUTER JOIN inv_cfg_init_item_category_other_expense ON (inv_cfg_init_item.category_code = inv_cfg_init_item_category_other_expense.category_code)"+
			" WHERE inv_cfg_init_item.code=? "+
			" AND inv_cfg_init_item_category_other_expense.sub_department_code=? ", ItemCode, SubDepartmentCode).Limit(1).Scan(&JournalAccountCode)
	var JournalAccountCodeExpense string
	if JournalAccountCode == "" {
		DB.Raw(
			"SELECT"+
				" inv_cfg_init_item_category.journal_account_code_expense "+
				"FROM"+
				" inv_cfg_init_item"+
				" LEFT OUTER JOIN inv_cfg_init_item_category ON (inv_cfg_init_item.category_code = inv_cfg_init_item_category.code)"+
				" WHERE inv_cfg_init_item.code=? ", ItemCode).Limit(1).Scan(&JournalAccountCodeExpense)
		return JournalAccountCodeExpense
	}
	return JournalAccountCode
}

func GetInventoryItemGroup(DB *gorm.DB, ItemCode string) string {
	var GroupCode string
	DB.Raw(
		"SELECT"+
			" inv_cfg_init_item_category.group_code "+
			"FROM"+
			" inv_cfg_init_item"+
			" LEFT OUTER JOIN inv_cfg_init_item_category ON (inv_cfg_init_item.category_code = inv_cfg_init_item_category.code)"+
			" WHERE inv_cfg_init_item.code=?", ItemCode).Limit(1).Scan(&GroupCode)
	return GroupCode
}

func IsJournalAlreadyImported(DB *gorm.DB, AuditDate time.Time) (bool, error) {
	var Id uint64
	if err := DB.Table("acc_import_journal_log").Select("id").Where("audit_date=?", General.FormatDate1(AuditDate)).Limit(1).Scan(&Id).Error; err != nil {
		return true, err
	}
	return Id > 0, nil
}

func IsMonthClosed(DB *gorm.DB, Month, Year int) (bool, error) {
	var Id uint64
	if err := DB.Table(DBVar.TableName.AccCloseMonth).Select("id").
		Where("month=?", Month).Where("year=?", Year).Limit(1).Scan(&Id).Error; err != nil {
		return true, err
	}
	return Id > 0, nil
}

func IsPurchaseOrderClosed(DB *gorm.DB, PONumberX string) (bool, error) {
	var PONumber string
	if err := DB.Table(DBVar.TableName.InvPurchaseOrderDetail).Select("po_number").Where("po_number=?", PONumberX).
		Where("quantity_not_received>0").
		Limit(1).
		Scan(&PONumber).Error; err != nil {
		return false, err
	}
	return PONumber != "", nil
}

func GetItemBasicUOMCode(DB *gorm.DB, ItemCode string) (string, error) {
	var DataOutput DBVar.Inv_cfg_init_item
	if err := DB.Table(DBVar.TableName.InvCfgInitItem).Select("uom_code").Where("code=?", ItemCode).Limit(1).Scan(&DataOutput).Error; err != nil {
		return "", err
	}

	return DataOutput.UomCode, nil
}

func GetItemMultiUOM(DB *gorm.DB, ItemCode, UomCode string) (DBVar.Inv_cfg_init_item_uom, error) {
	var DataOutput DBVar.Inv_cfg_init_item_uom
	if err := DB.Table(DBVar.TableName.InvCfgInitItemUom).Where("item_code=?", ItemCode).Where("item_code=?", ItemCode).Limit(1).Scan(&DataOutput).Error; err != nil {
		return DBVar.Inv_cfg_init_item_uom{}, err
	}

	return DataOutput, nil
}

func IsAPARHadPayment(DB *gorm.DB, APARNumber string) bool {
	var Number string
	DB.Raw(
		"SELECT ap_ar_number FROM acc_ap_ar_payment_detail"+
			" WHERE ap_ar_number=? "+
			"LIMIT 1", APARNumber).Scan(&Number)

	return Number != ""
}

func CheckDateRangeScheduledRate(DB *gorm.DB, FolioReservationNumber, ID uint64, FromDate, ToDate, ArrivalDate, DepartureDate time.Time, IsReservation bool) (bool, error) {
	ConditionNumber := "folio_number"
	TableName := DBVar.TableName.GuestScheduledRate
	if IsReservation {
		ConditionNumber = "reservation_number"
		TableName = DBVar.TableName.ReservationScheduledRate
	}

	query := DB.Table(TableName).Select("id").
		Where(ConditionNumber, FolioReservationNumber).
		Where(DB.Where("from_date<=?", General.FormatDate1(FromDate)).Where("to_date", General.FormatDate1(FromDate)).
			Or("from_date<=?", General.FormatDate1(ToDate)).Where("to_date>=?", General.FormatDate1(ToDate)))

	if ID > 0 {
		query.Where("id<>?", ID)
	}
	var Id uint64
	if err := query.Limit(1).Scan(&Id).Error; err != nil {
		return false, err
	}

	return Id > 0, nil
}

func GetRoomRateWeekdayRate1(DB *gorm.DB, RoomRateCode string) (WeekdayRate1 float64, err error) {
	if err := DB.Table(DBVar.TableName.CfgInitRoomRate).Select("weekday_rate1").Where("code", RoomRateCode).Limit(1).Scan(&WeekdayRate1).Error; err != nil {
		return 0, err
	}

	return WeekdayRate1, nil
}

func GenerateVoucher(c *gin.Context, DB *gorm.DB, Dataset *global_var.TDataset,
	StartDate, ExpireDate time.Time,
	TitleCode, FullName, Street, CompanyCode, TypeCode, ReasonCode, RequestByCode, AccomodationType, SetUpRequired, Description, RoomTypeCode string,
	Point, Price float64,
	Nights int, UserID string) (string, error) {
	result := ""

	prefix := "V" + TypeCode + StartDate.Format("06") + "-"
	if Price == 0 {
		Price = Dataset.ProgramConfiguration.VoucherDefaultPrice
	}

	rand.Seed(time.Now().UnixNano())

	for {
		var maxNumber int64
		if err := DB.Raw("SELECT CAST(RIGHT(number, LENGTH(number)-?) AS UNSIGNED) AS MaxNumber FROM voucher WHERE LEFT(number, ?) = ? ORDER BY MaxNumber DESC LIMIT 1",
			len(prefix), len(prefix), prefix).Scan(&maxNumber).Error; err != nil {
			return "", err
		}

		voucherNumber := ""
		code := ""

		if maxNumber == 0 {
			voucherNumber = "1"
		} else {
			voucherNumber = strconv.FormatInt(maxNumber+1, 10)
		}

		if len(voucherNumber) < Dataset.ProgramConfiguration.VoucherLength {
			voucherNumber = fmt.Sprintf("%0*d", Dataset.ProgramConfiguration.VoucherLength-len(voucherNumber), 0) + voucherNumber
		}

		voucherNumber = prefix + voucherNumber

		for i := 1; i <= 4; i++ {
			code += strconv.Itoa(rand.Intn(10))
		}

		var count int64
		if err := DB.Model(&DBVar.Voucher{}).Where("number = ?", voucherNumber).Count(&count).Error; err != nil {
			return "", err
		}
		if count == 0 {
			if err := DB.Create(&DBVar.Voucher{
				Number:            voucherNumber,
				Code:              code,
				TitleCode:         TitleCode,
				FullName:          FullName,
				Street:            Street,
				MemberTypeCode:    GlobalVar.MemberType.Room,
				TypeCode:          TypeCode,
				CompanyCode:       CompanyCode,
				ReasonCode:        ReasonCode,
				RequestByCode:     RequestByCode,
				StartDate:         StartDate,
				ExpireDate:        ExpireDate,
				AccommodationType: AccomodationType,
				SetUpRequired:     SetUpRequired,
				Description:       Description,
				Point:             Point,
				StatusCodeApprove: GlobalVar.VoucherStatusApprove.Unapprove,
				StatusCodeSold:    "",
				IsRoomChargeOnly:  1,
				StatusCode:        GlobalVar.VoucherStatus.Active,
				Price:             Price,
				Nights:            Nights,
				RoomTypeCode:      RoomTypeCode,
				IssuedDate:        time.Now(),
				CreatedBy:         UserID,
			}).Error; err != nil {
				return "", err
			}
			result = voucherNumber
			break
		}
	}
	return result, nil
}

func IsVoucherDateCanSoldRedeemCompliment(DB *gorm.DB, VoucherNumber string, AuditDate time.Time) (bool, error) {
	var ExpireDate time.Time
	if err := DB.Table(DBVar.TableName.Voucher).Select("expire_date").Where("number=?", VoucherNumber).Limit(1).Scan(&ExpireDate).Error; err != nil {
		return false, err
	}

	return ExpireDate.After(AuditDate), nil

}

func GetTotalMemberPointRedeem(DB *gorm.DB, MemberCode, MemberTypeCode string) (float64, error) {
	var TotalPoint float64
	if err := DB.Table(DBVar.TableName.MemberPoint).Select("IFNULL(SUM(point),0) AS TotalPoint").
		Where("member_code", MemberCode).
		Where("member_type_code", MemberTypeCode).
		Where("is_redeemed = 1 AND is_expired=0").
		Scan(&TotalPoint).Error; err != nil {
		return 0, err
	}
	return TotalPoint, nil
}

func GetTotalMemberPointRedeemActual(DB *gorm.DB, MemberCode, MemberTypeCode string) (float64, error) {
	var TotalPoint float64
	if err := DB.Table(DBVar.TableName.MemberPointRedeem).Select("IFNULL(SUM(point),0) AS TotalPoint").
		Where("member_code", MemberCode).
		Where("member_type_code", MemberTypeCode).
		Scan(&TotalPoint).Error; err != nil {
		return 0, err
	}
	return TotalPoint, nil
}

func SetAuditDate(c *gin.Context, DB *gorm.DB, Date time.Time, UserID string) error {
	if err := DB.Table(DBVar.TableName.AuditLog).Create(map[string]interface{}{
		"audit_date": Date,
		"created_by": UserID}).Error; err != nil {
		return err
	}
	newAuditDate := GetAuditDate(c, DB, true)
	websocket.SendMessage(c.GetString("UnitCode"), GlobalVar.WSMessageType.Client, nil, GlobalVar.WSDataType.AuditDateChanged, newAuditDate, UserID)
	return nil
}

func CheckDepreciation(DB *gorm.DB, Month, Year int) (bool, error) {
	var ID uint64
	if err := DB.Table(db_var.TableName.FaDepreciation).Select("id").Where("month", Month).Where("year", Year).Limit(1).Scan(&ID).Error; err != nil {
		return true, err
	}

	return ID > 0, nil
}

// func IsMonthClosed(DB *gorm.DB, Month, Year string) (bool, error){
// 	var ID uint64
// 	if err := DB.Table(db_var.TableName.AccCloseMonth).Select("id").Where("month", Month).Where("year", Year).Limit(1).Scan(&ID).Error;err!= nil {
// 		return true,err
// 	}

// 	return ID > 0, nil
// }

// func IsYearClosed(DB *gorm.DB, Year string) (bool, error){
// 	var ID uint64
// 	if err := DB.Table(db_var.TableName.AccCloseYear).Select("id").Where("year", Year).Limit(1).Scan(&ID).Error;err!= nil {
// 		return true,err
// 	}

// 	return ID > 0, nil
// }

func IsJournalExported(DB *gorm.DB, MonthX, YearX string) (bool, error) {
	var CountExported float64
	if err := DB.Raw(
		"SELECT" +
			" (IFNULL(COUNT(TransactionX.audit_date), 0) - IFNULL(COUNT(ExportedJournal.audit_date), 0)) AS CountExported " +
			"FROM" +
			" audit_log" +
			" LEFT OUTER JOIN (" +
			"SELECT DISTINCT audit_date FROM ((" +
			"SELECT DISTINCT audit_date FROM guest_deposit WHERE MONTH(guest_deposit.audit_date)='" + MonthX + "' AND YEAR(guest_deposit.audit_date)='" + YearX + "' AND void='0') UNION ALL (" +
			"SELECT DISTINCT audit_date FROM sub_folio WHERE MONTH(sub_folio.audit_date)='" + MonthX + "' AND YEAR(sub_folio.audit_date)='" + YearX + "' AND void='0')) AS Transaction)" +
			" AS TransactionX ON (audit_log.audit_date = TransactionX.audit_date)" +
			" LEFT OUTER JOIN (SELECT DISTINCT audit_date FROM acc_import_journal_log WHERE MONTH(acc_import_journal_log.audit_date)='" + MonthX + "' AND YEAR(acc_import_journal_log.audit_date)='" + YearX + "') AS ExportedJournal ON (audit_log.audit_date = ExportedJournal.audit_date)" +
			" WHERE MONTH(audit_log.audit_date)='" + MonthX + "'" +
			" AND YEAR(audit_log.audit_date)='" + YearX + "' " +
			"GROUP BY MONTH(audit_log.audit_date), YEAR(audit_log.audit_date)").Scan(&CountExported).Error; err != nil {
		return false, err
	}

	return CountExported <= 0, nil
}
