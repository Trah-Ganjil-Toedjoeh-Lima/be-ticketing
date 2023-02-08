package validation

type SeatDetailsResponse struct {
	SeatId uint   `json:"seat_id"`
	Name   string `json:"name"`
	Price  uint   `json:"price"`
	Status string `json:"status"`
}

type SeatResponse struct {
	Name  string `json:"name"`
	Price uint   `json:"price"`
}

type SeatReservationRequest struct {
	SeatIds []uint `json:"data" binding:"required,unique,max=5"`
}
