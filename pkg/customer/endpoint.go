package customer

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"reservations/pkg/storage"
)

type Endpoints struct {
	RegisterCustomerEndpoint   endpoint.Endpoint
	UnregisterCustomerEndpoint endpoint.Endpoint
	GetAllCustomersEndpoint    endpoint.Endpoint
	GetCustomerByIDEndpoint    endpoint.Endpoint
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		RegisterCustomerEndpoint:   MakeRegisterCustomerEndpoint(s),
		UnregisterCustomerEndpoint: MakeUnregisterCustomerEndpoint(s),
		GetAllCustomersEndpoint:    MakeGetAllCustomersEndpoint(s),
		GetCustomerByIDEndpoint:    MakeGetCustomerByIDEndpoint(s),
	}
}

type unregisterCustomerRequest struct {
	CustomerID int
}

type unregisterCustomerResponse struct {
	Err error `json:"err,omitempty"`
}

func (r unregisterCustomerResponse) HTTPError() error { return r.Err }

// UnregisterCustomer godoc
// @Summary Unregister an existing customer
// @Description Unregister an existing customer
// @Tags customer
// @Param id path string true "Customer ID"
// @Accept  json
// @Produce  json
// @Router /customer/{id} [delete]
func MakeUnregisterCustomerEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(unregisterCustomerRequest)
		e := s.UnregisterCustomer(ctx, req.CustomerID)
		return unregisterCustomerResponse{
			Err: e,
		}, nil
	}
}

type registerCustomerRequest struct {
	Customer *Customer
}

type registerCustomerResponse struct {
	Customer *Customer `json:"customer,omitempty"`
	Err      error     `json:"err,omitempty"`
}

func (r registerCustomerResponse) HTTPError() error { return r.Err }

// RegisterCustomer godoc
// @Summary Register a new Customer
// @Description Register a new Customer
// @Tags customer
// @Param customer body customer.Customer true "New Customer"
// @Accept  json
// @Produce  json
// @Success 200 {object} customer.Customer
// @Router /customer [post]
func MakeRegisterCustomerEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(registerCustomerRequest)
		c, e := s.RegisterCustomer(ctx, req.Customer)
		return registerCustomerResponse{
			Customer: c,
			Err:      e,
		}, nil
	}
}

type getAllCustomersRequest struct {
	Limit  uint
	Offset uint
}

type getAllCustomersResponse struct {
	Customers []Customer `json:"customers,omitempty"`
	Err       error      `json:"err,omitempty"`
}

func (r getAllCustomersResponse) HTTPError() error { return r.Err }

// GetAllCustomers godoc
// @Summary List existing customers
// @Description List existing customers
// @Tags customer
// @Param limit query int false "Customer count limit" default(100)
// @Param offset query int false "Customer count offset" default(0)
// @Accept  json
// @Produce  json
// @Success 200 {array} customer.Customer
// @Router /customers [get]
func MakeGetAllCustomersEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getAllCustomersRequest)
		cc, e := s.GetAllCustomers(ctx, &storage.QueryOptions{
			Limit:  req.Limit,
			Offset: req.Offset,
		})
		return getAllCustomersResponse{
			Customers: cc,
			Err:       e,
		}, nil
	}
}

type getCustomerByIDRequest struct {
	CustomerID int
}

type getCustomerByIDResponse struct {
	Customer Customer `json:"customer,omitempty"`
	Err      error    `json:"err,omitempty"`
}

func (r getCustomerByIDResponse) HTTPError() error { return r.Err }

// GetCustomerByID godoc
// @Summary Get an existing customer
// @Description Get an existing customer
// @Tags customer
// @Param id path string true "Customer ID"
// @Accept  json
// @Produce  json
// @Router /customer/{id} [get]
func MakeGetCustomerByIDEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getCustomerByIDRequest)
		c, e := s.GetCustomerByID(ctx, req.CustomerID)
		return getCustomerByIDResponse{
			Customer: c,
			Err:      e,
		}, nil
	}
}
