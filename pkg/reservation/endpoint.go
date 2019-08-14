package reservation

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"reservations/pkg/storage"
)

type Endpoints struct {
	BookReservationEndpoint                 endpoint.Endpoint
	DiscardReservationEndpoint              endpoint.Endpoint
	EditReservationEndpoint                 endpoint.Endpoint
	GetReservationHistoryByCustomerEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		BookReservationEndpoint:                 MakeBookReservationEndpoint(s),
		DiscardReservationEndpoint:              MakeDiscardReservationEndpoint(s),
		EditReservationEndpoint:                 MakeEditReservationEndpoint(s),
		GetReservationHistoryByCustomerEndpoint: MakeGetReservationHistoryPerCustomerEndpoint(s),
	}
}

type discardReservationRequest struct {
	ReservationID int
}

type discardReservationResponse struct {
	Err error `json:"err,omitempty"`
}

func (r discardReservationResponse) HTTPError() error { return r.Err }

// DiscardReservation godoc
// @Summary Discard an existing reservation
// @Description Discard an existing reservation
// @Tags reservation
// @Param id path string true "Reservation ID"
// @Accept  json
// @Produce  json
// @Router /reservation/{id} [delete]
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
	CustomerID  int
	Reservation *Reservation
}

type bookReservationResponse struct {
	Reservation *Reservation `json:"reservation,omitempty"`
	Err         error        `json:"err,omitempty"`
}

func (r bookReservationResponse) HTTPError() error { return r.Err }

// BookReservation godoc
// @Summary Book a new Reservation
// @Description Book a new Reservation
// @Tags reservation
// @Param reservation body reservation.Reservation true "New Reservation"
// @Accept  json
// @Produce  json
// @Success 200 {object} reservation.Reservation
// @Router /reservation [post]
func MakeBookReservationEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(bookReservationRequest)
		r, e := s.BookReservation(ctx, req.CustomerID, req.Reservation)
		return bookReservationResponse{
			Reservation: r,
			Err:         e,
		}, nil
	}
}

type getReservationHistoryPerCustomerRequest struct {
	CustomerID int
	Limit      uint
	Offset     uint
}

type getReservationHistoryPerCustomerResponse struct {
	Reservations []Reservation `json:"reservations,omitempty"`
	Err          error         `json:"err,omitempty"`
}

func (r getReservationHistoryPerCustomerResponse) HTTPError() error { return r.Err }

// GetReservationHistoryPerCustomer godoc
// @Summary List existing reservations per customer ordered by newest.
// @Description List existing reservations per customer ordered by newest.
// @Tags reservation
// @Param limit query int false "Reservation count limit" default(100)
// @Param offset query int false "Reservation count offset" default(0)
// @Param id path string true "Customer ID"
// @Accept  json
// @Produce  json
// @Success 200 {array} reservation.Reservation
// @Router /customer/{id}/reservations [get]
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

type editReservationRequest struct {
	ReservationID int
	Reservation   *Reservation
}

type editReservationResponse struct {
	Reservation Reservation `json:"reservation"`
	Err         error       `json:"err,omitempty"`
}

func (r editReservationResponse) HTTPError() error { return r.Err }

// EditReservation godoc
// @Summary Edit an existing reservation
// @Description Edit an existing reservation
// @Tags reservation
// @Param id path string true "Reservation ID"
// @Accept  json
// @Produce  json
// @Router /reservation/{id} [put]
func MakeEditReservationEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(editReservationRequest)
		r, e := s.EditReservation(ctx, req.ReservationID, req.Reservation)
		return editReservationResponse{
			Reservation: r,
			Err:         e,
		}, nil
	}
}
