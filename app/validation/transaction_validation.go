package validation

type SnapDetails struct {
	TransactionId uint64 `json:"transaction_id"`
	OrderId       string `json:"order_id"`
	UserId        uint64 `json:"user_id"`
	SeatId        uint   `json:"seat_id"`
	SeatName      string `json:"seat_name"`
	SeatPrice     uint   `json:"seat_price"`
	UserName      string `json:"user_name"`
	UserEmail     string `json:"user_email"`
	UserPhone     string `json:"user_phone"`
}
