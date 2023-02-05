package validation

type SeatResponse struct {
	SeatId uint   `json:"seat_id"`
	Name   string `json:"name"`
	Price  uint   `json:"price"`
	Status string `json:"status"`
}

type SeatReservationRequest struct {
	SeatIds []uint `json:"data" binding:"required,unique,max=5"`
}
