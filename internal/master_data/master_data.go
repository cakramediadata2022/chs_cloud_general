package master_data

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cakramediadata2022/chs_cloud_general/internal/config"
	"github.com/cakramediadata2022/chs_cloud_general/internal/db_var"
	"github.com/cakramediadata2022/chs_cloud_general/internal/general"
	"github.com/cakramediadata2022/chs_cloud_general/internal/global_var"
	"github.com/cakramediadata2022/chs_cloud_general/pkg/utils"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

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

func SendResponse(StatusCode uint, Message interface{}, Result interface{}, c *gin.Context) {
	var RequestResponse = global_var.TRequestResponse{
		StatusCode: StatusCode,
		Message:    Message,
		Result:     Result}

	if RequestResponse.StatusCode == global_var.ResponseCode.Successfully || RequestResponse.StatusCode == global_var.ResponseCode.SuccessfullyWithStatus {
		c.JSON(http.StatusOK, &RequestResponse)
	} else if (RequestResponse.StatusCode == global_var.ResponseCode.NotAuthorized) || (RequestResponse.StatusCode == global_var.ResponseCode.ErrorCreateToken) {
		c.JSON(http.StatusUnauthorized, &RequestResponse)
	} else if (RequestResponse.StatusCode == global_var.ResponseCode.InvalidDataFormat) || (RequestResponse.StatusCode == global_var.ResponseCode.DataNotFound) ||
		(RequestResponse.StatusCode == global_var.ResponseCode.InvalidDataValue) || (RequestResponse.StatusCode == global_var.ResponseCode.DatabaseValueChanged) ||
		(RequestResponse.StatusCode == global_var.ResponseCode.DatabaseError) || (RequestResponse.StatusCode == global_var.ResponseCode.DuplicateEntry) ||
		(RequestResponse.StatusCode == global_var.ResponseCode.OtherResult) || (RequestResponse.StatusCode == global_var.ResponseCode.Unregistered) || (RequestResponse.StatusCode == global_var.ResponseCode.SubscriptionExpired) {
		c.JSON(http.StatusBadRequest, &RequestResponse)
	} else if RequestResponse.StatusCode == global_var.ResponseCode.InternalServerError {
		c.JSON(http.StatusInternalServerError, &RequestResponse)
	} else {
		c.JSON(http.StatusBadRequest, &RequestResponse)
	}
}

func SendWebsocketResponse(StatusCode uint, Message interface{}, Result interface{}, con *websocket.Conn) {
	var RequestResponse = global_var.TRequestResponse{
		StatusCode: StatusCode,
		Message:    Message,
		Result:     Result}

	if RequestResponse.StatusCode == global_var.ResponseCode.Successfully || RequestResponse.StatusCode == global_var.ResponseCode.SuccessfullyWithStatus {
		con.WriteJSON(&RequestResponse)
	} else if (RequestResponse.StatusCode == global_var.ResponseCode.NotAuthorized) || (RequestResponse.StatusCode == global_var.ResponseCode.ErrorCreateToken) {
		con.WriteJSON(&RequestResponse)
	} else if (RequestResponse.StatusCode == global_var.ResponseCode.InvalidDataFormat) || (RequestResponse.StatusCode == global_var.ResponseCode.DataNotFound) ||
		(RequestResponse.StatusCode == global_var.ResponseCode.InvalidDataValue) || (RequestResponse.StatusCode == global_var.ResponseCode.DatabaseValueChanged) ||
		(RequestResponse.StatusCode == global_var.ResponseCode.DatabaseError) || (RequestResponse.StatusCode == global_var.ResponseCode.DuplicateEntry) ||
		(RequestResponse.StatusCode == global_var.ResponseCode.OtherResult) {
		con.WriteJSON(&RequestResponse)
	} else {
		con.WriteJSON(Message)
	}
}

// func GetUnitCode(c *gin.Context, DB *gorm.DB) string {
// 	var Code string
// 	DB.Table(db_var.TableName.HotelInformation).Select("code").Take(&Code)

// 	return Code
// }

func GetConfigurationAllP(c *gin.Context) {
	// Get Program Configuration
	val, exist := c.Get("pConfig")
	if !exist {
		SendResponse(global_var.ResponseCode.DatabaseError, "pConfig", nil, c)
		return
	}
	pConfig := val.(*config.CompanyDataConfiguration)
	DB := pConfig.DB
	var ConfigurationData []map[string]interface{}
	DB.Table(db_var.TableName.Configuration).Select("system_code", "category", "name", "value", "default_value").Scan(&ConfigurationData)
	SendResponse(global_var.ResponseCode.Successfully, "", ConfigurationData, c)
}

func GetConfiguration(DB *gorm.DB, SystemCode, Category, Name string, GetDefault bool) interface{} {
	var ConfigurationData db_var.Configuration
	DB.Table(db_var.TableName.Configuration).Where("system_code = ? AND category = ? AND name = ?", SystemCode, Category, Name).Find(&ConfigurationData)
	if GetDefault {
		return ConfigurationData.DefaultValue
	} else {
		return ConfigurationData.Value
	}
}

func GetConfigurationBool(DB *gorm.DB, SystemCode, Category, Name string, GetDefault bool) bool {
	t := GetConfiguration(DB, SystemCode, Category, Name, GetDefault)
	if strings.ToUpper(t.(string)) == "TRUE" {
		return true
	} else {
		return false
	}
}

func GetConfigurationString(DB *gorm.DB, SystemCode, Category, Name string, GetDefault bool) string {
	t := GetConfiguration(DB, SystemCode, Category, Name, GetDefault)
	return t.(string)
}

func GetConfigurationInt64(c *gin.Context, DB *gorm.DB, SystemCode, Category, Name string, GetDefault bool) int64 {
	t := GetConfiguration(DB, SystemCode, Category, Name, GetDefault)
	return t.(int64)
}

func GetConfigurationFloat64(c *gin.Context, DB *gorm.DB, SystemCode, Category, Name string, GetDefault bool) float64 {
	t := GetConfiguration(DB, SystemCode, Category, Name, GetDefault)
	return t.(float64)
}

func GetGeneralCodeName(DB *gorm.DB, TableName, OrderCondition string) []db_var.GeneralCodeNameStruct {
	var GeneralCodeName []db_var.GeneralCodeNameStruct
	DB.Table(TableName).Order(OrderCondition).Find(&GeneralCodeName)
	return GeneralCodeName
}

func GetGeneralCodeDescription(DB *gorm.DB, TableName, OrderCondition string) []db_var.GeneralCodeDescriptionStruct {
	var GeneralCodeDescription []db_var.GeneralCodeDescriptionStruct
	DB.Table(TableName).Order(OrderCondition).Find(&GeneralCodeDescription)
	return GeneralCodeDescription
}

func GetAccountField(DB *gorm.DB, Field, ConditionField string, ConditionValue interface{}) string {
	var Data string
	DB.Table(db_var.TableName.CfgInitAccount).Select(Field).Joins("JOIN cfg_init_account_sub_group ON (cfg_init_account.sub_group_code = cfg_init_account_sub_group.code)").
		Where(ConditionField+"=?", ConditionValue).Limit(1).Scan(&Data)

	return Data
}

func GetGeneralCodeNameWithParam(DB *gorm.DB, TableName, FieldConditions string, Condition1, Condition2, Condition3, Condition4, Condition5, OrderCondition interface{}, CountCondition byte) []db_var.GeneralCodeNameStruct {
	var GeneralCodeName []db_var.GeneralCodeNameStruct
	if CountCondition == 0 {
		DB.Table(TableName).Order(OrderCondition).Find(&GeneralCodeName)
	} else if CountCondition == 1 {
		DB.Table(TableName).Where(FieldConditions, Condition1).Order(OrderCondition).Find(&GeneralCodeName)
	} else if CountCondition == 2 {
		DB.Table(TableName).Where(FieldConditions, Condition1, Condition2).Order(OrderCondition).Find(&GeneralCodeName)
	} else if CountCondition == 3 {
		DB.Table(TableName).Where(FieldConditions, Condition1, Condition2, Condition3).Order(OrderCondition).Find(&GeneralCodeName)
	} else if CountCondition == 4 {
		DB.Table(TableName).Where(FieldConditions, Condition1, Condition2, Condition3, Condition4).Order(OrderCondition).Find(&GeneralCodeName)
	} else if CountCondition == 5 {
		DB.Table(TableName).Where(FieldConditions, Condition1, Condition2, Condition3, Condition4, Condition5).Order(OrderCondition).Find(&GeneralCodeName)
	}
	return GeneralCodeName
}

func GetGeneralCodeNameCondition(DB *gorm.DB, TableName, OrderCondition string, conds ...interface{}) []db_var.GeneralCodeNameStruct {
	var GeneralCodeName []db_var.GeneralCodeNameStruct
	DB.Table(TableName).Select("code,name").Order(OrderCondition).Find(&GeneralCodeName, conds...)
	return GeneralCodeName
}

func GetGeneralCodeNameQuery(DB *gorm.DB, Query string, Condition1, Condition2, Condition3, Condition4, Condition5 interface{}, CountCondition byte) []db_var.GeneralCodeNameStruct {
	var GeneralCodeName []db_var.GeneralCodeNameStruct
	if CountCondition == 0 {
		DB.Raw(Query).Scan(&GeneralCodeName)
	} else if CountCondition == 1 {
		DB.Raw(Query, Condition1).Scan(&GeneralCodeName)
	} else if CountCondition == 2 {
		DB.Raw(Query, Condition1, Condition2).Scan(&GeneralCodeName)
	} else if CountCondition == 3 {
		DB.Raw(Query, Condition1, Condition2, Condition3).Scan(&GeneralCodeName)
	} else if CountCondition == 4 {
		DB.Raw(Query, Condition1, Condition2, Condition3, Condition4, Condition5).Scan(&GeneralCodeName)
	} else if CountCondition == 5 {
		DB.Raw(Query, Condition1, Condition2, Condition3, Condition4, Condition5).Scan(&GeneralCodeName)
	}
	return GeneralCodeName
}

func GetFieldTime(DB *gorm.DB, TableName, Field, ConditionField, ConditionValue string, Default time.Time) time.Time {
	var Result []time.Time
	DB.Table(TableName).Select(Field).Where(ConditionField+" = ?", ConditionValue).Find(&Result)
	if len(Result) > 0 {
		return Result[0]
	}
	return Default
}

func GetFieldTimeQuery(DB *gorm.DB, Query string) time.Time {
	var Data []interface{}
	DB.Raw(Query).First(&Data)
	if len(Data) > 0 {
		return Data[0].(time.Time)
	} else {
		return time.Time{}
	}
}

func GetFieldString(DB *gorm.DB, TableName, Field, ConditionField string, ConditionValue interface{}, Default string) string {
	var Result string
	DB.Table(TableName).Select(Field).Where(ConditionField+" = ?", ConditionValue).Limit(1).Find(&Result)
	if Result != "" {
		return Result
	}
	return Default
}

func GetFieldStringQuery(DB *gorm.DB, Query string, Condition interface{}, Default string) string {
	var Result []string

	if Condition == nil {
		DB.Raw(Query).Find(&Result)
	} else {
		DB.Raw(Query, Condition).Find(&Result)
	}

	if len(Result) > 0 {
		return Result[0]
	}
	return Default
}

func GetFieldFloat(DB *gorm.DB, TableName, Field, ConditionField, ConditionValue string, Default float64) float64 {
	var Result []float64
	DB.Table(TableName).Select(Field).Where(ConditionField+" = ?", ConditionValue).Limit(1).Find(&Result)
	if len(Result) > 0 {
		return Result[0]
	}
	return Default
}

func GetFieldBool(DB *gorm.DB, TableName, Field, ConditionField string, ConditionValue interface{}, Default bool) bool {
	var Result map[string]interface{}
	DB.Table(TableName).Select(Field).Where(ConditionField+" = ?", ConditionValue).Limit(1).Scan(&Result)

	value, ok := Result[Field]
	if !ok {
		return Default
	}

	switch v := value.(type) {
	case int8:
		return v == 1
	case string:
		return v == "1"
	default:
		// Handle other data types or return a default value
		return Default
	}
}

func GetFieldFloatQuery(DB *gorm.DB, Query string, Default float64, Condition ...interface{}) float64 {
	var Result float64
	DB.Raw(Query, Condition...).Limit(1).Find(&Result)

	if Result > 0 {
		return Result
	}
	return Default
}

func GetFieldUintQuery(DB *gorm.DB, Query string, Condition ...interface{}) uint64 {
	var Result uint64
	DB.Raw(Query, Condition...).Limit(1).Find(&Result)

	return Result
}

func GetFieldUint(DB *gorm.DB, TableName, Field, ConditionField, ConditionValue string, Default uint64) uint64 {
	var Result uint64
	DB.Table(TableName).Select(Field).Where(ConditionField+" = ?", ConditionValue).Limit(1).Find(&Result)
	if Result > 0 {
		return Result
	}
	return Default
}

func ValidateRequestString(c *gin.Context) (uint, string) {
	var Result uint = global_var.ResponseCode.Successfully
	var StringValue string

	err := c.BindJSON(&StringValue)
	if err != nil {
		Result = global_var.ResponseCode.InvalidDataFormat
		return Result, ""
	} else {
		return Result, StringValue
	}
}

func ValidateRequestUint(c *gin.Context) (uint, uint64) {
	var Result uint = global_var.ResponseCode.Successfully
	var UintValue uint64

	err := c.BindJSON(&UintValue)
	if err != nil {
		Result = global_var.ResponseCode.InvalidDataFormat
		return Result, 0
	} else {
		return Result, UintValue
	}
}

func ValidateRequestBool(Token string, c *gin.Context) (uint, bool) {
	var Result uint = global_var.ResponseCode.Successfully
	var BoolValue bool

	err := c.BindJSON(&BoolValue)
	if err != nil {
		Result = global_var.ResponseCode.InvalidDataFormat
		return Result, false
	} else {
		return Result, BoolValue
	}
}

func CheckData(c *gin.Context, DB *gorm.DB, TableName, FieldName string, FieldCondition interface{}) bool {
	var DataOutput []interface{}
	RowsAffected := DB.Table(TableName).Select(FieldName).Where(FieldName+"=?", FieldCondition).Scan(&DataOutput).RowsAffected
	return RowsAffected > 0
}

func CheckData2(c *gin.Context, DB *gorm.DB, TableName, FieldName1, FieldName2 string, FieldCondition1, FieldCondition2 interface{}) bool {
	var DataOutput []interface{}

	RowsAffected := DB.Table(TableName).Select(FieldName1).Where(FieldName1+"=?", FieldCondition1).Where(FieldName2+"=?", FieldCondition2).Scan(&DataOutput).RowsAffected
	return RowsAffected > 0
}

func CheckData21(c *gin.Context, DB *gorm.DB, TableName, FieldName1, FieldName2 string, FieldCondition1, FieldCondition2 interface{}) bool {
	var DataOutput []interface{}

	RowsAffected := DB.Table(TableName).Select(FieldName1).Where(FieldName1+"=?", FieldCondition1).Or(FieldName2+"=?", FieldCondition2).Scan(&DataOutput).RowsAffected
	return RowsAffected > 0
}

func CheckData22(c *gin.Context, DB *gorm.DB, TableName, FieldName1, FieldName2 string, FieldCondition1, FieldCondition2 interface{}) bool {
	var DataOutput []interface{}

	RowsAffected := DB.Table(TableName).Select(FieldName1).Where(FieldName1+"<>?", FieldCondition1).Where(FieldName2+"=?", FieldCondition2).Scan(&DataOutput).RowsAffected
	return RowsAffected > 0
}

func CheckData31(c *gin.Context, DB *gorm.DB, TableName, FieldName1, FieldName2, FieldName3 string, FieldCondition1, FieldCondition2, FieldCondition3 interface{}) bool {
	var DataOutput []interface{}
	RowsAffected := DB.Table(TableName).Select(FieldName1).Where(FieldName1+"<>?", FieldCondition1).Where(FieldName2+"=?", FieldCondition2).Where(FieldName3+"=?", FieldCondition3).Scan(&DataOutput).RowsAffected
	return RowsAffected > 0
}

func CheckCode(c *gin.Context, DB *gorm.DB, TableName, Code string) bool {
	var DataOutput []string
	RowsAffected := DB.Table(TableName).Select("code").Where("code=?", Code).Scan(&DataOutput).RowsAffected
	return RowsAffected > 0
}

func CheckCodeField(DB *gorm.DB, TableName, FieldSelected string, cond interface{}, args ...interface{}) bool {
	var Data map[string]interface{}
	DB.Table(TableName).Select(FieldSelected).Where(cond, args...).Limit(1).Scan(&Data)
	return Data != nil
}

func CheckNumber(c *gin.Context, DB *gorm.DB, TableName, Number string) bool {
	var DataOutput []string
	RowsAffected := DB.Table(TableName).Select("number").Where("number=?", Number).Scan(&DataOutput).RowsAffected
	return RowsAffected > 0
}

func GetTableFieldId(DB *gorm.DB, TableName string, cond interface{}, args ...interface{}) uint64 {
	var Data uint64
	DB.Table(TableName).Select("id").Where(cond, args...).Limit(1).Scan(&Data)
	return Data
}

// =======================================================================================================================================================================//
// ========================================                MASTER DATA           =========================================================================================//
// =======================================================================================================================================================================//
func FilterGeneralCodeName(c *gin.Context, DB *gorm.DB, DataName string) interface{} {
	// Get Program Configuration
	val, exist := c.Get("pConfig")
	if !exist {
		SendResponse(global_var.ResponseCode.DatabaseError, "pConfig", nil, c)
		return nil
	}
	pConfig := val.(*config.CompanyDataConfiguration)
	//CHS MODULE
	if DataName == "JournalAccountCategory" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitJournalAccountCategory, "code")
	} else if DataName == "JournalAccountType" {
		return GetGeneralCodeName(DB, db_var.TableName.ConstJournalAccountType, "code")
	} else if DataName == "JournalAccountSubGroup" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitJournalAccountSubGroup, "code")
	} else if DataName == "JournalAccountGroup" {
		return GetGeneralCodeName(DB, db_var.TableName.ConstJournalAccountGroup, "code")
	} else if DataName == "JournalAccountSubGroupType" {
		return GetGeneralCodeName(DB, db_var.TableName.ConstJournalAccountSubGroupType, "code")
	} else if DataName == "JournalAccount" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitJournalAccount, "code")
	} else if DataName == "JournalAccountAP" {
		return GetGeneralCodeNameWithParam(DB, db_var.TableName.CfgInitJournalAccount, "type_code", global_var.GlobalJournalAccountType.AccountPayable, "", "", "", "", "code", 1)
	} else if DataName == "JournalAccountAR" {
		return GetGeneralCodeNameWithParam(DB, db_var.TableName.CfgInitJournalAccount, "type_code", global_var.GlobalJournalAccountType.AccountReceivable, "", "", "", "", "code", 1)
	} else if DataName == "JournalAccountIncome" {
		return GetGeneralCodeNameQuery(DB,
			"SELECT"+
				" cfg_init_journal_account.code,"+
				" cfg_init_journal_account.name,"+
				" cfg_init_journal_account_sub_group.name AS SubGroupName "+
				"FROM"+
				" cfg_init_journal_account"+
				" LEFT OUTER JOIN cfg_init_journal_account_sub_group ON (cfg_init_journal_account.sub_group_code = cfg_init_journal_account_sub_group.code)"+
				" WHERE cfg_init_journal_account_sub_group.group_code=?"+
				" OR cfg_init_journal_account_sub_group.group_code=? "+
				"ORDER BY cfg_init_journal_account.code", global_var.GlobalJournalAccountGroup.Income, global_var.GlobalJournalAccountGroup.OtherIncome, "", "", "", 2)
	} else if DataName == "JournalAccountExpense" {
		return GetGeneralCodeNameQuery(DB,
			"SELECT"+
				" cfg_init_journal_account.code,"+
				" cfg_init_journal_account.name,"+
				" cfg_init_journal_account_sub_group.name AS SubGroupName "+
				"FROM"+
				" cfg_init_journal_account"+
				" LEFT OUTER JOIN cfg_init_journal_account_sub_group ON (cfg_init_journal_account.sub_group_code = cfg_init_journal_account_sub_group.code)"+
				" WHERE cfg_init_journal_account_sub_group.group_code=?"+
				" OR cfg_init_journal_account_sub_group.group_code=?"+
				" OR cfg_init_journal_account_sub_group.group_code=? "+
				"ORDER BY cfg_init_journal_account.code", global_var.GlobalJournalAccountGroup.Expense1, global_var.GlobalJournalAccountGroup.Expense2, global_var.GlobalJournalAccountGroup.OtherExpense, "", "", 3)
	} else if DataName == "JournalAccountCosting" {
		return GetGeneralCodeNameQuery(DB,
			"SELECT"+
				" cfg_init_journal_account.code,"+
				" cfg_init_journal_account.name,"+
				" cfg_init_journal_account_sub_group.name AS SubGroupName "+
				"FROM"+
				" cfg_init_journal_account"+
				" LEFT OUTER JOIN cfg_init_journal_account_sub_group ON (cfg_init_journal_account.sub_group_code = cfg_init_journal_account_sub_group.code)"+
				" WHERE cfg_init_journal_account_sub_group.group_code=? "+
				"ORDER BY cfg_init_journal_account.code", global_var.GlobalJournalAccountGroup.Cost, "", "", "", "", 1)
	} else if DataName == "JournalAccountInventory" {
		return GetGeneralCodeNameQuery(DB,
			"SELECT"+
				" cfg_init_journal_account.code,"+
				" cfg_init_journal_account.name,"+
				" cfg_init_journal_account_sub_group.name AS SubGroupName "+
				"FROM"+
				" cfg_init_journal_account"+
				" LEFT OUTER JOIN cfg_init_journal_account_sub_group ON (cfg_init_journal_account.sub_group_code = cfg_init_journal_account_sub_group.code)"+
				" WHERE cfg_init_journal_account.sub_group_code=? "+
				"ORDER BY cfg_init_journal_account.code", pConfig.Dataset.GlobalJournalAccountSubGroup.Inventory, "", "", "", "", 1)
	} else if DataName == "BankAccountType" {
		return GetGeneralCodeName(DB, db_var.TableName.AccConstBankAccountType, "id_sort")
	} else if DataName == "BankAccount" {
		return GetGeneralCodeName(DB, db_var.TableName.AccCfgInitBankAccount, "code")
	} else if DataName == "BankAccountPayment" {
		return GetGeneralCodeNameWithParam(DB, db_var.TableName.AccCfgInitBankAccount, "for_payment", 1, nil, nil, nil, nil, "code", 1)
	} else if DataName == "DepartmentType" {
		return GetGeneralCodeName(DB, db_var.TableName.ConstDepartmentType, "id_sort")
	} else if DataName == "Department" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitDepartment, "id_sort")
	} else if DataName == "SubDepartment" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitSubDepartment, "id_sort")
	} else if DataName == "TaxAndService" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitTaxAndService, "code")
	} else if DataName == "ItemGroup" {
		return GetGeneralCodeName(DB, db_var.TableName.InvCfgInitItemGroup, "id_sort")
	} else if DataName == "SubFolioGroup" {
		return GetGeneralCodeName(DB, db_var.TableName.SubFolioGroup, "code")
	} else if DataName == "AccountGroup" {
		return GetGeneralCodeName(DB, db_var.TableName.ConstAccountGroup, "code")
	} else if DataName == "AccountSubGroup" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitAccountSubGroup, "id_sort")
	} else if DataName == "Account" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitAccount, "code")
	} else if DataName == "AccountEDC" {
		return GetGeneralCodeNameWithParam(DB, db_var.TableName.CfgInitAccount, "sub_group_code", global_var.GlobalAccountSubGroup.CreditDebitCard, "", "", "", "", "code", 1)
	} else if DataName == "AccountForCommision" {
		return GetGeneralCodeNameWithParam(DB, db_var.TableName.CfgInitAccount, "sub_group_code=? AND code<>? AND code<>?", global_var.GlobalAccountSubGroup.AccountPayable, pConfig.Dataset.Configuration[global_var.ConfigurationCategory.GlobalAccount][global_var.ConfigurationName.AccountAPRefundDeposit].(string), pConfig.Dataset.GlobalAccount.CreditCardAdm, "", "", "code", 3)
	} else if DataName == "CompanyType" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitCompanyType, "code")
	} else if DataName == "Company" {
		return GetGeneralCodeName(DB, db_var.TableName.Company, "code")
	} else if DataName == "GuestGroup" {
		return GetGeneralCodeNameCondition(DB, db_var.TableName.GuestGroup, "name", "is_active=1")
	} else if DataName == "BusinessSource" {
		return GetGeneralCodeNameWithParam(DB, db_var.TableName.Company, "is_business_source", 1, "", "", "", "", "code", 1)
	} else if DataName == "Currency" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitCurrency, "code")
	} else if DataName == "Market" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitMarket, "code")
	} else if DataName == "CommissionType" {
		return GetGeneralCodeName(DB, db_var.TableName.ConstCommissionType, "code")
	} else if DataName == "CommissionTypeForPackage" {
		return GetGeneralCodeNameWithParam(DB, db_var.TableName.ConstCommissionType, "code<>? AND code<>? AND code<>? AND code<>?", global_var.CommissionType.PercentFirstNightNettRate, global_var.CommissionType.PercentPerNightNettRate, global_var.CommissionType.PercentOfPriceFullPrice, global_var.CommissionType.PercentOfPriceNettPrice, "", "code", 4)
	} else if DataName == "ChargeFrequency" {
		return GetGeneralCodeName(DB, db_var.TableName.ConstChargeFrequency, "code")
	} else if DataName == "Continent" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitContinent, "code")
	} else if DataName == "RoomType" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitRoomType, "code")
	} else if DataName == "BedType" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitBedType, "code")
	} else if DataName == "RoomView" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitRoomView, "code")
	} else if DataName == "Package" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitPackage, "code")
	} else if DataName == "RoomRateCategory" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitRoomRateCategory, "code")
	} else if DataName == "RoomRateSubCategory" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitRoomRateSubCategory, "code")
	} else if DataName == "DynamicRateType" {
		return GetGeneralCodeName(DB, db_var.TableName.ConstDynamicRateType, "code")
	} else if DataName == "RoomRate" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitRoomRate, "code")
	} else if DataName == "RoomStatus" {
		return GetGeneralCodeName(DB, db_var.TableName.ConstRoomStatus, "code")
	} else if DataName == "Floor" {
		return GetGeneralCodeNameQuery(DB, "Select distinct floor as code, floor as name FROM cfg_init_room ORDER BY floor", nil, nil, nil, nil, nil, 0)
	} else if DataName == "Building" {
		return GetGeneralCodeNameQuery(DB, "Select distinct building as code, building as name FROM cfg_init_room ORDER BY building", nil, nil, nil, nil, nil, 0)
	} else if DataName == "RoomAmenities" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitRoomAmenities, "code")
	} else if DataName == "Room" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitRoom, "number")
	} else if DataName == "RoomBoy" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitRoomBoy, "code")
	} else if DataName == "Owner" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitOwner, "code")
	} else if DataName == "GuestTitle" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitTitle, "code")
	} else if DataName == "GuestType" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitGuestType, "code")
	} else if DataName == "Country" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitCountry, "code")
	} else if DataName == "State" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitState, "code")
	} else if DataName == "Regency" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitRegency, "code")
	} else if DataName == "City" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitCity, "code")
	} else if DataName == "Nationality" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitNationality, "code")
	} else if DataName == "Language" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitLanguage, "code")
	} else if DataName == "IDCardType" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitIdCardType, "code")
	} else if DataName == "PaymentGroup" {
		return GetGeneralCodeName(DB, db_var.TableName.ConstPaymentGroup, "code")
	} else if DataName == "PaymentType" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitPaymentType, "code")
	} else if DataName == "MarketCategory" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitMarketCategory, "code")
	} else if DataName == "Market" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitMarket, "code")
	} else if DataName == "BookingSource" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitBookingSource, "code")
	} else if DataName == "Sales" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitSales, "code")
	} else if DataName == "SalesSalary" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitSalesSalary, "code")
	} else if DataName == "PurposeOf" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitPurposeOf, "code")
	} else if DataName == "CardBank" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitCardBank, "code")
	} else if DataName == "CardType" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitCardType, "code")
	} else if DataName == "LoanItem" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitLoanItem, "code")
	} else if DataName == "PhoneBookType" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitPhoneBookType, "code")
	} else if DataName == "MemberPointType" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitMemberPointType, "code")
	} else if DataName == "VoucherReason" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitVoucherReason, "code")
	} else if DataName == "CompetitorCategory" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitCompetitorCategory, "code")
	} else if DataName == "Competitor" {
		return GetGeneralCodeName(DB, db_var.TableName.Competitor, "code")
	} else if DataName == "MemberType" {
		return GetGeneralCodeName(DB, db_var.TableName.ConstMemberType, "code")
	} else if DataName == "CustomLookupField01" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitCustomLookupField01, "code")
	} else if DataName == "CustomLookupField02" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitCustomLookupField02, "code")
	} else if DataName == "CustomLookupField03" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitCustomLookupField03, "code")
	} else if DataName == "CustomLookupField04" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitCustomLookupField04, "code")
	} else if DataName == "CustomLookupField05" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitCustomLookupField05, "code")
	} else if DataName == "CustomLookupField06" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitCustomLookupField06, "code")
	} else if DataName == "CustomLookupField07" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitCustomLookupField07, "code")
	} else if DataName == "CustomLookupField08" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitCustomLookupField08, "code")
	} else if DataName == "CustomLookupField09" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitCustomLookupField09, "code")
	} else if DataName == "CustomLookupField10" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitCustomLookupField10, "code")
	} else if DataName == "CustomLookupField11" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitCustomLookupField11, "code")
	} else if DataName == "CustomLookupField12" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitCustomLookupField12, "code")
	} else if DataName == "Timezone" {
		var DataOutput []db_var.GeneralCodeNameStruct
		DB.Table(db_var.TableName.CfgInitTimezone).Select("name, cfg_init_timezone.offset/3600 AS code").Scan(&DataOutput)

		return DataOutput
		// return GetGeneralCodeName(DB, db_var.TableName.CfgInitTimezone, "name")
		//POS MODULE
	} else if DataName == "Outlet" {
		DataOutput := make([]map[string]interface{}, 0)
		DB.Table(db_var.TableName.PosCfgInitOutlet).Select("code", "name", "sub_department_code").Scan(&DataOutput)

		return DataOutput
	} else if DataName == "Tenan" {
		return GetGeneralCodeName(DB, db_var.TableName.PosCfgInitTenan, "code")
	} else if DataName == "ProductCategory" {
		return GetGeneralCodeName(DB, db_var.TableName.PosCfgInitProductCategory, "code")
	} else if DataName == "ProductGroup" {
		return GetGeneralCodeName(DB, db_var.TableName.PosCfgInitProductGroup, "code")
	} else if DataName == "Product" {
		DataOutput := make([]map[string]interface{}, 0)
		DB.Table(db_var.TableName.PosCfgInitProduct).Select("pos_cfg_init_product.code", "pos_cfg_init_product.name", "pos_cfg_init_product_group.account_code").
			Joins("LEFT JOIN pos_cfg_init_product_group ON (pos_cfg_init_product.group_code=pos_cfg_init_product_group.code)").
			Scan(&DataOutput)

		return DataOutput
	} else if DataName == "POSMarket" {
		return GetGeneralCodeName(DB, db_var.TableName.PosCfgInitMarket, "code")
	} else if DataName == "SpaRoom" {
		return GetGeneralCodeName(DB, db_var.TableName.PosCfgInitSpaRoom, "number")
	} else if DataName == "Table" {
		return GetGeneralCodeName(DB, db_var.TableName.PosCfgInitTable, "number")
	} else if DataName == "TableType" {
		return GetGeneralCodeName(DB, db_var.TableName.PosCfgInitTableType, "code")
	} else if DataName == "Waitress" {
		return GetGeneralCodeName(DB, db_var.TableName.PosCfgInitWaitress, "code")
	} else if DataName == "Printer" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitPrinter, "code")
		//CAMS MODULE
	} else if DataName == "ShippingAddress" {
		return GetGeneralCodeName(DB, db_var.TableName.AstCfgInitShippingAddress, "code")
	} else if DataName == "UOM" {
		return GetGeneralCodeName(DB, db_var.TableName.InvCfgInitUom, "code")
	} else if DataName == "InventoryStore" {
		return GetGeneralCodeNameQuery(DB,
			"SELECT"+
				" code,"+
				" name "+
				"FROM"+
				" inv_cfg_init_store "+
				" WHERE is_room='0'", "", "", "", "", "", 0)
	} else if DataName == "Store" {
		return GetGeneralCodeName(DB, db_var.TableName.InvCfgInitStore, "code")
	} else if DataName == "ItemCategory" {
		return GetGeneralCodeName(DB, db_var.TableName.InvCfgInitItemCategory, "code")
	} else if DataName == "ItemGroupType" {
		return GetGeneralCodeName(DB, db_var.TableName.InvConstItemGroupType, "code")
	} else if DataName == "Item" {
		return GetGeneralCodeNameQuery(DB,
			"SELECT"+
				" code,"+
				" name "+
				"FROM"+
				" inv_cfg_init_item "+
				" WHERE is_active='1'", "", "", "", "", "", 0)
	} else if DataName == "ReturnStockReason" {
		return GetGeneralCodeName(DB, db_var.TableName.InvCfgInitReturnStockReason, "code")
	} else if DataName == "FAManufacture" {
		return GetGeneralCodeName(DB, db_var.TableName.FaCfgInitManufacture, "code")
	} else if DataName == "FALocationType" {
		return GetGeneralCodeName(DB, db_var.TableName.FaConstLocationType, "code")
	} else if DataName == "FALocation" {
		return GetGeneralCodeName(DB, db_var.TableName.FaCfgInitLocation, "code")
	} else if DataName == "FAItemCategory" {
		return GetGeneralCodeName(DB, db_var.TableName.FaCfgInitItemCategory, "code")
	} else if DataName == "FAItem" {
		return GetGeneralCodeName(DB, db_var.TableName.FaCfgInitItem, "code")
		//Global All Module
	} else if DataName == "UserGroup" {
		return GetGeneralCodeName(DB, db_var.TableName.UserGroupAccess, "code")
	} else if DataName == "User" {
		return GetGeneralCodeName(DB, db_var.TableName.User, "code")
	} else if DataName == "ReservationStatus" {
		return GetGeneralCodeName(DB, db_var.TableName.ConstReservationStatus, "code")
	} else if DataName == "Shift" {
		return GetGeneralCodeNameQuery(DB,
			"SELECT"+
				" working_shift.shift as `code`,"+
				" working_shift.shift as `name` "+
				"FROM"+
				" working_shift "+
				"ORDER BY working_shift.shift", "", "", "", "", "", 0)
	} else if DataName == "CfgInitSubAccount" {
		return GetGeneralCodeName(DB, db_var.TableName.CfgInitAccountSubGroup, "code")
	} else if DataName == "VoucherType" {
		return GetGeneralCodeName(DB, db_var.TableName.ConstVoucherType, "code")
	} else if DataName == "VoucherStatus" {
		return GetGeneralCodeName(DB, db_var.TableName.ConstVoucherStatus, "code")
	} else if DataName == "VoucherStatusApprove" {
		return GetGeneralCodeName(DB, db_var.TableName.ConstVoucherStatusApprove, "code")
	} else if DataName == "VoucherStatusSold" {
		return GetGeneralCodeName(DB, db_var.TableName.ConstVoucherStatusSold, "code")
	} else if DataName == "RoomBuilding" {
		return GetGeneralCodeNameQuery(DB,
			"SELECT"+
				" DISTINCT cfg_init_room.building as `code`,"+
				" cfg_init_room.building as `name` "+
				"FROM"+
				" cfg_init_room "+
				"ORDER BY cfg_init_room.building", "", "", "", "", "", 0)
	} else if DataName == "RoomFloor" {
		return GetGeneralCodeNameQuery(DB,
			"SELECT"+
				" DISTINCT cfg_init_room.floor as `code`,"+
				" cfg_init_room.floor as `name` "+
				"FROM"+
				" cfg_init_room "+
				"ORDER BY cfg_init_room.floor", "", "", "", "", "", 0)
	} else if DataName == "ComplimentType" {
		return GetGeneralCodeName(DB, db_var.TableName.PosConstComplimentType, "code")
	} else if DataName == "SalesSegment" {
		return GetGeneralCodeName(DB, db_var.TableName.SalCfgInitSegment, "code")
	} else if DataName == "SalesStatus" {
		return GetGeneralCodeName(DB, db_var.TableName.SalConstStatus, "code")
	} else if DataName == "SalesSource" {
		return GetGeneralCodeName(DB, db_var.TableName.SalCfgInitSource, "code")
	} else if DataName == "SalesTaskStatus" {
		return GetGeneralCodeName(DB, db_var.TableName.SalConstTaskStatus, "code")
	} else if DataName == "SalesTaskRepeat" {
		return GetGeneralCodeName(DB, db_var.TableName.SalCfgInitTaskRepeat, "code")
	} else if DataName == "SalesProposalStatus" {
		return GetGeneralCodeName(DB, db_var.TableName.SalConstProposalStatus, "code")
	} else if DataName == "ChannelManagerVendor" {
		return GetGeneralCodeName(DB, db_var.TableName.ConstChannelManagerVendor, "code")
	} else if DataName == "DirectBill" {
		return GetGeneralCodeNameQuery(DB,
			"SELECT"+
				" company.code,"+
				" company.name "+
				"FROM"+
				" company WHERE is_direct_bill = 1 "+
				"ORDER BY company.name", "", "", "", "", "", 0)
	} else if DataName == "Action" {
		return GetGeneralCodeNameQuery(DB,
			"SELECT id as `code`, name FROM log_user_action"+
				"ORDER BY name", "", "", "", "", "", 0)
	} else if DataName == "IpAddress" {
		return GetGeneralCodeNameQuery(DB,
			"SELECT DISTINCT ip_address FROM log_user_action"+
				"ORDER BY ip_address", "", "", "", "", "", 0)
	} else if DataName == "ComputerName" {
		return GetGeneralCodeNameQuery(DB,
			"SELECT DISTINCT computer_name FROM log_user_action"+
				"ORDER BY computer_name", "", "", "", "", "", 0)
	} else if DataName == "AccessName" {
		return GetGeneralCodeNameQuery(DB,
			"SELECT DISTINCT access_name FROM log_special_access"+
				"ORDER BY access_name", "", "", "", "", "", 0)
	}
	return nil
}

func FilterGeneralCodeDescription(c *gin.Context, DB *gorm.DB, DataName string) []db_var.GeneralCodeDescriptionStruct {
	if DataName == "RoomUnavailableReason" {
		return GetGeneralCodeDescription(DB, db_var.TableName.CfgInitRoomUnavailableReason, "code")
	}

	return nil
}

func GetMasterDataCodeNameArrayP(c *gin.Context) {
	var DataNameList []string
	Param := c.Query("DataNameList")
	err := json.Unmarshal([]byte(Param), &DataNameList)
	GeneralCodeNameArray := make(map[string]interface{})
	if err != nil {
		SendResponse(global_var.ResponseCode.InvalidDataFormat, "", nil, c)
	} else {
		// Get Program Configuration
		val, exist := c.Get("pConfig")
		if !exist {
			SendResponse(global_var.ResponseCode.DatabaseError, "pConfig", nil, c)
			return
		}
		pConfig := val.(*config.CompanyDataConfiguration)
		DB := pConfig.DB
		for _, DataName := range DataNameList {
			GeneralCodeName := FilterGeneralCodeName(c, DB, DataName)
			GeneralCodeNameArray[DataName] = GeneralCodeName
		}
	}
	SendResponse(global_var.ResponseCode.Successfully, "", GeneralCodeNameArray, c)
}

func GetMasterDataCodeNameP(c *gin.Context) {
	DataName := c.Param("DataName")
	// Get Program Configuration
	val, exist := c.Get("pConfig")
	if !exist {
		SendResponse(global_var.ResponseCode.DatabaseError, "pConfig", nil, c)
		return
	}
	pConfig := val.(*config.CompanyDataConfiguration)
	DB := pConfig.DB
	GeneralCodeName := FilterGeneralCodeName(c, DB, DataName)
	SendResponse(global_var.ResponseCode.Successfully, "", GeneralCodeName, c)
}

func GetMasterDataCodeDescriptionP(c *gin.Context) {
	DataName := c.Param("DataName")
	// Get Program Configuration
	val, exist := c.Get("pConfig")
	if !exist {
		SendResponse(global_var.ResponseCode.DatabaseError, "pConfig", nil, c)
		return
	}
	pConfig := val.(*config.CompanyDataConfiguration)
	DB := pConfig.DB
	GeneralCodeDescription := FilterGeneralCodeDescription(c, DB, DataName)
	SendResponse(global_var.ResponseCode.Successfully, "", GeneralCodeDescription, c)
}

func GetMasterDataTableName(DataName string) string {
	if DataName == "JournalAccount" {
		return db_var.TableName.CfgInitJournalAccount
	} else if DataName == "JournalAccountCategory" {
		return db_var.TableName.CfgInitJournalAccountSubGroup
	} else if DataName == "JournalAccountCategory" {
		return db_var.TableName.CfgInitJournalAccountCategory
	} else if DataName == "Department" {
		return db_var.TableName.CfgInitDepartment
	} else if DataName == "SubDepartment" {
		return db_var.TableName.CfgInitSubDepartment
	} else if DataName == "TaxAndService" {
		return db_var.TableName.CfgInitTaxAndService
	} else if DataName == "AccountSubGroup" {
		return db_var.TableName.CfgInitAccountSubGroup
	} else if DataName == "Account" {
		return db_var.TableName.CfgInitAccount
	} else if DataName == "BankAccount" {
		return db_var.TableName.AccCfgInitBankAccount
	} else if DataName == "CompanyType" {
		return db_var.TableName.CfgInitCompanyType
	} else if DataName == "Company" {
		return db_var.TableName.Company
	} else if DataName == "Currency" {
		return db_var.TableName.CfgInitCurrency
	} else if DataName == "CurrencyNominal" {
		return db_var.TableName.CfgInitCurrencyNominal
	} else if DataName == "RoomType" {
		return db_var.TableName.CfgInitRoomType
	} else if DataName == "BedType" {
		return db_var.TableName.CfgInitBedType
	} else if DataName == "Package" {
		return db_var.TableName.CfgInitPackage
	} else if DataName == "PackageBreakdown" {
		return db_var.TableName.CfgInitPackageBreakdown
	} else if DataName == "PackageBusinessSource" {
		return db_var.TableName.CfgInitPackageBusinessSource
	} else if DataName == "RoomRateCategory" {
		return db_var.TableName.CfgInitRoomRateCategory
	} else if DataName == "RoomRateSubCategory" {
		return db_var.TableName.CfgInitRoomRateSubCategory
	} else if DataName == "RoomRate" {
		return db_var.TableName.CfgInitRoomRate
	} else if DataName == "RoomRateBreakdown" {
		return db_var.TableName.CfgInitRoomRateBreakdown
	} else if DataName == "RoomRateBusinessSource" {
		return db_var.TableName.CfgInitRoomRateBusinessSource
	} else if DataName == "RoomRateDynamic" {
		return db_var.TableName.CfgInitRoomRateDynamic
	} else if DataName == "RoomRateCurrency" {
		return db_var.TableName.CfgInitRoomRateCurrency
	} else if DataName == "RoomView" {
		return db_var.TableName.CfgInitRoomView
	} else if DataName == "RoomAmenities" {
		return db_var.TableName.CfgInitRoomAmenities
	} else if DataName == "Room" {
		return db_var.TableName.CfgInitRoom
	} else if DataName == "RoomBoy" {
		return db_var.TableName.CfgInitRoomBoy
	} else if DataName == "RoomUnavailableReason" {
		return db_var.TableName.CfgInitRoomUnavailableReason
	} else if DataName == "Owner" {
		return db_var.TableName.CfgInitOwner
	} else if DataName == "Title" {
		return db_var.TableName.CfgInitTitle
	} else if DataName == "Continent" {
		return db_var.TableName.CfgInitContinent
	} else if DataName == "Country" {
		return db_var.TableName.CfgInitCountry
	} else if DataName == "State" {
		return db_var.TableName.CfgInitState
	} else if DataName == "Regency" {
		return db_var.TableName.CfgInitRegency
	} else if DataName == "City" {
		return db_var.TableName.CfgInitCity
	} else if DataName == "Nationality" {
		return db_var.TableName.CfgInitNationality
	} else if DataName == "Language" {
		return db_var.TableName.CfgInitLanguage
	} else if DataName == "IdCardType" {
		return db_var.TableName.CfgInitIdCardType
	} else if DataName == "PaymentType" {
		return db_var.TableName.CfgInitPaymentType
	} else if DataName == "MarketCategory" {
		return db_var.TableName.CfgInitMarketCategory
	} else if DataName == "Market" {
		return db_var.TableName.CfgInitMarket
	} else if DataName == "BookingSource" {
		return db_var.TableName.CfgInitBookingSource
	} else if DataName == "PurposeOf" {
		return db_var.TableName.CfgInitPurposeOf
	} else if DataName == "CardBank" {
		return db_var.TableName.CfgInitCardBank
	} else if DataName == "CardType" {
		return db_var.TableName.CfgInitCardType
	} else if DataName == "LoanItem" {
		return db_var.TableName.CfgInitLoanItem
	} else if DataName == "CreditCardCharge" {
		return db_var.TableName.CfgInitCreditCardCharge
	} else if DataName == "PhoneBookType" {
		return db_var.TableName.CfgInitPhoneBookType
	} else if DataName == "MemberOutletDiscount" {
		return db_var.TableName.PosCfgInitMemberOutletDiscount
	} else if DataName == "MemberOutletDiscountDetail" {
		return db_var.TableName.PosCfgInitMemberOutletDiscountDetail
	} else if DataName == "MemberPointType" {
		return db_var.TableName.CfgInitMemberPointType
	} else if DataName == "VoucherReason" {
		return db_var.TableName.CfgInitVoucherReason
	} else if DataName == "CompetitorCategory" {
		return db_var.TableName.CfgInitCompetitorCategory
	} else if DataName == "Competitor" {
		return db_var.TableName.Competitor
	} else if DataName == "Sales" {
		return db_var.TableName.CfgInitSales
	} else if DataName == "SalesSalary" {
		return db_var.TableName.CfgInitSalesSalary
	} else if DataName == "SalesSegment" {
		return db_var.TableName.SalCfgInitSegment
	} else if DataName == "SalesSource" {
		return db_var.TableName.SalCfgInitSource
	} else if DataName == "SalesTaskAction" {
		return db_var.TableName.SalCfgInitTaskAction
	} else if DataName == "SalesTaskRepeat" {
		return db_var.TableName.SalCfgInitTaskRepeat
	} else if DataName == "CustomLookupField01" {
		return db_var.TableName.CfgInitCustomLookupField01
	} else if DataName == "CustomLookupField02" {
		return db_var.TableName.CfgInitCustomLookupField02
	} else if DataName == "CustomLookupField03" {
		return db_var.TableName.CfgInitCustomLookupField03
	} else if DataName == "CustomLookupField04" {
		return db_var.TableName.CfgInitCustomLookupField04
	} else if DataName == "CustomLookupField05" {
		return db_var.TableName.CfgInitCustomLookupField05
	} else if DataName == "CustomLookupField06" {
		return db_var.TableName.CfgInitCustomLookupField06
	} else if DataName == "CustomLookupField07" {
		return db_var.TableName.CfgInitCustomLookupField07
	} else if DataName == "CustomLookupField08" {
		return db_var.TableName.CfgInitCustomLookupField08
	} else if DataName == "CustomLookupField09" {
		return db_var.TableName.CfgInitCustomLookupField09
	} else if DataName == "CustomLookupField10" {
		return db_var.TableName.CfgInitCustomLookupField10
	} else if DataName == "CustomLookupField11" {
		return db_var.TableName.CfgInitCustomLookupField11
	} else if DataName == "CustomLookupField12" {
		return db_var.TableName.CfgInitCustomLookupField12
	} else if DataName == "GuestType" {
		return db_var.TableName.CfgInitGuestType
	} else if DataName == "Timezone" {
		return db_var.TableName.CfgInitTimezone
		//POS MODULE
	} else if DataName == "Outlet" {
		return db_var.TableName.PosCfgInitOutlet
	} else if DataName == "Tenan" {
		return db_var.TableName.PosCfgInitTenan
	} else if DataName == "ProductCategory" {
		return db_var.TableName.PosCfgInitProductCategory
	} else if DataName == "ProductGroup" {
		return db_var.TableName.PosCfgInitProductGroup
	} else if DataName == "Product" {
		return db_var.TableName.PosCfgInitProduct
	} else if DataName == "POSMarket" {
		return db_var.TableName.PosCfgInitMarket
	} else if DataName == "POSPaymentGroup" {
		return db_var.TableName.PosCfgInitPaymentGroup
	} else if DataName == "SpaRoom" {
		return db_var.TableName.PosCfgInitSpaRoom
	} else if DataName == "Table" {
		return db_var.TableName.PosCfgInitTable
	} else if DataName == "TableType" {
		return db_var.TableName.PosCfgInitTableType
	} else if DataName == "Waitress" {
		return db_var.TableName.PosCfgInitWaitress
	} else if DataName == "Printer" {
		return db_var.TableName.CfgInitPrinter
	} else if DataName == "DiscountLimit" {
		return db_var.TableName.PosCfgInitDiscountLimit
		//CAMS MODULE
	} else if DataName == "ShippingAddress" {
		return db_var.TableName.AstCfgInitShippingAddress
	} else if DataName == "UOM" {
		return db_var.TableName.InvCfgInitUom
	} else if DataName == "Store" {
		return db_var.TableName.InvCfgInitStore
	} else if DataName == "ItemCategory" {
		return db_var.TableName.InvCfgInitItemCategory
	} else if DataName == "ItemCategoryOtherCOGS" {
		return db_var.TableName.InvCfgInitItemCategoryOtherCogs
	} else if DataName == "ItemCategoryOtherCOGS2" {
		return db_var.TableName.InvCfgInitItemCategoryOtherCogs2
	} else if DataName == "ItemCategoryOtherExpense" {
		return db_var.TableName.InvCfgInitItemCategoryOtherExpense
	} else if DataName == "Item" {
		return db_var.TableName.InvCfgInitItem
	} else if DataName == "ItemUOM" {
		return db_var.TableName.InvCfgInitItem
	} else if DataName == "ItemGroup" {
		return db_var.TableName.InvCfgInitItemGroup
	} else if DataName == "ReturnStockReason" {
		return db_var.TableName.InvCfgInitReturnStockReason
	} else if DataName == "MarketList" {
		return db_var.TableName.InvCfgInitMarketList
	} else if DataName == "FAManufacture" {
		return db_var.TableName.FaCfgInitManufacture
	} else if DataName == "FALocation" {
		return db_var.TableName.FaCfgInitLocation
	} else if DataName == "FAItemCategory" {
		return db_var.TableName.FaCfgInitItemCategory
	} else if DataName == "FAItem" {
		return db_var.TableName.FaCfgInitItem
		//Global All Module
	} else if DataName == "UserGroup" {
		return db_var.TableName.UserGroupAccess
	}
	return ""
}

func GetMasterDataValidationSearchField1(Index int64) string {
	var SearchField string
	switch Index {
	case 0:
		SearchField = "code"
	case 1:
		SearchField = "name"
	case 2:
		SearchField = "created_by"
	case 3:
		SearchField = "updated_by"
	}
	SearchField = SearchField + " LIKE ?"
	return SearchField
}

func GetMasterDataValidationSearchField2(Index int64) string {
	var SearchField string
	switch Index {
	case 0:
		SearchField = "code"
	case 1:
		SearchField = "description"
	case 2:
		SearchField = "created_by"
	case 3:
		SearchField = "updated_by"
	}
	SearchField = SearchField + " LIKE ?"
	return SearchField
}

func GetMasterDataValidationSearchField91(Index int64, TableName, CustomField string) string {
	var SearchField string
	switch Index {
	case 0:
		SearchField = TableName + ".code"
	case 1:
		SearchField = TableName + ".name"
	case 2:
		SearchField = CustomField
	case 3:
		SearchField = TableName + ".created_by"
	case 4:
		SearchField = TableName + ".updated_by"
	}
	SearchField = SearchField + " LIKE ?"
	return SearchField
}

func GetMasterDataValidationSearchField92(Index int64, TableName, CustomField1, CustomField2 string) string {
	var SearchField string
	switch Index {
	case 0:
		SearchField = TableName + ".code"
	case 1:
		SearchField = TableName + ".name"
	case 2:
		SearchField = CustomField1
	case 3:
		SearchField = CustomField2
	case 4:
		SearchField = TableName + ".created_by"
	case 5:
		SearchField = TableName + ".updated_by"
	}
	SearchField = SearchField + " LIKE ?"
	return SearchField
}

func GetMasterDataValidationSearchField93(Index int64, TableName, CustomField1, CustomField2, CustomField3 string) string {
	var SearchField string
	switch Index {
	case 0:
		SearchField = TableName + ".code"
	case 1:
		SearchField = TableName + ".name"
	case 2:
		SearchField = CustomField1
	case 3:
		SearchField = CustomField2
	case 4:
		SearchField = CustomField3
	case 5:
		SearchField = TableName + ".created_by"
	case 6:
		SearchField = TableName + ".updated_by"
	}
	SearchField = SearchField + " LIKE ?"
	return SearchField
}

func GetMasterDataValidationSearchField94(Index int64, TableName, CustomField1, CustomField2, CustomField3, CustomField4 string) string {
	var SearchField string
	switch Index {
	case 0:
		SearchField = TableName + ".code"
	case 1:
		SearchField = TableName + ".name"
	case 2:
		SearchField = CustomField1
	case 3:
		SearchField = CustomField2
	case 4:
		SearchField = CustomField3
	case 5:
		SearchField = CustomField4
	case 6:
		SearchField = TableName + ".created_by"
	case 7:
		SearchField = TableName + ".updated_by"
	}
	SearchField = SearchField + " LIKE ?"
	return SearchField
}

func GetMasterDataListValidation(DataName string, Index int64, Text string, Option1 int64, String1, String2 string) (string, string, string) {
	var Query, QueryCondition string
	var SearchField string
	TableName := GetMasterDataTableName(DataName)
	Query = ""
	QueryCondition = ""
	if DataName == "JournalAccount" {
		var SearchField string
		switch Index {
		case 0:
			SearchField = TableName + ".code"
		case 1:
			SearchField = TableName + ".name"
		case 2:
			SearchField = "CONCAT(cfg_init_journal_account.sub_group_code, ' - ', cfg_init_journal_account_sub_group.name)"
		case 3:
			SearchField = "CONCAT(cfg_init_journal_account_sub_group.group_code, ' - ', const_journal_account_group.name)"
		case 4:
			SearchField = "const_journal_account_type.name"
		case 5:
			SearchField = "inv_cfg_init_item_group.name"
		case 6:
			SearchField = "cfg_init_journal_account_category.name"
		case 7:
			SearchField = TableName + ".description"
		case 8:
			SearchField = TableName + ".created_by"
		case 9:
			SearchField = TableName + ".updated_by"
		}

		if Text != "" {
			QueryCondition = " WHERE " + SearchField + " LIKE ?"
		}
		Query =
			"SELECT" +
				" cfg_init_journal_account.*," +
				" CONCAT(cfg_init_journal_account.sub_group_code, ' - ', cfg_init_journal_account_sub_group.name) AS JournalAccountSubGroup," +
				" CONCAT(cfg_init_journal_account_sub_group.group_code, ' - ', const_journal_account_group.name) AS JournalAccountGroup," +
				" const_journal_account_type.name AS JournalAccountType," +
				" inv_cfg_init_item_group.name AS ItemGroup," +
				" cfg_init_journal_account_category.name AS Category " +
				"FROM" +
				" cfg_init_journal_account" +
				" LEFT OUTER JOIN cfg_init_journal_account_sub_group ON (cfg_init_journal_account.sub_group_code = cfg_init_journal_account_sub_group.code)" +
				" LEFT OUTER JOIN const_journal_account_group ON (cfg_init_journal_account_sub_group.group_code = const_journal_account_group.code)" +
				" LEFT OUTER JOIN const_journal_account_type ON (cfg_init_journal_account.type_code = const_journal_account_type.code)" +
				" LEFT OUTER JOIN inv_cfg_init_item_group ON (cfg_init_journal_account.item_group_code = inv_cfg_init_item_group.code)" +
				" LEFT OUTER JOIN cfg_init_journal_account_category ON (cfg_init_journal_account.category_code = cfg_init_journal_account_category.code)" +
				QueryCondition + " " +
				"ORDER BY cfg_init_journal_account.code;"
	} else if DataName == "JournalAccountSubGroup" {
		SearchField = GetMasterDataValidationSearchField92(Index, TableName, "const_journal_account_group.name", "const_journal_account_sub_group_type.name")

		if Text != "" {
			QueryCondition = " WHERE " + SearchField
		}
		Query = "SELECT" +
			" " + TableName + ".*," +
			" const_journal_account_group.name AS GroupName," +
			" const_journal_account_sub_group_type.name AS JournalAccountSubGroupTypeName " +
			"FROM" +
			" " + TableName +
			" LEFT OUTER JOIN const_journal_account_group ON (" + TableName + ".group_code = const_journal_account_group.code)" +
			" LEFT OUTER JOIN const_journal_account_sub_group_type ON (" + TableName + ".type_code = const_journal_account_sub_group_type.code)" +
			QueryCondition + " " +
			"ORDER BY " + TableName + ".code;"
	} else if DataName == "Department" {
		SearchField = GetMasterDataValidationSearchField91(Index, TableName, "const_department_type.name")

		if Text != "" {
			QueryCondition = " WHERE " + SearchField
		}
		Query = "SELECT" +
			" " + TableName + ".*," +
			" const_department_type.name AS TypeName " +
			"FROM" +
			" " + TableName +
			" LEFT OUTER JOIN const_department_type ON (" + TableName + ".type_code = const_department_type.code)" +
			QueryCondition + " " +
			"ORDER BY " + TableName + ".id_sort;"
	} else if DataName == "SubDepartment" {
		SearchField = GetMasterDataValidationSearchField91(Index, TableName, "cfg_init_department.name")

		if Text != "" {
			QueryCondition = " WHERE " + SearchField
		}
		Query = "SELECT" +
			" " + TableName + ".*," +
			" cfg_init_department.name AS DepartmentName " +
			"FROM" +
			" " + TableName +
			" LEFT OUTER JOIN cfg_init_department ON (" + TableName + ".department_code = cfg_init_department.code)" +
			QueryCondition + " " +
			"ORDER BY " + TableName + ".id_sort;"
	} else if DataName == "Account" {
		var SearchField string
		switch Index {
		case 0:
			SearchField = TableName + ".code"
		case 1:
			SearchField = TableName + ".name"
		case 2:
			SearchField = "const_transaction_type.name"
		case 3:
			SearchField = "cfg_init_account_sub_group.name"
		case 4:
			SearchField = "CONCAT(" + TableName + ".journal_account_code, ' - ', cfg_init_journal_account.name)"
		case 5:
			SearchField = "cfg_init_tax_and_service.name"
		case 6:
			SearchField = "inv_cfg_init_item_group.name"
		case 7:
			SearchField = TableName + ".created_by"
		case 8:
			SearchField = TableName + ".updated_by"
		}

		if Text != "" {
			QueryCondition = " WHERE " + SearchField + " LIKE ?"
		}
		Query = "SELECT" +
			" " + TableName + ".*," +
			" const_transaction_type.name AS TransactionTypeName," +
			" cfg_init_account_sub_group.name AS AccountSubGroupName," +
			" CONCAT(" + TableName + ".journal_account_code, ' - ', cfg_init_journal_account.name) as JournalAccount," +
			" cfg_init_tax_and_service.name AS TaxAndServiceName," +
			" inv_cfg_init_item_group.name AS ItemGroupName " +
			"FROM" +
			" " + TableName +
			" LEFT OUTER JOIN const_transaction_type ON (" + TableName + ".type_code = const_transaction_type.code)" +
			" LEFT OUTER JOIN cfg_init_account_sub_group ON (" + TableName + ".sub_group_code = cfg_init_account_sub_group.code)" +
			" LEFT OUTER JOIN cfg_init_journal_account ON (" + TableName + ".journal_account_code = cfg_init_journal_account.code)" +
			" LEFT OUTER JOIN cfg_init_tax_and_service ON (" + TableName + ".tax_and_service_code = cfg_init_tax_and_service.code)" +
			" LEFT OUTER JOIN inv_cfg_init_item_group ON (cfg_init_account.item_group_code = inv_cfg_init_item_group.code)" +
			QueryCondition + " " +
			"ORDER BY " + TableName + ".code;"
	} else if DataName == "AccountSubGroup" {
		SearchField = GetMasterDataValidationSearchField91(Index, TableName, "const_account_group.name")

		if Text != "" {
			QueryCondition = " WHERE " + SearchField
		}
		Query = "SELECT" +
			" " + TableName + ".*," +
			" const_account_group.name AS GroupName " +
			"FROM" +
			" " + TableName +
			" LEFT OUTER JOIN const_account_group ON (" + TableName + ".group_code = const_account_group.code)" +
			QueryCondition + " " +
			"ORDER BY " + TableName + ".id_sort;"
	} else if DataName == "BankAccount" {
		var SearchField string
		switch Index {
		case 0:
			SearchField = TableName + ".code"
		case 1:
			SearchField = TableName + ".name"
		case 2:
			SearchField = "CONCAT(" + TableName + ".journal_account_code, ' - ', cfg_init_journal_account.name)"
		case 3:
			SearchField = "acc_const_bank_account_type.name"
		case 4:
			SearchField = TableName + ".bank_name"
		case 5:
			SearchField = TableName + ".bank_account_number"
		case 6:
			SearchField = TableName + ".bank_address"
		case 7:
			SearchField = TableName + ".created_by"
		case 8:
			SearchField = TableName + ".updated_by"
		}

		if Text != "" {
			QueryCondition = " WHERE " + SearchField + " LIKE ?"
		}
		Query = "SELECT" +
			" " + TableName + ".*," +
			" CONCAT(" + TableName + ".journal_account_code, ' - ', cfg_init_journal_account.name) AS JournalAccount," +
			" acc_const_bank_account_type.name AS BankAccountTypeName " +
			"FROM" +
			" " + TableName +
			" LEFT OUTER JOIN cfg_init_journal_account ON (" + TableName + ".journal_account_code = cfg_init_journal_account.code)" +
			" LEFT OUTER JOIN acc_const_bank_account_type ON (" + TableName + ".type_code = acc_const_bank_account_type.code)" +
			QueryCondition + " " +
			"ORDER BY " + TableName + ".code;"
	} else if DataName == "CompanyType" {
		switch Index {
		case 0:
			SearchField = TableName + ".code"
		case 1:
			SearchField = TableName + ".name"
		case 2:
			SearchField = "CONCAT(" + TableName + ".journal_account_code_ap, ' - ', cfg_init_journal_account.name)"
		case 3:
			SearchField = "CONCAT(" + TableName + ".journal_account_code_ar, ' - ', cfg_init_journal_account1.name)"
		case 4:
			SearchField = TableName + ".created_by"
		case 5:
			SearchField = TableName + ".updated_by"
		}

		if Text != "" {
			QueryCondition = " WHERE " + SearchField + " LIKE ?"
		}
		Query = "SELECT" +
			" " + TableName + ".*," +
			" CONCAT(" + TableName + ".journal_account_code_ap, ' - ', cfg_init_journal_account.name) AS APAccount," +
			" CONCAT(" + TableName + ".journal_account_code_ar, ' - ', cfg_init_journal_account1.name) AS ARAccount " +
			"FROM" +
			" " + TableName +
			" LEFT OUTER JOIN cfg_init_journal_account ON (" + TableName + ".journal_account_code_ap = cfg_init_journal_account.code)" +
			" LEFT OUTER JOIN cfg_init_journal_account cfg_init_journal_account1 ON (" + TableName + ".journal_account_code_ar = cfg_init_journal_account1.code)" +
			QueryCondition + " " +
			"ORDER BY " + TableName + ".id_sort;"
	} else if DataName == "Company" {
		var SearchField string
		switch Index {
		case 0:
			SearchField = TableName + ".code"
		case 1:
			SearchField = TableName + ".name"
		case 2:
			SearchField = "cfg_init_company_type.name"
		case 3:
			SearchField = "cfg_init_sales.name"
		case 4:
			SearchField = TableName + ".contact_person"
		case 5:
			SearchField = "TRIM(CONCAT(company.street, ' ', company.city, ' ', IFNULL(cfg_init_state.name, ''), ' ', IFNULL(cfg_init_country.name, ''), ' ', company.postal_code))"
		case 6:
			SearchField = "IF(company.phone1 = '', company.phone2, CONCAT(company.phone1, ', ', company.phone2))"
		case 7:
			SearchField = TableName + ".fax"
		case 8:
			SearchField = TableName + ".email"
		case 9:
			SearchField = TableName + ".website"
		case 10:
			SearchField = TableName + ".created_by"
		case 11:
			SearchField = TableName + ".updated_by"
		}

		switch Option1 {
		case 1:
			QueryCondition = TableName + ".is_direct_bill=1"
		case 2:
			QueryCondition = TableName + ".is_business_source=1"
		case 3:
			QueryCondition = TableName + ".is_direct_bill=0 AND " + TableName + ".is_business_source=0"
		}

		if Text != "" {
			if QueryCondition == "" {
				QueryCondition = " WHERE " + SearchField + " LIKE ?"
			} else {
				QueryCondition = " WHERE " + QueryCondition + " AND " + SearchField + " LIKE ?"
			}
		} else {
			if QueryCondition != "" {
				QueryCondition = " WHERE " + QueryCondition
			}
		}

		Query = "SELECT" +
			" " + TableName + ".code," +
			" " + TableName + ".name," +
			" " + TableName + ".contact_person," +
			" TRIM(CONCAT(" + TableName + ".street, ' ', " + TableName + ".city, ' ', IFNULL(cfg_init_state.name, ''), ' ', IFNULL(cfg_init_country.name, ''), ' ', " + TableName + ".postal_code)) AS Address," +
			" IF(" + TableName + ".phone1 = '', " + TableName + ".phone2, CONCAT(" + TableName + ".phone1, ', ', " + TableName + ".phone2)) AS Phone," +
			" " + TableName + ".fax," +
			" " + TableName + ".email," +
			" " + TableName + ".website," +
			" " + TableName + ".ap_limit," +
			" " + TableName + ".ar_limit," +
			" " + TableName + ".is_direct_bill," +
			" " + TableName + ".is_business_source," +
			" " + TableName + ".invoice_due," +
			" " + TableName + ".created_at," +
			" " + TableName + ".created_by," +
			" " + TableName + ".updated_at," +
			" " + TableName + ".updated_by," +
			" " + TableName + ".id," +
			" cfg_init_sales.name AS SalesName," +
			" cfg_init_company_type.name AS CompanyTypeName " +
			"FROM" +
			" " + TableName +
			" LEFT OUTER JOIN cfg_init_sales ON (" + TableName + ".sales_code = cfg_init_sales.code)" +
			" LEFT OUTER JOIN cfg_init_company_type ON (" + TableName + ".type_code = cfg_init_company_type.code)" +
			" LEFT OUTER JOIN cfg_init_country ON (" + TableName + ".country_code = cfg_init_country.code)" +
			" LEFT OUTER JOIN cfg_init_state ON (" + TableName + ".state_code = cfg_init_state.code)" +
			QueryCondition + " " +
			"ORDER BY " + TableName + ".code;"
	} else if DataName == "CurrencyNominal" {
		var SearchField string
		switch Index {
		case 0:
			SearchField = "currency_sign"
		case 1:
			SearchField = "created_by"
		case 2:
			SearchField = "updated_by"
		}

		if Text != "" {
			QueryCondition = " WHERE " + SearchField + " LIKE ?"
		}
		Query = "SELECT * FROM " + TableName +
			QueryCondition + " " +
			"ORDER BY " + TableName + ".id_sort;"
	} else if DataName == "Package" {
		var SearchField string
		switch Index {
		case 0:
			SearchField = TableName + ".code"
		case 1:
			SearchField = TableName + ".name"
		case 2:
			SearchField = "pos_cfg_init_outlet.name"
		case 3:
			SearchField = "pos_cfg_init_product.name"
		case 4:
			SearchField = "cfg_init_sub_department.name"
		case 5:
			SearchField = "CONCAT(" + TableName + ".account_code, ' - ', cfg_init_account.name)"
		case 6:
			SearchField = "cfg_init_tax_and_service.name"
		case 7:
			SearchField = "const_charge_frequency.name"
		case 8:
			SearchField = TableName + ".created_by"
		case 9:
			SearchField = TableName + ".updated_by"
		}

		if Text != "" {
			QueryCondition = " WHERE " + SearchField + " LIKE ?"
		}
		Query = "SELECT" +
			" " + TableName + ".*," +
			" pos_cfg_init_outlet.name AS OutletName," +
			" pos_cfg_init_product.name AS ProductName," +
			" cfg_init_sub_department.name AS DepartmentName," +
			" CONCAT(" + TableName + ".account_code, ' - ', cfg_init_account.name) AS Account," +
			" cfg_init_tax_and_service.name AS TaxAndServiceName," +
			" const_charge_frequency.name AS FrequencyName " +
			"FROM" +
			" " + TableName +
			" LEFT OUTER JOIN pos_cfg_init_outlet ON (" + TableName + ".outlet_code = pos_cfg_init_outlet.code)" +
			" LEFT OUTER JOIN pos_cfg_init_product ON (" + TableName + ".product_code = pos_cfg_init_product.code)" +
			" LEFT OUTER JOIN cfg_init_sub_department ON (" + TableName + ".sub_department_code = cfg_init_sub_department.code)" +
			" LEFT OUTER JOIN cfg_init_account ON (" + TableName + ".account_code = cfg_init_account.code)" +
			" LEFT OUTER JOIN cfg_init_tax_and_service ON (" + TableName + ".tax_and_service_code = cfg_init_tax_and_service.code)" +
			" LEFT OUTER JOIN const_charge_frequency ON (" + TableName + ".charge_frequency_code = const_charge_frequency.code)" +
			QueryCondition + " " +
			"ORDER BY " + TableName + ".code;"
	} else if DataName == "PackageBreakdown" {
		Query = "SELECT" +
			" " + TableName + ".*," +
			" pos_cfg_init_outlet.name AS OutletName," +
			" pos_cfg_init_product.name AS ProductName," +
			" cfg_init_sub_department.name AS SubDepartmentName," +
			" CONCAT(" + TableName + ".account_code, ' - ', cfg_init_account.name) AS Account," +
			" company.name AS CompanyName," +
			" cfg_init_tax_and_service.name AS TaxAndServiceName," +
			" const_charge_frequency.name AS FrequencyName " +
			"FROM" +
			" " + TableName +
			" LEFT OUTER JOIN pos_cfg_init_outlet ON (" + TableName + ".outlet_code = pos_cfg_init_outlet.code)" +
			" LEFT OUTER JOIN pos_cfg_init_product ON (" + TableName + ".product_code = pos_cfg_init_product.code)" +
			" LEFT OUTER JOIN cfg_init_sub_department ON (" + TableName + ".sub_department_code = cfg_init_sub_department.code)" +
			" LEFT OUTER JOIN cfg_init_account ON (" + TableName + ".account_code = cfg_init_account.code)" +
			" LEFT OUTER JOIN company ON (cfg_init_package_breakdown.company_code = company.code)" +
			" LEFT OUTER JOIN cfg_init_tax_and_service ON (" + TableName + ".tax_and_service_code = cfg_init_tax_and_service.code)" +
			" LEFT OUTER JOIN const_charge_frequency ON (" + TableName + ".charge_frequency_code = const_charge_frequency.code)" +
			" WHERE " + TableName + ".package_code=?;"
	} else if DataName == "PackageBusinessSource" {
		Query = "SELECT" +
			" " + TableName + ".*," +
			" CONCAT(cfg_init_package_business_source.account_code, ' - ', cfg_init_account.name) AS Account," +
			" company.name AS CompanyName," +
			" const_commission_type.name AS CommissionTypeName " +
			"FROM" +
			" " + TableName +
			" LEFT OUTER JOIN cfg_init_account ON (cfg_init_package_business_source.account_code = cfg_init_account.code)" +
			" LEFT OUTER JOIN company ON (" + TableName + ".company_code = company.code)" +
			" LEFT OUTER JOIN const_commission_type ON (" + TableName + ".commission_type_code = const_commission_type.code)" +
			" WHERE " + TableName + ".package_code=?;"
	} else if DataName == "RoomRateSubCategory" {
		SearchField = GetMasterDataValidationSearchField91(Index, TableName, "cfg_init_room_rate_category.name")

		if Text != "" {
			QueryCondition = " WHERE " + SearchField
		}
		Query = "SELECT" +
			" " + TableName + ".*," +
			" cfg_init_room_rate_category.name AS CategoryName " +
			"FROM" +
			" " + TableName +
			" LEFT OUTER JOIN cfg_init_room_rate_category ON (" + TableName + ".category_code = cfg_init_room_rate_category.code)" +
			QueryCondition + " " +
			"ORDER BY " + TableName + ".id_sort;"
	} else if DataName == "RoomRate" {
		var SearchField string
		switch Index {
		case 0:
			SearchField = TableName + ".code"
		case 1:
			SearchField = TableName + ".name"
		case 2:
			SearchField = TableName + ".room_type_code"
		case 3:
			SearchField = "cfg_init_room_rate_sub_category.name"
		case 4:
			SearchField = "company.name"
		case 5:
			SearchField = "cfg_init_market.name"
		case 6:
			SearchField = "const_dynamic_rate_type.name"
		case 7:
			SearchField = TableName + ".cm_inv_code"
		case 8:
			SearchField = "cfg_init_tax_and_service.name"
		case 9:
			SearchField = "const_charge_frequency.name"
		case 10:
			SearchField = TableName + ".notes"
		case 11:
			SearchField = "created_by"
		case 12:
			SearchField = "updated_by"
		}

		if Text != "" {
			QueryCondition = " WHERE " + SearchField + " LIKE ?"
		} else {
			QueryCondition = " WHERE "
		}

		switch Option1 {
		case 0:
			if Text == "" {
				QueryCondition = ""
			}
		case 1:
			if Text != "" {
				QueryCondition += " AND "
			}
			QueryCondition = QueryCondition + TableName + ".is_active=1"
		case 2:
			if Text != "" {
				QueryCondition += " AND "
			}
			QueryCondition = QueryCondition + TableName + ".is_active=0"
		}

		Query = "SELECT" +
			" " + TableName + ".*," +
			" cfg_init_room_rate_sub_category.name AS SubCategoryName," +
			" company.name AS CompanyName," +
			" cfg_init_market.name MarketName," +
			" const_dynamic_rate_type.name AS DynamicRateType," +
			" cfg_init_tax_and_service.name TaxAndServiceName," +
			" const_charge_frequency.name AS ChargeFrequencyName " +
			"FROM" +
			" " + TableName +
			" LEFT OUTER JOIN cfg_init_room_rate_sub_category ON (cfg_init_room_rate.sub_category_code = cfg_init_room_rate_sub_category.code)" +
			" LEFT OUTER JOIN company ON (cfg_init_room_rate.company_code = company.code)" +
			" LEFT OUTER JOIN cfg_init_market ON (cfg_init_room_rate.market_code = cfg_init_market.code)" +
			" LEFT OUTER JOIN const_dynamic_rate_type ON (cfg_init_room_rate.dynamic_rate_type_code = const_dynamic_rate_type.code)" +
			" LEFT OUTER JOIN cfg_init_tax_and_service ON (cfg_init_room_rate.tax_and_service_code = cfg_init_tax_and_service.code)" +
			" LEFT OUTER JOIN const_charge_frequency ON (" + TableName + ".charge_frequency_code = const_charge_frequency.code)" +
			QueryCondition + " " +
			"ORDER BY " + TableName + ".id_sort;"
	} else if DataName == "RoomRateBreakdown" {
		Query = "SELECT" +
			" " + TableName + ".*," +
			" pos_cfg_init_outlet.name AS OutletName," +
			" pos_cfg_init_product.name AS ProductName," +
			" cfg_init_sub_department.name AS SubDepartmentName," +
			" CONCAT(" + TableName + ".account_code, ' - ', cfg_init_account.name) AS Account," +
			" company.name AS CompanyName," +
			" cfg_init_tax_and_service.name AS TaxAndServiceName," +
			" const_charge_frequency.name AS ChargeFrequencyName " +
			"FROM" +
			" " + TableName +
			" LEFT OUTER JOIN pos_cfg_init_outlet ON (" + TableName + ".outlet_code = pos_cfg_init_outlet.code)" +
			" LEFT OUTER JOIN pos_cfg_init_product ON (" + TableName + ".product_code = pos_cfg_init_product.code)" +
			" LEFT OUTER JOIN cfg_init_sub_department ON (" + TableName + ".sub_department_code = cfg_init_sub_department.code)" +
			" LEFT OUTER JOIN cfg_init_account ON (" + TableName + ".account_code = cfg_init_account.code)" +
			" LEFT OUTER JOIN company ON (cfg_init_room_rate_breakdown.company_code = company.code)" +
			" LEFT OUTER JOIN cfg_init_tax_and_service ON (" + TableName + ".tax_and_service_code = cfg_init_tax_and_service.code)" +
			" LEFT OUTER JOIN const_charge_frequency ON (" + TableName + ".charge_frequency_code = const_charge_frequency.code)" +
			" WHERE " + TableName + ".room_rate_code=?"
	} else if DataName == "RoomRateBusinessSource" {
		Query = "SELECT" +
			" " + TableName + ".*," +
			" company.name AS CompanyName," +
			" const_commission_type.name AS CommissionTypeName " +
			"FROM" +
			" " + TableName +
			" LEFT OUTER JOIN company ON (" + TableName + ".company_code = company.code)" +
			" LEFT OUTER JOIN const_commission_type ON (" + TableName + ".commission_type_code = const_commission_type.code)" +
			" WHERE " + TableName + ".room_rate_code=?"
	} else if DataName == "RoomRateDynamic" {
		Query = "SELECT" +
			" " + TableName + ".room_rate_code," +
			" " + TableName + ".name," +
			" " + TableName + ".occ_from," +
			" " + TableName + ".occ_to," +
			" " + TableName + ".is_percent," +
			" " + TableName + ".is_increase," +
			" " + TableName + ".amount," +
			" " + TableName + ".created_at," +
			" " + TableName + ".created_by," +
			" " + TableName + ".updated_at," +
			" " + TableName + ".updated_by," +
			" " + TableName + ".id " +
			"FROM" +
			" " + TableName +
			" WHERE " + TableName + ".room_rate_code=? " +
			"ORDER BY " + TableName + ".occ_from"
	} else if DataName == "RoomRateCurrency" {
		Query = "SELECT" +
			" " + TableName + ".*," +
			" cfg_init_currency.name AS CurrencyName " +
			"FROM" +
			" " + TableName +
			" LEFT OUTER JOIN cfg_init_currency ON (" + TableName + ".currency_code = cfg_init_currency.code)" +
			" WHERE " + TableName + ".room_rate_code=?"
	} else if DataName == "Room" {
		var SearchField string
		switch Index {
		case 0:
			SearchField = TableName + ".number"
		case 1:
			SearchField = TableName + ".lock_number"
		case 2:
			SearchField = TableName + ".name"
		case 3:
			SearchField = "CONCAT(cfg_init_room_type.name, ' ', cfg_init_bed_type.name)"
		case 4:
			SearchField = TableName + ".building"
		case 5:
			SearchField = TableName + ".floor"
		case 6:
			SearchField = "cfg_init_room_view.name"
		case 7:
			SearchField = TableName + ".description"
		case 8:
			SearchField = TableName + ".phone_number"
		case 9:
			SearchField = "cfg_init_owner.name"
		case 10:
			SearchField = "IF(" + TableName + ".is_smoking='0', 'No', 'Yes')"
		case 11:
			SearchField = TableName + ".created_by"
		case 12:
			SearchField = TableName + ".updated_by"
		}

		if Text != "" {
			QueryCondition = " WHERE " + SearchField + " LIKE ?"
		}
		Query = "SELECT" +
			" " + TableName + ".number," +
			" " + TableName + ".name," +
			" " + TableName + ".lock_number," +
			" IF(" + TableName + ".is_smoking='0', 'No', 'Yes') AS Smoking," +
			" " + TableName + ".building," +
			" " + TableName + ".floor," +
			" " + TableName + ".max_adult," +
			" " + TableName + ".description," +
			" " + TableName + ".phone_number," +
			" " + TableName + ".tv_quantity," +
			" " + TableName + ".start_date," +
			" " + TableName + ".id_sort," +
			" " + TableName + ".image," +
			" " + TableName + ".status_code," +
			" " + TableName + ".created_at," +
			" " + TableName + ".created_by," +
			" " + TableName + ".updated_at," +
			" " + TableName + ".updated_by," +
			" " + TableName + ".id," +
			" CONCAT(cfg_init_room_type.name, ' ', cfg_init_bed_type.name) AS RoomType," +
			" cfg_init_room_view.name AS ViewName," +
			" cfg_init_owner.name AS OwnerName " +
			"FROM" +
			" " + TableName +
			" LEFT OUTER JOIN cfg_init_room_type ON (" + TableName + ".room_type_code = cfg_init_room_type.code)" +
			" LEFT OUTER JOIN cfg_init_bed_type ON (" + TableName + ".bed_type_code = cfg_init_bed_type.code)" +
			" LEFT OUTER JOIN cfg_init_room_view ON (" + TableName + ".view_code = cfg_init_room_view.code)" +
			" LEFT OUTER JOIN cfg_init_owner ON (" + TableName + ".owner_code = cfg_init_owner.code)" +
			" LEFT OUTER JOIN const_room_status ON (cfg_init_room.status_code = const_room_status.code)" +
			QueryCondition + " " +
			"ORDER BY " + TableName + ".id_sort, " + TableName + ".number;"
	} else if DataName == "Country" {
		SearchField = GetMasterDataValidationSearchField92(Index, TableName, "cfg_init_continent.name", "iso_code")

		if Text != "" {
			QueryCondition = " WHERE " + SearchField
		}
		Query = "SELECT" +
			" " + TableName + ".*," +
			" cfg_init_continent.name AS ContinentName " +
			"FROM" +
			" " + TableName +
			" LEFT OUTER JOIN cfg_init_continent ON (" + TableName + ".continent_code = cfg_init_continent.code)" +
			QueryCondition + " " +
			"ORDER BY " + TableName + ".code;"
	} else if DataName == "State" {
		SearchField = GetMasterDataValidationSearchField91(Index, TableName, "cfg_init_country.name")

		if Text != "" {
			QueryCondition = " WHERE " + SearchField
		}
		Query = "SELECT" +
			" " + TableName + ".*," +
			" cfg_init_country.name AS CountryName " +
			"FROM" +
			" " + TableName +
			" LEFT OUTER JOIN cfg_init_country ON (" + TableName + ".country_code = cfg_init_country.code)" +
			QueryCondition + " " +
			"ORDER BY " + TableName + ".code;"
	} else if DataName == "City" {
		SearchField = GetMasterDataValidationSearchField92(Index, TableName, "cfg_init_state.name", "cfg_init_regency.name")

		if Text != "" {
			QueryCondition = " WHERE " + SearchField
		}
		Query = "SELECT" +
			" " + TableName + ".*," +
			" cfg_init_state.name AS CityName," +
			" cfg_init_regency.name AS RegencyName " +
			"FROM" +
			" " + TableName +
			" LEFT OUTER JOIN cfg_init_state ON (" + TableName + ".state_code = cfg_init_state.code)" +
			" LEFT OUTER JOIN cfg_init_regency ON (" + TableName + ".regency_code = cfg_init_regency.code)" +
			QueryCondition + " " +
			"ORDER BY " + TableName + ".code;"
	} else if DataName == "PaymentType" {
		SearchField = GetMasterDataValidationSearchField91(Index, TableName, "const_payment_group.name")

		if Text != "" {
			QueryCondition = " WHERE " + SearchField
		}
		Query = "SELECT" +
			" " + TableName + ".*," +
			" const_payment_group.name AS PaymentGroupName " +
			"FROM" +
			" " + TableName +
			" LEFT OUTER JOIN const_payment_group ON (" + TableName + ".group_code = const_payment_group.code)" +
			QueryCondition + " " +
			"ORDER BY " + TableName + ".code;"
	} else if DataName == "Market" {
		SearchField = GetMasterDataValidationSearchField91(Index, TableName, "cfg_init_market_category.name")

		if Text != "" {
			QueryCondition = " WHERE " + SearchField
		}
		Query = "SELECT" +
			" " + TableName + ".*," +
			" cfg_init_market_category.name AS CategoryName " +
			"FROM" +
			" " + TableName +
			" LEFT OUTER JOIN cfg_init_market_category ON (" + TableName + ".category_code = cfg_init_market_category.code)" +
			QueryCondition + " " +
			"ORDER BY " + TableName + ".code;"
	} else if DataName == "CreditCardCharge" {
		switch Index {
		case 0:
			SearchField = "cfg_init_account.name"
		case 1:
			SearchField = TableName + ".created_by"
		case 2:
			SearchField = TableName + ".updated_by"
		}

		if Text != "" {
			QueryCondition = " WHERE " + SearchField + " LIKE ?"
		}
		Query = "SELECT" +
			" " + TableName + ".*," +
			" cfg_init_account.name AS AccountName " +
			"FROM" +
			" " + TableName +
			" LEFT OUTER JOIN cfg_init_account ON (" + TableName + ".account_code = cfg_init_account.code)" +
			QueryCondition + " " +
			"ORDER BY " + TableName + ".account_code;"
	} else if DataName == "MemberPointType" {
		SearchField = GetMasterDataValidationSearchField92(Index, TableName, "const_member_type.name", "room_type_code")

		if Text != "" {
			QueryCondition = " WHERE " + SearchField
		}

		Query = "SELECT" +
			" " + TableName + ".*," +
			" const_member_type.name AS TypeName " +
			"FROM" +
			" " + TableName +
			" LEFT OUTER JOIN const_member_type ON (cfg_init_member_point_type.member_type_code = const_member_type.code)" +
			QueryCondition +
			"ORDER BY " + TableName + ".code;"
		//POS MODULE
	} else if DataName == "Outlet" {
		switch Index {
		case 0:
			SearchField = TableName + ".code"
		case 1:
			SearchField = TableName + ".name"
		case 2:
			SearchField = TableName + ".check_prefix"
		case 3:
			SearchField = "cfg_init_sub_department.name"
		case 4:
			SearchField = "inv_cfg_init_store.name"
		case 5:
			SearchField = "cfg_init_tax_and_service.name"
		case 6:
			SearchField = TableName + ".created_by"
		case 7:
			SearchField = TableName + ".updated_by"
		}

		if Text != "" {
			QueryCondition = " WHERE " + SearchField + " LIKE ?"
		}
		Query =
			"SELECT" +
				" " + TableName + ".*," +
				" cfg_init_sub_department.name AS DepartmentName," +
				" inv_cfg_init_store.name AS StoreName," +
				" cfg_init_tax_and_service.name AS TaxAndServiceName " +
				"FROM" +
				" " + TableName +
				" LEFT OUTER JOIN cfg_init_sub_department ON (" + TableName + ".sub_department_code = cfg_init_sub_department.code)" +
				" LEFT OUTER JOIN inv_cfg_init_store ON (pos_cfg_init_outlet.store_code = inv_cfg_init_store.code)" +
				" LEFT OUTER JOIN cfg_init_tax_and_service ON (" + TableName + ".tax_and_service_code = cfg_init_tax_and_service.code)" +
				QueryCondition + " " +
				"ORDER BY " + TableName + ".id_sort;"
	} else if DataName == "MemberOutletDiscount" {
		switch Index {
		case 0:
			SearchField = TableName + ".code"
		case 1:
			SearchField = "pos_cfg_init_outlet.name"
		case 2:
			SearchField = TableName + ".name"
		case 3:
			SearchField = TableName + ".created_by"
		case 4:
			SearchField = TableName + ".updated_by"
		}

		if Text != "" {
			QueryCondition = " WHERE " + SearchField + " LIKE ?"
		}
		Query = "SELECT" +
			" " + TableName + ".*," +
			" pos_cfg_init_outlet.name AS AccountName " +
			"FROM" +
			" " + TableName +
			" LEFT OUTER JOIN pos_cfg_init_outlet ON (" + TableName + ".outlet_code = pos_cfg_init_outlet.code)" +
			QueryCondition + " " +
			"ORDER BY " + TableName + ".code;"
	} else if DataName == "MemberOutletDiscountDetail" {
		QueryCondition = " WHERE member_outlet_discount_code=?"

		Query = "SELECT" +
			" " + TableName + ".*," +
			" pos_cfg_init_outlet.name AS OutletName," +
			" pos_cfg_init_product.name AS ProductName " +
			"FROM" +
			" " + TableName +
			" LEFT OUTER JOIN pos_cfg_init_outlet ON (" + TableName + ".outlet_code = pos_cfg_init_outlet.code)" +
			" LEFT OUTER JOIN pos_cfg_init_product ON (" + TableName + ".product_code = pos_cfg_init_product.code)" +
			QueryCondition + ";"
	} else if DataName == "Competitor" {
		SearchField = GetMasterDataValidationSearchField91(Index, TableName, "cfg_init_competitor_category.name")

		if Text != "" {
			QueryCondition = " WHERE " + SearchField
		}

		Query = "SELECT" +
			" " + TableName + ".*," +
			" cfg_init_competitor_category.name AS CategoryName " +
			"FROM" +
			" " + TableName +
			" LEFT OUTER JOIN cfg_init_competitor_category ON (" + TableName + ".category_code = cfg_init_competitor_category.code)" +
			QueryCondition + " " +
			"ORDER BY " + TableName + ".code;"
	} else if DataName == "Tenan" {
		switch Index {
		case 0:
			SearchField = TableName + ".code"
		case 1:
			SearchField = TableName + ".name"
		case 2:
			SearchField = "cfg_init_owner.name"
		case 3:
			SearchField = TableName + ".street"
		case 4:
			SearchField = TableName + ".city"
		case 5:
			SearchField = TableName + ".phone1"
		case 6:
			SearchField = TableName + ".phone2"
		case 7:
			SearchField = TableName + ".fax"
		case 8:
			SearchField = TableName + ".email"
		case 9:
			SearchField = TableName + ".created_by"
		case 10:
			SearchField = TableName + ".updated_by"
		}

		if Text != "" {
			QueryCondition = " WHERE " + SearchField + " LIKE ?"
		}

		Query = "SELECT" +
			" " + TableName + ".*," +
			" cfg_init_owner.name AS OwnerName " +
			"FROM" +
			" " + TableName +
			" LEFT OUTER JOIN cfg_init_owner ON (" + TableName + ".owner_code = cfg_init_owner.code)" +
			QueryCondition + " " +
			"ORDER BY " + TableName + ".code;"
	} else if DataName == "ProductGroup" {
		switch Index {
		case 0:
			SearchField = TableName + ".code"
		case 1:
			SearchField = TableName + ".name"
		case 2:
			SearchField = "CONCAT(" + TableName + ".account_code, ' - ', cfg_init_account.name)"
		case 3:
			SearchField = "inv_cfg_init_item_group.name"
		case 4:
			SearchField = TableName + ".created_by"
		case 5:
			SearchField = TableName + ".updated_by"
		}

		if Text != "" {
			QueryCondition = " WHERE " + SearchField + " LIKE ?"
		}

		Query = "SELECT" +
			" " + TableName + ".*," +
			" inv_cfg_init_item_group.name AS ItemGroupName," +
			" CONCAT(" + TableName + ".account_code, ' - ', cfg_init_account.name) AS Account " +
			"FROM" +
			" " + TableName +
			" LEFT OUTER JOIN cfg_init_account ON (" + TableName + ".account_code = cfg_init_account.code)" +
			" LEFT OUTER JOIN inv_cfg_init_item_group ON (pos_cfg_init_product_group.item_group_code = inv_cfg_init_item_group.code)" +
			QueryCondition + " " +
			"ORDER BY " + TableName + ".code;"
	} else if DataName == "Product" {
		switch Index {
		case 0:
			SearchField = TableName + ".code"
		case 1:
			SearchField = TableName + ".barcode"
		case 2:
			SearchField = TableName + ".name"
		case 3:
			SearchField = TableName + ".description"
		case 4:
			SearchField = "pos_cfg_init_product_category.name"
		case 5:
			SearchField = "pos_cfg_init_product_group.name"
		case 6:
			SearchField = TableName + ".outlet_code"
		case 7:
			SearchField = "cfg_init_tax_and_service.name"
		case 8:
			SearchField = "pos_cfg_init_tenan.name"
		case 9:
			SearchField = "cfg_init_package.name"
		case 10:
			SearchField = "cfg_init_printer.name"
		case 11:
			SearchField = "cfg_init_printer2.name"
		case 12:
			SearchField = TableName + ".created_by"
		case 13:
			SearchField = TableName + ".updated_by"
		}

		if Text != "" {
			QueryCondition = " WHERE " + SearchField + " LIKE ?"
		}

		// fmt.Println(String1)
		if String1 != "" {
			QueryCondition = QueryCondition + " AND " + TableName + ".outlet_code LIKE '%" + String1 + "%'"
		}

		// fmt.Println(Option1)
		switch Option1 {
		case 1:
			QueryCondition = QueryCondition + " AND (" + TableName + ".is_active=1 AND " + TableName + ".is_sold=0)"
		case 2:
			QueryCondition = QueryCondition + " AND " + TableName + ".is_active=0"
		case 3:
			QueryCondition = QueryCondition + " AND " + TableName + ".is_sold=1"
		}

		if QueryCondition != "" {
			if QueryCondition[:4] == " AND" {
				QueryCondition = " WHERE" + QueryCondition[4:len(QueryCondition)-0]
			}
		}

		Query = "SELECT" +
			" " + TableName + ".*," +
			" pos_cfg_init_product_category.name AS CategoryName," +
			" pos_cfg_init_product_group.name AS GroupName," +
			" pos_cfg_init_tenan.name AS TenanName," +
			" cfg_init_package.name AS PackageName," +
			" cfg_init_printer.name AS PrinterName," +
			" cfg_init_printer2.name AS PrinterName2," +
			" cfg_init_tax_and_service.name AS TaxAndServiceName " +
			"FROM" +
			" " + TableName +
			" LEFT OUTER JOIN pos_cfg_init_product_category ON (" + TableName + ".category_code = pos_cfg_init_product_category.code)" +
			" LEFT OUTER JOIN pos_cfg_init_product_group ON (" + TableName + ".group_code = pos_cfg_init_product_group.code)" +
			" LEFT OUTER JOIN pos_cfg_init_tenan ON (" + TableName + ".tenan_code = pos_cfg_init_tenan.code)" +
			" LEFT OUTER JOIN cfg_init_package ON (pos_cfg_init_product.package_code = cfg_init_package.code)" +
			" LEFT OUTER JOIN cfg_init_printer ON (" + TableName + ".printer_code = cfg_init_printer.code)" +
			" LEFT OUTER JOIN cfg_init_printer cfg_init_printer2 ON (" + TableName + ".printer_code2 = cfg_init_printer2.code)" +
			" LEFT OUTER JOIN cfg_init_tax_and_service ON (pos_cfg_init_product.tax_and_service_code = cfg_init_tax_and_service.code)" +
			QueryCondition + " " +
			"ORDER BY " + TableName + ".code;"
	} else if DataName == "POSPaymentGroup" {
		switch Index {
		case 0:
			SearchField = "account_code"
		case 1:
			SearchField = "name"
		case 2:
			SearchField = "created_by"
		case 3:
			SearchField = "updated_by"
		}
		SearchField = SearchField + " LIKE ?"
	} else if DataName == "SpaRoom" {
		switch Index {
		case 0:
			SearchField = TableName + ".number"
		case 1:
			SearchField = "pos_cfg_init_outlet.name"
		case 2:
			SearchField = TableName + ".description"
		case 3:
			SearchField = TableName + ".created_by"
		case 4:
			SearchField = TableName + ".updated_by"
		}

		if Text != "" {
			QueryCondition = " WHERE " + SearchField + " LIKE ?"
		}
		Query = "SELECT" +
			" pos_cfg_init_spa_room.*," +
			" pos_cfg_init_outlet.name AS OutletName " +
			"FROM" +
			" pos_cfg_init_spa_room" +
			" LEFT OUTER JOIN pos_cfg_init_outlet ON (pos_cfg_init_spa_room.outlet_code = pos_cfg_init_outlet.code)" +
			QueryCondition + " " +
			"ORDER BY " + TableName + ".number;"
	} else if DataName == "Table" {
		switch Index {
		case 0:
			SearchField = TableName + ".number"
		case 1:
			SearchField = "pos_cfg_init_outlet.name"
		case 2:
			SearchField = "pos_cfg_init_table_type.name"
		case 3:
			SearchField = TableName + ".status_code"
		case 4:
			SearchField = TableName + ".created_by"
		case 5:
			SearchField = TableName + ".updated_by"
		}

		if Text != "" {
			QueryCondition = " WHERE " + SearchField + " LIKE ?"
		}
		Query = "SELECT" +
			" " + TableName + ".*," +
			" pos_cfg_init_outlet.name AS OutletName," +
			" pos_cfg_init_table_type.name AS TableTypeName " +
			"FROM" +
			" " + TableName +
			" LEFT OUTER JOIN pos_cfg_init_outlet ON (" + TableName + ".outlet_code = pos_cfg_init_outlet.code)" +
			" LEFT OUTER JOIN pos_cfg_init_table_type ON (" + TableName + ".type_code = pos_cfg_init_table_type.code)" +
			QueryCondition + " " +
			"ORDER BY " + TableName + ".number;"
	} else if DataName == "DiscountLimit" {
		switch Index {
		case 0:
			SearchField = "pos_cfg_init_outlet.name"
		case 1:
			SearchField = TableName + ".user_group_code"
		case 2:
			SearchField = TableName + ".created_by"
		case 3:
			SearchField = TableName + ".updated_by"
		}

		if Text != "" {
			QueryCondition = " WHERE " + SearchField + " LIKE ?"
		}
		Query = "SELECT " +
			TableName + ".*," +
			"pos_cfg_init_outlet.name AS OutletName " +
			" FROM " +
			TableName +
			" LEFT OUTER JOIN pos_cfg_init_outlet ON (" + TableName + ".outlet_code = pos_cfg_init_outlet.code) " +
			QueryCondition + " " +
			"ORDER BY " + TableName + ".outlet_code;"
		//CAMS MODULE
	} else if DataName == "ShippingAddress" {
		switch Index {
		case 0:
			SearchField = TableName + ".code"
		case 1:
			SearchField = TableName + ".name"
		case 2:
			SearchField = TableName + ".contact_person"
		case 3:
			SearchField = "TRIM(CONCAT(" + TableName + ".street, ' ', " + TableName + ".city, ' ', IFNULL(cfg_init_state.name, ''), ' ', IFNULL(cfg_init_country.name, ''), ' ', " + TableName + ".postal_code))"
		case 4:
			SearchField = "(CASE WHEN " + TableName + ".phone1 = '' THEN " + TableName + ".phone2 } else { CONCAT(" + TableName + ".phone1, ', ', " + TableName + ".phone2) END)"
		case 5:
			SearchField = TableName + ".fax"
		case 6:
			SearchField = TableName + ".email"
		case 7:
			SearchField = TableName + ".website"
		case 8:
			SearchField = TableName + ".created_by"
		case 9:
			SearchField = TableName + ".updated_by"
		}

		if Text != "" {
			QueryCondition = " WHERE " + SearchField + " LIKE ?"
		}
		Query = "SELECT" +
			" " + TableName + ".*," +
			" TRIM(CONCAT(" + TableName + ".street, ' ', " + TableName + ".city, ' ', IFNULL(cfg_init_state.name, ''), ' ', IFNULL(cfg_init_country.name, ''), ' ', " + TableName + ".postal_code)) AS Address," +
			" (CASE WHEN " + TableName + ".phone1 = '' THEN " + TableName + ".phone2 ELSE CONCAT(" + TableName + ".phone1, ', ', " + TableName + ".phone2) END) AS Phone " +
			"FROM" +
			" " + TableName +
			" LEFT OUTER JOIN cfg_init_country ON (" + TableName + ".country_code = cfg_init_country.code)" +
			" LEFT OUTER JOIN cfg_init_state ON (" + TableName + ".state_code = cfg_init_state.code)" +
			QueryCondition + " " +
			"ORDER BY " + TableName + ".code;"
	} else if DataName == "ItemCategory" {
		switch Index {
		case 0:
			SearchField = TableName + ".code"
		case 1:
			SearchField = TableName + ".name"
		case 2:
			SearchField = "inv_cfg_init_item_group.name"
		case 3:
			SearchField = "CONCAT(" + TableName + ".journal_account_code, ' - ',  cfg_init_journal_account.name)"
		case 4:
			SearchField = "CONCAT(" + TableName + ".journal_account_code_cogs, ' - ',  cfg_init_journal_account1.name)"
		case 5:
			SearchField = "CONCAT(" + TableName + ".journal_account_code_expense, ' - ',  cfg_init_journal_account2.name)"
		case 6:
			SearchField = "CONCAT(" + TableName + ".journal_account_code_sell, ' - ',  cfg_init_journal_account3.name)"
		case 7:
			SearchField = "CONCAT(" + TableName + ".journal_account_code_adjustment, ' - ',  cfg_init_journal_account4.name)"
		case 8:
			SearchField = "CONCAT(" + TableName + ".journal_account_code_spoil, ' - ',  cfg_init_journal_account5.name)"
		case 9:
			SearchField = "CONCAT(" + TableName + ".journal_account_code_cogs2, ' - ',  cfg_init_journal_account6.name)"
		case 10:
			SearchField = TableName + ".created_by"
		case 11:
			SearchField = TableName + ".updated_by"
		}

		if Text != "" {
			QueryCondition = " WHERE " + SearchField + " LIKE ?"
		}
		Query = "SELECT" +
			" " + TableName + ".*," +
			" inv_cfg_init_item_group.name AS ItemGroupName," +
			" CONCAT(" + TableName + ".journal_account_code, ' - ',  cfg_init_journal_account.name) AS AccountInventory," +
			" CONCAT(" + TableName + ".journal_account_code_cogs, ' - ',  cfg_init_journal_account1.name) AS AccountCOGS," +
			" CONCAT(" + TableName + ".journal_account_code_expense, ' - ',  cfg_init_journal_account2.name) AS AccountExpense," +
			" CONCAT(" + TableName + ".journal_account_code_sell, ' - ',  cfg_init_journal_account3.name) AS AccountSell," +
			" CONCAT(" + TableName + ".journal_account_code_adjustment, ' - ',  cfg_init_journal_account4.name) AS AccountAdjustment," +
			" CONCAT(" + TableName + ".journal_account_code_spoil, ' - ',  cfg_init_journal_account5.name) AS AccountSpoil," +
			" CONCAT(" + TableName + ".journal_account_code_cogs2, ' - ',  cfg_init_journal_account6.name) AS AccountCOGS2 " +
			"FROM" +
			" " + TableName +
			" LEFT OUTER JOIN inv_cfg_init_item_group ON (" + TableName + ".group_code = inv_cfg_init_item_group.code)" +
			" LEFT OUTER JOIN cfg_init_journal_account ON (" + TableName + ".journal_account_code = cfg_init_journal_account.code)" +
			" LEFT OUTER JOIN cfg_init_journal_account cfg_init_journal_account1 ON (" + TableName + ".journal_account_code_cogs = cfg_init_journal_account1.code)" +
			" LEFT OUTER JOIN cfg_init_journal_account cfg_init_journal_account2 ON (" + TableName + ".journal_account_code_expense = cfg_init_journal_account2.code)" +
			" LEFT OUTER JOIN cfg_init_journal_account cfg_init_journal_account3 ON (" + TableName + ".journal_account_code_sell = cfg_init_journal_account3.code)" +
			" LEFT OUTER JOIN cfg_init_journal_account cfg_init_journal_account4 ON (" + TableName + ".journal_account_code_adjustment = cfg_init_journal_account4.code)" +
			" LEFT OUTER JOIN cfg_init_journal_account cfg_init_journal_account5 ON (" + TableName + ".journal_account_code_spoil = cfg_init_journal_account5.code)" +
			" LEFT OUTER JOIN cfg_init_journal_account cfg_init_journal_account6 ON (" + TableName + ".journal_account_code_cogs2 = cfg_init_journal_account6.code)" +
			QueryCondition + " " +
			"ORDER BY " + TableName + ".code;"
	} else if DataName == "ItemCategoryOtherCOGS" {
		Query = "SELECT" +
			" " + TableName + ".*," +
			" cfg_init_sub_department.name," +
			" CONCAT(" + TableName + ".journal_account_code, ' - ', cfg_init_journal_account.name) AS Account " +
			"FROM" +
			" " + TableName +
			" LEFT OUTER JOIN cfg_init_sub_department ON (" + TableName + ".sub_department_code = cfg_init_sub_department.code)" +
			" LEFT OUTER JOIN cfg_init_journal_account ON (" + TableName + ".journal_account_code = cfg_init_journal_account.code)" +
			" WHERE " + TableName + ".category_code=?"
	} else if DataName == "ItemCategoryOtherCOGS2" {
		Query = "SELECT" +
			" " + TableName + ".*," +
			" cfg_init_sub_department.name," +
			" CONCAT(" + TableName + ".journal_account_code, ' - ', cfg_init_journal_account.name) AS Account " +
			"FROM" +
			" " + TableName +
			" LEFT OUTER JOIN cfg_init_sub_department ON (" + TableName + ".sub_department_code = cfg_init_sub_department.code)" +
			" LEFT OUTER JOIN cfg_init_journal_account ON (" + TableName + ".journal_account_code = cfg_init_journal_account.code)" +
			" WHERE " + TableName + ".category_code=?"
	} else if DataName == "ItemCategoryOtherExpense" {
		Query = "SELECT" +
			" " + TableName + ".*," +
			" cfg_init_sub_department.name," +
			" CONCAT(" + TableName + ".journal_account_code, ' - ', cfg_init_journal_account.name) AS Account " +
			"FROM" +
			" " + TableName +
			" LEFT OUTER JOIN cfg_init_sub_department ON (" + TableName + ".sub_department_code = cfg_init_sub_department.code)" +
			" LEFT OUTER JOIN cfg_init_journal_account ON (" + TableName + ".journal_account_code = cfg_init_journal_account.code)" +
			" WHERE " + TableName + ".category_code=?"
	} else if DataName == "Item" {
		switch Index {
		case 0:
			SearchField = TableName + ".code"
		case 1:
			SearchField = TableName + ".name"
		case 2:
			SearchField = "inv_cfg_init_item_category.name"
		case 3:
			SearchField = "inv_cfg_init_uom.name"
		case 4:
			SearchField = TableName + ".remark"
		case 5:
			SearchField = TableName + ".created_by"
		case 6:
			SearchField = TableName + ".updated_by"
		}

		switch Option1 {
		case 0:
			QueryCondition = ""
		case 1:
			QueryCondition = " WHERE " + TableName + ".is_active=1 "
		case 2:
			QueryCondition = " WHERE " + TableName + ".is_active=0 "
		}

		if Text != "" {
			if QueryCondition != "" {
				QueryCondition += " AND " + SearchField + " LIKE ?"
			} else {
				QueryCondition = " WHERE " + SearchField + " LIKE ?"
			}
		}
		Query = "SELECT" +
			" " + TableName + ".*," +
			" inv_cfg_init_item_category.name AS CategoryName," +
			" inv_cfg_init_uom.name AS UOMName " +
			"FROM" +
			" " + TableName +
			" LEFT OUTER JOIN inv_cfg_init_item_category ON (" + TableName + ".category_code = inv_cfg_init_item_category.code)" +
			" LEFT OUTER JOIN inv_cfg_init_uom ON (" + TableName + ".uom_code = inv_cfg_init_uom.code) " +

			QueryCondition + " " +
			"ORDER BY " + TableName + ".code;"
	} else if DataName == "ItemUOM" {
		Query = "SELECT" +
			" " + TableName + ".*," +
			" inv_cfg_init_uom.name " +
			"FROM" +
			" " + TableName +
			" LEFT OUTER JOIN inv_cfg_init_uom ON (" + TableName + ".uom_code = inv_cfg_init_uom.code)" +
			" WHERE " + TableName + ".item_code=?"
	} else if DataName == "MarketList" {
		switch Index {
		case 0:
			SearchField = "company.name"
		case 1:
			SearchField = "inv_cfg_init_item.name"
		case 2:
			SearchField = "inv_cfg_init_uom.name"
		case 3:
			SearchField = TableName + ".created_by"
		case 4:
			SearchField = TableName + ".updated_by"
		}

		if Text != "" {
			QueryCondition = " WHERE " + SearchField + " LIKE ?"
		}
		Query = "SELECT" +
			" " + TableName + ".*," +
			" company.name AS SupplierName," +
			" inv_cfg_init_item.name AS ItemName," +
			" inv_cfg_init_uom.name AS UOMName " +
			"FROM" +
			" " + TableName +
			" LEFT OUTER JOIN company ON (" + TableName + ".company_code = company.code)" +
			" LEFT OUTER JOIN inv_cfg_init_item ON (" + TableName + ".item_code = inv_cfg_init_item.code)" +
			" LEFT OUTER JOIN inv_cfg_init_uom ON (" + TableName + ".uom_code = inv_cfg_init_uom.code)" +
			QueryCondition + " " +
			"ORDER BY company.name, inv_cfg_init_item.name;"
	} else if DataName == "FALocation" {
		SearchField = GetMasterDataValidationSearchField91(Index, TableName, "fa_const_location_type.name")

		if Text != "" {
			QueryCondition = " WHERE " + SearchField
		}

		Query = "SELECT" +
			" " + TableName + ".*," +
			" fa_const_location_type.name AS FALocationTypeName " +
			"FROM" +
			" " + TableName +
			" LEFT OUTER JOIN fa_const_location_type ON (" + TableName + ".location_type_code = fa_const_location_type.code)" +
			QueryCondition + " " +
			"ORDER BY " + TableName + ".code;"
	} else if DataName == "FAItemCategory" {
		switch Index {
		case 0:
			SearchField = TableName + ".code LIKE ?"
		case 1:
			SearchField = TableName + ".name LIKE ?"
		case 2:
			SearchField = "CONCAT(" + TableName + ".journal_account_code, ' - ',  cfg_init_journal_account.name) LIKE ?"
		case 3:
			SearchField = "CONCAT(" + TableName + ".journal_account_code_cogs, ' - ',  cfg_init_journal_account1.name) LIKE ?"
		case 4:
			SearchField = "CONCAT(" + TableName + ".journal_account_code_expense, ' - ',  cfg_init_journal_account2.name) LIKE ?"
		case 5:
			SearchField = "CONCAT(" + TableName + ".journal_account_code_sell, ' - ',  cfg_init_journal_account3.name) LIKE ?"
		case 6:
			SearchField = "CONCAT(" + TableName + ".journal_account_code_depreciation, ' - ',  cfg_init_journal_account4.name) LIKE ?"
		case 7:
			SearchField = "CONCAT(" + TableName + ".journal_account_code_spoil, ' - ',  cfg_init_journal_account5.name) LIKE ?"
		case 8:
			SearchField = TableName + ".created_by LIKE ?"
		case 9:
			SearchField = TableName + ".updated_by LIKE ?"
		}

		if Text != "" {
			QueryCondition = " WHERE " + SearchField
		}
		Query = "SELECT" +
			" " + TableName + ".*," +
			" CONCAT(" + TableName + ".journal_account_code, ' - ',  cfg_init_journal_account.name) AS AccountInventory," +
			" CONCAT(" + TableName + ".journal_account_code_cogs, ' - ',  cfg_init_journal_account1.name) AS AccountCOGS," +
			" CONCAT(" + TableName + ".journal_account_code_expense, ' - ',  cfg_init_journal_account2.name) AS AccountExpense," +
			" CONCAT(" + TableName + ".journal_account_code_sell, ' - ',  cfg_init_journal_account3.name) AS AccountSell," +
			" CONCAT(" + TableName + ".journal_account_code_depreciation, ' - ',  cfg_init_journal_account4.name) AS AccountDepreciation," +
			" CONCAT(" + TableName + ".journal_account_code_spoil, ' - ',  cfg_init_journal_account5.name) AS AccountSpoil " +
			"FROM" +
			" " + TableName +
			" LEFT OUTER JOIN cfg_init_journal_account ON (" + TableName + ".journal_account_code = cfg_init_journal_account.code)" +
			" LEFT OUTER JOIN cfg_init_journal_account cfg_init_journal_account1 ON (" + TableName + ".journal_account_code_cogs = cfg_init_journal_account1.code)" +
			" LEFT OUTER JOIN cfg_init_journal_account cfg_init_journal_account2 ON (" + TableName + ".journal_account_code_expense = cfg_init_journal_account2.code)" +
			" LEFT OUTER JOIN cfg_init_journal_account cfg_init_journal_account3 ON (" + TableName + ".journal_account_code_sell = cfg_init_journal_account3.code)" +
			" LEFT OUTER JOIN cfg_init_journal_account cfg_init_journal_account4 ON (" + TableName + ".journal_account_code_depreciation = cfg_init_journal_account4.code)" +
			" LEFT OUTER JOIN cfg_init_journal_account cfg_init_journal_account5 ON (" + TableName + ".journal_account_code_spoil = cfg_init_journal_account5.code)" +
			QueryCondition + " " +
			"ORDER BY " + TableName + ".code;"

	} else if DataName == "FAItem" {
		SearchField = GetMasterDataValidationSearchField93(Index, TableName, "inv_cfg_init_uom.name", "fa_cfg_init_item_category.name", "remark")

		if Text != "" {
			QueryCondition = " WHERE " + SearchField
		}
		Query = "SELECT" +
			" " + TableName + ".*," +
			" fa_cfg_init_item_category.name AS CategoryName," +
			" inv_cfg_init_uom.name AS UOMName " +
			"FROM" +
			" " + TableName +
			" LEFT OUTER JOIN fa_cfg_init_item_category ON (" + TableName + ".category_code = fa_cfg_init_item_category.code)" +
			" LEFT OUTER JOIN inv_cfg_init_uom ON (" + TableName + ".uom_code = inv_cfg_init_uom.code)" +
			QueryCondition + " " +
			"ORDER BY " + TableName + ".code;"
	} else if (DataName == "JournalAccountCategory") || (DataName == "TaxAndService") ||
		(DataName == "RoomType") || (DataName == "BedType") || (DataName == "RoomRateCategory") || (DataName == "RoomView") || (DataName == "RoomAmenities") ||
		(DataName == "RoomBoy") || (DataName == "Owner") || (DataName == "Title") ||
		(DataName == "Continent") || (DataName == "Regency") || (DataName == "Nationality") || (DataName == "Language") || (DataName == "IdCardType") ||
		(DataName == "MarketCategory") || (DataName == "BookingSource") || (DataName == "PurposeOf") ||
		(DataName == "CardBank") || (DataName == "CardType") || (DataName == "LoanItem") || (DataName == "PhoneBookType") || (DataName == "VoucherReason") ||
		(DataName == "CompetitorCategory") || (DataName == "SalesSegment") || (DataName == "SalesSource") || (DataName == "SalesTaskAction") ||
		(DataName == "SalesTaskRepeat") || (DataName == "SalesTaskTag") || (DataName == "POSMarket") || (DataName == "Waitress") || (DataName == "TableType") ||
		(DataName == "GuestType") ||
		//CAMS MODULE
		(DataName == "UOM") || (DataName == "Store") || (DataName == "ItemGroup") || (DataName == "ReturnStockReason") || (DataName == "FAManufacture") {
		SearchField = GetMasterDataValidationSearchField1(Index)
	} else if DataName == "ItemGroup" {
		SearchField = GetMasterDataValidationSearchField91(Index, TableName, "const_inv_item_group_type.name")
		if Text != "" {
			QueryCondition = " WHERE " + SearchField
		}
		Query = "SELECT" +
			" " + TableName + ".*," +
			" const_inv_item_group_type.name AS TypeName " +
			"FROM" +
			" " + TableName +
			" LEFT OUTER JOIN const_inv_item_group_type ON (" + TableName + ".item_group_type = const_inv_item_group_type.code)" +
			QueryCondition + " " +
			"ORDER BY " + TableName + ".code;"
	} else if DataName == "RoomUnavailableReason" {
		SearchField = GetMasterDataValidationSearchField2(Index)
	} else if DataName == "Printer" {
		SearchField = GetMasterDataValidationSearchField91(Index, TableName, "location")
	} else if DataName == "ProductCategory" {
		SearchField = GetMasterDataValidationSearchField92(Index, TableName, "description", "outlet_code")
	} else if DataName == "Currency" {
		SearchField = GetMasterDataValidationSearchField92(Index, TableName, "account_code", "symbol")
	} else if DataName == "Sales" {
		SearchField = GetMasterDataValidationSearchField93(Index, TableName, "email", "phone_number", "wa_number")
	}
	return TableName, Query, SearchField
}

func GetMasterDataListP(c *gin.Context) {
	DataName := c.Param("DataName")
	IndexP := c.Query("Index")
	ctx, span := global_var.Tracer.Start(utils.GetRequestCtx(c), "GetMasterDataListP", trace.WithAttributes(attribute.String("DataName", DataName)))
	defer span.End()

	var Index int64
	var err error
	if IndexP != "" {
		Index, err = strconv.ParseInt(IndexP, 10, 64)
	}
	Text := c.Query("Text")
	var Option1 int64 = 0
	String1 := ""
	String2 := ""
	if DataName == "Company" || DataName == "RoomRate" || DataName == "Item" {
		Option1P := c.Query("Option1")
		if Option1P != "" {
			Option1, err = strconv.ParseInt(Option1P, 10, 64)
		}
	} else if DataName == "Product" {
		Option1P := c.Query("Option1")
		String1 = c.Query("String1")
		String2 = c.Query("String2")
		if Option1P != "" {
			Option1, err = strconv.ParseInt(Option1P, 10, 64)
		}
	}

	if err != nil {
		SendResponse(global_var.ResponseCode.InvalidDataFormat, "", nil, c)
	} else {
		// Get Program Configuration
		val, exist := c.Get("pConfig")
		if !exist {
			SendResponse(global_var.ResponseCode.DatabaseError, "pConfig", nil, c)
			return
		}
		pConfig := val.(*config.CompanyDataConfiguration)
		DB := pConfig.DB

		TableName, Query, SearchField := GetMasterDataListValidation(DataName, Index, Text, Option1, String1, String2)
		fmt.Println("TableName", TableName)
		fmt.Println("Query", Query)
		fmt.Println("SearchField", SearchField)
		var DataOutput []map[string]interface{}
		if Query == "" {
			if SearchField == "" {
				result := DB.WithContext(ctx).Table(TableName).Order("id").Scan(&DataOutput)
				err = result.Error
			} else {
				if Text == "" {
					result := DB.WithContext(ctx).Table(TableName).Order("id").Scan(&DataOutput)
					err = result.Error
				} else {
					result := DB.WithContext(ctx).Table(TableName).Where(SearchField, "%"+Text+"%").Order("id").Scan(&DataOutput)
					err = result.Error
				}
			}
		} else {
			if Text == "" {
				result := DB.WithContext(ctx).Raw(Query).Scan(&DataOutput)
				err = result.Error
			} else {
				result := DB.WithContext(ctx).Raw(Query, "%"+Text+"%").Scan(&DataOutput)
				err = result.Error
			}

		}

		if err == nil {
			SendResponse(global_var.ResponseCode.Successfully, "", DataOutput, c)
		} else {
			SendResponse(global_var.ResponseCode.DatabaseError, "", nil, c)
		}
	}
}

func GetDetailDataListP(c *gin.Context) {
	DataName := c.Param("DataName")
	Param := c.Param("Param")
	ctx, span := global_var.Tracer.Start(utils.GetRequestCtx(c),
		"GetDetailDataListP",
		trace.WithAttributes(attribute.String("DataName", DataName),
			attribute.String("Param", Param)))
	defer span.End()

	var err error
	DataOutputArray := make(map[string][]map[string]interface{})

	// Get Program Configuration
	val, exist := c.Get("pConfig")
	if !exist {
		SendResponse(global_var.ResponseCode.DatabaseError, "pConfig", nil, c)
		return
	}
	pConfig := val.(*config.CompanyDataConfiguration)
	DB := pConfig.DB

	if DataName == "Package" {
		var DataOutput []map[string]interface{}
		_, Query, _ := GetMasterDataListValidation("PackageBreakdown", 0, "", 0, "", "")
		result := DB.WithContext(ctx).Raw(Query, Param).Scan(&DataOutput)
		if result.Error == nil {
			DataOutputArray["PackageBreakdown"] = DataOutput
		}

		DataOutput = nil
		_, Query, _ = GetMasterDataListValidation("PackageBusinessSource", 0, "", 0, "", "")
		result = DB.WithContext(ctx).Raw(Query, Param).Scan(&DataOutput)
		if result.Error == nil {
			DataOutputArray["PackageBusinessSource"] = DataOutput
		}
	} else if DataName == "RoomRate" {
		var DataOutput []map[string]interface{}
		_, Query, _ := GetMasterDataListValidation("RoomRateBreakdown", 0, "", 0, "", "")
		result := DB.WithContext(ctx).Raw(Query, Param).Scan(&DataOutput)
		if result.Error == nil {
			DataOutputArray["RoomRateBreakdown"] = DataOutput
		}

		DataOutput = nil
		_, Query, _ = GetMasterDataListValidation("RoomRateBusinessSource", 0, "", 0, "", "")
		result = DB.WithContext(ctx).Raw(Query, Param).Scan(&DataOutput)
		if result.Error == nil {
			DataOutputArray["RoomRateBusinessSource"] = DataOutput
		}

		DataOutput = nil
		_, Query, _ = GetMasterDataListValidation("RoomRateDynamic", 0, "", 0, "", "")
		result = DB.WithContext(ctx).Raw(Query, Param).Scan(&DataOutput)
		if result.Error == nil {
			DataOutputArray["RoomRateDynamic"] = DataOutput
		}

		DataOutput = nil
		_, Query, _ = GetMasterDataListValidation("RoomRateCurrency", 0, "", 0, "", "")
		result = DB.WithContext(ctx).Raw(Query, Param).Scan(&DataOutput)
		if result.Error == nil {
			DataOutputArray["RoomRateCurrency"] = DataOutput
		}
	} else if DataName == "ItemCategory" {
		var DataOutput []map[string]interface{}
		_, Query, _ := GetMasterDataListValidation("ItemCategoryOtherCOGS", 0, "", 0, "", "")
		result := DB.WithContext(ctx).Raw(Query, Param).Scan(&DataOutput)
		if result.Error == nil {
			DataOutputArray["ItemCategoryOtherCOGS"] = DataOutput
		}

		DataOutput = nil
		_, Query, _ = GetMasterDataListValidation("ItemCategoryOtherCOGS2", 0, "", 0, "", "")
		result = DB.WithContext(ctx).Raw(Query, Param).Scan(&DataOutput)
		if result.Error == nil {
			DataOutputArray["ItemCategoryOtherCOGS2"] = DataOutput
		}

		DataOutput = nil
		_, Query, _ = GetMasterDataListValidation("ItemCategoryOtherExpense", 0, "", 0, "", "")
		result = DB.WithContext(ctx).Raw(Query, Param).Scan(&DataOutput)
		if result.Error == nil {
			DataOutputArray["ItemCategoryOtherExpense"] = DataOutput
		}
	} else if DataName == "Item" {
		var DataOutput []map[string]interface{}
		_, Query, _ := GetMasterDataListValidation("ItemUOM", 0, "", 0, "", "")
		result := DB.WithContext(ctx).Raw(Query, Param).Scan(&DataOutput)
		if result.Error == nil {
			DataOutputArray["ItemUOM"] = DataOutput
		}
	} else if DataName == "MemberOutletDiscount" {
		var DataOutput []map[string]interface{}
		_, Query, _ := GetMasterDataListValidation("MemberOutletDiscountDetail", 0, "", 0, "", "")
		result := DB.WithContext(ctx).Raw(Query, Param).Scan(&DataOutput)
		if result.Error == nil {
			DataOutputArray["MemberOutletDiscountDetail"] = DataOutput
		}
	}

	if err == nil {
		SendResponse(global_var.ResponseCode.Successfully, "", DataOutputArray, c)
	} else {
		SendResponse(global_var.ResponseCode.DatabaseError, "", nil, c)
	}
}

func MasterDataValidation(ModeEditor byte, c *gin.Context, DataName string) (string, map[string]interface{}, error) {
	var DataInput map[string]interface{}
	var err error
	var DataInputMarshal []byte

	TableName := GetMasterDataTableName(DataName)

	// Get Program Configuration
	val, exist := c.Get("pConfig")
	if !exist {
		// return
	}
	pConfig := val.(*config.CompanyDataConfiguration)
	DB := pConfig.DB

	if DataName == "JournalAccount" {
		var DataInputX db_var.Cfg_init_journal_account
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "JournalAccountSubGroup" {
		var DataInputX db_var.Cfg_init_journal_account_sub_group
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "JournalAccountCategory" {
		var DataInputX db_var.Cfg_init_journal_account_category
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "Department" {
		var DataInputX db_var.Cfg_init_department
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "SubDepartment" {
		var DataInputX db_var.Cfg_init_sub_department
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "TaxAndService" {
		var DataInputX db_var.Cfg_init_tax_and_service
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "AccountSubGroup" {
		var DataInputX db_var.Cfg_init_account_sub_group
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "Account" {
		var DataInputX db_var.Cfg_init_account
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "BankAccount" {
		var DataInputX db_var.Acc_cfg_init_bank_account
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckData21(c, DB, TableName, "code", "journal_account_code", DataInputX.Code, DataInputX.JournalAccountCode) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				if CheckData22(c, DB, TableName, "id", "journal_account_code", DataInputX.Id, DataInputX.JournalAccountCode) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			}
		}
	} else if DataName == "CompanyType" {
		var DataInputX db_var.Cfg_init_company_type
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "Company" {
		var DataInputX db_var.Company
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "Currency" {
		var DataInputX db_var.Cfg_init_currency
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "CurrencyNominal" {
		var DataInputX db_var.Cfg_init_currency_nominal
		err = c.BindJSON(&DataInputX)
		if err == nil {
			DataInputMarshal, err = json.Marshal(DataInputX)
		}
	} else if DataName == "RoomType" {
		var DataInputX db_var.Cfg_init_room_type
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "BedType" {
		var DataInputX db_var.Cfg_init_bed_type
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "Package" {
		var DataInputX db_var.Cfg_init_package
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "PackageBreakdown" {
		var DataInputX db_var.Cfg_init_package_breakdown
		err = c.BindJSON(&DataInputX)
		if err == nil {
			DataInputMarshal, err = json.Marshal(DataInputX)
		}
	} else if DataName == "PackageBusinessSource" {
		var DataInputX db_var.Cfg_init_package_business_source
		err = c.BindJSON(&DataInputX)
		if err == nil {
			DataInputMarshal, err = json.Marshal(DataInputX)
		}
	} else if DataName == "RoomRateCategory" {
		var DataInputX db_var.Cfg_init_room_rate_category
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "RoomRateSubCategory" {
		var DataInputX db_var.Cfg_init_room_rate_sub_category
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "RoomRate" {
		var DataInputX db_var.Cfg_init_room_rate
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "RoomRateBreakdown" {
		var DataInputX db_var.Cfg_init_room_rate_breakdown
		err = c.BindJSON(&DataInputX)
		if err == nil {
			DataInputMarshal, err = json.Marshal(DataInputX)
		}
	} else if DataName == "RoomRateBusinessSource" {
		var DataInputX db_var.Cfg_init_room_rate_business_source
		err = c.BindJSON(&DataInputX)
		if err == nil {
			DataInputMarshal, err = json.Marshal(DataInputX)
		}
	} else if DataName == "RoomRateDynamic" {
		var DataInputX db_var.Cfg_init_room_rate_dynamic
		err = c.BindJSON(&DataInputX)
		if err == nil {
			DataInputMarshal, err = json.Marshal(DataInputX)
		}
	} else if DataName == "RoomRateCurrency" {
		var DataInputX db_var.Cfg_init_room_rate_currency
		err = c.BindJSON(&DataInputX)
		if err == nil {
			DataInputMarshal, err = json.Marshal(DataInputX)
		}
	} else if DataName == "RoomView" {
		var DataInputX db_var.Cfg_init_room_view
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "RoomAmenities" {
		var DataInputX db_var.Cfg_init_room_amenities
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "Room" {
		var DataInputX db_var.Cfg_init_room
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckNumber(c, DB, TableName, DataInputX.Number) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "RoomBoy" {
		var DataInputX db_var.Cfg_init_room_boy
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "RoomUnavailableReason" {
		var DataInputX db_var.Cfg_init_room_unavailable_reason
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "Owner" {
		var DataInputX db_var.Cfg_init_owner
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "Title" {
		var DataInputX db_var.Cfg_init_owner
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "Continent" {
		var DataInputX db_var.Cfg_init_continent
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "Country" {
		var DataInputX db_var.Cfg_init_country
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "State" {
		var DataInputX db_var.Cfg_init_state
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "Regency" {
		var DataInputX db_var.Cfg_init_regency
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "City" {
		var DataInputX db_var.Cfg_init_city
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "Nationality" {
		var DataInputX db_var.Cfg_init_nationality
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "Language" {
		var DataInputX db_var.Cfg_init_language
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "IdCardType" {
		var DataInputX db_var.Cfg_init_id_card_type
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "PaymentType" {
		var DataInputX db_var.Cfg_init_payment_type
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "MarketCategory" {
		var DataInputX db_var.Cfg_init_market_category
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "Market" {
		var DataInputX db_var.Cfg_init_market
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "BookingSource" {
		var DataInputX db_var.Cfg_init_booking_source
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "PurposeOf" {
		var DataInputX db_var.Cfg_init_purpose_of
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "CardBank" {
		var DataInputX db_var.Cfg_init_card_bank
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "CardType" {
		var DataInputX db_var.Cfg_init_card_type
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "LoanItem" {
		var DataInputX db_var.Cfg_init_loan_item
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "CreditCardCharge" {
		var DataInputX db_var.Cfg_init_credit_card_charge
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckData(c, DB, TableName, "account_code", DataInputX.AccountCode) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				if CheckData22(c, DB, TableName, "id", "account_code", DataInputX.Id, DataInputX.AccountCode) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			}
		}
	} else if DataName == "PhoneBookType" {
		var DataInputX db_var.Cfg_init_phone_book_type
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "MemberPointType" {
		var DataInputX db_var.Cfg_init_member_point_type
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "MemberOutletDiscount" {
		var DataInputX db_var.Pos_cfg_init_member_outlet_discount
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "MemberOutletDiscountDetail" {
		var DataInputX db_var.Pos_cfg_init_member_outlet_discount_detail
		err = c.BindJSON(&DataInputX)
		if err == nil {
			DataInputMarshal, err = json.Marshal(DataInputX)

		}
	} else if DataName == "VoucherReason" {
		var DataInputX db_var.Cfg_init_voucher_reason
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "CompetitorCategory" {
		var DataInputX db_var.Cfg_init_competitor_category
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "Competitor" {
		var DataInputX db_var.Competitor
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "Sales" {
		var DataInputX db_var.Cfg_init_sales
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "SalesSegment" {
		var DataInputX db_var.Sal_cfg_init_segment
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "SalesSource" {
		var DataInputX db_var.Sal_cfg_init_source
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "SalesTaskAction" {
		var DataInputX db_var.Sal_cfg_init_task_action
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "SalesTaskRepeat" {
		var DataInputX db_var.Sal_cfg_init_task_repeat
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "SalesTaskTag" {
		var DataInputX db_var.Sal_cfg_init_task_tag
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
		//POS MODULE
	} else if DataName == "Outlet" {
		var DataInputX db_var.Pos_cfg_init_outlet
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) || CheckData(c, DB, TableName, "check_prefix", DataInputX.CheckPrefix) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "ProductCategory" {
		var DataInputX db_var.Pos_cfg_init_product_category
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "Tenan" {
		var DataInputX db_var.Pos_cfg_init_tenan
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "ProductGroup" {
		var DataInputX db_var.Pos_cfg_init_product_group
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "Product" {
		var DataInputX db_var.Pos_cfg_init_product
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "POSMarket" {
		var DataInputX db_var.Pos_cfg_init_market
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "POSPaymentGroup" {
		var DataInputX db_var.Pos_cfg_init_payment_group
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckData(c, DB, TableName, "account_code", DataInputX.AccountCode) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "SpaRoom" {
		var DataInputX db_var.Pos_cfg_init_spa_room
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckNumber(c, DB, TableName, DataInputX.Number) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "TableType" {
		var DataInputX db_var.Pos_cfg_init_table_type
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "Table" {
		var DataInputX db_var.Pos_cfg_init_table
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckNumber(c, DB, TableName, DataInputX.Number) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "Waitress" {
		var DataInputX db_var.Pos_cfg_init_waitress
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "Printer" {
		var DataInputX db_var.Cfg_init_printer
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "DiscountLimit" {
		var DataInputX db_var.Pos_cfg_init_discount_limit
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckData2(c, DB, TableName, "outlet_code", "user_group_code", DataInputX.OutletCode, DataInputX.UserGroupCode) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				if CheckData31(c, DB, TableName, "id", "outlet_code", "user_group_code", DataInputX.Id, DataInputX.OutletCode, DataInputX.UserGroupCode) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			}
		}
		//CAMS MODULE
	} else if DataName == "ShippingAddress" {
		var DataInputX db_var.Ast_cfg_init_shipping_address
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "UOM" {
		var DataInputX db_var.Inv_cfg_init_uom
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "Store" {
		var DataInputX db_var.Inv_cfg_init_store
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "ItemCategory" {
		var DataInputX db_var.Inv_cfg_init_item_category
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "ItemCategoryOtherCOGS" {
		var DataInputX db_var.Inv_cfg_init_item_category_other_cogs
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckData2(c, DB, TableName, "sub_department_code", "category_code", DataInputX.SubDepartmentCode, DataInputX.CategoryCode) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				if CheckData31(c, DB, TableName, "id", "sub_department_code", "category_code", DataInputX.Id, DataInputX.SubDepartmentCode, DataInputX.CategoryCode) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			}
		}
	} else if DataName == "ItemCategoryOtherCOGS2" {
		var DataInputX db_var.Inv_cfg_init_item_category_other_cogs2
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckData2(c, DB, TableName, "sub_department_code", "category_code", DataInputX.SubDepartmentCode, DataInputX.CategoryCode) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				if CheckData31(c, DB, TableName, "id", "sub_department_code", "category_code", DataInputX.Id, DataInputX.SubDepartmentCode, DataInputX.CategoryCode) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			}
		}
	} else if DataName == "ItemCategoryOtherExpense" {
		var DataInputX db_var.Inv_cfg_init_item_category_other_expense
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckData2(c, DB, TableName, "sub_department_code", "category_code", DataInputX.SubDepartmentCode, DataInputX.CategoryCode) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				if CheckData31(c, DB, TableName, "id", "sub_department_code", "category_code", DataInputX.Id, DataInputX.SubDepartmentCode, DataInputX.CategoryCode) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			}
		}
	} else if DataName == "Item" {
		var DataInputX db_var.Inv_cfg_init_item
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "ItemUOM" {
		var DataInputX db_var.Inv_cfg_init_item_uom
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckData2(c, DB, TableName, "uom_code", "item_code", DataInputX.UomCode, DataInputX.ItemCode) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				if CheckData31(c, DB, TableName, "id", "uom_code", "item_code", DataInputX.Id, DataInputX.UomCode, DataInputX.ItemCode) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			}
		}
	} else if DataName == "ReturnStockReason" {
		var DataInputX db_var.Inv_cfg_init_return_stock_reason
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "MarketList" {
		var DataInputX db_var.Inv_cfg_init_market_list
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckData31(c, DB, TableName, "id", "company_code", "item_code", DataInputX.Id, DataInputX.CompanyCode, DataInputX.ItemCode) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "FAManufacture" {
		var DataInputX db_var.Fa_cfg_init_manufacture
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "FALocation" {
		var DataInputX db_var.Fa_cfg_init_location
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "FAItemCategory" {
		var DataInputX db_var.Fa_cfg_init_item_category
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "FAItem" {
		var DataInputX db_var.Fa_cfg_init_item
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "ItemGroup" {
		var DataInputX db_var.Inv_cfg_init_item_group
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else if DataName == "GuestType" {
		var DataInputX db_var.Cfg_init_guest_type
		err = c.BindJSON(&DataInputX)
		if err == nil {
			if ModeEditor == 0 {
				if CheckCode(c, DB, TableName, DataInputX.Code) {
					err = errors.New(global_var.ResponseText.DuplicateEntry)
				} else {
					DataInputMarshal, err = json.Marshal(DataInputX)
				}
			} else {
				DataInputMarshal, err = json.Marshal(DataInputX)
			}
		}
	} else {
		err = errors.New("Invalid Data Name")
	}
	// fmt.Println(DataInputMarshal)
	// fmt.Println("B")
	// fmt.Println(DataInput)
	if err == nil {
		json.Unmarshal(DataInputMarshal, &DataInput)
	}
	return TableName, DataInput, err
}

func InsertMasterDataP(c *gin.Context) {
	DataName := c.Param("DataName")
	ctx, span := global_var.Tracer.Start(utils.GetRequestCtx(c), "InsertMasterDataP",
		trace.WithAttributes(attribute.String("DataName", DataName)))
	defer span.End()

	if DataName != "" {
		TableName, DataInput, err := MasterDataValidation(0, c, DataName)
		if err != nil {
			if err.Error() == global_var.ResponseText.DuplicateEntry {
				SendResponse(global_var.ResponseCode.DuplicateEntry, "", nil, c)
			} else {
				StrError := general.GenerateValidateErrorMsg(c, err)
				SendResponse(global_var.ResponseCode.InvalidDataFormat, StrError, nil, c)
			}
		} else {

			// Get Program Configuration
			val, exist := c.Get("pConfig")
			if !exist {
				SendResponse(global_var.ResponseCode.DatabaseError, "pConfig", nil, c)
				return
			}
			pConfig := val.(*config.CompanyDataConfiguration)
			DB := pConfig.DB

			ValidUserCode := c.GetString("ValidUserCode")
			DataInput["created_by"] = ValidUserCode
			DataInput["updated_by"] = ""
			var result *gorm.DB
			if DataName == "Room" {
				var TotalRoom int64
				DB.WithContext(ctx).Table(db_var.TableName.CfgInitRoom).Count(&TotalRoom)

				if TotalRoom >= pConfig.Rooms {
					SendResponse(global_var.ResponseCode.OtherResult, "Cannot add more room", nil, c)
					return
				}

				DataInput["start_date"] = general.StrZToDate(DataInput["start_date"].(string))
				result = DB.WithContext(ctx).Table(TableName).Omit("status_code", "block_status_code", "temp_status_code", "remark", "pos_x", "pos_y", "width", "height", "created_at", "updated_at", "id").Create(&DataInput)
			} else if DataName == "SpaRoom" {
				result = DB.WithContext(ctx).Table(TableName).Omit("left", "top", "width", "height", "created_at", "updated_at", "id").Create(&DataInput)
			} else if DataName == "Table" {
				width := pConfig.Dataset.Configuration[global_var.ConfigurationCategoryPOS.TableView][global_var.ConfigurationNamePOS.MinRoomWidth].(string)
				height := pConfig.Dataset.Configuration[global_var.ConfigurationCategoryPOS.TableView][global_var.ConfigurationNamePOS.MinRoomHeight].(string)

				DataInput["width"] = width
				DataInput["height"] = height
				DataInput["start_date"] = time.Now()

				result = DB.WithContext(ctx).Table(TableName).Omit("left", "top", "created_at", "updated_at", "id").Create(&DataInput)
			} else if DataName == "RoomRate" {
				DataInput["cm_start_date"] = general.StrZToDate(DataInput["cm_start_date"].(string))
				DataInput["cm_end_date"] = general.StrZToDate(DataInput["cm_end_date"].(string))
				DataInput["from_date"] = general.StrZToDate(DataInput["from_date"].(string))
				DataInput["to_date"] = general.StrZToDate(DataInput["to_date"].(string))
				result = DB.WithContext(ctx).Table(TableName).Omit("is_cm_updated", "is_cm_updated_inclusion", "is_sent", "created_at", "updated_at", "id").Create(&DataInput)
			} else if DataName == "Product" {
				result = DB.WithContext(ctx).Table(TableName).Omit("is_sold", "created_at", "updated_at", "id").Create(&DataInput)
			} else if DataName == "MemberOutletDiscountDetail" {
				if DataInput["discount_percent"].(float64) > 100 || DataInput["discount_percent"].(float64) < 0 {
					SendResponse(global_var.ResponseCode.InvalidDataFormat, "Invalid discount value", nil, c)
					return
				}

				var ID uint64
				if err := DB.WithContext(ctx).Table(db_var.TableName.PosCfgInitMemberOutletDiscountDetail).Select("id").
					Where("product_code=?", DataInput["product_code"].(string)).
					Where("member_outlet_discount_code=?", DataInput["member_outlet_discount_code"].(string)).
					Limit(1).
					Scan(&ID).Error; err != nil {
					// utils.LogResponseError(c, ctx, global_var.AppLogger, errors.Wrap(err, "ProcessMemberProductDiscountP.query"))
					SendResponse(global_var.ResponseCode.DatabaseError, "", nil, c)
					return
				}

				if ID > 0 {
					DataInput["updated_by"] = ValidUserCode

					if err := DB.WithContext(ctx).Table(TableName).
						Where("id=?", ID).
						Omit("created_at, created_by,product_code,outlet_code, member_outlet_discount_code,id, updated_at").
						Updates(&DataInput).Error; err != nil {
						// utils.LogResponseError(c, ctx, global_var.AppLogger, errors.Wrap(err, "ProcessMemberProductDiscountP.UpdatePosCfgInitMemberProductDiscount"))
						SendResponse(global_var.ResponseCode.DatabaseError, "", nil, c)
						return
					}
				} else {
					if err := DB.WithContext(ctx).Table(TableName).Omit("created_at,updated_at,id").Create(&DataInput).Error; err != nil {
						// utils.LogResponseError(c, ctx, global_var.AppLogger, errors.Wrap(err, "ProcessMemberProductDiscountP.InsertPosCfgInitMemberProductDiscount"))
						SendResponse(global_var.ResponseCode.DatabaseError, "", nil, c)
						return
					}
				}
				SendResponse(global_var.ResponseCode.Successfully, "", nil, c)
				return
			} else {
				result = DB.WithContext(ctx).Table(TableName).Omit("created_at", "updated_at", "id").Create(&DataInput)
			}
			if result.Error == nil {
				SendResponse(global_var.ResponseCode.Successfully, "", nil, c)
			} else {
				SendResponse(global_var.ResponseCode.DatabaseError, "", nil, c)
			}
		}
	}
}

func GetMasterDataP(c *gin.Context) {
	DataName := c.Param("DataName")
	Param := c.Param("Param")
	TableName := GetMasterDataTableName(DataName)
	var result *gorm.DB
	var DataOutput map[string]interface{}

	// Get Program Configuration
	val, exist := c.Get("pConfig")
	if !exist {
		SendResponse(global_var.ResponseCode.DatabaseError, "pConfig", nil, c)
		return
	}
	pConfig := val.(*config.CompanyDataConfiguration)
	DB := pConfig.DB

	if (DataName == "Room") || (DataName == "SpaRoom") || (DataName == "Table") {
		result = DB.Table(TableName).Where("number=?", Param).Scan(&DataOutput)
	} else if (DataName == "PackageBreakdown") || (DataName == "PackageBusinessSource") ||
		(DataName == "RoomRateBreakdown") || (DataName == "RoomRateBusinessSource") || (DataName == "RoomRateDynamic") || (DataName == "RoomRateCurrency") ||
		(DataName == "CurrencyNominal") || (DataName == "DiscountLimit") || (DataName == "MarketList") ||
		(DataName == "ItemCategoryOtherCOGS") || (DataName == "ItemCategoryOtherCOGS2") || (DataName == "ItemCategoryOtherExpense") ||
		(DataName == "ItemUOM") || (DataName == "MemberOutletDiscountDetail") {
		result = DB.Table(TableName).Where("id=?", Param).Scan(&DataOutput)
	} else if (DataName == "POSPaymentGroup") || (DataName == "CreditCardCharge") {
		result = DB.Table(TableName).Where("account_code=?", Param).Scan(&DataOutput)
	} else {
		result = DB.Table(TableName).Where("code=?", Param).Scan(&DataOutput)
	}
	if result.Error == nil {
		SendResponse(global_var.ResponseCode.Successfully, "", DataOutput, c)
	} else {
		SendResponse(global_var.ResponseCode.DatabaseError, "", nil, c)
	}
}

func UpdateMasterDataP(c *gin.Context) {
	DataName := c.Param("DataName")
	if DataName != "" {
		// Get Program Configuration
		val, exist := c.Get("pConfig")
		if !exist {
			SendResponse(global_var.ResponseCode.DatabaseError, "pConfig", nil, c)
			return
		}
		pConfig := val.(*config.CompanyDataConfiguration)
		DB := pConfig.DB

		var KeyField string = "code"
		if (DataName == "Room") || (DataName == "SpaRoom") || (DataName == "Table") {
			KeyField = "number"
		} else if (DataName == "PackageBreakdown") || (DataName == "PackageBusinessSource") || (DataName == "CurrencyNominal") ||
			(DataName == "RoomRateBreakdown") || (DataName == "RoomRateBusinessSource") || (DataName == "RoomRateDynamic") || (DataName == "RoomRateCurrency") ||
			(DataName == "CreditCardCharge") || (DataName == "DiscountLimit") || (DataName == "MarketList") || (DataName == "ItemUOM") ||
			(DataName == "ItemCategoryOtherCOGS") || (DataName == "ItemCategoryOtherCOGS2") || (DataName == "ItemCategoryOtherExpense") {
			KeyField = "id"
		} else if DataName == "POSPaymentGroup" {
			KeyField = "account_code"
		}
		TableName, DataInput, err := MasterDataValidation(1, c, DataName)
		if err != nil {
			if err.Error() == global_var.ResponseText.DuplicateEntry {
				SendResponse(global_var.ResponseCode.DuplicateEntry, "", nil, c)
			} else {
				StrError := general.GenerateValidateErrorMsg(c, err)
				SendResponse(global_var.ResponseCode.InvalidDataFormat, StrError, nil, c)
			}
		} else {
			ValidUserCode := c.GetString("ValidUserCode")
			DataInput["updated_by"] = ValidUserCode
			var result *gorm.DB
			if DataName == "Room" {
				DataInput["start_date"] = general.StrZToDate(DataInput["start_date"].(string))
				result = DB.Table(TableName).Where(KeyField+"=?", DataInput[KeyField]).Omit(KeyField, "status_code", "block_status_code", "temp_status_code", "remark", "pos_x", "pos_y", "width", "height", "created_at", "created_by", "updated_at", "id").Updates(&DataInput)
			} else if DataName == "SpaRoom" {
				result = DB.Table(TableName).Where(KeyField+"=?", DataInput[KeyField]).Omit(KeyField, "left", "top", "width", "height", "created_at", "created_by", "updated_at", "id").Updates(&DataInput)
			} else if DataName == "Table" {
				result = DB.Table(TableName).Where(KeyField+"=?", DataInput[KeyField]).Omit(KeyField, "start_date", "left", "top", "width", "height", "created_at", "created_by", "updated_at", "id").Updates(&DataInput)
			} else if DataName == "RoomRate" {
				DataInput["cm_start_date"] = general.StrZToDate(DataInput["cm_start_date"].(string))
				DataInput["cm_end_date"] = general.StrZToDate(DataInput["cm_end_date"].(string))
				DataInput["from_date"] = general.StrZToDate(DataInput["from_date"].(string))
				DataInput["to_date"] = general.StrZToDate(DataInput["to_date"].(string))
				result = DB.Table(TableName).Where(KeyField+"=?", DataInput[KeyField]).Omit(KeyField, "is_cm_updated", "is_cm_updated_inclusion", "is_sent", "created_at", "created_by", "updated_at", "id").Updates(&DataInput)
			} else if (DataName == "PackageBreakdown") || (DataName == "PackageBusinessSource") {
				result = DB.Table(TableName).Where(KeyField+"=?", general.InterfaceToUint64(DataInput[KeyField])).Omit("package_code", "created_at", "created_by", "updated_at", KeyField).Updates(&DataInput)
			} else if (DataName == "RoomRateBreakdown") || (DataName == "RoomRateBusinessSource") || (DataName == "RoomRateDynamic") || (DataName == "RoomRateCurrency") {
				result = DB.Table(TableName).Where(KeyField+"=?", general.InterfaceToUint64(DataInput[KeyField])).Omit("room_rate_code", "created_at", "created_by", "updated_at", KeyField).Updates(&DataInput)
			} else if DataName == "CurrencyNominal" {
				result = DB.Table(TableName).Where(KeyField+"=?", general.InterfaceToUint64(DataInput[KeyField])).Omit("package_code", "created_at", "created_by", "updated_at", KeyField).Updates(&DataInput)
			} else if DataName == "Outlet" {
				result = DB.Table(TableName).Where(KeyField+"=?", DataInput[KeyField]).Omit("check_prefix", "created_at", "created_by", "updated_at", KeyField).Updates(&DataInput)
			} else if DataName == "Product" {
				result = DB.Table(TableName).Where(KeyField+"=?", DataInput[KeyField]).Omit("is_sold", "created_at", "created_by", "updated_at", KeyField).Updates(&DataInput)
			} else if (DataName == "ItemCategoryOtherCOGS") || (DataName == "ItemCategoryOtherCOGS2") || (DataName == "ItemCategoryOtherExpense") || (DataName == "ItemUOM") {
				result = DB.Table(TableName).Where(KeyField+"=?", general.InterfaceToUint64(DataInput[KeyField])).Omit("created_at", "created_by", "updated_at", KeyField).Updates(&DataInput)
			} else {
				result = DB.Table(TableName).Where(KeyField+"=?", DataInput[KeyField]).Omit(KeyField, "created_at", "created_by", "updated_at", "id").Updates(&DataInput)
			}

			if result.Error == nil {
				SendResponse(global_var.ResponseCode.Successfully, "", nil, c)
			} else {
				SendResponse(global_var.ResponseCode.DatabaseError, "", nil, c)
			}
		}
	}
}

func DeleteMasterDataP(c *gin.Context) {
	DataName := c.Param("DataName")
	Param := c.Param("Param")
	TableName := GetMasterDataTableName(DataName)

	// Get Program Configuration
	val, exist := c.Get("pConfig")
	if !exist {
		SendResponse(global_var.ResponseCode.DatabaseError, "pConfig", nil, c)
		return
	}
	pConfig := val.(*config.CompanyDataConfiguration)
	DB := pConfig.DB

	UpdatedBy := c.GetString("ValidUserCode")
	var KeyField string = "code"
	if (DataName == "Room") || (DataName == "SpaRoom") || (DataName == "Table") {
		KeyField = "number"
	} else if (DataName == "PackageBreakdown") || (DataName == "PackageBusinessSource") ||
		(DataName == "RoomRateBreakdown") || (DataName == "RoomRateBusinessSource") || (DataName == "RoomRateDynamic") || (DataName == "RoomRateCurrency") ||
		(DataName == "CurrencyNominal") || (DataName == "DiscountLimit") || (DataName == "MarketList") || (DataName == "ItemUOM") ||
		(DataName == "MemberOutletDiscountDetail") || (DataName == "ItemCategoryOtherCOGS") || (DataName == "ItemCategoryOtherCOGS2") || (DataName == "ItemCategoryOtherExpense") {
		KeyField = "id"
	} else if (DataName == "CreditCardCharge") || (DataName == "POSPaymentGroup") {
		KeyField = "account_code"
	}
	err := DB.Transaction(func(tx *gorm.DB) error {
		if err := DB.Table(TableName).Where(KeyField+"=?", Param).Updates(map[string]interface{}{"updated_by": UpdatedBy}).Error; err != nil {
			return err
		}
		if err := DB.Table(TableName).Where(KeyField+"=?", Param).Delete(&Param).Error; err != nil {
			return err
		}

		return nil
	})
	if err == nil {
		SendResponse(global_var.ResponseCode.Successfully, "", nil, c)
	} else {
		SendResponse(global_var.ResponseCode.DatabaseError, "", nil, c)
	}
}

//======================================== END MASTER DATA =================================================//

func GetMasterDataCodeName(DB *gorm.DB, TableName, Condition string, c *gin.Context) {
	Result, StrValue := ValidateRequestString(c)
	if Result != global_var.ResponseCode.Successfully {
		SendResponse(Result, "", nil, c)
	} else {
		var GeneralCodeName db_var.GeneralCodeNameStruct
		if Condition == "" {
			DB.Table(TableName).Select("code", "name").Where(Condition+" = ?", StrValue).First(&GeneralCodeName)
		} else {
			DB.Table(TableName).Select("code", "name").First(&GeneralCodeName)
		}
		SendResponse(global_var.ResponseCode.Successfully, "", GeneralCodeName, c)
	}
}

func GetMasterDataCodeNameList(DB *gorm.DB, TableName, Condition string, c *gin.Context) {
	// Result, StrValue := ValidateRequestString(c)
	StrValue := c.Param("Code")
	// if Result != global_var.ResponseCode.Successfully {
	// 	SendResponse(Result, "", nil, c)
	// } else {
	var GeneralCodeName []db_var.GeneralCodeNameStruct
	if Condition == "" {
		DB.Table(TableName).Select("code", "name").Scan(&GeneralCodeName)
	} else {
		DB.Table(TableName).Select("code", "name").Where(Condition+" = ?", StrValue).Scan(&GeneralCodeName)
	}
	SendResponse(global_var.ResponseCode.Successfully, "", GeneralCodeName, c)
	// }
}

func GetMasterDataCodeNameQuery(DB *gorm.DB, SQLQuery string, WithCondition bool, c *gin.Context) {
	Result, StrValue := ValidateRequestString(c)
	if Result != global_var.ResponseCode.Successfully {
		SendResponse(Result, "", nil, c)
	} else {
		var GeneralCodeName db_var.GeneralCodeNameStruct
		if WithCondition {
			DB.Raw(SQLQuery, StrValue).First(&GeneralCodeName)
		} else {
			DB.Raw(SQLQuery).First(&GeneralCodeName)
		}
		SendResponse(global_var.ResponseCode.Successfully, "", GeneralCodeName, c)
	}
}

func GetMasterDataCodeNameListQuery(DB *gorm.DB, SQLQuery string, WithCondition bool, c *gin.Context) {
	Result, StrValue := ValidateRequestString(c)
	if Result != global_var.ResponseCode.Successfully {
		SendResponse(Result, "", nil, c)
	} else {
		var GeneralCodeName []db_var.GeneralCodeNameStruct
		if WithCondition {
			DB.Raw(SQLQuery, StrValue).Scan(&GeneralCodeName)
		} else {
			DB.Raw(SQLQuery).Scan(&GeneralCodeName)
		}
		SendResponse(global_var.ResponseCode.Successfully, "", GeneralCodeName, c)
	}
}

func GetBusinessSourceCodeNameList(DB *gorm.DB) []db_var.GeneralCodeNameStruct {
	var GeneralCodeName []db_var.GeneralCodeNameStruct
	DB.Table(db_var.TableName.Company).Where("is_business_source = ?", 1).Find(&GeneralCodeName)
	return GeneralCodeName
}

func GetBedTypeByRoomTypeP(c *gin.Context) {
	// Get Program Configuration
	val, exist := c.Get("pConfig")
	if !exist {
		SendResponse(global_var.ResponseCode.DatabaseError, "pConfig", nil, c)
		return
	}
	pConfig := val.(*config.CompanyDataConfiguration)
	DB := pConfig.DB
	GetMasterDataCodeNameListQuery(DB,
		"SELECT"+
			" cfg_init_room.bed_type_code AS code,"+
			" cfg_init_bed_type.name "+
			"FROM"+
			" cfg_init_room"+
			" LEFT OUTER JOIN cfg_init_bed_type ON (cfg_init_room.bed_type_code = cfg_init_bed_type.code)"+
			" WHERE cfg_init_room.room_type_code = ? "+
			"GROUP BY cfg_init_room.bed_type_code;",
		true, c)
}

func GetBedTypeByRoomNumberP(c *gin.Context) {
	// Get Program Configuration
	val, exist := c.Get("pConfig")
	if !exist {
		SendResponse(global_var.ResponseCode.DatabaseError, "pConfig", nil, c)
		return
	}
	pConfig := val.(*config.CompanyDataConfiguration)
	DB := pConfig.DB
	GetMasterDataCodeNameQuery(DB,
		"SELECT"+
			" cfg_init_room.bed_type_code AS code,"+
			" cfg_init_bed_type.name "+
			"FROM"+
			" cfg_init_room"+
			" LEFT OUTER JOIN cfg_init_bed_type ON (cfg_init_room.bed_type_code = cfg_init_bed_type.code)"+
			" WHERE cfg_init_room.number = ? "+
			"GROUP BY cfg_init_room.bed_type_code;",
		true, c)
}

func GetStateByCountryP(c *gin.Context) {
	Code := c.Param("Code")
	// Get Program Configuration
	val, exist := c.Get("pConfig")
	if !exist {
		SendResponse(global_var.ResponseCode.DatabaseError, "pConfig", nil, c)
		return
	}
	pConfig := val.(*config.CompanyDataConfiguration)
	DB := pConfig.DB
	GeneralCodeName := GetGeneralCodeNameQuery(DB,
		"SELECT code, name FROM cfg_init_state"+
			" WHERE (country_code=?) "+
			"ORDER BY name;", Code, "", "", "", "", 1)

	SendResponse(global_var.ResponseCode.Successfully, "", GeneralCodeName, c)
}

func GetCityByStateP(c *gin.Context) {
	Code := c.Param("Code")
	// Get Program Configuration
	val, exist := c.Get("pConfig")
	if !exist {
		SendResponse(global_var.ResponseCode.DatabaseError, "pConfig", nil, c)
		return
	}
	pConfig := val.(*config.CompanyDataConfiguration)
	DB := pConfig.DB
	GeneralCodeName := GetGeneralCodeNameQuery(DB,
		"SELECT code, name FROM cfg_init_city"+
			" WHERE (state_code=? OR code =?) "+
			"ORDER BY name;", Code, global_var.ConstProgramVariable.CityOtherCode, "", "", "", 2)

	SendResponse(global_var.ResponseCode.Successfully, "", GeneralCodeName, c)
}

// func GetReservationComboListP(c *gin.Context) {
// 	var ReservationComboList db_var.ReservationComboListStruct
// 	GeneralCodeName := GetGeneralCodeName(DB,db_var.TableName.CfgInitRoomType, "")
// 	ReservationComboList.RoomType = GeneralCodeName
// 	GeneralCodeName = GetGeneralCodeName(DB,db_var.TableName.CfgInitCurrency, "")
// 	ReservationComboList.Currency = GeneralCodeName
// 	GeneralCodeName = GetBusinessSourceCodeNameList()
// 	ReservationComboList.BusinessSource = GeneralCodeName
// 	GeneralCodeName = GetGeneralCodeName(DB,db_var.TableName.ConstCommissionType, "")
// 	ReservationComboList.CommissionType = GeneralCodeName
// 	GeneralCodeName = GetGeneralCodeName(DB,db_var.TableName.CfgInitMarket, "")
// 	ReservationComboList.Market = GeneralCodeName
// 	GeneralCodeName = GetGeneralCodeName(DB,db_var.TableName.CfgInitBookingSource, "")
// 	ReservationComboList.BookingSource = GeneralCodeName
// 	GeneralCodeName = GetGeneralCodeName(DB,db_var.TableName.CfgInitPaymentType, "")
// 	ReservationComboList.PaymentType = GeneralCodeName
// 	GeneralCodeName = GetGeneralCodeName(DB,db_var.TableName.Member, "")
// 	ReservationComboList.Member = GeneralCodeName
// 	GeneralCodeName = GetGeneralCodeName(DB,db_var.TableName.CfgInitTitle, "")
// 	ReservationComboList.Title = GeneralCodeName
// 	GeneralCodeName = GetGeneralCodeName(DB,db_var.TableName.CfgInitCountry, "")
// 	ReservationComboList.Country = GeneralCodeName
// 	GeneralCodeName = GetGeneralCodeName(DB,db_var.TableName.CfgInitNationality, "")
// 	ReservationComboList.Nationality = GeneralCodeName
// 	GeneralCodeName = GetGeneralCodeName(DB,db_var.TableName.Company, "")
// 	ReservationComboList.Company = GeneralCodeName
// 	GeneralCodeName = GetGeneralCodeName(DB,db_var.TableName.CfgInitGuestType, "")
// 	ReservationComboList.GuestType = GeneralCodeName
// 	GeneralCodeName = GetGeneralCodeName(DB,db_var.TableName.CfgInitIdCardType, "")
// 	ReservationComboList.CardType = GeneralCodeName
// 	GeneralCodeName = GetGeneralCodeName(DB,db_var.TableName.CfgInitPurposeOf, "")
// 	ReservationComboList.PurposeOf = GeneralCodeName
// 	GeneralCodeName = GetGeneralCodeName(DB,db_var.TableName.GuestGroup, "")
// 	ReservationComboList.GuestGroup = GeneralCodeName
// 	GeneralCodeName = GetGeneralCodeName(DB,db_var.TableName.CfgInitSales, "")
// 	ReservationComboList.Sales = GeneralCodeName

// 	SendResponse(global_var.ResponseCode.Successfully, "", ReservationComboList, c)
// }

func GetReservationDepositComboListP(c *gin.Context) {
	var ReservationDepositComboList db_var.ReservationDepositComboListStruct
	// Get Program Configuration
	val, exist := c.Get("pConfig")
	if !exist {
		SendResponse(global_var.ResponseCode.DatabaseError, "pConfig", nil, c)
		return
	}
	pConfig := val.(*config.CompanyDataConfiguration)
	DB := pConfig.DB
	GeneralCodeName := GetGeneralCodeName(DB, db_var.TableName.CfgInitCurrency, "")
	ReservationDepositComboList.Currency = GeneralCodeName
	GeneralCodeName = GetGeneralCodeName(DB, db_var.TableName.SubFolioGroup, "")
	ReservationDepositComboList.SubFolioGroup = GeneralCodeName
	GeneralCodeName = GetGeneralCodeName(DB, db_var.TableName.CfgInitSubDepartment, "")
	ReservationDepositComboList.SubDepartment = GeneralCodeName

	SendResponse(global_var.ResponseCode.Successfully, "", ReservationDepositComboList, c)
}

func GetAccountCodeFromCurrency(c *gin.Context, DB *gorm.DB, CurrencyCode string) string {
	Result := "''"
	var DataOutput []map[string]interface{}
	QueryResult := DB.Raw(
		"SELECT account_code FROM cfg_init_currency" +
			" WHERE code='" + CurrencyCode + "';").Scan(&DataOutput)

	if QueryResult.RowsAffected > 0 {
		AccountCode := DataOutput[0]["account_code"].(string)
		if AccountCode != "" {
			AccountCodeArray := strings.Split(AccountCode, "|")
			Result := ""
			for I, S := range AccountCodeArray {
				if I < len(AccountCodeArray)-1 {
					Result = Result + AccountCode + "'" + S + "',"
					return Result
				} else {
					Result = "'" + S + "'"
					return Result
				}
			}
		}
	}
	return Result
}

func GetAccountBySubDepartmentTransactionEditor(c *gin.Context, DB *gorm.DB, ModeEditor byte, CurrencyCode, SubDepartmentCode string) []map[string]interface{} {
	//Mode Editor
	//0: Folio Charge
	//1: Folio Other Payment
	//2: Folio Cash Payment
	//3: Folio Card Payment
	//4: Folio Cash Refund
	//5: Folio Direct Bill
	//6: Deposit  Other
	//7: Deposit Cash
	//8: Deposit Card Payment
	//9: Deposit Cash Refund
	//10: Folio AP Transaction

	var DataOutput []map[string]interface{}
	if ModeEditor == 0 {
		DB.Raw(
			"SELECT" +
				" cfg_init_account.code," +
				" cfg_init_account.name," +
				" cfg_init_account.sub_folio_group_code " +
				"FROM" +
				" cfg_init_account" +
				" LEFT OUTER JOIN cfg_init_account_sub_group ON (cfg_init_account.sub_group_code = cfg_init_account_sub_group.code)" +
				" WHERE cfg_init_account_sub_group.group_code='" + global_var.GlobalAccountGroup.Charge + "'" +
				" AND cfg_init_account.sub_group_code<>'" + global_var.GlobalAccountSubGroup.AccountPayable + "'" +
				" AND cfg_init_account.sub_department_code LIKE '%" + SubDepartmentCode + "%' " +
				"ORDER BY cfg_init_account.code;").Scan(&DataOutput)
	} else if (ModeEditor == 1) || (ModeEditor == 6) {
		// S := "SELECT" +
		// 	" cfg_init_account.code," +
		// 	" cfg_init_account.name," +
		// 	" cfg_init_account.sub_folio_group_code " +
		// 	"FROM" +
		// 	" cfg_init_account" +
		// 	" LEFT OUTER JOIN cfg_init_account_sub_group ON (cfg_init_account.sub_group_code = cfg_init_account_sub_group.code)" +
		// 	" WHERE cfg_init_account_sub_group.group_code='" + global_var.GlobalAccountGroup.Payment + "'" +
		// 	" AND cfg_init_account.sub_group_code<>'" + global_var.GlobalAccountSubGroup.Payment + "'" +
		// 	" AND cfg_init_account.sub_group_code<>'" + global_var.GlobalAccountSubGroup.CreditDebitCard + "'" +
		// 	" AND cfg_init_account.sub_group_code<>'" + global_var.GlobalAccountSubGroup.AccountReceivable + "'" +
		// 	" AND cfg_init_account.sub_group_code<>'" + global_var.GlobalAccountSubGroup.Compliment + "'" +
		// 	" AND cfg_init_account.sub_department_code LIKE '%" + SubDepartmentCode + "%' " +
		// 	"ORDER BY cfg_init_account.code;"
		// // SaveTextToFile(S, "C:/Temp/A.txt")
		DB.Raw(
			"SELECT" +
				" cfg_init_account.code," +
				" cfg_init_account.name," +
				" cfg_init_account.sub_folio_group_code " +
				"FROM" +
				" cfg_init_account" +
				" LEFT OUTER JOIN cfg_init_account_sub_group ON (cfg_init_account.sub_group_code = cfg_init_account_sub_group.code)" +
				" WHERE cfg_init_account_sub_group.group_code='" + global_var.GlobalAccountGroup.Payment + "'" +
				" AND cfg_init_account.sub_group_code<>'" + global_var.GlobalAccountSubGroup.Payment + "'" +
				" AND cfg_init_account.sub_group_code<>'" + global_var.GlobalAccountSubGroup.CreditDebitCard + "'" +
				" AND cfg_init_account.sub_group_code<>'" + global_var.GlobalAccountSubGroup.AccountReceivable + "'" +
				" AND cfg_init_account.sub_group_code<>'" + global_var.GlobalAccountSubGroup.Compliment + "'" +
				" AND cfg_init_account.sub_department_code LIKE '%" + SubDepartmentCode + "%' " +
				"ORDER BY cfg_init_account.code;").Scan(&DataOutput)

	} else if (ModeEditor == 2) || (ModeEditor == 4) || (ModeEditor == 7) || (ModeEditor == 9) {
		AccountCodeList := GetAccountCodeFromCurrency(c, DB, CurrencyCode)
		if (ModeEditor == 2) || (ModeEditor == 7) {
			DB.Raw(
				"SELECT code, name, sub_folio_group_code FROM cfg_init_account" +
					" WHERE code IN (" + AccountCodeList + ")" +
					" AND sub_department_code LIKE '%" + SubDepartmentCode + "%'" +
					" AND is_payment = '1' " +
					"ORDER BY code;").Scan(&DataOutput)
		} else {
			DB.Raw(
				"SELECT code, name, sub_folio_group_code FROM cfg_init_account" +
					" WHERE code IN (" + AccountCodeList + ")" +
					" AND sub_department_code LIKE '%" + SubDepartmentCode + "%'" +
					" AND is_refund = '1' " +
					"ORDER BY code;").Scan(&DataOutput)
		}
	} else if (ModeEditor == 3) || (ModeEditor == 8) {
		DB.Raw(
			"SELECT code, name, sub_folio_group_code FROM cfg_init_account" +
				" WHERE sub_group_code='" + global_var.GlobalAccountSubGroup.CreditDebitCard + "'" +
				" AND sub_department_code LIKE '%" + SubDepartmentCode + "%' " +
				"ORDER BY code;").Scan(&DataOutput)
	} else if ModeEditor == 5 {
		DB.Raw(
			"SELECT code, name, sub_folio_group_code FROM cfg_init_account" +
				" WHERE sub_group_code='" + global_var.GlobalAccountSubGroup.AccountReceivable + "'" +
				" AND sub_department_code LIKE '%" + SubDepartmentCode + "%' " +
				"ORDER BY code;").Scan(&DataOutput)
	} else if ModeEditor == 10 {
		DB.Raw(
			"SELECT code, name, sub_folio_group_code FROM cfg_init_account" +
				" WHERE sub_group_code='" + global_var.GlobalAccountSubGroup.AccountPayable + "'" +
				" AND sub_department_code LIKE '%" + SubDepartmentCode + "%' " +
				"ORDER BY code;").Scan(&DataOutput)

		//fmt.Println(DataOutput)
	}
	return DataOutput
}

func GetAccountBySubDepartmentTransactionEditorP(c *gin.Context) {
	type DataInputStruct struct {
		ModeEditor                      byte
		CurrencyCode, SubDepartmentCode string `binding:"required"`
	}

	var DataInput DataInputStruct
	err := c.Bind(&DataInput)
	if err != nil {
		SendResponse(global_var.ResponseCode.InvalidDataFormat, "", nil, c)
	} else {
		// Get Program Configuration
		val, exist := c.Get("pConfig")
		if !exist {
			SendResponse(global_var.ResponseCode.DatabaseError, "pConfig", nil, c)
			return
		}
		pConfig := val.(*config.CompanyDataConfiguration)
		DB := pConfig.DB

		GeneralCodeName := GetAccountBySubDepartmentTransactionEditor(c, DB, DataInput.ModeEditor, DataInput.CurrencyCode, DataInput.SubDepartmentCode)
		SendResponse(global_var.ResponseCode.Successfully, "", GeneralCodeName, c)
	}
}

func GetAccountSubGroupByAccountCode1P(c *gin.Context) {
	var AccountCode string
	// Get Program Configuration
	val, exist := c.Get("pConfig")
	if !exist {
		SendResponse(global_var.ResponseCode.DatabaseError, "pConfig", nil, c)
		return
	}
	pConfig := val.(*config.CompanyDataConfiguration)
	DB := pConfig.DB

	err := c.Bind(&AccountCode)
	if err != nil {
		SendResponse(global_var.ResponseCode.InvalidDataFormat, "", nil, c)
	} else {
		var SubGroupCode string
		DB.Table(db_var.TableName.CfgInitAccount).Select("sub_group_code").Where("code = ?", AccountCode).Find(&SubGroupCode)
		SendResponse(global_var.ResponseCode.Successfully, "", SubGroupCode, c)
	}
}

func GetReservationDepositCardComboListP(c *gin.Context) {
	var ReservationDepositCardComboList db_var.ReservationDepositCardComboListStruct
	// Get Program Configuration
	val, exist := c.Get("pConfig")
	if !exist {
		SendResponse(global_var.ResponseCode.DatabaseError, "pConfig", nil, c)
		return
	}
	pConfig := val.(*config.CompanyDataConfiguration)
	DB := pConfig.DB

	GeneralCodeName := GetGeneralCodeName(DB, db_var.TableName.CfgInitCurrency, "")
	ReservationDepositCardComboList.Currency = GeneralCodeName
	GeneralCodeName = GetGeneralCodeName(DB, db_var.TableName.SubFolioGroup, "")
	ReservationDepositCardComboList.SubFolioGroup = GeneralCodeName
	GeneralCodeName = GetGeneralCodeName(DB, db_var.TableName.CfgInitSubDepartment, "")
	ReservationDepositCardComboList.SubDepartment = GeneralCodeName
	GeneralCodeName = GetGeneralCodeName(DB, db_var.TableName.CfgInitCardBank, "")
	ReservationDepositCardComboList.CardBank = GeneralCodeName
	GeneralCodeName = GetGeneralCodeName(DB, db_var.TableName.CfgInitCardType, "")
	ReservationDepositCardComboList.CardType = GeneralCodeName

	SendResponse(global_var.ResponseCode.Successfully, "", ReservationDepositCardComboList, c)
}

func GetInvCostingMethod(c *gin.Context) (CostingMethod string) {

	// Get Program Configuration
	val, exist := c.Get("pConfig")
	if !exist {
		SendResponse(global_var.ResponseCode.DatabaseError, "pConfig", nil, c)
		return
	}
	pConfig := val.(*config.CompanyDataConfiguration)
	CostingMethod = pConfig.Dataset.Configuration[global_var.ConfigurationCategory.Inventory][global_var.ConfigurationName.CostingMethod].(string)
	return
}

func GetBusinessSourceCommissionRateP(c *gin.Context) {
	type DataStruct struct {
		BusinessSourceCode string `binding:"required"`
		RoomRateCode       string `binding:"required"`
	}
	var DataInput DataStruct
	err := c.BindQuery(&DataInput)
	if err != nil {
		//fmt.Println(err.Error())
		errMsg := general.GenerateValidateErrorMsg(c, err)
		SendResponse(global_var.ResponseCode.InvalidDataFormat, errMsg, nil, c)
		return
	}

	// Get Program Configuration
	val, exist := c.Get("pConfig")
	if !exist {
		SendResponse(global_var.ResponseCode.DatabaseError, "pConfig", nil, c)
		return
	}
	pConfig := val.(*config.CompanyDataConfiguration)
	DB := pConfig.DB

	var DataOutput map[string]interface{}
	if err := DB.Table(db_var.TableName.CfgInitRoomRateBusinessSource).Select("commission_type_code, commission_value").
		Where("room_rate_code=?", DataInput.RoomRateCode).
		Where("company_code=?", DataInput.BusinessSourceCode).
		Limit(1).
		Scan(&DataOutput).Error; err != nil {
		SendResponse(global_var.ResponseCode.DatabaseError, "", "", c)
		return
	}

	SendResponse(global_var.ResponseCode.Successfully, "", DataOutput, c)
}
