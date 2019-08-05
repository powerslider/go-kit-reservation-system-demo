package httpjson

import (
	"context"
	"encoding/json"
	"net/http"
	errors "reservations/pkg/error"
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
