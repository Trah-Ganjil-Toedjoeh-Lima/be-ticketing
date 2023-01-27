package service

import (
	"github.com/frchandra/gmcgo/app/util"
)

type SnapService struct {
	txService   *TransactionService
	seatService *SeatService
	snapUtil    *util.SnapUtil
}

func NewSnapService(txService *TransactionService, seatService *SeatService, snapUtil *util.SnapUtil) *SnapService {
	return &SnapService{txService: txService, seatService: seatService, snapUtil: snapUtil}
}

func (s *SnapService) HandleSettlement(orderId string) error {
	//TODO: create qr, send email
	transactions, _ := s.txService.GetTxByOrderId(orderId)
	for _, tx := range transactions {
		if err := s.seatService.UpdateStatus(tx.SeatId, "sold"); err != nil {
			return err
		}
	}
	return nil
}

func (s *SnapService) HandleFailure(orderId string) error {
	transactions, _ := s.txService.GetTxByOrderId(orderId)
	for _, tx := range transactions {
		if err := s.seatService.UpdateStatus(tx.SeatId, "#"); err != nil {
			return err
		}
	}
	return nil
}
