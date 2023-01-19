package app

import (
	"fmt"
	"github.com/frchandra/gmcgo/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabase(appConfig *config.AppConfig) *gorm.DB {
	//appConfig := config.NewAppConfig()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", appConfig.DBHost, appConfig.DBUser, appConfig.DBPassword, appConfig.DBName, appConfig.DBPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("Failed on connecting to the migrator server")
	} else {
		fmt.Println("db connection established")
		fmt.Println("Using migrator " + db.Migrator().CurrentDatabase())
	}
	return db
}
