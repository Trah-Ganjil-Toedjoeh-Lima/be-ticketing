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

func (t TransactionRepository) InsertOne(tx *model.Transaction) *gorm.DB {
	result := t.db.Create(tx)
	return result
}

func (t *TransactionRepository) FindOne(transaction *model.Transaction) *gorm.DB {
	var txOut model.Transaction
	return t.db.Where("seat_id = ? AND user_id = ?", transaction.SeatId, transaction.UserId).Last(&txOut)
}

func (t *TransactionRepository) Upsert(transaction *model.Transaction) *gorm.DB {
	if result := t.FindOne(transaction); result.RowsAffected < 1 {
		return t.InsertOne(transaction)
	} else {
		return result.Updates(transaction)
	}
}
