package customer

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
	id, err := httpjson.ParseIntPathParam(r, "id", "customer ID")
	if err != nil {
		return nil, err
	}
	return unregisterCustomerRequest{CustomerID: id}, nil
}

func decodeGetCustomerByIDRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	id, err := httpjson.ParseIntPathParam(r, "id", "customer ID")
	if err != nil {
		return nil, err
	}
	return getCustomerByIDRequest{CustomerID: id}, nil
}

func decodeGetAllCustomersRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	return getAllCustomersRequest{
		Limit:  httpjson.ParseUintQueryParam(r, "limit"),
		Offset: httpjson.ParseUintQueryParam(r, "offset"),
	}, nil
}
