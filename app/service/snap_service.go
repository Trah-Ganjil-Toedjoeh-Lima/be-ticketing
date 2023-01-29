package service

import (
	"fmt"
	"github.com/frchandra/gmcgo/app/model"
	"github.com/frchandra/gmcgo/app/util"
)

type SnapService struct {
	txService   *TransactionService
	seatService *SeatService
	snapUtil    *util.SnapUtil
	userService *UserService
	emailUtil   *util.EmailUtil
	eticketUtil *util.ETicketUtil
}

func NewSnapService(txService *TransactionService, seatService *SeatService, snapUtil *util.SnapUtil, userService *UserService, emailUtil *util.EmailUtil, eticketUtil *util.ETicketUtil) *SnapService {
	return &SnapService{txService: txService, seatService: seatService, snapUtil: snapUtil, userService: userService, emailUtil: emailUtil, eticketUtil: eticketUtil}
}

func (s *SnapService) HandleSettlement(message map[string]any) error {
	transactions, _ := s.txService.GetTxByOrder(message["order_id"].(string))
	//update seats availability
	for _, tx := range transactions {
		if err := s.seatService.UpdateStatus(tx.SeatId, "sold"); err != nil {
			return err
		}
	}
	//update tx status
	if err := s.txService.UpdatePaymentStatus(message["order_id"].(string), message["payment_type"].(string), message["transaction_status"].(string)); err != nil {
		return err
	}

	return nil
}

func (s *SnapService) HandleFailure(message map[string]any) error {
	transactions, _ := s.txService.GetTxByOrder(message["order_id"].(string))
	for _, tx := range transactions {
		if err := s.seatService.UpdateStatus(tx.SeatId, "#"); err != nil {
			return err
		}
	}
	return nil
	//TODO: soft delete
}

func (s *SnapService) HandlePending(message map[string]any) error {
	//update tx status
	if err := s.txService.UpdatePaymentStatus(message["order_id"].(string), message["payment_type"].(string), message["transaction_status"].(string)); err != nil {
		return err
	}
	return nil
}

func (s *SnapService) PrepareTxDetailsByMsg(message map[string]any) ([]model.Seat, string, string) {
	var seats []model.Seat
	var userName string
	var userEmail string
	transactions, _ := s.txService.GetTxDetailsByOrder(message["order_id"].(string))
	for _, tx := range transactions {
		seats = append(seats, tx.Seat)
	}
	userName = transactions[0].User.Name
	userEmail = transactions[0].User.Email
	return seats, userName, userEmail
}

func (s *SnapService) SendInfoEmail(seats []model.Seat, receiverName, receiverEmail string) error {
	var seatsName []string
	for _, seat := range seats {
		seatsName = append(seatsName, seat.Name)
	}
	data := map[string]any{
		"Name":  receiverName,
		"Seats": seatsName,
	}
	fmt.Println("SENDING PENDING EMAIL")
	if err := s.emailUtil.SendEmail("./resource/template/info.gohtml", data, receiverEmail, "INFO EMAIL", []string{}); err != nil {
		return err
	}
	return nil
}

func (s *SnapService) SendTicketEmail(seats []model.Seat, receiverName, receiverEmail string) error {
	var attachementPath []string
	var seatsName []string
	for _, seat := range seats {
		if err := s.eticketUtil.GenerateETicket(seat.Name, seat.Link); err != nil {
			return err
		}
		attachementPath = append(attachementPath, "./storage/ticket/"+seat.Name+".png")
		seatsName = append(seatsName, seat.Name)
	}

	data := map[string]any{
		"Name":  receiverName,
		"Seats": seatsName,
	}
	fmt.Println(data)
	fmt.Println("SENDING OK EMAIL")
	if err := s.emailUtil.SendEmail("./resource/template/ticket.gohtml", data, receiverEmail, "TIKCET EMAIL", attachementPath); err != nil {
		return err
	}
	return nil
}
