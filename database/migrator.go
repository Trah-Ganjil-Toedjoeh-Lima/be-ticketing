package database

import (
	"fmt"
	"github.com/frchandra/gmcgo/database/factory"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

type Migrator struct {
	db        *gorm.DB
	migration Migration
}

func NewMigrator(db *gorm.DB, migration Migration) *Migrator {
	return &Migrator{
		db:        db,
		migration: migration,
	}
}

func (mi *Migrator) RunMigration(option string) {
	var err error
	m := gormigrate.New(mi.db, gormigrate.DefaultOptions, mi.migration.Migrations)

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

	err = mi.RunFactory()
	if err != nil {
		panic(err)
	}
	fmt.Println("Migration did run successfully")

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
