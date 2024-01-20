package global_query

import (
	"chs_cloud_general/internal/config"
	DBVar "chs_cloud_general/internal/db_var"
	"chs_cloud_general/internal/global_var"
	GlobalVar "chs_cloud_general/internal/global_var"
	MasterData "chs_cloud_general/internal/master_data"
	"chs_cloud_general/internal/utils/cache"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func GetCashRemittanceP(c *gin.Context) {
	// Get Program Configuration
	val, exist := c.Get("pConfig")
	if !exist {
		MasterData.SendResponse(global_var.ResponseCode.DatabaseError, "pConfig", nil, c)
		return
	}
	v, e := c.Get("UserInfo")
	if !e {
		MasterData.SendResponse(global_var.ResponseCode.InternalServerError, "UserInfo", nil, c)
		return
	}
	UserInfo := v.(global_var.TUserInfo)
	pConfig := val.(*config.CompanyDataConfiguration)
	DB := pConfig.DB

	var PaymentOutput []map[string]interface{}
	if err := DB.Table("((?)UNION ALL(?)) AS SubFolio",
		DB.Table(DBVar.TableName.SubFolio).Select(
			" sub_folio.account_code AS AccountCode,"+
				" cfg_init_account.name AS Account,"+
				" SUM(IF(sub_folio.type_code='"+GlobalVar.TransactionType.Credit+"', (sub_folio.quantity*sub_folio.amount), -(sub_folio.quantity*sub_folio.amount))) AS TotalAmount ").
			Joins("LEFT OUTER JOIN cfg_init_account ON (sub_folio.account_code = cfg_init_account.code)").
			Joins("LEFT OUTER JOIN cfg_init_account_sub_group ON (cfg_init_account.sub_group_code = cfg_init_account_sub_group.code)").
			Where("sub_folio.log_shift_id=?", UserInfo.LogShiftID).
			Where("sub_folio.void='0'").
			Where("cfg_init_account_sub_group.group_code=?", GlobalVar.GlobalAccountGroup.Payment).
			Group("sub_folio.breakdown1"),
		DB.Table(DBVar.TableName.GuestDeposit).Select(
			" guest_deposit.account_code AS AccountCode,"+
				" cfg_init_account.name AS Account,"+
				" SUM(IF(guest_deposit.type_code='"+GlobalVar.TransactionType.Credit+"', guest_deposit.amount, -guest_deposit.amount)) AS TotalAmount ").
			Joins("LEFT OUTER JOIN cfg_init_account ON (guest_deposit.account_code = cfg_init_account.code)").
			Joins("LEFT OUTER JOIN cfg_init_account_sub_group ON (cfg_init_account.sub_group_code = cfg_init_account_sub_group.code)").
			Where("guest_deposit.log_shift_id=?", UserInfo.LogShiftID).
			Where("guest_deposit.void='0'").
			Where("cfg_init_account_sub_group.group_code=?", GlobalVar.GlobalAccountGroup.Payment).
			Group("guest_deposit.id")).Select(
		" SubFolio.AccountCode," +
			" SubFolio.Account," +
			" SUM(SubFolio.TotalAmount) AS TotalAmount").
		Where("SubFolio.TotalAmount>0").
		Group("SubFolio.AccountCode").
		Order("SubFolio.AccountCode").Scan(&PaymentOutput).Error; err != nil {
		MasterData.SendResponse(GlobalVar.ResponseCode.DatabaseError, "", nil, c)
		return
	}

	if PaymentOutput == nil {
		PaymentOutput = make([]map[string]interface{}, 0)
	}

	var RefundOutput []map[string]interface{}
	if err := DB.Table("((?)UNION ALL(?)) AS SubFolio",
		DB.Table(DBVar.TableName.SubFolio).Select(
			" sub_folio.account_code AS AccountCode,"+
				" cfg_init_account.name AS Account,"+
				" SUM(IF(sub_folio.type_code='"+GlobalVar.TransactionType.Credit+"', (sub_folio.quantity*sub_folio.amount), -(sub_folio.quantity*sub_folio.amount))) AS TotalAmount ").
			Joins("LEFT OUTER JOIN cfg_init_account ON (sub_folio.account_code = cfg_init_account.code)").
			Joins("LEFT OUTER JOIN cfg_init_account_sub_group ON (cfg_init_account.sub_group_code = cfg_init_account_sub_group.code)").
			Where("sub_folio.log_shift_id=?", UserInfo.LogShiftID).
			Where("sub_folio.void='0'").
			Where("cfg_init_account_sub_group.group_code=?", GlobalVar.GlobalAccountGroup.Payment).
			Group("sub_folio.breakdown1"),
		DB.Table(DBVar.TableName.GuestDeposit).Select(
			" guest_deposit.account_code AS AccountCode,"+
				" cfg_init_account.name AS Account,"+
				" SUM(IF(guest_deposit.type_code='"+GlobalVar.TransactionType.Credit+"', guest_deposit.amount, -guest_deposit.amount)) AS TotalAmount ").
			Joins("LEFT OUTER JOIN cfg_init_account ON (guest_deposit.account_code = cfg_init_account.code)").
			Joins("LEFT OUTER JOIN cfg_init_account_sub_group ON (cfg_init_account.sub_group_code = cfg_init_account_sub_group.code)").
			Where("guest_deposit.log_shift_id=?", UserInfo.LogShiftID).
			Where("guest_deposit.void='0'").
			Where("cfg_init_account_sub_group.group_code=?", GlobalVar.GlobalAccountGroup.Payment).
			Group("guest_deposit.id")).
		Select(
			" SubFolio.AccountCode," +
				" SubFolio.Account," +
				" -SUM(SubFolio.TotalAmount) AS TotalAmount ").
		Where("SubFolio.TotalAmount<0").
		Group("SubFolio.AccountCode").
		Order("SubFolio.AccountCode").Scan(&RefundOutput).Error; err != nil {
		MasterData.SendResponse(GlobalVar.ResponseCode.DatabaseError, "", nil, c)
		return
	}
	if RefundOutput == nil {
		RefundOutput = make([]map[string]interface{}, 0)
	}

	var CashCountOutput []map[string]interface{}
	if err := DB.Table(DBVar.TableName.CfgInitCurrencyNominal).Select(
		" cfg_init_currency_nominal.currency_sign,"+
			" cfg_init_currency_nominal.nominal,"+
			" IFNULL(cash_count.id, 0) AS id,"+
			" 0 AS SubTotal,"+
			" IFNULL(cash_count.quantity, 0) AS quantity").
		Joins("LEFT OUTER JOIN (SELECT * FROM cash_count WHERE log_shift_id=?) AS cash_count ON (cfg_init_currency_nominal.currency_sign = cash_count.currency_sign AND cfg_init_currency_nominal.nominal = cash_count.nominal)", GlobalVar.UserInfo.LogShiftID).
		Where("cfg_init_currency_nominal.currency_sign='Rp'").
		Order("cfg_init_currency_nominal.id_sort").Scan(&CashCountOutput).Error; err != nil {
		MasterData.SendResponse(GlobalVar.ResponseCode.DatabaseError, "", nil, c)
		return
	}

	var CurrentBalanceOutput map[string]interface{}
	GACash := pConfig.Dataset.GlobalAccount.Cash
	if err := DB.Table("((?)UNION(?)) AS A",
		DB.Table(DBVar.TableName.SubFolio).Select(" SUM(IF(type_code='C', (quantity*amount), -(quantity*amount))) AS Balance").
			Where("log_shift_id=?", UserInfo.LogShiftID).
			Where("void='0'").
			Where("account_code=?", GACash).
			Group("account_code"),
		DB.Table(DBVar.TableName.GuestDeposit).Select(" SUM(IF(type_code='C', amount, -amount)) AS Balance").
			Where("log_shift_id=?", UserInfo.LogShiftID).
			Where("void='0'").
			Where("account_code=?", GACash).
			Group("account_code")).
		Select("SUM(Balance) AS Balance").
		Scan(&CurrentBalanceOutput).Error; err != nil {
		MasterData.SendResponse(GlobalVar.ResponseCode.DatabaseError, "", nil, c)
		return
	}

	var ShiftOutput map[string]interface{}
	if err := DB.Table(DBVar.TableName.LogShift).Where("id=?", UserInfo.LogShiftID).Limit(1).Scan(&ShiftOutput).Error; err != nil {
		MasterData.SendResponse(GlobalVar.ResponseCode.DatabaseError, "", nil, c)
		return
	}

	// Get the value of the cookie
	cookie, err := c.Cookie("sembarang")
	if err != nil {
		fmt.Println(cookie)
		// Handle cookie not found error
	}

	MasterData.SendResponse(GlobalVar.ResponseCode.Successfully, "", gin.H{
		"payment":           PaymentOutput,
		"refund":            RefundOutput,
		"cash_count":        CashCountOutput,
		"current_balance":   CurrentBalanceOutput,
		"shift_information": ShiftOutput,
	}, c)
}

func CloseShiftP(c *gin.Context) {

	// Get Program Configuration
	val, exist := c.Get("pConfig")
	if !exist {
		MasterData.SendResponse(global_var.ResponseCode.DatabaseError, "pConfig", nil, c)
		return
	}
	pConfig := val.(*config.CompanyDataConfiguration)
	DB := pConfig.DB
	v, ex := c.Get("UserInfo")
	if !ex {
		MasterData.SendResponse(global_var.ResponseCode.InternalServerError, "UserInfo", nil, c)
		return
	}
	UserInfo := v.(global_var.TUserInfo)
	UserID := c.GetString("ValidUserCode")

	if err := DB.Table(DBVar.TableName.LogShift).Where("id=?", UserInfo.LogShiftID).
		Updates(map[string]interface{}{
			"end_date": time.Now(),
			"is_open":  "0",
		}).Error; err != nil {
		MasterData.SendResponse(GlobalVar.ResponseCode.DatabaseError, "", nil, c)
		return
	}

	if err := cache.DataCache.Del(c, pConfig.CompanyCode, "USER_INFO_"+UserID); err != nil {
		MasterData.SendResponse(GlobalVar.ResponseCode.DatabaseError, "", nil, c)
		return
	}
	MasterData.SendResponse(GlobalVar.ResponseCode.Successfully, "", nil, c)
}

// func GetDetailListP(c *gin.Context) {
// 	Param := c.Param("ApArNumber")
// 	var DataOutput []map[string]interface{}
// 	PaymentTableName := DBVar.TableName.AccApArPayment
// 	PaymentDetailTableName := DBVar.TableName.AccApArPaymentDetail
// 	err := DB.Table(PaymentDetailTableName).Select(
// 		" "+PaymentDetailTableName+".ref_number,"+
// 			" "+PaymentDetailTableName+".amount,"+
// 			" "+PaymentTableName+".journal_account_code,"+
// 			" "+PaymentTableName+".`date`,"+
// 			" "+PaymentTableName+".remark,"+
// 			" "+PaymentDetailTableName+".user_id,"+
// 			" "+PaymentDetailTableName+".id,"+
// 			" acc_ap_ar_payment.create_ap_number,"+
// 			" acc_ap_ar_payment.source_code_ap_ar,"+
// 			" acc_ap_ar_payment.is_payment_ap_ar,"+
// 			" CONCAT("+PaymentTableName+".journal_account_code, ' - ', cfg_init_journal_account.name) AS JournalAccount ").
// 		Joins("LEFT OUTER JOIN "+PaymentTableName+" ON ("+PaymentDetailTableName+".ref_number = "+PaymentTableName+".ref_number)").
// 		Joins("LEFT OUTER JOIN cfg_init_journal_account ON ("+PaymentTableName+".journal_account_code = cfg_init_journal_account.code)").
// 		Where(PaymentDetailTableName+".ap_ar_number=?", Param).
// 		Scan(&DataOutput).Error

// 	if err == nil {
// 		MasterData.SendResponse(GlobalVar.ResponseCode.Successfully, "", DataOutput, c)
// 	} else {
// 		MasterData.SendResponse(GlobalVar.ResponseCode.DatabaseError, "", nil, c)
// 	}
// }

// func GetComboList(c *gin.Context) {
// 	GeneralCodeNameArray := make(map[string]interface{})
// 	var DataOutput []map[string]interface{}
// 	GeneralCodeNameArray["SubDepartment"] = MasterData.GetGeneralCodeName(DB, DBVar.TableName.CfgInitSubDepartment, "")
// 	GeneralCodeNameArray["Company"] = MasterData.GetGeneralCodeName(DB, DBVar.TableName.Company, "")
// 	Query := DB.Table(DBVar.TableName.AccCfgInitBankAccount).Select("bank_name,code, journal_account_code, name,bank_account_number ").
// 		Where("type_code=? AND for_payment=1", GlobalVar.BankAccountType.CashAccount).
// 		Or("type_code=? AND for_payment=1", GlobalVar.BankAccountType.SavingAccount)

// 	Query.Order("journal_account_code").Scan(&DataOutput)

// 	GeneralCodeNameArray["BankAccount"] = DataOutput

// 	MasterData.SendResponse(GlobalVar.ResponseCode.Successfully, "", GeneralCodeNameArray, c)
// }

// // Get Data Edit
// func GetP(c *gin.Context) {
// 	Param := c.Param("Param")
// 	var DataOutput map[string]interface{}
// 	Query := DB.Table(DBVar.TableName.LostAndFound).Where("id=?", Param).Take(&DataOutput)

// 	if Query.Error == nil {
// 		MasterData.SendResponse(GlobalVar.ResponseCode.Successfully, "", DataOutput, c)
// 	} else {
// 		MasterData.SendResponse(GlobalVar.ResponseCode.DataNotFound, "", nil, c)
// 	}
// }

// func InsertP(c *gin.Context) {
// 	Process(c, false, false)
// }

// func UpdateP(c *gin.Context) {
// 	Process(c, true, false)
// }

// func Process(c *gin.Context, IsUpdate, IsReceive bool) {
// 	type DataInputDetailStruct struct {
// 		Amount            float64 `json:"amount" binding:"required"`
// 		AccountCode       string  `json:"account_code" binding:"required"`
// 		Remark            string  `json:"remark"`
// 		SubDepartmentCode string  `json:"sub_department_code" binding:"required"`
// 	}
// 	type DataInputStruct struct {
// 		RefNumber         string                  `json:"ref_number"`
// 		Id                uint64                  `json:"id"`
// 		CompanyCode       string                  `json:"company_code" binding:"required"`
// 		SubDepartmentCode string                  `json:"sub_department_code" binding:"required"`
// 		BankAccountCode   string                  `json:"bank_account_code" binding:"required"`
// 		Remark            string                  `json:"remark"`
// 		Date              time.Time               `json:"date" binding:"required"`
// 		ItemDetails       []DataInputDetailStruct `json:"item_details" binding:"required"`
// 	}
// 	var DataInput DataInputStruct
// 	err := c.BindJSON(&DataInput)
// 	if err != nil || (DataInput.RefNumber == "" && IsUpdate) {
// 		var errMsg interface{}
// 		errMsg = "RefNumber is required"
// 		if err != nil {
// 			//fmt.Println(err.Error())
// 			errMsg = General.GenerateValidateErrorMsg(c, err)
// 		}
// 		MasterData.SendResponse(GlobalVar.ResponseCode.InvalidDataFormat, errMsg, nil, c)
// 		return
// 	}

// 	//TODO Add Validation out of date (check desktop version)
// 	UserID := c.GetString("ValidUserCode")
// 	ReceivePaymentNumber := ""
// 	RefNumber := ""
// 	UnitCode := MasterData.GetUnitCode()
// 	err = DB.Transaction(func(tx *gorm.DB) error {
// 		if IsUpdate {
// 			RefNumber = DataInput.RefNumber
// 			var OldDateJournalDetail time.Time
// 			tx.Table(DBVar.TableName.AccJournal).Select("date").Where("ref_number=?", DataInput.RefNumber).Take(&OldDateJournalDetail)
// 			tx.Exec("CALL delete_acc_journal_detail(?,?,?)", RefNumber, OldDateJournalDetail, UserID)
// 			err :=global_query.UpdateAccJournal(tx, RefNumber, DataInput.CompanyCode, DataInput.Date, DataInput.Remark, UserID)
// 			if err != nil {
// 				return err
// 			}
// 		} else {
// 			if IsReceive {
// 				RefNumber = global_query.GetJournalRefNumber(GlobalVar.JournalPrefix.Receive, DataInput.Date)
// 				ReceivePaymentNumber = GetReceiveNumber(DataInput.Date)
// 			} else {
// 				RefNumber = global_query.GetJournalRefNumber(GlobalVar.JournalPrefix.Disbursement, DataInput.Date)
// 				ReceivePaymentNumber = GetPaymentNumber(DataInput.Date)
// 			}
// 			err := global_query.InsertAccJournal(tx, RefNumber, UnitCode, ReceivePaymentNumber, DataInput.CompanyCode, GlobalVar.JournalType.Other, GlobalVar.JournalGroup.Other, DataInput.Date, DataInput.Remark, "", 0, UserID)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 		//Master
// 		var TransactionType1, TransactionType2 string
// 		if IsReceive {
// 			TransactionType1 = GlobalVar.TransactionType.Debit
// 			TransactionType2 = GlobalVar.TransactionType.Credit
// 		} else {
// 			TransactionType1 = GlobalVar.TransactionType.Credit
// 			TransactionType2 = GlobalVar.TransactionType.Debit
// 		}

// 		//Detail Credit
// 		TotalAmount := 0.00
// 		for _, detail := range DataInput.ItemDetails {
// 			err := global_query.InsertAccJournalDetail(tx, RefNumber, UnitCode, detail.SubDepartmentCode, detail.AccountCode, detail.Amount, TransactionType2, detail.Remark, "", UserID)
// 			if err != nil {
// 				return err
// 			}

// 			TotalAmount += detail.Amount
// 		}
// 		err = global_query.InsertAccJournalDetail(tx, RefNumber, UnitCode, DataInput.SubDepartmentCode, DataInput.BankAccountCode, TotalAmount, TransactionType1, DataInput.Remark, "", UserID)
// 		if err != nil {
// 			return err
// 		}
// 		if IsUpdate {
// 			if IsReceive {
// 				InsertLogUser(GlobalVar.SystemCode.Accounting, GlobalVar.LogUserAction.UpdateReceive,global_query.GetAuditDate(c,DB), "", "", "", RefNumber, "", "", "", UserID)
// 			} else {
// 				InsertLogUser(GlobalVar.SystemCode.Accounting, GlobalVar.LogUserAction.UpdateReceive,global_query.GetAuditDate(c,DB), "", "", "", RefNumber, "", "", "", UserID)
// 			}
// 		} else {
// 			if IsReceive {
// 				InsertLogUser(GlobalVar.SystemCode.Accounting, GlobalVar.LogUserAction.UpdateReceive,global_query.GetAuditDate(c,DB), "", "", "", RefNumber, "", "", "", UserID)
// 			} else {
// 				InsertLogUser(GlobalVar.SystemCode.Accounting, GlobalVar.LogUserAction.UpdateReceive,global_query.GetAuditDate(c,DB), "", "", "", RefNumber, "", "", "", UserID)
// 			}
// 		}

// 		return nil
// 	})

// 	if err == nil {
// 		MasterData.SendResponse(GlobalVar.ResponseCode.Successfully, "", nil, c)
// 	} else {
// 		MasterData.SendResponse(GlobalVar.ResponseCode.DatabaseError, "", nil, c)
// 	}
// }

// func DeleteP(c *gin.Context) {
// 	Param := c.Param("Number")
// 	UpdatedBy := c.GetString("ValidUserCode")
// 	RefNumber := ""
// 	err := DB.Transaction(func(tx *gorm.DB) error {
// 		type DataOutputStruct struct {
// 			Date      time.Time
// 			RefNumber string
// 		}
// 		var DataOutput DataOutputStruct
// 		if err := tx.Table(DBVar.TableName.AccApAr).Select("date,ref_number").Where("number=?", Param).Take(&DataOutput).Error; err != nil {
// 			return err
// 		}
// 		RefNumber = DataOutput.RefNumber
// 		err := DB.Table(DBVar.TableName.AccJournal).Where("ref_number=? AND date=?", Param, DataOutput.Date).Updates(map[string]interface{}{"updated_by": UpdatedBy}).Error
// 		if err != nil {
// 			return err
// 		}
// 		err = DB.Table(DBVar.TableName.AccJournal).Where("ref_number=? AND date=?", Param, DataOutput.Date).Delete(&Param).Error
// 		if err != nil {
// 			return err
// 		}
// 		err = DB.Table(DBVar.TableName.AccJournalDetail).Where("ref_number=? AND date=?", Param, DataOutput.Date).Updates(map[string]interface{}{"updated_by": UpdatedBy}).Error
// 		if err != nil {
// 			return err
// 		}
// 		err = DB.Table(DBVar.TableName.AccJournalDetail).Where("ref_number=? AND date=?", Param, DataOutput.Date).Delete(&Param).Error
// 		if err != nil {
// 			return err
// 		}

// 		err = tx.Exec("CALL delete_acc_ap_ar(?,?)", Param, UpdatedBy).Error
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// 	InsertLogUser(GlobalVar.SystemCode.Accounting, GlobalVar.LogUserActionCAS.DeleteAccountReceivable,global_query.GetAuditDate(c,DB), "", "", "", Param, RefNumber, "", "", UpdatedBy)
// 	if err == nil {
// 		MasterData.SendResponse(GlobalVar.ResponseCode.Successfully, "", nil, c)
// 	} else {
// 		MasterData.SendResponse(GlobalVar.ResponseCode.DatabaseError, "", nil, c)
// 	}
// }
