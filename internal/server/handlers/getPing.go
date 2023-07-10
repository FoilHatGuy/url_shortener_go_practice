package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"shortener/internal/storage"
)

// PingDatabase
// Ping server+database activity
func PingDatabase(dbController storage.DatabaseORM) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ping := dbController.Ping(ctx)
		if ping {
			ctx.Status(http.StatusOK)
		} else {
			ctx.Status(http.StatusInternalServerError)
		}
	}
}
