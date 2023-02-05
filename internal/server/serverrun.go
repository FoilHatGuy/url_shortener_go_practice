package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"shortener/internal/cfg"
	"shortener/internal/handlers"
)

func Run() {
	r := gin.Default()
	fmt.Print(cfg.Router.BaseURL)
	baseRouter := r.Group(cfg.Router.BaseURL)
	{
		baseRouter.GET("/:shortURL", handlers.GetShortURL)
		baseRouter.POST("/", handlers.PostURL)
		api := baseRouter.Group("/api")
		{
			api.POST("/shorten", handlers.PostApiURL)
		}
	}
	log.Fatal(r.Run(cfg.Server.Host + ":" + cfg.Server.Port))
}
