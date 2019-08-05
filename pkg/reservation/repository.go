package reservation

import "reservations/pkg/storage"

type Repository interface {
	AddReservation(r *Reservation) (*Reservation, error)
	RemoveReservation(rID string) error
	UpdateReservation(rID string) (Reservation, error)
	FindReservationsByCustomerID(cID string, opts *storage.QueryOptions) ([]Reservation, error)
}

type reservationRepository struct {
	db storage.Persistence
}

func NewReservationRepository(db storage.Persistence) Repository {
	return &reservationRepository{db: db}
}

func (r *reservationRepository) AddReservation(res *Reservation) (*Reservation, error) {
	return res, nil
}

func (r *reservationRepository) RemoveReservation(rID string) error {
	return nil
}

func (r *reservationRepository) UpdateReservation(rID string) (Reservation, error) {
	return Reservation{}, nil
}

func (r *reservationRepository) FindReservationsByCustomerID(cID string, opts *storage.QueryOptions) ([]Reservation, error) {
	return nil, nil
}
