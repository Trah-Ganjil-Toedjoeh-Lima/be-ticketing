package service

import (
	"fmt"
	"github.com/frchandra/ticketing-gmcgo/app/model"
	"github.com/frchandra/ticketing-gmcgo/app/repository"
	"github.com/frchandra/ticketing-gmcgo/app/util"
)

type SnapService struct {
	txService   *TransactionService
	seatService *SeatService
	txRepo      *repository.TransactionRepository
	snapUtil    *util.SnapUtil
	emailUtil   *util.EmailUtil
	eticketUtil *util.ETicketUtil
	log         *util.LogUtil
}

func NewSnapService(txService *TransactionService, seatService *SeatService, txRepo *repository.TransactionRepository, snapUtil *util.SnapUtil, emailUtil *util.EmailUtil, eticketUtil *util.ETicketUtil, log *util.LogUtil) *SnapService {
	return &SnapService{txService: txService, seatService: seatService, txRepo: txRepo, snapUtil: snapUtil, emailUtil: emailUtil, eticketUtil: eticketUtil, log: log}
}

func (s *SnapService) HandleSettlement(message map[string]any) error {
	transactions, _ := s.txService.GetByOrder(message["order_id"].(string))

	for _, tx := range transactions { //update seats availability
		if err := s.seatService.UpdateStatus(tx.SeatId, "purchased"); err != nil {
			return err
		}
	}

	if err := s.txService.UpdatePaymentStatus(message["order_id"].(string), message["payment_type"].(string), message["transaction_status"].(string)); err != nil { //update tx status
		return err
	}
	return nil
}

func (s *SnapService) HandleFailure(message map[string]any) error {
	transactions, _ := s.txService.GetByOrder(message["order_id"].(string))
	for _, tx := range transactions {
		if err := s.seatService.UpdateStatus(tx.SeatId, "available"); err != nil {
			return err
		}
	}

	if err := s.txService.UpdatePaymentStatus(message["order_id"].(string), message["payment_type"].(string), message["transaction_status"].(string)); err != nil { //update tx status
		return err
	}

	s.txRepo.SoftDeleteByOrder(message["order_id"].(string)) //soft delete tx status
	return nil

}

func (s *SnapService) HandlePending(message map[string]any) error {

	if err := s.txService.UpdatePaymentStatus(message["order_id"].(string), message["payment_type"].(string), message["transaction_status"].(string)); err != nil { //update tx status
		return err
	}
	return nil
}

func (s *SnapService) PrepareTxDetailsByMsg(message map[string]any) ([]model.Seat, string, string) {
	var seats []model.Seat
	var userName string
	var userEmail string
	transactions, _ := s.txService.GetDetailsByOrder(message["order_id"].(string))
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
	s.log.Log.Info("successfully sent pending email to " + receiverEmail)
	return nil
}

func (s *SnapService) SendTicketEmail(seats []model.Seat, receiverName, receiverEmail string) error {
	var attachmentPath []string
	var seatsName []string
	for _, seat := range seats {
		if err := s.eticketUtil.GenerateETicket(seat.Name, seat.Link); err != nil {
			return err
		}
		attachmentPath = append(attachmentPath, "./storage/ticket/"+seat.Name+".png")
		seatsName = append(seatsName, seat.Name)
	}

	data := map[string]any{
		"Name":  receiverName,
		"Seats": seatsName,
	}
	fmt.Println(data)
	fmt.Println("SENDING OK EMAIL")
	if err := s.emailUtil.SendEmail("./resource/template/ticket.gohtml", data, receiverEmail, "TICKET EMAIL", attachmentPath); err != nil {
		return err
	}
	s.log.Log.Info("successfully sent ticket email to " + receiverEmail)
	return nil
}
