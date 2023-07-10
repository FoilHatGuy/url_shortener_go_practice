package handlers

import (
	"net/http"

	"shortener/internal/cfg"
	"shortener/internal/storage"

	"github.com/gin-gonic/gin"

	utils "shortener/internal/server/handlers/utils"
)

// PostAPIURL
// Handler for batch shortening of urs.
// Takes the field "url" from request body and returns the result in "url" field for accessing original url.
func PostAPIURL(dbController storage.DatabaseORM, config *cfg.ConfigT) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var newReqBody struct {
			URL string `json:"url"`
		}
		owner, ok := ctx.Get("owner")
		if !ok {
			ctx.Status(http.StatusBadRequest)
			return
		}

		if err := ctx.BindJSON(&newReqBody); err != nil {
			return
		}

		result, added, err := utils.Shorten(ctx, dbController, newReqBody.URL, owner.(string), config)
		if err != nil {
			ctx.Status(http.StatusBadRequest)
			return
		}

		newResBody := struct {
			Result string `json:"result"`
		}{result}
		if added {
			ctx.IndentedJSON(http.StatusCreated, newResBody)
		} else {
			ctx.IndentedJSON(http.StatusConflict, newResBody)
		}
	}
}
