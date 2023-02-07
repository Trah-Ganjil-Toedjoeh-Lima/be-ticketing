package repository

import (
	"github.com/frchandra/ticketing-gmcgo/app/model"
	"github.com/frchandra/ticketing-gmcgo/app/util"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	db  *gorm.DB
	log *util.LogUtil
}

func NewTransactionRepository(db *gorm.DB, log *util.LogUtil) *TransactionRepository {
	return &TransactionRepository{db: db, log: log}
}

func (t *TransactionRepository) GetBySeat(transaction *model.Transaction, seatId uint) *gorm.DB {
	result := t.db.Where("seat_id = ?", seatId).Find(transaction)
	if result != nil {
		t.log.BasicLog(result.Error, "TransactionRepotisoty@GetBySeat")
	}
	return result
}

func (t *TransactionRepository) GetBySeatUser(transaction *model.Transaction, seatId uint, userId uint64) *gorm.DB {
	result := t.db.Where("seat_id = ? AND user_id = ?", seatId, userId).Find(transaction)
	if result != nil {
		t.log.BasicLog(result.Error, "TransactionRepotisoty@GetBySeatUser")
	}
	return result
}

func (t *TransactionRepository) GetByUser(transactions *[]model.Transaction, userId uint64) *gorm.DB {
	result := t.db.Where("user_id = ?", userId).Find(transactions)
	if result != nil {
		t.log.BasicLog(result.Error, "TransactionRepotisoty@GetByUser")
	}
	return result
}

func (t *TransactionRepository) GetDetailsByUser(transactions *[]model.Transaction, userId uint64) *gorm.DB {
	result := t.db.Joins("User").Joins("Seat").Where("transactions.user_id = ?", userId).Find(transactions)
	if result != nil {
		t.log.BasicLog(result.Error, "TransactionRepotisoty@GetDetailsByUser")
	}
	return result
}

func (t *TransactionRepository) GetByOrder(transactions *[]model.Transaction, orderId string) *gorm.DB {
	result := t.db.Where("order_id = ?", orderId).Find(transactions)
	if result != nil {
		t.log.BasicLog(result.Error, "TransactionRepotisoty@GetByOrder")
	}
	return result
}

func (t *TransactionRepository) GetDetailsByOrder(transactions *[]model.Transaction, orderId string) *gorm.DB {
	result := t.db.Joins("User").Joins("Seat").Where("transactions.order_id = ?", orderId).Find(transactions)
	if result != nil {
		t.log.BasicLog(result.Error, "TransactionRepotisoty@GetDetailsByOrder")
	}
	return result
}

func (t *TransactionRepository) UpdatePaymentStatus(orderId, vendor, confirmation string) *gorm.DB {
	result := t.db.Model(&model.Transaction{}).Where("order_id = ?", orderId).Updates(model.Transaction{Vendor: vendor, Confirmation: confirmation})
	if result != nil {
		t.log.BasicLog(result.Error, "TransactionRepotisoty@UpdatePaymentStatus")
	}
	return result
}

func (t *TransactionRepository) UpdateUserPaymentStatus(userId uint64, orderId, confirmation string) *gorm.DB {
	result := t.db.Model(&model.Transaction{}).Where("user_id = ? AND order_id = ?", userId, orderId).Update("confirmation", confirmation)
	if result != nil {
		t.log.BasicLog(result.Error, "TransactionRepotisoty@UpdateUserPaymentStatus")
	}
	return result
}

func (t *TransactionRepository) UpdateUserOrderId(userId uint64, orderId string) *gorm.DB {
	result := t.db.Model(&model.Transaction{}).Where("user_id = ? AND order_id = ?", userId, "").Update("order_id", orderId)
	if result != nil {
		t.log.BasicLog(result.Error, "TransactionRepotisoty@UpdateUserOrderId")
	}
	return result
}

func (t *TransactionRepository) InsertOne(tx *model.Transaction) *gorm.DB {
	result := t.db.Create(tx)
	if result != nil {
		t.log.BasicLog(result.Error, "TransactionRepotisoty@InsertOne")
	}
	return result
}

func (t *TransactionRepository) SoftDeleteBySeatUser(seatId uint, userId uint64) *gorm.DB {
	result := t.db.Where("seat_id = ? AND user_id = ?", seatId, userId).Delete(&model.Transaction{})
	if result != nil {
		t.log.BasicLog(result.Error, "TransactionRepotisoty@SoftDeleteBySeatUser")
	}
	return result
}

func (t *TransactionRepository) SoftDeleteByOrder(orderId string) *gorm.DB {
	result := t.db.Where("order_id = ?", orderId).Delete(&model.Transaction{})
	if result != nil {
		t.log.BasicLog(result.Error, "TransactionRepotisoty@SoftDeleteByOrder")
	}
	return result
}
