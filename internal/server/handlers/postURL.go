package handlers

import (
	"io"
	"net/http"

	"shortener/internal/cfg"
	"shortener/internal/storage"

	"github.com/gin-gonic/gin"

	"shortener/internal/server/handlers/utils"
)

// PostURL
// Handler for batch shortening of urs.
// Takes txt representation of request body and returns txt of url for accessing it.
func PostURL(dbController storage.DatabaseORM, config *cfg.ConfigT) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		buf, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			ctx.Status(http.StatusBadRequest)
			return
		}
		inputURL := string(buf)
		owner, ok := ctx.Get("owner")
		if !ok {
			ctx.Status(http.StatusBadRequest)
			return
		}

		result, added, err := utils.Shorten(ctx, dbController, inputURL, owner.(string), config)
		if err != nil {
			ctx.Status(http.StatusBadRequest)
			return
		}

		if added {
			ctx.String(http.StatusCreated, "%v", result)
		} else {
			ctx.String(http.StatusConflict, "%v", result)
		}
	}
}
