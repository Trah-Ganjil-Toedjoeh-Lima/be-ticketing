package database

import (
	"fmt"
	"github.com/frchandra/gmcgo/config"
	"github.com/frchandra/gmcgo/database/factory"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Migrator struct {
	Database   *gorm.DB
	Migrations Migration
}

func NewMigrator(appConfig *config.AppConfig, migration Migration) *Migrator {
	db, _ := initializeDb(appConfig)
	return &Migrator{
		Database:   db,
		Migrations: migration,
	}
}

func initializeDb(appConfig *config.AppConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", appConfig.DBHost, appConfig.DBUser, appConfig.DBPassword, appConfig.DBName, appConfig.DBPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("Failed on connecting to the migrator server")
	} else {
		fmt.Println("Database connection established")
		fmt.Println("Using migrator " + db.Migrator().CurrentDatabase())
	}
	return db, err
}

func (this *Migrator) RunMigration(option string) {
	var err error
	m := gormigrate.New(this.Database, gormigrate.DefaultOptions, this.Migrations.Migrations)

	fmt.Println("option " + option + " is chosen")

	switch option {
	case "migrate:fresh":
		err = m.RollbackTo("init")
		if err != nil {
			panic(err)
		}
		err = m.Migrate()
	default:
		panic("option " + option + " unknown")
	}

	err = this.RunFactory()
	if err != nil {
		panic(err)
	}
	fmt.Println("Migration did run successfully")

}

func (this *Migrator) GetFactory() []factory.Factory {
	return []factory.Factory{
		factory.NewUserFactory(this.Database),
	}
}

func (this *Migrator) RunFactory() error {
	for _, seeder := range this.GetFactory() {
		err := seeder.RunFactory(3)
		if err != nil {
			return err
		}
	}
	return nil
}
