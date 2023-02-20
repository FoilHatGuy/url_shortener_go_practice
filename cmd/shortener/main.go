package main

import (
	"shortener/internal/cfg"
	"shortener/internal/server"
	"shortener/internal/storage"
)

//func init() {
//}

func main() {
	cfg.Initialize()
	storage.RunAutosave()
	server.Run()
}
