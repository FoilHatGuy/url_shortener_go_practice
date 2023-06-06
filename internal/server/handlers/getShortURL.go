package handlers

import (
	"fmt"
	"net/http"

	"shortener/internal/cfg"

	"github.com/gin-gonic/gin"

	"shortener/internal/storage"
)

// GetShortURL
// Get the original url by the short url
func GetShortURL(dbController storage.DatabaseORM, config *cfg.ConfigT) gin.HandlerFunc {
	return func(c *gin.Context) {
		inputURL := c.Params.ByName("shortURL")
		fmt.Printf("Input url: %q\n", inputURL)
		if len(inputURL) != config.Shortener.URLLength {
			c.Status(http.StatusBadRequest)
			return
		}

		result, ok, err := dbController.GetURL(c, inputURL)
		fmt.Printf("Output url: %s, %t\n", result, err == nil)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		if result == "" && ok {
			c.Status(http.StatusGone)
			return
		}
		c.Redirect(307, result)
	}
}
