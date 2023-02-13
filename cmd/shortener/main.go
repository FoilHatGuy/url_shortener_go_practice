package main

import (
	"fmt"
	"shortener/internal/cfg"
	"shortener/internal/server"
)

func init() {
	//cfg.Initialize()
}

func main() {
	fmt.Println(cfg.Server.Host)
	server.Run()
}
