package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"shortener/internal/storage"
)

// PingDatabase
// Ping server+database activity
func PingDatabase(dbController storage.DatabaseORM) gin.HandlerFunc {
	return func(c *gin.Context) {
		ping := dbController.Ping(c)
		if ping {
			c.Status(http.StatusOK)
		} else {
			c.Status(http.StatusInternalServerError)
		}
	}
}
