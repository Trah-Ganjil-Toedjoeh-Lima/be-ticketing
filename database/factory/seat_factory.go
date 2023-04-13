package factory

import (
	"github.com/frchandra/ticketing-gmcgo/app/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SeatFactory struct {
	db *gorm.DB
}

func NewSeatFactory(db *gorm.DB) SeatFactory {
	return SeatFactory{db: db}
}

func (this SeatFactory) RunFactory() error {

	seat := &model.Seat{
		Name:     "A8",
		Price:    145000,
		Category: "Diamond",
		Link:     uuid.New().String(),
		Status:   "available",
		Row:      "A",
		Column:   8,
	}
	err := this.db.Debug().Create(seat).Error
	if err != nil {
		return err
	}

	seat = &model.Seat{
		Name:     "G40",
		Price:    145000,
		Category: "Diamond",
		Link:     uuid.New().String(),
		Status:   "available",
		Row:      "G",
		Column:   40,
	}
	err = this.db.Debug().Create(seat).Error
	if err != nil {
		return err
	}

	seat = &model.Seat{
		Name:     "A9",
		Price:    145000,
		Category: "Diamond",
		Link:     uuid.New().String(),
		Status:   "available",
		Row:      "A",
		Column:   9,
	}
	err = this.db.Debug().Create(seat).Error
	if err != nil {
		return err
	}

	seat = &model.Seat{
		Name:     "A22",
		Price:    145000,
		Category: "Gold",
		Link:     uuid.New().String(),
		Status:   "available",
		Row:      "A",
		Column:   22,
	}
	err = this.db.Debug().Create(seat).Error
	if err != nil {
		return err
	}

	seat = &model.Seat{
		Name:     "A23",
		Price:    145000,
		Category: "Gold",
		Link:     uuid.New().String(),
		Status:   "available",
		Row:      "A",
		Column:   23,
	}
	err = this.db.Debug().Create(seat).Error
	if err != nil {
		return err
	}

	seat = &model.Seat{
		Name:     "C1",
		Price:    165000,
		Category: "Iron",
		Link:     uuid.New().String(),
		Status:   "available",
		Row:      "C",
		Column:   1,
	}
	err = this.db.Debug().Create(seat).Error
	if err != nil {
		return err
	}

	seat = &model.Seat{
		Name:     "C2",
		Price:    165000,
		Category: "Iron",
		Link:     uuid.New().String(),
		Status:   "available",
		Row:      "C",
		Column:   2,
	}
	err = this.db.Debug().Create(seat).Error
	if err != nil {
		return err
	}

	seat = &model.Seat{
		Name:     "D1",
		Price:    165000,
		Category: "Wood",
		Link:     uuid.New().String(),
		Status:   "available",
		Row:      "D",
		Column:   1,
	}
	err = this.db.Debug().Create(seat).Error
	if err != nil {
		return err
	}

	seat = &model.Seat{
		Name:     "D2",
		Price:    165000,
		Category: "Wood",
		Link:     uuid.New().String(),
		Status:   "available",
		Row:      "D",
		Column:   2,
	}
	err = this.db.Debug().Create(seat).Error
	if err != nil {
		return err
	}

	return nil
}
