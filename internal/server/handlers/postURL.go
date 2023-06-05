package handlers

import (
	"io"
	"net/http"

	"shortener/internal/cfg"
	"shortener/internal/storage"

	"github.com/gin-gonic/gin"

	utils "shortener/internal/server/handlers/utils"
)

// PostURL
// Handler for batch shortening of urs.
// Takes txt representation of request body and returns txt of url for accessing it.
func PostURL(dbController storage.DatabaseORM, config *cfg.ConfigT) gin.HandlerFunc {
	return func(c *gin.Context) {
		buf, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		inputURL := string(buf)
		owner, ok := c.Get("owner")
		if !ok {
			c.Status(http.StatusBadRequest)
			return
		}

		result, added, err := utils.Shorten(c, dbController, inputURL, owner.(string), config)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		if added {
			c.String(http.StatusCreated, "%v", result)
		} else {
			c.String(http.StatusConflict, "%v", result)
		}
	}
}
