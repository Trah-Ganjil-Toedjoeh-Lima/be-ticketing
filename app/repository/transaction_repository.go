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

func (t *TransactionRepository) GetBasicsByLink(transaction *model.Transaction, link string) *gorm.DB {
	result := t.db.Select("users.name", "users.email", "users.phone", "seats.name").Joins("User").Joins("Seat", t.db.Where(&model.Seat{Link: link})).First(transaction)
	return result
}

func (t *TransactionRepository) GetByUser(transactions *[]model.Transaction, userId uint64) *gorm.DB {
	result := t.db.Where("transactions.user_id = ?", userId).Find(transactions)
	return result
}

func (t *TransactionRepository) GetDetailsByUserConfirmation(transactions *[]model.Transaction, userId uint64, confirmation string) *gorm.DB {
	result := t.db.Joins("User").Joins("Seat").Where("transactions.user_id = ?", userId).Where("confirmation = ?", confirmation).Find(transactions)
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
