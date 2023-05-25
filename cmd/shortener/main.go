package main

import (
	"shortener/internal/cfg"
	"shortener/internal/security"
	"shortener/internal/server"
	"shortener/internal/storage"
)

//func init() {
//}

func main() {
	cfgData := cfg.Initialize()
	security.Init(cfgData)
	storage.Initialize(cfgData)
	server.Run(cfgData)
}
