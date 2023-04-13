package factory

import (
	"github.com/frchandra/ticketing-gmcgo/app/model"
	"github.com/frchandra/ticketing-gmcgo/config"
	"gorm.io/gorm"
)

type UserFactory struct {
	db     *gorm.DB
	config *config.AppConfig
}

func NewUserFactory(db *gorm.DB, config *config.AppConfig) *UserFactory {
	return &UserFactory{db: db, config: config}
}

func (f *UserFactory) RunFactory() error {
	adminUser := model.User{
		Name:  f.config.AdminName,
		Email: f.config.AdminEmail,
		Phone: f.config.AdminPhone,
	}

	err := f.db.Debug().Create(&adminUser).Error
	if err != nil {
		return err
	}
	return nil
}
