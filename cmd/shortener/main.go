package main

import (
	"shortener/internal/cfg"
	"shortener/internal/server"
)

func main() {
	cfgData := cfg.New(cfg.FromDefaults(),
		cfg.FromFlags(),
		cfg.FromEnv(),
	)
	server.Run(cfgData)
}
