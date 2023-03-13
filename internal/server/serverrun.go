package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"shortener/internal/cfg"
)

func Run() {
	r := gin.Default()
	baseRouter := r.Group("")
	{
		baseRouter.Use(Gzip())
		baseRouter.Use(Gunzip())
		baseRouter.Use(Cooker())
		baseRouter.GET("/:shortURL", getShortURL)
		baseRouter.GET("/ping", pingDatabase)
		baseRouter.POST("/", postURL)
		api := baseRouter.Group("/api")
		{
			api.POST("/shorten", postAPIURL)
			api.POST("/shorten/batch", batchShorten)
			api.GET("/user/urls", getAllOwnedURL)
		}
	}

	fmt.Println("SERVER LISTENING ON", cfg.Server.Address)
	log.Fatal(r.Run(cfg.Server.Address))
}
