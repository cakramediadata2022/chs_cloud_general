package dataset

import (
	"fmt"

	"github.com/cakramediadata2022/chs_cloud_general/internal/db_var"
	"github.com/cakramediadata2022/chs_cloud_general/internal/general"
	"github.com/cakramediadata2022/chs_cloud_general/internal/global_var"
	"gorm.io/gorm"
)

func GenerateDataset(DB *gorm.DB) *global_var.TDataset {
	var Dataset *global_var.TDataset
	configuration := make(map[string]map[string]interface{})
	var DataOutput []db_var.Configuration
	DB.Table(db_var.TableName.Configuration).Scan(&DataOutput)
	if len(DataOutput) > 0 {

		for _, configurationX := range DataOutput {
			//bo0lean
			if configurationX.Category == global_var.ConfigurationCategory.General {
				if configuration[configurationX.Category] == nil {
					configuration[configurationX.Category] = make(map[string]interface{})
				}
				configuration[configurationX.Category][configurationX.Name] = configurationX.Value
				continue
			}
			if configurationX.Category == global_var.ConfigurationCategory.Accounting {
				if configuration[configurationX.Category] == nil {
					configuration[configurationX.Category] = make(map[string]interface{})
				}
				configuration[configurationX.Category][configurationX.Name] = configurationX.Value
				continue
			}
			if configurationX.Category == global_var.ConfigurationCategory.AmountPreset {
				if configuration[configurationX.Category] == nil {
					configuration[configurationX.Category] = make(map[string]interface{})
				}
				configuration[configurationX.Category][configurationX.Name] = configurationX.Value
				continue
			}
			if configurationX.Category == global_var.ConfigurationCategory.CompanyBankAccount {
				if configuration[configurationX.Category] == nil {
					configuration[configurationX.Category] = make(map[string]interface{})
				}
				configuration[configurationX.Category][configurationX.Name] = configurationX.Value
				continue
			}
			if configurationX.Category == global_var.ConfigurationCategory.GlobalAccount {
				if configuration[configurationX.Category] == nil {
					configuration[configurationX.Category] = make(map[string]interface{})
				}
				configuration[configurationX.Category][configurationX.Name] = configurationX.Value
				continue
			}
			if configurationX.Category == global_var.ConfigurationCategory.GlobalJournalAccount {
				if configuration[configurationX.Category] == nil {
					configuration[configurationX.Category] = make(map[string]interface{})
				}
				configuration[configurationX.Category][configurationX.Name] = configurationX.Value
				continue
			}
			if configurationX.Category == global_var.ConfigurationCategory.PaymentAccount {
				if configuration[configurationX.Category] == nil {
					configuration[configurationX.Category] = make(map[string]interface{})
				}
				configuration[configurationX.Category][configurationX.Name] = configurationX.Value
				continue
			}
			if configurationX.Category == global_var.ConfigurationCategory.Company {
				if configuration[configurationX.Category] == nil {
					configuration[configurationX.Category] = make(map[string]interface{})
				}
				configuration[configurationX.Category][configurationX.Name] = configurationX.Value
				continue
			}
			if configurationX.Category == global_var.ConfigurationCategory.CustomField {
				if configuration[configurationX.Category] == nil {
					configuration[configurationX.Category] = make(map[string]interface{})
				}
				configuration[configurationX.Category][configurationX.Name] = configurationX.Value
				continue
			}
			if configurationX.Category == global_var.ConfigurationCategory.CustomLookupField {
				if configuration[configurationX.Category] == nil {
					configuration[configurationX.Category] = make(map[string]interface{})
				}
				configuration[configurationX.Category][configurationX.Name] = configurationX.Value
				continue
			}
			if configurationX.Category == global_var.ConfigurationCategory.PaymentCityLedger {
				if configuration[configurationX.Category] == nil {
					configuration[configurationX.Category] = make(map[string]interface{})
				}
				configuration[configurationX.Category][configurationX.Name] = configurationX.Value
				continue
			}
			if configurationX.Category == global_var.ConfigurationCategory.RoomCosting {
				if configuration[configurationX.Category] == nil {
					configuration[configurationX.Category] = make(map[string]interface{})
				}
				configuration[configurationX.Category][configurationX.Name] = configurationX.Value
				continue
			}
			if configurationX.Category == global_var.ConfigurationCategory.RoomStatusColor {
				if configuration[configurationX.Category] == nil {
					configuration[configurationX.Category] = make(map[string]interface{})
				}
				configuration[configurationX.Category][configurationX.Name] = configurationX.Value
				continue
			}
			if configurationX.Category == global_var.ConfigurationCategory.ServiceCCMS {
				if configuration[configurationX.Category] == nil {
					configuration[configurationX.Category] = make(map[string]interface{})
				}
				configuration[configurationX.Category][configurationX.Name] = configurationX.Value
				continue
			}
			if configurationX.Category == global_var.ConfigurationCategory.SubDepartment {
				if configuration[configurationX.Category] == nil {
					configuration[configurationX.Category] = make(map[string]interface{})
				}
				configuration[configurationX.Category][configurationX.Name] = configurationX.Value
				continue
			}
			if configurationX.Category == global_var.ConfigurationCategory.GlobalDepartment {
				if configuration[configurationX.Category] == nil {
					configuration[configurationX.Category] = make(map[string]interface{})
				}
				configuration[configurationX.Category][configurationX.Name] = configurationX.Value
				continue
			}
			if configurationX.Category == global_var.ConfigurationCategory.GlobalSubDepartment {
				if configuration[configurationX.Category] == nil {
					configuration[configurationX.Category] = make(map[string]interface{})
				}
				configuration[configurationX.Category][configurationX.Name] = configurationX.Value
				continue
			}
			if configurationX.Category == global_var.ConfigurationCategory.GlobalJournalAccountSubGroup {
				if configuration[configurationX.Category] == nil {
					configuration[configurationX.Category] = make(map[string]interface{})
				}
				configuration[configurationX.Category][configurationX.Name] = configurationX.Value
				continue
			}
			if configurationX.Category == global_var.ConfigurationCategory.Inventory {
				if configuration[configurationX.Category] == nil {
					configuration[configurationX.Category] = make(map[string]interface{})
				}
				configuration[configurationX.Category][configurationX.Name] = configurationX.Value
				continue
			}

			if configurationX.Category == global_var.ConfigurationCategory.WeekendDay {
				if configuration[configurationX.Category] == nil {
					configuration[configurationX.Category] = make(map[string]interface{})
				}
				configuration[configurationX.Category][configurationX.Name] = configurationX.Value
				continue
			}

			if configurationX.Category == global_var.ConfigurationCategory.Invoice {
				if configuration[configurationX.Category] == nil {
					configuration[configurationX.Category] = make(map[string]interface{})
				}
				configuration[configurationX.Category][configurationX.Name] = configurationX.Value
				continue
			}

			if configurationX.Category == global_var.ConfigurationCategory.Other {
				if configuration[configurationX.Category] == nil {
					configuration[configurationX.Category] = make(map[string]interface{})
				}
				configuration[configurationX.Category][configurationX.Name] = configurationX.Value
				continue
			}
			if configurationX.Category == global_var.ConfigurationCategory.Reservation {
				if configuration[configurationX.Category] == nil {
					configuration[configurationX.Category] = make(map[string]interface{})
				}
				configuration[configurationX.Category][configurationX.Name] = configurationX.Value
				continue
			}
			if configurationX.Category == global_var.ConfigurationCategory.Folio {
				if configuration[configurationX.Category] == nil {
					configuration[configurationX.Category] = make(map[string]interface{})
				}
				configuration[configurationX.Category][configurationX.Name] = configurationX.Value
				continue
			}
			if configurationX.Category == global_var.ConfigurationCategory.FloorPlan {
				if configuration[configurationX.Category] == nil {
					configuration[configurationX.Category] = make(map[string]interface{})
				}
				configuration[configurationX.Category][configurationX.Name] = configurationX.Value
				continue
			}
			if configurationX.Category == global_var.ConfigurationCategory.OtherForm {
				if configuration[configurationX.Category] == nil {
					configuration[configurationX.Category] = make(map[string]interface{})
				}
				configuration[configurationX.Category][configurationX.Name] = configurationX.Value
				continue
			}
			if configurationX.Category == global_var.ConfigurationCategory.OtherHK {
				if configuration[configurationX.Category] == nil {
					configuration[configurationX.Category] = make(map[string]interface{})
				}
				configuration[configurationX.Category][configurationX.Name] = configurationX.Value
				continue
			}
			if configurationX.Category == global_var.ConfigurationCategoryPOS.Payment {
				if configuration[configurationX.Category] == nil {
					configuration[configurationX.Category] = make(map[string]interface{})
				}
				configuration[configurationX.Category][configurationX.Name] = configurationX.Value
				continue
			}
			if configurationX.Category == global_var.ConfigurationCategoryPOS.General {
				if configuration[configurationX.Category] == nil {
					configuration[configurationX.Category] = make(map[string]interface{})
				}
				configuration[configurationX.Category][configurationX.Name] = configurationX.Value
				continue
			}
			if configurationX.Category == global_var.ConfigurationCategoryCAMS.PurchaseRequestApp {
				if configuration[configurationX.Category] == nil {
					configuration[configurationX.Category] = make(map[string]interface{})
				}
				configuration[configurationX.Category][configurationX.Name] = configurationX.Value
				continue
			}

			if configurationX.Category == global_var.ConfigurationCategory.DefaultVariable {
				if configuration[configurationX.Category] == nil {
					configuration[configurationX.Category] = make(map[string]interface{})
				}
				configuration[configurationX.Category][configurationX.Name] = configurationX.Value
				continue
			}

			if configurationX.Category == global_var.ConfigurationCategory.DayendClosed {
				if configuration[configurationX.Category] == nil {
					configuration[configurationX.Category] = make(map[string]interface{})
				}
				configuration[configurationX.Category][configurationX.Name] = configurationX.Value
				continue
			}

			if configurationX.Category == global_var.ConfigurationCategoryPOS.TableView {
				if configuration[configurationX.Category] == nil {
					configuration[configurationX.Category] = make(map[string]interface{})
				}
				configuration[configurationX.Category][configurationX.Name] = configurationX.Value
				continue
			}
		}
		programConfiguration := global_var.TProgramConfiguration{
			AutoImportJournal:       general.StrToBool(configuration[global_var.ConfigurationCategory.DayendClosed][global_var.ConfigurationName.AutoImportJournal].(string)),
			SubDepartmentAllCCAdmin: configuration[global_var.ConfigurationCategory.Accounting][global_var.ConfigurationName.SubDepartmentAllCCAdmin].(string),
			FilterRateByMarket:      general.StrToBool(configuration[global_var.ConfigurationCategory.Reservation][global_var.ConfigurationName.FilterRateByMarket].(string)),
			AlwaysShowPublishRate:   general.StrToBool(configuration[global_var.ConfigurationCategory.Reservation][global_var.ConfigurationName.AlwaysShowPublishRate].(string)),
			FilterRateByCompany:     general.StrToBool(configuration[global_var.ConfigurationCategory.Reservation][global_var.ConfigurationName.FilterRateByCompany].(string)),

			ShowRate:                           general.StrToBool(configuration[global_var.ConfigurationCategory.Folio][global_var.ConfigurationName.ShowRate].(string)),
			SundayAsWeekend:                    general.StrToBool(configuration[global_var.ConfigurationCategory.WeekendDay][global_var.ConfigurationName.SundayAsWeekend].(string)),
			SaturdayAsWeekend:                  general.StrToBool(configuration[global_var.ConfigurationCategory.WeekendDay][global_var.ConfigurationName.SaturdayAsWeekend].(string)),
			FridayAsWeekend:                    general.StrToBool(configuration[global_var.ConfigurationCategory.WeekendDay][global_var.ConfigurationName.FridayAsWeekend].(string)),
			AutomaticCreateInvoiceCLAtCheckOut: general.StrToBool(configuration[global_var.ConfigurationCategory.Accounting][global_var.ConfigurationName.AutomaticCreateInvoiceCLAtCheckOut].(string)),
			IsCalculateAllRoomRevenueSubGroup:  general.StrToBool(configuration[global_var.ConfigurationCategory.Other][global_var.ConfigurationName.CalculateAllRoomRevenueSubGroup].(string)),
			AllowZeroAmount:                    general.StrToBool(configuration[global_var.ConfigurationCategory.Folio][global_var.ConfigurationName.AllowZeroAmount].(string)),
			ShowComplimentOnCashierReport:      general.StrToBool(configuration[global_var.ConfigurationCategory.Other][global_var.ConfigurationName.ShowComplimentOnCashierReport].(string)),
			ShowTransferOnCashierReport:        general.StrToBool(configuration[global_var.ConfigurationCategory.Other][global_var.ConfigurationName.ShowTransferOnCashierReport].(string)),
			UseChildRate:                       general.StrToBool(configuration[global_var.ConfigurationCategory.General][global_var.ConfigurationName.UseChildRate].(string)),
			IsRoomByName:                       general.StrToBool(configuration[global_var.ConfigurationCategory.General][global_var.ConfigurationName.IsRoomByName].(string)),
			PostDiscount:                       general.StrToBool(configuration[global_var.ConfigurationCategory.Reservation][global_var.ConfigurationName.PostDiscount].(string)),
			IsCompanyPRApplyPriceMoreThanOne:   general.StrToBool(configuration[global_var.ConfigurationCategoryCAMS.PurchaseRequestApp][global_var.ConfigurationName.IsCompanyPRApplyPriceMoreThanOne].(string)),
			ReceiveStockAPTwoDigitDecimal:      general.StrToBool(configuration[global_var.ConfigurationCategory.Inventory][global_var.ConfigurationName.ReceiveStockAPTwoDigitDecimal].(string)),
			CostingMethod:                      configuration[global_var.ConfigurationCategory.Inventory][global_var.ConfigurationName.CostingMethod].(string),
			CompanyTypeExpedition:              configuration[global_var.ConfigurationCategory.Other][global_var.ConfigurationName.CompanyTypeExpedition].(string),
			CompanyTypeSupplier:                configuration[global_var.ConfigurationCategory.Other][global_var.ConfigurationName.CompanyTypeSupplier].(string),
			CompanyTypeTravelAgent:             configuration[global_var.ConfigurationCategory.Other][global_var.ConfigurationName.CompanyTypeTravelAgent].(string),
			CheckOutLimit:                      configuration[global_var.ConfigurationCategory.Reservation][global_var.ConfigurationName.CheckOutLimit].(string),
			Timezone:                           configuration[global_var.ConfigurationCategory.General][global_var.ConfigurationName.Timezone].(string),
			//Template
			ProformaInvoiceDetail: configuration[global_var.ConfigurationCategory.OtherForm][global_var.ConfigurationName.ProformaInvoiceDetail].(string),
			FolioFooter:           configuration[global_var.ConfigurationCategory.Folio][global_var.ConfigurationName.FolioFooter].(string),
			DefaultFolio:          configuration[global_var.ConfigurationCategory.Folio][global_var.ConfigurationName.DefaultFolio].(string),

			//ChanneManager
			CCMSSMReservationAsAllotment:  general.StrToBool(configuration[global_var.ConfigurationCategory.ServiceCCMS][global_var.ConfigurationName.CCMSSMReservationAsAllotment].(string)),
			CCMSSMSynchronizeReservation:  general.StrToBool(configuration[global_var.ConfigurationCategory.ServiceCCMS][global_var.ConfigurationName.CCMSSMSynchronizeReservation].(string)),
			CCMSSMSynchronizeAvailability: general.StrToBool(configuration[global_var.ConfigurationCategory.ServiceCCMS][global_var.ConfigurationName.CCMSSMSynchronizeAvailability].(string)),
			CCMSSMSynchronizeRate:         general.StrToBool(configuration[global_var.ConfigurationCategory.ServiceCCMS][global_var.ConfigurationName.CCMSSMSynchronizeRate].(string)),
			CCMSVendor:                    configuration[global_var.ConfigurationCategory.ServiceCCMS][global_var.ConfigurationName.CCMSVendor].(string),
			CCMSSMUser:                    configuration[global_var.ConfigurationCategory.ServiceCCMS][global_var.ConfigurationName.CCMSSMUser].(string),
			CCMSSMPassword:                configuration[global_var.ConfigurationCategory.ServiceCCMS][global_var.ConfigurationName.CCMSSMPassword].(string),
			CCMSSMRequestorID:             configuration[global_var.ConfigurationCategory.ServiceCCMS][global_var.ConfigurationName.CCMSSMRequestorID].(string),
			CCMSSMHotelCode:               configuration[global_var.ConfigurationCategory.ServiceCCMS][global_var.ConfigurationName.CCMSSMHotelCode].(string),
			CCMSSMWSDL:                    configuration[global_var.ConfigurationCategory.ServiceCCMS][global_var.ConfigurationName.CCMSSMWSDL].(string),
		}

		globalAccount := global_var.TGlobalAccount{
			Breakfast:                         configuration[global_var.ConfigurationCategory.GlobalAccount][global_var.ConfigurationName.AccountBreakfast].(string),
			APCommission:                      configuration[global_var.ConfigurationCategory.GlobalAccount][global_var.ConfigurationName.AccountAPCommission].(string),
			RoomCharge:                        configuration[global_var.ConfigurationCategory.GlobalAccount][global_var.ConfigurationName.AccountRoomCharge].(string),
			APRefundDeposit:                   configuration[global_var.ConfigurationCategory.GlobalAccount][global_var.ConfigurationName.AccountAPRefundDeposit].(string),
			CreditCardAdm:                     configuration[global_var.ConfigurationCategory.GlobalAccount][global_var.ConfigurationName.AccountCreditCardAdm].(string),
			TransferDepositReservation:        configuration[global_var.ConfigurationCategory.GlobalAccount][global_var.ConfigurationName.AccountTransferDepositReservation].(string),
			Service:                           configuration[global_var.ConfigurationCategory.GlobalAccount][global_var.ConfigurationName.AccountService].(string),
			Tax:                               configuration[global_var.ConfigurationCategory.GlobalAccount][global_var.ConfigurationName.AccountTax].(string),
			CancellationFee:                   configuration[global_var.ConfigurationCategory.GlobalAccount][global_var.ConfigurationName.AccountCancellationFee].(string),
			NoShow:                            configuration[global_var.ConfigurationCategory.GlobalAccount][global_var.ConfigurationName.AccountNoShow].(string),
			Cash:                              configuration[global_var.ConfigurationCategory.GlobalAccount][global_var.ConfigurationName.AccountCash].(string),
			CityLedger:                        configuration[global_var.ConfigurationCategory.GlobalAccount][global_var.ConfigurationName.AccountCityLedger].(string),
			TransferDepositReservationToFolio: configuration[global_var.ConfigurationCategory.GlobalAccount][global_var.ConfigurationName.AccountTransferDepositReservationToFolio].(string),
			Telephone:                         configuration[global_var.ConfigurationCategory.GlobalAccount][global_var.ConfigurationName.AccountTelephone].(string),
			TransferCharge:                    configuration[global_var.ConfigurationCategory.GlobalAccount][global_var.ConfigurationName.AccountTransferCharge].(string),
			TransferPayment:                   configuration[global_var.ConfigurationCategory.GlobalAccount][global_var.ConfigurationName.AccountTransferPayment].(string),
			ExtraBed:                          configuration[global_var.ConfigurationCategory.GlobalAccount][global_var.ConfigurationName.AccountExtraBed].(string),
			VoucherCompliment:                 configuration[global_var.ConfigurationCategory.GlobalAccount][global_var.ConfigurationName.AccountVoucherCompliment].(string),
			Voucher:                           configuration[global_var.ConfigurationCategory.GlobalAccount][global_var.ConfigurationName.AccountVoucher].(string),
		}

		globalSubDepartment := global_var.TGlobalSubDepartment{
			Accounting:   configuration[global_var.ConfigurationCategory.GlobalSubDepartment][global_var.ConfigurationName.SDAccounting].(string),
			FrontOffice:  configuration[global_var.ConfigurationCategory.GlobalSubDepartment][global_var.ConfigurationName.SDFrontOffice].(string),
			HouseKeeping: configuration[global_var.ConfigurationCategory.GlobalSubDepartment][global_var.ConfigurationName.SDHouseKeeping].(string),
			Banquet:      configuration[global_var.ConfigurationCategory.GlobalSubDepartment][global_var.ConfigurationName.SDBanquet].(string),
		}

		globalDepartment := global_var.TGlobalDepartment{
			RoomDivision:  configuration[global_var.ConfigurationCategory.GlobalDepartment][global_var.ConfigurationName.DFoodBeverage].(string),
			Minor:         configuration[global_var.ConfigurationCategory.GlobalDepartment][global_var.ConfigurationName.DMinor].(string),
			Miscellaneous: configuration[global_var.ConfigurationCategory.GlobalDepartment][global_var.ConfigurationName.DMiscellaneous].(string),
			Banquet:       configuration[global_var.ConfigurationCategory.GlobalDepartment][global_var.ConfigurationName.DBanquet].(string),
			FoodBeverage:  configuration[global_var.ConfigurationCategory.GlobalDepartment][global_var.ConfigurationName.DFoodBeverage].(string),
		}
		globalJournalAccount := global_var.TGlobalJournalAccount{
			OverShortAsIncome:        configuration[global_var.ConfigurationCategory.GlobalJournalAccount][global_var.ConfigurationName.JAOverShortAsIncome].(string),
			OverShortAsExpense:       configuration[global_var.ConfigurationCategory.GlobalJournalAccount][global_var.ConfigurationName.JAOverShortAsExpense].(string),
			ServiceRevenue:           configuration[global_var.ConfigurationCategory.GlobalJournalAccount][global_var.ConfigurationName.JAServiceRevenue].(string),
			ExpensePurchasingTax:     configuration[global_var.ConfigurationCategory.GlobalJournalAccount][global_var.ConfigurationNameCAMS.JAPurchasingTax].(string),
			ExpenseShipping:          configuration[global_var.ConfigurationCategory.GlobalJournalAccount][global_var.ConfigurationNameCAMS.JAPurchasingShipping].(string),
			IncomePurchasingDiscount: configuration[global_var.ConfigurationCategory.GlobalJournalAccount][global_var.ConfigurationNameCAMS.JAPurchasingDiscount].(string),
			IncomeReturnStock:        configuration[global_var.ConfigurationCategory.GlobalJournalAccount][global_var.ConfigurationNameCAMS.JAIncomeReturnStock].(string),
			ExpenseReturnStock:       configuration[global_var.ConfigurationCategory.GlobalJournalAccount][global_var.ConfigurationNameCAMS.JAExpenseReturnStock].(string),
			APVoucher:                configuration[global_var.ConfigurationCategory.GlobalJournalAccount][global_var.ConfigurationName.JAAPVoucher].(string),
			ExpenseCreditCardAdm:     configuration[global_var.ConfigurationCategory.GlobalJournalAccount][global_var.ConfigurationName.JACreditCardAdm].(string),
			ExpenseInvoiceDiscount:   configuration[global_var.ConfigurationCategory.GlobalJournalAccount][global_var.ConfigurationName.JAInvoiceDiscount].(string),
			IncomeVoucherExpire:      configuration[global_var.ConfigurationCategory.GlobalJournalAccount][global_var.ConfigurationName.JAIncomeVoucherExpire].(string),
			GuestDepositReservation:  configuration[global_var.ConfigurationCategory.GlobalJournalAccount][global_var.ConfigurationName.JAGuestDepositReservation].(string),
			GuestDeposit:             configuration[global_var.ConfigurationCategory.GlobalJournalAccount][global_var.ConfigurationName.JAGuestDeposit].(string),
			GuestLedger:              configuration[global_var.ConfigurationCategory.GlobalJournalAccount][global_var.ConfigurationName.JAGuestLedger].(string),
			ProfitLossBeginningYear:  configuration[global_var.ConfigurationCategory.GlobalJournalAccount][global_var.ConfigurationName.JAPLBeginningYear].(string),
			ProfitLossCurrentYear:    configuration[global_var.ConfigurationCategory.GlobalJournalAccount][global_var.ConfigurationName.JAPLCurrentYear].(string),
			ProfitLossCurrency:       configuration[global_var.ConfigurationCategory.GlobalJournalAccount][global_var.ConfigurationName.JAPLCurrency].(string),
		}

		globalJournalAccountSubGroup := global_var.TGlobalJournalAccountSubGroup{
			AccumulatedDepreciation: configuration[global_var.ConfigurationCategory.GlobalJournalAccountSubGroup][global_var.ConfigurationName.JASGAccmDepreciation].(string),
			AccountPayable:          configuration[global_var.ConfigurationCategory.GlobalJournalAccountSubGroup][global_var.ConfigurationName.JASGAccountPayable].(string),
			Amortization:            configuration[global_var.ConfigurationCategory.GlobalJournalAccountSubGroup][global_var.ConfigurationName.JASGAmortization].(string),
			FixedAsset:              configuration[global_var.ConfigurationCategory.GlobalJournalAccountSubGroup][global_var.ConfigurationName.JASGFixedAsset].(string),
			Depreciation:            configuration[global_var.ConfigurationCategory.GlobalJournalAccountSubGroup][global_var.ConfigurationName.JASGDepreciation].(string),
			ManagementFee:           configuration[global_var.ConfigurationCategory.GlobalJournalAccountSubGroup][global_var.ConfigurationName.JASGManagementFee].(string),
			LoanInterest:            configuration[global_var.ConfigurationCategory.GlobalJournalAccountSubGroup][global_var.ConfigurationName.JASGLoanInterest].(string),
			Inventory:               configuration[global_var.ConfigurationCategory.GlobalJournalAccountSubGroup][global_var.ConfigurationName.JASGInventory].(string),
			IncomeTax:               configuration[global_var.ConfigurationCategory.GlobalJournalAccountSubGroup][global_var.ConfigurationName.JASGIncomeTax].(string),
		}

		var accountGroup []map[string]interface{}
		DB.Table(db_var.TableName.ConstJournalAccountGroup).Select("code", "name").Scan(&accountGroup)
		jaAccountGroup := func(code string, list []map[string]interface{}) string {
			for _, data := range list {
				if data["code"].(string) == code {
					return data["name"].(string)
				}
			}
			return ""
		}

		specialProduct := global_var.TSpecialProduct{
			FoodCode:            "[SPCFOOD]",
			BeverageCode:        "[SPCBVRG]",
			ProductCategoryCode: "[SPECIAL]",
		}

		globalJournalAccountGroupName := global_var.TGlobalJournalAccountGroupName{
			Assets:       jaAccountGroup("1", accountGroup),
			Liability:    jaAccountGroup("2", accountGroup),
			Equity:       jaAccountGroup("3", accountGroup),
			Income:       jaAccountGroup("4", accountGroup),
			Cost:         jaAccountGroup("5", accountGroup),
			Expense1:     jaAccountGroup("6", accountGroup),
			Expense2:     jaAccountGroup("7", accountGroup),
			OtherIncome:  jaAccountGroup("8", accountGroup),
			OtherExpense: jaAccountGroup("9", accountGroup),
		}

		Dataset = &global_var.TDataset{
			GlobalAccount:                 globalAccount,
			SpecialProduct:                specialProduct,
			Configuration:                 configuration,
			ProgramConfiguration:          programConfiguration,
			GlobalSubDepartment:           globalSubDepartment,
			GlobalDepartment:              globalDepartment,
			GlobalJournalAccount:          globalJournalAccount,
			GlobalJournalAccountSubGroup:  globalJournalAccountSubGroup,
			GlobalJournalAccountGroupName: globalJournalAccountGroupName,
		}
	}
	fmt.Println("dataset generated")
	return Dataset
}
