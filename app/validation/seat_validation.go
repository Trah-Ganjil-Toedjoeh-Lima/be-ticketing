package validation

type ReservationResponse struct {
	SeatId uint   `json:"seat_id"`
	Name   string `json:"name"`
	Price  uint   `json:"price"`
	Status string `json:"status"`
}

type BasicResponse struct {
	Name  string `json:"name"`
	Price uint   `json:"price"`
}

type ReservationRequest struct {
	SeatIds []uint `json:"data" binding:"required,unique,max=5"`
}
