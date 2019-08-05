package reservation

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"reservations/pkg/storage"
)

type Endpoints struct {
	BookReservationEndpoint                 endpoint.Endpoint
	DiscardReservationEndpoint              endpoint.Endpoint
	ChangeReservationEndpoint               endpoint.Endpoint
	GetReservationHistoryByCustomerEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		BookReservationEndpoint:                 MakeBookReservationEndpoint(s),
		DiscardReservationEndpoint:              MakeDiscardReservationEndpoint(s),
		ChangeReservationEndpoint:               MakeChangeReservationEndpoint(s),
		GetReservationHistoryByCustomerEndpoint: MakeGetReservationHistoryPerCustomerEndpoint(s),
	}
}

type discardReservationRequest struct {
	ReservationID string
}

type discardReservationResponse struct {
	Err error `json:"err,omitempty"`
}

func (r discardReservationResponse) HTTPError() error { return r.Err }

func MakeDiscardReservationEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(discardReservationRequest)
		e := s.DiscardReservation(ctx, req.ReservationID)
		return discardReservationResponse{
			Err: e,
		}, nil
	}
}

type bookReservationRequest struct {
	Reservation *Reservation
}

type bookReservationResponse struct {
	Reservation *Reservation `json:"reservation,omitempty"`
	Err         error        `json:"err,omitempty"`
}

func (r bookReservationResponse) HTTPError() error { return r.Err }

func MakeBookReservationEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(bookReservationRequest)
		r, e := s.BookReservation(ctx, req.Reservation)
		return bookReservationResponse{
			Reservation: r,
			Err:         e,
		}, nil
	}
}

type getReservationHistoryPerCustomerRequest struct {
	CustomerID string
	Limit      uint
	Offset     uint
}

type getReservationHistoryPerCustomerResponse struct {
	Reservations []Reservation `json:"reservations,omitempty"`
	Err          error         `json:"err,omitempty"`
}

func (r getReservationHistoryPerCustomerResponse) HTTPError() error { return r.Err }

func MakeGetReservationHistoryPerCustomerEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getReservationHistoryPerCustomerRequest)
		rr, e := s.GetReservationHistoryPerCustomer(ctx, req.CustomerID, &storage.QueryOptions{
			Limit:  req.Limit,
			Offset: req.Offset,
		})
		return getReservationHistoryPerCustomerResponse{
			Reservations: rr,
			Err:          e,
		}, nil
	}
}

type changeReservationRequest struct {
	ReservationID string
	Reservation   Reservation
}

type changeReservationResponse struct {
	Reservation Reservation `json:"reservation"`
	Err         error       `json:"err,omitempty"`
}

func (r changeReservationResponse) HTTPError() error { return r.Err }

func MakeChangeReservationEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(changeReservationRequest)
		r, e := s.ChangeReservation(ctx, req.ReservationID)
		return changeReservationResponse{
			Reservation: r,
			Err:         e,
		}, nil
	}
}
