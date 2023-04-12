package repository

import (
	"github.com/frchandra/ticketing-gmcgo/app/model"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (t *TransactionRepository) GetAllWithDetails(transactions *[]model.Transaction) *gorm.DB {
	result := t.db.Joins("User").Joins("Seat").Find(transactions)
	return result
}

func (t *TransactionRepository) GetBySeatTxn(txn *gorm.DB, transaction *model.Transaction, seatId uint) *gorm.DB {
	result := txn.Where("seat_id = ?", seatId).Find(transaction)
	return result
}

func (t *TransactionRepository) GetDetailsByLink(transaction *model.Transaction, link string) *gorm.DB {
	result := t.db.Joins("User").Joins("Seat", t.db.Where(&model.Seat{Link: link})).First(transaction)
	return result
}

type basicTransaction struct {
	UserName string
	Email    string
	Phone    string
	Link     string
	SeatName string
	Price    uint
}

func (t *TransactionRepository) GetBasicsByLink(transaction *model.Transaction, link string) *gorm.DB {
	var basic basicTransaction
	result := t.db.Model(transaction).Select("users.name AS user_name", "users.email", "users.phone", "seats.link", "seats.name AS seat_name", "seats.price").
		Joins("inner join users on users.user_id = transactions.user_id").
		Joins("inner join seats on seats.seat_id = transactions.seat_id").
		Where("seats.link = ?", link).
		Order("transaction_id").
		Limit(1).
		Scan(&basic)
	transaction.User.Name = basic.UserName
	transaction.User.Email = basic.Email
	transaction.User.Phone = basic.Phone
	transaction.Seat.Link = basic.Link
	transaction.Seat.Name = basic.SeatName
	return result
}

func (t *TransactionRepository) GetByUser(transactions *[]model.Transaction, userId uint64) *gorm.DB {
	result := t.db.Where("transactions.user_id = ?", userId).Find(transactions)
	return result
}

func (t *TransactionRepository) GetDetailsByUserConfirmation(transactions *[]model.Transaction, userId uint64, confirmation string) *gorm.DB {
	var basics []basicTransaction
	result := t.db.Table("transactions").Select("users.name AS user_name", "users.email", "users.phone", "seats.link", "seats.name AS seat_name", "seats.price").
		Joins("inner join users on users.user_id = transactions.user_id").
		Joins("inner join seats on seats.seat_id = transactions.seat_id").
		Where("transactions.user_id = ?", userId).
		Where("transactions.confirmation = ?", confirmation).
		Order("transaction_id").
		Scan(&basics)
	// for index, basic := range basics {

	// }

	//result := t.db.Joins("User").Joins("Seat").Where("transactions.user_id = ?", userId).Where("confirmation = ?", confirmation).Find(transactions)
	return result
}

func (t *TransactionRepository) GetByOrder(transactions *[]model.Transaction, orderId string) *gorm.DB {
	result := t.db.Where("order_id = ?", orderId).Find(transactions)
	return result
}

func (t *TransactionRepository) GetDetailsByOrder(transactions *[]model.Transaction, orderId string) *gorm.DB {
	result := t.db.Joins("User").Joins("Seat").Where("transactions.order_id = ?", orderId).Find(transactions)
	return result
}

func (t *TransactionRepository) UpdatePaymentStatus(orderId, vendor, confirmation string) *gorm.DB {
	result := t.db.Model(&model.Transaction{}).Where("order_id = ?", orderId).Updates(model.Transaction{Vendor: vendor, Confirmation: confirmation})
	return result
}

func (t *TransactionRepository) UpdatePaymentStatusByUser(userId uint64, confirmation string) *gorm.DB {
	result := t.db.Model(&model.Transaction{}).Where("user_id = ?", userId).Update("confirmation", confirmation)
	return result
}

func (t *TransactionRepository) UpdateUserOrderId(userId uint64, orderId string) *gorm.DB {
	result := t.db.Model(&model.Transaction{}).Where("user_id = ? AND order_id = ?", userId, "").Update("order_id", orderId)
	return result
}

func (t *TransactionRepository) InsertOne(tx *model.Transaction) *gorm.DB {
	result := t.db.Create(tx)
	return result
}

func (t *TransactionRepository) SoftDeleteBySeatUser(seatId uint, userId uint64) *gorm.DB {
	result := t.db.Where("seat_id = ? AND user_id = ?", seatId, userId).Delete(&model.Transaction{})
	return result
}

func (t *TransactionRepository) SoftDeleteByOrder(orderId string) *gorm.DB {
	result := t.db.Where("order_id = ?", orderId).Delete(&model.Transaction{})
	return result
}
