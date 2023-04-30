package repository

import (
	"github.com/frchandra/ticketing-gmcgo/app/model"
	"gorm.io/gorm"
	"time"
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

type transactionFields struct {
	TransactionId uint64
	UserId        uint64
	UserName      string
	Email         string
	Phone         string
	SeatId        uint
	SeatName      string
	Price         uint
	Category      string
	Link          string
	Confirmation  string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (t *TransactionRepository) GetBasicsByLink(transaction *model.Transaction, link string) *gorm.DB {
	var basic transactionFields
	result := t.db.Model(transaction).Select(
		"transactions.transaction_id",
		"users.user_id",
		"users.name AS user_name",
		"users.email",
		"users.phone",
		"seats.seat_id",
		"seats.name AS seat_name",
		"seats.price",
		"seats.category",
		"transactions.created_at",
		"transactions.updated_at").
		Joins("inner join users on users.user_id = transactions.user_id").
		Joins("inner join seats on seats.seat_id = transactions.seat_id").
		Where("transactions.deleted_at IS NULL").
		Where("seats.link = ?", link).
		Order("transaction_id").
		Limit(1).
		Scan(&basic)

	var transactionBuff model.Transaction = model.Transaction{
		TransactionId: basic.TransactionId,
		User:          model.User{UserId: basic.UserId, Name: basic.UserName, Phone: basic.Phone, Email: basic.Email},
		Seat:          model.Seat{SeatId: basic.SeatId, Name: basic.SeatName, Price: basic.Price, Category: basic.Category},
		CreatedAt:     basic.CreatedAt,
		UpdatedAt:     basic.UpdatedAt,
	}
	*transaction = transactionBuff
	return result
}

func (t *TransactionRepository) GetDetailsByLink(transaction *model.Transaction, link string) *gorm.DB {
	result := t.db.Joins("User").Joins("Seat", t.db.Where(&model.Seat{Link: link})).First(transaction)
	return result
}

func (t *TransactionRepository) GetByUser(transactions *[]model.Transaction, userId uint64) *gorm.DB {
	result := t.db.Where("user_id = ?", userId).Find(transactions)
	return result
}

func (t *TransactionRepository) GetDetailsByUser(transactions *[]model.Transaction, userId uint64) *gorm.DB {
	var basics []transactionFields
	result := t.db.Table("transactions").Select(
		"transactions.transaction_id",
		"users.user_id",
		"users.name AS user_name",
		"users.email",
		"users.phone",
		"seats.seat_id",
		"seats.name AS seat_name",
		"seats.price",
		"seats.category",
		"transactions.confirmation",
		"transactions.created_at",
		"transactions.updated_at").
		Joins("inner join users on users.user_id = transactions.user_id").
		Joins("inner join seats on seats.seat_id = transactions.seat_id").
		Where("transactions.deleted_at IS NULL").
		Where("transactions.user_id = ?", userId).
		Order("transaction_id").
		Scan(&basics)

	var transactionsBuff []model.Transaction
	var transactionBuff model.Transaction
	for _, basic := range basics {
		transactionBuff = model.Transaction{
			TransactionId: basic.TransactionId,
			User:          model.User{UserId: basic.UserId, Name: basic.UserName, Phone: basic.Phone, Email: basic.Email},
			Seat:          model.Seat{SeatId: basic.SeatId, Name: basic.SeatName, Price: basic.Price, Category: basic.Category},
			Confirmation:  basic.Confirmation,
			CreatedAt:     basic.CreatedAt,
			UpdatedAt:     basic.UpdatedAt,
		}
		transactionsBuff = append(transactionsBuff, transactionBuff)
	}
	*transactions = transactionsBuff
	return result
}

func (t *TransactionRepository) GetDetailsByUserConfirmation(transactions *[]model.Transaction, userId uint64, confirmation []string) *gorm.DB {
	var basics []transactionFields
	result := t.db.Table("transactions").Select(
		"transactions.transaction_id",
		"users.user_id",
		"users.name AS user_name",
		"users.email",
		"users.phone",
		"seats.seat_id",
		"seats.name AS seat_name",
		"seats.price",
		"seats.category",
		"seats.link",
		"transactions.confirmation",
		"transactions.created_at",
		"transactions.updated_at").
		Joins("inner join users on users.user_id = transactions.user_id").
		Joins("inner join seats on seats.seat_id = transactions.seat_id").
		Where("transactions.deleted_at IS NULL").
		Where("transactions.user_id = ?", userId).
		Where("transactions.confirmation = ?", confirmation[0]).
		Or("transactions.confirmation = ?", confirmation[1]).
		Order("transaction_id").
		Scan(&basics)

	var transactionsBuff []model.Transaction
	var transactionBuff model.Transaction
	for _, basic := range basics {
		transactionBuff = model.Transaction{
			TransactionId: basic.TransactionId,
			User:          model.User{UserId: basic.UserId, Name: basic.UserName, Phone: basic.Phone, Email: basic.Email},
			Seat:          model.Seat{SeatId: basic.SeatId, Name: basic.SeatName, Price: basic.Price, Category: basic.Category, Link: basic.Link},
			Confirmation:  basic.Confirmation,
			CreatedAt:     basic.CreatedAt,
			UpdatedAt:     basic.UpdatedAt,
		}
		transactionsBuff = append(transactionsBuff, transactionBuff)
	}
	*transactions = transactionsBuff
	return result
}

func (t *TransactionRepository) GetByOrder(transactions *[]model.Transaction, orderId string) *gorm.DB {
	result := t.db.Where("order_id = ?", orderId).Find(transactions)
	return result
}

func (t *TransactionRepository) GetDetailsByOrder(transactions *[]model.Transaction, orderId string) *gorm.DB {
	var basics []transactionFields
	result := t.db.Table("transactions").Select(
		"transactions.transaction_id",
		"users.user_id",
		"users.name AS user_name",
		"users.email",
		"users.phone",
		"seats.seat_id",
		"seats.name AS seat_name",
		"seats.price",
		"seats.category",
		"seats.link",
		"transactions.created_at",
		"transactions.updated_at").
		Joins("inner join users on users.user_id = transactions.user_id").
		Joins("inner join seats on seats.seat_id = transactions.seat_id").
		Where("transactions.order_id = ?", orderId).
		Order("transaction_id").
		Scan(&basics)
	var transactionsBuff []model.Transaction
	var transactionBuff model.Transaction
	for _, basic := range basics {
		transactionBuff = model.Transaction{
			TransactionId: basic.TransactionId,
			User:          model.User{UserId: basic.UserId, Name: basic.UserName, Phone: basic.Phone, Email: basic.Email},
			Seat:          model.Seat{SeatId: basic.SeatId, Name: basic.SeatName, Price: basic.Price, Category: basic.Category, Link: basic.Link},
			CreatedAt:     basic.CreatedAt,
			UpdatedAt:     basic.UpdatedAt,
		}
		transactionsBuff = append(transactionsBuff, transactionBuff)
	}
	*transactions = transactionsBuff
	return result
}

func (t *TransactionRepository) UpdatePaymentStatus(orderId, vendor, confirmation string) *gorm.DB {
	result := t.db.Model(&model.Transaction{}).Where("order_id = ?", orderId).Updates(model.Transaction{Vendor: vendor, Confirmation: confirmation})
	return result
}

func (t *TransactionRepository) UpdatePaymentStatusById(txId uint64, confirmation string) *gorm.DB {
	result := t.db.Model(&model.Transaction{}).Where("transaction_id = ?", txId).Update("confirmation", confirmation)
	return result
}

func (t *TransactionRepository) UpdateOrderIdById(txId uint64, orderId string) *gorm.DB {
	result := t.db.Model(&model.Transaction{}).Where("transaction_id = ?", txId).Updates(model.Transaction{OrderId: orderId, CreatedAt: time.Now()})
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

func (t *TransactionRepository) SoftDeleteByUserConfirmation(userId uint64, confirmation string) *gorm.DB {
	result := t.db.Where("user_id = ? AND confirmation = ?", userId, confirmation).Delete(&model.Transaction{})
	return result
}

func (t *TransactionRepository) SoftDeleteByOrder(orderId string) *gorm.DB {
	result := t.db.Where("order_id = ?", orderId).Delete(&model.Transaction{})
	return result
}

func (t *TransactionRepository) SoftDeletesById(transactions *[]model.Transaction) *gorm.DB {
	result := t.db.Delete(transactions)
	return result
}
