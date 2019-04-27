package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/davyj0nes/go-proxy/internal"
	"github.com/sirupsen/logrus"
)

func main() {
	port := flag.String("port", "8080", "port number to use with the proxy")
	target := flag.String("target", "http://localhost:80", "")
	flag.Parse()

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	handler := internal.NewHandler(logger, *target)

	srv := &http.Server{
		Addr:    ":" + *port,
		Handler: handler,
	}

	logger.Info("starting proxy...")
	go func(srv *http.Server) {
		if err := srv.ListenAndServe(); err != nil {
			logger.Fatalln(err)
		}
	}(srv)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	logger.Info("proxy has started...")

	killSignal := <-interrupt
	switch killSignal {
	case os.Interrupt:
		logger.Info("got SIGINT...")
	case syscall.SIGTERM:
		logger.Info("got SIGTERM...")
	}

	logger.Info("proxy shutting down...")
	if err := srv.Shutdown(context.Background()); err != nil {
		logger.Fatalln(err)
	}

	logger.Info("successfully shut down...")
}
