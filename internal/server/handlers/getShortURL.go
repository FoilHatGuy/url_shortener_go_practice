package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"shortener/internal/cfg"
	"shortener/internal/storage"
)

var config *cfg.ConfigT

// GetShortURL
// Get the original url by the short url
func GetShortURL(c *gin.Context) {
	if config == nil {
		config = cfg.Initialize()
	}
	inputURL := c.Params.ByName("shortURL")
	fmt.Printf("Input url: %q\n", inputURL)
	if len(inputURL) != config.Shortener.URLLength {
		c.Status(http.StatusBadRequest)
		return
	}

	result, ok, err := storage.Controller.GetURL(c, inputURL)
	fmt.Printf("Output url: %s, %t\n", result, err == nil)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	if result == "" && ok {
		c.Status(http.StatusGone)
		return
	}
	// fmt.Printf("get complete\n\n")
	c.Redirect(307, result)
}
