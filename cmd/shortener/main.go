package main

import (
	_ "github.com/sakirsensoy/genv/dotenv/autoload"
	"shortener/internal/server"
)

func main() {
	server.Run()

}
