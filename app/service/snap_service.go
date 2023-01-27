package service

import (
	"github.com/frchandra/gmcgo/app/repository"
	"github.com/frchandra/gmcgo/app/util"
)

type SnapService struct {
	txRepo   *repository.TransactionRepository
	seatRepo *repository.SeatRepository
	snapUtil *util.SnapUtil
}

func NewSnapService(txRepo *repository.TransactionRepository, seatRepo *repository.SeatRepository, snapUtil *util.SnapUtil) *SnapService {
	return &SnapService{txRepo: txRepo, seatRepo: seatRepo, snapUtil: snapUtil}
}
