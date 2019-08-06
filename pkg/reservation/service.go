package reservation

import (
	"context"
	"reservations/pkg/storage"
)

type Service interface {
	BookReservation(ctx context.Context, r *Reservation) (*Reservation, error)
	DiscardReservation(ctx context.Context, rID string) error
	ChangeReservation(ctx context.Context, rID string) (Reservation, error)
	GetReservationHistoryPerCustomer(ctx context.Context, cID string, opts *storage.QueryOptions) ([]Reservation, error)
}

type Reservation struct {
	ReservationID   int    `json:"reservationId" db:"rid" goqu:"skipinsert"`
	SeatCount       int    `json:"seatCount" db:"seat_count"`
	StartTime       string `json:"startTime" db:"start_time"`
	ReservationName string `json:"reservationName" db:"reservation_name"`
	Phone           string `json:"phone"`
	Comments        string `json:"comments"`
	Created         int64  `json:"created"`
}

type reservationService struct {
	resRepo Repository
}

func NewReservationService(repo Repository) Service {
	return &reservationService{
		resRepo: repo,
	}
}

func (s *reservationService) BookReservation(ctx context.Context, r *Reservation) (*Reservation, error) {
	res, err := s.resRepo.AddReservation(r)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *reservationService) DiscardReservation(ctx context.Context, rID string) error {
	err := s.resRepo.RemoveReservation(rID)
	if err != nil {
		return err
	}
	return nil
}

func (s *reservationService) ChangeReservation(ctx context.Context, rID string) (r Reservation, err error) {
	r, err = s.resRepo.UpdateReservation(rID)
	if err != nil {
		return r, err
	}
	return r, nil
}

func (s *reservationService) GetReservationHistoryPerCustomer(ctx context.Context, cID string, opts *storage.QueryOptions) ([]Reservation, error) {
	rr, err := s.resRepo.FindReservationsByCustomerID(cID, opts)
	if err != nil {
		return nil, err
	}
	return rr, nil
}
