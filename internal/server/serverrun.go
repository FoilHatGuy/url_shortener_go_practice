package server

import (
	"github.com/gin-gonic/gin"
	"log"
	"shortener/internal/handlers"
)

func Run() {
	r := gin.Default()

	r.GET("/:shortURL", handlers.ReceiveURL)
	r.POST("/", handlers.SendURL)
	log.Fatal(r.Run())
}
