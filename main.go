package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/davyj0nes/go-proxy/internal"
)

func main() {
	port := flag.String("port", "8080", "port number to use with the proxy")
	target := flag.String("target", "http://localhost:80", "")
	flag.Parse()

	handler := internal.NewHandler(*target)

	srv := &http.Server{
		Addr:    ":" + *port,
		Handler: handler,
	}

	log.Println("starting proxy...")
	go func(srv *http.Server) {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalln(err)
		}
	}(srv)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	log.Print("proxy has started...")

	killSignal := <-interrupt
	switch killSignal {
	case os.Interrupt:
		log.Print("got SIGINT...")
	case syscall.SIGTERM:
		log.Print("got SIGTERM...")
	}

	log.Print("proxy shutting down...")
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalln(err)
	}

	log.Print("successfully shut down...")
}
