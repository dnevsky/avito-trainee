package postgresDB

import (
	configApp "avito-trainee/internal/config"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

type Config struct {
	Host        string
	Port        string
	Username    string
	Password    string
	DBName      string
	SSLMode     string
	Environment string
}

func NewPostgresDB(config Config) (*gorm.DB, error) {
	dbLogger := gormLogger.Default.LogMode(gormLogger.Silent)
	if config.Environment == configApp.DevEnv {
		dbLogger = gormLogger.Default.LogMode(gormLogger.Info)
	}

	pg := postgres.New(postgres.Config{
		PreferSimpleProtocol: true,
		DSN:                  fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", config.Host, config.Port, config.Username, config.Password, config.DBName, config.SSLMode),
	})

	ormDB, err := gorm.Open(pg, &gorm.Config{
		Logger:                                   dbLogger,
		DisableForeignKeyConstraintWhenMigrating: true,
		SkipDefaultTransaction:                   true,
		PrepareStmt:                              false,
	})
	if err != nil {
		return nil, err
	}

	ormDB.Statement.RaiseErrorOnNotFound = true

	return ormDB, nil
}

func CloseDB(ormDB *gorm.DB) {
	db, err := ormDB.DB()
	if err != nil {
		log.Fatalf(fmt.Sprintf("cant close the connector: %s", err))
	}

	err = db.Close()
	if err != nil {
		log.Fatalf(fmt.Sprintf("cant close the connector: %s", err))
	}
}
