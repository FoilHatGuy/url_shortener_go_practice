package server

import (
	"shortener/internal/auth"
	"shortener/internal/cfg"
	"shortener/internal/storage"
)

func ExampleRun() {
	// First we need to New config, since it is needed in security, storage and server setup functions
	cfgData := cfg.New(cfg.FromDefaults())

	// Then we set up all additional modules used by server
	auth.New(cfgData)
	storage.New(cfgData)

	// Finally, we run the server
	Run(cfgData)

	// Be aware that server blocks further execution,
	// so if you need to perform further actions while the server is running, we use:
	go Run(cfgData)
}
