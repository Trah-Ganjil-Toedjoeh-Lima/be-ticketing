package service

import (
	"github.com/frchandra/gmcgo/app/util"
)

type SnapService struct {
	txService   *TransactionService
	seatService *SeatService
	snapUtil    *util.SnapUtil
	userService *UserService
	emailUtil   *util.EmailUtil
}

func NewSnapService(txService *TransactionService, seatService *SeatService, snapUtil *util.SnapUtil, userService *UserService, emailUtil *util.EmailUtil) *SnapService {
	return &SnapService{txService: txService, seatService: seatService, snapUtil: snapUtil, userService: userService, emailUtil: emailUtil}
}

func (s *SnapService) HandleSettlement(message map[string]any) error {
	//TODO: create ticket, send email
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

func (s *SnapService) PrepareTxDetailsByMsg(message map[string]any) ([]string, string, string) {
	var seats []string
	var userName string
	var userEmail string
	transactions, _ := s.txService.GetTxDetailsByOrder(message["order_id"].(string))
	for _, tx := range transactions {
		seats = append(seats, tx.Seat.Name)
	}
	userName = transactions[0].User.Name
	userEmail = transactions[0].User.Email
	return seats, userName, userEmail
}

func (s *SnapService) SendInfoEmail(seatsName []string, receiverName, receiverEmail string) error {
	data := map[string]any{
		"Name":  receiverName,
		"Seats": seatsName,
	}
	if err := s.emailUtil.SendEmail("./resource/template/info.gohtml", data, receiverEmail, "INFO EMAIL", []string{""}); err != nil {
		return err
	}
	return nil
}

func (s *SnapService) SendTicketEmail(seatsName []string, receiverName, receiverEmail string) error {
	data := map[string]any{
		"Name":  receiverName,
		"Seats": seatsName,
	}
	if err := s.emailUtil.SendEmail("./resource/template/ticket.gohtml", data, receiverEmail, "TIKCET EMAIL", []string{""}); err != nil {
		return err
	}
	return nil
}
