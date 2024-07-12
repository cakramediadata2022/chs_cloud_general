package global_query

import (
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

	"github.com/cakramediadata2022/chs_cloud_general/db_var"
	"github.com/cakramediadata2022/chs_cloud_general/general"
	"github.com/cakramediadata2022/chs_cloud_general/global_var"
	"github.com/cakramediadata2022/chs_cloud_general/internal/config"
	"github.com/cakramediadata2022/chs_cloud_general/internal/master_data"
	"github.com/cakramediadata2022/chs_cloud_general/internal/utils/cache"
	"github.com/cakramediadata2022/chs_cloud_general/internal/utils/websocket"
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
			return global_var.SigningKey, nil
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
	// Timezone := pConfig.Dataset.ProgramConfiguration.Timezone
	DateX, err := cache.DataCache.GetString(c, CompanyCode, "AUDIT_DATE")
	if err != nil || Reload {
		var Date time.Time
		DB.Table(db_var.TableName.AuditLog).Select("audit_date").Order("id DESC").Limit(1).Scan(&Date)

		cache.DataCache.Set(c, CompanyCode, "AUDIT_DATE", Date, 6*time.Hour)
		return Date
	}

	return general.StrZToDate(DateX)

}

func GetAuditDateP(c *gin.Context) {
	// Get Program Configuration
	val, exist := c.Get("pConfig")
	if !exist {
		master_data.SendResponse(global_var.ResponseCode.DatabaseError, "pConfig", nil, c)
		return
	}
	pConfig := val.(*config.CompanyDataConfiguration)
	DB := pConfig.DB
	AuditDate := GetAuditDate(c, DB, true)
	master_data.SendResponse(global_var.ResponseCode.Successfully, "", AuditDate, c)
}

func GetServerDateTime(c *gin.Context, DB *gorm.DB) time.Time {
	return master_data.GetFieldTimeQuery(DB, "SELECT NOW() AS DateServer;")
}

func GetServerDateTimeP(c *gin.Context) {
	// Get Program Configuration
	val, exist := c.Get("pConfig")
	if !exist {
		master_data.SendResponse(global_var.ResponseCode.DatabaseError, "pConfig", nil, c)
		return
	}
	pConfig := val.(*config.CompanyDataConfiguration)
	DB := pConfig.DB
	master_data.SendResponse(global_var.ResponseCode.Successfully, "", GetServerDateTime(c, DB), c)
}

func GetServerDate(c *gin.Context, DB *gorm.DB) time.Time {
	return master_data.GetFieldTimeQuery(DB, "SELECT DATE(NOW()) AS DateServer;")
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
		master_data.SendResponse(global_var.ResponseCode.DatabaseError, "pConfig", nil, c)
		return
	}
	pConfig := val.(*config.CompanyDataConfiguration)
	DB := pConfig.DB
	master_data.SendResponse(global_var.ResponseCode.Successfully, "", GetServerDate(c, DB), c)
}

func IsDiscountLimit(c *gin.Context, DB *gorm.DB, UserID string, RateOriginal, RateOverride float64) bool {
	var DiscountPercent float64
	DiscountAmount := RateOriginal - RateOverride
	if RateOriginal > 0 {
		DiscountPercent = (DiscountAmount) / RateOriginal * 100
	}
	var Code string
	if err := DB.Table(db_var.TableName.User).Select("IFNULL(user_group_access.code,'') AS Code").
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
	DB.Table(db_var.TableName.Reservation).Select("status_code").Where("number = ?", ReservationNumber).Take(&Status)

	return Status
}

func GetAvailableRoomCountByType(DB *gorm.DB, Dataset *global_var.TDataset, ArrivalDate, DepartureDate time.Time, RoomTypeCode, BedTypeCode string, ReservationNumber, FolioNumber, RoomUnavailableID, RoomAllotmentID uint64, ReadyOnly, AllotmentOnly bool) (int64, error) {
	var Result int64 = 0
	Timezone := Dataset.ProgramConfiguration.Timezone
	loc, _ := time.LoadLocation(Timezone)
	ArrivalDate = ArrivalDate.In(loc)
	DepartureDate = DepartureDate.In(loc)

	fmt.Println("ArrivalDate", ArrivalDate)
	fmt.Println("DepartureDate", DepartureDate)

	ArrivalDateStr := ArrivalDate.Format("2006-01-02")
	DepartureDateStr := DepartureDate.Format("2006-01-02")

	FieldCountRoomTotal := ""
	FieldCountRoom := ""
	FieldCountReservation := ""
	FieldCountFolio := ""
	FieldCountUnavailable := ""
	FieldCountAllotment := ""
	CountDay := general.DaysBetween(ArrivalDate, DepartureDate)
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
			CurrentDateStr := general.FormatDate1(general.IncDay(ArrivalDate, Count-1))
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
	QueryCondition1 = QueryCondition1 + " AND reservation.is_from_allotment='" + general.BoolToUint8String(AllotmentOnly) + "'"
	QueryNot := ""
	QueryRoom := ""
	QueryRoomUnavailable := ""
	QueryRoomAllotment := ""
	//TODO optimize query, not using prepare statement
	if AllotmentOnly {
		QueryCondition2A = QueryCondition2A + " AND folio.is_from_allotment='" + general.BoolToUint8String(AllotmentOnly) + "'"
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
			" WHERE reservation.status_code='" + global_var.ReservationStatus.New + "'" +
			" AND guest_detail.room_type_code='" + RoomTypeCode + "'" +
			QueryCondition4x +
			" AND DATE(convert_tz(guest_detail.arrival,'UTC','" + Timezone + "')) <'" + DepartureDateStr + "'" +
			" AND DATE(convert_tz(guest_detail.departure,'UTC','" + Timezone + "')) >'" + ArrivalDateStr + "'" +
			QueryCondition1 +
			")UNION(" +
			"SELECT" +
			" 'C' AS Code," +
			FieldCountFolio +
			"FROM" +
			" folio" +
			" LEFT OUTER JOIN guest_detail ON (folio.guest_detail_id = guest_detail.id)" +
			" LEFT OUTER JOIN cfg_init_room ON (guest_detail.room_number = cfg_init_room.number)" +
			" WHERE folio.status_code='" + global_var.FolioStatus.Open + "'" +
			" AND folio.type_code='" + global_var.FolioType.GuestFolio + "'" +
			" AND guest_detail.room_type_code='" + RoomTypeCode + "'" +
			QueryCondition4 +
			" AND DATE(convert_tz(guest_detail.arrival,'UTC','" + Timezone + "')) <'" + DepartureDateStr + "'" +
			" AND DATE(convert_tz(guest_detail.departure,'UTC','" + Timezone + "')) >'" + ArrivalDateStr + "'" +
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
	AvailableRoomCount = general.StrToInt64(AvailableRoomCountArray[0]["AvailableRoomCount1"].(string))
	if CountDay > 1 {
		for _, roomCount := range AvailableRoomCountArray {
			for count := range roomCount {
				var countX = general.StrToInt64(roomCount[count].(string))
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
				" AND status_code='" + global_var.RoomStatus.Ready + "'" +
				QueryCondition4 + " " +
				"AND number NOT IN(SELECT" +
				" guest_detail.room_number " +
				"FROM" +
				" reservation" +
				" LEFT OUTER JOIN guest_detail ON (reservation.guest_detail_id = guest_detail.id)" +
				" WHERE reservation.status_code='" + global_var.ReservationStatus.New + "'" +
				" AND guest_detail.room_type_code='" + RoomTypeCode + "'" +
				" AND guest_detail.room_number<>''" +
				" AND DATE(convert_tz(guest_detail.arrival,'UTC','" + Timezone + "')) <'" + DepartureDateStr + "'" +
				" AND DATE(convert_tz(guest_detail.departure,'UTC','" + Timezone + "')) >'" + ArrivalDateStr + "'" +
				QueryCondition1 + ") " +
				"AND number NOT IN(SELECT" +
				" guest_detail.room_number " +
				"FROM" +
				" folio" +
				" LEFT OUTER JOIN guest_detail ON (folio.guest_detail_id = guest_detail.id)" +
				" WHERE folio.status_code='" + global_var.FolioStatus.Open + "'" +
				" AND folio.type_code='" + global_var.FolioType.GuestFolio + "'" +
				" AND guest_detail.room_type_code='" + RoomTypeCode + "'" +
				" AND DATE(convert_tz(guest_detail.arrival,'UTC','" + Timezone + "')) <'" + DepartureDateStr + "'" +
				" AND DATE(convert_tz(guest_detail.departure,'UTC','" + Timezone + "')) >'" + ArrivalDateStr + "'" +
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

func GetAvailableRoomByType(DB *gorm.DB, Dataset *global_var.TDataset, ArrivalDate, DepartureDate time.Time, RoomTypeCode, BedTypeCode string, ReservationNumber, FolioNumber, RoomUnavailableID, RoomAllotmentID uint64, ReadyOnly, AllotmentOnly bool) []RoomNumberListStruct {
	Timezone := Dataset.ProgramConfiguration.Timezone
	loc, _ := time.LoadLocation(Timezone)
	ArrivalDateStr := ArrivalDate.In(loc).Format("2006-01-02")
	DepartureDateStr := DepartureDate.In(loc).Format("2006-01-02")

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
		QueryCondition5 = "status_code='" + global_var.RoomStatus.Ready + "' AND "
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
			" WHERE reservation.status_code='" + global_var.ReservationStatus.New + "'" +
			QueryConditionRoomType2 +
			" AND guest_detail.room_number<>''" +
			" AND DATE(convert_tz(guest_detail.arrival,'UTC','" + Timezone + "')) <'" + DepartureDateStr + "'" +
			" AND DATE(convert_tz(guest_detail.departure,'UTC','" + Timezone + "')) >'" + ArrivalDateStr + "'" +
			QueryCondition1 + ") " +
			"AND number NOT IN(SELECT" +
			" guest_detail.room_number " +
			"FROM" +
			" folio" +
			" LEFT OUTER JOIN guest_detail ON (folio.guest_detail_id = guest_detail.id)" +
			" WHERE folio.status_code='" + global_var.FolioStatus.Open + "'" +
			" AND folio.type_code='" + global_var.FolioType.GuestFolio + "'" +
			QueryConditionRoomType2 +
			" AND DATE(convert_tz(guest_detail.arrival,'UTC','" + Timezone + "'))<'" + DepartureDateStr + "'" +
			" AND DATE(convert_tz(guest_detail.departure,'UTC','" + Timezone + "'))>'" + ArrivalDateStr + "'" +
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
	ArrivalDateStr := general.FormatDate1(general.StrZToDate(ArrivalDate))
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
					////fmt.Println("3")
					QueryConditionMarket = " AND (IFNULL(cfg_init_room_rate.market_code, '')='" + MarketCode + "' OR IFNULL(cfg_init_room_rate.market_code, '')='')"
				} else {
					////fmt.Println("14")
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
					////fmt.Println("1")
					QueryConditionCompany = " AND (cfg_init_room_rate.company_code='" + CompanyCode + "' OR cfg_init_room_rate.company_code='')"
				} else {
					////fmt.Println("13")
					QueryConditionCompany = " AND cfg_init_room_rate.company_code='" + CompanyCode + "'"
				}
			}
		} else {
			if BusinessSourceCode == "" {
				QueryConditionCompany = " AND cfg_init_room_rate.company_code=''"
			} else {
				if Dataset.ProgramConfiguration.AlwaysShowPublishRate {
					////fmt.Println("2")
					QueryConditionCompany = " OR (cfg_init_room_rate.company_code='" + BusinessSourceCode + "' OR cfg_init_room_rate.company_code='')"
				} else {
					////fmt.Println("12")
					QueryConditionCompany = " OR cfg_init_room_rate.company_code='" + BusinessSourceCode + "'"
				}
			}
		}

		QueryConditionBusinessSource := ""
		if BusinessSourceCode == "" {
			QueryConditionBusinessSource = "IFNULL(cfg_init_room_rate_business_source.company_code,'')=''"
		} else {
			if Dataset.ProgramConfiguration.AlwaysShowPublishRate {
				////fmt.Println("1")
				QueryConditionBusinessSource = " (cfg_init_room_rate_business_source.company_code='" + BusinessSourceCode + "' OR IFNULL(cfg_init_room_rate_business_source.company_code,'')='')"
			} else {
				////fmt.Println("11")
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
				"ORDER BY cfg_init_room_rate.id_sort , cfg_init_room_rate.name;").Scan(&DataArray)
	}

	return DataArray
}

func GetRoomRateAmount(ctx context.Context, DB *gorm.DB, Dataset *global_var.TDataset, RoomRateCode, PostingDateStr string, Adult, Child int, IsWeekend bool) float64 {
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
			IsWeekend = general.IsWeekend(PostingDate, Dataset)
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

func CopyGuestProfileToContactPerson(GuestProfileData db_var.Guest_profile) db_var.Contact_person {
	var ContactPersonData db_var.Contact_person
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

func SaveGuestProfileContactPerson(Db *gorm.DB, ValidUserCode string, GuestProfileData db_var.Guest_profile, ContactPersonId, GuestProfileId uint64) (uint64, uint64, error) {
	GuestProfileData.TypeCode = global_var.CPType.Guest
	GuestProfileData.UpdatedBy = ValidUserCode
	GuestProfileData.CreatedBy = ValidUserCode
	GuestProfileData.IsActive = 1

	PhoneNumber := GuestProfileData.Phone1
	if PhoneNumber != "" {
		if PhoneNumber[:1] == "0" {
			GuestProfileData.Phone1 = "+62" + PhoneNumber[1:]
		}
	}
	////fmt.Println("GuestProfileId", GuestProfileId)
	////fmt.Println("GuestProfileData", GuestProfileData.Id)
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

		err := Db.Table(db_var.TableName.GuestProfile).Omit(OmitX, OmitSource).Save(&GuestProfileData).Error
		if err != nil {
			return 0, 0, err
		}
		err = Db.Table(db_var.TableName.ContactPerson).Omit(OmitY).Save(&ContactPersonData).Error
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
// 		master_data.SendResponse(global_var.ResponseCode.InvalidDataFormat, "", nil, c)
// 	} else {
// 		var FolioTransfer []string
// 		var IsFolioAutoTransferAccount bool
// 		DB.Table(db_var.TableName.FolioRouting).Select("folio_transfer").Where("folio_number = ? AND account_code = ?", DataInput.FolioNumber, DataInput.AccountCode).Find(&FolioTransfer)
// 		IsFolioAutoTransferAccount = len(FolioTransfer) > 0
// 		master_data.SendResponse(global_var.ResponseCode.Successfully, "", IsFolioAutoTransferAccount, c)
// 	}
// }

func GetDefaultCurrencyCode(DB *gorm.DB) string {
	Result := ""
	var Code []string
	DB.Table(db_var.TableName.CfgInitCurrency).Select("code").Where("is_default = ?", 1).Find(&Code)

	if len(Code) > 0 {
		Result = Code[0]
		return Result
	}

	return Result
}

func GetExchangeRateCurrency(DB *gorm.DB, CurrencyCode string) float64 {
	return master_data.GetFieldFloat(DB, db_var.TableName.CfgInitCurrency, "exchange_rate", "code", CurrencyCode, 1)
}

func GetGuestDepositCorrectionBreakDown(DB *gorm.DB) uint64 {
	var CorrectionBreakdown uint64
	DB.Table(db_var.TableName.GuestDeposit).Select("correction_breakdown").Order("correction_breakdown desc").Limit(1).Scan(&CorrectionBreakdown)

	return CorrectionBreakdown + 1
}

func GetAccountSubGroupCode(DB *gorm.DB, AccountCode string) string {
	return master_data.GetFieldString(DB, db_var.TableName.CfgInitAccount, "sub_group_code", "code", AccountCode, "")
}

func GetTotalDepositReservation(DB *gorm.DB, ReservationNumber uint64, SystemCode string) float64 {
	return master_data.GetFieldFloatQuery(DB,
		"SELECT"+
			" SUM(IF(type_code='C', amount, -amount)) AS TotalDeposit "+
			"FROM"+
			" guest_deposit"+
			" WHERE reservation_number=?"+
			" AND void='0'"+
			" AND system_code='"+SystemCode+"' "+
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
		RoomNumber = master_data.GetFieldString(DB, db_var.TableName.CfgInitRoom, "number", "name", RoomNumber, "")
	}
	// if master_data.GetConfigurationBool(DB, global_var.SystemCode.Hotel, global_var.ConfigurationCategory.General, global_var.ConfigurationName.IsRoomByName, false) {
	// 	RoomNumber = master_data.GetFieldString(db_var.TableName.CfgInitRoom, "number", "name", RoomNumber, "")
	// }

	DB.Table(db_var.TableName.GuestDetail).Select("room_number", "updated_by").Where("id = ?", GuestDetailID).Updates(&map[string]interface{}{
		"room_number": RoomNumber,
		"updated_by":  UpdatedBy,
	})
	return RoomNumber
}

func AssignRoom(c *gin.Context, DB *gorm.DB, Dataset *global_var.TDataset, ReservationNumber uint64, BedTypeCode string, ReadyOnly bool, UpdatedBy string) string {
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
			"ORDER BY guest_detail.room_number", ReservationNumber, global_var.ReservationStatus.New).Find(&DataOutput)

	if len(DataOutput) > 0 {
		RoomList := GetAvailableRoomByType(DB, Dataset, (DataOutput["DateArrival"].(time.Time)), (DataOutput["DateDeparture"].(time.Time)), DataOutput["room_type_code"].(string), BedTypeCode, 0, 0, 0, 0, ReadyOnly, general.InterfaceToBool(DataOutput["is_from_allotment"]))
		if len(RoomList) > 0 {
			RoomNumber = RoomList[0].RoomNumber
			if RoomNumber != "" {
				UpdateReservationRoomNumber(c, DB, general.InterfaceToUint64(DataOutput["id"]), RoomNumber, UpdatedBy)
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

	GroupCode := master_data.GetAccountField(DB, "group_code", "cfg_init_account.code", AccountCode)
	if GroupCode == global_var.GlobalAccountGroup.Charge {
		if TaxAndServiceCodeManual == "" {
			TaxServiceCode = master_data.GetAccountField(DB, "tax_and_service_code", "cfg_init_account.code", AccountCode)
		} else {
			TaxServiceCode = TaxAndServiceCodeManual
		}
		//fmt.Println("Amount", Amount)
		//fmt.Println("AccountCode", AccountCode)
		//fmt.Println("TaxServiceCode", TaxServiceCode)
		if TaxServiceCode != "" {
			var TaxAndServiceData db_var.Cfg_init_tax_and_service
			DB.Table(db_var.TableName.CfgInitTaxAndService).Where("code", TaxServiceCode).Limit(1).Scan(&TaxAndServiceData)

			TaxPercent = TaxAndServiceData.Tax                         //master_data.GetFieldFloat(DB, db_var.TableName.CfgInitTaxAndService, "tax", "code", TaxServiceCode, 0)
			ServicePercent = TaxAndServiceData.Service                 // master_data.GetFieldFloat(DB, db_var.TableName.CfgInitTaxAndService, "service", "code", TaxServiceCode, 0)
			ServiceTaxPercent = TaxAndServiceData.ServiceTax           //master_data.GetFieldFloat(DB, db_var.TableName.CfgInitTaxAndService, "service_tax", "code", TaxServiceCode, 0)
			IsTaxIncluded = TaxAndServiceData.IsTaxInclude > 0         //master_data.GetFieldBool(DB, db_var.TableName.CfgInitTaxAndService, "is_tax_include", "code", TaxServiceCode, false)
			IsServiceIncluded = TaxAndServiceData.IsServiceInclude > 0 //master_data.GetFieldBool(DB, db_var.TableName.CfgInitTaxAndService, "is_service_include", "code", TaxServiceCode, false)
			//Tax and Service Include
			if IsTaxIncluded && IsServiceIncluded {
				Tax = general.RoundToX3(Amount * (TaxPercent + (ServiceTaxPercent * ServicePercent / 100)) / (100 + TaxPercent + ServicePercent + (ServiceTaxPercent * ServicePercent / 100)))
				Service = general.RoundToX3(Amount * ServicePercent / (100 + TaxPercent + ServicePercent + (ServiceTaxPercent * ServicePercent / 100)))
				Basic = Amount - Tax - Service
				if isDebug {
					////fmt.Println("1")
					////fmt.Println(Basic)
					////fmt.Println(Amount)
					////fmt.Println(Tax)
					////fmt.Println(Service)
				}
				//Tax and Service Exclude
			} else if !IsTaxIncluded && !IsServiceIncluded {
				Basic = Amount
				Tax = general.RoundToX3(Amount * (TaxPercent + (ServiceTaxPercent * ServicePercent / 100)) / 100)
				Service = general.RoundToX3(Amount * ServicePercent / 100)
				if isDebug {
					////fmt.Println("2")
					////fmt.Println(Basic)
					////fmt.Println(Amount)
					////fmt.Println(Tax)
					////fmt.Println(Service)
				}
				//Tax Exclude and Service Include
			} else if !IsTaxIncluded && IsServiceIncluded {
				Service = general.RoundToX3(Amount / (100 + ServicePercent) * ServicePercent)
				Basic = Amount - Service
				Tax = general.RoundToX3(Basic * TaxPercent / 100)
				if ServiceTaxPercent > 0 {
					Tax = general.RoundToX3(Amount * TaxPercent / 100)
				}
				Tax = general.RoundToX3(Amount * (TaxPercent + (ServiceTaxPercent * ServicePercent / 100)) / 100)
				if isDebug {
					////fmt.Println("3")
					////fmt.Println(Basic)
					////fmt.Println(Amount)
					////fmt.Println(Tax)
					////fmt.Println(Service)
				}
				//Tax Include and Service Exclude
			} else if IsTaxIncluded && !IsServiceIncluded {
				Tax = general.RoundToX3(Amount / (100 + TaxPercent + (ServiceTaxPercent * ServicePercent / 100)) * (TaxPercent + (ServiceTaxPercent * ServicePercent / 100)))
				Basic = Amount - Tax
				Service = general.RoundToX3((Amount - (Amount / (100 + TaxPercent + (ServiceTaxPercent * ServicePercent / 100)) * (TaxPercent + (ServiceTaxPercent * ServicePercent / 100)))) * ServicePercent / 100)

				if isDebug {
					////fmt.Println("4")
					////fmt.Println(Basic)
					////fmt.Println(Amount)
					////fmt.Println(Tax)
					////fmt.Println(Service)
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
	DB.Table(db_var.TableName.CfgInitAccount).
		Select("cfg_init_account.tax_and_service_code",
			"cfg_init_account_sub_group.group_code as account_group_code",
			"cfg_init_tax_and_service.tax",
			"cfg_init_tax_and_service.service",
			"cfg_init_tax_and_service.service_tax").
		Joins("LEFT JOIN cfg_init_account_sub_group ON (cfg_init_account.sub_group_code=cfg_init_account_sub_group.code)").
		Joins("LEFT JOIN cfg_init_tax_and_service ON (cfg_init_account.tax_and_service_code=cfg_init_tax_and_service.code)").
		Where("cfg_init_account.code = ?", AccountCode).Take(&DataOutputAccount)

	if DataOutputAccount.AccountGroupCode == global_var.GlobalAccountGroup.Charge {
		if TaxServiceCodeManual == "" {
			TaxServiceCode = DataOutputAccount.TaxAndServiceCode
		} else {
			TaxServiceCode = TaxServiceCodeManual
		}
		//////fmt.Println(TaxServiceCode)
		//fmt.Println("TaxServiceCodeManual", AccountCode, TaxServiceCodeManual)
		if TaxServiceCode != "" {
			var TaxAndServiceData db_var.Cfg_init_tax_and_service
			DB.Table(db_var.TableName.CfgInitTaxAndService).Where("code", TaxServiceCode).Limit(1).Scan(&TaxAndServiceData)

			TaxPercent = TaxAndServiceData.Tax
			ServicePercent = TaxAndServiceData.Service
			ServiceTaxPercent = TaxAndServiceData.ServiceTax
			// IsTaxIncluded := TaxAndServiceData.IsTaxInclude > 0
			// IsServiceIncluded := TaxAndServiceData.IsServiceInclude > 0

			Tax = general.RoundToX3(Amount * (TaxPercent + (ServiceTaxPercent * ServicePercent / 100)) / (100 + TaxPercent + ServicePercent + (ServiceTaxPercent * ServicePercent / 100)))
			Service = general.RoundToX3(Amount * ServicePercent / (100 + TaxPercent + ServicePercent + (ServiceTaxPercent * ServicePercent / 100)))
			Basic = Amount - Tax - Service
		}
	}
	return Basic, Tax, Service
}

func GetBasicTaxServiceForeign(Tax, Service, AmountForeign, ExchangeRate float64) (BasicForeign float64, TaxForeign float64, ServiceForeign float64) {
	// var TaxForeign, ServiceForeign, BasicForeign float64
	if Tax > 0 {
		TaxForeign = general.RoundToX3(Tax / ExchangeRate)
	} else {
		TaxForeign = 0
	}

	if Service > 0 {
		ServiceForeign = general.RoundToX3(Service / ExchangeRate)
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
		"WHERE sub_folio.audit_date BETWEEN  DATE(DATE_ADD('" + general.FormatDate1(AuditDate) + "', INTERVAL -1 DAY)) AND DATE('" + general.FormatDate1(AuditDate) + "') " +
		") AS A " +
		"ORDER BY A.correction_breakdown DESC " +
		"LIMIT 1;").Scan(&CorrectionBreakdown)

	if CorrectionBreakdown <= 0 {
		DB.Table(db_var.TableName.SubFolio).Select("correction_breakdown").Order("correction_breakdown DESC").Limit(1).Scan(&CorrectionBreakdown)
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
		"WHERE sub_folio.audit_date BETWEEN  DATE(DATE_ADD('" + general.FormatDate1(AuditDate) + "', INTERVAL -1 DAY)) AND DATE('" + general.FormatDate1(AuditDate) + "') " +
		") AS A " +
		"ORDER BY A.breakdown1 DESC " +
		"LIMIT 1;").Scan(&Breakdown1)

	if Breakdown1 <= 0 {
		DB.Table(db_var.TableName.SubFolio).Select("breakdown1").Order("breakdown1 DESC").Limit(1).Scan(&Breakdown1)
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
		"AND sub_folio.audit_date BETWEEN  DATE(DATE_ADD('" + general.FormatDate1(AuditDate) + "', INTERVAL -1 DAY)) AND DATE('" + general.FormatDate1(AuditDate) + "') " +
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
	if general.StrToBool(pConfig.Dataset.Configuration[global_var.ConfigurationCategory.General][global_var.ConfigurationName.IsRoomByName].(string)) {
		Result = master_data.GetFieldString(DB, db_var.TableName.CfgInitRoom, "number", "name", RoomName, RoomName)
	}
	return Result
}

func GetRoomNameFromConfigurationIsRoomName(c *gin.Context, DB *gorm.DB, RoomNumber string) string {
	Result := RoomNumber

	// Get Program Configuration
	val, exist := c.Get("pConfig")
	if !exist {
		// master_data.SendResponse(global_var.ResponseCode.DatabaseError, "pConfig", nil, c)
		// return
	}
	pConfig := val.(*config.CompanyDataConfiguration)

	IsRoomByName := general.StrToBool(pConfig.Dataset.Configuration[global_var.ConfigurationCategory.General][global_var.ConfigurationName.IsRoomByName].(string))
	if IsRoomByName {
		Result = master_data.GetFieldString(DB, db_var.TableName.CfgInitRoom, "name", "number", RoomNumber, RoomNumber)
	}
	return Result
}

func GetAutoRouting(DB *gorm.DB, BelongsTo uint64, AccountCode string) (uint64, string) {
	var FolioRouting db_var.Folio_routing
	DB.Table(db_var.TableName.FolioRouting).Select("folio_transfer, sub_folio_transfer").Where("folio_number=?", BelongsTo).Where("account_code=?", AccountCode).Limit(1).Scan(&FolioRouting)

	return FolioRouting.FolioTransfer, FolioRouting.SubFolioTransfer
}

func InsertSubFolioX(c *gin.Context, Dataset *global_var.TDataset, IDCorrected uint64, DataInput db_var.Sub_folio, DB *gorm.DB) (uint64, error) {
	var err error
	var Id uint64

	DataInput.DefaultCurrencyCode = GetDefaultCurrencyCode(DB)
	if DataInput.CurrencyCode == "" {
		DataInput.CurrencyCode = DataInput.DefaultCurrencyCode
	}
	DataInput.ExchangeRate = GetExchangeRateCurrency(DB, DataInput.CurrencyCode)
	DataInput.AmountForeign = DataInput.Amount

	if DataInput.CurrencyCode != DataInput.DefaultCurrencyCode {
		DataInput.Amount = general.RoundToX3(DataInput.Amount * DataInput.ExchangeRate)
	}
	IsRoomByName := Dataset.ProgramConfiguration.IsRoomByName
	if IsRoomByName {
		DataInput.RoomNumber = GetRoomNumberFromConfigurationIsRoomName(c, DB, DataInput.RoomNumber)
	}
	if DataInput.AuditDate.IsZero() {
		AuditDate := GetAuditDate(c, DB, false)
		DataInput.AuditDate = AuditDate
	}
	Timezone := Dataset.ProgramConfiguration.Timezone
	loc, _ := time.LoadLocation(Timezone)
	DataInput.AuditDateUnixx = general.DateOf(DataInput.AuditDate.In(loc)).Unix()

	DataInput.Id = 0
	if err := DB.Table(db_var.TableName.SubFolio).Omit("id").Create(&DataInput).Error; err != nil {
		return Id, err
	}
	Id = DataInput.Id
	// Insert Foreign Cash
	if (GetAccountSubGroupCode(DB, DataInput.AccountCode) == global_var.GlobalAccountSubGroup.Payment || GetAccountSubGroupCode(DB, DataInput.AccountCode) == global_var.GlobalAccountSubGroup.CreditDebitCard || GetAccountSubGroupCode(DB, DataInput.AccountCode) == global_var.GlobalAccountSubGroup.BankTransfer) && DataInput.CurrencyCode != DataInput.DefaultCurrencyCode {
		RemarkForeignCash := "Payment for Folio: " + strconv.FormatUint(DataInput.FolioNumber, 10) + ", Room: " + DataInput.RoomNumber + ", Doc#: " + DataInput.DocumentNumber
		TypeCodeX := global_var.TransactionType.Debit
		if DataInput.TypeCode == global_var.TransactionType.Debit {
			TypeCodeX = global_var.TransactionType.Credit
		}
		if general.Uint8ToBool(DataInput.IsCorrection) {
			RemarkForeignCash = "Payment Correction  for Folio: " + strconv.FormatUint(DataInput.FolioNumber, 10) + ", Room: " + DataInput.RoomNumber + ", Doc#: " + DataInput.DocumentNumber
		}

		var ForeignCash db_var.Acc_foreign_cash
		ForeignCash.IdTransaction = DataInput.FolioNumber
		ForeignCash.IdCorrected = IDCorrected
		ForeignCash.IdChange = 0
		ForeignCash.IdTable = global_var.ForeignCashTableID.SubFolio
		ForeignCash.Breakdown = 0
		ForeignCash.RefNumber = ""
		ForeignCash.Date = global_var.ProgramVariable.AuditDate
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
		if err := DB.Table(db_var.TableName.AccForeignCash).Create(&ForeignCash).Error; err != nil {
			return Id, err
		}
	}
	return Id, err
}

func InsertSubFolio(c *gin.Context, DB *gorm.DB, Dataset *global_var.TDataset, IsCanAutoTransfer bool, AccountCodeTransfer, TaxServiceManual string, DataInput db_var.Sub_folio) (Result string, SubFolioID uint64, Error error) {
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

	}
	Timezone := Dataset.ProgramConfiguration.Timezone
	loc, _ := time.LoadLocation(Timezone)
	DataInput.AuditDateUnixx = general.DateOf(DataInput.AuditDate.In(loc)).Unix()

	if DataInput.CurrencyCode == "" {
		DataInput.CurrencyCode = GetDefaultCurrencyCode(DB)
		DataInput.ExchangeRate = GetExchangeRateCurrency(DB, DataInput.CurrencyCode)
	}
	if DataInput.ExchangeRate == 0 {
		DataInput.ExchangeRate = GetExchangeRateCurrency(DB, DataInput.CurrencyCode)
	}
	AllowZeroAmount := Dataset.Configuration[global_var.ConfigurationCategory.Folio][global_var.ConfigurationName.AllowZeroAmount].(string)
	if ((DataInput.Quantity*DataInput.Amount) > 0 || ((DataInput.Quantity*DataInput.Amount) <= 0) && AllowZeroAmount != "0") && DataInput.CurrencyCode != "" && DataInput.ExchangeRate > 0 {

		DataInput.DefaultCurrencyCode = GetDefaultCurrencyCode(DB)
		if DataInput.CurrencyCode != DataInput.DefaultCurrencyCode {
			DataInput.Amount = general.RoundToX3(DataInput.Amount * DataInput.ExchangeRate)
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
		DataInput.AmountForeign = general.RoundToX3(DataInput.Amount / DataInput.ExchangeRate)
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
			if AccountSubGroup == global_var.GlobalAccountSubGroup.AccountReceivable || AccountSubGroup == global_var.GlobalAccountSubGroup.AccountPayable || AccountSubGroup == global_var.GlobalAccountSubGroup.Compliment {
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
	if err := Db.Table(db_var.TableName.GuestDeposit).Where("id=?", GuestDepositID1).Updates(map[string]interface{}{"is_pair_with_folio": false, "transfer_pair_id": GuestDepositID2, "updated_by": ValidUserCode}).Error; err != nil {
		return err
	}

	if err := Db.Table(db_var.TableName.GuestDeposit).Where("id=?", GuestDepositID2).Updates(map[string]interface{}{"is_pair_with_folio": false, "transfer_pair_id": GuestDepositID1, "updated_by": ValidUserCode}).Error; err != nil {
		return err
	}

	return nil
}

func UpdateRoomStatus(tx *gorm.DB, ValidUserCode, RoomNumber, RoomStatusCode string) error {
	if err := tx.Table(db_var.TableName.CfgInitRoom).Where("number = ?", RoomNumber).Updates(&map[string]interface{}{
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
	DB.Table(db_var.TableName.Folio).Select("status_code").Where(Condition, Number).Scan(&FolioStatus)

	return FolioStatus.StatusCode
}

func GetRoomStatus(DB *gorm.DB, Number string) (StatusCode string, BlockStatusCode string) {
	var RoomStatus struct {
		StatusCode      string
		BlockStatusCode string
	}
	DB.Table(db_var.TableName.CfgInitRoom).Where("number = ?", Number).Limit(1).Scan(&RoomStatus)

	return RoomStatus.StatusCode, RoomStatus.BlockStatusCode
}

func GetRoomTypeCode(DB *gorm.DB, Number string) string {
	var Room struct {
		RoomTypeCode string
	}
	DB.Table(db_var.TableName.CfgInitRoom).Where("number = ?", Number).Scan(&Room)

	return Room.RoomTypeCode
}

func GetBedTypeCode(DB *gorm.DB, Number string) string {
	var BedTypeCode string
	DB.Table(db_var.TableName.CfgInitRoom).Select("bed_type_code").Where("number = ?", Number).Limit(1).Scan(&BedTypeCode)

	return BedTypeCode
}

func GetReservationStatus(DB *gorm.DB, Number uint64) (StatusCode string) {
	DB.Table(db_var.TableName.Reservation).Select("status_code").Where("number = ? ", Number).Limit(1).Scan(&StatusCode)
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
	RoomAvailableCount, err := GetAvailableRoomCountByType(DB, Dataset, ArrivalDate, DepartureDate, RoomTypeCode, BedTypeCode, ReservationNumber, FolioNumber, RoomUnavailableID, RoomAllotmentID, ReadyOnly, AllotmentOnly)
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
	Prefix := PrefixX + general.FormatDatePrefix(PostingDate) + "-"
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

	return Prefix + general.Uint64ToStr(DataOutput+1)

	// }
	// return fmt.Sprintf("%s%d", Prefix, general.StrToInt64(Number)+1)
}

func GetJournalRefNumberTemp(c *gin.Context, DB *gorm.DB, PrefixX string, PostingDate time.Time) string {
	Prefix := PrefixX + general.FormatDatePrefix(PostingDate) + "-"
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

	return Prefix + general.Uint64ToStr(DataOutput+1)

	// }
	// return fmt.Sprintf("%s%d", Prefix, general.StrToInt64(Number)+1)
}

func GetReceiveNumber(c *gin.Context, DB *gorm.DB, PostingDate time.Time) string {
	var Prefix string
	Prefix = global_var.ConstProgramVariable.ReceiveNumberPrefix + general.FormatDatePrefix(PostingDate) + "-"
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

	return Prefix + general.Uint64ToStr(DataOutput+1)

	// }
	// return fmt.Sprintf("%s%d", Prefix, general.StrToInt64(Number)+1)
}

func GetPaymentNumber(c *gin.Context, DB *gorm.DB, PostingDate time.Time) string {
	var Prefix string
	Prefix = global_var.ConstProgramVariable.PaymentNumberPrefix + general.FormatDatePrefix(PostingDate) + "-"
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

	return Prefix + general.Uint64ToStr(DataOutput+1)
	// }
	// return fmt.Sprintf("%s%d", Prefix, general.StrToInt64(Number)+1)
}

func GetJournalAccountTypeCode(DB *gorm.DB, JournalAccountCode string) string {
	var TypeCode string
	DB.Table(db_var.TableName.CfgInitAccount).Select("type_code").Where("code=?", JournalAccountCode).Limit(1).Scan(&TypeCode)
	return TypeCode
}

func GetAPNumber(c *gin.Context, DB *gorm.DB, PostingDate time.Time) string {
	var Prefix string
	Prefix = global_var.ConstProgramVariable.APNumberPrefix + general.FormatDatePrefix(PostingDate) + "-"

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

	return Prefix + general.Uint64ToStr(DataOutput+1)
	// }
	// return fmt.Sprintf("%s%d", Prefix, general.StrToInt64(Number)+1)

}

func GetARNumber(c *gin.Context, DB *gorm.DB, PostingDate time.Time) string {
	var Prefix string
	Prefix = global_var.ConstProgramVariable.ARNumberPrefix + general.FormatDatePrefix(PostingDate) + "-"
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

	return Prefix + general.Uint64ToStr(DataOutput+1)
	// }

	// return fmt.Sprintf("%s%d", Prefix, general.StrToInt64(Number)+1)
}

func GetCostingNumber(c *gin.Context, DB *gorm.DB, PostingDate time.Time) string {
	Prefix := global_var.ConstProgramVariable.CostingNumberPrefix + general.FormatDatePrefix(PostingDate) + "-"
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

	return Prefix + general.Uint64ToStr(DataOutput+1)
	// }
	// return fmt.Sprintf("%s%d", Prefix, general.StrToInt64(Number)+1)
}

func GetDepreciationNumber(c *gin.Context, DB *gorm.DB, PostingDate time.Time) string {
	Prefix := global_var.ConstProgramVariable.DepreciationNumberPrefix + general.FormatDatePrefix(PostingDate) + "-"
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

	return Prefix + general.Uint64ToStr(DataOutput+1)
	// }
	// return fmt.Sprintf("%s%d", Prefix, general.StrToInt64(Number)+1)
}

func GetSRNumber(c *gin.Context, DB *gorm.DB, PostingDate time.Time) string {
	Prefix := global_var.ConstProgramVariable.SRNumberPrefix + general.FormatDatePrefix(PostingDate) + "-"
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

	return Prefix + general.Uint64ToStr(DataOutput+1)
	// }
	// return fmt.Sprintf("%s%d", Prefix, general.StrToInt64(Number)+1)
}
func GetProductionNumber(c *gin.Context, DB *gorm.DB, PostingDate time.Time) string {
	Prefix := global_var.ConstProgramVariable.ProductionNumberPrefix + general.FormatDatePrefix(PostingDate) + "-"
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

	return Prefix + general.Uint64ToStr(DataOutput+1)
	// }
	// return fmt.Sprintf("%s%d", Prefix, general.StrToInt64(Number)+1)
}

func GetReturnStockNumber(c *gin.Context, DB *gorm.DB, PostingDate time.Time) string {
	Prefix := global_var.ConstProgramVariable.ReturnStockPrefix + general.FormatDatePrefix(PostingDate) + "-"
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

	return Prefix + general.Uint64ToStr(DataOutput+1)
	// }
	// return fmt.Sprintf("%s%d", Prefix, general.StrToInt64(Number)+1)
}

func GetOpnameNumber(c *gin.Context, DB *gorm.DB, PostingDate time.Time) string {
	Prefix := global_var.ConstProgramVariable.OpnameNumberPrefix + general.FormatDatePrefix(PostingDate) + "-"
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

	return Prefix + general.Uint64ToStr(DataOutput+1)
	// }
	// return fmt.Sprintf("%s%d", Prefix, general.StrToInt64(Number)+1)
}

func GetFAPONumber(c *gin.Context, DB *gorm.DB, PostingDate time.Time) string {
	Prefix := global_var.ConstProgramVariable.FAPONumberPrefix + general.FormatDatePrefix(PostingDate) + "-"
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

	return Prefix + general.Uint64ToStr(DataOutput+1)
	// }
	// return fmt.Sprintf("%s%d", Prefix, general.StrToInt64(Number)+1)
}

func GetFAReceiveNumber(c *gin.Context, DB *gorm.DB, PostingDate time.Time) string {
	Prefix := global_var.ConstProgramVariable.FAReceiveNumberPrefix + general.FormatDatePrefix(PostingDate) + "-"
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

	return Prefix + general.Uint64ToStr(DataOutput+1)
	// }
	// return fmt.Sprintf("%s%d", Prefix, general.StrToInt64(Number)+1)
}

func GetFACode(c *gin.Context, DB *gorm.DB, ItemCode string, PostingDate time.Time) (SortNumber uint64, Code string) {
	Prefix := general.FormatDatePrefix(PostingDate) + ItemCode
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

func GetProductItemGroup(DB *gorm.DB, ProductCode string) string {
	var ItemGroupCode string
	DB.Raw(
		"SELECT"+
			" pos_cfg_init_product_group.item_group_code "+
			"FROM"+
			" pos_cfg_init_product"+
			" LEFT OUTER JOIN pos_cfg_init_product_group ON (pos_cfg_init_product.group_code = pos_cfg_init_product_group.code)"+
			" WHERE pos_cfg_init_product.code=?", ProductCode).Scan(&ItemGroupCode)

	return ItemGroupCode
}

func GetFACondition(DB *gorm.DB, FACode string) string {
	ConditionCode := ""
	DB.Table(db_var.TableName.FaList).Select("condition_code").Where("code=?", FACode).Limit(1).Scan(&ConditionCode)
	return ConditionCode
}

func GetFAJournalAccount(DB *gorm.DB, ItemCode string) string {
	var JournalAccountCode string
	DB.Table(db_var.TableName.FaCfgInitItem).Select("fa_cfg_init_item_category.journal_account_code").
		Joins("LEFT OUTER JOIN fa_cfg_init_item_category ON (fa_cfg_init_item.category_code = fa_cfg_init_item_category.code)").
		Where("fa_cfg_init_item.code=?", ItemCode).
		Scan(&JournalAccountCode)
	return JournalAccountCode
}

func GetFAJournalAccountDepreciation(DB *gorm.DB, ItemCode string) string {
	var JournalAccountCodeDepreciation string
	DB.Table(db_var.TableName.FaCfgInitItem).Select("fa_cfg_init_item_category.journal_account_code_depreciation").
		Joins("LEFT OUTER JOIN fa_cfg_init_item_category ON (fa_cfg_init_item.category_code = fa_cfg_init_item_category.code)").
		Where("fa_cfg_init_item.code=?", ItemCode).
		Scan(&JournalAccountCodeDepreciation)
	return JournalAccountCodeDepreciation
}

func GetFACOGSAccount(DB *gorm.DB, ItemCode string) string {
	var JournalAccountCodeCogs string
	DB.Table(db_var.TableName.FaCfgInitItem).Select("fa_cfg_init_item_category.journal_account_code_cogs").
		Joins("LEFT OUTER JOIN fa_cfg_init_item_category ON (fa_cfg_init_item.category_code = fa_cfg_init_item_category.code)").
		Where("fa_cfg_init_item.code=?", ItemCode).
		Limit(1).
		Scan(&JournalAccountCodeCogs)
	return JournalAccountCodeCogs
}

func GetFAExpenseAccount(DB *gorm.DB, ItemCode string) string {
	var JournalAccountCodeExpense string
	DB.Table(db_var.TableName.FaCfgInitItem).Select("fa_cfg_init_item_category.journal_account_code_expense").
		Joins("LEFT OUTER JOIN fa_cfg_init_item_category ON (fa_cfg_init_item.category_code = fa_cfg_init_item_category.code)").
		Where("fa_cfg_init_item.code=?", ItemCode).
		Limit(1).
		Scan(&JournalAccountCodeExpense)
	return JournalAccountCodeExpense
}

func GetFASellAccount(DB *gorm.DB, ItemCode string) string {
	var JournalAccountCodeSell string
	DB.Table(db_var.TableName.FaCfgInitItem).Select("fa_cfg_init_item_category.journal_account_code_sell").
		Joins("LEFT OUTER JOIN fa_cfg_init_item_category ON (fa_cfg_init_item.category_code = fa_cfg_init_item_category.code)").
		Where("fa_cfg_init_item.code=?", ItemCode).
		Limit(1).
		Scan(&JournalAccountCodeSell)
	return JournalAccountCodeSell
}

func GetFASpoilJournalAccount(DB *gorm.DB, ItemCode string) string {
	var JournalAccountCodeSpoil string
	DB.Table(db_var.TableName.FaCfgInitItem).Select("fa_cfg_init_item_category.journal_account_code_spoil").
		Joins("LEFT OUTER JOIN fa_cfg_init_item_category ON (fa_cfg_init_item.category_code = fa_cfg_init_item_category.code)").
		Where("fa_cfg_init_item.code=?", ItemCode).
		Limit(1).
		Scan(&JournalAccountCodeSpoil)
	return JournalAccountCodeSpoil
}

func GetStockTransferNumber(c *gin.Context, DB *gorm.DB, PostingDate time.Time) string {
	var Prefix string
	Prefix = global_var.ConstProgramVariable.StockTransferNumberPrefix + general.FormatDatePrefix(PostingDate) + "-"
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

	return Prefix + general.Uint64ToStr(DataOutput+1)

	// }
	// return fmt.Sprintf("%s%d", Prefix, general.StrToInt64(Number)+1)
}

func IsYearClosed(DB *gorm.DB, Year uint64) bool {
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
		IsMonthClosed, _ := IsMonthClosed(DB, uint64(Month), uint64(Year))
		return IsMonthClosed || IsYearClosed(DB, uint64(Year))
	} else {
		MonthB4 := PostingDateB4.Month()
		YearB4 := PostingDateB4.Year()

		IsMonthClosed1, _ := IsMonthClosed(DB, uint64(Month), uint64(YearB4))
		IsMonthClosed2, _ := IsMonthClosed(DB, uint64(MonthB4), uint64(YearB4))

		return IsMonthClosed1 || IsYearClosed(DB, uint64(Year)) || IsMonthClosed2 || IsYearClosed(DB, uint64(YearB4))
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
			" WHERE id=?", global_var.TransactionType.Debit, SubFolioId, RefNumber, Dataset.GlobalAccount.APRefundDeposit, SubFolioId).Scan(&DataOutput)

	return DataOutput
}

func GetARCityLedgerInvoiceOutStanding(DB *gorm.DB, InvoiceNumber string) (float64, error) {
	var Amount float64
	if err := DB.Table(db_var.TableName.InvoiceItem).Select(
		"SUM(amount_charged - amount_paid) AS Amount",
	).Where("invoice_number=?", InvoiceNumber).
		Scan(&Amount).Error; err != nil {
		return 0, err
	}
	return Amount, nil
}
func IsFolioClosed(DB *gorm.DB, FolioNumber uint64) bool {
	var StatusCode string
	DB.Table(db_var.TableName.Folio).Select("status_code").Where("number = ? AND status_code=?", FolioNumber, global_var.FolioStatus.Open).Find(StatusCode)
	return StatusCode != global_var.FolioStatus.Open
}

// func GetGlobalAccount(Account string) string {
// 	return master_data.GetConfiguration(global_var.SystemCode.Hotel, global_var.ConfigurationCategory.GlobalAccount, Account, false).(string)
// }

// func GetGlobalSubDepartment(SubDepartment string) string {
// 	return master_data.GetConfiguration(global_var.SystemCode.Hotel, global_var.ConfigurationCategory.GlobalSubDepartment, SubDepartment, false).(string)
// }

// func GetGlobalDepartment(Department string) string {
// 	return master_data.GetConfiguration(global_var.SystemCode.Hotel, global_var.ConfigurationCategory.GlobalDepartment, Department, false).(string)
// }
// func GetGlobalJournalAccount(AccountName string) string {
// 	return master_data.GetConfiguration(global_var.SystemCode.Hotel, global_var.ConfigurationCategory.GlobalJournalAccount, AccountName, false).(string)
// }

func IsFolioHaveBreakfast(DB *gorm.DB, Dataset *global_var.TDataset, FolioNumber uint64) bool {
	Number := 0
	DB.Table(db_var.TableName.Folio).Select("folio.number").
		Joins("LEFT OUTER JOIN guest_detail ON (folio.guest_detail_id = guest_detail.id)").
		Joins("LEFT OUTER JOIN cfg_init_room_rate ON (guest_detail.room_rate_code = cfg_init_room_rate.code)").
		Joins("LEFT OUTER JOIN cfg_init_room_rate_breakdown ON (guest_detail.room_rate_code = cfg_init_room_rate_breakdown.room_rate_code)").
		Where("folio.number = ? AND (IFNULL(cfg_init_room_rate_breakdown.account_code, '') = ? OR cfg_init_room_rate.include_breakfast='1')", FolioNumber, Dataset.GlobalAccount.Breakfast).
		Find(&Number)

	return Number > 0
}

func IsReservationHaveBreakfast(DB *gorm.DB, Dataset *global_var.TDataset, ReservationNumber uint64) bool {
	Number := 0
	DB.Table(db_var.TableName.Reservation).Select("reservation.number").
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

	var DataOutput db_var.Cfg_init_room_rate
	DB.Table(db_var.TableName.CfgInitRoomRate).Select(" weekday_rate1,"+
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
		IsWeekendX = general.IsWeekend(PostingDate, Dataset)
	}

	if IsWeekendX {
		if general.StrToBool(Dataset.Configuration[global_var.ConfigurationCategory.General][global_var.ConfigurationName.UseChildRate].(string)) {
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
		if general.StrToBool(Dataset.Configuration[global_var.ConfigurationCategory.General][global_var.ConfigurationName.UseChildRate].(string)) {
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
	ArrivalDate = general.DateOf(ArrivalDate)
	WeeklyCharged = general.DaysBetween(ArrivalDate, AuditDate)%7 == 0
	MonthlyCharged = general.DaysBetween(ArrivalDate, AuditDate)%30 == 0

	return ((ChargeFrequencyCode == global_var.ChargeFrequency.OnceOnly && ArrivalDate == AuditDate) ||
		(ChargeFrequencyCode == global_var.ChargeFrequency.Daily) ||
		(ChargeFrequencyCode == global_var.ChargeFrequency.Weekly && WeeklyCharged) ||
		(ChargeFrequencyCode == global_var.ChargeFrequency.Monthly && MonthlyCharged))

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
			AmountExtra = BreakdownAmountExtra
		}
	}

	return Amount + AmountExtra
}

func GetCommission(c *gin.Context, DB *gorm.DB, CommissionTypeCode string, CommissionValue, RoomRateAmount, RoomRateBasicAmount float64, ArrivalDate time.Time) float64 {
	AuditDate := GetAuditDate(c, DB, false)
	Value := 0.00
	if CommissionTypeCode == global_var.CommissionType.PercentFirstNightFullRate {
		if general.DateOf(ArrivalDate).Equal(general.DateOf(AuditDate)) {
			Value = general.RoundToX3(CommissionValue * RoomRateAmount / 100)
			return Value
		}
	}

	if CommissionTypeCode == global_var.CommissionType.PercentPerNightFullRate {
		Value = general.RoundToX3(CommissionValue * RoomRateAmount / 100)
		return Value
	}

	if CommissionTypeCode == global_var.CommissionType.PercentFirstNightNettRate {
		if general.DateOf(ArrivalDate).Equal(general.DateOf(AuditDate)) {
			Value = general.RoundToX3(CommissionValue * RoomRateBasicAmount / 100)
			return Value
		}
	}

	if CommissionTypeCode == global_var.CommissionType.PercentPerNightNettRate {
		Value = general.RoundToX3(CommissionValue * RoomRateBasicAmount / 100)
		return Value
	}

	if CommissionTypeCode == global_var.CommissionType.FixAmountPerNight {
		Value = general.RoundToX3(CommissionValue)
		return Value
	}

	return Value
}

func GetCommissionPackage(c *gin.Context, DB *gorm.DB, CommissionTypeCode string, CommissionValue, PackageAmount float64, ArrivalDate time.Time) (value float64) {
	AuditDate := GetAuditDate(c, DB, false)
	if CommissionTypeCode == global_var.CommissionType.PercentFirstNightFullRate {
		if general.DateOf(ArrivalDate).Equal(general.DateOf(AuditDate)) {
			value = general.RoundToX3(CommissionValue * PackageAmount / 100)
		}
	} else if CommissionTypeCode == global_var.CommissionType.PercentPerNightFullRate {
		value = general.RoundToX3(CommissionValue * PackageAmount / 100)
	} else if CommissionTypeCode == global_var.CommissionType.FixAmountFirstNight {
		if general.DateOf(ArrivalDate).Equal(general.DateOf(AuditDate)) {
			value = CommissionValue
		}
	} else if CommissionTypeCode == global_var.CommissionType.FixAmountPerNight {
		value = CommissionValue
	}

	if CommissionTypeCode == global_var.CommissionType.FixAmountPerNight {
		value = general.RoundToX3(CommissionValue)
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

			Basic = general.RoundToX3(BasicForeign * ExchangeRate)
			Tax = general.RoundToX3(TaxForeign * ExchangeRate)
			Service = general.RoundToX3(ServiceForeign * ExchangeRate)
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

		var SubFolioData db_var.Sub_folio
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
	var GuestInHouseBreakdown db_var.Guest_in_house_breakdown
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

	err := DB.Table(db_var.TableName.GuestInHouseBreakdown).Create(&GuestInHouseBreakdown).Error
	if err != nil {
		return err
	}
	return nil
}

func IsInHousePosted(DB *gorm.DB, AuditDate time.Time, FolioNumber uint64) bool {
	Id := 0
	DB.Table(db_var.TableName.GuestInHouse).Select("id").Where("audit_date = ? AND folio_number=?", AuditDate, strconv.FormatUint(FolioNumber, 10)).Limit(1).Scan(&Id)

	return Id > 0
}

// func InsertGuestInHouse(DB *gorm.DB, AuditDate time.Time, FolioNumber uint64, GroupCode, RoomTypeCode, BedTypeCode,
// 	RoomNumber, RoomRateCode, BusinessSourceCode, CommissionTypeCode, PaymentTypeCode, MarketCode, TitleCode,
// 	FullName, Street, City, CityCode, CountryCode, StateCode, PostalCode, Phone1, Phone2, Fax, Email, Website, CompanyCode, GuestTypeCode, MarketingCode,
// 	ComplimentHu, Notes string, Adult, Child int, Rate, RateOriginal, Discount, CommissionValue float64,
// 	DiscountPercent, IsAdditional, IsScheduledRate, IsBreakfast uint8, UserID string) error {

// 	var GuestInHouseData db_var.Guest_in_house
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

// 	err := DB.Table(db_var.TableName.GuestInHouse).Create(&GuestInHouseData).Error

// 	if err != nil {
// 		return err
// 	}

// 	return nil

// }

func PostingRoomChargeManual(ctx context.Context, c *gin.Context, DB *gorm.DB, Dataset *global_var.TDataset, FolioNumber uint64, SubFolioGroupCode, SubDepartmentCode, CurrencyCode, Remark, DocumentNumber string, Amount, ExchangeRate float64, IsCanAutoTransferred bool, UserID string) error {
	ctx, span := global_var.Tracer.Start(ctx, "PostingRoomChargeManual")
	defer span.End()

	type DataOutputStruct struct {
		Folio         db_var.Folio              `gorm:"embedded"`
		ContactPerson db_var.Contact_person     `gorm:"embedded"`
		GuestDetail   db_var.Guest_detail       `gorm:"embedded"`
		GuestGeneral  db_var.Guest_general      `gorm:"embedded"`
		RoomRate      db_var.Cfg_init_room_rate `gorm:"embedded"`
		DateArrival   time.Time
	}

	var DataOutput DataOutputStruct
	var RoomRateAmountOriginal, RoomRateAmount, RoomChargeB4Breakdown, RoomChargeAfterBreakdown, RoomChargeBasic, RoomChargeTax, RoomChargeService, TotalBreakdown, Commission, BreakdownAmount, BreakdownBasic, BreakdownTax, BreakdownService float64
	var RoomNumber string
	var BusinessSourceCode string
	var IsBreakfast uint8
	var CorrectionBreakdown, BreakDown1 uint64

	if err := DB.Table(db_var.TableName.Folio).Select(
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

		var DataOutputGuestBreakdown []db_var.Guest_breakdown
		//Proses Query Breakdown
		DB.Table(db_var.TableName.GuestBreakdown).Select(
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
					CurrencyCode, Remark, DocumentNumber, "", global_var.TransactionType.Debit, "", "", CorrectionBreakdown,
					BreakDown1, global_var.SubFolioPostingType.Room, 0, RoomChargeBasic, RoomChargeTax, RoomChargeService,
					ExchangeRate, false, IsCanAutoTransferred, UserID)

				if err != nil {
					return err
				}

				AuditDate := GetAuditDate(c, tx, false)
				if !IsInHousePosted(tx, AuditDate, FolioNumber) {
					err = InsertGuestInHouse(tx, Dataset, AuditDate,
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
						RoomRateAmount,
						RoomRateAmountOriginal,
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
						SubFolioData := db_var.Sub_folio{}
						SubFolioData.FolioNumber = FolioNumber
						SubFolioData.GroupCode = SubFolioGroupCode
						SubFolioData.RoomNumber = RoomNumber
						SubFolioData.SubDepartmentCode = SubDepartmentCode
						SubFolioData.AccountCode = guestBreakdown.AccountCode
						SubFolioData.ProductCode = guestBreakdown.ProductCode
						SubFolioData.CurrencyCode = CurrencyCode
						SubFolioData.Remark = "Breakdown: " + guestBreakdown.Remark
						SubFolioData.TypeCode = global_var.TransactionType.Debit
						SubFolioData.CorrectionBreakdown = CorrectionBreakdown
						SubFolioData.Breakdown1 = BreakDown1
						SubFolioData.DirectBillCode = guestBreakdown.CompanyCode
						SubFolioData.PostingType = global_var.SubFolioPostingType.Room
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
					SubFolioData := db_var.Sub_folio{}
					SubFolioData.FolioNumber = FolioNumber
					SubFolioData.GroupCode = SubFolioGroupCode
					SubFolioData.RoomNumber = RoomNumber
					SubFolioData.SubDepartmentCode = Dataset.GlobalSubDepartment.FrontOffice
					SubFolioData.AccountCode = Dataset.GlobalAccount.APCommission
					SubFolioData.ProductCode = Dataset.GlobalAccount.RoomCharge
					SubFolioData.CurrencyCode = CurrencyCode
					SubFolioData.Remark = "Breakdown Commission"
					SubFolioData.TypeCode = global_var.TransactionType.Debit
					SubFolioData.CorrectionBreakdown = CorrectionBreakdown
					SubFolioData.Breakdown1 = BreakDown1
					SubFolioData.DirectBillCode = BusinessSourceCode
					SubFolioData.PostingType = global_var.SubFolioPostingType.Room
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
	err := DB.Table(db_var.TableName.SubFolio).
		Where("folio_number=? AND void='0'", FolioNumber).
		Update("group_code", SubFolioGroupCode).Update("updated_by", UserID).Error

	return err
}

func UpdateSubFolioTransferPairID(DB *gorm.DB, SubFolioID1, SubFolioID2 uint64, UserID string) (err error) {
	err = DB.Table(db_var.TableName.SubFolio).
		Where("id=?", SubFolioID1).Updates(&map[string]interface{}{
		"is_pair_with_deposit": 0,
		"transfer_pair_id":     SubFolioID2,
		"updated_by":           UserID,
	}).Error

	err = DB.Table(db_var.TableName.SubFolio).
		Where("id=?", SubFolioID2).Updates(&map[string]interface{}{
		"is_pair_with_deposit": 0,
		"transfer_pair_id":     SubFolioID1,
		"updated_by":           UserID,
	}).Error

	return err
}

func MoveSubFolioByFolioNumber(DB *gorm.DB, FolioNumberFrom, FolioNumberTo uint64, SubFolioGroupCode, UserID string) error {
	err := DB.Table(db_var.TableName.SubFolio).
		Select("folio_number", "group_code", "updated_by").
		Where("folio_number = ? AND void='0'", FolioNumberFrom).
		Updates(map[string]interface{}{"folio_number": FolioNumberTo, "group_code": SubFolioGroupCode, "updated_by": UserID}).Error

	return err
}

func MoveSubFolioByBreakdown(DB *gorm.DB, FolioNumber uint64, SubFolioGroupCode string, CorrectionBreakdown uint64, UserID string) error {
	err := DB.Table(db_var.TableName.SubFolio).
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
	DB.Table(db_var.TableName.Folio).Select(
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
			" AND folio.status_code='" + global_var.FolioStatus.Open + "' " +
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
			" AND folio.status_code='" + global_var.FolioStatus.Open + "' " +
			"GROUP BY sub_folio.belongs_to)) AS FolioReturn " +
			"GROUP BY number " +
			"ORDER BY id_sort, room_number, number").Scan(&DataOutput).Error

	return DataOutput, err
}

func GetInvoiceNumber(DB *gorm.DB, IssuedDate time.Time) string {
	ServerID := GetServerID(DB)
	Prefix := global_var.ConstProgramVariable.InvoiceNumberPrefix + general.FormatDatePrefix(IssuedDate) + strconv.Itoa(ServerID) + "-"
	var DataOutput uint64
	DB.Raw(
		"SELECT" +
			" CAST(RIGHT(number, LENGTH(number) - " + strconv.FormatInt(int64(len(Prefix)), 10) + ") AS UNSIGNED) AS MaxInvoiceNumber " +
			"FROM" +
			" invoice" +
			" WHERE LEFT(number," + strconv.FormatInt(int64(len(Prefix)), 10) + ")='" + Prefix + "' " +
			"ORDER BY MaxInvoiceNumber DESC " +
			"LIMIT 1;").Scan(&DataOutput)

	return Prefix + general.Uint64ToStr(DataOutput+1)
}
func GetReceiptNumber(DB *gorm.DB, IssuedDate time.Time) string {
	ServerID := GetServerID(DB)
	Prefix := general.FormatDatePrefix(IssuedDate) + strconv.Itoa(ServerID) + "-"
	var DataOutput uint64
	DB.Raw(
		"SELECT" +
			" CAST(RIGHT(number, LENGTH(number) - " + strconv.FormatInt(int64(len(Prefix)), 10) + ") AS UNSIGNED) AS MaxReceiptNumber " +
			"FROM" +
			" receipt" +
			" WHERE LEFT(number," + strconv.FormatInt(int64(len(Prefix)), 10) + ")='" + Prefix + "' " +
			"ORDER BY MaxReceiptNumber DESC " +
			"LIMIT 1;").Scan(&DataOutput)
	//////fmt.Println(ServerID)
	//////fmt.Println(Prefix)
	//////fmt.Println(DataOutput)
	//////fmt.Println(Prefix + general.Uint64ToStr(DataOutput+1))

	return Prefix + general.Uint64ToStr(DataOutput+1)
}

func GetPRNumber(DB *gorm.DB, IssuedDate time.Time) string {
	ServerID := GetServerID(DB)
	Prefix := global_var.ConstProgramVariable.PRNumberPrefix + general.FormatDatePrefix(IssuedDate) + strconv.Itoa(ServerID) + "-"
	var DataOutput uint64
	DB.Raw(
		"SELECT" +
			" CAST(RIGHT(number, LENGTH(number) - " + strconv.FormatInt(int64(len(Prefix)), 10) + ") AS UNSIGNED) AS MaxPRNumber " +
			"FROM" +
			" inv_purchase_request" +
			" WHERE LEFT(number," + strconv.FormatInt(int64(len(Prefix)), 10) + ")='" + Prefix + "' " +
			"ORDER BY MaxPRNumber DESC " +
			"LIMIT 1;").Scan(&DataOutput)

	return Prefix + general.Uint64ToStr(DataOutput+1)
}

func GetPONumber(DB *gorm.DB, IssuedDate time.Time) string {
	ServerID := GetServerID(DB)
	Prefix := global_var.ConstProgramVariable.PONumberPrefix + general.FormatDatePrefix(IssuedDate) + strconv.Itoa(ServerID) + "-"
	var DataOutput uint64
	DB.Raw(
		"SELECT" +
			" CAST(RIGHT(number, LENGTH(number) - " + strconv.FormatInt(int64(len(Prefix)), 10) + ") AS UNSIGNED) AS MaxPONumber " +
			"FROM" +
			" inv_purchase_order" +
			" WHERE LEFT(number," + strconv.FormatInt(int64(len(Prefix)), 10) + ")='" + Prefix + "' " +
			"ORDER BY MaxPONumber DESC " +
			"LIMIT 1;").Scan(&DataOutput)

	return Prefix + general.Uint64ToStr(DataOutput+1)
}

func GetJournalAccountCompanyAR(DB *gorm.DB, CompanyCode string) string {
	var DataOutput string
	DB.Table(db_var.TableName.Company).Select("cfg_init_company_type.journal_account_code_ar").
		Joins("LEFT OUTER JOIN cfg_init_company_type ON (company.type_code = cfg_init_company_type.code)").
		Where("company.code = ?", CompanyCode).Scan(&DataOutput)

	return DataOutput

}

func GetJournalIDSort(DB *gorm.DB, PostingDate time.Time) int {
	var IdSort int
	DB.Table(db_var.TableName.AccJournal).Select("id_sort").Where("date=?", PostingDate).Order("id_sort DESC").Limit(1).Scan(&IdSort)

	return IdSort + 1
}

func IsInvoiceHadPayment(DB *gorm.DB, InvoiceNumberX string) bool {
	var InvoiceNumber string
	DB.Table(db_var.TableName.InvoicePayment).Select("invoice_number").Where("invoice_number=?", InvoiceNumber).Take(&InvoiceNumber)
	return InvoiceNumber != ""
}

func IsSubFolioIdHadInvoice(DB *gorm.DB, SubFolioId uint64) bool {
	var SubFolioIdX uint64
	DB.Table(db_var.TableName.InvoiceItem).Select("sub_folio_id").Where("sub_folio_id=?", SubFolioId).Limit(1).Scan(&SubFolioIdX)

	return SubFolioIdX > 0
}

func GetFolioType(DB *gorm.DB, FolioNumber uint64) string {
	var Type string
	DB.Table(db_var.TableName.Folio).Select("type_code").Where("number=?", FolioNumber).Limit(1).Scan(&Type)

	return Type
}

func GetFolioSystemCode(DB *gorm.DB, FolioNumber uint64) string {
	var Code string
	DB.Table(db_var.TableName.Folio).Select("system_code").Where("number=?", FolioNumber).Limit(1).Scan(&Code)

	return Code
}
func IsThereCardActiveFolio(DB *gorm.DB, FolioNumber uint64) bool {
	return master_data.GetFieldBool(DB, db_var.TableName.LogKeylock, "id", "is_active = '1' AND folio_number", general.Uint64ToStr(FolioNumber), false)
}

func IsCanCreateInvoice(DB *gorm.DB, FolioNumber string) (DirectBillCode string, CanCreate bool) {
	DB.Table(db_var.TableName.SubFolio).Select("sub_folio.direct_bill_code").
		Joins("LEFT OUTER JOIN cfg_init_account ON (sub_folio.account_code = cfg_init_account.code)").
		Where("sub_folio.folio_number=?", FolioNumber).
		Where("sub_folio.direct_bill_code<>''").
		Where("cfg_init_account.sub_group_code=?", global_var.GlobalAccountSubGroup.AccountReceivable).
		Where("void='0'").
		Group("sub_folio.correction_breakdown").
		Limit(1).
		Scan(&DirectBillCode)

	return DirectBillCode, DirectBillCode != ""
}

func CheckFolioReceiveTransfer(DB *gorm.DB, FolioNumber uint64) []string {
	var FolioTransferMessage []string
	var FolioDetail []string

	DB.Table(db_var.TableName.FolioRouting).Select(" DISTINCT CONCAT(folio.number, '/Room: ', guest_detail.room_number, '/', contact_person.title_code, contact_person.full_name) AS FolioDetail").
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
	DB.Table(db_var.TableName.InvoicePayment).Select("SUM(amount_foreign) AS TotalAmount").
		Where("invoice_number=?", InvoiceNumber).
		Group("invoice_number").Scan(&TotalPayment)

	return TotalPayment
}

func IsPrepaidExpensePosted(DB *gorm.DB, PrepaidID, PrepaidPostedID uint64, PostingDate time.Time) bool {
	var Id uint64
	Year, Month, _ := PostingDate.Date()
	Query := DB.Table(db_var.TableName.AccPrepaidExpensePosted).Select("id").
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
	Query := DB.Table(db_var.TableName.AccDefferedIncomePosted).Select("id").
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
	DB.Table(db_var.TableName.CfgInitAccount).Select("journal_account_code").Where("code=?", AccountCode).Limit(1).Scan(&JournalAccountCode)

	return JournalAccountCode
}

func GetJournalBankAccountCode(DB *gorm.DB, BankAccountCode string) string {
	var JournalAccountCode string
	DB.Table(db_var.TableName.AccCfgInitBankAccount).Select("journal_account_code").
		Joins("LEFT JOIN cfg_init_journal_account ON acc_cfg_init_bank_account.journal_account_code =cfg_init_journal_account.code").
		Where("acc_cfg_init_bank_account.code=?", BankAccountCode).Limit(1).Scan(&JournalAccountCode)
	return JournalAccountCode
}

func GetJournalAccountCurrency(DB *gorm.DB, CurrencyCode string) (JournalAccountCode string) {
	DB.Table(db_var.TableName.CfgInitCurrency).Select("cfg_init_account.journal_account_code").
		Joins("LEFT OUTER JOIN cfg_init_account ON (cfg_init_currency.account_code = cfg_init_account.code)").
		Where("cfg_init_currency.code=?", CurrencyCode).
		Limit(1).
		Scan(&JournalAccountCode)

	return
}

func GetForeignCashBreakdown(DB *gorm.DB) (BreakdownId uint64) {
	DB.Table(db_var.TableName.AccForeignCash).Select("breakdown").
		Order("breakdown DESC").
		Limit(1).
		Scan(&BreakdownId)
	BreakdownId++
	return
}

func GetJournalAccountCodeFromTransaction(DB *gorm.DB, IdTable, IdTransaction uint64) (AccountCode string) {
	MainTableName := db_var.TableName.SubFolio
	if IdTable == 1 {
		MainTableName = db_var.TableName.GuestDeposit
	}
	DB.Table(MainTableName).Select("cfg_init_account.journal_account_code").
		Joins("LEFT OUTER JOIN cfg_init_account ON ("+MainTableName+".account_code = cfg_init_account.code)").
		Where(MainTableName+".id=?", IdTransaction).
		Limit(1).
		Scan(&AccountCode)

	return
}

func GetJournalAccountListByGroup(DB *gorm.DB, SubDepartmentCode, GroupCode1, GroupCode2, GroupCode3 string) (DataOutput []map[string]interface{}) {
	Query1 := ""
	Query2 := ""
	if GroupCode2 != "" {
		Query1 = " OR cfg_init_journal_account_sub_group.group_code='" + GroupCode2 + "' "
	}
	if GroupCode3 != "" {
		Query2 = " OR cfg_init_journal_account_sub_group.group_code='" + GroupCode3 + "' "
	}
	Query := DB.Table(db_var.TableName.CfgInitJournalAccount).Select(
		" cfg_init_journal_account.code,"+
			" cfg_init_journal_account.name").
		Joins("LEFT OUTER JOIN cfg_init_journal_account_sub_group ON (cfg_init_journal_account.sub_group_code = cfg_init_journal_account_sub_group.code)").
		Where("(cfg_init_journal_account_sub_group.group_code=? "+Query1+Query2+") ", GroupCode1)

	Query.Where("sub_department_code LIKE ?", "%"+SubDepartmentCode+"%").Order("cfg_init_journal_account.code").
		Scan(&DataOutput)
	return
}

func GetPurchaseRequestStatus(DB *gorm.DB, Id uint64) (Status string) {
	DB.Table(db_var.TableName.InvPurchaseRequest).Select("status_code").Where("id=?", Id).Limit(1).Scan(&Status)
	if Status == "" {
		Status = global_var.PurchaseRequestStatus.NotApproved
	}
	return
}

func IsPurchaseRequestApproved1(DB *gorm.DB, Id uint64) (IsApprove bool) {
	var Number string
	DB.Table(db_var.TableName.InvPurchaseRequest).Select("number").Where("id = ? AND is_user_approved1=1", Id).Limit(1).Scan(&Number)
	IsApprove = Number != ""
	return
}

func IsPurchaseRequestApproved12(DB *gorm.DB, Id uint64) (IsApprove bool) {
	var Number string
	DB.Table(db_var.TableName.InvPurchaseRequest).Select("number").Where("id = ? AND is_user_approved12=1", Id).Limit(1).Scan(&Number)
	IsApprove = Number != ""
	return
}

func IsPurchaseRequestApproved2(DB *gorm.DB, Id uint64) (IsApprove bool) {
	var Number string
	DB.Table(db_var.TableName.InvPurchaseRequest).Select("number").Where("id = ? AND is_user_approved2=1", Id).Limit(1).Scan(&Number)
	IsApprove = Number != ""
	return
}

func IsPurchaseRequestApproved3(DB *gorm.DB, Id uint64) (IsApprove bool) {
	var Number string
	DB.Table(db_var.TableName.InvPurchaseRequest).Select("number").Where("id = ? AND is_user_approved3=1", Id).Limit(1).Scan(&Number)
	IsApprove = Number != ""
	return
}
func IsPurchaseRequestPriceApplied(DB *gorm.DB, Number string) (IsApplied bool) {
	var PRNumber string
	DB.Table(db_var.TableName.InvPurchaseRequestDetail).Select("pr_number").Where("pr_number=? AND quantity_approved > 0", Number).Limit(1).Scan(&PRNumber)
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
	DB.Table(db_var.TableName.AstUserSubDepartment).Select("user_code").Where("user_code=? AND sub_department_code=? AND is_can_inv_pr_approve1=1", UserID, SubDepartmentCode).Limit(1).Scan(&UserCode)
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
	DB.Table(db_var.TableName.AstUserSubDepartment).Select("user_code").Where("user_code=? AND sub_department_code=? AND is_can_inv_pr_approve2=1", UserID, SubDepartmentCode).Limit(1).Scan(&UserCode)
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
	DB.Table(db_var.TableName.AstUserSubDepartment).Select("user_code").Where("user_code=? AND sub_department_code=? AND is_can_inv_pr_approve3=1", UserID, SubDepartmentCode).Limit(1).Scan(&UserCode)
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
	DB.Table(db_var.TableName.AstUserSubDepartment).Select("user_code").Where("user_code=? AND sub_department_code=? AND is_can_inv_pr_approve12=1", UserID, SubDepartmentCode).Limit(1).Scan(&UserCode)
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
	DB.Table(db_var.TableName.AstUserSubDepartment).Select("user_code").Where("user_code=? AND sub_department_code=? AND is_can_inv_pr_assign_price=1", UserID, SubDepartmentCode).Limit(1).Scan(&UserCode)
	// if IsUserPasswordValid(UserID, PasswordEncrypted) {
	//TODO please Review validation
	return UserCode != ""
}

func GetMarketListCheaperPrice(DB *gorm.DB, ItemCode, UOMCode string) map[string]interface{} {
	DataOutput2 := make(map[string]interface{})
	var DataOutput map[string]interface{}
	DB.Table(db_var.TableName.InvCfgInitMarketList).Select(
		" inv_cfg_init_market_list.company_code,"+
			" IF(inv_cfg_init_market_list.uom_code=inv_cfg_init_item.uom_code,IFNULL(inv_cfg_init_market_list.price,0)*IFNULL(inv_cfg_init_item_uom1.quantity,1),"+"IF(IFNULL(inv_cfg_init_item_uom.quantity,0)=0,0,IFNULL(inv_cfg_init_market_list.price,0)/inv_cfg_init_item_uom.quantity*IFNULL(inv_cfg_init_item_uom1.quantity,0))) AS MarketPrice,"+
			" IFNULL(company.name,'') AS name").
		Joins("INNER JOIN company ON (inv_cfg_init_market_list.company_code = company.code)").
		Joins("LEFT OUTER JOIN inv_cfg_init_item_uom ON (inv_cfg_init_item_uom.item_code = ? AND inv_cfg_init_market_list.uom_code = inv_cfg_init_item_uom.uom_code)", ItemCode).
		Joins("LEFT OUTER JOIN inv_cfg_init_item_uom inv_cfg_init_item_uom1 ON (inv_cfg_init_item_uom1.item_code = ? AND inv_cfg_init_item_uom1.uom_code = ?)", ItemCode, UOMCode).
		Joins("LEFT OUTER JOIN inv_cfg_init_item ON (inv_cfg_init_market_list.item_code = inv_cfg_init_item.code)").
		Where("inv_cfg_init_market_list.item_code=?", ItemCode).Limit(1).Scan(&DataOutput)
	////fmt.Println(DataOutput)
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
		DataOutput2["price"] = general.StrToFloat64(DataOutput["MarketPrice"].(string))
	}
	return DataOutput2
}

func GetReceiveLastPrice(DB *gorm.DB, ItemCode, UOMCode, StockDate string) map[string]interface{} {
	DataOutput2 := make(map[string]interface{})
	var DataOutput map[string]interface{}
	DB.Table(db_var.TableName.InvReceivingDetail).Select(
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
	////fmt.Println("uom", UOMCode)
	////fmt.Println("2", DataOutput2)
	////fmt.Println("3", DataOutput)
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
			DataOutput2["price"] = ConvertPriceFromBasic(DB, ItemCode, UOMCode, general.StrToFloat64(BasicPrice))
		}
	}
	return DataOutput2
}

func ConvertPriceFromBasic(DB *gorm.DB, ItemCode, UOMCode string, BasicPrice float64) (Price float64) {
	Price = BasicPrice
	//Ambil dari Receive
	Quantity := master_data.GetFieldFloatQuery(DB, "SELECT quantity FROM inv_cfg_init_item_uom"+
		" WHERE item_code=?"+
		" AND uom_code=?", 0, ItemCode, UOMCode)
	if Quantity > 0 {
		Price = BasicPrice * Quantity
	}
	return
}

func GetStockStoreUpdateStockTransfer(DB *gorm.DB, Dataset *global_var.TDataset, StoreCode, ItemCode, StockTransferNumber string, StockDate time.Time) (float64, error) {
	StockDateStr := general.FormatDate1(StockDate)
	StockInDate := GetStockInDate(DB, StoreCode, ItemCode, StockDate)
	StockInDateStr := general.FormatDate1(StockInDate)
	var Stock float64
	if Dataset.ProgramConfiguration.CostingMethod == global_var.InventoryCostingMethod.Average {
		if err := DB.Raw(
			"SELECT SUM(IFNULL(Stock.Quantity, 0)) FROM ((" +
				"SELECT" +
				" SUM(IF(inv_receiving_detail.store_code='" + StoreCode + "', inv_receiving_detail.basic_quantity, 0)) AS Quantity " +
				"FROM" +
				" inv_receiving_detail" +
				" LEFT OUTER JOIN inv_receiving ON (inv_receiving_detail.receive_number = inv_receiving.number)" +
				" WHERE inv_receiving_detail.item_code='" + ItemCode + "' " +
				" AND inv_receiving.date<='" + general.FormatDate1(StockDate) + "' " +
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
				" AND inv_costing_detail.store_code='" + StoreCode + "' " +
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
	StockDateStr = general.FormatDate1(StockDate)
	CostingMethod := Dataset.ProgramConfiguration.CostingMethod
	if CostingMethod == global_var.InventoryCostingMethod.Average {
		StockInDate = GetStockInDate(DB, StoreCode, ItemCode, StockDate)
		StockInDateStr = general.FormatDate1(StockInDate)

		LastItemClosedDate = GetLastItemStoreClosedDate(DB, ItemCode)
		LastItemClosedDateStr = general.FormatDate1(LastItemClosedDate)
		IncDay1ItemClosedDateStr = general.FormatDate1(LastItemClosedDate.AddDate(0, 0, 1))

		Query1 := DB.Table(db_var.TableName.InvReceivingDetail).
			Select("SUM(IF(inv_receiving_detail.store_code=?, inv_receiving_detail.basic_quantity, 0)) AS Quantity").
			Joins("LEFT OUTER JOIN inv_receiving ON (inv_receiving_detail.receive_number = inv_receiving.number)").
			Where("inv_receiving_detail.item_code=?", ItemCode).
			Where("(inv_receiving.date BETWEEN ? AND ?)", IncDay1ItemClosedDateStr, StockDateStr).
			Group("inv_receiving_detail.item_code")
		Query2 := DB.Table(db_var.TableName.InvStockTransferDetail).
			Select("SUM(inv_stock_transfer_detail.quantity) AS Quantity").
			Joins("LEFT OUTER JOIN inv_stock_transfer ON (inv_stock_transfer_detail.st_number = inv_stock_transfer.number)").
			Where("inv_stock_transfer_detail.item_code=?", ItemCode).
			Where("(inv_stock_transfer.date BETWEEN ? AND ?)", IncDay1ItemClosedDateStr, StockDateStr).
			Where("inv_stock_transfer_detail.to_store_code=?", StoreCode).
			Group("inv_stock_transfer_detail.item_code")
		Query3 := DB.Table(db_var.TableName.InvStockTransferDetail).
			Select("-SUM(inv_stock_transfer_detail.quantity) AS Quantity").
			Joins("LEFT OUTER JOIN inv_stock_transfer ON (inv_stock_transfer_detail.st_number = inv_stock_transfer.number)").
			Where("inv_stock_transfer_detail.item_code=?", ItemCode).
			Where("(inv_stock_transfer.date BETWEEN ? AND ?)", IncDay1ItemClosedDateStr, StockDateStr).
			Where("inv_stock_transfer_detail.from_store_code=?", StoreCode).
			Group("inv_stock_transfer_detail.item_code")
		Query4 := DB.Table(db_var.TableName.InvCostingDetail).
			Select("-SUM(inv_costing_detail.quantity) AS Quantity").
			Joins("LEFT OUTER JOIN inv_costing ON (inv_costing_detail.costing_number = inv_costing.number)").
			Where("inv_costing_detail.item_code=?", ItemCode).
			Where("inv_costing.date<=?", StockInDateStr).
			Where("inv_costing_detail.store_code=?", StoreCode).
			Group("inv_costing_detail.item_code")
		Query5 := DB.Table(db_var.TableName.InvCloseSummaryStore).
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

		if err := DB.Table(db_var.TableName.InvReceivingDetail).
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
	StockDateStr = general.FormatDate1(StockDate)
	CostingMethod := Dataset.Configuration[global_var.ConfigurationCategory.Inventory][global_var.ConfigurationName.CostingMethod].(string)
	if CostingMethod == global_var.InventoryCostingMethod.Average {
		StockInDate = GetStockInDate(DB, StoreCode, ItemCode, StockDate)
		StockInDateStr = general.FormatDate1(StockInDate)

		LastItemClosedDate = GetLastItemStoreClosedDate(DB, ItemCode)
		LastItemClosedDateStr = general.FormatDate1(LastItemClosedDate)
		IncDay1ItemClosedDateStr = general.FormatDate1(LastItemClosedDate.AddDate(0, 0, 1))

		Query1 := DB.Table(db_var.TableName.InvReceivingDetail).
			Select("SUM(IF(inv_receiving_detail.store_code=?, inv_receiving_detail.basic_quantity, 0)) AS Quantity").
			Joins("LEFT OUTER JOIN inv_receiving ON (inv_receiving_detail.receive_number = inv_receiving.number)").
			Where("inv_receiving_detail.item_code=?", ItemCode).
			Where("(inv_receiving.date BETWEEN ? AND ?)", IncDay1ItemClosedDateStr, StockDateStr).
			Group("inv_receiving_detail.item_code")
		Query2 := DB.Table(db_var.TableName.InvStockTransferDetail).
			Select("SUM(inv_stock_transfer_detail.quantity) AS Quantity ").
			Joins("LEFT OUTER JOIN inv_stock_transfer ON (inv_stock_transfer_detail.st_number = inv_stock_transfer.number)").
			Where("inv_stock_transfer_detail.item_code=?", ItemCode).
			Where("(inv_stock_transfer.date BETWEEN ? AND ?)", IncDay1ItemClosedDateStr, StockDateStr).
			Where("inv_stock_transfer_detail.to_store_code=?").
			Group("inv_stock_transfer_detail.item_code")
		Query3 := DB.Table(db_var.TableName.InvStockTransferDetail).
			Select("-SUM(inv_stock_transfer_detail.quantity) AS Quantity").
			Joins("LEFT OUTER JOIN inv_stock_transfer ON (inv_stock_transfer_detail.st_number = inv_stock_transfer.number)").
			Where("inv_stock_transfer_detail.item_code=?", ItemCode).
			Where("(inv_stock_transfer.date BETWEEN ? AND ?)", IncDay1ItemClosedDateStr, StockDateStr).
			Where("inv_stock_transfer_detail.from_store_code=?", StoreCode).
			Group("inv_stock_transfer_detail.item_code")
		Query4 := DB.Table(db_var.TableName.InvCostingDetail).
			Select("-SUM(inv_costing_detail.quantity) AS Quantity").
			Joins("LEFT OUTER JOIN inv_costing ON (inv_costing_detail.costing_number = inv_costing.number)").
			Where("inv_costing_detail.item_code=?", ItemCode).
			Where("inv_costing.date<=?", StockInDateStr).
			Where("inv_costing_detail.store_code=?", StoreCode).
			Where("inv_costing_detail.costing_number<>?", CostingNumber).
			Group("inv_costing_detail.item_code")
		Query5 := DB.Table(db_var.TableName.InvCloseSummaryStore).
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

		if err := DB.Table(db_var.TableName.InvReceivingDetail).
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
	DB.Table(db_var.TableName.InvCloseSummaryStore).Select("date").Where("item_code=?", ItemCode).Limit(1).Scan(&CloseDate)
	return
}

func GetLastInventoryClosedDate(DB *gorm.DB) (CloseDate time.Time) {
	DB.Table(db_var.TableName.InvCloseLog).Select("closed_date").Limit(1).Scan(&CloseDate)
	return
}

func GetStockInDate(DB *gorm.DB, StoreCode, ItemCode string, StockDate time.Time) time.Time {
	var StockDateStr string
	var ReceiveDate, TransferInDate, CostingDate time.Time
	var ReceiveDateQuery, TransferInDateQuery, CostingDateQuery time.Time
	ReceiveDate = StockDate
	TransferInDate = StockDate
	CostingDate = StockDate
	StockDateStr = general.FormatDate1(StockDate)

	DB.Table(db_var.TableName.InvReceivingDetail).Select("inv_receiving.`date`").
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

	DB.Table(db_var.TableName.InvStockTransferDetail).Select("inv_stock_transfer.`date`").
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

		DB.Table(db_var.TableName.InvCostingDetail).Select("inv_costing.`date`").
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
	DB.Table(db_var.TableName.InvCfgInitItem).Select("inv_cfg_init_item_category.journal_account_code").
		Joins("LEFT OUTER JOIN inv_cfg_init_item_category ON (inv_cfg_init_item.category_code = inv_cfg_init_item_category.code)").
		Where("inv_cfg_init_item.code=?", ItemCode).
		Limit(1).
		Scan(&JournalAccountCode)

	return
}

func GetStoreID(DB *gorm.DB, Code string) (ID uint64) {
	DB.Table(db_var.TableName.InvCfgInitStore).Select("id").
		Where("code=?", Code).
		Limit(1).
		Scan(&ID)

	return
}

func IsStockMinusStockTransfer(DB *gorm.DB, StoreCode, ItemCode, StockTransferNumber, CostingMethod string, StockDate time.Time) (IsMinus bool) {
	Query1 := DB.Table(db_var.TableName.InvCostingDetail).Distinct("inv_costing.`date").
		Joins("LEFT OUTER JOIN inv_costing ON (inv_costing_detail.costing_number = inv_costing.number)").
		Where("inv_costing_detail.item_code=?", ItemCode).
		Where("inv_costing_detail.store_code=?", StoreCode).
		Where("inv_costing.date>=?", StockDate)
	Query2 := DB.Table(db_var.TableName.InvStockTransferDetail).Distinct("inv_stock_transfer.`date").
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
	StockDateStr = general.FormatDate1(StockDate)
	if CostingMethod == global_var.InventoryCostingMethod.Average {
		StockInDate = GetStockInDate(DB, StoreCode, ItemCode, StockDate)
		StockInDateStr = general.FormatDate1(StockInDate)

		Query1 := DB.Table(db_var.TableName.InvReceivingDetail).
			Select("SUM(IF(inv_receiving_detail.store_code=?, inv_receiving_detail.basic_quantity, 0)) AS Quantity").
			Joins("LEFT OUTER JOIN inv_receiving ON (inv_receiving_detail.receive_number = inv_receiving.number)").
			Where("inv_receiving_detail.item_code=?", ItemCode).
			Where("(inv_receiving.date<=?", StockDateStr).
			Where("inv_receiving_detail.receive_number<>?", ReceiveNumber).
			Group("inv_receiving_detail.item_code")
		Query2 := DB.Table(db_var.TableName.InvStockTransferDetail).
			Select("SUM(inv_stock_transfer_detail.quantity) AS Quantity").
			Joins("LEFT OUTER JOIN inv_stock_transfer ON (inv_stock_transfer_detail.st_number = inv_stock_transfer.number)").
			Where("inv_stock_transfer_detail.item_code=?", ItemCode).
			Where("inv_stock_transfer.date<=?", StockDateStr).
			Where("inv_stock_transfer_detail.to_store_code=?", StoreCode).
			Where("inv_stock_transfer.number<>?", ReceiveNumber).
			Group("inv_stock_transfer_detail.item_code")
		Query3 := DB.Table(db_var.TableName.InvStockTransferDetail).
			Select("-SUM(inv_stock_transfer_detail.quantity) AS Quantity").
			Joins("LEFT OUTER JOIN inv_stock_transfer ON (inv_stock_transfer_detail.st_number = inv_stock_transfer.number)").
			Where("inv_stock_transfer_detail.item_code=?", ItemCode).
			Where("inv_stock_transfer.date<=?", StockDateStr).
			Where("inv_stock_transfer_detail.from_store_code=?", StoreCode).
			Group("inv_stock_transfer_detail.item_code")
		Query4 := DB.Table(db_var.TableName.InvCostingDetail).
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
		DetailTableName := db_var.TableName.InvStockTransferDetail
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

	StockDateStr = general.FormatDate1(StockDate)
	DB.Table(db_var.TableName.InvReceivingDetail).Select("(SUM(IF(inv_receiving_detail.store_code='"+StoreCode+"', inv_receiving_detail.basic_quantity, 0)) + SUM(IFNULL(StockTransfer.Quantity, 0)) - SUM(IFNULL(Costing.Quantity, 0))) AS Quantity").
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
	IsLocked = master_data.GetFieldString(DB, db_var.TableName.Reservation, "is_lock", "number", ReservationNumber, "") == "1"
	return
}

func IsThereCardActiveReservation(DB *gorm.DB, ReservationNumber interface{}) (IsThereCardActive bool) {
	IsThereCardActive = master_data.GetFieldString(DB, db_var.TableName.LogKeylock, "id", "is_active='1' AND reservation_number", ReservationNumber, "") != ""
	return
}

func IsGroupAlreadyUsed(DB *gorm.DB, GroupCode string) bool {
	IsUsedInReservation := master_data.GetFieldBool(DB, db_var.TableName.Reservation, "number", "group_code", GroupCode, false)
	IsUsedInFolio := master_data.GetFieldBool(DB, db_var.TableName.Folio, "number", "group_code", GroupCode, false)

	return IsUsedInReservation || IsUsedInFolio

}

func IsRoomOccupiedNow(c *gin.Context, DB *gorm.DB, RoomNumber string) bool {
	if master_data.GetConfigurationBool(DB, global_var.SystemCode.General, global_var.ConfigurationCategory.General, global_var.ConfigurationName.IsRoomByName, false) {
		RoomNumber = GetRoomNumberFromConfigurationIsRoomName(c, DB, RoomNumber)
	}

	var Number uint64
	DB.Table(db_var.TableName.Folio).Select("folio.number").
		Joins("LEFT OUTER JOIN guest_detail ON (folio.guest_detail_id = guest_detail.id)").
		Where("guest_detail.room_number=?", RoomNumber).
		Where("folio.status_code=?", global_var.FolioStatus.Open).
		Limit(1).
		Scan(&Number)

	return Number > 0
}

func GetRoomBlockStatus(c *gin.Context, DB *gorm.DB, RoomNumber string) string {
	if master_data.GetConfigurationBool(DB, global_var.SystemCode.General, global_var.ConfigurationCategory.General, global_var.ConfigurationName.IsRoomByName, false) {
		RoomNumber = GetRoomNumberFromConfigurationIsRoomName(c, DB, RoomNumber)
	}
	var Status sql.NullString
	DB.Table(db_var.TableName.CfgInitRoom).Select("block_status_code").
		Where("number=?", RoomNumber).
		Limit(1).
		Scan(&Status)
	return Status.String
}

func IsRoomBlockedNow(c *gin.Context, DB *gorm.DB, RoomNumber string) bool {
	if master_data.GetConfigurationBool(DB, global_var.SystemCode.General, global_var.ConfigurationCategory.General, global_var.ConfigurationName.IsRoomByName, false) {
		RoomNumber = GetRoomNumberFromConfigurationIsRoomName(c, DB, RoomNumber)
	}

	return GetRoomBlockStatus(c, DB, RoomNumber) != ""
}

func IsCheckIn(DB *gorm.DB, ReservationNumber uint64) bool {
	var Number uint64
	DB.Table(db_var.TableName.Folio).Select("reservation_number").
		Where("reservation_number=?", ReservationNumber).
		Limit(1).
		Scan(&Number)

	return Number > 0
}

func IsRoomReady(c *gin.Context, DB *gorm.DB, RoomNumber string) bool {
	if master_data.GetConfigurationBool(DB, global_var.SystemCode.General, global_var.ConfigurationCategory.General, global_var.ConfigurationName.IsRoomByName, false) {
		RoomNumber = GetRoomNumberFromConfigurationIsRoomName(c, DB, RoomNumber)
	}
	var Status string
	DB.Table(db_var.TableName.CfgInitRoom).Select("status_code").
		Where("number=?", RoomNumber).
		Where("status_code=?", global_var.RoomStatus.Ready).
		Limit(1).
		Scan(&Status)
	return Status != ""
}

func GetMemberCodeFromGuestProfile(DB *gorm.DB, GuestProfileID uint64) (Code string) {
	DB.Table(db_var.TableName.Member).Select("code").
		Where("guest_profile_id=?", GuestProfileID).
		Limit(1).
		Scan(&Code)
	return
}

func GetMemberIDFromGuestProfile(DB *gorm.DB, MemberCode string) (ID uint64) {
	DB.Table(db_var.TableName.Member).Select("guest_profile_id").
		Where("code=?", MemberCode).
		Limit(1).
		Scan(&ID)
	return
}

func IsScheduledRate(ctx context.Context, DB *gorm.DB, FolioNumber uint64, ADate time.Time) bool {
	var ID uint64
	DB.WithContext(ctx).Table(db_var.TableName.GuestScheduledRate).Select("id").
		Where("folio_number=?", FolioNumber).
		Where("from_date<=?", general.FormatDate1(ADate)).
		Where("to_date>=?", general.FormatDate1(ADate)).
		Limit(1).
		Scan(&ID)
	return ID > 0
}

func GetFolioComplimentHU(DB *gorm.DB, FolioNumber uint64) (ComplimentHU string) {
	DB.Table(db_var.TableName.Folio).Select("compliment_hu").Where("number=?", FolioNumber).Take(&ComplimentHU)
	return
}

func GetScheduledRateComplimentHU(ctx context.Context, DB *gorm.DB, FolioNumber uint64, PostingDate time.Time) (ComplimentHU string) {
	ctx, span := global_var.Tracer.Start(ctx, "GetScheduledRateComplimentHU")
	defer span.End()

	DB.Table(db_var.TableName.GuestScheduledRate).Select("compliment_hu").
		Where("folio_number=?", FolioNumber).
		Where("from_date<=?", general.FormatDate1(PostingDate)).
		Where("to_date>=?", general.FormatDate1(PostingDate)).
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
	Result := general.IncDay(general.DateOf(ArrivalDate), GetVoucherNights(ctx, DB, VoucherNumber)).Unix() > AuditDate.Unix()
	return Result
}

func GetVoucherNights(ctx context.Context, DB *gorm.DB, VoucherNumber string) (Night int) {
	DB.WithContext(ctx).Table(db_var.TableName.Voucher).Select("nights").Where("number=?", VoucherNumber).Limit(1).Scan(&Night)
	return
}

func GetVoucherType(ctx context.Context, DB *gorm.DB, VoucherNumber string) (Type string) {
	DB.WithContext(ctx).Table(db_var.TableName.Voucher).Select("type_code").Where("number=?", VoucherNumber).Limit(1).Scan(&Type)
	return
}

func GetScheduledRoomRateCode(DB *gorm.DB, FolioNumber uint64, PostingDate time.Time) (RoomRateCode string) {
	DB.Table(db_var.TableName.GuestScheduledRate).Select("room_rate_code").
		Where("folio_number=?", FolioNumber).
		Where("from_date<=?", general.FormatDate1(PostingDate)).
		Where("to_date>=?", general.FormatDate1(PostingDate)).
		Limit(1).Scan(&RoomRateCode)
	return
}

func GetRoomRateTaxAndServiceCode(DB *gorm.DB, RoomRateCode string) (Code string) {
	DB.Table(db_var.TableName.CfgInitRoomRate).Select("tax_and_service_code").Where("code=?", RoomRateCode).Limit(1).Scan(&Code)
	return
}

func GetScheduledRate(DB *gorm.DB, FolioNumber uint64, PostingDate time.Time) (RoomRate float64) {
	DB.Table(db_var.TableName.GuestScheduledRate).Select("rate").
		Where("folio_number=?", FolioNumber).
		Where("from_date<=?", general.FormatDate1(PostingDate)).
		Where("to_date>=?", general.FormatDate1(PostingDate)).
		Limit(1).Scan(&RoomRate)
	return
}

func GetVoucherPrice(DB *gorm.DB, VoucherNumber string) (Price float64) {
	DB.Table(db_var.TableName.Voucher).Select("price").Where("number=?", VoucherNumber).Limit(1).Scan(&Price)
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
		GroupCode               string    `json:"group_code"`
		CurrencyCode            string    `json:"currency_code"`
		AccountCode             string    `json:"account_code"`
		ExchangeRate            float64   `json:"exchange_rate"`
		PurposeOfCode           string    `json:"purpose_of_code"`
		SalesCode               string    `json:"sales_code"`
		VoucherNumber           string    `json:"voucher_number"`
		Notes                   string    `json:"notes"`
		ComplimentHU            string    `json:"compliment_hu"`
		DateArrival             time.Time `json:"DateArrival"`
		TitleCode               string    `json:"title_code"`
		FullName                string    `json:"full_name"`
		Street                  string    `json:"street"`
		CityCode                string    `json:"city_code"`
		City                    string    `json:"city"`
		NationalityCode         string    `json:"nationality_code"`
		CountryCode             string    `json:"country_code"`
		StateCode               string    `json:"state_code"`
		PostalCode              string    `json:"postal_code"`
		Phone1                  string    `json:"phone1"`
		Phone2                  string    `json:"phone2"`
		Fax                     string    `json:"fax"`
		Email                   string    `json:"email"`
		Website                 string    `json:"website"`
		CompanyCode             string    `json:"company_code"`
		GuestTypeCode           string    `json:"guest_type_code"`
		CustomLookupFieldCode01 string    `json:"custom_lookup_field_code01"`
		CustomLookupFieldCode02 string    `json:"custom_lookup_field_code02"`
		Adult                   int       `json:"adult"`
		Child                   int       `json:"child"`
		RoomTypeCode            string    `json:"room_type_code"`
		BedTypeCode             string    `json:"bed_type_code"`
		RoomNumber              string    `json:"room_number"`
		RoomRateCode            string    `json:"room_rate_code"`
		WeekdayRate             float64   `json:"weekday_rate"`
		WeekendRate             float64   `json:"weekend_rate"`
		DiscountPercent         uint8     `json:"discount_percent"`
		Discount                float64   `json:"discount"`
		BusinessSourceCode      string    `json:"business_source_code"`
		CommissionTypeCode      string    `json:"commission_type_code"`
		CommissionValue         float64   `json:"commission_value"`
		PaymentTypeCode         string    `json:"payment_type_code"`
		MarketCode              string    `json:"market_code"`
		BookingSourceCode       string    `json:"booking_source_code"`
		TaxAndServiceCode       string    `json:"tax_and_service_code"`
		ChargeFrequencyCode     string    `json:"charge_frequency_code"`
	}

	Timezone := Dataset.ProgramConfiguration.Timezone
	var DataOutput FolioStruct
	err = DB.Table(db_var.TableName.Folio).Select(
		" folio.group_code,"+
			" guest_detail.currency_code,"+
			" guest_detail.exchange_rate,"+
			" guest_general.purpose_of_code,"+
			" guest_general.sales_code,"+
			" folio.voucher_number,"+
			" guest_general.notes,"+
			" folio.compliment_hu,"+
			" DATE(convert_tz(guest_detail.arrival,'UTC','"+Timezone+"')) AS DateArrival,"+
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
			" cfg_init_room.revenue_account_code AS account_code,"+
			" cfg_init_room_rate.tax_and_service_code,"+
			" cfg_init_room_rate.charge_frequency_code ").
		Joins("LEFT JOIN contact_person ON (folio.contact_person_id1 = contact_person.id)").
		Joins("LEFT JOIN guest_detail ON (folio.guest_detail_id = guest_detail.id)").
		Joins("LEFT JOIN guest_general ON (folio.guest_general_id = guest_general.id)").
		Joins("LEFT OUTER JOIN cfg_init_room ON (guest_detail.room_number = cfg_init_room.number)").
		Joins("LEFT JOIN cfg_init_room_rate ON (guest_detail.room_rate_code = cfg_init_room_rate.code)").
		Where("folio.number=?", FolioNumber).Take(&DataOutput).Error
	if err != nil {
		return
	}
	//fmt.Println("TaxAndServiceCode", DataOutput.TaxAndServiceCode)
	var ScheduledRateData db_var.Guest_scheduled_rate
	if err := DB.Table(db_var.TableName.GuestScheduledRate).
		Where("folio_number=?", FolioNumber).
		Where("DATE(from_date)<=?", general.FormatDate1(PostingDate)).
		Where("DATE(to_date)>=?", general.FormatDate1(PostingDate)).
		Limit(1).
		Scan(&ScheduledRateData).Error; err != nil {
		return 0, err
	}

	// if err == nil {
	IsCanPostRoomCharge := IsCanPostCharge(c, DB, DataOutput.ChargeFrequencyCode, DataOutput.DateArrival)
	IsBreakfastX = IsFolioHaveBreakfast(DB, Dataset, FolioNumber)
	IsVoucherActiveX = IsVoucherComplimentStillActive(ctx, DB, DataOutput.VoucherNumber, DataOutput.DateArrival, PostingDate)
	VoucherTypeCode = GetVoucherType(ctx, DB, DataOutput.VoucherNumber)
	if IsVoucherActiveX && (VoucherTypeCode == global_var.VoucherType.Compliment) {
		ComplimentHU = global_var.RoomStatus.Compliment
	} else {
		ComplimentHU = DataOutput.ComplimentHU
		if ScheduledRateData.Id > 0 {
			ComplimentHU = ScheduledRateData.ComplimentHu
		}
	}

	SDFrontOffice := Dataset.GlobalSubDepartment.FrontOffice

	GARoomCharge := Dataset.GlobalAccount.RoomCharge
	if DataOutput.AccountCode != "" {
		GARoomCharge = DataOutput.AccountCode
	}
	if !(!AllowZeroAmount && (ComplimentHU == global_var.RoomStatus.HouseUseX || ComplimentHU == global_var.RoomStatus.Compliment)) {
		if IsCanPostRoomCharge {
			RoomNumber = DataOutput.RoomNumber

			CurrencyCode = DataOutput.CurrencyCode
			ExchangeRate = DataOutput.ExchangeRate
			if ExchangeRate <= 0 {
				ExchangeRate = 1
			}

			IsScheduledRateX = ScheduledRateData.Id > 0
			if IsScheduledRateX && !IsVoucherActiveX {
				RoomRateCode = ScheduledRateData.RoomRateCode
				if RoomRateCode == "" {
					RoomRateCode = DataOutput.RoomRateCode
					RoomRateTaxServiceCode = DataOutput.TaxAndServiceCode
				} else {
					RoomRateTaxServiceCode = GetRoomRateTaxAndServiceCode(DB, RoomRateCode)
				}

				RoomChargeB4Breakdown = *ScheduledRateData.Rate / ExchangeRate
				RoomRateAmount = RoomChargeB4Breakdown

				RoomRateAmountOriginal = GetRoomRateAmount(ctx, DB, Dataset, DataOutput.RoomRateCode, general.FormatDate1(PostingDate), DataOutput.Adult, DataOutput.Child, false)
				ComplimentHU = ScheduledRateData.ComplimentHu
				Discount = 0
			} else {
				RoomRateCode = DataOutput.RoomRateCode
				RoomRateTaxServiceCode = DataOutput.TaxAndServiceCode
				if general.IsWeekend(PostingDate, Dataset) {
					RoomChargeB4Breakdown = DataOutput.WeekendRate
				} else {
					RoomChargeB4Breakdown = DataOutput.WeekdayRate
				}
				RoomRateAmount = RoomChargeB4Breakdown
				RoomRateAmountOriginal = GetRoomRateAmount(ctx, DB, Dataset, DataOutput.RoomRateCode, general.FormatDate1(PostingDate), DataOutput.Adult, DataOutput.Child, false)
				ComplimentHU = DataOutput.ComplimentHU

				if (IsVoucherActiveX) && (VoucherTypeCode != global_var.VoucherType.Compliment) {
					Discount = GetVoucherPrice(DB, DataOutput.VoucherNumber)
				} else {
					if DataOutput.DiscountPercent > 0 {
						Discount = general.RoundTo(RoomChargeB4Breakdown * DataOutput.Discount / 100)
					} else {
						Discount = DataOutput.Discount
					}
				}

				if !Dataset.ProgramConfiguration.PostDiscount {
					RoomChargeB4Breakdown = RoomChargeB4Breakdown - Discount
					Discount = 0
				}
			}

			if (ComplimentHU == "H" || ComplimentHU == "P") && (RoomChargeB4Breakdown > 0) {
				if AllowZeroAmount {
					Result = 0
					//Post Room Charge
					CorrectionBreakdown = GetSubFolioCorrectionBreakdown(c, DB)
					BreakDown1 = GetSubFolioBreakdown1(c, DB)
					_, _, err = InsertSubFolio2(c, DB, Dataset, FolioNumber, SubFolioGroupCode, RoomNumber, SDFrontOffice, GARoomCharge, GARoomCharge, "", "", CurrencyCode, "Auto Room Charge", "", "", global_var.TransactionType.Debit, "", "", CorrectionBreakdown, BreakDown1, global_var.SubFolioPostingType.Room, 0, 0, 0, 0, ExchangeRate, AllowZeroAmount, true, UserID)
					if err != nil {
						return
					}
					if IsInHousePosted(DB, PostingDate, FolioNumber) {
						err = DeleteGuestInHouse(ctx, DB, PostingDate, FolioNumber)
						if err != nil {
							return
						}
					}
					err = InsertGuestInHouse(DB, Dataset, PostingDate, FolioNumber, DataOutput.GroupCode, DataOutput.RoomTypeCode, DataOutput.BedTypeCode, DataOutput.RoomNumber, DataOutput.RoomRateCode, DataOutput.BusinessSourceCode, DataOutput.CommissionTypeCode, DataOutput.PaymentTypeCode, DataOutput.MarketCode,
						DataOutput.TitleCode, DataOutput.FullName, DataOutput.Street, DataOutput.City, DataOutput.CityCode, DataOutput.CountryCode, DataOutput.StateCode, DataOutput.PostalCode, DataOutput.Phone1, DataOutput.Phone2, DataOutput.Fax, DataOutput.Email,
						DataOutput.Website, DataOutput.CompanyCode, DataOutput.GuestTypeCode, DataOutput.SalesCode, DataOutput.ComplimentHU, DataOutput.Notes, DataOutput.Adult, DataOutput.Child, 0, 0, 0, DataOutput.CommissionValue, DataOutput.DiscountPercent, 0, general.BoolToUint8(IsScheduledRateX), general.BoolToUint8(IsBreakfastX), DataOutput.BookingSourceCode,
						DataOutput.PurposeOfCode, DataOutput.CustomLookupFieldCode01, DataOutput.CustomLookupFieldCode02, 0, "", DataOutput.NationalityCode, UserID)
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
				//fmt.Println("1", GARoomCharge, RoomChargeBasic, RoomChargeTax, RoomChargeService)
				//Proses Query Breakdown
				Breakdown := []db_var.Cfg_init_room_rate_breakdown{}
				if IsScheduledRateX {
					err = DB.Table(db_var.TableName.CfgInitRoomRateBreakdown).Select(
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
					err = DB.Table(db_var.TableName.GuestBreakdown).Select(
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
					// fmt.Println("Amount", breakdown.Amount)
					// fmt.Println("Quantity", breakdown.Quantity)
					// fmt.Println("PerPax", breakdown.PerPax)
					if IsCanPostCharge(c, DB, breakdown.ChargeFrequencyCode, DataOutput.DateArrival) {
						if breakdown.IsAmountPercent > 0 {
							BreakdownAmount = GetTotalBreakdownAmount(breakdown.Quantity, RoomChargeB4Breakdown*breakdown.Amount/100, breakdown.ExtraPax, breakdown.PerPax > 0, breakdown.IncludeChild > 0, breakdown.PerPaxExtra > 0, breakdown.MaxPax, DataOutput.Adult, DataOutput.Child) / ExchangeRate
							//                  BreakdownAmount = GetTotalBreakdownAmount(MyQGuestBreakdownCalculatequantity.AsFloat, RoomRateAmount * MyQGuestBreakdownCalculateamount.AsFloat/100, MyQGuestBreakdownCalculateextra_pax.AsFloat, MyQGuestBreakdownCalculateper_pax.AsVariant, MyQGuestBreakdownCalculateinclude_child.AsVariant, MyQGuestBreakdownCalculateper_pax_extra.AsVariant, MyQGuestBreakdownCalculatemax_pax.AsInteger, MyQFolioCalculateadult.AsInteger, MyQFolioCalculatechild.AsInteger) / ExchangeRate
						} else {
							BreakdownAmount = GetTotalBreakdownAmount(breakdown.Quantity, breakdown.Amount, breakdown.ExtraPax, breakdown.PerPax > 0, breakdown.IncludeChild > 0, breakdown.PerPaxExtra > 0, breakdown.MaxPax, DataOutput.Adult, DataOutput.Child) / ExchangeRate
						}
						// fmt.Println("BreakdownAmount", BreakdownAmount)

						BreakdownBasic, BreakdownTax, BreakdownService = GetBasicTaxService(DB, breakdown.AccountCode, breakdown.TaxAndServiceCode, BreakdownAmount)
						BreakdownAmount = BreakdownBasic + BreakdownTax + BreakdownService
						TotalBreakdown = TotalBreakdown + BreakdownAmount
					}
				}

				//Get Commission from Business Source
				BusinessSourceCode = DataOutput.BusinessSourceCode
				Commission = 0
				if BusinessSourceCode != "" {
					Commission = GetCommission(c, DB, DataOutput.CommissionTypeCode, DataOutput.CommissionValue, RoomChargeB4Breakdown, RoomChargeBasic, DataOutput.DateArrival) / ExchangeRate
				}
				//Room Charge - Total Breakdown - Total Commission
				RoomChargeAfterBreakdown = RoomChargeB4Breakdown - TotalBreakdown - Commission
				if (RoomChargeAfterBreakdown > 0) || (AllowZeroAmount && (RoomChargeAfterBreakdown == 0)) {
					Result = 0
					//Cari Basic dari Room Charge Bersih (yang sudah dikurangi Breakdown dan Total Commission)
					//fmt.Println("2", "RoomRateTaxServiceCode", RoomRateTaxServiceCode)
					RoomChargeBasic, RoomChargeTax, RoomChargeService = GetBasicTaxService2(DB, GARoomCharge, RoomRateTaxServiceCode, RoomChargeAfterBreakdown)
					RoomChargeAfterBreakdown = RoomChargeBasic + RoomChargeTax + RoomChargeService
					//fmt.Println("2", GARoomCharge, RoomChargeBasic, RoomChargeTax, RoomChargeService)
					//Post Room Charge
					CorrectionBreakdown = GetSubFolioCorrectionBreakdown(c, DB)
					BreakDown1 = GetSubFolioBreakdown1(c, DB)
					_, _, err = InsertSubFolio2(c, DB, Dataset, FolioNumber, SubFolioGroupCode, RoomNumber, SDFrontOffice, GARoomCharge, GARoomCharge, "", "", CurrencyCode, "Auto Room Charge", "", "", global_var.TransactionType.Debit, "", "", CorrectionBreakdown, BreakDown1, global_var.SubFolioPostingType.Room, 0, RoomChargeBasic, RoomChargeTax, RoomChargeService, ExchangeRate, AllowZeroAmount, true, UserID)
					if err != nil {
						return
					}
					if IsInHousePosted(DB, PostingDate, FolioNumber) {
						DeleteGuestInHouse(ctx, DB, PostingDate, FolioNumber)
					}
					err = InsertGuestInHouse(DB, Dataset, PostingDate, FolioNumber, DataOutput.GroupCode, DataOutput.RoomTypeCode, DataOutput.BedTypeCode, DataOutput.RoomNumber, RoomRateCode, DataOutput.BusinessSourceCode, DataOutput.CommissionTypeCode, DataOutput.PaymentTypeCode, DataOutput.MarketCode,
						DataOutput.TitleCode, DataOutput.FullName, DataOutput.Street, DataOutput.City, DataOutput.CityCode, DataOutput.CountryCode, DataOutput.StateCode, DataOutput.PostalCode, DataOutput.Phone1, DataOutput.Phone2, DataOutput.Fax, DataOutput.Email,
						DataOutput.Website, DataOutput.CompanyCode, DataOutput.GuestTypeCode, DataOutput.SalesCode, ComplimentHU, DataOutput.Notes, DataOutput.Adult, DataOutput.Child, RoomRateAmount, RoomRateAmountOriginal, DataOutput.Discount, DataOutput.CommissionValue, DataOutput.DiscountPercent, 0, general.BoolToUint8(IsScheduledRateX), general.BoolToUint8(IsBreakfastX), DataOutput.BookingSourceCode,
						DataOutput.PurposeOfCode, DataOutput.CustomLookupFieldCode01, DataOutput.CustomLookupFieldCode02, 0, "", DataOutput.NationalityCode, UserID)
					if err != nil {
						return
					}
					//Posting Breakdown
					for _, breakdown := range Breakdown {
						if IsCanPostCharge(c, DB, breakdown.ChargeFrequencyCode, DataOutput.DateArrival) {
							if breakdown.IsAmountPercent > 0 {
								BreakdownAmount = GetTotalBreakdownAmount(breakdown.Quantity, RoomChargeB4Breakdown*breakdown.Amount/100, breakdown.ExtraPax, breakdown.PerPax > 0, breakdown.IncludeChild > 0, breakdown.PerPaxExtra > 0, breakdown.MaxPax, DataOutput.Adult, DataOutput.Child) / ExchangeRate
								//                  BreakdownAmount = GetTotalBreakdownAmount(MyQGuestBreakdownCalculatequantity.AsFloat, RoomRateAmount * MyQGuestBreakdownCalculateamount.AsFloat/100, MyQGuestBreakdownCalculateextra_pax.AsFloat, MyQGuestBreakdownCalculateper_pax.AsVariant, MyQGuestBreakdownCalculateinclude_child.AsVariant, MyQGuestBreakdownCalculateper_pax_extra.AsVariant, MyQGuestBreakdownCalculatemax_pax.AsInteger, MyQFolioCalculateadult.AsInteger, MyQFolioCalculatechild.AsInteger) / ExchangeRate
							} else {
								BreakdownAmount = GetTotalBreakdownAmount(breakdown.Quantity, breakdown.Amount, breakdown.ExtraPax, breakdown.PerPax > 0, breakdown.IncludeChild > 0, breakdown.PerPaxExtra > 0, breakdown.MaxPax, DataOutput.Adult, DataOutput.Child) / ExchangeRate
							}
							_, _, err = InsertSubFolio(c, DB, Dataset, true, GARoomCharge, breakdown.TaxAndServiceCode, db_var.Sub_folio{
								FolioNumber:         FolioNumber,
								RoomNumber:          RoomNumber,
								SubDepartmentCode:   breakdown.SubDepartmentCode,
								AccountCode:         breakdown.AccountCode,
								ProductCode:         breakdown.ProductCode,
								GroupCode:           SubFolioGroupCode,
								Remark:              "Breakdown: " + breakdown.Remark,
								TypeCode:            global_var.TransactionType.Debit,
								PostingType:         global_var.SubFolioPostingType.Room,
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
						_, _, err = InsertSubFolio(c, DB, Dataset, true, GARoomCharge, "", db_var.Sub_folio{
							FolioNumber:         FolioNumber,
							RoomNumber:          RoomNumber,
							SubDepartmentCode:   SDFrontOffice,
							AccountCode:         GAAPCommission,
							GroupCode:           SubFolioGroupCode,
							Remark:              "Breakdown Commission",
							TypeCode:            global_var.TransactionType.Debit,
							PostingType:         global_var.SubFolioPostingType.Room,
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
		ExtraCharge db_var.Guest_extra_charge `gorm:"embedded"`
		GuestDetail db_var.Guest_detail       `gorm:"embedded"`
		DateArrival time.Time
	}
	var DataOutput []ExtraChargeStruct
	Query := DB.WithContext(ctx).Table(db_var.TableName.GuestExtraCharge).Select(
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
			var GuestExtraChargeBreakdown []db_var.Guest_extra_charge_breakdown
			DB.WithContext(ctx).Table(db_var.TableName.GuestExtraChargeBreakdown).Where("guest_extra_charge_id=?", extraCharge.ExtraCharge.Id).Scan(&GuestExtraChargeBreakdown)

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
					AccountCodeMaster, *extraCharge.ExtraCharge.ProductCode, *extraCharge.ExtraCharge.PackageCode, "", "Extra Charge", "", "", global_var.TransactionType.Debit,
					"", "", CorrectionBreakdown, BreakDown1, global_var.SubFolioPostingType.ExtraCharge, extraCharge.ExtraCharge.Id, AmountBasic, AmountTax, AmountService, 0, false, true, UserID)
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
						_, _, err = InsertSubFolio(c, DB, Dataset, true, AccountCodeMaster, *breakdown.TaxAndServiceCode, db_var.Sub_folio{
							FolioNumber:         FolioNumber,
							RoomNumber:          RoomNumber,
							SubDepartmentCode:   breakdown.SubDepartmentCode,
							AccountCode:         breakdown.AccountCode,
							PackageCode:         *extraCharge.ExtraCharge.PackageCode,
							ProductCode:         *breakdown.ProductCode,
							GroupCode:           SubFolioGroupCode,
							Remark:              "Extra Charge Breakdown",
							TypeCode:            global_var.TransactionType.Debit,
							PostingType:         global_var.SubFolioPostingType.ExtraCharge,
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
	DB.Table(db_var.TableName.GuestDeposit).Select("SUM(IF(type_code='"+global_var.TransactionType.Debit+"', amount_foreign, -amount_foreign)) AS TotalAmount").
		Where("correction_breakdown=?", CorrectionBreakdown).
		Where("void=0").
		Group("correction_breakdown").
		Limit(1).
		Scan(&Amount)

	return
}

func GetSubFolioBreakdownAmountForeign(DB *gorm.DB, CorrectionBreakdown int64) (Amount float64) {
	DB.Table(db_var.TableName.SubFolio).Select("SUM(IF(type_code='"+global_var.TransactionType.Debit+"',(quantity * amount_foreign),-(quantity * amount_foreign))) AS TotalAmount ").
		Where("correction_breakdown=?", CorrectionBreakdown).
		Where("void=0").
		Group("correction_breakdown").
		Limit(1).
		Scan(&Amount)

	return
}

func GetSubFolioBreakdown1AmountForeign(DB *gorm.DB, Breakdown1 int64) (Amount float64) {
	DB.Table(db_var.TableName.SubFolio).Select("SUM(IF(type_code='"+global_var.TransactionType.Debit+"',(quantity * amount_foreign),-(quantity * amount_foreign))) AS TotalAmount").
		Where("breakdown1=?", Breakdown1).
		Where("void=0").
		Group("breakdown1").
		Limit(1).
		Scan(&Amount)

	return
}

func GetSubFolioBreakdownAmountForeign2(DB *gorm.DB, CorrectionBreakdown, Breakdown2 int64) (Amount float64) {
	DB.Table(db_var.TableName.SubFolio).Select("SUM(IF(type_code='"+global_var.TransactionType.Debit+"',(quantity * amount_foreign),-(quantity * amount_foreign))) AS TotalAmount").
		Where("correction_breakdown=?", CorrectionBreakdown).
		Where("breakdown2=?", Breakdown2).
		Where("void=0").
		Group("correction_breakdown,breakdown2").
		Limit(1).
		Scan(&Amount)

	return
}

func GetSubFolioBreakdown1AmountForeign2(DB *gorm.DB, Breakdown1, Breakdown2 int64) (Amount float64) {
	DB.Table(db_var.TableName.SubFolio).Select("SUM(IF(type_code='"+global_var.TransactionType.Debit+"',(quantity * amount_foreign),-(quantity * amount_foreign))) AS TotalAmount").
		Where("breakdown1=?", Breakdown1).
		Where("breakdown2=?", Breakdown2).
		Where("void=0").
		Group("breakdown1").
		Limit(1).
		Scan(&Amount)

	return
}

func GetCheckQuantityByCorrectionBreakdown(DB *gorm.DB, CorrectionBreakdown int64) (Amount float64) {
	DB.Table(db_var.TableName.SubFolio).Select("IF(cfg_init_account_sub_group.group_code='"+global_var.GlobalAccountGroup.Charge+"', IF(sub_folio.type_code='"+global_var.TransactionType.Debit+"', sub_folio.quantity, -sub_folio.quantity), IF(sub_folio.type_code='"+global_var.TransactionType.Credit+"', sub_folio.quantity, -sub_folio.quantity)) AS TotalQuantity ").
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
	DB.Table(db_var.TableName.PosCheckTransaction).Select("id").
		Where("(sub_folio_id=? OR sub_folio_id=?) AND sub_folio_id<>0", SubFolioID, TransferPairID).
		Scan(&Result)

	return Result > 0
}

func GetAccountGroupCode(DB *gorm.DB, AccountCode string) string {
	Result := ""
	DB.Table(db_var.TableName.CfgInitAccount).Select("group_code").
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
		Result = !master_data.GetFieldBool(DB, db_var.TableName.AccApRefundDepositPaymentDetail, "id", "sub_folio_id", TransactionID, false)
	} else if AccountSubGroupCode == global_var.GlobalAccountSubGroup.AccountPayable {
		Result = !master_data.GetFieldBool(DB, db_var.TableName.AccApCommissionPaymentDetail, "id", "sub_folio_id", TransactionID, false)
	} else if (AccountSubGroupCode == global_var.GlobalAccountSubGroup.CreditDebitCard) || (AccountSubGroupCode == global_var.GlobalAccountSubGroup.BankTransfer) {
		if Mode == 0 {
			Result = !master_data.GetFieldBool(DB, db_var.TableName.AccCreditCardReconDetail, "id", "guest_deposit_id", TransactionID, false)
		} else {
			Result = !master_data.GetFieldBool(DB, db_var.TableName.AccCreditCardReconDetail, "id", "sub_folio_id", TransactionID, false)
		}
	} else if AccountSubGroupCode == global_var.GlobalAccountSubGroup.AccountReceivable {
		Result = !master_data.GetFieldBool(DB, db_var.TableName.InvoiceItem, "id", "sub_folio_id", TransactionID, false)
	} else if AccountSubGroupCode == global_var.GlobalAccountSubGroup.Payment {
		if Mode == 0 {
			Result = !master_data.GetFieldBool(DB, db_var.TableName.AccForeignCash, "id", "id_transaction", general.Uint64ToStr(TransactionID)+"' AND id_change=0 AND stock<>amount_foreign AND id_table='"+strconv.Itoa(global_var.ForeignCashTableID.GuestDeposit), false)
		} else {
			Result = !master_data.GetFieldBool(DB, db_var.TableName.InvoiceItem, "id", "sub_folio_id", general.Uint64ToStr(TransactionID)+"' AND id_change=0 AND stock<>amount_foreign AND id_table='"+strconv.Itoa(global_var.ForeignCashTableID.SubFolio), false)
		}
	} else {
		Result = true
	}
	return
}

func GetCountAutoPosting(ctx context.Context, DB *gorm.DB, FolioNumber uint64, AccountCode, PostingType string, PostingDate time.Time) (Count int) {
	Query := DB.WithContext(ctx).Table(db_var.TableName.SubFolio).Select("account_code, is_correction").
		Where("belongs_to=?", FolioNumber).
		Where("audit_date=? AND void=0", general.FormatDate1(PostingDate))
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

func GetBreakdownAutoPosting(DB *gorm.DB, FolioNumber uint64, AccountCode, PostingType string, PostingDate time.Time) (CorrectionBreakdown []uint64) {
	Query := DB.Table(db_var.TableName.SubFolio).Select("account_code, correction_breakdown").
		Where("belongs_to=?", FolioNumber).
		Where("audit_date_unixx=UNIX_TIMESTAMP(?) AND void=0", general.FormatDate2("2006-01-02 00:00:00", PostingDate)).
		Where("posting_type", PostingType).
		Group("correction_breakdown")
	DB.Table("(?) AS ExtraCharge", Query).Select("correction_breakdown").Where("account_code=?", AccountCode).Scan(&CorrectionBreakdown)

	return

}

func IsCanCancelCheckInFolio(DB *gorm.DB, Dataset *global_var.TDataset, FolioNumber uint64) bool {
	var DataOutput []db_var.Sub_folio
	DB.Table(db_var.TableName.SubFolio).Select("account_code, id").Where("folio_number=? AND void=0", FolioNumber).Scan(&DataOutput)
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
	err := DB.Table(db_var.TableName.SubFolio).Select("CONCAT(sub_folio.belongs_to, '/Room: ', guest_detail.room_number, '/', contact_person.title_code, contact_person.full_name) AS FolioDetail").
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
	DB.Table(db_var.TableName.InvoiceItem).Select("invoice_number").Where("folio_number=?", FolioNumber).Limit(1).Scan(&InvoiceNumber)
	return
}

func GetCompanyDetailListP(c *gin.Context) {
	var DataOutputCompany []map[string]interface{}

	// Get Program Configuration
	val, exist := c.Get("pConfig")
	if !exist {
		master_data.SendResponse(global_var.ResponseCode.DatabaseError, "pConfig", nil, c)
		return
	}
	pConfig := val.(*config.CompanyDataConfiguration)
	DB := pConfig.DB
	DB.Table(db_var.TableName.Company).Select("company.*").
		Joins("LEFT JOIN cfg_init_company_type ON company.type_code=cfg_init_company_type.code").
		// Where("company.is_direct_bill = '1'").
		Order("company.name").
		Find(&DataOutputCompany)

	master_data.SendResponse(global_var.ResponseCode.Successfully, "", DataOutputCompany, c)
}

func GetARCompaniesP(c *gin.Context) {
	var DataOutputCompany []map[string]interface{}

	// Get Program Configuration
	val, exist := c.Get("pConfig")
	if !exist {
		master_data.SendResponse(global_var.ResponseCode.DatabaseError, "pConfig", nil, c)
		return
	}
	pConfig := val.(*config.CompanyDataConfiguration)
	DB := pConfig.DB
	DB.Table(db_var.TableName.Company).Select("company.*").
		Joins("LEFT JOIN cfg_init_company_type ON company.type_code=cfg_init_company_type.code").
		Where("company.is_direct_bill = '1'").
		Order("company.name").
		Find(&DataOutputCompany)

	master_data.SendResponse(global_var.ResponseCode.Successfully, "", DataOutputCompany, c)
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
	if err := DB.Table(db_var.TableName.AccApAr).Select("SUM(IFNULL(amount,0) - IFNULL(amount_paid,0)) AS Outstanding").
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
		DB.Table(db_var.TableName.SubFolio).Select(
			" sub_folio.id,"+
				" (SUM(IF(sub_folio.type_code='"+global_var.TransactionType.Debit+"', sub_folio.quantity*sub_folio.amount, -(sub_folio.quantity*sub_folio.amount))) - IFNULL(Payment.TotalPaid,0)) AS Amount ").
			Joins(" LEFT OUTER JOIN ("+
				"SELECT"+
				" sub_folio_id,"+
				" SUM(amount) AS TotalPaid "+
				"FROM"+
				"  acc_ap_commission_payment_detail"+
				" WHERE sub_folio_id="+general.Uint64ToStr(SubFolioID)+
				" AND ref_number<>'"+RefNumber+"' "+
				"GROUP BY sub_folio_id) AS Payment ON (sub_folio.id = Payment.sub_folio_id) ").
			Joins(" LEFT OUTER JOIN cfg_init_account ON (sub_folio.account_code = cfg_init_account.code)").
			Where("cfg_init_account.sub_group_code=?", global_var.GlobalAccountSubGroup.AccountPayable).
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
		return global_var.TimeSegment.Breakfast
	} else if Hour > 12 && Hour < 17 {
		return global_var.TimeSegment.Lunch
	}
	return global_var.TimeSegment.Dinner
}

func GetFirstMarketCodePOS(DB *gorm.DB) (string, error) {
	code := ""
	if err := DB.Table(db_var.TableName.PosCfgInitMarket).Select("code").Order("id_sort").Limit(1).Scan(&code).Error; err != nil {
		return "", err
	}
	return code, nil
}

func GetJournalAccountPayable(c *gin.Context, DB *gorm.DB) ([]map[string]interface{}, error) {
	var DataOutput []map[string]interface{}
	if err := DB.Table(db_var.TableName.CfgInitJournalAccount).Select(
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
	if err := DB.Table(db_var.TableName.CfgInitJournalAccount).Select(
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

func GetJournalAccountExpense(c *gin.Context, DB *gorm.DB, SubDepartmentCode string) ([]map[string]interface{}, error) {
	var DataOutput []map[string]interface{}
	if err := DB.Table(db_var.TableName.CfgInitJournalAccount).Select(
		"cfg_init_journal_account.code",
		"cfg_init_journal_account.name",
		"cfg_init_journal_account_sub_group.name AS SubGroupName").
		Joins("LEFT OUTER JOIN cfg_init_journal_account_sub_group ON (cfg_init_journal_account.sub_group_code = cfg_init_journal_account_sub_group.code)").
		Where("cfg_init_journal_account.sub_department_code LIKE ?", "%"+SubDepartmentCode+"%").
		Where("(cfg_init_journal_account_sub_group.group_code=?", 6).
		Or("cfg_init_journal_account_sub_group.group_code=?", 7).
		Or("cfg_init_journal_account_sub_group.group_code=?)", 9).
		Order("cfg_init_journal_account.code").Scan(&DataOutput).Error; err != nil {
		return DataOutput, err
	}

	return DataOutput, nil
}

func GetJournalAccountCosting(c *gin.Context, DB *gorm.DB) ([]map[string]interface{}, error) {
	var DataOutput []map[string]interface{}
	if err := DB.Table(db_var.TableName.CfgInitJournalAccount).Select(
		"cfg_init_journal_account.code",
		"cfg_init_journal_account.name",
		"cfg_init_journal_account_sub_group.name AS SubGroupName").
		Joins("LEFT OUTER JOIN cfg_init_journal_account_sub_group ON (cfg_init_journal_account.sub_group_code = cfg_init_journal_account_sub_group.code)").
		Where("cfg_init_journal_account_sub_group.group_code=?", 5).
		Order("cfg_init_journal_account.code").Scan(&DataOutput).Error; err != nil {
		return DataOutput, err
	}

	return DataOutput, nil
}

func GetBankAccountReceive(c *gin.Context, DB *gorm.DB) ([]map[string]interface{}, error) {
	var DataOutput []map[string]interface{}
	if err := DB.Table(db_var.TableName.AccCfgInitBankAccount).Select(
		"code", "name", "journal_account_code", "type_code", "bank_account_number").
		Where("for_receive=1").
		Order("journal_account_code").Scan(&DataOutput).Error; err != nil {
		return DataOutput, err
	}

	return DataOutput, nil
}

func GetBankAccountPayment(c *gin.Context, DB *gorm.DB) ([]map[string]interface{}, error) {
	var DataOutput []map[string]interface{}
	if err := DB.Table(db_var.TableName.AccCfgInitBankAccount).Select(
		"code", "name", "journal_account_code", "type_code", "bank_account_number").
		Where("for_payment=1").
		Order("journal_account_code").Scan(&DataOutput).Error; err != nil {
		return DataOutput, err
	}

	return DataOutput, nil
}

func GetCompany(c *gin.Context, DB *gorm.DB) ([]map[string]interface{}, error) {
	var DataOutput []map[string]interface{}
	if err := DB.Table(db_var.TableName.Company).Select(
		"code", "name").
		Order("name").Scan(&DataOutput).Error; err != nil {
		return DataOutput, err
	}

	return DataOutput, nil
}

func IsAPPaid(c *gin.Context, DB *gorm.DB, Number string) (IsPaid uint8, err error) {
	if err := DB.Table(db_var.TableName.AccApAr).
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
	DB.Table(db_var.TableName.InvReceivingDetail).Select("basic_quantity, quantity, total_price").Where("id=?", ReceiveID).Limit(1).Scan(&DataOutput)
	if Quantity == DataOutput.Quantity {
		if DataOutput.BasicQuantity > DataOutput.Quantity {
			return general.RoundToX3(DataOutput.TotalPrice - ((DataOutput.TotalPrice / DataOutput.BasicQuantity) * (DataOutput.BasicQuantity - DataOutput.Quantity)))
		} else {
			return general.RoundToX3(DataOutput.TotalPrice)
		}
	} else {
		return general.RoundToX3((DataOutput.TotalPrice / DataOutput.BasicQuantity)) * Quantity
	}
}

func IsFAReceiveUsed(DB *gorm.DB, ReceiveNumber string) (bool, error) {
	var Number string
	if err := DB.Table(db_var.TableName.FaRevaluation).Select("number").Joins("LEFT OUTER JOIN fa_list ON (fa_revaluation.fa_code = fa_list.code)").Where("fa_list.receive_number=?", ReceiveNumber).Limit(1).Scan(&Number).Error; err != nil {
		return false, err
	}
	var Code string
	if err := DB.Table(db_var.TableName.FaList).Select("code").Where("receive_number=?", ReceiveNumber).Where("condition_code<>?", global_var.FAItemCondition.Good).Limit(1).Scan(&Code).Error; err != nil {
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
			" AND inv_costing.date>='"+general.FormatDate1(StockDate)+"'"+
			") UNION ALL ("+
			"SELECT DISTINCT"+
			" inv_stock_transfer.`date` "+
			"FROM"+
			" inv_stock_transfer_detail"+
			" LEFT OUTER JOIN inv_stock_transfer ON (inv_stock_transfer_detail.st_number = inv_stock_transfer.number)"+
			" WHERE inv_stock_transfer_detail.item_code=? "+
			" AND inv_stock_transfer_detail.from_store_code=? "+
			" AND inv_stock_transfer.`date` >= '"+general.FormatDate1(StockDate)+"')) AS Stock "+
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
	DB.Table(db_var.TableName.InvStockTransferDetail).Select("receive_id").Where("receive_id IN (?)", ReceiveIDx).Limit(1).Scan(&ReceiveID)

	return ReceiveID > 0
}

func IsReceiveIDUsedInCosting(DB *gorm.DB, ReceiveIDx int64) bool {
	var ReceiveID uint64
	DB.Table(db_var.TableName.InvCostingDetail).Select("receive_id").Where("receive_id IN (?)", ReceiveIDx).Limit(1).Scan(&ReceiveID)

	return ReceiveID > 0
}

func IsFAPurchaseOrderUsedInReceive(DB *gorm.DB, PONumberx string) bool {
	var PONumber uint64
	DB.Table(db_var.TableName.FaReceive).Select("po_number").Where("po_number IN (?)", PONumberx).Limit(1).Scan(&PONumber)

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

func GetInventoryCOGS2JournalAccount(DB *gorm.DB, ItemCode, SubDepartmentCode string) string {
	var JournalAccountCode string
	DB.Raw(
		"SELECT"+
			" inv_cfg_init_item_category_other_cogs2.journal_account_code "+
			"FROM"+
			" inv_cfg_init_item"+
			" LEFT OUTER JOIN inv_cfg_init_item_category_other_cogs2 ON (inv_cfg_init_item.category_code = inv_cfg_init_item_category_other_cogs2.category_code)"+
			" WHERE inv_cfg_init_item.code=? "+
			" AND inv_cfg_init_item_category_other_cogs2.sub_department_code=? ", ItemCode, SubDepartmentCode).Limit(1).Scan(&JournalAccountCode)
	var JournalAccountCodeCOGS2 string
	if JournalAccountCode == "" {
		DB.Raw(
			"SELECT"+
				" inv_cfg_init_item_category.journal_account_code_cogs2 "+
				"FROM"+
				" inv_cfg_init_item"+
				" LEFT OUTER JOIN inv_cfg_init_item_category ON (inv_cfg_init_item.category_code = inv_cfg_init_item_category.code)"+
				" WHERE inv_cfg_init_item.code=? ", ItemCode).Limit(1).Scan(&JournalAccountCodeCOGS2)
		return JournalAccountCodeCOGS2
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
	if err := DB.Table("acc_import_journal_log").Select("id").Where("audit_date=?", general.FormatDate1(AuditDate)).Limit(1).Scan(&Id).Error; err != nil {
		return true, err
	}
	return Id > 0, nil
}

func IsMonthClosed(DB *gorm.DB, Month, Year uint64) (bool, error) {
	var Id uint64
	if err := DB.Table(db_var.TableName.AccCloseMonth).Select("id").
		Where("month=?", Month).Where("year=?", Year).Limit(1).Scan(&Id).Error; err != nil {
		return true, err
	}
	return Id > 0, nil
}

func IsPurchaseOrderClosed(DB *gorm.DB, PONumberX string) (bool, error) {
	var PONumber string
	if err := DB.Table(db_var.TableName.InvPurchaseOrderDetail).Select("po_number").Where("po_number=?", PONumberX).
		Where("quantity_not_received>0").
		Limit(1).
		Scan(&PONumber).Error; err != nil {
		return false, err
	}
	return PONumber != "", nil
}

func GetItemBasicUOMCode(DB *gorm.DB, ItemCode string) (string, error) {
	var DataOutput db_var.Inv_cfg_init_item
	if err := DB.Table(db_var.TableName.InvCfgInitItem).Select("uom_code").Where("code=?", ItemCode).Limit(1).Scan(&DataOutput).Error; err != nil {
		return "", err
	}

	return DataOutput.UomCode, nil
}

func GetItemMultiUOM(DB *gorm.DB, ItemCode, UomCode string) (db_var.Inv_cfg_init_item_uom, error) {
	var DataOutput db_var.Inv_cfg_init_item_uom
	if err := DB.Table(db_var.TableName.InvCfgInitItemUom).Where("item_code=?", ItemCode).Where("item_code=?", ItemCode).Limit(1).Scan(&DataOutput).Error; err != nil {
		return db_var.Inv_cfg_init_item_uom{}, err
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
	TableName := db_var.TableName.GuestScheduledRate
	if IsReservation {
		ConditionNumber = "reservation_number"
		TableName = db_var.TableName.ReservationScheduledRate
	}

	query := DB.Table(TableName).Select("id").
		Where(ConditionNumber, FolioReservationNumber).
		Where(DB.Where("from_date<=?", general.FormatDate1(FromDate)).Where("to_date", general.FormatDate1(FromDate)).
			Or("from_date<=?", general.FormatDate1(ToDate)).Where("to_date>=?", general.FormatDate1(ToDate)))

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
	if err := DB.Table(db_var.TableName.CfgInitRoomRate).Select("weekday_rate1").Where("code", RoomRateCode).Limit(1).Scan(&WeekdayRate1).Error; err != nil {
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
		if err := DB.Model(&db_var.Voucher{}).Where("number = ?", voucherNumber).Count(&count).Error; err != nil {
			return "", err
		}
		if count == 0 {
			if err := DB.Create(&db_var.Voucher{
				Number:            voucherNumber,
				Code:              code,
				TitleCode:         TitleCode,
				FullName:          FullName,
				Street:            Street,
				MemberTypeCode:    global_var.MemberType.Room,
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
				StatusCodeApprove: global_var.VoucherStatusApprove.Unapprove,
				StatusCodeSold:    "",
				IsRoomChargeOnly:  1,
				StatusCode:        global_var.VoucherStatus.Active,
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
	if err := DB.Table(db_var.TableName.Voucher).Select("expire_date").Where("number=?", VoucherNumber).Limit(1).Scan(&ExpireDate).Error; err != nil {
		return false, err
	}

	return ExpireDate.After(AuditDate), nil

}

func GetTotalMemberPointRedeem(DB *gorm.DB, MemberCode, MemberTypeCode string) (float64, error) {
	var TotalPoint float64
	if err := DB.Table(db_var.TableName.MemberPoint).Select("IFNULL(SUM(point),0) AS TotalPoint").
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
	if err := DB.Table(db_var.TableName.MemberPointRedeem).Select("IFNULL(SUM(point),0) AS TotalPoint").
		Where("member_code", MemberCode).
		Where("member_type_code", MemberTypeCode).
		Scan(&TotalPoint).Error; err != nil {
		return 0, err
	}
	return TotalPoint, nil
}

func SetAuditDate(c *gin.Context, DB *gorm.DB, Date time.Time, UserID string) error {
	if err := DB.Table(db_var.TableName.AuditLog).Create(map[string]interface{}{
		"audit_date": Date,
		"created_by": UserID}).Error; err != nil {
		return err
	}
	newAuditDate := GetAuditDate(c, DB, true)
	websocket.SendMessage(c.GetString("UnitCode"), global_var.WSMessageType.Client, nil, global_var.WSDataType.AuditDateChanged, newAuditDate, UserID)
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

func GetPaymentAccountName(DB *gorm.DB, PaymentAccountCode string) string {
	var PaymentAccount string
	DB.Table(db_var.TableName.CfgInitAccount).Select("name").Where("code", PaymentAccountCode).Limit(1).Scan(&PaymentAccount)
	return PaymentAccount
}

func GetJournalAccountName(DB *gorm.DB, JournalAccountCode string) (string, error) {
	var JournalAccountName string
	if err := DB.Table(db_var.TableName.CfgInitJournalAccount).Select("name").Where("code", JournalAccountCode).Limit(1).Scan(&JournalAccountName).Error; err != nil {
		return "", err
	}
	return JournalAccountName, nil
}

func IsJournalExportedYear(DB *gorm.DB, YearX string) (bool, error) {
	var CountExported float64
	if err := DB.Raw(
		"SELECT" +
			" (IFNULL(COUNT(TransactionX.audit_date), 0) - IFNULL(COUNT(ExportedJournal.audit_date), 0)) AS CountExported " +
			"FROM" +
			" audit_log" +
			" LEFT OUTER JOIN (" +
			"SELECT DISTINCT audit_date FROM ((" +
			"SELECT DISTINCT audit_date FROM guest_deposit WHERE YEAR(guest_deposit.audit_date)='" + YearX + "' AND void='0') UNION ALL (" +
			"SELECT DISTINCT audit_date FROM sub_folio WHERE YEAR(sub_folio.audit_date)='" + YearX + "' AND void='0')) AS Transaction)" +
			" AS TransactionX ON (audit_log.audit_date = TransactionX.audit_date)" +
			" LEFT OUTER JOIN (SELECT DISTINCT audit_date FROM acc_import_journal_log WHERE YEAR(acc_import_journal_log.audit_date)='" + YearX + "') AS ExportedJournal ON (audit_log.audit_date = ExportedJournal.audit_date)" +
			" WHERE YEAR(audit_log.audit_date)='" + YearX + "' " +
			"GROUP BY YEAR(audit_log.audit_date)").Scan(&CountExported).Error; err != nil {
		return false, err
	}
	return CountExported <= 0, nil
}

func GetJournalBankAccount(DB *gorm.DB, AccountCode string) (CodeName db_var.GeneralCodeNameStruct, err error) {
	if err := DB.Table(db_var.TableName.AccCfgInitBankAccount).
		Select("cfg_init_journal_account.code, cfg_init_journal_account.name").
		Joins("LEFT JOIN cfg_init_journal_account ON acc_cfg_init_bank_account.journal_account_code=cfg_init_journal_account.code").
		Where("acc_cfg_init_bank_account.code", AccountCode).
		Scan(&CodeName).Error; err != nil {
		return db_var.GeneralCodeNameStruct{}, err
	}

	return CodeName, nil
}

func GetBankAccountJournal(DB *gorm.DB, JournalAccountCode string) (BankAccount db_var.Acc_cfg_init_bank_account, err error) {
	if err := DB.Table(db_var.TableName.AccCfgInitBankAccount).
		Where("journal_account_code", JournalAccountCode).
		Scan(&BankAccount).Error; err != nil {
		return db_var.Acc_cfg_init_bank_account{}, err
	}

	return BankAccount, nil
}

func GetSubDepartmentName(DB *gorm.DB, SubDepartmentCode string) string {
	var SubDepartment string
	DB.Table(db_var.TableName.CfgInitSubDepartment).Select("name").Where("code", SubDepartmentCode).Limit(1).Scan(&SubDepartment)
	return SubDepartment
}

func IsFindBanReservationRemark(DB *gorm.DB, BookingNumber uint64, RemarkNumber int) (ID uint64, IsFound bool, err error) {
	if err := DB.Raw(
		"SELECT id FROM ban_reservation_remark"+
			" WHERE ban_reservation_remark.booking_number = ? "+
			" AND ban_reservation_remark.number = ? ", BookingNumber, RemarkNumber).Scan(&ID).Error; err != nil {
		return 0, false, err
	}
	return ID, ID > 0, nil
}

func GetCombineVenueNumber(DB *gorm.DB) (int64, error) {
	var Number int64
	if err := DB.Table(db_var.TableName.BanReservation).Select("venue_combine_number").Order("venue_combine_number DESC").Limit(1).Scan(&Number).Error; err != nil {
		return 0, err
	}

	if Number > 0 {
		return Number + 1, nil
	}
	return 1, nil
}

func UpdateBookingStatus(ctx context.Context, DB *gorm.DB, BookingNumber uint64, StatusCode, CancelBy, CancelReason string, AuditDate time.Time, UserID string) error {
	if (StatusCode == global_var.BookingStatus.Cancel) || (StatusCode == global_var.BookingStatus.NoShow) || (StatusCode == global_var.BookingStatus.Void) {
		if err := DB.Table(db_var.TableName.BanBooking).Where("number", BookingNumber).Updates(map[string]interface{}{
			"status_code":        StatusCode,
			"cancel_audit_date":  AuditDate,
			"cancel_date":        time.Now(),
			"cancel_by":          CancelBy,
			"cancel_reason":      CancelReason,
			"change_status_date": time.Now(),
			"change_status_by":   CancelBy,
			"updated_by":         UserID,
		}).Error; err != nil {
			return err
		}
		// CallProcedureSQL("update_ban_booking_status_cancel", strconv.FormatInt(BookingNumber)+", '"+StatusCode+"', '"+FormatDateTimeX(ProgramVariable.AuditDate)+"', '"+CancelBy+"', '"+CancelReason+"'", true)
	} else {
		if err := DB.Table(db_var.TableName.BanBooking).Where("number", BookingNumber).Updates(map[string]interface{}{
			"status_code": StatusCode,
			"updated_by":  UserID,
		}).Error; err != nil {
			return err
		}
		// CallProcedureSQL("update_ban_booking_status", strconv.FormatInt(BookingNumber)+", '"+StatusCode+"'", true)
	}
	return nil
}

func UpdateReservationStatusByBooking(ctx context.Context, DB *gorm.DB, BookingNumber uint64, StatusCode, CancelBy, CancelReason string, AuditDate time.Time, UserID string) error {
	if (StatusCode == global_var.BanquetReservationStatus.Canceled) || (StatusCode == global_var.BanquetReservationStatus.NoShow) || (StatusCode == global_var.BanquetReservationStatus.Void) {
		if err := DB.Table(db_var.TableName.BanReservation).Where("booking", BookingNumber).Updates(map[string]interface{}{
			"status_code":        StatusCode,
			"cancel_audit_date":  AuditDate,
			"cancel_date":        time.Now(),
			"cancel_by":          CancelBy,
			"cancel_reason":      CancelReason,
			"change_status_date": time.Now(),
			"change_status_by":   CancelBy,
			"updated_by":         UserID,
		}).Error; err != nil {
			return err
		}
	} else {
		if StatusCode == global_var.BanquetReservationStatus.Reservation {
			if err := DB.Table(db_var.TableName.BanReservation).Where("booking", BookingNumber).Updates(map[string]interface{}{
				"status_code":        StatusCode,
				"change_status_date": "0000-00-00 00:00:00",
				"updated_by":         UserID,
			}).Error; err != nil {
				return err
			}
		} else {
			if err := DB.Table(db_var.TableName.BanReservation).Where("booking", BookingNumber).Updates(map[string]interface{}{
				"status_code":        StatusCode,
				"change_status_date": time.Now(),
				"updated_by":         UserID,
			}).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func VoidReservationChargeByBookingNumber(DB *gorm.DB, BookingNumber, SubFolioID uint64, VoidBy, VoidReason, UserID string) error {
	if SubFolioID > 0 {
		if err := DB.Table(db_var.TableName.BanReservationCharge).Where("sub_folio_id", SubFolioID).Updates(map[string]interface{}{
			"void":        "1",
			"void_date":   time.Now(),
			"void_by":     VoidBy,
			"void_reason": VoidReason,
			"updated_by":  UserID,
		}).Error; err != nil {
			return err
		}
	} else {
		if err := DB.Table(db_var.TableName.BanReservationCharge).Where("booking_number", BookingNumber).Updates(map[string]interface{}{
			"void":        "1",
			"void_date":   time.Now(),
			"void_by":     VoidBy,
			"void_reason": VoidReason,
			"updated_by":  UserID,
		}).Error; err != nil {
			return err
		}
	}
	return nil
}

// begin
//   if SubFolioID <> '' then
//     CallProcedureSQL('update_ban_reservation_charge_void_by_sub_folio_id', '"' +SubFolioID+ '", "' +VoidBy+ '", "' +VoidReason+ '"', True)
//   else
//     CallProcedureSQL('update_ban_reservation_charge_void_by_booking_number', '"' +BookingNumber+ '", "' +VoidBy+ '", "' +VoidReason+ '"', True);
// end;

func UpdateReservationStatusByBookingX(ctx context.Context, DB *gorm.DB, BookingNumber uint64, ReservationStatusCode, UserID string) error {
	var Number []uint64
	if err := DB.Table(db_var.TableName.BanReservation).Select("number").Where("booking", BookingNumber).
		Where("status_code<>?", global_var.BanquetReservationStatus.Canceled).
		Where("status_code<>?", global_var.BanquetReservationStatus.NoShow).
		Where("status_code<>?", global_var.BanquetReservationStatus.Void).
		Scan(&Number).Error; err != nil {
		return err
	}

	for index := range Number {
		timeX := time.Now()
		if ReservationStatusCode == "R" {
			timeX = time.Time{}
		}
		if err := DB.Table(db_var.TableName.BanReservation).Where("number", Number[index]).Updates(map[string]interface{}{
			"status_code":        ReservationStatusCode,
			"change_status_date": timeX,
			"updated_by":         UserID,
		}).Error; err != nil {
			return err
		}
	}
	return nil
}
func UpdateBookingCheckNumber(ctx context.Context, DB *gorm.DB, BookingNumber uint64, CheckNumber, UserID string) error {
	if err := DB.Table(db_var.TableName.BanBooking).Where("number", BookingNumber).Updates(map[string]interface{}{
		"check_number": CheckNumber,
		"updated_by":   UserID,
	}).Error; err != nil {
		return err
	}
	return nil
}

func GetNumber2(DB *gorm.DB, Prefix, TableName, NumberField string, WithYear bool, AuditDate time.Time) (string, error) {
	// Prefix = Prefix // + ProgramVariable.ServerID
	if WithYear {
		Prefix = Prefix + general.FormatDatePrefix(AuditDate) + "-"
		// Prefix = Prefix + FormatDateTime("yy" + ProgramVariable.ServerID + "-", ProgramVariable.AuditDate)
	}

	var DataOutput uint64
	//From CPOS
	// DB.Raw("SELECT" +
	// 		" CAST(RIGHT(number, LENGTH(number) - " + strconv.FormatInt(int64(len(Prefix)), 10) + ") AS UNSIGNED) AS MaxNumber " +
	// 		"FROM" +
	// 		" pos_check" +
	// 		" WHERE LEFT(number," + strconv.FormatInt(int64(len(Prefix)), 10) + ")='" + Prefix + "' " +
	// 		"ORDER BY MaxNumber DESC " +
	// 		"LIMIT 1").Scan(&DataOutput)
	//From ban booking
	if err := DB.Raw("SELECT" +
		"  CAST(RIGHT(" + NumberField + ", LENGTH(" + NumberField + ") - " + strconv.FormatInt(int64(len(Prefix)), 10) + ") AS UNSIGNED) AS MaxNumber " +
		"FROM" +
		" " + TableName + " " +
		" WHERE LEFT(" + NumberField + "," + strconv.FormatInt(int64(len(Prefix)), 10) + ")='" + Prefix + "' " +
		"ORDER BY MaxNumber DESC " +
		"LIMIT 1").Scan(&DataOutput).Error; err != nil {
		return "0", err
	}
	return Prefix + general.Uint64ToStr(DataOutput+1), nil
}

func GetCheckNumber(DB *gorm.DB, Prefix string, AuditDate time.Time) (string, error) {
	Number, err := GetNumber2(DB, Prefix, "ban_booking", "check_number", true, AuditDate)
	if err != nil {
		return "", err
	}
	return Number, nil
}

func LoadOutletPOS(DB *gorm.DB, OutletCode string) (db_var.Pos_cfg_init_outlet, error) {
	var Outlet db_var.Pos_cfg_init_outlet
	if err := DB.Table(db_var.TableName.PosCfgInitOutlet).Where("code", OutletCode).Limit(1).Scan(&Outlet).Error; err != nil {
		return db_var.Pos_cfg_init_outlet{}, err
	}
	return Outlet, nil
}

// InsertCharge
func InsertBookingCharge(c *gin.Context, DB *gorm.DB, Dataset *global_var.TDataset, Outlet db_var.Pos_cfg_init_outlet, FolioNumber, ReservationNumber uint64, VenueCode, CurrencyCode string, ExchangeRate float64, UserID string) error {
	MainTableName := "ban_reservation_charge"
	TaxServiceManualOutlet := Outlet.TaxAndServiceCode
	DefaultCurrencyCode := GetDefaultCurrencyCode(DB)
	type MyQPostingChargeStruct struct {
		OutletCode    string
		ProductCode   string
		PackageCode   string
		CompanyCode   string
		AccountCode   string
		Description   string
		Quantity      float64
		PricePurchase float64
		PriceOriginal float64
		Price         float64
		Discount      float64
		Tax           float64
		Service       float64
		Remark        string
		TypeCode      string
		Name          string // Assuming this field is from pos_cfg_init_product
		ID            uint64
		VenueCode     string
		PackageRef    string
	}
	var MyQPostingCharge []MyQPostingChargeStruct
	if err := DB.Table(MainTableName).Select(
		" "+MainTableName+".outlet_code,"+
			" "+MainTableName+".product_code,"+
			" "+MainTableName+".package_code,"+
			" "+MainTableName+".company_code,"+
			" "+MainTableName+".account_code,"+
			" "+MainTableName+".description,"+
			" "+MainTableName+".quantity,"+
			" "+MainTableName+".price_purchase,"+
			" "+MainTableName+".price_original,"+ //8
			" "+MainTableName+".price,"+
			" "+MainTableName+".discount,"+
			" "+MainTableName+".tax,"+
			" "+MainTableName+".service,"+
			" "+MainTableName+".remark,"+
			" "+MainTableName+".type_code,"+
			" pos_cfg_init_product.name,"+ //15
			" "+MainTableName+".id,"+
			" IF(ban_reservation.venue_combine_code = '', ban_reservation.venue_code, ban_reservation.venue_combine_code) AS VenueCode,"+
			" "+MainTableName+".package_ref").
		Joins(" LEFT OUTER JOIN pos_cfg_init_product ON (ban_reservation_charge.product_code = pos_cfg_init_product.code)").
		Joins(" LEFT OUTER JOIN ban_reservation ON (ban_reservation_charge.reservation_number= ban_reservation.number)").
		Where(MainTableName+".void = '0'").
		Where(MainTableName+".booking_number = ?", ReservationNumber).
		Where(MainTableName+".input_of = ?", global_var.InputOf.Reservation).
		Where(MainTableName + ".is_posting = '1' ").
		Group(MainTableName + ".package_ref").
		Scan(&MyQPostingCharge).Error; err != nil {
		return err
	}

	for _, charge := range MyQPostingCharge {
		PackageCode := charge.PackageCode
		PackageRef := charge.PackageRef
		if PackageCode == "" {
			ProductCode := charge.ProductCode
			AccountCode := charge.AccountCode
			Remark := charge.Remark
			TransactionTypeCode := charge.TypeCode
			Quantity := charge.Quantity
			Amount := charge.PriceOriginal - charge.Discount
			IDLog := charge.ID
			VenueCode := charge.VenueCode
			_, SubFolioID, err := InsertSubFolio(c, DB, Dataset, false, AccountCode, TaxServiceManualOutlet, db_var.Sub_folio{
				FolioNumber:       FolioNumber,
				GroupCode:         global_var.SubFolioGroup.A,
				RoomNumber:        VenueCode,
				SubDepartmentCode: Dataset.GlobalSubDepartment.Banquet,
				AccountCode:       AccountCode,
				ProductCode:       ProductCode,
				PackageCode:       PackageCode,
				CurrencyCode:      CurrencyCode,
				Remark:            Remark,
				DocumentNumber:    "",
				VoucherNumber:     "",
				TypeCode:          TransactionTypeCode,
				PostingType:       global_var.SubFolioPostingType.None,
				Quantity:          Quantity,
				Amount:            Amount,
				ExchangeRate:      ExchangeRate,
				IsCorrection:      0,
				CreatedBy:         UserID,
			})
			if err != nil {
				return err
			}
			// InsertSubFolio(FolioNumber, SubFolioGroup.A, VenueCode, GlobalSubDepartment.Banquet, AccountCode, AccountCode, ProductCode, PackageCode, CurrencyCode, Remark, "", "", TransactionTypeCode, "", "", "", "", "", SubFolioPostingType.None, TaxServiceManualOutlet, 0, Quantity, Amount, ExchangeRate, False, False, SubFolioID)
			UpdateBanquetReservationChargeSubFolioID(DB, IDLog, SubFolioID, UserID)
		} else {
			//Package
			var MyQPostingChargeDetail []MyQPostingChargeStruct
			if err := DB.Table(MainTableName).Select(
				" "+MainTableName+".outlet_code,"+
					" "+MainTableName+".product_code,"+
					" "+MainTableName+".package_code,"+
					" "+MainTableName+".company_code,"+
					" "+MainTableName+".account_code,"+
					" "+MainTableName+".description,"+
					" "+MainTableName+".quantity,"+
					" "+MainTableName+".price_purchase,"+
					" "+MainTableName+".price_original,"+ //8
					" "+MainTableName+".price,"+
					" "+MainTableName+".discount,"+
					" "+MainTableName+".tax,"+
					" "+MainTableName+".service,"+
					" "+MainTableName+".remark,"+
					" "+MainTableName+".type_code,"+
					" pos_cfg_init_product.name,"+ //15
					" "+MainTableName+".id_log,"+
					" IF(ban_reservation.venue_combine_code = '', ban_reservation.venue_code, ban_reservation.venue_combine_code) AS VenueCode,"+
					" "+MainTableName+".package_ref ").
				Joins(" LEFT OUTER JOIN pos_cfg_init_product ON (ban_reservation_charge.product_code = pos_cfg_init_product.code)").
				Joins(" LEFT OUTER JOIN ban_reservation ON (ban_reservation_charge.reservation_number= ban_reservation.number)").
				Where(MainTableName+".void = '0'").
				Where(MainTableName+".input_of = ?", global_var.InputOf.Reservation).
				Where(MainTableName+".is_posting = '1'").
				Where(MainTableName+".booking_number = ?", ReservationNumber).
				Where(MainTableName+".package_ref = ?", PackageRef).
				Scan(&MyQPostingChargeDetail).Error; err != nil {
				return err
			}

			var BreakDown1, CorrectionBreakdown uint64
			if len(MyQPostingChargeDetail) > 0 {
				CorrectionBreakdown = GetSubFolioCorrectionBreakdown(c, DB)
				BreakDown1 = GetSubFolioBreakdown1(c, DB)
			}
			for _, detail := range MyQPostingChargeDetail {
				OutletCode := detail.OutletCode
				ProductCode := detail.ProductCode
				CompanyCode := detail.CompanyCode
				AccountCode := detail.AccountCode
				Remark := detail.Remark
				TransactionTypeCode := detail.TypeCode
				Quantity := detail.Quantity
				Price := detail.Price
				// Amount := detail.PriceOriginal - detail.Discount
				Tax := detail.Tax
				Service := detail.Service
				// IDLog := detail.ID
				VenueCode := detail.VenueCode

				SubDepartmentCode := OutletCode
				if SubDepartmentCode == Outlet.Code {
					SubDepartmentCode = Outlet.SubDepartmentCode

					BreakDown2 := GetSubFolioBreakdown2(c, DB, BreakDown1)
					if _, err := InsertSubFolioX(c, Dataset, 0, db_var.Sub_folio{
						FolioNumber:         FolioNumber,
						BelongsTo:           FolioNumber,
						GroupCode:           global_var.SubFolioGroup.A,
						RoomNumber:          VenueCode,
						SubDepartmentCode:   SubDepartmentCode,
						AccountCode:         AccountCode,
						ProductCode:         ProductCode,
						PackageCode:         PackageCode,
						CurrencyCode:        CurrencyCode,
						DefaultCurrencyCode: DefaultCurrencyCode,
						Remark:              Remark,
						DocumentNumber:      "",
						VoucherNumber:       "",
						TypeCode:            TransactionTypeCode,
						CorrectionBreakdown: CorrectionBreakdown,
						Breakdown1:          BreakDown1,
						Breakdown2:          BreakDown2,
						DirectBillCode:      CompanyCode,
						PostingType:         global_var.SubFolioPostingType.None,
						ExtraChargeId:       0,
						Quantity:            Quantity,
						Amount:              Price,
						AmountForeign:       Price,
						ExchangeRate:        ExchangeRate,
						CreatedBy:           UserID,
					}, DB); err != nil {
						return err
					}
					if Tax > 0 {
						if _, err := InsertSubFolioX(c, Dataset, 0, db_var.Sub_folio{
							FolioNumber:         FolioNumber,
							BelongsTo:           FolioNumber,
							GroupCode:           global_var.SubFolioGroup.A,
							RoomNumber:          VenueCode,
							SubDepartmentCode:   SubDepartmentCode,
							AccountCode:         Dataset.GlobalAccount.Tax,
							ProductCode:         ProductCode,
							PackageCode:         PackageCode,
							CurrencyCode:        CurrencyCode,
							DefaultCurrencyCode: DefaultCurrencyCode,
							Remark:              Remark,
							DocumentNumber:      "",
							VoucherNumber:       "",
							TypeCode:            TransactionTypeCode,
							CorrectionBreakdown: CorrectionBreakdown,
							Breakdown1:          BreakDown1,
							Breakdown2:          BreakDown2,
							DirectBillCode:      CompanyCode,
							PostingType:         global_var.SubFolioPostingType.None,
							ExtraChargeId:       0,
							Quantity:            Quantity,
							Amount:              Tax,
							AmountForeign:       Tax,
							ExchangeRate:        ExchangeRate,
							CreatedBy:           UserID,
						}, DB); err != nil {
							return err
						}
					}
					if Service > 0 {
						if _, err := InsertSubFolioX(c, Dataset, 0, db_var.Sub_folio{
							FolioNumber:         FolioNumber,
							BelongsTo:           FolioNumber,
							GroupCode:           global_var.SubFolioGroup.A,
							RoomNumber:          VenueCode,
							SubDepartmentCode:   SubDepartmentCode,
							AccountCode:         Dataset.GlobalAccount.Service,
							ProductCode:         ProductCode,
							PackageCode:         PackageCode,
							CurrencyCode:        CurrencyCode,
							DefaultCurrencyCode: DefaultCurrencyCode,
							Remark:              Remark,
							DocumentNumber:      "",
							VoucherNumber:       "",
							TypeCode:            TransactionTypeCode,
							CorrectionBreakdown: CorrectionBreakdown,
							Breakdown1:          BreakDown1,
							Breakdown2:          BreakDown2,
							DirectBillCode:      CompanyCode,
							PostingType:         global_var.SubFolioPostingType.None,
							ExtraChargeId:       0,
							Quantity:            Quantity,
							Amount:              Service,
							AmountForeign:       Service,
							ExchangeRate:        ExchangeRate,
							CreatedBy:           UserID,
						}, DB); err != nil {
							return err
						}
					}
				}
			}
		}
	}
	return nil
}

// UpdateReservationChargeSubFolioID
func UpdateBanquetReservationChargeSubFolioID(DB *gorm.DB, ID, SubFolioID uint64, UserID string) error {
	if err := DB.Table(db_var.TableName.BanReservationCharge).Where("id", ID).Updates(map[string]interface{}{
		"sub_folio_id": SubFolioID,
		"updated_by":   UserID,
	}).Error; err != nil {
		return err
	}
	return nil
}

func UpdateBookingReservationType(DB *gorm.DB, BookingNumber uint64, ResTypeCode, UserID string) error {
	if err := DB.Table(db_var.TableName.BanBooking).Where("number", BookingNumber).Updates(map[string]interface{}{
		"reservation_type": ResTypeCode,
		"updated_by":       UserID,
	}).Error; err != nil {
		return err
	}
	return nil
}

func IsReservationHaveEstimateCharge(DB *gorm.DB, ReservationNumber uint64) (bool, error) {
	var ID uint64
	if err := DB.Table(db_var.TableName.BanReservationCharge).Select("id").Where("reservation_number", ReservationNumber).Where("void='0'").Limit(1).Scan(&ID).Error; err != nil {
		return true, err
	}
	return ID > 0, nil
}

// UpdateReservationStatus
func UpdateBanquetReservationStatus(DB *gorm.DB, ReservationNumber uint64, StatusCode, CancelBy, CancelReason string, AuditDate time.Time, UserID string) error {
	//Get Combine Venue Number
	type MyQGetCombineNumberStruct struct {
		VenueCombineNumber, VenueCombineCode string
	}
	var MyQGetCombineNumber MyQGetCombineNumberStruct
	if err := DB.Table(db_var.TableName.BanReservation).Select("venue_combine_number, venue_combine_code").Where("ban_reservation.number", ReservationNumber).Limit(1).Scan(&MyQGetCombineNumber).Error; err != nil {
		return err
	}

	if MyQGetCombineNumber.VenueCombineCode != "" {
		//Get Reservation
		var MyQGetReservationNumber []uint64
		if err := DB.Table(db_var.TableName.BanReservation).Select("number").Where("ban_reservation.venue_combine_number", MyQGetCombineNumber.VenueCombineNumber).Scan(&MyQGetReservationNumber).Error; err != nil {
			return err
		}

		if len(MyQGetReservationNumber) > 0 {
			for _, number := range MyQGetReservationNumber {
				if (StatusCode == global_var.BanquetReservationStatus.Canceled) || (StatusCode == global_var.BanquetReservationStatus.NoShow) || (StatusCode == global_var.BanquetReservationStatus.Void) {
					if err := DB.Table(db_var.TableName.BanReservation).Where("number", number).Updates(map[string]interface{}{
						"status_code":        StatusCode,
						"cancel_audit_date":  AuditDate,
						"cancel_date":        time.Now(),
						"cancel_by":          CancelBy,
						"cancel_reason":      CancelReason,
						"change_status_date": time.Now(),
						"change_status_by":   CancelBy,
						"updated_by":         UserID,
					}).Error; err != nil {
						return err
					}
				} else {
					if StatusCode == global_var.BanquetReservationStatus.Reservation {
						if err := DB.Table(db_var.TableName.BanReservation).Where("number", number).Updates(map[string]interface{}{
							"status_code":        StatusCode,
							"change_status_date": "0000-00-00 00:00:00",
							"updated_by":         UserID,
						}).Error; err != nil {
							return err
						}
					} else {
						if err := DB.Table(db_var.TableName.BanReservation).Where("number", number).Updates(map[string]interface{}{
							"status_code":        StatusCode,
							"change_status_date": time.Now(),
							"updated_by":         UserID,
						}).Error; err != nil {
							return err
						}
					}
				}
			}
		}

	} else {
		if (StatusCode == global_var.BanquetReservationStatus.Canceled) || (StatusCode == global_var.BanquetReservationStatus.NoShow) || (StatusCode == global_var.BanquetReservationStatus.Void) {
			if err := DB.Table(db_var.TableName.BanReservation).Where("number", ReservationNumber).Updates(map[string]interface{}{
				"status_code":        StatusCode,
				"cancel_audit_date":  AuditDate,
				"cancel_date":        time.Now(),
				"cancel_by":          CancelBy,
				"cancel_reason":      CancelReason,
				"change_status_date": time.Now(),
				"change_status_by":   CancelBy,
				"updated_by":         UserID,
			}).Error; err != nil {
				return err
			}
		} else {
			if StatusCode == global_var.BanquetReservationStatus.Reservation {
				if err := DB.Table(db_var.TableName.BanReservation).Where("number", ReservationNumber).Updates(map[string]interface{}{
					"status_code":        StatusCode,
					"change_status_date": "0000-00-00 00:00:00",
					"updated_by":         UserID,
				}).Error; err != nil {
					return err
				}
			} else {
				if err := DB.Table(db_var.TableName.BanReservation).Where("number", ReservationNumber).Updates(map[string]interface{}{
					"status_code":        StatusCode,
					"change_status_date": time.Now(),
					"updated_by":         UserID,
				}).Error; err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func GetBanReservationChargePackageRef(DB *gorm.DB) uint64 {
	var PackageRef uint64
	DB.Table(db_var.TableName.BanReservationCharge).Select("package_ref").Order("package_ref DESC").Limit(1).Scan(&PackageRef)

	return PackageRef + 1
}

func GetFilePath(UnitCode, FolderName string, Name string) (fullPath, fileName string, err error) {
	basePath := global_var.PublicPath
	// datePath := int(time.Now().Month())
	// yearPath := int(time.Now().Year())
	// dateYearPath := strconv.Itoa(datePath) + strconv.Itoa(yearPath)
	// Path := basePath + "/" + FolderName + "/"
	Path := fmt.Sprintf("%s/%s/%s/", basePath, FolderName, UnitCode)
	if err := CreateDirectoryIfNotExist(Path); err != nil {
		return "", "", err
	}
	now := time.Now()
	Path += fmt.Sprintf("%s_%d", Name, now.Unix())
	return Path, strings.Replace(Path, basePath, "", 1), nil
}

func CreateDirectoryIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetIssuedCardCount(DB *gorm.DB, ReservationNumber, FolioNumber uint64, ActiveOnly bool) (int64, error) {
	QueryCondition := ""
	if ActiveOnly {
		QueryCondition = " AND is_active='1'"
	} else {
		QueryCondition = ""
	}

	var Count int64
	if ReservationNumber == 0 {
		if err := DB.Raw(
			"SELECT COUNT(id) AS IssuedCount FROM log_keylock"+
				" WHERE folio_number=? "+
				QueryCondition+"", FolioNumber).Scan(&Count).Error; err != nil {
			return 0, err
		}
	} else {
		if err := DB.Raw(
			"SELECT COUNT(id) AS IssuedCount FROM log_keylock"+
				" WHERE reservation_number=? "+
				QueryCondition+"", ReservationNumber).Scan(&Count).Error; err != nil {
			return 0, err
		}
	}

	return Count, nil
}

func GetOutletDetail(DB *gorm.DB, OutletCode string) db_var.Pos_cfg_init_outlet {
	DataOutput := db_var.Pos_cfg_init_outlet{}
	err := DB.Table(db_var.TableName.PosCfgInitOutlet).Where("code=?", OutletCode).Take(&DataOutput)
	if err != nil {
		fmt.Println(err)
	}

	return DataOutput
}

func LogRegistrationFormFieldChange(c *gin.Context, DB *gorm.DB, Number uint64, NewData db_var.RegistrationFormStruct, IsReservation bool, Reason string) error {
	type DataOutputStruct struct {
		// ReservationData   db_var.Reservation    `gorm:"embedded"`
		GuestGeneralData  db_var.Guest_general  `gorm:"embedded"`
		GuestDetailData   db_var.Guest_detail   `gorm:"embedded"`
		GuestProfileData1 db_var.Contact_person `gorm:"embedded"`
		GuestProfileData2 db_var.Contact_person `gorm:"embedded"`
		GuestProfileData3 db_var.Contact_person `gorm:"embedded"`
		GuestProfileData4 db_var.Contact_person `gorm:"embedded"`
	}
	var OldData DataOutputStruct
	Remark := ""
	if IsReservation {
		Remark = "R"
		if err := DB.Table(db_var.TableName.Reservation).Select(
			// "reservation.*",
			"guest_general.*", "guest_detail.*", "contact_person1.*", "contact_person2.*", "contact_person3.*", "contact_person4.*",
			"IFNULL(SUM(IF(guest_deposit.type_code='"+global_var.TransactionType.Debit+"', guest_deposit.amount, -guest_deposit.amount)), 0) AS Balance", "folio.number AS FolioNumber").
			Joins("LEFT JOIN folio ON reservation.number = folio.reservation_number").
			Joins("LEFT JOIN contact_person AS contact_person1 ON reservation.contact_person_id1 = contact_person1.id").
			Joins("LEFT JOIN contact_person AS contact_person2 ON reservation.contact_person_id2 = contact_person2.id").
			Joins("LEFT JOIN contact_person AS contact_person3 ON reservation.contact_person_id3 = contact_person3.id").
			Joins("LEFT JOIN contact_person AS contact_person4 ON reservation.contact_person_id4 = contact_person4.id").
			Joins("LEFT JOIN guest_detail  ON reservation.guest_detail_id = guest_detail.id").
			Joins("LEFT JOIN guest_general  ON reservation.guest_general_id = guest_general.id").
			Joins("LEFT JOIN guest_deposit  ON reservation.number = guest_deposit.reservation_number AND guest_deposit.void='0' AND guest_deposit.system_code='"+global_var.ConstProgramVariable.DefaultSystemCode+"'").
			Joins("LEFT JOIN cfg_init_room ON guest_detail.room_number = cfg_init_room.number").
			Where("reservation.number = ?", Number).
			Group("reservation.number").
			Take(&OldData).Error; err != nil {
			return err
		}
	} else {

		// type DataOutputStruct struct {
		// 	FolioData         db_var.Folio          `gorm:"embedded"`
		// 	GuestProfileData1 db_var.Contact_person `gorm:"embedded"`
		// 	GuestProfileData2 db_var.Contact_person `gorm:"embedded"`
		// 	GuestProfileData3 db_var.Contact_person `gorm:"embedded"`
		// 	GuestProfileData4 db_var.Contact_person `gorm:"embedded"`
		// 	GuestGeneralData  db_var.Guest_general  `gorm:"embedded"`
		// 	GuestDetailData   db_var.Guest_detail   `gorm:"embedded"`
		// }
		if err := DB.Table(db_var.TableName.Folio).Select(
			// "folio.*",
			"guest_general.*", "guest_detail.*", "contact_person1.*", "contact_person2.*", "contact_person3.*", "contact_person4.*",
			"IFNULL(SUM(IF(sub_folio.type_code='"+global_var.TransactionType.Debit+"', (sub_folio.quantity * sub_folio.amount), -(sub_folio.quantity * sub_folio.amount))), 0) AS Balance", "cfg_init_room.bed_type_code").
			Joins("LEFT JOIN contact_person AS contact_person1 ON folio.contact_person_id1 = contact_person1.id").
			Joins("LEFT JOIN contact_person AS contact_person2 ON folio.contact_person_id2 = contact_person2.id").
			Joins("LEFT JOIN contact_person AS contact_person3 ON folio.contact_person_id3 = contact_person3.id").
			Joins("LEFT JOIN contact_person AS contact_person4 ON folio.contact_person_id4 = contact_person4.id").
			Joins("LEFT JOIN guest_detail  ON folio.guest_detail_id = guest_detail.id").
			Joins("LEFT JOIN guest_general  ON folio.guest_general_id = guest_general.id").
			Joins("LEFT JOIN sub_folio  ON folio.number = sub_folio.folio_number AND sub_folio.void='0'").
			Joins("LEFT JOIN cfg_init_room ON guest_detail.room_number = cfg_init_room.number").
			Where("folio.number = ?", Number).
			Group("folio.number").
			Take(&OldData).Error; err != nil {
			return err
		}
	}
	//Stay Information
	if NewData.GuestDetailData.Arrival != OldData.GuestDetailData.Arrival {
		InsertLogUserX(c, DB, global_var.SystemCode.Hotel, global_var.LogUserAction.URSIArrival, Number, OldData.GuestDetailData.Arrival, NewData.GuestDetailData.Arrival, "")
	}
	// if cxSpinEditNights.Value != NightsB4 {
	// InsertLogUserX(c,DB,global_var.SystemCode.Hotel, global_var.LogUserAction.URSINights, Number,  strconv.FormatInt(NightsB4), strconv.FormatInt(cxSpinEditNights.Value), Remark)
	if NewData.GuestDetailData.Departure != OldData.GuestDetailData.Departure {
		InsertLogUserX(c, DB, global_var.SystemCode.Hotel, global_var.LogUserAction.URSIDeparture, Number, OldData.GuestDetailData.Departure, NewData.GuestDetailData.Departure, "")
	}
	if NewData.GuestDetailData.Adult != OldData.GuestDetailData.Adult {
		InsertLogUserX(c, DB, global_var.SystemCode.Hotel, global_var.LogUserAction.URSIAdult, Number, OldData.GuestDetailData.Adult, NewData.GuestDetailData.Adult, Remark)
	}
	if *NewData.GuestDetailData.Child != *OldData.GuestDetailData.Child {
		InsertLogUserX(c, DB, global_var.SystemCode.Hotel, global_var.LogUserAction.URSIChild, Number, OldData.GuestDetailData.Child, NewData.GuestDetailData.Child, Remark)
	}
	if NewData.GuestDetailData.RoomTypeCode != OldData.GuestDetailData.RoomTypeCode {
		InsertLogUserX(c, DB, global_var.SystemCode.Hotel, global_var.LogUserAction.URSIRoomType, Number, OldData.GuestDetailData.RoomTypeCode, NewData.GuestDetailData.RoomRateCode, Remark)
	}
	if NewData.GuestDetailData.BedTypeCode != OldData.GuestDetailData.BedTypeCode {
		InsertLogUserX(c, DB, global_var.SystemCode.Hotel, global_var.LogUserAction.URSIBedType, Number, OldData.GuestDetailData.BedTypeCode, NewData.GuestDetailData.BedTypeCode, Remark)
	}
	if *NewData.GuestDetailData.RoomNumber != *OldData.GuestDetailData.RoomNumber {
		InsertLogUserX(c, DB, global_var.SystemCode.Hotel, global_var.LogUserAction.URSIRoom, Number, OldData.GuestDetailData.RoomNumber, NewData.GuestDetailData.RoomNumber, Remark)
	}
	if NewData.GuestDetailData.CurrencyCode != OldData.GuestDetailData.CurrencyCode {
		InsertLogUserX(c, DB, global_var.SystemCode.Hotel, global_var.LogUserAction.URSICurrency, Number, OldData.GuestDetailData.CurrencyCode, NewData.GuestDetailData.CurrencyCode, Remark)
	}
	if NewData.GuestDetailData.ExchangeRate != OldData.GuestDetailData.ExchangeRate {
		InsertLogUserX(c, DB, global_var.SystemCode.Hotel, global_var.LogUserAction.URSIExchangeRate, Number, OldData.GuestDetailData.ExchangeRate, NewData.GuestDetailData.ExchangeRate, Remark)
	}
	if *NewData.GuestDetailData.BusinessSourceCode != *OldData.GuestDetailData.BusinessSourceCode {
		InsertLogUserX(c, DB, global_var.SystemCode.Hotel, global_var.LogUserAction.URSIBusinessSource, Number, OldData.GuestDetailData.BusinessSourceCode, NewData.GuestDetailData.BusinessSourceCode, Reason)
	}
	if *NewData.GuestDetailData.CommissionTypeCode != *OldData.GuestDetailData.CommissionTypeCode {
		// fmt.Println("CommissionTypeCode1", NewData.GuestDetailData.CommissionTypeCode)
		// fmt.Println("CommissionTypeCode2", OldData.GuestDetailData.CommissionTypeCode)
		InsertLogUserX(c, DB, global_var.SystemCode.Hotel, global_var.LogUserAction.URSICommissionType, Number, OldData.GuestDetailData.CommissionTypeCode, NewData.GuestDetailData.CommissionTypeCode, Reason)
	}
	if *NewData.GuestDetailData.CommissionValue != *OldData.GuestDetailData.CommissionValue {
		InsertLogUserX(c, DB, global_var.SystemCode.Hotel, global_var.LogUserAction.URSICommissionValue, Number, OldData.GuestDetailData.CommissionValue, NewData.GuestDetailData.CommissionValue, Reason)
	}
	if NewData.GuestDetailData.RoomRateCode != OldData.GuestDetailData.RoomRateCode {
		InsertLogUserX(c, DB, global_var.SystemCode.Hotel, global_var.LogUserAction.URSIRoomRate, Number, OldData.GuestDetailData.RoomRateCode, NewData.GuestDetailData.RoomRateCode, Reason)
	}
	if *NewData.GuestDetailData.WeekdayRate != *OldData.GuestDetailData.WeekdayRate {
		InsertLogUserX(c, DB, global_var.SystemCode.Hotel, global_var.LogUserAction.URSIWeekdayRate, Number, OldData.GuestDetailData.WeekdayRate, NewData.GuestDetailData.WeekdayRate, Reason)
	}
	if *NewData.GuestDetailData.WeekendRate != *OldData.GuestDetailData.WeekendRate {
		InsertLogUserX(c, DB, global_var.SystemCode.Hotel, global_var.LogUserAction.URSIWeekendRate, Number, OldData.GuestDetailData.WeekendRate, NewData.GuestDetailData.WeekendRate, Reason)
	}
	if (*NewData.GuestDetailData.DiscountPercent != *OldData.GuestDetailData.DiscountPercent) || (*NewData.GuestDetailData.Discount != *OldData.GuestDetailData.Discount) {
		DiscountRemarkB4 := ""
		DiscountRemarkAfter := ""
		if *OldData.GuestDetailData.DiscountPercent > 0 {
			DiscountRemarkB4 = fmt.Sprintf("%.2f%s", *OldData.GuestDetailData.Discount, " (Percent)")
		} else {
			DiscountRemarkB4 = fmt.Sprintf("%f", *OldData.GuestDetailData.Discount)
		}
		if *NewData.GuestDetailData.DiscountPercent > 0 {
			DiscountRemarkAfter = fmt.Sprintf("%.2f%s", *NewData.GuestDetailData.Discount, " (Percent)")
		} else {
			DiscountRemarkAfter = fmt.Sprintf("%.2f", *NewData.GuestDetailData.Discount)
		}
		InsertLogUserX(c, DB, global_var.SystemCode.Hotel, global_var.LogUserAction.URSIDiscount, Number, DiscountRemarkB4, DiscountRemarkAfter, Reason)

	}
	if NewData.GuestDetailData.PaymentTypeCode != OldData.GuestDetailData.PaymentTypeCode {
		InsertLogUserX(c, DB, global_var.SystemCode.Hotel, global_var.LogUserAction.URSIPaymentType, Number, OldData.GuestDetailData.PaymentTypeCode, NewData.GuestDetailData.PaymentTypeCode, Remark)
	}
	if *NewData.GuestDetailData.MarketCode != *OldData.GuestDetailData.MarketCode {
		InsertLogUserX(c, DB, global_var.SystemCode.Hotel, global_var.LogUserAction.URSIMarket, Number, OldData.GuestDetailData.MarketCode, NewData.GuestDetailData.MarketCode, Remark)
	}
	if *NewData.GuestDetailData.BillInstruction != *OldData.GuestDetailData.BillInstruction {
		InsertLogUserX(c, DB, global_var.SystemCode.Hotel, global_var.LogUserAction.URSIBillInstruction, Number, OldData.GuestDetailData.BillInstruction, NewData.GuestDetailData.BillInstruction, Remark)

		//Personal Information
	}
	// if NewData.GuestProfileData1.Mem != MemberB4 {
	// 	InsertLogUserX(c, DB, global_var.SystemCode.Hotel, global_var.LogUserAction.URPIMember, Number,  MemberB4, VarToStr(cxLookupComboBoxMember.EditValue), Remark)
	// }
	if NewData.GuestProfileData1.TitleCode != *OldData.GuestProfileData1.TitleCode {
		InsertLogUserX(c, DB, global_var.SystemCode.Hotel, global_var.LogUserAction.URPITitle, Number, OldData.GuestProfileData1.TitleCode, NewData.GuestProfileData1.TitleCode, Remark)
	}
	if NewData.GuestProfileData1.FullName != *OldData.GuestProfileData1.FullName {
		InsertLogUserX(c, DB, global_var.SystemCode.Hotel, global_var.LogUserAction.URPIFullName, Number, OldData.GuestProfileData1.FullName, NewData.GuestProfileData1.FullName, Remark)
	}
	// if cxTextEditReservationBy != ReservationByB4 {
	// 	InsertLogUserX(c, DB, global_var.SystemCode.Hotel, global_var.LogUserAction.URPIReservationBy, Number,  ReservationByB4, cxTextEditReservationBy, Remark)
	// }
	if NewData.GuestProfileData1.Street != *OldData.GuestProfileData1.Street {
		InsertLogUserX(c, DB, global_var.SystemCode.Hotel, global_var.LogUserAction.URPIStreet, Number, OldData.GuestProfileData1.Street, NewData.GuestProfileData1.Street, Remark)
	}
	if NewData.GuestProfileData1.CityCode != *OldData.GuestProfileData1.CityCode {
		InsertLogUserX(c, DB, global_var.SystemCode.Hotel, global_var.LogUserAction.URPICity, Number, OldData.GuestProfileData1.CityCode, NewData.GuestProfileData1.CityCode, Remark)
	}
	if NewData.GuestProfileData1.City != *OldData.GuestProfileData1.City {
		InsertLogUserX(c, DB, global_var.SystemCode.Hotel, global_var.LogUserAction.URPICity, Number, OldData.GuestProfileData1.City, NewData.GuestProfileData1.City, Remark)
	}
	if NewData.GuestProfileData1.CountryCode != *OldData.GuestProfileData1.CountryCode {
		InsertLogUserX(c, DB, global_var.SystemCode.Hotel, global_var.LogUserAction.URPICountry, Number, OldData.GuestProfileData1.CountryCode, NewData.GuestProfileData1.CountryCode, Remark)
	}
	if NewData.GuestProfileData1.StateCode != *OldData.GuestProfileData1.StateCode {
		InsertLogUserX(c, DB, global_var.SystemCode.Hotel, global_var.LogUserAction.URPIState, Number, OldData.GuestProfileData1.StateCode, NewData.GuestProfileData1.StateCode, Remark)
	}
	if NewData.GuestProfileData1.PostalCode != *OldData.GuestProfileData1.PostalCode {
		InsertLogUserX(c, DB, global_var.SystemCode.Hotel, global_var.LogUserAction.URPIPostCode, Number, OldData.GuestProfileData1.PostalCode, NewData.GuestProfileData1.PostalCode, Remark)
	}
	if NewData.GuestProfileData1.Phone1 != *OldData.GuestProfileData1.Phone1 {
		InsertLogUserX(c, DB, global_var.SystemCode.Hotel, global_var.LogUserAction.URPIPhone1, Number, OldData.GuestProfileData1.Phone1, NewData.GuestProfileData1.Phone1, Remark)
	}
	if NewData.GuestProfileData1.Phone2 != *OldData.GuestProfileData1.Phone2 {
		InsertLogUserX(c, DB, global_var.SystemCode.Hotel, global_var.LogUserAction.URPIPhone2, Number, OldData.GuestProfileData1.Phone2, NewData.GuestProfileData1.Phone2, Remark)
	}
	if NewData.GuestProfileData1.Fax != *OldData.GuestProfileData1.Fax {
		InsertLogUserX(c, DB, global_var.SystemCode.Hotel, global_var.LogUserAction.URPIFax, Number, OldData.GuestProfileData1.Fax, NewData.GuestProfileData1.Fax, Remark)
	}
	if NewData.GuestProfileData1.Email != *OldData.GuestProfileData1.Email {
		InsertLogUserX(c, DB, global_var.SystemCode.Hotel, global_var.LogUserAction.URPIEmail, Number, OldData.GuestProfileData1.Email, NewData.GuestProfileData1.Email, Remark)
	}
	if NewData.GuestProfileData1.Website != *OldData.GuestProfileData1.Website {
		InsertLogUserX(c, DB, global_var.SystemCode.Hotel, global_var.LogUserAction.URPIWebsite, Number, OldData.GuestProfileData1.Website, NewData.GuestProfileData1.Website, Remark)
	}
	if NewData.GuestProfileData1.CompanyCode != *OldData.GuestProfileData1.CompanyCode {
		InsertLogUserX(c, DB, global_var.SystemCode.Hotel, global_var.LogUserAction.URPICompany, Number, OldData.GuestProfileData1.CompanyCode, NewData.GuestProfileData1.CompanyCode, Remark)
	}
	if NewData.GuestProfileData1.GuestTypeCode != *OldData.GuestProfileData1.GuestTypeCode {
		InsertLogUserX(c, DB, global_var.SystemCode.Hotel, global_var.LogUserAction.URPIGuestType, Number, OldData.GuestProfileData1.GuestTypeCode, NewData.GuestProfileData1.GuestTypeCode, Remark)
	}
	if NewData.GuestProfileData1.IdCardCode != *OldData.GuestProfileData1.IdCardCode {
		InsertLogUserX(c, DB, global_var.SystemCode.Hotel, global_var.LogUserAction.URPIIDCardType, Number, OldData.GuestProfileData1.IdCardCode, NewData.GuestProfileData1.IdCardCode, Remark)
	}
	if NewData.GuestProfileData1.IdCardNumber != *OldData.GuestProfileData1.IdCardNumber {
		InsertLogUserX(c, DB, global_var.SystemCode.Hotel, global_var.LogUserAction.URPIIDCardNumber, Number, OldData.GuestProfileData1.IdCardNumber, NewData.GuestProfileData1.IdCardNumber, Remark)
	}
	if NewData.GuestProfileData1.BirthPlace != *OldData.GuestProfileData1.BirthPlace {
		InsertLogUserX(c, DB, global_var.SystemCode.Hotel, global_var.LogUserAction.URPIBirthdayPlace, Number, OldData.GuestProfileData1.BirthPlace, NewData.GuestProfileData1.BirthPlace, Remark)
	}
	if NewData.GuestProfileData1.BirthDate != *OldData.GuestProfileData1.BirthDate {
		InsertLogUserX(c, DB, global_var.SystemCode.Hotel, global_var.LogUserAction.URPIBirthdate, Number, OldData.GuestProfileData1.BirthDate, NewData.GuestProfileData1.BirthDate, Remark)

		//General Information
	}
	if *NewData.GuestGeneralData.PurposeOfCode != *OldData.GuestGeneralData.PurposeOfCode {
		InsertLogUserX(c, DB, global_var.SystemCode.Hotel, global_var.LogUserAction.URGIPurposeOf, Number, OldData.GuestGeneralData.PurposeOfCode, NewData.GuestGeneralData.PurposeOfCode, Remark)
	}
	// if NewData.GuestGeneralData. != GroupB4 {
	// 	InsertLogUserX(c, DB, global_var.SystemCode.Hotel, global_var.LogUserAction.URGIGroup, Number,  GroupB4, VarToStr(cxLookupComboBoxGroup.EditValue), Remark)
	// }
	if *NewData.GuestGeneralData.SalesCode != *OldData.GuestGeneralData.SalesCode {
		InsertLogUserX(c, DB, global_var.SystemCode.Hotel, global_var.LogUserAction.URGIMarketing, Number, OldData.GuestGeneralData.SalesCode, NewData.GuestGeneralData.SalesCode, Remark)
	}
	if *NewData.GuestGeneralData.DocumentNumber != *OldData.GuestGeneralData.DocumentNumber {
		InsertLogUserX(c, DB, global_var.SystemCode.Hotel, global_var.LogUserAction.URGIDocumentNumber, Number, OldData.GuestGeneralData.DocumentNumber, NewData.GuestGeneralData.DocumentNumber, Remark)
	}
	if *NewData.GuestGeneralData.VoucherNumberTa != *OldData.GuestGeneralData.VoucherNumberTa {
		InsertLogUserX(c, DB, global_var.SystemCode.Hotel, global_var.LogUserAction.URGITAVoucherNumber, Number, OldData.GuestGeneralData.VoucherNumberTa, NewData.GuestGeneralData.VoucherNumberTa, Remark)
	}
	if *NewData.GuestGeneralData.FlightNumber != *OldData.GuestGeneralData.FlightNumber {
		InsertLogUserX(c, DB, global_var.SystemCode.Hotel, global_var.LogUserAction.URGIFlightNumber, Number, OldData.GuestGeneralData.FlightNumber, NewData.GuestGeneralData.FlightNumber, Remark)
	}
	if *NewData.GuestGeneralData.FlightArrival != *OldData.GuestGeneralData.FlightArrival {
		InsertLogUserX(c, DB, global_var.SystemCode.Hotel, global_var.LogUserAction.URGIFlightArrival, Number, OldData.GuestGeneralData.FlightArrival, NewData.GuestGeneralData.FlightArrival, Remark)
	}
	if *NewData.GuestGeneralData.FlightDeparture != *OldData.GuestGeneralData.FlightDeparture {
		InsertLogUserX(c, DB, global_var.SystemCode.Hotel, global_var.LogUserAction.URGIFlightDeparture, Number, OldData.GuestGeneralData.FlightDeparture, NewData.GuestGeneralData.FlightDeparture, Remark)
	}
	if *NewData.GuestGeneralData.Notes != *OldData.GuestGeneralData.Notes {
		InsertLogUserX(c, DB, global_var.SystemCode.Hotel, global_var.LogUserAction.URGINotes, Number, OldData.GuestGeneralData.Notes, NewData.GuestGeneralData.Notes, Remark)
	}
	if *NewData.GuestGeneralData.HkNote != *OldData.GuestGeneralData.HkNote {
		InsertLogUserX(c, DB, global_var.SystemCode.Hotel, global_var.LogUserAction.URGIHKNotes, Number, OldData.GuestGeneralData.HkNote, NewData.GuestGeneralData.HkNote, Remark)
	}
	return nil
}

func ProcessProductCosting(c *gin.Context, pConfig *config.CompanyDataConfiguration, DB *gorm.DB, CheckNumber string, CheckTransactionID uint64, OutletCode, ProductCode, ProductName, SubDepartmentCompliment, RefNumberCosting, CostingNumber string, ProductQuantity float64, PostingDate time.Time, IsCompliment bool, UserID string) (RefNumberCostingX string, CostingNumberX string, Error error) {

	type InventoryCostRecipe struct {
		StoreCode    string  `json:"store_code"`
		StoreID      int64   `json:"store_id"`
		ItemCode     string  `json:"item_code"`
		ItemID       int64   `json:"item_id"`
		Quantity     float64 `json:"quantity"`
		UOMCode      string  `json:"uom_code"`
		Name         string  `json:"name"`
		BasicUOMCode string  `json:"basic_uom_code"`
		GroupCode    string  `json:"group_code"`
	}
	type InventoryReceiveDetail struct {
		ItemCode         string    `json:"item_code"`
		BasicQuantity    float64   `json:"basic_quantity"`
		Quantity         float64   `json:"quantity"` // This is the calculated field
		AllStoreQuantity float64   `json:"all_store_quantity"`
		TotalPrice       float64   `json:"total_price"`
		ID               int64     `json:"id"`
		ReceiveDate      time.Time `json:"receive_date"`
	}
	UnitCode := c.GetString("UnitCode")
	OutletDetail := GetOutletDetail(DB, OutletCode)
	QueryCondition := " WHERE IFNULL(inv_cost_recipe.outlet_code,'')='' AND inv_cost_recipe.product_code='" + ProductCode + "'"
	if master_data.CheckCodeField(DB, "inv_cost_recipe", "id", "outlet_code", OutletCode) {
		QueryCondition = " WHERE inv_cost_recipe.outlet_code='" + OutletCode + "' AND inv_cost_recipe.product_code='" + ProductCode + "'"
	}

	var MyQCostRecipe []InventoryCostRecipe
	DB.Debug().Raw(
		"SELECT" +
			" inv_cost_recipe.store_code," +
			" inv_cfg_init_store.id AS StoreID," +
			" inv_cost_recipe.item_code," +
			" inv_cfg_init_item.id AS ItemID," +
			" inv_cost_recipe.quantity," +
			" inv_cost_recipe.uom_code," +
			" inv_cfg_init_item.name," +
			" inv_cfg_init_item.uom_code AS basic_uom_code," +
			" inv_cfg_init_item_category.group_code " +
			"FROM" +
			" inv_cost_recipe" +
			" LEFT OUTER JOIN inv_cfg_init_store ON (inv_cost_recipe.store_code = inv_cfg_init_store.code)" +
			" LEFT OUTER JOIN inv_cfg_init_item ON (inv_cost_recipe.item_code = inv_cfg_init_item.code)" +
			" LEFT OUTER JOIN inv_cfg_init_item_category ON (inv_cfg_init_item.category_code = inv_cfg_init_item_category.code)" +
			" " + QueryCondition + "").Scan(&MyQCostRecipe)

	if len(MyQCostRecipe) > 0 {
		// var CostingDetail []db_var.Inv_costing_detail
		var JournalDetailDebit []db_var.Acc_journal_detail
		var JournalDetailCredit []db_var.Acc_journal_detail
		var ProductCosting []db_var.Pos_product_costing

		ProductItemGroup := GetProductItemGroup(DB, ProductCode)
		if err := DB.Transaction(func(tx *gorm.DB) error {
			for _, data := range MyQCostRecipe {
				StoreCode := data.StoreCode
				StoreID := data.StoreID
				if data.StoreCode == "" {
					StoreCode = OutletDetail.StoreCode
					StoreID = int64(GetStoreID(DB, StoreCode))
				}

				ItemCode := data.ItemCode
				ItemID := data.ItemID
				ItemName := data.Name
				UOMCode := data.UOMCode
				BasicUOMCode := data.BasicUOMCode
				CostRecipeQuantity := data.Quantity * ProductQuantity

				InventoryJournalAccount := GetInventoryJournalAccount(DB, ItemCode)
				InventoryCOGSJournalAccount := ""
				if IsCompliment {
					InventoryCOGSJournalAccount = GetInventoryExpenseJournalAccount(DB, ItemCode, SubDepartmentCompliment)
				} else {
					if ProductItemGroup == data.GroupCode {
						InventoryCOGSJournalAccount = GetInventoryCOGSJournalAccount(DB, ItemCode, OutletDetail.SubDepartmentCode)
					} else {
						InventoryCOGSJournalAccount = GetInventoryCOGS2JournalAccount(DB, ItemCode, OutletDetail.SubDepartmentCode)
					}
				}

				var MyQUOM db_var.Inv_cfg_init_item_uom
				tx.Raw(
					"SELECT * FROM inv_cfg_init_item_uom WHERE item_code=? AND uom_code=? ", ItemCode, UOMCode).Scan(&MyQUOM)

				BasicQuantity := CostRecipeQuantity
				if MyQUOM.UomCode != "" {
					BasicQuantity = MyQUOM.Quantity * CostRecipeQuantity
				}

				if (pConfig.Dataset.ProgramConfiguration.CostingMethod == global_var.InventoryCostingMethod.FIFO) || (pConfig.Dataset.ProgramConfiguration.CostingMethod == global_var.InventoryCostingMethod.LIFO) {
					StockOrderBy := ""
					StockDateStr := general.FormatDate1(PostingDate)
					if pConfig.Dataset.ProgramConfiguration.CostingMethod == global_var.InventoryCostingMethod.LIFO {
						StockOrderBy = "DESC"
					}

					var MyQStockID []InventoryReceiveDetail

					tx.Raw("SELECT * FROM (" +
						"SELECT" +
						" inv_receiving_detail.item_code," +
						" inv_receiving_detail.basic_quantity," +
						" (SUM(IF(inv_receiving_detail.store_code='" + StoreCode + "', inv_receiving_detail.basic_quantity, 0)) + IFNULL(StockTransfer.Quantity, 0) - IFNULL(Costing.Quantity, 0)) AS Quantity," +
						" inv_receiving_detail.quantity AS AllStoreQuantity," +
						" inv_receiving_detail.total_price," +
						" inv_receiving_detail.id," +
						" inv_receiving.date " +
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
						" AND inv_stock_transfer_detail.item_code='" + ItemCode + "' " +
						" AND (inv_stock_transfer_detail.from_store_code='" + StoreCode + "' OR (inv_stock_transfer_detail.to_store_code='" + StoreCode + "' AND inv_stock_transfer.date<='" + StockDateStr + "')) " +
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
						" AND inv_costing_detail.item_code='" + ItemCode + "' " +
						" AND inv_costing_detail.store_code='" + StoreCode + "' " +
						"GROUP BY inv_costing_detail.receive_id) AS Costing ON (inv_receiving_detail.id = Costing.receive_id)" +
						" WHERE inv_receiving_detail.quantity>0" +
						" AND inv_receiving.date<='" + StockDateStr + "'" +
						" AND inv_receiving_detail.item_code='" + ItemCode + "' " +
						"GROUP BY inv_receiving_detail.id) AS Stock" +
						" WHERE Stock.Quantity>0 " +
						"ORDER BY Stock.date " + StockOrderBy + "").Scan(&MyQStockID)

					TempQuantity := BasicQuantity
					var TotalPrice float64
					var CostingDetailID uint64
					if len(MyQStockID) > 0 {
						//Master
						if RefNumberCosting == "" {
							RefNumberCosting = GetJournalRefNumber(c, tx, global_var.JournalPrefix.InventoryCPOS, PostingDate)
							CostingNumber = GetCostingNumber(c, DB, PostingDate)
							InsertAccJournal(tx, RefNumberCosting, "", UnitCode, "", "", global_var.JournalType.Other, global_var.JournalGroup.Costing, PostingDate, "Product Costing: "+CostingNumber+", Check Number: "+CheckNumber, "", 0, UserID)
							InsertInvCosting(tx, CostingNumber, RefNumberCosting, CheckNumber, SubDepartmentCompliment, StoreCode, PostingDate, "", "Product Costing: "+CheckNumber, 0, 0, 0, 0, 0, 1, UserID)
						}
						var Quantity, Price float64
						for _, stock := range MyQStockID {
							if TempQuantity > 0 {
								if TempQuantity > general.RoundToX3(stock.Quantity) {
									Quantity = general.RoundToX3(stock.Quantity)
									TempQuantity = TempQuantity - general.RoundToX3(stock.Quantity)
									//Jika stock terakhir
									if general.RoundToX3(Quantity) == general.RoundToX3(stock.AllStoreQuantity) {
										//Jika pernah ada yg di costing ambil sisanya
										if general.RoundToX3(stock.BasicQuantity) > general.RoundToX3(stock.AllStoreQuantity) {
											Price = stock.TotalPrice - (general.RoundToX3(general.RoundToX3(stock.TotalPrice)/general.RoundToX3(stock.BasicQuantity)) * (general.RoundToX3(stock.BasicQuantity) - general.RoundToX3(stock.AllStoreQuantity)))
										} else {
											Price = general.RoundToX3(stock.TotalPrice)
										}
									} else {
										Price = general.RoundToX3(general.RoundToX3(stock.TotalPrice)/general.RoundToX3(stock.BasicQuantity)) * Quantity
									}
								} else {
									Quantity = TempQuantity
									TempQuantity = 0
									//Jika stock terakhir
									if general.RoundToX3(Quantity) == general.RoundToX3(stock.AllStoreQuantity) {
										//Jika pernah ada yg di costing ambil sisanya
										if general.RoundToX3(stock.BasicQuantity) > general.RoundToX3(stock.AllStoreQuantity) {
											Price = stock.TotalPrice - (general.RoundToX3(general.RoundToX3(stock.TotalPrice)/general.RoundToX3(stock.BasicQuantity)) * (general.RoundToX3(stock.BasicQuantity) - general.RoundToX3(stock.AllStoreQuantity)))
										} else {
											Price = general.RoundToX3(stock.TotalPrice)
										}
									} else {
										Price = general.RoundToX3(general.RoundToX3(stock.TotalPrice)/general.RoundToX3(stock.BasicQuantity)) * Quantity
									}
								}
								TotalPrice = TotalPrice + Price
								var err error
								IsCOGS := 0
								if !IsCompliment {
									IsCOGS = 1
								}
								CostingDetailID, err = InsertInvCostingDetail(tx, pConfig.Dataset, CostingNumber, StoreCode, uint64(StoreID), ItemCode, uint64(ItemID), PostingDate, Quantity, BasicUOMCode, Price, uint64(stock.ID), InventoryCOGSJournalAccount, ProductItemGroup, "", 0, uint8(IsCOGS), UserID)
								if err != nil {
									return err
								}
							} else {
								break
							}
						}
					}

					if TempQuantity < BasicQuantity {
						CostingQuantity := BasicQuantity
						if TempQuantity > 0 {
							CostingQuantity = BasicQuantity - TempQuantity
						}

						//Journal Detail Costing Debit
						InsertAccJournalDetail(tx, RefNumberCosting, UnitCode, SubDepartmentCompliment, InventoryCOGSJournalAccount, TotalPrice, global_var.TransactionType.Debit, "Product Costing: "+CostingNumber+", Check Number: "+CheckNumber+", Product: "+ProductName+", Store Code:"+StoreCode+": "+ItemName+" = "+general.FloatToStrX3(CostingQuantity)+" "+BasicUOMCode+" @"+general.FloatToStrX2(TotalPrice/CostingQuantity), "", false, UserID)
						//Journal Detail Costing Credit
						InsertAccJournalDetail(tx, RefNumberCosting, UnitCode, SubDepartmentCompliment, InventoryJournalAccount, TotalPrice, global_var.TransactionType.Credit, "Product Costing: "+CostingNumber+", Check Number: "+CheckNumber+", Product: "+ProductName+", Store Code:"+StoreCode+": "+ItemName+" = "+general.FloatToStrX3(CostingQuantity)+" "+BasicUOMCode+" @"+general.FloatToStrX3(TotalPrice/CostingQuantity), "", false, UserID)

						InsertPosProductCosting(tx, CheckNumber, CheckTransactionID, CostingNumber, CostingDetailID, ProductCode, StoreCode, ItemCode, CostRecipeQuantity, UOMCode, BasicQuantity, BasicUOMCode, CostingQuantity, UserID)
					}
				} else if pConfig.Dataset.ProgramConfiguration.CostingMethod == global_var.InventoryCostingMethod.Average {
					CostingQuantity := BasicQuantity
					StoreStock, err := GetStockStore(tx, pConfig.Dataset, StoreCode, ItemCode, PostingDate)
					if err != nil {
						return err
					}
					if StoreStock < BasicQuantity {
						CostingQuantity = StoreStock
					}

					Quantity := CostingQuantity
					var Price, TotalPrice float64

					if CostingQuantity > 0 {
						//Master
						if RefNumberCosting == "" {
							RefNumberCosting = GetJournalRefNumber(c, tx, global_var.JournalPrefix.InventoryCPOS, PostingDate)
							CostingNumber = GetCostingNumber(c, tx, PostingDate)
							InsertAccJournal(tx, RefNumberCosting, "", UnitCode, "", "", global_var.JournalType.Other, global_var.JournalGroup.Costing, PostingDate, "Product Costing: "+CostingNumber+", Check Number: "+CheckNumber, "", 0, UserID)
							InsertInvCosting(tx, CostingNumber, RefNumberCosting, CheckNumber, SubDepartmentCompliment, StoreCode, PostingDate, "", "Product Costing: "+CheckNumber, 0, 0, 0, 0, 0, 1, UserID)
						}

						CostingDetailID, err := InsertInvCostingDetail(tx, pConfig.Dataset, CostingNumber, StoreCode, uint64(StoreID), ItemCode, uint64(ItemID), PostingDate, Quantity, BasicUOMCode, Price, 0, InventoryCOGSJournalAccount, ProductItemGroup, "", 0, 1, UserID)
						if err != nil {
							return err
						}
						//Journal Detail Costing Debit
						JournalDetailDebit = append(JournalDetailDebit, db_var.Acc_journal_detail{
							RefNumber:         RefNumberCosting,
							Date:              PostingDate,
							UnitCode:          UnitCode,
							SubDepartmentCode: SubDepartmentCompliment,
							AccountCode:       InventoryCOGSJournalAccount,
							Amount:            general.RoundToX3(TotalPrice),
							TypeCode:          global_var.TransactionType.Debit,
							Remark:            "Product: " + ProductName + ", " + StoreCode + ": " + ItemName + " = " + general.FloatToStrX3(CostingQuantity) + " " + BasicUOMCode + " @" + general.FloatToStrX3(TotalPrice/CostingQuantity),
							CreatedBy:         UserID,
						})
						//InsertJournalDetail(RefNumberCosting, SubDepartmentCompliment, InventoryCOGSJournalAccount, TransactionType.Debit, "Product Costing: " + CostingNumber + ", Check Number: " + CheckNumber + ", Product: " + ProductName + ", " + StoreCode + ": " + ItemName + " = " + FormatFloatX(CostingQuantity) + " " + BasicUOMCode + " @" + FormatFloatX(TotalPrice / CostingQuantity), "", TotalPrice)
						//Journal Detail Costing Credit
						JournalDetailCredit = append(JournalDetailCredit, db_var.Acc_journal_detail{
							RefNumber:         RefNumberCosting,
							Date:              PostingDate,
							UnitCode:          UnitCode,
							SubDepartmentCode: SubDepartmentCompliment,
							AccountCode:       InventoryJournalAccount,
							Amount:            general.RoundToX3(TotalPrice),
							TypeCode:          global_var.TransactionType.Credit,
							Remark:            "Product Costing: " + CostingNumber + ", Check Number: " + CheckNumber + ", Product: " + ProductName + ", " + StoreCode + ": " + ItemName + " = " + general.FloatToStrX2(CostingQuantity) + " " + BasicUOMCode + " @" + general.FloatToStrX3(TotalPrice/CostingQuantity),
							CreatedBy:         UserID,
						})
						//InsertJournalDetail(RefNumberCosting, SubDepartmentCompliment, InventoryJournalAccount, TransactionType.Credit, "Product Costing: " + CostingNumber + ", Check Number: " + CheckNumber + ", Product: " + ProductName + ", " + StoreCode + ": " + ItemName + " = " + FormatFloatX(CostingQuantity) + " " + BasicUOMCode + " @" + FormatFloatX(TotalPrice / CostingQuantity), "", TotalPrice)
						ProductCosting = append(ProductCosting, db_var.Pos_product_costing{
							CheckNumber:        CheckNumber,
							CheckTransactionId: CheckTransactionID,
							CostingNumber:      CostingNumber,
							CostingDetailId:    CostingDetailID,
							ProductCode:        ProductCode,
							StoreCode:          StoreCode,
							ItemCode:           ItemCode,
							Quantity:           general.RoundToX3(CostRecipeQuantity),
							UomCode:            UOMCode,
							BasicQuantity:      general.RoundToX3(BasicQuantity),
							BasicUomCode:       BasicUOMCode,
							CostingQuantity:    general.RoundToX3(CostingQuantity),
							CreatedBy:          UserID,
						})
						//InsertProductCosting(CheckNumber, CheckTransactionID, CostingNumber, CostingDetailID, ProductCode, StoreCode, ItemCode, UOMCode, BasicUOMCode, CostRecipeQuantity, BasicQuantity, CostingQuantity)
					}
				}
			}
			if len(JournalDetailCredit) > 0 {
				if err := BatchInsertAccJournalDetail(tx, RefNumberCosting, JournalDetailCredit, false); err != nil {
					return err
				}
			}
			if len(JournalDetailDebit) > 0 {
				if err := BatchInsertAccJournalDetail(tx, RefNumberCosting, JournalDetailDebit, false); err != nil {
					return err
				}
			}
			if len(ProductCosting) > 0 {
				if err := BatchInsertPosProductCosting(tx, ProductCosting); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			return "", "", err
		}
	}
	return RefNumberCosting, CostingNumber, nil
}

func GetDatabaseVersion(DB *gorm.DB) (Version int, err error) {
	if err := DB.Table("schema_migrations").Select("version").Limit(1).Scan(&Version).Error; err != nil {
		return 0, err
	}
	return Version, nil
}

func GetCurrentGlobalDatabaseVersion(DB *gorm.DB) (Version int, err error) {
	if err := DB.Table("version").Select("version").Order("version DESC").Limit(1).Scan(&Version).Error; err != nil {
		return 0, err
	}
	return Version, nil
}

func GetInventoryLastClosedDate(DB *gorm.DB) (time.Time, error) {
	var lastClosedDate time.Time
	if err := DB.Raw(`SELECT closed_date FROM inv_close_log ORDER BY closed_date DESC LIMIT 1`).Scan(&lastClosedDate).Error; err != nil {
		return time.Time{}, err
	}

	return lastClosedDate, nil
}

func IsInventoryClosedAverageMethod(DB *gorm.DB, costingDate time.Time) (bool, error) {
	lastClosedDate, err := GetInventoryLastClosedDate(DB)
	if err != nil {
		return false, err
	}
	// Check if lastClosedDate is not the zero value of time.Time
	if !lastClosedDate.IsZero() {
		if costingDate.Before(lastClosedDate) || costingDate.Equal(lastClosedDate) {
			return true, nil
		}
	}
	return false, nil
}

func IsVenueAvailable(DB *gorm.DB, Timezone string, Mode byte, ReservationNumber uint64, StartDate, EndDate time.Time, VenueCode string) (bool, error) {
	loc, _ := time.LoadLocation(Timezone)
	StartDate = StartDate.In(loc)
	EndDate = EndDate.In(loc)
	CountBeetween := general.DaysBetween(StartDate, EndDate)
	TempDate := StartDate
	QueryQonditionMode := ""

	//Check Venue Combine
	var VenueCodeData []string
	DB.Raw(
		"SELECT" +
			" venue_code " +
			"FROM" +
			" ban_cfg_init_venue_combine_detail " +
			"WHERE ban_cfg_init_venue_combine_detail.combine_venue_code = '" + VenueCode + "'").Scan(&VenueCodeData)

	CountCheck := 0
	if len(VenueCodeData) <= 0 {
		CountCheck = 1
	} else {
		CountCheck = len(VenueCodeData)
	}
	//=================================================
	for CountX := 1; CountX <= CountCheck; CountX++ {
		if CountCheck > 1 {
			VenueCode = VenueCodeData[0]
		}
		if (ReservationNumber > 0) && ((Mode == 1) || (Mode == 3)) { //Mode Update & CheckIn
			if CountCheck > 0 {
				//Get Combine Venue Number
				var VenueCombineNumber uint64
				if err := DB.Raw(
					"SELECT"+
						" venue_combine_number "+
						"FROM"+
						" ban_reservation "+
						"WHERE ban_reservation.number = ?", ReservationNumber).Scan(&VenueCombineNumber).Error; err != nil {
					return false, err
				}
				//Get Reservation
				var ReservationNumberX uint64
				if err := DB.Raw(
					"SELECT"+
						" number "+
						"FROM"+
						" ban_reservation "+
						"WHERE ban_reservation.venue_combine_number = ? "+
						" AND ban_reservation.venue_code = '"+VenueCode+"'", VenueCombineNumber).Scan(&ReservationNumberX).Error; err != nil {
					return false, err
				}
				//        ShowMessage(DataModuleMain.MyQGetReservationNumber.Fields.Fields[0].AsString)
				QueryQonditionMode = " AND ban_reservation.number <> " + general.Uint64ToStr(ReservationNumberX) + " "
			} else {
				QueryQonditionMode = " AND ban_reservation.number <> " + general.Uint64ToStr(ReservationNumber) + " "
			}
		}

		for Count := 0; Count <= CountBeetween; Count++ {
			var CountResv uint64
			if err := DB.Raw(
				"SELECT" +
					" COUNT(number) AS CountResv " +
					"FROM ban_reservation " +
					" LEFT OUTER JOIN guest_detail ON (ban_reservation.guest_detail_id = guest_detail.id) " +
					"WHERE (ban_reservation.status_code = '" + global_var.ReservationStatus.Reservation + "' OR ban_reservation.status_code = '" + global_var.ReservationStatus.InHouse + "') " +
					" AND DATE(guest_detail.arrival) = '" + general.FormatDate1(TempDate) + "' " +
					" AND ban_reservation.venue_code = '" + VenueCode + "' " +
					" AND (TIME(convert_tz(guest_detail.arrival,'UTC','" + Timezone + "')) <= '" + general.FormatTime1(EndDate) + "' AND TIME(convert_tz(guest_detail.departure,'UTC','" + Timezone + "')) >= '" + general.FormatTime1(StartDate) + "') " +
					QueryQonditionMode +
					"LIMIT 1").Scan(&CountResv).Error; err != nil {
				return false, err
			}

			if CountResv > 0 {
				return false, nil
			}
		}
		TempDate = general.IncDay(TempDate, 1)
	}
	return true, nil
}

func GetPosCheckTransactionBreakdown1() uint64 {
	uuid := time.Now().UnixMilli()
	return uint64(uuid)
	// var BreakDown1 float64
	//     DeleteSQLX("temp_pos_check_transaction_breakdown1", "'" +ProgramVariable.IPAddress+ "'", False)

	//     while Result = "" do
	// //      //Added this code for optimized
	// //      ChangeQueryString(MyQGeneral,
	// //        "SELECT" +
	// //        " A.breakdown1 " +
	// //        "FROM " +
	// //          "(SELECT breakdown1 FROM sub_folio " +
	// //          "WHERE sub_folio.audit_date BETWEEN  DATE(DATE_ADD('" +FormatDateTimeX(ProgramVariable.AuditDate)+ "', INTERVAL -1 DAY)) AND DATE('" +FormatDateTimeX(ProgramVariable.AuditDate)+ "') " +
	// //          ") AS A " +
	// //        "ORDER BY A.breakdown1 DESC " +
	// //        "LIMIT 1",
	// //        "", "", "", "", "", "", "", "", "", "")
	// //      //End Optimize
	// //
	// //      if MyQGeneral.IsEmpty {
	//       ChangeQueryString(MyQGeneral,
	//         "SELECT breakdown1 FROM pos_check_transaction " +
	//         "ORDER BY breakdown1 DESC " +
	//         "LIMIT 1",
	//         "", "", "", "", "", "", "", "", "", "")

	//       if MyQGeneral.IsEmpty {
	//         BreakDown1 = 1
	//       } else {
	//         BreakDown1 = MyQGeneral.Fields.Fields[0].AsLargeInt + 1

	//       if CheckCode("temp_pos_check_transaction_breakdown1", "breakdown1", "breakdown1", strconv.FormatInt(BreakDown1)) {
	//         ChangeQueryString(MyQGeneral,
	//           "SELECT breakdown1 FROM temp_pos_check_transaction_breakdown1 " +
	//           "ORDER BY breakdown1 DESC " +
	//           "LIMIT 1",
	//           "", "", "", "", "", "", "", "", "", "")

	//         if not MyQGeneral.IsEmpty {
	//           BreakDown1 = MyQGeneral.Fields.Fields[0].AsLargeInt + 1
	//       }

	//       try
	//         InsertSQLX("temp_pos_check_transaction_breakdown1", strconv.FormatInt(BreakDown1)+ ", '" +ProgramVariable.IPAddress+ "'", False)
	//       finally
	//         ChangeQueryString(MyQGeneral,
	//           "SELECT breakdown1 FROM temp_pos_check_transaction_breakdown1" +
	//           " WHERE ip_address = '" +ProgramVariable.IPAddress+ "' " +
	//           "ORDER BY breakdown1 DESC " +
	//           "LIMIT 1",
	//           "", "", "", "", "", "", "", "", "", "")

	//         if not MyQGeneral.IsEmpty {
	//           Result = MyQGeneral.Fields.Fields[0].AsString
	//       }
	//     }
	//   }
}

func GetPrepaidExpenseOutstanding(DB *gorm.DB, PrepaidID, PrepaidPostedID uint64) (float64, error) {
	QueryCondition := ""
	if PrepaidPostedID > 0 {
		QueryCondition = " AND acc_prepaid_expense_posted.id<>" + general.Uint64ToStr(PrepaidPostedID)
	}
	var AmountPosted float64
	if err := DB.Raw(
		"SELECT"+
			" (acc_prepaid_expense.amount - IFNULL(SUM(acc_prepaid_expense_posted.amount), 0)) AS AmountPosted "+
			"FROM"+
			" acc_prepaid_expense"+
			" LEFT OUTER JOIN acc_prepaid_expense_posted ON (acc_prepaid_expense.id = acc_prepaid_expense_posted.prepaid_id "+QueryCondition+")"+
			" WHERE acc_prepaid_expense.id=? "+
			//      QueryCondition+ " " +
			" GROUP BY acc_prepaid_expense.id", PrepaidID).Scan(&AmountPosted).Error; err != nil {
		return 0, err
	}
	return AmountPosted, nil
}

func GetDifferedIncomeOutstanding(DB *gorm.DB, DefferedID, DefferedPostedID uint64) (float64, error) {
	QueryCondition := ""
	if DefferedPostedID > 0 {
		QueryCondition = " AND acc_deffered_income_posted.id<>" + general.Uint64ToStr(DefferedPostedID)
	}
	var AmountPosted float64
	if err := DB.Raw(
		"SELECT"+
			" (acc_deffered_income.amount - IFNULL(SUM(acc_deffered_income_posted.amount), 0)) AS AmountPosted "+
			"FROM"+
			" acc_deffered_income"+
			" LEFT OUTER JOIN acc_deffered_income_posted ON (acc_deffered_income.id = acc_deffered_income_posted.deffered_id "+QueryCondition+")"+
			" WHERE acc_deffered_income.id=? "+
			//      QueryCondition+ " " +
			" GROUP BY acc_deffered_income.id", DefferedID).Scan(&AmountPosted).Error; err != nil {
		return 0, err
	}
	return AmountPosted, nil
}
