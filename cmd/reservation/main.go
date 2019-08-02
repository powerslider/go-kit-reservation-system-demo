package main

import (
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	"net/http"
	"os"
	"os/signal"
	"reservations/pkg/customer"
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

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var r customer.Repository
	{
		r = customer.NewCustomerRepository(*db)
	}

	var s customer.Service
	{
		s = customer.NewCustomerService(r)
		s = customer.LoggingMiddleware(logger)(s)
	}

	var h http.Handler
	{
		h = customer.MakeHTTPHandler(s, log.With(logger, "component", "HTTP"))
	}

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logger.Log("transport", "HTTP", "addr", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, h)
	}()

	logger.Log("exit", <-errs)
}
