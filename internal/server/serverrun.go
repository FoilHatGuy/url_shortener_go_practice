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
		baseRouter.GET("/:shortURL", getShortURL)
		baseRouter.POST("/", postURL)
		api := baseRouter.Group("/api")
		{
			api.POST("/shorten", postAPIURL)
		}
	}

	fmt.Println("SERVER LISTENING ON", cfg.Server.Address)
	log.Fatal(r.Run(cfg.Server.Address))
}
