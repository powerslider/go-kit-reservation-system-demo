package customer

import (
	"github.com/doug-martin/goqu/v7"
	"github.com/doug-martin/goqu/v7/exec"
	errors "reservations/pkg/error"
	"reservations/pkg/storage"
	"time"
)

const (
	defaultLimit  uint = 100
	defaultOffset uint = 0
)

type Repository interface {
	AddCustomer(c *Customer) (*Customer, error)
	RemoveCustomer(cID int) error
	FindAllCustomers(opts *storage.QueryOptions) ([]Customer, error)
	FindCustomerByID(cID int) (Customer, error)
}

type customerRepository struct {
	db storage.Persistence
}

func NewCustomerRepository(db storage.Persistence) Repository {
	return &customerRepository{db: db}
}

func (r *customerRepository) AddCustomer(c *Customer) (*Customer, error) {
	created := time.Now().Unix()

	result, err := r.db.Tx(func(tx *goqu.TxDatabase) exec.QueryExecutor {
		c.Created = created
		c.LastUpdated = created
		return tx.From("customer").Insert(c)
	});
	if err != nil {
		return nil, errors.DBError.Wrap(err, "error adding new customer")
	}

	cID, _ := result.LastInsertId()
	c.CustomerID = int(cID)

	return c, nil
}

func (r *customerRepository) RemoveCustomer(cID int) error {
	_, err := r.db.Tx(func(tx *goqu.TxDatabase) exec.QueryExecutor {
		return tx.From("customer").Where(goqu.Ex{"cid": cID}).Delete()
	})

	if err != nil {
		return errors.DBError.Wrapf(err, "error deleting customers with id %d", cID)
	}
	return nil
}

func (r *customerRepository) FindAllCustomers(opts *storage.QueryOptions) (cc []Customer, err error) {
	if opts.Limit == 0 {
		opts.Limit = defaultLimit
	}

	err = r.db.DB.From("customer").
		Limit(opts.Limit).
		Offset(opts.Offset).
		ScanStructs(&cc)

	if err != nil {
		return nil, errors.DBError.Wrapf(err, "error getting all customers")
	}
	return cc, nil
}

func (r *customerRepository) FindCustomerByID(cID int) (c Customer, err error) {
	found, err := r.db.DB.From("customer").Where(
		goqu.C("cid").Eq(cID),
	).ScanStruct(&c)

	if !found {
		return c, errors.NotFound.Newf("customer with ID %d not found", cID).
			AddContext("CustomerID", "non existent ID")
	}

	if err != nil {
		return c, errors.DBError.Wrapf(err, "error getting customer with ID %d", cID)
	}

	return c, nil
}
