package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shortener/internal/storage"
)

// PingDatabase
// Ping server+database activity
func PingDatabase(c *gin.Context) {
	ping := storage.Controller.Ping(c)
	if ping {
		c.Status(http.StatusOK)
	} else {
		c.Status(http.StatusInternalServerError)
	}

}
