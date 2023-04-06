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
	/*	temp := []string{"A", "B", "C", "D", "E"}
		for i := 1; i <= 5; i++ {
			for j := 10; j <= 21; j++ {
				name := temp[i-1] + strconv.Itoa(j)
				seat := &model.Seat{
					Name:   name,
					Price:  165000,
					Link:   uuid.New().String(),
					Status: "available",
				}
				err := this.db.Debug().Create(seat).Error
				if err != nil {
					return err
				}
			}
		}

		temp = []string{"H", "I", "J", "K", "L"}
		for i := 1; i <= 5; i++ {
			for j := 10; j <= 21; j++ {
				name := temp[i-1] + strconv.Itoa(j)
				seat := &model.Seat{
					Name:   name,
					Price:  165000,
					Link:   uuid.New().String(),
					Status: "available",
				}
				err := this.db.Debug().Create(seat).Error
				if err != nil {
					return err
				}
			}
		}

		temp = []string{"E", "F", "G", "H", "I", "J", "K", "L", "M"}
		for i := 1; i <= 9; i++ {
			for j := 1; j <= 9; j++ {
				name := temp[i-1] + strconv.Itoa(j)
				seat := &model.Seat{
					Name:   name,
					Price:  145000,
					Link:   uuid.New().String(),
					Status: "available",
				}
				err := this.db.Debug().Create(seat).Error
				if err != nil {
					return err
				}
			}
		}

		temp = []string{"E", "F", "G", "H", "I", "J", "K", "L", "M"}
		for i := 1; i <= 9; i++ {
			for j := 22; j <= 30; j++ {
				name := temp[i-1] + strconv.Itoa(j)
				seat := &model.Seat{
					Name:   name,
					Price:  145000,
					Link:   uuid.New().String(),
					Status: "available",
				}
				err := this.db.Debug().Create(seat).Error
				if err != nil {
					return err
				}
			}
		}

		for i := 10; i <= 15; i++ {
			seat := &model.Seat{
				Name:   "M" + strconv.Itoa(i),
				Price:  165000,
				Link:   uuid.New().String(),
				Status: "available",
			}
			err := this.db.Debug().Create(seat).Error
			if err != nil {
				return err
			}
		}

		temp = []string{"B", "C", "D"}
		for i := 1; i <= 3; i++ {
			for j := 3; j <= 9; j++ {
				name := temp[i-1] + strconv.Itoa(j)
				seat := &model.Seat{
					Name:   name,
					Price:  145000,
					Link:   uuid.New().String(),
					Status: "available",
				}
				err := this.db.Debug().Create(seat).Error
				if err != nil {
					return err
				}
			}
		}

		temp = []string{"B", "C", "D"}
		for i := 1; i <= 3; i++ {
			for j := 22; j <= 28; j++ {
				name := temp[i-1] + strconv.Itoa(j)
				seat := &model.Seat{
					Name:   name,
					Price:  145000,
					Link:   uuid.New().String(),
					Status: "available",
				}
				err := this.db.Debug().Create(seat).Error
				if err != nil {
					return err
				}
			}
		}

		for i := 1; i <= 15; i++ {
			seat := &model.Seat{
				Name:   "O" + strconv.Itoa(i),
				Price:  145000,
				Link:   uuid.New().String(),
				Status: "available",
			}
			err := this.db.Debug().Create(seat).Error
			if err != nil {
				return err
			}
		}

		for i := 23; i <= 30; i++ {
			seat := &model.Seat{
				Name:   "O" + strconv.Itoa(i),
				Price:  145000,
				Link:   uuid.New().String(),
				Status: "available",
			}
			err := this.db.Debug().Create(seat).Error
			if err != nil {
				return err
			}
		}
	*/
	seat := &model.Seat{
		Name:   "A8",
		Price:  145000,
		Link:   uuid.New().String(),
		Status: "available",
	}
	err := this.db.Debug().Create(seat).Error
	if err != nil {
		return err
	}

	seat = &model.Seat{
		Name:   "A9",
		Price:  145000,
		Link:   uuid.New().String(),
		Status: "available",
	}
	err = this.db.Debug().Create(seat).Error
	if err != nil {
		return err
	}

	seat = &model.Seat{
		Name:   "A22",
		Price:  145000,
		Link:   uuid.New().String(),
		Status: "available",
	}
	err = this.db.Debug().Create(seat).Error
	if err != nil {
		return err
	}

	seat = &model.Seat{
		Name:   "A23",
		Price:  145000,
		Link:   uuid.New().String(),
		Status: "available",
	}
	err = this.db.Debug().Create(seat).Error
	if err != nil {
		return err
	}

	seat = &model.Seat{
		Name:   "C1",
		Price:  165000,
		Link:   uuid.New().String(),
		Status: "available",
	}
	err = this.db.Debug().Create(seat).Error
	if err != nil {
		return err
	}

	seat = &model.Seat{
		Name:   "C2",
		Price:  165000,
		Link:   uuid.New().String(),
		Status: "available",
	}
	err = this.db.Debug().Create(seat).Error
	if err != nil {
		return err
	}

	seat = &model.Seat{
		Name:   "D1",
		Price:  165000,
		Link:   uuid.New().String(),
		Status: "available",
	}
	err = this.db.Debug().Create(seat).Error
	if err != nil {
		return err
	}

	seat = &model.Seat{
		Name:   "D2",
		Price:  165000,
		Link:   uuid.New().String(),
		Status: "available",
	}
	err = this.db.Debug().Create(seat).Error
	if err != nil {
		return err
	}

	/*	seat = &model.Seat{
			Name:   "C29",
			Price:  165000,
			Link:   uuid.New().String(),
			Status: "available",
		}
		err = this.db.Debug().Create(seat).Error
		if err != nil {
			return err
		}

		seat = &model.Seat{
			Name:   "C30",
			Price:  165000,
			Link:   uuid.New().String(),
			Status: "available",
		}
		err = this.db.Debug().Create(seat).Error
		if err != nil {
			return err
		}

		seat = &model.Seat{
			Name:   "D29",
			Price:  165000,
			Link:   uuid.New().String(),
			Status: "available",
		}
		err = this.db.Debug().Create(seat).Error
		if err != nil {
			return err
		}

		seat = &model.Seat{
			Name:   "D30",
			Price:  165000,
			Link:   uuid.New().String(),
			Status: "available",
		}
		err = this.db.Debug().Create(seat).Error
		if err != nil {
			return err
		}*/

	/*	temp = []string{"F", "G"}
		for i := 1; i <= 2; i++ {
			for j := 10; j <= 21; j++ {
				name := temp[i-1] + strconv.Itoa(j)
				seat = &model.Seat{
					Name:   name,
					Price:  165000,
					Link:   uuid.New().String(),
					Status: "available",
				}
				err = this.db.Debug().Create(seat).Error
				if err != nil {
					return err
				}
			}
		}*/

	return nil

}
