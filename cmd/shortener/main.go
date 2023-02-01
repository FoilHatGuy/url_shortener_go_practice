package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"shortener/internal/handlers"
)

const ( //config
	host = "localhost"
	port = 8080
)

func main() {
	r := chi.NewRouter()
	server := http.Server{
		Addr: fmt.Sprintf("%s:%d", host, port),
	}

	r.Post("/", handlers.SendURL)
	r.Get("/{shortURL:[a-zA-Z]{10}}", handlers.ReceiveURL)

	//http.HandleFunc("/", handlers.ReceiveURL)

	log.Fatal(server.ListenAndServe())
}
