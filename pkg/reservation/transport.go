package reservation

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
	errors "reservations/pkg/error"
	"reservations/pkg/transport"
	"strconv"
)

func MakeHTTPHandler(r *mux.Router, s Service, logger log.Logger) *mux.Router {
	e := MakeServerEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(httpjson.EncodeError),
	}

	r.Methods("POST").Path("/reservation").
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
			e.ChangeReservationEndpoint,
			decodeChangeReservationRequest,
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
	if e := json.NewDecoder(r.Body).Decode(&req.Reservation); e != nil {
		return nil, e
	}
	return req, nil
}

func decodeDiscardReservationRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errors.ValidationError.Newf("missing or invalid reservation ID %s", id)
	}
	return discardReservationRequest{ReservationID: id}, nil
}

func decodeChangeReservationRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errors.ValidationError.Newf("missing or invalid reservation ID %s", id)
	}

	var req changeReservationRequest
	if e := json.NewDecoder(r.Body).Decode(&req.Reservation); e != nil {
		return nil, e
	}
	return req, nil
}

func decodeGetReservationHistoryPerCustomerRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errors.ValidationError.Newf("missing or invalid customer ID %s", id)
	}

	q := r.URL.Query()
	limit, _ := strconv.ParseUint(q.Get("limit"), 10, 64)
	offset, _ := strconv.ParseUint(q.Get("offset"), 10, 64)

	return getReservationHistoryPerCustomerRequest{
		CustomerID: id,
		Limit:      uint(limit),
		Offset:     uint(offset),
	}, nil
}
