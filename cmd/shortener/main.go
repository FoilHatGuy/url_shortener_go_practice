package main

import (
	"fmt"
	_ "github.com/sakirsensoy/genv/dotenv/autoload"
	"shortener/internal/cfg"
	"shortener/internal/server"
)

func init() {
	cfg.Initialize()
}

func main() {
	fmt.Println(cfg.Server.Host)
	server.Run()
}
