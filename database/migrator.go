package database

import (
	"fmt"
	"github.com/frchandra/ticketing-gmcgo/app/model"
	"github.com/frchandra/ticketing-gmcgo/database/factory"
	"gorm.io/gorm"
)

type Migrator struct {
	db *gorm.DB
}

func NewMigrator(db *gorm.DB) *Migrator {
	return &Migrator{
		db: db,
	}
}

func (mi *Migrator) RunMigration(option string) {
	if err := mi.db.Migrator().DropTable(&model.User{}, &model.Seat{}, &model.Transaction{}); err != nil {
		panic(err)
	}
	if err := mi.db.AutoMigrate(&model.User{}, &model.Seat{}, &model.Transaction{}); err != nil {
		panic(err)
	}
	if err := mi.RunFactory(); err != nil {
		panic(err)
	}
	fmt.Println("Migration run successfully")
}

func (mi *Migrator) GetFactory() []factory.Factory {
	return []factory.Factory{
		factory.NewUserFactory(mi.db),
		factory.NewSeatFactory(mi.db),
	}
}

func (mi *Migrator) RunFactory() error {
	for _, seeder := range mi.GetFactory() {
		err := seeder.RunFactory()
		if err != nil {
			return err
		}
	}
	return nil
}
