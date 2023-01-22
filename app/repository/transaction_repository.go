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

func (t *TransactionRepository) GetLastTxBySeatIdUserId(transaction *model.Transaction, seatId uint64, userId uint) *gorm.DB {
	return t.db.Where("seat_id = ? AND user_id = ?", seatId, userId).Last(transaction)
}

func (t TransactionRepository) CreateTx(tx *model.Transaction) *gorm.DB {
	result := t.db.Create(tx)
	return result
}
