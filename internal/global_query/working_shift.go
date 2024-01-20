package global_query

import (
	"chs/internal/config"
	"chs/internal/db_var"
	General "chs/internal/general"
	"chs/internal/global_var"
	GlobalVar "chs/internal/global_var"
	"chs/internal/master_data"
	MasterData "chs/internal/master_data"
	"chs/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func GetShiftInformationP(c *gin.Context) {
	UserID := c.Param("UserID")
	// Get Program Configuration
	val, exist := c.Get("pConfig")
	if !exist {
		MasterData.SendResponse(global_var.ResponseCode.DatabaseError, "pConfig", nil, c)
		return
	}
	pConfig := val.(*config.CompanyDataConfiguration)
	DB := pConfig.DB
	var DataOutput db_var.Log_shift
	if err := DB.Table(db_var.TableName.LogShift).Where("created_by=?", UserID).Where("is_open=1").Limit(1).Scan(&DataOutput).Error; err != nil {
		master_data.SendResponse(global_var.ResponseCode.DatabaseError, "", nil, c)
		return
	}
	master_data.SendResponse(global_var.ResponseCode.Successfully, "", DataOutput, c)

}

func GetWorkingShiftP(c *gin.Context) {
	// Get Program Configuration
	val, exist := c.Get("pConfig")
	if !exist {
		MasterData.SendResponse(global_var.ResponseCode.DatabaseError, "pConfig", nil, c)
		return
	}
	pConfig := val.(*config.CompanyDataConfiguration)
	DB := pConfig.DB
	var DataOutput []map[string]interface{}
	if err := DB.Table(db_var.TableName.WorkingShift).Scan(&DataOutput).Error; err != nil {
		master_data.SendResponse(global_var.ResponseCode.DatabaseError, "", nil, c)
		return
	}
	master_data.SendResponse(global_var.ResponseCode.Successfully, "", DataOutput, c)
}

func SetOpeningBalanceP(c *gin.Context) {
	ctx, span := global_var.Tracer.Start(c, "MoveSubFolioP")
	defer span.End()

	type DataInputStruct struct {
		OpeningBalance float64 `json:"opening_balance"`
	}
	// Get Program Configuration
	val, exist := c.Get("pConfig")
	if !exist {
		MasterData.SendResponse(global_var.ResponseCode.DatabaseError, "pConfig", nil, c)
		return
	}
	pConfig := val.(*config.CompanyDataConfiguration)
	DB := pConfig.DB
	var DataInput DataInputStruct
	err := c.BindJSON(&DataInput)
	if err != nil {
		errMsg := General.GenerateValidateErrorMsg(c, err)
		MasterData.SendResponse(GlobalVar.ResponseCode.InvalidDataFormat, errMsg, nil, c)
		return
	}

	UserID := c.GetString("ValidUserCode")

	// Get User Info
	UserInfo, err := loadUserInfo(c, UserID, pConfig.CompanyCode)
	if err != nil {
		utils.LogResponseError(c, ctx, global_var.AppLogger, errors.Wrap(err, "UpdateUserGroupAccessP.loadUserInfo"))
		master_data.SendResponse(global_var.ResponseCode.InternalServerError, "Failed load User", nil, c)
		return
	}

	if DataInput.OpeningBalance > 0 {
		if err := DB.Table(db_var.TableName.LogShift).Where("id=?", UserInfo.LogShiftID).
			Update("opening_balance", DataInput.OpeningBalance).Error; err != nil {
			master_data.SendResponse(global_var.ResponseCode.DatabaseError, "", nil, c)
			return
		}
	}
	master_data.SendResponse(global_var.ResponseCode.Successfully, "", nil, c)
}
