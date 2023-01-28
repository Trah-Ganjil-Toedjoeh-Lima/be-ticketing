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

func (s *SnapService) HandleSettlement(message map[string]any) error {
	//TODO: create qr, send email
	transactions, _ := s.txService.GetTxByOrderId(message["order_id"].(string))
	//update seats availability
	for _, tx := range transactions {
		if err := s.seatService.UpdateStatus(tx.SeatId, "sold"); err != nil {
			return err
		}
	}
	//update tx status
	s.txService.UpdatePaymentStatus(message["order_id"].(string), message["payment_type"].(string), message["transaction_status"].(string))
	return nil
}

func (s *SnapService) HandleFailure(message map[string]any) error {
	transactions, _ := s.txService.GetTxByOrderId(message["order_id"].(string))
	for _, tx := range transactions {
		if err := s.seatService.UpdateStatus(tx.SeatId, "#"); err != nil {
			return err
		}
	}
	return nil
}

func (s *SnapService) HandlePending(message map[string]any) error {
	//update tx status
	s.txService.UpdatePaymentStatus(message["order_id"].(string), message["payment_type"].(string), message["transaction_status"].(string))
	return nil
}
