package httpjson

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

// HTTPErrorer is implemented by all concrete response types that may contain
// errors. It allows us to change the HTTP response code without needing to
// trigger an endpoint (transport-level) error.
type HTTPErrorer interface {
	HTTPError() error
}

// EncodeResponse is the common method to encode all response types to the
// client. Since we're using JSON, there's no reason to provide anything more specific.
// There is also the option to specialize on a per-response (per-method) basis.
func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	e, ok := response.(HTTPErrorer)
	resErr := e.HTTPError()

	if ok && resErr != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		EncodeError(ctx, resErr, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func EncodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))

	if e := json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err,
	}); e != nil {
		panic(err)
	}
}

func HTTPRequestFinalizer(logger log.Logger) httptransport.ServerFinalizerFunc {
	return func(ctx context.Context, code int, r *http.Request) {
		route := r.URL.Path
		query := r.URL.RawQuery

		var keyvals []interface{}
		keyvals = append(keyvals, "proto", r.Proto, "method", r.Method, "route", route, "status_code", code)
		if len(query) > 0 {
			keyvals = append(keyvals, "query", query)
		}

		logger.Log(keyvals...)
	}
}

func DefaultServerOptions(logger log.Logger) []httptransport.ServerOption {
	return []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(EncodeError),
		httptransport.ServerFinalizer(HTTPRequestFinalizer(logger)),
	}
}

func ParseIntPathParam(req *http.Request, paramName string, paramDesc string) (int, error) {
	vars := mux.Vars(req)
	id, ok := vars[paramName]
	if !ok {
		return 0, errors.ValidationError.Newf("missing or invalid %s %s", paramDesc, id)
	}
	p, _ := strconv.ParseInt(id, 10, 64)

	return int(p), nil
}

func ParseUintQueryParam(req *http.Request, paramName string) uint {
	q := req.URL.Query()
	p, _ := strconv.ParseUint(q.Get(paramName), 10, 64)
	return uint(p)
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
