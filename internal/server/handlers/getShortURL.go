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
	return func(ctx *gin.Context) {
		inputURL := ctx.Params.ByName("shortURL")
		fmt.Printf("Input url: %q\n", inputURL)
		if len(inputURL) != config.Shortener.URLLength {
			ctx.Status(http.StatusBadRequest)
			return
		}

		result, ok, err := dbController.GetURL(ctx, inputURL)
		fmt.Printf("Output url: %s, %t\n", result, err == nil)
		if err != nil {
			ctx.Status(http.StatusBadRequest)
			return
		}
		if result == "" && ok {
			ctx.Status(http.StatusGone)
			return
		}
		ctx.Redirect(307, result)
	}
}
