package main

import (
	"fmt"
	"shortener/internal/cfg"
	"shortener/internal/server"
)

const (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	fmt.Printf("buildVersion\t= %q\n", buildVersion)
	fmt.Printf("buildDate\t= %q\n", buildDate)
	fmt.Printf("buildCommit\t= %q\n", buildCommit)

	cfgData := cfg.New(cfg.FromDefaults(),
		cfg.FromFlags(),
		cfg.FromEnv(),
	)
	server.Run(cfgData)
}
