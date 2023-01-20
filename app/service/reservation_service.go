package service

import (
	"github.com/frchandra/gmcgo/app/model"
	"github.com/frchandra/gmcgo/app/repository"
)

type ReservationService struct {
	resRepo *repository.ReservationRepository
}

func NewReservationService(resRepo *repository.ReservationRepository) *ReservationService {
	return &ReservationService{resRepo: resRepo}
}

func (s *ReservationService) GetAllSeats() ([]model.Seat, error) {
	seats := []model.Seat{}
	if err := s.resRepo.GetAllSeats(&seats).Error; err != nil {
		return nil, err
	}
	return seats, nil
}
