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
func DeleteLine(c *gin.Context) {
	owner, ok := c.Get("owner")
	if !ok {
		fmt.Println("NO OWNER CONTEXT")
		c.Status(http.StatusBadRequest)
		return
	}
	var urls []string
	if err := c.BindJSON(&urls); err != nil {
		c.Status(http.StatusInternalServerError)
	}

	c.Status(http.StatusAccepted)
	go func() {
		err := storage.Controller.Delete(c, urls, owner.(string))
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
	}()
}
