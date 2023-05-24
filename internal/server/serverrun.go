package server

import (
	"fmt"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"log"
	"shortener/internal/cfg"
	"shortener/internal/server/handlers"
	"shortener/internal/server/middleware"
)

func Run() {
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
	fmt.Println("SERVER LISTENING ON", cfg.Server.Address)
	log.Fatal(r.Run(cfg.Server.Address))
}
