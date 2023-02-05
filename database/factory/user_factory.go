package factory

import (
	"github.com/bxcodec/faker/v4"
	"github.com/frchandra/ticketing-gmcgo/app/model"
	"gorm.io/gorm"
	"time"
)

type UserFactory struct {
	db *gorm.DB
}

func NewUserFactory(db *gorm.DB) UserFactory {
	return UserFactory{db: db}
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

func (this UserFactory) RunFactory() error {
	count := 3
	for i := 0; i < count; i++ {
		err := this.db.Debug().Create(this.GetData()).Error
		if err != nil {
			return err
		}
	}
	return nil
}
