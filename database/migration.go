package database

import (
	"github.com/frchandra/gmcgo/app/model"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

type Migration struct {
	Migrations []*gormigrate.Migration
}

func NewMigration() Migration {
	migration := []*gormigrate.Migration{
		{
			ID: "init",
			Migrate: func(tx *gorm.DB) error {
				return tx.Debug().AutoMigrate()
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Debug().Migrator().DropTable()
			},
		},
		{
			ID: "create_user_table",
			Migrate: func(tx *gorm.DB) error {
				return tx.Debug().AutoMigrate(model.User{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Debug().Migrator().DropTable(model.User{})
			},
		},
		{
			ID: "create_seats_table",
			Migrate: func(tx *gorm.DB) error {
				return tx.Debug().AutoMigrate(model.Seat{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Debug().Migrator().DropTable(model.Seat{})
			},
		},
	}

	return Migration{
		Migrations: migration,
	}
}
