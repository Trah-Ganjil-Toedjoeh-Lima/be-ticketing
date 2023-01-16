package database

import (
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
	}

	return Migration{
		Migrations: migration,
	}
}
