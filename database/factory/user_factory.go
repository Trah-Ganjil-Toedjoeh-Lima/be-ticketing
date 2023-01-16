package factory

import (
	"github.com/bxcodec/faker/v4"
	"github.com/frchandra/gmcgo/app/model"
	"gorm.io/gorm"
	"time"
)

type UserFactory struct {
}

func NewUserFactory() *UserFactory {
	return &UserFactory{}
}

func (this *UserFactory) GetData() interface{} {
	return &model.User{
		Name:      faker.Name(),
		Email:     faker.Email(),
		Phone:     faker.E164PhoneNumber(),
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		DeletedAt: gorm.DeletedAt{},
	}
}
