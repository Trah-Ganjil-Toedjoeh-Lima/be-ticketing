package repository

import (
	"github.com/frchandra/gmcgo/app/model"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (t *TransactionRepository) GetLastTxBySeatIdUserId(transaction *model.Transaction, seatId uint, userId uint64) *gorm.DB {
	return t.db.Where("seat_id = ? AND user_id = ?", seatId, userId).Find(transaction)
}

func (t *TransactionRepository) GetLastTxByUserId(transactions *[]model.Transaction, userId uint64) *gorm.DB {
	return t.db.Where("user_id = ?", userId).Find(transactions)
}

func (t *TransactionRepository) InsertOne(tx *model.Transaction) *gorm.DB {
	result := t.db.Create(tx)
	return result
}

func (t *TransactionRepository) SoftDeleteTransaction(seatId uint, userId uint64) *gorm.DB {
	return t.db.Where("seat_id = ? AND user_id = ?", seatId, userId).Delete(&model.Transaction{})
}

func (t *TransactionRepository) GetUserTransactionDetails(transactions *[]model.Transaction, userId uint64) *gorm.DB {
	return t.db.Joins("User").Joins("Seat").Where("transactions.user_id = ?", userId).Find(transactions)
}

func (t *TransactionRepository) GetTxByOrderId(transactions *[]model.Transaction, orderId string) *gorm.DB {
	return t.db.Where("order_id = ?", orderId).Find(transactions)
}
