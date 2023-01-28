package main

import (
	"fmt"
	"github.com/FoilHatGuy/url_shortener_go_practice/cmd/internal/handlers"
	"log"
	"net/http"
)

const ( //config
	host = "localhost"
	port = 8080
)

func main() {
	server := http.Server{
		Addr: fmt.Sprintf("%s:%d", host, port),
	}

	http.HandleFunc("/", handlers.ReceiveUrl)

	log.Fatal(server.ListenAndServe())
}
