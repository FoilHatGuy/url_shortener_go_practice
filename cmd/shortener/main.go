package main

import (
	"shortener/internal/server"
)

const ( //config
	urlLength = 10
	host      = "localhost"
	port      = 8080
)

func main() {
	server.Run()

}
