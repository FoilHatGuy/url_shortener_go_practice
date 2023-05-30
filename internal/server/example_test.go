package server

import (
	"shortener/internal/cfg"
	"shortener/internal/security"
	"shortener/internal/storage"
)

func ExampleRun() {
	// First we need to initialize config, since it is needed in security, storage and server setup functions
	cfgData := cfg.Initialize()

	// Then we setup all additional modules used by server
	security.Init(cfgData)
	storage.Initialize(cfgData)

	// Finally, we run the server
	Run(cfgData)

	// Be aware that server blocks further execution,
	// so if you need to perform further actions while the server is running, we use:
	go Run(cfgData)
}
