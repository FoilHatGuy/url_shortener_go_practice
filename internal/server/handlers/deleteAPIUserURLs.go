package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"shortener/internal/storage"
)

// DeleteLine
// Make a url unavailable. Can only delete owned urls.
// Owner is being calculated via cookie of the requester
func DeleteLine(dbController storage.DatabaseORM) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		owner, ok := ctx.Get("owner")
		if !ok {
			fmt.Println("NO OWNER CONTEXT")
			ctx.Status(http.StatusBadRequest)
			return
		}
		var urls []string
		if err := ctx.BindJSON(&urls); err != nil {
			ctx.Status(http.StatusInternalServerError)
		}

		ctx.Status(http.StatusAccepted)
		go func() {
			err := dbController.Delete(ctx, urls, owner.(string))
			if err != nil {
				ctx.Status(http.StatusInternalServerError)
				return
			}
		}()
	}
}
