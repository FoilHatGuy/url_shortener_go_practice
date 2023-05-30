package server

import (
	"fmt"
	"log"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"

	"shortener/internal/cfg"
	"shortener/internal/server/handlers"
	"shortener/internal/server/middleware"
)

// Run
// Performs initial setup of server router and launches it.
func Run(config *cfg.ConfigT) {
	r := gin.Default()
	baseRouter := r.Group("")
	{
		baseRouter.Use(middleware.Gzip())
		baseRouter.Use(middleware.Gunzip())
		baseRouter.Use(middleware.Cooker())
		baseRouter.GET("/:shortURL", handlers.GetShortURL)
		baseRouter.GET("/ping", handlers.PingDatabase)
		baseRouter.POST("/", handlers.PostURL)
		api := baseRouter.Group("/api")
		{
			api.POST("/shorten", handlers.PostAPIURL)
			api.POST("/shorten/batch", handlers.BatchShorten)
			api.GET("/user/urls", handlers.GetAllOwnedURL)
			api.DELETE("/user/urls", handlers.DeleteLine)
		}
	}
	pprof.Register(r)
	fmt.Println("SERVER LISTENING ON", config.Server.Address)
	log.Fatal(r.Run(config.Server.Address))
}
