package reservation

import (
	"github.com/doug-martin/goqu/v7"
	"github.com/doug-martin/goqu/v7/exec"
	errors "reservations/pkg/error"
	"reservations/pkg/storage"
	"time"
)

type Repository interface {
	AddReservation(cID int, r *Reservation) (*Reservation, error)
	RemoveReservation(rID int) error
	UpdateReservation(rID int, r *Reservation) (Reservation, error)
	FindReservationsByCustomerID(cID int, opts *storage.QueryOptions) ([]Reservation, error)
}

type reservationRepository struct {
	db storage.Persistence
}

func NewReservationRepository(db storage.Persistence) Repository {
	return &reservationRepository{db: db}
}

func (r *reservationRepository) AddReservation(cID int, res *Reservation) (*Reservation, error) {
	created := time.Now().Unix()

	result, err := r.db.Tx(func(tx *goqu.TxDatabase) exec.QueryExecutor {
		res.Created = created
		res.LastUpdated = created
		res.CustomerID = cID
		return tx.From("reservation").Insert(res)
	})

	if err != nil {
		return nil, errors.DBError.Wrap(err, "error adding new reservation")
	}

	rID, _ := result.LastInsertId()
	res.ReservationID = int(rID)

	return res, nil
}

func (r *reservationRepository) RemoveReservation(rID int) error {
	_, err := r.db.Tx(func(tx *goqu.TxDatabase) exec.QueryExecutor {
		return tx.From("reservation").Where(goqu.Ex{"rid": rID}).Delete()
	})

	if err != nil {
		return errors.DBError.Wrapf(err, "error deleting reservation with ID %d", rID)
	}
	return nil
}

func (r *reservationRepository) UpdateReservation(rID int, res *Reservation) (result Reservation, err error) {
	lastUpdated := time.Now().Unix()

	_, err = r.db.Tx(func(tx *goqu.TxDatabase) exec.QueryExecutor {
		res.LastUpdated = lastUpdated
		return tx.From("reservation").Update(res)
	})

	_, err = r.db.DB.From("reservation").Where(
		goqu.C("rid").Eq(rID),
	).ScanStruct(&result)

	if err != nil {
		return result, errors.DBError.Wrapf(err, "error updating reservation with ID %d", rID)
	}

	return result, nil
}

func (r *reservationRepository) FindReservationsByCustomerID(cID int, opts *storage.QueryOptions) ([]Reservation, error) {
	return nil, nil
}
