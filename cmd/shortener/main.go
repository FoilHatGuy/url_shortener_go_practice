package main

import (
	"shortener/internal/cfg"
	"shortener/internal/server"
)

//func init() {
//}

func main() {
	cfgData := cfg.New(cfg.FromDefaults(),
		cfg.FromFlags(),
		cfg.FromEnv(),
	)
	server.Run(cfgData)
}
