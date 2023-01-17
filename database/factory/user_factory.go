package factory

import (
	"github.com/bxcodec/faker/v4"
	"github.com/frchandra/gmcgo/app/model"
	"gorm.io/gorm"
	"time"
)

type UserFactory struct {
	Database *gorm.DB
}

func NewUserFactory(db *gorm.DB) UserFactory {
	return UserFactory{Database: db}
}

func (this UserFactory) GetData() interface{} {
	return &model.User{
		Name:      faker.Name(),
		Email:     faker.Email(),
		Phone:     faker.E164PhoneNumber(),
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		DeletedAt: gorm.DeletedAt{},
	}
}

func (this UserFactory) RunFactory(count int) error {
	for i := 0; i < count; i++ {
		err := this.Database.Debug().Create(this.GetData()).Error
		if err != nil {
			return err
		}
	}
	return nil
}
