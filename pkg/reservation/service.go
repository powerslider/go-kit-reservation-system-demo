package reservation

import "context"

type Service interface {
	Create(ctx context.Context, customerID int)
	Delete(ctx context.Context, customerID int)
}

type Reservation struct {
	ReservationID   int    `json:"reservationId" db:"rid"`
	SeatCount       int    `json:"seatCount" db:"seat_count"`
	StartTime       string `json:"startTime" db:"start_time"`
	ReservationName string `json:"reservationName" db:"reservation_name"`
	Phone           string `json:"phone"`
	Comments        string `json:"comments"`
}

type reservationService struct {

}
