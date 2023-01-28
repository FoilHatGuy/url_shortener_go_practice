package main

import (
	"fmt"
	"log"
	"net/http"
	"shortener/internal/handlers"
)

const ( //config
	host = "localhost"
	port = 8080
)

func main() {
	server := http.Server{
		Addr: fmt.Sprintf("%s:%d", host, port),
	}

	http.HandleFunc("/", handlers.ReceiveURL)

	log.Fatal(server.ListenAndServe())
}