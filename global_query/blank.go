package global_query

// func GetList(c *gin.Context) {
// 	type QueryParamStruct struct {
// 		Index              int
// 		Text               string
// 		StartDate, EndDate time.Time
// 		IsAP               bool
// 		IsPaid             int
// 	}

// 	var QueryParam QueryParamStruct
// 	err := c.BindQuery(&QueryParam)

// 	if err != nil {
// 		MasterData.SendResponse(GlobalVar.ResponseCode.InvalidDataFormat, "", nil, c)
// 		return
// 	}

// 	MainTableName := DBVar.TableName.AccApAr
// 	QueryCondition := ""

// 	if QueryParam.Text != "" {
// 		switch QueryParam.Index {
// 		case 0:
// 			QueryCondition = MainTableName + ".number LIKE ?"
// 		case 1:
// 			QueryCondition = MainTableName + ".document_number LIKE ?"
// 		case 2:
// 			QueryCondition = MainTableName + ".ref_number LIKE ?"
// 		case 3:
// 			QueryCondition = "company.name LIKE ?"
// 		case 4:
// 			QueryCondition = "CONCAT(" + MainTableName + ".journal_account_debit, ' - ', cfg_init_journal_account.name) LIKE ?"
// 		case 5:
// 			QueryCondition = "CONCAT(" + MainTableName + ".journal_account_credit, ' - ', cfg_init_journal_account1.name) LIKE ?"
// 		case 6:
// 			QueryCondition = MainTableName + ".remark LIKE ?"
// 		case 7:
// 			QueryCondition = MainTableName + ".created_by LIKE ?"
// 		case 8:
// 			QueryCondition = MainTableName + ".updated_by LIKE ?"
// 		}
// 	}

// 	var DataOutput []map[string]interface{}
// 	Query := DB.Table(MainTableName).Select(
// 		" " + MainTableName + ".*," +
// 			" (" + MainTableName + ".amount-" + MainTableName + ".amount_paid) AS Outstanding," +
// 			" company.name As Company," +
// 			" CONCAT(" + MainTableName + ".journal_account_debit, ' - ', cfg_init_journal_account.name) AS JournalAccountDebit," +
// 			" CONCAT(" + MainTableName + ".journal_account_credit, ' - ', cfg_init_journal_account1.name) AS JournalAccountCredit ").
// 		Joins("LEFT OUTER JOIN company ON (" + MainTableName + ".company_code = company.code)").
// 		Joins("LEFT OUTER JOIN cfg_init_journal_account ON (" + MainTableName + ".journal_account_debit = cfg_init_journal_account.code)").
// 		Joins("LEFT OUTER JOIN cfg_init_journal_account cfg_init_journal_account1 ON (" + MainTableName + ".journal_account_credit = cfg_init_journal_account1.code)")

// 	if QueryCondition != "" {
// 		Query.Where(QueryCondition, "%"+QueryParam.Text+"%")
// 	}
// 	if !QueryParam.StartDate.IsZero() && !QueryParam.StartDate.IsZero() && General.DateOf(QueryParam.StartDate).Unix() <=General.DateOf(QueryParam.EndDate).Unix() {
// 		Query.Where(MainTableName+".date BETWEEN ? AND ?", QueryParam.StartDate, QueryParam.EndDate)
// 	}
// 	Query.Where(MainTableName+".is_ap=?", QueryParam.IsAP).Where(MainTableName+".is_paid=?", QueryParam.IsPaid).Order(MainTableName + ".date DESC").Order(MainTableName + ".number").Scan(&DataOutput)

// 	if Query.Error == nil {
// 		MasterData.SendResponse(GlobalVar.ResponseCode.Successfully, "", DataOutput, c)
// 	} else {
// 		MasterData.SendResponse(GlobalVar.ResponseCode.DatabaseError, "", nil, c)
// 	}
// }

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
// 			" "+PaymentDetailTableName+".id_log,"+
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
// 			err := UpdateAccJournal(tx, RefNumber, "", DataInput.CompanyCode, DataInput.Date, DataInput.Remark, UserID)
// 			if err != nil {
// 				return err
// 			}
// 		} else {
// 			if IsReceive {
// 				RefNumber = GetJournalRefNumber(GlobalVar.JournalPrefix.Receive, DataInput.Date)
// 				ReceivePaymentNumber = GetReceiveNumber(DataInput.Date)
// 			} else {
// 				RefNumber = GetJournalRefNumber(GlobalVar.JournalPrefix.Disbursement, DataInput.Date)
// 				ReceivePaymentNumber = GetPaymentNumber(DataInput.Date)
// 			}
// 			err := InsertAccJournal(tx, RefNumber, "", UnitCode, ReceivePaymentNumber, DataInput.CompanyCode, GlobalVar.JournalType.Other, GlobalVar.JournalGroup.Other, DataInput.Date, DataInput.Remark, "", 0, UserID)
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
// 			err := InsertAccJournalDetail(tx, RefNumber, UnitCode, detail.SubDepartmentCode, detail.AccountCode, detail.Amount, TransactionType2, detail.Remark, "", UserID)
// 			if err != nil {
// 				return err
// 			}

// 			TotalAmount += detail.Amount
// 		}
// 		err = InsertAccJournalDetail(tx, RefNumber, UnitCode, DataInput.SubDepartmentCode, DataInput.BankAccountCode, TotalAmount, TransactionType1, DataInput.Remark, "", UserID)
// 		if err != nil {
// 			return err
// 		}
// 		if IsUpdate {
// 			if IsReceive {
// 				InsertLogUser(GlobalVar.SystemCode.Accounting, GlobalVar.LogUserAction.UpdateReceive,  GetAuditDate(c, DB, false), "", "", "", RefNumber, "", "", "", UserID)
// 			} else {
// 				InsertLogUser(GlobalVar.SystemCode.Accounting, GlobalVar.LogUserAction.UpdateReceive,  GetAuditDate(c, DB, false), "", "", "", RefNumber, "", "", "", UserID)
// 			}
// 		} else {
// 			if IsReceive {
// 				InsertLogUser(GlobalVar.SystemCode.Accounting, GlobalVar.LogUserAction.UpdateReceive,  GetAuditDate(c, DB, false), "", "", "", RefNumber, "", "", "", UserID)
// 			} else {
// 				InsertLogUser(GlobalVar.SystemCode.Accounting, GlobalVar.LogUserAction.UpdateReceive,  GetAuditDate(c, DB, false), "", "", "", RefNumber, "", "", "", UserID)
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
// 	InsertLogUser(GlobalVar.SystemCode.Accounting, GlobalVar.LogUserActionCAS.DeleteAccountReceivable,  GetAuditDate(c, DB, false), "", "", "", Param, RefNumber, "", "", UpdatedBy)
// 	if err == nil {
// 		MasterData.SendResponse(GlobalVar.ResponseCode.Successfully, "", nil, c)
// 	} else {
// 		MasterData.SendResponse(GlobalVar.ResponseCode.DatabaseError, "", nil, c)
// 	}
// }
