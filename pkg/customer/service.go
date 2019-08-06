package customer

import (
	"context"
	"reservations/pkg/storage"
)

type Service interface {
	RegisterCustomer(ctx context.Context, c *Customer) (*Customer, error)
	UnregisterCustomer(ctx context.Context, cID string) error
	GetAllCustomers(ctx context.Context, opts *storage.QueryOptions) ([]Customer, error)
	GetCustomerByID(ctx context.Context, cID string) (Customer, error)
}

type Customer struct {
	CustomerID string `json:"customerId" db:"cid" goqu:"skipinsert"`
	FirstName  string `json:"firstName" db:"first_name"`
	LastName   string `json:"lastName" db:"last_name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Created    int64  `json:"created"`
}

type customerService struct {
	custRepo Repository
}

func NewCustomerService(repo Repository) Service {
	return &customerService{
		custRepo: repo,
	}
}

func (s *customerService) RegisterCustomer(ctx context.Context, c *Customer) (*Customer, error) {
	res, err := s.custRepo.AddCustomer(c)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *customerService) UnregisterCustomer(ctx context.Context, cID string) error {
	err := s.custRepo.RemoveCustomer(cID)
	if err != nil {
		return err
	}
	return nil
}

func (s *customerService) GetAllCustomers(ctx context.Context, opts *storage.QueryOptions) ([]Customer, error) {
	cc, err := s.custRepo.FindAllCustomers(opts)
	if err != nil {
		return nil, err
	}
	return cc, nil
}

func (s *customerService) GetCustomerByID(ctx context.Context, cID string) (Customer, error) {
	c, err := s.custRepo.FindCustomerByID(cID)
	if err != nil {
		return c, err
	}
	return c, nil
}
