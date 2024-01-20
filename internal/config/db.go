package config

import (
	"chs/config"
	"chs/internal/general"
	"chs/internal/global_var"
	"chs/internal/init/dataset"
	"chs/internal/utils/cache"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var gormConfig = &gorm.Config{
	NamingStrategy: schema.NamingStrategy{
		SingularTable: true,
	},
}

// DatabaseConfig holds the configuration for a client database
type DatabaseConfig struct {
	Host         string
	Port         int
	Username     string
	Password     string
	Database     string
	MaxOpenConns int
	MaxIdleConns int
}

type ConnectionInfo struct {
	Username string
	Database string
	Host     string
	Port     int
}

// DBPool is a custom connection pool
type TDBPool struct {
	companyDataConfiguration []CompanyDataConfiguration
	// pool                 []*gorm.DB
	config        config.MySQL
	mu            sync.Mutex
	connectionIdx int // Index of the next connection to use
}

// Company represents a company record in the database
type Company struct {
	Name         string
	Subdomain    string
	DatabaseName string
	// Add other fields as needed
}

type CompanyDataConfiguration struct {
	CompanyCode        string
	CompanyID          uint64
	Name               string
	Domain             string
	Subdomain          string
	DatabaseName       string
	MaxUser            int64
	Rooms              int64
	DB                 *gorm.DB
	Dataset            *global_var.TDataset
	MxGeneral          *sync.RWMutex
	MxFolio            *sync.RWMutex
	MxRoomAvailability *sync.RWMutex
	MxStoreStock       *sync.RWMutex
	MxSubFolio         *sync.RWMutex
	MxDayendClose      *sync.RWMutex
	MxReceiving        *sync.RWMutex
	MxJournal          *sync.RWMutex
	MxPurchaseRequest  *sync.RWMutex
}

var DBPool *TDBPool
var mainDB *gorm.DB
var ctx = context.Background()

func InitDB(configX *config.Config, appLogger *otelzap.Logger) (*gorm.DB, *TDBPool, error) {
	dbConfig := configX.MySQL
	global_var.DatabaseInfo.HostName = dbConfig.MySqlHost
	global_var.DatabaseInfo.Port = dbConfig.MySqlPort
	global_var.DatabaseInfo.UserName = dbConfig.MySqlUser
	global_var.DatabaseInfo.Password = dbConfig.MySqlPassword
	global_var.DatabaseInfo.DatabaseName = dbConfig.MySqlDatabase
	DecryptedPassword, _ := general.DecryptString(global_var.EncryptKey, global_var.DatabaseInfo.Password)
	DecryptedDatabaseName, _ := general.DecryptString(global_var.EncryptKey, global_var.DatabaseInfo.DatabaseName)

	// pass, _ := general.EncryptString(global_var.EncryptKey, "cakratendados")
	// db, _ := general.EncryptString(global_var.EncryptKey, "db_chs_cloud")

	// fmt.Println("db", db)
	// fmt.Println("pass", pass)

	dbConfig.MySqlPassword = DecryptedPassword
	// if configX.Server.MultiDatabase {
	DecryptedDatabaseName = "db_extranet"
	if configX.Server.Staging {
		DecryptedDatabaseName = "db_staging_extranet"
	}
	// }
	connection := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=UTC",
		global_var.DatabaseInfo.UserName,
		DecryptedPassword,
		global_var.DatabaseInfo.HostName,
		global_var.DatabaseInfo.Port,
		DecryptedDatabaseName,
	)
	var err error
	mainDB, err = gorm.Open(mysql.Open(connection), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},

		Logger: debugLogger(configX),
	})
	if err != nil {
		appLogger.Error("Connection Failed to Open", zap.Error(err))
	}

	if configX.Server.MultiDatabase {
		fmt.Println("masuk Multi")
		global_var.DBMain = mainDB
	} else {
		fmt.Println("masuk Single")
		global_var.Db = mainDB
	}

	if configX.Server.MultiDatabase {
		// Create a new database connection pool for multidatabase
		dbPool, err := NewDBPool(10, configX)
		if err != nil {
			return nil, nil, err
		}

		DBPool = dbPool
		return mainDB, dbPool, nil
	}

	return mainDB, nil, nil
}

func debugLogger(config *config.Config) gormLogger.Interface {
	if config.Logger.LogFileEnabled {
		return nil
	}
	// return logger.Default.LogMode(logger.Info)
	return nil
}

// NewDBPool creates a new connection pool
func NewDBPool(size int, config *config.Config) (*TDBPool, error) {
	var companyDataConfiguration []CompanyDataConfiguration
	type CompanyStruct struct {
		CompanyCode, DatabaseName, Subdomain, Domain string
		CompanyID                                    uint64
		MaxUser, Rooms                               int64
	}
	var Company []CompanyStruct
	if err := mainDB.Table("company").Select("company.company_code as CompanyCode, company.company_id, database_name,domain,subdomain, max_user, rooms").
		Joins("LEFT JOIN company_database ON company.company_id = company_database.company_id").
		Joins("LEFT JOIN subscription ON company.company_id = subscription.company_id").
		// Where("subscription.start_date<=NOW()").
		// Where("subscription.end_date>NOW()").
		Scan(&Company).Error; err != nil {
		return nil, err
	}

	DecryptedPassword, _ := general.DecryptString(global_var.EncryptKey, global_var.DatabaseInfo.Password)
	for i := 0; i < len(Company); i++ {
		fmt.Println("count", i)
		cache.DataCache.Del(ctx, "database:", Company[i].Subdomain)
		cache.DataCache.Del(ctx, "database:", Company[i].Domain)
		if Company[i].Subdomain != "" || Company[i].Domain != "" {
			DBName := Company[i].DatabaseName
			if config.Server.Staging {
				DBName = "staging_" + Company[i].DatabaseName
			}
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=UTC",
				global_var.DatabaseInfo.UserName,
				DecryptedPassword,
				global_var.DatabaseInfo.HostName,
				global_var.DatabaseInfo.Port,
				DBName)
			db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			}})
			if err != nil {
				continue
			}
			if config.Jaeger.LogQuery {
				if err := db.Use(otelgorm.NewPlugin()); err != nil {
					panic(err)
				}
			}
			sqlDB, err := db.DB()
			if err != nil {
				continue
			}

			sqlDB.SetMaxOpenConns(config.MySQL.MaxOpenConns)
			sqlDB.SetMaxIdleConns(config.MySQL.MaxIdleConns)

			cache.DataCache.Del(ctx, Company[i].CompanyCode, "AUDIT_DATE")
			dataset := dataset.GenerateDataset(db)

			// Checking the connection status
			sqlDB, errXX := db.DB()
			if errXX != nil {
				fmt.Println("Error getting DB connection", errXX)
				return nil, errXX
			}

			errX := sqlDB.Ping()
			if errX != nil {
				fmt.Println("Database connection is not established.")
				return nil, errX
			}

			companyDataConfiguration = append(companyDataConfiguration, CompanyDataConfiguration{
				CompanyCode:        Company[i].CompanyCode,
				CompanyID:          Company[i].CompanyID,
				Domain:             Company[i].Domain,
				Subdomain:          Company[i].Subdomain,
				DatabaseName:       Company[i].DatabaseName,
				MaxUser:            Company[i].MaxUser,
				Rooms:              Company[i].Rooms,
				DB:                 db,
				Dataset:            dataset,
				MxGeneral:          &sync.RWMutex{},
				MxFolio:            &sync.RWMutex{},
				MxRoomAvailability: &sync.RWMutex{},
				MxStoreStock:       &sync.RWMutex{},
				MxSubFolio:         &sync.RWMutex{},
				MxDayendClose:      &sync.RWMutex{},
				MxReceiving:        &sync.RWMutex{},
				MxJournal:          &sync.RWMutex{},
				MxPurchaseRequest:  &sync.RWMutex{},
			})
			fmt.Println(Company[i].CompanyCode)
			fmt.Println(Company[i].Subdomain)
		}
	}

	return &TDBPool{
		companyDataConfiguration: companyDataConfiguration,
		// pool:   pool,
		config: config.MySQL,
	}, nil
}

// GetConnection retrieves a database connection from the pool
func (p *TDBPool) GetConnection(Hostname string, CompanyCode string) (*CompanyDataConfiguration, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Check if the connection information is cached in Redis
	RedisKey := Hostname
	if CompanyCode != "" {
		RedisKey = CompanyCode
	}
	if val, err := cache.DataCache.Get(ctx, "database:", RedisKey); err == nil {
		// If the connection information is cached, obtain the connection from the pool
		connInfo := p.config
		err = json.Unmarshal(val, &connInfo)
		if err != nil {
			return nil, err
		}

		// Find the connection with the matching database configuration
		for _, pConfig := range p.companyDataConfiguration {
			if pConfig.Subdomain == Hostname || pConfig.Domain == Hostname || pConfig.CompanyCode == CompanyCode {
				err = cache.DataCache.Set(ctx, "database:", RedisKey, connInfo, 0)
				if err != nil {
					return nil, err
				}
				global_var.Db = pConfig.DB
				return &pConfig, nil
			}

			// dbConfig := db.Config
			// // Compare the connection configuration with the cached configuration
			// if dbConfig.Dialector.Name() == "mysql" &&
			// 	dbConfig.Dialector.(*mysql.Dialector).DSN == fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=UTC", connInfo.MySqlUser, connInfo.MySqlPassword, connInfo.MySqlHost, connInfo.MySqlPort, connInfo.MySqlDatabase) {

			// 	// Set the maximum connection limits
			// 	sqlDB, err := db.DB()
			// 	if err != nil {
			// 		return nil, err
			// 	}
			// 	sqlDB.SetMaxOpenConns(connInfo.MaxOpenConns)
			// 	sqlDB.SetMaxIdleConns(connInfo.MaxIdleConns)

			// 	// Cache the connection information in Redis
			// 	err = cache.DataCache.Set(ctx, "database:", Hostname, connInfo, 0)
			// 	if err != nil {
			// 		return nil, err
			// 	}

			// 	return db, nil
			// }
		}
		// return nil, fmt.Errorf("matching database connection not found in the pool")
	}

	// If the connection information is not cached, retrieve it from the original configuration
	connInfo := p.config
	// Retrieve the database name for the company
	var DatabaseName string
	query := mainDB.Table("company").Select("database_name").
		Joins("LEFT JOIN company_database ON company.company_id = company_database.company_id").
		Joins("LEFT JOIN subscription ON company.company_id = subscription.company_id")
		// Where("subscription.start_date<=NOW()").
		// Where("subscription.end_date>NOW()")

	if CompanyCode != "" {
		query.Where("company.company_code", CompanyCode)
	} else if Hostname != "" {
		query.Where("company.subdomain=?", Hostname).
			Or("company.domain=?", Hostname)
	}
	if err := query.Take(&DatabaseName).Error; err != nil {
		return nil, err
	}

	connInfo.MySqlPassword, _ = general.DecryptString(global_var.EncryptKey, global_var.DatabaseInfo.Password)
	connInfo.MySqlDatabase = DatabaseName
	if global_var.Config.Server.Staging {
		connInfo.MySqlDatabase = "staging_" + DatabaseName
	}
	// Obtain a connection from the pool that matches the configuration
	for _, pConfig := range p.companyDataConfiguration {
		if pConfig.Subdomain == Hostname || pConfig.Domain == Hostname || pConfig.CompanyCode == CompanyCode {
			// Checking the connection status
			if pConfig.DB == nil {
				return nil, fmt.Errorf("database connection not found in the pool")
			}
			sqlDB, errXX := pConfig.DB.DB()
			if errXX != nil {
				fmt.Println("Error getting DB connection", errXX)
				return nil, errXX
			}

			errX := sqlDB.Ping()
			if errX != nil {
				fmt.Println("Database connection is not established.")
				return nil, errX
			}

			err := cache.DataCache.Set(ctx, "database:", RedisKey, connInfo, 0)
			if err != nil {
				return nil, err
			}
			global_var.Db = pConfig.DB
			return &pConfig, nil
		}

		// if db == nil {
		// 	continue
		// }
		// dbConfig := db.Config
		// // Compare the connection configuration with the cached configuration
		// if dbConfig.Dialector.Name() == "mysql" &&
		// 	dbConfig.Dialector.(*mysql.Dialector).DSN == fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=UTC", connInfo.MySqlUser, connInfo.MySqlPassword, connInfo.MySqlHost, connInfo.MySqlPort, connInfo.MySqlDatabase) {

		// 	// Set the maximum connection limits
		// 	sqlDB, err := db.DB()
		// 	if err != nil {
		// 		return nil, err
		// 	}
		// 	sqlDB.SetMaxOpenConns(connInfo.MaxOpenConns)
		// 	sqlDB.SetMaxIdleConns(connInfo.MaxIdleConns)

		// 	// Cache the connection information in Redis
		// 	err = cache.DataCache.Set(ctx, "database:", Hostname, connInfo, 0)
		// 	if err != nil {
		// 		return nil, err
		// 	}

		// 	return db, nil

	}

	return nil, fmt.Errorf("matching database connection not found in the pool")
}

// ReleaseConnection releases a database connection back to the pool
func (p *TDBPool) ReleaseConnection(db *gorm.DB) {
	// Close the GORM DB instance
	sql, err := db.DB()
	if err != nil {
		// log.Println("Failed to close GORM DB instance:", err)
	}
	sql.Close()

	// If you have any other connection cleanup/reset tasks, you can perform them here

	// Note: You may want to consider implementing a connection reuse strategy
	// where you keep track of released connections and reuse them instead of
	// creating new ones every time. This can help optimize performance.

	// For now, we'll simply log a message to indicate that the connection has been released
	// log.Println("Released database connection to the pool")
}

func (p *TDBPool) RegenerateDataset(DB *gorm.DB, pConfig CompanyDataConfiguration) {
	// var Dataset *global_var.TDataset
	// configuration := make(map[string]map[string]interface{})
	// var DataOutput []db_var.Configuration
	// DB.Table(db_var.TableName.Configuration).Scan(&DataOutput)
	// if len(DataOutput) > 0 {
	// 	for _, configurationX := range DataOutput {
	// 		//bo0lean
	// 		if configurationX.Category == global_var.ConfigurationCategory.General {
	// 			if configuration[configurationX.Category] == nil {
	// 				configuration[configurationX.Category] = make(map[string]interface{})
	// 			}
	// 			configuration[configurationX.Category][configurationX.Name] = configurationX.Value
	// 			continue
	// 		}
	// 		if configurationX.Category == global_var.ConfigurationCategory.Accounting {
	// 			if configuration[configurationX.Category] == nil {
	// 				configuration[configurationX.Category] = make(map[string]interface{})
	// 			}
	// 			configuration[configurationX.Category][configurationX.Name] = configurationX.Value
	// 			continue
	// 		}
	// 		if configurationX.Category == global_var.ConfigurationCategory.AmountPreset {
	// 			if configuration[configurationX.Category] == nil {
	// 				configuration[configurationX.Category] = make(map[string]interface{})
	// 			}
	// 			configuration[configurationX.Category][configurationX.Name] = configurationX.Value
	// 			continue
	// 		}
	// 		if configurationX.Category == global_var.ConfigurationCategory.CompanyBankAccount {
	// 			if configuration[configurationX.Category] == nil {
	// 				configuration[configurationX.Category] = make(map[string]interface{})
	// 			}
	// 			configuration[configurationX.Category][configurationX.Name] = configurationX.Value
	// 			continue
	// 		}
	// 		if configurationX.Category == global_var.ConfigurationCategory.GlobalAccount {
	// 			if configuration[configurationX.Category] == nil {
	// 				configuration[configurationX.Category] = make(map[string]interface{})
	// 			}
	// 			configuration[configurationX.Category][configurationX.Name] = configurationX.Value
	// 			continue
	// 		}
	// 		if configurationX.Category == global_var.ConfigurationCategory.GlobalJournalAccount {
	// 			if configuration[configurationX.Category] == nil {
	// 				configuration[configurationX.Category] = make(map[string]interface{})
	// 			}
	// 			configuration[configurationX.Category][configurationX.Name] = configurationX.Value
	// 			continue
	// 		}
	// 		if configurationX.Category == global_var.ConfigurationCategory.PaymentAccount {
	// 			if configuration[configurationX.Category] == nil {
	// 				configuration[configurationX.Category] = make(map[string]interface{})
	// 			}
	// 			configuration[configurationX.Category][configurationX.Name] = configurationX.Value
	// 			continue
	// 		}
	// 		if configurationX.Category == global_var.ConfigurationCategory.Company {
	// 			if configuration[configurationX.Category] == nil {
	// 				configuration[configurationX.Category] = make(map[string]interface{})
	// 			}
	// 			configuration[configurationX.Category][configurationX.Name] = configurationX.Value
	// 			continue
	// 		}
	// 		if configurationX.Category == global_var.ConfigurationCategory.CustomField {
	// 			if configuration[configurationX.Category] == nil {
	// 				configuration[configurationX.Category] = make(map[string]interface{})
	// 			}
	// 			configuration[configurationX.Category][configurationX.Name] = configurationX.Value
	// 			continue
	// 		}
	// 		if configurationX.Category == global_var.ConfigurationCategory.CustomLookupField {
	// 			if configuration[configurationX.Category] == nil {
	// 				configuration[configurationX.Category] = make(map[string]interface{})
	// 			}
	// 			configuration[configurationX.Category][configurationX.Name] = configurationX.Value
	// 			continue
	// 		}
	// 		if configurationX.Category == global_var.ConfigurationCategory.PaymentCityLedger {
	// 			if configuration[configurationX.Category] == nil {
	// 				configuration[configurationX.Category] = make(map[string]interface{})
	// 			}
	// 			configuration[configurationX.Category][configurationX.Name] = configurationX.Value
	// 			continue
	// 		}
	// 		if configurationX.Category == global_var.ConfigurationCategory.RoomCosting {
	// 			if configuration[configurationX.Category] == nil {
	// 				configuration[configurationX.Category] = make(map[string]interface{})
	// 			}
	// 			configuration[configurationX.Category][configurationX.Name] = configurationX.Value
	// 			continue
	// 		}
	// 		if configurationX.Category == global_var.ConfigurationCategory.RoomStatusColor {
	// 			if configuration[configurationX.Category] == nil {
	// 				configuration[configurationX.Category] = make(map[string]interface{})
	// 			}
	// 			configuration[configurationX.Category][configurationX.Name] = configurationX.Value
	// 			continue
	// 		}
	// 		if configurationX.Category == global_var.ConfigurationCategory.ServiceCCMS {
	// 			if configuration[configurationX.Category] == nil {
	// 				configuration[configurationX.Category] = make(map[string]interface{})
	// 			}
	// 			configuration[configurationX.Category][configurationX.Name] = configurationX.Value
	// 			continue
	// 		}
	// 		if configurationX.Category == global_var.ConfigurationCategory.SubDepartment {
	// 			if configuration[configurationX.Category] == nil {
	// 				configuration[configurationX.Category] = make(map[string]interface{})
	// 			}
	// 			configuration[configurationX.Category][configurationX.Name] = configurationX.Value
	// 			continue
	// 		}
	// 		if configurationX.Category == global_var.ConfigurationCategory.GlobalDepartment {
	// 			if configuration[configurationX.Category] == nil {
	// 				configuration[configurationX.Category] = make(map[string]interface{})
	// 			}
	// 			configuration[configurationX.Category][configurationX.Name] = configurationX.Value
	// 			continue
	// 		}
	// 		if configurationX.Category == global_var.ConfigurationCategory.GlobalSubDepartment {
	// 			if configuration[configurationX.Category] == nil {
	// 				configuration[configurationX.Category] = make(map[string]interface{})
	// 			}
	// 			configuration[configurationX.Category][configurationX.Name] = configurationX.Value
	// 			continue
	// 		}
	// 		if configurationX.Category == global_var.ConfigurationCategory.GlobalJournalAccountSubGroup {
	// 			if configuration[configurationX.Category] == nil {
	// 				configuration[configurationX.Category] = make(map[string]interface{})
	// 			}
	// 			configuration[configurationX.Category][configurationX.Name] = configurationX.Value
	// 			continue
	// 		}
	// 	}
	// 	globalAccount := global_var.TGlobalAccount{
	// 		APCommission:                      configuration[global_var.ConfigurationCategory.GlobalAccount][global_var.ConfigurationName.AccountAPCommission].(string),
	// 		RoomCharge:                        configuration[global_var.ConfigurationCategory.GlobalAccount][global_var.ConfigurationName.AccountRoomCharge].(string),
	// 		APRefundDeposit:                   configuration[global_var.ConfigurationCategory.GlobalAccount][global_var.ConfigurationName.AccountAPRefundDeposit].(string),
	// 		CreditCardAdm:                     configuration[global_var.ConfigurationCategory.GlobalAccount][global_var.ConfigurationName.AccountCreditCardAdm].(string),
	// 		TransferDepositReservation:        configuration[global_var.ConfigurationCategory.GlobalAccount][global_var.ConfigurationName.AccountTransferDepositReservation].(string),
	// 		Service:                           configuration[global_var.ConfigurationCategory.GlobalAccount][global_var.ConfigurationName.AccountService].(string),
	// 		Tax:                               configuration[global_var.ConfigurationCategory.GlobalAccount][global_var.ConfigurationName.AccountTax].(string),
	// 		CancellationFee:                   configuration[global_var.ConfigurationCategory.GlobalAccount][global_var.ConfigurationName.AccountCancellationFee].(string),
	// 		NoShow:                            configuration[global_var.ConfigurationCategory.GlobalAccount][global_var.ConfigurationName.AccountNoShow].(string),
	// 		Cash:                              configuration[global_var.ConfigurationCategory.GlobalAccount][global_var.ConfigurationName.AccountCash].(string),
	// 		CityLedger:                        configuration[global_var.ConfigurationCategory.GlobalAccount][global_var.ConfigurationName.AccountCityLedger].(string),
	// 		TransferDepositReservationToFolio: configuration[global_var.ConfigurationCategory.GlobalAccount][global_var.ConfigurationName.AccountTransferDepositReservationToFolio].(string),
	// 		Telephone:                         configuration[global_var.ConfigurationCategory.GlobalAccount][global_var.ConfigurationName.AccountTelephone].(string),
	// 		TransferCharge:                    configuration[global_var.ConfigurationCategory.GlobalAccount][global_var.ConfigurationName.AccountTransferCharge].(string),
	// 		TransferPayment:                   configuration[global_var.ConfigurationCategory.GlobalAccount][global_var.ConfigurationName.AccountTransferPayment].(string),
	// 		ExtraBed:                          configuration[global_var.ConfigurationCategory.GlobalAccount][global_var.ConfigurationName.AccountExtraBed].(string),
	// 		VoucherCompliment:                 configuration[global_var.ConfigurationCategory.GlobalAccount][global_var.ConfigurationName.AccountVoucherCompliment].(string),
	// 		Voucher:                           configuration[global_var.ConfigurationCategory.GlobalAccount][global_var.ConfigurationName.AccountVoucher].(string),
	// 	}

	// 	globalSubDepartment := global_var.TGlobalSubDepartment{
	// 		Accounting:   configuration[global_var.ConfigurationCategory.GlobalSubDepartment][global_var.ConfigurationName.SDAccounting].(string),
	// 		FrontOffice:  configuration[global_var.ConfigurationCategory.GlobalSubDepartment][global_var.ConfigurationName.SDFrontOffice].(string),
	// 		HouseKeeping: configuration[global_var.ConfigurationCategory.GlobalSubDepartment][global_var.ConfigurationName.SDHouseKeeping].(string),
	// 		Banquet:      configuration[global_var.ConfigurationCategory.GlobalSubDepartment][global_var.ConfigurationName.SDBanquet].(string),
	// 	}

	// 	globalDepartment := global_var.TGlobalDepartment{
	// 		RoomDivision:  configuration[global_var.ConfigurationCategory.GlobalDepartment][global_var.ConfigurationName.DFoodBeverage].(string),
	// 		Minor:         configuration[global_var.ConfigurationCategory.GlobalDepartment][global_var.ConfigurationName.DMinor].(string),
	// 		Miscellaneous: configuration[global_var.ConfigurationCategory.GlobalDepartment][global_var.ConfigurationName.DMiscellaneous].(string),
	// 		Banquet:       configuration[global_var.ConfigurationCategory.GlobalDepartment][global_var.ConfigurationName.DBanquet].(string),
	// 		FoodBeverage:  configuration[global_var.ConfigurationCategory.GlobalDepartment][global_var.ConfigurationName.DFoodBeverage].(string),
	// 	}

	// 	globalJournalAccount := global_var.TGlobalJournalAccount{
	// 		APVoucher:               configuration[global_var.ConfigurationCategory.GlobalJournalAccount][global_var.ConfigurationName.JAAPVoucher].(string),
	// 		ExpenseCreditCardAdm:    configuration[global_var.ConfigurationCategory.GlobalJournalAccount][global_var.ConfigurationName.JACreditCardAdm].(string),
	// 		ExpenseInvoiceDiscount:  configuration[global_var.ConfigurationCategory.GlobalJournalAccount][global_var.ConfigurationName.JAInvoiceDiscount].(string),
	// 		IncomeVoucherExpire:     configuration[global_var.ConfigurationCategory.GlobalJournalAccount][global_var.ConfigurationName.JAIncomeVoucherExpire].(string),
	// 		GuestDepositReservation: configuration[global_var.ConfigurationCategory.GlobalJournalAccount][global_var.ConfigurationName.JAGuestDepositReservation].(string),
	// 		GuestDeposit:            configuration[global_var.ConfigurationCategory.GlobalJournalAccount][global_var.ConfigurationName.JAGuestDeposit].(string),
	// 		GuestLedger:             configuration[global_var.ConfigurationCategory.GlobalJournalAccount][global_var.ConfigurationName.JAGuestLedger].(string),
	// 		ProfitLossBeginningYear: configuration[global_var.ConfigurationCategory.GlobalJournalAccount][global_var.ConfigurationName.JAPLBeginningYear].(string),
	// 		ProfitLossCurrentYear:   configuration[global_var.ConfigurationCategory.GlobalJournalAccount][global_var.ConfigurationName.JAPLCurrentYear].(string),
	// 		ProfitLossCurrency:      configuration[global_var.ConfigurationCategory.GlobalJournalAccount][global_var.ConfigurationName.JAPLCurrency].(string),
	// 	}

	// 	globalJournalAccountSubGroup := global_var.TGlobalJournalAccountSubGroup{
	// 		AccumulatedDepreciation: configuration[global_var.ConfigurationCategory.GlobalJournalAccountSubGroup][global_var.ConfigurationName.JASGAccmDepreciation].(string),
	// 		AccountPayable:          configuration[global_var.ConfigurationCategory.GlobalJournalAccountSubGroup][global_var.ConfigurationName.JASGAccountPayable].(string),
	// 		Amortization:            configuration[global_var.ConfigurationCategory.GlobalJournalAccountSubGroup][global_var.ConfigurationName.JASGAmortization].(string),
	// 		FixedAsset:              configuration[global_var.ConfigurationCategory.GlobalJournalAccountSubGroup][global_var.ConfigurationName.JASGFixedAsset].(string),
	// 		Depreciation:            configuration[global_var.ConfigurationCategory.GlobalJournalAccountSubGroup][global_var.ConfigurationName.JASGDepreciation].(string),
	// 		ManagementFee:           configuration[global_var.ConfigurationCategory.GlobalJournalAccountSubGroup][global_var.ConfigurationName.JASGManagementFee].(string),
	// 		LoanInterest:            configuration[global_var.ConfigurationCategory.GlobalJournalAccountSubGroup][global_var.ConfigurationName.JASGLoanInterest].(string),
	// 		Inventory:               configuration[global_var.ConfigurationCategory.GlobalJournalAccountSubGroup][global_var.ConfigurationName.JASGInventory].(string),
	// 		IncomeTax:               configuration[global_var.ConfigurationCategory.GlobalJournalAccountSubGroup][global_var.ConfigurationName.JASGIncomeTax].(string),
	// 	}

	// 	var accountGroup []map[string]interface{}
	// 	DB.Table(db_var.TableName.ConstJournalAccountGroup).Select("code", "name").Scan(&accountGroup)
	// 	jaAccountGroup := func(code string, list []map[string]interface{}) string {
	// 		for _, data := range list {
	// 			if data["code"].(string) == code {
	// 				return data["name"].(string)
	// 			}
	// 		}
	// 		return ""
	// 	}

	// 	specialProduct := global_var.TSpecialProduct{
	// 		FoodCode:            "[SPCFOOD]",
	// 		BeverageCode:        "[SPCBVRG]",
	// 		ProductCategoryCode: "[SPECIAL]",
	// 	}

	// 	globalJournalAccountGroupName := global_var.TGlobalJournalAccountGroupName{
	// 		Assets:       jaAccountGroup("1", accountGroup),
	// 		Liability:    jaAccountGroup("2", accountGroup),
	// 		Equity:       jaAccountGroup("3", accountGroup),
	// 		Income:       jaAccountGroup("4", accountGroup),
	// 		Cost:         jaAccountGroup("5", accountGroup),
	// 		Expense1:     jaAccountGroup("6", accountGroup),
	// 		Expense2:     jaAccountGroup("7", accountGroup),
	// 		OtherIncome:  jaAccountGroup("8", accountGroup),
	// 		OtherExpense: jaAccountGroup("9", accountGroup),
	// 	}

	// 	Dataset = &global_var.TDataset{
	// 		GlobalAccount:                 globalAccount,
	// 		SpecialProduct:                specialProduct,
	// 		Configuration:                 configuration,
	// 		GlobalSubDepartment:           globalSubDepartment,
	// 		GlobalDepartment:              globalDepartment,
	// 		GlobalJournalAccount:          globalJournalAccount,
	// 		GlobalJournalAccountSubGroup:  globalJournalAccountSubGroup,
	// 		GlobalJournalAccountGroupName: globalJournalAccountGroupName,
	// 	}
	// }

	Dataset := dataset.GenerateDataset(DB)
	for i, pConfigX := range p.companyDataConfiguration {
		if pConfigX.Domain == pConfig.Domain || pConfigX.Subdomain == pConfig.Subdomain {
			p.companyDataConfiguration[i].Dataset = Dataset
			break
		}
	}
}

// Split SQL statements by semicolon
func (p *TDBPool) GenerateDataConfig(c *gin.Context, DatabaseName string) (*gorm.DB, error) {
	type CompanyStruct struct {
		CompanyCode, DatabaseName, Subdomain, Domain string
		CompanyID                                    uint64
		MaxUser, Rooms                               int64
	}
	var Company CompanyStruct
	if err := global_var.DBMain.Table("company").Select("company.company_code as CompanyCode, company.company_id, database_name,domain,subdomain, max_user, rooms").
		Joins("LEFT JOIN company_database ON company.company_id = company_database.company_id").
		Joins("LEFT JOIN subscription ON company.company_id = subscription.company_id").
		// Where("subscription.start_date<=NOW()").
		// Where("subscription.end_date>NOW()").
		Where("company_database.database_name", DatabaseName).
		Limit(1).
		Scan(&Company).Error; err != nil {
		return nil, err
	}

	if Company.Subdomain != "" || Company.Domain != "" {
		DecryptedPassword, _ := general.DecryptString(global_var.EncryptKey, global_var.DatabaseInfo.Password)
		// DecryptedDatabaseName, _ := general.DecryptString(global_var.EncryptKey, global_var.DatabaseInfo.DatabaseName)
		DBName := Company.DatabaseName
		if global_var.Config.Server.Staging {
			DBName = "staging_" + Company.DatabaseName
		}
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=UTC",
			global_var.DatabaseInfo.UserName,
			DecryptedPassword,
			global_var.DatabaseInfo.HostName,
			global_var.DatabaseInfo.Port,
			DBName)
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		}})
		if err != nil {
			return nil, err
		}
		sqlDB, err := db.DB()
		if err != nil {
			return nil, err
		}

		sqlDB.SetMaxOpenConns(global_var.Config.MySQL.MaxOpenConns)
		sqlDB.SetMaxIdleConns(global_var.Config.MySQL.MaxIdleConns)

		cache.DataCache.Del(c, Company.CompanyCode, "AUDIT_DATE")
		dataset := dataset.GenerateDataset(db)

		companyDataConfiguration := CompanyDataConfiguration{
			CompanyCode:        Company.CompanyCode,
			CompanyID:          Company.CompanyID,
			Domain:             Company.Domain,
			Subdomain:          Company.Subdomain,
			DatabaseName:       Company.DatabaseName,
			MaxUser:            Company.MaxUser,
			Rooms:              Company.Rooms,
			DB:                 db,
			Dataset:            dataset,
			MxGeneral:          &sync.RWMutex{},
			MxFolio:            &sync.RWMutex{},
			MxRoomAvailability: &sync.RWMutex{},
			MxStoreStock:       &sync.RWMutex{},
			MxSubFolio:         &sync.RWMutex{},
			MxDayendClose:      &sync.RWMutex{},
			MxReceiving:        &sync.RWMutex{},
			MxJournal:          &sync.RWMutex{},
			MxPurchaseRequest:  &sync.RWMutex{},
		}

		DBPool.companyDataConfiguration = append(DBPool.companyDataConfiguration, companyDataConfiguration)

		return db, nil
	}
	return nil, errors.New("Data not found")
}
