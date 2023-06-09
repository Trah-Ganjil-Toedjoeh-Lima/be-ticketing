package validation

type ReservationResponse struct {
	SeatId   uint   `json:"seat_id"`
	Name     string `json:"name"`
	Price    uint   `json:"price"`
	Status   string `json:"status"`
	Row      string `json:"row"`
	Column   uint   `json:"column"`
	Category string `json:"category"`
}

type BasicSeatResponse struct {
	Name     string `json:"name"`
	Price    uint   `json:"price"`
	Category string `json:"category"`
	Link     string `json:"link"`
}

type ReservationRequest struct {
	SeatIds []uint `json:"data" binding:"required,unique,max=5"`
}
