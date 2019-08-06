package reservation

import (
	"github.com/doug-martin/goqu/v7"
	"github.com/doug-martin/goqu/v7/exec"
	errors "reservations/pkg/error"
	"reservations/pkg/storage"
	"time"
)

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
	created := time.Now().Unix()

	if err := r.db.Tx(func(tx *goqu.TxDatabase) exec.QueryExecutor {
		res.Created = created
		return tx.From("reservation").Insert(res)
	}); err != nil {
		return nil, errors.DBError.Wrap(err, "error adding new reservation")
	}

	return res, nil
}

func (r *reservationRepository) RemoveReservation(rID string) error {
	if err := r.db.Tx(func(tx *goqu.TxDatabase) exec.QueryExecutor {
		return tx.From("reservation").Where(goqu.Ex{"rid": rID}).Delete()
	}); err != nil {
		return errors.DBError.Wrapf(err, "error deleting reservation with ID %s", rID)
	}
	return nil
}

func (r *reservationRepository) UpdateReservation(rID string) (Reservation, error) {
	return Reservation{}, nil
}

func (r *reservationRepository) FindReservationsByCustomerID(cID string, opts *storage.QueryOptions) ([]Reservation, error) {
	return nil, nil
}
