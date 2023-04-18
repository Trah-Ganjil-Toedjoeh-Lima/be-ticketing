package factory

import (
	"github.com/frchandra/ticketing-gmcgo/app/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type SeatFactory struct {
	db *gorm.DB
}

func NewSeatFactory(db *gorm.DB) *SeatFactory {
	return &SeatFactory{db: db}
}

func (f *SeatFactory) RunFactory() error {
	red := "red"
	pink := "pink"
	yellow := "yellow"
	blue := "blue"
	green := "green"

	var redPrice uint = 120000
	yellowPrice := 170000
	greenPrice := 145000
	pinkPrice := 85000
	bluePrice := 60000

	// ROW A
	rowName := "A"
	for j := 1; j <= 16; j++ { //column
		name := rowName + strconv.Itoa(j)
		seat := &model.Seat{
			Name:     name,
			Price:    0,
			Category: red,
			Link:     uuid.New().String(),
			Status:   "available",
			Row:      rowName,
			Column:   uint(j),
		}
		err := f.db.Debug().Create(seat).Error
		if err != nil {
			return err
		}
	}
	// ROW B
	rowName = "B"
	for j := 1; j <= 44; j++ { //column
		name := rowName + strconv.Itoa(j)
		var category string
		if j <= 13 {
			category = pink
		} else if j >= 14 && j <= 31 {
			category = red
		} else {
			category = pink
		}
		seat := &model.Seat{
			Name:     name,
			Price:    0,
			Category: category,
			Link:     uuid.New().String(),
			Status:   "available",
			Row:      rowName,
			Column:   uint(j),
		}
		err := f.db.Debug().Create(seat).Error
		if err != nil {
			return err
		}
	}
	// ROW C
	rowName = "C"
	for j := 1; j <= 46; j++ { //column
		name := rowName + strconv.Itoa(j)
		var category string
		if j <= 15 {
			category = pink
		} else if j >= 16 && j <= 33 {
			category = red
		} else {
			category = pink
		}
		seat := &model.Seat{
			Name:     name,
			Price:    0,
			Category: category,
			Link:     uuid.New().String(),
			Status:   "available",
			Row:      rowName,
			Column:   uint(j),
		}
		err := f.db.Debug().Create(seat).Error
		if err != nil {
			return err
		}
	}
	// ROW D
	rowName = "D"
	for j := 1; j <= 50; j++ { //column
		name := rowName + strconv.Itoa(j)
		var category string
		if j <= 15 {
			category = pink
		} else if j >= 16 && j <= 36 {
			category = red
		} else {
			category = pink
		}
		seat := &model.Seat{
			Name:     name,
			Price:    0,
			Category: category,
			Link:     uuid.New().String(),
			Status:   "available",
			Row:      rowName,
			Column:   uint(j),
		}
		err := f.db.Debug().Create(seat).Error
		if err != nil {
			return err
		}
	}
	// ROW E
	rowName = "E"
	for j := 1; j <= 54; j++ { //column
		name := rowName + strconv.Itoa(j)
		var category string
		if j <= 16 {
			category = pink
		} else if j >= 17 && j <= 38 {
			category = red
		} else {
			category = pink
		}
		seat := &model.Seat{
			Name:     name,
			Price:    0,
			Category: category,
			Link:     uuid.New().String(),
			Status:   "available",
			Row:      rowName,
			Column:   uint(j),
		}
		err := f.db.Debug().Create(seat).Error
		if err != nil {
			return err
		}
	}
	// ROW F
	rowName = "F"
	for j := 1; j <= 58; j++ { //column
		name := rowName + strconv.Itoa(j)
		var category string
		if j <= 17 {
			category = pink
		} else if j >= 18 && j <= 41 {
			category = red
		} else {
			category = pink
		}
		seat := &model.Seat{
			Name:     name,
			Price:    0,
			Category: category,
			Link:     uuid.New().String(),
			Status:   "available",
			Row:      rowName,
			Column:   uint(j),
		}
		err := f.db.Debug().Create(seat).Error
		if err != nil {
			return err
		}
	}
	// ROW G
	rowName = "G"
	for j := 1; j <= 58; j++ { //column
		name := rowName + strconv.Itoa(j)
		var category string
		if j <= 17 {
			category = pink
		} else if j >= 18 && j <= 41 {
			category = red
		} else {
			category = pink
		}
		seat := &model.Seat{
			Name:     name,
			Price:    0,
			Category: category,
			Link:     uuid.New().String(),
			Status:   "available",
			Row:      rowName,
			Column:   uint(j),
		}
		err := f.db.Debug().Create(seat).Error
		if err != nil {
			return err
		}
	}
	// ROW H
	rowName = "H"
	for j := 1; j <= 57; j++ { //column
		name := rowName + strconv.Itoa(j)
		var category string
		if j <= 16 {
			category = pink
		} else if j >= 17 && j <= 42 {
			category = yellow
		} else {
			category = pink
		}
		seat := &model.Seat{
			Name:     name,
			Price:    0,
			Category: category,
			Link:     uuid.New().String(),
			Status:   "available",
			Row:      rowName,
			Column:   uint(j),
		}
		err := f.db.Debug().Create(seat).Error
		if err != nil {
			return err
		}
	}
	// ROW I
	rowName = "I"
	for j := 1; j <= 56; j++ { //column
		name := rowName + strconv.Itoa(j)
		var category string
		if j <= 15 {
			category = pink
		} else if j >= 16 && j <= 43 {
			category = yellow
		} else {
			category = pink
		}
		seat := &model.Seat{
			Name:     name,
			Price:    0,
			Category: category,
			Link:     uuid.New().String(),
			Status:   "available",
			Row:      rowName,
			Column:   uint(j),
		}
		err := f.db.Debug().Create(seat).Error
		if err != nil {
			return err
		}
	}
	// ROW J
	rowName = "J"
	for j := 1; j <= 55; j++ { //column
		name := rowName + strconv.Itoa(j)
		var category string
		if j <= 14 {
			category = pink
		} else if j >= 15 && j <= 42 {
			category = yellow
		} else {
			category = pink
		}
		seat := &model.Seat{
			Name:     name,
			Price:    0,
			Category: category,
			Link:     uuid.New().String(),
			Status:   "available",
			Row:      rowName,
			Column:   uint(j),
		}
		err := f.db.Debug().Create(seat).Error
		if err != nil {
			return err
		}
	}
	// ROW K
	rowName = "K"
	for j := 1; j <= 56; j++ { //column
		name := rowName + strconv.Itoa(j)
		var category string
		if j <= 13 {
			category = blue
		} else if j >= 14 && j <= 43 {
			category = yellow
		} else {
			category = blue
		}
		seat := &model.Seat{
			Name:     name,
			Price:    0,
			Category: category,
			Link:     uuid.New().String(),
			Status:   "available",
			Row:      rowName,
			Column:   uint(j),
		}
		err := f.db.Debug().Create(seat).Error
		if err != nil {
			return err
		}
	}
	// ROW L
	rowName = "L"
	for j := 1; j <= 54; j++ { //column
		name := rowName + strconv.Itoa(j)
		var category string
		if j <= 12 {
			category = blue
		} else if j >= 13 && j <= 42 {
			category = yellow
		} else {
			category = blue
		}
		seat := &model.Seat{
			Name:     name,
			Price:    0,
			Category: category,
			Link:     uuid.New().String(),
			Status:   "available",
			Row:      rowName,
			Column:   uint(j),
		}
		err := f.db.Debug().Create(seat).Error
		if err != nil {
			return err
		}
	}
	// ROW M
	rowName = "M"
	for j := 1; j <= 54; j++ { //column
		name := rowName + strconv.Itoa(j)
		var category string
		if j <= 11 {
			category = blue
		} else if j >= 12 && j <= 43 {
			category = green
		} else {
			category = blue
		}
		seat := &model.Seat{
			Name:     name,
			Price:    0,
			Category: category,
			Link:     uuid.New().String(),
			Status:   "available",
			Row:      rowName,
			Column:   uint(j),
		}
		err := f.db.Debug().Create(seat).Error
		if err != nil {
			return err
		}
	}
	// ROW N
	rowName = "N"
	for j := 1; j <= 52; j++ { //column
		name := rowName + strconv.Itoa(j)
		var category string
		if j <= 10 {
			category = blue
		} else if j >= 11 && j <= 42 {
			category = green
		} else {
			category = blue
		}
		seat := &model.Seat{
			Name:     name,
			Price:    0,
			Category: category,
			Link:     uuid.New().String(),
			Status:   "available",
			Row:      rowName,
			Column:   uint(j),
		}
		err := f.db.Debug().Create(seat).Error
		if err != nil {
			return err
		}
	}
	// ROW O
	rowName = "O"
	for j := 1; j <= 54; j++ { //column
		name := rowName + strconv.Itoa(j)
		var category string
		if j <= 10 {
			category = blue
		} else if j >= 11 && j <= 45 {
			category = green
		} else {
			category = blue
		}
		seat := &model.Seat{
			Name:     name,
			Price:    0,
			Category: category,
			Link:     uuid.New().String(),
			Status:   "available",
			Row:      rowName,
			Column:   uint(j),
		}
		err := f.db.Debug().Create(seat).Error
		if err != nil {
			return err
		}
	}
	// ROW P
	rowName = "P"
	for j := 1; j <= 53; j++ { //column
		name := rowName + strconv.Itoa(j)
		var category string
		if j <= 9 {
			category = blue
		} else if j >= 10 && j <= 44 {
			category = green
		} else {
			category = blue
		}
		seat := &model.Seat{
			Name:     name,
			Price:    0,
			Category: category,
			Link:     uuid.New().String(),
			Status:   "available",
			Row:      rowName,
			Column:   uint(j),
		}
		err := f.db.Debug().Create(seat).Error
		if err != nil {
			return err
		}
	}
	// ROW Q
	rowName = "Q"
	for j := 1; j <= 53; j++ { //column
		name := rowName + strconv.Itoa(j)
		var category string
		if j <= 8 {
			category = blue
		} else if j >= 9 && j <= 44 {
			category = green
		} else {
			category = blue
		}
		seat := &model.Seat{
			Name:     name,
			Price:    0,
			Category: category,
			Link:     uuid.New().String(),
			Status:   "available",
			Row:      rowName,
			Column:   uint(j),
		}
		err := f.db.Debug().Create(seat).Error
		if err != nil {
			return err
		}
	}
	// ROW R
	rowName = "R"
	for j := 1; j <= 53; j++ { //column
		name := rowName + strconv.Itoa(j)
		var category string
		if j <= 7 {
			category = blue
		} else if j >= 8 && j <= 45 {
			category = green
		} else {
			category = blue
		}
		seat := &model.Seat{
			Name:     name,
			Price:    0,
			Category: category,
			Link:     uuid.New().String(),
			Status:   "available",
			Row:      rowName,
			Column:   uint(j),
		}
		err := f.db.Debug().Create(seat).Error
		if err != nil {
			return err
		}
	}
	// ROW S
	rowName = "S"
	for j := 1; j <= 52; j++ { //column
		name := rowName + strconv.Itoa(j)
		var category string
		if j <= 7 {
			category = blue
		} else if j >= 8 && j <= 45 {
			category = green
		} else {
			category = blue
		}
		seat := &model.Seat{
			Name:     name,
			Price:    0,
			Category: category,
			Link:     uuid.New().String(),
			Status:   "available",
			Row:      rowName,
			Column:   uint(j),
		}
		err := f.db.Debug().Create(seat).Error
		if err != nil {
			return err
		}
	}
	// ROW T
	rowName = "T"
	for j := 1; j <= 7; j++ { //column
		name := rowName + strconv.Itoa(j)
		seat := &model.Seat{
			Name:     name,
			Price:    0,
			Category: blue,
			Link:     uuid.New().String(),
			Status:   "available",
			Row:      rowName,
			Column:   uint(j),
		}
		err := f.db.Debug().Create(seat).Error
		if err != nil {
			return err
		}
	}
	// ROW U
	rowName = "U"
	for j := 1; j <= 63; j++ { //column
		name := rowName + strconv.Itoa(j)
		seat := &model.Seat{
			Name:     name,
			Price:    0,
			Category: pink,
			Link:     uuid.New().String(),
			Status:   "available",
			Row:      rowName,
			Column:   uint(j),
		}
		err := f.db.Debug().Create(seat).Error
		if err != nil {
			return err
		}
	}
	// ROW V
	rowName = "V"
	for j := 1; j <= 62; j++ { //column
		name := rowName + strconv.Itoa(j)
		seat := &model.Seat{
			Name:     name,
			Price:    0,
			Category: pink,
			Link:     uuid.New().String(),
			Status:   "available",
			Row:      rowName,
			Column:   uint(j),
		}
		err := f.db.Debug().Create(seat).Error
		if err != nil {
			return err
		}
	}

	//
	// SET PRICE
	//
	// RED
	err := f.db.Model(model.Seat{}).Where("category = ?", red).Updates(model.Seat{Price: redPrice, CreatedAt: time.Now()}).Error
	if err != nil {
		return err
	}
	// YELLOW
	err = f.db.Model(model.Seat{}).Where("category = ?", yellow).Update("price", yellowPrice).Error
	if err != nil {
		return err
	}
	// GREEN
	err = f.db.Model(model.Seat{}).Where("category = ?", green).Update("price", greenPrice).Error
	if err != nil {
		return err
	}
	// PINK
	err = f.db.Model(model.Seat{}).Where("category = ?", pink).Update("price", pinkPrice).Error
	if err != nil {
		return err
	}
	// BLUE
	err = f.db.Model(model.Seat{}).Where("category = ?", blue).Update("price", bluePrice).Error
	if err != nil {
		return err
	}

	return nil
}
