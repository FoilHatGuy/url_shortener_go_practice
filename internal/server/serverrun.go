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
	baseRouter := r.Group("")
	{
		baseRouter.Use(handlers.ArchiveData())
		baseRouter.GET("/:shortURL", handlers.GetShortURL)
		baseRouter.POST("/", handlers.PostURL)
		api := baseRouter.Group("/api")
		{
			api.POST("/shorten", handlers.PostAPIURL)
		}
	}

	t := time.NewTicker(time.Duration(cfg.Storage.AutosaveInterval) * time.Second)
	storage.Database.LoadData()
	go func() {
		for range t.C {
			//fmt.Print("AUTOSAVE\n")
			storage.Database.SaveData()
		}
	}()
	fmt.Println("SERVER LISTENING ON", cfg.Server.Host+":"+cfg.Server.Port)
	log.Fatal(r.Run(cfg.Server.Host + ":" + cfg.Server.Port))
}
