package reservation

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
	"reservations/pkg/transport"
)

func MakeHTTPHandler(r *mux.Router, s Service, logger log.Logger) *mux.Router {
	e := MakeServerEndpoints(s)

	options := httpjson.DefaultServerOptions(logger)

	r.Methods("POST").Path("/customer/{id}/reservation").
		Handler(httptransport.NewServer(
			e.BookReservationEndpoint,
			decodeBookReservationRequest,
			httpjson.EncodeResponse,
			options...,
		))

	r.Methods("DELETE").Path("/reservation/{id}").
		Handler(httptransport.NewServer(
			e.DiscardReservationEndpoint,
			decodeDiscardReservationRequest,
			httpjson.EncodeResponse,
			options...,
		))

	r.Methods("PUT").Path("/reservation/{id}").
		Handler(httptransport.NewServer(
			e.EditReservationEndpoint,
			decodeEditReservationRequest,
			httpjson.EncodeResponse,
			options...,
		))

	r.Methods("GET").Path("/customer/{id}/reservations").
		Handler(httptransport.NewServer(
			e.GetReservationHistoryByCustomerEndpoint,
			decodeGetReservationHistoryPerCustomerRequest,
			httpjson.EncodeResponse,
			options...,
		))

	return r
}

func decodeBookReservationRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req bookReservationRequest

	id, err := httpjson.ParseIntPathParam(r, "id", "customer ID")
	if err != nil {
		return nil, err
	}
	req.CustomerID = id

	if e := json.NewDecoder(r.Body).Decode(&req.Reservation); e != nil {
		return nil, e
	}
	return req, nil
}

func decodeDiscardReservationRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	id, err := httpjson.ParseIntPathParam(r, "id", "reservation ID")
	if err != nil {
		return nil, err
	}
	return discardReservationRequest{ReservationID: id}, nil
}

func decodeEditReservationRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req editReservationRequest

	id, err := httpjson.ParseIntPathParam(r, "id", "reservation ID")
	if err != nil {
		return nil, err
	}
	req.ReservationID = id

	if e := json.NewDecoder(r.Body).Decode(&req.Reservation); e != nil {
		return nil, e
	}
	return req, nil
}

func decodeGetReservationHistoryPerCustomerRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	id, err := httpjson.ParseIntPathParam(r, "id", "customer ID")
	if err != nil {
		return nil, err
	}

	return getReservationHistoryPerCustomerRequest{
		CustomerID: id,
		Limit:      httpjson.ParseUintQueryParam(r, "limit"),
		Offset:     httpjson.ParseUintQueryParam(r, "offset"),
	}, nil
}
