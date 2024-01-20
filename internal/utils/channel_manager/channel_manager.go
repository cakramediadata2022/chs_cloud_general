package channel_manager

import (
	"bytes"
	initConfig "chs_cloud_general/config"
	"chs_cloud_general/internal/config"
	"chs_cloud_general/internal/db_var"
	"chs_cloud_general/internal/general"
	"chs_cloud_general/internal/global_query"
	"chs_cloud_general/internal/global_var"
	"chs_cloud_general/pkg/utils"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

var baseURL string

var username = "ketut"
var password = "ketut123456"
var encodedCredentials string
var credentials string
var cmConfig *initConfig.CM

func CMInit(config *initConfig.CM) {
	cmConfig = config
	baseURL = config.CXURL

	credentials = config.Username + ":" + config.Password
	encodedCredentials = base64.StdEncoding.EncodeToString([]byte(credentials))
}

func CMScheduler(c *gin.Context) {

	// // Get Program Configuration
	// 	val, exist := c.Get("pConfig")
	// 	if !exist {
	// 		master_data.SendResponse(global_var.ResponseCode.DatabaseError, "pConfig", nil, c)
	// 		return
	// 	}
	// 	pConfig := val.(*config.CompanyDataConfiguration)
	// 	DB := pConfig.DB
}

func CMPushAvailability(c *gin.Context, DB *gorm.DB, StartDate, EndDate time.Time, RoomTypeCode string, SyncAll bool) error {
	ctx, span := global_var.Tracer.Start(c, "MoveSubFolioP")
	defer span.End()

	// Get Program Configuration
	val, exist := c.Get("pConfig")
	if !exist {
		utils.LogResponseError(c, ctx, global_var.AppLogger, errors.Wrap(nil, "CMPushAvailability.pConfig Not Found"))
		return errors.New("pConfig Not Found")
	}
	pConfig := val.(*config.CompanyDataConfiguration)
	url := ""
	if pConfig.Dataset.ProgramConfiguration.CCMSVendor != "CNX" {
		baseURL = cmConfig.RGURL
		url = fmt.Sprintf("%s/RGOTAHotelAvailNotifRQ", baseURL)
	} else {
		baseURL = cmConfig.CXURL
		url = fmt.Sprintf("%s/UpdateAvailability", baseURL)
	}

	CCMSSMHotelCode := pConfig.Dataset.ProgramConfiguration.CCMSSMHotelCode
	CCMSSMPassword := pConfig.Dataset.ProgramConfiguration.CCMSSMPassword
	CCMSSMRequestorID := pConfig.Dataset.ProgramConfiguration.CCMSSMRequestorID
	CCMSSMWSDL := pConfig.Dataset.ProgramConfiguration.CCMSSMWSDL
	CCMSSMUser := pConfig.Dataset.ProgramConfiguration.CCMSSMUser
	CCMSSMSynchronizeAvailability := pConfig.Dataset.ProgramConfiguration.CCMSSMSynchronizeAvailability

	if !CCMSSMSynchronizeAvailability {
		return nil
	}

	RoomType := []string{}

	type AvailabilityStruct struct {
		RoomTypeCode string
		Availability int64
		BookingLimit int64
		StartDate    time.Time
		EndDate      time.Time
	}

	Query := DB.Table(db_var.TableName.CfgInitRoomType).Distinct("code").
		Joins("INNER JOIN cfg_init_room ON cfg_init_room_type.code = cfg_init_room.room_type_code")

	if !SyncAll {
		Query.Where("cfg_init_room_type.code", RoomTypeCode)
	}

	Query.Scan(&RoomType)

	Date := StartDate
	Data := make([]map[string]interface{}, 0)
	DataMap := make(map[string]AvailabilityStruct)
	Nights := general.DaysBetween(StartDate, EndDate)
	for _, roomType := range RoomType {
		for i := 0; i < Nights; i++ {
			ArrivalDate := Date.AddDate(0, 0, i)
			DepartureDate := ArrivalDate.AddDate(0, 0, 1)
			Avail, _ := global_query.GetAvailableRoomCountByType(DB, ArrivalDate, DepartureDate, roomType, "", 0, 0, 0, 0, false, false)

			if (i > 0 && DataMap[roomType].Availability != Avail) || i >= Nights-1 {
				Data = append(Data, map[string]interface{}{
					"start_date":     DataMap[roomType].StartDate,
					"end_date":       DepartureDate.AddDate(0, 0, -1),
					"hotel_code":     c.GetString("UnitCode"),
					"room_type_code": roomType,
					"booking_limit":  DataMap[roomType].Availability,
					"availability":   DataMap[roomType].Availability,
				})

				DataMap[roomType] = AvailabilityStruct{
					StartDate:    ArrivalDate,
					Availability: Avail,
					BookingLimit: Avail,
					RoomTypeCode: roomType,
				}
			}

			if i == 0 {
				DataMap[roomType] = AvailabilityStruct{
					StartDate:    ArrivalDate,
					Availability: Avail,
					BookingLimit: Avail,
					RoomTypeCode: roomType,
				}
			}

		}
	}

	data := gin.H{
		"hotel_code":   CCMSSMHotelCode,
		"hotel_name":   "",
		"user":         CCMSSMUser,
		"password":     CCMSSMPassword,
		"requestor_id": CCMSSMRequestorID,
		"wsdl":         CCMSSMWSDL,
		"details":      Data,
	}
	// Marshal the data into a JSON string.
	jsonData, err := json.Marshal(data)
	if err != nil {
		utils.LogResponseError(c, ctx, global_var.AppLogger, errors.Wrap(err, "CMPushAvailability.json.Marshal"))
		return err
	}

	fmt.Println("listJson2", bytes.NewBuffer(jsonData))
	// master_data.SendResponse(global_var.ResponseCode.Successfully, "", Data, c)
	// Create a request with the JSON body.
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		utils.LogResponseError(c, ctx, global_var.AppLogger, errors.Wrap(err, "CMPushAvailability.NewRequest"))
		return err
	}

	// Set the request headers.
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+encodedCredentials)

	// Create an HTTP client and send the request.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		utils.LogResponseError(c, ctx, global_var.AppLogger, errors.Wrap(err, "CMPushAvailability.client.Do"))
		return err
	}
	defer resp.Body.Close()

	// Read the response body.
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		// Handle the error when reading the response body.
		utils.LogResponseError(c, ctx, global_var.AppLogger, errors.Wrap(err, "CMPushAvailability.ReadAll"))
		return err
	}

	// Convert the response body to a string.
	responseString := string(responseBody)

	// Now you can work with the responseString as needed.
	fmt.Println("Response body:", responseString)

	if resp.Status == "200 OK" {
		fmt.Println("Request was successful!")
	} else {
		utils.LogResponseError(c, ctx, global_var.AppLogger, errors.Wrap(nil, "CMPushAvailability.resp"))
		return errors.New(responseString)
	}
	return nil

}

func CMPushRate(c *gin.Context, RateCode string, SyncAll bool) error {
	ctx, span := global_var.Tracer.Start(c, "MoveSubFolioP")
	defer span.End()

	type Rate struct {
		StartDate         time.Time
		EndDate           time.Time
		HotelCode         string
		RatePlanCode      string
		Rate              string
		RoomTypeCode      string
		RatePlanName      string
		Day1              int
		Day2              int
		Day3              int
		Day4              int
		Day5              int
		Day6              int
		Day7              int
		StopSell          int
		ClosedToArrival   int
		ClosedToDeparture int
	}
	// Get Program Configuration
	val, exist := c.Get("pConfig")
	if !exist {
		utils.LogResponseError(c, ctx, global_var.AppLogger, errors.Wrap(nil, "CMPushRate.pConfig Not Found"))
		return errors.New("pConfig Not Found")
	}
	pConfig := val.(*config.CompanyDataConfiguration)
	DB := pConfig.DB

	url := ""
	if pConfig.Dataset.ProgramConfiguration.CCMSVendor != "CNX" {
		baseURL = cmConfig.RGURL
		url = fmt.Sprintf("%s/RGOTAHotelRateAmountNotifRQ", baseURL)
	} else {
		baseURL = cmConfig.CXURL
		url = fmt.Sprintf("%s/UpdateRatePlan", baseURL)
	}

	CCMSSMHotelCode := pConfig.Dataset.ProgramConfiguration.CCMSSMHotelCode
	CCMSSMPassword := pConfig.Dataset.ProgramConfiguration.CCMSSMPassword
	CCMSSMRequestorID := pConfig.Dataset.ProgramConfiguration.CCMSSMRequestorID
	CCMSSMWSDL := pConfig.Dataset.ProgramConfiguration.CCMSSMWSDL
	CCMSSMUser := pConfig.Dataset.ProgramConfiguration.CCMSSMUser
	CCMSSMSynchronizeRate := pConfig.Dataset.ProgramConfiguration.CCMSSMSynchronizeRate

	if !CCMSSMSynchronizeRate {
		return nil
	}

	RoomRate := []db_var.Cfg_init_room_rate{}
	Query := DB.Table(db_var.TableName.CfgInitRoomRate).Where("is_online", "1").Where("cm_inv_code <> ?", "")

	if !SyncAll {
		Query.Where("code", RateCode).Limit(1)
	}

	if err := Query.Scan(&RoomRate).Error; err != nil {
		utils.LogResponseError(c, ctx, global_var.AppLogger, errors.Wrap(err, "CMPushRate.Scan"))
		return err
	}

	if len(RoomRate) <= 0 {
		return errors.New("Rate code not found")
	}

	var rateList = make([]map[string]interface{}, 0)
	for _, rate := range RoomRate {
		rateList = append(rateList, map[string]interface{}{
			"start_date":          rate.CmStartDate,
			"end_date":            rate.CmEndDate,
			"hotel_code":          CCMSSMHotelCode,
			"rate_plan_code":      rate.Code,
			"inv_code":            rate.CmInvCode,
			"rate":                general.FloatToStrX3(rate.WeekdayRate1),
			"room_type_code":      rate.CmInvCode,
			"rate_plan_name":      rate.Name,
			"day1":                rate.Day1,
			"day2":                rate.Day2,
			"day3":                rate.Day3,
			"day4":                rate.Day4,
			"day5":                rate.Day5,
			"day6":                rate.Day6,
			"day7":                rate.Day7,
			"stop_sell":           rate.CmStopSell,
			"closed_to_arrival":   0,
			"closed_to_departure": 0,
		})
	}

	data := gin.H{
		"hotel_code":   CCMSSMHotelCode,
		"hotel_name":   "",
		"user":         CCMSSMUser,
		"password":     CCMSSMPassword,
		"requestor_id": CCMSSMRequestorID,
		"wsdl":         CCMSSMWSDL,
		"details":      rateList,
	}

	// Marshal the data into a JSON string.
	jsonData, err := json.Marshal(data)
	if err != nil {
		utils.LogResponseError(c, ctx, global_var.AppLogger, errors.Wrap(err, "CMPushRate.json.Marshal"))
		return err
	}

	fmt.Println("listJson2", bytes.NewBuffer(jsonData))
	// Create a request with the JSON body.
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		utils.LogResponseError(c, ctx, global_var.AppLogger, errors.Wrap(err, "CMPushRate.NewRequest"))
		return err
	}

	// Set the request headers.
	req.Header.Set("Content-Type", "application/json")

	// Create an HTTP client and send the request.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		utils.LogResponseError(c, ctx, global_var.AppLogger, errors.Wrap(err, "CMPushRate.client.Do"))
		return err
	}
	defer resp.Body.Close()

	// Read the response body.
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		// Handle the error when reading the response body.
		utils.LogResponseError(c, ctx, global_var.AppLogger, errors.Wrap(err, "CMPushRate.ReadAll"))
		return err
	}

	// Convert the response body to a string.
	responseString := string(responseBody)

	// Now you can work with the responseString as needed.
	fmt.Println("Response body:", responseString)

	if resp.Status == "200 OK" {
		fmt.Println("Request was successful!")

		// Read the response body.
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			// Handle the error when reading the response body.
			utils.LogResponseError(c, ctx, global_var.AppLogger, errors.Wrap(err, "CMPushRate.ReadAll"))
			return err
		}

		// Convert the response body to a string.
		responseString := string(responseBody)

		// Now you can work with the responseString as needed.
		fmt.Println("Response body:", responseString)
	} else {
		utils.LogResponseError(c, ctx, global_var.AppLogger, errors.Wrap(nil, "CMPushRate.resp"))
		return errors.New(resp.Status)
	}
	return nil
}

func UpdateAvailability(c *gin.Context, DB *gorm.DB, RoomTypeCode string, NewStartDate, OldStartDate time.Time, NewEndDate, OldEndDate time.Time, isNew bool) error {
	StartDate := NewStartDate
	EndDate := NewEndDate
	if OldStartDate != (time.Time{}) && NewStartDate.After(OldStartDate) {
		StartDate = OldStartDate
	}

	if OldEndDate != (time.Time{}) && OldEndDate.After(NewEndDate) {
		EndDate = OldStartDate
	}

	CMPushAvailability(c, DB, StartDate, EndDate, RoomTypeCode, false)

	return nil
}
