package main

import (
	"fmt"
	"github.com/FoilHatGuy/url_shortener_go_practice/cmd/handlers"
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

	http.HandleFunc("/", handlers.ReceiveURL)

	log.Fatal(server.ListenAndServe())
}
