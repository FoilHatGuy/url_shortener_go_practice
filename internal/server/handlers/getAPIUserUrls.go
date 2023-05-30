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
func GetAllOwnedURL(c *gin.Context) {
	owner, ok := c.Get("owner")
	if !ok {
		fmt.Println("NO OWNER CONTEXT")
		c.Status(http.StatusBadRequest)
		return
	}

	result, err := storage.Controller.GetURLByOwner(c, owner.(string))
	if err != nil {
		fmt.Println("ERROR WHILE GETTING DATA FROM DB")
		c.Status(http.StatusBadRequest)
		return
	}
	if result != nil {
		c.IndentedJSON(http.StatusOK, result)
	} else {
		c.Status(http.StatusNoContent)
	}
}
