package main

import (
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"shortener/internal/handlers"
	"strconv"
)

const ( //config
	urlLength = 10
	host      = "localhost"
	port      = 8080
)

func main() {
	r := chi.NewRouter()

	r.Post("/", handlers.SendURL)
	r.Get("/{shortURL:[a-zA-Z]{"+strconv.FormatInt(urlLength, 10)+"}}", handlers.ReceiveURL)

	log.Fatal(http.ListenAndServe("localhost:8080", r))
}
