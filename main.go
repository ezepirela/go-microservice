package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ezepirela/go-microservice/handlers"
)

func main() {

	l := log.New(os.Stdout, "product-api", log.LstdFlags)
	// hh := handlers.NewHello(l)
	productHandler := handlers.NewProducts(l)

	sm := http.NewServeMux()
	// sm.Handle("/", hh)
	sm.Handle("/", productHandler)

	s := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	signChannel := make(chan os.Signal)
	signal.Notify(signChannel, os.Interrupt)
	signal.Notify(signChannel, os.Kill)

	sig := <-signChannel
	l.Println("Received terminate, graceful shutdown", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}
