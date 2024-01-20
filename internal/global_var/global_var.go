package global_var

import (
	"chs/config"
	"sync"
	"time"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/otel"
	"gorm.io/gorm"
)

type databaseinfo struct {
	HostName, UserName, Password, DatabaseName string
	Port                                       int
}

type TProgramVariable struct {
	DatabaseVersion, ProgramPath, ProgramBackupPath, ProgramReportGeneralPath, ProgramReportPath, ProgramReportPathCPOS, ProgramDataAnalysisPath, ServerID, IPAddress, ComputerName, MACAddress, WallpaperPath, NotificationPosition, ProvenosToken, OutletCode string
	WallpaperIndex, WallpaperPosition, WallpaperColor, DebugCount, NotificationState, NotificationLastState, AutoLogOffCountDown                                                                                                                                int
	ShowSplashScreen, SQLDebugActive, ForceLogOff, AutoLogOffMessageShowed                                                                                                                                                                                      bool
	AuditDate                                                                                                                                                                                                                                                   time.Time
}

type TTimeSegment struct {
	Breakfast, Lunch, Dinner, CoffeBreak string
}

type TFNBStructureRate struct {
	Portion, Revenue, Average string
}

type TReportTemplateName struct {
	DailyRevenueReport1,
	DailyRevenueReport2,
	DailyRevenueReport3,
	DailyRevenueReport4,
	DailyRevenueReport5,
	DailyRevenueReport5Point2D,
	DailyRevenueReport6,
	DailyRevenueReport7,
	DailyRevenueReport8,
	DailyRevenueReport9,
	DailyRevenueReport10,
	DailyRevenueReport11 string
}

type TCustomDateOptions struct {
	Condition21,
	Condition22,
	ConditionAuditDate2 string
}

type TReportCodeName struct {
	Reservation,
	FrontDesk,
	MemberVoucherAndGift,
	Room,
	CakrasoftPointOfSales,
	Profile,
	MarketingGraphicAndAnalisis,
	SalesActivity,
	CityLedger,
	ApRefundDepositAndCommission,
	Log,
	DailyReport,
	MonthlyReport,
	YearlyReport,
	FavoriteReport,
	ReservationList,
	CancelledReservation,
	NoShowReservation,
	VoidedReservation,
	GroupReservation,
	ExpectedArrival,
	ArrivalList,
	SamedayReservation,
	AdvancedPaymentDeposit,
	BalanceDeposit,
	WaitListReservation,
	CurrentInHouse,
	GuestInHouse,
	GuestInHouseByBusinessSource,
	GuestInHouseByMarket,
	GuestInHouseByGuestType,
	GuestInHouseByCountry,
	GuestInHouseByState,
	GuestInHouseListing,
	GuestInHouseByCity,
	GuestInHouseByBookingSource,
	MasterFolio,
	DeskFolio,
	GuestInForecast,
	RepeaterGuest,
	FolioList,
	GuestInHouseByNationality,
	GuestList,
	GuestInHouseBreakfast,
	IncognitoGuest,
	ComplimentGuest,
	HouseUseGuest,
	EarlyCheckIn,
	DayUse,
	EarlyDeparture,
	ExpectedDeparture,
	ExtendedDeparture,
	DepartureList,
	ActualDepartureGuestList,
	FolioTransaction,
	DailyFolioTransaction,
	MonthlyFolioTransaction,
	YearlyFolioTransaction,
	ChargeList,
	DailyChargeList,
	MonthlyChargeList,
	YearlyChargeList,
	CashierReport,
	PaymentList,
	DailyPaymentList,
	MonthlyPaymentList,
	YearlyPaymentList,
	ExportCsvByDepartureDate,
	GuestLedger,
	GuestDeposit,
	GuestAccount,
	DailySales,
	DailyRevenueReport,
	DailyRevenueReportSummary,
	PaymentBySubDepartment,
	PaymentByAccount,
	RevenueBySubDepartment,
	FolioOpenBalance,
	Correction,
	VoidList,
	CancelCheckIn,
	LostAndFound,
	CashierReportReprint,
	PackageSales,
	CashSummaryReport,
	TransactionReportByStaff,
	TaxBreakdownDetailed,
	TodayRoomRevenueBreakdown,
	CancelCheckOut,
	Member,
	Voucher,
	VoucherSoldRedemeedAndComplimented,
	RoomList,
	RoomType,
	RoomRate,
	RoomCountSheet,
	RoomUnavailable,
	RoomSales,
	RoomHistory,
	RoomTypeAvailability,
	RoomTypeAvailabilityDetail,
	RoomStatus,
	RoomCountSheetByBuildingFloorRoomType,
	RoomCountSheetByRoomTypeBedType,
	RoomSalesByRoomNumber,
	RoomRateBreakdown,
	Package,
	PackageBreakdown,
	RoomRateStructure,
	RoomRevenue,
	Sales,
	SalesSummary,
	FrequentlySales,
	CaptainOrderList,
	CancelledCaptainOrder,
	VoidedCheckList,
	CancelledCaptainOrderDetail,
	BreakfastControl,
	GuestProfile,
	FrequentlyGuest,
	Company,
	PhoneBook,
	ContractRate,
	EventList,
	ReservationChart,
	ReservationGraphic,
	OccupiedGraphic,
	OccupiedByBusinessSourceGraphic,
	OccupiedByMarketGraphic,
	OccupiedByGuestTypeGraphic,
	OccupiedByCountryGraphic,
	OccupiedByStateGraphic,
	OccupancyGraphic,
	RoomAvailabilityGraphic,
	RoomUnvailabilityGraphic,
	RevenueGraphic,
	PaymentGraphic,
	RoomStatistic,
	GuestForecastReport,
	CityLedgerContributionAnalysis,
	FoodAndBeverageStatistic,
	RoomProduction,
	BusinessSourceProductivity,
	GuestForecastReportYearly,
	MarketStatistic,
	GuestForecastComparison,
	DailyFlashReport,
	DailyHotelCompetitor,
	DailyStatisticReport,
	RateCodeAnalysis,
	SalesContributionAnalysis,
	LeadList,
	TaskList,
	ProposalList,
	ActivityLog,
	SalesActivityDetail,
	CityLedgerList,
	CityLedgerAgingReport,
	CityLedgerAgingReportDetail,
	CityLedgerInvoice,
	CityLedgerInvoiceDetail,
	CityLedgerInvoicePayment,
	CityLedgerInvoiceMutation,
	BankReconciliation,
	BankTransactionList,
	BankTransactionAgingReport,
	BankTransactionAgingReportDetail,
	BankTransactionMutation,
	ApRefundDepositList,
	ApRefundDepositAgingReport,
	ApRefundDepositAgingReportDetail,
	ApRefundDepositPayment,
	ApRefundDepositMutation,
	ApCommissionAndOtherList,
	ApCommissionAndOtherAgingReport,
	ApCommissionAndOtherAgingReportDetail,
	ApCommissionAndOtherPayment,
	ApCommissionAndOtherMutation,
	LogUser,
	LogMoveRoom,
	LogTransferTransaction,
	LogSpecialAccess,
	KeylockHistory,
	LogVoidTransaction,
	LogHouseKeeping,
	LogPabx,
	LogShift,
	Product,
	ProductSales,
	DayendCloseReprint,
	ProductCategory,
	ProductList,
	ProductCosting,
	ProductCostingItem,
	FAndBRateStructure,
	RemovedProductCaptainOrder,
	RealizationCostOfGoodSold,
	MasterData,
	GeneralAndSales,
	Venue,
	Theme,
	SeatingPlan,
	BanquetBookingDetail,
	BanquetCalender,
	BanquetForecast,
	BanquetAdvancedDeposit,
	BanquetBalanceDeposit,
	BanquetChargeList,
	BanquetDailySales,
	CancelBanquetReservation,
	VoidBanquetReservation,
	BanquetBooking,
	AccountReceivable,
	AccountPayable,
	GeneralLedgerAndBank,
	AccountReceivableList,
	AccountReceivableAgingReport,
	AccountReceivableAgingReportDetail,
	AccountReceivablePayment,
	AccountReceivableMutation,
	AccountPayableList,
	AccountPayableAgingReport,
	AccountPayableAgingReportDetail,
	AccountPayablePayment,
	AccountPayableMutation,
	OperationalAccount,
	ChartOfAccount,
	Journal,
	GeneralLedger,
	TrialBalance,
	WorkSheet,
	BalanceSheet,
	IncomeStatement,
	ProfitAndLoss,
	ProfitAndLossDetail,
	ProfitAndLossByDepartment,
	ProfitAndLossDetailByDepartment,
	BankBookAccount,
	CurrentAssetAccount,
	CashOnHandAccount,
	FixedAssetAccount,
	OtherAssetAccount,
	ARLedgerAccount,
	CurrentLiabilityAccount,
	APLedgerAccount,
	LongTermLiabilityAccount,
	OtherLiabilityAccount,
	BankBookAccountGroup,
	BankBookAccountSummary,
	ProfitAndLossBySubDepartment,
	ProfitAndLossDetailBySubDepartment,
	ExportedJournal,
	BalanceMultiPeriod,
	CashFlow,
	ProfitAndLossMultiPeriodDetail,
	ProfitAndLossGraphic,
	Inventory,
	FixedAsset,
	Uom,
	InventoryItem,
	FixedAssetItem,
	InventoryPurchaseOrder,
	InventoryPurchaseOrderDetail,
	ReceiveStock,
	ReceiveStockDetail,
	StockTransfer,
	StockTransferDetail,
	Costing,
	CostingDetail,
	StockOpname,
	StockOpnameDetail,
	StoreStock,
	StoreStockCard,
	LowLevelStoreStock,
	HighLevelStoreStock,
	AllStoreStock,
	AllStoreStockCard,
	LowLevelAllStoreStock,
	HighLevelAllStoreStock,
	Production,
	ProductionDetail,
	CostRecipe,
	ReturnStock,
	ReturnStockDetail,
	DailyInventoryReconciliation,
	MonthlyInventoryReconciliation,
	InventoryReconciliation,
	AverageItemPricePurchase,
	ItemPurchasePriceGraphic,
	InventoryPurchaseRequest,
	InventoryPurchaseRequestDetail,
	StoreRequisition,
	StoreRequisitionDetail,
	RecapitulationFoodAndBeverage,
	RealizationCostOfGoodsSold,
	ComparisonCostSalesFbGraphic,
	FixedAssetPurchaseOrder,
	FixedAssetPurchaseOrderDetail,
	FixedAssetReceive,
	FixedAssetReceiveDetail,
	FixedAssetList,
	FixedAssetDepreciation int
}

type TProgramConfiguration struct {
	MemberPointExpire, VoucherLength, VoucherExpire, VoucherPointRedeem, LimitGrid, GridRowHeight, DueDateColor, LogoWidth, ReportHeaderAlignment, PMSServicePort, NotificationDuration, NotificationDelay, NotificationWidth, NotificationHeight, NotificationStyle, NotificationPosition int
	IsReplication, UseChildRate, IsRoomByName, PostFirstNightCharge, PostDiscount, IsRoomNumberRequired, IsBusinessSourceRequired, IsMarketRequired, IsTitleRequired, IsStateRequired, IsCityRequired, IsNationalityRequired, IsPhone1Required, IsEmailRequired, IsCompanyRequired, IsPurposeOfRequired, IsTAVoucherRequired, IsHKNoteRequired, FilterRateByMarket, FilterRateByCompany, AlwaysShowPublishRate, LockFolioOnCheckIn, IsAccrualBase, FormMDIChild, ShowGridGroupByBox, GridCanEditCell,
	ShowTaxService, AllowZeroAmount, PrintRegFormAfterCheckIn, ShowRate, ShowTransferOnCashierReport, ShowComplimentOnCashierReport, IncomeBudgetCalculateEachDay, InsertServerIDOnNumber, VRVCRoomStatus, ShowRemarkOnRoomAvailability, IsCalculateAllRoomRevenueSubGroup,
	MemberAutoUpdateMemberProfile, PABXAutoChargeToFolio, CCMSSMReservationAsAllotment, CCMSSMSynchronizeReservation, CCMSSMSynchronizeAvailability, CCMSSMSynchronizeRate, MikrotikAutoCreateUser, MikrotikSameUserPassword, MikrotikAskForPassword, NotificationActive, CanCLOverLimit, CompanyContactPersonRequired, CompanyStreetRequired, CompanyCityRequired, CompanyCountryRequired, CompanyStateRequired, CompanyPostalCodeRequired, CompanyPhone1Required, AutoGenerateCompanyCode, VentazaUseLift,
	EmailReminderReservation, EmailOnCheckIn, EmailOnCheckOut, EmailOnGuestBirthDay, EmailSalesActivitySendReminder, FridayAsWeekend, SaturdayAsWeekend, SundayAsWeekend,
	ShowPurchasePricePRSR,
	LockTransactionDateInventory,
	SynchronizePOAndReceive,
	IsPurchasingApproval,
	IsCompanyPRApplyPriceMoreThanOne,
	WAReminderReservation, WAOnCheckIn, WAOnCheckOut, WAOnGuestBirthDay, WASalesActivitySendReminder, AutomaticCreateInvoiceCLAtCheckOut, ReceiveStockAPTwoDigitDecimal, AutoImportJournal bool
	UnitCode, BillFileName,
	DefaultFolio, RegistrationFormReservation, RegistrationFormInHouse, ConfirmationLetter, ConfirmationLetterSelected, GuaranteedLetter, ProformaInvoice, ProformaInvoiceDetail, FolioFooter, KeylockVendor, PMSServiceAddress, SureLockHotelID, KimaOperatorCode, VingCardKeyType, VingCardUserGroup,
	InvoiceRemark, InvoiceNote, InvoiceAPNote, InvoiceTemplate, TaxServiceRemark, GuaranteeLetterRemark, BankName1, BankAccount1, HolderName1, BankName2, BankAccount2, HolderName2,
	ActiveStoreCode, CompanyTypeTravelAgent, VoucherDescription, VoucherTemplate, VoucherSaleDiscountTemplate, VingCardSource, VingCardDestination, BeTechServerName, BeTechUserName, BeTechPassword, OnityServerName, OnityEncoderNumber,
	SaflokToStation, SaflokFromStation, SaflokRequestNumber, SaflokPassword, SaflokEncoderStation, VisionlinePMSAddress, ColcomClientNumber,
	SD01, SD02, SD03, SD04, SD05, SD06, SD07, SD08, SD09, SD10, SD11, SD12, SDN01, SDN02, SDN03, SDN04, SDN05, SDN06, SDN07, SDN08, SDN09, SDN10, SDN11, SDN12,
	PYAccount01, PYAccount02, PYAccount03, PYAccount04, PYAccount05, PYAccount06, PYAccount07, PYAccount08, PYAccount09, PYAccount10, PYAccount11, PYAccount12,
	PYNAccount01, PYNAccount02, PYNAccount03, PYNAccount04, PYNAccount05, PYNAccount06, PYNAccount07, PYNAccount08, PYNAccount09, PYNAccount10, PYNAccount11, PYNAccount12,
	PABXLANSMDRPassword, IPTVendor, IPTVPassword, IPTVAddress, CCMSVendor, CCMSSMUser, CCMSSMPassword, CCMSSMRequestorID, CCMSSMHotelCode, CCMSSMWSDL, MikrotikVersion, MikrotikVendor, MikrotikAddress, MikrotikUser, MikrotikHotspotServer, MikrotikUserProfile, MikrotikServiceAddress, MikrotikPassword, SubDepartmentAllCCAdmin, DownsAuthCode, RTIncomeStatement, PABXLocal, PABXSLJJ, PABXSLI, EmailServer, WAAPIVendor, WAAPIURL, WAAPIKEY string
	CostingMethod,
	DefaultShippingAddress string
	PABXChargePerSecond, PABXChargePerSecondOther, VoucherDefaultPrice                                                                                 float64
	DefaultBill, BeTechPort, BeTechReaderType, BeTechDatabaseType, OnityDBType, OnitySoftwareType, OnityPort, MikrotikPasswordOption, CompanyCodeDigit byte
	KeylockLimit, CheckOutLimit, Timezone                                                                                                              string

	//CAMS
	PRDHColor, PRCCColor, PRFNColor, PRRJColor, LowStockColor int
	ActiveStore, CompanyTypeSupplier, CompanyTypeExpedition, RFPurchaseRequest, RFPurchaseOrder, RFReceiveStock, RFStoreRequisition, RFStockTransfer, RFCosting, RFFAPurchaseOrder, RFFAReceive,
	RTJournalVoucher, RTJournalVoucherPaymentForm, RTJournalVoucherReceiveForm, DefaultStore, RTInventoryReconciliation, RTDailyInventoryReconciliation, RTMonthlyInventoryReconciliation string
	AutoRefreshOnFormActivate byte
}

type TProgramConfigurationCAMS struct {
	LimitGrid, GridRowHeight, LogoWidth, ReportHeaderAlignment, NotificationDuration, NotificationDelay, NotificationWidth, NotificationHeight, NotificationStyle, NotificationPosition, PRDHColor, PRCCColor, PRFNColor, PRRJColor, LowStockColor int
	UnitCode, ActiveStoreCode, ActiveStore, InvoiceAPNote, CompanyTypeSupplier, CompanyTypeExpedition, RFPurchaseRequest, RFPurchaseOrder, RFReceiveStock, RFStoreRequisition, RFStockTransfer, RFCosting, RFFAPurchaseOrder, RFFAReceive,
	RTJournalVoucher, RTJournalVoucherPaymentForm, RTJournalVoucherReceiveForm, DefaultShippingAddress, DefaultStore, RTIncomeStatement, RTInventoryReconciliation, RTDailyInventoryReconciliation, RTMonthlyInventoryReconciliation string
	CostingMethod string
	IsAccrualBase, FormMDIChild, ShowGridGroupByBox, GridCanEditCell, ShowTaxService, AllowZeroAmount, InsertServerIDOnNumber, AutoRefreshOnFormActivate, SynchronizePOAndReceive, NotificationActive, AutoGenerateCompanyCode, ReceiveStockAPTwoDigitDecimal, ShowPurchasePricePRSR,
	IsPurchasingApproval, IsCompanyPRApplyPriceMoreThanOne, CompanyContactPersonRequired, CompanyStreetRequired, CompanyCityRequired, CompanyCountryRequired, CompanyStateRequired, CompanyPostalCodeRequired, CompanyPhone1Required bool
	CompanyCodeDigit             byte
	LockTransactionDateInventory time.Time
	// FormatSettingX TFormatSettings;
}

type TConstProgramVariable struct {
	DefaultSystemCode, DateTimeFormatYYYYMMDD string
	DatabaseVersion, SettingFileName, WallpaperListFile, TempImageFile,
	InvoiceNumberPrefix, APNumberPrefix, DepreciationNumberPrefix, ARNumberPrefix, PaymentNumberPrefix, PRNumberPrefix, PONumberPrefix, ReceiveNumberPrefix, SRNumberPrefix, ProductionNumberPrefix, ReturnStockPrefix, OpnameNumberPrefix,
	FAPONumberPrefix, FAReceiveNumberPrefix, StockTransferNumberPrefix, CostingNumberPrefix, OTPPrefix, CompanyCodePrefix, CityOtherCode string
	AddKey, MinPasswordLength, MaxLoginTry, MaxOTPGenerate, ConfigurationAccess byte
}

type TConstTableName struct {
	Report, ReportDefaultField, ReportTemplate, ReportTemplateField, ReportGroupField, ReportOrderField, ReportGroupingField string
}
type TDepreciationType struct {
	None, StraightLine, DecliningBalance, DoubleDecliningBalance, SumOfYearDigit string
}

type TGlobalJournalAccountGroupName struct {
	Assets, Liability, Equity, Income, Cost, Expense1, Expense2, OtherIncome, OtherExpense string
}

type TStoreRequisitionStatus struct {
	NotApproved, Approved, Rejected string
}

type TFAItemCondition struct {
	Good, Broken, Repaired, Sold, FullyDepreciated, Transfered string
}

type TFALocationType struct {
	HKStore, Laundry, Room, Other string
}

type TItemGroupTypeCode struct {
	Food, Beverage string
}
type TItemGroupCode struct {
	Food, Beverage string
}

type TVariableDLL struct {
	KeyDLL, KeyOtherDLL, AddWordDLL, ActivateKeyDLL, StartExpireDateDLL, EndExpireDateDLL, ProductIDDLL, ProductKey string
	DemoDLL, GenerateHDDDLL                                                                                         bool
}

type TReceivedStatus struct {
	UnReceived, ReceivedPartial, Received int
}
type TSystemCode struct {
	General, Hotel, Pos, Banquet, Accounting, Asset, Payroll, Corporate, Report, Tools string
}

type TConfigurationCategory struct {
	General,
	Form,
	FormatSetting,
	Grid,
	Report,
	WeekendDay,
	ReportSignature,
	ReportTemplate,
	Other,
	Reservation,
	DefaultVariable,
	AmountPreset,
	CustomField,
	CustomLookupField,
	MemberVoucherGift,
	Accounting,
	GlobalAccount,
	GlobalJournalAccount,
	PaymentAccount,
	SubDepartment,
	GlobalDepartment,
	GlobalSubDepartment,
	GlobalJournalAccountSubGroup,
	GlobalOther,
	Personal,
	Invoice,
	Folio,
	OtherForm,
	CompanyBankAccount,
	PaymentCityLedger,
	Company,
	FloorPlan,
	RoomStatusColor,
	Keylock,
	PABX,
	RoomCosting,
	OtherHK,
	ServiceIPTV,
	ServiceCCMS,
	Mikrotik,
	Notification,
	WAAPI,
	Email,
	Inventory,
	DayendClosed string
}

type TConfigurationCategoryCAMS struct {
	General,
	Form,
	FormatSetting,
	Grid,
	Report,
	DefaultVariable,
	Other,
	Company,
	Reservation,
	Accounting,
	Personal,
	Invoice,
	Inventory,
	GlobalJournalAccount,
	GlobalJournalAccountSubGroup,
	GlobalSubDepartment,
	ReportForm,
	PurchaseRequestApp,
	ReportTemplate,
	Notification string
}

type TConfigurationName struct {
	//General
	Timezone,
	DatabaseVersion,
	IsReplication,
	UseChildRate,
	IsRoomByName,
	LogOffTriggerAddress,
	//Form
	FormMDIChild,
	RequiredFieldColor,
	//FormatSetting
	ShortDateFormat, DateSeparator, CurrencyFormat, DecimalSeparator, ThousandsSeparator,
	//Grid
	ShowGridGroupByBox, GridCanEditCell, LimitGrid, GridRowHeight, DueDateColor,
	//Report
	LogoWidth,
	ReportHeaderAligment,
	//Weekend Day
	FridayAsWeekend,
	SaturdayAsWeekend,
	SundayAsWeekend,
	//Report Signature
	ShowPreparedBy,
	PreparedBy,
	ShowCheckedBy,
	CheckedBy,
	ShowApprovedBy,
	ApprovedBy,
	ShowInDailyRevenueReport,
	ShowInRoomStatisticReport,
	//Report Template
	DepositReceipt,
	DepositReceiptRefund,
	FolioReceipt,
	FolioReceiptRefund,
	MiscellaneousCharge,
	CashierReport,
	CashRemittance,
	BreakfastList,
	PickUpService,
	DropService,
	HKChecklist,
	HKRoomStatus,
	HKRoomStatusSummary,
	HKRoomDiscrepancy,
	HKRoomAttendantControlSheet,
	DailyRevenueReport,
	DailyRevenueReportSummary,
	RevenueBySubDepartment,
	GuestInHouseListing,
	RoomStatistic,
	DailyFlashReport,
	RoomProduction,
	GuestForecast,
	RTIncomeStatement,
	//Other
	ShowTransferOnCashierReport,
	ShowComplimentOnCashierReport,
	IncomeBudgetCalculateEachDay,
	InsertServerIDOnNumber,
	CalculateAllRoomRevenueSubGroup,
	CompanyTypeTravelAgent, CompanyTypeSupplier,
	//Reservation
	PostFirstNightCharge, PostDiscount, KeylockLimit, CheckOutLimit,
	IsRoomNumberRequired, IsBusinessSourceRequired, IsMarketRequired, IsTitleRequired, IsStateRequired, IsCityRequired, IsNationalityRequired, IsPhone1Required, IsEmailRequired, IsCompanyRequired, IsPurposeOfRequired, IsTAVoucherRequired, IsHKNoteRequired, FilterRateByMarket, FilterRateByCompany,
	AlwaysShowPublishRate, LockFolioOnCheckIn, AutoGenerateCompanyCode, CompanyCodeDigit,
	//Default Variable
	DVRoomType,
	DVRoomRate,
	DVSubDepartment,
	DVPaymentType,
	DVComplimentRate,
	DVHouseUseRate,
	DVMarket,
	DVIndividualMarket,
	//Amount Preset
	APPreset1,
	APPreset2,
	APPreset3,
	APPreset4,
	APPreset5,
	APPreset6,
	APPreset7,
	APPreset8,
	//Custom Field
	CustomFieldName01,
	CustomFieldName02,
	CustomFieldName03,
	CustomFieldName04,
	CustomFieldName05,
	CustomFieldName06,
	CustomFieldName07,
	CustomFieldName08,
	CustomFieldName09,
	CustomFieldName10,
	CustomFieldName11,
	CustomFieldName12,
	CustomFieldDefaultValue01,
	CustomFieldDefaultValue02,
	CustomFieldDefaultValue03,
	CustomFieldDefaultValue04,
	CustomFieldDefaultValue05,
	CustomFieldDefaultValue06,
	CustomFieldDefaultValue07,
	CustomFieldDefaultValue08,
	CustomFieldDefaultValue09,
	CustomFieldDefaultValue10,
	CustomFieldDefaultValue11,
	CustomFieldDefaultValue12,
	//CustomLookupField
	CustomLookupFieldName01,
	CustomLookupFieldName02,
	CustomLookupFieldName03,
	CustomLookupFieldName04,
	CustomLookupFieldName05,
	CustomLookupFieldName06,
	CustomLookupFieldName07,
	CustomLookupFieldName08,
	CustomLookupFieldName09,
	CustomLookupFieldName10,
	CustomLookupFieldName11,
	CustomLookupFieldName12,
	CustomLookupFieldDefaultValue01,
	CustomLookupFieldDefaultValue02,
	CustomLookupFieldDefaultValue03,
	CustomLookupFieldDefaultValue04,
	CustomLookupFieldDefaultValue05,
	CustomLookupFieldDefaultValue06,
	CustomLookupFieldDefaultValue07,
	CustomLookupFieldDefaultValue08,
	CustomLookupFieldDefaultValue09,
	CustomLookupFieldDefaultValue10,
	CustomLookupFieldDefaultValue11,
	CustomLookupFieldDefaultValue12,
	//Member, Voucher and Gift
	MemberPointExpire,
	MemberAutoUpdateMemberProfile,
	VoucherLength,
	VoucherExpire,
	VoucherPointRedeem,
	VoucherDefaultPrice,
	VoucherDescription,
	VoucherTemplate,
	VoucherSaleDiscountTemplate,
	//Accounting
	IsAccrualBase,
	SubDepartmentAllCCAdmin,
	AutomaticCreateInvoiceCLAtCheckOut,
	//Global Account
	AccountRoomCharge,
	AccountExtraBed,
	AccountCancellationFee,
	AccountNoShow,
	AccountBreakfast,
	AccountTelephone,
	AccountAPRefundDeposit,
	AccountAPCommission,
	AccountCreditCardAdm,
	AccountCash,
	AccountCityLedger,
	AccountVoucher,
	AccountVoucherCompliment,
	AccountTax,
	AccountService,
	AccountTransferDepositReservation,
	AccountTransferDepositReservationToFolio,
	AccountTransferCharge,
	AccountTransferPayment,
	//Global Journal Account
	JAGuestLedger,
	JAAPVoucher,
	JAGuestDeposit,
	JAGuestDepositReservation,
	JAPLBeginningYear,
	JAPLCurrentYear,
	JAPLCurrency,
	JAIncomeVoucherExpire,
	JAInvoiceDiscount,
	JACreditCardAdm,
	JAAPPaymentTemporer,
	JAIncomeVoucherExpired,
	JAOverShortAsIncome,
	JAOverShortAsExpense,
	JAServiceRevenue,
	//Global Journal Account Sub Group
	JASGInventory,
	JASGFixedAsset,
	JASGAccmDepreciation,
	JASGAccountPayable,
	JASGManagementFee,
	JASGDepreciation,
	JASGAmortization,
	JASGLoanInterest,
	JASGIncomeTax,
	//Payment Account
	PYAccount01,
	PYAccount02,
	PYAccount03,
	PYAccount04,
	PYAccount05,
	PYAccount06,
	PYAccount07,
	PYAccount08,
	PYAccount09,
	PYAccount10,
	PYAccount11,
	PYAccount12,
	PYNAccount01,
	PYNAccount02,
	PYNAccount03,
	PYNAccount04,
	PYNAccount05,
	PYNAccount06,
	PYNAccount07,
	PYNAccount08,
	PYNAccount09,
	PYNAccount10,
	PYNAccount11,
	PYNAccount12,
	//Sub Department
	SD01,
	SD02,
	SD03,
	SD04,
	SD05,
	SD06,
	SD07,
	SD08,
	SD09,
	SD10,
	SD11,
	SD12,
	SDN01,
	SDN02,
	SDN03,
	SDN04,
	SDN05,
	SDN06,
	SDN07,
	SDN08,
	SDN09,
	SDN10,
	SDN11,
	SDN12,
	//Global Department
	DRoomDivision,
	DFoodBeverage,
	DBanquet,
	DMinor,
	DMiscellaneous,
	//Global Sub Department
	SDFrontOffice,
	SDHouseKeeping,
	SDBanquet,
	SDAccounting,
	//Global Other
	GOPaymentType,
	//Personal
	DRName,
	DRPosition,
	HMName,
	HMPosition,
	FHName,
	FHPosition,
	IAName,
	IAPosition,
	GCName,
	GCPosition,
	ARName,
	ARPosition,
	APName,
	APPosition,
	//Invoice
	InvoiceRemark,
	InvoiceNote,
	InvoiceAPNote,
	InvoiceTemplate,

	//Folio
	DefaultFolio, FolioFooter, AllowZeroAmount, PrintRegFormAfterCheckIn, ShowRate,
	//Other Form
	RegistrationFormReservation, RegistrationFormInHouse, ConfirmationLetter, ConfirmationLetterSelected, GuaranteedLetter, ProformaInvoice, ProformaInvoiceDetail, TaxServiceRemark, GuaranteeLetterRemark,
	//Company Bank Account
	BankName1, BankAccount1, HolderName1,
	BankName2, BankAccount2, HolderName2,
	//Payment City Ledger
	CanCLOverLimit,
	//Company
	CompanyContactPersonRequired, CompanyStreetRequired, CompanyCityRequired, CompanyCountryRequired, CompanyStateRequired, CompanyPostalCodeRequired, CompanyPhone1Required,
	//Floor Plan
	MinRoomWidth, MinRoomHeight, LeftMargin, TopMargin, TileColumn, TileDistance, SnapGrid,
	//Room Status Color
	Reserved, Occupied, HouseUse, Compliment, OutOfOrder, OfficeUse, UnderConstruction, Available,
	//Keylock
	KeylockVendor, PMSServiceAddress, PMSServicePort, SureLockHotelID, KimaOperatorCode,
	VingCardKeyType, VingCardUserGroup, VingCardSource, VingCardDestination,
	BeTechPort, BeTechReaderType, BeTechDatabaseType, BeTechServerName, BeTechUserName, BeTechPassword,
	OnityDBType, OnityServerName, OnitySoftwareType, OnityPort,
	SaflokToStation, SaflokFromStation, SaflokRequestNumber, SaflokPassword, SaflokEncoderStation,
	VisionlinePMSAddress, DownsAuthCode, VentazaUseLift,
	//PABX
	PABXChargePerSecond, PABXChargePerSecondOther, PABXAutoPostToFolio, PABXLANSMDRPassword, PABXLocal, PABXSLJJ, PABXSLI,
	//Room Costing
	DefaultStore,
	//Other HK
	VRVCRoomStatus,
	ShowRemarkOnRoomAvailability,
	//Service IPTV
	IPTVVendor,
	IPTVPassword,
	IPTVSubFolio,
	IPTVAddress,
	//Service CCMS
	CCMSVendor, CCMSSMUser, CCMSSMPassword, CCMSSMRequestorID, CCMSSMHotelCode, CCMSSMWSDL, CCMSSMReservationAsAllotment, CCMSSMSynchronizeReservation, CCMSSMSynchronizeAvailability, CCMSSMSynchronizeRate, CCMSSTAAHGlobalPercentAvailability, CCMSSTAAHGlobalMinRoomLeft,
	//Mikrotik
	MikrotikVersion, MikrotikVendor, MikrotikAddress, MikrotikUser, MikrotikPassword, MikrotikHotspotServer, MikrotikUserProfile, MikrotikServiceAddress, MikrotikAutoCreateUser, MikrotikSameUserPassword, MikrotikAskForPassword, MikrotikPasswordOption,
	//Notification
	NotificationActive, NotificationDuration, NotificationDelay, NotificationWidth, NotificationHeight, NotificationStyle, NotificationPosition,
	//WAAPI
	WAAPIVendor, WAAPIKey, WAAPIURL, WAAPITemplateCheckIn, WAAPITemplateCheckOut, WAOnReminderReservation, WAOnCheckIn, WAOnCheckOut, WAOnGuestBirthday, WAOnSalesReminder,
	//Email
	EmailServer, EmailAuthentication, EmailPort, EmailUser, EmailPassword, EmailUserSSL, EmailUseStartTLS,
	EmailReminderReservation, EmailOnCheckIn, EmailOnCheckOut, EmailOnGuestBirthday, EmailSalesActivitySendreminder,
	//Inventory
	CostingMethod,
	SynchronizePOAndReceive,
	DefaultShippingAddress,
	ReceiveStockAPTwoDigitDecimal,
	ShowPurchasePricePRSR,
	LockTransactionDateInventory,
	IsPurchasingApproval,
	IsCompanyPRApplyPriceMoreThanOne,
	// Dayend Closed
	AutoImportJournal string
	//General
	//Personal
	CCName,
	CCPosition,
	PMName,
	PMPosition,
	//Invoice
	//Purchase Request Approval
	UserApproval1,
	UserApproval2,
	UserApproval3,
	//SR & PR Color
	DHColor,
	CCColor,
	FNColor,
	RJColor,
	//All Store Stock & Store Stock
	LowStockColor,
	//Inventory
	//Company
	//Global Journal Accounting
	JAAPSupplier,
	JAPurchasingDiscount,
	JAPurchasingTax,
	JAPurchasingShipping,
	JAIncomeReturnStock,
	JAExpenseReturnStock,
	//Global Journal Sub Group
	//Global Sub Department
	//Other
	CompanyTypeExpedition,
	//Report Form
	PurchaseRequest,
	PurchaseOrder,
	ReceiveStock,
	StoreRequitition,
	StockTransfer,
	Costing,
	FAPurchaseOrder,
	FAReceive,
	//Report Template
	RTJournalVoucher,
	RTJournalVoucherPaymentForm,
	RTJournalVoucherReceiveForm,
	RTInventoryReconciliation,
	RTDailyInventoryReconciliation,
	RTMonthlyInventoryReconciliation string
}
type TConfigurationNamePOS struct {
	//General
	DatabaseVersion,
	RoomServiceOutlet,
	LogOffTriggerAddress,
	IsRoomByName,
	//Payment
	CanPaymentIfChargeZero,
	//Other POS
	MemberMustScanFinger,
	//Form
	FormMDIChild,
	RequiredFieldColor,
	//FormatSetting
	ShortDateFormat, DateSeparator, CurrencyFormat, DecimalSeparator, ThousandsSeparator,
	//Grid
	ShowGridGroupByBox, GridCanEditCell, LimitGrid, GridRowHeight,
	//Report
	LogoWidth,
	ReportHeaderAligment,
	//Default Variable
	DVSubDepartment,
	DVIPTVMarket,
	//Other
	ShowTransferOnCashierReport,
	ShowComplimentOnCashierReport,
	InsertServerIDOnNumber,
	//Bill
	DefaultBill,
	BillFileName,
	CaptainOrderFileName,
	CaptainOrderStationFileName,
	//Report Template
	RTCashierReport,
	RTCashRemittance,
	RTFNBRateStructure,
	RTIncomeStatement,
	//Accounting
	IsAccrualBase,
	SubDepartmentAllCCAdmin,
	//Other
	ShowTaxService, PostDiscount, AutoCostingCostRecipeonCloseTransaction, AutoGenerateCompanyCode, CompanyCodeDigit, CompanyTypeSPATherapist,
	//Company
	CompanyContactPersonRequired, CompanyStreetRequired, CompanyCityRequired, CompanyCountryRequired, CompanyStateRequired, CompanyPostalCodeRequired, CompanyPhone1Required,
	//Global Account
	AccountRoomCharge,
	AccountAPRefundDeposit,
	AccountAPCommission,
	AccountCreditCardAdm,
	AccountCash,
	AccountCityLedger,
	AccountTax,
	AccountService,
	AccountTransferDepositReservation,
	AccountTransferDepositReservationToFolio,
	AccountTransferCharge,
	AccountTransferPayment,
	//Global Journal Account
	JAGuestLedger,
	JAGuestDeposit,
	JAGuestDepositReservation,
	//Global Sub Department
	SDAccounting,
	//Table View
	MinRoomWidth, MinRoomHeight, LeftMargin, TopMargin, TileColumn, TileDistance, SnapGrid,
	TableColor, COColor1, COColor2, COColor3, COColor4, COColor5,
	//Reservation
	RVRangeTimeCODefault,
	//Notification
	NotificationActive, NotificationDuration, NotificationDelay, NotificationWidth, NotificationHeight, NotificationStyle, NotificationPosition,
	//Cash Drawer
	CashDrawerComport, CashDrawerOpenCode1, CashDrawerOpenCode2, CashDrawerOpenCode3, CashDrawerOpenCode4, CashDrawerOpenCode5, CashDrawerOpenCode6,
	//Customer Display
	CDVendor, CDPort, CDBaundRate, CDParity, CDDataWidth, StopBit,
	//Kitchen Printer
	PrintCOAfterChangeRemoveCancel,
	//Inventory
	CostingMethod,
	SynchronizePOAndReceive,
	DefaultShippingAddress,
	ReceiveStockAPTwoDigitDecimal,
	ShowPurchasePricePRSR,
	LockTransactionDateInventory,
	IsPurchasingApproval,
	IsCompanyPRApplyPriceMoreThanOne string
}

type TCacheKey struct {
	LastAPNumber,
	LastARNumber,
	LastRefNumber,
	LastManualRefNumber, LastTransactionRefNumber, LastDisbursementRefNumber, LastReceiveRefNumber, LastInventoryRefNumber, LastAdjustmentRefNumber,
	LastFixedAssetRefNumber, LastBeginningYearRefNumber, LastCostingNumber string

	LastInvoiceNumber, LastDepreciationNumber, LastPaymentNumber, LastPRNumber, LastPONumber, LastReceiveNumber, LastSRNumber, LastProductionNumber, LastReturnStock, LastOpnameNumber,
	LastFAPONumber, LastFAReceiveNumber, LastStockTransferNumber string
}
type TConfigurationNameCAMS struct {
	//General
	DatabaseVersion,
	//Form
	FormMDIChild,
	RequiredFieldColor,
	//FormatSetting
	ShortDateFormat, DateSeparator, CurrencyFormat, DecimalSeparator, ThousandsSeparator,
	//Grid
	ShowGridGroupByBox, GridCanEditCell, LimitGrid, GridRowHeight,
	//Report
	LogoWidth,
	ReportHeaderAligment,
	//Default Variable
	DVSubDepartment,
	//Reservation
	AutoGenerateCompanyCode, CompanyCodeDigit,
	//Other
	InsertServerIDOnNumber,
	//Accounting
	IsAccrualBase,
	//Personal
	DRName,
	DRPosition,
	HMName,
	HMPosition,
	FHName,
	FHPosition,
	GCName,
	GCPosition,
	CCName,
	CCPosition,
	PMName,
	PMPosition,
	//Invoice
	InvoiceAPNote,
	//Purchase Request Approval
	UserApproval1,
	UserApproval2,
	UserApproval3,
	//SR & PR Color
	DHColor,
	CCColor,
	FNColor,
	RJColor,
	//All Store Stock & Store Stock
	LowStockColor,
	//Inventory
	CostingMethod,
	SynchronizePOAndReceive,
	DefaultShippingAddress,
	ReceiveStockAPTwoDigitDecimal,
	ShowPurchasePricePRSR,
	LockTransactionDateInventory,
	IsPurchasingApproval,
	IsCompanyPRApplyPriceMoreThanOne,
	//Company
	CompanyContactPersonRequired, CompanyStreetRequired, CompanyCityRequired, CompanyCountryRequired, CompanyStateRequired, CompanyPostalCodeRequired, CompanyPhone1Required,
	//Global Journal Accounting
	JAAPSupplier,
	JAPurchasingDiscount,
	JAPurchasingTax,
	JAPurchasingShipping,
	JAIncomeReturnStock,
	JAExpenseReturnStock,
	//Global Journal Sub Group
	JASGInventory,
	JASGFixedAsset,
	JASGAccmDepreciation,
	JASGAccountPayable,
	//Global Sub Department
	SDAccounting,
	//Other
	CompanyTypeSupplier,
	CompanyTypeExpedition,
	//Report Form
	PurchaseRequest,
	PurchaseOrder,
	ReceiveStock,
	StoreRequitition,
	StockTransfer,
	Costing,
	FAPurchaseOrder,
	FAReceive,
	//Report Template
	RTJournalVoucher,
	RTJournalVoucherPaymentForm,
	RTJournalVoucherReceiveForm,
	RTIncomeStatement,
	RTInventoryReconciliation,
	RTDailyInventoryReconciliation,
	RTMonthlyInventoryReconciliation,
	//Notification
	NotificationActive, NotificationDuration, NotificationDelay, NotificationWidth, NotificationHeight, NotificationStyle, NotificationPosition string
}
type TConfigurationCategoryPOS struct {
	General,
	Payment,
	Form,
	FormatSetting,
	Grid,
	Report,
	DefaultVariable,
	Other,
	Company,
	Bill,
	ReportTemplate,
	Accounting,
	GlobalAccount,
	GlobalJournalAccount,
	GlobalSubDepartment,
	TableView,
	Reservation,
	Notification,
	CashDrawer,
	CustomerDisplay,
	KitchenPrinter,
	Inventory string
}

type PrimaryKey struct {
	Id     uint64 `json:"id"`
	Number string `json:"number"`
}
type TCalenderLabelColor struct {
	Reservation, InHouse, CheckOut int
}
type TTempVariable struct {
	GuestProfileID, StringValue1, CardHolder, CardNumber, ValidMonth, ValidYear, BookingNumberGlobal, ReservationStatusCode                                                      string
	EmailMessage                                                                                                                                                                 string
	CharValue1                                                                                                                                                                   string
	FormWalkInSuccess, CancelReservationSuccess, CheckOutSuccess, AutoPostingSuccess, ExtendSuccess, ChangeArrivalSuccess, FolioBalanceSuccess, ReservationDepositBalanceSuccess bool
	SelectDate1, AuditDateImportX                                                                                                                                                time.Time
	ModalResultX1, ModalResultX2, ModalResultX3                                                                                                                                  int
}

type TFloorPlanVariable struct {
	MinRoomWidth, MinRoomHeight, LeftMargin, TopMargin, TileColumn, TileDistance, SnapGrid int
}

type TDefaultVariable struct {
	RoomType, RoomRate, SubDepartment, PaymentType, ComplimentRate, HouseUseRate, Market, IndividualMarket string
}

type TFormMessage struct {
	GuestInHouse        []string
	Housekeeping        []string
	Unavailable         []string
	FolioHistory        []string
	PostingRoomCharge   []string
	PostingExtraCharge  []string
	TransactionTerminal []string
	RoomRate            []string
	APARPayment         []string
}

type TProgramMessage struct {
	DataNotComplete, DeleteData, VoidData, AreYouSure, DuplicateEntry, EmptyData, CannotEditDelete,
	Confirmation, Information, Attention, Error, CannotAccess, ProcessCompleted, NoDataFound, NoDataDetail, PleaseReopenThisForm,
	InvalidDateRange, ReportFileNotFound, JournalPeriodIsClosed, PaymentDateCannotLowerThanIssuedDate, PaymentDateIsGreaterThanTodayDate, CannotUpdateDataAlreadyUsed, CannotDeleteDataAlreadyUsed string
}

type TSQLParameterName struct {
	Param01, Param02, Param03, Param04, Param05, Param06, Param07, Param08, Param09, Param10 string
}

type TUserInfo struct {
	ID, Name, Password, GroupCode, WorkingShift, SessionID, CompanyCode string
	AccessLevel                                                         int
	LogShiftID                                                          uint64
}

type TUserAccessType struct {
	Form, Report, Special, Keylock, Reservation, Deposit, InHouse, WalkIn, Folio, FolioHistory, FloorPlan, Company, Invoice, MemberVoucherGift, PreviewReport, PaymentByAPAR string
}
type TReportAccessType struct {
	Form, Preview, FrontDesk, Pos, Banquet, Accounting, Asset string
}
type TGlobalJournalAccountSubGroup struct {
	Inventory, FixedAsset, AccumulatedDepreciation, AccountPayable, ManagementFee, Depreciation, Amortization, LoanInterest, IncomeTax string
}

type TUserGroup struct {
	System, SuperAdmin, Admin string
}

type TUserFormAccessOrder struct {
	Summary, FloorPlan, RoomAvailability, RoomTypeAvailability, RoomAllotment, GuestProfile, GuestGroup, Reservation, WalkIn, GuestInHouse, MasterFolio, DeskFolio, FoliosAndTransaction, FolioHistory,
	HouseKeeping, RoomCosting, LostAndFound, CashierReport, GlobalPostTransaction, AutoPostTransaction, DayendClose, Company, Package, RoomRate, EventList, PhoneBook, Member, Voucher, Gift,
	APRefundDeposit, APCommission, ARCityLedger, ARCityLedgerInvoice, Receipt, BankTransaction, BankReconciliation, Cheque, ExportJournal, IncomeBudget, BudgetStatistic, Report, DataAnalysis,
	Configuration, UserSetting, OneTimePassword, UserActivityLog, PABXSMDRViewer, BackupDatabase, RestoreDatabase, OptimizeData, Notification, BreakfastControl, CompetitorData, GuestLoanItem, SalesActivity, DynamicRate, RoomRateLastDeal byte
}

type TUserReportAccessOrder struct {
	ReservationList, CancelledReservation, NoShowReservation, VoidedReservation, GroupReservation, ExpectedArrival, ArrivalList, SamedayReservation, AdvancedPaymentDeposit, BalanceDeposit, WaitListReservation,
	CurrentInHouse, GuestInHouse, GuestInHousebyBusinessSource, GuestInHousebyMarket, GuestInHousebyGuestType, GuestInHousebyCountry, GuestInHousebyState, MasterFolio, DeskFolio,
	IncognitoGuest, ComplimentGuest, HouseUseGuest, EarlyCheckIn, DayUse, EarlyDeparture, ExpectedDeparture, ExtendedDeparture, DepartureList,
	FolioTransaction, DailyFolioTransaction, MonthlyFolioTransaction, YearlyTransaction,
	ChargeList, DailyChargeList, MonthlyChargeList, YearlyChargeList, CashierReport,
	PaymentList, DailyPaymentList, MonthlyPaymentList, YearlyPaymentList,
	ExportCSVbyDepartureDate, GuestLedger, GuestDeposit, GuestAccount, DailySales, DailyRevenueReport, DailyRevenueReportSummary,
	FolioOpenBalance, Correction, VoidList, CancelCheckIn, LostandFound,
	RoomList, RoomType, RoomRate, RoomCountSheet, RoomCountSheetByBuildingFloorRoomType, RoomCountSheetByRoomTypeBedType, RoomUnavailable,
	RoomSales, RoomHistory, RoomTypeAvailability, RoomTypeAvailabilityDetail, RoomStatus,
	Sales, SalesSummary, FrequentlySales, CaptainOrderList, CancelledCaptainOrder, VoidedCheckList,
	GuestProfile, FrequentlyGuest, Company, PhoneBook,
	ContractRate, EventList, ReservationChart, ReservationGraphic,
	OccupiedGraphic, OccupiedbyBusinessSourceGraphic, OccupiedbyMarketGraphic, OccupiedbyGuestTypeGraphic, OccupiedbyCountryGraphic, OccupiedbyStateGraphic, OccupancyGraphic,
	RoomAvailabilityGraphic, RoomUnvailabilityGraphic, RevenueGraphic, PaymentGraphic, RoomStatistic, GuestForecastReport, CityLedgerContributionAnalysis,
	CityLedgerList, CityLedgerAgingReport, CityLedgerAgingReportetail, CityLedgerInvoice, CityLedgerInvoiceDetail, CityLedgerPayment, CityLedgerMutation, BankReconciliation,
	APRefundDepositList, APRefundDepositAgingReport, APRefundDepositAgingReportetail, APRefundDepositPayment, APRefundDepositMutation,
	APCommissionList, APCommissionAgingReport, APCommissionAgingReportetail, APCommissionPayment, APCommissionMutation,
	LogUser, LogMoveRoom, LogTransferTransaction, LogSpecialAccess, KeyLockHistory, LogVoidTransaction, LogHouseKeeping, RoomRateBreakdown, Package, PackageBreakdown, FBStatistic, Member, Voucher, VoucherSRC, CancelledCaptainOrderDetail, RoomRateStructure,
	PaymentBySubDepartment, PaymentByAccount, RoomProduction, RoomRevenue, OTAProductivity, GuestForecastReportYearly, CashierReportReprint, GuestInHouseListing, GuestInHouseForecast, BankTransactionList, RepeaterGuest, MarketStatistic,
	PackageSales, LogPABX, BankTransactionAgingReport, BankTransactionAgingReportDetail, BankTransactionMutation, FolioList, LogShift, GuestForecastComparison, CashSummaryReport, TransactionByStaff, TaxBreakDownDetailed, DailyFlashReport, BreakfastControl, DailyHotelCompetitor, ActualDepartureGuestList,
	TodayRoomRevenueBreakdown, RoomSalesByRoomNumber, DailyStatisticReport, RateCodeAnalysis, GuestInHouseByCity, CancelCheckOut, SalesContributionAnalysis, GuestInHousebyBookingSource, GuestInHouseByNationality, GuestList,
	LeadList, TaskList, ProposalList, ActivityLog, SalesActivityDetail, RevenueBySubDepartment, GuestInHouseBreakfast byte
}

type TUserSpecialAccessOrder struct {
	UnlockReservation, VoidReservation, VoidDeposit, CorrectDeposit, DecreaseStay, BusinessSource, OverrideRateDiscount, ModifyScheduleRate, ModifyBreakdown, ModifyExtraCharge, ComplimentGuest, HouseUseGuest,
	MoveRoom, VoidSubFolio, CorrectSubFolio, CancelCheckIn, CancelCheckOut, CreateMasterFolio, PrintInvoice, ModifyClosedJournal, TransferToDeskFolio, TransferToMasterFolio, UpdateGuestName byte
}

type TUserKeylockAccessOrder struct {
	CheckInWithoutCard, CheckOutWithoutCard, IssuedCardMoreThanTwice, ModifyArrivalDate, ModifyDepartureDate, ModifyDepartureTime, DepartureDate1Night, ShowAccessIssuedCardMoreThanOne, CanIssuedCardMoreThanOne byte
}

type TUserReservationAccessOrder struct {
	Insert, Update, Duplicate, Deposit, Cancel, Void, NoShow, AutoAssign, Lock, CheckIn, InsertFromAllotment, Keylock byte
}

type TUserDepositAccessOrder struct {
	Insert, Cash, Card, Refund, Void, Correction, Transfer, UpdateSubDepartment, UpdateRemark, UpdateDocumentNumber byte
}

type TUserInHouseAccessOrder struct {
	Transaction, Update, Keylock, Compliment, HouseUse, MoveRoom, SwitchRoom, LockFolio, CancelCheckIn, GueestMessage, ToDo, CheckOut byte
}

type TUserWalkInAccessOrder struct {
	ScheduleRate, Breakdown, ExtraCharge byte
}

type TUserFolioAccessOrder struct {
	Charge, Cash, Card, DirectBill, UpdateDirectBill, CashRefund, OtherPayment, Void, Correction, Transfer, UpdateSubDepartment, UpdateRemark, UpdateDocumentNumber, CheckOut, PrintFolio byte
}

type TUserFolioHistoryAccessOrder struct {
	Transaction, PrintFolio, CancelCheckOut byte
}

type TUserFloorPlanAccessOrder struct {
	Reception, HouseKeeping, ModifyFloorPlan byte
}

type TUserCompanyAccessOrder struct {
	Insert, Update, Delete, APLimit, ARLimit, DirectBill, BusinessSource byte
}

type TUserInvoiceAccessOrder struct {
	Insert, Update, Delete, InsertPayment, DeletePayment, PrintReceipt, ExchangeRate byte
}

type TUserMemberVoucherGift struct {
	RedeemPoint, InsertVoucher, DeleteVoucher, ApproveVoucher, SoldVoucher, RedeemVoucher, ComplimentVoucher byte
}

type TUserPreviewReportAccessOrder struct {
	EditReport, ExportReport, CustomizeReport byte
}

type TUserPaymentByAPARAccessOrder struct {
	PaymentByAPAR byte
}

type TOTPStatus struct {
	Active, Used, Expire, NotActive string
}

type TLogUserAction struct {
	//Reservation
	InsertReservation,
	VoidReservation,
	CheckInReservation,
	CancelReservation,
	NoShowReservation,
	InsertDeposit,
	VoidDeposit,
	CorrectDeposit,
	RefundDeposit,
	TransferDeposit,
	InsertReservationScheduledRate,
	UpdateReservationScheduledRate,
	DeleteReservationScheduledRate,
	InsertReservationExtraCharge,
	UpdateReservationExtraCharge,
	DeleteReservationExtraCharge,
	InsertReservationExtraChargeBreakdown,
	UpdateReservationExtraChargeBreakdown,
	DeleteReservationExtraChargeBreakdown,
	//Update Reservation Stay Information
	URSIArrival,
	URSINights,
	URSIDeparture,
	URSIAdult,
	URSIChild,
	URSIRoomType,
	URSIRoom,
	URSIRoomRate,
	URSIBusinessSource,
	URSICommissionType,
	URSICommissionValue,
	URSIWeekdayRate,
	URSIWeekendRate,
	URSIDiscount,
	URSIPaymentType,
	URSIMarket,
	URSIBillInstruction,
	URSICurrency,
	URSIExchangeRate,
	//Update Reservation Personal Information
	URPIMember,
	URPITitle,
	URPIFullName,
	URPIReservationBy,
	URPIStreet,
	URPICity,
	URPICountry,
	URPIState,
	URPIPostCode,
	URPIPhone1,
	URPIPhone2,
	URPIFax,
	URPIEmail,
	URPIWebsite,
	URPICompany,
	URPIGuestType,
	URPIIDCardType,
	URPIIDCardNumber,
	URPIBirthdayPlace,
	URPIBirthdate,
	//Update Reservation General Information
	URGIPurposeOf,
	URGIGroup,
	URGIDocumentNumber,
	URGIFlightNumber,
	URGIFlightArrival,
	URGIFlightDeparture,
	URGINotes,
	URGIHKNotes,
	URGIMarketing,
	URGITAVoucherNumber,
	//Group Reservation
	InsertGroupReservation,
	UpdateGroupReservation,
	DeleteGroupReservation,
	//Folio
	WalkIn,
	InsertMasterFolio,
	InsertDeskFolio,
	InsertFolioScheduledRate,
	UpdateFolioScheduledRate,
	DeleteFolioScheduledRate,
	InsertFolioExtraCharge,
	UpdateFolioExtraCharge,
	DeleteFolioExtraCharge,
	InsertFolioExtraChargeBreakdown,
	UpdateFolioExtraChargeBreakdown,
	DeleteFolioExtraChargeBreakdown,
	InsertTransaction,
	VoidTransaction,
	CorrectTransaction,
	TransferTransaction,
	RoutingFolio,
	ReturnTransfer,
	RemoveRouting,
	MoveRoom,
	SwitchRoom,
	InsertMessage,
	UpdateMessage,
	DeleteMessage,
	MarkAsDeliveredMessage,
	MarkAsUndeliveredMessage,
	InsertToDo,
	UpdateToDo,
	DeleteToDo,
	MarkAsDoneToDo,
	MarkAsNotDoneToDo,
	CancelCheckIn,
	CheckOutFolio,
	CancelCheckOut,
	ComplimentGuest,
	HouseUseGuest,
	DefaultGuest,
	UpdateFolioVoucher,
	InsertVoucherPayment,
	IssuedCard,
	ReplaceCard,
	EraseCard,
	DeactivateWithoutCard,
	ForceEraseCard,
	//Update Folio Stay Information
	UFSIArrival,
	UFSINights,
	UFSIDeparture,
	UFSIAdult,
	UFSIChild,
	UFSIRoomType,
	UFSIRoom,
	UFSIRoomRate,
	UFSIBusinessSource,
	UFSICommissionType,
	UFSICommissionValue,
	UFSIWeekdayRate,
	UFSIWeekendRate,
	UFSIDiscount,
	UFSIPaymentType,
	UFSIMarket,
	UFSIBillInstruction,
	UFSICurrency,
	UFSIExchangeRate,
	//Update Folio Personal Information
	UFPIMember,
	UFPITitle,
	UFPIFullName,
	UFPIStreet,
	UFPICity,
	UFPICountry,
	UFPIState,
	UFPIPostCode,
	UFPIPhone1,
	UFPIPhone2,
	UFPIFax,
	UFPIEmail,
	UFPIWebsite,
	UFPICompany,
	UFPIGuestType,
	UFPIIDCardType,
	UFPIIDCardNumber,
	UFPIBirthdayPlace,
	UFPIBirthdate,
	//Update Folio General Information
	UFGIPurposeOf,
	UFGIGroup,
	UFGIDocumentNumber,
	UFGIFlightNumber,
	UFGIFlightArrival,
	UFGIFlightDeparture,
	UFGINotes,
	UFGIHKNotes,
	UFGIMarketing,
	UFGITAVoucherNumber,
	//Member Voucher Gift
	MVInsertMember,
	MVUpdateMember,
	MVDeleteMember,
	MVRedeemMemberPoint,
	MVInsertVoucher,
	MVDeleteVoucher,
	MVApproveVoucher,
	MVNotApproveVoucher,
	MVUnapproveVoucher,
	MVSoldVoucher,
	MVRedeemVoucher,
	MVComplimentVoucher,
	MVUnsoldVoucher,
	//Sales Activity
	InsertLead,
	UpdateLead,
	VoidLead,
	InsertProposal,
	UpdateProposal,
	VoidProposal,
	InsertTask,
	UpdateTask,
	VoidTask,
	InsertSendReminder,
	UpdateSendReminder,
	DeleteSendReminder,
	InsertNotes,
	UpdateNotes,
	DeleteNotes,
	InsertActivityLog,
	UpdateActivityLog,
	DeleteActivityLog,
	InsertContact,
	UpdateContact,
	DeleteContact,

	//House Keeping
	ReadyRoom,
	CleanRoom,
	DirtyRoom,
	OutOfOrder,
	OfficeUse,
	UnderConstruction,
	DontDisturb,
	DoubleLock,
	SleepOut,
	PleaseClean,
	InsertRoomCosting,
	DeleteRoomCostingTransfer,
	DeleteRoomCosting,
	//Report
	InsertReportTemplate,
	UpdateReportTemplate,
	DeleteReportTemplate,
	SetDefaultReportTemplate,
	//Configuration
	CompanyInformation,
	Configuration,
	InsertMasterData,
	UpdateMasterData,
	DeleteMasterData,
	ModifyFloorPlan,
	//Accounting Tool
	InsertReceive,
	UpdateReceive,
	DeleteReceive,
	InsertPayment,
	UpdatePayment,
	DeletePayment,
	InsertReceipt,
	UpdateReceipt,
	DeleteReceipt,
	InsertIncomeBudget,
	UpdateIncomeBudget,
	DeleteIncomeBudget,
	InsertBudgetStatistic,
	UpdateBudgetStatistic,
	DeleteBudgetStatistic,
	//Account Payable
	InsertAPRefundDepositPayment,
	UpdateAPRefundDepositPayment,
	DeleteAPRefundDepositPayment,
	InsertAPCommissionPayment,
	UpdateAPCommissionPayment,
	DeleteAPCommissionPayment,
	//Account Receivable
	InsertInvoiceCityLedger,
	UpdateInvoiceCityLedger,
	DeleteInvoiceCityLedger,
	InsertPaymentInvoiceCityLedger,
	UpdatePaymentInvoiceCityLedger,
	DeletePaymentInvoiceCityLedger,
	InsertBankReconciliation,
	UpdateBankReconciliation,
	DeleteBankReconciliation,
	PrintInvoice,
	//User Setting
	InsertUser,
	UpdateUser,
	DeleteUser,
	InsertUserGroup,
	UpdateUserGroup,
	DeleteUserGroup,
	//Database
	BackupDatabase,
	RestoreDatabase,
	OptimizingDatabase,
	//Login
	Login,
	Logout,
	ChangePassword,
	LoginDenied int
	//Reservation
	InsertReservationX,
	VoidReservationX,
	CheckInReservationX,
	CancelReservationX,
	NoShowReservationX,
	InsertDepositX,
	VoidDepositX,
	CorrectDepositX,
	RefundDepositX,
	TransferDepositX,
	InsertReservationScheduledRateX,
	UpdateReservationScheduledRateX,
	DeleteReservationScheduledRateX,
	InsertReservationExtraChargeX,
	UpdateReservationExtraChargeX,
	DeleteReservationExtraChargeX,
	InsertReservationExtraChargeBreakdownX,
	UpdateReservationExtraChargeBreakdownX,
	DeleteReservationExtraChargeBreakdownX,
	//Update Reservation Stay Information
	URSIArrivalX,
	URSINightsX,
	URSIDepartureX,
	URSIAdultX,
	URSIChildX,
	URSIRoomTypeX,
	URSIRoomX,
	URSIRoomRateX,
	URSIBusinessSourceX,
	URSICommissionTypeX,
	URSICommissionValueX,
	URSIWeekdayRateX,
	URSIWeekendRateX,
	URSIDiscountX,
	URSIPaymentTypeX,
	URSIMarketX,
	URSIBillInstructionX,
	URSICurrencyX,
	URSIExchangeRateX,
	//Update Reservation Personal Information
	URPIMemberX,
	URPITitleX,
	URPIFullNameX,
	URPIReservationByX,
	URPIStreetX,
	URPICityX,
	URPICountryX,
	URPIStateX,
	URPIPostCodeX,
	URPIPhone1X,
	URPIPhone2X,
	URPIFaxX,
	URPIEmailX,
	URPIWebsiteX,
	URPICompanyX,
	URPIGuestTypeX,
	URPIIDCardTypeX,
	URPIIDCardNumberX,
	URPIBirthdayPlaceX,
	URPIBirthdateX,
	//Update Reservation General Information
	URGIPurposeOfX,
	URGIGroupX,
	URGIDocumentNumberX,
	URGIFlightNumberX,
	URGIFlightArrivalX,
	URGIFlightDepartureX,
	URGINotesX,
	URGIHKNotesX,
	URGIMarketingX,
	URGITAVoucherNumberX,
	//Group Reservation
	InsertGroupReservationX,
	UpdateGroupReservationX,
	DeleteGroupReservationX,
	//Folio
	WalkInX,
	InsertMasterFolioX,
	InsertDeskFolioX,
	InsertFolioScheduledRateX,
	UpdateFolioScheduledRateX,
	DeleteFolioScheduledRateX,
	InsertFolioExtraChargeX,
	UpdateFolioExtraChargeX,
	DeleteFolioExtraChargeX,
	InsertFolioExtraChargeBreakdownX,
	UpdateFolioExtraChargeBreakdownX,
	DeleteFolioExtraChargeBreakdownX,
	InsertTransactionX,
	VoidTransactionX,
	CorrectTransactionX,
	TransferTransactionX,
	RoutingFolioX,
	ReturnTransferX,
	RemoveRoutingX,
	MoveRoomX,
	SwitchRoomX,
	InsertMessageX,
	UpdateMessageX,
	DeleteMessageX,
	MarkAsDeliveredMessageX,
	MarkAsUndeliveredMessageX,
	InsertToDoX,
	UpdateToDoX,
	DeleteToDoX,
	MarkAsDoneToDoX,
	MarkAsNotDoneToDoX,
	CancelCheckInX,
	CheckOutFolioX,
	CancelCheckOutX,
	ComplimentGuestX,
	HouseUseGuestX,
	DefaultGuestX,
	UpdateFolioVoucherX,
	InsertVoucherPaymentX,
	IssuedCardX,
	ReplaceCardX,
	EraseCardX,
	DeactivateWithoutCardX,
	ForceEraseCardX,
	//Update Folio Stay Information
	UFSIArrivalX,
	UFSINightsX,
	UFSIDepartureX,
	UFSIAdultX,
	UFSIChildX,
	UFSIRoomTypeX,
	UFSIRoomX,
	UFSIRoomRateX,
	UFSIBusinessSourceX,
	UFSICommissionTypeX,
	UFSICommissionValueX,
	UFSIWeekdayRateX,
	UFSIWeekendRateX,
	UFSIDiscountX,
	UFSIPaymentTypeX,
	UFSIMarketX,
	UFSIBillInstructionX,
	UFSICurrencyX,
	UFSIExchangeRateX,
	//Update Folio Personal Information
	UFPIMemberX,
	UFPITitleX,
	UFPIFullNameX,
	UFPIStreetX,
	UFPICityX,
	UFPICountryX,
	UFPIStateX,
	UFPIPostCodeX,
	UFPIPhone1X,
	UFPIPhone2X,
	UFPIFaxX,
	UFPIEmailX,
	UFPIWebsiteX,
	UFPICompanyX,
	UFPIGuestTypeX,
	UFPIIDCardTypeX,
	UFPIIDCardNumberX,
	UFPIBirthdayPlaceX,
	UFPIBirthdateX,
	//Update Folio General Information
	UFGIPurposeOfX,
	UFGIGroupX,
	UFGIDocumentNumberX,
	UFGIFlightNumberX,
	UFGIFlightArrivalX,
	UFGIFlightDepartureX,
	UFGINotesX,
	UFGIHKNotesX,
	UFGIMarketingX,
	UFGITAVoucherNumberX,
	//Member Voucher Gift
	MVInsertMemberX,
	MVUpdateMemberX,
	MVDeleteMemberX,
	MVRedeemMemberPointX,
	MVInsertVoucherX,
	MVDeleteVoucherX,
	MVApproveVoucherX,
	MVNotApproveVoucherX,
	MVUnapproveVoucherX,
	MVSoldVoucherX,
	MVRedeemVoucherX,
	MVComplimentVoucherX,
	MVUnsoldVoucherX,

	//Sales Activity
	InsertLeadX,
	UpdateLeadX,
	VoidLeadX,
	InsertProposalX,
	UpdateProposalX,
	VoidProposalX,
	InsertTaskX,
	UpdateTaskX,
	VoidTaskX,
	InsertSendReminderX,
	UpdateSendReminderX,
	DeleteSendReminderX,
	InsertNotesX,
	UpdateNotesX,
	DeleteNotesX,
	InsertActivityLogX,
	UpdateActivityLogX,
	DeleteActivityLogX,
	InsertContactX,
	UpdateContactX,
	DeleteContactX,

	//House Keeping
	ReadyRoomX,
	CleanRoomX,
	DirtyRoomX,
	OutOfOrderX,
	OfficeUseX,
	UnderConstructionX,
	DontDisturbX,
	DoubleLockX,
	SleepOutX,
	PleaseCleanX,
	InsertRoomCostingX,
	DeleteRoomCostingTransferX,
	DeleteRoomCostingX,
	//Report
	InsertReportTemplateX,
	UpdateReportTemplateX,
	DeleteReportTemplateX,
	SetDefaultReportTemplateX,
	//Configuration
	CompanyInformationX,
	ConfigurationX,
	InsertMasterDataX,
	UpdateMasterDataX,
	DeleteMasterDataX,
	ModifyFloorPlanX,
	//Accounting Tool
	InsertReceiptX,
	UpdateReceiptX,
	DeleteReceiptX,
	InsertIncomeBudgetX,
	UpdateIncomeBudgetX,
	DeleteIncomeBudgetX,
	InsertBudgetStatisticX,
	UpdateBudgetStatisticX,
	DeleteBudgetStatisticX,
	//Account Payable
	InsertAPRefundDepositPaymentX,
	UpdateAPRefundDepositPaymentX,
	DeleteAPRefundDepositPaymentX,
	InsertAPCommissionPaymentX,
	UpdateAPCommissionPaymentX,
	DeleteAPCommissionPaymentX,
	//Account Receivable
	InsertInvoiceCityLedgerX,
	UpdateInvoiceCityLedgerX,
	DeleteInvoiceCityLedgerX,
	InsertPaymentInvoiceCityLedgerX,
	UpdatePaymentInvoiceCityLedgerX,
	DeletePaymentInvoiceCityLedgerX,
	InsertBankReconciliationX,
	UpdateBankReconciliationX,
	DeleteBankReconciliationX,
	PrintInvoiceX,
	//User Setting
	InsertUserX,
	UpdateUserX,
	DeleteUserX,
	InsertUserGroupX,
	UpdateUserGroupX,
	DeleteUserGroupX,
	//Database
	BackupDatabaseX,
	RestoreDatabaseX,
	OptimizingDatabaseX,
	//Login
	LoginX,
	LogoutX,
	ChangePasswordX,
	LoginDeniedX bool
}
type TLogUserActionPOS struct {
	//Customer
	InsertCustomer,
	UpdateCustomer,
	DeleteCustomer,
	//Reservation
	InsertReservationPOS,
	UpdateReservationPOS,
	ChekInReservationPOS,
	CancelReservationPOS,
	NoShowReservationPOS,
	VoidReservationPOS,
	//POS Terminal and Table View
	InsertCaptainOrder,
	UpdateCaptainOrder,
	TransferCaptainOrder,
	CancelCaptainOrder,
	ChangeQuantity,
	UpdateRemark,
	Discount,
	OverridePrice,
	ModifyPriceZero,
	ModifyPriceRemoveTaxAndService,
	RemoveItem,
	FinishSale,
	ModifyTableView,
	VoidCheck,
	//Member Voucher Gift
	MVInsertMember,
	MVUpdateMember,
	MVDeleteMember,
	MVRedeemMemberPoint,
	MVInsertVoucher,
	MVDeleteVoucher,
	MVApproveVoucher,
	MVNotApproveVoucher,
	MVUnapproveVoucher,
	MVSoldVoucher,
	MVRedeemVoucher,
	MVComplimentVoucher,
	MVUnsoldVoucher,
	//Budget
	InsertFBBudget,
	UpdateFBBudget,
	DeleteFBBudget,
	//Report
	InsertReportTemplate,
	UpdateReportTemplate,
	DeleteReportTemplate,
	SetDefaultReportTemplate,
	//Configuration
	CompanyInformation,
	Configuration,
	InsertMasterData,
	UpdateMasterData,
	DeleteMasterData,
	//User Setting
	InsertUser,
	UpdateUser,
	DeleteUser,
	InsertUserGroup,
	UpdateUserGroup,
	DeleteUserGroup,
	//Database
	BackupDatabase,
	RestoreDatabase,
	OptimizingDatabase,
	//Login
	Login,
	Logout,
	ChangePassword,
	LoginDenied,
	//CHS Parameter Dummy
	CancelReservation,
	NoShowReservation,
	CleanRoom,
	DirtyRoom int
	//POS Terminal and Table View
	InsertCaptainOrderX,
	UpdateCaptainOrderX,
	TransferCaptainOrderX,
	CancelCaptainOrderX,
	ChangeQuantityX,
	UpdateRemarkX,
	DiscountX,
	OverridePriceX,
	ModifyPriceZeroX,
	ModifyPriceRemoveTaxAndServiceX,
	RemoveItemX,
	FinishSaleX,
	ModifyTableViewX,
	VoidCheckX,
	//Member Voucher Gift
	MVInsertMemberX,
	MVUpdateMemberX,
	MVDeleteMemberX,
	MVRedeemMemberPointX,
	MVInsertVoucherX,
	MVDeleteVoucherX,
	MVApproveVoucherX,
	MVNotApproveVoucherX,
	MVUnapproveVoucherX,
	MVSoldVoucherX,
	MVRedeemVoucherX,
	MVComplimentVoucherX,
	MVUnsoldVoucherX,
	//Budget
	InsertFBBudgetX,
	UpdateFBBudgetX,
	DeleteFBBudgetX,
	//Report
	InsertReportTemplateX,
	UpdateReportTemplateX,
	DeleteReportTemplateX,
	SetDefaultReportTemplateX,
	//Configuration
	CompanyInformationX,
	ConfigurationX,
	InsertMasterDataX,
	UpdateMasterDataX,
	DeleteMasterDataX,
	//User Setting
	InsertUserX,
	UpdateUserX,
	DeleteUserX,
	InsertUserGroupX,
	UpdateUserGroupX,
	DeleteUserGroupX,
	//Database
	BackupDatabaseX,
	RestoreDatabaseX,
	OptimizingDatabaseX,
	//Login
	LoginX,
	LogoutX,
	ChangePasswordX,
	LoginDeniedX,
	//CHS Parameter Dummy
	CancelReservationX,
	NoShowReservationX,
	CleanRoomX,
	DirtyRoomX bool
}

type TLogUserActionCAMS struct {
	//Inventory
	InsertPurchaseRequest,
	UpdatePurchaseRequest,
	DeletePurchaseRequest,
	InsertPurchaseOrder,
	UpdatePurchaseOrder,
	DeletePurchaseOrder,
	InsertReceiveStock,
	UpdateReceiveStock,
	DeleteReceiveStock,
	InsertStockTransfer,
	UpdateStockTransfer,
	DeleteStockTransfer,
	InsertCosting,
	UpdateCosting,
	DeleteCosting,
	InsertProduction,
	UpdateProduction,
	DeleteProduction,
	InsertCostRecipe,
	UpdateCostRecipe,
	DeleteCostRecipe,
	InsertReturnStock,
	UpdateReturnStock,
	DeleteReturnStock,
	InsertStockOpname,
	DeleteStockOpname,
	SetActiveStore,
	//Fixed Asset
	InsertFAPurchaseOrder,
	UpdateFAPurchaseOrder,
	DeleteFAPurchaseOrder,
	InsertFAReceiveStock,
	UpdateFAReceiveStock,
	DeleteFAReceiveStock,
	InsertFixedAssetList,
	UpdateFixedAssetList,
	DeleteFixedAssetList,
	InsertDepreciation,
	DeleteDepreciation,
	//Report
	InsertReportTemplate,
	UpdateReportTemplate,
	DeleteReportTemplate,
	SetDefaultReportTemplate,
	//Configuration
	CompanyInformation,
	Configuration,
	InsertMasterData,
	UpdateMasterData,
	DeleteMasterData,
	//User Setting
	InsertUser,
	UpdateUser,
	DeleteUser,
	InsertUserGroup,
	UpdateUserGroup,
	DeleteUserGroup,
	//Database
	BackupDatabase,
	RestoreDatabase,
	OptimizingDatabase,
	//Login
	Login,
	Logout,
	ChangePassword,
	LoginDenied,
	//Dummy
	CancelCheckOut int
	//Inventory
	InsertPurchaseRequestX,
	UpdatePurchaseRequestX,
	DeletePurchaseRequestX,
	InsertPurchaseOrderX,
	UpdatePurchaseOrderX,
	DeletePurchaseOrderX,
	InsertReceiveStockX,
	UpdateReceiveStockX,
	DeleteReceiveStockX,
	InsertStockTransferX,
	UpdateStockTransferX,
	DeleteStockTransferX,
	InsertCostingX,
	UpdateCostingX,
	DeleteCostingX,
	InsertProductionX,
	UpdateProductionX,
	DeleteProductionX,
	InsertCostRecipeX,
	UpdateCostRecipeX,
	DeleteCostRecipeX,
	InsertReturnStockX,
	UpdateReturnStockX,
	DeleteReturnStockX,
	InsertStockOpnameX,
	DeleteStockOpnameX,
	SetActiveStoreX,
	//Fixed Asset
	InsertFAPurchaseOrderX,
	UpdateFAPurchaseOrderX,
	DeleteFAPurchaseOrderX,
	InsertFAReceiveStockX,
	UpdateFAReceiveStockX,
	DeleteFAReceiveStockX,
	InsertFixedAssetListX,
	UpdateFixedAssetListX,
	DeleteFixedAssetListX,
	InsertDepreciationX,
	DeleteDepreciationX,
	//Report
	InsertReportTemplateX,
	UpdateReportTemplateX,
	DeleteReportTemplateX,
	SetDefaultReportTemplateX,
	//Configuration
	CompanyInformationX,
	ConfigurationX,
	InsertMasterDataX,
	UpdateMasterDataX,
	DeleteMasterDataX,
	ModifyFloorPlanX,
	//User Setting
	InsertUserX,
	UpdateUserX,
	DeleteUserX,
	InsertUserGroupX,
	UpdateUserGroupX,
	DeleteUserGroupX,
	//Database
	BackupDatabaseX,
	RestoreDatabaseX,
	OptimizingDatabaseX,
	//Login
	LoginX,
	LogoutX,
	ChangePasswordX,
	LoginDeniedX,
	//Dummy
	CancelCheckOutX bool
}
type TLogUserActionCAS struct {
	//Accounting Tool
	InsertReceive,
	UpdateReceive,
	DeleteReceive,
	InsertPayment,
	UpdatePayment,
	DeletePayment,
	InsertReceipt,
	UpdateReceipt,
	DeleteReceipt,
	InsertJournal,
	UpdateJournal,
	DeleteJournal,
	InsertIncomeBudget,
	UpdateIncomeBudget,
	DeleteIncomeBudget,
	InsertExpenseBudget,
	UpdateExpenseBudget,
	DeleteExpenseBudget,
	InsertBudgetStatistic,
	UpdateBudgetStatistic,
	DeleteBudgetStatistic,
	CloseMonth,
	CancelCloseMonth,
	CloseYear,
	CancelCloseYear,
	//Account Payable
	InsertAccountPayable,
	UpdateAccountPayable,
	DeleteAccountPayable,
	InsertAccountPayablePayment,
	UpdateAccountPayablePayment,
	DeleteAccountPayablePayment,
	InsertAPRefundDepositPayment,
	UpdateAPRefundDepositPayment,
	DeleteAPRefundDepositPayment,
	InsertAPCommissionPayment,
	UpdateAPCommissionPayment,
	DeleteAPCommissionPayment,
	//Account Receivable
	InsertAccountReceivable,
	UpdateAccountReceivable,
	DeleteAccountReceivable,
	InsertAccountReceivablePayment,
	UpdateAccountReceivablePayment,
	DeleteAccountReceivablePayment,
	InsertInvoiceCityLedger,
	UpdateInvoiceCityLedger,
	DeleteInvoiceCityLedger,
	InsertPaymentInvoiceCityLedger,
	UpdatePaymentInvoiceCityLedger,
	DeletePaymentInvoiceCityLedger,
	InsertBankReconciliation,
	UpdateBankReconciliation,
	DeleteBankReconciliation,
	PrintInvoice,
	//Report
	InsertReportTemplate,
	UpdateReportTemplate,
	DeleteReportTemplate,
	SetDefaultReportTemplate,
	//Configuration
	CompanyInformation,
	Configuration,
	InsertMasterData,
	UpdateMasterData,
	DeleteMasterData,
	//User Setting
	InsertUser,
	UpdateUser,
	DeleteUser,
	InsertUserGroup,
	UpdateUserGroup,
	DeleteUserGroup,
	//Database
	BackupDatabase,
	RestoreDatabase,
	OptimizingDatabase,
	//Login
	Login,
	Logout,
	ChangePassword,
	LoginDenied,
	//Dummy
	CancelCheckOut int
	//Accounting Tool
	InsertReceiveX,
	UpdateReceiveX,
	DeleteReceiveX,
	InsertPaymentX,
	UpdatePaymentX,
	DeletePaymentX,
	InsertReceiptX,
	UpdateReceiptX,
	DeleteReceiptX,
	InsertJournalX,
	UpdateJournalX,
	DeleteJournalX,
	InsertIncomeBudgetX,
	UpdateIncomeBudgetX,
	DeleteIncomeBudgetX,
	InsertExpenseBudgetX,
	UpdateExpenseBudgetX,
	DeleteExpenseBudgetX,
	InsertBudgetStatisticX,
	UpdateBudgetStatisticX,
	DeleteBudgetStatisticX,
	CloseMonthX,
	CancelCloseMonthX,
	CloseYearX,
	CancelCloseYearX,
	//Account Payable
	InsertAccountPayableX,
	UpdateAccountPayableX,
	DeleteAccountPayableX,
	InsertAccountPayablePaymentX,
	UpdateAccountPayablePaymentX,
	DeleteAccountPayablePaymentX,
	InsertAPRefundDepositPaymentX,
	UpdateAPRefundDepositPaymentX,
	DeleteAPRefundDepositPaymentX,
	InsertAPCommissionPaymentX,
	UpdateAPCommissionPaymentX,
	DeleteAPCommissionPaymentX,
	//Account Receivable
	InsertAccountReceivableX,
	UpdateAccountReceivableX,
	DeleteAccountReceivableX,
	InsertAccountReceivablePaymentX,
	UpdateAccountReceivablePaymentX,
	DeleteAccountReceivablePaymentX,
	InsertInvoiceCityLedgerX,
	UpdateInvoiceCityLedgerX,
	DeleteInvoiceCityLedgerX,
	InsertPaymentInvoiceCityLedgerX,
	UpdatePaymentInvoiceCityLedgerX,
	DeletePaymentInvoiceCityLedgerX,
	InsertBankReconciliationX,
	UpdateBankReconciliationX,
	DeleteBankReconciliationX,
	PrintInvoiceX,
	//Report
	InsertReportTemplateX,
	UpdateReportTemplateX,
	DeleteReportTemplateX,
	SetDefaultReportTemplateX,
	//Configuration
	CompanyInformationX,
	ConfigurationX,
	InsertMasterDataX,
	UpdateMasterDataX,
	DeleteMasterDataX,
	//User Setting
	InsertUserX,
	UpdateUserX,
	DeleteUserX,
	InsertUserGroupX,
	UpdateUserGroupX,
	DeleteUserGroupX,
	//Database
	BackupDatabaseX,
	RestoreDatabaseX,
	OptimizingDatabaseX,
	//Login
	LoginX,
	LogoutX,
	ChangePasswordX,
	LoginDeniedX,
	//Dummy
	CancelCheckOutX bool
}

type TTransactionType struct {
	Debit, Credit string
}

type TFolioType struct {
	GuestFolio, MasterFolio, DeskFolio string
}

type TGlobalAccountGroup struct {
	Charge, Payment, TaxService, Transfer string
}

type TGlobalAccountSubGroup struct {
	RoomCharge, AccountPayable, Payment, CreditDebitCard, BankTransfer, AccountReceivable, Compliment string
}

type TGlobalAccount struct {
	Breakfast, RoomCharge, ExtraBed, CancellationFee, NoShow, Telephone, APRefundDeposit, APCommission, CreditCardAdm,
	Cash, CityLedger, Voucher, VoucherCompliment, Tax, Service, TransferDepositReservation, TransferDepositReservationToFolio, TransferCharge, TransferPayment string
}

type TGlobalDepartment struct {
	FoodBeverage, RoomDivision, Banquet, Minor, Miscellaneous string
}

type TGlobalSubDepartment struct {
	FrontOffice, HouseKeeping, Banquet, Accounting string
}

type TGlobalPaymentType struct {
	TransferDeposit string
}

type TWeekendDay struct {
	Friday, Saturday, Sunday                                                                     bool
	FridayFromTime, FridayToTime, SaturdayFromTime, SaturdayToTime, SundayFromTime, SundayToTime time.Time
}

type TRoomStatus struct {
	Ready, Clean, Dirty, Vacant, Occupied, HouseUseX, Compliment, DontDisturb, DoubleLock, SleepOut, PleaseClean, OutOfOrder, OfficeUse, UnderConstruction string
}

type TRoomBlockStatus struct {
	GeneralCleaning, ShowingRoom string
}

type TRoomStatusColor struct {
	Reserved, Occupied, HouseUse, Compliment, OutOfOrder, OfficeUse, UnderConstruction, Available int
}

type TReservationStatus struct {
	New, WaitList, InHouse, Canceled, NoShow, Void, CheckOut string
}

type TReservationBlockType struct {
	Reservation, BlockOnly, ChekIn string
}

type TReservationStatus2 struct {
	Tentative, Confirm, Guaranteed string
}

type TReservationType struct {
	Guaranteed, NonGuaranteed string
}

type TFolioStatus struct {
	Open, Closed, CancelCheckIn string
}

type TFolioTransferBy struct {
	NoTransfer, ByAccount, ByAccountSubGroup string
}

type TSubFolioGroup struct {
	A, B, C, D string
}

type TSubFolioGroupName struct {
	A, B, C, D string
}

type ReportAccessForm struct {
	CashierReport, FrontDeskReport, PointOfSalesReport, BanquetReport, AccountingReport, InventoryAssetReport int
}

type TReportAccessOrder struct {
	AccessForm ReportAccessForm
}

type TSubFolioPostingType struct {
	None, Deposit, Transfer, Room, ExtraCharge string
}

type TCommissionType struct {
	PercentFirstNightFullRate, PercentPerNightFullRate, PercentFirstNightNettRate, PercentPerNightNettRate, FixAmountFirstNight, FixAmountPerNight, PercentOfPriceFullPrice, PercentOfPriceNettPrice string
}

type TCPType struct {
	Guest, Company, Invoice, RoomOwner string
}

type TChargeFrequency struct {
	OnceOnly, Daily, Weekly, Monthly string
}

type TKeylockVendor struct {
	Ventaza, UltraLock, DLock, DLockOldVersion, Kiara, Colcom, SureLock, Deluns, Kima, HuneLock, VingCard, SunVitio, Rafles, BeTech, OnityOld, Saflok, ColcomOld, Tesa, VingcardVisionline, CLock, Downs, Onity, OnityOldVersion, VingCardSerial, Kaba, RaflesUSB string
}

type TChannelManager struct {
	BookNLink, Stah, SiteMinder string
}

type TMikrotikVendor struct {
	Mikrotik, Megalos, Provenos, Coova string
}

type TPMSCommand struct {
	Issued, Verified, Replaced, Erase, ForceErase string
}

// Accounting
type TJournalPrefix struct {
	Manual, Transaction, Disbursement, Receive, Inventory, Adjustment, FixedAsset, BeginningYear, AccountPayable, AccountReceivable string
}

type TJournalType struct {
	CashIn, CashOut, CashTransfer, BankIn, BankOut, BankTransfer, CreditCardReconciliation, Cheque, Other string
}

type TJournalGroup struct {
	CapitalReceipt, Receiving, ARPayment, Loan, Investment, BankWithdrawal,
	AssetPurchasing, InventoryPurchasing, OtherPurchasing, APPayment, Expense, BankDeposit, CashTransfer, BankTransfer,
	AccountPayable, AccountReceivable, Costing, FixedAsset, Other string
}

type TGlobalJournalAccountType struct {
	Bank, AccountReceivable, CreditCard, OtherCurrentAsset, FixAsset, OtherAsset, AccountPayable, OtherCurrentLiability, LongTermLiability, OtherLiability, Equity, Income, Cost, Expense, OtherIncome, OtherExpense string
}

type TGlobalJournalAccountGroup struct {
	Assets, Liability, Equity, Income, Cost, Expense1, Expense2, OtherIncome, OtherExpense string
}

type TGlobalJournalAccount struct {
	APSupplier, OverShortAsIncome, OverShortAsExpense, ServiceRevenue, IncomePurchasingDiscount, ExpensePurchasingTax, ExpenseShipping, IncomeReturnStock, ExpenseReturnStock                            string
	GuestLedger, APVoucher, GuestDeposit, GuestDepositReservation, ProfitLossBeginningYear, ProfitLossCurrentYear, ProfitLossCurrency, IncomeVoucherExpire, ExpenseInvoiceDiscount, ExpenseCreditCardAdm string
}

type TBankAccountType struct {
	CashAccount, SavingAccount, CreditAccount, ChequeAccount string
}

type TBudgetType struct {
	Manual, Average, Percentage, Daily string
}

type TAmountPreset struct {
	Preset1, Preset2, Preset3, Preset4, Preset5, Preset6, Preset7, Preset8 int
}

type TForeignCashTableID struct {
	GuestDeposit, SubFolio, InvoicePayment int
}

type TReportSignature struct {
	ShowPreparedBy, ShowCheckedBy, ShowApprovedBy, ShowInDailySalesReport, ShowInDailyRevenueReport, ShowInRoomStatisticReport bool
	PreparedBy, CheckedBy, ApprovedBy                                                                                          string
}

type TReportTemplate struct {
	DepositReceipt, DepositreceiptRefund, FolioReceipt, FolioReceiptRefund, MiscellaneousCharge, CashierReport, CashRemittance, BreakfastList, PickUpService, DropService, HKCheckList, HKRoomStatus, HKRoomStatusSummary, HKRoomDiscrepancy, HKRoomAttendantControlSheet, DailyRevenueReport, DailyRevenueReportSummary, GuestInHouseListing, RoomStatistic, DailyFlashReport, RoomProduction, GuestForecast, RevenueBySubDepartment string
}

type TPersonal struct {
	DirectorName, DirectorPosition, HotelManagerName, HotelManagerPosition, FinanceHead, FinanceHeadPosition, IncomeAudit, IncomeAuditPosition, GeneralCashier, GeneralCashierPosition,
	AccountReceivableName, AccountReceivablePosition, AccountPayableName, AccountPayablePosition string
}

type TStatisticAccount struct {
	TotalRoom, OutOfOrder, OfficeUse, UnderConstruction, HouseUse, Compliment, RoomSold, DayUse, WalkIn, CheckInByReservation, NoShow, ReservationMade, CancelationReservation, EarlyCheckOut, CheckOut, RevenueGross, RevenueNonPackage, RevenueNett, RevenueWithCompliment, Adult, Child, AdultSold, ChildSold, AdultDayUse, ChildDayUse, AdultCompliment, ChildCompliment, AdultHouseUse, ChildHouseUse, NumberOfCover, NettFoodSales, NettBeverageSales, NettBreakfastSales, NettBanquetSales, NettWeddingSales, NettGatheringSales, BreakfastCover, BeverageCover, FoodCover, BanquetCover, WeddingCover, GatheringCover string
}

type TPaymentGroup struct {
	Cash, Bank, DirectBill, Other, None string
}

type TInventoryItemGroup struct {
	Food, Beverage string
}

type TDynamicRateType struct {
	None, BaseOccupancy, BaseSession, BaseScale, BaseCompetitor, BaseWeekly string
}

type TMemberType struct {
	Room, Outlet, Banquet string
}

type TVoucherType struct {
	Compliment, Discount, Sale string
}

type TVoucherStatus struct {
	Active, Used, Expire string
}

type TVoucherStatusApprove struct {
	Unapprove, Approved, NotApproved string
}

type TVoucherStatusSold struct {
	Sold, Redeemed, Compliment string
}

type TSMSEvent struct {
	OnInsertReservation, OnWalkIn, OnCheckIn, OnCheckOut, OnDayendCloseFinish int
}

type TComplimentType struct {
	None, Compliment, OfficerCheck, EntertainCheck string
}

type TCaptainOrderType struct {
	DineIn, TakeAway, Delivery, Reservation string
}

type TSpecialProduct struct {
	ProductCategoryCode, FoodCode, BeverageCode string
}

type TGuestProfileSource struct {
	Hotel, Pos, Banquet string
}

type TCustomField struct {
	CFName01, CFName02, CFName03, CFName04, CFName05, CFName06, CFName07, CFName08, CFName09, CFName10, CFName11, CFName12,
	CFDefaultValue01, CFDefaultValue02, CFDefaultValue03, CFDefaultValue04, CFDefaultValue05, CFDefaultValue06, CFDefaultValue07, CFDefaultValue08, CFDefaultValue09, CFDefaultValue10, CFDefaultValue11, CFDefaultValue12 string
}

type TCustomLookupField struct {
	CLFName01, CLFName02, CLFName03, CLFName04, CLFName05, CLFName06, CLFName07, CLFName08, CLFName09, CLFName10, CLFName11, CLFName12,
	CLFDefaultValue01, CLFDefaultValue02, CLFDefaultValue03, CLFDefaultValue04, CLFDefaultValue05, CLFDefaultValue06, CLFDefaultValue07, CLFDefaultValue08, CLFDefaultValue09, CLFDefaultValue10, CLFDefaultValue11, CLFDefaultValue12 string
}

type TCMUpdateType struct {
	Reservation, Folio, RoomAllotment, Rate, Availability string
}

type TInventoryCostingMethod struct {
	FIFO, LIFO, Average string
}

type TSalesActivityStatus struct {
	Deal, Finished, New, Qualified, ProposalSend, Working, Void string
}

type TSalesActivityProposalStatus struct {
	Draft, Send, Revised, Contacted, Declined, Accepted, Void string
}

type TSalesActivityTaskPriority struct {
	Low, Medium, High, Urgent string
}

type TSalesActivityTaskStatus struct {
	NotStarted, InProgress, Testing, AwaitingFeedback, Complete, Void string
}

type TNotifThirdPartyTemplateID struct {
	EmailReminderreservation, EmailOnCheckIn, EmailOnCheckOut, EmailOnGusetBirthday, EmailSalesActivitySendRemider,
	WAReminderreservation, WAOnCheckIn, WAOnCheckOut, WAOnGusetBirthday, WASalesActivitySendRemider int
}

type TNotifThirdPartySourceCode struct {
	Reservation, Folio, SalesActivitySendRemider, Other string
}

type TPaymentSourceAPAR struct {
	None, AccountPayable, AccountReceivable, CityLedger, APCommission string
}

type TResponseCode struct {
	Successfully, InternalServerError, NotAuthorized, Unregistered, InvalidDataFormat, SubscriptionExpired, ErrorCreateToken, DataNotFound, InvalidDataValue, DatabaseValueChanged, DatabaseError, EmptyRequireField,
	PaymentDateCannotLowerThanIssuedDate, SuccessfullyWithStatus, DuplicateEntry, OtherResult uint
}

type TPaginationResponds struct {
	List       interface{}
	Page       int
	Limit      int
	TotalPages int
	TotalCount int64
	HasMore    bool
}

type TResponseText struct {
	Successfully, NotAuthorized, InvalidDataFormat, ErrorCreateToken, DataNotFound, InvalidDataValue, DatabaseValueChanged, DatabaseError, EmptyRequireField,
	PaymentDateCannotLowerThanIssuedDate, SuccessfullyWithStatus, DuplicateEntry, OtherResult string
}
type TDepartmentType struct {
	Income, OtherIncome, NonIncome string
}

type TJournalAccountSubGroupType struct {
	PayrollRelated, OtherExpense, EnergyCost, SalesEnergy string
}
type TPurchaseRequestStatus struct {
	NotApproved, Approved, Rejected string
}

type TRequestResponse struct {
	StatusCode uint
	Message    interface{}
	Result     interface{}
}

type TWSMessageType struct {
	Broadcast  int
	Client     int
	Channel    int
	Room       int
	Connection int
}

type TWSDataType struct {
	ServerStatus                   int
	DayendCloseStatus              int
	ModifiedRoomAvailabilityStatus int
	AuditDateChanged               int
}

type TDataset struct {
	// ProgramVariable      global_var.TProgramVariable
	Configuration        map[string]map[string]interface{}
	ProgramConfiguration TProgramConfiguration
	SpecialProduct       TSpecialProduct
	// VariableDLL: TVariableDLL;
	// TempVariable: TTempVariable;
	DefaultVariable TDefaultVariable
	// ProgramMessage: TProgramMessage;
	// DateCheckMessage: array [0..1] of PChar;
	// UserInfo: TUserInfo;
	GlobalAccount                 TGlobalAccount
	GlobalDepartment              TGlobalDepartment
	GlobalSubDepartment           TGlobalSubDepartment
	GlobalJournalAccount          TGlobalJournalAccount
	GlobalJournalAccountSubGroup  TGlobalJournalAccountSubGroup
	WeekendDay                    TWeekendDay
	GlobalJournalAccountGroupName TGlobalJournalAccountGroupName
}

const PasswordKeyString = "^%^#@JHGHFsd56hs123^93g$0"

var Tracer = otel.Tracer("PMS Service")
var PublicPath string
var Config *config.Config
var Mutex = &sync.RWMutex{}
var MxSocket = &sync.RWMutex{}
var localTime, _ = time.LoadLocation("Asia/Makassar")
var AppLogger *otelzap.Logger

var SigningKey = []byte("H6$#%123hjsdf(765398$$")
var EncryptKey = []byte("bas5%$#8123jhsjHFHjs7426%54238$#")
var EncryptKeyGlobal = []byte("JHGD123DJH&^*^&SDJJHSDGHGR%@#^@%&@&")

var AESSecretKey = []byte("^%^#@JHGHFsd5123d63^93g$0gtr#vf3")

var AESiv = []byte("JHGHFsd5123d63^9")

var DatabaseInfo databaseinfo

var Db *gorm.DB
var DBMain *gorm.DB
var EmptyFloat64 *float64
var EmptyFloat32 *float32
var EmptyUint64 *uint64
var EmptyUint32 *uint32
var EmptyUint *uint
var EmptyUint8 *uint8
var EmptyString *string
var EmptyBool *bool
var EmptyInt *int

var SQLParameterName = TSQLParameterName{
	Param01: ":Param01",
	Param02: ":Param02",
	Param03: ":Param03",
	Param04: ":Param04",
	Param05: ":Param05",
	Param06: ":Param06",
	Param07: ":Param07",
	Param08: ":Param08",
	Param09: ":Param09",
	Param10: ":Param10"}
var EditorFormCaption = []string{"Insert ", "Update ", "Duplicate "}
var ConstProgramVariable = TConstProgramVariable{
	DefaultSystemCode:         "H",
	DateTimeFormatYYYYMMDD:    "2006-01-02",
	DatabaseVersion:           "07.05.00.00",
	SettingFileName:           "Setting.ini",
	WallpaperListFile:         "WallpaperH.txt",
	TempImageFile:             "TempImage.jpg",
	InvoiceNumberPrefix:       "INV",
	APNumberPrefix:            "AP",
	ARNumberPrefix:            "AR",
	PaymentNumberPrefix:       "PY",
	ReceiveNumberPrefix:       "RC",
	StockTransferNumberPrefix: "ST",
	CostingNumberPrefix:       "CO",
	OTPPrefix:                 "OTP",
	CompanyCodePrefix:         "C",
	PRNumberPrefix:            "PR",
	PONumberPrefix:            "PO",
	SRNumberPrefix:            "SR",
	ProductionNumberPrefix:    "PR",
	ReturnStockPrefix:         "RS",
	OpnameNumberPrefix:        "OP",
	FAPONumberPrefix:          "FP",
	FAReceiveNumberPrefix:     "FR",
	DepreciationNumberPrefix:  "DN",
	CityOtherCode:             "OTH",
	AddKey:                    8,
	MinPasswordLength:         6,
	MaxLoginTry:               3,
	MaxOTPGenerate:            50,
	ConfigurationAccess:       37}
var ConstTableName = TConstTableName{
	Report:              "report",
	ReportDefaultField:  "report_default_field",
	ReportTemplate:      "report_template",
	ReportTemplateField: "report_template_field",
	ReportGroupField:    "report_group_field",
	ReportOrderField:    "report_order_field",
	ReportGroupingField: "report_grouping_field"}
var SystemCode = TSystemCode{
	General:    "G",
	Hotel:      "H",
	Pos:        "P",
	Banquet:    "B",
	Accounting: "A",
	Asset:      "I",
	Payroll:    "R",
	Corporate:  "0",
	Report:     "R",
	Tools:      "T"}
var ConfigurationCategory = TConfigurationCategory{
	General:                      "GENERAL",
	Form:                         "FORM",
	FormatSetting:                "FORMAT_SETTING",
	Grid:                         "GRID",
	Report:                       "REPORT",
	WeekendDay:                   "WEEKEND_DAY",
	ReportSignature:              "REPORT_SIGNATURE",
	ReportTemplate:               "REPORT_TEMPLATE",
	Other:                        "OTHER",
	Reservation:                  "RESERVATION",
	DefaultVariable:              "DEFAULT_VARIABLE",
	AmountPreset:                 "AMOUNT_PRESET",
	CustomField:                  "CUSTOM_FIELD",
	CustomLookupField:            "CUSTOM_LOOKUP_FIELD",
	MemberVoucherGift:            "MEMBER_VOUCHER_GIFT",
	Accounting:                   "ACCOUNTING",
	GlobalAccount:                "GLOBAL_ACCOUNT",
	GlobalJournalAccount:         "GLOBAL_GL_ACCOUNT",
	GlobalJournalAccountSubGroup: "GLOBAL_GL_ACCOUNT_SG",
	PaymentAccount:               "PAYMENT_ACCOUNT",
	SubDepartment:                "SUB_DEPARTMENT",
	GlobalDepartment:             "GLOBAL_DEPARTMENT",
	GlobalSubDepartment:          "GLOBAL_SUB_DEPT",
	GlobalOther:                  "GLOBAL_OTHER",
	Personal:                     "PERSONAL",
	Invoice:                      "INVOICE",
	Folio:                        "FOLIO",
	OtherForm:                    "OTHER_FORM",
	CompanyBankAccount:           "COMPANY_BANK_ACCOUNT",
	PaymentCityLedger:            "PAYMENT_CITY_LEDGER",
	Company:                      "COMPANY",
	FloorPlan:                    "FLOOR_PLAN",
	RoomStatusColor:              "ROOM_STATUS_COLOR",
	Keylock:                      "KEYLOCK",
	PABX:                         "PABX",
	RoomCosting:                  "ROOM_COSTING",
	OtherHK:                      "OTHER_HK",
	ServiceIPTV:                  "SERVICE_IPTV",
	ServiceCCMS:                  "SERVICE_CCMS",
	Mikrotik:                     "MIKROTIK",
	Notification:                 "NOTIFICATION",
	WAAPI:                        "WA_API",
	Email:                        "EMAIL",
	Inventory:                    "INVENTORY",
	DayendClosed:                 "DAYEND_CLOSED"}
var ConfigurationCategoryCAMS = TConfigurationCategoryCAMS{
	General:                      "GENERAL",
	Form:                         "FORM",
	FormatSetting:                "FORMAT_SETTING",
	Grid:                         "GRID",
	Report:                       "REPORT",
	DefaultVariable:              "DEFAULT_VARIABLE",
	Other:                        "OTHER",
	Company:                      "COMPANY",
	Reservation:                  "RESERVATION",
	Accounting:                   "ACCOUNTING",
	Personal:                     "PERSONAL",
	Invoice:                      "INVOICE",
	Inventory:                    "INVENTORY",
	GlobalJournalAccount:         "GLOBAL_GL_ACCOUNT",
	GlobalJournalAccountSubGroup: "GLOBAL_GL_ACCOUNT_SG",
	GlobalSubDepartment:          "GLOBAL_SUB_DEPT",
	ReportForm:                   "REPORT_FORM",
	PurchaseRequestApp:           "PURCHASE_REQUEST_APP",
	ReportTemplate:               "REPORT_TEMPLATE",
	Notification:                 "NOTIFICATION"}
var ConfigurationName = TConfigurationName{
	//General
	Timezone:             "TIMEZONE",
	DatabaseVersion:      "DATABASE_VERSION",
	IsReplication:        "IS_REPLICATION",
	UseChildRate:         "USE_RATE_CHILD",
	IsRoomByName:         "IS_ROOM_BY_NAME",
	LogOffTriggerAddress: "LOG_OFF_TRIGGER_ADDRESS",
	//Form
	FormMDIChild:       "FORM_MDI_CHILD",
	RequiredFieldColor: "REQUIRED_FIELD_COLOR",
	//Format Setting
	ShortDateFormat:    "SHORT_DATE_FORMAT",
	DateSeparator:      "DATE_SEPARATOR",
	CurrencyFormat:     "CURRENCY_FORMAT",
	DecimalSeparator:   "DECIMAL_SEPARATOR",
	ThousandsSeparator: "THOUSANDS_SEPARATOR",
	//Grid
	ShowGridGroupByBox: "SHOW_GRID_GROUP_BY_BOX",
	GridCanEditCell:    "GRID_CAN_EDIT_CELL",
	LimitGrid:          "LIMIT_GRID",
	GridRowHeight:      "GRID_ROW_HEIGHT",
	DueDateColor:       "DUE_DATE_COLOR",
	//Report
	LogoWidth:            "LOGO_WIDTH",
	ReportHeaderAligment: "REPORT_HEADER_ALIGNMENT",
	//Weekend Day
	FridayAsWeekend:   "FRIDAY_AS_WEEKEND",
	SaturdayAsWeekend: "SATURDAY_AS_WEEKEND",
	SundayAsWeekend:   "SUNDAY_AS_WEEKEND",
	//Report Signature
	ShowPreparedBy:            "SHOW_PREPARED_BY",
	PreparedBy:                "PREPARED_BY",
	ShowCheckedBy:             "SHOW_CHECKED_BY",
	CheckedBy:                 "CHECKED_BY",
	ShowApprovedBy:            "SHOW_APPROVED_BY",
	ApprovedBy:                "APPROVED_BY",
	ShowInDailyRevenueReport:  "SHOW_IN_DAILY_REVENUE_REPORT",
	ShowInRoomStatisticReport: "SHOW_IN_ROOM_STATISTIC_REPORT",
	//Report Template
	DepositReceipt:              "DEPOSIT_RECEIPT",
	DepositReceiptRefund:        "DEPOSIT_RECEIPT_REFUND",
	FolioReceipt:                "FOLIO_RECEIPT",
	FolioReceiptRefund:          "FOLIO_RECEIPT_REFUND",
	MiscellaneousCharge:         "MISCELLANEOUS_CHARGE",
	CashierReport:               "CASHIER_REPORT",
	CashRemittance:              "CASH_REMITTANCE",
	BreakfastList:               "BREAKFAST_LIST",
	PickUpService:               "PICK_UP_SERVICE",
	DropService:                 "DROP_SERVICE",
	HKChecklist:                 "HK_CHECK_LIST",
	HKRoomStatus:                "HK_ROOM_STATUS",
	HKRoomStatusSummary:         "HK_ROOM_STATUS_SUMMARY",
	HKRoomDiscrepancy:           "HK_ROOM_DISCREPANCY",
	HKRoomAttendantControlSheet: "HK_ROOM_ATTENDANT_CONTROL_SHEET",
	DailyRevenueReport:          "DAILY_REVENUE_REPORT",
	DailyRevenueReportSummary:   "DAILY_REVENUE_REPORT_SUMMARY",
	RevenueBySubDepartment:      "REVENUE_BY_SUB_DEPARTMENT",
	GuestInHouseListing:         "GUEST_IN_HOUSE_LISTING",
	RoomStatistic:               "ROOM_STATISTIC",
	DailyFlashReport:            "DAILY_FLASH_REPORT",
	RoomProduction:              "ROOM_PRODUCTION",
	GuestForecast:               "GUEST_FORECAST",
	RTIncomeStatement:           "INCOME_STATEMENT",
	//Other
	ShowTransferOnCashierReport:     "SHOW_TRANSFER_ON_CASHIER_REPORT",
	ShowComplimentOnCashierReport:   "SHOW_COMPLIMENT_ON_CASHIER_REPORT",
	IncomeBudgetCalculateEachDay:    "INCOME_BUDGET_CALCULATE_EACH_DAY",
	InsertServerIDOnNumber:          "INSERT_SERVER_ID_ON__NUMBER",
	CalculateAllRoomRevenueSubGroup: "CALCULATE_ALL_ROOM_REVENUE",
	CompanyTypeTravelAgent:          "COMPANY_TYPE_TRAVEL_AGENT",
	//Reservation
	PostFirstNightCharge:     "POST_FIRST_NIGHT_CHARGE",
	PostDiscount:             "POST_DISCOUNT",
	KeylockLimit:             "KEYLOCK_LIMIT",
	CheckOutLimit:            "CHECK_OUT_LIMIT",
	IsRoomNumberRequired:     "IS_ROOM_NUMBER_REQUIRED",
	IsBusinessSourceRequired: "IS_BUSINESS_SOURCE_REQUIRED",
	IsMarketRequired:         "IS_MARKET_REQUIRED",
	IsTitleRequired:          "IS_TITLE_REQUIRED",
	IsStateRequired:          "IS_STATE_REQUIRED",
	IsCityRequired:           "IS_CITY_REQUIRED",
	IsNationalityRequired:    "IS_NATIONALITY_REQUIRED",
	IsPhone1Required:         "IS_PHONE1_REQUIRED",
	IsEmailRequired:          "IS_EMAIL_REQUIRED",
	IsCompanyRequired:        "IS_COMPANY_REQUIRED",
	IsPurposeOfRequired:      "IS_PURPOSE_OF_REQUIRED",
	IsTAVoucherRequired:      "IS_TA_VOUCHER_REQUIRED",
	IsHKNoteRequired:         "IS_HK_NOTE_REQUIRED",
	FilterRateByMarket:       "FILTER_RATE_BY_MARKET",
	FilterRateByCompany:      "FILTER_RATE_BY_COMPANY",
	AlwaysShowPublishRate:    "ALWAYS_SHOW_PUBLISH_RATE",
	LockFolioOnCheckIn:       "LOCK_FOLIO_ON_CHECK_IN",
	AutoGenerateCompanyCode:  "AUTO_GENERATE_COMPANY_CODE",
	CompanyCodeDigit:         "COMPANY_CODE_DIGIT",
	//Default Variable
	DVRoomType:         "DV_ROOM_TYPE",
	DVRoomRate:         "DV_ROOM_RATE",
	DVSubDepartment:    "DV_SUB_DEPARTMENT",
	DVPaymentType:      "DV_PAYMENT_TYPE",
	DVComplimentRate:   "DV_COMPLIMENT_RATE",
	DVHouseUseRate:     "DV_HOUSE_USE_RATE",
	DVMarket:           "DV_MARKET",
	DVIndividualMarket: "DV_INDIVIDUAL_MARKET",
	//Amount Preset
	APPreset1: "PRESET1",
	APPreset2: "PRESET2",
	APPreset3: "PRESET3",
	APPreset4: "PRESET4",
	APPreset5: "PRESET5",
	APPreset6: "PRESET6",
	APPreset7: "PRESET7",
	APPreset8: "PRESET8",
	//Custom Field
	CustomFieldName01:         "CUSTOM_FIELD_NAME01",
	CustomFieldName02:         "CUSTOM_FIELD_NAME02",
	CustomFieldName03:         "CUSTOM_FIELD_NAME03",
	CustomFieldName04:         "CUSTOM_FIELD_NAME04",
	CustomFieldName05:         "CUSTOM_FIELD_NAME05",
	CustomFieldName06:         "CUSTOM_FIELD_NAME06",
	CustomFieldName07:         "CUSTOM_FIELD_NAME07",
	CustomFieldName08:         "CUSTOM_FIELD_NAME08",
	CustomFieldName09:         "CUSTOM_FIELD_NAME09",
	CustomFieldName10:         "CUSTOM_FIELD_NAME10",
	CustomFieldName11:         "CUSTOM_FIELD_NAME11",
	CustomFieldName12:         "CUSTOM_FIELD_NAME12",
	CustomFieldDefaultValue01: "CUSTOM_FIELD_DEFAULT_VALUE01",
	CustomFieldDefaultValue02: "CUSTOM_FIELD_DEFAULT_VALUE02",
	CustomFieldDefaultValue03: "CUSTOM_FIELD_DEFAULT_VALUE03",
	CustomFieldDefaultValue04: "CUSTOM_FIELD_DEFAULT_VALUE04",
	CustomFieldDefaultValue05: "CUSTOM_FIELD_DEFAULT_VALUE05",
	CustomFieldDefaultValue06: "CUSTOM_FIELD_DEFAULT_VALUE06",
	CustomFieldDefaultValue07: "CUSTOM_FIELD_DEFAULT_VALUE07",
	CustomFieldDefaultValue08: "CUSTOM_FIELD_DEFAULT_VALUE08",
	CustomFieldDefaultValue09: "CUSTOM_FIELD_DEFAULT_VALUE09",
	CustomFieldDefaultValue10: "CUSTOM_FIELD_DEFAULT_VALUE10",
	CustomFieldDefaultValue11: "CUSTOM_FIELD_DEFAULT_VALUE11",
	CustomFieldDefaultValue12: "CUSTOM_FIELD_DEFAULT_VALUE12",
	//Custom Lookup Field
	CustomLookupFieldName01:         "CUSTOM_LOOKUP_FIELD_NAME01",
	CustomLookupFieldName02:         "CUSTOM_LOOKUP_FIELD_NAME02",
	CustomLookupFieldName03:         "CUSTOM_LOOKUP_FIELD_NAME03",
	CustomLookupFieldName04:         "CUSTOM_LOOKUP_FIELD_NAME04",
	CustomLookupFieldName05:         "CUSTOM_LOOKUP_FIELD_NAME05",
	CustomLookupFieldName06:         "CUSTOM_LOOKUP_FIELD_NAME06",
	CustomLookupFieldName07:         "CUSTOM_LOOKUP_FIELD_NAME07",
	CustomLookupFieldName08:         "CUSTOM_LOOKUP_FIELD_NAME08",
	CustomLookupFieldName09:         "CUSTOM_LOOKUP_FIELD_NAME09",
	CustomLookupFieldName10:         "CUSTOM_LOOKUP_FIELD_NAME10",
	CustomLookupFieldName11:         "CUSTOM_LOOKUP_FIELD_NAME11",
	CustomLookupFieldName12:         "CUSTOM_LOOKUP_FIELD_NAME12",
	CustomLookupFieldDefaultValue01: "CUSTOM_LOOKUP_FIELD_DEFAULT_VALUE01",
	CustomLookupFieldDefaultValue02: "CUSTOM_LOOKUP_FIELD_DEFAULT_VALUE02",
	CustomLookupFieldDefaultValue03: "CUSTOM_LOOKUP_FIELD_DEFAULT_VALUE03",
	CustomLookupFieldDefaultValue04: "CUSTOM_LOOKUP_FIELD_DEFAULT_VALUE04",
	CustomLookupFieldDefaultValue05: "CUSTOM_LOOKUP_FIELD_DEFAULT_VALUE05",
	CustomLookupFieldDefaultValue06: "CUSTOM_LOOKUP_FIELD_DEFAULT_VALUE06",
	CustomLookupFieldDefaultValue07: "CUSTOM_LOOKUP_FIELD_DEFAULT_VALUE07",
	CustomLookupFieldDefaultValue08: "CUSTOM_LOOKUP_FIELD_DEFAULT_VALUE08",
	CustomLookupFieldDefaultValue09: "CUSTOM_LOOKUP_FIELD_DEFAULT_VALUE09",
	CustomLookupFieldDefaultValue10: "CUSTOM_LOOKUP_FIELD_DEFAULT_VALUE10",
	CustomLookupFieldDefaultValue11: "CUSTOM_LOOKUP_FIELD_DEFAULT_VALUE11",
	CustomLookupFieldDefaultValue12: "CUSTOM_LOOKUP_FIELD_DEFAULT_VALUE12",
	//Member Voucher and Gift
	MemberPointExpire:             "MEMBER_POINT_EXPIRE",
	MemberAutoUpdateMemberProfile: "MEMBER_AUTO_UPDATE_MEMBER_PROFILE",
	VoucherLength:                 "VOUCHER_LENGTH",
	VoucherExpire:                 "VOUCHER_EXPIRE",
	VoucherPointRedeem:            "VOUCHER_POINT_REDEEM",
	VoucherDefaultPrice:           "VOUCHER_DEFAULT_PRICE",
	VoucherDescription:            "VOUCHER_DESCRIPTION",
	VoucherTemplate:               "VOUCHER_TEMPLATE",
	VoucherSaleDiscountTemplate:   "VOUCHER_TEMPLATE_SALE_DISCOUNT",
	//Accounting
	IsAccrualBase:                      "IS_ACCRUAL_BASE",
	SubDepartmentAllCCAdmin:            "SUB_DEPARTMENT_ALL_CC_ADMIN",
	AutomaticCreateInvoiceCLAtCheckOut: "AUTOMATIC_CREATE_INVOICE_CL_AT_CHECK_OUT",
	//Global Account
	AccountRoomCharge:                        "ACCOUNT_ROOM_CHARGE",
	AccountExtraBed:                          "ACCOUNT_EXTRA_BED",
	AccountCancellationFee:                   "ACCOUNT_CANCELLATION_FEE",
	AccountNoShow:                            "ACCOUNT_NO_SHOW",
	AccountBreakfast:                         "ACCOUNT_BREAKFAST",
	AccountTelephone:                         "ACCOUNT_TELEPHONE",
	AccountAPRefundDeposit:                   "ACCOUNT_AP_REFUND_DEPOSIT",
	AccountAPCommission:                      "ACCOUNT_AP_COMMISSION",
	AccountCreditCardAdm:                     "ACCOUNT_CC_ADM",
	AccountCash:                              "ACCOUNT_CASH",
	AccountCityLedger:                        "ACCOUNT_CITY_LEDGER",
	AccountVoucher:                           "ACCOUNT_VOUCHER",
	AccountVoucherCompliment:                 "ACCOUNT_VOUCHER_COMPLIMENT",
	AccountTax:                               "ACCOUNT_TAX",
	AccountService:                           "ACCOUNT_SERVICE",
	AccountTransferDepositReservation:        "ACCOUNT_TRANSFER_DEPOSIT_RESERVATION",
	AccountTransferDepositReservationToFolio: "ACCOUNT_TRANSFER_DEPOSIT_RESERVATION_TO_FOLIO",
	AccountTransferCharge:                    "ACCOUNT_TRANSFER_CHARGE",
	AccountTransferPayment:                   "ACCOUNT_TRANSFER_PAYMENT",
	//Global Journal Account
	JAGuestLedger:             "GUEST_LEDGER",
	JAAPVoucher:               "AP_VOUCHER",
	JAGuestDeposit:            "GUEST_DEPOSIT",
	JAGuestDepositReservation: "GUEST_DEPOSIT_RESERVATION",
	JAPLBeginningYear:         "PL_BEGINNING_YEAR",
	JAPLCurrentYear:           "PL_CURRENT_YEAR",
	JAPLCurrency:              "PL_CURRENCY",
	JAIncomeVoucherExpire:     "INCOME_VOUCHER_EXPIRE",
	JAInvoiceDiscount:         "INVOICE_DISCOUNT",
	JACreditCardAdm:           "CREDIT_CARD_ADMINISTRATION",
	JAAPPaymentTemporer:       "AP_PAYMENT_TEMPORER",
	JAIncomeVoucherExpired:    "INCOME_VOUCHER_EXPIRE",
	JAOverShortAsIncome:       "OVER_SHORT_AS_INCOME",
	JAOverShortAsExpense:      "OVER_SHORT_AS_EXPENSE",
	JAServiceRevenue:          "SERVICE_REVENUE",
	//Global Journal Account Sub Group
	JASGInventory:        "INVENTORY",
	JASGFixedAsset:       "FIXED_ASSET",
	JASGAccmDepreciation: "ACCUMULATED_DEPRECIATION",
	JASGAccountPayable:   "ACCOUNT_PAYABLE",
	JASGManagementFee:    "MANAGEMENT_FEE",
	JASGDepreciation:     "DEPRECIATION",
	JASGAmortization:     "AMORTIZATION",
	JASGLoanInterest:     "LOAN_INTEREST",
	JASGIncomeTax:        "INCOME_TAX",
	//Global Department
	DRoomDivision:  "DEPARTMENT_ROOM_DIVISION",
	DFoodBeverage:  "DEPARTMENT_FOOD_BEVERAGE",
	DBanquet:       "DEPARTMENT_BANQUET",
	DMinor:         "DEPARTMENT_MINOR",
	DMiscellaneous: "DEPARTMENT_MISCELLANEOUS",
	//Payment Account
	PYAccount01: "ACCOUNT01",
	PYAccount02: "ACCOUNT02",
	PYAccount03: "ACCOUNT03",
	PYAccount04: "ACCOUNT04",
	PYAccount05: "ACCOUNT05",
	PYAccount06: "ACCOUNT06",
	PYAccount07: "ACCOUNT07",
	PYAccount08: "ACCOUNT08",
	PYAccount09: "ACCOUNT09",
	PYAccount10: "ACCOUNT10",
	PYAccount11: "ACCOUNT11",
	PYAccount12: "ACCOUNT12",
	//Payment Account Name
	PYNAccount01: "ACCOUNT_NAME01",
	PYNAccount02: "ACCOUNT_NAME02",
	PYNAccount03: "ACCOUNT_NAME03",
	PYNAccount04: "ACCOUNT_NAME04",
	PYNAccount05: "ACCOUNT_NAME05",
	PYNAccount06: "ACCOUNT_NAME06",
	PYNAccount07: "ACCOUNT_NAME07",
	PYNAccount08: "ACCOUNT_NAME08",
	PYNAccount09: "ACCOUNT_NAME09",
	PYNAccount10: "ACCOUNT_NAME10",
	PYNAccount11: "ACCOUNT_NAME11",
	PYNAccount12: "ACCOUNT_NAME12",
	//Sub Department
	SD01: "SD01",
	SD02: "SD02",
	SD03: "SD03",
	SD04: "SD04",
	SD05: "SD05",
	SD06: "SD06",
	SD07: "SD07",
	SD08: "SD08",
	SD09: "SD09",
	SD10: "SD10",
	SD11: "SD11",
	SD12: "SD12",
	//Sub Department Name
	SDN01: "SD_NAME01",
	SDN02: "SD_NAME02",
	SDN03: "SD_NAME03",
	SDN04: "SD_NAME04",
	SDN05: "SD_NAME05",
	SDN06: "SD_NAME06",
	SDN07: "SD_NAME07",
	SDN08: "SD_NAME08",
	SDN09: "SD_NAME09",
	SDN10: "SD_NAME10",
	SDN11: "SD_NAME11",
	SDN12: "SD_NAME12",
	//Global Sub Department
	SDFrontOffice:  "SUB_DEPARTMENT_FRONT_OFFICE",
	SDHouseKeeping: "SUB_DEPARTMENT_HOUSE_KEEPING",
	SDBanquet:      "SUB_DEPARTMENT_BANQUET",
	SDAccounting:   "SUB_DEPARTMENT_ACCOUNTING",
	//Global Other
	GOPaymentType: "PAYMENT_TYPE",
	//Personal
	DRName:     "DR_NAME",
	DRPosition: "DR_POSITION",
	HMName:     "HM_NAME",
	HMPosition: "HM_POSITION",
	FHName:     "FH_NAME",
	FHPosition: "FH_POSITION",
	IAName:     "IA_NAME",
	IAPosition: "IA_POSITION",
	GCName:     "GC_NAME",
	GCPosition: "GC_POSITION",
	ARName:     "AR_NAME",
	ARPosition: "AR_POSITION",
	APName:     "AP_NAME",
	APPosition: "AP_POSITION",
	//Invoice
	InvoiceRemark:   "INVOICE_REMARK",
	InvoiceNote:     "INVOICE_NOTE",
	InvoiceAPNote:   "INVOICE_AP_NOTE",
	InvoiceTemplate: "INVOICE_TEMPLATE",
	//Folio
	DefaultFolio:             "DEFAULT_FOLIO",
	FolioFooter:              "FOLIO_FOOTER",
	AllowZeroAmount:          "ALLOW_ZERO_AMOUNT",
	PrintRegFormAfterCheckIn: "PRINT_REG_FORM_AFTER_CHECK_IN",
	ShowRate:                 "SHOW_RATE",
	//Other Form
	RegistrationFormReservation: "REGISTRATION_FORM_RESERVATION",
	RegistrationFormInHouse:     "REGISTRATION_FORM_IN_HOUSE",
	ConfirmationLetter:          "CONFIRMATION_LETTER",
	ConfirmationLetterSelected:  "CONFIRMATION_LETTER_SELECTED",
	GuaranteedLetter:            "GUARANTEE_LETTER",
	ProformaInvoice:             "PROFORMA_INVOICE",
	ProformaInvoiceDetail:       "PROFORMA_INVOICE_DETAIL",
	TaxServiceRemark:            "TAX_AND_SERVICE_REMARK",
	GuaranteeLetterRemark:       "GUARANTEE_LETTER_REMARK",
	//Company Bank Account
	BankName1:    "BANK_NAME",
	BankAccount1: "BANK_ACCOUNT",
	HolderName1:  "HOLDER_NAME",
	BankName2:    "BANK_NAME2",
	BankAccount2: "BANK_ACCOUNT2",
	HolderName2:  "HOLDER_NAME2",
	//Payment City Ledger
	CanCLOverLimit: "CAN_CL_OVER_LIMIT",
	//Company
	CompanyContactPersonRequired: "COMPANY_CONTACT_PERSON_REQUIRED",
	CompanyStreetRequired:        "COMPANY_STREET_REQUIRED",
	CompanyCityRequired:          "COMPANY_CITY_REQUIRED",
	CompanyCountryRequired:       "COMPANY_COUNTRY_REQUIRED",
	CompanyStateRequired:         "COMPANY_STATE_REQUIRED",
	CompanyPostalCodeRequired:    "COMPANY_POSTAL_CODE_REQUIRED",
	CompanyPhone1Required:        "COMPANY_PHONE1_REQUIRED",
	//Floor Plan
	MinRoomWidth:  "MIN_ROOM_WIDTH",
	MinRoomHeight: "MIN_ROOM_HEIGHT",
	LeftMargin:    "LEFT_MARGIN",
	TopMargin:     "TOP_MARGIN",
	TileColumn:    "TILE_COLUMN",
	TileDistance:  "TILE_DISTANCE",
	SnapGrid:      "SNAP_GRID",
	//Room Status Color
	Reserved:          "RESERVED",
	Occupied:          "OCCUPIED",
	HouseUse:          "HOUSE_USE",
	Compliment:        "COMPLIMENT",
	OutOfOrder:        "OUT_OF_ORDER",
	OfficeUse:         "OFFICE_USE",
	UnderConstruction: "UNDER_CONSTRUCTION",
	Available:         "AVAILABLE",
	//Keylock
	KeylockVendor:        "KEYLOCK_VENDOR",
	PMSServiceAddress:    "PMS_SERVICE_ADDRESS",
	PMSServicePort:       "PMS_SERVICE_PORT",
	SureLockHotelID:      "SURE_LOCK_HOTEL_ID",
	KimaOperatorCode:     "KIMA_OPERATOR_CODE",
	VingCardKeyType:      "VING_CARD_KEY_TYPE",
	VingCardUserGroup:    "VING_CARD_USER_GROUP",
	VingCardSource:       "VING_CARD_SOURCE",
	VingCardDestination:  "VING_CARD_DESTINATION",
	BeTechPort:           "BE_TECH_PORT",
	BeTechReaderType:     "BE_TECH_READER_TYPE",
	BeTechDatabaseType:   "BE_TECH_DATABASE_TYPE",
	BeTechServerName:     "BE_TECH_SERVER_NAME",
	BeTechUserName:       "BE_TECH_USER_NAME",
	BeTechPassword:       "BE_TECH_PASSWORD",
	OnityDBType:          "ONITY_DATABASE_TYPE",
	OnityServerName:      "ONITY_SERVER_NAME",
	OnitySoftwareType:    "ONITY_SOFTWARE_TYPE",
	OnityPort:            "ONITY_PORT",
	SaflokToStation:      "SAFLOK_TO_STATION",
	SaflokFromStation:    "SAFLOK_FROM_STATION",
	SaflokRequestNumber:  "SAFLOK_REQUEST_NUMBER",
	SaflokPassword:       "SAFLOK_PASSWORD",
	SaflokEncoderStation: "SAFLOK_ENCODER_STATION",
	VisionlinePMSAddress: "VISIONLINE_PMS_ADDRESS",
	DownsAuthCode:        "DOWNS_AUTH_CODE",
	VentazaUseLift:       "VENTAZA_USE_LIFT",
	//PABX
	PABXChargePerSecond:      "CHARGE_PER_SECOND",
	PABXChargePerSecondOther: "CHARGE_PER_SECOND_OTHER",
	PABXAutoPostToFolio:      "AUTO_CHARGE_TO_FOLIO",
	PABXLANSMDRPassword:      "LAN_SMDR_PASSWORD",
	PABXLocal:                "_LOCAL",
	PABXSLJJ:                 "_SLJJ",
	PABXSLI:                  "_SLI",
	//Room Costing
	DefaultStore: "DEFAULT_STORE",
	//Other HK
	VRVCRoomStatus:               "VR_VC_ROOM_STATUS",
	ShowRemarkOnRoomAvailability: "SHOW_REMARK_ON_ROOM_AVAILABILITY",
	//IPTV
	IPTVVendor:   "IPTV_VENDOR",
	IPTVPassword: "IPTV_PASSWORD",
	IPTVSubFolio: "IPTV_SUB_FOLIO",
	IPTVAddress:  "IPTV_ADDRESS",
	//CCMS
	CCMSVendor:                         "CHANNEL_MANAGER_VENDOR",
	CCMSSMUser:                         "SITE_MINDER_USER",
	CCMSSMPassword:                     "SITE_MINDER_PASSWORD",
	CCMSSMRequestorID:                  "SITE_MINDER_REQUESTOR_ID",
	CCMSSMHotelCode:                    "SITE_MINDER_HOTEL_CODE",
	CCMSSMWSDL:                         "SITE_MINDER_WSDL",
	CCMSSMReservationAsAllotment:       "SITE_MINDER_RESERVATION_AS_ALLOTMENT",
	CCMSSMSynchronizeReservation:       "SITE_MINDER_SYNCHRONIZE_RESERVATION",
	CCMSSMSynchronizeAvailability:      "SITE_MINDER_SYNCHRONIZE_AVAILABILITY",
	CCMSSMSynchronizeRate:              "SITE_MINDER_SYNCHRONIZE_RATE",
	CCMSSTAAHGlobalPercentAvailability: "CHANNEL_MANAGER_GLOBAL_PERCENT_AVAILABILITY",
	CCMSSTAAHGlobalMinRoomLeft:         "CHANNEL_MANAGER_GLOBAL_MIN_ROOM_LEFT",
	//Mikrotik
	MikrotikVersion:          "MIKROTIK_VERSION",
	MikrotikVendor:           "MIKROTIK_VENDOR",
	MikrotikAddress:          "MIKROTIK_ADDRESS",
	MikrotikUser:             "MIKROTIK_USER",
	MikrotikPassword:         "MIKROTIK_PASSWORD",
	MikrotikHotspotServer:    "MIKROTIK_HOTSPOT_SERVER",
	MikrotikUserProfile:      "MIKROTIK_USER_PROFILE",
	MikrotikServiceAddress:   "MIKROTIK_SERVICE_ADDRESS",
	MikrotikAutoCreateUser:   "MIKROTIK_AUTO_CREATE_USER",
	MikrotikSameUserPassword: "MIKROTIK_SAME_USER_PASSWORD",
	MikrotikAskForPassword:   "MIKROTIK_ASK_FOR_PASSWORD",
	MikrotikPasswordOption:   "MIKROTIK_PASSWORD_OPTION",
	//Notification
	NotificationActive:   "NOTIFICATION_ACTIVE",
	NotificationDuration: "NOTIFICATION_DURATION",
	NotificationDelay:    "NOTIFICATION_DELAY",
	NotificationWidth:    "NOTIFICATION_WIDTH",
	NotificationHeight:   "NOTIFICATION_HEIGHT",
	NotificationStyle:    "NOTIFICATION_STYLE",
	NotificationPosition: "NOTIFICATION_POSITION",
	//WA API
	WAAPIVendor:             "WA_VENDOR",
	WAAPIKey:                "WA_API_KEY",
	WAAPIURL:                "WA_API_URL",
	WAAPITemplateCheckIn:    "WA_API_TEMPLATE_CHECK_IN",
	WAAPITemplateCheckOut:   "WA_API_TEMPLATE_CHECK_OUT",
	WAOnReminderReservation: "REMINDER_RESERVATION",
	WAOnCheckIn:             "ON_CHECK_IN",
	WAOnCheckOut:            "ON_CHECK_OUT",
	WAOnGuestBirthday:       "ON_GUEST_BIRTHDAY",
	WAOnSalesReminder:       "SALES_ACTIVITY_SEND_REMINDER",
	//Email
	EmailServer:                    "SERVER",
	EmailAuthentication:            "AUTHENTICATION",
	EmailPort:                      "PORT",
	EmailUser:                      "USER",
	EmailPassword:                  "PASSWORD",
	EmailUserSSL:                   "USE_SSL",
	EmailUseStartTLS:               "USE_START_TLS",
	EmailReminderReservation:       "REMINDER_RESERVATION",
	EmailOnCheckIn:                 "ON_CHECK_IN",
	EmailOnCheckOut:                "ON_CHECK_OUT",
	EmailOnGuestBirthday:           "ON_GUEST_BIRTHDAY",
	EmailSalesActivitySendreminder: "SALES_ACTIVITY_SEND_REMINDER",
	//Inventory
	CostingMethod:                    "COSTING_METHOD",
	SynchronizePOAndReceive:          "SYNCHRONIZE_PO_AND_RECEIVE",
	DefaultShippingAddress:           "DEFAULT_SHIPPING_ADDRESS",
	ReceiveStockAPTwoDigitDecimal:    "RECEIVE_STOCK_AP_TWO_DIGIT_DECIMAL",
	ShowPurchasePricePRSR:            "SHOW_PURCHASE_PRICE_PR_SR",
	LockTransactionDateInventory:     "LOCK_TRANSACTION_DATE_INVENTORY",
	IsPurchasingApproval:             "IS_PURCHASING_APPROVAL",
	IsCompanyPRApplyPriceMoreThanOne: "IS_COMPANY_PR_APLLY_PRICE_MORE_THAN_ONE",
	//DayendClosed
	AutoImportJournal: "AUTO_IMPORT_JOURNAL",

	//Personal
	CCName:     "CC_NAME",
	CCPosition: "CC_POSITION",
	PMName:     "PM_NAME",
	PMPosition: "PM_POSITION",
	//Invoice
	//Purchase Request Approval
	UserApproval1: "USER_APPROVAL1",
	UserApproval2: "USER_APPROVAL2",
	UserApproval3: "USER_APPROVAL3",
	//SR & PR Color
	DHColor: "DH_COLOR",
	CCColor: "CC_COLOR",
	FNColor: "FN_COLOR",
	RJColor: "RJ_COLOR",
	//All Store Stock & Store Stock
	LowStockColor: "LOW_STOCK_COLOR",
	//Inventory
	//Company
	//Global Journal Accounting
	JAAPSupplier:         "AP_SUPPLIER",
	JAPurchasingDiscount: "PURCHASING_DISCOUNT",
	JAPurchasingTax:      "PURCHASING_TAX",
	JAPurchasingShipping: "PURCHASING_SHIPPING",
	JAIncomeReturnStock:  "INCOME_RETURN_STOCK",
	JAExpenseReturnStock: "EXPENSE_RETURN_STOCK",
	//Global Sub Department
	CompanyTypeSupplier:   "COMPANY_TYPE_SUPPLIER",
	CompanyTypeExpedition: "COMPANY_TYPE_EXPEDITION",
	//Report Form
	PurchaseRequest:  "PURCHASE_REQUEST",
	PurchaseOrder:    "PURCHASE_ORDER",
	ReceiveStock:     "RECEIVE_STOCK",
	StoreRequitition: "STORE_REQUISITION",
	StockTransfer:    "STOCK_TRANSFER",
	Costing:          "COSTING",
	FAPurchaseOrder:  "FA_PURCHASE_ORDER",
	FAReceive:        "FA_RECEIVE",
	//Report Template
}

var ConfigurationCategoryPOS = TConfigurationCategoryPOS{
	General:              "GENERAL",
	Payment:              "PAYMENT",
	Form:                 "FORM",
	FormatSetting:        "FORMAT_SETTING",
	Grid:                 "GRID",
	Report:               "REPORT",
	DefaultVariable:      "DEFAULT_VARIABLE",
	Other:                "OTHER",
	Company:              "COMPANY",
	Bill:                 "BILL",
	ReportTemplate:       "REPORT_TEMPLATE",
	Accounting:           "ACCOUNTING",
	GlobalAccount:        "GLOBAL_ACCOUNT",
	GlobalJournalAccount: "GLOBAL_GL_ACCOUNT",
	GlobalSubDepartment:  "GLOBAL_SUB_DEPT",
	TableView:            "TABLE_VIEW",
	Reservation:          "POS_RESERVATION",
	Notification:         "NOTIFICATION",
	CashDrawer:           "CASH_DRAWER",
	CustomerDisplay:      "CUSTOMER_DISPLAY",
	KitchenPrinter:       "KITCHEN_PRINTER",
	Inventory:            "INVENTORY"}
var ConfigurationNamePOS = TConfigurationNamePOS{
	//General
	DatabaseVersion:      "DATABASE_VERSION",
	RoomServiceOutlet:    "ROOM_SERVICE_OUTLET",
	IsRoomByName:         "IS_ROOM_BY_NAME",
	LogOffTriggerAddress: "LOG_OFF_TRIGGER_ADDRESS",
	//Payment
	CanPaymentIfChargeZero: "CAN_PAYMENT_IF_CHARGE_ZERO",
	//Other POS
	MemberMustScanFinger: "MEMBER_MUST_SCAN_FINGER",
	//Form
	FormMDIChild:       "FORM_MDI_CHILD",
	RequiredFieldColor: "REQUIRED_FIELD_COLOR",
	//Format Setting
	ShortDateFormat:    "SHORT_DATE_FORMAT",
	DateSeparator:      "DATE_SEPARATOR",
	CurrencyFormat:     "CURRENCY_FORMAT",
	DecimalSeparator:   "DECIMAL_SEPARATOR",
	ThousandsSeparator: "THOUSANDS_SEPARATOR",
	//Grid
	ShowGridGroupByBox: "SHOW_GRID_GROUP_BY_BOX",
	GridCanEditCell:    "GRID_CAN_EDIT_CELL",
	LimitGrid:          "LIMIT_GRID",
	GridRowHeight:      "GRID_ROW_HEIGHT",
	//Report
	LogoWidth:            "LOGO_WIDTH",
	ReportHeaderAligment: "REPORT_HEADER_ALIGNMENT",
	//Default Variable
	DVSubDepartment: "DV_SUB_DEPARTMENT",
	DVIPTVMarket:    "DV_IPTV_MARKET",
	//Other
	ShowTransferOnCashierReport:   "SHOW_TRANSFER_ON_CASHIER_REPORT",
	ShowComplimentOnCashierReport: "SHOW_COMPLIMENT_ON_CASHIER_REPORT",
	InsertServerIDOnNumber:        "INSERT_SERVER_ID_ON__NUMBER",
	//Bill
	DefaultBill:                 "DEFAULT_BILL",
	BillFileName:                "BILL_FILE_NAME",
	CaptainOrderFileName:        "CAPTAIN_ORDER_FILE_NAME",
	CaptainOrderStationFileName: "CAPTAIN_ORDER_STATION_FILE_NAME",
	//Report Template
	RTCashierReport:    "CASHIER_REPORT",
	RTCashRemittance:   "CASH_REMITTANCE",
	RTFNBRateStructure: "FNB_RATE_STRUCTURE",
	RTIncomeStatement:  "INCOME_STATEMENT",
	//Accounting
	IsAccrualBase:           "IS_ACCRUAL_BASE",
	SubDepartmentAllCCAdmin: "SUB_DEPARTMENT_ALL_CC_ADMIN",
	//Other
	ShowTaxService:                          "SHOW_TAX_SERVICE",
	PostDiscount:                            "POST_DISCOUNT",
	AutoCostingCostRecipeonCloseTransaction: "AUTO_COSTING_COST_RECIPE_ON_CLOSE_TRANSACTION",
	AutoGenerateCompanyCode:                 "AUTO_GENERATE_COMPANY_CODE",
	CompanyCodeDigit:                        "COMPANY_CODE_DIGIT",
	CompanyTypeSPATherapist:                 "COMPANY_TYPE_SPA_THERAPIST",
	//Company
	CompanyContactPersonRequired: "COMPANY_CONTACT_PERSON_REQUIRED",
	CompanyStreetRequired:        "COMPANY_STREET_REQUIRED",
	CompanyCityRequired:          "COMPANY_CITY_REQUIRED",
	CompanyCountryRequired:       "COMPANY_COUNTRY_REQUIRED",
	CompanyStateRequired:         "COMPANY_STATE_REQUIRED",
	CompanyPostalCodeRequired:    "COMPANY_POSTAL_CODE_REQUIRED",
	CompanyPhone1Required:        "COMPANY_PHONE1_REQUIRED",
	//Global Account
	AccountRoomCharge:                        "ACCOUNT_ROOM_CHARGE",
	AccountAPRefundDeposit:                   "ACCOUNT_AP_REFUND_DEPOSIT",
	AccountAPCommission:                      "ACCOUNT_AP_COMMISSION",
	AccountCreditCardAdm:                     "ACCOUNT_CC_ADM",
	AccountCash:                              "ACCOUNT_CASH",
	AccountCityLedger:                        "ACCOUNT_CITY_LEDGER",
	AccountTax:                               "ACCOUNT_TAX",
	AccountService:                           "ACCOUNT_SERVICE",
	AccountTransferDepositReservation:        "ACCOUNT_TRANSFER_DEPOSIT_RESERVATION",
	AccountTransferDepositReservationToFolio: "ACCOUNT_TRANSFER_DEPOSIT_RESERVATION_TO_FOLIO",
	AccountTransferCharge:                    "ACCOUNT_TRANSFER_CHARGE",
	AccountTransferPayment:                   "ACCOUNT_TRANSFER_PAYMENT",
	//Global Journal Account
	JAGuestLedger:             "GUEST_LEDGER",
	JAGuestDeposit:            "GUEST_DEPOSIT",
	JAGuestDepositReservation: "GUEST_DEPOSIT_RESERVATION",
	//Global Account
	SDAccounting: "SUB_DEPARTMENT_ACCOUNTING",
	//Table View
	MinRoomWidth:         "MIN_ROOM_WIDTH",
	MinRoomHeight:        "MIN_ROOM_HEIGHT",
	LeftMargin:           "LEFT_MARGIN",
	TopMargin:            "TOP_MARGIN",
	TileColumn:           "TILE_COLUMN",
	TileDistance:         "TILE_DISTANCE",
	SnapGrid:             "SNAP_GRID",
	TableColor:           "TABLE_COLOR",
	COColor1:             "CAPTAIN_ORDER COLOR_1",
	COColor2:             "CAPTAIN_ORDER COLOR_2",
	COColor3:             "CAPTAIN_ORDER COLOR_3",
	COColor4:             "CAPTAIN_ORDER COLOR_4",
	COColor5:             "CAPTAIN_ORDER COLOR_5",
	RVRangeTimeCODefault: "RANGE_TIME_CHECKOUT_DEFAULT",
	//Notification
	NotificationActive:   "NOTIFICATION_ACTIVE",
	NotificationDuration: "NOTIFICATION_DURATION",
	NotificationDelay:    "NOTIFICATION_DELAY",
	NotificationWidth:    "NOTIFICATION_WIDTH",
	NotificationHeight:   "NOTIFICATION_HEIGHT",
	NotificationStyle:    "NOTIFICATION_STYLE",
	NotificationPosition: "NOTIFICATION_POSITION",
	//Cash Drawer
	CashDrawerComport:   "CASH_DRAWER_COMPORT",
	CashDrawerOpenCode1: "CASH_DRAWER_OPEN_CODE1",
	CashDrawerOpenCode2: "CASH_DRAWER_OPEN_CODE2",
	CashDrawerOpenCode3: "CASH_DRAWER_OPEN_CODE3",
	CashDrawerOpenCode4: "CASH_DRAWER_OPEN_CODE4",
	CashDrawerOpenCode5: "CASH_DRAWER_OPEN_CODE5",
	CashDrawerOpenCode6: "CASH_DRAWER_OPEN_CODE6",
	//Customer Display
	CDVendor:    "CD_VENDOR",
	CDPort:      "CD_PORT",
	CDBaundRate: "CD_BAUNDRATE",
	CDParity:    "CD_PARITY",
	CDDataWidth: "CD_DATAWIDTH",
	StopBit:     "CD_STOPBIT",
	//Kitchen printer
	PrintCOAfterChangeRemoveCancel: "PRINT_CO_AFTER_CHANGE_REMOVE_CANCEL",
	//Inventory
	CostingMethod:                    "COSTING_METHOD",
	SynchronizePOAndReceive:          "SYNCHRONIZE_PO_AND_RECEIVE",
	DefaultShippingAddress:           "DEFAULT_SHIPPING_ADDRESS",
	ReceiveStockAPTwoDigitDecimal:    "RECEIVE_STOCK_AP_TWO_DIGIT_DECIMAL",
	ShowPurchasePricePRSR:            "SHOW_PURCHASE_PRICE_PR_SR",
	LockTransactionDateInventory:     "LOCK_TRANSACTION_DATE_INVENTORY",
	IsPurchasingApproval:             "IS_PURCHASING_APPROVAL",
	IsCompanyPRApplyPriceMoreThanOne: "IS_COMPANY_PR_APLLY_PRICE_MORE_THAN_ONE",
}
var ConfigurationNameCAMS = TConfigurationNameCAMS{

	//General
	DatabaseVersion: "DATABASE_VERSION",
	//Form
	FormMDIChild:       "FORM_MDI_CHILD",
	RequiredFieldColor: "REQUIRED_FIELD_COLOR",
	//Format Setting
	ShortDateFormat:    "SHORT_DATE_FORMAT",
	DateSeparator:      "DATE_SEPARATOR",
	CurrencyFormat:     "CURRENCY_FORMAT",
	DecimalSeparator:   "DECIMAL_SEPARATOR",
	ThousandsSeparator: "THOUSANDS_SEPARATOR",
	//Grid
	ShowGridGroupByBox: "SHOW_GRID_GROUP_BY_BOX",
	GridCanEditCell:    "GRID_CAN_EDIT_CELL",
	LimitGrid:          "LIMIT_GRID",
	GridRowHeight:      "GRID_ROW_HEIGHT",
	//Report
	LogoWidth:            "LOGO_WIDTH",
	ReportHeaderAligment: "REPORT_HEADER_ALIGNMENT",
	//Default Variable
	DVSubDepartment: "DV_SUB_DEPARTMENT",
	//Reservation
	AutoGenerateCompanyCode: "AUTO_GENERATE_COMPANY_CODE",
	CompanyCodeDigit:        "COMPANY_CODE_DIGIT",
	//Other
	InsertServerIDOnNumber: "INSERT_SERVER_ID_ON__NUMBER",
	//Accounting
	IsAccrualBase: "IS_ACCRUAL_BASE",
	//Personal
	DRName:     "DR_NAME",
	DRPosition: "DR_POSITION",
	HMName:     "HM_NAME",
	HMPosition: "HM_POSITION",
	FHName:     "FH_NAME",
	FHPosition: "FH_POSITION",
	GCName:     "GC_NAME",
	GCPosition: "GC_POSITION",
	CCName:     "CC_NAME",
	CCPosition: "CC_POSITION",
	PMName:     "PM_NAME",
	PMPosition: "PM_POSITION",
	//Invoice
	InvoiceAPNote: "INVOICE_AP_NOTE",
	//Purchase Request Approval
	UserApproval1: "USER_APPROVAL1",
	UserApproval2: "USER_APPROVAL2",
	UserApproval3: "USER_APPROVAL3",
	//SR & PR Color
	DHColor: "DH_COLOR",
	CCColor: "CC_COLOR",
	FNColor: "FN_COLOR",
	RJColor: "RJ_COLOR",
	//All Store Stock & Store Stock
	LowStockColor: "LOW_STOCK_COLOR",
	//Inventory
	CostingMethod:                    "COSTING_METHOD",
	SynchronizePOAndReceive:          "SYNCHRONIZE_PO_AND_RECEIVE",
	DefaultShippingAddress:           "DEFAULT_SHIPPING_ADDRESS",
	ReceiveStockAPTwoDigitDecimal:    "RECEIVE_STOCK_AP_TWO_DIGIT_DECIMAL",
	ShowPurchasePricePRSR:            "SHOW_PURCHASE_PRICE_PR_SR",
	LockTransactionDateInventory:     "LOCK_TRANSACTION_DATE_INVENTORY",
	IsPurchasingApproval:             "IS_PURCHASING_APPROVAL",
	IsCompanyPRApplyPriceMoreThanOne: "IS_COMPANY_PR_APLLY_PRICE_MORE_THAN_ONE",
	//Company
	CompanyContactPersonRequired: "COMPANY_CONTACT_PERSON_REQUIRED",
	CompanyStreetRequired:        "COMPANY_STREET_REQUIRED",
	CompanyCityRequired:          "COMPANY_CITY_REQUIRED",
	CompanyCountryRequired:       "COMPANY_COUNTRY_REQUIRED",
	CompanyStateRequired:         "COMPANY_STATE_REQUIRED",
	CompanyPostalCodeRequired:    "COMPANY_POSTAL_CODE_REQUIRED",
	CompanyPhone1Required:        "COMPANY_PHONE1_REQUIRED",
	//Global Journal Accounting
	JAAPSupplier:         "AP_SUPPLIER",
	JAPurchasingDiscount: "PURCHASING_DISCOUNT",
	JAPurchasingTax:      "PURCHASING_TAX",
	JAPurchasingShipping: "PURCHASING_SHIPPING",
	JAIncomeReturnStock:  "INCOME_RETURN_STOCK",
	JAExpenseReturnStock: "EXPENSE_RETURN_STOCK",
	JASGInventory:        "INVENTORY",
	JASGFixedAsset:       "FIXED_ASSET",
	JASGAccmDepreciation: "ACCUMULATED_DEPRECIATION",
	JASGAccountPayable:   "ACCOUNT_PAYABLE",
	//Global Sub Department
	SDAccounting:          "SUB_DEPARTMENT_ACCOUNTING",
	CompanyTypeSupplier:   "COMPANY_TYPE_SUPPLIER",
	CompanyTypeExpedition: "COMPANY_TYPE_EXPEDITION",
	//Report Form
	PurchaseRequest:  "PURCHASE_REQUEST",
	PurchaseOrder:    "PURCHASE_ORDER",
	ReceiveStock:     "RECEIVE_STOCK",
	StoreRequitition: "STORE_REQUISITION",
	StockTransfer:    "STOCK_TRANSFER",
	Costing:          "COSTING",
	FAPurchaseOrder:  "FA_PURCHASE_ORDER",
	FAReceive:        "FA_RECEIVE",
	//Report Template
	RTJournalVoucher:                 "JOURNAL_VOUCHER",
	RTJournalVoucherPaymentForm:      "JOURNAL_VOUCHER_PAYMENT_FORM",
	RTJournalVoucherReceiveForm:      "JOURNAL_VOUCHER_RECEIVE_FORM",
	RTIncomeStatement:                "INCOME_STATEMENT",
	RTInventoryReconciliation:        "INVENTORY_RECONCILIATION",
	RTDailyInventoryReconciliation:   "DAILY_INVENTORY_RECONCILIATION",
	RTMonthlyInventoryReconciliation: "MONTHLY_INVENTORY_RECONCILIATION",
	//Notification
	NotificationActive:   "NOTIFICATION_ACTIVE",
	NotificationDuration: "NOTIFICATION_DURATION",
	NotificationDelay:    "NOTIFICATION_DELAY",
	NotificationWidth:    "NOTIFICATION_WIDTH",
	NotificationHeight:   "NOTIFICATION_HEIGHT",
	NotificationStyle:    "NOTIFICATION_STYLE",
	NotificationPosition: "NOTIFICATION_POSITION"}
var UserAccessType = TUserAccessType{
	Form:              "F",
	Report:            "R",
	Special:           "S",
	Keylock:           "K",
	Reservation:       "V",
	Deposit:           "D",
	InHouse:           "I",
	WalkIn:            "W",
	Folio:             "L",
	FolioHistory:      "H",
	FloorPlan:         "P",
	Company:           "C",
	Invoice:           "O",
	MemberVoucherGift: "M",
	PreviewReport:     "T",
	PaymentByAPAR:     "A"}
var ReportAccessType = TReportAccessType{
	Form:       "F",
	Preview:    "V",
	Pos:        "P",
	FrontDesk:  "H",
	Banquet:    "B",
	Asset:      "I",
	Accounting: "A",
}
var LogUserAction = TLogUserAction{
	//Reservation
	InsertReservation:                     10101,
	VoidReservation:                       10102,
	CheckInReservation:                    10103,
	CancelReservation:                     10104,
	NoShowReservation:                     10105,
	InsertDeposit:                         10106,
	VoidDeposit:                           10107,
	CorrectDeposit:                        10108,
	RefundDeposit:                         10109,
	TransferDeposit:                       10110,
	InsertReservationScheduledRate:        10151,
	UpdateReservationScheduledRate:        10152,
	DeleteReservationScheduledRate:        10153,
	InsertReservationExtraCharge:          10154,
	UpdateReservationExtraCharge:          10155,
	DeleteReservationExtraCharge:          10156,
	InsertReservationExtraChargeBreakdown: 10157,
	UpdateReservationExtraChargeBreakdown: 10158,
	DeleteReservationExtraChargeBreakdown: 10159,
	//Update Reservation Stay Information
	URSIArrival:         10201,
	URSINights:          10202,
	URSIDeparture:       10203,
	URSIAdult:           10204,
	URSIChild:           10205,
	URSIRoomType:        10206,
	URSIRoom:            10207,
	URSIRoomRate:        10208,
	URSIBusinessSource:  10209,
	URSICommissionType:  10210,
	URSICommissionValue: 10211,
	URSIWeekdayRate:     10212,
	URSIWeekendRate:     10213,
	URSIDiscount:        10214,
	URSIPaymentType:     10215,
	URSIMarket:          10216,
	URSIBillInstruction: 10217,
	URSICurrency:        10218,
	URSIExchangeRate:    10219,
	//Update Reservation Personal Information
	URPIMember:        10300,
	URPITitle:         10301,
	URPIFullName:      10302,
	URPIReservationBy: 10303,
	URPIStreet:        10304,
	URPICity:          10305,
	URPICountry:       10306,
	URPIState:         10307,
	URPIPostCode:      10308,
	URPIPhone1:        10309,
	URPIPhone2:        10310,
	URPIFax:           10311,
	URPIEmail:         10312,
	URPIWebsite:       10313,
	URPICompany:       10314,
	URPIGuestType:     10315,
	URPIIDCardType:    10316,
	URPIIDCardNumber:  10317,
	URPIBirthdayPlace: 10318,
	URPIBirthdate:     10319,
	//Update Reservation General Information
	URGIPurposeOf:       10410,
	URGIGroup:           10401,
	URGIDocumentNumber:  10402,
	URGIFlightNumber:    10403,
	URGIFlightArrival:   10404,
	URGIFlightDeparture: 10405,
	URGINotes:           10406,
	URGIHKNotes:         10407,
	URGIMarketing:       10408,
	URGITAVoucherNumber: 10409,
	//Group Reservation
	InsertGroupReservation: 20101,
	UpdateGroupReservation: 20102,
	DeleteGroupReservation: 20103,
	//m9	WalkIn:                          30101,
	InsertMasterFolio:               30102,
	InsertDeskFolio:                 30103,
	InsertFolioScheduledRate:        30104,
	UpdateFolioScheduledRate:        30105,
	DeleteFolioScheduledRate:        30106,
	InsertFolioExtraCharge:          30107,
	UpdateFolioExtraCharge:          30108,
	DeleteFolioExtraCharge:          30109,
	InsertFolioExtraChargeBreakdown: 30110,
	UpdateFolioExtraChargeBreakdown: 30111,
	DeleteFolioExtraChargeBreakdown: 30112,
	InsertTransaction:               30113,
	VoidTransaction:                 30114,
	CorrectTransaction:              30115,
	TransferTransaction:             30116,
	RoutingFolio:                    30117,
	ReturnTransfer:                  30118,
	RemoveRouting:                   30119,
	MoveRoom:                        30120,
	SwitchRoom:                      30121,
	InsertMessage:                   30122,
	UpdateMessage:                   30123,
	DeleteMessage:                   30124,
	MarkAsDeliveredMessage:          30125,
	MarkAsUndeliveredMessage:        30126,
	InsertToDo:                      30127,
	UpdateToDo:                      30128,
	DeleteToDo:                      30129,
	MarkAsDoneToDo:                  30130,
	MarkAsNotDoneToDo:               30131,
	CancelCheckIn:                   30132,
	CheckOutFolio:                   30133,
	CancelCheckOut:                  30134,
	ComplimentGuest:                 30135,
	HouseUseGuest:                   30136,
	DefaultGuest:                    30137,
	UpdateFolioVoucher:              30138,
	InsertVoucherPayment:            30139,
	IssuedCard:                      30161,
	ReplaceCard:                     30162,
	EraseCard:                       30163,
	DeactivateWithoutCard:           30164,
	ForceEraseCard:                  30165,
	//Update Folio Stay Information
	UFSIArrival:         30201,
	UFSINights:          30202,
	UFSIDeparture:       30203,
	UFSIAdult:           30204,
	UFSIChild:           30205,
	UFSIRoomType:        30206,
	UFSIRoom:            30207,
	UFSIRoomRate:        30208,
	UFSIBusinessSource:  30209,
	UFSICommissionType:  30210,
	UFSICommissionValue: 30211,
	UFSIWeekdayRate:     30212,
	UFSIWeekendRate:     30213,
	UFSIDiscount:        30214,
	UFSIPaymentType:     30215,
	UFSIMarket:          30216,
	UFSIBillInstruction: 30217,
	UFSICurrency:        30218,
	UFSIExchangeRate:    30219,
	//Update Folio Personal Information
	UFPIMember:        30300,
	UFPITitle:         30301,
	UFPIFullName:      30302,
	UFPIStreet:        30303,
	UFPICity:          30304,
	UFPICountry:       30305,
	UFPIState:         30306,
	UFPIPostCode:      30307,
	UFPIPhone1:        30308,
	UFPIPhone2:        30309,
	UFPIFax:           30310,
	UFPIEmail:         30311,
	UFPIWebsite:       30312,
	UFPICompany:       30313,
	UFPIGuestType:     30314,
	UFPIIDCardType:    30315,
	UFPIIDCardNumber:  30316,
	UFPIBirthdayPlace: 30317,
	UFPIBirthdate:     30318,
	//Update Folio General Information
	UFGIPurposeOf:       30410,
	UFGIGroup:           30401,
	UFGIDocumentNumber:  30402,
	UFGIFlightNumber:    30403,
	UFGIFlightArrival:   30404,
	UFGIFlightDeparture: 30405,
	UFGINotes:           30406,
	UFGIHKNotes:         30407,
	UFGIMarketing:       30408,
	UFGITAVoucherNumber: 30409,
	//Member Voucher Gift
	MVInsertMember:      30501,
	MVUpdateMember:      30502,
	MVDeleteMember:      30503,
	MVRedeemMemberPoint: 30504,
	MVInsertVoucher:     30511,
	MVDeleteVoucher:     30512,
	MVApproveVoucher:    30513,
	MVNotApproveVoucher: 30514,
	MVUnapproveVoucher:  30519,
	MVSoldVoucher:       30515,
	MVRedeemVoucher:     30516,
	MVComplimentVoucher: 30517,
	MVUnsoldVoucher:     30518,
	//Sales Activity
	InsertLead:         30601,
	UpdateLead:         30602,
	VoidLead:           30603,
	InsertProposal:     30604,
	UpdateProposal:     30605,
	VoidProposal:       30606,
	InsertTask:         30607,
	UpdateTask:         30608,
	VoidTask:           30609,
	InsertSendReminder: 30610,
	UpdateSendReminder: 30611,
	DeleteSendReminder: 30612,
	InsertNotes:        30613,
	UpdateNotes:        30614,
	DeleteNotes:        30615,
	InsertActivityLog:  30616,
	UpdateActivityLog:  30617,
	DeleteActivityLog:  30618,
	InsertContact:      30619,
	UpdateContact:      30620,
	DeleteContact:      30621,
	//House Keeping
	ReadyRoom:                 40101,
	CleanRoom:                 40102,
	DirtyRoom:                 40103,
	OutOfOrder:                40104,
	OfficeUse:                 40105,
	UnderConstruction:         40107,
	DontDisturb:               40108,
	DoubleLock:                40109,
	SleepOut:                  40110,
	PleaseClean:               40114,
	InsertRoomCosting:         40111,
	DeleteRoomCostingTransfer: 40112,
	DeleteRoomCosting:         40113,
	//Report
	InsertReportTemplate:     50101,
	UpdateReportTemplate:     50102,
	DeleteReportTemplate:     50103,
	SetDefaultReportTemplate: 50104,
	//Configuration
	CompanyInformation: 60101,
	Configuration:      60102,
	InsertMasterData:   60103,
	UpdateMasterData:   60104,
	DeleteMasterData:   60105,
	ModifyFloorPlan:    60106,
	//Accounting Tool
	InsertReceive:         43101,
	UpdateReceive:         43102,
	DeleteReceive:         43103,
	InsertPayment:         43104,
	UpdatePayment:         43105,
	DeletePayment:         43106,
	InsertReceipt:         43107,
	UpdateReceipt:         43108,
	DeleteReceipt:         43109,
	InsertIncomeBudget:    43113,
	UpdateIncomeBudget:    43114,
	DeleteIncomeBudget:    43115,
	InsertBudgetStatistic: 43119,
	UpdateBudgetStatistic: 43120,
	DeleteBudgetStatistic: 43121,
	//Account Payable
	InsertAPRefundDepositPayment: 43207,
	UpdateAPRefundDepositPayment: 43208,
	DeleteAPRefundDepositPayment: 4209,
	InsertAPCommissionPayment:    43210,
	UpdateAPCommissionPayment:    43211,
	DeleteAPCommissionPayment:    43212,
	//Account Receivable
	InsertInvoiceCityLedger:        43307,
	UpdateInvoiceCityLedger:        43308,
	DeleteInvoiceCityLedger:        43309,
	InsertPaymentInvoiceCityLedger: 43310,
	UpdatePaymentInvoiceCityLedger: 43311,
	DeletePaymentInvoiceCityLedger: 43312,
	InsertBankReconciliation:       43313,
	UpdateBankReconciliation:       43314,
	DeleteBankReconciliation:       43315,
	PrintInvoice:                   43316,
	//User Setting
	InsertUser:      70101,
	UpdateUser:      70102,
	DeleteUser:      70103,
	InsertUserGroup: 70104,
	UpdateUserGroup: 70105,
	DeleteUserGroup: 70106,
	//Database
	BackupDatabase:     80101,
	RestoreDatabase:    80102,
	OptimizingDatabase: 80103,
	//Login
	Login:          90101,
	Logout:         90102,
	ChangePassword: 90103,
	LoginDenied:    90104,
	//Reservation
	InsertReservationX:                     true,
	VoidReservationX:                       true,
	CheckInReservationX:                    true,
	CancelReservationX:                     true,
	NoShowReservationX:                     true,
	InsertDepositX:                         true,
	VoidDepositX:                           true,
	RefundDepositX:                         true,
	TransferDepositX:                       true,
	InsertReservationScheduledRateX:        true,
	UpdateReservationScheduledRateX:        true,
	DeleteReservationScheduledRateX:        true,
	InsertReservationExtraChargeX:          true,
	UpdateReservationExtraChargeX:          true,
	DeleteReservationExtraChargeX:          true,
	InsertReservationExtraChargeBreakdownX: true,
	UpdateReservationExtraChargeBreakdownX: true,
	DeleteReservationExtraChargeBreakdownX: true,
	//Update Reservation Stay Information
	URSIArrivalX:         true,
	URSINightsX:          true,
	URSIDepartureX:       true,
	URSIAdultX:           true,
	URSIChildX:           true,
	URSIRoomTypeX:        true,
	URSIRoomX:            true,
	URSIRoomRateX:        true,
	URSIBusinessSourceX:  true,
	URSICommissionTypeX:  true,
	URSICommissionValueX: true,
	URSIWeekdayRateX:     true,
	URSIWeekendRateX:     true,
	URSIDiscountX:        true,
	URSIPaymentTypeX:     true,
	URSIMarketX:          true,
	URSIBillInstructionX: true,
	URSICurrencyX:        true,
	URSIExchangeRateX:    true,
	//Update Reservation Personal Information
	URPIMemberX:        true,
	URPITitleX:         true,
	URPIFullNameX:      true,
	URPIReservationByX: true,
	URPIStreetX:        true,
	URPICityX:          true,
	URPICountryX:       true,
	URPIStateX:         true,
	URPIPostCodeX:      true,
	URPIPhone1X:        true,
	URPIPhone2X:        true,
	URPIFaxX:           true,
	URPIEmailX:         true,
	URPIWebsiteX:       true,
	URPICompanyX:       true,
	URPIGuestTypeX:     true,
	URPIIDCardTypeX:    true,
	URPIIDCardNumberX:  true,
	URPIBirthdayPlaceX: true,
	URPIBirthdateX:     true,
	//Update Reservation General Information
	URGIPurposeOfX:       true,
	URGIGroupX:           true,
	URGIDocumentNumberX:  true,
	URGIFlightNumberX:    true,
	URGIFlightArrivalX:   true,
	URGIFlightDepartureX: true,
	URGINotesX:           true,
	URGIHKNotesX:         true,
	URGIMarketingX:       true,
	URGITAVoucherNumberX: true,
	//Group Reservation
	InsertGroupReservationX: true,
	UpdateGroupReservationX: true,
	DeleteGroupReservationX: true,
	//Folio
	WalkInX:                          true,
	InsertMasterFolioX:               true,
	InsertDeskFolioX:                 true,
	InsertFolioScheduledRateX:        true,
	UpdateFolioScheduledRateX:        true,
	DeleteFolioScheduledRateX:        true,
	InsertFolioExtraChargeX:          true,
	UpdateFolioExtraChargeX:          true,
	DeleteFolioExtraChargeX:          true,
	InsertFolioExtraChargeBreakdownX: true,
	UpdateFolioExtraChargeBreakdownX: true,
	DeleteFolioExtraChargeBreakdownX: true,
	InsertTransactionX:               true,
	VoidTransactionX:                 true,
	TransferTransactionX:             true,
	RoutingFolioX:                    true,
	ReturnTransferX:                  true,
	RemoveRoutingX:                   true,
	MoveRoomX:                        true,
	SwitchRoomX:                      true,
	InsertMessageX:                   true,
	UpdateMessageX:                   true,
	DeleteMessageX:                   true,
	MarkAsDeliveredMessageX:          true,
	MarkAsUndeliveredMessageX:        true,
	InsertToDoX:                      true,
	UpdateToDoX:                      true,
	DeleteToDoX:                      true,
	MarkAsDoneToDoX:                  true,
	MarkAsNotDoneToDoX:               true,
	CancelCheckInX:                   true,
	CheckOutFolioX:                   true,
	CancelCheckOutX:                  true,
	ComplimentGuestX:                 true,
	HouseUseGuestX:                   true,
	DefaultGuestX:                    true,
	UpdateFolioVoucherX:              true,
	InsertVoucherPaymentX:            true,
	IssuedCardX:                      true,
	ReplaceCardX:                     true,
	EraseCardX:                       true,
	DeactivateWithoutCardX:           true,
	ForceEraseCardX:                  true,
	//Update Folio Stay Information
	UFSIArrivalX:         true,
	UFSINightsX:          true,
	UFSIDepartureX:       true,
	UFSIAdultX:           true,
	UFSIChildX:           true,
	UFSIRoomTypeX:        true,
	UFSIRoomX:            true,
	UFSIRoomRateX:        true,
	UFSIBusinessSourceX:  true,
	UFSICommissionTypeX:  true,
	UFSICommissionValueX: true,
	UFSIWeekdayRateX:     true,
	UFSIWeekendRateX:     true,
	UFSIDiscountX:        true,
	UFSIPaymentTypeX:     true,
	UFSIMarketX:          true,
	UFSIBillInstructionX: true,
	UFSICurrencyX:        true,
	UFSIExchangeRateX:    true,
	//Update Folio Personal Information
	UFPIMemberX:        true,
	UFPITitleX:         true,
	UFPIFullNameX:      true,
	UFPIStreetX:        true,
	UFPICityX:          true,
	UFPICountryX:       true,
	UFPIStateX:         true,
	UFPIPostCodeX:      true,
	UFPIPhone1X:        true,
	UFPIPhone2X:        true,
	UFPIFaxX:           true,
	UFPIEmailX:         true,
	UFPIWebsiteX:       true,
	UFPICompanyX:       true,
	UFPIGuestTypeX:     true,
	UFPIIDCardTypeX:    true,
	UFPIIDCardNumberX:  true,
	UFPIBirthdayPlaceX: true,
	UFPIBirthdateX:     true,
	//Update Folio General Information
	UFGIPurposeOfX:       true,
	UFGIGroupX:           true,
	UFGIDocumentNumberX:  true,
	UFGIFlightNumberX:    true,
	UFGIFlightArrivalX:   true,
	UFGIFlightDepartureX: true,
	UFGINotesX:           true,
	UFGIHKNotesX:         true,
	UFGIMarketingX:       true,
	UFGITAVoucherNumberX: true,
	//Member Voucher Gift
	MVInsertMemberX:      true,
	MVUpdateMemberX:      true,
	MVDeleteMemberX:      true,
	MVRedeemMemberPointX: true,
	MVInsertVoucherX:     true,
	MVDeleteVoucherX:     true,
	MVApproveVoucherX:    true,
	MVNotApproveVoucherX: true,
	MVUnapproveVoucherX:  true,
	MVSoldVoucherX:       true,
	MVRedeemVoucherX:     true,
	MVComplimentVoucherX: true,
	MVUnsoldVoucherX:     true,

	//Sales Activity
	InsertLeadX:         true,
	UpdateLeadX:         true,
	VoidLeadX:           true,
	InsertProposalX:     true,
	UpdateProposalX:     true,
	VoidProposalX:       true,
	InsertTaskX:         true,
	UpdateTaskX:         true,
	VoidTaskX:           true,
	InsertSendReminderX: true,
	UpdateSendReminderX: true,
	DeleteSendReminderX: true,
	InsertNotesX:        true,
	UpdateNotesX:        true,
	DeleteNotesX:        true,
	InsertActivityLogX:  true,
	UpdateActivityLogX:  true,
	DeleteActivityLogX:  true,
	InsertContactX:      true,
	UpdateContactX:      true,
	DeleteContactX:      true,

	//House Keeping,
	ReadyRoomX:                 true,
	CleanRoomX:                 true,
	DirtyRoomX:                 true,
	OutOfOrderX:                true,
	OfficeUseX:                 true,
	UnderConstructionX:         true,
	DontDisturbX:               true,
	DoubleLockX:                true,
	SleepOutX:                  true,
	PleaseCleanX:               true,
	InsertRoomCostingX:         true,
	DeleteRoomCostingTransferX: true,
	DeleteRoomCostingX:         true,
	//Report

	//User Setting
	InsertReportTemplateX:     true,
	UpdateReportTemplateX:     true,
	DeleteReportTemplateX:     true,
	SetDefaultReportTemplateX: true,
	//Configuration
	CompanyInformationX: true,
	ConfigurationX:      true,
	InsertMasterDataX:   true,
	UpdateMasterDataX:   true,
	DeleteMasterDataX:   true,
	ModifyFloorPlanX:    true,
	//Accounting Tool
	InsertReceiptX:         true,
	UpdateReceiptX:         true,
	DeleteReceiptX:         true,
	InsertIncomeBudgetX:    true,
	UpdateIncomeBudgetX:    true,
	DeleteIncomeBudgetX:    true,
	InsertBudgetStatisticX: true,
	UpdateBudgetStatisticX: true,
	DeleteBudgetStatisticX: true,
	//Account Payable
	InsertAPRefundDepositPaymentX: true,
	UpdateAPRefundDepositPaymentX: true,
	DeleteAPRefundDepositPaymentX: true,
	InsertAPCommissionPaymentX:    true,
	UpdateAPCommissionPaymentX:    true,
	DeleteAPCommissionPaymentX:    true,
	//Account Receivable
	InsertInvoiceCityLedgerX:        true,
	UpdateInvoiceCityLedgerX:        true,
	DeleteInvoiceCityLedgerX:        true,
	InsertPaymentInvoiceCityLedgerX: true,
	UpdatePaymentInvoiceCityLedgerX: true,
	DeletePaymentInvoiceCityLedgerX: true,
	InsertBankReconciliationX:       true,
	UpdateBankReconciliationX:       true,
	DeleteBankReconciliationX:       true,
	PrintInvoiceX:                   true,
	//User Setting
	InsertUserX:      true,
	UpdateUserX:      true,
	DeleteUserX:      true,
	InsertUserGroupX: true,
	UpdateUserGroupX: true,
	DeleteUserGroupX: true,
	//Database
	BackupDatabaseX:     true,
	RestoreDatabaseX:    true,
	OptimizingDatabaseX: true,
	//Login
	LoginX:          true,
	LogoutX:         true,
	ChangePasswordX: true,
	LoginDeniedX:    true}

var LogUserActionPOS = TLogUserActionPOS{
	//Customer
	InsertCustomer: 41115,
	UpdateCustomer: 41116,
	DeleteCustomer: 41117,
	//Reservation
	InsertReservationPOS: 41118,
	UpdateReservationPOS: 41119,
	ChekInReservationPOS: 41120,
	CancelReservationPOS: 41121,
	NoShowReservationPOS: 41122,
	VoidReservationPOS:   41123,
	//POS Terminal and Table View
	InsertCaptainOrder:             41101,
	UpdateCaptainOrder:             41102,
	TransferCaptainOrder:           41103,
	CancelCaptainOrder:             41104,
	ChangeQuantity:                 41105,
	UpdateRemark:                   41106,
	Discount:                       41107,
	OverridePrice:                  41108,
	ModifyPriceZero:                41109,
	ModifyPriceRemoveTaxAndService: 41110,
	RemoveItem:                     41111,
	FinishSale:                     41112,
	ModifyTableView:                41113,
	VoidCheck:                      41114,
	//Member Voucher Gift
	MVInsertMember:      30501,
	MVUpdateMember:      30502,
	MVDeleteMember:      30503,
	MVRedeemMemberPoint: 30504,
	MVInsertVoucher:     30511,
	MVDeleteVoucher:     30512,
	MVApproveVoucher:    30513,
	MVNotApproveVoucher: 30514,
	MVUnapproveVoucher:  30519,
	MVSoldVoucher:       30515,
	MVRedeemVoucher:     30516,
	MVComplimentVoucher: 30517,
	MVUnsoldVoucher:     30518,
	//Budget
	InsertFBBudget: 43126,
	UpdateFBBudget: 43127,
	DeleteFBBudget: 43128,
	//Report
	InsertReportTemplate:     50101,
	UpdateReportTemplate:     50102,
	DeleteReportTemplate:     50103,
	SetDefaultReportTemplate: 50104,
	//Configuration
	CompanyInformation: 60101,
	Configuration:      60102,
	InsertMasterData:   60103,
	UpdateMasterData:   60104,
	DeleteMasterData:   60105,
	//User Setting
	InsertUser:      70101,
	UpdateUser:      70102,
	DeleteUser:      70103,
	InsertUserGroup: 70104,
	UpdateUserGroup: 70105,
	DeleteUserGroup: 70106,
	//Database
	BackupDatabase:     80101,
	RestoreDatabase:    80102,
	OptimizingDatabase: 80103,
	//Login
	Login:          90101,
	Logout:         90102,
	ChangePassword: 90103,
	LoginDenied:    90104,
	//CHS Parameter Dummy
	CancelReservation: 10104,
	NoShowReservation: 10105,
	CleanRoom:         40102,
	DirtyRoom:         40103,
	//POS Terminal and Table View
	InsertCaptainOrderX:             true,
	UpdateCaptainOrderX:             true,
	TransferCaptainOrderX:           true,
	CancelCaptainOrderX:             true,
	ChangeQuantityX:                 true,
	UpdateRemarkX:                   true,
	DiscountX:                       true,
	OverridePriceX:                  true,
	ModifyPriceZeroX:                true,
	ModifyPriceRemoveTaxAndServiceX: true,
	RemoveItemX:                     true,
	FinishSaleX:                     true,
	ModifyTableViewX:                true,
	VoidCheckX:                      true,
	//Member Voucher Gift
	MVInsertMemberX:      true,
	MVUpdateMemberX:      true,
	MVDeleteMemberX:      true,
	MVRedeemMemberPointX: true,
	MVInsertVoucherX:     true,
	MVDeleteVoucherX:     true,
	MVApproveVoucherX:    true,
	MVNotApproveVoucherX: true,
	MVUnapproveVoucherX:  true,
	MVSoldVoucherX:       true,
	MVRedeemVoucherX:     true,
	MVComplimentVoucherX: true,
	MVUnsoldVoucherX:     true,
	//Budget
	InsertFBBudgetX: true,
	UpdateFBBudgetX: true,
	DeleteFBBudgetX: true,
	//Report
	InsertReportTemplateX:     true,
	UpdateReportTemplateX:     true,
	DeleteReportTemplateX:     true,
	SetDefaultReportTemplateX: true,
	//Configuration
	CompanyInformationX: true,
	ConfigurationX:      true,
	InsertMasterDataX:   true,
	UpdateMasterDataX:   true,
	DeleteMasterDataX:   true,
	//User Setting
	InsertUserX:      true,
	UpdateUserX:      true,
	DeleteUserX:      true,
	InsertUserGroupX: true,
	UpdateUserGroupX: true,
	DeleteUserGroupX: true,
	//Database
	BackupDatabaseX:     true,
	RestoreDatabaseX:    true,
	OptimizingDatabaseX: true,
	//Login
	LoginX:          true,
	LogoutX:         true,
	ChangePasswordX: true,
	LoginDeniedX:    true,
	//CHS Parameter Dummy
	CancelReservationX: true,
	NoShowReservationX: true,
	CleanRoomX:         true,
	DirtyRoomX:         true}
var UserGroup = TUserGroup{
	System:     "SYSTEM",
	SuperAdmin: "SUPERADMIN",
	Admin:      "ADMIN"}
var UserFormAccessOrder = TUserFormAccessOrder{
	Summary:               1,
	FloorPlan:             2,
	RoomAvailability:      3,
	RoomTypeAvailability:  4,
	RoomAllotment:         48,
	GuestProfile:          5,
	GuestGroup:            6,
	Reservation:           7,
	WalkIn:                8,
	GuestInHouse:          9,
	MasterFolio:           10,
	DeskFolio:             11,
	FoliosAndTransaction:  12,
	FolioHistory:          13,
	HouseKeeping:          14,
	RoomCosting:           15,
	LostAndFound:          16,
	CashierReport:         17,
	GlobalPostTransaction: 18,
	AutoPostTransaction:   19,
	DayendClose:           20,
	Company:               21,
	Package:               22,
	RoomRate:              23,
	EventList:             24,
	PhoneBook:             25,
	Member:                44,
	Voucher:               45,
	Gift:                  46,
	APRefundDeposit:       26,
	APCommission:          27,
	ARCityLedger:          28,
	ARCityLedgerInvoice:   29,
	Receipt:               30,
	BankTransaction:       31,
	BankReconciliation:    32,
	Cheque:                33,
	ExportJournal:         34,
	IncomeBudget:          35,
	BudgetStatistic:       36,
	Report:                37,
	DataAnalysis:          51,
	Configuration:         38,
	UserSetting:           39,
	OneTimePassword:       50,
	UserActivityLog:       40,
	PABXSMDRViewer:        49,
	BackupDatabase:        41,
	RestoreDatabase:       42,
	OptimizeData:          43,
	Notification:          47,
	BreakfastControl:      52,
	CompetitorData:        53,
	GuestLoanItem:         54,
	SalesActivity:         55,
	DynamicRate:           56,
	RoomRateLastDeal:      57}
var LogUserActionCAMS = TLogUserActionCAMS{
	//Inventory
	InsertPurchaseRequest: 46125,
	UpdatePurchaseRequest: 46126,
	DeletePurchaseRequest: 46127,
	InsertPurchaseOrder:   46101,
	UpdatePurchaseOrder:   46102,
	DeletePurchaseOrder:   46103,
	InsertReceiveStock:    46104,
	UpdateReceiveStock:    46105,
	DeleteReceiveStock:    46106,
	InsertStockTransfer:   46107,
	UpdateStockTransfer:   46108,
	DeleteStockTransfer:   46109,
	InsertCosting:         46110,
	UpdateCosting:         46111,
	DeleteCosting:         46112,
	InsertProduction:      46116,
	UpdateProduction:      46117,
	DeleteProduction:      46118,
	InsertCostRecipe:      46119,
	UpdateCostRecipe:      46120,
	DeleteCostRecipe:      46121,
	InsertReturnStock:     46122,
	UpdateReturnStock:     46123,
	DeleteReturnStock:     46124,
	InsertStockOpname:     46113,
	DeleteStockOpname:     46114,
	SetActiveStore:        46115,
	//Fixed Asset
	InsertFAPurchaseOrder: 46201,
	UpdateFAPurchaseOrder: 46202,
	DeleteFAPurchaseOrder: 46203,
	InsertFAReceiveStock:  46204,
	UpdateFAReceiveStock:  46205,
	DeleteFAReceiveStock:  46206,
	InsertFixedAssetList:  46207,
	UpdateFixedAssetList:  46208,
	DeleteFixedAssetList:  46209,
	InsertDepreciation:    46210,
	DeleteDepreciation:    46211,
	//Report
	InsertReportTemplate:     50101,
	UpdateReportTemplate:     50102,
	DeleteReportTemplate:     50103,
	SetDefaultReportTemplate: 50104,
	//Configuration
	CompanyInformation: 60101,
	Configuration:      60102,
	InsertMasterData:   60103,
	UpdateMasterData:   60104,
	DeleteMasterData:   60105,
	//User Setting
	InsertUser:      70101,
	UpdateUser:      70102,
	DeleteUser:      70103,
	InsertUserGroup: 70104,
	UpdateUserGroup: 70105,
	DeleteUserGroup: 70106,
	//Database
	BackupDatabase:     80101,
	RestoreDatabase:    80102,
	OptimizingDatabase: 80103,
	//Login
	Login:          90101,
	Logout:         90102,
	ChangePassword: 90103,
	LoginDenied:    90104,
	//Inventory
	InsertPurchaseRequestX: true,
	UpdatePurchaseRequestX: true,
	DeletePurchaseRequestX: true,
	InsertPurchaseOrderX:   true,
	UpdatePurchaseOrderX:   true,
	DeletePurchaseOrderX:   true,
	InsertReceiveStockX:    true,
	UpdateReceiveStockX:    true,
	DeleteReceiveStockX:    true,
	InsertStockTransferX:   true,
	UpdateStockTransferX:   true,
	DeleteStockTransferX:   true,
	InsertCostingX:         true,
	UpdateCostingX:         true,
	DeleteCostingX:         true,
	InsertProductionX:      true,
	UpdateProductionX:      true,
	DeleteProductionX:      true,
	InsertCostRecipeX:      true,
	UpdateCostRecipeX:      true,
	DeleteCostRecipeX:      true,
	InsertReturnStockX:     true,
	UpdateReturnStockX:     true,
	DeleteReturnStockX:     true,
	InsertStockOpnameX:     true,
	DeleteStockOpnameX:     true,
	SetActiveStoreX:        true,
	//Fixed Asset
	InsertFAPurchaseOrderX: true,
	UpdateFAPurchaseOrderX: true,
	DeleteFAPurchaseOrderX: true,
	InsertFAReceiveStockX:  true,
	UpdateFAReceiveStockX:  true,
	DeleteFAReceiveStockX:  true,
	InsertFixedAssetListX:  true,
	UpdateFixedAssetListX:  true,
	DeleteFixedAssetListX:  true,
	InsertDepreciationX:    true,
	DeleteDepreciationX:    true,
	//User Setting
	InsertReportTemplateX:     true,
	UpdateReportTemplateX:     true,
	DeleteReportTemplateX:     true,
	SetDefaultReportTemplateX: true,
	//Configuration
	CompanyInformationX: true,
	ConfigurationX:      true,
	InsertMasterDataX:   true,
	UpdateMasterDataX:   true,
	DeleteMasterDataX:   true,
	//User Setting
	InsertUserX:      true,
	UpdateUserX:      true,
	DeleteUserX:      true,
	InsertUserGroupX: true,
	UpdateUserGroupX: true,
	DeleteUserGroupX: true,
	//Database
	BackupDatabaseX:     true,
	RestoreDatabaseX:    true,
	OptimizingDatabaseX: true,
	//Login
	LoginX:          true,
	LogoutX:         true,
	ChangePasswordX: true,
	LoginDeniedX:    true}
var LogUserActionCAS = TLogUserActionCAS{
	//Accounting Tool
	InsertReceive:         43101,
	UpdateReceive:         43102,
	DeleteReceive:         43103,
	InsertPayment:         43104,
	UpdatePayment:         43105,
	DeletePayment:         43106,
	InsertReceipt:         43107,
	UpdateReceipt:         43108,
	DeleteReceipt:         43109,
	InsertJournal:         43110,
	UpdateJournal:         43111,
	DeleteJournal:         43112,
	InsertIncomeBudget:    43113,
	UpdateIncomeBudget:    43114,
	DeleteIncomeBudget:    43115,
	InsertExpenseBudget:   43116,
	UpdateExpenseBudget:   43117,
	DeleteExpenseBudget:   43118,
	InsertBudgetStatistic: 43119,
	UpdateBudgetStatistic: 43120,
	DeleteBudgetStatistic: 43121,
	CloseMonth:            43122,
	CancelCloseMonth:      43123,
	CloseYear:             43124,
	CancelCloseYear:       43125,
	//Account Payable
	InsertAccountPayable:         43201,
	UpdateAccountPayable:         43202,
	DeleteAccountPayable:         43203,
	InsertAccountPayablePayment:  43204,
	UpdateAccountPayablePayment:  43205,
	DeleteAccountPayablePayment:  43206,
	InsertAPRefundDepositPayment: 43207,
	UpdateAPRefundDepositPayment: 43208,
	DeleteAPRefundDepositPayment: 43209,
	InsertAPCommissionPayment:    43210,
	UpdateAPCommissionPayment:    43211,
	DeleteAPCommissionPayment:    43212,
	//Account Receivable
	InsertAccountReceivable:        43301,
	UpdateAccountReceivable:        43302,
	DeleteAccountReceivable:        43303,
	InsertAccountReceivablePayment: 43304,
	UpdateAccountReceivablePayment: 43305,
	DeleteAccountReceivablePayment: 43306,
	InsertInvoiceCityLedger:        43307,
	UpdateInvoiceCityLedger:        43308,
	DeleteInvoiceCityLedger:        43309,
	InsertPaymentInvoiceCityLedger: 43310,
	UpdatePaymentInvoiceCityLedger: 43311,
	DeletePaymentInvoiceCityLedger: 43312,
	InsertBankReconciliation:       43313,
	UpdateBankReconciliation:       43314,
	DeleteBankReconciliation:       43315,
	PrintInvoice:                   43316,
	//Report
	InsertReportTemplate:     50101,
	UpdateReportTemplate:     50102,
	DeleteReportTemplate:     50103,
	SetDefaultReportTemplate: 50104,
	//Configuration
	CompanyInformation: 60101,
	Configuration:      60102,
	InsertMasterData:   60103,
	UpdateMasterData:   60104,
	DeleteMasterData:   60105,
	//User Setting
	InsertUser:      70101,
	UpdateUser:      70102,
	DeleteUser:      70103,
	InsertUserGroup: 70104,
	UpdateUserGroup: 70105,
	DeleteUserGroup: 70106,
	//Database
	BackupDatabase:     80101,
	RestoreDatabase:    80102,
	OptimizingDatabase: 80103,
	//Login
	Login:          90101,
	Logout:         90102,
	ChangePassword: 90103,
	LoginDenied:    90104,
	//Accounting Tool
	InsertReceiveX:         true,
	UpdateReceiveX:         true,
	DeleteReceiveX:         true,
	InsertPaymentX:         true,
	UpdatePaymentX:         true,
	DeletePaymentX:         true,
	InsertReceiptX:         true,
	UpdateReceiptX:         true,
	DeleteReceiptX:         true,
	InsertJournalX:         true,
	UpdateJournalX:         true,
	DeleteJournalX:         true,
	InsertIncomeBudgetX:    true,
	UpdateIncomeBudgetX:    true,
	DeleteIncomeBudgetX:    true,
	InsertExpenseBudgetX:   true,
	UpdateExpenseBudgetX:   true,
	DeleteExpenseBudgetX:   true,
	InsertBudgetStatisticX: true,
	UpdateBudgetStatisticX: true,
	DeleteBudgetStatisticX: true,
	CloseMonthX:            true,
	CancelCloseMonthX:      true,
	CloseYearX:             true,
	CancelCloseYearX:       true,
	//Account Payable
	InsertAccountPayableX:         true,
	UpdateAccountPayableX:         true,
	DeleteAccountPayableX:         true,
	InsertAccountPayablePaymentX:  true,
	UpdateAccountPayablePaymentX:  true,
	DeleteAccountPayablePaymentX:  true,
	InsertAPRefundDepositPaymentX: true,
	UpdateAPRefundDepositPaymentX: true,
	DeleteAPRefundDepositPaymentX: true,
	InsertAPCommissionPaymentX:    true,
	UpdateAPCommissionPaymentX:    true,
	DeleteAPCommissionPaymentX:    true,
	//Account Receivable
	InsertAccountReceivableX:        true,
	UpdateAccountReceivableX:        true,
	DeleteAccountReceivableX:        true,
	InsertAccountReceivablePaymentX: true,
	UpdateAccountReceivablePaymentX: true,
	DeleteAccountReceivablePaymentX: true,
	InsertInvoiceCityLedgerX:        true,
	UpdateInvoiceCityLedgerX:        true,
	DeleteInvoiceCityLedgerX:        true,
	InsertPaymentInvoiceCityLedgerX: true,
	UpdatePaymentInvoiceCityLedgerX: true,
	DeletePaymentInvoiceCityLedgerX: true,
	InsertBankReconciliationX:       true,
	UpdateBankReconciliationX:       true,
	DeleteBankReconciliationX:       true,
	PrintInvoiceX:                   true,
	//Report
	InsertReportTemplateX:     true,
	UpdateReportTemplateX:     true,
	DeleteReportTemplateX:     true,
	SetDefaultReportTemplateX: true,
	//Configuration
	CompanyInformationX: true,
	ConfigurationX:      true,
	InsertMasterDataX:   true,
	UpdateMasterDataX:   true,
	DeleteMasterDataX:   true,
	//User Setting
	InsertUserX:      true,
	UpdateUserX:      true,
	DeleteUserX:      true,
	InsertUserGroupX: true,
	UpdateUserGroupX: true,
	DeleteUserGroupX: true,
	//Database
	BackupDatabaseX:     true,
	RestoreDatabaseX:    true,
	OptimizingDatabaseX: true,
	//Login
	LoginX:          true,
	LogoutX:         true,
	ChangePasswordX: true,
	LoginDeniedX:    true}
var UserReportAccessOrder = TUserReportAccessOrder{
	ReservationList:                       1,
	CancelledReservation:                  2,
	NoShowReservation:                     3,
	VoidedReservation:                     4,
	GroupReservation:                      5,
	ExpectedArrival:                       6,
	ArrivalList:                           7,
	SamedayReservation:                    8,
	AdvancedPaymentDeposit:                9,
	BalanceDeposit:                        10,
	WaitListReservation:                   11,
	CurrentInHouse:                        12,
	GuestInHouse:                          13,
	GuestInHousebyBusinessSource:          14,
	GuestInHousebyMarket:                  15,
	GuestInHousebyGuestType:               16,
	GuestInHousebyCountry:                 17,
	GuestInHousebyState:                   18,
	MasterFolio:                           19,
	DeskFolio:                             20,
	IncognitoGuest:                        21,
	ComplimentGuest:                       22,
	HouseUseGuest:                         23,
	EarlyCheckIn:                          24,
	DayUse:                                25,
	EarlyDeparture:                        26,
	ExpectedDeparture:                     27,
	ExtendedDeparture:                     28,
	DepartureList:                         29,
	FolioTransaction:                      30,
	DailyFolioTransaction:                 31,
	MonthlyFolioTransaction:               32,
	YearlyTransaction:                     33,
	ChargeList:                            34,
	DailyChargeList:                       35,
	MonthlyChargeList:                     36,
	YearlyChargeList:                      37,
	CashierReport:                         38,
	PaymentList:                           39,
	DailyPaymentList:                      40,
	MonthlyPaymentList:                    41,
	YearlyPaymentList:                     42,
	ExportCSVbyDepartureDate:              43,
	GuestLedger:                           44,
	GuestDeposit:                          45,
	GuestAccount:                          46,
	DailySales:                            47,
	DailyRevenueReport:                    48,
	DailyRevenueReportSummary:             49,
	FolioOpenBalance:                      50,
	Correction:                            51,
	VoidList:                              52,
	CancelCheckIn:                         53,
	LostandFound:                          54,
	RoomList:                              55,
	RoomType:                              56,
	RoomRate:                              57,
	RoomCountSheet:                        58,
	RoomCountSheetByBuildingFloorRoomType: 59,
	RoomCountSheetByRoomTypeBedType:       60,
	RoomUnavailable:                       61,
	RoomSales:                             62,
	RoomHistory:                           63,
	RoomTypeAvailability:                  64,
	RoomTypeAvailabilityDetail:            65,
	RoomStatus:                            66,
	Sales:                                 67,
	SalesSummary:                          68,
	FrequentlySales:                       69,
	CaptainOrderList:                      70,
	CancelledCaptainOrder:                 71,
	VoidedCheckList:                       72,
	GuestProfile:                          73,
	FrequentlyGuest:                       74,
	Company:                               75,
	PhoneBook:                             76,
	ContractRate:                          77,
	EventList:                             78,
	ReservationChart:                      79,
	ReservationGraphic:                    80,
	OccupiedGraphic:                       81,
	OccupiedbyBusinessSourceGraphic:       82,
	OccupiedbyMarketGraphic:               83,
	OccupiedbyGuestTypeGraphic:            84,
	OccupiedbyCountryGraphic:              85,
	OccupiedbyStateGraphic:                86,
	OccupancyGraphic:                      87,
	RoomAvailabilityGraphic:               88,
	RoomUnvailabilityGraphic:              89,
	RevenueGraphic:                        90,
	PaymentGraphic:                        91,
	RoomStatistic:                         92,
	GuestForecastReport:                   93,
	CityLedgerContributionAnalysis:        94,
	CityLedgerList:                        95,
	CityLedgerAgingReport:                 96,
	CityLedgerAgingReportetail:            97,
	CityLedgerInvoice:                     98,
	CityLedgerInvoiceDetail:               99,
	CityLedgerPayment:                     100,
	CityLedgerMutation:                    101,
	BankReconciliation:                    102,
	APRefundDepositList:                   103,
	APRefundDepositAgingReport:            104,
	APRefundDepositAgingReportetail:       105,
	APRefundDepositPayment:                106,
	APRefundDepositMutation:               107,
	APCommissionList:                      108,
	APCommissionAgingReport:               109,
	APCommissionAgingReportetail:          110,
	APCommissionPayment:                   111,
	APCommissionMutation:                  112,
	LogUser:                               113,
	LogMoveRoom:                           114,
	LogTransferTransaction:                115,
	LogSpecialAccess:                      116,
	KeyLockHistory:                        117,
	LogVoidTransaction:                    118,
	LogHouseKeeping:                       119,
	RoomRateBreakdown:                     120,
	Package:                               121,
	PackageBreakdown:                      122,
	FBStatistic:                           123,
	Member:                                124,
	Voucher:                               125,
	VoucherSRC:                            126,
	CancelledCaptainOrderDetail:           127,
	RoomRateStructure:                     128,
	PaymentBySubDepartment:                129,
	PaymentByAccount:                      130,
	RoomProduction:                        131,
	RoomRevenue:                           132,
	OTAProductivity:                       133,
	GuestForecastReportYearly:             134,
	CashierReportReprint:                  135,
	GuestInHouseListing:                   136,
	GuestInHouseForecast:                  137,
	BankTransactionList:                   138,
	RepeaterGuest:                         139,
	MarketStatistic:                       140,
	PackageSales:                          141,
	LogPABX:                               142,
	BankTransactionAgingReport:            143,
	BankTransactionAgingReportDetail:      144,
	BankTransactionMutation:               145,
	FolioList:                             146,
	LogShift:                              147,
	GuestForecastComparison:               148,
	CashSummaryReport:                     149,
	TransactionByStaff:                    150,
	TaxBreakDownDetailed:                  151,
	DailyFlashReport:                      152,
	BreakfastControl:                      153,
	DailyHotelCompetitor:                  154,
	ActualDepartureGuestList:              155,
	TodayRoomRevenueBreakdown:             156,
	RoomSalesByRoomNumber:                 157,
	DailyStatisticReport:                  158,
	RateCodeAnalysis:                      159,
	GuestInHouseByCity:                    160,
	CancelCheckOut:                        161,
	SalesContributionAnalysis:             162,
	GuestInHousebyBookingSource:           163,
	GuestInHouseByNationality:             164,
	GuestList:                             165,
	LeadList:                              166,
	TaskList:                              167,
	ProposalList:                          168,
	ActivityLog:                           169,
	SalesActivityDetail:                   170,
	RevenueBySubDepartment:                171,
	GuestInHouseBreakfast:                 172}
var UserSpecialAccessOrder = TUserSpecialAccessOrder{
	UnlockReservation:     1,
	VoidReservation:       2,
	VoidDeposit:           3,
	CorrectDeposit:        4,
	DecreaseStay:          5,
	BusinessSource:        21,
	OverrideRateDiscount:  22,
	ModifyScheduleRate:    6,
	ModifyBreakdown:       7,
	ModifyExtraCharge:     8,
	ComplimentGuest:       9,
	HouseUseGuest:         10,
	MoveRoom:              11,
	VoidSubFolio:          12,
	CorrectSubFolio:       13,
	CancelCheckIn:         14,
	CancelCheckOut:        15,
	CreateMasterFolio:     16,
	PrintInvoice:          17,
	ModifyClosedJournal:   18,
	TransferToDeskFolio:   19,
	TransferToMasterFolio: 20,
	UpdateGuestName:       23}
var UserKeylockAccessOrder = TUserKeylockAccessOrder{
	CheckInWithoutCard:              1,
	CheckOutWithoutCard:             2,
	IssuedCardMoreThanTwice:         3,
	ModifyArrivalDate:               4,
	ModifyDepartureDate:             5,
	ModifyDepartureTime:             6,
	DepartureDate1Night:             7,
	ShowAccessIssuedCardMoreThanOne: 8,
	CanIssuedCardMoreThanOne:        9}
var UserReservationAccessOrder = TUserReservationAccessOrder{
	Insert:              1,
	Update:              2,
	Duplicate:           3,
	Deposit:             4,
	Cancel:              5,
	Void:                6,
	NoShow:              7,
	AutoAssign:          8,
	Lock:                9,
	CheckIn:             10,
	InsertFromAllotment: 11,
	Keylock:             12}
var UserDepositAccessOrder = TUserDepositAccessOrder{
	Insert:               1,
	Cash:                 2,
	Card:                 3,
	Refund:               4,
	Void:                 5,
	Correction:           6,
	Transfer:             7,
	UpdateSubDepartment:  8,
	UpdateRemark:         9,
	UpdateDocumentNumber: 10}
var UserInHouseAccessOrder = TUserInHouseAccessOrder{
	Transaction:   1,
	Update:        2,
	Keylock:       3,
	Compliment:    4,
	HouseUse:      5,
	MoveRoom:      6,
	SwitchRoom:    7,
	LockFolio:     8,
	CancelCheckIn: 9,
	GueestMessage: 10,
	ToDo:          11,
	CheckOut:      12}
var UserWalkInAccessOrder = TUserWalkInAccessOrder{
	ScheduleRate: 1,
	Breakdown:    2,
	ExtraCharge:  3}
var UserFolioAccessOrder = TUserFolioAccessOrder{
	Charge:               1,
	Cash:                 2,
	Card:                 3,
	DirectBill:           4,
	UpdateDirectBill:     5,
	CashRefund:           6,
	OtherPayment:         7,
	Void:                 8,
	Correction:           9,
	Transfer:             10,
	UpdateSubDepartment:  11,
	UpdateRemark:         12,
	UpdateDocumentNumber: 13,
	CheckOut:             14,
	PrintFolio:           15}
var UserFolioHistoryAccessOrder = TUserFolioHistoryAccessOrder{
	Transaction:    1,
	PrintFolio:     2,
	CancelCheckOut: 3}
var UserFloorPlanAccessOrder = TUserFloorPlanAccessOrder{
	Reception:       1,
	HouseKeeping:    2,
	ModifyFloorPlan: 3}
var UserCompanyAccessOrder = TUserCompanyAccessOrder{
	Insert:         1,
	Update:         2,
	Delete:         3,
	APLimit:        4,
	ARLimit:        5,
	DirectBill:     6,
	BusinessSource: 7}
var UserInvoiceAccessOrder = TUserInvoiceAccessOrder{
	Insert:        1,
	Update:        2,
	Delete:        3,
	InsertPayment: 4,
	DeletePayment: 5,
	PrintReceipt:  6,
	ExchangeRate:  7}
var UserMemberVoucherGift = TUserMemberVoucherGift{
	RedeemPoint:       1,
	InsertVoucher:     2,
	DeleteVoucher:     3,
	ApproveVoucher:    4,
	SoldVoucher:       5,
	RedeemVoucher:     6,
	ComplimentVoucher: 7}
var UserPreviewReportAccessOrder = TUserPreviewReportAccessOrder{
	EditReport:      1,
	ExportReport:    2,
	CustomizeReport: 3}
var UserPaymentByAPARAccessOrder = TUserPaymentByAPARAccessOrder{
	PaymentByAPAR: 1}
var OTPStatus = TOTPStatus{
	Active:    "A",
	Used:      "U",
	Expire:    "E",
	NotActive: "N"}
var TransactionType = TTransactionType{
	Debit:  "D",
	Credit: "C"}
var FolioType = TFolioType{
	GuestFolio:  "F",
	MasterFolio: "M",
	DeskFolio:   "D"}
var GlobalAccountGroup = TGlobalAccountGroup{
	Charge:     "1",
	Payment:    "2",
	TaxService: "3",
	Transfer:   "4"}
var GlobalAccountSubGroup = TGlobalAccountSubGroup{
	RoomCharge:        "ROOM",
	AccountPayable:    "ACPY",
	Payment:           "PYMT",
	CreditDebitCard:   "CRDB",
	BankTransfer:      "BKTR",
	AccountReceivable: "ACRV",
	Compliment:        "CMPL"}
var RoomStatus = TRoomStatus{
	Ready:             "R",
	Clean:             "C",
	Dirty:             "D",
	Vacant:            "V",
	Occupied:          "O",
	HouseUseX:         "H",
	Compliment:        "P",
	DontDisturb:       "DD",
	DoubleLock:        "DL",
	SleepOut:          "SO",
	PleaseClean:       "PC",
	OutOfOrder:        "OO",
	OfficeUse:         "OF",
	UnderConstruction: "UC"}
var RoomBlockStatus = TRoomBlockStatus{
	GeneralCleaning: "G",
	ShowingRoom:     "S"}
var ReservationStatus = TReservationStatus{
	New:      "N",
	WaitList: "W",
	InHouse:  "I",
	Canceled: "C",
	NoShow:   "S",
	Void:     "V",
	CheckOut: "O"}

var ReservationBlockType = TReservationBlockType{
	Reservation: "R",
	BlockOnly:   "B",
	ChekIn:      "I"}
var ReservationStatus2 = TReservationStatus2{
	Tentative:  "T",
	Confirm:    "C",
	Guaranteed: "G"}
var ReservationType = TReservationType{
	Guaranteed:    "GRTD",
	NonGuaranteed: "NGRD"}
var FolioStatus = TFolioStatus{
	Open:          "O",
	Closed:        "C",
	CancelCheckIn: "I"}
var FolioTransferBy = TFolioTransferBy{
	NoTransfer:        "N",
	ByAccount:         "A",
	ByAccountSubGroup: "S"}
var SubFolioGroup = TSubFolioGroup{
	A: "A",
	B: "B",
	C: "C",
	D: "D"}
var SubFolioPostingType = TSubFolioPostingType{
	None:        "N",
	Deposit:     "D",
	Transfer:    "T",
	Room:        "R",
	ExtraCharge: "E"}
var CommissionType = TCommissionType{
	PercentFirstNightFullRate: "1",
	PercentPerNightFullRate:   "2",
	PercentFirstNightNettRate: "3",
	PercentPerNightNettRate:   "4",
	FixAmountFirstNight:       "5",
	FixAmountPerNight:         "6",
	PercentOfPriceFullPrice:   "7",
	PercentOfPriceNettPrice:   "8"}

var DepreciationType = TDepreciationType{
	None:                   "N",
	StraightLine:           "S",
	DecliningBalance:       "F",
	DoubleDecliningBalance: "D",
	SumOfYearDigit:         "Y"}
var PurchaseRequestStatus = TPurchaseRequestStatus{
	NotApproved: "N",
	Approved:    "A",
	Rejected:    "R"}
var StoreRequisitionStatus = TStoreRequisitionStatus{
	NotApproved: "N",
	Approved:    "A",
	Rejected:    "R"}
var FAItemCondition = TFAItemCondition{
	Good:             "G",
	Broken:           "B",
	Repaired:         "R",
	Sold:             "S",
	Transfered:       "T",
	FullyDepreciated: "D"}
var FALocationType = TFALocationType{
	HKStore: "H",
	Laundry: "L",
	Room:    "R",
	Other:   "O"}
var ItemGroupTypeCode = TItemGroupTypeCode{
	Food:     "F",
	Beverage: "B"}
var ItemGroupCode = TItemGroupCode{
	Food:     "FOOD",
	Beverage: "BVRG"}

var SpecialAccessName = []string{
	"Unlock Reservation",
	"Void Reservation",
	"Void Deposit",
	"Correct Deposit",
	"Decrease Stay",
	"Scheduled Rate",
	"Breakdown",
	"Extra Charge",
	"Compliment Guest",
	"House Use Guest",
	"Move Room",
	"Void Sub Folio",
	"Correct Sub Folio",
	"Cancel Check In",
	"Cancel Check Out",
	"Create Master Folio",
	"Print More Invoice",
	"Modify Closed Journal",
	"Transfer to Desk Folio",
	"Transfer to Master Folio",
	"Business Source",
	"Override Rate/Discount"}
var CPType = TCPType{
	Guest:     "G",
	Company:   "C",
	Invoice:   "I",
	RoomOwner: "R"}
var ChargeFrequency = TChargeFrequency{
	OnceOnly: "0",
	Daily:    "1",
	Weekly:   "2",
	Monthly:  "3"}
var KeylockVendor = TKeylockVendor{
	Ventaza:            "VENT",
	UltraLock:          "ULTR",
	DLock:              "DLOC",
	DLockOldVersion:    "DLOL",
	Kiara:              "KIAR",
	Colcom:             "COLC",
	SureLock:           "SURE",
	Deluns:             "DLUN",
	Kima:               "KIMA",
	HuneLock:           "HUNE",
	VingCard:           "VING",
	SunVitio:           "SUNV",
	Rafles:             "RAFE",
	BeTech:             "BTEC",
	OnityOld:           "ONOL",
	Saflok:             "SAFL",
	ColcomOld:          "COLO",
	Tesa:               "TESA",
	VingcardVisionline: "VISI",
	CLock:              "CAKR",
	Downs:              "DOWN",
	Onity:              "ONIT",
	OnityOldVersion:    "ONIO",
	VingCardSerial:     "VISE",
	Kaba:               "KABA",
	RaflesUSB:          "RAUS"}
var MikrotikVendor = TMikrotikVendor{
	Mikrotik: "MIKR",
	Megalos:  "MEGA",
	Provenos: "PROV",
	Coova:    "COOV"}
var ChannelManager = TChannelManager{
	BookNLink:  "BKNL",
	Stah:       "STAH",
	SiteMinder: "STMD"}
var PMSCommand = TPMSCommand{
	Issued:     "I",
	Verified:   "V",
	Replaced:   "R",
	Erase:      "E",
	ForceErase: "F"}
var JournalPrefix = TJournalPrefix{
	Manual:            "JM",
	Transaction:       "JT",
	Disbursement:      "JD",
	Receive:           "JR",
	Inventory:         "JI",
	Adjustment:        "JA",
	FixedAsset:        "JF",
	BeginningYear:     "JB",
	AccountPayable:    "AP",
	AccountReceivable: "AR"}
var JournalType = TJournalType{
	CashIn:                   "1",
	CashOut:                  "2",
	CashTransfer:             "3",
	BankIn:                   "4",
	BankOut:                  "5",
	BankTransfer:             "6",
	CreditCardReconciliation: "7",
	Cheque:                   "8",
	Other:                    "9"}
var JournalGroup = TJournalGroup{
	CapitalReceipt:      "CAPT",
	Receiving:           "RECV",
	ARPayment:           "PYAR",
	Loan:                "LOAN",
	Investment:          "INVS",
	BankWithdrawal:      "BKWD",
	AssetPurchasing:     "PCAS",
	InventoryPurchasing: "PCIV",
	OtherPurchasing:     "PCOT",
	APPayment:           "PYAP",
	Expense:             "EXPE",
	BankDeposit:         "BKDP",
	CashTransfer:        "TRCS",
	BankTransfer:        "TRBK",
	AccountPayable:      "ACPY",
	AccountReceivable:   "ACRV",
	Costing:             "COST",
	FixedAsset:          "FAST",
	Other:               "OTHE"}
var GlobalJournalAccountGroup = TGlobalJournalAccountGroup{
	Assets:       "1",
	Liability:    "2",
	Equity:       "3",
	Income:       "4",
	Cost:         "5",
	Expense1:     "6",
	Expense2:     "7",
	OtherIncome:  "8",
	OtherExpense: "9"}
var GlobalJournalAccountType = TGlobalJournalAccountType{
	Bank:                  "1",
	AccountReceivable:     "2",
	CreditCard:            "3",
	OtherCurrentAsset:     "4",
	FixAsset:              "5",
	OtherAsset:            "6",
	AccountPayable:        "7",
	OtherCurrentLiability: "8",
	LongTermLiability:     "9",
	OtherLiability:        "10",
	Equity:                "11",
	Income:                "12",
	Cost:                  "13",
	Expense:               "14",
	OtherIncome:           "15",
	OtherExpense:          "16"}
var BankAccountType = TBankAccountType{
	CashAccount:   "C",
	SavingAccount: "S",
	CreditAccount: "R",
	ChequeAccount: "Q"}
var BudgetType = TBudgetType{
	Manual:     "M",
	Average:    "A",
	Percentage: "P",
	Daily:      "D"}
var CaptainOrderType = TCaptainOrderType{
	DineIn:      "I",
	TakeAway:    "T",
	Delivery:    "D",
	Reservation: "R"}
var ForeignCashTableID = TForeignCashTableID{
	GuestDeposit:   1,
	SubFolio:       2,
	InvoicePayment: 31}
var StatisticAccount = TStatisticAccount{
	TotalRoom:              "0101",
	OutOfOrder:             "0102",
	OfficeUse:              "0103",
	UnderConstruction:      "0104",
	HouseUse:               "0105",
	Compliment:             "0201",
	RoomSold:               "0202",
	DayUse:                 "0203",
	WalkIn:                 "0204",
	CheckInByReservation:   "0205",
	NoShow:                 "0206",
	ReservationMade:        "0207",
	CancelationReservation: "0208",
	EarlyCheckOut:          "0209",
	CheckOut:               "0210",
	RevenueGross:           "0301",
	RevenueNonPackage:      "0302",
	RevenueNett:            "0303",
	RevenueWithCompliment:  "0304",
	Adult:                  "0401",
	Child:                  "0402",
	AdultHouseUse:          "0403",
	ChildHouseUse:          "0404",
	NumberOfCover:          "1101",
	NettFoodSales:          "1102",
	NettBeverageSales:      "1103",
	NettBreakfastSales:     "1104",
	NettBanquetSales:       "1105",
	NettWeddingSales:       "1112",
	NettGatheringSales:     "1113",
	BreakfastCover:         "1106",
	BeverageCover:          "1107",
	FoodCover:              "1108",
	BanquetCover:           "1109",
	WeddingCover:           "1110",
	GatheringCover:         "1111"}
var PaymentGroup = TPaymentGroup{
	Cash:       "C",
	Bank:       "B",
	DirectBill: "D",
	Other:      "O",
	None:       "N"}
var InventoryItemGroup = TInventoryItemGroup{
	Food:     "FOOD",
	Beverage: "BVRG"}
var DynamicRateType = TDynamicRateType{
	None:           "N",
	BaseOccupancy:  "O",
	BaseSession:    "S",
	BaseScale:      "L",
	BaseCompetitor: "C",
	BaseWeekly:     "W"}
var MemberType = TMemberType{
	Room:    "R",
	Outlet:  "O",
	Banquet: "B"}
var VoucherType = TVoucherType{
	Compliment: "C",
	Discount:   "D",
	Sale:       "S"}
var VoucherStatus = TVoucherStatus{
	Active: "A",
	Used:   "U",
	Expire: "E"}
var VoucherStatusApprove = TVoucherStatusApprove{
	Unapprove:   "U",
	Approved:    "A",
	NotApproved: "N"}
var VoucherStatusSold = TVoucherStatusSold{
	Sold:       "S",
	Redeemed:   "R",
	Compliment: "C"}
var SMSevent = TSMSEvent{
	OnInsertReservation: 100001,
	OnWalkIn:            101001,
	OnCheckIn:           101002,
	OnCheckOut:          101003,
	OnDayendCloseFinish: 105001}
var ComplimentType = TComplimentType{
	None:           "N",
	Compliment:     "C",
	OfficerCheck:   "O",
	EntertainCheck: "E"}

var DepartmentType = TDepartmentType{
	Income:      "I",
	OtherIncome: "O",
	NonIncome:   "N"}
var GuestProfileSource = TGuestProfileSource{
	Hotel:   "H",
	Pos:     "P",
	Banquet: "B"}
var CMUpdateType = TCMUpdateType{
	Reservation:   "R",
	Folio:         "F",
	RoomAllotment: "A",
	Rate:          "T",
	Availability:  "U"}
var InventoryCostingMethod = TInventoryCostingMethod{
	FIFO:    "F",
	LIFO:    "L",
	Average: "A"}

var JournalAccountSubGroupType = TJournalAccountSubGroupType{
	PayrollRelated: "1",
	OtherExpense:   "2",
	EnergyCost:     "3",
	SalesEnergy:    "4"}
var SalesActivityStatus = TSalesActivityStatus{
	Deal:         "D",
	Finished:     "F",
	New:          "N",
	Qualified:    "Q",
	ProposalSend: "S",
	Working:      "W",
	Void:         "V"}

var SalesActivityProposalStatus = TSalesActivityProposalStatus{
	Draft:     "D",
	Send:      "S",
	Revised:   "R",
	Contacted: "T",
	Declined:  "C",
	Accepted:  "A",
	Void:      "V"}

var SalesActivityTaskPriority = TSalesActivityTaskPriority{
	Low:    "L",
	Medium: "M",
	High:   "H",
	Urgent: "U"}

var SalesActivityTaskStatus = TSalesActivityTaskStatus{
	NotStarted:       "N",
	InProgress:       "I",
	Testing:          "T",
	AwaitingFeedback: "A",
	Complete:         "C",
	Void:             "V"}
var NotifThirdPartyTemplateID = TNotifThirdPartyTemplateID{
	EmailReminderreservation:      1,
	EmailOnCheckIn:                2,
	EmailOnCheckOut:               3,
	EmailOnGusetBirthday:          4,
	EmailSalesActivitySendRemider: 5,

	WAReminderreservation:      20,
	WAOnCheckIn:                21,
	WAOnCheckOut:               22,
	WAOnGusetBirthday:          23,
	WASalesActivitySendRemider: 24}
var NotifThirdPartySourceCode = TNotifThirdPartySourceCode{
	Reservation:              "R",
	Folio:                    "F",
	SalesActivitySendRemider: "S",
	Other:                    "O"}

var CalenderLabelColor = TCalenderLabelColor{
	Reservation: 16711680,
	InHouse:     32896,
	CheckOut:    65535}
var PaymentSourceAPAR = TPaymentSourceAPAR{
	None:              "N",
	AccountPayable:    "P",
	AccountReceivable: "R",
	CityLedger:        "C",
	APCommission:      "A"}

var ResponseCode = TResponseCode{
	Successfully:                         0,
	NotAuthorized:                        1,
	InvalidDataFormat:                    2,
	ErrorCreateToken:                     3,
	DataNotFound:                         4,
	InvalidDataValue:                     5,
	DatabaseValueChanged:                 6,
	DatabaseError:                        7,
	EmptyRequireField:                    8,
	PaymentDateCannotLowerThanIssuedDate: 9,
	SuccessfullyWithStatus:               10,
	DuplicateEntry:                       11,
	SubscriptionExpired:                  12,
	Unregistered:                         13,
	InternalServerError:                  500,
	OtherResult:                          999}

var ResponseText = TResponseText{
	Successfully:                         "Successfully",
	NotAuthorized:                        "NotAuthorized",
	InvalidDataFormat:                    "InvalidDataFormat",
	ErrorCreateToken:                     "ErrorCreateToken",
	DataNotFound:                         "DataNotFound",
	InvalidDataValue:                     "InvalidDataValue",
	DatabaseValueChanged:                 "DatabaseValueChanged",
	DatabaseError:                        "DatabaseError",
	EmptyRequireField:                    "EmptyRequireField",
	PaymentDateCannotLowerThanIssuedDate: "PaymentDateCannotLowerThanIssuedDate",
	SuccessfullyWithStatus:               "SuccessfullyWithStatus",
	DuplicateEntry:                       "DuplicateEntry",
	OtherResult:                          "OtherResult"}

var ReceivedStatus = TReceivedStatus{
	UnReceived:      0,
	ReceivedPartial: 1,
	Received:        2,
}

var GlobalPaymentType = TGlobalPaymentType{
	TransferDeposit: "TRDP",
}

var FormMessage = TFormMessage{
	GuestInHouse: []string{
		0:  "Can not Cancel Check In, because Arrival Date different with Audit Date",
		1:  "There is Guest Card active, check out without card?",
		2:  "Room Charge already posting, are you want to repost room charge?",
		3:  "No need to create new profile for this guest.",
		4:  "Cannot posting compliment or house use guest.",
		5:  "Voucher only can update on first night.",
		6:  "Room Charge cannot posted, out of stay date",
		7:  "Folio have transaction that already paid or invoice created or foreign cash already changed. Please delete payment and invoice first or delete foreign cash changed to void or correct it.",
		8:  "Please input phone number as hotspot password",
		9:  "Please input ID card number as hotspot password",
		10: "Posting successful",
		11: "Breakdown is too large",
		12: "Zero amount not allowed",
		13: "No room charge for today",
		14: "No room charge for compliment",
		15: "Posting room charge failed",
		16: "Please input hotspot password",
	},
	Housekeeping: []string{
		0: "This data conflict with another data",
		1: "This room is occupied",
		2: "This room is already block for general cleaning",
		3: "This room is already block for showing room",
		4: "This room is already unblock",
		5: "This room is vacant",
	},
	Unavailable: []string{
		0: "Date must be greater or equal than audit date",
		1: "Date you entered conflict with another data",
		2: "Date range is not valid",
		3: "Cannot delete data. Start date is lower than audit date",
	},
	FolioHistory: []string{
		0: "Cannot cancel check out, because room is not available anymore.",
		1: "Cannot cancel check out, because departure date is different with audit date.",
		2: "Only guest folio can cancel check out",
		3: "Cannot cancel check out, this folio was cancel check in",
		4: "Cannot cancel check out, this folio was created invoice",
	},
	PostingRoomCharge: []string{
		0:   "Reposting successful",
		1:   "Breakdown is too large",
		2:   "Zero amount not allowed",
		3:   "No room charge for today",
		4:   "No room charge for compliment",
		255: "Posting room charge failed",
	},
	PostingExtraCharge: []string{
		0:   "Reposting successful",
		1:   "Breakdown is too large",
		2:   "Zero amount not allowed",
		3:   "No extra charge for today",
		4:   "No extra charge for compliment",
		255: "Posting extra charge failed",
	},
	TransactionTerminal: []string{
		0:  "Are you want to close captain order?",
		1:  "Please input transaction",
		2:  "Data not changed beacuse value is zero",
		3:  "Discount is too large",
		4:  "There are item that had price lower then discount, this items will not apply new discount",
		5:  "Price cannot lower than discount",
		6:  "There are item that had price lower then discount, this items will not apply new price",
		7:  "Folio transfered is already closed, please select another folio for transfer",
		8:  "Cannot input this type of payment, because there is compliment payment",
		9:  "Cannot input this type of payment, because there is non compliment payment",
		10: "Cannot input this compliment, beacause there is other type of compliment payment",
		11: "Cannot discount, override or modify price for free product",
		12: "Transfer charge cannot greater than bill amount",
		13: "This Captain Order already close",
		14: "Discount is large then discount limit",
		15: "There are item that had discount large then discount limit, this items will not apply new discount",
		16: "Comport not installed",
		17: "Are you want to print to kitchen?",
		18: "All order already printed on kitchen",
		19: "All order already printed",
		20: "Product can not discount, product set disable discount",
		21: "This is not spa transaction",
		22: "Massage already started",
		23: "Massage already stopped",
		24: "Therapist fingerprint mismatch",
		25: "Captain order closed by other user/computer",
		26: "Cannot insert payment on compliment captain order."},
	RoomRate: []string{
		0: "Cannot delete this room rate, because it use on reservation or guest folio",
		1: "There is no guest in house using this rate",
		2: "Please fill channel manager inventory code",
	},
	APARPayment: []string{
		0: "Value cannot below zero",
		1: "AP/AR Already Paid",
		2: "AP/AR Number already exist, please select another",
		3: "Discount cannot be greater than Total",
		4: "Total expense cannot be greater than Total",
		5: "Amount cannot be greater than Outstanding",
		6: "If Payment By AP/AR, Total Amount Must be Equal with Total AP/AR Amount Use",
		7: "Payment Date not valid, There is Payment AP/AR Date Greater than Payment Date",
	},
}

var UserInfo = TUserInfo{}

var WSMessageType = TWSMessageType{
	Broadcast:  1,
	Channel:    2,
	Client:     3,
	Connection: 4,
	Room:       5,
}

var WSDataType = TWSDataType{
	ServerStatus:                   1,
	DayendCloseStatus:              2,
	ModifiedRoomAvailabilityStatus: 3,
	AuditDateChanged:               4,
}

var ReportCodeName = TReportCodeName{
	Reservation:                           110,
	FrontDesk:                             115,
	MemberVoucherAndGift:                  117,
	Room:                                  120,
	CakrasoftPointOfSales:                 125,
	Profile:                               130,
	MarketingGraphicAndAnalisis:           135,
	SalesActivity:                         136,
	CityLedger:                            140,
	ApRefundDepositAndCommission:          150,
	Log:                                   160,
	DailyReport:                           181,
	MonthlyReport:                         182,
	YearlyReport:                          183,
	FavoriteReport:                        190,
	ReservationList:                       11001,
	CancelledReservation:                  11002,
	NoShowReservation:                     11003,
	VoidedReservation:                     11004,
	GroupReservation:                      11005,
	ExpectedArrival:                       11006,
	ArrivalList:                           11007,
	SamedayReservation:                    11008,
	AdvancedPaymentDeposit:                11009,
	BalanceDeposit:                        11010,
	WaitListReservation:                   11011,
	CurrentInHouse:                        11501,
	GuestInHouse:                          11502,
	GuestInHouseByBusinessSource:          11503,
	GuestInHouseByMarket:                  11504,
	GuestInHouseByGuestType:               11505,
	GuestInHouseByCountry:                 11506,
	GuestInHouseByState:                   11507,
	GuestInHouseListing:                   11508,
	GuestInHouseByCity:                    11509,
	GuestInHouseByBookingSource:           11510,
	MasterFolio:                           11511,
	DeskFolio:                             11512,
	GuestInForecast:                       11513,
	RepeaterGuest:                         11514,
	FolioList:                             11515,
	GuestInHouseByNationality:             11516,
	GuestList:                             11517,
	GuestInHouseBreakfast:                 11518,
	IncognitoGuest:                        11520,
	ComplimentGuest:                       11521,
	HouseUseGuest:                         11522,
	EarlyCheckIn:                          11523,
	DayUse:                                11524,
	EarlyDeparture:                        11525,
	ExpectedDeparture:                     11526,
	ExtendedDeparture:                     11527,
	DepartureList:                         11528,
	ActualDepartureGuestList:              11529,
	FolioTransaction:                      11540,
	DailyFolioTransaction:                 11541,
	MonthlyFolioTransaction:               11542,
	YearlyFolioTransaction:                11543,
	ChargeList:                            11544,
	DailyChargeList:                       11545,
	MonthlyChargeList:                     11546,
	YearlyChargeList:                      11547,
	CashierReport:                         11548,
	PaymentList:                           11549,
	DailyPaymentList:                      11550,
	MonthlyPaymentList:                    11551,
	YearlyPaymentList:                     11552,
	ExportCsvByDepartureDate:              11560,
	GuestLedger:                           11561,
	GuestDeposit:                          11562,
	GuestAccount:                          11563,
	DailySales:                            11564,
	DailyRevenueReport:                    11565,
	DailyRevenueReportSummary:             11566,
	PaymentBySubDepartment:                11567,
	PaymentByAccount:                      11568,
	RevenueBySubDepartment:                11569,
	FolioOpenBalance:                      11580,
	Correction:                            11581,
	VoidList:                              11582,
	CancelCheckIn:                         11583,
	LostAndFound:                          11584,
	CashierReportReprint:                  11585,
	PackageSales:                          11586,
	CashSummaryReport:                     11587,
	TransactionReportByStaff:              11588,
	TaxBreakdownDetailed:                  11589,
	TodayRoomRevenueBreakdown:             11590,
	CancelCheckOut:                        11591,
	Member:                                11701,
	Voucher:                               11702,
	VoucherSoldRedemeedAndComplimented:    11703,
	RoomList:                              12001,
	RoomType:                              12002,
	RoomRate:                              12003,
	RoomCountSheet:                        12004,
	RoomUnavailable:                       12005,
	RoomSales:                             12006,
	RoomHistory:                           12007,
	RoomTypeAvailability:                  12008,
	RoomTypeAvailabilityDetail:            12009,
	RoomStatus:                            12010,
	RoomCountSheetByBuildingFloorRoomType: 12011,
	RoomCountSheetByRoomTypeBedType:       12012,
	RoomSalesByRoomNumber:                 12013,
	RoomRateBreakdown:                     12021,
	Package:                               12022,
	PackageBreakdown:                      12023,
	RoomRateStructure:                     12024,
	RoomRevenue:                           12025,
	Sales:                                 12501,
	SalesSummary:                          12502,
	FrequentlySales:                       12503,
	CaptainOrderList:                      12504,
	CancelledCaptainOrder:                 12505,
	VoidedCheckList:                       12506,
	CancelledCaptainOrderDetail:           12507,
	BreakfastControl:                      12508,
	GuestProfile:                          13001,
	FrequentlyGuest:                       13002,
	Company:                               13003,
	PhoneBook:                             13004,
	ContractRate:                          13501,
	EventList:                             13502,
	ReservationChart:                      13503,
	ReservationGraphic:                    13504,
	OccupiedGraphic:                       13505,
	OccupiedByBusinessSourceGraphic:       13506,
	OccupiedByMarketGraphic:               13507,
	OccupiedByGuestTypeGraphic:            13508,
	OccupiedByCountryGraphic:              13509,
	OccupiedByStateGraphic:                13510,
	OccupancyGraphic:                      13511,
	RoomAvailabilityGraphic:               13531,
	RoomUnvailabilityGraphic:              13532,
	RevenueGraphic:                        13533,
	PaymentGraphic:                        13534,
	RoomStatistic:                         13535,
	GuestForecastReport:                   13536,
	CityLedgerContributionAnalysis:        13537,
	FoodAndBeverageStatistic:              13538,
	RoomProduction:                        13539,
	BusinessSourceProductivity:            13540,
	GuestForecastReportYearly:             13541,
	MarketStatistic:                       13542,
	GuestForecastComparison:               13543,
	DailyFlashReport:                      13544,
	DailyHotelCompetitor:                  13545,
	DailyStatisticReport:                  13546,
	RateCodeAnalysis:                      13547,
	SalesContributionAnalysis:             13548,
	LeadList:                              13601,
	TaskList:                              13602,
	ProposalList:                          13603,
	ActivityLog:                           13604,
	SalesActivityDetail:                   13605,
	CityLedgerList:                        14001,
	CityLedgerAgingReport:                 14002,
	CityLedgerAgingReportDetail:           14003,
	CityLedgerInvoice:                     14004,
	CityLedgerInvoiceDetail:               14005,
	CityLedgerInvoicePayment:              14006,
	CityLedgerInvoiceMutation:             14007,
	BankReconciliation:                    14008,
	BankTransactionList:                   14009,
	BankTransactionAgingReport:            14010,
	BankTransactionAgingReportDetail:      14011,
	BankTransactionMutation:               14012,
	ApRefundDepositList:                   15001,
	ApRefundDepositAgingReport:            15002,
	ApRefundDepositAgingReportDetail:      15003,
	ApRefundDepositPayment:                15004,
	ApRefundDepositMutation:               15005,
	ApCommissionAndOtherList:              15006,
	ApCommissionAndOtherAgingReport:       15007,
	ApCommissionAndOtherAgingReportDetail: 15008,
	ApCommissionAndOtherPayment:           15009,
	ApCommissionAndOtherMutation:          15010,
	LogUser:                               16501,
	LogMoveRoom:                           16502,
	LogTransferTransaction:                16503,
	LogSpecialAccess:                      16504,
	KeylockHistory:                        16505,
	LogVoidTransaction:                    16506,
	LogHouseKeeping:                       16507,
	LogPabx:                               16508,
	LogShift:                              16509,
	Product:                               201,
	ProductSales:                          202,
	DayendCloseReprint:                    203,
	ProductCategory:                       21001,
	ProductList:                           21002,
	ProductCosting:                        22007,
	ProductCostingItem:                    22008,
	FAndBRateStructure:                    22010,
	RemovedProductCaptainOrder:            22011,
	RealizationCostOfGoodSold:             22012,
	MasterData:                            301,
	GeneralAndSales:                       302,
	Venue:                                 31001,
	Theme:                                 31002,
	SeatingPlan:                           31003,
	BanquetBookingDetail:                  32001,
	BanquetCalender:                       32002,
	BanquetForecast:                       32003,
	BanquetAdvancedDeposit:                32004,
	BanquetBalanceDeposit:                 32005,
	BanquetChargeList:                     32006,
	BanquetDailySales:                     32007,
	CancelBanquetReservation:              32008,
	VoidBanquetReservation:                32009,
	BanquetBooking:                        32010,
	AccountReceivable:                     401,
	AccountPayable:                        402,
	GeneralLedgerAndBank:                  403,
	AccountReceivableList:                 41011,
	AccountReceivableAgingReport:          41012,
	AccountReceivableAgingReportDetail:    41013,
	AccountReceivablePayment:              41014,
	AccountReceivableMutation:             41015,
	AccountPayableList:                    42011,
	AccountPayableAgingReport:             42012,
	AccountPayableAgingReportDetail:       42013,
	AccountPayablePayment:                 42014,
	AccountPayableMutation:                42015,
	OperationalAccount:                    43001,
	ChartOfAccount:                        43002,
	Journal:                               43003,
	GeneralLedger:                         43004,
	TrialBalance:                          43005,
	WorkSheet:                             43006,
	BalanceSheet:                          43007,
	IncomeStatement:                       43008,
	ProfitAndLoss:                         43009,
	ProfitAndLossDetail:                   43010,
	ProfitAndLossByDepartment:             43011,
	ProfitAndLossDetailByDepartment:       43012,
	BankBookAccount:                       43013,
	CurrentAssetAccount:                   43014,
	CashOnHandAccount:                     43015,
	FixedAssetAccount:                     43016,
	OtherAssetAccount:                     43017,
	ARLedgerAccount:                       43018,
	CurrentLiabilityAccount:               43019,
	APLedgerAccount:                       43020,
	LongTermLiabilityAccount:              43021,
	OtherLiabilityAccount:                 43022,
	BankBookAccountGroup:                  43023,
	BankBookAccountSummary:                43024,
	ProfitAndLossBySubDepartment:          43025,
	ProfitAndLossDetailBySubDepartment:    43026,
	ExportedJournal:                       43027,
	BalanceMultiPeriod:                    43028,
	CashFlow:                              43029,
	ProfitAndLossMultiPeriodDetail:        43030,
	ProfitAndLossGraphic:                  43031,
	Inventory:                             502,
	FixedAsset:                            503,
	Uom:                                   51001,
	InventoryItem:                         51002,
	FixedAssetItem:                        51003,
	InventoryPurchaseOrder:                52001,
	InventoryPurchaseOrderDetail:          52002,
	ReceiveStock:                          52003,
	ReceiveStockDetail:                    52004,
	StockTransfer:                         52005,
	StockTransferDetail:                   52006,
	Costing:                               52007,
	CostingDetail:                         52008,
	StockOpname:                           52009,
	StockOpnameDetail:                     52010,
	StoreStock:                            52011,
	StoreStockCard:                        52012,
	LowLevelStoreStock:                    52013,
	HighLevelStoreStock:                   52014,
	AllStoreStock:                         52015,
	AllStoreStockCard:                     52016,
	LowLevelAllStoreStock:                 52017,
	HighLevelAllStoreStock:                52018,
	Production:                            52019,
	ProductionDetail:                      52020,
	CostRecipe:                            52021,
	ReturnStock:                           52022,
	ReturnStockDetail:                     52023,
	DailyInventoryReconciliation:          52024,
	MonthlyInventoryReconciliation:        52025,
	InventoryReconciliation:               52026,
	AverageItemPricePurchase:              52027,
	ItemPurchasePriceGraphic:              52028,
	InventoryPurchaseRequest:              52029,
	InventoryPurchaseRequestDetail:        52030,
	StoreRequisition:                      52031,
	StoreRequisitionDetail:                52032,
	RecapitulationFoodAndBeverage:         52033,
	RealizationCostOfGoodsSold:            52034,
	ComparisonCostSalesFbGraphic:          52035,
	FixedAssetPurchaseOrder:               53001,
	FixedAssetPurchaseOrderDetail:         53002,
	FixedAssetReceive:                     53003,
	FixedAssetReceiveDetail:               53004,
	FixedAssetList:                        53005,
	FixedAssetDepreciation:                53006,
}

var ReportTemplateName = TReportTemplateName{
	DailyRevenueReport1:        "DailyRevenueReport01.fr3",
	DailyRevenueReport2:        "DailyRevenueReport02.fr3",
	DailyRevenueReport3:        "DailyRevenueReport03.fr3",
	DailyRevenueReport4:        "DailyRevenueReport04.fr3",
	DailyRevenueReport5:        "DailyRevenueReport05.fr3",
	DailyRevenueReport5Point2D: "DailyRevenueReport05Point2D.fr3",
	DailyRevenueReport6:        "DailyRevenueReport06.fr3",
	DailyRevenueReport7:        "DailyRevenueReport07.fr3",
	DailyRevenueReport8:        "DailyRevenueReport08.fr3",
	DailyRevenueReport9:        "DailyRevenueReport09.fr3",
	DailyRevenueReport10:       "DailyRevenueReport10.fr3",
	DailyRevenueReport11:       "DailyRevenueReport11.fr3",
}

var CustomDateOptions = TCustomDateOptions{
	Condition21:         "21",
	Condition22:         "22",
	ConditionAuditDate2: "AU2",
}

var FNBStructureRate = TFNBStructureRate{
	Portion: "1",
	Revenue: "2",
	Average: "3",
}

var TimeSegment = TTimeSegment{
	Breakfast:  "B",
	Lunch:      "L",
	Dinner:     "D",
	CoffeBreak: "C",
}

var CacheKey = TCacheKey{
	LastAPNumber:               "LAST_AP_NUMBER",
	LastARNumber:               "LAST_AR_NUMBER",
	LastRefNumber:              "LAST_REF_NUMBER",
	LastManualRefNumber:        "LAST_MANUAL_RN",
	LastTransactionRefNumber:   "LAST_TRANSACTION_RN",
	LastDisbursementRefNumber:  "LAST_DISBURSEMENT_RN",
	LastReceiveRefNumber:       "LAST_RECEIVE_RN",
	LastInventoryRefNumber:     "LAST_INVENTORY_RN",
	LastAdjustmentRefNumber:    "LAST_ADJUSTMENT_RN",
	LastFixedAssetRefNumber:    "LAST_FIXED_ASSET_RN",
	LastBeginningYearRefNumber: "LAST_BEGINNING_YEAR_RN",
	LastCostingNumber:          "LAST_CO_NUMBER",
	LastInvoiceNumber:          "LAST_INVOICE_NUMBER",
	LastDepreciationNumber:     "LAST_DEP_NUMBER",
	LastPaymentNumber:          "LAST_PAYMENT_NUMBER",
	LastPRNumber:               "LAST_PR_NUMBER",
	LastPONumber:               "LAST_PO_NUMBER",
	LastReceiveNumber:          "LAST_RV_NUMBER",
	LastSRNumber:               "LAST_SR_NUMBER",
	LastProductionNumber:       "LAST_PROD_NUMBER",
	LastReturnStock:            "LAST_RS_NUMBER",
	LastOpnameNumber:           "LAST_OP_NUMBER",
	LastFAPONumber:             "LAST_FAOP_NUMBER",
	LastFAReceiveNumber:        "LAST_FARV_NUMBER",
	LastStockTransferNumber:    "LAST_ST_NUMBER",
}

var ReportTemplate = TReportTemplate{}
var ReportAccessOrder = TReportAccessOrder{
	AccessForm: ReportAccessForm{
		CashierReport:        0,
		FrontDeskReport:      1,
		PointOfSalesReport:   2,
		BanquetReport:        3,
		AccountingReport:     4,
		InventoryAssetReport: 5,
	},
}

// TODO need to set on init
// TODO Load data configuration
var ProgramVariable TProgramVariable
var ProgramConfiguration TProgramConfiguration
var SpecialProduct = TSpecialProduct{
	FoodCode:     "[SPCFOOD]",
	BeverageCode: "[SPCBVRG]",
}

// VariableDLL: TVariableDLL;
// TempVariable: TTempVariable;
var DefaultVariable = TDefaultVariable{}

// ProgramMessage: TProgramMessage;
// DateCheckMessage: array [0..1] of PChar;
// UserInfo: TUserInfo;
var GlobalAccount = TGlobalAccount{}
var GlobalDepartment = TGlobalDepartment{}
var GlobalSubDepartment = TGlobalSubDepartment{}
var GlobalJournalAccount = TGlobalJournalAccount{}
var GlobalJournalAccountSubGroup = TGlobalJournalAccountSubGroup{}
var WeekendDay = TWeekendDay{}
var GlobalJournalAccountGroupName = TGlobalJournalAccountGroupName{}

func Main() {

}

type TResponseCodeCM struct {
	InvalidDate, InvalidEndDate, InvalidStartDate, SystemCurrentlyUnavailable, InvalidRateCode, InvalidValue, RequiredFieldMissing, InvalidActionOrStatusCode, InvalidHotel, InvalidHotelCode,
	InvalidRoomTypeCode, InvalidStartDateEndDate, UnableToUpdate, SystemError, UnableToProcess uint
}

var ResponseCodeCM = TResponseCodeCM{
	InvalidDate:                15,
	InvalidEndDate:             135,
	InvalidStartDate:           136,
	SystemCurrentlyUnavailable: 187,
	InvalidRateCode:            249,
	InvalidValue:               320,
	RequiredFieldMissing:       321,
	InvalidActionOrStatusCode:  356,
	InvalidHotel:               361,
	InvalidHotelCode:           400,
	InvalidRoomTypeCode:        402,
	InvalidStartDateEndDate:    404,
	UnableToUpdate:             447,
	SystemError:                448,
	UnableToProcess:            450,
}
