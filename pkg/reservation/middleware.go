package reservation

import (
	"context"
	"github.com/go-kit/kit/log"
	"reservations/pkg/storage"
	"time"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(Service) Service

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

func (mw loggingMiddleware) BookReservation(ctx context.Context, r *Reservation) (result *Reservation, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "BookReservation", "id", r.ReservationID, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.BookReservation(ctx, r)
}

func (mw loggingMiddleware) DiscardReservation(ctx context.Context, rID string) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "DiscardReservation", "id", rID, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.DiscardReservation(ctx, rID)
}

func (mw loggingMiddleware) ChangeReservation(ctx context.Context, rID string) (r Reservation, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "ChangeReservation", "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.ChangeReservation(ctx, rID)
}

func (mw loggingMiddleware) GetReservationHistoryPerCustomer(ctx context.Context, cID string, opts *storage.QueryOptions) (result []Reservation, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetReservationHistoryPerCustomer", "id", cID, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetReservationHistoryPerCustomer(ctx, cID, opts)
}
