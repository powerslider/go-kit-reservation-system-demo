package main

import (
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"os/signal"
	"reservations/pkg/customer"
	"reservations/pkg/reservation"
	"reservations/pkg/storage"
	"syscall"
)

func main() {

	var (
		httpAddr = flag.String("http.addr", ":8080", "HTTP listen address")
	)
	flag.Parse()

	db, err := storage.NewDB("reservations")
	if err != nil {
		panic(err)
	}

	logger := log.NewLogfmtLogger(os.Stdout)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	r := mux.NewRouter()
	r = initCustomerHandler(r, db, logger)
	r = initReservationHandler(r, db, logger)

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logger.Log("transport", "HTTP", "addr", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, r)
	}()

	logger.Log("exit", <-errs)
}

func initCustomerHandler(router *mux.Router, db *storage.Persistence, logger log.Logger) *mux.Router {
	r := customer.NewCustomerRepository(*db)
	s := customer.NewCustomerService(r)
	s = customer.LoggingMiddleware(logger)(s)
	return customer.MakeHTTPHandler(router, s, logger)
}

func initReservationHandler(router *mux.Router, db *storage.Persistence, logger log.Logger) *mux.Router {
	r := reservation.NewReservationRepository(*db)
	s := reservation.NewReservationService(r)
	s = reservation.LoggingMiddleware(logger)(s)
	return reservation.MakeHTTPHandler(router, s, logger)
}
