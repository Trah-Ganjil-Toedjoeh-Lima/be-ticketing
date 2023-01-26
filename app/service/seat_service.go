package service

import "github.com/frchandra/gmcgo/app/repository"

type SeatService struct {
	seatRepo *repository.SeatRepository
}

func NewSeatService(seatRepo *repository.SeatRepository) *SeatService {
	return &SeatService{seatRepo: seatRepo}
}

func (s *SeatService) UpdateStatus(seatId uint, status string) error {
	if result := s.seatRepo.UpdateStatus(seatId, status); result.Error != nil {
		return result.Error
	}
	return nil
}
