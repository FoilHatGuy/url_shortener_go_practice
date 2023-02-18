package main

import (
	"shortener/internal/cfg"
	"shortener/internal/server"
	"shortener/internal/storage"
)

func init() {
	cfg.Initialize()
}

func main() {
	storage.RunAutosave()
	server.Run()
}
