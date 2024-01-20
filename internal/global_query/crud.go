package global_query

import (
	"chs_cloud_general/internal/db_var"
	DBVar "chs_cloud_general/internal/db_var"
	General "chs_cloud_general/internal/general"
	"chs_cloud_general/internal/global_var"
	GlobalVar "chs_cloud_general/internal/global_var"
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func UpdateInvoiceRefNumber(DB *gorm.DB, InvoiceNumber, RefNumber string) error {
	err := DB.Table(DBVar.TableName.Invoice).Where("number=?", InvoiceNumber).Update("ref_number", RefNumber).Error
	return err
}

func InsertAccApRefundDepositPayment(DB *gorm.DB, RefNumber string, JournalAccountCode string, DiscountJournalAccountCode string, BaJournalAccountCode string, OeJournalAccountCode string, TotalAmount float64, Discount float64, BankAdministration float64, OtherExpense float64, Date time.Time, Remark string, PaidByApAr string, CreatedBy string) error {
	var AccApRefundDepositPayment = DBVar.Acc_ap_refund_deposit_payment{
		RefNumber:                  RefNumber,
		JournalAccountCode:         JournalAccountCode,
		DiscountJournalAccountCode: DiscountJournalAccountCode,
		BaJournalAccountCode:       BaJournalAccountCode,
		OeJournalAccountCode:       OeJournalAccountCode,
		TotalAmount:                TotalAmount,
		Discount:                   Discount,
		BankAdministration:         BankAdministration,
		OtherExpense:               OtherExpense,
		Date:                       Date,
		Remark:                     Remark,
		PaidByApAr:                 PaidByApAr,
		CreatedBy:                  CreatedBy,
	}
	result := DB.Table(DBVar.TableName.AccApRefundDepositPayment).Create(&AccApRefundDepositPayment)
	return result.Error
}

func InsertAccApRefundDepositPaymentDetail(DB *gorm.DB, SubFolioId uint64, RefNumber string, Amount float64, CreatedBy string) error {
	var AccApRefundDepositPaymentDetail = DBVar.Acc_ap_refund_deposit_payment_detail{
		SubFolioId: SubFolioId,
		RefNumber:  RefNumber,
		Amount:     Amount,
		CreatedBy:  CreatedBy,
	}
	err := DB.Table(DBVar.TableName.AccApRefundDepositPaymentDetail).Create(&AccApRefundDepositPaymentDetail).Error

	return err
}

func InsertAccJournal(DB *gorm.DB, RefNumber, DocumentNumber string, UnitCode string, ApArNumber string, CompanyCode string, TypeCode string, GroupCode string,
	Date time.Time, Memo string, ChequeNumber string, IdSort int, CreatedBy string) error {
	if !(IdSort > 0) {
		IdSort = GetJournalIDSort(DB, Date)
	}
	var AccJournal = DBVar.Acc_journal{
		RefNumber:      RefNumber,
		DocumentNumber: DocumentNumber,
		UnitCode:       UnitCode,
		ApArNumber:     ApArNumber,
		CompanyCode:    CompanyCode,
		TypeCode:       TypeCode,
		GroupCode:      GroupCode,
		Date:           Date,
		DateUnixx:      Date.Unix(),
		Memo:           Memo,
		ChequeNumber:   ChequeNumber,
		IdSort:         IdSort,
		CreatedBy:      CreatedBy,
	}
	err := DB.Table(DBVar.TableName.AccJournal).Create(&AccJournal).Error
	return err
}

func UpdateAccJournal(DB *gorm.DB, RefNumber, DocumentNumber string, CompanyCode string, Date time.Time, Memo string, UpdatedBy string) error {
	var AccJournal = DBVar.Acc_journal{
		RefNumber:      RefNumber,
		DocumentNumber: DocumentNumber,
		CompanyCode:    CompanyCode,
		Date:           Date,
		DateUnixx:      Date.Unix(),
		Memo:           Memo,
		UpdatedBy:      UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.AccJournal).Omit("created_at", "created_by", "updated_at", "id").Updates(&AccJournal)
	return result.Error
}

func InsertAccJournalDetail(DB *gorm.DB, RefNumber string, UnitCode string, SubDepartmentCode string, AccountCode string, Amount float64, TypeCode string, Remark string, IdData string, IsTwoDigitDecimal bool, CreatedBy string) (Error error) {
	var Date time.Time
	if IsTwoDigitDecimal {
		Amount = General.RoundToX2(Amount)
	} else {
		Amount = General.RoundToX3(Amount)
	}
	if err := DB.Table(DBVar.TableName.AccJournal).Select("date").Where("ref_number=?", RefNumber).Take(&Date).Error; err != nil {
		return err
	}
	var AccJournalDetail = DBVar.Acc_journal_detail{
		RefNumber:         RefNumber,
		Date:              Date,
		UnitCode:          UnitCode,
		SubDepartmentCode: SubDepartmentCode,
		AccountCode:       AccountCode,
		Amount:            Amount,
		TypeCode:          TypeCode,
		Remark:            Remark,
		IdData:            IdData,
		CreatedBy:         CreatedBy,
	}
	err := DB.Table(DBVar.TableName.AccJournalDetail).Create(&AccJournalDetail).Error
	return err
}

func UpdateAccJournalDetail(DB *gorm.DB, RefNumber string, Date time.Time, UnitCode string, SubDepartmentCode string, AccountCode string, Amount float64, TypeCode string, Remark string, IdData string, UpdatedBy string, IdHolding uint64) error {
	var AccJournalDetail = DBVar.Acc_journal_detail{
		RefNumber:         RefNumber,
		Date:              Date,
		UnitCode:          UnitCode,
		SubDepartmentCode: SubDepartmentCode,
		AccountCode:       AccountCode,
		Amount:            Amount,
		TypeCode:          TypeCode,
		Remark:            Remark,
		IdData:            IdData,
		UpdatedBy:         UpdatedBy,
		IdHolding:         IdHolding,
	}
	result := DB.Table(DBVar.TableName.AccJournalDetail).Omit("created_at", "created_by", "updated_at", "id").Updates(&AccJournalDetail)
	return result.Error
}

func InsertLogUser(DB *gorm.DB, SystemCode string, ActionId int, AuditDate time.Time, IpAddress string, ComputerName string, MacAddress string, DataLink1 interface{}, DataLink2 interface{}, DataLink3 interface{}, Remark interface{}, CreatedBy string) error {
	result := DB.Table(DBVar.TableName.LogUser).Omit("actual_date").Create(map[string]interface{}{
		"system_code":   SystemCode,
		"action_id":     ActionId,
		"audit_date":    AuditDate,
		"ip_address":    IpAddress,
		"computer_name": ComputerName,
		"mac_address":   MacAddress,
		"data_link1":    DataLink1,
		"data_link2":    DataLink2,
		"data_link3":    DataLink3,
		"remark":        Remark,
		"created_by":    CreatedBy,
	})
	return result.Error
}
func InsertAccCashSaleRecon(DB *gorm.DB, JournalAccountCode string, JournalAccountCodeShortOver string, RefNumber string, Date time.Time, DateRecon time.Time, Amount float64, AmountShortOver float64, AmountDetail float64, Remark string, ReconBy string, IsOver uint8, CreatedBy string) error {
	var AccCashSaleRecon = DBVar.Acc_cash_sale_recon{
		JournalAccountCode:          JournalAccountCode,
		JournalAccountCodeShortOver: JournalAccountCodeShortOver,
		RefNumber:                   RefNumber,
		Date:                        Date,
		DateRecon:                   DateRecon,
		Amount:                      Amount,
		AmountShortOver:             AmountShortOver,
		AmountDetail:                AmountDetail,
		Remark:                      Remark,
		ReconBy:                     ReconBy,
		IsOver:                      &IsOver,
		CreatedBy:                   CreatedBy,
	}
	result := DB.Table(DBVar.TableName.AccCashSaleRecon).Create(&AccCashSaleRecon)
	return result.Error
}

func UpdateAccCashSaleRecon(DB *gorm.DB, Id uint64, JournalAccountCode string, JournalAccountCodeShortOver string, RefNumber string, DateRecon time.Time, Amount float64, AmountShortOver float64, AmountDetail float64, Remark string, ReconBy string, IsOver uint8, UpdatedBy string) error {
	var AccCashSaleRecon = DBVar.Acc_cash_sale_recon{
		Id:                          Id,
		JournalAccountCode:          JournalAccountCode,
		JournalAccountCodeShortOver: JournalAccountCodeShortOver,
		RefNumber:                   RefNumber,
		DateRecon:                   DateRecon,
		Amount:                      Amount,
		AmountShortOver:             AmountShortOver,
		AmountDetail:                AmountDetail,
		Remark:                      Remark,
		ReconBy:                     ReconBy,
		IsOver:                      &IsOver,
		UpdatedBy:                   UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.AccCashSaleRecon).Omit("created_at", "created_by", "updated_at", "id").Updates(&AccCashSaleRecon)
	return result.Error
}

func InsertFolioRouting(DB *gorm.DB, FolioNumber uint64, AccountCode string, FolioTransfer uint64, SubFolioTransfer string, CreatedBy string) error {
	var FolioRouting = DBVar.Folio_routing{
		FolioNumber:      FolioNumber,
		AccountCode:      AccountCode,
		FolioTransfer:    FolioTransfer,
		SubFolioTransfer: SubFolioTransfer,
		CreatedBy:        CreatedBy,
	}
	result := DB.Table(DBVar.TableName.FolioRouting).Create(&FolioRouting)
	return result.Error
}

func UpdateFolioRouting(DB *gorm.DB, FolioNumber uint64, AccountCode string, FolioTransfer uint64, SubFolioTransfer string, UpdatedBy string) error {
	var FolioRouting = DBVar.Folio_routing{
		FolioNumber:      FolioNumber,
		AccountCode:      AccountCode,
		FolioTransfer:    FolioTransfer,
		SubFolioTransfer: SubFolioTransfer,
		UpdatedBy:        UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.FolioRouting).Omit("created_at", "created_by", "updated_at", "id").Updates(&FolioRouting)
	return result.Error
}

// PROCEDURE =================================================
func DeleteFolioRouting(DB *gorm.DB, FolioNumber uint64, UserID string) error {
	err := DB.Exec("CALL delete_folio_routing(?,?)", FolioNumber, UserID).Error

	return err
}

func InsertInvoice(DB *gorm.DB, Number string, CompanyCode string, ContactPersonId uint64, IssuedDate time.Time, DueDate time.Time, Remark string, IsPaid uint8, RefNumber string, PrintCount int, CreatedBy string) error {
	var Invoice = DBVar.Invoice{
		Number:          Number,
		CompanyCode:     CompanyCode,
		ContactPersonId: ContactPersonId,
		IssuedDate:      IssuedDate,
		DueDate:         DueDate,
		Remark:          Remark,
		IsPaid:          IsPaid,
		RefNumber:       RefNumber,
		PrintCount:      PrintCount,
		CreatedBy:       CreatedBy,
	}
	result := DB.Table(DBVar.TableName.Invoice).Create(&Invoice)
	return result.Error
}

func UpdateInvoice(DB *gorm.DB, Number string, CompanyCode string, ContactPersonId uint64, IssuedDate time.Time, DueDate time.Time, Remark string, IsPaid uint8, RefNumber string, PrintCount int, UpdatedBy string) error {
	var Invoice = DBVar.Invoice{
		Number:          Number,
		CompanyCode:     CompanyCode,
		ContactPersonId: ContactPersonId,
		IssuedDate:      IssuedDate,
		DueDate:         DueDate,
		Remark:          Remark,
		IsPaid:          IsPaid,
		RefNumber:       RefNumber,
		PrintCount:      PrintCount,
		UpdatedBy:       UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.Invoice).Where("number", Number).Omit("created_at", "created_by", "id").Updates(&Invoice)
	return result.Error
}

func InsertContactPerson(DB *gorm.DB, TitleCode string, FullName string, Street string, CountryCode string, StateCode string,
	CityCode string, City string, NationalityCode string, PostalCode string, Phone1 string, Phone2 string, Fax string,
	Email string, Website string, CompanyCode string, GuestTypeCode string, IdCardCode string, IdCardNumber string,
	IsMale uint8, BirthPlace string, BirthDate time.Time, TypeCode string, CustomField01 string, CustomField02 string,
	CustomField03 string, CustomField04 string, CustomField05 string, CustomField06 string, CustomField07 string,
	CustomField08 string, CustomField09 string, CustomField10 string, CustomField11 string, CustomField12 string,
	CustomLookupFieldCode01 string, CustomLookupFieldCode02 string, CustomLookupFieldCode03 string, CustomLookupFieldCode04 string,
	CustomLookupFieldCode05 string, CustomLookupFieldCode06 string, CustomLookupFieldCode07 string, CustomLookupFieldCode08 string,
	CustomLookupFieldCode09 string, CustomLookupFieldCode10 string, CustomLookupFieldCode11 string, CustomLookupFieldCode12 string,
	CreatedBy string) (ID uint64, Error error) {
	var ContactPerson = DBVar.Contact_person{
		TitleCode:               &TitleCode,
		FullName:                &FullName,
		Street:                  &Street,
		CountryCode:             &CountryCode,
		StateCode:               &StateCode,
		CityCode:                &CityCode,
		City:                    &City,
		NationalityCode:         &NationalityCode,
		PostalCode:              &PostalCode,
		Phone1:                  &Phone1,
		Phone2:                  &Phone2,
		Fax:                     &Fax,
		Email:                   &Email,
		Website:                 &Website,
		CompanyCode:             &CompanyCode,
		GuestTypeCode:           &GuestTypeCode,
		IdCardCode:              &IdCardCode,
		IdCardNumber:            &IdCardNumber,
		IsMale:                  &IsMale,
		BirthPlace:              &BirthPlace,
		BirthDate:               &BirthDate,
		TypeCode:                TypeCode,
		CustomField01:           &CustomField01,
		CustomField02:           &CustomField02,
		CustomField03:           &CustomField03,
		CustomField04:           &CustomField04,
		CustomField05:           &CustomField05,
		CustomField06:           &CustomField06,
		CustomField07:           &CustomField07,
		CustomField08:           &CustomField08,
		CustomField09:           &CustomField09,
		CustomField10:           &CustomField10,
		CustomField11:           &CustomField11,
		CustomField12:           &CustomField12,
		CustomLookupFieldCode01: &CustomLookupFieldCode01,
		CustomLookupFieldCode02: &CustomLookupFieldCode02,
		CustomLookupFieldCode03: &CustomLookupFieldCode03,
		CustomLookupFieldCode04: &CustomLookupFieldCode04,
		CustomLookupFieldCode05: &CustomLookupFieldCode05,
		CustomLookupFieldCode06: &CustomLookupFieldCode06,
		CustomLookupFieldCode07: &CustomLookupFieldCode07,
		CustomLookupFieldCode08: &CustomLookupFieldCode08,
		CustomLookupFieldCode09: &CustomLookupFieldCode09,
		CustomLookupFieldCode10: &CustomLookupFieldCode10,
		CustomLookupFieldCode11: &CustomLookupFieldCode11,
		CustomLookupFieldCode12: &CustomLookupFieldCode12,
		CreatedBy:               CreatedBy,
	}
	result := DB.Table(DBVar.TableName.ContactPerson).Create(&ContactPerson)
	return ContactPerson.Id, result.Error
}

func UpdateContactPerson(DB *gorm.DB, Id uint64, TitleCode string, FullName string, Street string, CountryCode string, StateCode string, CityCode string, City string,
	NationalityCode string, PostalCode string, Phone1 string, Phone2 string, Fax string, Email string, Website string, CompanyCode string, GuestTypeCode string,
	IdCardCode string, IdCardNumber string, IsMale uint8, BirthPlace string, BirthDate time.Time, TypeCode string, CustomField01 string, CustomField02 string,
	CustomField03 string, CustomField04 string, CustomField05 string, CustomField06 string, CustomField07 string, CustomField08 string, CustomField09 string,
	CustomField10 string, CustomField11 string, CustomField12 string, CustomLookupFieldCode01 string, CustomLookupFieldCode02 string, CustomLookupFieldCode03 string,
	CustomLookupFieldCode04 string, CustomLookupFieldCode05 string, CustomLookupFieldCode06 string, CustomLookupFieldCode07 string, CustomLookupFieldCode08 string,
	CustomLookupFieldCode09 string, CustomLookupFieldCode10 string, CustomLookupFieldCode11 string, CustomLookupFieldCode12 string, UpdatedBy string) error {
	var ContactPerson = DBVar.Contact_person{
		TitleCode:               &TitleCode,
		FullName:                &FullName,
		Street:                  &Street,
		CountryCode:             &CountryCode,
		StateCode:               &StateCode,
		CityCode:                &CityCode,
		City:                    &City,
		NationalityCode:         &NationalityCode,
		PostalCode:              &PostalCode,
		Phone1:                  &Phone1,
		Phone2:                  &Phone2,
		Fax:                     &Fax,
		Email:                   &Email,
		Website:                 &Website,
		CompanyCode:             &CompanyCode,
		GuestTypeCode:           &GuestTypeCode,
		IdCardCode:              &IdCardCode,
		IdCardNumber:            &IdCardNumber,
		IsMale:                  &IsMale,
		BirthPlace:              &BirthPlace,
		BirthDate:               &BirthDate,
		TypeCode:                TypeCode,
		CustomField01:           &CustomField01,
		CustomField02:           &CustomField02,
		CustomField03:           &CustomField03,
		CustomField04:           &CustomField04,
		CustomField05:           &CustomField05,
		CustomField06:           &CustomField06,
		CustomField07:           &CustomField07,
		CustomField08:           &CustomField08,
		CustomField09:           &CustomField09,
		CustomField10:           &CustomField10,
		CustomField11:           &CustomField11,
		CustomField12:           &CustomField12,
		CustomLookupFieldCode01: &CustomLookupFieldCode01,
		CustomLookupFieldCode02: &CustomLookupFieldCode02,
		CustomLookupFieldCode03: &CustomLookupFieldCode03,
		CustomLookupFieldCode04: &CustomLookupFieldCode04,
		CustomLookupFieldCode05: &CustomLookupFieldCode05,
		CustomLookupFieldCode06: &CustomLookupFieldCode06,
		CustomLookupFieldCode07: &CustomLookupFieldCode07,
		CustomLookupFieldCode08: &CustomLookupFieldCode08,
		CustomLookupFieldCode09: &CustomLookupFieldCode09,
		CustomLookupFieldCode10: &CustomLookupFieldCode10,
		CustomLookupFieldCode11: &CustomLookupFieldCode11,
		CustomLookupFieldCode12: &CustomLookupFieldCode12,
		UpdatedBy:               UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.ContactPerson).Where("id", Id).Omit("created_at", "updated_at", "created_by", "id").Updates(&ContactPerson)
	return result.Error
}

func InsertInvoiceItem(DB *gorm.DB, InvoiceNumber string, SubFolioId uint64, FolioNumber uint64, CorrectionBreakdown uint64, Amount float64, AmountCharged float64, DefaultCurrencyCode string, AmountChargedForeign float64, ExchangeRate float64, CurrencyCode string, AmountPaid float64, RefNumber string, Remark string, TypeCode string, CreatedBy string) error {
	var InvoiceItem = DBVar.Invoice_item{
		InvoiceNumber:        InvoiceNumber,
		SubFolioId:           SubFolioId,
		FolioNumber:          FolioNumber,
		CorrectionBreakdown:  CorrectionBreakdown,
		Amount:               Amount,
		AmountCharged:        AmountCharged,
		DefaultCurrencyCode:  DefaultCurrencyCode,
		AmountChargedForeign: AmountChargedForeign,
		ExchangeRate:         ExchangeRate,
		CurrencyCode:         CurrencyCode,
		AmountPaid:           AmountPaid,
		RefNumber:            RefNumber,
		Remark:               Remark,
		TypeCode:             TypeCode,
		CreatedBy:            CreatedBy,
	}
	result := DB.Table(DBVar.TableName.InvoiceItem).Create(&InvoiceItem)
	return result.Error
}

func UpdateInvoiceItem(DB *gorm.DB, InvoiceNumber string, SubFolioId uint64, FolioNumber uint64, CorrectionBreakdown uint64, Amount float64, AmountCharged float64, DefaultCurrencyCode string, AmountChargedForeign float64, ExchangeRate float64, CurrencyCode string, AmountPaid float64, RefNumber string, Remark string, TypeCode string, UpdatedBy string) error {
	var InvoiceItem = DBVar.Invoice_item{
		InvoiceNumber:        InvoiceNumber,
		SubFolioId:           SubFolioId,
		FolioNumber:          FolioNumber,
		CorrectionBreakdown:  CorrectionBreakdown,
		Amount:               Amount,
		AmountCharged:        AmountCharged,
		DefaultCurrencyCode:  DefaultCurrencyCode,
		AmountChargedForeign: AmountChargedForeign,
		ExchangeRate:         ExchangeRate,
		CurrencyCode:         CurrencyCode,
		AmountPaid:           AmountPaid,
		RefNumber:            RefNumber,
		Remark:               Remark,
		TypeCode:             TypeCode,
		UpdatedBy:            UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.InvoiceItem).Omit("created_at", "created_by", "updated_at", "id").Updates(&InvoiceItem)
	return result.Error
}

func InsertCmUpdate(DB *gorm.DB, TypeCode string, Number uint64, RoomTypeCode string, BedTypeCode string, RoomRateCode string, RateAmount float64, StartDate time.Time, EndDate time.Time) error {
	if TypeCode != GlobalVar.CMUpdateType.RoomAllotment && TypeCode != GlobalVar.CMUpdateType.Availability {
		EndDate = General.IncDay(EndDate, -1)
	}
	var CmUpdate = DBVar.Cm_update{
		TypeCode:     TypeCode,
		Number:       Number,
		RoomTypeCode: RoomTypeCode,
		BedTypeCode:  BedTypeCode,
		RoomRateCode: RoomRateCode,
		RateAmount:   RateAmount,
		StartDate:    StartDate,
		EndDate:      EndDate,
		PostingDate:  time.Now(),
	}
	var result error
	if StartDate.Unix() <= EndDate.Unix() {
		result = DB.Table(DBVar.TableName.CmUpdate).Create(&CmUpdate).Error
	}
	return result
}

func UpdateCmUpdate(DB *gorm.DB, TypeCode string, Number uint64, RoomTypeCode string, BedTypeCode string, RoomRateCode string, RateAmount float64, StartDate time.Time, EndDate time.Time, PostingDate time.Time, IsUpdated uint8) error {
	var CmUpdate = DBVar.Cm_update{
		TypeCode:     TypeCode,
		Number:       Number,
		RoomTypeCode: RoomTypeCode,
		BedTypeCode:  BedTypeCode,
		RoomRateCode: RoomRateCode,
		RateAmount:   RateAmount,
		StartDate:    StartDate,
		EndDate:      EndDate,
		PostingDate:  PostingDate,
		IsUpdated:    IsUpdated,
	}
	result := DB.Table(DBVar.TableName.CmUpdate).Omit("id").Updates(&CmUpdate)
	return result.Error
}

func InsertAccCreditCardRecon(DB *gorm.DB, JournalAccountCode string, RefNumber string, Date time.Time, AmountReceived float64, CreatedBy string) (Id uint64, err error) {
	var AccCreditCardRecon = DBVar.Acc_credit_card_recon{
		JournalAccountCode: JournalAccountCode,
		RefNumber:          RefNumber,
		Date:               Date,
		AmountReceived:     AmountReceived,
		CreatedBy:          CreatedBy,
	}
	result := DB.Table(DBVar.TableName.AccCreditCardRecon).Create(&AccCreditCardRecon)
	return AccCreditCardRecon.Id, result.Error
}

func UpdateAccCreditCardRecon(DB *gorm.DB, JournalAccountCode string, RefNumber string, Date time.Time, AmountReceived float64, UpdatedBy string) error {
	var AccCreditCardRecon = DBVar.Acc_credit_card_recon{
		JournalAccountCode: JournalAccountCode,
		RefNumber:          RefNumber,
		Date:               Date,
		AmountReceived:     AmountReceived,
		UpdatedBy:          UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.AccCreditCardRecon).Omit("created_at", "created_by", "updated_at", "id").Updates(&AccCreditCardRecon)
	return result.Error
}

func InsertAccCreditCardReconDetail(DB *gorm.DB, AccCreditCardReconId uint64, GuestDepositId uint64, SubFolioId uint64, Amount float64, Remark string, CreatedBy string) error {
	var AccCreditCardReconDetail = DBVar.Acc_credit_card_recon_detail{
		AccCreditCardReconId: AccCreditCardReconId,
		GuestDepositId:       GuestDepositId,
		SubFolioId:           SubFolioId,
		Amount:               Amount,
		Remark:               Remark,
		CreatedBy:            CreatedBy,
	}
	result := DB.Table(DBVar.TableName.AccCreditCardReconDetail).Create(&AccCreditCardReconDetail)
	return result.Error
}

func UpdateAccCreditCardReconDetail(DB *gorm.DB, AccCreditCardReconId uint64, GuestDepositId uint64, SubFolioId uint64, Amount float64, Remark string, UpdatedBy string) error {
	var AccCreditCardReconDetail = DBVar.Acc_credit_card_recon_detail{
		AccCreditCardReconId: AccCreditCardReconId,
		GuestDepositId:       GuestDepositId,
		SubFolioId:           SubFolioId,
		Amount:               Amount,
		Remark:               Remark,
		UpdatedBy:            UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.AccCreditCardReconDetail).Omit("created_at", "created_by", "updated_at", "id").Updates(&AccCreditCardReconDetail)
	return result.Error
}

func InsertInvoicePayment(DB *gorm.DB, InvoiceNumber string, SubFolioId uint64, FolioNumber uint64, RefNumber string, Amount float64, DefaultCurrencyCode string, AmountForeign float64, ExchangeRate float64, CurrencyCode string, AmountActual float64, ExchangeRateActual float64, Date time.Time, Remark string, CreatedBy string) error {
	var InvoicePayment = DBVar.Invoice_payment{
		InvoiceNumber:       InvoiceNumber,
		SubFolioId:          SubFolioId,
		FolioNumber:         FolioNumber,
		RefNumber:           RefNumber,
		Amount:              Amount,
		DefaultCurrencyCode: DefaultCurrencyCode,
		AmountForeign:       AmountForeign,
		ExchangeRate:        ExchangeRate,
		CurrencyCode:        CurrencyCode,
		AmountActual:        AmountActual,
		ExchangeRateActual:  ExchangeRateActual,
		Date:                Date,
		Remark:              Remark,
		CreatedBy:           CreatedBy,
	}
	result := DB.Table(DBVar.TableName.InvoicePayment).Create(&InvoicePayment)
	if err := DB.Table(DBVar.TableName.InvoiceItem).Where("sub_folio_id=?", SubFolioId).Updates(&map[string]interface{}{
		"amount_paid": gorm.Expr("amount_paid + ?", Amount),
		"updated_by":  CreatedBy,
	}).Error; err != nil {
		return err
	}
	return result.Error
}

func UpdateInvoicePayment(DB *gorm.DB, InvoiceNumber string, SubFolioId uint64, FolioNumber uint64, RefNumber string, Amount float64, DefaultCurrencyCode string, AmountForeign float64, ExchangeRate float64, CurrencyCode string, AmountActual float64, ExchangeRateActual float64, Date time.Time, Remark string, UpdatedBy string) error {
	var InvoicePayment = DBVar.Invoice_payment{
		InvoiceNumber:       InvoiceNumber,
		SubFolioId:          SubFolioId,
		FolioNumber:         FolioNumber,
		RefNumber:           RefNumber,
		Amount:              Amount,
		DefaultCurrencyCode: DefaultCurrencyCode,
		AmountForeign:       AmountForeign,
		ExchangeRate:        ExchangeRate,
		CurrencyCode:        CurrencyCode,
		AmountActual:        AmountActual,
		ExchangeRateActual:  ExchangeRateActual,
		Date:                Date,
		Remark:              Remark,
		UpdatedBy:           UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.InvoicePayment).Omit("created_at", "created_by", "updated_at", "id").Updates(&InvoicePayment)
	return result.Error
}

func InsertAccForeignCash(DB *gorm.DB, IdTransaction uint64, IdCorrected uint64, IdChange uint64, IdTable int, Breakdown uint64, RefNumber string, Date time.Time, TypeCode string, Amount float64, DefaultCurrencyCode string, AmountForeign float64, ExchangeRate float64, CurrencyCode string, Remark string, IsCorrection uint8, CreatedBy string) error {
	Amount = General.RoundToX3(AmountForeign * ExchangeRate)
	var AccForeignCash = DBVar.Acc_foreign_cash{
		IdTransaction:       IdTransaction,
		IdCorrected:         IdCorrected,
		IdChange:            IdChange,
		IdTable:             IdTable,
		Breakdown:           Breakdown,
		RefNumber:           RefNumber,
		Date:                Date,
		TypeCode:            TypeCode,
		Amount:              Amount,
		Stock:               AmountForeign,
		DefaultCurrencyCode: DefaultCurrencyCode,
		AmountForeign:       AmountForeign,
		ExchangeRate:        ExchangeRate,
		CurrencyCode:        CurrencyCode,
		Remark:              Remark,
		IsCorrection:        IsCorrection,
		CreatedBy:           CreatedBy,
	}
	result := DB.Table(DBVar.TableName.AccForeignCash).Create(&AccForeignCash)
	return result.Error
}

func UpdateAccForeignCash(DB *gorm.DB, Id uint64, Date time.Time, Amount float64, AmountForeign float64, ExchangeRate float64, CurrencyCode string, Remark string, UpdatedBy string) error {
	Amount = General.RoundToX3(AmountForeign * ExchangeRate)
	var AccForeignCash = DBVar.Acc_foreign_cash{
		Id:            Id,
		Date:          Date,
		Amount:        Amount,
		AmountForeign: AmountForeign,
		ExchangeRate:  ExchangeRate,
		CurrencyCode:  CurrencyCode,
		Remark:        Remark,
		UpdatedBy:     UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.AccForeignCash).Omit("created_at", "created_by", "updated_at", "id").Updates(&AccForeignCash)
	return result.Error
}

func InsertAccApAr(DB *gorm.DB, Number string, DocumentNumber string, RefNumber string, CompanyCode string, JournalAccountDebit string, JournalAccountCredit string, Amount float64, Date time.Time, DueDate time.Time, Remark string, IsAp uint8, IsAccrued uint8, IsAuto uint8, CreatedBy string) error {
	var AccApAr = DBVar.Acc_ap_ar{
		Number:               Number,
		DocumentNumber:       DocumentNumber,
		RefNumber:            RefNumber,
		CompanyCode:          CompanyCode,
		JournalAccountDebit:  JournalAccountDebit,
		JournalAccountCredit: JournalAccountCredit,
		Amount:               Amount,
		Date:                 Date,
		DueDate:              DueDate,
		Remark:               Remark,
		IsAp:                 IsAp,
		IsAccrued:            IsAccrued,
		IsAuto:               IsAuto,
		CreatedBy:            CreatedBy,
	}
	result := DB.Table(DBVar.TableName.AccApAr).Create(&AccApAr)
	return result.Error
}

func UpdateAccApAr(DB *gorm.DB, Number string, DocumentNumber string, CompanyCode string, JournalAccountDebit string, JournalAccountCredit string, Amount float64, Date time.Time, DueDate time.Time, Remark string, UpdatedBy string) error {
	var AccApAr = DBVar.Acc_ap_ar{
		Number:               Number,
		DocumentNumber:       DocumentNumber,
		CompanyCode:          CompanyCode,
		JournalAccountDebit:  JournalAccountDebit,
		JournalAccountCredit: JournalAccountCredit,
		Amount:               Amount,
		Date:                 Date,
		DueDate:              DueDate,
		Remark:               Remark,
		UpdatedBy:            UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.AccApAr).Where("number=?", Number).Omit("created_at", "created_by", "updated_at", "id").Updates(&AccApAr)
	return result.Error
}
func InsertAccPrepaidExpense(DB *gorm.DB, Date time.Time, RefNumber string, Description string, Amount float64, CompanyCode string, PrepaidAccountCode string, AmountPayment float64, BankAccountCode string, IsCreateAp *uint8, ApArNumber string, SubDepartmentExpenseCode string, ExpenseAccountCode string, Month int, IsNextMonth *uint8, Remark string, CreatedBy string) (Id uint64, err error) {
	var AccPrepaidExpense = DBVar.Acc_prepaid_expense{
		Date:                     Date,
		RefNumber:                RefNumber,
		Description:              Description,
		Amount:                   Amount,
		CompanyCode:              CompanyCode,
		PrepaidAccountCode:       PrepaidAccountCode,
		AmountPayment:            AmountPayment,
		BankAccountCode:          BankAccountCode,
		IsCreateAp:               IsCreateAp,
		ApArNumber:               ApArNumber,
		SubDepartmentExpenseCode: SubDepartmentExpenseCode,
		ExpenseAccountCode:       ExpenseAccountCode,
		Month:                    Month,
		IsNextMonth:              IsNextMonth,
		Remark:                   Remark,
		CreatedBy:                CreatedBy,
	}
	result := DB.Table(DBVar.TableName.AccPrepaidExpense).Create(&AccPrepaidExpense)
	return AccPrepaidExpense.Id, result.Error
}

func UpdateAccPrepaidExpense(DB *gorm.DB, Id uint64, Date time.Time, RefNumber string, Description string, Amount float64, CompanyCode string, PrepaidAccountCode string, AmountPayment float64, BankAccountCode string, IsCreateAp *uint8, ApArNumber string, SubDepartmentExpenseCode string, ExpenseAccountCode string, Month int, IsNextMonth *uint8, Remark string, UpdatedBy string) (IdData uint64, err error) {
	var AccPrepaidExpense = DBVar.Acc_prepaid_expense{
		Id:                       Id,
		Date:                     Date,
		RefNumber:                RefNumber,
		Description:              Description,
		Amount:                   Amount,
		CompanyCode:              CompanyCode,
		PrepaidAccountCode:       PrepaidAccountCode,
		AmountPayment:            AmountPayment,
		BankAccountCode:          BankAccountCode,
		IsCreateAp:               IsCreateAp,
		ApArNumber:               ApArNumber,
		SubDepartmentExpenseCode: SubDepartmentExpenseCode,
		ExpenseAccountCode:       ExpenseAccountCode,
		Month:                    Month,
		IsNextMonth:              IsNextMonth,
		Remark:                   Remark,
		UpdatedBy:                UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.AccPrepaidExpense).Omit("created_at", "created_by", "updated_at", "id").Updates(&AccPrepaidExpense)
	return AccPrepaidExpense.Id, result.Error
}

func InsertAccPrepaidExpensePosted(DB *gorm.DB, PrepaidId uint64, RefNumber string, PostingDate time.Time, Amount float64, SubDepartmentCode string, ExpenseAccountCode string, Remark string, CreatedBy string) error {
	var AccPrepaidExpensePosted = DBVar.Acc_prepaid_expense_posted{
		PrepaidId:          PrepaidId,
		RefNumber:          RefNumber,
		PostingDate:        PostingDate,
		Amount:             Amount,
		SubDepartmentCode:  SubDepartmentCode,
		ExpenseAccountCode: ExpenseAccountCode,
		Remark:             Remark,
		CreatedBy:          CreatedBy,
	}
	result := DB.Table(DBVar.TableName.AccPrepaidExpensePosted).Create(&AccPrepaidExpensePosted)
	return result.Error
}

func UpdateAccPrepaidExpensePosted(DB *gorm.DB, Id, PrepaidId uint64, RefNumber string, PostingDate time.Time, Amount float64, SubDepartmentCode string, ExpenseAccountCode string, Remark string, UpdatedBy string) error {
	var AccPrepaidExpensePosted = DBVar.Acc_prepaid_expense_posted{
		Id:                 Id,
		PrepaidId:          PrepaidId,
		RefNumber:          RefNumber,
		PostingDate:        PostingDate,
		Amount:             Amount,
		SubDepartmentCode:  SubDepartmentCode,
		ExpenseAccountCode: ExpenseAccountCode,
		Remark:             Remark,
		UpdatedBy:          UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.AccPrepaidExpensePosted).Omit("created_at", "created_by", "updated_at", "id").Updates(&AccPrepaidExpensePosted)
	return result.Error
}
func InsertAccDifferedIncome(DB *gorm.DB, Date time.Time, RefNumber string, Description string, Amount float64, CompanyCode string, DefferedAccountCode string, AmountPayment float64, BankAccountCode string, IsCreateAr *uint8, ApArNumber string, SubDepartmentIncomeCode string, IncomeAccountCode string, Month int, IsNextMonth *uint8, Remark string, CreatedBy string) error {
	var AccDefferedIncome = DBVar.Acc_deffered_income{
		Date:                    Date,
		RefNumber:               RefNumber,
		Description:             Description,
		Amount:                  Amount,
		CompanyCode:             CompanyCode,
		DefferedAccountCode:     DefferedAccountCode,
		AmountPayment:           AmountPayment,
		BankAccountCode:         BankAccountCode,
		IsCreateAr:              IsCreateAr,
		ApArNumber:              ApArNumber,
		SubDepartmentIncomeCode: SubDepartmentIncomeCode,
		IncomeAccountCode:       IncomeAccountCode,
		Month:                   Month,
		IsNextMonth:             IsNextMonth,
		Remark:                  Remark,
		CreatedBy:               CreatedBy,
	}
	result := DB.Table(DBVar.TableName.AccDefferedIncome).Create(&AccDefferedIncome)
	return result.Error
}

func UpdateAccDifferedIncome(DB *gorm.DB, Id uint64, Date time.Time, RefNumber string, Description string, Amount float64, CompanyCode string, DefferedAccountCode string, AmountPayment float64, BankAccountCode string, IsCreateAr *uint8, ApArNumber string, SubDepartmentIncomeCode string, IncomeAccountCode string, Month int, IsNextMonth *uint8, Remark string, UpdatedBy string) error {
	var AccDefferedIncome = DBVar.Acc_deffered_income{
		Id:                      Id,
		Date:                    Date,
		RefNumber:               RefNumber,
		Description:             Description,
		Amount:                  Amount,
		CompanyCode:             CompanyCode,
		DefferedAccountCode:     DefferedAccountCode,
		AmountPayment:           AmountPayment,
		BankAccountCode:         BankAccountCode,
		IsCreateAr:              IsCreateAr,
		ApArNumber:              ApArNumber,
		SubDepartmentIncomeCode: SubDepartmentIncomeCode,
		IncomeAccountCode:       IncomeAccountCode,
		Month:                   Month,
		IsNextMonth:             IsNextMonth,
		Remark:                  Remark,
		UpdatedBy:               UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.AccDefferedIncome).Omit("created_at", "created_by", "updated_at", "id").Updates(&AccDefferedIncome)
	return result.Error
}

func InsertAccApArPayment(DB *gorm.DB, RefNumber string, JournalAccountCode string, ApJournalAccountCode string, CreateApNumber string, DiscountJournalAccountCode string, BaJournalAccountCode string, OeJournalAccountCode string, TotalAmount float64, Discount float64, BankAdministration float64, OtherExpense float64, Date time.Time, Remark string, SourceCodeApAr string, IsPaymentApAr *uint8, CreatedBy string) error {
	var AccApArPayment = DBVar.Acc_ap_ar_payment{
		RefNumber:                  RefNumber,
		JournalAccountCode:         JournalAccountCode,
		ApJournalAccountCode:       ApJournalAccountCode,
		CreateApNumber:             CreateApNumber,
		DiscountJournalAccountCode: DiscountJournalAccountCode,
		BaJournalAccountCode:       BaJournalAccountCode,
		OeJournalAccountCode:       OeJournalAccountCode,
		TotalAmount:                TotalAmount,
		Discount:                   Discount,
		BankAdministration:         BankAdministration,
		OtherExpense:               OtherExpense,
		Date:                       Date,
		Remark:                     Remark,
		SourceCodeApAr:             SourceCodeApAr,
		IsPaymentApAr:              IsPaymentApAr,
		CreatedBy:                  CreatedBy,
	}
	result := DB.Table(DBVar.TableName.AccApArPayment).Create(&AccApArPayment)
	return result.Error
}

func UpdateAccApArPayment(DB *gorm.DB, RefNumber string, JournalAccountCode string, ApJournalAccountCode string, CreateApNumber string, DiscountJournalAccountCode string, BaJournalAccountCode string, OeJournalAccountCode string, TotalAmount float64, Discount float64, BankAdministration float64, OtherExpense float64, Date time.Time, Remark string, SourceCodeApAr string, IsPaymentApAr *uint8, UpdatedBy string) error {
	var AccApArPayment = DBVar.Acc_ap_ar_payment{
		RefNumber:                  RefNumber,
		JournalAccountCode:         JournalAccountCode,
		ApJournalAccountCode:       ApJournalAccountCode,
		CreateApNumber:             CreateApNumber,
		DiscountJournalAccountCode: DiscountJournalAccountCode,
		BaJournalAccountCode:       BaJournalAccountCode,
		OeJournalAccountCode:       OeJournalAccountCode,
		TotalAmount:                TotalAmount,
		Discount:                   Discount,
		BankAdministration:         BankAdministration,
		OtherExpense:               OtherExpense,
		Date:                       Date,
		Remark:                     Remark,
		SourceCodeApAr:             SourceCodeApAr,
		IsPaymentApAr:              IsPaymentApAr,
		UpdatedBy:                  UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.AccApArPayment).Omit("created_at", "created_by", "updated_at", "id").Updates(&AccApArPayment)
	return result.Error
}

func InsertReceipt(DB *gorm.DB, Number string, ReceiveFrom string, Amount float64, IssuedDate time.Time, ForPayment string, CreatedBy string) error {
	var Receipt = DBVar.Receipt{
		Number:      Number,
		ReceiveFrom: ReceiveFrom,
		Amount:      Amount,
		IssuedDate:  IssuedDate,
		ForPayment:  ForPayment,
		CreatedBy:   CreatedBy,
	}
	result := DB.Table(DBVar.TableName.Receipt).Create(&Receipt)
	return result.Error
}

func UpdateReceipt(DB *gorm.DB, Number string, ReceiveFrom string, Amount float64, IssuedDate time.Time, ForPayment string, UpdatedBy string) error {
	var Receipt = DBVar.Receipt{
		Number:      Number,
		ReceiveFrom: ReceiveFrom,
		Amount:      Amount,
		IssuedDate:  IssuedDate,
		ForPayment:  ForPayment,
		UpdatedBy:   UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.Receipt).Omit("created_at", "created_by", "updated_at", "id").Updates(&Receipt)
	return result.Error
}

func InsertAccDefferedIncomePosted(DB *gorm.DB, DefferedId uint64, RefNumber string, PostingDate time.Time, Amount float64, SubDepartmentCode string, IncomeAccountCode string, Remark string, CreatedBy string) error {
	var AccDefferedIncomePosted = DBVar.Acc_deffered_income_posted{
		DefferedId:        DefferedId,
		RefNumber:         RefNumber,
		PostingDate:       PostingDate,
		Amount:            Amount,
		SubDepartmentCode: SubDepartmentCode,
		IncomeAccountCode: IncomeAccountCode,
		Remark:            Remark,
		CreatedBy:         CreatedBy,
	}
	result := DB.Table(DBVar.TableName.AccDefferedIncomePosted).Create(&AccDefferedIncomePosted)
	return result.Error
}

func UpdateAccDefferedIncomePosted(DB *gorm.DB, Id, DefferedId uint64, RefNumber string, PostingDate time.Time, Amount float64, SubDepartmentCode string, IncomeAccountCode string, Remark string, UpdatedBy string) error {
	var AccDefferedIncomePosted = DBVar.Acc_deffered_income_posted{
		Id:                Id,
		DefferedId:        DefferedId,
		RefNumber:         RefNumber,
		PostingDate:       PostingDate,
		Amount:            Amount,
		SubDepartmentCode: SubDepartmentCode,
		IncomeAccountCode: IncomeAccountCode,
		Remark:            Remark,
		UpdatedBy:         UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.AccDefferedIncomePosted).Omit("created_at", "created_by", "updated_at", "id").Updates(&AccDefferedIncomePosted)
	return result.Error
}

func InsertAccApCommissionPayment(DB *gorm.DB, RefNumber string, JournalAccountCode string, DiscountJournalAccountCode string, BaJournalAccountCode string, OeJournalAccountCode string, TotalAmount float64, Discount float64, BankAdministration float64, OtherExpense float64, Date time.Time, Remark string, SourceCodeApAr string, CreatedBy string) error {
	var AccApCommissionPayment = DBVar.Acc_ap_commission_payment{
		RefNumber:                  RefNumber,
		JournalAccountCode:         JournalAccountCode,
		DiscountJournalAccountCode: DiscountJournalAccountCode,
		BaJournalAccountCode:       BaJournalAccountCode,
		OeJournalAccountCode:       OeJournalAccountCode,
		TotalAmount:                TotalAmount,
		Discount:                   Discount,
		BankAdministration:         BankAdministration,
		OtherExpense:               OtherExpense,
		Date:                       Date,
		Remark:                     Remark,
		SourceCodeApAr:             SourceCodeApAr,
		CreatedBy:                  CreatedBy,
	}
	result := DB.Table(DBVar.TableName.AccApCommissionPayment).Create(&AccApCommissionPayment)
	return result.Error
}

func UpdateAccApCommissionPayment(DB *gorm.DB, RefNumber string, JournalAccountCode string, DiscountJournalAccountCode string, BaJournalAccountCode string, OeJournalAccountCode string, TotalAmount float64, Discount float64, BankAdministration float64, OtherExpense float64, Date time.Time, Remark string, SourceCodeApAr string, UpdatedBy string) error {
	var AccApCommissionPayment = DBVar.Acc_ap_commission_payment{
		RefNumber:                  RefNumber,
		JournalAccountCode:         JournalAccountCode,
		DiscountJournalAccountCode: DiscountJournalAccountCode,
		BaJournalAccountCode:       BaJournalAccountCode,
		OeJournalAccountCode:       OeJournalAccountCode,
		TotalAmount:                TotalAmount,
		Discount:                   Discount,
		BankAdministration:         BankAdministration,
		OtherExpense:               OtherExpense,
		Date:                       Date,
		Remark:                     Remark,
		SourceCodeApAr:             SourceCodeApAr,
		UpdatedBy:                  UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.AccApCommissionPayment).Omit("created_at", "created_by", "updated_at", "id").Updates(&AccApCommissionPayment)
	return result.Error
}

func InsertAccApCommissionPaymentDetail(DB *gorm.DB, SubFolioId uint64, RefNumber string, Amount float64, CreatedBy string) error {
	var AccApCommissionPaymentDetail = DBVar.Acc_ap_commission_payment_detail{
		SubFolioId: SubFolioId,
		RefNumber:  RefNumber,
		Amount:     Amount,
		CreatedBy:  CreatedBy,
	}
	result := DB.Table(DBVar.TableName.AccApCommissionPaymentDetail).Create(&AccApCommissionPaymentDetail)
	return result.Error
}

func UpdateAccApCommissionPaymentDetail(DB *gorm.DB, SubFolioId uint64, RefNumber string, Amount float64, UpdatedBy string) error {
	var AccApCommissionPaymentDetail = DBVar.Acc_ap_commission_payment_detail{
		SubFolioId: SubFolioId,
		RefNumber:  RefNumber,
		Amount:     Amount,
		UpdatedBy:  UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.AccApCommissionPaymentDetail).Omit("created_at", "created_by", "updated_at", "id").Updates(&AccApCommissionPaymentDetail)
	return result.Error
}

func InsertInvPurchaseOrder(ctx context.Context, DB *gorm.DB, Number string, CompanyCode string, ExpeditionCode string, PrNumber string, Date time.Time, ShippingAddressCode string, ContactPerson string, Street string, City string, CountryCode string, StateCode string, PostalCode string, Phone1 string, Phone2 string, Fax string, Email string, RequestBy string, Remark string, IsReceived uint8, CreatedBy string) error {
	var InvPurchaseOrder = DBVar.Inv_purchase_order{
		Number:              Number,
		CompanyCode:         CompanyCode,
		ExpeditionCode:      ExpeditionCode,
		PrNumber:            PrNumber,
		Date:                Date,
		ShippingAddressCode: ShippingAddressCode,
		ContactPerson:       ContactPerson,
		Street:              Street,
		City:                City,
		CountryCode:         CountryCode,
		StateCode:           StateCode,
		PostalCode:          PostalCode,
		Phone1:              Phone1,
		Phone2:              Phone2,
		Fax:                 Fax,
		Email:               Email,
		RequestBy:           RequestBy,
		Remark:              Remark,
		IsReceived:          IsReceived,
		CreatedBy:           CreatedBy,
	}
	result := DB.WithContext(ctx).Table(DBVar.TableName.InvPurchaseOrder).Create(&InvPurchaseOrder)
	return result.Error
}

func UpdateInvPurchaseOrder(DB *gorm.DB, Id uint64, Number string, CompanyCode string, ExpeditionCode string, PrNumber string, Date time.Time, ShippingAddressCode string, ContactPerson string, Street string, City string, CountryCode string, StateCode string, PostalCode string, Phone1 string, Phone2 string, Fax string, Email string, RequestBy string, Remark string, UpdatedBy string) error {
	var InvPurchaseOrder = DBVar.Inv_purchase_order{
		Id:                  Id,
		Number:              Number,
		CompanyCode:         CompanyCode,
		ExpeditionCode:      ExpeditionCode,
		PrNumber:            PrNumber,
		Date:                Date,
		ShippingAddressCode: ShippingAddressCode,
		ContactPerson:       ContactPerson,
		Street:              Street,
		City:                City,
		CountryCode:         CountryCode,
		StateCode:           StateCode,
		PostalCode:          PostalCode,
		Phone1:              Phone1,
		Phone2:              Phone2,
		Fax:                 Fax,
		Email:               Email,
		RequestBy:           RequestBy,
		Remark:              Remark,
		UpdatedBy:           UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.InvPurchaseOrder).Omit("created_at", "created_by", "updated_at", "id").Updates(&InvPurchaseOrder)
	return result.Error
}

func InsertInvPurchaseOrderDetail(ctx context.Context, DB *gorm.DB, PoNumber string, ItemCode string, StoreCode string, Quantity float64, QuantityReceived float64, QuantityNotReceived float64, Convertion float64, UomCode string, Price float64, Remark string, CreatedBy string) error {
	var InvPurchaseOrderDetail = DBVar.Inv_purchase_order_detail{
		PoNumber:            PoNumber,
		ItemCode:            ItemCode,
		StoreCode:           StoreCode,
		Quantity:            Quantity,
		QuantityReceived:    QuantityReceived,
		QuantityNotReceived: QuantityNotReceived,
		Convertion:          Convertion,
		UomCode:             UomCode,
		Price:               Price,
		Remark:              Remark,
		CreatedBy:           CreatedBy,
	}
	result := DB.WithContext(ctx).Table(DBVar.TableName.InvPurchaseOrderDetail).Create(&InvPurchaseOrderDetail)
	return result.Error
}

func UpdateInvPurchaseOrderDetail(DB *gorm.DB, PoNumber string, ItemCode string, StoreCode string, Quantity float64, QuantityReceived float64, QuantityNotReceived float64, Convertion float64, UomCode string, Price float64, Remark string, UpdatedBy string) error {
	var InvPurchaseOrderDetail = DBVar.Inv_purchase_order_detail{
		PoNumber:            PoNumber,
		ItemCode:            ItemCode,
		StoreCode:           StoreCode,
		Quantity:            Quantity,
		QuantityReceived:    QuantityReceived,
		QuantityNotReceived: QuantityNotReceived,
		Convertion:          Convertion,
		UomCode:             UomCode,
		Price:               Price,
		Remark:              Remark,
		UpdatedBy:           UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.InvPurchaseOrderDetail).Omit("created_at", "created_by", "updated_at", "id").Updates(&InvPurchaseOrderDetail)
	return result.Error
}

func InsertInvCostRecipe(DB *gorm.DB, ProductCode string, StoreCode string, ItemCode string, Quantity float64, UomCode string, Remark string, CreatedBy string) error {
	var InvCostRecipe = DBVar.Inv_cost_recipe{
		ProductCode: ProductCode,
		StoreCode:   StoreCode,
		ItemCode:    ItemCode,
		Quantity:    Quantity,
		UomCode:     UomCode,
		Remark:      Remark,
		CreatedBy:   CreatedBy,
	}
	result := DB.Table(DBVar.TableName.InvCostRecipe).Create(&InvCostRecipe)
	return result.Error
}

func UpdateInvCostRecipe(DB *gorm.DB, ProductCode string, StoreCode string, ItemCode string, Quantity float64, UomCode string, Remark string, UpdatedBy string) error {
	var InvCostRecipe = DBVar.Inv_cost_recipe{
		ProductCode: ProductCode,
		StoreCode:   StoreCode,
		ItemCode:    ItemCode,
		Quantity:    Quantity,
		UomCode:     UomCode,
		Remark:      Remark,
		UpdatedBy:   UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.InvCostRecipe).Omit("created_at", "created_by", "updated_at", "id").Updates(&InvCostRecipe)
	return result.Error
}

func InsertInvPurchaseRequest(DB *gorm.DB, Number string, SubDepartmentCode string, Date time.Time, NeedDate time.Time, RequestBy string, Remark string, CreatedBy string) error {
	var InvPurchaseRequest = DBVar.Inv_purchase_request{
		Number:            Number,
		ContactPerson:     GlobalVar.EmptyString,
		SubDepartmentCode: SubDepartmentCode,
		Date:              Date,
		NeedDate:          NeedDate,
		RequestBy:         &RequestBy,
		Remark:            &Remark,
		StatusCode:        GlobalVar.PurchaseRequestStatus.NotApproved,
		CreatedBy:         CreatedBy,
	}
	result := DB.Table(DBVar.TableName.InvPurchaseRequest).Create(&InvPurchaseRequest)
	return result.Error
}

func UpdateInvPurchaseRequest(DB *gorm.DB, Id uint64, Number string, SubDepartmentCode string, Date time.Time, NeedDate time.Time, ShippingAddressCode string, ContactPerson string, Street string, City string, CountryCode string, PostalCode string, Phone1 string, Phone2 string, Fax string, Email string, RequestBy string, Remark string, UpdatedBy string) error {
	var InvPurchaseRequest = DBVar.Inv_purchase_request{
		Id:                  Id,
		Number:              Number,
		SubDepartmentCode:   SubDepartmentCode,
		Date:                Date,
		NeedDate:            NeedDate,
		ShippingAddressCode: ShippingAddressCode,
		ContactPerson:       &ContactPerson,
		Street:              &Street,
		City:                &City,
		CountryCode:         &CountryCode,
		PostalCode:          &PostalCode,
		Phone1:              &Phone1,
		Phone2:              &Phone2,
		Fax:                 &Fax,
		Email:               &Email,
		RequestBy:           &RequestBy,
		Remark:              &Remark,
		UpdatedBy:           UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.InvPurchaseRequest).Omit("created_at", "created_by", "updated_at", "id").Updates(&InvPurchaseRequest)
	return result.Error
}

func InsertInvPurchaseRequestDetail(DB *gorm.DB, PrNumber string, ItemCode string, Quantity float64, Convertion float64, UomCode string, CompanyCode string, Price float64, CompanyCode2 string, Price2 float64, CompanyCode3 string, Price3 float64, EstimatePrice float64, StoreCode string, Remark string, CreatedBy string) error {
	var InvPurchaseRequestDetail = DBVar.Inv_purchase_request_detail{
		PrNumber: PrNumber,
		ItemCode: ItemCode,
		Quantity: Quantity,
		// QuantityApproved: QuantityApproved,
		Convertion:    Convertion,
		UomCode:       UomCode,
		CompanyCode:   CompanyCode,
		Price:         Price,
		CompanyCode2:  &CompanyCode2,
		Price2:        &Price2,
		CompanyCode3:  &CompanyCode3,
		Price3:        &Price3,
		EstimatePrice: EstimatePrice,
		StoreCode:     StoreCode,
		Remark:        &Remark,
		CreatedBy:     CreatedBy,
	}
	result := DB.Table(DBVar.TableName.InvPurchaseRequestDetail).Create(&InvPurchaseRequestDetail)
	return result.Error
}

func UpdateInvPurchaseRequestDetail(DB *gorm.DB, PrNumber string, ItemCode string, Quantity float64, QuantityApproved *float64, Convertion float64, UomCode string, CompanyCode string, Price float64, CompanyCode2 *string, Price2 *float64, CompanyCode3 *string, Price3 *float64, EstimatePrice float64, StoreCode string, Remark *string, UpdatedBy string) error {
	var InvPurchaseRequestDetail = DBVar.Inv_purchase_request_detail{
		PrNumber:         PrNumber,
		ItemCode:         ItemCode,
		Quantity:         Quantity,
		QuantityApproved: QuantityApproved,
		Convertion:       Convertion,
		UomCode:          UomCode,
		CompanyCode:      CompanyCode,
		Price:            Price,
		CompanyCode2:     CompanyCode2,
		Price2:           Price2,
		CompanyCode3:     CompanyCode3,
		Price3:           Price3,
		EstimatePrice:    EstimatePrice,
		StoreCode:        StoreCode,
		Remark:           Remark,
		UpdatedBy:        UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.InvPurchaseRequestDetail).Omit("created_at", "created_by", "updated_at", "id").Updates(&InvPurchaseRequestDetail)
	return result.Error
}

func InsertInvCosting(DB *gorm.DB, Number string, RefNumber string, DocumentNumber string, SubDepartmentCode string, StoreCode string, Date time.Time, RequestBy string, Remark string, IsStoreRequisition uint8, IsOpname uint8, IsProduction uint8, IsReturn uint8, IsRoom uint8, IsCostRecipe uint8, CreatedBy string) error {
	var InvCosting = DBVar.Inv_costing{
		Number:             Number,
		RefNumber:          RefNumber,
		DocumentNumber:     DocumentNumber,
		SubDepartmentCode:  SubDepartmentCode,
		StoreCode:          StoreCode,
		Date:               Date,
		RequestBy:          RequestBy,
		Remark:             &Remark,
		IsStoreRequisition: IsStoreRequisition,
		IsOpname:           IsOpname,
		IsProduction:       IsProduction,
		IsReturn:           IsReturn,
		IsRoom:             IsRoom,
		IsCostRecipe:       IsCostRecipe,
		CreatedBy:          CreatedBy,
	}
	result := DB.Table(DBVar.TableName.InvCosting).Create(&InvCosting)
	return result.Error
}

func UpdateInvCosting(DB *gorm.DB, Number string, DocumentNumber string, SubDepartmentCode string, StoreCode string, Date time.Time, RequestBy string, Remark string, UpdatedBy string) error {
	var InvCosting = DBVar.Inv_costing{
		Number:            Number,
		DocumentNumber:    DocumentNumber,
		SubDepartmentCode: SubDepartmentCode,
		StoreCode:         StoreCode,
		Date:              Date,
		RequestBy:         RequestBy,
		Remark:            &Remark,
		UpdatedBy:         UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.InvCosting).Omit("created_at", "created_by", "updated_at", "id").Where("number=?", Number).Updates(&InvCosting)
	return result.Error
}

func InsertInvCostingDetail(DB *gorm.DB, Dataset *GlobalVar.TDataset, CostingNumber string, StoreCode string, StoreId uint64, ItemCode string, ItemId uint64, Date time.Time, Quantity float64, UomCode string, TotalPrice float64, ReceiveId uint64, JournalAccountCode string, ItemGroupCode string, ReasonCode string, IsSpoil uint8, IsCogs uint8, CreatedBy string) error {
	Price := TotalPrice / Quantity
	if Dataset.ProgramConfiguration.ReceiveStockAPTwoDigitDecimal {
		Price = General.RoundToX2(Price)
		TotalPrice = General.RoundToX2(TotalPrice)
	} else {
		Price = General.RoundToX3(Price)
		TotalPrice = General.RoundToX3(TotalPrice)
	}
	var InvCostingDetail = DBVar.Inv_costing_detail{
		CostingNumber:      CostingNumber,
		StoreCode:          StoreCode,
		StoreId:            StoreId,
		ItemCode:           ItemCode,
		ItemId:             ItemId,
		Date:               Date,
		Quantity:           Quantity,
		UomCode:            UomCode,
		Price:              Price,
		TotalPrice:         TotalPrice,
		ReceiveId:          ReceiveId,
		JournalAccountCode: JournalAccountCode,
		ItemGroupCode:      ItemGroupCode,
		ReasonCode:         ReasonCode,
		IsSpoil:            IsSpoil,
		IsCogs:             IsCogs,
		CreatedBy:          CreatedBy,
	}
	result := DB.Table(DBVar.TableName.InvCostingDetail).Create(&InvCostingDetail)
	return result.Error
}

func UpdateInvCostingDetail(DB *gorm.DB, CostingNumber string, StoreCode string, StoreId uint64, ItemCode string, ItemId uint64, Date time.Time, Quantity float64, UomCode string, TotalPrice float64, ReceiveId uint64, JournalAccountCode string, ItemGroupCode string, ReasonCode string, IsSpoil uint8, IsCogs uint8, UpdatedBy string) error {
	Price := TotalPrice / Quantity
	var InvCostingDetail = DBVar.Inv_costing_detail{
		CostingNumber:      CostingNumber,
		StoreCode:          StoreCode,
		StoreId:            StoreId,
		ItemCode:           ItemCode,
		ItemId:             ItemId,
		Date:               Date,
		Quantity:           Quantity,
		UomCode:            UomCode,
		Price:              Price,
		TotalPrice:         TotalPrice,
		ReceiveId:          ReceiveId,
		JournalAccountCode: JournalAccountCode,
		ItemGroupCode:      ItemGroupCode,
		ReasonCode:         ReasonCode,
		IsSpoil:            IsSpoil,
		IsCogs:             IsCogs,
		UpdatedBy:          UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.InvCostingDetail).Omit("created_at", "created_by", "updated_at", "id").Updates(&InvCostingDetail)
	return result.Error
}

func InsertReservation(DB *gorm.DB, PostingDate time.Time, ContactPersonId1 uint64, ContactPersonId2 uint64, ContactPersonId3 uint64, ContactPersonId4 uint64, GuestDetailId uint64, GuestProfileId1 uint64, GuestProfileId2 uint64, GuestProfileId3 uint64, GuestProfileId4 uint64, GuestGeneralId uint64, ReservationBy string, GroupCode string, MemberCode string, IsWaitList uint8, IsIncognito uint8, BookingCode string, OtaId string, CmResStatus string, CreatedBy string) (uint64, error) {
	StatusCode := GlobalVar.ReservationStatus.New
	if IsWaitList > 0 {
		StatusCode = GlobalVar.ReservationStatus.WaitList
	}
	var Reservation = DBVar.Reservation{
		ContactPersonId1: ContactPersonId1,
		ContactPersonId2: ContactPersonId2,
		ContactPersonId3: ContactPersonId3,
		ContactPersonId4: ContactPersonId4,
		GuestDetailId:    GuestDetailId,
		GuestProfileId1:  GuestProfileId1,
		GuestProfileId2:  GuestProfileId2,
		GuestProfileId3:  GuestProfileId3,
		GuestProfileId4:  GuestProfileId4,
		GuestGeneralId:   GuestGeneralId,
		ReservationBy:    ReservationBy,
		AuditDate:        PostingDate,
		GroupCode:        GroupCode,
		MemberCode:       MemberCode,
		IsIncognito:      &IsIncognito,
		OtaId:            OtaId,
		BookingCode:      BookingCode,
		CmResStatus:      CmResStatus,
		StatusCode:       StatusCode,
		StatusCode2:      GlobalVar.ReservationStatus2.Tentative,
		CreatedBy:        CreatedBy,
	}

	result := DB.Table(DBVar.TableName.Reservation).Create(&Reservation)
	return Reservation.Number, result.Error
}

func UpdateReservation(DB *gorm.DB, Number uint64, ReservationBy string, GroupCode string, MemberCode string, IsWaitList uint8, IsIncognito uint8, UpdatedBy string) error {
	StatusCode := GlobalVar.ReservationStatus.New
	if IsWaitList > 0 {
		StatusCode = GlobalVar.ReservationStatus.WaitList
	}
	var Reservation = DBVar.Reservation{
		Number:        Number,
		ReservationBy: ReservationBy,
		GroupCode:     GroupCode,
		MemberCode:    MemberCode,
		StatusCode:    StatusCode,
		IsIncognito:   &IsIncognito,
		UpdatedBy:     UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.Reservation).Omit("created_at", "created_by", "updated_at").Updates(&Reservation)

	return result.Error
}

func InsertGuestDetail(DB *gorm.DB, Arrival time.Time, Departure time.Time, Adult int, Child int, RoomTypeCode string, BedTypeCode string, RoomNumber string, CurrencyCode string, ExchangeRate float64, IsConstantCurrency uint8, RoomRateCode string, IsOverrideRate uint8, WeekdayRate float64, WeekendRate float64, DiscountPercent uint8, Discount float64, BusinessSourceCode string, IsOverrideCommission uint8, CommissionTypeCode string, CommissionValue float64, PaymentTypeCode string, MarketCode string, BookingSourceCode string, BillInstruction string, CreatedBy string) (Id uint64, err error) {
	if ExchangeRate <= 0 {
		ExchangeRate = GetExchangeRateCurrency(DB, CurrencyCode)
	}
	var DepartureUnix, ArrivalUnixx int64
	if !Departure.IsZero() {
		DepartureUnix = General.DateOf(Departure).Unix()
	}

	if !Arrival.IsZero() {
		ArrivalUnixx = General.DateOf(Arrival).Unix()
	}

	var GuestDetail = DBVar.Guest_detail{
		Arrival:              Arrival,
		ArrivalUnixx:         ArrivalUnixx,
		ArrivalRes:           Arrival,
		Departure:            Departure,
		DepartureUnixx:       DepartureUnix,
		DepartureRes:         Departure,
		Adult:                Adult,
		Child:                General.PtrInt(Child),
		RoomTypeCode:         RoomTypeCode,
		BedTypeCode:          BedTypeCode,
		RoomNumber:           General.PtrString(RoomNumber),
		CurrencyCode:         CurrencyCode,
		ExchangeRate:         ExchangeRate,
		IsConstantCurrency:   IsConstantCurrency,
		RoomRateCode:         RoomRateCode,
		IsOverrideRate:       General.PtrUint8(IsOverrideRate),
		WeekdayRate:          General.PtrFloat64(WeekdayRate),
		WeekendRate:          General.PtrFloat64(WeekendRate),
		DiscountPercent:      General.PtrUint8(DiscountPercent),
		Discount:             General.PtrFloat64(Discount),
		BusinessSourceCode:   General.PtrString(BusinessSourceCode),
		IsOverrideCommission: General.PtrUint8(IsOverrideCommission),
		CommissionTypeCode:   General.PtrString(CommissionTypeCode),
		CommissionValue:      General.PtrFloat64(CommissionValue),
		PaymentTypeCode:      PaymentTypeCode,
		MarketCode:           General.PtrString(MarketCode),
		BookingSourceCode:    General.PtrString(BookingSourceCode),
		BillInstruction:      General.PtrString(BillInstruction),
		CreatedBy:            CreatedBy,
	}
	result := DB.Table(DBVar.TableName.GuestDetail).Create(&GuestDetail)
	return GuestDetail.Id, result.Error
}

func UpdateGuestDetail(DB *gorm.DB, Id uint64, Arrival time.Time, Departure time.Time, Adult int, Child int, RoomTypeCode string, BedTypeCode string, RoomNumber string,
	CurrencyCode string, ExchangeRate float64, IsConstantCurrency uint8, RoomRateCode string, IsOverrideRate uint8, WeekdayRate float64, WeekendRate float64, DiscountPercent uint8, Discount float64, BusinessSourceCode string,
	IsOverrideCommission uint8, CommissionTypeCode string, CommissionValue float64, PaymentTypeCode string, MarketCode string, BookingSourceCode string, BillInstruction string, UpdatedBy string) error {

	var DepartureUnix, ArrivalUnixx int64
	if !Departure.IsZero() {
		DepartureUnix = General.DateOf(Departure).Unix()
	}

	if !Arrival.IsZero() {
		ArrivalUnixx = General.DateOf(Arrival).Unix()
	}

	var GuestDetail = DBVar.Guest_detail{
		Id:                   Id,
		Arrival:              Arrival,
		ArrivalUnixx:         ArrivalUnixx,
		Departure:            Departure,
		DepartureUnixx:       DepartureUnix,
		Adult:                Adult,
		Child:                General.PtrInt(Child),
		RoomTypeCode:         RoomTypeCode,
		BedTypeCode:          BedTypeCode,
		RoomNumber:           General.PtrString(RoomNumber),
		CurrencyCode:         CurrencyCode,
		ExchangeRate:         ExchangeRate,
		IsConstantCurrency:   IsConstantCurrency,
		RoomRateCode:         RoomRateCode,
		IsOverrideRate:       General.PtrUint8(IsOverrideRate),
		WeekdayRate:          General.PtrFloat64(WeekdayRate),
		WeekendRate:          General.PtrFloat64(WeekendRate),
		DiscountPercent:      General.PtrUint8(DiscountPercent),
		Discount:             General.PtrFloat64(Discount),
		BusinessSourceCode:   General.PtrString(BusinessSourceCode),
		IsOverrideCommission: General.PtrUint8(IsOverrideCommission),
		CommissionTypeCode:   General.PtrString(CommissionTypeCode),
		CommissionValue:      General.PtrFloat64(CommissionValue),
		PaymentTypeCode:      PaymentTypeCode,
		MarketCode:           General.PtrString(MarketCode),
		BookingSourceCode:    General.PtrString(BookingSourceCode),
		BillInstruction:      General.PtrString(BillInstruction),
		UpdatedBy:            UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.GuestDetail).Omit("created_at", "created_by", "updated_at", "id").Updates(&GuestDetail)
	return result.Error
}

func InsertGuestInHouse(DB *gorm.DB, AuditDate time.Time, FolioNumber uint64, GroupCode, RoomTypeCode, BedTypeCode,
	RoomNumber, RoomRateCode, BusinessSourceCode, CommissionTypeCode, PaymentTypeCode, MarketCode, TitleCode,
	FullName, Street, City, CityCode, CountryCode, StateCode, PostalCode, Phone1, Phone2, Fax, Email, Website, CompanyCode, GuestTypeCode, SalesCode,
	ComplimentHu, Notes string, Adult, Child int, Rate, RateOriginal, Discount, CommissionValue float64,
	DiscountPercent, IsAdditional, IsScheduledRate, IsBreakfast uint8, BookingSourceCode, PurposeOfCode, CustomLookupFieldCode01, CustomLookupFieldCode02 string,
	PaxBreakfast int, BreakfastVoucherNumber, NationalityCode string, CreatedBy string) error {
	var GuestInHouse = DBVar.Guest_in_house{
		AuditDate:               AuditDate,
		AuditDateUnixx:          AuditDate.Unix(),
		FolioNumber:             FolioNumber,
		GroupCode:               GroupCode,
		Adult:                   Adult,
		Child:                   Child,
		RoomTypeCode:            RoomTypeCode,
		BedTypeCode:             BedTypeCode,
		RoomNumber:              RoomNumber,
		RoomRateCode:            RoomRateCode,
		RateOriginal:            RateOriginal,
		Rate:                    Rate,
		DiscountPercent:         DiscountPercent,
		Discount:                Discount,
		BusinessSourceCode:      BusinessSourceCode,
		CommissionTypeCode:      CommissionTypeCode,
		CommissionValue:         CommissionValue,
		PaymentTypeCode:         PaymentTypeCode,
		MarketCode:              MarketCode,
		BookingSourceCode:       BookingSourceCode,
		TitleCode:               TitleCode,
		FullName:                FullName,
		Street:                  Street,
		CountryCode:             CountryCode,
		StateCode:               StateCode,
		CityCode:                CityCode,
		City:                    City,
		NationalityCode:         NationalityCode,
		PostalCode:              PostalCode,
		Phone1:                  Phone1,
		Phone2:                  Phone2,
		Fax:                     Fax,
		Email:                   Email,
		Website:                 Website,
		CompanyCode:             CompanyCode,
		GuestTypeCode:           GuestTypeCode,
		PurposeOfCode:           PurposeOfCode,
		SalesCode:               SalesCode,
		CustomLookupFieldCode01: CustomLookupFieldCode01,
		CustomLookupFieldCode02: CustomLookupFieldCode02,
		ComplimentHu:            ComplimentHu,
		IsAdditional:            IsAdditional,
		IsScheduledRate:         IsScheduledRate,
		IsBreakfast:             IsBreakfast,
		PaxBreakfast:            PaxBreakfast,
		BreakfastVoucherNumber:  BreakfastVoucherNumber,
		Notes:                   Notes,
		CreatedBy:               CreatedBy,
	}
	result := DB.Table(DBVar.TableName.GuestInHouse).Create(&GuestInHouse)
	return result.Error
}

func UpdateGuestInHouse(DB *gorm.DB, AuditDate time.Time, FolioNumber uint64, GroupCode string, Adult int, Child int, RoomTypeCode string, BedTypeCode string, RoomNumber string, RoomRateCode string, RateOriginal float64, Rate float64, DiscountPercent uint8, Discount float64, BusinessSourceCode string, CommissionTypeCode string, CommissionValue float64, PaymentTypeCode string, MarketCode string, BookingSourceCode string, TitleCode string, FullName string, Street string, CountryCode string, StateCode string, CityCode string, City string, NationalityCode string, PostalCode string, Phone1 string, Phone2 string, Fax string, Email string, Website string, CompanyCode string, GuestTypeCode string, PurposeOfCode string, SalesCode string, CustomLookupFieldCode01 string, CustomLookupFieldCode02 string, ComplimentHu string, IsAdditional uint8, IsScheduledRate uint8, IsBreakfast uint8, Notes string, UpdatedBy string) error {
	var GuestInHouse = DBVar.Guest_in_house{
		GroupCode:               GroupCode,
		Adult:                   Adult,
		Child:                   Child,
		RoomTypeCode:            RoomTypeCode,
		BedTypeCode:             BedTypeCode,
		RoomNumber:              RoomNumber,
		RoomRateCode:            RoomRateCode,
		RateOriginal:            RateOriginal,
		Rate:                    Rate,
		DiscountPercent:         DiscountPercent,
		Discount:                Discount,
		BusinessSourceCode:      BusinessSourceCode,
		CommissionTypeCode:      CommissionTypeCode,
		CommissionValue:         CommissionValue,
		PaymentTypeCode:         PaymentTypeCode,
		MarketCode:              MarketCode,
		BookingSourceCode:       BookingSourceCode,
		TitleCode:               TitleCode,
		FullName:                FullName,
		Street:                  Street,
		CountryCode:             CountryCode,
		StateCode:               StateCode,
		CityCode:                CityCode,
		City:                    City,
		NationalityCode:         NationalityCode,
		PostalCode:              PostalCode,
		Phone1:                  Phone1,
		Phone2:                  Phone2,
		Fax:                     Fax,
		Email:                   Email,
		Website:                 Website,
		CompanyCode:             CompanyCode,
		GuestTypeCode:           GuestTypeCode,
		PurposeOfCode:           PurposeOfCode,
		SalesCode:               SalesCode,
		CustomLookupFieldCode01: CustomLookupFieldCode01,
		CustomLookupFieldCode02: CustomLookupFieldCode02,
		ComplimentHu:            ComplimentHu,
		IsAdditional:            IsAdditional,
		IsScheduledRate:         IsScheduledRate,
		IsBreakfast:             IsBreakfast,
		Notes:                   Notes,
		UpdatedBy:               UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.GuestInHouse).Where("audit_date=?", AuditDate).Where("folio_number=?", FolioNumber).Omit("created_at", "created_by", "updated_at", "id").Updates(&GuestInHouse)
	return result.Error
}

func UpdateGuestInHouseWithoutRate(DB *gorm.DB, AuditDate time.Time, FolioNumber uint64, GroupCode string, Adult int, Child int,
	RoomNumber string, BusinessSourceCode string, PaymentTypeCode string, MarketCode string, BookingSourceCode string,
	TitleCode string, FullName string, Street string, CountryCode string, StateCode string, CityCode string, City string, NationalityCode string,
	PostalCode string, Phone1 string, Phone2 string, Fax string, Email string, Website string, CompanyCode string, GuestTypeCode string, PurposeOfCode string,
	SalesCode string, CustomLookupFieldCode01 string, CustomLookupFieldCode02 string, IsAdditional uint8, Notes string, UpdatedBy string) error {
	var GuestInHouse = DBVar.Guest_in_house{
		GroupCode:               GroupCode,
		Adult:                   Adult,
		Child:                   Child,
		RoomNumber:              RoomNumber,
		BusinessSourceCode:      BusinessSourceCode,
		PaymentTypeCode:         PaymentTypeCode,
		MarketCode:              MarketCode,
		BookingSourceCode:       BookingSourceCode,
		TitleCode:               TitleCode,
		FullName:                FullName,
		Street:                  Street,
		CountryCode:             CountryCode,
		StateCode:               StateCode,
		CityCode:                CityCode,
		City:                    City,
		NationalityCode:         NationalityCode,
		PostalCode:              PostalCode,
		Phone1:                  Phone1,
		Phone2:                  Phone2,
		Fax:                     Fax,
		Email:                   Email,
		Website:                 Website,
		CompanyCode:             CompanyCode,
		GuestTypeCode:           GuestTypeCode,
		PurposeOfCode:           PurposeOfCode,
		SalesCode:               SalesCode,
		CustomLookupFieldCode01: CustomLookupFieldCode01,
		CustomLookupFieldCode02: CustomLookupFieldCode02,
		IsAdditional:            IsAdditional,
		Notes:                   Notes,
		UpdatedBy:               UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.GuestInHouse).Where("audit_date=?", AuditDate).Where("folio_number=?", FolioNumber).Omit("created_at", "created_by", "updated_at", "id").Updates(&GuestInHouse)
	return result.Error
}

func InsertGuestGeneral(DB *gorm.DB, PurposeOfCode string, SalesCode string, VoucherNumberTa string, FlightNumber string, FlightArrival time.Time, FlightDeparture time.Time, Notes string, ShowNotes uint8, HkNote string, DocumentNumber string, CreatedBy string) (Id uint64, err error) {
	var GuestGeneral = DBVar.Guest_general{
		PurposeOfCode:   General.PtrString(PurposeOfCode),
		SalesCode:       General.PtrString(SalesCode),
		VoucherNumberTa: General.PtrString(VoucherNumberTa),
		FlightNumber:    General.PtrString(FlightNumber),
		FlightArrival:   General.PtrTime(FlightArrival),
		FlightDeparture: General.PtrTime(FlightDeparture),
		Notes:           General.PtrString(Notes),
		ShowNotes:       General.PtrUint8(ShowNotes),
		HkNote:          General.PtrString(HkNote),
		DocumentNumber:  General.PtrString(DocumentNumber),
		CreatedBy:       CreatedBy,
	}
	result := DB.Table(DBVar.TableName.GuestGeneral).Create(&GuestGeneral)
	return GuestGeneral.Id, result.Error
}

func UpdateGuestGeneral(DB *gorm.DB, Id uint64, PurposeOfCode *string, SalesCode *string, VoucherNumberTa *string, FlightNumber *string, FlightArrival *time.Time, FlightDeparture *time.Time, Notes *string, ShowNotes *uint8, HkNote *string, DocumentNumber *string, UpdatedBy string) error {
	var GuestGeneral = DBVar.Guest_general{
		Id:              Id,
		PurposeOfCode:   PurposeOfCode,
		SalesCode:       SalesCode,
		VoucherNumberTa: VoucherNumberTa,
		FlightNumber:    FlightNumber,
		FlightArrival:   FlightArrival,
		FlightDeparture: FlightDeparture,
		Notes:           Notes,
		ShowNotes:       ShowNotes,
		HkNote:          HkNote,
		DocumentNumber:  DocumentNumber,
		UpdatedBy:       UpdatedBy,
	}
	result := DB.Model(&GuestGeneral).Where("id=?", Id).Omit("created_at", "created_by", "updated_at", "id").Updates(&GuestGeneral)
	return result.Error
}

func InsertCreditCard(DB *gorm.DB, GuestDepositId uint64, SubFolioId uint64, CardNumber string, CardHolder string, ValidMonth string, ValidYear string, CreatedBy string) error {
	var CreditCard = DBVar.Credit_card{
		GuestDepositId: GuestDepositId,
		SubFolioId:     SubFolioId,
		CardNumber:     CardNumber,
		CardHolder:     CardHolder,
		ValidMonth:     ValidMonth,
		ValidYear:      ValidYear,
		CreatedBy:      CreatedBy,
	}
	result := DB.Table(DBVar.TableName.CreditCard).Create(&CreditCard)
	return result.Error
}

func UpdateCreditCard(DB *gorm.DB, GuestDepositId uint64, SubFolioId uint64, CardNumber string, CardHolder string, ValidMonth string, ValidYear string, UpdatedBy string) error {
	var CreditCard = DBVar.Credit_card{
		GuestDepositId: GuestDepositId,
		SubFolioId:     SubFolioId,
		CardNumber:     CardNumber,
		CardHolder:     CardHolder,
		ValidMonth:     ValidMonth,
		ValidYear:      ValidYear,
		UpdatedBy:      UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.CreditCard).Omit("created_at", "created_by", "updated_at", "id").Updates(&CreditCard)
	return result.Error
}

func InsertGuestDeposit(c *gin.Context, DB *gorm.DB, ReservationNumber uint64, IDCorrected uint64, SubDepartmentCode string, AccountCode string, Amount float64, DefaultCurrencyCode string, ExchangeRate float64, CurrencyCode string, AuditDate time.Time, Remark string, DocumentNumber string, TypeCode string, CardBankCode string, CardTypeCode string, IsCorrection uint8, CorrectionBy string, CorrectionReason string, CorrectionBreakdown uint64, Shift string, LogShiftId uint64, SystemCode string, CreatedBy string) (Id uint64, err error) {
	AmountForeign := Amount
	if AuditDate.IsZero() {
		AuditDate = GetAuditDate(c, DB, false)
	}
	if CurrencyCode == "" {
		CurrencyCode = GetDefaultCurrencyCode(DB)
		ExchangeRate = GetExchangeRateCurrency(DB, CurrencyCode)
	}

	if DefaultCurrencyCode == "" {
		DefaultCurrencyCode = GetDefaultCurrencyCode(DB)
	}

	if CurrencyCode != DefaultCurrencyCode {
		Amount = General.RoundToX3(Amount * ExchangeRate)
	}

	if CorrectionBreakdown == 0 {
		CorrectionBreakdown = GetGuestDepositCorrectionBreakDown(DB)
	}
	fmt.Println("doc", DocumentNumber)

	var GuestDeposit = DBVar.Guest_deposit{
		ReservationNumber:   ReservationNumber,
		SubDepartmentCode:   SubDepartmentCode,
		AccountCode:         AccountCode,
		Amount:              Amount,
		DefaultCurrencyCode: DefaultCurrencyCode,
		AmountForeign:       AmountForeign,
		ExchangeRate:        ExchangeRate,
		CurrencyCode:        CurrencyCode,
		AuditDate:           AuditDate,
		Remark:              Remark,
		DocumentNumber:      DocumentNumber,
		TypeCode:            TypeCode,
		CardBankCode:        CardBankCode,
		CardTypeCode:        CardTypeCode,
		IsCorrection:        IsCorrection,
		CorrectionBy:        CorrectionBy,
		CorrectionReason:    CorrectionReason,
		CorrectionBreakdown: CorrectionBreakdown,
		Shift:               Shift,
		LogShiftId:          LogShiftId,
		SystemCode:          SystemCode,
		CreatedBy:           CreatedBy,
	}
	result := DB.Table(DBVar.TableName.GuestDeposit).Create(&GuestDeposit)

	// Insert Foreign Cash
	if (GetAccountSubGroupCode(DB, AccountCode) == GlobalVar.GlobalAccountSubGroup.Payment || GetAccountSubGroupCode(DB, AccountCode) == GlobalVar.GlobalAccountSubGroup.CreditDebitCard || GetAccountSubGroupCode(DB, AccountCode) == GlobalVar.GlobalAccountSubGroup.BankTransfer) && (CurrencyCode != DefaultCurrencyCode) {
		var RemarkForeignCash, TypeCodeX string
		if TypeCode == GlobalVar.TransactionType.Debit {
			TypeCodeX = GlobalVar.TransactionType.Credit
		} else {
			TypeCodeX = GlobalVar.TransactionType.Debit
		}
		if IsCorrection > 0 {
			RemarkForeignCash = "Guest Deposit Correction for Reservation: " + General.Uint64ToStr(ReservationNumber) + ", Doc#: " + DocumentNumber
		} else {
			RemarkForeignCash = "Guest Deposit for Reservation: " + General.Uint64ToStr(ReservationNumber) + ", Doc#: " + DocumentNumber
		}
		if err := InsertAccForeignCash(DB, GuestDeposit.Id, IDCorrected, 0, GlobalVar.ForeignCashTableID.GuestDeposit, 0, "", AuditDate, TypeCodeX, Amount, DefaultCurrencyCode, AmountForeign, ExchangeRate, CurrencyCode, RemarkForeignCash, IsCorrection, CreatedBy); err != nil {
			return 0, err
		}
	}
	return GuestDeposit.Id, result.Error
}

func UpdateGuestDeposit(DB *gorm.DB, ReservationNumber uint64, SubDepartmentCode string, AccountCode string, Amount float64, DefaultCurrencyCode string, AmountForeign float64, ExchangeRate float64, CurrencyCode string, AuditDate time.Time, Remark string, DocumentNumber string, TypeCode string, CardBankCode string, CardTypeCode string, RefNumber uint64, Void uint8, VoidDate time.Time, VoidBy string, VoidReason string, IsCorrection uint8, CorrectionBy string, CorrectionReason string, CorrectionBreakdown uint64, Shift string, LogShiftId uint64, IsPairWithFolio uint8, TransferPairId uint64, SystemCode string, UpdatedBy string) error {
	var GuestDeposit = DBVar.Guest_deposit{
		ReservationNumber:   ReservationNumber,
		SubDepartmentCode:   SubDepartmentCode,
		AccountCode:         AccountCode,
		Amount:              Amount,
		DefaultCurrencyCode: DefaultCurrencyCode,
		AmountForeign:       AmountForeign,
		ExchangeRate:        ExchangeRate,
		CurrencyCode:        CurrencyCode,
		AuditDate:           AuditDate,
		Remark:              Remark,
		DocumentNumber:      DocumentNumber,
		TypeCode:            TypeCode,
		CardBankCode:        CardBankCode,
		CardTypeCode:        CardTypeCode,
		RefNumber:           RefNumber,
		Void:                Void,
		VoidDate:            VoidDate,
		VoidBy:              VoidBy,
		VoidReason:          VoidReason,
		IsCorrection:        IsCorrection,
		CorrectionBy:        CorrectionBy,
		CorrectionReason:    CorrectionReason,
		CorrectionBreakdown: CorrectionBreakdown,
		Shift:               Shift,
		LogShiftId:          LogShiftId,
		IsPairWithFolio:     IsPairWithFolio,
		TransferPairId:      TransferPairId,
		SystemCode:          SystemCode,
		UpdatedBy:           UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.GuestDeposit).Omit("created_at", "created_by", "updated_at", "id").Updates(&GuestDeposit)
	return result.Error
}

func InsertReservationScheduledRate(DB *gorm.DB, ReservationNumber uint64, FromDate time.Time, ToDate time.Time, RoomRateCode string, Rate float64, ComplimentHu string, CreatedBy string) error {
	var ReservationScheduledRate = DBVar.Reservation_scheduled_rate{
		ReservationNumber: ReservationNumber,
		FromDate:          FromDate,
		ToDate:            ToDate,
		RoomRateCode:      RoomRateCode,
		Rate:              &Rate,
		ComplimentHu:      ComplimentHu,
		CreatedBy:         CreatedBy,
	}
	result := DB.Table(DBVar.TableName.ReservationScheduledRate).Create(&ReservationScheduledRate)
	return result.Error
}

func UpdateReservationScheduledRate(DB *gorm.DB, Id uint64, FromDate time.Time, ToDate time.Time, RoomRateCode string, Rate float64, ComplimentHu string, UpdatedBy string) error {
	var ReservationScheduledRate = DBVar.Reservation_scheduled_rate{
		Id:           Id,
		FromDate:     FromDate,
		ToDate:       ToDate,
		RoomRateCode: RoomRateCode,
		Rate:         &Rate,
		ComplimentHu: ComplimentHu,
		UpdatedBy:    UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.ReservationScheduledRate).Omit("created_at", "created_by", "updated_at", "id").Updates(&ReservationScheduledRate)
	return result.Error
}

func InsertGuestGroup(DB *gorm.DB, Code string, Name string, ContactPerson *string, Street *string, CountryCode *string, StateCode *string, CityCode *string, City *string, PostalCode *string, Phone1 *string, Phone2 *string, Fax *string, Email *string, Website *string, IsActive uint8, CreatedBy string) error {
	var GuestGroup = DBVar.Guest_group{
		Code:          Code,
		Name:          Name,
		ContactPerson: ContactPerson,
		Street:        Street,
		CountryCode:   CountryCode,
		StateCode:     StateCode,
		CityCode:      CityCode,
		City:          City,
		PostalCode:    PostalCode,
		Phone1:        Phone1,
		Phone2:        Phone2,
		Fax:           Fax,
		Email:         Email,
		Website:       Website,
		IsActive:      IsActive,
		CreatedBy:     CreatedBy,
	}
	result := DB.Table(DBVar.TableName.GuestGroup).Create(&GuestGroup)
	return result.Error
}

func UpdateGuestGroup(DB *gorm.DB, Id uint64, Name string, ContactPerson *string, Street *string, CountryCode *string, StateCode *string, CityCode *string, City *string, PostalCode *string, Phone1 *string, Phone2 *string, Fax *string, Email *string, Website *string, UpdatedBy string) error {
	var GuestGroup = DBVar.Guest_group{
		Id:            Id,
		Name:          Name,
		ContactPerson: ContactPerson,
		Street:        Street,
		CountryCode:   CountryCode,
		StateCode:     StateCode,
		CityCode:      CityCode,
		City:          City,
		PostalCode:    PostalCode,
		Phone1:        Phone1,
		Phone2:        Phone2,
		Fax:           Fax,
		Email:         Email,
		Website:       Website,
		UpdatedBy:     UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.GuestGroup).Omit("created_at", "created_by", "updated_at", "id").Updates(&GuestGroup)
	return result.Error
}

func InsertFolio(DB *gorm.DB, Dataset *GlobalVar.TDataset, TypeCode string, CoNumber string, ReservationNumber uint64, ContactPersonId1 uint64, ContactPersonId2 uint64, ContactPersonId3 uint64, ContactPersonId4 uint64, GuestDetailId uint64, GuestProfileId1 uint64, GuestProfileId2 uint64, GuestProfileId3 uint64, GuestProfileId4 uint64, GuestGeneralId uint64, StatusCode string, GroupCode string, ComplimentHu string, IsFromAllotment uint8, IsIncognito uint8, CreatedBy string) (FolioNumberX uint64, err error) {
	var LockOnCheckIn uint8
	if General.StrToBool(Dataset.Configuration[GlobalVar.ConfigurationCategory.Reservation][GlobalVar.ConfigurationName.LockFolioOnCheckIn].(string)) {
		LockOnCheckIn = 1
	}

	var Folio = DBVar.Folio{
		TypeCode:          TypeCode,
		CoNumber:          CoNumber,
		ReservationNumber: ReservationNumber,
		ContactPersonId1:  ContactPersonId1,
		ContactPersonId2:  ContactPersonId2,
		ContactPersonId3:  ContactPersonId3,
		ContactPersonId4:  ContactPersonId4,
		GuestDetailId:     GuestDetailId,
		GuestProfileId1:   GuestProfileId1,
		GuestProfileId2:   GuestProfileId2,
		GuestProfileId3:   GuestProfileId3,
		GuestProfileId4:   GuestProfileId4,
		GuestGeneralId:    GuestGeneralId,
		GroupCode:         GroupCode,
		StatusCode:        StatusCode,
		// VoucherNumber:     VoucherNumber,
		ComplimentHu:    ComplimentHu,
		IsLock:          LockOnCheckIn,
		IsIncognito:     IsIncognito,
		IsFromAllotment: IsFromAllotment,
		SystemCode:      "",
		CreatedBy:       CreatedBy,
	}
	result := DB.Table(DBVar.TableName.Folio).Create(&Folio)
	return Folio.Number, result.Error
}

func InsertFolioClose(DB *gorm.DB, TypeCode string, CoNumber string, ReservationNumber uint64, ContactPersonId1 uint64, ContactPersonId2 uint64, ContactPersonId3 uint64, ContactPersonId4 uint64, GuestDetailId uint64, GuestProfileId1 uint64, GuestProfileId2 uint64, GuestProfileId3 uint64, GuestProfileId4 uint64, GuestGeneralId uint64, GroupCode string, ComplimentHu string, IsFromAllotment uint8, CreatedBy string) (FolioNumberX uint64, err error) {

	var Folio = DBVar.Folio{
		TypeCode:          TypeCode,
		CoNumber:          CoNumber,
		ReservationNumber: ReservationNumber,
		ContactPersonId1:  ContactPersonId1,
		ContactPersonId2:  ContactPersonId2,
		ContactPersonId3:  ContactPersonId3,
		ContactPersonId4:  ContactPersonId4,
		GuestDetailId:     GuestDetailId,
		GuestProfileId1:   GuestProfileId1,
		GuestProfileId2:   GuestProfileId2,
		GuestProfileId3:   GuestProfileId3,
		GuestProfileId4:   GuestProfileId4,
		GuestGeneralId:    GuestGeneralId,
		GroupCode:         GroupCode,
		StatusCode:        GlobalVar.FolioStatus.Closed,
		// VoucherNumber:     VoucherNumber,
		ComplimentHu:    ComplimentHu,
		SystemCode:      "",
		IsFromAllotment: IsFromAllotment,
		CreatedBy:       CreatedBy,
	}
	result := DB.Table(DBVar.TableName.Folio).Create(&Folio)
	return Folio.Number, result.Error
}

func UpdateFolio(DB *gorm.DB, Number uint64, TypeCode string, CoNumber string, ReservationNumber uint64, ContactPersonId1 uint64, ContactPersonId2 uint64, ContactPersonId3 uint64, ContactPersonId4 uint64, GuestDetailId uint64, GuestProfileId1 uint64, GuestProfileId2 uint64, GuestProfileId3 uint64, GuestProfileId4 uint64, GuestGeneralId uint64, GroupCode string, RoomStatusCode string, StatusCode string, VoucherNumber string, ComplimentHu string, UpdatedBy string) error {
	var Folio = DBVar.Folio{
		Number:            Number,
		TypeCode:          TypeCode,
		CoNumber:          CoNumber,
		ReservationNumber: ReservationNumber,
		ContactPersonId1:  ContactPersonId1,
		ContactPersonId2:  ContactPersonId2,
		ContactPersonId3:  ContactPersonId3,
		ContactPersonId4:  ContactPersonId4,
		GuestDetailId:     GuestDetailId,
		GuestProfileId1:   GuestProfileId1,
		GuestProfileId2:   GuestProfileId2,
		GuestProfileId3:   GuestProfileId3,
		GuestProfileId4:   GuestProfileId4,
		GuestGeneralId:    GuestGeneralId,
		GroupCode:         GroupCode,
		RoomStatusCode:    RoomStatusCode,
		StatusCode:        StatusCode,
		VoucherNumber:     VoucherNumber,
		ComplimentHu:      ComplimentHu,
		UpdatedBy:         UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.Folio).Omit("created_at", "created_by", "updated_at").Updates(&Folio)
	return result.Error
}

func InsertGuestScheduledRate(DB *gorm.DB, FolioNumber uint64, FromDate time.Time, ToDate time.Time, RoomRateCode string, Rate float64, ComplimentHu string, CreatedBy string) error {
	var GuestScheduledRate = DBVar.Guest_scheduled_rate{
		FolioNumber:  FolioNumber,
		FromDate:     FromDate,
		ToDate:       ToDate,
		RoomRateCode: RoomRateCode,
		Rate:         &Rate,
		ComplimentHu: ComplimentHu,
		CreatedBy:    CreatedBy,
	}
	result := DB.Table(DBVar.TableName.GuestScheduledRate).Create(&GuestScheduledRate)
	return result.Error
}

func UpdateGuestScheduledRate(DB *gorm.DB, Id uint64, FromDate time.Time, ToDate time.Time, RoomRateCode string, Rate float64, ComplimentHu string, UpdatedBy string) error {
	var GuestScheduledRate = DBVar.Guest_scheduled_rate{
		Id:           Id,
		FromDate:     FromDate,
		ToDate:       ToDate,
		RoomRateCode: RoomRateCode,
		Rate:         &Rate,
		ComplimentHu: ComplimentHu,
		UpdatedBy:    UpdatedBy,
	}
	result := DB.Debug().Table(DBVar.TableName.GuestScheduledRate).Omit("created_at", "created_by", "updated_at", "id").Updates(&GuestScheduledRate)
	return result.Error
}

func InsertReservationExtraCharge(DB *gorm.DB, ReservationNumber uint64, PackageName string, OutletCode string, ProductCode string, PackageCode string, GroupCode string, SubDepartmentCode string, AccountCode string, Quantity float64, Amount float64, PerPax uint8, IncludeChild uint8, TaxAndServiceCode string, ChargeFrequencyCode string, MaxPax int, ExtraPax float64, PerPaxExtra uint8, CreatedBy string) (uint64, error) {
	var ReservationExtraCharge = DBVar.Reservation_extra_charge{
		ReservationNumber:   ReservationNumber,
		PackageName:         &PackageName,
		OutletCode:          &OutletCode,
		ProductCode:         &ProductCode,
		PackageCode:         &PackageCode,
		GroupCode:           GroupCode,
		SubDepartmentCode:   SubDepartmentCode,
		AccountCode:         AccountCode,
		Quantity:            Quantity,
		Amount:              Amount,
		PerPax:              &PerPax,
		IncludeChild:        &IncludeChild,
		TaxAndServiceCode:   &TaxAndServiceCode,
		ChargeFrequencyCode: ChargeFrequencyCode,
		MaxPax:              MaxPax,
		ExtraPax:            &ExtraPax,
		PerPaxExtra:         &PerPaxExtra,
		CreatedBy:           CreatedBy,
	}
	result := DB.Table(DBVar.TableName.ReservationExtraCharge).Create(&ReservationExtraCharge)
	return ReservationExtraCharge.Id, result.Error
}

func UpdateReservationExtraCharge(DB *gorm.DB, Id uint64, PackageName string, OutletCode string, ProductCode string, PackageCode string, GroupCode string, SubDepartmentCode string, AccountCode string, Quantity float64, Amount float64, PerPax uint8, IncludeChild uint8, TaxAndServiceCode string, ChargeFrequencyCode string, MaxPax int, ExtraPax float64, PerPaxExtra uint8, UpdatedBy string) error {
	var ReservationExtraCharge = DBVar.Reservation_extra_charge{
		Id:                  Id,
		PackageName:         &PackageName,
		OutletCode:          &OutletCode,
		ProductCode:         &ProductCode,
		PackageCode:         &PackageCode,
		GroupCode:           GroupCode,
		SubDepartmentCode:   SubDepartmentCode,
		AccountCode:         AccountCode,
		Quantity:            Quantity,
		Amount:              Amount,
		PerPax:              &PerPax,
		IncludeChild:        &IncludeChild,
		TaxAndServiceCode:   &TaxAndServiceCode,
		ChargeFrequencyCode: ChargeFrequencyCode,
		MaxPax:              MaxPax,
		ExtraPax:            &ExtraPax,
		PerPaxExtra:         &PerPaxExtra,
		UpdatedBy:           UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.ReservationExtraCharge).Omit("created_at", "created_by", "updated_at", "id").Updates(&ReservationExtraCharge)
	return result.Error
}

func InsertReservationExtraChargeBreakdown(DB *gorm.DB, ReservationExtraChargeId uint64, OutletCode string, ProductCode string, SubDepartmentCode string, AccountCode string, CompanyCode string, Quantity float64, IsAmountPercent uint8, Amount float64, PerPax uint8, IncludeChild uint8, Remark string, TaxAndServiceCode string, ChargeFrequencyCode string, MaxPax int, ExtraPax float64, PerPaxExtra uint8, CreatedBy string) error {
	var ReservationExtraChargeBreakdown = DBVar.Reservation_extra_charge_breakdown{
		ReservationExtraChargeId: ReservationExtraChargeId,
		OutletCode:               &OutletCode,
		ProductCode:              &ProductCode,
		SubDepartmentCode:        SubDepartmentCode,
		AccountCode:              AccountCode,
		CompanyCode:              &CompanyCode,
		Quantity:                 Quantity,
		IsAmountPercent:          &IsAmountPercent,
		Amount:                   Amount,
		PerPax:                   &PerPax,
		IncludeChild:             &IncludeChild,
		Remark:                   &Remark,
		TaxAndServiceCode:        &TaxAndServiceCode,
		ChargeFrequencyCode:      ChargeFrequencyCode,
		MaxPax:                   MaxPax,
		ExtraPax:                 &ExtraPax,
		PerPaxExtra:              &PerPaxExtra,
		CreatedBy:                CreatedBy,
	}
	result := DB.Table(DBVar.TableName.ReservationExtraChargeBreakdown).Create(&ReservationExtraChargeBreakdown)
	return result.Error
}

func UpdateReservationExtraChargeBreakdown(DB *gorm.DB, Id uint64, OutletCode string, ProductCode string, SubDepartmentCode string, AccountCode string, CompanyCode string, Quantity float64, IsAmountPercent uint8, Amount float64, PerPax uint8, IncludeChild uint8, Remark string, TaxAndServiceCode string, ChargeFrequencyCode string, MaxPax int, ExtraPax float64, PerPaxExtra uint8, UpdatedBy string) error {
	var ReservationExtraChargeBreakdown = DBVar.Reservation_extra_charge_breakdown{
		Id:                  Id,
		OutletCode:          &OutletCode,
		ProductCode:         &ProductCode,
		SubDepartmentCode:   SubDepartmentCode,
		AccountCode:         AccountCode,
		CompanyCode:         &CompanyCode,
		Quantity:            Quantity,
		IsAmountPercent:     &IsAmountPercent,
		Amount:              Amount,
		PerPax:              &PerPax,
		IncludeChild:        &IncludeChild,
		Remark:              &Remark,
		TaxAndServiceCode:   &TaxAndServiceCode,
		ChargeFrequencyCode: ChargeFrequencyCode,
		MaxPax:              MaxPax,
		ExtraPax:            &ExtraPax,
		PerPaxExtra:         &PerPaxExtra,
		UpdatedBy:           UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.ReservationExtraChargeBreakdown).Omit("created_at", "created_by", "updated_at", "id").Updates(&ReservationExtraChargeBreakdown)
	return result.Error
}

func InsertGuestExtraCharge(DB *gorm.DB, FolioNumber uint64, PackageName string, OutletCode string, ProductCode string, PackageCode string, GroupCode string, SubDepartmentCode string, AccountCode string, Quantity float64, Amount float64, PerPax uint8, IncludeChild uint8, TaxAndServiceCode string, ChargeFrequencyCode string, MaxPax int, ExtraPax float64, PerPaxExtra uint8, CreatedBy string) (uint64, error) {
	var GuestExtraCharge = DBVar.Guest_extra_charge{
		FolioNumber:         FolioNumber,
		PackageName:         &PackageName,
		OutletCode:          &OutletCode,
		ProductCode:         &ProductCode,
		PackageCode:         &PackageCode,
		GroupCode:           GroupCode,
		SubDepartmentCode:   SubDepartmentCode,
		AccountCode:         AccountCode,
		Quantity:            Quantity,
		Amount:              Amount,
		PerPax:              &PerPax,
		IncludeChild:        &IncludeChild,
		TaxAndServiceCode:   &TaxAndServiceCode,
		ChargeFrequencyCode: ChargeFrequencyCode,
		MaxPax:              MaxPax,
		ExtraPax:            &ExtraPax,
		PerPaxExtra:         &PerPaxExtra,
		CreatedBy:           CreatedBy,
	}
	result := DB.Table(DBVar.TableName.GuestExtraCharge).Create(&GuestExtraCharge)
	return GuestExtraCharge.Id, result.Error
}

func UpdateGuestExtraCharge(DB *gorm.DB, Id uint64, PackageName string, OutletCode string, ProductCode string, PackageCode string, GroupCode string, SubDepartmentCode string, AccountCode string, Quantity float64, Amount float64, PerPax uint8, IncludeChild uint8, TaxAndServiceCode string, ChargeFrequencyCode string, MaxPax int, ExtraPax float64, PerPaxExtra uint8, UpdatedBy string) error {
	var GuestExtraCharge = DBVar.Guest_extra_charge{
		Id:                  Id,
		PackageName:         &PackageName,
		OutletCode:          &OutletCode,
		ProductCode:         &ProductCode,
		PackageCode:         &PackageCode,
		GroupCode:           GroupCode,
		SubDepartmentCode:   SubDepartmentCode,
		AccountCode:         AccountCode,
		Quantity:            Quantity,
		Amount:              Amount,
		PerPax:              &PerPax,
		IncludeChild:        &IncludeChild,
		TaxAndServiceCode:   &TaxAndServiceCode,
		ChargeFrequencyCode: ChargeFrequencyCode,
		MaxPax:              MaxPax,
		ExtraPax:            &ExtraPax,
		PerPaxExtra:         &PerPaxExtra,
		UpdatedBy:           UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.GuestExtraCharge).Omit("created_at", "created_by", "updated_at", "id").Updates(&GuestExtraCharge)
	return result.Error
}

func InsertGuestExtraChargeBreakdown(DB *gorm.DB, GuestExtraChargeId uint64, OutletCode string, ProductCode string, SubDepartmentCode string, AccountCode string, CompanyCode string, Quantity float64, IsAmountPercent uint8, Amount float64, PerPax uint8, IncludeChild uint8, Remark string, TaxAndServiceCode string, ChargeFrequencyCode string, MaxPax int, ExtraPax float64, PerPaxExtra uint8, CreatedBy string) error {
	var GuestExtraChargeBreakdown = DBVar.Guest_extra_charge_breakdown{
		GuestExtraChargeId:  GuestExtraChargeId,
		OutletCode:          &OutletCode,
		ProductCode:         &ProductCode,
		SubDepartmentCode:   SubDepartmentCode,
		AccountCode:         AccountCode,
		CompanyCode:         &CompanyCode,
		Quantity:            Quantity,
		IsAmountPercent:     &IsAmountPercent,
		Amount:              Amount,
		PerPax:              &PerPax,
		IncludeChild:        &IncludeChild,
		Remark:              &Remark,
		TaxAndServiceCode:   &TaxAndServiceCode,
		ChargeFrequencyCode: ChargeFrequencyCode,
		MaxPax:              MaxPax,
		ExtraPax:            &ExtraPax,
		PerPaxExtra:         &PerPaxExtra,
		CreatedBy:           CreatedBy,
	}
	result := DB.Table(DBVar.TableName.GuestExtraChargeBreakdown).Create(&GuestExtraChargeBreakdown)
	return result.Error
}

func UpdateGuestExtraChargeBreakdown(DB *gorm.DB, Id uint64, OutletCode string, ProductCode string, SubDepartmentCode string, AccountCode string, CompanyCode string, Quantity float64, IsAmountPercent uint8, Amount float64, PerPax uint8, IncludeChild uint8, Remark string, TaxAndServiceCode string, ChargeFrequencyCode string, MaxPax int, ExtraPax float64, PerPaxExtra uint8, UpdatedBy string) error {
	var GuestExtraChargeBreakdown = DBVar.Guest_extra_charge_breakdown{
		Id:                  Id,
		OutletCode:          &OutletCode,
		ProductCode:         &ProductCode,
		SubDepartmentCode:   SubDepartmentCode,
		AccountCode:         AccountCode,
		CompanyCode:         &CompanyCode,
		Quantity:            Quantity,
		IsAmountPercent:     &IsAmountPercent,
		Amount:              Amount,
		PerPax:              &PerPax,
		IncludeChild:        &IncludeChild,
		Remark:              &Remark,
		TaxAndServiceCode:   &TaxAndServiceCode,
		ChargeFrequencyCode: ChargeFrequencyCode,
		MaxPax:              MaxPax,
		ExtraPax:            &ExtraPax,
		PerPaxExtra:         &PerPaxExtra,
		UpdatedBy:           UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.GuestExtraChargeBreakdown).Omit("created_at", "created_by", "updated_at", "id").Updates(&GuestExtraChargeBreakdown)
	return result.Error
}

func InsertCfgInitRoomRate(DB *gorm.DB, Code string, Name string, RoomTypeCode string, FromDate time.Time, ToDate time.Time, SubCategoryCode string, CompanyCode string, MarketCode string, DynamicRateTypeCode string, IsLastDeal uint8, IsRateStructure uint8, IsCompliment uint8, IncludeBreakfast uint8, WeekdayRate1 float64, WeekdayRate2 float64, WeekdayRate3 float64, WeekdayRate4 float64, WeekendRate1 float64, WeekendRate2 float64, WeekendRate3 float64, WeekendRate4 float64, WeekdayRateChild1 float64, WeekdayRateChild2 float64, WeekdayRateChild3 float64, WeekdayRateChild4 float64, WeekendRateChild1 float64, WeekendRateChild2 float64, WeekendRateChild3 float64, WeekendRateChild4 float64, TaxAndServiceCode string, ChargeFrequencyCode string, ExtraPax float64, PerPax uint8, IncludeChild uint8, Day1 uint8, Day2 uint8, Day3 uint8, Day4 uint8, Day5 uint8, Day6 uint8, Day7 uint8, Notes string, IdSort int, IsActive uint8, CmInvCode string, CmStopSell uint8, IsCmUpdated uint8, IsCmUpdatedInclusion uint8, CmStartDate time.Time, CmEndDate time.Time, IsSent uint8, IsOnline uint8, CreatedBy string) error {
	var CfgInitRoomRate = DBVar.Cfg_init_room_rate{
		Code:                 Code,
		Name:                 Name,
		RoomTypeCode:         RoomTypeCode,
		FromDate:             FromDate,
		ToDate:               ToDate,
		SubCategoryCode:      SubCategoryCode,
		CompanyCode:          CompanyCode,
		MarketCode:           MarketCode,
		DynamicRateTypeCode:  DynamicRateTypeCode,
		IsLastDeal:           IsLastDeal,
		IsRateStructure:      IsRateStructure,
		IsCompliment:         IsCompliment,
		IncludeBreakfast:     IncludeBreakfast,
		WeekdayRate1:         WeekdayRate1,
		WeekdayRate2:         WeekdayRate2,
		WeekdayRate3:         WeekdayRate3,
		WeekdayRate4:         WeekdayRate4,
		WeekendRate1:         WeekendRate1,
		WeekendRate2:         WeekendRate2,
		WeekendRate3:         WeekendRate3,
		WeekendRate4:         WeekendRate4,
		WeekdayRateChild1:    WeekdayRateChild1,
		WeekdayRateChild2:    WeekdayRateChild2,
		WeekdayRateChild3:    WeekdayRateChild3,
		WeekdayRateChild4:    WeekdayRateChild4,
		WeekendRateChild1:    WeekendRateChild1,
		WeekendRateChild2:    WeekendRateChild2,
		WeekendRateChild3:    WeekendRateChild3,
		WeekendRateChild4:    WeekendRateChild4,
		TaxAndServiceCode:    TaxAndServiceCode,
		ChargeFrequencyCode:  ChargeFrequencyCode,
		ExtraPax:             ExtraPax,
		PerPax:               PerPax,
		IncludeChild:         IncludeChild,
		Day1:                 Day1,
		Day2:                 Day2,
		Day3:                 Day3,
		Day4:                 Day4,
		Day5:                 Day5,
		Day6:                 Day6,
		Day7:                 Day7,
		Notes:                Notes,
		IdSort:               IdSort,
		IsActive:             IsActive,
		CmInvCode:            CmInvCode,
		CmStopSell:           CmStopSell,
		IsCmUpdated:          IsCmUpdated,
		IsCmUpdatedInclusion: IsCmUpdatedInclusion,
		IsSent:               IsSent,
		IsOnline:             IsOnline,
		CreatedBy:            CreatedBy,
	}
	result := DB.Table(DBVar.TableName.CfgInitRoomRate).Create(&CfgInitRoomRate)
	return result.Error
}

func UpdateCfgInitRoomRate(DB *gorm.DB, Code string, Name string, RoomTypeCode string, FromDate time.Time, ToDate time.Time, SubCategoryCode string, CompanyCode string, MarketCode string, DynamicRateTypeCode string, IsLastDeal uint8, IsRateStructure uint8, IsCompliment uint8, IncludeBreakfast uint8, WeekdayRate1 float64, WeekdayRate2 float64, WeekdayRate3 float64, WeekdayRate4 float64, WeekendRate1 float64, WeekendRate2 float64, WeekendRate3 float64, WeekendRate4 float64, WeekdayRateChild1 float64, WeekdayRateChild2 float64, WeekdayRateChild3 float64, WeekdayRateChild4 float64, WeekendRateChild1 float64, WeekendRateChild2 float64, WeekendRateChild3 float64, WeekendRateChild4 float64, TaxAndServiceCode string, ChargeFrequencyCode string, ExtraPax float64, PerPax uint8, IncludeChild uint8, Day1 uint8, Day2 uint8, Day3 uint8, Day4 uint8, Day5 uint8, Day6 uint8, Day7 uint8, Notes string, IdSort int, IsActive uint8, CmInvCode string, CmStopSell uint8, IsCmUpdated uint8, IsCmUpdatedInclusion uint8, CmStartDate time.Time, CmEndDate time.Time, IsSent uint8, IsOnline uint8, UpdatedBy string) error {
	var CfgInitRoomRate = DBVar.Cfg_init_room_rate{
		Code:                 Code,
		Name:                 Name,
		RoomTypeCode:         RoomTypeCode,
		FromDate:             FromDate,
		ToDate:               ToDate,
		SubCategoryCode:      SubCategoryCode,
		CompanyCode:          CompanyCode,
		MarketCode:           MarketCode,
		DynamicRateTypeCode:  DynamicRateTypeCode,
		IsLastDeal:           IsLastDeal,
		IsRateStructure:      IsRateStructure,
		IsCompliment:         IsCompliment,
		IncludeBreakfast:     IncludeBreakfast,
		WeekdayRate1:         WeekdayRate1,
		WeekdayRate2:         WeekdayRate2,
		WeekdayRate3:         WeekdayRate3,
		WeekdayRate4:         WeekdayRate4,
		WeekendRate1:         WeekendRate1,
		WeekendRate2:         WeekendRate2,
		WeekendRate3:         WeekendRate3,
		WeekendRate4:         WeekendRate4,
		WeekdayRateChild1:    WeekdayRateChild1,
		WeekdayRateChild2:    WeekdayRateChild2,
		WeekdayRateChild3:    WeekdayRateChild3,
		WeekdayRateChild4:    WeekdayRateChild4,
		WeekendRateChild1:    WeekendRateChild1,
		WeekendRateChild2:    WeekendRateChild2,
		WeekendRateChild3:    WeekendRateChild3,
		WeekendRateChild4:    WeekendRateChild4,
		TaxAndServiceCode:    TaxAndServiceCode,
		ChargeFrequencyCode:  ChargeFrequencyCode,
		ExtraPax:             ExtraPax,
		PerPax:               PerPax,
		IncludeChild:         IncludeChild,
		Day1:                 Day1,
		Day2:                 Day2,
		Day3:                 Day3,
		Day4:                 Day4,
		Day5:                 Day5,
		Day6:                 Day6,
		Day7:                 Day7,
		Notes:                Notes,
		IdSort:               IdSort,
		IsActive:             IsActive,
		CmInvCode:            CmInvCode,
		CmStopSell:           CmStopSell,
		IsCmUpdated:          IsCmUpdated,
		IsCmUpdatedInclusion: IsCmUpdatedInclusion,
		IsSent:               IsSent,
		IsOnline:             IsOnline,
		UpdatedBy:            UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.CfgInitRoomRate).Omit("created_at", "created_by", "updated_at", "id").Updates(&CfgInitRoomRate)
	return result.Error
}

func InsertCfgInitRoomRateBreakdown(DB *gorm.DB, RoomRateCode string, OutletCode string, ProductCode string, SubDepartmentCode string, AccountCode string, CompanyCode string, Quantity float64, IsAmountPercent uint8, Amount float64, PerPax uint8, IncludeChild uint8, Remark string, TaxAndServiceCode string, ChargeFrequencyCode string, MaxPax int, ExtraPax float64, PerPaxExtra uint8, CreatedBy string) error {
	var CfgInitRoomRateBreakdown = DBVar.Cfg_init_room_rate_breakdown{
		RoomRateCode:        RoomRateCode,
		OutletCode:          OutletCode,
		ProductCode:         ProductCode,
		SubDepartmentCode:   SubDepartmentCode,
		AccountCode:         AccountCode,
		CompanyCode:         CompanyCode,
		Quantity:            Quantity,
		IsAmountPercent:     IsAmountPercent,
		Amount:              Amount,
		PerPax:              PerPax,
		IncludeChild:        IncludeChild,
		Remark:              Remark,
		TaxAndServiceCode:   TaxAndServiceCode,
		ChargeFrequencyCode: ChargeFrequencyCode,
		MaxPax:              MaxPax,
		ExtraPax:            ExtraPax,
		PerPaxExtra:         PerPaxExtra,
		CreatedBy:           CreatedBy,
	}
	result := DB.Table(DBVar.TableName.CfgInitRoomRateBreakdown).Create(&CfgInitRoomRateBreakdown)
	return result.Error
}

func UpdateCfgInitRoomRateBreakdown(DB *gorm.DB, RoomRateCode string, OutletCode string, ProductCode string, SubDepartmentCode string, AccountCode string, CompanyCode string, Quantity float64, IsAmountPercent uint8, Amount float64, PerPax uint8, IncludeChild uint8, Remark string, TaxAndServiceCode string, ChargeFrequencyCode string, MaxPax int, ExtraPax float64, PerPaxExtra uint8, UpdatedBy string) error {
	var CfgInitRoomRateBreakdown = DBVar.Cfg_init_room_rate_breakdown{
		RoomRateCode:        RoomRateCode,
		OutletCode:          OutletCode,
		ProductCode:         ProductCode,
		SubDepartmentCode:   SubDepartmentCode,
		AccountCode:         AccountCode,
		CompanyCode:         CompanyCode,
		Quantity:            Quantity,
		IsAmountPercent:     IsAmountPercent,
		Amount:              Amount,
		PerPax:              PerPax,
		IncludeChild:        IncludeChild,
		Remark:              Remark,
		TaxAndServiceCode:   TaxAndServiceCode,
		ChargeFrequencyCode: ChargeFrequencyCode,
		MaxPax:              MaxPax,
		ExtraPax:            ExtraPax,
		PerPaxExtra:         PerPaxExtra,
		UpdatedBy:           UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.CfgInitRoomRateBreakdown).Omit("created_at", "created_by", "updated_at", "id").Updates(&CfgInitRoomRateBreakdown)
	return result.Error
}

func InsertGuestBreakdown(DB *gorm.DB, FolioNumber uint64, OutletCode string, ProductCode string, SubDepartmentCode string, AccountCode string, CompanyCode string, Quantity float64, IsAmountPercent uint8, Amount float64, PerPax uint8, IncludeChild uint8, Remark string, TaxAndServiceCode string, ChargeFrequencyCode string, MaxPax int, ExtraPax float64, PerPaxExtra uint8, CreatedBy string) error {
	var GuestBreakdown = DBVar.Guest_breakdown{
		FolioNumber:         FolioNumber,
		OutletCode:          OutletCode,
		ProductCode:         ProductCode,
		SubDepartmentCode:   SubDepartmentCode,
		AccountCode:         AccountCode,
		CompanyCode:         CompanyCode,
		Quantity:            Quantity,
		IsAmountPercent:     IsAmountPercent,
		Amount:              Amount,
		PerPax:              PerPax,
		IncludeChild:        IncludeChild,
		Remark:              Remark,
		TaxAndServiceCode:   TaxAndServiceCode,
		ChargeFrequencyCode: ChargeFrequencyCode,
		MaxPax:              MaxPax,
		ExtraPax:            ExtraPax,
		PerPaxExtra:         PerPaxExtra,
		CreatedBy:           CreatedBy,
	}
	result := DB.Table(DBVar.TableName.GuestBreakdown).Create(&GuestBreakdown)
	return result.Error
}

func UpdateGuestBreakdown(DB *gorm.DB, FolioNumber uint64, OutletCode string, ProductCode string, SubDepartmentCode string, AccountCode string, CompanyCode string, Quantity float64, IsAmountPercent uint8, Amount float64, PerPax uint8, IncludeChild uint8, Remark string, TaxAndServiceCode string, ChargeFrequencyCode string, MaxPax int, ExtraPax float64, PerPaxExtra uint8, UpdatedBy string) error {
	var GuestBreakdown = DBVar.Guest_breakdown{
		FolioNumber:         FolioNumber,
		OutletCode:          OutletCode,
		ProductCode:         ProductCode,
		SubDepartmentCode:   SubDepartmentCode,
		AccountCode:         AccountCode,
		CompanyCode:         CompanyCode,
		Quantity:            Quantity,
		IsAmountPercent:     IsAmountPercent,
		Amount:              Amount,
		PerPax:              PerPax,
		IncludeChild:        IncludeChild,
		Remark:              Remark,
		TaxAndServiceCode:   TaxAndServiceCode,
		ChargeFrequencyCode: ChargeFrequencyCode,
		MaxPax:              MaxPax,
		ExtraPax:            ExtraPax,
		PerPaxExtra:         PerPaxExtra,
		UpdatedBy:           UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.GuestBreakdown).Omit("created_at", "created_by", "updated_at", "id").Updates(&GuestBreakdown)
	return result.Error
}

func DeleteGuestInHouse(ctx context.Context, DB *gorm.DB, PostingDate time.Time, FolioNumber uint64) error {
	PostingDateString := General.FormatDate1(PostingDate)
	result := DB.WithContext(ctx).Table(DBVar.TableName.GuestInHouseBreakdown).Where("audit_date=?", PostingDateString).Where("folio_number=?", FolioNumber).Delete(&FolioNumber)
	if result.Error != nil {
		return result.Error
	}
	result = DB.WithContext(ctx).Table(DBVar.TableName.GuestInHouse).Where("audit_date=?", PostingDateString).Where("folio_number=?", FolioNumber).Delete(&FolioNumber)
	return result.Error
}

func InsertPosCheckTransaction(DB *gorm.DB, CheckNumber string, CaptainOrderTransactionId uint64, SubFolioId uint64, InventoryCode string, TenanCode string, SeatNumber int, SpaRoomNumber string, SpaStartDate time.Time, SpaEndDate time.Time, ProductCode string, PricePurchase float64, PriceOriginal float64, Price float64, Discount float64, EstimationCost float64, Tax float64, Service float64, CompanyCode string, CompanyCode2 string, CardCharge float64, FolioTransfer uint64, IsCompliment uint8, IsFree uint8, CreatedBy string) error {
	var PosCheckTransaction = DBVar.Pos_check_transaction{
		CheckNumber:               CheckNumber,
		CaptainOrderTransactionId: CaptainOrderTransactionId,
		SubFolioId:                SubFolioId,
		InventoryCode:             InventoryCode,
		TenanCode:                 TenanCode,
		SeatNumber:                SeatNumber,
		SpaRoomNumber:             SpaRoomNumber,
		SpaStartDate:              SpaStartDate,
		SpaEndDate:                SpaEndDate,
		ProductCode:               ProductCode,
		PricePurchase:             PricePurchase,
		PriceOriginal:             PriceOriginal,
		Price:                     Price,
		Discount:                  Discount,
		EstimationCost:            EstimationCost,
		Tax:                       Tax,
		Service:                   Service,
		CompanyCode:               CompanyCode,
		CompanyCode2:              CompanyCode2,
		CardCharge:                CardCharge,
		FolioTransfer:             FolioTransfer,
		IsCompliment:              IsCompliment,
		IsFree:                    IsFree,
		CreatedBy:                 CreatedBy,
	}
	result := DB.Table(DBVar.TableName.PosCheckTransaction).Create(&PosCheckTransaction)
	return result.Error
}

func UpdatePosCheckTransaction(DB *gorm.DB, CheckNumber string, CaptainOrderTransactionId uint64, SubFolioId uint64, InventoryCode string, TenanCode string, SeatNumber int, SpaRoomNumber string, SpaStartDate time.Time, SpaEndDate time.Time, ProductCode string, PricePurchase float64, PriceOriginal float64, Price float64, Discount float64, EstimationCost float64, Tax float64, Service float64, CompanyCode string, CompanyCode2 string, CardCharge float64, FolioTransfer uint64, IsCompliment uint8, IsFree uint8, UpdatedBy string) error {
	var PosCheckTransaction = DBVar.Pos_check_transaction{
		CheckNumber:               CheckNumber,
		CaptainOrderTransactionId: CaptainOrderTransactionId,
		SubFolioId:                SubFolioId,
		InventoryCode:             InventoryCode,
		TenanCode:                 TenanCode,
		SeatNumber:                SeatNumber,
		SpaRoomNumber:             SpaRoomNumber,
		SpaStartDate:              SpaStartDate,
		SpaEndDate:                SpaEndDate,
		ProductCode:               ProductCode,
		PricePurchase:             PricePurchase,
		PriceOriginal:             PriceOriginal,
		Price:                     Price,
		Discount:                  Discount,
		EstimationCost:            EstimationCost,
		Tax:                       Tax,
		Service:                   Service,
		CompanyCode:               CompanyCode,
		CompanyCode2:              CompanyCode2,
		CardCharge:                CardCharge,
		FolioTransfer:             FolioTransfer,
		IsCompliment:              IsCompliment,
		IsFree:                    IsFree,
		UpdatedBy:                 UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.PosCheckTransaction).Omit("created_at", "created_by", "updated_at", "id").Updates(&PosCheckTransaction)
	return result.Error
}

func UpdateVoucherStatusVoidSubFolio(ctx context.Context, DB *gorm.DB, VoucherNumber string, UserID string) error {
	result := DB.WithContext(ctx).Table(DBVar.TableName.Voucher).Where("number=?", VoucherNumber).Updates(map[string]interface{}{
		"status_code":  "A",
		"used_date":    "0000-00-00",
		"folio_number": 0,
		"sub_folio_id": 0,
		"updated_by":   UserID,
	})

	return result.Error
}

func InsertRoomUnavailable(DB *gorm.DB, RoomNumber string, StartDate time.Time, EndDate time.Time, StatusCode string, ReasonCode string, Note string, CreatedBy string) error {
	var RoomUnavailable = DBVar.Room_unavailable{
		RoomNumber: RoomNumber,
		StartDate:  StartDate,
		EndDate:    EndDate,
		StatusCode: StatusCode,
		ReasonCode: ReasonCode,
		Note:       General.PtrString(Note),
		CreatedBy:  CreatedBy,
	}
	result := DB.Table(DBVar.TableName.RoomUnavailable).Omit("updated_at", "id").Create(&RoomUnavailable)
	return result.Error
}

func UpdateRoomUnavailable(DB *gorm.DB, Id uint64, StartDate time.Time, EndDate time.Time, StatusCode string, ReasonCode string, Note string, UpdatedBy string) error {
	var RoomUnavailable = DBVar.Room_unavailable{
		Id:         Id,
		StartDate:  StartDate,
		EndDate:    EndDate,
		StatusCode: StatusCode,
		ReasonCode: ReasonCode,
		Note:       General.PtrString(Note),
		UpdatedBy:  UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.RoomUnavailable).Omit("created_at", "created_by", "updated_at", "id").Updates(&RoomUnavailable)
	return result.Error
}

func InsertMemberPoint(DB *gorm.DB, MemberCode string, AuditDate time.Time, PointTypeCode string, MemberTypeCode string, FolioNumber uint64, IsFromRate uint8, RoomTypeCode string, RateAmount float64, Point float64, CreatedBy string) error {
	var MemberPoint = DBVar.Member_point{
		MemberCode:     MemberCode,
		AuditDate:      AuditDate,
		PointTypeCode:  PointTypeCode,
		MemberTypeCode: MemberTypeCode,
		FolioNumber:    FolioNumber,
		IsFromRate:     IsFromRate,
		RoomTypeCode:   RoomTypeCode,
		RateAmount:     RateAmount,
		Point:          Point,
		CreatedBy:      CreatedBy,
	}
	result := DB.Table(DBVar.TableName.MemberPoint).Create(&MemberPoint)
	return result.Error
}

func UpdateMemberPoint(DB *gorm.DB, MemberCode string, AuditDate time.Time, PointTypeCode string, MemberTypeCode string, FolioNumber uint64, IsFromRate uint8, RoomTypeCode string, RateAmount float64, Point float64, UpdatedBy string) error {
	var MemberPoint = DBVar.Member_point{
		MemberCode:     MemberCode,
		AuditDate:      AuditDate,
		PointTypeCode:  PointTypeCode,
		MemberTypeCode: MemberTypeCode,
		FolioNumber:    FolioNumber,
		IsFromRate:     IsFromRate,
		RoomTypeCode:   RoomTypeCode,
		RateAmount:     RateAmount,
		Point:          Point,
		UpdatedBy:      UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.MemberPoint).Omit("created_at", "created_by", "updated_at", "id").Updates(&MemberPoint)
	return result.Error
}

func InsertRoomUnavailableHistory(ctx context.Context, DB *gorm.DB, AuditDate time.Time, RoomNumber string, StatusCode string, ReasonCode string, Note string, CreatedBy string) error {
	var RoomUnavailableHistory = DBVar.Room_unavailable_history{
		AuditDate:  AuditDate,
		RoomNumber: RoomNumber,
		StatusCode: StatusCode,
		ReasonCode: ReasonCode,
		Note:       Note,
		CreatedBy:  CreatedBy,
	}
	result := DB.WithContext(ctx).Table(DBVar.TableName.RoomUnavailableHistory).Create(&RoomUnavailableHistory)
	return result.Error
}

func UpdateRoomUnavailableHistory(DB *gorm.DB, AuditDate time.Time, RoomNumber string, StatusCode string, ReasonCode string, Note string, UpdatedBy string) error {
	var RoomUnavailableHistory = DBVar.Room_unavailable_history{
		AuditDate:  AuditDate,
		RoomNumber: RoomNumber,
		StatusCode: StatusCode,
		ReasonCode: ReasonCode,
		Note:       Note,
		UpdatedBy:  UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.RoomUnavailableHistory).Omit("created_at", "created_by", "updated_at", "id").Updates(&RoomUnavailableHistory)
	return result.Error
}

func InsertRoomStatistic(ctx context.Context, DB *gorm.DB, Date time.Time, TotalRoom int, OutOfOrder int, OfficeUse int, UnderConstruction int, HouseUse int, Compliment int, RoomSold int, DayUse int, RevenueGross float64, RevenueWithCompliment float64, RevenueNonPackage float64, RevenueNett float64, Adult int, Child int, AdultSold int, ChildSold int, ChildDayUse int, AdultDayUse int, AdultCompliment int, ChildCompliment int, AdultHu int, ChildHu int, PaxSingle int, WalkIn int, WalkInForeign int, CheckIn int, PersonCheckIn int, CheckInTomorrow int, CheckInPersonTomorrow int, CheckInForeign int, Reservation int, CancelReservation int, NoShowReservation int, CheckOut int, PersonCheckOut int, EarlyCheckOut int, CheckOutTomorrow int, CheckOutPersonTomorrow int, BreakfastCover int, FoodCover int, BeverageCover int, BanquetCover int, WeddingCover int, GatheringCover int, SegmentCoverBreakfast int, SegmentCoverLunch int, SegmentCoverDinner int, SegmentCoverCoffeeBreak int, RevenueBreakfast float64, RevenueFood float64, RevenueBeverage float64, RevenueBanquet float64, RevenueWedding float64, RevenueGathering float64, GuestLedger float64, GuestDeposit float64, UnitCode string) error {
	var RoomStatistic = DBVar.Room_statistic{
		Date:                    Date,
		TotalRoom:               TotalRoom,
		OutOfOrder:              OutOfOrder,
		OfficeUse:               OfficeUse,
		UnderConstruction:       UnderConstruction,
		HouseUse:                HouseUse,
		Compliment:              Compliment,
		RoomSold:                RoomSold,
		DayUse:                  DayUse,
		RevenueGross:            RevenueGross,
		RevenueWithCompliment:   RevenueWithCompliment,
		RevenueNonPackage:       RevenueNonPackage,
		RevenueNett:             RevenueNett,
		Adult:                   Adult,
		Child:                   Child,
		AdultSold:               AdultSold,
		ChildSold:               ChildSold,
		ChildDayUse:             ChildDayUse,
		AdultDayUse:             AdultDayUse,
		AdultCompliment:         AdultCompliment,
		ChildCompliment:         ChildCompliment,
		AdultHu:                 AdultHu,
		ChildHu:                 ChildHu,
		PaxSingle:               PaxSingle,
		WalkIn:                  WalkIn,
		WalkInForeign:           WalkInForeign,
		CheckIn:                 CheckIn,
		PersonCheckIn:           PersonCheckIn,
		CheckInTomorrow:         CheckInTomorrow,
		CheckInPersonTomorrow:   CheckInPersonTomorrow,
		CheckInForeign:          CheckInForeign,
		Reservation:             Reservation,
		CancelReservation:       CancelReservation,
		NoShowReservation:       NoShowReservation,
		CheckOut:                CheckOut,
		PersonCheckOut:          PersonCheckOut,
		EarlyCheckOut:           EarlyCheckOut,
		CheckOutTomorrow:        CheckOutTomorrow,
		CheckOutPersonTomorrow:  CheckOutPersonTomorrow,
		BreakfastCover:          BreakfastCover,
		FoodCover:               FoodCover,
		BeverageCover:           BeverageCover,
		BanquetCover:            BanquetCover,
		WeddingCover:            WeddingCover,
		GatheringCover:          GatheringCover,
		SegmentCoverBreakfast:   SegmentCoverBreakfast,
		SegmentCoverLunch:       SegmentCoverLunch,
		SegmentCoverDinner:      SegmentCoverDinner,
		SegmentCoverCoffeeBreak: SegmentCoverCoffeeBreak,
		RevenueBreakfast:        RevenueBreakfast,
		RevenueFood:             RevenueFood,
		RevenueBeverage:         RevenueBeverage,
		RevenueBanquet:          RevenueBanquet,
		RevenueWedding:          RevenueWedding,
		RevenueGathering:        RevenueGathering,
		GuestLedger:             GuestLedger,
		GuestDeposit:            GuestDeposit,
		UnitCode:                UnitCode,
	}
	result := DB.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "date"}, {Name: "unit_code"}},
		UpdateAll: true,
	}).Create(&RoomStatistic)

	return result.Error
}

func UpdateRoomStatisticRevenue(DB *gorm.DB, Date time.Time, RevenueGross, RevenueWithCompliment, RevenueNonPackage, RevenueNett, RevenueBreakfast, RevenueFood, RevenueBeverage, RevenueBanquet, RevenueWedding, RevenueGathering float64, UnitCode string) error {
	var RoomStatistic = DBVar.Room_statistic{
		Date:                  Date,
		RevenueGross:          RevenueGross,
		RevenueWithCompliment: RevenueWithCompliment,
		RevenueNonPackage:     RevenueNonPackage,
		RevenueNett:           RevenueNett,
		RevenueBreakfast:      RevenueBreakfast,
		RevenueFood:           RevenueFood,
		RevenueBeverage:       RevenueBeverage,
		RevenueBanquet:        RevenueBanquet,
		RevenueWedding:        RevenueWedding,
		RevenueGathering:      RevenueGathering,
	}
	result := DB.Table(DBVar.TableName.RoomStatistic).Where("date=?", General.FormatDate1(Date)).Where("unit_code=?", UnitCode).Updates(&RoomStatistic)
	return result.Error
}

func InsertFbStatistic(ctx context.Context, DB *gorm.DB, Date time.Time, OutletCode string, Adult int, Child int, AdultBeverage int, ChildBeverage int, FoodNett float64, BeverageNett float64) error {
	var FbStatistic = DBVar.Fb_statistic{
		Date:          Date,
		OutletCode:    OutletCode,
		Adult:         Adult,
		Child:         Child,
		AdultBeverage: AdultBeverage,
		ChildBeverage: ChildBeverage,
		FoodNett:      FoodNett,
		BeverageNett:  BeverageNett,
	}
	result := DB.WithContext(ctx).Table(DBVar.TableName.FbStatistic).Create(&FbStatistic)
	return result.Error
}

func UpdateFbStatistic(DB *gorm.DB, Date time.Time, OutletCode string, Adult int, Child int, AdultBeverage int, ChildBeverage int, FoodNett float64, BeverageNett float64) error {
	var FbStatistic = DBVar.Fb_statistic{
		Date:          Date,
		OutletCode:    OutletCode,
		Adult:         Adult,
		Child:         Child,
		AdultBeverage: AdultBeverage,
		ChildBeverage: ChildBeverage,
		FoodNett:      FoodNett,
		BeverageNett:  BeverageNett,
	}
	result := DB.Table(DBVar.TableName.FbStatistic).Omit("id").Updates(&FbStatistic)
	return result.Error
}

func InsertMarketStatistic(DB *gorm.DB, AuditDate time.Time, MarketCategoryCode string, MarketCode string, MarketCompanyCode string, RoomMarket int, RoomMarketCompliment int, PaxMarket int, PaxMarketCompliment int, RevenueNettMarket float64, RevenueNettMarketCompliment float64, RevenueGrossMarket float64, RevenueGrossMarketCompliment float64, RevenueNonPackageMarket float64, RevenueNonPackageMarketCompliment float64, BusinessSourceCode string, BusinessSourceCompanyCode string, RoomBusinessSource int, RoomBusinessSourceCompliment int, PaxBusinessSource int, PaxBusinessSourceCompliment int, RevenueNettBusinessSource float64, RevenueNettBusinessSourceCompliment float64, RevenueGrossBusinessSource float64, RevenueGrossBusinessSourceCompliment float64, RevenueNonPackageBusinessSource float64, RevenueNonPackageBusinessSourceCompliment float64, RoomTypeCode string, RoomTypeCompanyCode string, RoomRoomType int, RoomRoomTypeCompliment int, PaxRoomType int, PaxRoomTypeCompliment int, RevenueNettRoomType float64, RevenueNettRoomTypeCompliment float64, RevenueGrossRoomType float64, RevenueGrossRoomTypeCompliment float64, RevenueNonPackageRoomType float64, RevenueNonPackageRoomTypeCompliment float64, RoomRateCode string, RoomRateCompanyCode string, RoomRoomRate int, RoomRoomRateCompliment int, PaxRoomRate int, PaxRoomRateCompliment int, RevenueNettRoomRate float64, RevenueNettRoomRateCompliment float64, RevenueGrossRoomRate float64, RevenueGrossRoomRateCompliment float64, RevenueNonPackageRoomRate float64, RevenueNonPackageRoomRateCompliment float64, MarketingCode string, MarketingBusinessSourceCode string, MarketingCompanyCode string, RoomMarketing int, RoomMarketingCompliment int, PaxMarketing int, PaxMarketingCompliment int, RevenueNettMarketing float64, RevenueNettMarketingCompliment float64, RevenueGrossMarketing float64, RevenueGrossMarketingCompliment float64, RevenueNonPackageMarketing float64, RevenueNonPackageMarketingCompliment float64, RevenueAllNettMarketing float64, RevenueAllGrossMarketing float64, CountryCode string, CountryStateCode string, CountryCityCode string, RoomCountry int, RoomCountryCompliment int, PaxCountry int, PaxCountryCompliment int, RevenueNettCountry float64, RevenueNettCountryCompliment float64, RevenueGrossCountry float64, RevenueGrossCountryCompliment float64, RevenueNonPackageCountry float64, RevenueNonPackageCountryCompliment float64, RevenueAllNettCountry float64, RevenueAllGrossCountry float64, NationalityCode string, NationalityCountryCode string, RoomNationality int, RoomNationalityCompliment int, PaxNationality int, PaxNationalityCompliment int, RevenueNettNationality float64, RevenueNettNationalityCompliment float64, RevenueGrossNationality float64, RevenueGrossNationalityCompliment float64, RevenueNonPackageNationality float64, RevenueNonPackageNationalityCompliment float64, RevenueAllNettNationality float64, RevenueAllGrossNationality float64, BookingSourceCode string, RoomBookingSource int, RoomBookingSourceCompliment int, PaxBookingSource int, PaxBookingSourceCompliment int, RevenueNettBookingSource float64, RevenueNettBookingSourceCompliment float64, RevenueGrossBookingSource float64, RevenueGrossBookingSourceCompliment float64, RevenueNonPackageBookingSource float64, RevenueNonPackageBookingSourceCompliment float64, RevenueAllNettBookingSource float64, RevenueAllGrossBookingSource float64, PurposeOfCode string, RoomPurposeOf int, RoomPurposeOfCompliment int, PaxPurposeOf int, PaxPurposeOfCompliment int, RevenueNettPurposeOf float64, RevenueNettPurposeOfCompliment float64, RevenueGrossPurposeOf float64, RevenueGrossPurposeOfCompliment float64, RevenueNonPackagePurposeOf float64, RevenueNonPackagePurposeOfCompliment float64, RevenueAllNettPurposeOf float64, RevenueAllGrossPurposeOf float64) error {
	var MarketStatistic = DBVar.Market_statistic{
		AuditDate:                                 AuditDate,
		MarketCategoryCode:                        MarketCategoryCode,
		MarketCode:                                MarketCode,
		MarketCompanyCode:                         MarketCompanyCode,
		RoomMarket:                                RoomMarket,
		RoomMarketCompliment:                      RoomMarketCompliment,
		PaxMarket:                                 PaxMarket,
		PaxMarketCompliment:                       PaxMarketCompliment,
		RevenueNettMarket:                         RevenueNettMarket,
		RevenueNettMarketCompliment:               RevenueNettMarketCompliment,
		RevenueGrossMarket:                        RevenueGrossMarket,
		RevenueGrossMarketCompliment:              RevenueGrossMarketCompliment,
		RevenueNonPackageMarket:                   RevenueNonPackageMarket,
		RevenueNonPackageMarketCompliment:         RevenueNonPackageMarketCompliment,
		BusinessSourceCode:                        BusinessSourceCode,
		BusinessSourceCompanyCode:                 BusinessSourceCompanyCode,
		RoomBusinessSource:                        RoomBusinessSource,
		RoomBusinessSourceCompliment:              RoomBusinessSourceCompliment,
		PaxBusinessSource:                         PaxBusinessSource,
		PaxBusinessSourceCompliment:               PaxBusinessSourceCompliment,
		RevenueNettBusinessSource:                 RevenueNettBusinessSource,
		RevenueNettBusinessSourceCompliment:       RevenueNettBusinessSourceCompliment,
		RevenueGrossBusinessSource:                RevenueGrossBusinessSource,
		RevenueGrossBusinessSourceCompliment:      RevenueGrossBusinessSourceCompliment,
		RevenueNonPackageBusinessSource:           RevenueNonPackageBusinessSource,
		RevenueNonPackageBusinessSourceCompliment: RevenueNonPackageBusinessSourceCompliment,
		RoomTypeCode:                              RoomTypeCode,
		RoomTypeCompanyCode:                       RoomTypeCompanyCode,
		RoomRoomType:                              RoomRoomType,
		RoomRoomTypeCompliment:                    RoomRoomTypeCompliment,
		PaxRoomType:                               PaxRoomType,
		PaxRoomTypeCompliment:                     PaxRoomTypeCompliment,
		RevenueNettRoomType:                       RevenueNettRoomType,
		RevenueNettRoomTypeCompliment:             RevenueNettRoomTypeCompliment,
		RevenueGrossRoomType:                      RevenueGrossRoomType,
		RevenueGrossRoomTypeCompliment:            RevenueGrossRoomTypeCompliment,
		RevenueNonPackageRoomType:                 RevenueNonPackageRoomType,
		RevenueNonPackageRoomTypeCompliment:       RevenueNonPackageRoomTypeCompliment,
		RoomRateCode:                              RoomRateCode,
		RoomRateCompanyCode:                       RoomRateCompanyCode,
		RoomRoomRate:                              RoomRoomRate,
		RoomRoomRateCompliment:                    RoomRoomRateCompliment,
		PaxRoomRate:                               PaxRoomRate,
		PaxRoomRateCompliment:                     PaxRoomRateCompliment,
		RevenueNettRoomRate:                       RevenueNettRoomRate,
		RevenueNettRoomRateCompliment:             RevenueNettRoomRateCompliment,
		RevenueGrossRoomRate:                      RevenueGrossRoomRate,
		RevenueGrossRoomRateCompliment:            RevenueGrossRoomRateCompliment,
		RevenueNonPackageRoomRate:                 RevenueNonPackageRoomRate,
		RevenueNonPackageRoomRateCompliment:       RevenueNonPackageRoomRateCompliment,
		MarketingCode:                             MarketingCode,
		MarketingBusinessSourceCode:               MarketingBusinessSourceCode,
		MarketingCompanyCode:                      MarketingCompanyCode,
		RoomMarketing:                             RoomMarketing,
		RoomMarketingCompliment:                   RoomMarketingCompliment,
		PaxMarketing:                              PaxMarketing,
		PaxMarketingCompliment:                    PaxMarketingCompliment,
		RevenueNettMarketing:                      RevenueNettMarketing,
		RevenueNettMarketingCompliment:            RevenueNettMarketingCompliment,
		RevenueGrossMarketing:                     RevenueGrossMarketing,
		RevenueGrossMarketingCompliment:           RevenueGrossMarketingCompliment,
		RevenueNonPackageMarketing:                RevenueNonPackageMarketing,
		RevenueNonPackageMarketingCompliment:      RevenueNonPackageMarketingCompliment,
		RevenueAllNettMarketing:                   RevenueAllNettMarketing,
		RevenueAllGrossMarketing:                  RevenueAllGrossMarketing,
		CountryCode:                               CountryCode,
		CountryStateCode:                          CountryStateCode,
		CountryCityCode:                           CountryCityCode,
		RoomCountry:                               RoomCountry,
		RoomCountryCompliment:                     RoomCountryCompliment,
		PaxCountry:                                PaxCountry,
		PaxCountryCompliment:                      PaxCountryCompliment,
		RevenueNettCountry:                        RevenueNettCountry,
		RevenueNettCountryCompliment:              RevenueNettCountryCompliment,
		RevenueGrossCountry:                       RevenueGrossCountry,
		RevenueGrossCountryCompliment:             RevenueGrossCountryCompliment,
		RevenueNonPackageCountry:                  RevenueNonPackageCountry,
		RevenueNonPackageCountryCompliment:        RevenueNonPackageCountryCompliment,
		RevenueAllNettCountry:                     RevenueAllNettCountry,
		RevenueAllGrossCountry:                    RevenueAllGrossCountry,
		NationalityCode:                           NationalityCode,
		NationalityCountryCode:                    NationalityCountryCode,
		RoomNationality:                           RoomNationality,
		RoomNationalityCompliment:                 RoomNationalityCompliment,
		PaxNationality:                            PaxNationality,
		PaxNationalityCompliment:                  PaxNationalityCompliment,
		RevenueNettNationality:                    RevenueNettNationality,
		RevenueNettNationalityCompliment:          RevenueNettNationalityCompliment,
		RevenueGrossNationality:                   RevenueGrossNationality,
		RevenueGrossNationalityCompliment:         RevenueGrossNationalityCompliment,
		RevenueNonPackageNationality:              RevenueNonPackageNationality,
		RevenueNonPackageNationalityCompliment:    RevenueNonPackageNationalityCompliment,
		RevenueAllNettNationality:                 RevenueAllNettNationality,
		RevenueAllGrossNationality:                RevenueAllGrossNationality,
		BookingSourceCode:                         BookingSourceCode,
		RoomBookingSource:                         RoomBookingSource,
		RoomBookingSourceCompliment:               RoomBookingSourceCompliment,
		PaxBookingSource:                          PaxBookingSource,
		PaxBookingSourceCompliment:                PaxBookingSourceCompliment,
		RevenueNettBookingSource:                  RevenueNettBookingSource,
		RevenueNettBookingSourceCompliment:        RevenueNettBookingSourceCompliment,
		RevenueGrossBookingSource:                 RevenueGrossBookingSource,
		RevenueGrossBookingSourceCompliment:       RevenueGrossBookingSourceCompliment,
		RevenueNonPackageBookingSource:            RevenueNonPackageBookingSource,
		RevenueNonPackageBookingSourceCompliment:  RevenueNonPackageBookingSourceCompliment,
		RevenueAllNettBookingSource:               RevenueAllNettBookingSource,
		RevenueAllGrossBookingSource:              RevenueAllGrossBookingSource,
		PurposeOfCode:                             PurposeOfCode,
		RoomPurposeOf:                             RoomPurposeOf,
		RoomPurposeOfCompliment:                   RoomPurposeOfCompliment,
		PaxPurposeOf:                              PaxPurposeOf,
		PaxPurposeOfCompliment:                    PaxPurposeOfCompliment,
		RevenueNettPurposeOf:                      RevenueNettPurposeOf,
		RevenueNettPurposeOfCompliment:            RevenueNettPurposeOfCompliment,
		RevenueGrossPurposeOf:                     RevenueGrossPurposeOf,
		RevenueGrossPurposeOfCompliment:           RevenueGrossPurposeOfCompliment,
		RevenueNonPackagePurposeOf:                RevenueNonPackagePurposeOf,
		RevenueNonPackagePurposeOfCompliment:      RevenueNonPackagePurposeOfCompliment,
		RevenueAllNettPurposeOf:                   RevenueAllNettPurposeOf,
		RevenueAllGrossPurposeOf:                  RevenueAllGrossPurposeOf,
	}
	result := DB.Table(DBVar.TableName.MarketStatistic).Create(&MarketStatistic)
	return result.Error
}

func UpdateMarketStatistic(DB *gorm.DB, AuditDate time.Time, MarketCategoryCode string, MarketCode string, MarketCompanyCode string, RoomMarket int, RoomMarketCompliment int, PaxMarket int, PaxMarketCompliment int, RevenueNettMarket float64, RevenueNettMarketCompliment float64, RevenueGrossMarket float64, RevenueGrossMarketCompliment float64, RevenueNonPackageMarket float64, RevenueNonPackageMarketCompliment float64, BusinessSourceCode string, BusinessSourceCompanyCode string, RoomBusinessSource int, RoomBusinessSourceCompliment int, PaxBusinessSource int, PaxBusinessSourceCompliment int, RevenueNettBusinessSource float64, RevenueNettBusinessSourceCompliment float64, RevenueGrossBusinessSource float64, RevenueGrossBusinessSourceCompliment float64, RevenueNonPackageBusinessSource float64, RevenueNonPackageBusinessSourceCompliment float64, RoomTypeCode string, RoomTypeCompanyCode string, RoomRoomType int, RoomRoomTypeCompliment int, PaxRoomType int, PaxRoomTypeCompliment int, RevenueNettRoomType float64, RevenueNettRoomTypeCompliment float64, RevenueGrossRoomType float64, RevenueGrossRoomTypeCompliment float64, RevenueNonPackageRoomType float64, RevenueNonPackageRoomTypeCompliment float64, RoomRateCode string, RoomRateCompanyCode string, RoomRoomRate int, RoomRoomRateCompliment int, PaxRoomRate int, PaxRoomRateCompliment int, RevenueNettRoomRate float64, RevenueNettRoomRateCompliment float64, RevenueGrossRoomRate float64, RevenueGrossRoomRateCompliment float64, RevenueNonPackageRoomRate float64, RevenueNonPackageRoomRateCompliment float64, MarketingCode string, MarketingBusinessSourceCode string, MarketingCompanyCode string, RoomMarketing int, RoomMarketingCompliment int, PaxMarketing int, PaxMarketingCompliment int, RevenueNettMarketing float64, RevenueNettMarketingCompliment float64, RevenueGrossMarketing float64, RevenueGrossMarketingCompliment float64, RevenueNonPackageMarketing float64, RevenueNonPackageMarketingCompliment float64, RevenueAllNettMarketing float64, RevenueAllGrossMarketing float64, CountryCode string, CountryStateCode string, CountryCityCode string, RoomCountry int, RoomCountryCompliment int, PaxCountry int, PaxCountryCompliment int, RevenueNettCountry float64, RevenueNettCountryCompliment float64, RevenueGrossCountry float64, RevenueGrossCountryCompliment float64, RevenueNonPackageCountry float64, RevenueNonPackageCountryCompliment float64, RevenueAllNettCountry float64, RevenueAllGrossCountry float64, NationalityCode string, NationalityCountryCode string, RoomNationality int, RoomNationalityCompliment int, PaxNationality int, PaxNationalityCompliment int, RevenueNettNationality float64, RevenueNettNationalityCompliment float64, RevenueGrossNationality float64, RevenueGrossNationalityCompliment float64, RevenueNonPackageNationality float64, RevenueNonPackageNationalityCompliment float64, RevenueAllNettNationality float64, RevenueAllGrossNationality float64, BookingSourceCode string, RoomBookingSource int, RoomBookingSourceCompliment int, PaxBookingSource int, PaxBookingSourceCompliment int, RevenueNettBookingSource float64, RevenueNettBookingSourceCompliment float64, RevenueGrossBookingSource float64, RevenueGrossBookingSourceCompliment float64, RevenueNonPackageBookingSource float64, RevenueNonPackageBookingSourceCompliment float64, RevenueAllNettBookingSource float64, RevenueAllGrossBookingSource float64, PurposeOfCode string, RoomPurposeOf int, RoomPurposeOfCompliment int, PaxPurposeOf int, PaxPurposeOfCompliment int, RevenueNettPurposeOf float64, RevenueNettPurposeOfCompliment float64, RevenueGrossPurposeOf float64, RevenueGrossPurposeOfCompliment float64, RevenueNonPackagePurposeOf float64, RevenueNonPackagePurposeOfCompliment float64, RevenueAllNettPurposeOf float64, RevenueAllGrossPurposeOf float64) error {
	var MarketStatistic = DBVar.Market_statistic{
		AuditDate:                                 AuditDate,
		MarketCategoryCode:                        MarketCategoryCode,
		MarketCode:                                MarketCode,
		MarketCompanyCode:                         MarketCompanyCode,
		RoomMarket:                                RoomMarket,
		RoomMarketCompliment:                      RoomMarketCompliment,
		PaxMarket:                                 PaxMarket,
		PaxMarketCompliment:                       PaxMarketCompliment,
		RevenueNettMarket:                         RevenueNettMarket,
		RevenueNettMarketCompliment:               RevenueNettMarketCompliment,
		RevenueGrossMarket:                        RevenueGrossMarket,
		RevenueGrossMarketCompliment:              RevenueGrossMarketCompliment,
		RevenueNonPackageMarket:                   RevenueNonPackageMarket,
		RevenueNonPackageMarketCompliment:         RevenueNonPackageMarketCompliment,
		BusinessSourceCode:                        BusinessSourceCode,
		BusinessSourceCompanyCode:                 BusinessSourceCompanyCode,
		RoomBusinessSource:                        RoomBusinessSource,
		RoomBusinessSourceCompliment:              RoomBusinessSourceCompliment,
		PaxBusinessSource:                         PaxBusinessSource,
		PaxBusinessSourceCompliment:               PaxBusinessSourceCompliment,
		RevenueNettBusinessSource:                 RevenueNettBusinessSource,
		RevenueNettBusinessSourceCompliment:       RevenueNettBusinessSourceCompliment,
		RevenueGrossBusinessSource:                RevenueGrossBusinessSource,
		RevenueGrossBusinessSourceCompliment:      RevenueGrossBusinessSourceCompliment,
		RevenueNonPackageBusinessSource:           RevenueNonPackageBusinessSource,
		RevenueNonPackageBusinessSourceCompliment: RevenueNonPackageBusinessSourceCompliment,
		RoomTypeCode:                              RoomTypeCode,
		RoomTypeCompanyCode:                       RoomTypeCompanyCode,
		RoomRoomType:                              RoomRoomType,
		RoomRoomTypeCompliment:                    RoomRoomTypeCompliment,
		PaxRoomType:                               PaxRoomType,
		PaxRoomTypeCompliment:                     PaxRoomTypeCompliment,
		RevenueNettRoomType:                       RevenueNettRoomType,
		RevenueNettRoomTypeCompliment:             RevenueNettRoomTypeCompliment,
		RevenueGrossRoomType:                      RevenueGrossRoomType,
		RevenueGrossRoomTypeCompliment:            RevenueGrossRoomTypeCompliment,
		RevenueNonPackageRoomType:                 RevenueNonPackageRoomType,
		RevenueNonPackageRoomTypeCompliment:       RevenueNonPackageRoomTypeCompliment,
		RoomRateCode:                              RoomRateCode,
		RoomRateCompanyCode:                       RoomRateCompanyCode,
		RoomRoomRate:                              RoomRoomRate,
		RoomRoomRateCompliment:                    RoomRoomRateCompliment,
		PaxRoomRate:                               PaxRoomRate,
		PaxRoomRateCompliment:                     PaxRoomRateCompliment,
		RevenueNettRoomRate:                       RevenueNettRoomRate,
		RevenueNettRoomRateCompliment:             RevenueNettRoomRateCompliment,
		RevenueGrossRoomRate:                      RevenueGrossRoomRate,
		RevenueGrossRoomRateCompliment:            RevenueGrossRoomRateCompliment,
		RevenueNonPackageRoomRate:                 RevenueNonPackageRoomRate,
		RevenueNonPackageRoomRateCompliment:       RevenueNonPackageRoomRateCompliment,
		MarketingCode:                             MarketingCode,
		MarketingBusinessSourceCode:               MarketingBusinessSourceCode,
		MarketingCompanyCode:                      MarketingCompanyCode,
		RoomMarketing:                             RoomMarketing,
		RoomMarketingCompliment:                   RoomMarketingCompliment,
		PaxMarketing:                              PaxMarketing,
		PaxMarketingCompliment:                    PaxMarketingCompliment,
		RevenueNettMarketing:                      RevenueNettMarketing,
		RevenueNettMarketingCompliment:            RevenueNettMarketingCompliment,
		RevenueGrossMarketing:                     RevenueGrossMarketing,
		RevenueGrossMarketingCompliment:           RevenueGrossMarketingCompliment,
		RevenueNonPackageMarketing:                RevenueNonPackageMarketing,
		RevenueNonPackageMarketingCompliment:      RevenueNonPackageMarketingCompliment,
		RevenueAllNettMarketing:                   RevenueAllNettMarketing,
		RevenueAllGrossMarketing:                  RevenueAllGrossMarketing,
		CountryCode:                               CountryCode,
		CountryStateCode:                          CountryStateCode,
		CountryCityCode:                           CountryCityCode,
		RoomCountry:                               RoomCountry,
		RoomCountryCompliment:                     RoomCountryCompliment,
		PaxCountry:                                PaxCountry,
		PaxCountryCompliment:                      PaxCountryCompliment,
		RevenueNettCountry:                        RevenueNettCountry,
		RevenueNettCountryCompliment:              RevenueNettCountryCompliment,
		RevenueGrossCountry:                       RevenueGrossCountry,
		RevenueGrossCountryCompliment:             RevenueGrossCountryCompliment,
		RevenueNonPackageCountry:                  RevenueNonPackageCountry,
		RevenueNonPackageCountryCompliment:        RevenueNonPackageCountryCompliment,
		RevenueAllNettCountry:                     RevenueAllNettCountry,
		RevenueAllGrossCountry:                    RevenueAllGrossCountry,
		NationalityCode:                           NationalityCode,
		NationalityCountryCode:                    NationalityCountryCode,
		RoomNationality:                           RoomNationality,
		RoomNationalityCompliment:                 RoomNationalityCompliment,
		PaxNationality:                            PaxNationality,
		PaxNationalityCompliment:                  PaxNationalityCompliment,
		RevenueNettNationality:                    RevenueNettNationality,
		RevenueNettNationalityCompliment:          RevenueNettNationalityCompliment,
		RevenueGrossNationality:                   RevenueGrossNationality,
		RevenueGrossNationalityCompliment:         RevenueGrossNationalityCompliment,
		RevenueNonPackageNationality:              RevenueNonPackageNationality,
		RevenueNonPackageNationalityCompliment:    RevenueNonPackageNationalityCompliment,
		RevenueAllNettNationality:                 RevenueAllNettNationality,
		RevenueAllGrossNationality:                RevenueAllGrossNationality,
		BookingSourceCode:                         BookingSourceCode,
		RoomBookingSource:                         RoomBookingSource,
		RoomBookingSourceCompliment:               RoomBookingSourceCompliment,
		PaxBookingSource:                          PaxBookingSource,
		PaxBookingSourceCompliment:                PaxBookingSourceCompliment,
		RevenueNettBookingSource:                  RevenueNettBookingSource,
		RevenueNettBookingSourceCompliment:        RevenueNettBookingSourceCompliment,
		RevenueGrossBookingSource:                 RevenueGrossBookingSource,
		RevenueGrossBookingSourceCompliment:       RevenueGrossBookingSourceCompliment,
		RevenueNonPackageBookingSource:            RevenueNonPackageBookingSource,
		RevenueNonPackageBookingSourceCompliment:  RevenueNonPackageBookingSourceCompliment,
		RevenueAllNettBookingSource:               RevenueAllNettBookingSource,
		RevenueAllGrossBookingSource:              RevenueAllGrossBookingSource,
		PurposeOfCode:                             PurposeOfCode,
		RoomPurposeOf:                             RoomPurposeOf,
		RoomPurposeOfCompliment:                   RoomPurposeOfCompliment,
		PaxPurposeOf:                              PaxPurposeOf,
		PaxPurposeOfCompliment:                    PaxPurposeOfCompliment,
		RevenueNettPurposeOf:                      RevenueNettPurposeOf,
		RevenueNettPurposeOfCompliment:            RevenueNettPurposeOfCompliment,
		RevenueGrossPurposeOf:                     RevenueGrossPurposeOf,
		RevenueGrossPurposeOfCompliment:           RevenueGrossPurposeOfCompliment,
		RevenueNonPackagePurposeOf:                RevenueNonPackagePurposeOf,
		RevenueNonPackagePurposeOfCompliment:      RevenueNonPackagePurposeOfCompliment,
		RevenueAllNettPurposeOf:                   RevenueAllNettPurposeOf,
		RevenueAllGrossPurposeOf:                  RevenueAllGrossPurposeOf,
	}
	result := DB.Table(DBVar.TableName.MarketStatistic).Omit("id").Updates(&MarketStatistic)
	return result.Error
}

func InsertRoomStatus(DB *gorm.DB, AuditDate time.Time, RoomNumber string, Status string) error {
	var RoomStatus = DBVar.Room_status{
		AuditDate:  AuditDate,
		RoomNumber: RoomNumber,
		Status:     Status,
	}
	result := DB.Table(DBVar.TableName.RoomStatus).Create(&RoomStatus)
	return result.Error
}

func InsertPosCaptainOrder(DB *gorm.DB, ReservationNumber uint64, OutletCode string, TableNumber string, WaitressCode string, CustomerCode string, MemberCode string, TitleCode string, FullName string, Adult int, Child int, DocumentNumber string, Remark string, AuditDate time.Time, MarketCode string, CompanyCode string, MarketingCode string, TimeSegmentCode string, TypeCode string, ComplimentTypeCode string, SubDepartmentCode string, CreatedBy string) (DBVar.Pos_captain_order, error) {
	var PosCaptainOrder = DBVar.Pos_captain_order{
		ReservationNumber:  ReservationNumber,
		OutletCode:         OutletCode,
		TableNumber:        TableNumber,
		WaitressCode:       WaitressCode,
		CustomerCode:       CustomerCode,
		MemberCode:         MemberCode,
		TitleCode:          TitleCode,
		FullName:           FullName,
		Adult:              Adult,
		Child:              Child,
		DocumentNumber:     DocumentNumber,
		Remark:             Remark,
		AuditDate:          AuditDate,
		MarketCode:         MarketCode,
		CompanyCode:        CompanyCode,
		SalesCode:          MarketingCode,
		TimeSegmentCode:    TimeSegmentCode,
		TypeCode:           TypeCode,
		ComplimentTypeCode: ComplimentTypeCode,
		SubDepartmentCode:  SubDepartmentCode,
		IsOpen:             1,
		CreatedBy:          CreatedBy,
	}
	result := DB.Table(DBVar.TableName.PosCaptainOrder).Create(&PosCaptainOrder)
	return PosCaptainOrder, result.Error
}

func UpdatePosCaptainOrder(DB *gorm.DB, Id uint64, TableNumber string, WaitressCode string, TitleCode string, FullName string, Adult int, Child int, DocumentNumber string, Remark string, MarketCode string, CompanyCode string, MarketingCode string, TimeSegmentCode string, TypeCode string, ComplimentTypeCode string, SubDepartmentCode string, UpdatedBy string) (DBVar.Pos_captain_order, error) {
	var PosCaptainOrder = DBVar.Pos_captain_order{
		TableNumber:        TableNumber,
		WaitressCode:       WaitressCode,
		TitleCode:          TitleCode,
		FullName:           FullName,
		Adult:              Adult,
		Child:              Child,
		DocumentNumber:     DocumentNumber,
		Remark:             Remark,
		MarketCode:         MarketCode,
		TypeCode:           TypeCode,
		CompanyCode:        CompanyCode,
		SalesCode:          MarketingCode,
		TimeSegmentCode:    TimeSegmentCode,
		ComplimentTypeCode: ComplimentTypeCode,
		SubDepartmentCode:  SubDepartmentCode,
		Id:                 Id,
		UpdatedBy:          UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.PosCaptainOrder).Omit("created_at", "created_by", "updated_at", "id").Updates(&PosCaptainOrder)
	return PosCaptainOrder, result.Error
}

func InsertPosReservationTable(DB *gorm.DB, ReservationNumber uint64, Start time.Time, Finish time.Time, TableNumber string, ParentId uint64, EventType int, Options int, Caption string, RecurrenceIndex int, RecurrenceInfo string, Message string, ReminderDate time.Time, ReminderMinutes int, State int, LabelColor int, SincId string, ReminderResource string, BlockType string, CaptainOrderId uint64, CreatedBy string) error {
	var PosReservationTable = DBVar.Pos_reservation_table{
		ReservationNumber: ReservationNumber,
		Start:             Start,
		Finish:            Finish,
		TableNumber:       TableNumber,
		ParentId:          ParentId,
		EventType:         EventType,
		Options:           Options,
		Caption:           Caption,
		RecurrenceIndex:   RecurrenceIndex,
		RecurrenceInfo:    RecurrenceInfo,
		Message:           Message,
		ReminderDate:      ReminderDate,
		ReminderMinutes:   ReminderMinutes,
		State:             State,
		LabelColor:        LabelColor,
		SincId:            SincId,
		ReminderResource:  ReminderResource,
		BlockType:         BlockType,
		CaptainOrderId:    CaptainOrderId,
		CreatedBy:         CreatedBy,
	}
	result := DB.Table(DBVar.TableName.PosReservationTable).Create(&PosReservationTable)
	return result.Error
}

func UpdatePosReservationTable(DB *gorm.DB, TableNumber string, ParentId uint64, EventType int, Options int, Caption string, RecurrenceIndex int, RecurrenceInfo string, Message string, State int, LabelColor int, SincId string, ReminderResource string, BlockType string, CaptainOrderId uint64, UpdatedBy string) error {
	var PosReservationTable = DBVar.Pos_reservation_table{
		TableNumber:      TableNumber,
		ParentId:         ParentId,
		EventType:        EventType,
		Options:          Options,
		Caption:          Caption,
		RecurrenceIndex:  RecurrenceIndex,
		RecurrenceInfo:   RecurrenceInfo,
		Message:          Message,
		State:            State,
		LabelColor:       LabelColor,
		SincId:           SincId,
		ReminderResource: ReminderResource,
		BlockType:        BlockType,
		CaptainOrderId:   CaptainOrderId,
		UpdatedBy:        UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.PosReservationTable).Where("captain_order_id=?", CaptainOrderId).Omit("created_at", "created_by", "updated_at", "id").Updates(&PosReservationTable)
	return result.Error
}

func InsertPosCaptainOrderTransaction(ctx context.Context, c *gin.Context, DB *gorm.DB, CaptainOrderId uint64, InventoryCode string, TenanCode string, SeatNumber int, SpaRoomNumber string, SpaStartDate time.Time, SpaEndDate time.Time, ProductCode string, AccountCode string, Description string, Quantity float64, PricePurchase float64, PriceOriginal float64, Price float64, Discount float64, DiscountTemp float64, Tax float64, Service float64, DefaultCurrencyCode string, CurrencyCode string, ExchangeRate float64, Remark string, TypeCode string, AuditDate time.Time, CompanyCode string, CompanyCode2 string, CardBankCode string, CardTypeCode string, CardCharge float64, CardNumber string, CardHolder string, ValidMonth string, ValidYear string, FolioTransfer uint64, SubFolioTransfer string, IsCompliment uint8, IsFree uint8, Shift string, LogShiftId uint64, CreatedBy string) error {
	ctx, span := global_var.Tracer.Start(ctx, "InsertPosCaptainOrderTransaction")
	defer span.End()

	if AuditDate == (time.Time{}) {
		AuditDate = GetAuditDate(c, DB, false)
	}

	var PosCaptainOrderTransaction = DBVar.Pos_captain_order_transaction{
		CaptainOrderId:      CaptainOrderId,
		InventoryCode:       InventoryCode,
		TenanCode:           TenanCode,
		SeatNumber:          SeatNumber,
		SpaRoomNumber:       SpaRoomNumber,
		SpaStartDate:        SpaStartDate,
		SpaEndDate:          SpaEndDate,
		ProductCode:         ProductCode,
		AccountCode:         AccountCode,
		Description:         Description,
		Quantity:            Quantity,
		PricePurchase:       PricePurchase,
		PriceOriginal:       PriceOriginal,
		Price:               Price,
		Discount:            Discount,
		DiscountTemp:        DiscountTemp,
		Tax:                 Tax,
		Service:             Service,
		DefaultCurrencyCode: DefaultCurrencyCode,
		CurrencyCode:        CurrencyCode,
		ExchangeRate:        ExchangeRate,
		Remark:              Remark,
		TypeCode:            TypeCode,
		AuditDate:           AuditDate,
		CompanyCode:         CompanyCode,
		CompanyCode2:        CompanyCode2,
		CardBankCode:        CardBankCode,
		CardTypeCode:        CardTypeCode,
		CardCharge:          CardCharge,
		CardNumber:          CardNumber,
		CardHolder:          CardHolder,
		ValidMonth:          ValidMonth,
		ValidYear:           ValidYear,
		FolioTransfer:       FolioTransfer,
		SubFolioTransfer:    SubFolioTransfer,
		IsCompliment:        IsCompliment,
		IsFree:              IsFree,
		Shift:               Shift,
		LogShiftId:          LogShiftId,
		CreatedBy:           CreatedBy,
	}
	result := DB.WithContext(ctx).Table(DBVar.TableName.PosCaptainOrderTransaction).Create(&PosCaptainOrderTransaction)
	return result.Error
}

func UpdatePosCaptainOrderTransaction(DB *gorm.DB, CaptainOrderId uint64, InventoryCode string, TenanCode string, SeatNumber int, SpaRoomNumber string, SpaStartDate time.Time, SpaEndDate time.Time, ProductCode string, AccountCode string, Description string, Quantity float64, QuantityPrinted float64, QuantityPrintedCheck float64, PricePurchase float64, PriceOriginal float64, Price float64, Discount float64, DiscountTemp float64, Tax float64, Service float64, DefaultCurrencyCode string, CurrencyCode string, ExchangeRate float64, Remark string, TypeCode string, AuditDate time.Time, PostingDate time.Time, CompanyCode string, CompanyCode2 string, CardBankCode string, CardTypeCode string, CardCharge float64, CardNumber string, CardHolder string, ValidMonth string, ValidYear string, FolioTransfer uint64, SubFolioTransfer string, IsCompliment uint8, IsFree uint8, IsRemove uint8, RemoveDate time.Time, RemoveBy string, Shift string, LogShiftId uint64, UpdatedBy string) error {
	var PosCaptainOrderTransaction = DBVar.Pos_captain_order_transaction{
		CaptainOrderId:       CaptainOrderId,
		InventoryCode:        InventoryCode,
		TenanCode:            TenanCode,
		SeatNumber:           SeatNumber,
		SpaRoomNumber:        SpaRoomNumber,
		SpaStartDate:         SpaStartDate,
		SpaEndDate:           SpaEndDate,
		ProductCode:          ProductCode,
		AccountCode:          AccountCode,
		Description:          Description,
		Quantity:             Quantity,
		QuantityPrinted:      QuantityPrinted,
		QuantityPrintedCheck: QuantityPrintedCheck,
		PricePurchase:        PricePurchase,
		PriceOriginal:        PriceOriginal,
		Price:                Price,
		Discount:             Discount,
		DiscountTemp:         DiscountTemp,
		Tax:                  Tax,
		Service:              Service,
		DefaultCurrencyCode:  DefaultCurrencyCode,
		CurrencyCode:         CurrencyCode,
		ExchangeRate:         ExchangeRate,
		Remark:               Remark,
		TypeCode:             TypeCode,
		AuditDate:            AuditDate,
		PostingDate:          PostingDate,
		CompanyCode:          CompanyCode,
		CompanyCode2:         CompanyCode2,
		CardBankCode:         CardBankCode,
		CardTypeCode:         CardTypeCode,
		CardCharge:           CardCharge,
		CardNumber:           CardNumber,
		CardHolder:           CardHolder,
		ValidMonth:           ValidMonth,
		ValidYear:            ValidYear,
		FolioTransfer:        FolioTransfer,
		SubFolioTransfer:     SubFolioTransfer,
		IsCompliment:         IsCompliment,
		IsFree:               IsFree,
		IsRemove:             IsRemove,
		RemoveDate:           RemoveDate,
		RemoveBy:             RemoveBy,
		Shift:                Shift,
		LogShiftId:           LogShiftId,
		UpdatedBy:            UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.PosCaptainOrderTransaction).Omit("created_at", "created_by", "updated_at", "id").Updates(&PosCaptainOrderTransaction)
	return result.Error
}

func InsertPosCheck(DB *gorm.DB, Number string, TypeCode string, CaptainOrderId uint64, FolioNumber uint64, ContactPersonId uint64, OutletCode string, TableNumber string, WaitressCode string, MemberCode string, ComplimentTypeCode string, SubDepartmentCode string, Remark string, AuditDate time.Time, MarketCode string, TimeSegmentCode string, Void uint8, VoidDate time.Time, VoidBy string, VoidReason string, CreatedBy string) error {
	var PosCheck = DBVar.Pos_check{
		Number:             Number,
		TypeCode:           TypeCode,
		CaptainOrderId:     CaptainOrderId,
		FolioNumber:        FolioNumber,
		ContactPersonId:    ContactPersonId,
		OutletCode:         OutletCode,
		TableNumber:        TableNumber,
		WaitressCode:       WaitressCode,
		MemberCode:         MemberCode,
		ComplimentTypeCode: ComplimentTypeCode,
		SubDepartmentCode:  SubDepartmentCode,
		Remark:             Remark,
		AuditDate:          AuditDate,
		MarketCode:         MarketCode,
		TimeSegmentCode:    TimeSegmentCode,
		Void:               Void,
		VoidDate:           VoidDate,
		VoidBy:             VoidBy,
		VoidReason:         VoidReason,
		CreatedBy:          CreatedBy,
	}
	result := DB.Table(DBVar.TableName.PosCheck).Create(&PosCheck)
	return result.Error
}

func UpdatePosCheck(DB *gorm.DB, Id uint64, Number string, TypeCode string, CaptainOrderId uint64, FolioNumber uint64, ContactPersonId uint64, OutletCode string, TableNumber string, WaitressCode string, MemberCode string, ComplimentTypeCode string, SubDepartmentCode string, Remark string, MarketCode string, TimeSegmentCode string, Void uint8, VoidDate time.Time, VoidBy string, VoidReason string, UpdatedBy string) error {
	var PosCheck = DBVar.Pos_check{
		Number:             Number,
		TypeCode:           TypeCode,
		CaptainOrderId:     CaptainOrderId,
		FolioNumber:        FolioNumber,
		ContactPersonId:    ContactPersonId,
		OutletCode:         OutletCode,
		TableNumber:        TableNumber,
		WaitressCode:       WaitressCode,
		MemberCode:         MemberCode,
		ComplimentTypeCode: ComplimentTypeCode,
		SubDepartmentCode:  SubDepartmentCode,
		Remark:             Remark,
		MarketCode:         MarketCode,
		TimeSegmentCode:    TimeSegmentCode,
		Void:               Void,
		VoidDate:           VoidDate,
		VoidBy:             VoidBy,
		VoidReason:         VoidReason,
		UpdatedBy:          UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.PosCheck).Where("id=?", Id).Omit("created_at", "created_by", "updated_at", "id").Updates(&PosCheck)
	return result.Error
}
func InsertLogShift(DB *gorm.DB, CreatedBy string, Shift string, StartDate time.Time, AuditDate time.Time, OpeningBalance float64, Remark string, IpAddress string, ComputerName string, MacAddress string) (uint64, error) {
	var LogShift = DBVar.Log_shift{
		CreatedBy:      CreatedBy,
		Shift:          Shift,
		StartDate:      StartDate,
		AuditDate:      AuditDate,
		OpeningBalance: OpeningBalance,
		Remark:         Remark,
		IpAddress:      IpAddress,
		ComputerName:   ComputerName,
		MacAddress:     MacAddress,
		IsOpen:         1,
	}
	result := DB.Table(DBVar.TableName.LogShift).Create(&LogShift)
	return LogShift.Id, result.Error
}

func UpdateLogShift(DB *gorm.DB, UpdatedBy string, Shift string, StartDate time.Time, EndDate time.Time, AuditDate time.Time, OpeningBalance float64, Remark string, IpAddress string, ComputerName string, MacAddress string, IsOpen uint8) error {
	var LogShift = DBVar.Log_shift{
		UpdatedBy:      UpdatedBy,
		Shift:          Shift,
		StartDate:      StartDate,
		EndDate:        EndDate,
		AuditDate:      AuditDate,
		OpeningBalance: OpeningBalance,
		Remark:         Remark,
		IpAddress:      IpAddress,
		ComputerName:   ComputerName,
		MacAddress:     MacAddress,
		IsOpen:         IsOpen,
	}
	result := DB.Table(DBVar.TableName.LogShift).Omit("id", "created_at", "created_by", "updated_at").Updates(&LogShift)
	return result.Error
}
func InsertUserGroup(DB *gorm.DB, Id chan uint64, AccessForm string, AccessSpecial string, AccessKeylock string, AccessReservation string, AccessDeposit string, AccessInHouse string, AccessWalkIn string, AccessFolio string, AccessFolioHistory string, AccessFloorPlan string, AccessMemberVoucherGift string, SaMaxDiscountPercent int, SaMaxDiscountAmount float64, CreatedBy string) (uint64, error) {
	var UserGroup = DBVar.User_group{
		AccessForm:              AccessForm,
		AccessSpecial:           AccessSpecial,
		AccessKeylock:           AccessKeylock,
		AccessReservation:       AccessReservation,
		AccessDeposit:           AccessDeposit,
		AccessInHouse:           AccessInHouse,
		AccessWalkIn:            AccessWalkIn,
		AccessFolio:             AccessFolio,
		AccessFolioHistory:      AccessFolioHistory,
		AccessFloorPlan:         AccessFloorPlan,
		AccessMemberVoucherGift: AccessMemberVoucherGift,
		SaMaxDiscountPercent:    SaMaxDiscountPercent,
		SaMaxDiscountAmount:     SaMaxDiscountAmount,
		IsActive:                1,
		CreatedBy:               CreatedBy,
	}
	result := DB.Table(DBVar.TableName.UserGroup).Create(&UserGroup)
	Id <- UserGroup.Id
	return UserGroup.Id, result.Error
}

func UpdateUserGroup(DB *gorm.DB, Id uint64, AccessForm string, AccessSpecial string, AccessKeylock string, AccessReservation string, AccessDeposit string, AccessInHouse string, AccessWalkIn string, AccessFolio string, AccessFolioHistory string, AccessFloorPlan string, AccessMemberVoucherGift string, SaMaxDiscountPercent int, SaMaxDiscountAmount float64, UpdatedBy string) error {
	var UserGroup = DBVar.User_group{
		Id:                      Id,
		AccessForm:              AccessForm,
		AccessSpecial:           AccessSpecial,
		AccessKeylock:           AccessKeylock,
		AccessReservation:       AccessReservation,
		AccessDeposit:           AccessDeposit,
		AccessInHouse:           AccessInHouse,
		AccessWalkIn:            AccessWalkIn,
		AccessFolio:             AccessFolio,
		AccessFolioHistory:      AccessFolioHistory,
		AccessFloorPlan:         AccessFloorPlan,
		AccessMemberVoucherGift: AccessMemberVoucherGift,
		SaMaxDiscountPercent:    SaMaxDiscountPercent,
		SaMaxDiscountAmount:     SaMaxDiscountAmount,
		UpdatedBy:               UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.UserGroup).Omit("created_at", "created_by", "updated_at", "id").Where("id=?", Id).Updates(&UserGroup)
	return result.Error
}

func InsertBanUserGroup(DB *gorm.DB, Id chan uint64, AccessForm string, AccessSpecial string, AccessReservation string, AccessDeposit string, AccessInHouse string, AccessFolio string, AccessFolioHistory string, CreatedBy string) (uint64, error) {
	var BanUserGroup = DBVar.Ban_user_group{
		AccessForm:         AccessForm,
		AccessSpecial:      AccessSpecial,
		AccessReservation:  AccessReservation,
		AccessDeposit:      AccessDeposit,
		AccessInHouse:      AccessInHouse,
		AccessFolio:        AccessFolio,
		AccessFolioHistory: AccessFolioHistory,
		CreatedBy:          CreatedBy,
	}
	result := DB.Table(DBVar.TableName.BanUserGroup).Create(&BanUserGroup)
	Id <- BanUserGroup.Id
	return BanUserGroup.Id, result.Error
}

func UpdateBanUserGroup(DB *gorm.DB, Id uint64, AccessForm string, AccessSpecial string, AccessReservation string, AccessDeposit string, AccessInHouse string, AccessFolio string, AccessFolioHistory string, UpdatedBy string) error {
	var BanUserGroup = DBVar.Ban_user_group{
		Id:                 Id,
		AccessForm:         AccessForm,
		AccessSpecial:      AccessSpecial,
		AccessReservation:  AccessReservation,
		AccessDeposit:      AccessDeposit,
		AccessInHouse:      AccessInHouse,
		AccessFolio:        AccessFolio,
		AccessFolioHistory: AccessFolioHistory,
		UpdatedBy:          UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.BanUserGroup).Omit("created_at", "created_by", "updated_at", "id").Updates(&BanUserGroup)
	return result.Error
}

func InsertUserGroupAccess(DB *gorm.DB, Code string, GeneralUserGroupId uint64, UserGroupId uint64, PosUserGroupId uint64, BanUserGroupId uint64, AccUserGroupId uint64, AstUserGroupId uint64, PyrUserGroupId uint64, CorUserGroupId uint64, ReportUserGroupId uint64, ToolsUserGroupId uint64, UserAccessLevelCode int, IsActive uint8, CreatedBy string) error {
	var UserGroupAccess = DBVar.User_group_access{
		Code:                Code,
		GeneralUserGroupId:  GeneralUserGroupId,
		UserGroupId:         UserGroupId,
		PosUserGroupId:      PosUserGroupId,
		BanUserGroupId:      BanUserGroupId,
		AccUserGroupId:      AccUserGroupId,
		AstUserGroupId:      AstUserGroupId,
		PyrUserGroupId:      PyrUserGroupId,
		CorUserGroupId:      CorUserGroupId,
		ReportUserGroupId:   ReportUserGroupId,
		ToolsUserGroupId:    ToolsUserGroupId,
		UserAccessLevelCode: UserAccessLevelCode,
		IsActive:            IsActive,
		CreatedBy:           CreatedBy,
	}
	result := DB.Table(DBVar.TableName.UserGroupAccess).Create(&UserGroupAccess)
	return result.Error
}

func UpdateUserGroupAccess(DB *gorm.DB, Id string, GeneralUserGroupId uint64, UserGroupId uint64, PosUserGroupId uint64, BanUserGroupId uint64, AccUserGroupId uint64, AstUserGroupId uint64, PyrUserGroupId uint64, CorUserGroupId uint64, ReportUserGroupId uint64, ToolsUserGroupId uint64, UserAccessLevelCode int, IsActive uint8, UpdatedBy string) error {
	var UserGroupAccess = DBVar.User_group_access{
		GeneralUserGroupId:  GeneralUserGroupId,
		UserGroupId:         UserGroupId,
		PosUserGroupId:      PosUserGroupId,
		BanUserGroupId:      BanUserGroupId,
		AccUserGroupId:      AccUserGroupId,
		AstUserGroupId:      AstUserGroupId,
		PyrUserGroupId:      PyrUserGroupId,
		CorUserGroupId:      CorUserGroupId,
		ReportUserGroupId:   ReportUserGroupId,
		ToolsUserGroupId:    ToolsUserGroupId,
		UserAccessLevelCode: UserAccessLevelCode,
		IsActive:            IsActive,
		UpdatedBy:           UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.UserGroupAccess).Where("id=?", Id).Omit("id", "created_at", "created_by", "updated_at").Updates(&UserGroupAccess)
	return result.Error
}

func InsertGeneralUserGroup(DB *gorm.DB, Id chan uint64, AccessModule string, IsActive uint8, CreatedBy string) (uint64, error) {
	var GeneralUserGroup = DBVar.General_user_group{
		AccessModule: AccessModule,
		IsActive:     IsActive,
		CreatedBy:    CreatedBy,
	}
	result := DB.Table(DBVar.TableName.GeneralUserGroup).Create(&GeneralUserGroup)
	Id <- GeneralUserGroup.Id
	return GeneralUserGroup.Id, result.Error
}

func UpdateGeneralUserGroup(DB *gorm.DB, Id uint64, AccessModule string, IsActive uint8, UpdatedBy string) error {
	var GeneralUserGroup = DBVar.General_user_group{
		AccessModule: AccessModule,
		IsActive:     IsActive,
		UpdatedBy:    UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.GeneralUserGroup).Where("id=?", Id).Omit("created_at", "created_by", "updated_at", "id").Updates(&GeneralUserGroup)
	return result.Error
}

func InsertAstUserGroup(DB *gorm.DB, Id chan uint64, AccessForm string, AccessInventoryReceive string, AccessFixedAssetReceive string, AccessSpecial string, CreatedBy string) (uint64, error) {
	var AstUserGroup = DBVar.Ast_user_group{
		AccessForm:              AccessForm,
		AccessInventoryReceive:  AccessInventoryReceive,
		AccessFixedAssetReceive: AccessFixedAssetReceive,
		AccessSpecial:           AccessSpecial,
		CreatedBy:               CreatedBy,
	}
	result := DB.Table(DBVar.TableName.AstUserGroup).Create(&AstUserGroup)
	Id <- AstUserGroup.Id
	return AstUserGroup.Id, result.Error
}

func UpdateAstUserGroup(DB *gorm.DB, Id uint64, AccessForm string, AccessInventoryReceive string, AccessFixedAssetReceive string, AccessSpecial string, UpdatedBy string) error {
	var AstUserGroup = DBVar.Ast_user_group{
		AccessForm:              AccessForm,
		AccessInventoryReceive:  AccessInventoryReceive,
		AccessFixedAssetReceive: AccessFixedAssetReceive,
		AccessSpecial:           AccessSpecial,
		UpdatedBy:               UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.AstUserGroup).Where("id=?", Id).Omit("created_at", "created_by", "updated_at", "id").Updates(&AstUserGroup)
	return result.Error
}

func InsertUser(DB *gorm.DB, Code string, Name string, Password string, UserGroupAccessCode string, CreatedBy string) error {
	var User = DBVar.User{
		Code:                Code,
		Name:                Name,
		Password:            Password,
		UserGroupAccessCode: UserGroupAccessCode,
		IsActive:            1,
		CreatedBy:           CreatedBy,
	}
	result := DB.Table(DBVar.TableName.User).Create(&User)
	return result.Error
}

func UpdateUser(DB *gorm.DB, Id uint64, Name string, Password string, UserGroupAccessCode string, PasswordChanged bool, UpdatedBy string) error {

	if !PasswordChanged {
		Password = ""
	}
	var User = DBVar.User{
		Id:                  Id,
		Name:                Name,
		Password:            Password,
		UserGroupAccessCode: UserGroupAccessCode,
		UpdatedBy:           UpdatedBy,
	}
	result := DB.Debug().Table(DBVar.TableName.User).Omit("created_at", "created_by", "updated_at", "id").Updates(&User)

	return result.Error
}

func InsertAccUserGroup(DB *gorm.DB, Id chan uint64, AccessForm string, AccessSpecial string, AccessInvoice string, PrintInvoiceCount int, CreatedBy string) (uint64, error) {
	var AccUserGroup = DBVar.Acc_user_group{
		AccessForm:        AccessForm,
		AccessSpecial:     AccessSpecial,
		AccessInvoice:     AccessInvoice,
		PrintInvoiceCount: PrintInvoiceCount,
		CreatedBy:         CreatedBy,
	}
	result := DB.Table(DBVar.TableName.AccUserGroup).Create(&AccUserGroup)
	Id <- AccUserGroup.Id
	return AccUserGroup.Id, result.Error
}

func UpdateAccUserGroup(DB *gorm.DB, Id uint64, AccessForm string, AccessSpecial string, AccessInvoice string, PrintInvoiceCount int, UpdatedBy string) error {
	var AccUserGroup = DBVar.Acc_user_group{
		Id:                Id,
		AccessForm:        AccessForm,
		AccessSpecial:     AccessSpecial,
		AccessInvoice:     AccessInvoice,
		PrintInvoiceCount: PrintInvoiceCount,
		UpdatedBy:         UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.AccUserGroup).Where("id=?", Id).Omit("created_at", "created_by", "updated_at", "id").Updates(&AccUserGroup)
	return result.Error
}

func InsertPosUserGroup(DB *gorm.DB, Id chan uint64, AccessForm string, AccessSpecial string, AccessTransactionTerminal string, AccessTableView string, AccessReservation string, CreatedBy string) (uint64, error) {
	var PosUserGroup = DBVar.Pos_user_group{
		AccessForm:                AccessForm,
		AccessSpecial:             AccessSpecial,
		AccessTransactionTerminal: AccessTransactionTerminal,
		AccessTableView:           AccessTableView,
		AccessReservation:         AccessReservation,
		CreatedBy:                 CreatedBy,
	}
	result := DB.Table(DBVar.TableName.PosUserGroup).Create(&PosUserGroup)
	Id <- PosUserGroup.Id
	return PosUserGroup.Id, result.Error
}

func UpdatePosUserGroup(DB *gorm.DB, Id uint64, AccessForm string, AccessSpecial string, AccessTransactionTerminal string, AccessTableView string, AccessReservation string, UpdatedBy string) error {
	var PosUserGroup = DBVar.Pos_user_group{
		Id:                        Id,
		AccessForm:                AccessForm,
		AccessSpecial:             AccessSpecial,
		AccessTransactionTerminal: AccessTransactionTerminal,
		AccessTableView:           AccessTableView,
		AccessReservation:         AccessReservation,
		UpdatedBy:                 UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.PosUserGroup).Where("id=?", Id).Omit("created_at", "created_by", "updated_at", "id").Updates(&PosUserGroup)
	return result.Error
}

func InsertToolsUserGroup(DB *gorm.DB, Id chan uint64, AccessForm string, AccessConfiguration string, AccessCompany string, CreatedBy string) (uint64, error) {
	var ToolsUserGroup = DBVar.Tools_user_group{
		AccessForm:          AccessForm,
		AccessConfiguration: AccessConfiguration,
		AccessCompany:       AccessCompany,
		CreatedBy:           CreatedBy,
	}
	result := DB.Table(DBVar.TableName.ToolsUserGroup).Create(&ToolsUserGroup)
	Id <- ToolsUserGroup.Id
	return ToolsUserGroup.Id, result.Error
}

func UpdateToolsUserGroup(DB *gorm.DB, Id uint64, AccessForm string, AccessConfiguration string, AccessCompany string, UpdatedBy string) error {
	var ToolsUserGroup = DBVar.Tools_user_group{
		AccessForm:          AccessForm,
		AccessConfiguration: AccessConfiguration,
		AccessCompany:       AccessCompany,
		UpdatedBy:           UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.ToolsUserGroup).Where("id=?", Id).Omit("created_at", "created_by", "updated_at", "id").Updates(&ToolsUserGroup)
	return result.Error
}

func InsertCompetitorData(DB *gorm.DB, CompetitorCode string, Date time.Time, AvailableRoom uint, RoomSold uint, AverageRoomRate float64, CreatedBy string) error {
	var CompetitorData = DBVar.Competitor_data{
		CompetitorCode:  CompetitorCode,
		Date:            Date,
		AvailableRoom:   AvailableRoom,
		RoomSold:        RoomSold,
		AverageRoomRate: AverageRoomRate,
		CreatedBy:       CreatedBy,
	}
	result := DB.Table(DBVar.TableName.CompetitorData).Create(&CompetitorData)
	return result.Error
}

func UpdateCompetitorData(DB *gorm.DB, Id uint64, CompetitorCode string, Date time.Time, AvailableRoom uint, RoomSold uint, AverageRoomRate float64, UpdatedBy string) error {
	var CompetitorData = DBVar.Competitor_data{
		Id:              Id,
		CompetitorCode:  CompetitorCode,
		Date:            Date,
		AvailableRoom:   AvailableRoom,
		RoomSold:        RoomSold,
		AverageRoomRate: AverageRoomRate,
		UpdatedBy:       UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.CompetitorData).Omit("created_at", "created_by", "updated_at", "id").Updates(&CompetitorData)
	return result.Error
}

func InsertReportUserGroup(DB *gorm.DB, Id chan uint64, AccessForm string, AccessFoReport string, AccessPosReport string, AccessBanReport string, AccessAccReport string, AccessAstReport string, AccessPyrReport string, AccessCorReport string, AccessPreviewReport string, CreatedBy string) (uint64, error) {
	var ReportUserGroup = DBVar.Report_user_group{
		AccessForm:          AccessForm,
		AccessFoReport:      AccessFoReport,
		AccessPosReport:     AccessPosReport,
		AccessBanReport:     AccessBanReport,
		AccessAccReport:     AccessAccReport,
		AccessAstReport:     AccessAstReport,
		AccessPyrReport:     AccessPyrReport,
		AccessCorReport:     AccessCorReport,
		AccessPreviewReport: AccessPreviewReport,
		CreatedBy:           CreatedBy,
	}
	result := DB.Table(DBVar.TableName.ReportUserGroup).Create(&ReportUserGroup)
	Id <- ReportUserGroup.Id
	return ReportUserGroup.Id, result.Error
}

func UpdateReportUserGroup(DB *gorm.DB, Id uint64, AccessForm string, AccessFoReport string, AccessPosReport string, AccessBanReport string, AccessAccReport string, AccessAstReport string, AccessPyrReport string, AccessCorReport string, AccessPreviewReport string, UpdatedBy string) error {
	var ReportUserGroup = DBVar.Report_user_group{
		AccessForm:          AccessForm,
		AccessFoReport:      AccessFoReport,
		AccessPosReport:     AccessPosReport,
		AccessBanReport:     AccessBanReport,
		AccessAccReport:     AccessAccReport,
		AccessAstReport:     AccessAstReport,
		AccessPyrReport:     AccessPyrReport,
		AccessCorReport:     AccessCorReport,
		AccessPreviewReport: AccessPreviewReport,
		UpdatedBy:           UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.ReportUserGroup).Where("id=?", Id).Omit("created_at", "created_by", "updated_at", "id").Updates(&ReportUserGroup)
	return result.Error
}

func InsertProformaInvoiceDetail(DB *gorm.DB, ReservationNumber uint64, ArrivalDate time.Time, DepartureDate time.Time, Datex time.Time, RoomTypeCode string, RoomRate float64, IsWeekend string, ChargeFrequency string, Userid string) (uint64, error) {
	var ProformaInvoiceDetail = DBVar.Proforma_invoice_detail{
		ReservationNumber: ReservationNumber,
		ArrivalDate:       ArrivalDate,
		DepartureDate:     DepartureDate,
		Datex:             Datex,
		RoomTypeCode:      RoomTypeCode,
		RoomRate:          RoomRate,
		IsWeekend:         IsWeekend,
		ChargeFrequency:   ChargeFrequency,
		Userid:            Userid,
	}
	result := DB.Table(DBVar.TableName.ProformaInvoiceDetail).Create(&ProformaInvoiceDetail)
	return ProformaInvoiceDetail.Id, result.Error
}

func UpdateProformaInvoiceDetail(DB *gorm.DB, ReservationNumber uint64, ArrivalDate time.Time, DepartureDate time.Time, Datex time.Time, RoomTypeCode string, RoomRate float64, IsWeekend string, ChargeFrequency string, Userid string) (uint64, error) {
	var ProformaInvoiceDetail = DBVar.Proforma_invoice_detail{
		ReservationNumber: ReservationNumber,
		ArrivalDate:       ArrivalDate,
		DepartureDate:     DepartureDate,
		Datex:             Datex,
		RoomTypeCode:      RoomTypeCode,
		RoomRate:          RoomRate,
		IsWeekend:         IsWeekend,
		ChargeFrequency:   ChargeFrequency,
		Userid:            Userid,
	}
	result := DB.Table(DBVar.TableName.ProformaInvoiceDetail).Omit("id").Updates(&ProformaInvoiceDetail)
	return ProformaInvoiceDetail.Id, result.Error
}

func InsertGuestProfile(DB *gorm.DB, TitleCode string, FullName string, Street string, CountryCode string, StateCode string, CityCode string, City string, NationalityCode string, PostalCode string, Phone1 string, Phone2 string, Fax string, Email string,
	Website string, CompanyCode string, GuestTypeCode string, IdCardCode string, IdCardNumber string, IsMale uint8, BirthPlace string, BirthDate time.Time, TypeCode string, CustomField01 string, CustomField02 string, CustomField03 string, CustomField04 string,
	CustomField05 string, CustomField06 string, CustomField07 string, CustomField08 string, CustomField09 string, CustomField10 string, CustomField11 string, CustomField12 string, CustomLookupFieldCode01 string, CustomLookupFieldCode02 string,
	CustomLookupFieldCode03 string, CustomLookupFieldCode04 string, CustomLookupFieldCode05 string, CustomLookupFieldCode06 string, CustomLookupFieldCode07 string, CustomLookupFieldCode08 string, CustomLookupFieldCode09 string, CustomLookupFieldCode10 string,
	CustomLookupFieldCode11 string, CustomLookupFieldCode12 string, IsActive uint8, IsBlacklist uint8, CustomerCode string, Source string, CreatedBy string) (uint64, error) {
	var GuestProfile = DBVar.Guest_profile{
		TitleCode:               TitleCode,
		FullName:                FullName,
		Street:                  Street,
		CountryCode:             CountryCode,
		StateCode:               StateCode,
		CityCode:                CityCode,
		City:                    City,
		NationalityCode:         NationalityCode,
		PostalCode:              PostalCode,
		Phone1:                  Phone1,
		Phone2:                  Phone2,
		Fax:                     Fax,
		Email:                   Email,
		Website:                 Website,
		CompanyCode:             CompanyCode,
		GuestTypeCode:           GuestTypeCode,
		IdCardCode:              IdCardCode,
		IdCardNumber:            IdCardNumber,
		IsMale:                  IsMale,
		BirthPlace:              BirthPlace,
		BirthDate:               BirthDate,
		TypeCode:                TypeCode,
		CustomField01:           CustomField01,
		CustomField02:           CustomField02,
		CustomField03:           CustomField03,
		CustomField04:           CustomField04,
		CustomField05:           CustomField05,
		CustomField06:           CustomField06,
		CustomField07:           CustomField07,
		CustomField08:           CustomField08,
		CustomField09:           CustomField09,
		CustomField10:           CustomField10,
		CustomField11:           CustomField11,
		CustomField12:           CustomField12,
		CustomLookupFieldCode01: CustomLookupFieldCode01,
		CustomLookupFieldCode02: CustomLookupFieldCode02,
		CustomLookupFieldCode03: CustomLookupFieldCode03,
		CustomLookupFieldCode04: CustomLookupFieldCode04,
		CustomLookupFieldCode05: CustomLookupFieldCode05,
		CustomLookupFieldCode06: CustomLookupFieldCode06,
		CustomLookupFieldCode07: CustomLookupFieldCode07,
		CustomLookupFieldCode08: CustomLookupFieldCode08,
		CustomLookupFieldCode09: CustomLookupFieldCode09,
		CustomLookupFieldCode10: CustomLookupFieldCode10,
		CustomLookupFieldCode11: CustomLookupFieldCode11,
		CustomLookupFieldCode12: CustomLookupFieldCode12,
		IsActive:                IsActive,
		IsBlacklist:             IsBlacklist,
		CustomerCode:            CustomerCode,
		Source:                  Source,
		CreatedBy:               CreatedBy,
	}
	result := DB.Table(DBVar.TableName.GuestProfile).Create(&GuestProfile)
	return GuestProfile.Id, result.Error
}

func UpdateGuestProfile(DB *gorm.DB, Id uint64, TitleCode string, FullName string, Street string, CountryCode string, StateCode string, CityCode string, City string, NationalityCode string, PostalCode string, Phone1 string, Phone2 string, Fax string, Email string, Website string, CompanyCode string, GuestTypeCode string, IdCardCode string, IdCardNumber string, IsMale uint8, BirthPlace string, BirthDate time.Time, TypeCode string, CustomField01 string, CustomField02 string, CustomField03 string, CustomField04 string, CustomField05 string, CustomField06 string, CustomField07 string, CustomField08 string, CustomField09 string, CustomField10 string, CustomField11 string, CustomField12 string, CustomLookupFieldCode01 string, CustomLookupFieldCode02 string, CustomLookupFieldCode03 string, CustomLookupFieldCode04 string, CustomLookupFieldCode05 string, CustomLookupFieldCode06 string, CustomLookupFieldCode07 string, CustomLookupFieldCode08 string, CustomLookupFieldCode09 string, CustomLookupFieldCode10 string, CustomLookupFieldCode11 string, CustomLookupFieldCode12 string, IsActive uint8, IsBlacklist uint8, CustomerCode string, Source string, UpdatedBy string) error {
	var GuestProfile = DBVar.Guest_profile{
		TitleCode:               TitleCode,
		FullName:                FullName,
		Street:                  Street,
		CountryCode:             CountryCode,
		StateCode:               StateCode,
		CityCode:                CityCode,
		City:                    City,
		NationalityCode:         NationalityCode,
		PostalCode:              PostalCode,
		Phone1:                  Phone1,
		Phone2:                  Phone2,
		Fax:                     Fax,
		Email:                   Email,
		Website:                 Website,
		CompanyCode:             CompanyCode,
		GuestTypeCode:           GuestTypeCode,
		IdCardCode:              IdCardCode,
		IdCardNumber:            IdCardNumber,
		IsMale:                  IsMale,
		BirthPlace:              BirthPlace,
		BirthDate:               BirthDate,
		TypeCode:                TypeCode,
		CustomField01:           CustomField01,
		CustomField02:           CustomField02,
		CustomField03:           CustomField03,
		CustomField04:           CustomField04,
		CustomField05:           CustomField05,
		CustomField06:           CustomField06,
		CustomField07:           CustomField07,
		CustomField08:           CustomField08,
		CustomField09:           CustomField09,
		CustomField10:           CustomField10,
		CustomField11:           CustomField11,
		CustomField12:           CustomField12,
		CustomLookupFieldCode01: CustomLookupFieldCode01,
		CustomLookupFieldCode02: CustomLookupFieldCode02,
		CustomLookupFieldCode03: CustomLookupFieldCode03,
		CustomLookupFieldCode04: CustomLookupFieldCode04,
		CustomLookupFieldCode05: CustomLookupFieldCode05,
		CustomLookupFieldCode06: CustomLookupFieldCode06,
		CustomLookupFieldCode07: CustomLookupFieldCode07,
		CustomLookupFieldCode08: CustomLookupFieldCode08,
		CustomLookupFieldCode09: CustomLookupFieldCode09,
		CustomLookupFieldCode10: CustomLookupFieldCode10,
		CustomLookupFieldCode11: CustomLookupFieldCode11,
		CustomLookupFieldCode12: CustomLookupFieldCode12,
		IsActive:                IsActive,
		IsBlacklist:             IsBlacklist,
		CustomerCode:            CustomerCode,
		Source:                  Source,
		UpdatedBy:               UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.GuestProfile).Where("id=?", Id).Omit("created_at", "created_by", "updated_at", "id").Updates(&GuestProfile)
	return result.Error
}
func InsertAccApArPaymentDetail(DB *gorm.DB, ApArNumber string, RefNumber string, Amount float64, Remark string, CreatedBy string) error {
	var AccApArPaymentDetail = DBVar.Acc_ap_ar_payment_detail{
		ApArNumber: ApArNumber,
		RefNumber:  RefNumber,
		Amount:     Amount,
		Remark:     Remark,
		CreatedBy:  CreatedBy,
	}
	result := DB.Table(DBVar.TableName.AccApArPaymentDetail).Create(&AccApArPaymentDetail)
	return result.Error
}

func UpdateAccApArPaymentDetail(DB *gorm.DB, ApArNumber string, RefNumber string, Amount float64, Remark *string, UpdatedBy string) error {
	var AccApArPaymentDetail = DBVar.Acc_ap_ar_payment_detail{
		ApArNumber: ApArNumber,
		RefNumber:  RefNumber,
		Amount:     Amount,
		Remark:     *Remark,
		UpdatedBy:  UpdatedBy,
	}
	result := DB.Table(DBVar.TableName.AccApArPaymentDetail).Omit("created_at", "created_by", "updated_at", "id").Updates(&AccApArPaymentDetail)
	return result.Error
}

func DeleteAPAR(c *gin.Context, DB *gorm.DB, Number string, IsAP bool, UserID string) error {
	if err := DB.Transaction(func(tx *gorm.DB) error {
		type DataOutputStruct struct {
			Date      time.Time
			RefNumber string
		}
		var DataOutput DataOutputStruct
		if err := tx.Table(DBVar.TableName.AccApAr).Select("date,ref_number").Where("number=?", Number).Take(&DataOutput).Error; err != nil {
			return err
		}
		RefNumber := DataOutput.RefNumber
		err := DB.Table(DBVar.TableName.AccJournal).Where("ref_number=? AND date=?", RefNumber, DataOutput.Date).Updates(map[string]interface{}{"updated_by": UserID}).Error
		if err != nil {
			return err
		}
		err = DB.Table(DBVar.TableName.AccJournal).Where("ref_number=? AND date=?", RefNumber, DataOutput.Date).Delete(&RefNumber).Error
		if err != nil {
			return err
		}
		err = DB.Table(DBVar.TableName.AccJournalDetail).Where("ref_number=? AND date=?", RefNumber, DataOutput.Date).Updates(map[string]interface{}{"updated_by": UserID}).Error
		if err != nil {
			return err
		}
		err = DB.Table(DBVar.TableName.AccJournalDetail).Where("ref_number=? AND date=?", RefNumber, DataOutput.Date).Delete(&RefNumber).Error
		if err != nil {
			return err
		}

		err = tx.Exec("CALL delete_acc_ap_ar(?,?)", Number, UserID).Error
		if err != nil {
			return err
		}

		LogAction := GlobalVar.LogUserActionCAS.DeleteAccountReceivable
		if IsAP {
			LogAction = GlobalVar.LogUserActionCAS.DeleteAccountPayable
		}
		InsertLogUser(tx, GlobalVar.SystemCode.Accounting, LogAction, GetAuditDate(c, DB, false), "", "", "", Number, RefNumber, "", "", UserID)

		return nil
	}); err != nil {
		return err
	}
	return nil
}
func InsertBudgetStatistic(DBTrx *gorm.DB, Period int, SubDepartmentCode string, Code string, Remark string, Amount float64, TypeCode string, M01 float64, M02 float64, M03 float64, M04 float64, M05 float64, M06 float64, M07 float64, M08 float64, M09 float64, M10 float64, M11 float64, M12 float64, UnitCode string, CreatedBy string, IdHolding uint64) error {
	var BudgetStatistic = DBVar.Budget_statistic{
		Period:            Period,
		SubDepartmentCode: SubDepartmentCode,
		Code:              Code,
		Remark:            Remark,
		Amount:            Amount,
		TypeCode:          TypeCode,
		M01:               &M01,
		M02:               &M02,
		M03:               &M03,
		M04:               &M04,
		M05:               &M05,
		M06:               &M06,
		M07:               &M07,
		M08:               &M08,
		M09:               &M09,
		M10:               &M10,
		M11:               &M11,
		M12:               &M12,
		UnitCode:          UnitCode,
		CreatedBy:         CreatedBy,
		IdHolding:         IdHolding,
	}
	result := DBTrx.Table(DBVar.TableName.BudgetStatistic).Create(&BudgetStatistic)
	return result.Error
}

func UpdateBudgetStatistic(DBTrx *gorm.DB, Id uint64, Period int, SubDepartmentCode string, Code string, Remark string, Amount float64, TypeCode string, M01 float64, M02 float64, M03 float64, M04 float64, M05 float64, M06 float64, M07 float64, M08 float64, M09 float64, M10 float64, M11 float64, M12 float64, UnitCode string, UpdatedBy string, IdHolding uint64) error {
	var BudgetStatistic = DBVar.Budget_statistic{
		Period:            Period,
		SubDepartmentCode: SubDepartmentCode,
		Code:              Code,
		Remark:            Remark,
		Amount:            Amount,
		TypeCode:          TypeCode,
		M01:               &M01,
		M02:               &M02,
		M03:               &M03,
		M04:               &M04,
		M05:               &M05,
		M06:               &M06,
		M07:               &M07,
		M08:               &M08,
		M09:               &M09,
		M10:               &M10,
		M11:               &M11,
		M12:               &M12,
		UnitCode:          UnitCode,
		UpdatedBy:         UpdatedBy,
		IdHolding:         IdHolding,
	}
	result := DBTrx.Table(DBVar.TableName.BudgetStatistic).Where("id=?", Id).Omit("created_at", "created_by", "updated_at", "id").Updates(&BudgetStatistic)
	return result.Error
}
func InsertBudgetIncome(DBTrx *gorm.DB, Period int, SubDepartmentCode string, AccountCode string, Remark string, Amount float64, TypeCode string, M01 float64, M02 float64, M03 float64, M04 float64, M05 float64, M06 float64, M07 float64, M08 float64, M09 float64, M10 float64, M11 float64, M12 float64, UnitCode string, CreatedBy string, IdHolding uint64) error {
	var BudgetIncome = DBVar.Budget_income{
		Period:            Period,
		SubDepartmentCode: SubDepartmentCode,
		AccountCode:       AccountCode,
		Remark:            Remark,
		Amount:            Amount,
		TypeCode:          TypeCode,
		M01:               &M01,
		M02:               &M02,
		M03:               &M03,
		M04:               &M04,
		M05:               &M05,
		M06:               &M06,
		M07:               &M07,
		M08:               &M08,
		M09:               &M09,
		M10:               &M10,
		M11:               &M11,
		M12:               &M12,
		UnitCode:          UnitCode,
		CreatedBy:         CreatedBy,
		IdHolding:         IdHolding,
	}
	result := DBTrx.Table(DBVar.TableName.BudgetIncome).Create(&BudgetIncome)
	return result.Error
}

func UpdateBudgetIncome(DBTrx *gorm.DB, Id uint64, Period int, SubDepartmentCode string, AccountCode string, Remark string, Amount float64, TypeCode string, M01 float64, M02 float64, M03 float64, M04 float64, M05 float64, M06 float64, M07 float64, M08 float64, M09 float64, M10 float64, M11 float64, M12 float64, UnitCode string, UpdatedBy string, IdHolding uint64) error {
	var BudgetIncome = DBVar.Budget_income{
		Period:            Period,
		SubDepartmentCode: SubDepartmentCode,
		AccountCode:       AccountCode,
		Remark:            Remark,
		Amount:            Amount,
		TypeCode:          TypeCode,
		M01:               &M01,
		M02:               &M02,
		M03:               &M03,
		M04:               &M04,
		M05:               &M05,
		M06:               &M06,
		M07:               &M07,
		M08:               &M08,
		M09:               &M09,
		M10:               &M10,
		M11:               &M11,
		M12:               &M12,
		UnitCode:          UnitCode,
		UpdatedBy:         UpdatedBy,
		IdHolding:         IdHolding,
	}
	result := DBTrx.Table(DBVar.TableName.BudgetIncome).Where("id=?", Id).Omit("created_at", "created_by", "updated_at", "id").Updates(&BudgetIncome)
	return result.Error
}

func InsertBudgetExpense(DBTrx *gorm.DB, Period int, SubDepartmentCode string, JournalAccountCode string, Remark string, Amount float64, TypeCode string, M01 float64, M02 float64, M03 float64, M04 float64, M05 float64, M06 float64, M07 float64, M08 float64, M09 float64, M10 float64, M11 float64, M12 float64, UnitCode string, CreatedBy string, IdHolding uint64) error {
	var BudgetExpense = DBVar.Budget_expense{
		Period:             Period,
		SubDepartmentCode:  SubDepartmentCode,
		JournalAccountCode: JournalAccountCode,
		Remark:             Remark,
		Amount:             Amount,
		TypeCode:           TypeCode,
		M01:                &M01,
		M02:                &M02,
		M03:                &M03,
		M04:                &M04,
		M05:                &M05,
		M06:                &M06,
		M07:                &M07,
		M08:                &M08,
		M09:                &M09,
		M10:                &M10,
		M11:                &M11,
		M12:                &M12,
		UnitCode:           UnitCode,
		CreatedBy:          CreatedBy,
		IdHolding:          IdHolding,
	}
	result := DBTrx.Table(DBVar.TableName.BudgetExpense).Create(&BudgetExpense)
	return result.Error
}

func UpdateBudgetExpense(DBTrx *gorm.DB, Id uint64, Period int, SubDepartmentCode string, JournalAccountCode string, Remark string, Amount float64, TypeCode string, M01 float64, M02 float64, M03 float64, M04 float64, M05 float64, M06 float64, M07 float64, M08 float64, M09 float64, M10 float64, M11 float64, M12 float64, UnitCode string, UpdatedBy string, IdHolding uint64) error {
	var BudgetExpense = DBVar.Budget_expense{
		Period:             Period,
		SubDepartmentCode:  SubDepartmentCode,
		JournalAccountCode: JournalAccountCode,
		Remark:             Remark,
		Amount:             Amount,
		TypeCode:           TypeCode,
		M01:                &M01,
		M02:                &M02,
		M03:                &M03,
		M04:                &M04,
		M05:                &M05,
		M06:                &M06,
		M07:                &M07,
		M08:                &M08,
		M09:                &M09,
		M10:                &M10,
		M11:                &M11,
		M12:                &M12,
		UnitCode:           UnitCode,
		UpdatedBy:          UpdatedBy,
		IdHolding:          IdHolding,
	}
	result := DBTrx.Table(DBVar.TableName.BudgetExpense).Where("id=?", Id).Omit("created_at", "created_by", "updated_at", "id").Updates(&BudgetExpense)
	return result.Error
}

func InsertBudgetFb(DBTrx *gorm.DB, Period int, OutletCode string, Code string, Remark string, Amount int, TypeCode string, M01 float64, M02 float64, M03 float64, M04 float64, M05 float64, M06 float64, M07 float64, M08 float64, M09 float64, M10 float64, M11 float64, M12 float64, CreatedBy string) error {
	var BudgetFb = DBVar.Budget_fb{
		Period:     Period,
		OutletCode: OutletCode,
		Code:       Code,
		Remark:     Remark,
		Amount:     Amount,
		TypeCode:   TypeCode,
		M01:        M01,
		M02:        M02,
		M03:        M03,
		M04:        M04,
		M05:        M05,
		M06:        M06,
		M07:        M07,
		M08:        M08,
		M09:        M09,
		M10:        M10,
		M11:        M11,
		M12:        M12,
		CreatedBy:  CreatedBy,
	}
	result := DBTrx.Table(DBVar.TableName.BudgetFb).Create(&BudgetFb)
	return result.Error
}

func UpdateBudgetFb(DBTrx *gorm.DB, Id uint64, Period int, OutletCode string, Code string, Remark string, Amount int, TypeCode string, M01 float64, M02 float64, M03 float64, M04 float64, M05 float64, M06 float64, M07 float64, M08 float64, M09 float64, M10 float64, M11 float64, M12 float64, UpdatedBy string) error {
	var BudgetFb = DBVar.Budget_fb{
		Period:     Period,
		OutletCode: OutletCode,
		Code:       Code,
		Remark:     Remark,
		Amount:     Amount,
		TypeCode:   TypeCode,
		M01:        M01,
		M02:        M02,
		M03:        M03,
		M04:        M04,
		M05:        M05,
		M06:        M06,
		M07:        M07,
		M08:        M08,
		M09:        M09,
		M10:        M10,
		M11:        M11,
		M12:        M12,
		UpdatedBy:  UpdatedBy,
	}
	result := DBTrx.Table(DBVar.TableName.BudgetFb).Where("id=?").Omit("created_at", "created_by", "updated_at", "id").Updates(&BudgetFb)
	return result.Error
}

func InsertInvProduction(DBTrx *gorm.DB, Number string, RefNumber string, ReceiveNumber string, CostingNumber string, DocumentNumber string, Date time.Time, Remark string, CreatedBy string) error {
	var InvProduction = DBVar.Inv_production{
		Number:         Number,
		RefNumber:      RefNumber,
		ReceiveNumber:  ReceiveNumber,
		CostingNumber:  CostingNumber,
		DocumentNumber: DocumentNumber,
		Date:           Date,
		Remark:         &Remark,
		CreatedBy:      CreatedBy,
	}
	result := DBTrx.Table(DBVar.TableName.InvProduction).Create(&InvProduction)
	return result.Error
}

func UpdateInvProduction(DBTrx *gorm.DB, Number string, DocumentNumber string, Date time.Time, Remark string, UpdatedBy string) error {
	var InvProduction = DBVar.Inv_production{
		Number:         Number,
		DocumentNumber: DocumentNumber,
		Date:           Date,
		Remark:         &Remark,
		UpdatedBy:      UpdatedBy,
	}
	result := DBTrx.Table(DBVar.TableName.InvProduction).Where("number=?", Number).Omit("created_at", "created_by", "updated_at", "id").Updates(&InvProduction)
	return result.Error
}

func InsertInvReceiving(DBTrx *gorm.DB, Dataset *GlobalVar.TDataset, Number string, RefNumber string, PoNumber string, ApNumber string, CostingNumber string, CompanyCode string, InvoiceNumber string, BankAccountCode string, AmountPayment float64, Date time.Time, IsConsignment uint8, Remark string, IsSeparate uint8, IsDiscountIncome uint8, IsTaxExpense uint8, IsShippingExpense uint8, IsCredit uint8, DueDate time.Time, IsPaid uint8, IsOpname uint8, IsProduction uint8, CreatedBy string) error {
	if Dataset.ProgramConfiguration.ReceiveStockAPTwoDigitDecimal {
		AmountPayment = General.RoundToX2(AmountPayment)
	} else {
		AmountPayment = General.RoundToX3(AmountPayment)
	}
	var InvReceiving = DBVar.Inv_receiving{
		Number:            Number,
		RefNumber:         RefNumber,
		PoNumber:          PoNumber,
		ApNumber:          ApNumber,
		CostingNumber:     CostingNumber,
		CompanyCode:       CompanyCode,
		InvoiceNumber:     InvoiceNumber,
		BankAccountCode:   BankAccountCode,
		AmountPayment:     AmountPayment,
		Date:              Date,
		IsConsignment:     IsConsignment,
		Remark:            &Remark,
		IsSeparate:        IsSeparate,
		IsDiscountIncome:  IsDiscountIncome,
		IsTaxExpense:      IsTaxExpense,
		IsShippingExpense: IsShippingExpense,
		IsCredit:          IsCredit,
		DueDate:           DueDate,
		IsPaid:            IsPaid,
		IsOpname:          IsOpname,
		IsProduction:      IsProduction,
		CreatedBy:         CreatedBy,
	}
	result := DBTrx.Table(DBVar.TableName.InvReceiving).Create(&InvReceiving)
	return result.Error
}

func UpdateInvReceiving(DBTrx *gorm.DB, Dataset *GlobalVar.TDataset, Number string, PoNumber string, ApNumber string, CostingNumber string, CompanyCode string, InvoiceNumber string, BankAccountCode string, AmountPayment float64, Date time.Time, IsConsignment uint8, Remark string, IsSeparate uint8, IsDiscountIncome uint8, IsTaxExpense uint8, IsShippingExpense uint8, IsCredit uint8, DueDate time.Time, IsOpname uint8, IsProduction uint8, UpdatedBy string) error {
	if Dataset.ProgramConfiguration.ReceiveStockAPTwoDigitDecimal {
		AmountPayment = General.RoundToX2(AmountPayment)
	} else {
		AmountPayment = General.RoundToX3(AmountPayment)
	}
	var InvReceiving = DBVar.Inv_receiving{
		Number:            Number,
		PoNumber:          PoNumber,
		ApNumber:          ApNumber,
		CostingNumber:     CostingNumber,
		CompanyCode:       CompanyCode,
		InvoiceNumber:     InvoiceNumber,
		BankAccountCode:   BankAccountCode,
		AmountPayment:     AmountPayment,
		Date:              Date,
		IsConsignment:     IsConsignment,
		Remark:            &Remark,
		IsSeparate:        IsSeparate,
		IsDiscountIncome:  IsDiscountIncome,
		IsTaxExpense:      IsTaxExpense,
		IsShippingExpense: IsShippingExpense,
		IsCredit:          IsCredit,
		DueDate:           DueDate,
		IsOpname:          IsOpname,
		IsProduction:      IsProduction,
		UpdatedBy:         UpdatedBy,
	}
	result := DBTrx.Table(DBVar.TableName.InvReceiving).Where("number=?", Number).Omit("created_at", "created_by", "updated_at", "id").Updates(&InvReceiving)
	return result.Error
}
func InsertInvReceivingDetail(DBTrx *gorm.DB, Dataset *GlobalVar.TDataset, ReceiveNumber string, StoreCode string, StoreId uint64, ItemCode string, ItemId uint64, Date time.Time, PoId uint64, PoQuantity float64, ReceiveQuantity float64, ReceiveUomCode string, ReceivePrice float64, BasicQuantity float64, BasicUomCode string, BasicPrice float64, Quantity float64, TotalPrice float64, Discount float64, Tax float64, Shipping float64, Remark string, ExpireDate time.Time, IsCogs uint8, JournalAccountCode string, ItemGroupCode string, CreatedBy string) (uint64, error) {
	Round := func(value float64, IsTwoDecimal bool) float64 {
		if IsTwoDecimal {
			return General.RoundToX2(value)
		}
		return General.RoundToX3(value)
	}
	if Dataset.ProgramConfiguration.ReceiveStockAPTwoDigitDecimal {
		ReceiveQuantity = Round(ReceiveQuantity, Dataset.ProgramConfiguration.ReceiveStockAPTwoDigitDecimal)
		ReceivePrice = Round(ReceivePrice, Dataset.ProgramConfiguration.ReceiveStockAPTwoDigitDecimal)
		BasicQuantity = Round(BasicQuantity, Dataset.ProgramConfiguration.ReceiveStockAPTwoDigitDecimal)
		BasicPrice = Round(BasicPrice, Dataset.ProgramConfiguration.ReceiveStockAPTwoDigitDecimal)
		TotalPrice = Round(TotalPrice, Dataset.ProgramConfiguration.ReceiveStockAPTwoDigitDecimal)
		Quantity = Round(Quantity, Dataset.ProgramConfiguration.ReceiveStockAPTwoDigitDecimal)
		Discount = Round(Discount, Dataset.ProgramConfiguration.ReceiveStockAPTwoDigitDecimal)
		Tax = Round(Tax, Dataset.ProgramConfiguration.ReceiveStockAPTwoDigitDecimal)
	}

	var InvReceivingDetail = DBVar.Inv_receiving_detail{
		ReceiveNumber:      ReceiveNumber,
		StoreCode:          StoreCode,
		StoreId:            StoreId,
		ItemCode:           ItemCode,
		ItemId:             ItemId,
		Date:               Date,
		PoId:               PoId,
		PoQuantity:         PoQuantity,
		ReceiveQuantity:    ReceiveQuantity,
		ReceiveUomCode:     ReceiveUomCode,
		ReceivePrice:       ReceivePrice,
		BasicQuantity:      BasicQuantity,
		BasicUomCode:       BasicUomCode,
		BasicPrice:         BasicPrice,
		Quantity:           Quantity,
		TotalPrice:         TotalPrice,
		Discount:           Discount,
		Tax:                Tax,
		Shipping:           Shipping,
		Remark:             &Remark,
		ExpireDate:         ExpireDate,
		IsCogs:             IsCogs,
		JournalAccountCode: JournalAccountCode,
		ItemGroupCode:      ItemGroupCode,
		CreatedBy:          CreatedBy,
	}
	result := DBTrx.Table(DBVar.TableName.InvReceivingDetail).Create(&InvReceivingDetail)
	return InvReceivingDetail.Id, result.Error
}

func UpdateInvReceivingDetail(DBTrx *gorm.DB, ReceiveNumber string, StoreCode string, StoreId uint64, ItemCode string, ItemId uint64, Date time.Time, PoId uint64, PoQuantity float64, ReceiveQuantity float64, ReceiveUomCode string, ReceivePrice float64, BasicQuantity float64, BasicUomCode string, BasicPrice float64, Quantity float64, TotalPrice float64, Discount float64, Tax float64, Shipping float64, Remark string, ExpireDate time.Time, IsCogs uint8, JournalAccountCode string, ItemGroupCode string, UpdatedBy string) error {
	var InvReceivingDetail = DBVar.Inv_receiving_detail{
		ReceiveNumber:      ReceiveNumber,
		StoreCode:          StoreCode,
		StoreId:            StoreId,
		ItemCode:           ItemCode,
		ItemId:             ItemId,
		Date:               Date,
		PoId:               PoId,
		PoQuantity:         PoQuantity,
		ReceiveQuantity:    ReceiveQuantity,
		ReceiveUomCode:     ReceiveUomCode,
		ReceivePrice:       ReceivePrice,
		BasicQuantity:      BasicQuantity,
		BasicUomCode:       BasicUomCode,
		BasicPrice:         BasicPrice,
		Quantity:           Quantity,
		TotalPrice:         TotalPrice,
		Discount:           Discount,
		Tax:                Tax,
		Shipping:           Shipping,
		Remark:             &Remark,
		ExpireDate:         ExpireDate,
		IsCogs:             IsCogs,
		JournalAccountCode: JournalAccountCode,
		ItemGroupCode:      ItemGroupCode,
		UpdatedBy:          UpdatedBy,
	}
	result := DBTrx.Table(DBVar.TableName.InvReceivingDetail).Omit("created_at", "created_by", "updated_at", "id").Updates(&InvReceivingDetail)
	return result.Error
}
func InsertInvReturnStock(DBTrx *gorm.DB, Number string, RefNumber string, CostingNumber string, ArNumber string, CompanyCode string, DocumentNumber string, BankAccountCode string, TotalReturn float64, DueDate time.Time, PaymentRemark string, CreatedBy string) error {
	var InvReturnStock = DBVar.Inv_return_stock{
		Number:          Number,
		RefNumber:       RefNumber,
		CostingNumber:   CostingNumber,
		ArNumber:        ArNumber,
		CompanyCode:     CompanyCode,
		DocumentNumber:  DocumentNumber,
		BankAccountCode: BankAccountCode,
		TotalReturn:     TotalReturn,
		DueDate:         DueDate,
		PaymentRemark:   &PaymentRemark,
		CreatedBy:       CreatedBy,
	}
	result := DBTrx.Table(DBVar.TableName.InvReturnStock).Create(&InvReturnStock)
	return result.Error
}

func UpdateInvReturnStock(DBTrx *gorm.DB, Number string, RefNumber string, ArNumber string, CompanyCode string, DocumentNumber string, BankAccountCode string, TotalReturn float64, DueDate time.Time, PaymentRemark string, UpdatedBy string) error {
	var InvReturnStock = DBVar.Inv_return_stock{
		Number:          Number,
		RefNumber:       RefNumber,
		ArNumber:        ArNumber,
		CompanyCode:     CompanyCode,
		DocumentNumber:  DocumentNumber,
		BankAccountCode: BankAccountCode,
		TotalReturn:     TotalReturn,
		DueDate:         DueDate,
		PaymentRemark:   &PaymentRemark,
		UpdatedBy:       UpdatedBy,
	}
	result := DBTrx.Table(DBVar.TableName.InvReturnStock).Where("number=?", Number).Omit("created_at", "created_by", "updated_at", "id").Updates(&InvReturnStock)
	return result.Error
}

func InsertInvOpname(DBTrx *gorm.DB, Number string, RefNumber string, ReceiveNumber string, CostingNumber string, Date time.Time, CreatedBy string) error {
	var InvOpname = DBVar.Inv_opname{
		Number:        Number,
		RefNumber:     RefNumber,
		ReceiveNumber: ReceiveNumber,
		CostingNumber: CostingNumber,
		Date:          Date,
		CreatedBy:     CreatedBy,
	}
	result := DBTrx.Table(DBVar.TableName.InvOpname).Create(&InvOpname)
	return result.Error
}

func UpdateInvOpname(DBTrx *gorm.DB, Number string, RefNumber string, ReceiveNumber string, CostingNumber string, Date time.Time, UpdatedBy string) error {
	var InvOpname = DBVar.Inv_opname{
		Number:        Number,
		RefNumber:     RefNumber,
		ReceiveNumber: ReceiveNumber,
		CostingNumber: CostingNumber,
		Date:          Date,
		UpdatedBy:     UpdatedBy,
	}
	result := DBTrx.Table(DBVar.TableName.InvOpname).Omit("created_at", "created_by", "updated_at", "id").Updates(&InvOpname)
	return result.Error
}
func InsertInvStockTransfer(DBTrx *gorm.DB, Number string, DocumentNumber string, RequestBy string, StoreCode string, Date time.Time, Remark string, IsStoreRequisition uint8, CreatedBy string) error {
	var InvStockTransfer = DBVar.Inv_stock_transfer{
		Number:             Number,
		DocumentNumber:     DocumentNumber,
		RequestBy:          RequestBy,
		StoreCode:          StoreCode,
		Date:               Date,
		Remark:             Remark,
		IsStoreRequisition: IsStoreRequisition,
		CreatedBy:          CreatedBy,
	}
	result := DBTrx.Table(DBVar.TableName.InvStockTransfer).Create(&InvStockTransfer)
	return result.Error
}

func UpdateInvStockTransfer(DBTrx *gorm.DB, Number string, DocumentNumber string, RequestBy string, StoreCode string, Date time.Time, Remark string, IsStoreRequisition uint8, UpdatedBy string) error {
	var InvStockTransfer = DBVar.Inv_stock_transfer{
		Number:             Number,
		DocumentNumber:     DocumentNumber,
		RequestBy:          RequestBy,
		StoreCode:          StoreCode,
		Date:               Date,
		Remark:             Remark,
		IsStoreRequisition: IsStoreRequisition,
		UpdatedBy:          UpdatedBy,
	}
	result := DBTrx.Table(DBVar.TableName.InvStockTransfer).Where("number=?", Number).Omit("created_at", "created_by", "updated_at", "id").Updates(&InvStockTransfer)
	return result.Error
}

func InsertInvStockTransferDetail(DBTrx *gorm.DB, StNumber string, FromStoreCode string, ToStoreCode string, ItemCode string, Quantity float64, UomCode string, ReceiveId uint64, CreatedBy string) error {
	var InvStockTransferDetail = DBVar.Inv_stock_transfer_detail{
		StNumber:      StNumber,
		FromStoreCode: FromStoreCode,
		ToStoreCode:   ToStoreCode,
		ItemCode:      ItemCode,
		Quantity:      Quantity,
		UomCode:       UomCode,
		ReceiveId:     ReceiveId,
		CreatedBy:     CreatedBy,
	}
	result := DBTrx.Table(DBVar.TableName.InvStockTransferDetail).Create(&InvStockTransferDetail)
	return result.Error
}

func UpdateInvStockTransferDetail(DBTrx *gorm.DB, StNumber string, FromStoreCode string, ToStoreCode string, ItemCode string, Quantity float64, UomCode string, Price float64, TotalPrice float64, ReceiveId uint64, UpdatedBy string) error {
	var InvStockTransferDetail = DBVar.Inv_stock_transfer_detail{
		StNumber:      StNumber,
		FromStoreCode: FromStoreCode,
		ToStoreCode:   ToStoreCode,
		ItemCode:      ItemCode,
		Quantity:      Quantity,
		UomCode:       UomCode,
		Price:         Price,
		TotalPrice:    TotalPrice,
		ReceiveId:     ReceiveId,
		UpdatedBy:     UpdatedBy,
	}
	result := DBTrx.Table(DBVar.TableName.InvStockTransferDetail).Omit("created_at", "created_by", "updated_at", "id").Updates(&InvStockTransferDetail)
	return result.Error
}
func InsertInvStoreRequisition(DBTrx *gorm.DB, Number string, SubDepartmentCode string, StoreCode string, Date time.Time, DocumentNumber string, RequestBy string, Remark string, CreatedBy string) error {
	var InvStoreRequisition = DBVar.Inv_store_requisition{
		Number:            Number,
		SubDepartmentCode: SubDepartmentCode,
		StoreCode:         StoreCode,
		Date:              Date,
		DocumentNumber:    DocumentNumber,
		RequestBy:         RequestBy,
		Remark:            &Remark,
		StatusCode:        GlobalVar.StoreRequisitionStatus.NotApproved,
		CreatedBy:         CreatedBy,
	}
	result := DBTrx.Table(DBVar.TableName.InvStoreRequisition).Create(&InvStoreRequisition)
	return result.Error
}

func UpdateInvStoreRequisition(DBTrx *gorm.DB, Number string, SubDepartmentCode string, StoreCode string, Date time.Time, DocumentNumber string, RequestBy string, Remark string, UpdatedBy string) error {

	result := DBTrx.Table(DBVar.TableName.InvStoreRequisition).Where("number=?", Number).Updates(map[string]interface{}{
		"number":              Number,
		"sub_department_code": SubDepartmentCode,
		"store_code":          StoreCode,
		"document_number":     DocumentNumber,
		"request_by":          RequestBy,
		"remark":              Remark,
		"date":                Date,
		"updated_by":          UpdatedBy,
	})
	return result.Error
}

func InsertInvStoreRequisitionDetail(DBTrx *gorm.DB, SrNumber string, FromStoreCode string, ToStoreCode string, ItemCode string, Quantity float64, QuantityApproved float64, Convertion float64, UomCode string, EstimatePrice float64, CreatedBy string) error {
	var InvStoreRequisitionDetail = DBVar.Inv_store_requisition_detail{
		SrNumber:         SrNumber,
		FromStoreCode:    FromStoreCode,
		ToStoreCode:      ToStoreCode,
		ItemCode:         ItemCode,
		Quantity:         Quantity,
		QuantityApproved: QuantityApproved,
		Convertion:       Convertion,
		UomCode:          UomCode,
		EstimatePrice:    EstimatePrice,
		CreatedBy:        CreatedBy,
	}
	result := DBTrx.Table(DBVar.TableName.InvStoreRequisitionDetail).Create(&InvStoreRequisitionDetail)
	return result.Error
}

func UpdateInvStoreRequisitionDetail(DBTrx *gorm.DB, SrNumber string, FromStoreCode string, ToStoreCode string, ItemCode string, Quantity float64, QuantityApproved float64, Convertion float64, UomCode string, EstimatePrice float64, UpdatedBy string) error {
	var InvStoreRequisitionDetail = DBVar.Inv_store_requisition_detail{
		SrNumber:         SrNumber,
		FromStoreCode:    FromStoreCode,
		ToStoreCode:      ToStoreCode,
		ItemCode:         ItemCode,
		Quantity:         Quantity,
		QuantityApproved: QuantityApproved,
		Convertion:       Convertion,
		UomCode:          UomCode,
		EstimatePrice:    EstimatePrice,
		UpdatedBy:        UpdatedBy,
	}
	result := DBTrx.Table(DBVar.TableName.InvStoreRequisitionDetail).Omit("created_at", "created_by", "updated_at", "id").Updates(&InvStoreRequisitionDetail)
	return result.Error
}

func UpdateFAPurchaseOrderIsReceived(DB *gorm.DB, PONumber, UserID string) error {
	if err := DB.Table(DBVar.TableName.FaPurchaseOrder).Where("number=?", PONumber).Updates(map[string]interface{}{"is_received": 1, "updated_by": UserID}).Error; err != nil {
		return err
	}
	return nil
}
func InsertFaReceive(DBTrx *gorm.DB, Number string, RefNumber string, RefNumberOldAsset string, PoNumber string, ApNumber string, CompanyCode string, InvoiceNumber string, BankAccountCode string, AmountPayment float64, Date time.Time, Remark string, IsSeparate uint8, IsDiscountIncome uint8, IsTaxExpense uint8, IsShippingExpense uint8, IsCredit uint8, DueDate time.Time, IsPaid uint8, CreatedBy string) error {
	var FaReceive = DBVar.Fa_receive{
		Number:            Number,
		RefNumber:         RefNumber,
		RefNumberOldAsset: RefNumberOldAsset,
		PoNumber:          PoNumber,
		ApNumber:          ApNumber,
		CompanyCode:       CompanyCode,
		InvoiceNumber:     InvoiceNumber,
		BankAccountCode:   BankAccountCode,
		AmountPayment:     AmountPayment,
		Date:              Date,
		Remark:            Remark,
		IsSeparate:        IsSeparate,
		IsDiscountIncome:  IsDiscountIncome,
		IsTaxExpense:      IsTaxExpense,
		IsShippingExpense: IsShippingExpense,
		IsCredit:          IsCredit,
		DueDate:           DueDate,
		IsPaid:            IsPaid,
		CreatedBy:         CreatedBy,
	}
	result := DBTrx.Table(DBVar.TableName.FaReceive).Create(&FaReceive)
	return result.Error
}

func UpdateFaReceive(DBTrx *gorm.DB, Number string, RefNumberOldAsset string, PoNumber string, ApNumber string, CompanyCode string, InvoiceNumber string, BankAccountCode string, AmountPayment float64, Date time.Time, Remark string, IsSeparate uint8, IsDiscountIncome uint8, IsTaxExpense uint8, IsShippingExpense uint8, IsCredit uint8, DueDate time.Time, UpdatedBy string) error {
	var FaReceive = DBVar.Fa_receive{
		Number:            Number,
		RefNumberOldAsset: RefNumberOldAsset,
		PoNumber:          PoNumber,
		ApNumber:          ApNumber,
		CompanyCode:       CompanyCode,
		InvoiceNumber:     InvoiceNumber,
		BankAccountCode:   BankAccountCode,
		AmountPayment:     AmountPayment,
		Date:              Date,
		Remark:            Remark,
		IsSeparate:        IsSeparate,
		IsDiscountIncome:  IsDiscountIncome,
		IsTaxExpense:      IsTaxExpense,
		IsShippingExpense: IsShippingExpense,
		IsCredit:          IsCredit,
		DueDate:           DueDate,
		UpdatedBy:         UpdatedBy,
	}
	result := DBTrx.Table(DBVar.TableName.FaReceive).Where("number=?", Number).Omit("created_at", "created_by", "updated_at", "id").Updates(&FaReceive)
	return result.Error
}
func InsertFaReceiveDetail(DBTrx *gorm.DB, ReceiveNumber string, ItemCode string, DetailName string, PoQuantity int, PoPrice float64, ReceiveQuantity int, ReceiveUomCode string, ReceivePrice float64, Quantity int, TotalPrice float64, Discount float64, Tax float64, Shipping float64, IsOldAsset uint8, DepreciatedMonth int, DepreciatedValue float64, Remark string, CreatedBy string) (uint64, error) {
	var FaReceiveDetail = DBVar.Fa_receive_detail{
		ReceiveNumber:    ReceiveNumber,
		ItemCode:         ItemCode,
		DetailName:       DetailName,
		PoQuantity:       PoQuantity,
		PoPrice:          PoPrice,
		ReceiveQuantity:  ReceiveQuantity,
		ReceiveUomCode:   ReceiveUomCode,
		ReceivePrice:     ReceivePrice,
		Quantity:         Quantity,
		TotalPrice:       TotalPrice,
		Discount:         Discount,
		Tax:              Tax,
		Shipping:         Shipping,
		IsOldAsset:       &IsOldAsset,
		DepreciatedMonth: DepreciatedMonth,
		DepreciatedValue: DepreciatedValue,
		Remark:           Remark,
		CreatedBy:        CreatedBy,
	}
	result := DBTrx.Table(DBVar.TableName.FaReceiveDetail).Create(&FaReceiveDetail)
	return FaReceiveDetail.Id, result.Error
}

func UpdateFaReceiveDetail(DBTrx *gorm.DB, ReceiveNumber string, ItemCode string, DetailName string, PoQuantity int, PoPrice float64, ReceiveQuantity int, ReceiveUomCode string, ReceivePrice float64, Quantity int, TotalPrice float64, Discount float64, Tax float64, Shipping float64, IsOldAsset uint8, DepreciatedMonth int, DepreciatedValue float64, Remark string, UpdatedBy string) error {
	var FaReceiveDetail = DBVar.Fa_receive_detail{
		ReceiveNumber:    ReceiveNumber,
		ItemCode:         ItemCode,
		DetailName:       DetailName,
		PoQuantity:       PoQuantity,
		PoPrice:          PoPrice,
		ReceiveQuantity:  ReceiveQuantity,
		ReceiveUomCode:   ReceiveUomCode,
		ReceivePrice:     ReceivePrice,
		Quantity:         Quantity,
		TotalPrice:       TotalPrice,
		Discount:         Discount,
		Tax:              Tax,
		Shipping:         Shipping,
		IsOldAsset:       &IsOldAsset,
		DepreciatedMonth: DepreciatedMonth,
		DepreciatedValue: DepreciatedValue,
		Remark:           Remark,
		UpdatedBy:        UpdatedBy,
	}
	result := DBTrx.Table(DBVar.TableName.FaReceiveDetail).Omit("created_at", "created_by", "updated_at", "id").Updates(&FaReceiveDetail)
	return result.Error
}

func DeleteFAReceiveDetail(DB *gorm.DB, ReceiveNumber string, UserID string) error {
	if err := DB.Table(DBVar.TableName.FaReceiveDetail).Where("receive_number=?", ReceiveNumber).Update("updated_by", UserID).Error; err != nil {
		return err
	}
	if err := DB.Table(DBVar.TableName.FaReceiveDetail).Where("receive_number=?", ReceiveNumber).Delete(ReceiveNumber).Error; err != nil {
		return err
	}
	return nil
}

func DeleteFAListByReceiveNumber(DB *gorm.DB, ReceiveNumber string, UserID string) error {
	if err := DB.Table(DBVar.TableName.FaList).Where("receive_number=?", ReceiveNumber).Update("updated_by", UserID).Error; err != nil {
		return err
	}
	if err := DB.Table(DBVar.TableName.FaList).Where("receive_number=?", ReceiveNumber).Delete(ReceiveNumber).Error; err != nil {
		return err
	}
	return nil
}
func InsertFaList(DBTrx *gorm.DB, Code string, Barcode string, ReceiveNumber string, ReceiveId uint64, ItemCode string, SortNumber uint64, Name string, AcquisitionDate time.Time, DepreciationDate time.Time, DepreciationTypeCode string, DepreciationSubDepartmentCode string, DepreciationExpenseAccountCode string, PurchasePrice float64, CurrentValue float64, ResidualValue float64, SerialNumber string, ManufactureCode string, Trademark string, WarrantyDate time.Time, LocationCode string, UsefulLife int, ConditionCode string, Remark string, DepreciationRate float64, FoNumber string, RefNumber1 string, DoNotRevenueJournal uint8, RefNumber2 string, IsOldAsset uint8, DepreciatedMonth int, DepreciatedValue float64, CreatedBy string) error {
	var FaList = DBVar.Fa_list{
		Code:                           Code,
		Barcode:                        &Barcode,
		ReceiveNumber:                  ReceiveNumber,
		ReceiveId:                      ReceiveId,
		ItemCode:                       ItemCode,
		SortNumber:                     &SortNumber,
		Name:                           Name,
		AcquisitionDate:                AcquisitionDate,
		DepreciationDate:               DepreciationDate,
		DepreciationTypeCode:           DepreciationTypeCode,
		DepreciationSubDepartmentCode:  &DepreciationSubDepartmentCode,
		DepreciationExpenseAccountCode: &DepreciationExpenseAccountCode,
		PurchasePrice:                  PurchasePrice,
		CurrentValue:                   CurrentValue,
		ResidualValue:                  ResidualValue,
		SerialNumber:                   &SerialNumber,
		ManufactureCode:                &ManufactureCode,
		Trademark:                      &Trademark,
		WarrantyDate:                   &WarrantyDate,
		LocationCode:                   &LocationCode,
		UsefulLife:                     &UsefulLife,
		ConditionCode:                  ConditionCode,
		Remark:                         &Remark,
		DepreciationRate:               &DepreciationRate,
		FoNumber:                       FoNumber,
		RefNumber1:                     RefNumber1,
		DoNotRevenueJournal:            &DoNotRevenueJournal,
		RefNumber2:                     RefNumber2,
		IsOldAsset:                     &IsOldAsset,
		DepreciatedMonth:               &DepreciatedMonth,
		DepreciatedValue:               &DepreciatedValue,
		CreatedBy:                      CreatedBy,
	}
	result := DBTrx.Table(DBVar.TableName.FaList).Create(&FaList)
	return result.Error
}

func UpdateFaList(DBTrx *gorm.DB, Code string, Barcode string, ReceiveNumber string, ReceiveId uint64, ItemCode string, SortNumber uint64, Name string, AcquisitionDate time.Time, DepreciationDate time.Time, DepreciationTypeCode string, DepreciationSubDepartmentCode string, DepreciationExpenseAccountCode string, PurchasePrice float64, CurrentValue float64, ResidualValue float64, SerialNumber string, ManufactureCode string, Trademark string, WarrantyDate time.Time, LocationCode string, UsefulLife int, ConditionCode string, Remark string, DepreciationRate float64, FoNumber string, RefNumber1 string, DoNotRevenueJournal uint8, RefNumber2 string, IsOldAsset uint8, DepreciatedMonth int, DepreciatedValue float64, UpdatedBy string) error {
	var FaList = DBVar.Fa_list{
		Code:                           Code,
		Barcode:                        &Barcode,
		ReceiveNumber:                  ReceiveNumber,
		ReceiveId:                      ReceiveId,
		ItemCode:                       ItemCode,
		SortNumber:                     &SortNumber,
		Name:                           Name,
		AcquisitionDate:                AcquisitionDate,
		DepreciationDate:               DepreciationDate,
		DepreciationTypeCode:           DepreciationTypeCode,
		DepreciationSubDepartmentCode:  &DepreciationSubDepartmentCode,
		DepreciationExpenseAccountCode: &DepreciationExpenseAccountCode,
		PurchasePrice:                  PurchasePrice,
		CurrentValue:                   CurrentValue,
		ResidualValue:                  ResidualValue,
		SerialNumber:                   &SerialNumber,
		ManufactureCode:                &ManufactureCode,
		Trademark:                      &Trademark,
		WarrantyDate:                   &WarrantyDate,
		LocationCode:                   &LocationCode,
		UsefulLife:                     &UsefulLife,
		ConditionCode:                  ConditionCode,
		Remark:                         &Remark,
		DepreciationRate:               &DepreciationRate,
		FoNumber:                       FoNumber,
		RefNumber1:                     RefNumber1,
		DoNotRevenueJournal:            &DoNotRevenueJournal,
		RefNumber2:                     RefNumber2,
		IsOldAsset:                     &IsOldAsset,
		DepreciatedMonth:               &DepreciatedMonth,
		DepreciatedValue:               &DepreciatedValue,
		UpdatedBy:                      UpdatedBy,
	}
	result := DBTrx.Table(DBVar.TableName.FaList).Omit("created_at", "created_by", "updated_at", "id").Updates(&FaList)
	return result.Error
}

func UpdatePurchaseOrderDetailReceive(DB *gorm.DB, Id uint64, Quantity float64, UserID string) error {
	if err := DB.Debug().Table(DBVar.TableName.InvPurchaseOrderDetail).Where("id=?", Id).Updates(map[string]interface{}{
		"quantity_received": Quantity,
		"updated_by":        UserID,
	}).Error; err != nil {
		return err
	}
	return nil
}

func InsertAccImportJournalLog(DBTrx *gorm.DB, RefNumber string, AuditDate time.Time, CreatedBy string) error {
	var AccImportJournalLog = DBVar.Acc_import_journal_log{
		RefNumber:   RefNumber,
		AuditDate:   AuditDate,
		PostingDate: time.Now(),
		CreatedBy:   CreatedBy,
	}
	result := DBTrx.Table(DBVar.TableName.AccImportJournalLog).Create(&AccImportJournalLog)
	return result.Error
}

func UpdateAccImportJournalLog(DBTrx *gorm.DB, RefNumber string, AuditDate time.Time, PostingDate time.Time, UpdatedBy string) error {
	var AccImportJournalLog = DBVar.Acc_import_journal_log{
		RefNumber:   RefNumber,
		AuditDate:   AuditDate,
		PostingDate: PostingDate,
		UpdatedBy:   UpdatedBy,
	}
	result := DBTrx.Table(DBVar.TableName.AccImportJournalLog).Omit("created_at", "created_by", "updated_at", "id").Updates(&AccImportJournalLog)
	return result.Error
}

func UpdateSubFolioRefNumber(DB *gorm.DB, Id uint64, RefNumber string, UserID string) error {
	if err := DB.Table(DBVar.TableName.SubFolio).Where("id=?", Id).Updates(map[string]interface{}{
		"ref_number": RefNumber,
		"updated_by": UserID,
	}).Error; err != nil {
		return err
	}
	return nil
}

func UpdatePurchaseOrderIsReceived(DB *gorm.DB, PONumber string, UserID string) error {
	if err := DB.Debug().Table(DBVar.TableName.InvPurchaseOrder).Where("number=?", PONumber).Updates(map[string]interface{}{
		"is_received": gorm.Expr("IF(IFNULL(( "+
			"SELECT "+
			"COUNT(po_number) AS A "+
			"FROM "+
			"inv_purchase_order_detail "+
			"WHERE inv_purchase_order_detail.po_number = ? "+
			"AND (quantity_received + quantity_not_received) < quantity),0)>0, "+
			"IF(IFNULL(( "+
			"SELECT "+
			"SUM(inv_purchase_order_detail.quantity_received) "+
			"FROM "+
			"inv_purchase_order_detail "+
			"WHERE inv_purchase_order_detail.po_number = ?),0)>0, '1', '0'), "+
			"'2')", PONumber, PONumber),
		"updated_by": UserID,
	}).Error; err != nil {
		return err
	}
	return nil
}
func InsertFaPurchaseOrder(DBTrx *gorm.DB, Number string, CompanyCode string, ExpeditionCode string, ContactPersonId uint64, ShippingCompany string, ContactPersonShippingId uint64, Date time.Time, RequestBy string, Remark string, CreatedBy string) error {
	var FaPurchaseOrder = DBVar.Fa_purchase_order{
		Number:                  Number,
		CompanyCode:             CompanyCode,
		ExpeditionCode:          &ExpeditionCode,
		ContactPersonId:         ContactPersonId,
		ShippingCompany:         ShippingCompany,
		ContactPersonShippingId: ContactPersonShippingId,
		Date:                    Date,
		RequestBy:               RequestBy,
		Remark:                  &Remark,
		IsReceived:              0,
		CreatedBy:               CreatedBy,
	}
	result := DBTrx.Table(DBVar.TableName.FaPurchaseOrder).Create(&FaPurchaseOrder)
	return result.Error
}

func UpdateFaPurchaseOrder(DBTrx *gorm.DB, Number string, CompanyCode string, ExpeditionCode string, ShippingCompany string, Date time.Time, RequestBy string, Remark string, UpdatedBy string) error {
	var FaPurchaseOrder = DBVar.Fa_purchase_order{
		Number:          Number,
		CompanyCode:     CompanyCode,
		ExpeditionCode:  &ExpeditionCode,
		ShippingCompany: ShippingCompany,
		Date:            Date,
		RequestBy:       RequestBy,
		Remark:          &Remark,
		UpdatedBy:       UpdatedBy,
	}
	result := DBTrx.Table(DBVar.TableName.FaPurchaseOrder).Where("number=?", Number).Omit("created_at", "created_by", "updated_at", "id").Updates(&FaPurchaseOrder)
	return result.Error
}

func UpdateFolioComplimentHU(DBTrx *gorm.DB, FolioNumber uint64, ComplimentHU, UserID string) error {
	if err := DBTrx.Table(DBVar.TableName.Folio).Where("number=?", FolioNumber).Updates(map[string]interface{}{
		"compliment_hu": ComplimentHU,
		"updated_by":    UserID,
	}).Error; err != nil {
		return err
	}
	return nil
}

func UpdateCfgInitRoomRateSession(DBTrx *gorm.DB, RoomRateCode string, FromDate time.Time, ToDate time.Time, Amount float64, UpdatedBy string) error {
	var CfgInitRoomRateSession = DBVar.Cfg_init_room_rate_session{
		RoomRateCode: RoomRateCode,
		FromDate:     FromDate,
		ToDate:       ToDate,
		Amount:       Amount,
		UpdatedBy:    UpdatedBy,
	}
	result := DBTrx.Table(DBVar.TableName.CfgInitRoomRateSession).Omit("created_at", "created_by", "updated_at", "id").Updates(&CfgInitRoomRateSession)
	return result.Error
}

func InsertMember(DBTrx *gorm.DB, Code string, GuestProfileId uint64, IsForRoom uint8, RoomPointTypeCode string, IsForOutlet uint8, OutletPointTypeCode string, IsForBanquet uint8, BanquetPointTypeCode string, OutletDiscountCode string, ExpireDate time.Time, FingerprintTemplate []byte, CreatedBy string) error {
	var Member = DBVar.Member{
		Code:                 Code,
		GuestProfileId:       GuestProfileId,
		IsForRoom:            &IsForRoom,
		RoomPointTypeCode:    &RoomPointTypeCode,
		IsForOutlet:          &IsForOutlet,
		OutletPointTypeCode:  &OutletPointTypeCode,
		IsForBanquet:         &IsForBanquet,
		BanquetPointTypeCode: &BanquetPointTypeCode,
		OutletDiscountCode:   OutletDiscountCode,
		ExpireDate:           ExpireDate,
		FingerprintTemplate:  FingerprintTemplate,
		CreatedBy:            CreatedBy,
	}
	result := DBTrx.Table(DBVar.TableName.Member).Create(&Member)
	return result.Error
}

func UpdateMember(DBTrx *gorm.DB, Code string, GuestProfileId uint64, IsForRoom uint8, RoomPointTypeCode string, IsForOutlet uint8, OutletPointTypeCode string, IsForBanquet uint8, BanquetPointTypeCode string, OutletDiscountCode string, ExpireDate time.Time, FingerprintTemplate []byte, UpdatedBy string) error {
	var Member = DBVar.Member{
		Code:                 Code,
		GuestProfileId:       GuestProfileId,
		IsForRoom:            &IsForRoom,
		RoomPointTypeCode:    &RoomPointTypeCode,
		IsForOutlet:          &IsForOutlet,
		OutletPointTypeCode:  &OutletPointTypeCode,
		IsForBanquet:         &IsForBanquet,
		BanquetPointTypeCode: &BanquetPointTypeCode,
		OutletDiscountCode:   OutletDiscountCode,
		ExpireDate:           ExpireDate,
		FingerprintTemplate:  FingerprintTemplate,
		UpdatedBy:            UpdatedBy,
	}
	result := DBTrx.Table(DBVar.TableName.Member).Where("code", Code).Omit("created_at", "created_by", "updated_at", "id").Updates(&Member)
	return result.Error
}

func InsertPosCfgInitMemberProductDiscount(DBTrx *gorm.DB, OutletCode string, MemberCode string, ProductCode string, DiscountPercent float64, CreatedBy string) error {
	var PosCfgInitMemberProductDiscount = DBVar.Pos_cfg_init_member_product_discount{
		OutletCode:      OutletCode,
		MemberCode:      MemberCode,
		ProductCode:     ProductCode,
		DiscountPercent: DiscountPercent,
		CreatedBy:       CreatedBy,
	}
	result := DBTrx.Table(DBVar.TableName.PosCfgInitMemberProductDiscount).Create(&PosCfgInitMemberProductDiscount)
	return result.Error
}

func UpdatePosCfgInitMemberProductDiscount(DBTrx *gorm.DB, MemberCode string, ProductCode string, DiscountPercent float64, UpdatedBy string) error {
	var PosCfgInitMemberProductDiscount = DBVar.Pos_cfg_init_member_product_discount{
		MemberCode:      MemberCode,
		ProductCode:     ProductCode,
		DiscountPercent: DiscountPercent,
		UpdatedBy:       UpdatedBy,
	}
	result := DBTrx.Table(DBVar.TableName.PosCfgInitMemberProductDiscount).Where("product_code=?", ProductCode).Where("member_code=?", MemberCode).Omit("created_at", "created_by", "updated_at", "id").Updates(&PosCfgInitMemberProductDiscount)
	return result.Error
}

func UpdateMemberPointRedeemed(DB *gorm.DB, Id uint64, Redeem bool) error {
	if err := DB.Table(DBVar.TableName.MemberPoint).Where("id=?", Id).Update("is_redeemed", Redeem).Error; err != nil {
		return err
	}
	return nil
}
func InsertHotelInformation(DBTrx *gorm.DB, Code string, Name string, Street string, City string, CountryCode string, StateCode string, PostalCode string, Phone1 string, Phone2 string, Fax string, Email string, Website string, ImageUrl string) error {
	if err := DBTrx.Exec("TRUNCATE TABLE hotel_information").Error; err != nil {
		return err
	}
	var HotelInformation = DBVar.Hotel_information{
		Code:        Code,
		Name:        Name,
		Street:      Street,
		City:        City,
		CountryCode: CountryCode,
		StateCode:   StateCode,
		PostalCode:  PostalCode,
		Phone1:      Phone1,
		Phone2:      Phone2,
		Fax:         Fax,
		Email:       Email,
		Website:     Website,
		ImageUrl:    ImageUrl,
	}
	result := DBTrx.Table(DBVar.TableName.HotelInformation).Create(&HotelInformation)
	return result.Error
}

func UpdateHotelInformation(DBTrx *gorm.DB, Code string, Name string, Street string, City string, CountryCode string, StateCode string, PostalCode string, Phone1 string, Phone2 string, Fax string, Email string, Website string, ImageUrl string) error {
	var HotelInformation = DBVar.Hotel_information{
		Code:        Code,
		Name:        Name,
		Street:      Street,
		City:        City,
		CountryCode: CountryCode,
		StateCode:   StateCode,
		PostalCode:  PostalCode,
		Phone1:      Phone1,
		Phone2:      Phone2,
		Fax:         Fax,
		Email:       Email,
		Website:     Website,
		ImageUrl:    ImageUrl,
	}
	result := DBTrx.Table(DBVar.TableName.HotelInformation).Updates(&HotelInformation)
	return result.Error
}

func InsertCmUpdateAvailability(DBTrx *gorm.DB, StartDate time.Time, EndDate time.Time, RoomTypeCode string, Availability int, Status string) error {
	var CmUpdateAvailability = DBVar.Cm_update_availability{
		StartDate:    StartDate,
		EndDate:      EndDate,
		RoomTypeCode: RoomTypeCode,
		Availability: Availability,
		Status:       Status,
	}
	result := DBTrx.Table(DBVar.TableName.CmUpdateAvailability).Create(&CmUpdateAvailability)
	return result.Error
}

func UpdateCmUpdateAvailability(DBTrx *gorm.DB, StartDate time.Time, EndDate time.Time, RoomTypeCode string, Availability int, Status string) error {
	var CmUpdateAvailability = DBVar.Cm_update_availability{
		StartDate:    StartDate,
		EndDate:      EndDate,
		RoomTypeCode: RoomTypeCode,
		Availability: Availability,
		Status:       Status,
	}
	result := DBTrx.Table(DBVar.TableName.CmUpdateAvailability).Omit("id", "created_at", "updated_at").Updates(&CmUpdateAvailability)
	return result.Error
}

func InsertCmUpdateRate(DBTrx *gorm.DB, StartDate time.Time, EndDate time.Time, RoomRateCode string, RateAmount float64, RoomTypeCode string, BedTypeCode string, Day1 uint8, Day2 uint8, Day3 uint8, Day4 uint8, Day5 uint8, Day6 uint8, Day7 uint8, StopSell uint8, ClosedToArrival uint8, ClosedToDeparture uint8, Status string) error {
	var CmUpdateRate = DBVar.Cm_update_rate{
		StartDate:         StartDate,
		EndDate:           EndDate,
		RoomRateCode:      RoomRateCode,
		RateAmount:        RateAmount,
		RoomTypeCode:      RoomTypeCode,
		BedTypeCode:       BedTypeCode,
		Day1:              Day1,
		Day2:              Day2,
		Day3:              Day3,
		Day4:              Day4,
		Day5:              Day5,
		Day6:              Day6,
		Day7:              Day7,
		StopSell:          StopSell,
		ClosedToArrival:   ClosedToArrival,
		ClosedToDeparture: ClosedToDeparture,
		Status:            Status,
	}
	result := DBTrx.Table(DBVar.TableName.CmUpdateRate).Create(&CmUpdateRate)
	return result.Error
}

func UpdateCmUpdateRate(DBTrx *gorm.DB, StartDate time.Time, EndDate time.Time, RoomRateCode string, RateAmount float64, RoomTypeCode string, BedTypeCode string, Day1 uint8, Day2 uint8, Day3 uint8, Day4 uint8, Day5 uint8, Day6 uint8, Day7 uint8, StopSell uint8, ClosedToArrival uint8, ClosedToDeparture uint8, Status string) error {
	var CmUpdateRate = DBVar.Cm_update_rate{
		StartDate:         StartDate,
		EndDate:           EndDate,
		RoomRateCode:      RoomRateCode,
		RateAmount:        RateAmount,
		RoomTypeCode:      RoomTypeCode,
		BedTypeCode:       BedTypeCode,
		Day1:              Day1,
		Day2:              Day2,
		Day3:              Day3,
		Day4:              Day4,
		Day5:              Day5,
		Day6:              Day6,
		Day7:              Day7,
		StopSell:          StopSell,
		ClosedToArrival:   ClosedToArrival,
		ClosedToDeparture: ClosedToDeparture,
		Status:            Status,
	}
	result := DBTrx.Table(DBVar.TableName.CmUpdateRate).Omit("id", "created_at", "updated_at").Updates(&CmUpdateRate)
	return result.Error
}

func InsertPosCfgInitMemberOutletDiscountDetail(DBTrx *gorm.DB, MemberOutletDiscountCode string, OutletCode string, ProductCode string, DiscountPercent float64, CreatedBy string) error {
	var PosCfgInitMemberOutletDiscountDetail = DBVar.Pos_cfg_init_member_outlet_discount_detail{
		MemberOutletDiscountCode: MemberOutletDiscountCode,
		OutletCode:               OutletCode,
		ProductCode:              ProductCode,
		DiscountPercent:          DiscountPercent,
		CreatedBy:                CreatedBy,
	}
	result := DBTrx.Table(DBVar.TableName.PosCfgInitMemberOutletDiscountDetail).Create(&PosCfgInitMemberOutletDiscountDetail)
	return result.Error
}

func UpdatePosCfgInitMemberOutletDiscountDetail(DBTrx *gorm.DB, Id uint64, MemberOutletDiscountCode string, ProductCode string, DiscountPercent float64, UpdatedBy string) error {
	var PosCfgInitMemberOutletDiscountDetail = DBVar.Pos_cfg_init_member_outlet_discount_detail{
		MemberOutletDiscountCode: MemberOutletDiscountCode,
		ProductCode:              ProductCode,
		DiscountPercent:          DiscountPercent,
		UpdatedBy:                UpdatedBy,
	}
	result := DBTrx.Table(DBVar.TableName.PosCfgInitMemberOutletDiscountDetail).Omit("created_at", "created_by", "updated_at", "id").Where("id", Id).Updates(&PosCfgInitMemberOutletDiscountDetail)
	return result.Error
}

// CALL PROCEDURE========================================================================================================================

func InsertCfgInitRoomRateSession(DBTrx *gorm.DB, RoomRateCode string, FromDate time.Time, ToDate time.Time, Amount float64, IsDefault uint8, CreatedBy string) error {
	result := DBTrx.Debug().Exec("CALL insert_cfg_init_room_rate_session(?,?,?,?,?,?)", RoomRateCode, General.FormatDate1(FromDate), General.FormatDate1(ToDate), Amount, fmt.Sprintf("%d", IsDefault), CreatedBy)
	return result.Error
}

func InsertRoomRateLastDeal(DBTrx *gorm.DB, RoomRateCode string, StartTime string, EndTime string, Percentage float64, IsDefault uint8, CreatedBy string) error {
	result := DBTrx.Debug().Exec("CALL insert_cfg_init_room_rate_last_deal(?,?,?,?,?,?)", RoomRateCode, StartTime, EndTime, Percentage, fmt.Sprintf("%d", IsDefault), CreatedBy)
	return result.Error
}

func DeleteStoreRequisitionDetail(DB *gorm.DB, Number, UserID string) error {
	err := DB.Exec("CALL delete_inv_store_requisition_detail(?,?)", Number, UserID).Error
	return err
}

func DeleteStoreRequisition(DB *gorm.DB, Number, UserID string) error {
	err := DB.Exec("CALL delete_inv_store_requisition(?,?)", Number, UserID).Error
	return err
}

func DeleteJournal(ctx context.Context, DB *gorm.DB, RefNumber, UserID string) error {
	if RefNumber == "" {
		return nil
	}
	err := DB.Exec("CALL delete_acc_journal(?,?)", RefNumber, UserID).Error
	return err
}

func DeleteJournalDetail(DB *gorm.DB, RefNumber string, UserID string) error {
	var DateX time.Time
	if err := DB.Table(DBVar.TableName.AccJournal).Select("date").Where("ref_number=?", RefNumber).Limit(1).Scan(&DateX).Error; err != nil {
		return err
	}
	if !DateX.IsZero() {
		if err := DB.Exec("CALL delete_acc_journal_detail(?,?,?)", RefNumber, General.FormatDate1(DateX), UserID).Error; err != nil {
			return err
		}
	}
	return nil
}

func DeleteStockTransfer(DB *gorm.DB, Number, UserID string) error {
	err := DB.Exec("CALL delete_inv_stock_transfer(?,?)", Number, UserID).Error
	return err
}

func DeleteStockTransferDetail(DB *gorm.DB, Number, UserID string) error {
	err := DB.Exec("CALL delete_inv_stock_transfer_detail(?,?)", Number, UserID).Error
	return err
}
func DeleteCosting(DB *gorm.DB, Number, UserID string) error {
	err := DB.Exec("CALL delete_inv_costing(?,?)", Number, UserID).Error
	return err
}

func DeleteCostingDetail(DB *gorm.DB, Number, UserID string) error {
	err := DB.Exec("CALL delete_inv_costing_detail(?,?)", Number, UserID).Error
	return err
}

func VoidSubFolioByBreakdown1(ctx context.Context, DB *gorm.DB, BelongsTo uint64, Breakdown1 uint64, VoidBy, VoidReason, UserID string) error {
	result := DB.WithContext(ctx).Exec("CALL update_sub_folio_void_by_breakdown1(?,?,?,?,?)", BelongsTo, Breakdown1, VoidBy, VoidReason, UserID)
	return result.Error
}

func DeleteAPARProc(DB *gorm.DB, Number, UserID string) error {
	err := DB.Exec("CALL delete_acc_ap_ar(?,?)", Number, UserID).Error
	return err
}

func DeleteAPARPayment(DB *gorm.DB, RefNumber string, UserID string) error {
	if err := DB.Exec("CALL delete_acc_ap_ar_payment_by_ref_number(?,?)", RefNumber, UserID).Error; err != nil {
		return err
	}
	return nil
}

func DeleteImportJournalLog(DB *gorm.DB, Date time.Time, UserID string) error {
	if err := DB.Exec("CALL delete_acc_import_journal_log(?,?,?)", Date, GlobalVar.JournalPrefix.Transaction, UserID).Error; err != nil {
		return err
	}
	return nil
}
func VoidSubFolioByCorrectionBreakdown(ctx context.Context, DB *gorm.DB, BelongsTo uint64, CorrectionBreakdown uint64, VoidBy, VoidReason, UserID string) error {
	ctx, span := global_var.Tracer.Start(ctx, "VoidSubFolioByCorrectionBreakdown")
	defer span.End()

	// result := DB.WithContext(ctx).Exec("CALL update_sub_folio_void_by_correction_breakdown(?,?,?,?,?)", BelongsTo, CorrectionBreakdown, VoidBy, VoidReason, UserID)
	// return result.Error

	var subFolios []db_var.Sub_folio
	if err := DB.WithContext(ctx).Table(db_var.TableName.SubFolio).Where("correction_breakdown = ?", CorrectionBreakdown).Find(&subFolios).Error; err != nil {
		return err
	}

	for _, subFolio := range subFolios {
		// Delete acc_foreign_cash_by_transaction
		// Assuming you have a function deleteAccForeignCashByTransaction in your code
		if err := DeleteAccForeignCashByTransaction(ctx, DB, subFolio.Id, 2, UserID); err != nil {
			return err
		}

		if subFolio.IsPairWithDeposit == 0 {
			if err := UpdateSubFolioVoid(ctx, DB, subFolio.TransferPairId, VoidBy, VoidReason, UserID); err != nil {
				return err
			}
		} else {
			if err := UpdateGuestDepositVoid(ctx, DB, subFolio.TransferPairId, VoidBy, VoidReason, UserID); err != nil {
				return err
			}
		}
	}

	// Update sub_folio
	if err := DB.WithContext(ctx).Table(db_var.TableName.SubFolio).
		Where("correction_breakdown = ?", CorrectionBreakdown).
		Updates(map[string]interface{}{
			"void":        1,
			"void_date":   time.Now(),
			"void_by":     VoidBy,
			"void_reason": VoidReason,
			"updated_by":  UserID,
		}).Error; err != nil {
		return err
	}

	return nil
}

func UpdateGuestDepositVoid(ctx context.Context, DB *gorm.DB, DepositID uint64, VoidBy, VoidReason, UserID string) error {
	ctx, span := global_var.Tracer.Start(ctx, "UpdateGuestDepositVoid")
	defer span.End()

	if DepositID == 0 {
		return nil
	}
	// Fetching data from guest_deposit
	var GuestDeposit []db_var.Guest_deposit
	if err := DB.WithContext(ctx).Table(db_var.TableName.GuestDeposit).Select("is_pair_with_folio, transfer_pair_id, void").
		Where("id = ?", DepositID).Scan(&GuestDeposit).Error; err != nil {
		return err
	}

	for _, deposit := range GuestDeposit {
		if deposit.IsPairWithFolio == 0 && deposit.Void == 0 {
			if err := DB.WithContext(ctx).Table(db_var.TableName.GuestDeposit).Where("id = ?", deposit.Id).Updates(map[string]interface{}{
				"void":        1,
				"void_date":   time.Now(),
				"void_by":     VoidBy,
				"void_reason": VoidReason,
				"updated_by":  UserID,
			}).Error; err != nil {
				return err
			}

			queryFC := DB.WithContext(ctx).Table(db_var.TableName.AccForeignCash).Where("id_transaction = ? AND id_table = ?", DepositID, 1)
			if err := queryFC.Update("updated_by", UserID).Error; err != nil {
				return err
			}

			if err := queryFC.Delete(DepositID).Error; err != nil {
				return err
			}

			if err := UpdateGuestDepositVoid(ctx, DB, deposit.TransferPairId, VoidBy, VoidReason, UserID); err != nil {
				return err
			}
		} else {
			if err := UpdateSubFolioVoid(ctx, DB, deposit.TransferPairId, VoidBy, VoidReason, UserID); err != nil {
				return err
			}
		}
	}

	// Update guest_deposit and delete related acc_foreign_cash
	if err := DB.WithContext(ctx).Table(db_var.TableName.GuestDeposit).Where("id = ?", DepositID).Updates(map[string]interface{}{
		"void":        1,
		"void_date":   time.Now(),
		"void_by":     VoidBy,
		"void_reason": VoidReason,
		"updated_by":  UserID,
	}).Error; err != nil {
		return err
	}

	if err := DeleteAccForeignCashByTransaction(ctx, DB, uint64(DepositID), 1, UserID); err != nil {
		return err
	}

	return nil
}

func UpdateSubFolioVoid(ctx context.Context, DB *gorm.DB, SubFolioID uint64, VoidBy string, VoidReason string, UserID string) error {
	ctx, span := global_var.Tracer.Start(ctx, "UpdateSubFolioVoid")
	defer span.End()

	if SubFolioID == 0 {
		return nil
	}
	if err := DeleteAccForeignCashByTransaction(ctx, DB, SubFolioID, 2, UserID); err != nil {
		return err
	}
	if err := DB.WithContext(ctx).Table(db_var.TableName.SubFolio).Where("id", SubFolioID).Updates(map[string]interface{}{
		"void":        1,
		"void_date":   time.Now(),
		"void_by":     VoidBy,
		"void_reason": VoidReason,
		"updated_by":  UserID,
	}).Error; err != nil {
		return err
	}
	return nil
}

func DeleteAccForeignCashByTransaction(ctx context.Context, DB *gorm.DB, IDTransaction uint64, IDTable uint64, UserID string) error {
	ctx, span := global_var.Tracer.Start(ctx, "DeleteAccForeignCashByTransaction")
	defer span.End()

	Query := DB.WithContext(ctx).Table(db_var.TableName.AccForeignCash).Where("id_transaction=? AND id_table=?", IDTransaction, IDTable)
	if err := Query.Update("updated_by", UserID).Error; err != nil {
		return err
	}
	if err := Query.Delete(IDTransaction).Error; err != nil {
		return err
	}

	return nil
}

func DeleteAccForeignCashByRefNumber(ctx context.Context, DB *gorm.DB, RefNumber string, IDTable uint64, UserID string) error {
	ctx, span := global_var.Tracer.Start(ctx, "DeleteAccForeignCashByRefNumber")
	defer span.End()

	Query := DB.WithContext(ctx).Table(db_var.TableName.AccForeignCash).Where("ref_number=? AND id_table=?", RefNumber, IDTable)
	if err := Query.Update("updated_by", UserID).Error; err != nil {
		return err
	}
	if err := Query.Delete(RefNumber).Error; err != nil {
		return err
	}

	return nil
}

func CheckOutFolio(DB *gorm.DB, FolioNumber uint64, Departure time.Time, UserID string) error {
	result := DB.Exec("CALL update_folio_status_closed(?,?,?)", FolioNumber, Departure, UserID)
	return result.Error
}

func CancelCheckOutFolio(DB *gorm.DB, FolioNumber uint64, GuestDetailId uint64, UserID string) error {
	result := DB.Table(DBVar.TableName.Folio).Where("number=?", FolioNumber).Updates(map[string]interface{}{
		"status_code":  "O",
		"check_out_at": "0000-00-00 00:00:00",
		"check_out_by": "",
		"updated_by":   UserID,
	}).Error

	if result != nil {
		return result
	}
	result = DB.Table(DBVar.TableName.GuestDetail).Where("id=?", GuestDetailId).Updates(map[string]interface{}{
		"departure":       DB.Raw(`ADDDATE(DATE(departure), INTERVAL 1 DAY)`),
		"departure_unixx": DB.Raw("UNIX_TIMESTAMP(ADDDATE(DATE(departure), INTERVAL 1 DAY))"),
	}).Error

	return result
}

func DeleteInventoryReceiveDetail(DB *gorm.DB, ReceiveNumber string, UserID string) error {
	if err := DB.Exec("CALL delete_inv_receive_detail(?,?)", ReceiveNumber, UserID).Error; err != nil {
		return err
	}
	return nil
}

func DeleteInvoicePaymentByRefNumber(ctx context.Context, DB *gorm.DB, RefNumber string, UserID string) error {
	ctx, span := global_var.Tracer.Start(ctx, "DeleteInvoicePaymentByRefNumber")
	defer span.End()

	type InvoicePayment struct {
		SubFolioID uint64
		Amount     float64
	}
	var invoicePayments []InvoicePayment
	DB.WithContext(ctx).Table(db_var.TableName.InvoicePayment).Where("ref_number = ?", RefNumber).Scan(&invoicePayments)

	// Process each invoice payment
	for i, payment := range invoicePayments {
		if i == 1 {
			if err := DeleteJournal(ctx, DB, RefNumber, UserID); err != nil {
				return err
			}
		}
		if err := UpdateInvoiceItemAmountPaid(ctx, DB, payment.SubFolioID, -payment.Amount, UserID); err != nil {
			return err
		}
	}

	if err := DeleteAccForeignCashByRefNumber(ctx, DB, RefNumber, 31, UserID); err != nil {
		return err
	}

	if err := DB.WithContext(ctx).Table(db_var.TableName.InvoicePayment).Where("ref_number = ?", RefNumber).Update("updated_by", UserID).Error; err != nil {
		return err
	}
	if err := DB.WithContext(ctx).Table(db_var.TableName.InvoicePayment).Where("ref_number = ?", RefNumber).Delete(RefNumber).Error; err != nil {
		return err
	}
	return nil
}

func UpdateAPARPaid(ctx context.Context, DB *gorm.DB, APARNumber string, Outstanding, Amount float64, UserID string) error {
	IsPaid := Outstanding == Amount
	if err := DB.WithContext(ctx).Table(DBVar.TableName.AccApAr).Where("number=?", APARNumber).Updates(map[string]interface{}{
		"updated_by":  UserID,
		"amount_paid": gorm.Expr("amount_paid + ?", Amount),
		"is_paid":     IsPaid,
	}).Error; err != nil {
		return err
	}

	if IsPaid {
		if err := DB.WithContext(ctx).Table(db_var.TableName.FaReceive).Where("ap_number", APARNumber).Updates(map[string]interface{}{
			"is_paid":    IsPaid,
			"updated_by": UserID,
		}).Error; err != nil {
			return err
		}
		if err := DB.WithContext(ctx).Table(db_var.TableName.InvReceiving).Where("ap_number", APARNumber).Updates(map[string]interface{}{
			"is_paid":    IsPaid,
			"updated_by": UserID,
		}).Error; err != nil {
			return err
		}
	}

	return nil
}

func UpdateInvoiceItemAmountPaid(ctx context.Context, DB *gorm.DB, SubFolioID uint64, Amount float64, UserID string) error {
	if err := DB.WithContext(ctx).Table(db_var.TableName.InvoiceItem).Where("sub_folio_id", SubFolioID).
		Updates(map[string]interface{}{
			"amount_paid": gorm.Expr("amount_paid + ?", Amount),
			"updated_by":  UserID}).Error; err != nil {
		return err
	}
	return nil
}

func DeleteAccAPARPaymentByRefNumber(ctx context.Context, DB *gorm.DB, RefNumber string, UserID string) error {
	ctx, span := global_var.Tracer.Start(ctx, "DeleteAccAPARPaymentByRefNumber", trace.WithAttributes(attribute.String("RefNumber", RefNumber)))
	defer span.End()

	type AccAPARPaymentDetail struct {
		APARNumber string
		Amount     float64
	}
	var accAPARPaymentDetails []AccAPARPaymentDetail
	DB.WithContext(ctx).Table(db_var.TableName.AccApArPaymentDetail).Where("ref_number = ?", RefNumber).Scan(&accAPARPaymentDetails)

	// Process each APAR payment detail
	for _, paymentDetail := range accAPARPaymentDetails {
		DB.WithContext(ctx).Table(db_var.TableName.AccApAr).Where("number = ?", paymentDetail.APARNumber).Updates(map[string]interface{}{
			"amount_paid": gorm.Expr("amount_paid - ?", paymentDetail.Amount),
			"is_paid":     "0",
			"updated_by":  UserID,
		})
	}

	if err := DeleteJournal(ctx, DB, RefNumber, UserID); err != nil {
		return err
	}

	if err := DB.WithContext(ctx).Table(db_var.TableName.AccApArPayment).Where("ref_number = ?", RefNumber).Update("updated_by", UserID).Error; err != nil {
		return err
	}
	if err := DB.WithContext(ctx).Table(db_var.TableName.AccApArPayment).Where("ref_number = ?", RefNumber).Delete(RefNumber).Error; err != nil {
		return err
	}

	if err := DB.WithContext(ctx).Table(db_var.TableName.AccApArPaymentDetail).Where("ref_number = ?", RefNumber).Update("updated_by", UserID).Error; err != nil {
		return err
	}
	if err := DB.WithContext(ctx).Table(db_var.TableName.AccApArPaymentDetail).Where("ref_number = ?", RefNumber).Delete(RefNumber).Error; err != nil {
		return err
	}
	return nil
}

func DeleteAccAPCommissionPaymentByRefNumber(ctx context.Context, DB *gorm.DB, RefNumber string, UserID string) error {
	ctx, span := global_var.Tracer.Start(ctx, "DeleteAccAPCommissionPaymentByRefNumber", trace.WithAttributes(attribute.String("RefNumber", RefNumber)))
	defer span.End()

	if err := DeleteJournal(ctx, DB, RefNumber, UserID); err != nil {
		return err
	}

	if err := DB.WithContext(ctx).Table(db_var.TableName.AccApCommissionPayment).Where("ref_number = ?", RefNumber).Update("updated_by", UserID).Error; err != nil {
		return err
	}
	if err := DB.WithContext(ctx).Table(db_var.TableName.AccApCommissionPayment).Where("ref_number = ?", RefNumber).Delete(RefNumber).Error; err != nil {
		return err
	}

	if err := DB.WithContext(ctx).Table(db_var.TableName.AccApCommissionPaymentDetail).Where("ref_number = ?", RefNumber).Update("updated_by", UserID).Error; err != nil {
		return err
	}
	if err := DB.WithContext(ctx).Table(db_var.TableName.AccApCommissionPaymentDetail).Where("ref_number = ?", RefNumber).Delete(RefNumber).Error; err != nil {
		return err
	}
	return nil
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

// func DeleteInventoryCloseSummaryOnCloseMonth(AuditDate time.Time) error {
// 	DeleteDate := IncDayX(DeleteDate, -1)
// 	DeleteSQLX("inv_close_summary_on_close_month", "'"+FormatDateTimeX(DeleteDate)+"'", False)
// }
