package db_var

import (
	"time"

	"gorm.io/gorm"
)

type TTableName struct {
	AccApAr, AccApArPayment, AccApArPaymentDetail, AccApCommissionPayment, AccApCommissionPaymentDetail, AccApRefundDepositPayment, AccApRefundDepositPaymentDetail,
	AccCashSaleRecon, AccCfgInitBankAccount, AccCloseMonth, AccCloseYear, AccConstBankAccountType, AccConstJournalGroup, AccConstJournalType, AccConstUnit,
	AccCreditCardRecon, AccCreditCardReconDetail, AccDefferedIncome, AccDefferedIncomePosted, AccForeignCash, AccImportJournalLog, AccJournal, CmLog,
	AccJournalDetail, AccPrepaidExpense, AccPrepaidExpensePosted, AccReport, AccReportDefaultField, AccReportGroupField, AccReportGroupingField,
	AccReportOrderField, AccReportTemplate, AccReportTemplateField, AccUserGroup, AstCfgInitShippingAddress, AstConstPurchaseRequestStatus,
	AstConstStoreRequisitionStatus, AstReport, AstReportDefaultField, AstReportGroupField, AstReportGroupingField, AstReportOrderField, AstReportTemplate,
	AstReportTemplateField, AstUserGroup, AstUserSubDepartment, AuditLog, BanBooking, BanCfgInitSeatingPlan, BanCfgInitTheme, BanCfgInitVenue,
	BanCfgInitVenueCombine, BanCfgInitLayout, BanCfgInitVenueCombineDetail, BanCfgInitVenueGroup, BanConstBookingStatus, BanConstReservationStatus, BanConstReservationType,
	BanConstVenueLocation, BanReport, BanReportDefaultField, BanReportGroupField, BanReportGroupingField, BanReportOrderField, BanReportTemplate,
	BanReportTemplateField, BanReservation, BanReservationCharge, BanReservationRemark, BanUserGroup, BreakfastListTemp, BudgetExpense, BudgetFb,
	BudgetIncome, BudgetStatistic, CashCount, CfgInitAccount, CfgInitAccountSubGroup, CfgInitBedType, CfgInitBookingSource, CfgInitCardBank,
	CfgInitCardType, CfgInitCity, CfgInitCompanyType, CfgInitCompetitorCategory, CfgInitContinent, CfgInitCountry, CfgInitCreditCardCharge,
	CfgInitCurrency, CfgInitCurrencyNominal, CfgInitCustomLookupField01, CfgInitCustomLookupField02, CfgInitCustomLookupField03, InvConstItemGroupType,
	CfgInitCustomLookupField04, CfgInitCustomLookupField05, CfgInitCustomLookupField06, CfgInitCustomLookupField07, CfgInitCustomLookupField08,
	CfgInitCustomLookupField09, CfgInitCustomLookupField10, CfgInitCustomLookupField11, CfgInitCustomLookupField12, CfgInitDepartment, CfgInitGuestType,
	CfgInitIdCardType, CfgInitIsFbSubDepartmentGroup, CfgInitIsFbSubDepartmentGroupDetail, CfgInitJournalAccount, CfgInitJournalAccountCategory,
	CfgInitJournalAccountSubGroup, CfgInitLanguage, CfgInitLoanItem, CfgInitMarket, CfgInitMarketCategory, CfgInitMemberPointType,
	CfgInitNationality, CfgInitOwner, CfgInitPabxRate, CfgInitPackage, CfgInitPackageBreakdown, CfgInitPackageBusinessSource, CfgInitPaymentType,
	CfgInitPhoneBookType, CfgInitPrinter, CfgInitPurposeOf, CfgInitRegency, CfgInitReservationMark, CfgInitRoom, CfgInitRoomAllotmentType,
	CfgInitRoomAmenities, CfgInitRoomBoy, CfgInitRoomRate, CfgInitRoomRateBreakdown, CfgInitRoomRateBusinessSource, CfgInitRoomRateCategory,
	CfgInitRoomRateCompetitor, CfgInitRoomRateCurrency, CfgInitRoomRateDynamic, CfgInitRoomRateLastDeal, CfgInitRoomRateScale, CfgInitRoomRateSession,
	CfgInitRoomRateSubCategory, CfgInitRoomRateWeekly, CfgInitRoomType, CfgInitRoomUnavailableReason, CfgInitRoomView, CfgInitSales, CfgInitSalesSalary, CfgInitState, CfgInitSubDepartment,
	CfgInitTaxAndService, CfgInitTitle, CfgInitVoucherReason, CmNotification, CmUpdate, Company, Competitor, CompetitorData, Configuration,
	ConstAccountGroup, ConstBudgetType, ConstChannelManagerVendor, ConstChargeFrequency, ConstChargeType, ConstCommissionType, ConstCustomerDisplayVendor,
	ConstDepartmentType, ConstDynamicRateType, ConstFolioStatus, ConstFolioType, ConstForecastDay, ConstForecastMonth, ConstForeignCashTableId,
	ConstGuestStatus, ConstImage, ConstIptvVendor, ConstJournalAccountGroup, ConstJournalAccountSubGroupType, ConstJournalAccountType, ConstJournalPrefix,
	ConstKeylockVendor, ConstMemberType, ConstMikrotikVendor, ConstNotificationType, ConstOtherIcon, ConstOtpStatus, ConstPabxRateType, ConstPaymentGroup,
	ConstReportFont, ConstReportFormat, ConstReservationStatus, ConstRoomBlockStatus, ConstRoomStatus, ConstSmsDestinationType, ConstSmsRepeatType,
	ConstStatisticAccount, ConstSystem, ConstTransactionType, ConstVoucherStatus, ConstVoucherStatusApprove, ConstVoucherStatusSold, ConstVoucherType,
	ContactPerson, CorCfgInitUnit, CorReport, CorReportDefaultField, CorReportGroupField, CorReportGroupingField, CorReportOrderField, CorReportTemplate,
	CorReportTemplateField, CorUserGroup, CreditCard, DataAnalysis, DataAnalysisQuery, DataAnalysisQueryList, Events, FaCfgInitItem, FaCfgInitItemCategory,
	FaCfgInitLocation, FaCfgInitManufacture, FaConstDepreciationType, FaConstItemCondition, FaConstLocationType, FaDepreciation, FaList,
	FaLocationHistory, FaPurchaseOrder, FaPurchaseOrderDetail, FaReceive, FaReceiveDetail, FaRepair, FaRevaluation, FbStatistic, Folio, FolioList,
	FolioRouting, ForecastInHouseChangePax, ForecastMonthlyDay, ForecastMonthlyDayPrevious, GridProperties, GuestBreakdown, GuestDeposit, GuestDepositTax,
	GuestDetail, GuestExtraCharge, GuestExtraChargeBreakdown, GuestGeneral, GuestGroup, GuestInHouse, GuestInHouseBreakdown, GuestLoanItem, GuestMessage,
	GuestProfile, GuestScheduledRate, GuestToDo, HotelInformation, InvCfgInitItem, InvCfgInitItemCategory, InvCfgInitItemCategoryOtherCogs,
	InvCfgInitItemCategoryOtherCogs2, InvCfgInitItemCategoryOtherExpense, InvCfgInitItemGroup, InvCfgInitItemUom, InvCfgInitMarketList,
	InvCfgInitReturnStockReason, InvCfgInitStore, InvCfgInitUom, InvCloseLog, InvCloseSummary, InvCloseSummaryStore, InvCostRecipe, InvCosting,
	InvCostingDetail, InvOpname, InvProduction, InvPurchaseOrder, InvPurchaseOrderDetail, InvPurchaseRequest, InvPurchaseRequestDetail, InvReceiving,
	InvReceivingDetail, InvReturnStock, InvStockTransfer, InvStockTransferDetail, InvStoreRequisition, InvStoreRequisitionDetail, Invoice, InvoiceItem,
	InvoicePayment, Log, LogBackup, LogKeylock, LogMode, LogShift, LogSpecialAccess, LogUser, LogUserAction, LogUserActionGroup, LostAndFound,
	MarketStatistic, Member, MemberGift, MemberPoint, MemberPointRedeem, NotifTp, NotifTpCfgInitTemplate, NotifTpConstEvent, NotifTpConstVariable,
	NotifTpConstVendor, Notification, OneTimePassword, PabxSmdr, PhoneBook, PosCaptainOrder, PosCaptainOrderTransaction, PosCfgInitDiscountLimit,
	PosCfgInitMarket, PosCfgInitMemberOutletDiscount, PosCfgInitMemberOutletDiscountDetail, PosCfgInitMemberProductDiscount, PosCfgInitOutlet,
	PosCfgInitPaymentGroup, PosCfgInitProduct, PosCfgInitProductCategory, PosCfgInitProductGroup, PosCfgInitRoomBoy, PosCfgInitSpaRoom, PosCfgInitTable,
	PosCfgInitTableType, PosCfgInitTenan, PosCfgInitTherapistFingerprint, PosCfgInitWaitress, PosCheck, PosCheckTransaction, PosConstCheckType,
	PosConstComplimentType, PosConstDiscount, PosConstTimeSegment, PosInformation, PosIptvMenuOrder, PosMember, PosProductCosting, PosReport,
	PosReportDefaultField, PosReportGroupField, PosReportGroupingField, PosReportOrderField, PosReportTemplate, PosReportTemplateField, PosReservation,
	PosReservationTable, PosTableUnavailable, PosUserGroup, PosUserGroupOutlet, ProformaInvoiceDetail, Receipt, Report, ReportCustom, ReportCustomFavorite,
	ReportDefaultField, ReportGroupField, ReportGroupingField, ReportOrderField, ReportPivotTemp, ReportRoomRateStructureTemp, ReportRoomSales,
	ReportTemplate, ReportTemplateField, Reservation, ReservationExtraCharge, ReservationExtraChargeBreakdown, ReservationScheduledRate, RoomAllotment,
	RoomStatistic, RoomStatus, RoomUnavailable, RoomUnavailableHistory, SalActivity, SalActivityLog, SalCfgInitSegment, SalCfgInitSource,
	SalCfgInitTaskAction, SalCfgInitTaskRepeat, SalCfgInitTaskTag, SalCfgInitTemplate, SalConstProposalStatus, SalConstResource, SalConstStatus,
	SalConstTaskPriority, SalConstTaskStatus, SalContact, SalNotes, SalProposal, SalSendReminder, SalTask, SmsEvent, SmsOutbox, SmsSchedule, SubFolio,
	SubFolioGroup, SubFolioGrouping, SubFolioTax, TempSubFolioBreakdown1, TempSubFolioCorrectionBreakdown, TransactionList, TransactionListTax, User,
	UserGroup, UserGroupAccess, Voucher, WorkingShift, GeneralUserGroup, ReportUserGroup, ToolsUserGroup, ConstUserAccessLevel, CfgInitTimezone, CmUpdateAvailability, CmUpdateRate, BanBookingSchedulePayment string
}

var TableNameArray = []string{"acc_ap_ar",
	"acc_ap_ar_payment",
	"acc_ap_ar_payment_detail",
	"acc_ap_commission_payment",
	"acc_ap_commission_payment_detail",
	"acc_ap_refund_deposit_payment",
	"acc_ap_refund_deposit_payment_detail",
	"acc_cash_sale_recon",
	"acc_cfg_init_bank_account",
	"acc_close_month",
	"acc_close_year",
	"acc_const_bank_account_type",
	"acc_const_journal_group",
	"acc_const_journal_type",
	"acc_const_unit",
	"acc_credit_card_recon",
	"acc_credit_card_recon_detail",
	"acc_deffered_income",
	"acc_deffered_income_posted",
	"acc_foreign_cash",
	"acc_import_journal_log",
	"acc_journal",
	"acc_journal_detail",
	"acc_prepaid_expense",
	"acc_prepaid_expense_posted",
	"acc_report",
	"acc_report_default_field",
	"acc_report_grouping_field",
	"acc_report_group_field",
	"acc_report_order_field",
	"acc_report_template",
	"acc_report_template_field",
	"acc_user_group",
	"ast_cfg_init_shipping_address",
	"ast_const_purchase_request_status",
	"ast_const_store_requisition_status",
	"ast_report",
	"ast_report_default_field",
	"ast_report_grouping_field",
	"ast_report_group_field",
	"ast_report_order_field",
	"ast_report_template",
	"ast_report_template_field",
	"ast_user_group",
	"ast_user_sub_department",
	"audit_log",
	"ban_booking",
	"ban_cfg_init_seating_plan",
	"ban_cfg_init_theme",
	"ban_cfg_init_venue",
	"ban_cfg_init_venue_combine",
	"ban_cfg_init_venue_combine_detail",
	"ban_cfg_init_venue_group",
	"ban_const_booking_status",
	"ban_const_reservation_status",
	"ban_const_reservation_type",
	"ban_const_venue_location",
	"ban_report",
	"ban_report_default_field",
	"ban_report_grouping_field",
	"ban_report_group_field",
	"ban_report_order_field",
	"ban_report_template",
	"ban_report_template_field",
	"ban_reservation",
	"ban_reservation_charge",
	"ban_reservation_remark",
	"ban_user_group",
	"breakfast_list_temp",
	"budget_expense",
	"budget_fb",
	"budget_income",
	"budget_statistic",
	"cash_count",
	"cfg_init_account",
	"cfg_init_account_sub_group",
	"cfg_init_bed_type",
	"cfg_init_booking_source",
	"cfg_init_card_bank",
	"cfg_init_card_type",
	"cfg_init_city",
	"cfg_init_company_type",
	"cfg_init_competitor_category",
	"cfg_init_continent",
	"cfg_init_country",
	"cfg_init_credit_card_charge",
	"cfg_init_currency",
	"cfg_init_currency_nominal",
	"cfg_init_custom_lookup_field01",
	"cfg_init_custom_lookup_field02",
	"cfg_init_custom_lookup_field03",
	"cfg_init_custom_lookup_field04",
	"cfg_init_custom_lookup_field05",
	"cfg_init_custom_lookup_field06",
	"cfg_init_custom_lookup_field07",
	"cfg_init_custom_lookup_field08",
	"cfg_init_custom_lookup_field09",
	"cfg_init_custom_lookup_field10",
	"cfg_init_custom_lookup_field11",
	"cfg_init_custom_lookup_field12",
	"cfg_init_department",
	"cfg_init_guest_type",
	"cfg_init_id_card_type",
	"cfg_init_is_fb_sub_department_group",
	"cfg_init_is_fb_sub_department_group_detail",
	"cfg_init_journal_account",
	"cfg_init_journal_account_category",
	"cfg_init_journal_account_sub_group",
	"cfg_init_language",
	"cfg_init_loan_item",
	"cfg_init_market",
	"cfg_init_market_category",
	"cfg_init_member_point_type",
	"cfg_init_nationality",
	"cfg_init_owner",
	"cfg_init_pabx_rate",
	"cfg_init_package",
	"cfg_init_package_breakdown",
	"cfg_init_package_business_source",
	"cfg_init_payment_type",
	"cfg_init_phone_book_type",
	"cfg_init_printer",
	"cfg_init_purpose_of",
	"cfg_init_regency",
	"cfg_init_reservation_mark",
	"cfg_init_room",
	"cfg_init_room_allotment_type",
	"cfg_init_room_amenities",
	"cfg_init_room_boy",
	"cfg_init_room_rate",
	"cfg_init_room_rate_breakdown",
	"cfg_init_room_rate_business_source",
	"cfg_init_room_rate_category",
	"cfg_init_room_rate_competitor",
	"cfg_init_room_rate_currency",
	"cfg_init_room_rate_dynamic",
	"cfg_init_room_rate_last_deal",
	"cfg_init_room_rate_scale",
	"cfg_init_room_rate_session",
	"cfg_init_room_rate_sub_category",
	"cfg_init_room_rate_weekly",
	"cfg_init_room_type",
	"cfg_init_room_unavailable_reason",
	"cfg_init_room_view",
	"cfg_init_sales",
	"cfg_init_state",
	"cfg_init_sub_department",
	"cfg_init_tax_and_service",
	"cfg_init_title",
	"cfg_init_voucher_reason",
	"cm_notification",
	"cm_update",
	"company",
	"competitor",
	"competitor_data",
	"configuration",
	"const_account_group",
	"const_budget_type",
	"const_channel_manager_vendor",
	"const_charge_frequency",
	"const_charge_type",
	"const_commission_type",
	"const_customer_display_vendor",
	"const_department_type",
	"const_dynamic_rate_type",
	"const_folio_status",
	"const_folio_type",
	"const_forecast_day",
	"const_forecast_month",
	"const_foreign_cash_table_id",
	"const_guest_status",
	"const_image",
	"const_iptv_vendor",
	"const_journal_account_group",
	"const_journal_account_sub_group_type",
	"const_journal_account_type",
	"const_journal_prefix",
	"const_keylock_vendor",
	"const_member_type",
	"const_mikrotik_vendor",
	"const_notification_type",
	"const_other_icon",
	"const_otp_status",
	"const_pabx_rate_type",
	"const_payment_group",
	"const_report_font",
	"const_report_format",
	"const_reservation_status",
	"const_room_block_status",
	"const_room_status",
	"const_sms_destination_type",
	"const_sms_repeat_type",
	"const_statistic_account",
	"const_system",
	"const_transaction_type",
	"const_user_access_level",
	"const_voucher_status",
	"const_voucher_status_approve",
	"const_voucher_status_sold",
	"const_voucher_type",
	"contact_person",
	"cor_cfg_init_unit",
	"cor_report",
	"cor_report_default_field",
	"cor_report_grouping_field",
	"cor_report_group_field",
	"cor_report_order_field",
	"cor_report_template",
	"cor_report_template_field",
	"cor_user_group",
	"credit_card",
	"data_analysis",
	"data_analysis_query",
	"data_analysis_query_list",
	"events",
	"fa_cfg_init_item",
	"fa_cfg_init_item_category",
	"fa_cfg_init_location",
	"fa_cfg_init_manufacture",
	"fa_const_depreciation_type",
	"fa_const_item_condition",
	"fa_const_location_type",
	"fa_depreciation",
	"fa_list",
	"fa_location_history",
	"fa_purchase_order",
	"fa_purchase_order_detail",
	"fa_receive",
	"fa_receive_detail",
	"fa_repair",
	"fa_revaluation",
	"fb_statistic",
	"folio",
	"folio_routing",
	"forecast_in_house_change_pax",
	"forecast_monthly_day",
	"forecast_monthly_day_previous",
	"general_user_group",
	"grid_properties",
	"guest_breakdown",
	"guest_deposit",
	"guest_detail",
	"guest_extra_charge",
	"guest_extra_charge_breakdown",
	"guest_general",
	"guest_group",
	"guest_in_house",
	"guest_in_house_breakdown",
	"guest_loan_item",
	"guest_message",
	"guest_profile",
	"guest_registration",
	"guest_scheduled_rate",
	"guest_to_do",
	"hotel_information",
	"invoice",
	"invoice_item",
	"invoice_payment",
	"inv_cfg_init_item",
	"inv_cfg_init_item_category",
	"inv_cfg_init_item_category_other_cogs",
	"inv_cfg_init_item_category_other_cogs2",
	"inv_cfg_init_item_category_other_expense",
	"inv_cfg_init_item_group",
	"inv_cfg_init_item_uom",
	"inv_cfg_init_market_list",
	"inv_cfg_init_return_stock_reason",
	"inv_cfg_init_store",
	"inv_cfg_init_uom",
	"inv_close_log",
	"inv_close_summary",
	"inv_close_summary_store",
	"inv_costing",
	"inv_costing_detail",
	"inv_cost_recipe",
	"inv_opname",
	"inv_production",
	"inv_purchase_order",
	"inv_purchase_order_detail",
	"inv_purchase_request",
	"inv_purchase_request_detail",
	"inv_receiving",
	"inv_receiving_detail",
	"inv_return_stock",
	"inv_stock_transfer",
	"inv_stock_transfer_detail",
	"inv_store_requisition",
	"inv_store_requisition_detail",
	"log",
	"log_backup",
	"log_keylock",
	"log_mode",
	"log_shift",
	"log_special_access",
	"log_user",
	"log_user_action",
	"log_user_action_group",
	"lost_and_found",
	"market_statistic",
	"member",
	"member_gift",
	"member_point",
	"member_point_redeem",
	"notification",
	"notif_tp",
	"notif_tp_cfg_init_template",
	"notif_tp_const_event",
	"notif_tp_const_variable",
	"notif_tp_const_vendor",
	"one_time_password",
	"pabx_smdr",
	"phone_book",
	"pos_captain_order",
	"pos_captain_order_transaction",
	"pos_cfg_init_discount_limit",
	"pos_cfg_init_market",
	"pos_cfg_init_member_outlet_discount",
	"pos_cfg_init_member_outlet_discount_detail",
	"pos_cfg_init_member_product_discount",
	"pos_cfg_init_outlet",
	"pos_cfg_init_payment_group",
	"pos_cfg_init_product",
	"pos_cfg_init_product_category",
	"pos_cfg_init_product_group",
	"pos_cfg_init_room_boy",
	"pos_cfg_init_spa_room",
	"pos_cfg_init_table",
	"pos_cfg_init_table_type",
	"pos_cfg_init_tenan",
	"pos_cfg_init_therapist_fingerprint",
	"pos_cfg_init_waitress",
	"pos_check",
	"pos_check_transaction",
	"pos_const_check_type",
	"pos_const_compliment_type",
	"pos_const_discount",
	"pos_const_time_segment",
	"pos_information",
	"pos_iptv_menu_order",
	"pos_member",
	"pos_product_costing",
	"pos_report",
	"pos_report_default_field",
	"pos_report_grouping_field",
	"pos_report_group_field",
	"pos_report_order_field",
	"pos_report_template",
	"pos_report_template_field",
	"pos_reservation",
	"pos_reservation_table",
	"pos_table_unavailable",
	"pos_user_group",
	"pos_user_group_outlet",
	"proforma_invoice_detail",
	"receipt",
	"report",
	"report_custom",
	"report_custom_favorite",
	"report_default_field",
	"report_grouping_field",
	"report_group_field",
	"report_order_field",
	"report_pivot_temp",
	"report_room_rate_structure_temp",
	"report_room_sales",
	"report_template",
	"report_template_field",
	"report_user_group",
	"reservation",
	"reservation_extra_charge",
	"reservation_extra_charge_breakdown",
	"reservation_scheduled_rate",
	"room_allotment",
	"room_statistic",
	"room_status",
	"room_unavailable",
	"room_unavailable_history",
	"sal_activity",
	"sal_activity_log",
	"sal_cfg_init_segment",
	"sal_cfg_init_source",
	"sal_cfg_init_task_action",
	"sal_cfg_init_task_repeat",
	"sal_cfg_init_task_tag",
	"sal_cfg_init_template",
	"sal_const_proposal_status",
	"sal_const_resource",
	"sal_const_status",
	"sal_const_task_priority",
	"sal_const_task_status",
	"sal_contact",
	"sal_notes",
	"sal_proposal",
	"sal_send_reminder",
	"sal_task",
	"sms_event",
	"sms_outbox",
	"sms_schedule",
	"sub_folio",
	"sub_folio_group",
	"temp_sub_folio_breakdown1",
	"temp_sub_folio_correction_breakdown",
	"tools_user_group",
	"user",
	"user_group",
	"user_group_access",
	"voucher",
	"working_shift",
}

type Date string

func (d *Date) UnmarshalJSON(bytes []byte) error {
	dd, err := time.Parse(`"2006-01-02T15:04:05Z"`, string(bytes))
	if err != nil {
		de, err := time.Parse(`"2006-01-02T15:04:05+08:00"`, string(bytes))
		if err != nil {
			dt, err := time.Parse(`"2006-01-02"`, string(bytes))
			if err != nil {
				return err
			}
			*d = Date(dt.Format("2006-01-02"))
			return nil
		}
		*d = Date(de.Format("2006-01-02"))
		return nil
	}
	*d = Date(dd.Format("2006-01-02"))

	return nil
}

type DateTime string

func (d *DateTime) UnmarshalJSON(bytes []byte) error {
	dd, err := time.Parse(`"2006-01-02T15:04:05Z"`, string(bytes))
	if err != nil {
		de, err := time.Parse(`"2006-01-02T15:04:05+08:00"`, string(bytes))
		if err != nil {
			return err
		}
		*d = DateTime(de.Format("2006-01-02 15:04:05"))
		return nil
	}
	*d = DateTime(dd.Format("2006-01-02 15:04:05"))

	return nil
}

type Acc_ap_ar struct {
	Number               string    `json:"number" binding:"required" gorm:"primaryKey"`
	DocumentNumber       string    `json:"document_number"`
	RefNumber            string    `json:"ref_number" binding:"required"`
	CompanyCode          string    `json:"company_code"`
	JournalAccountDebit  string    `json:"journal_account_debit" binding:"required"`
	JournalAccountCredit string    `json:"journal_account_credit" binding:"required"`
	Amount               float64   `json:"amount" binding:"required"`
	AmountPaid           float64   `json:"amount_paid" binding:"required"`
	Date                 time.Time `json:"date" binding:"required"`
	DueDate              time.Time `json:"due_date" binding:"required"`
	Remark               string    `json:"remark"`
	IsAp                 uint8     `json:"is_ap" binding:"required"`
	IsAccrued            uint8     `json:"is_accrued" binding:"required"`
	IsAuto               uint8     `json:"is_auto" binding:"required"`
	IsPaid               uint8     `json:"is_paid" binding:"required"`
	CreatedAt            time.Time `json:"created_at"`
	CreatedBy            string    `json:"created_by"`
	UpdatedAt            time.Time `json:"updated_at"`
	UpdatedBy            string    `json:"updated_by"`
	Id                   uint64    `json:"id"`
}

type Acc_ap_ar_payment struct {
	RefNumber                  string    `json:"ref_number" binding:"required"`
	JournalAccountCode         string    `json:"journal_account_code" binding:"required"`
	ApJournalAccountCode       string    `json:"ap_journal_account_code"`
	CreateApNumber             string    `json:"create_ap_number"`
	DiscountJournalAccountCode string    `json:"discount_journal_account_code"`
	BaJournalAccountCode       string    `json:"ba_journal_account_code"`
	OeJournalAccountCode       string    `json:"oe_journal_account_code"`
	TotalAmount                float64   `json:"total_amount" binding:"required"`
	Discount                   float64   `json:"discount"`
	BankAdministration         float64   `json:"bank_administration"`
	OtherExpense               float64   `json:"other_expense"`
	Date                       time.Time `json:"date" binding:"required"`
	Remark                     string    `json:"remark"`
	SourceCodeApAr             string    `json:"source_code_ap_ar" binding:"required"`
	IsPaymentApAr              *uint8    `json:"is_payment_ap_ar" binding:"required"`
	CreatedAt                  time.Time `json:"created_at"`
	CreatedBy                  string    `json:"created_by"`
	UpdatedAt                  time.Time `json:"updated_at"`
	UpdatedBy                  string    `json:"updated_by"`
	Id                         uint64    `json:"id"`
}

type Acc_ap_ar_payment_detail struct {
	ApArNumber string    `json:"ap_ar_number" binding:"required"`
	RefNumber  string    `json:"ref_number" binding:"required"`
	Amount     float64   `json:"amount" binding:"required"`
	Remark     string    `json:"remark"`
	CreatedAt  time.Time `json:"created_at"`
	CreatedBy  string    `json:"created_by"`
	UpdatedAt  time.Time `json:"updated_at"`
	UpdatedBy  string    `json:"updated_by"`
	Id         uint64    `json:"id"`
}

type Acc_ap_commission_payment struct {
	RefNumber                  string    `json:"ref_number" binding:"required"`
	JournalAccountCode         string    `json:"journal_account_code" binding:"required"`
	DiscountJournalAccountCode string    `json:"discount_journal_account_code"`
	BaJournalAccountCode       string    `json:"ba_journal_account_code"`
	OeJournalAccountCode       string    `json:"oe_journal_account_code"`
	TotalAmount                float64   `json:"total_amount" binding:"required"`
	Discount                   float64   `json:"discount"`
	BankAdministration         float64   `json:"bank_administration"`
	OtherExpense               float64   `json:"other_expense"`
	Date                       time.Time `json:"date" binding:"required"`
	Remark                     string    `json:"remark"`
	SourceCodeApAr             string    `json:"source_code_ap_ar" binding:"required"`
	CreatedAt                  time.Time `json:"created_at"`
	CreatedBy                  string    `json:"created_by"`
	UpdatedAt                  time.Time `json:"updated_at"`
	UpdatedBy                  string    `json:"updated_by"`
	Id                         uint64    `json:"id"`
}

type Acc_ap_commission_payment_detail struct {
	SubFolioId uint64    `json:"sub_folio_id" binding:"required"`
	RefNumber  string    `json:"ref_number" binding:"required"`
	Amount     float64   `json:"amount" binding:"required"`
	CreatedAt  time.Time `json:"created_at"`
	CreatedBy  string    `json:"created_by"`
	UpdatedAt  time.Time `json:"updated_at"`
	UpdatedBy  string    `json:"updated_by"`
	Id         uint64    `json:"id"`
}

type Acc_ap_refund_deposit_payment struct {
	RefNumber                  string    `json:"ref_number" binding:"required"`
	JournalAccountCode         string    `json:"journal_account_code" binding:"required"`
	DiscountJournalAccountCode string    `json:"discount_journal_account_code"`
	BaJournalAccountCode       string    `json:"ba_journal_account_code"`
	OeJournalAccountCode       string    `json:"oe_journal_account_code"`
	TotalAmount                float64   `json:"total_amount" binding:"required"`
	Discount                   float64   `json:"discount"`
	BankAdministration         float64   `json:"bank_administration"`
	OtherExpense               float64   `json:"other_expense"`
	Date                       time.Time `json:"date" binding:"required"`
	Remark                     string    `json:"remark"`
	PaidByApAr                 string    `json:"paid_by_ap_ar"`
	CreatedAt                  time.Time `json:"created_at"`
	CreatedBy                  string    `json:"created_by"`
	UpdatedAt                  time.Time `json:"updated_at"`
	UpdatedBy                  string    `json:"updated_by"`
	Id                         uint64    `json:"id"`
}

type Acc_ap_refund_deposit_payment_detail struct {
	SubFolioId uint64    `json:"sub_folio_id" binding:"required"`
	RefNumber  string    `json:"ref_number" binding:"required"`
	Amount     float64   `json:"amount" binding:"required"`
	CreatedAt  time.Time `json:"created_at"`
	CreatedBy  string    `json:"created_by"`
	UpdatedAt  time.Time `json:"updated_at"`
	UpdatedBy  string    `json:"updated_by"`
	Id         uint64    `json:"id"`
}

type Acc_cash_sale_recon struct {
	JournalAccountCode          string    `json:"journal_account_code" binding:"required"`
	JournalAccountCodeShortOver string    `json:"journal_account_code_short_over"`
	RefNumber                   string    `json:"ref_number" binding:"required"`
	Date                        time.Time `json:"date" binding:"required"`
	DateRecon                   time.Time `json:"date_recon"`
	Amount                      float64   `json:"amount" binding:"required"`
	AmountShortOver             float64   `json:"amount_short_over"`
	AmountDetail                float64   `json:"amount_detail"`
	Remark                      string    `json:"remark"`
	ReconBy                     string    `json:"recon_by" binding:"required"`
	IsOver                      *uint8    `json:"is_over"`
	CreatedAt                   time.Time `json:"created_at"`
	CreatedBy                   string    `json:"created_by"`
	UpdatedAt                   time.Time `json:"updated_at"`
	UpdatedBy                   string    `json:"updated_by"`
	Id                          uint64    `json:"id" gorm:"primaryKey"`
}

type Acc_cfg_init_bank_account struct {
	Code               string    `json:"code" binding:"required"`
	Name               string    `json:"name" binding:"required"`
	JournalAccountCode string    `json:"journal_account_code" binding:"required"`
	TypeCode           string    `json:"type_code" binding:"required"`
	BankName           string    `json:"bank_name"`
	BankAccountNumber  string    `json:"bank_account_number"`
	BankAddress        string    `json:"bank_address"`
	ForReceive         uint8     `json:"for_receive"`
	ForPayment         uint8     `json:"for_payment"`
	CreatedAt          time.Time `json:"created_at"`
	CreatedBy          string    `json:"created_by"`
	UpdatedAt          time.Time `json:"updated_at"`
	UpdatedBy          string    `json:"updated_by"`
	Id                 uint64    `json:"id"`
}

type Acc_close_month struct {
	Month     uint8     `json:"month" binding:"required"`
	Year      uint64    `json:"year" binding:"required"`
	CloseTime time.Time `json:"close_time" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Acc_close_year struct {
	Year      uint64    `json:"year" binding:"required"`
	CloseTime time.Time `json:"close_time" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Acc_const_bank_account_type struct {
	Code   string `json:"code" binding:"required"`
	Name   string `json:"name" binding:"required"`
	IdSort int    `json:"id_sort" binding:"required"`
}

type Acc_const_journal_group struct {
	Code   string `json:"code" binding:"required"`
	Name   string `json:"name" binding:"required"`
	IdSort int    `json:"id_sort" binding:"required"`
}

type Acc_const_journal_type struct {
	Code string `json:"code" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type Acc_const_unit struct {
	Code   string `json:"code" binding:"required"`
	Name   string `json:"name"`
	IdSort uint8  `json:"id_sort"`
}

type Acc_credit_card_recon struct {
	JournalAccountCode string    `json:"journal_account_code" binding:"required"`
	RefNumber          string    `json:"ref_number" binding:"required"`
	Date               time.Time `json:"date" binding:"required"`
	AmountReceived     float64   `json:"amount_received" binding:"required"`
	CreatedAt          time.Time `json:"created_at"`
	CreatedBy          string    `json:"created_by"`
	UpdatedAt          time.Time `json:"updated_at"`
	UpdatedBy          string    `json:"updated_by"`
	Id                 uint64    `json:"id"`
}

type Acc_credit_card_recon_detail struct {
	AccCreditCardReconId uint64    `json:"acc_credit_card_recon_id" binding:"required"`
	GuestDepositId       uint64    `json:"guest_deposit_id" binding:"required"`
	SubFolioId           uint64    `json:"sub_folio_id" binding:"required"`
	Amount               float64   `json:"amount" binding:"required"`
	Remark               string    `json:"remark" binding:"required"`
	CreatedAt            time.Time `json:"created_at"`
	CreatedBy            string    `json:"created_by"`
	UpdatedAt            time.Time `json:"updated_at"`
	UpdatedBy            string    `json:"updated_by"`
	Id                   uint64    `json:"id" gorm:"primaryKey"`
}

type Acc_deffered_income struct {
	Date                    time.Time `json:"date" binding:"required"`
	RefNumber               string    `json:"ref_number" binding:"required"`
	Description             string    `json:"description" binding:"required"`
	Amount                  float64   `json:"amount" binding:"required"`
	CompanyCode             string    `json:"company_code" binding:"required"`
	DefferedAccountCode     string    `json:"deffered_account_code" binding:"required"`
	AmountPayment           float64   `json:"amount_payment" binding:"required"`
	BankAccountCode         string    `json:"bank_account_code" binding:"required"`
	IsCreateAr              *uint8    `json:"is_create_ar"`
	ApArNumber              string    `json:"ap_ar_number"`
	SubDepartmentIncomeCode string    `json:"sub_department_income_code" binding:"required"`
	IncomeAccountCode       string    `json:"income_account_code" binding:"required"`
	Month                   int       `json:"month" binding:"required"`
	IsNextMonth             *uint8    `json:"is_next_month" binding:"required"`
	Remark                  string    `json:"remark"`
	CreatedAt               time.Time `json:"created_at"`
	CreatedBy               string    `json:"created_by"`
	UpdatedAt               time.Time `json:"updated_at"`
	UpdatedBy               string    `json:"updated_by"`
	Id                      uint64    `json:"id"`
}

type Acc_deffered_income_posted struct {
	DefferedId        uint64    `json:"deffered_id" binding:"required"`
	RefNumber         string    `json:"ref_number" binding:"required"`
	PostingDate       time.Time `json:"posting_date" binding:"required"`
	Amount            float64   `json:"amount" binding:"required"`
	SubDepartmentCode string    `json:"sub_department_code" binding:"required"`
	IncomeAccountCode string    `json:"income_account_code" binding:"required"`
	Remark            string    `json:"remark"`
	CreatedAt         time.Time `json:"created_at"`
	CreatedBy         string    `json:"created_by"`
	UpdatedAt         time.Time `json:"updated_at"`
	UpdatedBy         string    `json:"updated_by"`
	Id                uint64    `json:"id"`
}

type Acc_foreign_cash struct {
	IdTransaction       uint64    `json:"id_transaction" binding:"required"`
	IdCorrected         uint64    `json:"id_corrected" binding:"required"`
	IdChange            uint64    `json:"id_change" binding:"required"`
	IdTable             int       `json:"id_table" binding:"required"`
	Breakdown           uint64    `json:"breakdown" binding:"required"`
	RefNumber           string    `json:"ref_number" binding:"required"`
	Date                time.Time `json:"date" binding:"required"`
	TypeCode            string    `json:"type_code" binding:"required"`
	Amount              float64   `json:"amount" binding:"required"`
	Stock               float64   `json:"stock" binding:"required"`
	DefaultCurrencyCode string    `json:"default_currency_code" binding:"required"`
	AmountForeign       float64   `json:"amount_foreign" binding:"required"`
	ExchangeRate        float64   `json:"exchange_rate" binding:"required"`
	CurrencyCode        string    `json:"currency_code" binding:"required"`
	Remark              string    `json:"remark"`
	IsCorrection        uint8     `json:"is_correction" binding:"required"`
	CreatedAt           time.Time `json:"created_at"`
	CreatedBy           string    `json:"created_by"`
	UpdatedAt           time.Time `json:"updated_at"`
	UpdatedBy           string    `json:"updated_by"`
	Id                  uint64    `json:"id" gorm:"primaryKey"`
}

type Acc_import_journal_log struct {
	RefNumber   string    `json:"ref_number" binding:"required"`
	AuditDate   time.Time `json:"audit_date" binding:"required"`
	PostingDate time.Time `json:"posting_date" binding:"required"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   string    `json:"created_by"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedBy   string    `json:"updated_by"`
	Id          uint64    `json:"id"`
}

type Acc_journal struct {
	RefNumber      string    `json:"ref_number" gorm:"primaryKey" binding:"required"`
	UnitCode       string    `json:"unit_code" binding:"required"`
	ApArNumber     string    `json:"ap_ar_number"`
	CompanyCode    string    `json:"company_code"`
	TypeCode       string    `json:"type_code" binding:"required"`
	GroupCode      string    `json:"group_code" binding:"required"`
	Date           time.Time `json:"date" binding:"required"`
	DateUnixx      int64     `json:"date_unixx" binding:"required"`
	DocumentNumber string    `json:"document_number"`
	Memo           string    `json:"memo" binding:"required"`
	ChequeNumber   string    `json:"cheque_number"`
	IdSort         int       `json:"id_sort"`
	CreatedAt      time.Time `json:"created_at"`
	CreatedBy      string    `json:"created_by"`
	UpdatedAt      time.Time `json:"updated_at"`
	UpdatedBy      string    `json:"updated_by"`
	Id             uint64    `json:"id"`
	IdHolding      uint64    `json:"id_holding"`
}

type Acc_journal_detail struct {
	RefNumber         string    `json:"ref_number" binding:"required"`
	Date              time.Time `json:"date" binding:"required"`
	UnitCode          string    `json:"unit_code"`
	SubDepartmentCode string    `json:"sub_department_code" binding:"required"`
	AccountCode       string    `json:"account_code" binding:"required"`
	Amount            float64   `json:"amount" binding:"required"`
	TypeCode          string    `json:"type_code" binding:"required"`
	Remark            string    `json:"remark"`
	IdData            string    `json:"id_data" binding:"required"`
	CreatedAt         time.Time `json:"created_at"`
	CreatedBy         string    `json:"created_by"`
	UpdatedAt         time.Time `json:"updated_at"`
	UpdatedBy         string    `json:"updated_by"`
	Id                uint64    `json:"id"`
	IdHolding         uint64    `json:"id_holding"`
}

type Acc_prepaid_expense struct {
	Date                     time.Time `json:"date" binding:"required"`
	RefNumber                string    `json:"ref_number" binding:"required"`
	Description              string    `json:"description" binding:"required"`
	Amount                   float64   `json:"amount" binding:"required"`
	CompanyCode              string    `json:"company_code" binding:"required"`
	PrepaidAccountCode       string    `json:"prepaid_account_code" binding:"required"`
	AmountPayment            float64   `json:"amount_payment" binding:"required"`
	BankAccountCode          string    `json:"bank_account_code" binding:"required"`
	IsCreateAp               *uint8    `json:"is_create_ap"`
	ApArNumber               string    `json:"ap_ar_number"`
	SubDepartmentExpenseCode string    `json:"sub_department_expense_code" binding:"required"`
	ExpenseAccountCode       string    `json:"expense_account_code" binding:"required"`
	Month                    int       `json:"month" binding:"required"`
	IsNextMonth              *uint8    `json:"is_next_month" binding:"required"`
	Remark                   string    `json:"remark"`
	CreatedAt                time.Time `json:"created_at"`
	CreatedBy                string    `json:"created_by"`
	UpdatedAt                time.Time `json:"updated_at"`
	UpdatedBy                string    `json:"updated_by"`
	Id                       uint64    `json:"id" gorm:"primaryKey"`
}

type Acc_prepaid_expense_posted struct {
	PrepaidId          uint64    `json:"prepaid_id" binding:"required"`
	RefNumber          string    `json:"ref_number" binding:"required"`
	PostingDate        time.Time `json:"posting_date" binding:"required"`
	Amount             float64   `json:"amount" binding:"required"`
	SubDepartmentCode  string    `json:"sub_department_code" binding:"required"`
	ExpenseAccountCode string    `json:"expense_account_code" binding:"required"`
	Remark             string    `json:"remark"`
	CreatedAt          time.Time `json:"created_at"`
	CreatedBy          string    `json:"created_by"`
	UpdatedAt          time.Time `json:"updated_at"`
	UpdatedBy          string    `json:"updated_by"`
	Id                 uint64    `json:"id"`
}

type Acc_report struct {
	Code        int       `json:"code" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	ReportQuery string    `json:"report_query" binding:"required"`
	ParentId    uint64    `json:"parent_id" binding:"required"`
	IsSystem    uint8     `json:"is_system" binding:"required"`
	IdSort      int       `json:"id_sort" binding:"required"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   string    `json:"created_by"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedBy   string    `json:"updated_by"`
	Id          uint64    `json:"id"`
}

type Acc_report_default_field struct {
	ReportCode int    `json:"report_code" binding:"required"`
	FieldQuery string `json:"field_query" binding:"required"`
	FieldName  string `json:"field_name" binding:"required"`
	IdSort     int    `json:"id_sort" binding:"required"`
}

type Acc_report_group_field struct {
	TemplateId uint64 `json:"template_id" binding:"required"`
	FieldName  string `json:"field_name" binding:"required"`
	IdSort     int    `json:"id_sort" binding:"required"`
}

type Acc_report_grouping_field struct {
	TemplateId uint64 `json:"template_id" binding:"required"`
	FieldName  string `json:"field_name" binding:"required"`
	IdSort     int    `json:"id_sort" binding:"required"`
}

type Acc_report_order_field struct {
	TemplateId  uint64 `json:"template_id" binding:"required"`
	FieldName   string `json:"field_name" binding:"required"`
	IsAscending uint8  `json:"is_ascending" binding:"required"`
	IdSort      int    `json:"id_sort" binding:"required"`
}

type Acc_report_template struct {
	ReportCode       int       `json:"report_code" binding:"required"`
	Name             string    `json:"name" binding:"required"`
	GroupLevel       int       `json:"group_level" binding:"required"`
	HeaderRemark     string    `json:"header_remark"`
	ShowFooter       uint8     `json:"show_footer" binding:"required"`
	ShowPageNumber   string    `json:"show_page_number" binding:"required"`
	PaperSize        int       `json:"paper_size" binding:"required"`
	PaperWidth       float64   `json:"paper_width" binding:"required"`
	PaperHeight      float64   `json:"paper_height" binding:"required"`
	IsPortrait       uint8     `json:"is_portrait" binding:"required"`
	HeaderRowHeight  int       `json:"header_row_height" binding:"required"`
	RowHeight        int       `json:"row_height" binding:"required"`
	HorizontalBorder uint8     `json:"horizontal_border" binding:"required"`
	VerticalBorder   uint8     `json:"vertical_border" binding:"required"`
	SignName1        string    `json:"sign_name1"`
	SignPosition1    string    `json:"sign_position1"`
	SignName2        string    `json:"sign_name2"`
	SignPosition2    string    `json:"sign_position2"`
	SignName3        string    `json:"sign_name3"`
	SignPosition3    string    `json:"sign_position3"`
	SignName4        string    `json:"sign_name4"`
	SignPosition4    string    `json:"sign_position4"`
	IsDefault        uint8     `json:"is_default" binding:"required"`
	IsSystem         uint8     `json:"is_system" binding:"required"`
	IdSort           int       `json:"id_sort" binding:"required"`
	CreatedAt        time.Time `json:"created_at"`
	CreatedBy        string    `json:"created_by"`
	UpdatedAt        time.Time `json:"updated_at"`
	UpdatedBy        string    `json:"updated_by"`
	Id               uint64    `json:"id"`
}

type Acc_report_template_field struct {
	TemplateId      uint64 `json:"template_id" binding:"required"`
	FieldName       string `json:"field_name" binding:"required"`
	HeaderName      string `json:"header_name" binding:"required"`
	FooterType      int    `json:"footer_type" binding:"required"`
	FormatCode      int    `json:"format_code" binding:"required"`
	Width           int    `json:"width" binding:"required"`
	Font            int    `json:"font" binding:"required"`
	FontSize        int    `json:"font_size" binding:"required"`
	FontColor       int    `json:"font_color" binding:"required"`
	Alignment       string `json:"alignment" binding:"required"`
	HeaderFontSize  int    `json:"header_font_size" binding:"required"`
	HeaderFontColor int    `json:"header_font_color" binding:"required"`
	HeaderAlignment string `json:"header_alignment" binding:"required"`
	IdSort          int    `json:"id_sort" binding:"required"`
}
type Acc_user_group struct {
	Code              string    `json:"code" binding:"required"`
	AccessForm        string    `json:"access_form" binding:"required"`
	AccessSpecial     string    `json:"access_special" binding:"required"`
	AccessInvoice     string    `json:"access_invoice" binding:"required"`
	PrintInvoiceCount int       `json:"print_invoice_count" binding:"required"`
	CreatedAt         time.Time `json:"created_at"`
	CreatedBy         string    `json:"created_by"`
	UpdatedAt         time.Time `json:"updated_at"`
	UpdatedBy         string    `json:"updated_by"`
	Id                uint64    `json:"id" gorm:"primaryKey"`
}

type Ast_cfg_init_shipping_address struct {
	Code          string    `json:"code" binding:"required"`
	Name          string    `json:"name" binding:"required"`
	ContactPerson string    `json:"contact_person"`
	Street        string    `json:"street" binding:"required"`
	City          string    `json:"city" binding:"required"`
	CountryCode   string    `json:"country_code"`
	StateCode     string    `json:"state_code"`
	PostalCode    string    `json:"postal_code"`
	Phone1        string    `json:"phone1"`
	Phone2        string    `json:"phone2"`
	Fax           string    `json:"fax"`
	Email         string    `json:"email"`
	Website       string    `json:"website"`
	IsActive      uint8     `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	CreatedBy     string    `json:"created_by"`
	UpdatedAt     time.Time `json:"updated_at"`
	UpdatedBy     string    `json:"updated_by"`
	Id            uint64    `json:"id"`
}

type Ast_const_purchase_request_status struct {
	Code string `json:"code" binding:"required"`
	Name string `json:"name"`
}

type Ast_const_store_requisition_status struct {
	Code string `json:"code" binding:"required"`
	Name string `json:"name"`
}

type Ast_report struct {
	Code        int       `json:"code" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	ReportQuery string    `json:"report_query" binding:"required"`
	ParentId    uint64    `json:"parent_id" binding:"required"`
	IsSystem    uint8     `json:"is_system" binding:"required"`
	IdSort      int       `json:"id_sort" binding:"required"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   string    `json:"created_by"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedBy   string    `json:"updated_by"`
	Id          uint64    `json:"id"`
}

type Ast_report_default_field struct {
	ReportCode int    `json:"report_code" binding:"required"`
	FieldQuery string `json:"field_query" binding:"required"`
	FieldName  string `json:"field_name" binding:"required"`
	IdSort     int    `json:"id_sort" binding:"required"`
}

type Ast_report_group_field struct {
	TemplateId uint64 `json:"template_id" binding:"required"`
	FieldName  string `json:"field_name" binding:"required"`
	IdSort     int    `json:"id_sort" binding:"required"`
}

type Ast_report_grouping_field struct {
	TemplateId uint64 `json:"template_id" binding:"required"`
	FieldName  string `json:"field_name" binding:"required"`
	IdSort     int    `json:"id_sort" binding:"required"`
}

type Ast_report_order_field struct {
	TemplateId  uint64 `json:"template_id" binding:"required"`
	FieldName   string `json:"field_name" binding:"required"`
	IsAscending uint8  `json:"is_ascending" binding:"required"`
	IdSort      int    `json:"id_sort" binding:"required"`
}

type Ast_report_template struct {
	ReportCode       int       `json:"report_code" binding:"required"`
	Name             string    `json:"name" binding:"required"`
	GroupLevel       int       `json:"group_level" binding:"required"`
	HeaderRemark     string    `json:"header_remark"`
	ShowFooter       uint8     `json:"show_footer" binding:"required"`
	ShowPageNumber   string    `json:"show_page_number" binding:"required"`
	PaperSize        int       `json:"paper_size" binding:"required"`
	PaperWidth       float64   `json:"paper_width" binding:"required"`
	PaperHeight      float64   `json:"paper_height" binding:"required"`
	IsPortrait       uint8     `json:"is_portrait" binding:"required"`
	HeaderRowHeight  int       `json:"header_row_height" binding:"required"`
	RowHeight        int       `json:"row_height" binding:"required"`
	HorizontalBorder uint8     `json:"horizontal_border" binding:"required"`
	VerticalBorder   uint8     `json:"vertical_border" binding:"required"`
	SignName1        string    `json:"sign_name1"`
	SignPosition1    string    `json:"sign_position1"`
	SignName2        string    `json:"sign_name2"`
	SignPosition2    string    `json:"sign_position2"`
	SignName3        string    `json:"sign_name3"`
	SignPosition3    string    `json:"sign_position3"`
	SignName4        string    `json:"sign_name4"`
	SignPosition4    string    `json:"sign_position4"`
	IsDefault        uint8     `json:"is_default" binding:"required"`
	IsSystem         uint8     `json:"is_system" binding:"required"`
	IdSort           int       `json:"id_sort" binding:"required"`
	CreatedAt        time.Time `json:"created_at"`
	CreatedBy        string    `json:"created_by"`
	UpdatedAt        time.Time `json:"updated_at"`
	UpdatedBy        string    `json:"updated_by"`
	Id               uint64    `json:"id"`
}

type Ast_report_template_field struct {
	TemplateId      uint64 `json:"template_id" binding:"required"`
	FieldName       string `json:"field_name" binding:"required"`
	HeaderName      string `json:"header_name" binding:"required"`
	FooterType      int    `json:"footer_type" binding:"required"`
	FormatCode      int    `json:"format_code" binding:"required"`
	Width           int    `json:"width" binding:"required"`
	Font            int    `json:"font" binding:"required"`
	FontSize        int    `json:"font_size" binding:"required"`
	FontColor       int    `json:"font_color" binding:"required"`
	Alignment       string `json:"alignment" binding:"required"`
	HeaderFontSize  int    `json:"header_font_size" binding:"required"`
	HeaderFontColor int    `json:"header_font_color" binding:"required"`
	HeaderAlignment string `json:"header_alignment" binding:"required"`
	IdSort          int    `json:"id_sort" binding:"required"`
}
type Ast_user_group struct {
	Code                    string    `json:"code" binding:"required"`
	AccessForm              string    `json:"access_form" binding:"required"`
	AccessInventoryReceive  string    `json:"access_inventory_receive" binding:"required"`
	AccessFixedAssetReceive string    `json:"access_fixed_asset_receive" binding:"required"`
	AccessSpecial           string    `json:"access_special" binding:"required"`
	CreatedAt               time.Time `json:"created_at"`
	CreatedBy               string    `json:"created_by"`
	UpdatedAt               time.Time `json:"updated_at"`
	UpdatedBy               string    `json:"updated_by"`
	Id                      uint64    `json:"id" gorm:"primaryKey"`
}

type Ast_user_sub_department struct {
	UserCode              string    `json:"user_code" binding:"required"`
	SubDepartmentCode     string    `json:"sub_department_code" binding:"required"`
	IsCanInvPrApprove1    uint8     `json:"is_can_inv_pr_approve1" binding:"required"`
	IsCanInvPrApprove12   string    `json:"is_can_inv_pr_approve12" binding:"required"`
	IsCanInvPrApprove2    uint8     `json:"is_can_inv_pr_approve2" binding:"required"`
	IsCanInvPrApprove3    uint8     `json:"is_can_inv_pr_approve3" binding:"required"`
	IsCanInvPrAssignPrice uint8     `json:"is_can_inv_pr_assign_price" binding:"required"`
	IsCanInvSrApprove1    uint8     `json:"is_can_inv_sr_approve1" binding:"required"`
	IsCanInvSrApprove2    uint8     `json:"is_can_inv_sr_approve2" binding:"required"`
	CreatedAt             time.Time `json:"created_at"`
	CreatedBy             string    `json:"created_by"`
	UpdatedAt             time.Time `json:"updated_at"`
	UpdatedBy             string    `json:"updated_by"`
	Id                    uint64    `json:"id"`
}

type Audit_log struct {
	Id          uint64    `json:"id"`
	AuditDate   time.Time `json:"audit_date" binding:"required"`
	PostingDate time.Time `json:"posting_date" binding:"required"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   string    `json:"created_by"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedBy   string    `json:"updated_by"`
}

type Ban_booking struct {
	Number             uint64    `json:"number" gorm:"primaryKey"`
	CheckNumber        string    `json:"check_number"`
	IsContinueEvent    uint8     `json:"is_continue_event"`
	ContactPersonId    uint64    `json:"contact_person_id"`
	GuestDetailId      uint64    `json:"guest_detail_id"`
	GuestProfileId     uint64    `json:"guest_profile_id"`
	CurrencyCode       string    `json:"currency_code" `
	ExchangeRate       float64   `json:"exchange_rate" `
	IsConstantCurrency uint8     `json:"is_constant_currency" `
	ReservationBy      string    `json:"reservation_by"`
	ThemeCode          string    `json:"theme_code"`
	SeatingPlanCode    string    `json:"seating_plan_code"`
	LocationCode       string    `json:"location_code"`
	VenueCode          string    `json:"venue_code"`
	GroupCode          string    `json:"group_code"`
	MarketingCode      string    `json:"marketing_code"`
	DocumentNumber     string    `json:"document_number"`
	Notes              string    `json:"notes"`
	EstimateRevenue    float64   `json:"estimate_revenue"`
	BeoNote            string    `json:"beo_note" binding:"required"`
	ShowNotes          uint8     `json:"show_notes"`
	AuditDate          time.Time `json:"audit_date"`
	CancelAuditDate    time.Time `json:"cancel_audit_date"`
	CancelDate         time.Time `json:"cancel_date"`
	CancelBy           string    `json:"cancel_by"`
	CancelReason       string    `json:"cancel_reason"`
	StatusCode         string    `json:"status_code"`
	ReservationType    string    `json:"reservation_type" binding:"required"`
	IsLock             uint8     `json:"is_lock"`
	IsPublic           uint8     `json:"is_public"`
	ChangeStatusDate   time.Time `json:"change_status_date"`
	ChangeStatusBy     string    `json:"change_status_by"`
	FolioTransfer      uint64    `json:"folio_transfer"`
	CreatedAt          time.Time `json:"created_at"`
	CreatedBy          string    `json:"created_by"`
	UpdatedAt          time.Time `json:"updated_at"`
	UpdatedBy          string    `json:"updated_by"`
}

type Ban_cfg_init_seating_plan struct {
	Code            string    `json:"code" binding:"required"`
	Name            string    `json:"name"`
	AssignVenueCode string    `json:"assign_venue_code"`
	Image           []byte    `json:"image"`
	Description     string    `json:"description"`
	CreatedAt       time.Time `json:"created_at"`
	CreatedBy       string    `json:"created_by"`
	UpdatedAt       time.Time `json:"updated_at"`
	UpdatedBy       string    `json:"updated_by"`
	Id              uint64    `json:"id"`
}

type Ban_cfg_init_theme struct {
	Code        string    `json:"code" binding:"required"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   string    `json:"created_by"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedBy   string    `json:"updated_by"`
	Id          uint64    `json:"id"`
}

type Ban_cfg_init_venue struct {
	Code           string    `json:"code" binding:"required"`
	Name           string    `json:"name"`
	VenueGroupCode string    `json:"venue_group_code" binding:"required"`
	LocationCode   string    `json:"location_code"`
	IdSort         int       `json:"id_sort" binding:"required"`
	CreatedAt      time.Time `json:"created_at"`
	CreatedBy      string    `json:"created_by"`
	UpdatedAt      time.Time `json:"updated_at"`
	UpdatedBy      string    `json:"updated_by"`
	Id             uint64    `json:"id"`
}

type Ban_cfg_init_venue_combine struct {
	Code         string    `json:"code" binding:"required"`
	Name         string    `json:"name" binding:"required"`
	GroupNumber  uint8     `json:"group_number"`
	LocationCode string    `json:"location_code" binding:"required"`
	CreatedAt    time.Time `json:"created_at"`
	CreatedBy    string    `json:"created_by"`
	UpdatedAt    time.Time `json:"updated_at"`
	UpdatedBy    string    `json:"updated_by"`
	Id           uint64    `json:"id"`
}

type Ban_cfg_init_venue_combine_detail struct {
	CombineVenueCode string    `json:"combine_venue_code" binding:"required"`
	VenueCode        string    `json:"venue_code" binding:"required"`
	CombineGroup     string    `json:"combine_group"`
	CreatedAt        time.Time `json:"created_at"`
	CreatedBy        string    `json:"created_by"`
	UpdatedAt        time.Time `json:"updated_at"`
	UpdatedBy        string    `json:"updated_by"`
	Id               uint64    `json:"id"`
}

type Ban_cfg_init_venue_group struct {
	Code      string    `json:"code" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	IdSort    int       `json:"id_sort"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Ban_const_booking_status struct {
	Code string `json:"code" binding:"required"`
	Name string `json:"name"`
}

type Ban_const_reservation_status struct {
	Code   string `json:"code" binding:"required"`
	Name   string `json:"name" binding:"required"`
	IdSort string `json:"id_sort" binding:"required"`
	Icon   []byte `json:"icon" binding:"required"`
}

type Ban_const_reservation_type struct {
	Code     string `json:"code" binding:"required"`
	Name     string `json:"name"`
	IsActive string `json:"is_active"`
	IdSort   int    `json:"id_sort"`
}

type Ban_const_venue_location struct {
	Code string `json:"code" binding:"required"`
	Name string `json:"name"`
}

type Ban_report struct {
	Code        int       `json:"code" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	ReportQuery string    `json:"report_query" binding:"required"`
	ParentId    uint64    `json:"parent_id" binding:"required"`
	IsSystem    uint8     `json:"is_system" binding:"required"`
	IdSort      int       `json:"id_sort" binding:"required"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   string    `json:"created_by"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedBy   string    `json:"updated_by"`
	Id          uint64    `json:"id"`
}

type Ban_report_default_field struct {
	ReportCode int    `json:"report_code" binding:"required"`
	FieldQuery string `json:"field_query" binding:"required"`
	FieldName  string `json:"field_name" binding:"required"`
	IdSort     int    `json:"id_sort" binding:"required"`
}

type Ban_report_group_field struct {
	TemplateId uint64 `json:"template_id" binding:"required"`
	FieldName  string `json:"field_name" binding:"required"`
	IdSort     int    `json:"id_sort" binding:"required"`
}

type Ban_report_grouping_field struct {
	TemplateId uint64 `json:"template_id" binding:"required"`
	FieldName  string `json:"field_name" binding:"required"`
	IdSort     int    `json:"id_sort" binding:"required"`
}

type Ban_report_order_field struct {
	TemplateId  uint64 `json:"template_id" binding:"required"`
	FieldName   string `json:"field_name" binding:"required"`
	IsAscending uint8  `json:"is_ascending" binding:"required"`
	IdSort      int    `json:"id_sort" binding:"required"`
}

type Ban_report_template struct {
	ReportCode       int       `json:"report_code" binding:"required"`
	Name             string    `json:"name" binding:"required"`
	GroupLevel       int       `json:"group_level" binding:"required"`
	HeaderRemark     string    `json:"header_remark"`
	ShowFooter       uint8     `json:"show_footer" binding:"required"`
	PaperSize        int       `json:"paper_size" binding:"required"`
	PaperWidth       float64   `json:"paper_width" binding:"required"`
	PaperHeight      float64   `json:"paper_height" binding:"required"`
	IsPortrait       uint8     `json:"is_portrait" binding:"required"`
	HeaderRowHeight  int       `json:"header_row_height" binding:"required"`
	RowHeight        int       `json:"row_height" binding:"required"`
	HorizontalBorder uint8     `json:"horizontal_border" binding:"required"`
	VerticalBorder   uint8     `json:"vertical_border" binding:"required"`
	SignName1        string    `json:"sign_name1"`
	SignPosition1    string    `json:"sign_position1"`
	SignName2        string    `json:"sign_name2"`
	SignPosition2    string    `json:"sign_position2"`
	SignName3        string    `json:"sign_name3"`
	SignPosition3    string    `json:"sign_position3"`
	SignName4        string    `json:"sign_name4"`
	SignPosition4    string    `json:"sign_position4"`
	IsDefault        uint8     `json:"is_default" binding:"required"`
	IsSystem         uint8     `json:"is_system" binding:"required"`
	IdSort           int       `json:"id_sort" binding:"required"`
	CreatedAt        time.Time `json:"created_at"`
	CreatedBy        string    `json:"created_by"`
	UpdatedAt        time.Time `json:"updated_at"`
	UpdatedBy        string    `json:"updated_by"`
	Id               uint64    `json:"id"`
}

type Ban_report_template_field struct {
	TemplateId      uint64 `json:"template_id" binding:"required"`
	FieldName       string `json:"field_name" binding:"required"`
	HeaderName      string `json:"header_name" binding:"required"`
	FooterType      int    `json:"footer_type" binding:"required"`
	FormatCode      int    `json:"format_code" binding:"required"`
	Width           int    `json:"width" binding:"required"`
	Font            int    `json:"font" binding:"required"`
	FontSize        int    `json:"font_size" binding:"required"`
	FontColor       int    `json:"font_color" binding:"required"`
	Alignment       string `json:"alignment" binding:"required"`
	HeaderFontSize  int    `json:"header_font_size" binding:"required"`
	HeaderFontColor int    `json:"header_font_color" binding:"required"`
	HeaderAlignment string `json:"header_alignment" binding:"required"`
	IdSort          int    `json:"id_sort" binding:"required"`
}

type Ban_reservation struct {
	Number             uint64    `json:"number" binding:"required" gorm:"primaryKey"`
	Booking            uint64    `json:"booking" binding:"required"`
	CheckNumber        string    `json:"check_number"`
	IsContinueEvent    uint8     `json:"is_continue_event"`
	ContactPersonId    uint64    `json:"contact_person_id"`
	GuestDetailId      uint64    `json:"guest_detail_id"`
	GuestProfileId     uint64    `json:"guest_profile_id"`
	CurrencyCode       string    `json:"currency_code" binding:"required"`
	ExchangeRate       float64   `json:"exchange_rate" binding:"required"`
	IsConstantCurrency uint8     `json:"is_constant_currency" binding:"required"`
	ReservationBy      string    `json:"reservation_by"`
	ThemeCode          string    `json:"theme_code"`
	SeatingPlanCode    string    `json:"seating_plan_code"`
	LocationCode       string    `json:"location_code"`
	VenueCombineCode   string    `json:"venue_combine_code"`
	VenueCombineNumber uint64    `json:"venue_combine_number"`
	VenueCode          string    `json:"venue_code"`
	GroupCode          string    `json:"group_code"`
	MarketingCode      string    `json:"marketing_code"`
	DocumentNumber     string    `json:"document_number"`
	Notes              string    `json:"notes"`
	ShowNotes          uint8     `json:"show_notes" binding:"required"`
	AuditDate          time.Time `json:"audit_date"`
	CancelAuditDate    time.Time `json:"cancel_audit_date"`
	CancelDate         time.Time `json:"cancel_date"`
	CancelBy           string    `json:"cancel_by"`
	CancelReason       string    `json:"cancel_reason"`
	StatusCode         string    `json:"status_code"`
	ReservationType    string    `json:"reservation_type" binding:"required"`
	IsLock             uint8     `json:"is_lock"`
	IsPublic           uint8     `json:"is_public"`
	ChangeStatusDate   time.Time `json:"change_status_date"`
	ChangeStatusBy     string    `json:"change_status_by"`
	FolioTransfer      uint64    `json:"folio_transfer" binding:"required"`
	CreatedAt          time.Time `json:"created_at"`
	CreatedBy          string    `json:"created_by"`
	UpdatedAt          time.Time `json:"updated_at"`
	UpdatedBy          string    `json:"updated_by"`
}

type Ban_reservation_charge struct {
	BookingNumber     uint64    `json:"booking_number" binding:"required"`
	ReservationNumber uint64    `json:"reservation_number" binding:"required"`
	ServedTime        time.Time `json:"served_time"`
	ServedEndTime     time.Time `json:"served_end_time"`
	OutletCode        string    `json:"outlet_code"`
	ProductCode       string    `json:"product_code"`
	VenueCode         string    `json:"venue_code"`
	SeatingPlanCode   string    `json:"seating_plan_code"`
	PackageCode       string    `json:"package_code"`
	PackageRef        uint64    `json:"package_ref"`
	CompanyCode       string    `json:"company_code" binding:"required"`
	AccountCode       string    `json:"account_code"`
	Description       string    `json:"description"`
	Quantity          float64   `json:"quantity"`
	PricePurchase     float64   `json:"price_purchase"`
	PriceOriginal     float64   `json:"price_original"`
	Price             float64   `json:"price"`
	Discount          float64   `json:"discount"`
	Tax               float64   `json:"tax"`
	Service           float64   `json:"service"`
	Remark            string    `json:"remark"`
	TaxAndServiceCode string    `json:"tax_and_service_code" binding:"required"`
	TypeCode          string    `json:"type_code"`
	AuditDate         time.Time `json:"audit_date"`
	PostingDate       time.Time `json:"posting_date"`
	Void              uint8     `json:"void"`
	VoidDate          time.Time `json:"void_date" binding:"required"`
	VoidBy            string    `json:"void_by" binding:"required"`
	VoidReason        string    `json:"void_reason" binding:"required"`
	LayoutID          uint64    `json:"layout_id" binding:"required"`
	InputOf           string    `json:"input_of"`
	SubFolioId        uint64    `json:"sub_folio_id"`
	IsPosting         uint8     `json:"is_posting"`
	IsBeo             uint8     `json:"is_beo"`
	CreatedAt         time.Time `json:"created_at"`
	CreatedBy         string    `json:"created_by"`
	UpdatedAt         time.Time `json:"updated_at"`
	UpdatedBy         string    `json:"updated_by"`
	Id                uint64    `json:"id"`
}

type Ban_reservation_remark struct {
	BookingNumber uint64    `json:"booking_number" binding:"required"`
	Number        int       `json:"number"`
	Header        string    `json:"header"`
	Remark        string    `json:"remark"`
	CreatedAt     time.Time `json:"created_at"`
	CreatedBy     string    `json:"created_by"`
	UpdatedAt     time.Time `json:"updated_at"`
	UpdatedBy     string    `json:"updated_by"`
	Id            uint64    `json:"id"`
}

type Ban_user_group struct {
	Code               string    `json:"code" binding:"required"`
	AccessForm         string    `json:"access_form" binding:"required"`
	AccessSpecial      string    `json:"access_special" binding:"required"`
	AccessReservation  string    `json:"access_reservation"`
	AccessDeposit      string    `json:"access_deposit"`
	AccessInHouse      string    `json:"access_in_house" binding:"required"`
	AccessFolio        string    `json:"access_folio"`
	AccessFolioHistory string    `json:"access_folio_history"`
	CreatedAt          time.Time `json:"created_at"`
	CreatedBy          string    `json:"created_by"`
	UpdatedAt          time.Time `json:"updated_at"`
	UpdatedBy          string    `json:"updated_by"`
	Id                 uint64    `json:"id" gorm:"primaryKey"`
}

type Breakfast_list_temp struct {
	GuestName  string `json:"guest_name" binding:"required"`
	RoomNumber string `json:"room_number" binding:"required"`
	GroupName  string `json:"group_name" binding:"required"`
	IdSort     int    `json:"id_sort"`
}

type Budget_expense struct {
	Period             int       `json:"period" binding:"required"`
	SubDepartmentCode  string    `json:"sub_department_code" binding:"required"`
	JournalAccountCode string    `json:"journal_account_code" binding:"required"`
	Remark             string    `json:"remark"`
	Amount             float64   `json:"amount" binding:"required"`
	TypeCode           string    `json:"type_code" binding:"required"`
	M01                *float64  `json:"M01" binding:"required"`
	M02                *float64  `json:"M02" binding:"required"`
	M03                *float64  `json:"M03" binding:"required"`
	M04                *float64  `json:"M04" binding:"required"`
	M05                *float64  `json:"M05" binding:"required"`
	M06                *float64  `json:"M06" binding:"required"`
	M07                *float64  `json:"M07" binding:"required"`
	M08                *float64  `json:"M08" binding:"required"`
	M09                *float64  `json:"M09" binding:"required"`
	M10                *float64  `json:"M10" binding:"required"`
	M11                *float64  `json:"M11" binding:"required"`
	M12                *float64  `json:"M12" binding:"required"`
	UnitCode           string    `json:"unit_code" binding:"required"`
	CreatedAt          time.Time `json:"created_at"`
	CreatedBy          string    `json:"created_by"`
	UpdatedAt          time.Time `json:"updated_at"`
	UpdatedBy          string    `json:"updated_by"`
	Id                 uint64    `json:"id"`
	IdHolding          uint64    `json:"id_log_holding"`
}

type Budget_fb struct {
	Period     int       `json:"period" binding:"required"`
	OutletCode string    `json:"outlet_code" binding:"required"`
	Code       string    `json:"code" binding:"required"`
	Remark     string    `json:"remark"`
	Amount     int       `json:"amount" binding:"required"`
	TypeCode   string    `json:"type_code" binding:"required"`
	M01        float64   `json:"M01" binding:"required"`
	M02        float64   `json:"M02" binding:"required"`
	M03        float64   `json:"M03" binding:"required"`
	M04        float64   `json:"M04" binding:"required"`
	M05        float64   `json:"M05" binding:"required"`
	M06        float64   `json:"M06" binding:"required"`
	M07        float64   `json:"M07" binding:"required"`
	M08        float64   `json:"M08" binding:"required"`
	M09        float64   `json:"M09" binding:"required"`
	M10        float64   `json:"M10" binding:"required"`
	M11        float64   `json:"M11" binding:"required"`
	M12        float64   `json:"M12" binding:"required"`
	CreatedAt  time.Time `json:"created_at"`
	CreatedBy  string    `json:"created_by"`
	UpdatedAt  time.Time `json:"updated_at"`
	UpdatedBy  string    `json:"updated_by"`
	Id         uint64    `json:"id"`
}

type Budget_income struct {
	Period            int       `json:"period" binding:"required" gorm:"primaryKey"`
	SubDepartmentCode string    `json:"sub_department_code" binding:"required"`
	AccountCode       string    `json:"account_code" binding:"required"`
	Remark            string    `json:"remark"`
	Amount            float64   `json:"amount" binding:"required"`
	TypeCode          string    `json:"type_code" binding:"required"`
	M01               *float64  `json:"M01" binding:"required"`
	M02               *float64  `json:"M02" binding:"required"`
	M03               *float64  `json:"M03" binding:"required"`
	M04               *float64  `json:"M04" binding:"required"`
	M05               *float64  `json:"M05" binding:"required"`
	M06               *float64  `json:"M06" binding:"required"`
	M07               *float64  `json:"M07" binding:"required"`
	M08               *float64  `json:"M08" binding:"required"`
	M09               *float64  `json:"M09" binding:"required"`
	M10               *float64  `json:"M10" binding:"required"`
	M11               *float64  `json:"M11" binding:"required"`
	M12               *float64  `json:"M12" binding:"required"`
	UnitCode          string    `json:"unit_code" binding:"required"`
	CreatedAt         time.Time `json:"created_at"`
	CreatedBy         string    `json:"created_by"`
	UpdatedAt         time.Time `json:"updated_at"`
	UpdatedBy         string    `json:"updated_by"`
	Id                uint64    `json:"id"`
	IdHolding         uint64    `json:"id_log_holding"`
}

type Budget_statistic struct {
	Period            int       `json:"period" binding:"required"`
	SubDepartmentCode string    `json:"sub_department_code" binding:"required"`
	Code              string    `json:"code" binding:"required"`
	Remark            string    `json:"remark"`
	Amount            float64   `json:"amount" binding:"required"`
	TypeCode          string    `json:"type_code" binding:"required"`
	M01               *float64  `json:"M01" binding:"required"`
	M02               *float64  `json:"M02" binding:"required"`
	M03               *float64  `json:"M03" binding:"required"`
	M04               *float64  `json:"M04" binding:"required"`
	M05               *float64  `json:"M05" binding:"required"`
	M06               *float64  `json:"M06" binding:"required"`
	M07               *float64  `json:"M07" binding:"required"`
	M08               *float64  `json:"M08" binding:"required"`
	M09               *float64  `json:"M09" binding:"required"`
	M10               *float64  `json:"M10" binding:"required"`
	M11               *float64  `json:"M11" binding:"required"`
	M12               *float64  `json:"M12" binding:"required"`
	UnitCode          string    `json:"unit_code" binding:"required"`
	CreatedAt         time.Time `json:"created_at"`
	CreatedBy         string    `json:"created_by"`
	UpdatedAt         time.Time `json:"updated_at"`
	UpdatedBy         string    `json:"updated_by"`
	Id                uint64    `json:"id"`
	IdHolding         uint64    `json:"id_log_holding"`
}

type Cash_count struct {
	Id           uint64    `json:"id"`
	LogShiftId   uint64    `json:"log_shift_id" binding:"required"`
	CurrencySign string    `json:"currency_sign" binding:"required"`
	Nominal      float64   `json:"nominal" binding:"required"`
	Quantity     int       `json:"quantity" binding:"required"`
	CreatedAt    time.Time `json:"created_at"`
	CreatedBy    string    `json:"created_by"`
	UpdatedAt    time.Time `json:"updated_at"`
	UpdatedBy    string    `json:"updated_by"`
}

type Cfg_init_account struct {
	Code               string    `json:"code" gorm:"primaryKey" binding:"required"`
	Name               string    `json:"name" binding:"required"`
	TypeCode           string    `json:"type_code" binding:"required"`
	SubGroupCode       string    `json:"sub_group_code" binding:"required"`
	SubDepartmentCode  string    `json:"sub_department_code"`
	JournalAccountCode string    `json:"journal_account_code" binding:"required"`
	TaxAndServiceCode  string    `json:"tax_and_service_code" binding:"required"`
	ItemGroupCode      string    `json:"item_group_code"`
	SubFolioGroupCode  string    `json:"sub_folio_group_code" binding:"required"`
	IsPayment          uint8     `json:"is_payment"`
	IsRefund           uint8     `json:"is_refund"`
	IsRoomCharge       uint8     `json:"is_room_charge"`
	IdSort             int       `json:"id_sort"`
	CreatedAt          time.Time `json:"created_at"`
	CreatedBy          string    `json:"created_by"`
	UpdatedAt          time.Time `json:"updated_at"`
	UpdatedBy          string    `json:"updated_by"`
	Id                 uint64    `json:"id"`
}

type Cfg_init_account_sub_group struct {
	Code      string    `json:"code"  gorm:"primaryKey" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	GroupCode string    `json:"group_code" binding:"required"`
	IdSort    int       `json:"id_sort"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Cfg_init_bed_type struct {
	Code      string    `json:"code" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	IdSort    int       `json:"id_sort"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Cfg_init_booking_source struct {
	Code      string    `json:"code" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Cfg_init_card_bank struct {
	Code      string    `json:"code" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	IdSort    int       `json:"id_sort"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Cfg_init_card_type struct {
	Code      string    `json:"code" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	IdSort    int       `json:"id_sort"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Cfg_init_city struct {
	Code        string    `json:"code" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	StateCode   string    `json:"state_code" binding:"required"`
	RegencyCode string    `json:"regency_code" binding:"required"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   string    `json:"created_by"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedBy   string    `json:"updated_by"`
	Id          uint64    `json:"id"`
}

type Cfg_init_company_type struct {
	Code                 string    `json:"code" binding:"required"`
	Name                 string    `json:"name" binding:"required"`
	JournalAccountCodeAp string    `json:"journal_account_code_ap" binding:"required"`
	JournalAccountCodeAr string    `json:"journal_account_code_ar" binding:"required"`
	IsPersonal           uint8     `json:"is_personal"`
	IdSort               int       `json:"id_sort"`
	CreatedAt            time.Time `json:"created_at"`
	CreatedBy            string    `json:"created_by"`
	UpdatedAt            time.Time `json:"updated_at"`
	UpdatedBy            string    `json:"updated_by"`
	Id                   uint64    `json:"id"`
}

type Cfg_init_competitor_category struct {
	Code      string    `json:"code" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	IdSort    int       `json:"id_sort"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Cfg_init_continent struct {
	Code      string    `json:"code" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Cfg_init_country struct {
	Code          string    `json:"code" binding:"required"`
	Name          string    `json:"name" binding:"required"`
	ContinentCode string    `json:"continent_code" binding:"required"`
	IsoCode       string    `json:"iso_code"`
	CreatedAt     time.Time `json:"created_at"`
	CreatedBy     string    `json:"created_by"`
	UpdatedAt     time.Time `json:"updated_at"`
	UpdatedBy     string    `json:"updated_by"`
	Id            uint64    `json:"id"`
}

type Cfg_init_credit_card_charge struct {
	AccountCode   string    `json:"account_code" binding:"required"`
	ChargePercent float64   `json:"charge_percent" binding:"required"`
	CreatedAt     time.Time `json:"created_at"`
	CreatedBy     string    `json:"created_by"`
	UpdatedAt     time.Time `json:"updated_at"`
	UpdatedBy     string    `json:"updated_by"`
	Id            uint64    `json:"id"`
}

type Cfg_init_currency struct {
	Code         string    `json:"code" binding:"required"`
	Name         string    `json:"name" binding:"required"`
	AccountCode  string    `json:"account_code" binding:"required"`
	ExchangeRate float64   `json:"exchange_rate" binding:"required"`
	Symbol       string    `json:"symbol" binding:"required"`
	Format       string    `json:"format" binding:"required"`
	IsDefault    uint8     `json:"is_default"`
	CreatedAt    time.Time `json:"created_at"`
	CreatedBy    string    `json:"created_by"`
	UpdatedAt    time.Time `json:"updated_at"`
	UpdatedBy    string    `json:"updated_by"`
	Id           uint64    `json:"id"`
}

type Cfg_init_currency_nominal struct {
	CurrencySign string    `json:"currency_sign" binding:"required"`
	Nominal      float64   `json:"nominal" binding:"required"`
	IdSort       int       `json:"id_sort"`
	CreatedAt    time.Time `json:"created_at"`
	CreatedBy    string    `json:"created_by"`
	UpdatedAt    time.Time `json:"updated_at"`
	UpdatedBy    string    `json:"updated_by"`
	Id           uint64    `json:"id"`
}

type Cfg_init_custom_lookup_field01 struct {
	Code      string    `json:"code" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Cfg_init_custom_lookup_field02 struct {
	Code      string    `json:"code" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Cfg_init_custom_lookup_field03 struct {
	Code      string    `json:"code" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Cfg_init_custom_lookup_field04 struct {
	Code      string    `json:"code" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Cfg_init_custom_lookup_field05 struct {
	Code      string    `json:"code" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Cfg_init_custom_lookup_field06 struct {
	Code      string    `json:"code" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Cfg_init_custom_lookup_field07 struct {
	Code      string    `json:"code" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Cfg_init_custom_lookup_field08 struct {
	Code      string    `json:"code" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Cfg_init_custom_lookup_field09 struct {
	Code      string    `json:"code" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Cfg_init_custom_lookup_field10 struct {
	Code      string    `json:"code" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Cfg_init_custom_lookup_field11 struct {
	Code      string    `json:"code" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Cfg_init_custom_lookup_field12 struct {
	Code      string    `json:"code" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Cfg_init_department struct {
	Code      string    `json:"code" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	TypeCode  string    `json:"type_code" binding:"required"`
	IdSort    int       `json:"id_sort"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Cfg_init_guest_type struct {
	Code      string    `json:"code" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	Color     string    `json:"color" binding:"required"`
	IdSort    int       `json:"id_sort"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Cfg_init_id_card_type struct {
	Code      string    `json:"code" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Cfg_init_is_fb_sub_department_group struct {
	Code      string    `json:"code" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	IsDetail  string    `json:"is_detail"`
	IdSort    int       `json:"id_sort"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Cfg_init_is_fb_sub_department_group_detail struct {
	GroupCode         string    `json:"group_code" binding:"required"`
	SubDepartmentCode string    `json:"sub_department_code" binding:"required"`
	CreatedAt         time.Time `json:"created_at"`
	CreatedBy         string    `json:"created_by"`
	UpdatedAt         time.Time `json:"updated_at"`
	UpdatedBy         string    `json:"updated_by"`
	Id                uint64    `json:"id"`
}

type Cfg_init_journal_account struct {
	Code              string    `json:"code" binding:"required"`
	Name              string    `json:"name" binding:"required"`
	SubGroupCode      string    `json:"sub_group_code" binding:"required"`
	SubDepartmentCode string    `json:"sub_department_code"`
	TypeCode          string    `json:"type_code" binding:"required"`
	ItemGroupCode     string    `json:"item_group_code"`
	CategoryCode      string    `json:"category_code"`
	Description       string    `json:"description"`
	IsTaxExpense      uint8     `json:"is_tax_expense"`
	IsNew             uint8     `json:"is_new"`
	IdSort            int       `json:"id_sort"`
	CreatedAt         time.Time `json:"created_at"`
	CreatedBy         string    `json:"created_by"`
	UpdatedAt         time.Time `json:"updated_at"`
	UpdatedBy         string    `json:"updated_by"`
	Id                uint64    `json:"id"`
}

type Cfg_init_journal_account_category struct {
	Code      string    `json:"code" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	IdSort    int       `json:"id_sort"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Cfg_init_journal_account_sub_group struct {
	Code       string    `json:"code" binding:"required"`
	Name       string    `json:"name" binding:"required"`
	GroupCode  string    `json:"group_code" binding:"required"`
	TypeCode   string    `json:"type_code"`
	IsHeader   uint8     `json:"is_header"`
	ParentCode string    `json:"parent_code"`
	CreatedAt  time.Time `json:"created_at"`
	CreatedBy  string    `json:"created_by"`
	UpdatedAt  time.Time `json:"updated_at"`
	UpdatedBy  string    `json:"updated_by"`
	Id         uint64    `json:"id"`
}

type Cfg_init_language struct {
	Code      string    `json:"code" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Cfg_init_loan_item struct {
	Code      string    `json:"code" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	IdSort    int       `json:"id_sort"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Cfg_init_market struct {
	Code         string    `json:"code" binding:"required"`
	Name         string    `json:"name" binding:"required"`
	CategoryCode string    `json:"category_code"`
	CreatedAt    time.Time `json:"created_at"`
	CreatedBy    string    `json:"created_by"`
	UpdatedAt    time.Time `json:"updated_at"`
	UpdatedBy    string    `json:"updated_by"`
	Id           uint64    `json:"id"`
}

type Cfg_init_market_category struct {
	Code      string    `json:"code" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	IdSort    int       `json:"id_sort"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Cfg_init_member_point_type struct {
	Code           string    `json:"code" binding:"required"`
	Name           string    `json:"name" binding:"required"`
	MemberTypeCode string    `json:"member_type_code" binding:"required"`
	IsFromRate     uint8     `json:"is_from_rate"`
	RoomTypeCode   string    `json:"room_type_code"`
	RateAmount     float64   `json:"rate_amount" binding:"required"`
	Point          float64   `json:"point" binding:"required"`
	CreatedAt      time.Time `json:"created_at"`
	CreatedBy      string    `json:"created_by"`
	UpdatedAt      time.Time `json:"updated_at"`
	UpdatedBy      string    `json:"updated_by"`
	Id             uint64    `json:"id"`
}

type Cfg_init_nationality struct {
	Code      string    `json:"code" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Cfg_init_owner struct {
	Code      string    `json:"code" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Cfg_init_pabx_rate struct {
	Code        string    `json:"code" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	Rate        float64   `json:"rate" binding:"required"`
	FreeSeconds int       `json:"free_seconds" binding:"required"`
	TypeCode    string    `json:"type_code" binding:"required"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   string    `json:"created_by"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedBy   string    `json:"updated_by"`
	Id          uint64    `json:"id"`
}

type Cfg_init_package struct {
	Code                string    `json:"code" binding:"required"`
	Name                string    `json:"name" binding:"required"`
	OutletCode          string    `json:"outlet_code"`
	ProductCode         string    `json:"product_code"`
	SubDepartmentCode   string    `json:"sub_department_code" binding:"required"`
	AccountCode         string    `json:"account_code" binding:"required"`
	Quantity            float64   `json:"quantity" binding:"required"`
	Amount              float64   `json:"amount" binding:"required"`
	PerPax              uint8     `json:"per_pax"`
	IncludeChild        uint8     `json:"include_child"`
	TaxAndServiceCode   string    `json:"tax_and_service_code" binding:"required"`
	ChargeFrequencyCode string    `json:"charge_frequency_code" binding:"required"`
	MaxPax              int       `json:"max_pax"`
	ExtraPax            float64   `json:"extra_pax"`
	PerPaxExtra         uint8     `json:"per_pax_extra"`
	ShowInTransaction   uint8     `json:"show_in_transaction"`
	IsSent              uint8     `json:"is_sent"`
	IsOnline            uint8     `json:"is_online"`
	IsActive            uint8     `json:"is_active"`
	CreatedAt           time.Time `json:"created_at"`
	CreatedBy           string    `json:"created_by"`
	UpdatedAt           time.Time `json:"updated_at"`
	UpdatedBy           string    `json:"updated_by"`
	Id                  uint64    `json:"id"`
}

type Cfg_init_package_breakdown struct {
	PackageCode         string    `json:"package_code" binding:"required"`
	OutletCode          string    `json:"outlet_code"`
	ProductCode         string    `json:"product_code"`
	SubDepartmentCode   string    `json:"sub_department_code" binding:"required"`
	AccountCode         string    `json:"account_code" binding:"required"`
	CompanyCode         string    `json:"company_code"`
	Quantity            float64   `json:"quantity" binding:"required"`
	IsAmountPercent     uint8     `json:"is_amount_percent"`
	Amount              float64   `json:"amount" binding:"required"`
	PerPax              uint8     `json:"per_pax"`
	IncludeChild        uint8     `json:"include_child"`
	Remark              string    `json:"remark"`
	TaxAndServiceCode   string    `json:"tax_and_service_code" binding:"required"`
	ChargeFrequencyCode string    `json:"charge_frequency_code" binding:"required"`
	MaxPax              int       `json:"max_pax"`
	ExtraPax            float64   `json:"extra_pax"`
	PerPaxExtra         uint8     `json:"per_pax_extra"`
	CreatedAt           time.Time `json:"created_at"`
	CreatedBy           string    `json:"created_by"`
	UpdatedAt           time.Time `json:"updated_at"`
	UpdatedBy           string    `json:"updated_by"`
	Id                  uint64    `json:"id"`
}

type Cfg_init_package_business_source struct {
	PackageCode        string    `json:"package_code" binding:"required"`
	AccountCode        string    `json:"account_code" binding:"required"`
	CompanyCode        string    `json:"company_code" binding:"required"`
	CommissionTypeCode string    `json:"commission_type_code" binding:"required"`
	CommissionValue    float64   `json:"commission_value" binding:"required"`
	CreatedAt          time.Time `json:"created_at"`
	CreatedBy          string    `json:"created_by"`
	UpdatedAt          time.Time `json:"updated_at"`
	UpdatedBy          string    `json:"updated_by"`
	Id                 uint64    `json:"id"`
}

type Cfg_init_payment_type struct {
	Code      string    `json:"code" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	GroupCode string    `json:"group_code" binding:"required"`
	IdSort    int       `json:"id_sort"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Cfg_init_phone_book_type struct {
	Code      string    `json:"code" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	IdSort    int       `json:"id_sort"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Cfg_init_printer struct {
	Code      string    `json:"code" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	Location  string    `json:"location"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Cfg_init_purpose_of struct {
	Code      string    `json:"code" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	IdSort    int       `json:"id_sort"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Cfg_init_regency struct {
	Code      string    `json:"code" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Cfg_init_reservation_mark struct {
	Code      string    `json:"code" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Cfg_init_room struct {
	Number             string    `json:"number" binding:"required"`
	Name               string    `json:"name"`
	LockNumber         string    `json:"lock_number"`
	RoomTypeCode       string    `json:"room_type_code" binding:"required"`
	BedTypeCode        string    `json:"bed_type_code" binding:"required"`
	ViewCode           string    `json:"view_code"`
	IsSmoking          uint8     `json:"is_smoking"`
	Building           string    `json:"building" binding:"required"`
	Floor              string    `json:"floor" binding:"required"`
	MaxAdult           int       `json:"max_adult" binding:"required"`
	Description        string    `json:"description"`
	PhoneNumber        string    `json:"phone_number"`
	TvQuantity         int       `json:"tv_quantity" binding:"required"`
	StartDate          time.Time `json:"start_date" binding:"required"`
	OwnerCode          string    `json:"owner_code"`
	IdSort             int       `json:"id_sort"`
	Image              string    `json:"image"`
	StatusCode         string    `json:"status_code"`
	BlockStatusCode    string    `json:"block_status_code"`
	TempStatusCode     string    `json:"temp_status_code"`
	RevenueAccountCode string    `json:"revenue_account_code"`
	Remark             string    `json:"remark"`
	PosX               *int      `json:"pos_x"`
	PosY               *int      `json:"pos_y"`
	Width              *int      `json:"width"`
	Height             *int      `json:"height"`
	CreatedAt          time.Time `json:"created_at"`
	CreatedBy          string    `json:"created_by"`
	UpdatedAt          time.Time `json:"updated_at"`
	UpdatedBy          string    `json:"updated_by"`
	Id                 uint64    `json:"id"`
}

type Cfg_init_room_allotment_type struct {
	Code      string    `json:"code" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	IdSort    int       `json:"id_sort"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	IdLog     uint64    `json:"id_log" binding:"required"`
}

type Cfg_init_room_amenities struct {
	Code        string    `json:"code" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"`
	Image       []byte    `json:"image"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   string    `json:"created_by"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedBy   string    `json:"updated_by"`
	Id          uint64    `json:"id"`
}

type Cfg_init_room_boy struct {
	Code      string    `json:"code" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	IsActive  uint8     `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Cfg_init_room_rate struct {
	Code                 string    `json:"code" binding:"required"`
	Name                 string    `json:"name" binding:"required"`
	RoomTypeCode         string    `json:"room_type_code" binding:"required"`
	FromDate             time.Time `json:"from_date" binding:"required"`
	ToDate               time.Time `json:"to_date" binding:"required"`
	SubCategoryCode      string    `json:"sub_category_code"`
	CompanyCode          string    `json:"company_code"`
	MarketCode           string    `json:"market_code"`
	DynamicRateTypeCode  string    `json:"dynamic_rate_type_code" binding:"required"`
	IsLastDeal           uint8     `json:"is_last_deal"`
	IsRateStructure      uint8     `json:"is_rate_structure"`
	IsCompliment         uint8     `json:"is_compliment"`
	IncludeBreakfast     uint8     `json:"include_breakfast"`
	WeekdayRate1         float64   `json:"weekday_rate1"`
	WeekdayRate2         float64   `json:"weekday_rate2"`
	WeekdayRate3         float64   `json:"weekday_rate3"`
	WeekdayRate4         float64   `json:"weekday_rate4"`
	WeekendRate1         float64   `json:"weekend_rate1"`
	WeekendRate2         float64   `json:"weekend_rate2"`
	WeekendRate3         float64   `json:"weekend_rate3"`
	WeekendRate4         float64   `json:"weekend_rate4"`
	WeekdayRateChild1    float64   `json:"weekday_rate_child1"`
	WeekdayRateChild2    float64   `json:"weekday_rate_child2"`
	WeekdayRateChild3    float64   `json:"weekday_rate_child3"`
	WeekdayRateChild4    float64   `json:"weekday_rate_child4"`
	WeekendRateChild1    float64   `json:"weekend_rate_child1"`
	WeekendRateChild2    float64   `json:"weekend_rate_child2"`
	WeekendRateChild3    float64   `json:"weekend_rate_child3"`
	WeekendRateChild4    float64   `json:"weekend_rate_child4"`
	TaxAndServiceCode    string    `json:"tax_and_service_code"`
	ChargeFrequencyCode  string    `json:"charge_frequency_code" binding:"required"`
	ExtraPax             float64   `json:"extra_pax"`
	PerPax               uint8     `json:"per_pax"`
	IncludeChild         uint8     `json:"include_child"`
	Day1                 uint8     `json:"day1"`
	Day2                 uint8     `json:"day2"`
	Day3                 uint8     `json:"day3"`
	Day4                 uint8     `json:"day4"`
	Day5                 uint8     `json:"day5"`
	Day6                 uint8     `json:"day6"`
	Day7                 uint8     `json:"day7"`
	Notes                string    `json:"notes"`
	IdSort               int       `json:"id_sort"`
	IsActive             uint8     `json:"is_active"`
	CmInvCode            string    `json:"cm_inv_code"`
	CmStopSell           uint8     `json:"cm_stop_sell"`
	IsCmUpdated          uint8     `json:"is_cm_updated"`
	IsCmUpdatedInclusion uint8     `json:"is_cm_updated_inclusion"`
	CmStartDate          time.Time `json:"cm_start_date" binding:"required"`
	CmEndDate            time.Time `json:"cm_end_date" binding:"required"`
	IsSent               uint8     `json:"is_sent"`
	IsOnline             uint8     `json:"is_online"`
	CreatedAt            time.Time `json:"created_at"`
	CreatedBy            string    `json:"created_by"`
	UpdatedAt            time.Time `json:"updated_at"`
	UpdatedBy            string    `json:"updated_by"`
	Id                   uint64    `json:"id"`
}

type Cfg_init_room_rate_breakdown struct {
	RoomRateCode        string    `json:"room_rate_code" binding:"required"`
	OutletCode          string    `json:"outlet_code"`
	ProductCode         string    `json:"product_code"`
	SubDepartmentCode   string    `json:"sub_department_code" binding:"required"`
	AccountCode         string    `json:"account_code" binding:"required"`
	CompanyCode         string    `json:"company_code"`
	Quantity            float64   `json:"quantity" binding:"required"`
	IsAmountPercent     uint8     `json:"is_amount_percent"`
	Amount              float64   `json:"amount" binding:"required"`
	PerPax              uint8     `json:"per_pax"`
	IncludeChild        uint8     `json:"include_child"`
	Remark              string    `json:"remark"`
	TaxAndServiceCode   string    `json:"tax_and_service_code"`
	ChargeFrequencyCode string    `json:"charge_frequency_code" binding:"required"`
	MaxPax              int       `json:"max_pax"`
	ExtraPax            float64   `json:"extra_pax"`
	PerPaxExtra         uint8     `json:"per_pax_extra"`
	CreatedAt           time.Time `json:"created_at"`
	CreatedBy           string    `json:"created_by"`
	UpdatedAt           time.Time `json:"updated_at"`
	UpdatedBy           string    `json:"updated_by"`
	Id                  uint64    `json:"id"`
}

type Cfg_init_room_rate_business_source struct {
	RoomRateCode       string    `json:"room_rate_code" binding:"required"`
	CompanyCode        string    `json:"company_code" binding:"required"`
	CommissionTypeCode string    `json:"commission_type_code" binding:"required"`
	CommissionValue    float64   `json:"commission_value" binding:"required"`
	CreatedAt          time.Time `json:"created_at"`
	CreatedBy          string    `json:"created_by"`
	UpdatedAt          time.Time `json:"updated_at"`
	UpdatedBy          string    `json:"updated_by"`
	Id                 uint64    `json:"id"`
}

type Cfg_init_room_rate_category struct {
	Code      string    `json:"code" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	IdSort    int       `json:"id_sort"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Cfg_init_room_rate_competitor struct {
	RoomRateCode   string    `json:"room_rate_code" binding:"required"`
	Name           string    `json:"name" binding:"required"`
	CompetitorCode string    `json:"competitor_code" binding:"required"`
	RateAmount     float64   `json:"rate_amount" binding:"required"`
	CreatedAt      time.Time `json:"created_at"`
	CreatedBy      string    `json:"created_by"`
	UpdatedAt      time.Time `json:"updated_at"`
	UpdatedBy      string    `json:"updated_by"`
	Id             uint64    `json:"id"`
}

type Cfg_init_room_rate_currency struct {
	RoomRateCode string    `json:"room_rate_code" binding:"required"`
	CurrencyCode string    `json:"currency_code" binding:"required"`
	ExchangeRate float64   `json:"exchange_rate" binding:"required"`
	CreatedAt    time.Time `json:"created_at"`
	CreatedBy    string    `json:"created_by"`
	UpdatedAt    time.Time `json:"updated_at"`
	UpdatedBy    string    `json:"updated_by"`
	Id           uint64    `json:"id"`
}

type Cfg_init_room_rate_dynamic struct {
	RoomRateCode string    `json:"room_rate_code" binding:"required"`
	Name         string    `json:"name" binding:"required"`
	OccFrom      int       `json:"occ_from" binding:"required"`
	OccTo        int       `json:"occ_to" binding:"required"`
	IsPercent    uint8     `json:"is_percent"`
	IsIncrease   uint8     `json:"is_increase"`
	Amount       float64   `json:"amount"`
	CreatedAt    time.Time `json:"created_at"`
	CreatedBy    string    `json:"created_by"`
	UpdatedAt    time.Time `json:"updated_at"`
	UpdatedBy    string    `json:"updated_by"`
	Id           uint64    `json:"id"`
}

type Cfg_init_room_rate_last_deal struct {
	RoomRateCode string    `json:"room_rate_code" binding:"required"`
	StartTime    string    `json:"start_time" binding:"required"`
	EndTime      string    `json:"end_time" binding:"required"`
	Percentage   float64   `json:"percentage" binding:"required"`
	CreatedAt    time.Time `json:"created_at"`
	CreatedBy    string    `json:"created_by"`
	UpdatedAt    time.Time `json:"updated_at"`
	UpdatedBy    string    `json:"updated_by"`
	Id           uint64    `json:"id"`
}

type Cfg_init_room_rate_scale struct {
	RoomRateCode    string    `json:"room_rate_code" binding:"required"`
	FromDate        time.Time `json:"from_date" binding:"required"`
	ToDate          time.Time `json:"to_date" binding:"required"`
	ScalePercentage float64   `json:"scale_percentage" binding:"required"`
	CreatedAt       time.Time `json:"created_at"`
	CreatedBy       string    `json:"created_by"`
	UpdatedAt       time.Time `json:"updated_at"`
	UpdatedBy       string    `json:"updated_by"`
	Id              uint64    `json:"id"`
}

type Cfg_init_room_rate_session struct {
	RoomRateCode string    `json:"room_rate_code" binding:"required"`
	FromDate     time.Time `json:"from_date" binding:"required"`
	ToDate       time.Time `json:"to_date" binding:"required"`
	Amount       float64   `json:"amount" binding:"required"`
	CreatedAt    time.Time `json:"created_at"`
	CreatedBy    string    `json:"created_by"`
	UpdatedAt    time.Time `json:"updated_at"`
	UpdatedBy    string    `json:"updated_by"`
	Id           uint64    `json:"id"`
}

type Cfg_init_room_rate_sub_category struct {
	Code         string    `json:"code" binding:"required"`
	Name         string    `json:"name" binding:"required"`
	CategoryCode string    `json:"category_code" binding:"required"`
	IdSort       int       `json:"id_sort" binding:"required"`
	CreatedAt    time.Time `json:"created_at"`
	CreatedBy    string    `json:"created_by"`
	UpdatedAt    time.Time `json:"updated_at"`
	UpdatedBy    string    `json:"updated_by"`
	Id           uint64    `json:"id"`
}

type Cfg_init_room_rate_weekly struct {
	RoomRateCode string    `json:"room_rate_code" binding:"required"`
	Day1         float64   `json:"day1" binding:"required"`
	Day2         float64   `json:"day2" binding:"required"`
	Day3         float64   `json:"day3" binding:"required"`
	Day4         float64   `json:"day4" binding:"required"`
	Day5         float64   `json:"day5" binding:"required"`
	Day6         float64   `json:"day6" binding:"required"`
	Day7         float64   `json:"day7" binding:"required"`
	CreatedAt    time.Time `json:"created_at"`
	CreatedBy    string    `json:"created_by"`
	UpdatedAt    time.Time `json:"updated_at"`
	UpdatedBy    string    `json:"updated_by"`
	Id           uint64    `json:"id"`
}

type Cfg_init_room_type struct {
	Code                  string    `json:"code" binding:"required"`
	Name                  string    `json:"name" binding:"required"`
	CmPercentAvailability int       `json:"cm_percent_availability"`
	MinRoomLeft           int       `json:"min_room_left"`
	IsCmFromGlobal        uint8     `json:"is_cm_from_global"`
	IsSent                uint8     `json:"is_sent"`
	IdSort                int       `json:"id_sort"`
	CreatedAt             time.Time `json:"created_at"`
	CreatedBy             string    `json:"created_by"`
	UpdatedAt             time.Time `json:"updated_at"`
	UpdatedBy             string    `json:"updated_by"`
	Id                    uint64    `json:"id"`
}

type Cfg_init_room_unavailable_reason struct {
	Code        string    `json:"code" binding:"required"`
	Description string    `json:"description" binding:"required"`
	IdSort      int       `json:"id_sort"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   string    `json:"created_by"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedBy   string    `json:"updated_by"`
	Id          uint64    `json:"id"`
}

type Cfg_init_room_view struct {
	Code      string    `json:"code" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	IdSort    int       `json:"id_sort"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Cfg_init_sales struct {
	Code        string    `json:"code" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	WaNumber    string    `json:"wa_number"`
	IdSort      int       `json:"id_sort"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   string    `json:"created_by"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedBy   string    `json:"updated_by"`
	Id          uint64    `json:"id"`
}

type Cfg_init_state struct {
	Code        string    `json:"code" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	CountryCode string    `json:"country_code" binding:"required"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   string    `json:"created_by"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedBy   string    `json:"updated_by"`
	Id          uint64    `json:"id"`
}

type Cfg_init_sub_department struct {
	Code           string    `json:"code" binding:"required"`
	Name           string    `json:"name" binding:"required"`
	DepartmentCode string    `json:"department_code" binding:"required"`
	GroupNumber    int       `json:"group_number"`
	GroupName      string    `json:"group_name"`
	IsCompliment   uint8     `json:"is_compliment"`
	IdSort         int       `json:"id_sort"`
	CreatedAt      time.Time `json:"created_at"`
	CreatedBy      string    `json:"created_by"`
	UpdatedAt      time.Time `json:"updated_at"`
	UpdatedBy      string    `json:"updated_by"`
	Id             uint64    `json:"id"`
}

type Cfg_init_tax_and_service struct {
	Code             string    `json:"code" gorm:"primaryKey" binding:"required"`
	Name             string    `json:"name" binding:"required"`
	Tax              float64   `json:"tax"`
	Service          float64   `json:"service"`
	ServiceTax       float64   `json:"service_tax"`
	IsTaxInclude     uint8     `json:"is_tax_include"`
	IsServiceInclude uint8     `json:"is_service_include"`
	CreatedAt        time.Time `json:"created_at"`
	CreatedBy        string    `json:"created_by"`
	UpdatedAt        time.Time `json:"updated_at"`
	UpdatedBy        string    `json:"updated_by"`
	Id               uint64    `json:"id"`
}

type Cfg_init_title struct {
	Code      string    `json:"code" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Cfg_init_voucher_reason struct {
	Code      string    `json:"code" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Cm_notification struct {
	Number            string    `json:"number" binding:"required"`
	ReservationNumber uint64    `json:"reservation_number" binding:"required"`
	AuditDate         time.Time `json:"audit_date" binding:"required"`
	PostingDate       time.Time `json:"posting_date" binding:"required"`
	NotifType         string    `json:"notif_type" binding:"required"`
	IsSent            uint8     `json:"is_sent" binding:"required"`
	Id                uint64    `json:"id"`
}

type Cm_update struct {
	Id           uint64    `json:"id"`
	TypeCode     string    `json:"type_code" binding:"required"`
	Number       uint64    `json:"number" binding:"required"`
	RoomTypeCode string    `json:"room_type_code" binding:"required"`
	BedTypeCode  string    `json:"bed_type_code" binding:"required"`
	RoomRateCode string    `json:"room_rate_code"`
	RateAmount   float64   `json:"rate_amount"`
	StartDate    time.Time `json:"start_date" binding:"required"`
	EndDate      time.Time `json:"end_date" binding:"required"`
	PostingDate  time.Time `json:"posting_date" binding:"required"`
	IsUpdated    uint8     `json:"is_updated" binding:"required"`
}

type Company struct {
	Code             string    `json:"code" binding:"required"`
	Name             string    `json:"name" binding:"required"`
	TypeCode         string    `json:"type_code" binding:"required"`
	SalesCode        string    `json:"sales_code"`
	ContactPerson    string    `json:"contact_person"`
	Street           string    `json:"street"`
	City             string    `json:"city"`
	CountryCode      string    `json:"country_code"`
	StateCode        string    `json:"state_code"`
	PostalCode       string    `json:"postal_code"`
	Phone1           string    `json:"phone1"`
	Phone2           string    `json:"phone2"`
	Fax              string    `json:"fax"`
	Email            string    `json:"email"`
	Website          string    `json:"website"`
	Birthday         Date      `json:"birthday"`
	ApLimit          float64   `json:"ap_limit"`
	ArLimit          float64   `json:"ar_limit"`
	IsDirectBill     uint8     `json:"is_direct_bill"`
	IsBusinessSource uint8     `json:"is_business_source"`
	InvoiceDue       int       `json:"invoice_due"`
	CreatedAt        time.Time `json:"created_at"`
	CreatedBy        string    `json:"created_by"`
	UpdatedAt        time.Time `json:"updated_at"`
	UpdatedBy        string    `json:"updated_by"`
	Id               uint64    `json:"id"`
}

type Competitor struct {
	Code         string    `json:"code" binding:"required"`
	Name         string    `json:"name" binding:"required"`
	CategoryCode string    `json:"category_code"`
	TotalRoom    int       `json:"total_room" binding:"required"`
	IdSort       int       `json:"id_sort" binding:"required"`
	CreatedAt    time.Time `json:"created_at"`
	CreatedBy    string    `json:"created_by"`
	UpdatedAt    time.Time `json:"updated_at"`
	UpdatedBy    string    `json:"updated_by"`
	Id           uint64    `json:"id"`
}

type Competitor_data struct {
	CompetitorCode  string    `json:"competitor_code" binding:"required"`
	Date            time.Time `json:"date" binding:"required"`
	AvailableRoom   uint      `json:"available_room" binding:"required"`
	RoomSold        uint      `json:"room_sold"`
	AverageRoomRate float64   `json:"average_room_rate"`
	CreatedAt       time.Time `json:"created_at"`
	CreatedBy       string    `json:"created_by"`
	UpdatedAt       time.Time `json:"updated_at"`
	UpdatedBy       string    `json:"updated_by"`
	Id              uint64    `json:"id" gorm:"primaryKey"`
}

type Configuration struct {
	SystemCode   string    `json:"system_code" binding:"required"`
	Category     string    `json:"category" binding:"required"`
	Name         string    `json:"name" binding:"required"`
	Value        string    `json:"value" binding:"required"`
	DefaultValue string    `json:"default_value" binding:"required"`
	CreatedAt    time.Time `json:"created_at"`
	CreatedBy    string    `json:"created_by"`
	UpdatedAt    time.Time `json:"updated_at"`
	UpdatedBy    string    `json:"updated_by"`
	Id           uint64    `json:"id"`
}

type Const_account_group struct {
	Code string `json:"code"  gorm:"primaryKey" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type Const_budget_type struct {
	Code string `json:"code" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type Const_channel_manager_vendor struct {
	Code   string `json:"code" binding:"required"`
	Name   string `json:"name" binding:"required"`
	IdSort int    `json:"id_sort" binding:"required"`
}

type Const_charge_frequency struct {
	Code string `json:"code" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type Const_charge_type struct {
	Code string `json:"code" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type Const_commission_type struct {
	Code string `json:"code" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type Const_customer_display_vendor struct {
	Code   string `json:"code" binding:"required"`
	Name   string `json:"name" binding:"required"`
	IdSort int    `json:"id_sort" binding:"required"`
}

type Const_department_type struct {
	Code   string `json:"code" binding:"required"`
	Name   string `json:"name" binding:"required"`
	IdSort int    `json:"id_sort" binding:"required"`
}

type Const_dynamic_rate_type struct {
	Code string `json:"code" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type Const_folio_status struct {
	Code string `json:"code" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type Const_folio_type struct {
	Code   string `json:"code" binding:"required"`
	Name   string `json:"name" binding:"required"`
	IdSort int    `json:"id_sort" binding:"required"`
	Icon   []byte `json:"icon" binding:"required"`
}

type Const_forecast_day struct {
	Id    int `json:"id"`
	Day   int `json:"day" binding:"required"`
	Month int `json:"month" binding:"required"`
}

type Const_forecast_month struct {
	Id   uint8  `json:"id"`
	Name string `json:"name"`
}

type Const_foreign_cash_table_id struct {
	Id          int    `json:"id"`
	Description string `json:"description" binding:"required"`
}

type Const_guest_status struct {
	Code   string `json:"code" binding:"required"`
	Name   string `json:"name" binding:"required"`
	IdSort string `json:"id_sort" binding:"required"`
	Icon   []byte `json:"icon" binding:"required"`
}

type Const_image struct {
	Code  string `json:"code" binding:"required"`
	Name  string `json:"name" binding:"required"`
	Image []byte `json:"image" binding:"required"`
}

type Const_iptv_vendor struct {
	Code   string `json:"code" binding:"required"`
	Name   string `json:"name" binding:"required"`
	IdSort int    `json:"id_sort" binding:"required"`
}

type Const_journal_account_group struct {
	Code string `json:"code" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type Const_journal_account_sub_group_type struct {
	Code string `json:"code" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type Const_journal_account_type struct {
	Code      string `json:"code" binding:"required"`
	Name      string `json:"name"`
	GroupCode string `json:"group_code" binding:"required"`
}

type Const_journal_prefix struct {
	Code   string `json:"code" binding:"required"`
	Name   string `json:"name"`
	IdSort int    `json:"id_sort" binding:"required"`
}

type Const_keylock_vendor struct {
	Code   string `json:"code" binding:"required"`
	Name   string `json:"name" binding:"required"`
	IdSort int    `json:"id_sort" binding:"required"`
}

type Const_member_type struct {
	Code string `json:"code" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type Const_mikrotik_vendor struct {
	Code   string `json:"code" binding:"required"`
	Name   string `json:"name" binding:"required"`
	IdSort int    `json:"id_sort" binding:"required"`
}

type Const_notification_type struct {
	Code   string `json:"code" binding:"required"`
	Name   string `json:"name" binding:"required"`
	IdSort int    `json:"id_sort" binding:"required"`
}

type Const_other_icon struct {
	Code string `json:"code" binding:"required"`
	Icon []byte `json:"icon" binding:"required"`
}

type Const_otp_status struct {
	Code string `json:"code" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type Const_pabx_rate_type struct {
	Code string `json:"code" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type Const_payment_group struct {
	Code string `json:"code" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type Const_report_font struct {
	Code int    `json:"code" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type Const_report_format struct {
	Code      int    `json:"code" binding:"required"`
	Name      string `json:"name" binding:"required"`
	IsNumeric uint8  `json:"is_numeric" binding:"required"`
}

type Const_reservation_status struct {
	Code   string `json:"code" binding:"required"`
	Name   string `json:"name" binding:"required"`
	IdSort string `json:"id_sort" binding:"required"`
	Icon   []byte `json:"icon" binding:"required"`
}

type Const_room_block_status struct {
	Code string `json:"code" binding:"required"`
	Name string `json:"name"`
	Icon []byte `json:"icon" binding:"required"`
}

type Const_room_status struct {
	Code     string `json:"code" binding:"required"`
	Name     string `json:"name" binding:"required"`
	TypeCode string `json:"type_code" binding:"required"`
	Icon     []byte `json:"icon" binding:"required"`
	IdSort   int    `json:"id_sort" binding:"required"`
}

type Const_sms_destination_type struct {
	Code string `json:"code" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type Const_sms_repeat_type struct {
	Code string `json:"code" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type Const_statistic_account struct {
	Code string `json:"code" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type Const_system struct {
	Code string `json:"code" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type Const_transaction_type struct {
	Code string `json:"code" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type Const_voucher_status struct {
	Code string `json:"code" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type Const_voucher_status_approve struct {
	Code string `json:"code" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type Const_voucher_status_sold struct {
	Code string `json:"code" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type Const_voucher_type struct {
	Code string `json:"code" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type Contact_person struct {
	TitleCode               *string    `json:"title_code"`
	FullName                *string    `json:"full_name" binding:"required"`
	Street                  *string    `json:"street"`
	CountryCode             *string    `json:"country_code"`
	StateCode               *string    `json:"state_code"`
	CityCode                *string    `json:"city_code"`
	City                    *string    `json:"city"`
	NationalityCode         *string    `json:"nationality_code"`
	PostalCode              *string    `json:"postal_code"`
	Phone1                  *string    `json:"phone1"`
	Phone2                  *string    `json:"phone2"`
	Fax                     *string    `json:"fax"`
	Email                   *string    `json:"email"`
	Website                 *string    `json:"website"`
	CompanyCode             *string    `json:"company_code"`
	GuestTypeCode           *string    `json:"guest_type_code"`
	IdCardCode              *string    `json:"id_card_code"`
	IdCardNumber            *string    `json:"id_card_number"`
	IsMale                  *uint8     `json:"is_male"`
	BirthPlace              *string    `json:"birth_place"`
	BirthDate               *time.Time `json:"birth_date"`
	TypeCode                string     `json:"type_code"`
	CustomField01           *string    `json:"custom_field01"`
	CustomField02           *string    `json:"custom_field02"`
	CustomField03           *string    `json:"custom_field03"`
	CustomField04           *string    `json:"custom_field04"`
	CustomField05           *string    `json:"custom_field05"`
	CustomField06           *string    `json:"custom_field06"`
	CustomField07           *string    `json:"custom_field07"`
	CustomField08           *string    `json:"custom_field08"`
	CustomField09           *string    `json:"custom_field09"`
	CustomField10           *string    `json:"custom_field10"`
	CustomField11           *string    `json:"custom_field11"`
	CustomField12           *string    `json:"custom_field12"`
	CustomLookupFieldCode01 *string    `json:"custom_lookup_field_code01"`
	CustomLookupFieldCode02 *string    `json:"custom_lookup_field_code02"`
	CustomLookupFieldCode03 *string    `json:"custom_lookup_field_code03"`
	CustomLookupFieldCode04 *string    `json:"custom_lookup_field_code04"`
	CustomLookupFieldCode05 *string    `json:"custom_lookup_field_code05"`
	CustomLookupFieldCode06 *string    `json:"custom_lookup_field_code06"`
	CustomLookupFieldCode07 *string    `json:"custom_lookup_field_code07"`
	CustomLookupFieldCode08 *string    `json:"custom_lookup_field_code08"`
	CustomLookupFieldCode09 *string    `json:"custom_lookup_field_code09"`
	CustomLookupFieldCode10 *string    `json:"custom_lookup_field_code10"`
	CustomLookupFieldCode11 *string    `json:"custom_lookup_field_code11"`
	CustomLookupFieldCode12 *string    `json:"custom_lookup_field_code12"`
	CreatedAt               time.Time  `json:"created_at"`
	CreatedBy               string     `json:"created_by"`
	UpdatedAt               time.Time  `json:"updated_at"`
	UpdatedBy               string     `json:"updated_by"`
	Id                      uint64     `json:"id" gorm:"primaryKey"`
}

type Cor_cfg_init_unit struct {
	Code      string    `json:"code" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	IdSort    int       `json:"id_sort"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Cor_report struct {
	Code        int       `json:"code" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	ReportQuery string    `json:"report_query" binding:"required"`
	ParentId    uint64    `json:"parent_id" binding:"required"`
	IsSystem    uint8     `json:"is_system" binding:"required"`
	IdSort      int       `json:"id_sort" binding:"required"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   string    `json:"created_by"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedBy   string    `json:"updated_by"`
	Id          uint64    `json:"id"`
}

type Cor_report_default_field struct {
	ReportCode int    `json:"report_code" binding:"required"`
	FieldQuery string `json:"field_query" binding:"required"`
	FieldName  string `json:"field_name" binding:"required"`
	IdSort     int    `json:"id_sort" binding:"required"`
}

type Cor_report_group_field struct {
	TemplateId uint64 `json:"template_id" binding:"required"`
	FieldName  string `json:"field_name" binding:"required"`
	IdSort     int    `json:"id_sort" binding:"required"`
}

type Cor_report_grouping_field struct {
	TemplateId uint64 `json:"template_id" binding:"required"`
	FieldName  string `json:"field_name" binding:"required"`
	IdSort     int    `json:"id_sort" binding:"required"`
}

type Cor_report_order_field struct {
	TemplateId  uint64 `json:"template_id" binding:"required"`
	FieldName   string `json:"field_name" binding:"required"`
	IsAscending uint8  `json:"is_ascending" binding:"required"`
	IdSort      int    `json:"id_sort" binding:"required"`
}

type Cor_report_template struct {
	ReportCode       int       `json:"report_code" binding:"required"`
	Name             string    `json:"name" binding:"required"`
	GroupLevel       int       `json:"group_level" binding:"required"`
	HeaderRemark     string    `json:"header_remark"`
	ShowFooter       uint8     `json:"show_footer" binding:"required"`
	PaperSize        int       `json:"paper_size" binding:"required"`
	PaperWidth       float64   `json:"paper_width" binding:"required"`
	PaperHeight      float64   `json:"paper_height" binding:"required"`
	IsPortrait       uint8     `json:"is_portrait" binding:"required"`
	HeaderRowHeight  int       `json:"header_row_height" binding:"required"`
	RowHeight        int       `json:"row_height" binding:"required"`
	HorizontalBorder uint8     `json:"horizontal_border" binding:"required"`
	VerticalBorder   uint8     `json:"vertical_border" binding:"required"`
	SignName1        string    `json:"sign_name1"`
	SignPosition1    string    `json:"sign_position1"`
	SignName2        string    `json:"sign_name2"`
	SignPosition2    string    `json:"sign_position2"`
	SignName3        string    `json:"sign_name3"`
	SignPosition3    string    `json:"sign_position3"`
	SignName4        string    `json:"sign_name4"`
	SignPosition4    string    `json:"sign_position4"`
	IsDefault        uint8     `json:"is_default" binding:"required"`
	IsSystem         uint8     `json:"is_system" binding:"required"`
	IdSort           int       `json:"id_sort" binding:"required"`
	CreatedAt        time.Time `json:"created_at"`
	CreatedBy        string    `json:"created_by"`
	UpdatedAt        time.Time `json:"updated_at"`
	UpdatedBy        string    `json:"updated_by"`
	Id               uint64    `json:"id"`
}

type Cor_report_template_field struct {
	TemplateId      uint64 `json:"template_id" binding:"required"`
	FieldName       string `json:"field_name" binding:"required"`
	HeaderName      string `json:"header_name" binding:"required"`
	FooterType      int    `json:"footer_type" binding:"required"`
	FormatCode      int    `json:"format_code" binding:"required"`
	Width           int    `json:"width" binding:"required"`
	Font            int    `json:"font" binding:"required"`
	FontSize        int    `json:"font_size" binding:"required"`
	FontColor       int    `json:"font_color" binding:"required"`
	Alignment       string `json:"alignment" binding:"required"`
	HeaderFontSize  int    `json:"header_font_size" binding:"required"`
	HeaderFontColor int    `json:"header_font_color" binding:"required"`
	HeaderAlignment string `json:"header_alignment" binding:"required"`
	IdSort          int    `json:"id_sort" binding:"required"`
}

type Cor_user_group struct {
	Code                string    `json:"code" binding:"required"`
	AccessForm          string    `json:"access_form" binding:"required"`
	AccessReport        string    `json:"access_report" binding:"required"`
	AccessSpecial       string    `json:"access_special" binding:"required"`
	AccessPreviewReport string    `json:"access_preview_report" binding:"required"`
	CreatedAt           time.Time `json:"created_at"`
	CreatedBy           string    `json:"created_by"`
	UpdatedAt           time.Time `json:"updated_at"`
	UpdatedBy           string    `json:"updated_by"`
	Id                  uint64    `json:"id"`
}

type Credit_card struct {
	GuestDepositId uint64    `json:"guest_deposit_id" binding:"required"`
	SubFolioId     uint64    `json:"sub_folio_id" binding:"required"`
	CardNumber     string    `json:"card_number" binding:"required"`
	CardHolder     string    `json:"card_holder" binding:"required"`
	ValidMonth     string    `json:"valid_month" binding:"required"`
	ValidYear      string    `json:"valid_year" binding:"required"`
	CreatedAt      time.Time `json:"created_at"`
	CreatedBy      string    `json:"created_by"`
	UpdatedAt      time.Time `json:"updated_at"`
	UpdatedBy      string    `json:"updated_by"`
	Id             uint64    `json:"id"`
}

type Data_analysis struct {
	Code         int    `json:"code" binding:"required"`
	Name         string `json:"name" binding:"required"`
	SystemCode   string `json:"system_code" binding:"required"`
	CategoryCode int    `json:"category_code" binding:"required"`
	IdSort       int    `json:"id_sort" binding:"required"`
}

type Data_analysis_query struct {
	Code        int    `json:"code" binding:"required"`
	Name        string `json:"name" binding:"required"`
	QueryString string `json:"query_string" binding:"required"`
}

type Data_analysis_query_list struct {
	Id               uint64 `json:"id"`
	DataAnalysisCode int    `json:"data_analysis_code" binding:"required"`
	QueryCode        int    `json:"query_code" binding:"required"`
	SheetName        string `json:"sheet_name" binding:"required"`
}

type Events struct {
	ID                int       `json:"ID" binding:"required"`
	ParentID          int       `json:"ParentID"`
	Type              uint8     `json:"Type"`
	Start             time.Time `json:"Start"`
	Finish            time.Time `json:"Finish"`
	Options           uint8     `json:"Options"`
	Caption           string    `json:"Caption"`
	RecurrenceIndex   int       `json:"RecurrenceIndex"`
	RecurrenceInfo    string    `json:"RecurrenceInfo"`
	ResourceID        string    `json:"ResourceID"`
	Location          string    `json:"Location"`
	Message           string    `json:"Message"`
	ReminderDate      time.Time `json:"ReminderDate"`
	ReminderMinutes   int       `json:"ReminderMinutes"`
	State             int       `json:"State"`
	LabelColor        int       `json:"LabelColor"`
	ActualStart       time.Time `json:"ActualStart"`
	ActualFinish      time.Time `json:"ActualFinish"`
	SyncID            string    `json:"SyncID"`
	SportID           int       `json:"SportID"`
	ReminderResources string    `json:"ReminderResources"`
	TaskComplete      int       `json:"task_complete"`
}

type Fa_cfg_init_item struct {
	Code         string    `json:"code" binding:"required"`
	Name         string    `json:"name" binding:"required"`
	CategoryCode string    `json:"category_code" binding:"required"`
	UomCode      string    `json:"uom_code" binding:"required"`
	Remark       string    `json:"remark"`
	CreatedAt    time.Time `json:"created_at"`
	CreatedBy    string    `json:"created_by"`
	UpdatedAt    time.Time `json:"updated_at"`
	UpdatedBy    string    `json:"updated_by"`
	Id           uint64    `json:"id"`
}

type Fa_cfg_init_item_category struct {
	Code                           string    `json:"code" binding:"required"`
	Name                           string    `json:"name" binding:"required"`
	JournalAccountCode             string    `json:"journal_account_code" binding:"required"`
	JournalAccountCodeCogs         string    `json:"journal_account_code_cogs" binding:"required"`
	JournalAccountCodeExpense      string    `json:"journal_account_code_expense" binding:"required"`
	JournalAccountCodeSell         string    `json:"journal_account_code_sell" binding:"required"`
	JournalAccountCodeDepreciation string    `json:"journal_account_code_depreciation" binding:"required"`
	JournalAccountCodeSpoil        string    `json:"journal_account_code_spoil" binding:"required"`
	IsLinen                        uint8     `json:"is_linen"`
	IsIntangible                   uint8     `json:"is_intangible"`
	CreatedAt                      time.Time `json:"created_at"`
	CreatedBy                      string    `json:"created_by"`
	UpdatedAt                      time.Time `json:"updated_at"`
	UpdatedBy                      string    `json:"updated_by"`
	Id                             uint64    `json:"id"`
}

type Fa_cfg_init_location struct {
	Code             string    `json:"code" binding:"required"`
	Name             string    `json:"name" binding:"required"`
	LocationTypeCode string    `json:"location_type_code"`
	CreatedAt        time.Time `json:"created_at"`
	CreatedBy        string    `json:"created_by"`
	UpdatedAt        time.Time `json:"updated_at"`
	UpdatedBy        string    `json:"updated_by"`
	Id               uint64    `json:"id"`
}

type Cm_log struct {
	Id         uint64    `json:"id"`
	BookingId  string    `json:"booking_id" binding:"required"`
	RevisionId string    `json:"revision_id" binding:"required"`
	Message    string    `json:"message"`
	CreatedAt  time.Time `json:"created_at"`
	CreatedBy  string    `json:"created_by"`
}

type Fa_cfg_init_manufacture struct {
	Code      string    `json:"code" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Fa_const_depreciation_type struct {
	Code   string `json:"code" binding:"required"`
	Name   string `json:"name" binding:"required"`
	IdSort int    `json:"id_sort" binding:"required"`
}

type Fa_const_item_condition struct {
	Code string `json:"code" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type Fa_const_location_type struct {
	Code string `json:"code" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type Fa_depreciation struct {
	Number    string    `json:"number" binding:"required"`
	RefNumber string    `json:"ref_number" binding:"required"`
	Month     int       `json:"month" binding:"required"`
	Year      int       `json:"year" binding:"required"`
	Date      time.Time `json:"date" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Fa_list struct {
	Code                           string     `json:"code" binding:"required"`
	Barcode                        *string    `json:"barcode"`
	ReceiveNumber                  string     `json:"receive_number" binding:"required"`
	ReceiveId                      uint64     `json:"receive_id" binding:"required"`
	ItemCode                       string     `json:"item_code" binding:"required"`
	SortNumber                     *uint64    `json:"sort_number" binding:"required"`
	Name                           string     `json:"name" binding:"required"`
	AcquisitionDate                time.Time  `json:"acquisition_date" binding:"required"`
	DepreciationDate               time.Time  `json:"depreciation_date" binding:"required"`
	DepreciationTypeCode           string     `json:"depreciation_type_code" binding:"required"`
	DepreciationSubDepartmentCode  *string    `json:"depreciation_sub_department_code"`
	DepreciationExpenseAccountCode *string    `json:"depreciation_expense_account_code"`
	PurchasePrice                  float64    `json:"purchase_price" binding:"required"`
	CurrentValue                   float64    `json:"current_value" binding:"required"`
	ResidualValue                  *float64   `json:"residual_value" binding:"required"`
	SerialNumber                   *string    `json:"serial_number"`
	ManufactureCode                *string    `json:"manufacture_code"`
	Trademark                      *string    `json:"trademark"`
	WarrantyDate                   *time.Time `json:"warranty_date"`
	LocationCode                   *string    `json:"location_code"`
	UsefulLife                     *int       `json:"useful_life" binding:"required"`
	ConditionCode                  string     `json:"condition_code" binding:"required"`
	Remark                         *string    `json:"remark"`
	DepreciationRate               *float64   `json:"depreciation_rate"`
	FoNumber                       string     `json:"fo_number"`
	RefNumber1                     string     `json:"ref_number1"`
	DoNotRevenueJournal            *uint8     `json:"do_not_revenue_journal"`
	RefNumber2                     string     `json:"ref_number2"`
	IsOldAsset                     *uint8     `json:"is_old_asset"`
	DepreciatedMonth               *int       `json:"depreciated_month"`
	DepreciatedValue               *float64   `json:"depreciated_value"`
	CreatedAt                      time.Time  `json:"created_at"`
	CreatedBy                      string     `json:"created_by"`
	UpdatedAt                      time.Time  `json:"updated_at"`
	UpdatedBy                      string     `json:"updated_by"`
	Id                             uint64     `json:"id"`
}

type Fa_location_history struct {
	FaCode           string    `json:"fa_code" binding:"required"`
	FromLocationCode string    `json:"from_location_code" binding:"required"`
	ToLocationCode   string    `json:"to_location_code" binding:"required"`
	PostingDate      time.Time `json:"posting_date" binding:"required"`
	Remark           string    `json:"remark"`
	CreatedAt        time.Time `json:"created_at"`
	CreatedBy        string    `json:"created_by"`
	UpdatedAt        time.Time `json:"updated_at"`
	UpdatedBy        string    `json:"updated_by"`
	Id               uint64    `json:"id"`
}

type Fa_purchase_order struct {
	Number                  string    `json:"number" binding:"required"`
	CompanyCode             string    `json:"company_code" binding:"required"`
	ExpeditionCode          *string   `json:"expedition_code"`
	ContactPersonId         uint64    `json:"contact_person_id" binding:"required"`
	ShippingCompany         string    `json:"shipping_company"`
	ContactPersonShippingId uint64    `json:"contact_person_shipping_id" binding:"required"`
	Date                    time.Time `json:"date" binding:"required"`
	RequestBy               string    `json:"request_by" binding:"required"`
	Remark                  *string   `json:"remark"`
	IsReceived              uint8     `json:"is_received" binding:"required"`
	CreatedAt               time.Time `json:"created_at"`
	CreatedBy               string    `json:"created_by"`
	UpdatedAt               time.Time `json:"updated_at"`
	UpdatedBy               string    `json:"updated_by"`
	Id                      uint64    `json:"id"`
}

type Fa_purchase_order_detail struct {
	PoNumber  string    `json:"po_number" binding:"required"`
	ItemCode  string    `json:"item_code" binding:"required"`
	Quantity  int       `json:"quantity" binding:"required"`
	UomCode   string    `json:"uom_code" binding:"required"`
	Price     float64   `json:"price" binding:"required"`
	Remark    *string   `json:"remark"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Fa_receive struct {
	Number            string    `json:"number" binding:"required"`
	RefNumber         string    `json:"ref_number" binding:"required"`
	RefNumberOldAsset string    `json:"ref_number_old_asset"`
	PoNumber          string    `json:"po_number" binding:"required"`
	ApNumber          string    `json:"ap_number"`
	CompanyCode       string    `json:"company_code" binding:"required"`
	InvoiceNumber     string    `json:"invoice_number"`
	BankAccountCode   *string   `json:"bank_account_code"`
	AmountPayment     *float64  `json:"amount_payment"`
	Date              time.Time `json:"date" binding:"required"`
	Remark            *string   `json:"remark"`
	IsSeparate        *uint8    `json:"is_separate" binding:"required"`
	IsDiscountIncome  *uint8    `json:"is_discount_income" binding:"required"`
	IsTaxExpense      *uint8    `json:"is_tax_expense" binding:"required"`
	IsShippingExpense *uint8    `json:"is_shipping_expense" binding:"required"`
	IsCredit          *uint8    `json:"is_credit" binding:"required"`
	DueDate           time.Time `json:"due_date"`
	IsPaid            *uint8    `json:"is_paid" binding:"required"`
	CreatedAt         time.Time `json:"created_at"`
	CreatedBy         string    `json:"created_by"`
	UpdatedAt         time.Time `json:"updated_at"`
	UpdatedBy         string    `json:"updated_by"`
	Id                uint64    `json:"id"`
}

type Fa_receive_detail struct {
	ReceiveNumber    string    `json:"receive_number" binding:"required"`
	ItemCode         string    `json:"item_code" binding:"required"`
	DetailName       string    `json:"detail_name" binding:"required"`
	PoQuantity       int       `json:"po_quantity" binding:"required"`
	PoPrice          float64   `json:"po_price" binding:"required"`
	ReceiveQuantity  int       `json:"receive_quantity" binding:"required"`
	ReceiveUomCode   string    `json:"receive_uom_code" binding:"required"`
	ReceivePrice     float64   `json:"receive_price" binding:"required"`
	Quantity         int       `json:"quantity" binding:"required"`
	TotalPrice       float64   `json:"total_price" binding:"required"`
	Discount         float64   `json:"discount"`
	Tax              float64   `json:"tax"`
	Shipping         float64   `json:"shipping"`
	IsOldAsset       *uint8    `json:"is_old_asset"`
	DepreciatedMonth int       `json:"depreciated_month"`
	DepreciatedValue float64   `json:"depreciated_value"`
	Remark           string    `json:"remark"`
	CreatedAt        time.Time `json:"created_at"`
	CreatedBy        string    `json:"created_by"`
	UpdatedAt        time.Time `json:"updated_at"`
	UpdatedBy        string    `json:"updated_by"`
	Id               uint64    `json:"id" gorm:"primaryKey"`
}

type Fa_repair struct {
	FaCode        string    `json:"fa_code" binding:"required"`
	Date          time.Time `json:"date" binding:"required"`
	TypeCode      string    `json:"type_code" binding:"required"`
	LaborCost     float64   `json:"labor_cost"`
	SparepartCost float64   `json:"sparepart_cost"`
	EstimatedDate time.Time `json:"estimated_date"`
	CreatedAt     time.Time `json:"created_at"`
	CreatedBy     string    `json:"created_by"`
	UpdatedAt     time.Time `json:"updated_at"`
	UpdatedBy     string    `json:"updated_by"`
	Id            uint64    `json:"id"`
}

type Fa_revaluation struct {
	Number         string    `json:"number" binding:"required"`
	FaCode         string    `json:"fa_code" binding:"required"`
	RefNumber      string    `json:"ref_number" binding:"required"`
	Date           time.Time `json:"date" binding:"required"`
	Amount         float64   `json:"amount" binding:"required"`
	IsIncrease     uint8     `json:"is_increase" binding:"required"`
	IsDepreciation uint8     `json:"is_depreciation" binding:"required"`
	CreatedAt      time.Time `json:"created_at"`
	CreatedBy      string    `json:"created_by"`
	UpdatedAt      time.Time `json:"updated_at"`
	UpdatedBy      string    `json:"updated_by"`
	Id             uint64    `json:"id"`
}

type Fb_statistic struct {
	Date          time.Time `json:"date" binding:"required"`
	OutletCode    string    `json:"outlet_code" binding:"required"`
	Adult         int       `json:"adult" binding:"required"`
	Child         int       `json:"child" binding:"required"`
	AdultBeverage int       `json:"adult_beverage" binding:"required"`
	ChildBeverage int       `json:"child_beverage" binding:"required"`
	FoodNett      float64   `json:"food_nett" binding:"required"`
	BeverageNett  float64   `json:"beverage_nett" binding:"required"`
	Id            uint64    `json:"id"`
}

type Folio struct {
	Number            uint64    `json:"number"  gorm:"primaryKey"`
	TypeCode          string    `json:"type_code"`
	CoNumber          string    `json:"co_number"`
	ReservationNumber uint64    `json:"reservation_number"`
	ContactPersonId1  uint64    `json:"contact_person_id1"`
	ContactPersonId2  uint64    `json:"contact_person_id2"`
	ContactPersonId3  uint64    `json:"contact_person_id3"`
	ContactPersonId4  uint64    `json:"contact_person_id4"`
	GuestDetailId     uint64    `json:"guest_detail_id"`
	GuestProfileId1   uint64    `json:"guest_profile_id1"`
	GuestProfileId2   uint64    `json:"guest_profile_id2"`
	GuestProfileId3   uint64    `json:"guest_profile_id3"`
	GuestProfileId4   uint64    `json:"guest_profile_id4"`
	GuestGeneralId    uint64    `json:"guest_general_id"`
	GroupCode         string    `json:"group_code"`
	RoomStatusCode    string    `json:"room_status_code"`
	StatusCode        string    `json:"status_code" binding:"required"`
	IsIncognito       uint8     `json:"is_incognito"`
	VoucherNumber     string    `json:"voucher_number"`
	ComplimentHu      string    `json:"compliment_hu"`
	IsCtlApproved     uint8     `json:"is_ctl_approved"`
	IsLock            uint8     `json:"is_lock"`
	IsPrinted         uint8     `json:"is_printed"`
	IsFromAllotment   uint8     `json:"is_from_allotment"`
	CheckOutAt        time.Time `json:"check_out_at"`
	CheckOutBy        string    `json:"check_out_by"`
	CancelledAt       time.Time `json:"cancelled_at"`
	CancelledBy       string    `json:"cancelled_by"`
	CancelReason      string    `json:"cancel_reason"`
	SystemCode        string    `json:"system_code"`
	CreatedAt         time.Time `json:"created_at"`
	CreatedBy         string    `json:"created_by"`
	UpdatedAt         time.Time `json:"updated_at"`
	UpdatedBy         string    `json:"updated_by"`
}

type Folio_routing struct {
	FolioNumber      uint64    `json:"folio_number" binding:"required"`
	AccountCode      string    `json:"account_code" binding:"required"`
	FolioTransfer    uint64    `json:"folio_transfer" binding:"required"`
	SubFolioTransfer string    `json:"sub_folio_transfer" binding:"required"`
	CreatedAt        time.Time `json:"created_at"`
	CreatedBy        string    `json:"created_by"`
	UpdatedAt        time.Time `json:"updated_at"`
	UpdatedBy        string    `json:"updated_by"`
	Id               uint64    `json:"id"`
}

type Forecast_in_house_change_pax struct {
	AuditDate     time.Time `json:"audit_date" binding:"required"`
	FolioNumber   uint64    `json:"folio_number" binding:"required"`
	Adult         int       `json:"adult" binding:"required"`
	Child         int       `json:"child" binding:"required"`
	IsReservation uint8     `json:"is_reservation" binding:"required"`
	CreatedAt     time.Time `json:"created_at"`
	CreatedBy     string    `json:"created_by"`
	UpdatedAt     time.Time `json:"updated_at"`
	UpdatedBy     string    `json:"updated_by"`
	Id            uint64    `json:"id"`
}

type Forecast_monthly_day struct {
	Id                           uint64    `json:"id"`
	AuditDate                    time.Time `json:"audit_date" binding:"required"`
	TotalRoom                    int       `json:"total_room" binding:"required"`
	TotalOo                      int       `json:"total_oo" binding:"required"`
	RoomTypeCode                 string    `json:"room_type_code"`
	BedTypeCode                  string    `json:"bed_type_code"`
	Pax                          int       `json:"pax" binding:"required"`
	PaxGroupArrivalByRes         int       `json:"pax_group_arrival_by_res" binding:"required"`
	PaxGroupArrivalByWalk        int       `json:"pax_group_arrival_by_walk" binding:"required"`
	PaxIndividualArrivalByRes    int       `json:"pax_individual_arrival_by_res" binding:"required"`
	PaxIndividualArrivalByWalk   int       `json:"pax_individual_arrival_by_walk" binding:"required"`
	PaxGroupDeparture            int       `json:"pax_group_departure" binding:"required"`
	PaxIndividualDeparture       int       `json:"pax_individual_departure" binding:"required"`
	PaxDayUse                    int       `json:"pax_day_use" binding:"required"`
	PaxGroupCompliment           int       `json:"pax_group_compliment" binding:"required"`
	PaxIndividualCompliment      int       `json:"pax_individual_compliment" binding:"required"`
	PaxHouseUse                  int       `json:"pax_house_use" binding:"required"`
	Rooms                        int       `json:"rooms" binding:"required"`
	RoomsDayUse                  int       `json:"rooms_day_use" binding:"required"`
	RoomsGroup                   int       `json:"rooms_group" binding:"required"`
	RoomsGroupArrival            int       `json:"rooms_group_arrival" binding:"required"`
	RoomsGroupArrivalByRes       int       `json:"rooms_group_arrival_by_res" binding:"required"`
	RoomsGroupArrivalByWalk      int       `json:"rooms_group_arrival_by_walk" binding:"required"`
	RoomsGroupDeparture          int       `json:"rooms_group_departure" binding:"required"`
	RoomsGroupStay               int       `json:"rooms_group_stay" binding:"required"`
	RoomsIndividual              int       `json:"rooms_individual" binding:"required"`
	RoomsIndividualArrival       int       `json:"rooms_individual_arrival" binding:"required"`
	RoomsIndividualArrivalByRes  int       `json:"rooms_individual_arrival_by_res" binding:"required"`
	RoomsIndividualArrivalByWalk int       `json:"rooms_individual_arrival_by_walk" binding:"required"`
	RoomsIndividualDeparture     int       `json:"rooms_individual_departure" binding:"required"`
	RoomsIndividualStay          int       `json:"rooms_individual_stay" binding:"required"`
	RoomsCompliment              int       `json:"rooms_compliment" binding:"required"`
	RoomsComplimentIndividual    int       `json:"rooms_compliment_individual" binding:"required"`
	RoomsComplimentGroup         int       `json:"rooms_compliment_group" binding:"required"`
	RoomsHouseUse                int       `json:"rooms_house_use" binding:"required"`
	RoomAvailable                int       `json:"room_available" binding:"required"`
	RevenuePrediction            float64   `json:"revenue_prediction" binding:"required"`
	Revenue                      float64   `json:"revenue" binding:"required"`
	RevenueNett                  float64   `json:"revenue_nett" binding:"required"`
	RevenueBreakfast             float64   `json:"revenue_breakfast"`
	RevenueBreakfastNett         float64   `json:"revenue_breakfast_nett"`
	IpAddress                    string    `json:"ip_address" binding:"required"`
}

type Forecast_monthly_day_previous struct {
	Id                    uint64 `json:"id"`
	RoomGroup             int    `json:"room_group" binding:"required"`
	RoomNonGroup          int    `json:"room_non_group" binding:"required"`
	RoomDepartureGroup    int    `json:"room_departure_group" binding:"required"`
	RoomDepartureNonGroup int    `json:"room_departure_non_group" binding:"required"`
	RoomDeparture         int    `json:"room_departure" binding:"required"`
	Rooms                 int    `json:"rooms" binding:"required"`
	Adult                 int    `json:"adult" binding:"required"`
	Child                 int    `json:"child" binding:"required"`
	DayUseRoomGroup       int    `json:"day_use_room_group" binding:"required"`
	DayUseRoomNonGroup    int    `json:"day_use_room_non_group" binding:"required"`
	DayUseRooms           int    `json:"day_use_rooms" binding:"required"`
	DayUseAdult           int    `json:"day_use_adult" binding:"required"`
	DayUseChild           int    `json:"day_use_child" binding:"required"`
	IpAddress             string `json:"ip_address" binding:"required"`
}

type Grid_properties struct {
	SystemCode string    `json:"system_code" binding:"required"`
	GridName   string    `json:"grid_name" binding:"required"`
	Properties []byte    `json:"properties" binding:"required"`
	CreatedAt  time.Time `json:"created_at"`
	CreatedBy  string    `json:"created_by"`
	UpdatedAt  time.Time `json:"updated_at"`
	UpdatedBy  string    `json:"updated_by"`
	Id         uint64    `json:"id"`
}

type Guest_breakdown struct {
	FolioNumber         uint64    `json:"folio_number" binding:"required"`
	OutletCode          string    `json:"outlet_code"`
	ProductCode         string    `json:"product_code"`
	SubDepartmentCode   string    `json:"sub_department_code" binding:"required"`
	AccountCode         string    `json:"account_code" binding:"required"`
	CompanyCode         string    `json:"company_code"`
	Quantity            float64   `json:"quantity" binding:"required"`
	IsAmountPercent     uint8     `json:"is_amount_percent" binding:"required"`
	Amount              float64   `json:"amount" binding:"required"`
	PerPax              uint8     `json:"per_pax" binding:"required"`
	IncludeChild        uint8     `json:"include_child" binding:"required"`
	Remark              string    `json:"remark"`
	TaxAndServiceCode   string    `json:"tax_and_service_code"`
	ChargeFrequencyCode string    `json:"charge_frequency_code" binding:"required"`
	MaxPax              int       `json:"max_pax"`
	ExtraPax            float64   `json:"extra_pax"`
	PerPaxExtra         uint8     `json:"per_pax_extra" binding:"required"`
	CreatedAt           time.Time `json:"created_at"`
	CreatedBy           string    `json:"created_by"`
	UpdatedAt           time.Time `json:"updated_at"`
	UpdatedBy           string    `json:"updated_by"`
	Id                  uint64    `json:"id"`
}

type Guest_deposit struct {
	ReservationNumber   uint64    `json:"reservation_number" binding:"required"`
	SubDepartmentCode   string    `json:"sub_department_code" binding:"required"`
	AccountCode         string    `json:"account_code" binding:"required"`
	Amount              float64   `json:"amount" binding:"required"`
	DefaultCurrencyCode string    `json:"default_currency_code"`
	AmountForeign       float64   `json:"amount_foreign" binding:"required"`
	ExchangeRate        float64   `json:"exchange_rate" binding:"required"`
	CurrencyCode        string    `json:"currency_code" binding:"required"`
	AuditDate           time.Time `json:"audit_date"`
	Remark              string    `json:"remark"`
	DocumentNumber      string    `json:"document_number"`
	TypeCode            string    `json:"type_code"`
	CardBankCode        string    `json:"card_bank_code"`
	CardTypeCode        string    `json:"card_type_code"`
	RefNumber           uint64    `json:"ref_number"`
	Void                uint8     `json:"void"`
	VoidDate            time.Time `json:"void_date"`
	VoidBy              string    `json:"void_by"`
	VoidReason          string    `json:"void_reason"`
	IsCorrection        uint8     `json:"is_correction"`
	CorrectionBy        string    `json:"correction_by"`
	CorrectionReason    string    `json:"correction_reason"`
	CorrectionBreakdown uint64    `json:"correction_breakdown" `
	Shift               string    `json:"shift" `
	LogShiftId          uint64    `json:"log_shift_id"  `
	IsPairWithFolio     uint8     `json:"is_pair_with_folio"  `
	TransferPairId      uint64    `json:"transfer_pair_id"`
	SystemCode          string    `json:"system_code" `
	CreatedAt           time.Time `json:"created_at"`
	CreatedBy           string    `json:"created_by"`
	UpdatedAt           time.Time `json:"updated_at"`
	UpdatedBy           string    `json:"updated_by"`
	Id                  uint64    `json:"id" gorm:"primaryKey"`
}

type Guest_detail struct {
	Arrival              time.Time `json:"arrival" binding:"required"`
	ArrivalUnixx         int64     `json:"arrival_unixx"`
	ArrivalRes           time.Time `json:"arrival_res" `
	Departure            time.Time `json:"departure" binding:"required"`
	DepartureUnixx       int64     `json:"departure_unixx"`
	DepartureRes         time.Time `json:"departure_res"`
	Adult                int       `json:"adult" binding:"required"`
	Child                *int      `json:"child"`
	RoomTypeCode         string    `json:"room_type_code"`
	BedTypeCode          string    `json:"bed_type_code"`
	RoomNumber           *string   `json:"room_number"`
	CurrencyCode         string    `json:"currency_code" binding:"required"`
	ExchangeRate         float64   `json:"exchange_rate"`
	IsConstantCurrency   uint8     `json:"is_constant_currency"`
	RoomRateCode         string    `json:"room_rate_code"`
	IsOverrideRate       *uint8    `json:"is_override_rate"`
	WeekdayRate          *float64  `json:"weekday_rate"`
	WeekendRate          *float64  `json:"weekend_rate"`
	DiscountPercent      *uint8    `json:"discount_percent"`
	Discount             *float64  `json:"discount"`
	BusinessSourceCode   *string   `json:"business_source_code"`
	IsOverrideCommission *uint8    `json:"is_override_commission"`
	CommissionTypeCode   *string   `json:"commission_type_code"`
	CommissionValue      *float64  `json:"commission_value"`
	PaymentTypeCode      string    `json:"payment_type_code" binding:"required"`
	MarketCode           *string   `json:"market_code"`
	BookingSourceCode    *string   `json:"booking_source_code"`
	BillInstruction      *string   `json:"bill_instruction"`
	CreatedAt            time.Time `json:"created_at"`
	CreatedBy            string    `json:"created_by"`
	UpdatedAt            time.Time `json:"updated_at"`
	UpdatedBy            string    `json:"updated_by"`
	Id                   uint64    `json:"id" gorm:"primaryKey"`
}

type Guest_extra_charge struct {
	FolioNumber         uint64    `json:"folio_number" binding:"required"`
	PackageName         *string   `json:"package_name" binding:"required"`
	OutletCode          *string   `json:"outlet_code"`
	ProductCode         *string   `json:"product_code"`
	PackageCode         *string   `json:"package_code"`
	GroupCode           string    `json:"group_code" binding:"required"`
	SubDepartmentCode   string    `json:"sub_department_code" binding:"required"`
	AccountCode         string    `json:"account_code" binding:"required"`
	Quantity            float64   `json:"quantity" binding:"required"`
	Amount              float64   `json:"amount" binding:"required"`
	PerPax              *uint8    `json:"per_pax" binding:"required"`
	IncludeChild        *uint8    `json:"include_child" binding:"required"`
	TaxAndServiceCode   *string   `json:"tax_and_service_code"`
	ChargeFrequencyCode string    `json:"charge_frequency_code" binding:"required"`
	MaxPax              int       `json:"max_pax"`
	ExtraPax            *float64  `json:"extra_pax"`
	PerPaxExtra         *uint8    `json:"per_pax_extra" binding:"required"`
	CreatedAt           time.Time `json:"created_at"`
	CreatedBy           string    `json:"created_by"`
	UpdatedAt           time.Time `json:"updated_at"`
	UpdatedBy           string    `json:"updated_by"`
	Id                  uint64    `json:"id" gorm:"primaryKey"`
}

type Guest_extra_charge_breakdown struct {
	GuestExtraChargeId  uint64    `json:"guest_extra_charge_id" binding:"required"`
	OutletCode          *string   `json:"outlet_code"`
	ProductCode         *string   `json:"product_code"`
	SubDepartmentCode   string    `json:"sub_department_code" binding:"required"`
	AccountCode         string    `json:"account_code" binding:"required"`
	CompanyCode         *string   `json:"company_code"`
	Quantity            float64   `json:"quantity" binding:"required"`
	IsAmountPercent     *uint8    `json:"is_amount_percent" binding:"required"`
	Amount              float64   `json:"amount" binding:"required"`
	PerPax              *uint8    `json:"per_pax" binding:"required"`
	IncludeChild        *uint8    `json:"include_child" binding:"required"`
	Remark              *string   `json:"remark"`
	TaxAndServiceCode   *string   `json:"tax_and_service_code"`
	ChargeFrequencyCode string    `json:"charge_frequency_code" binding:"required"`
	MaxPax              int       `json:"max_pax"`
	ExtraPax            *float64  `json:"extra_pax"`
	PerPaxExtra         *uint8    `json:"per_pax_extra" binding:"required"`
	CreatedAt           time.Time `json:"created_at"`
	CreatedBy           string    `json:"created_by"`
	UpdatedAt           time.Time `json:"updated_at"`
	UpdatedBy           string    `json:"updated_by"`
	Id                  uint64    `json:"id" gorm:"primaryKey"`
}

type Guest_general struct {
	PurposeOfCode   *string    `json:"purpose_of_code"`
	SalesCode       *string    `json:"sales_code"`
	VoucherNumberTa *string    `json:"voucher_number_ta"`
	FlightNumber    *string    `json:"flight_number"`
	FlightArrival   *time.Time `json:"flight_arrival"`
	FlightDeparture *time.Time `json:"flight_departure"`
	Notes           *string    `json:"notes"`
	ShowNotes       *uint8     `json:"show_notes"`
	HkNote          *string    `json:"hk_note"`
	DocumentNumber  *string    `json:"document_number"`
	CreatedAt       time.Time  `json:"created_at"`
	CreatedBy       string     `json:"created_by"`
	UpdatedAt       time.Time  `json:"updated_at"`
	UpdatedBy       string     `json:"updated_by"`
	Id              uint64     `json:"id" gorm:"primaryKey"`
}

type Guest_group struct {
	Code          string    `json:"code" binding:"required" gorm:"uniqueIndex"`
	Name          string    `json:"name" binding:"required"`
	ContactPerson *string   `json:"contact_person"`
	Street        *string   `json:"street"`
	CountryCode   *string   `json:"country_code"`
	StateCode     *string   `json:"state_code"`
	CityCode      *string   `json:"city_code"`
	City          *string   `json:"city"`
	PostalCode    *string   `json:"postal_code"`
	Phone1        *string   `json:"phone1"`
	Phone2        *string   `json:"phone2"`
	Fax           *string   `json:"fax"`
	Email         *string   `json:"email"`
	Website       *string   `json:"website"`
	IsActive      uint8     `json:"is_active" binding:"required"`
	CreatedAt     time.Time `json:"created_at"`
	CreatedBy     string    `json:"created_by"`
	UpdatedAt     time.Time `json:"updated_at"`
	UpdatedBy     string    `json:"updated_by"`
	Id            uint64    `json:"id" gorm:"primaryKey"`
}

type Guest_in_house struct {
	AuditDate               time.Time `json:"audit_date" binding:"required"`
	AuditDateUnixx          int64     `json:"audit_date_unixx" binding:"required"`
	FolioNumber             uint64    `json:"folio_number" binding:"required"`
	GroupCode               string    `json:"group_code"`
	Adult                   int       `json:"adult" binding:"required"`
	Child                   int       `json:"child" binding:"required"`
	RoomTypeCode            string    `json:"room_type_code" binding:"required"`
	BedTypeCode             string    `json:"bed_type_code" binding:"required"`
	RoomNumber              string    `json:"room_number"`
	RoomRateCode            string    `json:"room_rate_code" binding:"required"`
	RateOriginal            float64   `json:"rate_original" binding:"required"`
	Rate                    float64   `json:"rate" binding:"required"`
	DiscountPercent         uint8     `json:"discount_percent" binding:"required"`
	Discount                float64   `json:"discount"`
	BusinessSourceCode      string    `json:"business_source_code"`
	CommissionTypeCode      string    `json:"commission_type_code"`
	CommissionValue         float64   `json:"commission_value"`
	PaymentTypeCode         string    `json:"payment_type_code" binding:"required"`
	MarketCode              string    `json:"market_code"`
	BookingSourceCode       string    `json:"booking_source_code"`
	TitleCode               string    `json:"title_code"`
	FullName                string    `json:"full_name" binding:"required"`
	Street                  string    `json:"street"`
	CountryCode             string    `json:"country_code"`
	StateCode               string    `json:"state_code"`
	CityCode                string    `json:"city_code"`
	City                    string    `json:"city"`
	NationalityCode         string    `json:"nationality_code"`
	PostalCode              string    `json:"postal_code"`
	Phone1                  string    `json:"phone1"`
	Phone2                  string    `json:"phone2"`
	Fax                     string    `json:"fax"`
	Email                   string    `json:"email"`
	Website                 string    `json:"website"`
	CompanyCode             string    `json:"company_code"`
	GuestTypeCode           string    `json:"guest_type_code"`
	PurposeOfCode           string    `json:"purpose_of_code"`
	SalesCode               string    `json:"sales_code"`
	CustomLookupFieldCode01 string    `json:"custom_lookup_field_code01"`
	CustomLookupFieldCode02 string    `json:"custom_lookup_field_code02"`
	ComplimentHu            string    `json:"compliment_hu" binding:"required"`
	IsAdditional            uint8     `json:"is_additional" binding:"required"`
	IsScheduledRate         uint8     `json:"is_scheduled_rate" binding:"required"`
	IsBreakfast             uint8     `json:"is_breakfast" binding:"required"`
	PaxBreakfast            int       `json:"pax_breakfast"`
	BreakfastVoucherNumber  string    `json:"breakfast_voucher_number"`
	Notes                   string    `json:"notes"`
	CreatedAt               time.Time `json:"created_at"`
	CreatedBy               string    `json:"created_by"`
	UpdatedAt               time.Time `json:"updated_at"`
	UpdatedBy               string    `json:"updated_by"`
	Id                      uint64    `json:"id"`
}

type Guest_in_house_breakdown struct {
	AuditDate           time.Time `json:"audit_date" binding:"required"`
	FolioNumber         uint64    `json:"folio_number" binding:"required"`
	OutletCode          string    `json:"outlet_code"`
	ProductCode         string    `json:"product_code"`
	SubDepartmentCode   string    `json:"sub_department_code" binding:"required"`
	AccountCode         string    `json:"account_code" binding:"required"`
	Quantity            float64   `json:"quantity" binding:"required"`
	Amount              float64   `json:"amount" binding:"required"`
	PerPax              uint8     `json:"per_pax" binding:"required"`
	IncludeChild        uint8     `json:"include_child" binding:"required"`
	Remark              string    `json:"remark"`
	TaxAndServiceCode   string    `json:"tax_and_service_code"`
	ChargeFrequencyCode string    `json:"charge_frequency_code" binding:"required"`
	MaxPax              int       `json:"max_pax"`
	ExtraPax            float64   `json:"extra_pax"`
	PerPaxExtra         uint8     `json:"per_pax_extra" binding:"required"`
	CreatedAt           time.Time `json:"created_at"`
	CreatedBy           string    `json:"created_by"`
	UpdatedAt           time.Time `json:"updated_at"`
	UpdatedBy           string    `json:"updated_by"`
	Id                  uint64    `json:"id"`
}

type Guest_loan_item struct {
	ItemCode    string    `json:"item_code" binding:"required"`
	Quantity    int       `json:"quantity" binding:"required"`
	FolioNumber uint64    `json:"folio_number" binding:"required"`
	Date        time.Time `json:"date" binding:"required"`
	ReturnDate  time.Time `json:"return_date"`
	IsReturned  *uint8    `json:"is_returned"`
	Remark      *string   `json:"remark"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   string    `json:"created_by"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedBy   string    `json:"updated_by"`
	Id          uint64    `json:"id"`
}

type Guest_message struct {
	FolioNumber  uint64    `json:"folio_number" binding:"required"`
	Title        string    `json:"title" binding:"required"`
	Body         string    `json:"body"`
	PostingDate  time.Time `json:"posting_date" binding:"required"`
	ReminderDate time.Time `json:"reminder_date"`
	IsDelivered  uint8     `json:"is_delivered" binding:"required"`
	CreatedAt    time.Time `json:"created_at"`
	CreatedBy    string    `json:"created_by"`
	UpdatedAt    time.Time `json:"updated_at"`
	UpdatedBy    string    `json:"updated_by"`
	Id           uint64    `json:"id"`
}

type Guest_profile struct {
	TitleCode               string    `json:"title_code"`
	FullName                string    `json:"full_name"`
	Street                  string    `json:"street"`
	CountryCode             string    `json:"country_code"`
	StateCode               string    `json:"state_code"`
	CityCode                string    `json:"city_code"`
	City                    string    `json:"city"`
	NationalityCode         string    `json:"nationality_code"`
	PostalCode              string    `json:"postal_code"`
	Phone1                  string    `json:"phone1"`
	Phone2                  string    `json:"phone2"`
	Fax                     string    `json:"fax"`
	Email                   string    `json:"email"`
	Website                 string    `json:"website"`
	CompanyCode             string    `json:"company_code"`
	GuestTypeCode           string    `json:"guest_type_code"`
	IdCardCode              string    `json:"id_card_code"`
	IdCardNumber            string    `json:"id_card_number"`
	IsMale                  uint8     `json:"is_male"`
	BirthPlace              string    `json:"birth_place"`
	BirthDate               time.Time `json:"birth_date"`
	TypeCode                string    `json:"type_code"`
	CustomField01           string    `json:"custom_field01"`
	CustomField02           string    `json:"custom_field02"`
	CustomField03           string    `json:"custom_field03"`
	CustomField04           string    `json:"custom_field04"`
	CustomField05           string    `json:"custom_field05"`
	CustomField06           string    `json:"custom_field06"`
	CustomField07           string    `json:"custom_field07"`
	CustomField08           string    `json:"custom_field08"`
	CustomField09           string    `json:"custom_field09"`
	CustomField10           string    `json:"custom_field10"`
	CustomField11           string    `json:"custom_field11"`
	CustomField12           string    `json:"custom_field12"`
	CustomLookupFieldCode01 string    `json:"custom_lookup_field_code01"`
	CustomLookupFieldCode02 string    `json:"custom_lookup_field_code02"`
	CustomLookupFieldCode03 string    `json:"custom_lookup_field_code03"`
	CustomLookupFieldCode04 string    `json:"custom_lookup_field_code04"`
	CustomLookupFieldCode05 string    `json:"custom_lookup_field_code05"`
	CustomLookupFieldCode06 string    `json:"custom_lookup_field_code06"`
	CustomLookupFieldCode07 string    `json:"custom_lookup_field_code07"`
	CustomLookupFieldCode08 string    `json:"custom_lookup_field_code08"`
	CustomLookupFieldCode09 string    `json:"custom_lookup_field_code09"`
	CustomLookupFieldCode10 string    `json:"custom_lookup_field_code10"`
	CustomLookupFieldCode11 string    `json:"custom_lookup_field_code11"`
	CustomLookupFieldCode12 string    `json:"custom_lookup_field_code12"`
	IsActive                uint8     `json:"is_active"`
	IsBlacklist             uint8     `json:"is_blacklist"`
	CustomerCode            string    `json:"customer_code"`
	TpMemberCode            string    `json:"tp_member_code"`
	TpTypeCode              string    `json:"tp_type_code"`
	Source                  string    `json:"source"`
	CreatedAt               time.Time `json:"created_at"`
	CreatedBy               string    `json:"created_by"`
	UpdatedAt               time.Time `json:"updated_at"`
	UpdatedBy               string    `json:"updated_by"`
	Id                      uint64    `json:"id" gorm:"primaryKey"`
}

type Guest_profile_required struct {
	TitleCode               string    `json:"title_code"`
	FullName                string    `json:"full_name" binding:"required"`
	Street                  string    `json:"street"`
	CountryCode             string    `json:"country_code"`
	StateCode               string    `json:"state_code"`
	CityCode                string    `json:"city_code"`
	City                    string    `json:"city"`
	NationalityCode         string    `json:"nationality_code"`
	PostalCode              string    `json:"postal_code"`
	Phone1                  string    `json:"phone1"`
	Phone2                  string    `json:"phone2"`
	Fax                     string    `json:"fax"`
	Email                   string    `json:"email"`
	Website                 string    `json:"website"`
	CompanyCode             string    `json:"company_code"`
	GuestTypeCode           string    `json:"guest_type_code"`
	IdCardCode              string    `json:"id_card_code"`
	IdCardNumber            string    `json:"id_card_number"`
	IsMale                  uint8     `json:"is_male"`
	BirthPlace              string    `json:"birth_place"`
	BirthDate               time.Time `json:"birth_date"`
	TypeCode                string    `json:"type_code"`
	CustomField01           string    `json:"custom_field01"`
	CustomField02           string    `json:"custom_field02"`
	CustomField03           string    `json:"custom_field03"`
	CustomField04           string    `json:"custom_field04"`
	CustomField05           string    `json:"custom_field05"`
	CustomField06           string    `json:"custom_field06"`
	CustomField07           string    `json:"custom_field07"`
	CustomField08           string    `json:"custom_field08"`
	CustomField09           string    `json:"custom_field09"`
	CustomField10           string    `json:"custom_field10"`
	CustomField11           string    `json:"custom_field11"`
	CustomField12           string    `json:"custom_field12"`
	CustomLookupFieldCode01 string    `json:"custom_lookup_field_code01"`
	CustomLookupFieldCode02 string    `json:"custom_lookup_field_code02"`
	CustomLookupFieldCode03 string    `json:"custom_lookup_field_code03"`
	CustomLookupFieldCode04 string    `json:"custom_lookup_field_code04"`
	CustomLookupFieldCode05 string    `json:"custom_lookup_field_code05"`
	CustomLookupFieldCode06 string    `json:"custom_lookup_field_code06"`
	CustomLookupFieldCode07 string    `json:"custom_lookup_field_code07"`
	CustomLookupFieldCode08 string    `json:"custom_lookup_field_code08"`
	CustomLookupFieldCode09 string    `json:"custom_lookup_field_code09"`
	CustomLookupFieldCode10 string    `json:"custom_lookup_field_code10"`
	CustomLookupFieldCode11 string    `json:"custom_lookup_field_code11"`
	CustomLookupFieldCode12 string    `json:"custom_lookup_field_code12"`
	IsActive                uint8     `json:"is_active"`
	IsBlacklist             uint8     `json:"is_blacklist"`
	CustomerCode            string    `json:"customer_code"`
	Source                  string    `json:"source"`
	CreatedAt               time.Time `json:"created_at"`
	CreatedBy               string    `json:"created_by"`
	UpdatedAt               time.Time `json:"updated_at"`
	UpdatedBy               string    `json:"updated_by"`
	Id                      uint64    `json:"id" gorm:"primaryKey"`
}

type Guest_scheduled_rate struct {
	FolioNumber  uint64    `json:"folio_number" binding:"required"`
	FromDate     time.Time `json:"from_date" binding:"required"`
	ToDate       time.Time `json:"to_date" binding:"required"`
	RoomRateCode string    `json:"room_rate_code" binding:"required"`
	Rate         *float64  `json:"rate" binding:"required"`
	ComplimentHu string    `json:"compliment_hu"`
	CreatedAt    time.Time `json:"created_at"`
	CreatedBy    string    `json:"created_by"`
	UpdatedAt    time.Time `json:"updated_at"`
	UpdatedBy    string    `json:"updated_by"`
	Id           uint64    `json:"id"`
}

type Guest_to_do struct {
	FolioNumber  uint64    `json:"folio_number" binding:"required"`
	Description  string    `json:"description" binding:"required"`
	PostingDate  time.Time `json:"posting_date" binding:"required"`
	ReminderDate time.Time `json:"reminder_date"`
	IsDone       uint8     `json:"is_done" binding:"required"`
	CreatedAt    time.Time `json:"created_at"`
	CreatedBy    string    `json:"created_by"`
	UpdatedAt    time.Time `json:"updated_at"`
	UpdatedBy    string    `json:"updated_by"`
	Id           uint64    `json:"id"`
}

type Hotel_information struct {
	Code        string `json:"code" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Street      string `json:"street" binding:"required"`
	City        string `json:"city" binding:"required"`
	CountryCode string `json:"country_code"`
	StateCode   string `json:"state_code"`
	PostalCode  string `json:"postal_code"`
	Phone1      string `json:"phone1"`
	Phone2      string `json:"phone2"`
	Fax         string `json:"fax"`
	Email       string `json:"email"`
	Website     string `json:"website"`
	ImageUrl    string `json:"image_url"`
}

type Inv_cfg_init_item struct {
	Code          string    `json:"code" binding:"required"`
	Name          string    `json:"name" binding:"required"`
	CategoryCode  string    `json:"category_code" binding:"required"`
	UomCode       string    `json:"uom_code" binding:"required"`
	PurchasePrice float64   `json:"purchase_price"`
	SellPrice     float64   `json:"sell_price"`
	Barcode       string    `json:"barcode"`
	StockMinimum  float64   `json:"stock_minimum"`
	StockMaximum  float64   `json:"stock_maximum"`
	Remark        string    `json:"remark"`
	IsActive      uint8     `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	CreatedBy     string    `json:"created_by"`
	UpdatedAt     time.Time `json:"updated_at"`
	UpdatedBy     string    `json:"updated_by"`
	Id            uint64    `json:"id"`
}

type Inv_cfg_init_item_category struct {
	Code                         string    `json:"code" binding:"required"`
	Name                         string    `json:"name" binding:"required"`
	GroupCode                    string    `json:"group_code" binding:"required"`
	JournalAccountCode           string    `json:"journal_account_code" binding:"required"`
	JournalAccountCodeCogs       string    `json:"journal_account_code_cogs" binding:"required"`
	JournalAccountCodeCogs2      string    `json:"journal_account_code_cogs2" binding:"required"`
	JournalAccountCodeExpense    string    `json:"journal_account_code_expense" binding:"required"`
	JournalAccountCodeSell       string    `json:"journal_account_code_sell" binding:"required"`
	JournalAccountCodeAdjustment string    `json:"journal_account_code_adjustment" binding:"required"`
	JournalAccountCodeSpoil      string    `json:"journal_account_code_spoil" binding:"required"`
	CreatedAt                    time.Time `json:"created_at"`
	CreatedBy                    string    `json:"created_by"`
	UpdatedAt                    time.Time `json:"updated_at"`
	UpdatedBy                    string    `json:"updated_by"`
	Id                           uint64    `json:"id"`
}

type Inv_cfg_init_item_category_other_cogs struct {
	CategoryCode       string    `json:"category_code" binding:"required"`
	SubDepartmentCode  string    `json:"sub_department_code" binding:"required"`
	JournalAccountCode string    `json:"journal_account_code" binding:"required"`
	CreatedAt          time.Time `json:"created_at"`
	CreatedBy          string    `json:"created_by"`
	UpdatedAt          time.Time `json:"updated_at"`
	UpdatedBy          string    `json:"updated_by"`
	Id                 uint64    `json:"id"`
}

type Inv_cfg_init_item_category_other_cogs2 struct {
	CategoryCode       string    `json:"category_code" binding:"required"`
	SubDepartmentCode  string    `json:"sub_department_code" binding:"required"`
	JournalAccountCode string    `json:"journal_account_code" binding:"required"`
	CreatedAt          time.Time `json:"created_at"`
	CreatedBy          string    `json:"created_by"`
	UpdatedAt          time.Time `json:"updated_at"`
	UpdatedBy          string    `json:"updated_by"`
	Id                 uint64    `json:"id"`
}

type Inv_cfg_init_item_category_other_expense struct {
	CategoryCode       string    `json:"category_code" binding:"required"`
	SubDepartmentCode  string    `json:"sub_department_code" binding:"required"`
	JournalAccountCode string    `json:"journal_account_code" binding:"required"`
	CreatedAt          time.Time `json:"created_at"`
	CreatedBy          string    `json:"created_by"`
	UpdatedAt          time.Time `json:"updated_at"`
	UpdatedBy          string    `json:"updated_by"`
	Id                 uint64    `json:"id"`
}

type Inv_cfg_init_item_group struct {
	Code           string    `json:"code" binding:"required"`
	Name           string    `json:"name" binding:"required"`
	PercentageCost float64   `json:"percentage_cost"`
	ItemGroupType  string    `json:"item_group_type" binding:"required"`
	IdSort         int       `json:"id_sort"`
	CreatedAt      time.Time `json:"created_at"`
	CreatedBy      string    `json:"created_by"`
	UpdatedAt      time.Time `json:"updated_at"`
	UpdatedBy      string    `json:"updated_by"`
	Id             uint64    `json:"id"`
}

type Inv_cfg_init_item_uom struct {
	ItemCode      string    `json:"item_code" binding:"required"`
	UomCode       string    `json:"uom_code" binding:"required"`
	Quantity      float64   `json:"quantity" binding:"required"`
	PurchasePrice float64   `json:"purchase_price"`
	SellPrice     float64   `json:"sell_price"`
	Barcode       string    `json:"barcode"`
	CreatedAt     time.Time `json:"created_at"`
	CreatedBy     string    `json:"created_by"`
	UpdatedAt     time.Time `json:"updated_at"`
	UpdatedBy     string    `json:"updated_by"`
	Id            uint64    `json:"id"`
}

type Inv_cfg_init_market_list struct {
	CompanyCode string    `json:"company_code" binding:"required"`
	ItemCode    string    `json:"item_code" binding:"required"`
	UomCode     string    `json:"uom_code" binding:"required"`
	Price       float64   `json:"price" binding:"required"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   string    `json:"created_by"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedBy   string    `json:"updated_by"`
	Id          uint64    `json:"id"`
}

type Inv_cfg_init_return_stock_reason struct {
	Code      string    `json:"code" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Inv_cfg_init_store struct {
	Code      string    `json:"code" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	IsRoom    uint8     `json:"is_room"`
	IdSort    int       `json:"id_sort"`
	UnitCode  string    `json:"unit_code"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Inv_cfg_init_uom struct {
	Code      string    `json:"code" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Inv_close_log struct {
	Id          uint64    `json:"id"`
	ClosedDate  time.Time `json:"closed_date" binding:"required"`
	PostingDate time.Time `json:"posting_date" binding:"required"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   string    `json:"created_by"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedBy   string    `json:"updated_by"`
}

type Inv_close_summary struct {
	Date        time.Time `json:"date" binding:"required"`
	ItemCode    string    `json:"item_code" binding:"required"`
	QuantityAll float64   `json:"quantity_all" binding:"required"`
	PriceAll    float64   `json:"price_all" binding:"required"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   string    `json:"created_by"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedBy   string    `json:"updated_by"`
	Id          uint64    `json:"id"`
}

type Inv_close_summary_store struct {
	Date        time.Time `json:"date" binding:"required"`
	StoreCode   string    `json:"store_code" binding:"required"`
	ItemCode    string    `json:"item_code" binding:"required"`
	QuantityAll float64   `json:"quantity_all" binding:"required"`
	PriceAll    float64   `json:"price_all" binding:"required"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   string    `json:"created_by"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedBy   string    `json:"updated_by"`
	Id          uint64    `json:"id"`
}

type Inv_cost_recipe struct {
	ProductCode string    `json:"product_code" binding:"required"`
	StoreCode   string    `json:"store_code"`
	ItemCode    string    `json:"item_code" binding:"required"`
	Quantity    float64   `json:"quantity" binding:"required"`
	UomCode     string    `json:"uom_code" binding:"required"`
	Remark      string    `json:"remark" binding:"required"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   string    `json:"created_by"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedBy   string    `json:"updated_by"`
	Id          uint64    `json:"id"`
}

type Inv_costing struct {
	Number             string    `json:"number" binding:"required"`
	RefNumber          string    `json:"ref_number" binding:"required"`
	DocumentNumber     string    `json:"document_number"`
	SubDepartmentCode  string    `json:"sub_department_code" binding:"required"`
	StoreCode          string    `json:"store_code" binding:"required"`
	Date               time.Time `json:"date" binding:"required"`
	RequestBy          string    `json:"request_by" binding:"required"`
	Remark             *string   `json:"remark"`
	IsStoreRequisition uint8     `json:"is_store_requisition" binding:"required"`
	IsOpname           uint8     `json:"is_opname" binding:"required"`
	IsProduction       uint8     `json:"is_production" binding:"required"`
	IsReturn           uint8     `json:"is_return" binding:"required"`
	IsRoom             uint8     `json:"is_room" binding:"required"`
	IsCostRecipe       uint8     `json:"is_cost_recipe" binding:"required"`
	CreatedAt          time.Time `json:"created_at"`
	CreatedBy          string    `json:"created_by"`
	UpdatedAt          time.Time `json:"updated_at"`
	UpdatedBy          string    `json:"updated_by"`
	Id                 uint64    `json:"id"`
}

type Inv_costing_detail struct {
	CostingNumber      string    `json:"costing_number" binding:"required"`
	StoreCode          string    `json:"store_code" binding:"required"`
	StoreId            uint64    `json:"store_id" binding:"required"`
	ItemCode           string    `json:"item_code" binding:"required"`
	ItemId             uint64    `json:"item_id" binding:"required"`
	Date               time.Time `json:"date" binding:"required"`
	Quantity           float64   `json:"quantity" binding:"required"`
	UomCode            string    `json:"uom_code" binding:"required"`
	Price              float64   `json:"price" binding:"required"`
	TotalPrice         float64   `json:"total_price" binding:"required"`
	ReceiveId          uint64    `json:"receive_id" binding:"required"`
	JournalAccountCode string    `json:"journal_account_code" binding:"required"`
	ItemGroupCode      string    `json:"item_group_code" binding:"required"`
	ReasonCode         string    `json:"reason_code"`
	IsSpoil            uint8     `json:"is_spoil" binding:"required"`
	IsCogs             uint8     `json:"is_cogs" binding:"required"`
	CreatedAt          time.Time `json:"created_at"`
	CreatedBy          string    `json:"created_by"`
	UpdatedAt          time.Time `json:"updated_at"`
	UpdatedBy          string    `json:"updated_by"`
	Id                 uint64    `json:"id" gorm:"primaryKey"`
}

type Inv_opname struct {
	Number        string    `json:"number" binding:"required"`
	RefNumber     string    `json:"ref_number" binding:"required"`
	ReceiveNumber string    `json:"receive_number" binding:"required"`
	CostingNumber string    `json:"costing_number" binding:"required"`
	Date          time.Time `json:"date" binding:"required"`
	CreatedAt     time.Time `json:"created_at"`
	CreatedBy     string    `json:"created_by"`
	UpdatedAt     time.Time `json:"updated_at"`
	UpdatedBy     string    `json:"updated_by"`
	Id            uint64    `json:"id"`
}

type Inv_production struct {
	Number         string    `json:"number" binding:"required"`
	RefNumber      string    `json:"ref_number" binding:"required"`
	ReceiveNumber  string    `json:"receive_number" binding:"required"`
	CostingNumber  string    `json:"costing_number" binding:"required"`
	DocumentNumber string    `json:"document_number"`
	Date           time.Time `json:"date" binding:"required"`
	Remark         *string   `json:"remark"`
	CreatedAt      time.Time `json:"created_at"`
	CreatedBy      string    `json:"created_by"`
	UpdatedAt      time.Time `json:"updated_at"`
	UpdatedBy      string    `json:"updated_by"`
	Id             uint64    `json:"id"`
}

type Inv_purchase_order struct {
	Number              string    `json:"number" binding:"required"`
	CompanyCode         string    `json:"company_code" binding:"required"`
	ExpeditionCode      string    `json:"expedition_code"`
	PrNumber            string    `json:"pr_number" binding:"required"`
	Date                time.Time `json:"date" binding:"required"`
	ShippingAddressCode string    `json:"shipping_address_code" binding:"required"`
	ContactPerson       string    `json:"contact_person" binding:"required"`
	Street              string    `json:"street" binding:"required"`
	City                string    `json:"city"`
	CountryCode         string    `json:"country_code"`
	StateCode           string    `json:"state_code"`
	PostalCode          string    `json:"postal_code"`
	Phone1              string    `json:"phone1"`
	Phone2              string    `json:"phone2"`
	Fax                 string    `json:"fax"`
	Email               string    `json:"email"`
	RequestBy           string    `json:"request_by" binding:"required"`
	Remark              string    `json:"remark"`
	IsReceived          uint8     `json:"is_received" binding:"required"`
	CreatedAt           time.Time `json:"created_at"`
	CreatedBy           string    `json:"created_by"`
	UpdatedAt           time.Time `json:"updated_at"`
	UpdatedBy           string    `json:"updated_by"`
	Id                  uint64    `json:"id" gorm:"primaryKey"`
}

type Inv_purchase_order_detail struct {
	PoNumber            string    `json:"po_number" binding:"required"`
	ItemCode            string    `json:"item_code" binding:"required"`
	StoreCode           string    `json:"store_code" binding:"required"`
	Quantity            float64   `json:"quantity" binding:"required"`
	QuantityReceived    float64   `json:"quantity_received" binding:"required"`
	QuantityNotReceived float64   `json:"quantity_not_received" binding:"required"`
	Convertion          float64   `json:"convertion" binding:"required"`
	UomCode             string    `json:"uom_code" binding:"required"`
	Price               float64   `json:"price" binding:"required"`
	Remark              string    `json:"remark"`
	CreatedAt           time.Time `json:"created_at"`
	CreatedBy           string    `json:"created_by"`
	UpdatedAt           time.Time `json:"updated_at"`
	UpdatedBy           string    `json:"updated_by"`
	Id                  uint64    `json:"id"`
}

type Inv_purchase_request struct {
	Number              string         `json:"number" binding:"required"`
	SubDepartmentCode   string         `json:"sub_department_code" binding:"required"`
	Date                time.Time      `json:"date" binding:"required"`
	NeedDate            time.Time      `json:"need_date" binding:"required"`
	ShippingAddressCode string         `json:"shipping_address_code" binding:"required"`
	ContactPerson       *string        `json:"contact_person" binding:"required"`
	Street              *string        `json:"street" binding:"required"`
	City                *string        `json:"city"`
	CountryCode         *string        `json:"country_code"`
	StateCode           *string        `json:"state_code"`
	PostalCode          *string        `json:"postal_code"`
	Phone1              *string        `json:"phone1"`
	Phone2              *string        `json:"phone2"`
	Fax                 *string        `json:"fax"`
	Email               *string        `json:"email"`
	RequestBy           *string        `json:"request_by" binding:"required"`
	Remark              *string        `json:"remark"`
	IsUserApproved1     uint8          `json:"is_user_approved1" binding:"required"`
	ApproveBy1          string         `json:"approve_by1"`
	ApproveDate1        time.Time      `json:"approve_date1"`
	IsUserApproved12    uint8          `json:"is_user_approved12" binding:"required"`
	ApproveBy12         string         `json:"approve_by12"`
	ApproveDate12       time.Time      `json:"approve_date12"`
	IsUserApproved2     uint8          `json:"is_user_approved2" binding:"required"`
	ApproveBy2          string         `json:"approve_by2"`
	ApproveDate2        time.Time      `json:"approve_date2"`
	IsUserApproved3     uint8          `json:"is_user_approved3" binding:"required"`
	ApproveBy3          string         `json:"approve_by3"`
	ApproveDate3        time.Time      `json:"approve_date3"`
	RejectBy            string         `json:"reject_by"`
	RejectDate          time.Time      `json:"reject_date"`
	StatusCode          string         `json:"status_code" binding:"required"`
	CreatedAt           time.Time      `json:"created_at"`
	CreatedBy           string         `json:"created_by"`
	UpdatedAt           time.Time      `json:"updated_at"`
	UpdatedBy           string         `json:"updated_by"`
	DeletedAt           gorm.DeletedAt `json:"deleted_at"`
	DeletedBy           string         `json:"deleted_by" gorm:"->:false;<-:create"`
	Id                  uint64         `json:"id" gorm:"primaryKey"`
}

type Inv_purchase_request_detail struct {
	PrNumber         string    `json:"pr_number" binding:"required"`
	ItemCode         string    `json:"item_code" binding:"required"`
	Quantity         float64   `json:"quantity" binding:"required"`
	QuantityApproved *float64  `json:"quantity_approved"`
	Convertion       float64   `json:"convertion" binding:"required"`
	UomCode          string    `json:"uom_code" binding:"required"`
	CompanyCode      string    `json:"company_code"`
	Price            float64   `json:"price"`
	CompanyCode2     *string   `json:"company_code2"`
	Price2           *float64  `json:"price2"`
	CompanyCode3     *string   `json:"company_code3"`
	Price3           *float64  `json:"price3"`
	EstimatePrice    float64   `json:"estimate_price" binding:"required"`
	StoreCode        string    `json:"store_code"`
	Remark           *string   `json:"remark"`
	CreatedAt        time.Time `json:"created_at"`
	CreatedBy        string    `json:"created_by"`
	UpdatedAt        time.Time `json:"updated_at"`
	UpdatedBy        string    `json:"updated_by"`
	Id               uint64    `json:"id"`
}

type Inv_receiving struct {
	Number            string    `json:"number" binding:"required"`
	RefNumber         string    `json:"ref_number" binding:"required"`
	PoNumber          *string   `json:"po_number" binding:"required"`
	ApNumber          *string   `json:"ap_number"`
	CostingNumber     *string   `json:"costing_number"`
	CompanyCode       string    `json:"company_code" binding:"required"`
	InvoiceNumber     string    `json:"invoice_number"`
	BankAccountCode   *string   `json:"bank_account_code"`
	AmountPayment     *float64  `json:"amount_payment"`
	Date              time.Time `json:"date" binding:"required"`
	IsConsignment     uint8     `json:"is_consignment" binding:"required"`
	Remark            *string   `json:"remark"`
	IsSeparate        *uint8    `json:"is_separate"`
	IsDiscountIncome  *uint8    `json:"is_discount_income"`
	IsTaxExpense      *uint8    `json:"is_tax_expense"`
	IsShippingExpense *uint8    `json:"is_shipping_expense"`
	IsCredit          *uint8    `json:"is_credit"`
	DueDate           time.Time `json:"due_date"`
	IsPaid            uint8     `json:"is_paid" binding:"required"`
	IsOpname          uint8     `json:"is_opname" binding:"required"`
	IsProduction      uint8     `json:"is_production" binding:"required"`
	CreatedAt         time.Time `json:"created_at"`
	CreatedBy         string    `json:"created_by"`
	UpdatedAt         time.Time `json:"updated_at"`
	UpdatedBy         string    `json:"updated_by"`
	Id                uint64    `json:"id"`
}

type Inv_receiving_detail struct {
	ReceiveNumber      string    `json:"receive_number" binding:"required"`
	StoreCode          string    `json:"store_code" binding:"required"`
	StoreId            uint64    `json:"store_id" binding:"required"`
	ItemCode           string    `json:"item_code" binding:"required"`
	ItemId             uint64    `json:"item_id" binding:"required"`
	Date               time.Time `json:"date" binding:"required"`
	PoId               uint64    `json:"po_id" binding:"required"`
	PoQuantity         float64   `json:"po_quantity" binding:"required"`
	ReceiveQuantity    float64   `json:"receive_quantity" binding:"required"`
	ReceiveUomCode     string    `json:"receive_uom_code" binding:"required"`
	ReceivePrice       float64   `json:"receive_price" binding:"required"`
	BasicQuantity      float64   `json:"basic_quantity" binding:"required"`
	BasicUomCode       string    `json:"basic_uom_code" binding:"required"`
	BasicPrice         float64   `json:"basic_price" binding:"required"`
	Quantity           float64   `json:"quantity" binding:"required"`
	TotalPrice         float64   `json:"total_price" binding:"required"`
	Discount           float64   `json:"discount"`
	Tax                float64   `json:"tax"`
	Shipping           float64   `json:"shipping"`
	Remark             *string   `json:"remark"`
	ExpireDate         time.Time `json:"expire_date"`
	IsCogs             uint8     `json:"is_cogs" binding:"required"`
	JournalAccountCode string    `json:"journal_account_code"`
	ItemGroupCode      string    `json:"item_group_code"`
	CreatedAt          time.Time `json:"created_at"`
	CreatedBy          string    `json:"created_by"`
	UpdatedAt          time.Time `json:"updated_at"`
	UpdatedBy          string    `json:"updated_by"`
	Id                 uint64    `json:"id" gorm:"primaryKey"`
}

type Inv_return_stock struct {
	Number          string    `json:"number" binding:"required"`
	RefNumber       string    `json:"ref_number" binding:"required"`
	CostingNumber   string    `json:"costing_number" binding:"required"`
	ArNumber        string    `json:"ar_number"`
	CompanyCode     string    `json:"company_code" binding:"required"`
	DocumentNumber  string    `json:"document_number" binding:"required"`
	BankAccountCode string    `json:"bank_account_code" binding:"required"`
	TotalReturn     float64   `json:"total_return" binding:"required"`
	DueDate         time.Time `json:"due_date"`
	PaymentRemark   *string   `json:"payment_remark"`
	CreatedAt       time.Time `json:"created_at"`
	CreatedBy       string    `json:"created_by"`
	UpdatedAt       time.Time `json:"updated_at"`
	UpdatedBy       string    `json:"updated_by"`
	Id              uint64    `json:"id"`
}

type Inv_stock_transfer struct {
	Number             string    `json:"number" binding:"required"`
	DocumentNumber     string    `json:"document_number" binding:"required"`
	RequestBy          string    `json:"request_by" binding:"required"`
	StoreCode          string    `json:"store_code" binding:"required"`
	Date               time.Time `json:"date" binding:"required"`
	Remark             string    `json:"remark"`
	IsStoreRequisition uint8     `json:"is_store_requisition" binding:"required"`
	CreatedAt          time.Time `json:"created_at"`
	CreatedBy          string    `json:"created_by"`
	UpdatedAt          time.Time `json:"updated_at"`
	UpdatedBy          string    `json:"updated_by"`
	Id                 uint64    `json:"id"`
}

type Inv_stock_transfer_detail struct {
	StNumber      string    `json:"st_number" binding:"required"`
	FromStoreCode string    `json:"from_store_code" binding:"required"`
	ToStoreCode   string    `json:"to_store_code" binding:"required"`
	ItemCode      string    `json:"item_code" binding:"required"`
	Quantity      float64   `json:"quantity" binding:"required"`
	UomCode       string    `json:"uom_code" binding:"required"`
	Price         float64   `json:"price" binding:"required"`
	TotalPrice    float64   `json:"total_price" binding:"required"`
	ReceiveId     uint64    `json:"receive_id" binding:"required"`
	CreatedAt     time.Time `json:"created_at"`
	CreatedBy     string    `json:"created_by"`
	UpdatedAt     time.Time `json:"updated_at"`
	UpdatedBy     string    `json:"updated_by"`
	Id            uint64    `json:"id"`
}

type Inv_store_requisition struct {
	Number            string    `json:"number" binding:"required"`
	SubDepartmentCode string    `json:"sub_department_code" binding:"required"`
	StoreCode         string    `json:"store_code" binding:"required"`
	Date              time.Time `json:"date" binding:"required"`
	DocumentNumber    string    `json:"document_number" binding:"required"`
	RequestBy         string    `json:"request_by" binding:"required"`
	Remark            *string   `json:"remark"`
	IsUserApproved1   uint8     `json:"is_user_approved1" binding:"required"`
	ApproveBy1        string    `json:"approve_by1"`
	ApproveDate1      time.Time `json:"approve_date1"`
	IsUserApproved2   uint8     `json:"is_user_approved2" binding:"required"`
	ApproveBy2        string    `json:"approve_by2"`
	ApproveDate2      time.Time `json:"approve_date2"`
	RejectBy          string    `json:"reject_by"`
	RejectDate        time.Time `json:"reject_date"`
	StNumber          string    `json:"st_number"`
	CostingNumber     string    `json:"costing_number"`
	StatusCode        string    `json:"status_code" binding:"required"`
	CreatedAt         time.Time `json:"created_at"`
	CreatedBy         string    `json:"created_by"`
	UpdatedAt         time.Time `json:"updated_at"`
	UpdatedBy         string    `json:"updated_by"`
	Id                uint64    `json:"id"`
}

type Inv_store_requisition_detail struct {
	SrNumber         string    `json:"sr_number" binding:"required"`
	FromStoreCode    string    `json:"from_store_code" binding:"required"`
	ToStoreCode      string    `json:"to_store_code" binding:"required"`
	ItemCode         string    `json:"item_code" binding:"required"`
	Quantity         float64   `json:"quantity" binding:"required"`
	QuantityApproved float64   `json:"quantity_approved" binding:"required"`
	Convertion       float64   `json:"convertion" binding:"required"`
	UomCode          string    `json:"uom_code" binding:"required"`
	EstimatePrice    float64   `json:"estimate_price" binding:"required"`
	CreatedAt        time.Time `json:"created_at"`
	CreatedBy        string    `json:"created_by"`
	UpdatedAt        time.Time `json:"updated_at"`
	UpdatedBy        string    `json:"updated_by"`
	Id               uint64    `json:"id"`
}

type Invoice struct {
	Number          string    `json:"number" binding:"required" gorm:"primaryKey"`
	CompanyCode     string    `json:"company_code" binding:"required"`
	ContactPersonId uint64    `json:"contact_person_id" binding:"required"`
	IssuedDate      time.Time `json:"issued_date" binding:"required"`
	DueDate         time.Time `json:"due_date" binding:"required"`
	Remark          string    `json:"remark"`
	IsPaid          uint8     `json:"is_paid" binding:"required"`
	RefNumber       string    `json:"ref_number"`
	PrintCount      int       `json:"print_count" binding:"required"`
	CreatedAt       time.Time `json:"created_at"`
	CreatedBy       string    `json:"created_by"`
	UpdatedAt       time.Time `json:"updated_at"`
	UpdatedBy       string    `json:"updated_by"`
	Id              uint64    `json:"id"`
}

type Invoice_item struct {
	InvoiceNumber        string    `json:"invoice_number" binding:"required"`
	SubFolioId           uint64    `json:"sub_folio_id" binding:"required"`
	FolioNumber          uint64    `json:"folio_number" binding:"required"`
	CorrectionBreakdown  uint64    `json:"correction_breakdown" binding:"required"`
	Amount               float64   `json:"amount" binding:"required"`
	AmountCharged        float64   `json:"amount_charged"`
	DefaultCurrencyCode  string    `json:"default_currency_code" binding:"required"`
	AmountChargedForeign float64   `json:"amount_charged_foreign" binding:"required"`
	ExchangeRate         float64   `json:"exchange_rate" binding:"required"`
	CurrencyCode         string    `json:"currency_code" binding:"required"`
	AmountPaid           float64   `json:"amount_paid"`
	RefNumber            string    `json:"ref_number" binding:"required"`
	Remark               string    `json:"remark" binding:"required"`
	TypeCode             string    `json:"type_code" binding:"required"`
	CreatedAt            time.Time `json:"created_at"`
	CreatedBy            string    `json:"created_by"`
	UpdatedAt            time.Time `json:"updated_at"`
	UpdatedBy            string    `json:"updated_by"`
	Id                   uint64    `json:"id"`
}

type Invoice_payment struct {
	InvoiceNumber       string    `json:"invoice_number" binding:"required"`
	SubFolioId          uint64    `json:"sub_folio_id" binding:"required"`
	FolioNumber         uint64    `json:"folio_number" binding:"required"`
	RefNumber           string    `json:"ref_number" binding:"required"`
	Amount              float64   `json:"amount" binding:"required"`
	DefaultCurrencyCode string    `json:"default_currency_code" binding:"required"`
	AmountForeign       float64   `json:"amount_foreign" binding:"required"`
	ExchangeRate        float64   `json:"exchange_rate" binding:"required"`
	CurrencyCode        string    `json:"currency_code" binding:"required"`
	AmountActual        float64   `json:"amount_actual" binding:"required"`
	ExchangeRateActual  float64   `json:"exchange_rate_actual" binding:"required"`
	Date                time.Time `json:"date" binding:"required"`
	Remark              string    `json:"remark"`
	CreatedAt           time.Time `json:"created_at"`
	CreatedBy           string    `json:"created_by"`
	UpdatedAt           time.Time `json:"updated_at"`
	UpdatedBy           string    `json:"updated_by"`
	Id                  uint64    `json:"id"`
}

type Log struct {
	Id           uint64    `json:"id"`
	IdTable      uint64    `json:"id_table" binding:"required"`
	Mode         string    `json:"mode" binding:"required"`
	TableNameLog string    `json:"table_name_log" binding:"required"`
	CreatedAt    time.Time `json:"created_at"`
	CreatedBy    string    `json:"created_by"`
	UpdatedAt    time.Time `json:"updated_at"`
	UpdatedBy    string    `json:"updated_by"`
	LogDate      time.Time `json:"log_date" binding:"required"`
	DataQuery    string    `json:"data_query"`
}

type Log_backup struct {
	Id        uint64    `json:"id"`
	AuditDate time.Time `json:"audit_date" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
}

type Log_keylock struct {
	Id                uint64    `json:"id"`
	ReservationNumber uint64    `json:"reservation_number" binding:"required"`
	FolioNumber       uint64    `json:"folio_number" binding:"required"`
	RoomNumber1       string    `json:"room_number1" binding:"required"`
	RoomNumber2       string    `json:"room_number2"`
	RoomNumber3       string    `json:"room_number3"`
	RoomNumber4       string    `json:"room_number4"`
	ArrivalDate       time.Time `json:"arrival_date" binding:"required"`
	DepartureDate     time.Time `json:"departure_date" binding:"required"`
	GuestName         string    `json:"guest_name" binding:"required"`
	IssuedDate        time.Time `json:"issued_date" binding:"required"`
	IssuedBy          string    `json:"issued_by" binding:"required"`
	ErasedDate        time.Time `json:"erased_date" binding:"required"`
	ErasedBy          string    `json:"erased_by" binding:"required"`
	KeylockVendorCode string    `json:"keylock_vendor_code" binding:"required"`
	CardNumber        string    `json:"card_number"`
	IsActive          uint8     `json:"is_active" binding:"required"`
}

type Log_mode struct {
	Code string `json:"code" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type Log_shift struct {
	Id             uint64    `json:"id" gorm:"primaryKey"`
	CreatedAt      time.Time `json:"created_at"`
	CreatedBy      string    `json:"created_by"`
	UpdatedAt      time.Time `json:"updated_at"`
	UpdatedBy      string    `json:"updated_by"`
	Shift          string    `json:"shift" binding:"required"`
	StartDate      time.Time `json:"start_date" binding:"required"`
	EndDate        time.Time `json:"end_date" binding:"required"`
	AuditDate      time.Time `json:"audit_date" binding:"required"`
	OpeningBalance float64   `json:"opening_balance" binding:"required"`
	Remark         string    `json:"remark" binding:"required"`
	IpAddress      string    `json:"ip_address" binding:"required"`
	ComputerName   string    `json:"computer_name" binding:"required"`
	MacAddress     string    `json:"mac_address" binding:"required"`
	IsOpen         uint8     `json:"is_open" binding:"required" `
}

type Log_special_access struct {
	Id           uint64    `json:"id"`
	SystemCode   string    `json:"system_code" binding:"required"`
	AccessName   string    `json:"access_name" binding:"required"`
	LogDate      time.Time `json:"log_date" binding:"required"`
	AuditDate    time.Time `json:"audit_date" binding:"required"`
	IpAddress    string    `json:"ip_address" binding:"required"`
	ComputerName string    `json:"computer_name" binding:"required"`
	MacAddress   string    `json:"mac_address" binding:"required"`
	AccessDenied uint8     `json:"access_denied" binding:"required"`
	CreatedAt    time.Time `json:"created_at"`
	CreatedBy    string    `json:"created_by"`
	UpdatedAt    time.Time `json:"updated_at"`
	UpdatedBy    string    `json:"updated_by"`
}

type Log_user struct {
	Id           uint64    `json:"id"`
	SystemCode   string    `json:"system_code" binding:"required"`
	ActionId     int       `json:"action_id" binding:"required"`
	ActualDate   time.Time `json:"actual_date"`
	AuditDate    time.Time `json:"audit_date" binding:"required"`
	IpAddress    string    `json:"ip_address" binding:"required"`
	ComputerName string    `json:"computer_name" binding:"required"`
	MacAddress   string    `json:"mac_address" binding:"required"`
	DataLink1    string    `json:"data_link1" binding:"required"`
	DataLink2    string    `json:"data_link2"`
	DataLink3    string    `json:"data_link3"`
	Remark       string    `json:"remark"`
	CreatedBy    string    `json:"created_by" binding:"required"`
}

type Log_user_action struct {
	Id      int    `json:"id"`
	Name    string `json:"name" binding:"required"`
	GroupId int    `json:"group_id" binding:"required"`
	IdSort  int    `json:"id_sort" binding:"required"`
}

type Log_user_action_group struct {
	Id   int    `json:"id"`
	Name string `json:"name" binding:"required"`
}

type Lost_and_found struct {
	IsLost          *uint8    `json:"is_lost"`
	Item            string    `json:"item" binding:"required"`
	Color           *string   `json:"color"`
	Location        *string   `json:"location"`
	Who             *string   `json:"who"`
	Value           *float64  `json:"value"`
	CurrentLocation *string   `json:"current_location"`
	DatePosting     time.Time `json:"date_posting" binding:"required"`
	IsReturn        *uint8    `json:"is_return"`
	DateReturn      time.Time `json:"date_return"`
	ReturnBy        *string   `json:"return_by"`
	Owner           *string   `json:"owner"`
	Phone           *string   `json:"phone"`
	Notes           *string   `json:"notes"`
	IsActive        uint8     `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	CreatedBy       string    `json:"created_by"`
	UpdatedAt       time.Time `json:"updated_at"`
	UpdatedBy       string    `json:"updated_by"`
	Id              uint64    `json:"id"`
}

type Market_statistic struct {
	AuditDate                                 time.Time `json:"audit_date" binding:"required"`
	MarketCategoryCode                        string    `json:"market_category_code"`
	MarketCode                                string    `json:"market_code"`
	MarketCompanyCode                         string    `json:"market_company_code"`
	RoomMarket                                int       `json:"room_market"`
	RoomMarketCompliment                      int       `json:"room_market_compliment"`
	PaxMarket                                 int       `json:"pax_market"`
	PaxMarketCompliment                       int       `json:"pax_market_compliment"`
	RevenueNettMarket                         float64   `json:"revenue_nett_market"`
	RevenueNettMarketCompliment               float64   `json:"revenue_nett_market_compliment"`
	RevenueGrossMarket                        float64   `json:"revenue_gross_market"`
	RevenueGrossMarketCompliment              float64   `json:"revenue_gross_market_compliment"`
	RevenueNonPackageMarket                   float64   `json:"revenue_non_package_market"`
	RevenueNonPackageMarketCompliment         float64   `json:"revenue_non_package_market_compliment"`
	BusinessSourceCode                        string    `json:"business_source_code"`
	BusinessSourceCompanyCode                 string    `json:"business_source_company_code"`
	RoomBusinessSource                        int       `json:"room_business_source"`
	RoomBusinessSourceCompliment              int       `json:"room_business_source_compliment"`
	PaxBusinessSource                         int       `json:"pax_business_source"`
	PaxBusinessSourceCompliment               int       `json:"pax_business_source_compliment"`
	RevenueNettBusinessSource                 float64   `json:"revenue_nett_business_source"`
	RevenueNettBusinessSourceCompliment       float64   `json:"revenue_nett_business_source_compliment"`
	RevenueGrossBusinessSource                float64   `json:"revenue_gross_business_source"`
	RevenueGrossBusinessSourceCompliment      float64   `json:"revenue_gross_business_source_compliment"`
	RevenueNonPackageBusinessSource           float64   `json:"revenue_non_package_business_source"`
	RevenueNonPackageBusinessSourceCompliment float64   `json:"revenue_non_package_business_source_compliment"`
	RoomTypeCode                              string    `json:"room_type_code"`
	RoomTypeCompanyCode                       string    `json:"room_type_company_code"`
	RoomRoomType                              int       `json:"room_room_type"`
	RoomRoomTypeCompliment                    int       `json:"room_room_type_compliment"`
	PaxRoomType                               int       `json:"pax_room_type"`
	PaxRoomTypeCompliment                     int       `json:"pax_room_type_compliment"`
	RevenueNettRoomType                       float64   `json:"revenue_nett_room_type"`
	RevenueNettRoomTypeCompliment             float64   `json:"revenue_nett_room_type_compliment"`
	RevenueGrossRoomType                      float64   `json:"revenue_gross_room_type"`
	RevenueGrossRoomTypeCompliment            float64   `json:"revenue_gross_room_type_compliment"`
	RevenueNonPackageRoomType                 float64   `json:"revenue_non_package_room_type"`
	RevenueNonPackageRoomTypeCompliment       float64   `json:"revenue_non_package_room_type_compliment"`
	RoomRateCode                              string    `json:"room_rate_code"`
	RoomRateCompanyCode                       string    `json:"room_rate_company_code"`
	RoomRoomRate                              int       `json:"room_room_rate"`
	RoomRoomRateCompliment                    int       `json:"room_room_rate_compliment"`
	PaxRoomRate                               int       `json:"pax_room_rate"`
	PaxRoomRateCompliment                     int       `json:"pax_room_rate_compliment"`
	RevenueNettRoomRate                       float64   `json:"revenue_nett_room_rate"`
	RevenueNettRoomRateCompliment             float64   `json:"revenue_nett_room_rate_compliment"`
	RevenueGrossRoomRate                      float64   `json:"revenue_gross_room_rate"`
	RevenueGrossRoomRateCompliment            float64   `json:"revenue_gross_room_rate_compliment"`
	RevenueNonPackageRoomRate                 float64   `json:"revenue_non_package_room_rate"`
	RevenueNonPackageRoomRateCompliment       float64   `json:"revenue_non_package_room_rate_compliment"`
	MarketingCode                             string    `json:"marketing_code"`
	MarketingBusinessSourceCode               string    `json:"marketing_business_source_code" binding:"required"`
	MarketingCompanyCode                      string    `json:"marketing_company_code"`
	RoomMarketing                             int       `json:"room_marketing"`
	RoomMarketingCompliment                   int       `json:"room_marketing_compliment"`
	PaxMarketing                              int       `json:"pax_marketing"`
	PaxMarketingCompliment                    int       `json:"pax_marketing_compliment"`
	RevenueNettMarketing                      float64   `json:"revenue_nett_marketing"`
	RevenueNettMarketingCompliment            float64   `json:"revenue_nett_marketing_compliment"`
	RevenueGrossMarketing                     float64   `json:"revenue_gross_marketing"`
	RevenueGrossMarketingCompliment           float64   `json:"revenue_gross_marketing_compliment"`
	RevenueNonPackageMarketing                float64   `json:"revenue_non_package_marketing"`
	RevenueNonPackageMarketingCompliment      float64   `json:"revenue_non_package_marketing_compliment"`
	RevenueAllNettMarketing                   float64   `json:"revenue_all_nett_marketing" binding:"required"`
	RevenueAllGrossMarketing                  float64   `json:"revenue_all_gross_marketing" binding:"required"`
	CountryCode                               string    `json:"country_code"`
	CountryStateCode                          string    `json:"country_state_code" binding:"required"`
	CountryCityCode                           string    `json:"country_city_code" binding:"required"`
	RoomCountry                               int       `json:"room_country"`
	RoomCountryCompliment                     int       `json:"room_country_compliment"`
	PaxCountry                                int       `json:"pax_country"`
	PaxCountryCompliment                      int       `json:"pax_country_compliment"`
	RevenueNettCountry                        float64   `json:"revenue_nett_country"`
	RevenueNettCountryCompliment              float64   `json:"revenue_nett_country_compliment"`
	RevenueGrossCountry                       float64   `json:"revenue_gross_country"`
	RevenueGrossCountryCompliment             float64   `json:"revenue_gross_country_compliment"`
	RevenueNonPackageCountry                  float64   `json:"revenue_non_package_country"`
	RevenueNonPackageCountryCompliment        float64   `json:"revenue_non_package_country_compliment"`
	RevenueAllNettCountry                     float64   `json:"revenue_all_nett_country" binding:"required"`
	RevenueAllGrossCountry                    float64   `json:"revenue_all_gross_country" binding:"required"`
	NationalityCode                           string    `json:"nationality_code"`
	NationalityCountryCode                    string    `json:"nationality_country_code" binding:"required"`
	RoomNationality                           int       `json:"room_nationality"`
	RoomNationalityCompliment                 int       `json:"room_nationality_compliment"`
	PaxNationality                            int       `json:"pax_nationality"`
	PaxNationalityCompliment                  int       `json:"pax_nationality_compliment"`
	RevenueNettNationality                    float64   `json:"revenue_nett_nationality"`
	RevenueNettNationalityCompliment          float64   `json:"revenue_nett_nationality_compliment"`
	RevenueGrossNationality                   float64   `json:"revenue_gross_nationality"`
	RevenueGrossNationalityCompliment         float64   `json:"revenue_gross_nationality_compliment"`
	RevenueNonPackageNationality              float64   `json:"revenue_non_package_nationality"`
	RevenueNonPackageNationalityCompliment    float64   `json:"revenue_non_package_nationality_compliment"`
	RevenueAllNettNationality                 float64   `json:"revenue_all_nett_nationality" binding:"required"`
	RevenueAllGrossNationality                float64   `json:"revenue_all_gross_nationality" binding:"required"`
	BookingSourceCode                         string    `json:"booking_source_code"`
	RoomBookingSource                         int       `json:"room_booking_source"`
	RoomBookingSourceCompliment               int       `json:"room_booking_source_compliment"`
	PaxBookingSource                          int       `json:"pax_booking_source"`
	PaxBookingSourceCompliment                int       `json:"pax_booking_source_compliment"`
	RevenueNettBookingSource                  float64   `json:"revenue_nett_booking_source"`
	RevenueNettBookingSourceCompliment        float64   `json:"revenue_nett_booking_source_compliment"`
	RevenueGrossBookingSource                 float64   `json:"revenue_gross_booking_source"`
	RevenueGrossBookingSourceCompliment       float64   `json:"revenue_gross_booking_source_compliment"`
	RevenueNonPackageBookingSource            float64   `json:"revenue_non_package_booking_source"`
	RevenueNonPackageBookingSourceCompliment  float64   `json:"revenue_non_package_booking_source_compliment"`
	RevenueAllNettBookingSource               float64   `json:"revenue_all_nett_booking_source" binding:"required"`
	RevenueAllGrossBookingSource              float64   `json:"revenue_all_gross_booking_source" binding:"required"`
	PurposeOfCode                             string    `json:"purpose_of_code"`
	RoomPurposeOf                             int       `json:"room_purpose_of"`
	RoomPurposeOfCompliment                   int       `json:"room_purpose_of_compliment"`
	PaxPurposeOf                              int       `json:"pax_purpose_of"`
	PaxPurposeOfCompliment                    int       `json:"pax_purpose_of_compliment"`
	RevenueNettPurposeOf                      float64   `json:"revenue_nett_purpose_of"`
	RevenueNettPurposeOfCompliment            float64   `json:"revenue_nett_purpose_of_compliment"`
	RevenueGrossPurposeOf                     float64   `json:"revenue_gross_purpose_of"`
	RevenueGrossPurposeOfCompliment           float64   `json:"revenue_gross_purpose_of_compliment"`
	RevenueNonPackagePurposeOf                float64   `json:"revenue_non_package_purpose_of"`
	RevenueNonPackagePurposeOfCompliment      float64   `json:"revenue_non_package_purpose_of_compliment"`
	RevenueAllNettPurposeOf                   float64   `json:"revenue_all_nett_purpose_of" binding:"required"`
	RevenueAllGrossPurposeOf                  float64   `json:"revenue_all_gross_purpose_of" binding:"required"`
	Id                                        uint64    `json:"id"`
}

type Member struct {
	Code                 string    `json:"code" binding:"required"`
	GuestProfileId       uint64    `json:"guest_profile_id" binding:"required"`
	IsForRoom            *uint8    `json:"is_for_room" binding:"required"`
	RoomPointTypeCode    *string   `json:"room_point_type_code"`
	IsForOutlet          *uint8    `json:"is_for_outlet" binding:"required"`
	OutletPointTypeCode  *string   `json:"outlet_point_type_code"`
	IsForBanquet         *uint8    `json:"is_for_banquet" binding:"required"`
	BanquetPointTypeCode *string   `json:"banquet_point_type_code"`
	OutletDiscountCode   string    `json:"outlet_discount_code" binding:"required"`
	ExpireDate           time.Time `json:"expire_date"`
	FingerprintTemplate  []byte    `json:"fingerprint_template"`
	CreatedAt            time.Time `json:"created_at"`
	CreatedBy            string    `json:"created_by"`
	UpdatedAt            time.Time `json:"updated_at"`
	UpdatedBy            string    `json:"updated_by"`
	Id                   uint64    `json:"id"`
}

type Member_gift struct {
	Number         string    `json:"number" binding:"required"`
	IssuedDate     time.Time `json:"issued_date" binding:"required"`
	MemberTypeCode string    `json:"member_type_code" binding:"required"`
	StatusCode     string    `json:"status_code" binding:"required"`
	ExpireDate     time.Time `json:"expire_date" binding:"required"`
	Description    string    `json:"description" binding:"required"`
	CreatedAt      time.Time `json:"created_at"`
	CreatedBy      string    `json:"created_by"`
	UpdatedAt      time.Time `json:"updated_at"`
	UpdatedBy      string    `json:"updated_by"`
	Id             uint64    `json:"id"`
}

type Member_point struct {
	MemberCode     string    `json:"member_code" binding:"required"`
	AuditDate      time.Time `json:"audit_date" binding:"required"`
	PointTypeCode  string    `json:"point_type_code" binding:"required"`
	MemberTypeCode string    `json:"member_type_code" binding:"required"`
	FolioNumber    uint64    `json:"folio_number" binding:"required"`
	IsFromRate     uint8     `json:"is_from_rate" binding:"required"`
	RoomTypeCode   string    `json:"room_type_code" binding:"required"`
	RateAmount     float64   `json:"rate_amount" binding:"required"`
	Point          float64   `json:"point" binding:"required"`
	ExpireDate     time.Time `json:"expire_date" binding:"required"`
	IsRedeemed     uint8     `json:"is_redeemed" binding:"required"`
	IsExpired      uint8     `json:"is_expired" binding:"required"`
	CreatedAt      time.Time `json:"created_at"`
	CreatedBy      string    `json:"created_by"`
	UpdatedAt      time.Time `json:"updated_at"`
	UpdatedBy      string    `json:"updated_by"`
	Id             uint64    `json:"id"`
}

type Member_point_redeem struct {
	MemberCode       string    `json:"member_code" binding:"required"`
	AuditDate        time.Time `json:"audit_date" binding:"required"`
	MemberTypeCode   string    `json:"member_type_code" binding:"required"`
	Point            float64   `json:"point" binding:"required"`
	IsVoucher        uint8     `json:"is_voucher" binding:"required"`
	GiftSerialNumber string    `json:"gift_serial_number" binding:"required"`
	CreatedAt        time.Time `json:"created_at"`
	CreatedBy        string    `json:"created_by"`
	UpdatedAt        time.Time `json:"updated_at"`
	UpdatedBy        string    `json:"updated_by"`
	Id               uint64    `json:"id"`
}

type Notif_tp struct {
	TemplateId      int       `json:"template_id" binding:"required"`
	SourceId        uint64    `json:"source_id"`
	SourceCode      string    `json:"source_code"`
	ScheduleDate    time.Time `json:"schedule_date"`
	SentTo          string    `json:"sent_to" binding:"required"`
	SentTo2         string    `json:"sent_to2"`
	TemplateSubject string    `json:"template_subject" binding:"required"`
	TemplateMessage string    `json:"template_message" binding:"required"`
	AttachFile      string    `json:"attach_file"`
	StatusCode      string    `json:"status_code" binding:"required"`
	SentDate        time.Time `json:"sent_date" binding:"required"`
	ResponApi       string    `json:"respon_api"`
	TypeCode        string    `json:"type_code" binding:"required"`
	FirstInsert     time.Time `json:"first_insert" binding:"required"`
	InsertBy        string    `json:"insert_by" binding:"required"`
	CreatedAt       time.Time `json:"created_at"`
	CreatedBy       string    `json:"created_by"`
	UpdatedAt       time.Time `json:"updated_at"`
	UpdatedBy       string    `json:"updated_by"`
	Id              uint64    `json:"id"`
}

type Notif_tp_cfg_init_template struct {
	Id         int       `json:"id"`
	EventCode  string    `json:"event_code"`
	Subject    string    `json:"subject" binding:"required"`
	Message    string    `json:"message" binding:"required"`
	PathFile   string    `json:"path_file"`
	FormatType string    `json:"format_type"`
	TypeCode   string    `json:"type_code"`
	CreatedAt  time.Time `json:"created_at"`
	CreatedBy  string    `json:"created_by"`
	UpdatedAt  time.Time `json:"updated_at"`
	UpdatedBy  string    `json:"updated_by"`
	IdLog      uint64    `json:"id_log" binding:"required"`
}

type Notif_tp_const_event struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type Notif_tp_const_variable struct {
	Name   string `json:"name" binding:"required"`
	IdSort int    `json:"id_sort" binding:"required"`
}

type Notif_tp_const_vendor struct {
	Code string `json:"code" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type Notification struct {
	TypeCode    string    `json:"type_code" binding:"required"`
	Message     string    `json:"message" binding:"required"`
	IsChs       uint8     `json:"is_chs" binding:"required"`
	IsCpos      uint8     `json:"is_cpos" binding:"required"`
	IsCas       uint8     `json:"is_cas" binding:"required"`
	IsCams      uint8     `json:"is_cams" binding:"required"`
	DateStart   time.Time `json:"date_start" binding:"required"`
	DateEnd     time.Time `json:"date_end" binding:"required"`
	PostingDate time.Time `json:"posting_date" binding:"required"`
	IsAllDay    uint8     `json:"is_all_day" binding:"required"`
	IsSystem    uint8     `json:"is_system" binding:"required"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   string    `json:"created_by"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedBy   string    `json:"updated_by"`
	Id          uint64    `json:"id"`
}

type One_time_password struct {
	Code             string    `json:"code" binding:"required"`
	Password         string    `json:"password" binding:"required"`
	PasswordTemp     string    `json:"password_temp" binding:"required"`
	UserGroupCode    string    `json:"user_group_code" binding:"required"`
	CreatedDate      time.Time `json:"created_date" binding:"required"`
	CreatedBy        string    `json:"created_by"`
	Interval         int       `json:"interval" binding:"required"`
	StatusCode       string    `json:"status_code" binding:"required"`
	SmsTo            string    `json:"sms_to"`
	MailTo           string    `json:"mail_to"`
	ChangeStatusDate time.Time `json:"change_status_date"`
	ChangeStatusBy   string    `json:"change_status_by"`
	Id               uint64    `json:"id"`
}

type Pabx_smdr struct {
	Date           time.Time `json:"date" binding:"required"`
	PostingDate    time.Time `json:"posting_date" binding:"required"`
	Extension      string    `json:"extension" binding:"required"`
	CoLine         string    `json:"co_line" binding:"required"`
	DialNumber     string    `json:"dial_number" binding:"required"`
	Duration       time.Time `json:"duration" binding:"required"`
	DurationSecond int       `json:"duration_second" binding:"required"`
	Charge         float64   `json:"charge" binding:"required"`
	Code           string    `json:"code" binding:"required"`
	CreatedAt      time.Time `json:"created_at"`
	CreatedBy      string    `json:"created_by"`
	UpdatedAt      time.Time `json:"updated_at"`
	UpdatedBy      string    `json:"updated_by"`
	Id             uint64    `json:"id"`
}

type Phone_book struct {
	TypeCode       string    `json:"type_code" binding:"required"`
	TitleCode      string    `json:"title_code"`
	FullName       string    `json:"full_name" binding:"required"`
	Phone1         string    `json:"phone1"`
	Phone2         string    `json:"phone2"`
	Facebook       string    `json:"facebook"`
	Twitter        string    `json:"twitter"`
	YahooMessenger string    `json:"yahoo_messenger"`
	PinBb          string    `json:"pin_bb"`
	Email          string    `json:"email"`
	Website        string    `json:"website"`
	Street         string    `json:"street"`
	City           string    `json:"city"`
	CountryCode    string    `json:"country_code"`
	StateCode      string    `json:"state_code"`
	PostalCode     string    `json:"postal_code"`
	JobTitle       string    `json:"job_title"`
	CompanyCode    string    `json:"company_code"`
	Remark         string    `json:"remark"`
	IsActive       uint8     `json:"is_active" binding:"required"`
	CreatedAt      time.Time `json:"created_at"`
	CreatedBy      string    `json:"created_by"`
	UpdatedAt      time.Time `json:"updated_at"`
	UpdatedBy      string    `json:"updated_by"`
	Id             uint64    `json:"id"`
}

type Pos_captain_order struct {
	ReservationNumber  uint64    `json:"reservation_number"`
	OutletCode         string    `json:"outlet_code" binding:"required"`
	TableNumber        string    `json:"table_number"`
	WaitressCode       string    `json:"waitress_code"`
	CustomerCode       string    `json:"customer_code"`
	MemberCode         string    `json:"member_code"`
	TitleCode          string    `json:"title_code"`
	FullName           string    `json:"full_name"`
	Adult              int       `json:"adult"`
	Child              int       `json:"child"`
	DocumentNumber     string    `json:"document_number"`
	Remark             string    `json:"remark"`
	AuditDate          time.Time `json:"audit_date" binding:"required"`
	PostingDate        time.Time `json:"posting_date" binding:"required"`
	MarketCode         string    `json:"market_code" binding:"required"`
	CompanyCode        string    `json:"company_code"`
	SalesCode          string    `json:"sales_code"`
	TimeSegmentCode    string    `json:"time_segment_code" binding:"required"`
	TypeCode           string    `json:"type_code" binding:"required"`
	ComplimentTypeCode string    `json:"compliment_type_code" binding:"required"`
	SubDepartmentCode  string    `json:"sub_department_code"`
	IsOpen             uint8     `json:"is_open" binding:"required"`
	IsPrinted          uint8     `json:"is_printed" binding:"required"`
	IsCancel           uint8     `json:"is_cancel" binding:"required"`
	CancelledAt        time.Time `json:"cancelled_at"`
	CancelledBy        string    `json:"cancelled_by"`
	CancelReason       string    `json:"cancel_reason"`
	CreatedAt          time.Time `json:"created_at"`
	CreatedBy          string    `json:"created_by"`
	UpdatedAt          time.Time `json:"updated_at"`
	UpdatedBy          string    `json:"updated_by"`
	Id                 uint64    `json:"id" gorm:"primaryKey"`
}

type Pos_captain_order_transaction struct {
	CaptainOrderId       uint64    `json:"captain_order_id" binding:"required"`
	InventoryCode        string    `json:"inventory_code" binding:"required"`
	TenanCode            string    `json:"tenan_code" binding:"required"`
	SeatNumber           int       `json:"seat_number" binding:"required"`
	SpaRoomNumber        string    `json:"spa_room_number"`
	SpaStartDate         time.Time `json:"spa_start_date"`
	SpaEndDate           time.Time `json:"spa_end_date"`
	ProductCode          string    `json:"product_code" binding:"required"`
	AccountCode          string    `json:"account_code" binding:"required"`
	Description          string    `json:"description" binding:"required"`
	Quantity             float64   `json:"quantity" binding:"required"`
	QuantityPrinted      float64   `json:"quantity_printed" binding:"required"`
	QuantityPrintedCheck float64   `json:"quantity_printed_check" binding:"required"`
	PricePurchase        float64   `json:"price_purchase" binding:"required"`
	PriceOriginal        float64   `json:"price_original" binding:"required"`
	Price                float64   `json:"price" binding:"required"`
	Discount             float64   `json:"discount"`
	DiscountTemp         float64   `json:"discount_temp" binding:"required"`
	Tax                  float64   `json:"tax"`
	Service              float64   `json:"service"`
	DefaultCurrencyCode  string    `json:"default_currency_code" binding:"required"`
	CurrencyCode         string    `json:"currency_code" binding:"required"`
	ExchangeRate         float64   `json:"exchange_rate" binding:"required"`
	Remark               string    `json:"remark"`
	TypeCode             string    `json:"type_code" binding:"required"`
	AuditDate            time.Time `json:"audit_date" binding:"required"`
	PostingDate          time.Time `json:"posting_date" binding:"required"`
	CompanyCode          string    `json:"company_code"`
	CompanyCode2         string    `json:"company_code2"`
	CardBankCode         string    `json:"card_bank_code"`
	CardTypeCode         string    `json:"card_type_code"`
	CardCharge           float64   `json:"card_charge" binding:"required"`
	CardNumber           string    `json:"card_number"`
	CardHolder           string    `json:"card_holder"`
	ValidMonth           string    `json:"valid_month"`
	ValidYear            string    `json:"valid_year"`
	FolioTransfer        uint64    `json:"folio_transfer" binding:"required"`
	SubFolioTransfer     string    `json:"sub_folio_transfer"`
	IsCompliment         uint8     `json:"is_compliment" binding:"required"`
	IsFree               uint8     `json:"is_free" binding:"required"`
	IsRemove             uint8     `json:"is_remove" binding:"required"`
	RemoveDate           time.Time `json:"remove_date"`
	RemoveBy             string    `json:"remove_by"`
	Shift                string    `json:"shift" binding:"required"`
	LogShiftId           uint64    `json:"log_shift_id" binding:"required"`
	CreatedAt            time.Time `json:"created_at"`
	CreatedBy            string    `json:"created_by"`
	UpdatedAt            time.Time `json:"updated_at"`
	UpdatedBy            string    `json:"updated_by"`
	Id                   uint64    `json:"id"`
}

type Pos_cfg_init_discount_limit struct {
	OutletCode    string    `json:"outlet_code" binding:"required"`
	UserGroupCode string    `json:"user_group_code"`
	MaxDiscount   float64   `json:"max_discount"`
	CreatedAt     time.Time `json:"created_at"`
	CreatedBy     string    `json:"created_by"`
	UpdatedAt     time.Time `json:"updated_at"`
	UpdatedBy     string    `json:"updated_by"`
	Id            uint64    `json:"id"`
}

type Pos_cfg_init_market struct {
	Code      string    `json:"code" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	IdSort    int       `json:"id_sort"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Pos_cfg_init_member_outlet_discount struct {
	Code            string    `json:"code" binding:"required"`
	Name            string    `json:"name" binding:"required"`
	OutletCode      string    `json:"outlet_code" binding:"required"`
	MinimumSale     float64   `json:"minimum_sale"`
	MaximumDiscount float64   `json:"maximum_discount"`
	IsForAllProduct uint8     `json:"is_for_all_product"`
	DiscountPercent float64   `json:"discount_percent"`
	CreatedAt       time.Time `json:"created_at"`
	CreatedBy       string    `json:"created_by"`
	UpdatedAt       time.Time `json:"updated_at"`
	UpdatedBy       string    `json:"updated_by"`
	Id              uint64    `json:"id"`
}

type Pos_cfg_init_member_outlet_discount_detail struct {
	MemberOutletDiscountCode string    `json:"member_outlet_discount_code" binding:"required"`
	OutletCode               string    `json:"outlet_code"  binding:"required"`
	ProductCode              string    `json:"product_code"  binding:"required"`
	DiscountPercent          float64   `json:"discount_percent"`
	CreatedAt                time.Time `json:"created_at"`
	CreatedBy                string    `json:"created_by"`
	UpdatedAt                time.Time `json:"updated_at"`
	UpdatedBy                string    `json:"updated_by"`
	Id                       uint64    `json:"id"`
}

type Pos_cfg_init_member_product_discount struct {
	OutletCode      string    `json:"outlet_code" binding:"required"`
	MemberCode      string    `json:"member_code" binding:"required"`
	ProductCode     string    `json:"product_code" binding:"required"`
	DiscountPercent float64   `json:"discount_percent"`
	CreatedAt       time.Time `json:"created_at"`
	CreatedBy       string    `json:"created_by"`
	UpdatedAt       time.Time `json:"updated_at"`
	UpdatedBy       string    `json:"updated_by"`
	Id              uint64    `json:"id"`
}

type Pos_cfg_init_outlet struct {
	Code                       string    `json:"code" binding:"required"`
	Name                       string    `json:"name" binding:"required"`
	CheckPrefix                string    `json:"check_prefix" binding:"required"`
	SubDepartmentCode          string    `json:"sub_department_code" binding:"required"`
	StoreCode                  string    `json:"store_code" binding:"required"`
	TaxAndServiceCode          string    `json:"tax_and_service_code" binding:"required"`
	TaxAndServiceCodeSpecial   string    `json:"tax_and_service_code_special" binding:"required"`
	PrinterCodeSpecial         string    `json:"printer_code_special"`
	PrinterCodeSpecialBeverage string    `json:"printer_code_special_beverage"`
	PrinterCodeCheck           string    `json:"printer_code_check"`
	AccountCode                string    `json:"account_code"`
	CompanyCode                string    `json:"company_code"`
	CommissionTypeCode         string    `json:"commission_type_code"`
	CommissionPercent          float64   `json:"commission_percent"`
	IdSort                     int       `json:"id_sort"`
	IsActive                   uint8     `json:"is_active"`
	Foto                       []byte    `json:"foto"`
	IsCheckAvailableTable      uint8     `json:"is_check_available_table"`
	IsForIptv                  uint8     `json:"is_for_iptv"`
	ImageLink                  string    `json:"image_link"`
	CreatedAt                  time.Time `json:"created_at"`
	CreatedBy                  string    `json:"created_by"`
	UpdatedAt                  time.Time `json:"updated_at"`
	UpdatedBy                  string    `json:"updated_by"`
	Id                         uint64    `json:"id"`
}

type Pos_cfg_init_payment_group struct {
	AccountCode string    `json:"account_code" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	OutletCode  string    `json:"outlet_code"`
	IdSort      int       `json:"id_sort"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   string    `json:"created_by"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedBy   string    `json:"updated_by"`
	Id          uint64    `json:"id"`
}

type Pos_cfg_init_product struct {
	Code              string    `json:"code" binding:"required"`
	Name              string    `json:"name" binding:"required"`
	Barcode           string    `json:"barcode"`
	Description       string    `json:"description"`
	CategoryCode      string    `json:"category_code" binding:"required"`
	GroupCode         string    `json:"group_code" binding:"required"`
	OutletCode        string    `json:"outlet_code" binding:"required"`
	TenanCode         string    `json:"tenan_code"`
	PackageCode       string    `json:"package_code"`
	TaxAndServiceCode string    `json:"tax_and_service_code"`
	PrinterCode       string    `json:"printer_code"`
	PrinterCode2      string    `json:"printer_code2"`
	Price             float64   `json:"price"`
	DisableDiscount   uint8     `json:"disable_discount"`
	Discount          float64   `json:"discount"`
	EstimationCost    float64   `json:"estimation_cost"`
	Buy               int       `json:"buy" binding:"required"`
	Free              int       `json:"free"`
	IsActive          uint8     `json:"is_active"`
	IsForIptv         uint8     `json:"is_for_iptv"`
	IsUsingSpaRoom    uint8     `json:"is_using_spa_room"`
	IsSold            uint8     `json:"is_sold"`
	ImageLink         string    `json:"image_link"`
	CreatedAt         time.Time `json:"created_at"`
	CreatedBy         string    `json:"created_by"`
	UpdatedAt         time.Time `json:"updated_at"`
	UpdatedBy         string    `json:"updated_by"`
	Id                uint64    `json:"id"`
}

type Pos_cfg_init_product_category struct {
	Code        string    `json:"code" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"`
	OutletCode  string    `json:"outlet_code" binding:"required"`
	IdSort      int       `json:"id_sort"`
	IsForIptv   uint8     `json:"is_for_iptv"`
	ImageLink   string    `json:"image_link"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   string    `json:"created_by"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedBy   string    `json:"updated_by"`
	Id          uint64    `json:"id"`
}

type Pos_cfg_init_product_group struct {
	Code          string    `json:"code" binding:"required"`
	Name          string    `json:"name" binding:"required"`
	ItemGroupCode string    `json:"item_group_code"`
	AccountCode   string    `json:"account_code" binding:"required"`
	CreatedAt     time.Time `json:"created_at"`
	CreatedBy     string    `json:"created_by"`
	UpdatedAt     time.Time `json:"updated_at"`
	UpdatedBy     string    `json:"updated_by"`
	Id            uint64    `json:"id"`
}

type Pos_cfg_init_room_boy struct {
	Code       string    `json:"code" binding:"required"`
	OutletCode string    `json:"outlet_code"`
	Name       string    `json:"name" binding:"required"`
	CreatedAt  time.Time `json:"created_at"`
	CreatedBy  string    `json:"created_by"`
	UpdatedAt  time.Time `json:"updated_at"`
	UpdatedBy  string    `json:"updated_by"`
	Id         uint64    `json:"id"`
}

type Pos_cfg_init_spa_room struct {
	Number      string    `json:"number" binding:"required"`
	OutletCode  string    `json:"outlet_code" binding:"required"`
	Description string    `json:"description"`
	Left        int       `json:"left"`
	Top         int       `json:"top"`
	Width       int       `json:"width"`
	Height      int       `json:"height"`
	IdSort      int       `json:"id_sort"`
	IsActive    uint8     `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   string    `json:"created_by"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedBy   string    `json:"updated_by"`
	Id          uint64    `json:"id"`
}

type Pos_cfg_init_table struct {
	Number      string    `json:"number" binding:"required"`
	OutletCode  string    `json:"outlet_code" binding:"required"`
	TypeCode    string    `json:"type_code"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date"`
	StatusCode  string    `json:"status_code"`
	Left        *int      `json:"left"`
	Top         *int      `json:"top"`
	Width       *int      `json:"width"`
	Height      *int      `json:"height"`
	IdSort      int       `json:"id_sort"`
	IsActive    uint8     `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   string    `json:"created_by"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedBy   string    `json:"updated_by"`
	Id          uint64    `json:"id"`
}

type Pos_cfg_init_table_type struct {
	Code      string    `json:"code" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	IdSort    int       `json:"id_sort"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Pos_cfg_init_tenan struct {
	Code       string    `json:"code" binding:"required"`
	Name       string    `json:"name" binding:"required"`
	OwnerCode  string    `json:"owner_code" binding:"required"`
	Commission float64   `json:"commission"`
	Street     string    `json:"street"`
	City       string    `json:"city"`
	Phone1     string    `json:"phone1"`
	Phone2     string    `json:"phone2"`
	Fax        string    `json:"fax"`
	Email      string    `json:"email"`
	CreatedAt  time.Time `json:"created_at"`
	CreatedBy  string    `json:"created_by"`
	UpdatedAt  time.Time `json:"updated_at"`
	UpdatedBy  string    `json:"updated_by"`
	Id         uint64    `json:"id"`
}

type Pos_cfg_init_therapist_fingerprint struct {
	TherapistCode       string    `json:"therapist_code" binding:"required"`
	FingerprintTemplate []byte    `json:"fingerprint_template" binding:"required"`
	CreatedAt           time.Time `json:"created_at"`
	CreatedBy           string    `json:"created_by"`
	UpdatedAt           time.Time `json:"updated_at"`
	UpdatedBy           string    `json:"updated_by"`
	Id                  uint64    `json:"id"`
}

type Pos_cfg_init_waitress struct {
	Code       string    `json:"code" binding:"required"`
	Name       string    `json:"name" binding:"required"`
	OutletCode string    `json:"outlet_code"`
	IsActive   uint8     `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	CreatedBy  string    `json:"created_by"`
	UpdatedAt  time.Time `json:"updated_at"`
	UpdatedBy  string    `json:"updated_by"`
	Id         uint64    `json:"id"`
}

type Pos_check struct {
	Number             string    `json:"number" binding:"required"`
	TypeCode           string    `json:"type_code" binding:"required"`
	CaptainOrderId     uint64    `json:"captain_order_id" binding:"required"`
	FolioNumber        uint64    `json:"folio_number" binding:"required"`
	ContactPersonId    uint64    `json:"contact_person_id" binding:"required"`
	OutletCode         string    `json:"outlet_code" binding:"required"`
	TableNumber        string    `json:"table_number"`
	WaitressCode       string    `json:"waitress_code"`
	MemberCode         string    `json:"member_code"`
	ComplimentTypeCode string    `json:"compliment_type_code" binding:"required"`
	SubDepartmentCode  string    `json:"sub_department_code"`
	Remark             string    `json:"remark"`
	AuditDate          time.Time `json:"audit_date" binding:"required"`
	PostingDate        time.Time `json:"posting_date" binding:"required"`
	MarketCode         string    `json:"market_code" binding:"required"`
	TimeSegmentCode    string    `json:"time_segment_code" binding:"required"`
	Void               uint8     `json:"void" binding:"required"`
	VoidDate           time.Time `json:"void_date" binding:"required"`
	VoidBy             string    `json:"void_by"`
	VoidReason         string    `json:"void_reason"`
	CreatedAt          time.Time `json:"created_at"`
	CreatedBy          string    `json:"created_by"`
	UpdatedAt          time.Time `json:"updated_at"`
	UpdatedBy          string    `json:"updated_by"`
	Id                 uint64    `json:"id"`
}

type Pos_check_transaction struct {
	CheckNumber               string    `json:"check_number" binding:"required"`
	CaptainOrderTransactionId uint64    `json:"captain_order_transaction_id" binding:"required"`
	SubFolioId                uint64    `json:"sub_folio_id" binding:"required"`
	InventoryCode             string    `json:"inventory_code" binding:"required"`
	TenanCode                 string    `json:"tenan_code" binding:"required"`
	SeatNumber                int       `json:"seat_number" binding:"required"`
	SpaRoomNumber             string    `json:"spa_room_number"`
	SpaStartDate              time.Time `json:"spa_start_date"`
	SpaEndDate                time.Time `json:"spa_end_date"`
	ProductCode               string    `json:"product_code" binding:"required"`
	PricePurchase             float64   `json:"price_purchase" binding:"required"`
	PriceOriginal             float64   `json:"price_original" binding:"required"`
	Price                     float64   `json:"price" binding:"required"`
	Discount                  float64   `json:"discount"`
	EstimationCost            float64   `json:"estimation_cost" binding:"required"`
	Tax                       float64   `json:"tax"`
	Service                   float64   `json:"service"`
	CompanyCode               string    `json:"company_code"`
	CompanyCode2              string    `json:"company_code2"`
	CardCharge                float64   `json:"card_charge"`
	FolioTransfer             uint64    `json:"folio_transfer" binding:"required"`
	IsCompliment              uint8     `json:"is_compliment" binding:"required"`
	IsFree                    uint8     `json:"is_free" binding:"required"`
	CreatedAt                 time.Time `json:"created_at"`
	CreatedBy                 string    `json:"created_by"`
	UpdatedAt                 time.Time `json:"updated_at"`
	UpdatedBy                 string    `json:"updated_by"`
	Id                        uint64    `json:"id" gorm:"primaryKey"`
}

type Pos_const_check_type struct {
	Code string `json:"code" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type Pos_const_compliment_type struct {
	Code string `json:"code" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type Pos_const_discount struct {
	Id    int     `json:"id"`
	Value float64 `json:"value"`
}

type Pos_const_time_segment struct {
	Code string `json:"code" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type Pos_information struct {
	Code        string `json:"code" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Street      string `json:"street" binding:"required"`
	City        string `json:"city" binding:"required"`
	CountryCode string `json:"country_code"`
	StateCode   string `json:"state_code"`
	PostalCode  string `json:"postal_code"`
	Phone1      string `json:"phone1"`
	Phone2      string `json:"phone2"`
	Fax         string `json:"fax"`
	Email       string `json:"email"`
	Website     string `json:"website"`
	Foto        []byte `json:"foto"`
}

type Pos_iptv_menu_order struct {
	RoomNumber  string    `json:"room_number" binding:"required"`
	OutletCode  string    `json:"outlet_code" binding:"required"`
	ProductCode string    `json:"product_code" binding:"required"`
	Quantity    float64   `json:"quantity" binding:"required"`
	Remark      string    `json:"remark"`
	IsProcessed uint8     `json:"is_processed" binding:"required"`
	IsCancelled uint8     `json:"is_cancelled" binding:"required"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   string    `json:"created_by"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedBy   string    `json:"updated_by"`
	Id          uint64    `json:"id"`
}

type Pos_member struct {
	Code      string    `json:"code" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Pos_product_costing struct {
	CheckNumber        string    `json:"check_number" binding:"required"`
	CheckTransactionId uint64    `json:"check_transaction_id" binding:"required"`
	CostingNumber      string    `json:"costing_number" binding:"required"`
	CostingDetailId    uint64    `json:"costing_detail_id" binding:"required"`
	ProductCode        string    `json:"product_code" binding:"required"`
	StoreCode          string    `json:"store_code" binding:"required"`
	ItemCode           string    `json:"item_code" binding:"required"`
	Quantity           float64   `json:"quantity" binding:"required"`
	UomCode            string    `json:"uom_code" binding:"required"`
	BasicQuantity      float64   `json:"basic_quantity" binding:"required"`
	BasicUomCode       string    `json:"basic_uom_code" binding:"required"`
	CostingQuantity    float64   `json:"costing_quantity" binding:"required"`
	CreatedAt          time.Time `json:"created_at"`
	CreatedBy          string    `json:"created_by"`
	UpdatedAt          time.Time `json:"updated_at"`
	UpdatedBy          string    `json:"updated_by"`
	Id                 uint64    `json:"id"`
}

type Pos_report struct {
	Code        int       `json:"code" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	ReportQuery string    `json:"report_query" binding:"required"`
	ParentId    uint64    `json:"parent_id" binding:"required"`
	IsSystem    uint8     `json:"is_system" binding:"required"`
	IdSort      int       `json:"id_sort" binding:"required"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   string    `json:"created_by"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedBy   string    `json:"updated_by"`
	Id          uint64    `json:"id"`
}

type Pos_report_default_field struct {
	ReportCode int    `json:"report_code" binding:"required"`
	FieldQuery string `json:"field_query" binding:"required"`
	FieldName  string `json:"field_name" binding:"required"`
	IdSort     int    `json:"id_sort" binding:"required"`
}

type Pos_report_group_field struct {
	TemplateId uint64 `json:"template_id" binding:"required"`
	FieldName  string `json:"field_name" binding:"required"`
	IdSort     int    `json:"id_sort" binding:"required"`
}

type Pos_report_grouping_field struct {
	TemplateId uint64 `json:"template_id" binding:"required"`
	FieldName  string `json:"field_name" binding:"required"`
	IdSort     int    `json:"id_sort" binding:"required"`
}

type Pos_report_order_field struct {
	TemplateId  uint64 `json:"template_id" binding:"required"`
	FieldName   string `json:"field_name" binding:"required"`
	IsAscending uint8  `json:"is_ascending" binding:"required"`
	IdSort      int    `json:"id_sort" binding:"required"`
}

type Pos_report_template struct {
	ReportCode       int       `json:"report_code" binding:"required"`
	Name             string    `json:"name" binding:"required"`
	GroupLevel       int       `json:"group_level" binding:"required"`
	HeaderRemark     string    `json:"header_remark"`
	ShowFooter       uint8     `json:"show_footer" binding:"required"`
	PaperSize        int       `json:"paper_size" binding:"required"`
	PaperWidth       float64   `json:"paper_width" binding:"required"`
	PaperHeight      float64   `json:"paper_height" binding:"required"`
	IsPortrait       uint8     `json:"is_portrait" binding:"required"`
	HeaderRowHeight  int       `json:"header_row_height" binding:"required"`
	RowHeight        int       `json:"row_height" binding:"required"`
	HorizontalBorder uint8     `json:"horizontal_border" binding:"required"`
	VerticalBorder   uint8     `json:"vertical_border" binding:"required"`
	SignName1        string    `json:"sign_name1"`
	SignPosition1    string    `json:"sign_position1"`
	SignName2        string    `json:"sign_name2"`
	SignPosition2    string    `json:"sign_position2"`
	SignName3        string    `json:"sign_name3"`
	SignPosition3    string    `json:"sign_position3"`
	SignName4        string    `json:"sign_name4"`
	SignPosition4    string    `json:"sign_position4"`
	IsDefault        uint8     `json:"is_default" binding:"required"`
	IsSystem         uint8     `json:"is_system" binding:"required"`
	IdSort           int       `json:"id_sort" binding:"required"`
	CreatedAt        time.Time `json:"created_at"`
	CreatedBy        string    `json:"created_by"`
	UpdatedAt        time.Time `json:"updated_at"`
	UpdatedBy        string    `json:"updated_by"`
	Id               uint64    `json:"id"`
}

type Pos_report_template_field struct {
	TemplateId      uint64 `json:"template_id" binding:"required"`
	FieldName       string `json:"field_name" binding:"required"`
	HeaderName      string `json:"header_name" binding:"required"`
	FooterType      int    `json:"footer_type" binding:"required"`
	FormatCode      int    `json:"format_code" binding:"required"`
	Width           int    `json:"width" binding:"required"`
	Font            int    `json:"font" binding:"required"`
	FontSize        int    `json:"font_size" binding:"required"`
	FontColor       int    `json:"font_color" binding:"required"`
	Alignment       string `json:"alignment" binding:"required"`
	HeaderFontSize  int    `json:"header_font_size" binding:"required"`
	HeaderFontColor int    `json:"header_font_color" binding:"required"`
	HeaderAlignment string `json:"header_alignment" binding:"required"`
	IdSort          int    `json:"id_sort" binding:"required"`
}

type Pos_reservation struct {
	Number             uint64    `json:"number" binding:"required"`
	CheckIn            time.Time `json:"check_in"`
	CheckInRes         time.Time `json:"check_in_res"`
	CheckOut           time.Time `json:"check_out"`
	CheckOutRes        time.Time `json:"check_out_res"`
	OutletCode         string    `json:"outlet_code" binding:"required"`
	TableNumber        string    `json:"table_number"`
	WaitressCode       string    `json:"waitress_code"`
	CustomerCode       string    `json:"customer_code"`
	MemberCode         string    `json:"member_code"`
	TitleCode          string    `json:"title_code"`
	FullName           string    `json:"full_name"`
	Adult              int       `json:"adult"`
	Child              int       `json:"child"`
	DocumentNumber     string    `json:"document_number"`
	Remark             string    `json:"remark"`
	AuditDate          time.Time `json:"audit_date" binding:"required"`
	PostingDate        time.Time `json:"posting_date" binding:"required"`
	ComplimentTypeCode string    `json:"compliment_type_code" binding:"required"`
	SubDepartmentCode  string    `json:"sub_department_code"`
	StatusCode         string    `json:"status_code" binding:"required"`
	CancelAuditDate    time.Time `json:"cancel_audit_date"`
	CancelDate         time.Time `json:"cancel_date"`
	CancelBy           string    `json:"cancel_by"`
	CancelReason       string    `json:"cancel_reason"`
	CreatedAt          time.Time `json:"created_at"`
	CreatedBy          string    `json:"created_by"`
	UpdatedAt          time.Time `json:"updated_at"`
	UpdatedBy          string    `json:"updated_by"`
}

type Pos_reservation_table struct {
	ReservationNumber uint64    `json:"reservation_number"`
	Start             time.Time `json:"start"`
	Finish            time.Time `json:"finish"`
	TableNumber       string    `json:"table_number"`
	ParentId          uint64    `json:"parent_id"`
	EventType         int       `json:"event_type"`
	Options           int       `json:"options"`
	Caption           string    `json:"caption"`
	RecurrenceIndex   int       `json:"recurrence_index"`
	RecurrenceInfo    string    `json:"recurrence_info"`
	Message           string    `json:"message"`
	ReminderDate      time.Time `json:"reminder_date"`
	ReminderMinutes   int       `json:"reminder_minutes"`
	State             int       `json:"state"`
	LabelColor        int       `json:"label_color"`
	ActualStart       time.Time `json:"actual_start"`
	ActualFinish      time.Time `json:"actual_finish"`
	SincId            string    `json:"sinc_id"`
	ReminderResource  string    `json:"reminder_resource"`
	BlockType         string    `json:"block_type"`
	CaptainOrderId    uint64    `json:"captain_order_id"`
	CreatedAt         time.Time `json:"created_at"`
	CreatedBy         string    `json:"created_by"`
	UpdatedAt         time.Time `json:"updated_at"`
	UpdatedBy         string    `json:"updated_by"`
	Id                uint64    `json:"id"`
}

type Pos_table_unavailable struct {
	TableNumber string    `json:"table_number" binding:"required"`
	OutletCode  string    `json:"outlet_code"`
	FromDate    time.Time `json:"from_date" binding:"required"`
	ToDate      time.Time `json:"to_date" binding:"required"`
	StatusCode  string    `json:"status_code" binding:"required"`
	ReasonCode  string    `json:"reason_code" binding:"required"`
	Note        string    `json:"note"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   string    `json:"created_by"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedBy   string    `json:"updated_by"`
	Id          uint64    `json:"id"`
}

type Pos_user_group struct {
	Code                      string    `json:"code" binding:"required"`
	AccessForm                string    `json:"access_form" binding:"required"`
	AccessSpecial             string    `json:"access_special" binding:"required"`
	AccessTransactionTerminal string    `json:"access_transaction_terminal" binding:"required"`
	AccessTableView           string    `json:"access_table_view" binding:"required"`
	AccessReservation         string    `json:"access_reservation" binding:"required"`
	CreatedAt                 time.Time `json:"created_at"`
	CreatedBy                 string    `json:"created_by"`
	UpdatedAt                 time.Time `json:"updated_at"`
	UpdatedBy                 string    `json:"updated_by"`
	Id                        uint64    `json:"id"`
}

type Pos_user_group_outlet struct {
	UserGroupCode string    `json:"user_group_code" binding:"required"`
	OutletCode    string    `json:"outlet_code" binding:"required"`
	CreatedAt     time.Time `json:"created_at"`
	CreatedBy     string    `json:"created_by"`
	UpdatedAt     time.Time `json:"updated_at"`
	UpdatedBy     string    `json:"updated_by"`
	Id            uint64    `json:"id"`
}

type Proforma_invoice_detail struct {
	Id                uint64    `json:"id" gorm:"primaryKey"`
	ReservationNumber uint64    `json:"reservation_number" binding:"required"`
	ArrivalDate       time.Time `json:"arrival_date" binding:"required"`
	DepartureDate     time.Time `json:"departure_date" binding:"required"`
	Datex             time.Time `json:"datex" binding:"required"`
	RoomTypeCode      string    `json:"room_type_code" binding:"required"`
	RoomRate          float64   `json:"room_rate" binding:"required"`
	IsWeekend         string    `json:"is_weekend" binding:"required"`
	ChargeFrequency   string    `json:"charge_frequency" binding:"required"`
	Userid            string    `json:"userid" binding:"required"`
}

type Receipt struct {
	Number      string    `json:"number" binding:"required" gorm:"primaryKey"`
	ReceiveFrom string    `json:"receive_from" binding:"required"`
	Amount      float64   `json:"amount" binding:"required"`
	IssuedDate  time.Time `json:"issued_date" binding:"required"`
	ForPayment  string    `json:"for_payment" binding:"required"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   string    `json:"created_by"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedBy   string    `json:"updated_by"`
	Id          uint64    `json:"id"`
}

type Report struct {
	Code        int       `json:"code" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	ReportQuery string    `json:"report_query" binding:"required"`
	ParentId    uint64    `json:"parent_id" binding:"required"`
	IsSystem    uint8     `json:"is_system" binding:"required"`
	IdSort      int       `json:"id_sort" binding:"required"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   string    `json:"created_by"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedBy   string    `json:"updated_by"`
	Id          uint64    `json:"id"`
}

type Report_custom struct {
	Code          int    `json:"code" binding:"required"`
	ParentId      uint64 `json:"parent_id" binding:"required"`
	SystemCode    string `json:"system_code" binding:"required"`
	UserGroupCode string `json:"user_group_code" binding:"required"`
	Id            uint64 `json:"id" gorm:"primaryKey"`
}

type Report_custom_favorite struct {
	Code       int       `json:"code" binding:"required"`
	ParentId   uint64    `json:"parent_id" binding:"required"`
	SystemCode string    `json:"system_code" binding:"required"`
	IdSort     uint64    `json:"id_sort" binding:"required"`
	CreatedAt  time.Time `json:"created_at"`
	CreatedBy  string    `json:"created_by"`
	UpdatedAt  time.Time `json:"updated_at"`
	UpdatedBy  string    `json:"updated_by"`
}

type Report_default_field struct {
	ReportCode int    `json:"report_code" binding:"required"`
	FieldQuery string `json:"field_query" binding:"required"`
	FieldName  string `json:"field_name" binding:"required"`
	IdSort     int    `json:"id_sort" binding:"required"`
}

type Report_group_field struct {
	TemplateId uint64 `json:"template_id" binding:"required"`
	FieldName  string `json:"field_name" binding:"required"`
	IdSort     int    `json:"id_sort" binding:"required"`
}

type Report_grouping_field struct {
	TemplateId uint64 `json:"template_id" binding:"required"`
	FieldName  string `json:"field_name" binding:"required"`
	IdSort     int    `json:"id_sort" binding:"required"`
}

type Report_order_field struct {
	TemplateId  uint64 `json:"template_id" binding:"required"`
	FieldName   string `json:"field_name" binding:"required"`
	IsAscending uint8  `json:"is_ascending" binding:"required"`
	IdSort      int    `json:"id_sort" binding:"required"`
}

type Report_pivot_temp struct {
	ItemGroupCode   string  `json:"item_group_code" binding:"required"`
	FbStructureCode string  `json:"fb_structure_code" binding:"required"`
	PosMarketCode   string  `json:"pos_market_code" binding:"required"`
	D01             float64 `json:"d01"`
	D02             float64 `json:"d02"`
	D03             float64 `json:"d03"`
	D04             float64 `json:"d04"`
	D05             float64 `json:"d05"`
	D06             float64 `json:"d06"`
	D07             float64 `json:"d07"`
	D08             float64 `json:"d08"`
	D09             float64 `json:"d09"`
	D10             float64 `json:"d10"`
	D11             float64 `json:"d11"`
	D12             float64 `json:"d12"`
	D13             float64 `json:"d13"`
	D14             float64 `json:"d14"`
	D15             float64 `json:"d15"`
	D16             float64 `json:"d16"`
	D17             float64 `json:"d17"`
	D18             float64 `json:"d18"`
	D19             float64 `json:"d19"`
	D20             float64 `json:"d20"`
	Subtotal01      float64 `json:"subtotal01"`
	Subtotal02      float64 `json:"subtotal02"`
	Subtotal03      float64 `json:"subtotal03"`
	Subtotal04      float64 `json:"subtotal04"`
	Subtotal05      float64 `json:"subtotal05"`
	Total1          float64 `json:"total1"`
	Id              uint64  `json:"id"`
}

type Report_room_rate_structure_temp struct {
	RoomRateSubCategoryCode string  `json:"room_rate_sub_category_code" binding:"required"`
	RoomTypeCode            string  `json:"room_type_code" binding:"required"`
	RateAmount              float64 `json:"rate_amount" binding:"required"`
	Id                      uint64  `json:"id"`
}

type Report_room_sales struct {
	Id         int       `json:"id"`
	AuditDate  time.Time `json:"audit_date" binding:"required"`
	RoomType01 int       `json:"room_type01"`
	RoomType02 int       `json:"room_type02"`
	RoomType03 int       `json:"room_type03"`
	RoomType04 int       `json:"room_type04"`
	RoomType05 int       `json:"room_type05"`
	RoomType06 int       `json:"room_type06"`
	RoomType07 int       `json:"room_type07"`
	RoomType08 int       `json:"room_type08"`
	RoomType09 int       `json:"room_type09"`
	RoomType10 int       `json:"room_type10"`
	Adult      int       `json:"adult"`
	Child      int       `json:"child"`
	IpAddress  string    `json:"ip_address"`
}

type Report_template struct {
	ReportCode       int       `json:"report_code" binding:"required"`
	Name             string    `json:"name" binding:"required"`
	GroupLevel       *int      `json:"group_level" binding:"required"`
	HeaderRemark     string    `json:"header_remark"`
	ShowFooter       *uint8    `json:"show_footer" binding:"required"`
	ShowPageNumber   *uint8    `json:"show_page_number" binding:"required"`
	PaperSize        int       `json:"paper_size" binding:"required"`
	PaperWidth       float64   `json:"paper_width" binding:"required"`
	PaperHeight      float64   `json:"paper_height" binding:"required"`
	IsPortrait       *uint8    `json:"is_portrait" binding:"required"`
	HeaderRowHeight  int       `json:"header_row_height" binding:"required"`
	RowHeight        int       `json:"row_height" binding:"required"`
	HorizontalBorder *uint8    `json:"horizontal_border" binding:"required"`
	VerticalBorder   *uint8    `json:"vertical_border" binding:"required"`
	SignName1        string    `json:"sign_name1"`
	SignPosition1    string    `json:"sign_position1"`
	SignName2        string    `json:"sign_name2"`
	SignPosition2    string    `json:"sign_position2"`
	SignName3        string    `json:"sign_name3"`
	SignPosition3    string    `json:"sign_position3"`
	SignName4        string    `json:"sign_name4"`
	SignPosition4    string    `json:"sign_position4"`
	IsDefault        *uint8    `json:"is_default" binding:"required"`
	IsSystem         uint8     `json:"is_system" binding:"required"`
	IdSort           int       `json:"id_sort" binding:"required"`
	CreatedAt        time.Time `json:"created_at"`
	CreatedBy        string    `json:"created_by"`
	UpdatedAt        time.Time `json:"updated_at"`
	UpdatedBy        string    `json:"updated_by"`
	Id               uint64    `json:"id" gorm:"primaryKey"`
}

type Report_template_field struct {
	TemplateId      uint64 `json:"template_id" binding:"required"`
	FieldName       string `json:"field_name" binding:"required"`
	HeaderName      string `json:"header_name" binding:"required"`
	FooterType      int    `json:"footer_type" binding:"required"`
	FormatCode      int    `json:"format_code" binding:"required"`
	Width           int    `json:"width" binding:"required"`
	Font            int    `json:"font" binding:"required"`
	FontSize        int    `json:"font_size" binding:"required"`
	FontColor       string `json:"font_color" binding:"required"`
	Alignment       string `json:"alignment" binding:"required"`
	HeaderFontSize  int    `json:"header_font_size" binding:"required"`
	HeaderFontColor string `json:"header_font_color" binding:"required"`
	HeaderAlignment string `json:"header_alignment" binding:"required"`
	IdSort          int    `json:"id_sort" binding:"required"`
}

type Reservation struct {
	Number           uint64    `json:"number" gorm:"primaryKey"`
	ContactPersonId1 uint64    `json:"contact_person_id1" gorm:"index"`
	ContactPersonId2 uint64    `json:"contact_person_id2"`
	ContactPersonId3 uint64    `json:"contact_person_id3"`
	ContactPersonId4 uint64    `json:"contact_person_id4"`
	GuestDetailId    uint64    `json:"guest_detail_id"`
	GuestProfileId1  uint64    `json:"guest_profile_id1"`
	GuestProfileId2  uint64    `json:"guest_profile_id2"`
	GuestProfileId3  uint64    `json:"guest_profile_id3"`
	GuestProfileId4  uint64    `json:"guest_profile_id4"`
	GuestGeneralId   uint64    `json:"guest_general_id"`
	ReservationBy    string    `json:"reservation_by" binding:"required"`
	GroupCode        string    `json:"group_code"`
	MemberCode       string    `json:"member_code"`
	StatusCode       string    `json:"status_code" `
	StatusCode2      string    `json:"status_code2" `
	IsIncognito      *uint8    `json:"is_incognito"`
	IsLock           *uint8    `json:"is_lock"`
	IsFromAllotment  uint8     `json:"is_from_allotment"`
	BookingCode      string    `json:"booking_code"`
	OtaId            string    `json:"ota_id"`
	CmResStatus      string    `json:"cm_res_status"`
	CmRevId          string    `json:"cm_rev_id"`
	IsCmConfirmed    uint8     `json:"is_cm_confirmed"`
	ChangeStatusAt   time.Time `json:"change_status_at"`
	ChangeStatusBy   string    `json:"change_status_by"`
	AuditDate        time.Time `json:"audit_date"`
	CancelledAt      time.Time `json:"cancelled_at"`
	CancelledBy      string    `json:"cancelled_by"`
	CancelReason     string    `json:"cancel_reason"`
	CreatedAt        time.Time `json:"created_at"`
	CreatedBy        string    `json:"created_by"`
	UpdatedAt        time.Time `json:"updated_at"`
	UpdatedBy        string    `json:"updated_by"`
}

type Reservation_extra_charge struct {
	ReservationNumber   uint64    `json:"reservation_number" binding:"required"`
	PackageName         *string   `json:"package_name" binding:"required"`
	OutletCode          *string   `json:"outlet_code"`
	ProductCode         *string   `json:"product_code"`
	PackageCode         *string   `json:"package_code"`
	GroupCode           string    `json:"group_code" binding:"required"`
	SubDepartmentCode   string    `json:"sub_department_code" binding:"required"`
	AccountCode         string    `json:"account_code" binding:"required"`
	Quantity            float64   `json:"quantity" binding:"required"`
	Amount              float64   `json:"amount" binding:"required"`
	PerPax              *uint8    `json:"per_pax" binding:"required"`
	IncludeChild        *uint8    `json:"include_child" binding:"required"`
	TaxAndServiceCode   *string   `json:"tax_and_service_code"`
	ChargeFrequencyCode string    `json:"charge_frequency_code" binding:"required"`
	MaxPax              int       `json:"max_pax"`
	ExtraPax            *float64  `json:"extra_pax"`
	PerPaxExtra         *uint8    `json:"per_pax_extra" binding:"required"`
	CreatedAt           time.Time `json:"created_at"`
	CreatedBy           string    `json:"created_by"`
	UpdatedAt           time.Time `json:"updated_at"`
	UpdatedBy           string    `json:"updated_by"`
	Id                  uint64    `json:"id" gorm:"primaryKey"`
}

type Reservation_extra_charge_breakdown struct {
	ReservationExtraChargeId uint64    `json:"reservation_extra_charge_id" binding:"required"`
	OutletCode               *string   `json:"outlet_code"`
	ProductCode              *string   `json:"product_code"`
	SubDepartmentCode        string    `json:"sub_department_code" binding:"required"`
	AccountCode              string    `json:"account_code" binding:"required"`
	CompanyCode              *string   `json:"company_code"`
	Quantity                 float64   `json:"quantity" binding:"required"`
	IsAmountPercent          *uint8    `json:"is_amount_percent" binding:"required"`
	Amount                   float64   `json:"amount" binding:"required"`
	PerPax                   *uint8    `json:"per_pax" binding:"required"`
	IncludeChild             *uint8    `json:"include_child" binding:"required"`
	Remark                   *string   `json:"remark"`
	TaxAndServiceCode        *string   `json:"tax_and_service_code"`
	ChargeFrequencyCode      string    `json:"charge_frequency_code" binding:"required"`
	MaxPax                   int       `json:"max_pax"`
	ExtraPax                 *float64  `json:"extra_pax"`
	PerPaxExtra              *uint8    `json:"per_pax_extra" binding:"required"`
	CreatedAt                time.Time `json:"created_at"`
	CreatedBy                string    `json:"created_by"`
	UpdatedAt                time.Time `json:"updated_at"`
	UpdatedBy                string    `json:"updated_by"`
	Id                       uint64    `json:"id" gorm:"primaryKey"`
}

type Reservation_scheduled_rate struct {
	ReservationNumber uint64    `json:"reservation_number" binding:"required"`
	FromDate          time.Time `json:"from_date" binding:"required"`
	ToDate            time.Time `json:"to_date" binding:"required"`
	RoomRateCode      string    `json:"room_rate_code" binding:"required"`
	Rate              *float64  `json:"rate" binding:"required"`
	ComplimentHu      string    `json:"compliment_hu"`
	CreatedAt         time.Time `json:"created_at"`
	CreatedBy         string    `json:"created_by"`
	UpdatedAt         time.Time `json:"updated_at"`
	UpdatedBy         string    `json:"updated_by"`
	Id                uint64    `json:"id" gorm:"primaryKey"`
}

type Room_allotment struct {
	RoomNumber  string    `json:"room_number" binding:"required"`
	FromDate    time.Time `json:"from_date" binding:"required"`
	ToDate      time.Time `json:"to_date" binding:"required"`
	TypeCode    string    `json:"type_code" binding:"required"`
	CompanyCode string    `json:"company_code"`
	Note        string    `json:"note"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   string    `json:"created_by"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedBy   string    `json:"updated_by"`
	Id          uint64    `json:"id" gorm:"primaryKey"`
}

type Room_statistic struct {
	Date                    time.Time `json:"date" binding:"required"`
	TotalRoom               int       `json:"total_room" binding:"required"`
	OutOfOrder              int       `json:"out_of_order" binding:"required"`
	OfficeUse               int       `json:"office_use" binding:"required"`
	UnderConstruction       int       `json:"under_construction" binding:"required"`
	HouseUse                int       `json:"house_use" binding:"required"`
	Compliment              int       `json:"compliment" binding:"required"`
	RoomSold                int       `json:"room_sold" binding:"required"`
	DayUse                  int       `json:"day_use" binding:"required"`
	RevenueGross            float64   `json:"revenue_gross" binding:"required"`
	RevenueWithCompliment   float64   `json:"revenue_with_compliment" binding:"required"`
	RevenueNonPackage       float64   `json:"revenue_non_package" binding:"required"`
	RevenueNett             float64   `json:"revenue_nett" binding:"required"`
	Adult                   int       `json:"adult" binding:"required"`
	Child                   int       `json:"child" binding:"required"`
	AdultSold               int       `json:"adult_sold" binding:"required"`
	ChildSold               int       `json:"child_sold" binding:"required"`
	ChildDayUse             int       `json:"child_day_use" binding:"required"`
	AdultDayUse             int       `json:"adult_day_use" binding:"required"`
	AdultCompliment         int       `json:"adult_compliment" binding:"required"`
	ChildCompliment         int       `json:"child_compliment" binding:"required"`
	AdultHu                 int       `json:"adult_hu" binding:"required"`
	ChildHu                 int       `json:"child_hu" binding:"required"`
	PaxSingle               int       `json:"pax_single" binding:"required"`
	WalkIn                  int       `json:"walk_in" binding:"required"`
	WalkInForeign           int       `json:"walk_in_foreign" binding:"required"`
	CheckIn                 int       `json:"check_in" binding:"required"`
	PersonCheckIn           int       `json:"person_check_in" binding:"required"`
	CheckInTomorrow         int       `json:"check_in_tomorrow" binding:"required"`
	CheckInPersonTomorrow   int       `json:"check_in_person_tomorrow" binding:"required"`
	CheckInForeign          int       `json:"check_in_foreign" binding:"required"`
	Reservation             int       `json:"reservation" binding:"required"`
	CancelReservation       int       `json:"cancel_reservation" binding:"required"`
	NoShowReservation       int       `json:"no_show_reservation" binding:"required"`
	CheckOut                int       `json:"check_out" binding:"required"`
	PersonCheckOut          int       `json:"person_check_out" binding:"required"`
	EarlyCheckOut           int       `json:"early_check_out" binding:"required"`
	CheckOutTomorrow        int       `json:"check_out_tomorrow" binding:"required"`
	CheckOutPersonTomorrow  int       `json:"check_out_person_tomorrow" binding:"required"`
	BreakfastCover          int       `json:"breakfast_cover" binding:"required"`
	FoodCover               int       `json:"food_cover" binding:"required"`
	BeverageCover           int       `json:"beverage_cover" binding:"required"`
	BanquetCover            int       `json:"banquet_cover" binding:"required"`
	WeddingCover            int       `json:"wedding_cover" binding:"required"`
	GatheringCover          int       `json:"gathering_cover" binding:"required"`
	SegmentCoverBreakfast   int       `json:"segment_cover_breakfast" binding:"required"`
	SegmentCoverLunch       int       `json:"segment_cover_lunch" binding:"required"`
	SegmentCoverDinner      int       `json:"segment_cover_dinner" binding:"required"`
	SegmentCoverCoffeeBreak int       `json:"segment_cover_coffee_break" binding:"required"`
	RevenueBreakfast        float64   `json:"revenue_breakfast" binding:"required"`
	RevenueFood             float64   `json:"revenue_food" binding:"required"`
	RevenueBeverage         float64   `json:"revenue_beverage" binding:"required"`
	RevenueBanquet          float64   `json:"revenue_banquet" binding:"required"`
	RevenueWedding          float64   `json:"revenue_wedding"`
	RevenueGathering        float64   `json:"revenue_gathering"`
	GuestLedger             float64   `json:"guest_ledger" binding:"required"`
	GuestDeposit            float64   `json:"guest_deposit" binding:"required"`
	UnitCode                string    `json:"unit_code" binding:"required"`
	Id                      uint64    `json:"id" gorm:"primaryKey"`
	IdHolding               uint64    `json:"id_holding"`
}

type Room_status struct {
	AuditDate  time.Time `json:"audit_date" binding:"required"`
	RoomNumber string    `json:"room_number" binding:"required"`
	Status     string    `json:"status" binding:"required"`
	Id         uint64    `json:"id"`
}

type Room_unavailable struct {
	RoomNumber string    `json:"room_number" binding:"required"`
	StartDate  time.Time `json:"start_date" binding:"required"`
	EndDate    time.Time `json:"end_date" binding:"required"`
	StatusCode string    `json:"status_code" binding:"required"`
	ReasonCode string    `json:"reason_code" binding:"required"`
	Note       *string   `json:"note"`
	CreatedAt  time.Time `json:"created_at"`
	CreatedBy  string    `json:"created_by"`
	UpdatedAt  time.Time `json:"updated_at"`
	UpdatedBy  string    `json:"updated_by"`
	Id         uint64    `json:"id" gorm:"primaryKey"`
}

type Room_unavailable_history struct {
	AuditDate  time.Time `json:"audit_date" binding:"required"`
	RoomNumber string    `json:"room_number" binding:"required"`
	StatusCode string    `json:"status_code" binding:"required"`
	ReasonCode string    `json:"reason_code" binding:"required"`
	Note       string    `json:"note"`
	CreatedAt  time.Time `json:"created_at"`
	CreatedBy  string    `json:"created_by"`
	UpdatedAt  time.Time `json:"updated_at"`
	UpdatedBy  string    `json:"updated_by"`
	Id         uint64    `json:"id"`
}

type Sal_activity struct {
	Code        string    `json:"code" binding:"required"`
	CompanyCode string    `json:"company_code" binding:"required"`
	Date        time.Time `json:"date" binding:"required"`
	StatusCode  string    `json:"status_code" binding:"required"`
	SegmentCode string    `json:"segment_code"`
	SourceCode  string    `json:"source_code" binding:"required"`
	AssignedBy  string    `json:"assigned_by" binding:"required"`
	Subject     string    `json:"subject" binding:"required"`
	Notes       string    `json:"notes" binding:"required"`
	FirstInsert time.Time `json:"first_insert" binding:"required"`
	InsertBy    string    `json:"insert_by" binding:"required"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   string    `json:"created_by"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedBy   string    `json:"updated_by"`
	Id          uint64    `json:"id"`
}

type Sal_activity_log struct {
	SalesActivityId uint64    `json:"sales_activity_id" binding:"required"`
	Date            time.Time `json:"date" binding:"required"`
	Activity        string    `json:"activity" binding:"required"`
	CreatedAt       time.Time `json:"created_at"`
	CreatedBy       string    `json:"created_by"`
	UpdatedAt       time.Time `json:"updated_at"`
	UpdatedBy       string    `json:"updated_by"`
	Id              uint64    `json:"id"`
}

type Sal_cfg_init_segment struct {
	Code      string    `json:"code" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Sal_cfg_init_source struct {
	Code      string    `json:"code" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Sal_cfg_init_task_action struct {
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Sal_cfg_init_task_repeat struct {
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Days      int       `json:"days"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Sal_cfg_init_task_tag struct {
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Id        uint64    `json:"id"`
}

type Sal_cfg_init_template struct {
	Id        int       `json:"id"`
	Template  string    `json:"template" binding:"required"`
	IsEmail   string    `json:"is_email" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	IdLog     uint64    `json:"id_log" binding:"required"`
}

type Sal_const_proposal_status struct {
	Code string `json:"code" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type Sal_const_resource struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Sal_const_status struct {
	Code string `json:"code" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type Sal_const_task_priority struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type Sal_const_task_status struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type Sal_contact struct {
	SalesActivityId uint64    `json:"sales_activity_id" binding:"required"`
	PhoneBookId     uint64    `json:"phone_book_id" binding:"required"`
	CreatedAt       time.Time `json:"created_at"`
	CreatedBy       string    `json:"created_by"`
	UpdatedAt       time.Time `json:"updated_at"`
	UpdatedBy       string    `json:"updated_by"`
	Id              uint64    `json:"id"`
}

type Sal_notes struct {
	SalesActivityId uint64    `json:"sales_activity_id" binding:"required"`
	Date            time.Time `json:"date" binding:"required"`
	Notes           string    `json:"notes" binding:"required"`
	CreatedAt       time.Time `json:"created_at"`
	CreatedBy       string    `json:"created_by"`
	UpdatedAt       time.Time `json:"updated_at"`
	UpdatedBy       string    `json:"updated_by"`
	Id              uint64    `json:"id"`
}

type Sal_proposal struct {
	Number          string    `json:"number" binding:"required"`
	SalesActivityId uint64    `json:"sales_activity_id" binding:"required"`
	Date            time.Time `json:"date" binding:"required"`
	ValidDate       time.Time `json:"valid_date" binding:"required"`
	CurrencyCode    string    `json:"currency_code" binding:"required"`
	Value           float64   `json:"value"`
	Description     string    `json:"description"`
	StatusCode      string    `json:"status_code"`
	DealValue       float64   `json:"deal_value"`
	ActualDealValue float64   `json:"actual_deal_value" binding:"required"`
	Notes           string    `json:"notes"`
	FirstInsert     time.Time `json:"first_insert" binding:"required"`
	InsertBy        string    `json:"insert_by" binding:"required"`
	CreatedAt       time.Time `json:"created_at"`
	CreatedBy       string    `json:"created_by"`
	UpdatedAt       time.Time `json:"updated_at"`
	UpdatedBy       string    `json:"updated_by"`
	Id              uint64    `json:"id"`
}

type Sal_send_reminder struct {
	SalesActivityId uint64    `json:"sales_activity_id" binding:"required"`
	Date            time.Time `json:"date" binding:"required"`
	IsSentMe        string    `json:"is_sent_me" binding:"required"`
	IsSentEmail     string    `json:"is_sent_email"`
	IsSentWa        string    `json:"is_sent_wa" binding:"required"`
	CreatedAt       time.Time `json:"created_at"`
	CreatedBy       string    `json:"created_by"`
	UpdatedAt       time.Time `json:"updated_at"`
	UpdatedBy       string    `json:"updated_by"`
	Id              uint64    `json:"id"`
}

type Sal_task struct {
	SalesActivityId uint64    `json:"sales_activity_id"`
	Subject         string    `json:"subject"`
	StartDate       time.Time `json:"start_date"`
	DueDate         time.Time `json:"due_date"`
	PriorityCode    string    `json:"priority_code"`
	RepeatCode      string    `json:"repeat_code"`
	TaskDescription string    `json:"task_description"`
	TagCode         string    `json:"tag_code"`
	Description     string    `json:"description"`
	IsVisit         string    `json:"is_visit"`
	VisitTo         string    `json:"visit_to"`
	ActionCode      string    `json:"action_code"`
	StatusCode      string    `json:"status_code"`
	FirstInsert     time.Time `json:"first_insert"`
	InsertBy        string    `json:"insert_by"`
	CreatedAt       time.Time `json:"created_at"`
	CreatedBy       string    `json:"created_by"`
	UpdatedAt       time.Time `json:"updated_at"`
	UpdatedBy       string    `json:"updated_by"`
	Id              uint64    `json:"id"`
}

type Sms_event struct {
	Code         int    `json:"code" binding:"required"`
	Name         string `json:"name" binding:"required"`
	Query        string `json:"query" binding:"required"`
	DefaultField string `json:"default_field" binding:"required"`
}

type Sms_outbox struct {
	SmsEventCode int       `json:"sms_event_code" binding:"required"`
	Number       string    `json:"number" binding:"required"`
	Smsc         string    `json:"smsc"`
	Text         string    `json:"text" binding:"required"`
	Date         time.Time `json:"date" binding:"required"`
	Sent         time.Time `json:"sent" binding:"required"`
	IsSent       uint8     `json:"is_sent" binding:"required"`
	IsFailed     uint8     `json:"is_failed" binding:"required"`
	IsDelivered  uint8     `json:"is_delivered" binding:"required"`
	CmgsId       string    `json:"cmgs_id" binding:"required"`
	CreatedAt    time.Time `json:"created_at"`
	CreatedBy    string    `json:"created_by"`
	UpdatedAt    time.Time `json:"updated_at"`
	UpdatedBy    string    `json:"updated_by"`
	Id           uint64    `json:"id"`
}

type Sms_schedule struct {
	SmsEventCode        int       `json:"sms_event_code" binding:"required"`
	DestinationTypeCode string    `json:"destination_type_code" binding:"required"`
	Number              string    `json:"number" binding:"required"`
	SmsFormat           string    `json:"sms_format" binding:"required"`
	SendDate            time.Time `json:"send_date" binding:"required"`
	RepeatTypeCode      string    `json:"repeat_type_code" binding:"required"`
	IsActive            uint8     `json:"is_active" binding:"required"`
	CreatedAt           time.Time `json:"created_at"`
	CreatedBy           string    `json:"created_by"`
	UpdatedAt           time.Time `json:"updated_at"`
	UpdatedBy           string    `json:"updated_by"`
	Id                  uint64    `json:"id"`
}

type Sub_folio struct {
	FolioNumber         uint64    `json:"folio_number" binding:"required"`
	BelongsTo           uint64    `json:"belongs_to" binding:"required"`
	GroupCode           string    `json:"group_code" binding:"required"`
	RoomNumber          string    `json:"room_number" binding:"required"`
	SubDepartmentCode   string    `json:"sub_department_code" binding:"required"`
	AccountCode         string    `json:"account_code" binding:"required"`
	ProductCode         string    `json:"product_code" binding:"required"`
	PackageCode         string    `json:"package_code"`
	Quantity            float64   `json:"quantity" binding:"required"`
	Amount              float64   `json:"amount" binding:"required"`
	DefaultCurrencyCode string    `json:"default_currency_code" binding:"required"`
	AmountForeign       float64   `json:"amount_foreign" binding:"required"`
	ExchangeRate        float64   `json:"exchange_rate" binding:"required"`
	CurrencyCode        string    `json:"currency_code" binding:"required"`
	AuditDate           time.Time `json:"audit_date" binding:"required"`
	AuditDateUnixx      int64     `json:"audit_date_unixx" binding:"required"`
	Remark              string    `json:"remark"`
	DocumentNumber      string    `json:"document_number"`
	VoucherNumber       string    `json:"voucher_number"`
	TypeCode            string    `json:"type_code" binding:"required"`
	CardBankCode        string    `json:"card_bank_code"`
	CardTypeCode        string    `json:"card_type_code"`
	Void                uint8     `json:"void" binding:"required"`
	VoidDate            time.Time `json:"void_date"`
	VoidBy              string    `json:"void_by"`
	VoidReason          string    `json:"void_reason"`
	IsCorrection        uint8     `json:"is_correction" binding:"required"`
	CorrectionBy        string    `json:"correction_by"`
	CorrectionReason    string    `json:"correction_reason"`
	CorrectionBreakdown uint64    `json:"correction_breakdown" binding:"required"`
	Breakdown1          uint64    `json:"breakdown1" binding:"required"`
	Breakdown2          int       `json:"breakdown2" binding:"required"`
	DirectBillCode      string    `json:"direct_bill_code"`
	PostingType         string    `json:"posting_type" binding:"required"`
	ExtraChargeId       uint64    `json:"extra_charge_id" binding:"required"`
	RefNumber           string    `json:"ref_number"`
	InsertBy            string    `json:"insert_by" binding:"required"`
	FirstInsert         time.Time `json:"first_insert" binding:"required"`
	Shift               string    `json:"shift" binding:"required"`
	LogShiftId          uint64    `json:"log_shift_id" binding:"required"`
	IsPairWithDeposit   uint8     `json:"is_pair_with_deposit" binding:"required"`
	TransferPairId      uint64    `json:"transfer_pair_id"`
	SystemCode          string    `json:"system_code" binding:"required"`
	CreatedAt           time.Time `json:"created_at"`
	CreatedBy           string    `json:"created_by"`
	UpdatedAt           time.Time `json:"updated_at"`
	UpdatedBy           string    `json:"updated_by"`
	Id                  uint64    `json:"id" gorm:"primaryKey"`
}

type Sub_folio_group struct {
	Code string `json:"code" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type Temp_sub_folio_breakdown1 struct {
	Breakdown1 uint64 `json:"breakdown1" binding:"required"`
	IpAddress  string `json:"ip_address" binding:"required"`
}

type Temp_sub_folio_correction_breakdown struct {
	CorrectionBreakdown uint64 `json:"correction_breakdown" binding:"required"`
	IpAddress           string `json:"ip_address" binding:"required"`
}

type User struct {
	Code                string    `json:"code" binding:"required"`
	Name                string    `json:"name" binding:"required"`
	Password            string    `json:"password" binding:"required"`
	UserGroupAccessCode string    `json:"user_group_access_code" binding:"required"`
	IsActive            uint8     `json:"is_active" binding:"required"`
	CreatedAt           time.Time `json:"created_at"`
	CreatedBy           string    `json:"created_by"`
	UpdatedAt           time.Time `json:"updated_at"`
	UpdatedBy           string    `json:"updated_by"`
	Id                  uint64    `json:"id" gorm:"primaryKey"`
}
type User_group struct {
	Code                    string    `json:"code"`
	AccessForm              string    `json:"access_form" binding:"required"`
	AccessSpecial           string    `json:"access_special" binding:"required"`
	AccessKeylock           string    `json:"access_keylock" binding:"required"`
	AccessReservation       string    `json:"access_reservation" binding:"required"`
	AccessDeposit           string    `json:"access_deposit" binding:"required"`
	AccessInHouse           string    `json:"access_in_house" binding:"required"`
	AccessWalkIn            string    `json:"access_walk_in" binding:"required"`
	AccessFolio             string    `json:"access_folio" binding:"required"`
	AccessFolioHistory      string    `json:"access_folio_history" binding:"required"`
	AccessFloorPlan         string    `json:"access_floor_plan" binding:"required"`
	AccessMemberVoucherGift string    `json:"access_member_voucher_gift" binding:"required"`
	SaMaxDiscountPercent    int       `json:"sa_max_discount_percent" binding:"required"`
	SaMaxDiscountAmount     float64   `json:"sa_max_discount_amount" binding:"required"`
	IsActive                uint8     `json:"is_active"`
	CreatedAt               time.Time `json:"created_at"`
	CreatedBy               string    `json:"created_by"`
	UpdatedAt               time.Time `json:"updated_at"`
	UpdatedBy               string    `json:"updated_by"`
	Id                      uint64    `json:"id"`
}
type User_group_access struct {
	Id                  uint64    `json:"id"`
	Code                string    `json:"code" binding:"required"`
	GeneralUserGroupId  uint64    `json:"general_user_group_id"`
	UserGroupId         uint64    `json:"user_group_id"`
	PosUserGroupId      uint64    `json:"pos_user_group_id"`
	BanUserGroupId      uint64    `json:"ban_user_group_id"`
	AccUserGroupId      uint64    `json:"acc_user_group_id"`
	AstUserGroupId      uint64    `json:"ast_user_group_id"`
	PyrUserGroupId      uint64    `json:"pyr_user_group_id"`
	CorUserGroupId      uint64    `json:"cor_user_group_id"`
	ReportUserGroupId   uint64    `json:"report_user_group_id"`
	ToolsUserGroupId    uint64    `json:"tools_user_group_id"`
	UserAccessLevelCode int       `json:"user_access_level_code"`
	IsActive            uint8     `json:"is_active"`
	CreatedAt           time.Time `json:"created_at"`
	CreatedBy           string    `json:"created_by"`
	UpdatedAt           time.Time `json:"updated_at"`
	UpdatedBy           string    `json:"updated_by"`
}

type Tools_user_group struct {
	Code                string    `json:"code" binding:"required"`
	AccessForm          string    `json:"access_form" binding:"required"`
	AccessConfiguration string    `json:"access_configuration" binding:"required"`
	AccessCompany       string    `json:"access_company" binding:"required"`
	CreatedAt           time.Time `json:"created_at"`
	CreatedBy           string    `json:"created_by"`
	UpdatedAt           time.Time `json:"updated_at"`
	UpdatedBy           string    `json:"updated_by"`
	Id                  uint64    `json:"id"`
}

type Report_user_group struct {
	Code                string    `json:"code" binding:"required"`
	AccessForm          string    `json:"access_form"`
	AccessFoReport      string    `json:"access_fo_report"`
	AccessPosReport     string    `json:"access_pos_report"`
	AccessBanReport     string    `json:"access_ban_report"`
	AccessAccReport     string    `json:"access_acc_report"`
	AccessAstReport     string    `json:"access_ast_report"`
	AccessPyrReport     string    `json:"access_pyr_report"`
	AccessCorReport     string    `json:"access_cor_report"`
	AccessPreviewReport string    `json:"access_preview_report"`
	CreatedAt           time.Time `json:"created_at"`
	CreatedBy           string    `json:"created_by"`
	UpdatedAt           time.Time `json:"updated_at"`
	UpdatedBy           string    `json:"updated_by"`
	Id                  uint64    `json:"id"`
}
type General_user_group struct {
	Code         string    `json:"code"`
	AccessModule string    `json:"access_module" binding:"required"`
	IsActive     uint8     `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	CreatedBy    string    `json:"created_by"`
	UpdatedAt    time.Time `json:"updated_at"`
	UpdatedBy    string    `json:"updated_by"`
	Id           uint64    `json:"id"`
}

type Voucher struct {
	Number            string    `json:"number" binding:"required"`
	Code              string    `json:"code" binding:"required"`
	RefNumber         string    `json:"ref_number" binding:"required"`
	AccountCode       string    `json:"account_code"`
	IssuedDate        time.Time `json:"issued_date" binding:"required"`
	UsedDate          time.Time `json:"used_date" binding:"required"`
	TitleCode         string    `json:"title_code"`
	FullName          string    `json:"full_name"`
	Street            string    `json:"street"`
	CompanyCode       string    `json:"company_code"`
	MemberTypeCode    string    `json:"member_type_code" binding:"required"`
	TypeCode          string    `json:"type_code" binding:"required"`
	StatusCode        string    `json:"status_code" binding:"required"`
	StatusCodeApprove string    `json:"status_code_approve" binding:"required"`
	StatusCodeSold    string    `json:"status_code_sold" binding:"required"`
	ReasonCode        string    `json:"reason_code"`
	RequestByCode     string    `json:"request_by_code"`
	AccommodationType string    `json:"accommodation_type"`
	SetUpRequired     string    `json:"set_up_required"`
	ApproveBy         string    `json:"approve_by" binding:"required"`
	StartDate         time.Time `json:"start_date" binding:"required"`
	ExpireDate        time.Time `json:"expire_date" binding:"required"`
	Description       string    `json:"description" binding:"required"`
	Point             float64   `json:"point" binding:"required"`
	Price             float64   `json:"price" binding:"required"`
	IsRoomChargeOnly  uint8     `json:"is_room_charge_only" binding:"required"`
	IsPerDay          uint8     `json:"is_per_day" binding:"required"`
	Nights            int       `json:"nights" binding:"required"`
	RoomTypeCode      string    `json:"room_type_code" binding:"required"`
	FolioNumber       uint64    `json:"folio_number"`
	SubFolioId        uint64    `json:"sub_folio_id"`
	CreatedAt         time.Time `json:"created_at"`
	CreatedBy         string    `json:"created_by"`
	UpdatedAt         time.Time `json:"updated_at"`
	UpdatedBy         string    `json:"updated_by"`
	Id                uint64    `json:"id"`
}

type Working_shift struct {
	Shift     string    `json:"shift" binding:"required"`
	StartTime time.Time `json:"start_time" binding:"required"`
	EndTime   time.Time `json:"end_time" binding:"required"`
}

// =================== CUSTOM STRUCT ========================================================================================================================//
// type GeneralCodeStruct struct {
// 	Code string `json:"code"`
// }

type GeneralCodeNameStruct struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type GeneralIDNameStruct struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type GeneralCodeDescriptionStruct struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

type ReservationComboListStruct struct {
	RoomType, Currency, BusinessSource, CommissionType, Market, BookingSource, PaymentType, Member, Title, Country, Nationality, Company, GuestType, CardType,
	PurposeOf, GuestGroup, Sales []GeneralCodeNameStruct
}

type GuestDepositPostStruct struct {
	GuestDeposit Guest_deposit `json:"guest_deposit" binding:"required,dive" `
	IDCorrected  uint64        `json:"id_corrected"`
	IsCard       bool          `json:"is_card"`
	SystemCode   string        `json:"system_code"`
	ChargeAmount float64       `json:"charge_amount"`
	CardNumber   string        `json:"card_number"`
	CardHolder   string        `json:"card_holder"`
	ValidMonth   string        `json:"valid_month"`
	ValidYear    string        `json:"valid_year"`
}

type ReservationPostStruct struct {
	GuestProfileData1    Guest_profile
	GuestProfileData2    Guest_profile
	GuestProfileData3    Guest_profile
	GuestProfileData4    Guest_profile
	GuestDetailData      Guest_detail
	GuestGeneralData     Guest_general
	ReservationData      Reservation
	IsWaitList           uint8
	IsIncognito          uint8
	PostFirstNightCharge uint8
}

type GuestInHousePostStruct struct {
	GuestProfileData1    Guest_profile
	GuestProfileData2    Guest_profile
	GuestProfileData3    Guest_profile
	GuestProfileData4    Guest_profile
	GuestDetailData      Guest_detail
	GuestGeneralData     Guest_general
	FolioData            Folio
	IsIncognito          uint8
	PostFirstNightCharge uint8
}

type RegistrationFormStruct struct {
	GuestProfileData1 Guest_profile
	GuestProfileData2 Guest_profile
	GuestProfileData3 Guest_profile
	GuestProfileData4 Guest_profile
	GuestDetailData   Guest_detail
	GuestGeneralData  Guest_general
}

type MasterDeskFolioPostStruct struct {
	ContactPerson    Contact_person
	GuestDetailData  Guest_detail
	GuestGeneralData Guest_general
	FolioData        Folio
}
type ReservationDepositComboListStruct struct {
	Currency, SubFolioGroup, SubDepartment []GeneralCodeNameStruct
}

type ReservationDepositCardComboListStruct struct {
	Currency, SubFolioGroup, SubDepartment, CardBank, CardType []GeneralCodeNameStruct
}

type GeneralNumberStruct struct {
	Number string
}

type CaptainOrderTransactionStruct struct {
	CaptainOrderId       uint64
	PriceOriginalForeign float64
	SubTotal             float64
	Total                float64
	AccountCode          string
	Buy                  int
	CompanyCode1         string
	CompanyCode2         string
	CompanyName          string
	CurrencyCode         string
	DefaultCurrencyCode  string
	Description          string
	DisableDiscount      string
	Discount             float64
	ExchangeRate         float64
	Free                 int
	Id                   uint64
	IsFree               string
	PackageCode          string
	Price                float64
	PriceOriginal        float64
	ProductCode          string
	ProductName          string
	Quantity             float64
	Remark               string
	Service              float64
	SpaEndDate           time.Time
	SpaRoomNumber        string
	SpaStartDate         time.Time
	Tax                  float64
	TaxAndServiceCode    string
	TypeCode             string
	CategoryCode         string
}

// type Get_room_list_struct struct {
// 	Number             string
// 	RoomTypeCode       string
// 	StatusRoom         string
// 	StatusRoomOccupied string
// 	FullName           string
// }

var TableName = TTableName{
	AccApAr:                              "acc_ap_ar",
	AccApArPayment:                       "acc_ap_ar_payment",
	AccApArPaymentDetail:                 "acc_ap_ar_payment_detail",
	AccApCommissionPayment:               "acc_ap_commission_payment",
	AccApCommissionPaymentDetail:         "acc_ap_commission_payment_detail",
	AccApRefundDepositPayment:            "acc_ap_refund_deposit_payment",
	AccApRefundDepositPaymentDetail:      "acc_ap_refund_deposit_payment_detail",
	AccCashSaleRecon:                     "acc_cash_sale_recon",
	AccCfgInitBankAccount:                "acc_cfg_init_bank_account",
	AccCloseMonth:                        "acc_close_month",
	AccCloseYear:                         "acc_close_year",
	AccConstBankAccountType:              "acc_const_bank_account_type",
	AccConstJournalGroup:                 "acc_const_journal_group",
	AccConstJournalType:                  "acc_const_journal_type",
	AccConstUnit:                         "acc_const_unit",
	AccCreditCardRecon:                   "acc_credit_card_recon",
	AccCreditCardReconDetail:             "acc_credit_card_recon_detail",
	AccDefferedIncome:                    "acc_deffered_income",
	AccDefferedIncomePosted:              "acc_deffered_income_posted",
	AccForeignCash:                       "acc_foreign_cash",
	AccImportJournalLog:                  "acc_import_journal_log",
	AccJournal:                           "acc_journal",
	AccJournalDetail:                     "acc_journal_detail",
	AccPrepaidExpense:                    "acc_prepaid_expense",
	AccPrepaidExpensePosted:              "acc_prepaid_expense_posted",
	AccReport:                            "acc_report",
	AccReportDefaultField:                "acc_report_default_field",
	AccReportGroupField:                  "acc_report_group_field",
	AccReportGroupingField:               "acc_report_grouping_field",
	AccReportOrderField:                  "acc_report_order_field",
	AccReportTemplate:                    "acc_report_template",
	AccReportTemplateField:               "acc_report_template_field",
	AccUserGroup:                         "acc_user_group",
	AstCfgInitShippingAddress:            "ast_cfg_init_shipping_address",
	AstConstPurchaseRequestStatus:        "ast_const_purchase_request_status",
	AstConstStoreRequisitionStatus:       "ast_const_store_requisition_status",
	AstReport:                            "ast_report",
	AstReportDefaultField:                "ast_report_default_field",
	AstReportGroupField:                  "ast_report_group_field",
	AstReportGroupingField:               "ast_report_grouping_field",
	AstReportOrderField:                  "ast_report_order_field",
	AstReportTemplate:                    "ast_report_template",
	AstReportTemplateField:               "ast_report_template_field",
	AstUserGroup:                         "ast_user_group",
	AstUserSubDepartment:                 "ast_user_sub_department",
	AuditLog:                             "audit_log",
	BanBooking:                           "ban_booking",
	BanBookingSchedulePayment:            "ban_booking_schedule_payment",
	BanCfgInitSeatingPlan:                "ban_cfg_init_seating_plan",
	BanCfgInitTheme:                      "ban_cfg_init_theme",
	BanCfgInitVenue:                      "ban_cfg_init_venue",
	BanCfgInitVenueCombine:               "ban_cfg_init_venue_combine",
	BanCfgInitVenueCombineDetail:         "ban_cfg_init_venue_combine_detail",
	BanCfgInitVenueGroup:                 "ban_cfg_init_venue_group",
	BanConstBookingStatus:                "ban_const_booking_status",
	BanConstReservationStatus:            "ban_const_reservation_status",
	BanConstReservationType:              "ban_const_reservation_type",
	BanConstVenueLocation:                "ban_const_venue_location",
	BanReport:                            "ban_report",
	BanReportDefaultField:                "ban_report_default_field",
	BanReportGroupField:                  "ban_report_group_field",
	BanReportGroupingField:               "ban_report_grouping_field",
	BanReportOrderField:                  "ban_report_order_field",
	BanReportTemplate:                    "ban_report_template",
	BanReportTemplateField:               "ban_report_template_field",
	BanReservation:                       "ban_reservation",
	BanReservationCharge:                 "ban_reservation_charge",
	BanReservationRemark:                 "ban_reservation_remark",
	BanUserGroup:                         "ban_user_group",
	BanCfgInitLayout:                     "ban_cfg_init_layout",
	BreakfastListTemp:                    "breakfast_list_temp",
	BudgetExpense:                        "budget_expense",
	BudgetFb:                             "budget_fb",
	BudgetIncome:                         "budget_income",
	BudgetStatistic:                      "budget_statistic",
	CashCount:                            "cash_count",
	CfgInitAccount:                       "cfg_init_account",
	CfgInitAccountSubGroup:               "cfg_init_account_sub_group",
	CfgInitBedType:                       "cfg_init_bed_type",
	CfgInitBookingSource:                 "cfg_init_booking_source",
	CfgInitCardBank:                      "cfg_init_card_bank",
	CfgInitCardType:                      "cfg_init_card_type",
	CfgInitCity:                          "cfg_init_city",
	CfgInitCompanyType:                   "cfg_init_company_type",
	CfgInitCompetitorCategory:            "cfg_init_competitor_category",
	CfgInitContinent:                     "cfg_init_continent",
	CfgInitCountry:                       "cfg_init_country",
	CfgInitCreditCardCharge:              "cfg_init_credit_card_charge",
	CfgInitCurrency:                      "cfg_init_currency",
	CfgInitCurrencyNominal:               "cfg_init_currency_nominal",
	CfgInitCustomLookupField01:           "cfg_init_custom_lookup_field01",
	CfgInitCustomLookupField02:           "cfg_init_custom_lookup_field02",
	CfgInitCustomLookupField03:           "cfg_init_custom_lookup_field03",
	CfgInitCustomLookupField04:           "cfg_init_custom_lookup_field04",
	CfgInitCustomLookupField05:           "cfg_init_custom_lookup_field05",
	CfgInitCustomLookupField06:           "cfg_init_custom_lookup_field06",
	CfgInitCustomLookupField07:           "cfg_init_custom_lookup_field07",
	CfgInitCustomLookupField08:           "cfg_init_custom_lookup_field08",
	CfgInitCustomLookupField09:           "cfg_init_custom_lookup_field09",
	CfgInitCustomLookupField10:           "cfg_init_custom_lookup_field10",
	CfgInitCustomLookupField11:           "cfg_init_custom_lookup_field11",
	CfgInitCustomLookupField12:           "cfg_init_custom_lookup_field12",
	CfgInitDepartment:                    "cfg_init_department",
	CfgInitGuestType:                     "cfg_init_guest_type",
	CfgInitIdCardType:                    "cfg_init_id_card_type",
	CfgInitIsFbSubDepartmentGroup:        "cfg_init_is_fb_sub_department_group",
	CfgInitIsFbSubDepartmentGroupDetail:  "cfg_init_is_fb_sub_department_group_detail",
	CfgInitJournalAccount:                "cfg_init_journal_account",
	CfgInitJournalAccountCategory:        "cfg_init_journal_account_category",
	CfgInitJournalAccountSubGroup:        "cfg_init_journal_account_sub_group",
	CfgInitLanguage:                      "cfg_init_language",
	CfgInitLoanItem:                      "cfg_init_loan_item",
	CfgInitMarket:                        "cfg_init_market",
	CfgInitMarketCategory:                "cfg_init_market_category",
	CfgInitMemberPointType:               "cfg_init_member_point_type",
	CfgInitNationality:                   "cfg_init_nationality",
	CfgInitOwner:                         "cfg_init_owner",
	CfgInitPabxRate:                      "cfg_init_pabx_rate",
	CfgInitPackage:                       "cfg_init_package",
	CfgInitPackageBreakdown:              "cfg_init_package_breakdown",
	CfgInitPackageBusinessSource:         "cfg_init_package_business_source",
	CfgInitPaymentType:                   "cfg_init_payment_type",
	CfgInitPhoneBookType:                 "cfg_init_phone_book_type",
	CfgInitPrinter:                       "cfg_init_printer",
	CfgInitPurposeOf:                     "cfg_init_purpose_of",
	CfgInitRegency:                       "cfg_init_regency",
	CfgInitReservationMark:               "cfg_init_reservation_mark",
	CfgInitRoom:                          "cfg_init_room",
	CfgInitRoomAllotmentType:             "cfg_init_room_allotment_type",
	CfgInitRoomAmenities:                 "cfg_init_room_amenities",
	CfgInitRoomBoy:                       "cfg_init_room_boy",
	CfgInitRoomRate:                      "cfg_init_room_rate",
	CfgInitRoomRateBreakdown:             "cfg_init_room_rate_breakdown",
	CfgInitRoomRateBusinessSource:        "cfg_init_room_rate_business_source",
	CfgInitRoomRateCategory:              "cfg_init_room_rate_category",
	CfgInitRoomRateCompetitor:            "cfg_init_room_rate_competitor",
	CfgInitRoomRateCurrency:              "cfg_init_room_rate_currency",
	CfgInitRoomRateDynamic:               "cfg_init_room_rate_dynamic",
	CfgInitRoomRateLastDeal:              "cfg_init_room_rate_last_deal",
	CfgInitRoomRateScale:                 "cfg_init_room_rate_scale",
	CfgInitRoomRateSession:               "cfg_init_room_rate_session",
	CfgInitRoomRateSubCategory:           "cfg_init_room_rate_sub_category",
	CfgInitRoomRateWeekly:                "cfg_init_room_rate_weekly",
	CfgInitRoomType:                      "cfg_init_room_type",
	CfgInitRoomUnavailableReason:         "cfg_init_room_unavailable_reason",
	CfgInitRoomView:                      "cfg_init_room_view",
	CfgInitSales:                         "cfg_init_sales",
	CfgInitSalesSalary:                   "cfg_init_sales_salary",
	CfgInitState:                         "cfg_init_state",
	CfgInitSubDepartment:                 "cfg_init_sub_department",
	CfgInitTaxAndService:                 "cfg_init_tax_and_service",
	CfgInitTitle:                         "cfg_init_title",
	CfgInitVoucherReason:                 "cfg_init_voucher_reason",
	CmNotification:                       "cm_notification",
	CmUpdate:                             "cm_update",
	Company:                              "company",
	Competitor:                           "competitor",
	CompetitorData:                       "competitor_data",
	Configuration:                        "configuration",
	ConstAccountGroup:                    "const_account_group",
	CmLog:                                "cm_log",
	ConstBudgetType:                      "const_budget_type",
	ConstChannelManagerVendor:            "const_channel_manager_vendor",
	ConstChargeFrequency:                 "const_charge_frequency",
	ConstChargeType:                      "const_charge_type",
	ConstCommissionType:                  "const_commission_type",
	ConstCustomerDisplayVendor:           "const_customer_display_vendor",
	ConstDepartmentType:                  "const_department_type",
	ConstDynamicRateType:                 "const_dynamic_rate_type",
	ConstFolioStatus:                     "const_folio_status",
	ConstFolioType:                       "const_folio_type",
	ConstForecastDay:                     "const_forecast_day",
	ConstForecastMonth:                   "const_forecast_month",
	ConstForeignCashTableId:              "const_foreign_cash_table_id",
	ConstGuestStatus:                     "const_guest_status",
	ConstImage:                           "const_image",
	ConstIptvVendor:                      "const_iptv_vendor",
	ConstJournalAccountGroup:             "const_journal_account_group",
	ConstJournalAccountSubGroupType:      "const_journal_account_sub_group_type",
	ConstJournalAccountType:              "const_journal_account_type",
	ConstJournalPrefix:                   "const_journal_prefix",
	ConstKeylockVendor:                   "const_keylock_vendor",
	ConstMemberType:                      "const_member_type",
	ConstMikrotikVendor:                  "const_mikrotik_vendor",
	ConstNotificationType:                "const_notification_type",
	ConstOtherIcon:                       "const_other_icon",
	ConstOtpStatus:                       "const_otp_status",
	ConstPabxRateType:                    "const_pabx_rate_type",
	ConstPaymentGroup:                    "const_payment_group",
	ConstReportFont:                      "const_report_font",
	ConstReportFormat:                    "const_report_format",
	ConstReservationStatus:               "const_reservation_status",
	ConstRoomBlockStatus:                 "const_room_block_status",
	ConstRoomStatus:                      "const_room_status",
	ConstSmsDestinationType:              "const_sms_destination_type",
	ConstSmsRepeatType:                   "const_sms_repeat_type",
	ConstStatisticAccount:                "const_statistic_account",
	ConstSystem:                          "const_system",
	ConstTransactionType:                 "const_transaction_type",
	ConstVoucherStatus:                   "const_voucher_status",
	ConstVoucherStatusApprove:            "const_voucher_status_approve",
	ConstVoucherStatusSold:               "const_voucher_status_sold",
	ConstVoucherType:                     "const_voucher_type",
	ContactPerson:                        "contact_person",
	CorCfgInitUnit:                       "cor_cfg_init_unit",
	CorReport:                            "cor_report",
	CorReportDefaultField:                "cor_report_default_field",
	CorReportGroupField:                  "cor_report_group_field",
	CorReportGroupingField:               "cor_report_grouping_field",
	CorReportOrderField:                  "cor_report_order_field",
	CorReportTemplate:                    "cor_report_template",
	CorReportTemplateField:               "cor_report_template_field",
	CorUserGroup:                         "cor_user_group",
	CreditCard:                           "credit_card",
	DataAnalysis:                         "data_analysis",
	DataAnalysisQuery:                    "data_analysis_query",
	DataAnalysisQueryList:                "data_analysis_query_list",
	Events:                               "events",
	FaCfgInitItem:                        "fa_cfg_init_item",
	FaCfgInitItemCategory:                "fa_cfg_init_item_category",
	FaCfgInitLocation:                    "fa_cfg_init_location",
	FaCfgInitManufacture:                 "fa_cfg_init_manufacture",
	FaConstDepreciationType:              "fa_const_depreciation_type",
	FaConstItemCondition:                 "fa_const_item_condition",
	FaConstLocationType:                  "fa_const_location_type",
	FaDepreciation:                       "fa_depreciation",
	FaList:                               "fa_list",
	FaLocationHistory:                    "fa_location_history",
	FaPurchaseOrder:                      "fa_purchase_order",
	FaPurchaseOrderDetail:                "fa_purchase_order_detail",
	FaReceive:                            "fa_receive",
	FaReceiveDetail:                      "fa_receive_detail",
	FaRepair:                             "fa_repair",
	FaRevaluation:                        "fa_revaluation",
	FbStatistic:                          "fb_statistic",
	Folio:                                "folio",
	FolioList:                            "folio_list",
	FolioRouting:                         "folio_routing",
	ForecastInHouseChangePax:             "forecast_in_house_change_pax",
	ForecastMonthlyDay:                   "forecast_monthly_day",
	ForecastMonthlyDayPrevious:           "forecast_monthly_day_previous",
	GeneralUserGroup:                     "general_user_group",
	ConstUserAccessLevel:                 "const_user_access_level",
	GridProperties:                       "grid_properties",
	GuestBreakdown:                       "guest_breakdown",
	GuestDeposit:                         "guest_deposit",
	GuestDepositTax:                      "guest_deposit_tax",
	GuestDetail:                          "guest_detail",
	GuestExtraCharge:                     "guest_extra_charge",
	GuestExtraChargeBreakdown:            "guest_extra_charge_breakdown",
	GuestGeneral:                         "guest_general",
	GuestGroup:                           "guest_group",
	GuestInHouse:                         "guest_in_house",
	GuestInHouseBreakdown:                "guest_in_house_breakdown",
	GuestLoanItem:                        "guest_loan_item",
	GuestMessage:                         "guest_message",
	GuestProfile:                         "guest_profile",
	GuestScheduledRate:                   "guest_scheduled_rate",
	GuestToDo:                            "guest_to_do",
	HotelInformation:                     "hotel_information",
	InvCfgInitItem:                       "inv_cfg_init_item",
	InvCfgInitItemCategory:               "inv_cfg_init_item_category",
	InvCfgInitItemCategoryOtherCogs:      "inv_cfg_init_item_category_other_cogs",
	InvCfgInitItemCategoryOtherCogs2:     "inv_cfg_init_item_category_other_cogs2",
	InvCfgInitItemCategoryOtherExpense:   "inv_cfg_init_item_category_other_expense",
	InvCfgInitItemGroup:                  "inv_cfg_init_item_group",
	InvCfgInitItemUom:                    "inv_cfg_init_item_uom",
	InvCfgInitMarketList:                 "inv_cfg_init_market_list",
	InvCfgInitReturnStockReason:          "inv_cfg_init_return_stock_reason",
	InvCfgInitStore:                      "inv_cfg_init_store",
	InvCfgInitUom:                        "inv_cfg_init_uom",
	InvCloseLog:                          "inv_close_log",
	InvConstItemGroupType:                "inv_const_item_group_type",
	InvCloseSummary:                      "inv_close_summary",
	InvCloseSummaryStore:                 "inv_close_summary_store",
	InvCostRecipe:                        "inv_cost_recipe",
	InvCosting:                           "inv_costing",
	InvCostingDetail:                     "inv_costing_detail",
	InvOpname:                            "inv_opname",
	InvProduction:                        "inv_production",
	InvPurchaseOrder:                     "inv_purchase_order",
	InvPurchaseOrderDetail:               "inv_purchase_order_detail",
	InvPurchaseRequest:                   "inv_purchase_request",
	InvPurchaseRequestDetail:             "inv_purchase_request_detail",
	InvReceiving:                         "inv_receiving",
	InvReceivingDetail:                   "inv_receiving_detail",
	InvReturnStock:                       "inv_return_stock",
	InvStockTransfer:                     "inv_stock_transfer",
	InvStockTransferDetail:               "inv_stock_transfer_detail",
	InvStoreRequisition:                  "inv_store_requisition",
	InvStoreRequisitionDetail:            "inv_store_requisition_detail",
	Invoice:                              "invoice",
	InvoiceItem:                          "invoice_item",
	InvoicePayment:                       "invoice_payment",
	Log:                                  "log",
	LogBackup:                            "log_backup",
	LogKeylock:                           "log_keylock",
	LogMode:                              "log_mode",
	LogShift:                             "log_shift",
	LogSpecialAccess:                     "log_special_access",
	LogUser:                              "log_user",
	LogUserAction:                        "log_user_action",
	LogUserActionGroup:                   "log_user_action_group",
	LostAndFound:                         "lost_and_found",
	MarketStatistic:                      "market_statistic",
	Member:                               "member",
	MemberGift:                           "member_gift",
	MemberPoint:                          "member_point",
	MemberPointRedeem:                    "member_point_redeem",
	NotifTp:                              "notif_tp",
	NotifTpCfgInitTemplate:               "notif_tp_cfg_init_template",
	NotifTpConstEvent:                    "notif_tp_const_event",
	NotifTpConstVariable:                 "notif_tp_const_variable",
	NotifTpConstVendor:                   "notif_tp_const_vendor",
	Notification:                         "notification",
	OneTimePassword:                      "one_time_password",
	PabxSmdr:                             "pabx_smdr",
	PhoneBook:                            "phone_book",
	PosCaptainOrder:                      "pos_captain_order",
	PosCaptainOrderTransaction:           "pos_captain_order_transaction",
	PosCfgInitDiscountLimit:              "pos_cfg_init_discount_limit",
	PosCfgInitMarket:                     "pos_cfg_init_market",
	PosCfgInitMemberOutletDiscount:       "pos_cfg_init_member_outlet_discount",
	PosCfgInitMemberOutletDiscountDetail: "pos_cfg_init_member_outlet_discount_detail",
	PosCfgInitMemberProductDiscount:      "pos_cfg_init_member_product_discount",
	PosCfgInitOutlet:                     "pos_cfg_init_outlet",
	PosCfgInitPaymentGroup:               "pos_cfg_init_payment_group",
	PosCfgInitProduct:                    "pos_cfg_init_product",
	PosCfgInitProductCategory:            "pos_cfg_init_product_category",
	PosCfgInitProductGroup:               "pos_cfg_init_product_group",
	PosCfgInitRoomBoy:                    "pos_cfg_init_room_boy",
	PosCfgInitSpaRoom:                    "pos_cfg_init_spa_room",
	PosCfgInitTable:                      "pos_cfg_init_table",
	PosCfgInitTableType:                  "pos_cfg_init_table_type",
	PosCfgInitTenan:                      "pos_cfg_init_tenan",
	PosCfgInitTherapistFingerprint:       "pos_cfg_init_therapist_fingerprint",
	PosCfgInitWaitress:                   "pos_cfg_init_waitress",
	PosCheck:                             "pos_check",
	PosCheckTransaction:                  "pos_check_transaction",
	PosConstCheckType:                    "pos_const_check_type",
	PosConstComplimentType:               "pos_const_compliment_type",
	PosConstDiscount:                     "pos_const_discount",
	PosConstTimeSegment:                  "pos_const_time_segment",
	PosInformation:                       "pos_information",
	PosIptvMenuOrder:                     "pos_iptv_menu_order",
	PosMember:                            "pos_member",
	PosProductCosting:                    "pos_product_costing",
	PosReport:                            "pos_report",
	PosReportDefaultField:                "pos_report_default_field",
	PosReportGroupField:                  "pos_report_group_field",
	PosReportGroupingField:               "pos_report_grouping_field",
	PosReportOrderField:                  "pos_report_order_field",
	PosReportTemplate:                    "pos_report_template",
	PosReportTemplateField:               "pos_report_template_field",
	PosReservation:                       "pos_reservation",
	PosReservationTable:                  "pos_reservation_table",
	PosTableUnavailable:                  "pos_table_unavailable",
	PosUserGroup:                         "pos_user_group",
	PosUserGroupOutlet:                   "pos_user_group_outlet",
	ProformaInvoiceDetail:                "proforma_invoice_detail",
	Receipt:                              "receipt",
	Report:                               "report",
	ReportCustom:                         "report_custom",
	ReportCustomFavorite:                 "report_custom_favorite",
	ReportDefaultField:                   "report_default_field",
	ReportGroupField:                     "report_group_field",
	ReportGroupingField:                  "report_grouping_field",
	ReportOrderField:                     "report_order_field",
	ReportPivotTemp:                      "report_pivot_temp",
	ReportRoomRateStructureTemp:          "report_room_rate_structure_temp",
	ReportRoomSales:                      "report_room_sales",
	ReportTemplate:                       "report_template",
	ReportTemplateField:                  "report_template_field",
	Reservation:                          "reservation",
	ReservationExtraCharge:               "reservation_extra_charge",
	ReservationExtraChargeBreakdown:      "reservation_extra_charge_breakdown",
	ReservationScheduledRate:             "reservation_scheduled_rate",
	RoomAllotment:                        "room_allotment",
	RoomStatistic:                        "room_statistic",
	RoomStatus:                           "room_status",
	RoomUnavailable:                      "room_unavailable",
	RoomUnavailableHistory:               "room_unavailable_history",
	SalActivity:                          "sal_activity",
	SalActivityLog:                       "sal_activity_log",
	SalCfgInitSegment:                    "sal_cfg_init_segment",
	SalCfgInitSource:                     "sal_cfg_init_source",
	SalCfgInitTaskAction:                 "sal_cfg_init_task_action",
	SalCfgInitTaskRepeat:                 "sal_cfg_init_task_repeat",
	SalCfgInitTaskTag:                    "sal_cfg_init_task_tag",
	SalCfgInitTemplate:                   "sal_cfg_init_template",
	SalConstProposalStatus:               "sal_const_proposal_status",
	SalConstResource:                     "sal_const_resource",
	SalConstStatus:                       "sal_const_status",
	SalConstTaskPriority:                 "sal_const_task_priority",
	SalConstTaskStatus:                   "sal_const_task_status",
	SalContact:                           "sal_contact",
	SalNotes:                             "sal_notes",
	SalProposal:                          "sal_proposal",
	SalSendReminder:                      "sal_send_reminder",
	SalTask:                              "sal_task",
	SmsEvent:                             "sms_event",
	SmsOutbox:                            "sms_outbox",
	SmsSchedule:                          "sms_schedule",
	SubFolio:                             "sub_folio",
	SubFolioGroup:                        "sub_folio_group",
	SubFolioGrouping:                     "sub_folio_grouping",
	SubFolioTax:                          "sub_folio_tax",
	TempSubFolioBreakdown1:               "temp_sub_folio_breakdown1",
	TempSubFolioCorrectionBreakdown:      "temp_sub_folio_correction_breakdown",
	TransactionList:                      "transaction_list",
	TransactionListTax:                   "transaction_list_tax",
	User:                                 "user",
	UserGroupAccess:                      "user_group_access",
	UserGroup:                            "user_group",
	Voucher:                              "voucher",
	WorkingShift:                         "working_shift",
	CfgInitTimezone:                      "cfg_init_timezone",
	CmUpdateRate:                         "cm_update_rate",
	CmUpdateAvailability:                 "cm_update_availability",
}

type Cm_update_availability struct {
	Id           uint64    `json:"id"`
	StartDate    time.Time `json:"start_date" binding:"required"`
	EndDate      time.Time `json:"end_date" binding:"required"`
	RoomTypeCode string    `json:"room_type_code" binding:"required"`
	BedTypeCode  string    `json:"bed_type_code"`
	Availability int       `json:"availability"`
	Status       string    `json:"status" binding:"required"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Cm_update_allotment struct {
	Id           uint64    `json:"id"`
	TypeCode     string    `json:"type_code" binding:"required"`
	Number       uint64    `json:"number" binding:"required"`
	RoomTypeCode string    `json:"room_type_code" binding:"required"`
	BedTypeCode  string    `json:"bed_type_code" binding:"required"`
	RoomRateCode string    `json:"room_rate_code"`
	RateAmount   float64   `json:"rate_amount"`
	StartDate    time.Time `json:"start_date" binding:"required"`
	EndDate      time.Time `json:"end_date" binding:"required"`
	PostingDate  time.Time `json:"posting_date" binding:"required"`
	IsUpdated    uint8     `json:"is_updated" binding:"required"`
}

type Ban_booking_schedule_payment struct {
	Id             uint64    `json:"id"`
	BookingNumber  uint64    `json:"booking_number" binding:"required"`
	GuestDepositId uint64    `json:"guest_deposit_id" binding:"required"`
	Name           string    `json:"name" binding:"required"`
	Date           time.Time `json:"date" binding:"required"`
	Amount         float64   `json:"amount" binding:"required"`
	Remark         string    `json:"remark"`
	IsPaid         uint8     `json:"is_paid" binding:"required"`
	PaymentRemark  string    `json:"payment_remark" binding:"required"`
	CreatedAt      time.Time `json:"created_at"`
	CreatedBy      string    `json:"created_by"`
	UpdatedAt      time.Time `json:"updated_at"`
	UpdatedBy      string    `json:"updated_by"`
}

type Cm_update_inhouse struct {
	Id           uint64    `json:"id"`
	TypeCode     string    `json:"type_code" binding:"required"`
	Number       uint64    `json:"number" binding:"required"`
	RoomTypeCode string    `json:"room_type_code" binding:"required"`
	BedTypeCode  string    `json:"bed_type_code" binding:"required"`
	RoomRateCode string    `json:"room_rate_code"`
	RateAmount   float64   `json:"rate_amount"`
	StartDate    time.Time `json:"start_date" binding:"required"`
	EndDate      time.Time `json:"end_date" binding:"required"`
	PostingDate  time.Time `json:"posting_date" binding:"required"`
	IsUpdated    uint8     `json:"is_updated" binding:"required"`
}

type Cm_update_rate struct {
	Id                uint64    `json:"id"`
	StartDate         time.Time `json:"start_date" binding:"required"`
	EndDate           time.Time `json:"end_date" binding:"required"`
	RoomRateCode      string    `json:"room_rate_code"`
	RateAmount        float64   `json:"rate_amount"`
	RoomTypeCode      string    `json:"room_type_code" binding:"required"`
	BedTypeCode       string    `json:"bed_type_code" binding:"required"`
	Day1              uint8     `json:"day1"`
	Day2              uint8     `json:"day2"`
	Day3              uint8     `json:"day3"`
	Day4              uint8     `json:"day4"`
	Day5              uint8     `json:"day5"`
	Day6              uint8     `json:"day6"`
	Day7              uint8     `json:"day7"`
	StopSell          uint8     `json:"stop_sell"`
	ClosedToArrival   uint8     `json:"closed_to_arrival"`
	ClosedToDeparture uint8     `json:"closed_to_departure"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type Cm_update_reservation struct {
	Id           uint64    `json:"id"`
	TypeCode     string    `json:"type_code" binding:"required"`
	Number       uint64    `json:"number" binding:"required"`
	RoomTypeCode string    `json:"room_type_code" binding:"required"`
	BedTypeCode  string    `json:"bed_type_code" binding:"required"`
	RoomRateCode string    `json:"room_rate_code"`
	RateAmount   float64   `json:"rate_amount"`
	StartDate    time.Time `json:"start_date" binding:"required"`
	EndDate      time.Time `json:"end_date" binding:"required"`
	PostingDate  time.Time `json:"posting_date" binding:"required"`
	IsUpdated    uint8     `json:"is_updated" binding:"required"`
}
