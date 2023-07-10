package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"shortener/internal/storage"
)

// GetAllOwnedURL
// Get all owned urls.
// Owner is being calculated via cookie of the requester
func GetAllOwnedURL(dbController storage.DatabaseORM) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		owner, ok := ctx.Get("owner")
		if !ok {
			fmt.Println("NO OWNER CONTEXT")
			ctx.Status(http.StatusBadRequest)
			return
		}
		result, err := dbController.GetURLByOwner(ctx, owner.(string))
		if err != nil {
			fmt.Println("ERROR WHILE GETTING DATA FROM DB")
			ctx.Status(http.StatusBadRequest)
			return
		}
		if result != nil {
			ctx.IndentedJSON(http.StatusOK, result)
		} else {
			ctx.Status(http.StatusNoContent)
		}
	}
}
