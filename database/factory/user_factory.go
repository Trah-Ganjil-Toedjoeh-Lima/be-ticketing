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

func (this *UserFactory) RunFactory() error {
	adminUser := model.User{
		Name:  this.config.AdminName,
		Email: this.config.AdminEmail,
		Phone: this.config.AdminPhone,
	}

	err := this.db.Debug().Create(&adminUser).Error
	if err != nil {
		return err
	}
	return nil
}
