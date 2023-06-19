package server

import (
	"fmt"
	"log"

	"shortener/internal/auth"
	"shortener/internal/storage"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"

	"shortener/internal/cfg"
	"shortener/internal/server/handlers"
	"shortener/internal/server/middleware"
)

// Run
// Performs initial setup of server router and launches it.
func Run(config *cfg.ConfigT) {
	dbController := storage.New(config)

	r := gin.Default()
	baseRouter := r.Group("")

	baseRouter.Use(middleware.Gzip())
	baseRouter.Use(middleware.Gunzip())
	baseRouter.Use(middleware.Cooker(config, auth.New(config)))
	baseRouter.GET("/:shortURL",
		handlers.GetShortURL(dbController, config))
	baseRouter.GET("/ping", handlers.PingDatabase(dbController))
	baseRouter.POST("/", handlers.PostURL(dbController, config))

	api := baseRouter.Group("/api")
	api.POST("/shorten", handlers.PostAPIURL(dbController, config))
	api.POST("/shorten/batch", handlers.BatchShorten(dbController, config))
	api.GET("/user/urls", handlers.GetAllOwnedURL(dbController))
	api.DELETE("/user/urls", handlers.DeleteLine(dbController))

	pprof.Register(r)
	fmt.Println("SERVER LISTENING ON", config.Server.Address)
	if config.Server.IsHTTPS {
		certPEM, certKey := auth.GetCertificate()
		log.Fatal(r.RunTLS(config.Server.Address, certPEM, certKey))
	} else {
		log.Fatal(r.Run(config.Server.Address))
	}
}
