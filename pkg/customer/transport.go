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
	"strconv"
)

func MakeHTTPHandler(s Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeServerEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Methods("POST").Path("/customer").
		Handler(httptransport.NewServer(
			e.RegisterCustomerEndpoint,
			decodeRegisterCustomerRequest,
			encodeResponse,
			options...,
		))

	r.Methods("DELETE").Path("/customer/{id}").
		Handler(httptransport.NewServer(
			e.RegisterCustomerEndpoint,
			decodeUnregisterCustomerRequest,
			encodeResponse,
			options...,
		))

	r.Methods("GET").Path("/customer/{id}").
		Handler(httptransport.NewServer(
			e.GetCustomerByIDEndpoint,
			decodeGetCustomerByIDRequest,
			encodeResponse,
			options...,
		))

	r.Methods("GET").Path("/customers").
		Handler(httptransport.NewServer(
			e.GetAllCustomersEndpoint,
			decodeGetAllCustomersRequest,
			encodeResponse,
			options...,
		))

	return r
}

// errorer is implemented by all concrete response types that may contain
// errors. It allows us to change the HTTP response code without needing to
// trigger an endpoint (transport-level) error. For more information, read the
// big comment in endpoints.go.
type errorer interface {
	error() error
}

// encodeResponse is the common method to encode all response types to the
// client. I chose to do it this way because, since we're using JSON, there's no
// reason to provide anything more specific. It's certainly possible to
// specialize on a per-response (per-method) basis.
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	e, ok := response.(errorer)
	resErr := e.error()

	if ok && resErr != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, resErr, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	// jsonErr, _ := json.Marshal(err)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err,
	})
	// json.NewEncoder(w).Encode(err)
}

func codeFrom(err error) int {
	customErr := err.(errors.AppError)

	switch customErr.ErrorType {
	case errors.NotFound:
		return http.StatusNotFound
	// case ErrAlreadyExists, ErrInconsistentIDs:
	// 	return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
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
