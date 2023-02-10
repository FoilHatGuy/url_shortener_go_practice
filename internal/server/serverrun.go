package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"shortener/internal/cfg"
	"shortener/internal/handlers"
	"shortener/internal/storage"
	"time"
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

	t := time.NewTicker(time.Duration(cfg.Storage.AutosaveInterval) * time.Second)
	storage.Database.LoadData()
	go func() {
		for {
			select {
			case <-t.C:
				fmt.Print("AUTOSAVE\n")
				storage.Database.SaveData()
			}
		}
	}()
	log.Fatal(r.Run(cfg.Server.Host + ":" + cfg.Server.Port))
}
