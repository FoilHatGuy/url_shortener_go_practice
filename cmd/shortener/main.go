package main

import (
	_ "github.com/sakirsensoy/genv/dotenv/autoload"
	"shortener/internal/cfg"
	"shortener/internal/server"
)

func main() {
	cfg.Init()
	server.Run()
}
