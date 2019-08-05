package customer

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

func MakeHTTPHandler(s Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeServerEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(httpjson.EncodeError),
	}

	r.Methods("POST").Path("/customer").
		Handler(httptransport.NewServer(
			e.RegisterCustomerEndpoint,
			decodeRegisterCustomerRequest,
			httpjson.EncodeResponse,
			options...,
		))

	r.Methods("DELETE").Path("/customer/{id}").
		Handler(httptransport.NewServer(
			e.RegisterCustomerEndpoint,
			decodeUnregisterCustomerRequest,
			httpjson.EncodeResponse,
			options...,
		))

	r.Methods("GET").Path("/customer/{id}").
		Handler(httptransport.NewServer(
			e.GetCustomerByIDEndpoint,
			decodeGetCustomerByIDRequest,
			httpjson.EncodeResponse,
			options...,
		))

	r.Methods("GET").Path("/customers").
		Handler(httptransport.NewServer(
			e.GetAllCustomersEndpoint,
			decodeGetAllCustomersRequest,
			httpjson.EncodeResponse,
			options...,
		))

	return r
}

func decodeRegisterCustomerRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req registerCustomerRequest
	if e := json.NewDecoder(r.Body).Decode(&req.Customer); e != nil {
		return nil, e
	}
	return req, nil
}

func decodeUnregisterCustomerRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errors.ValidationError.Newf("missing or invalid customer ID %s", id)
	}
	return unregisterCustomerRequest{CustomerID: id}, nil
}

func decodeGetCustomerByIDRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errors.ValidationError.Newf("missing or invalid customer ID %s", id)
	}
	return getCustomerByIDRequest{CustomerID: id}, nil
}

func decodeGetAllCustomersRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	q := r.URL.Query()
	limit, _ := strconv.ParseUint(q.Get("limit"), 10, 64)
	offset, _ := strconv.ParseUint(q.Get("offset"), 10, 64)

	return getAllCustomersRequest{Limit: uint(limit), Offset: uint(offset)}, nil
}
