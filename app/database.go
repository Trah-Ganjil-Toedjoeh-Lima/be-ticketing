package app

import (
	"fmt"
	"github.com/frchandra/ticketing-gmcgo/config"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabase(appConfig *config.AppConfig, log *logrus.Logger) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", appConfig.DBHost, appConfig.DBUser, appConfig.DBPassword, appConfig.DBName, appConfig.DBPort)

	var gormConfig *gorm.Config
	if appConfig.IsProduction == "false" {
		gormConfig = &gorm.Config{Logger: logger.Default.LogMode(logger.Info)}
	} else {
		gormConfig = &gorm.Config{Logger: logger.Default.LogMode(logger.Error)}
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		log.Panic("failed on connecting to the database server")
	} else {
		log.Info("application is successfully connected to the database " + db.Migrator().CurrentDatabase())
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(appConfig.DBMaxIdleConnection)
	sqlDB.SetMaxOpenConns(appConfig.DBMaxOpenConnection)
	sqlDB.SetConnMaxLifetime(appConfig.DBConnectionMaxLifeMinute)

	return db
}
