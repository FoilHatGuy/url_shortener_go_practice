package handlers

import (
	"net/http"

	"shortener/internal/cfg"
	"shortener/internal/storage"

	"github.com/gin-gonic/gin"

	"shortener/internal/server/handlers/utils"
)

// BatchShorten
// Handler for batch shortening of urs.
// Accepts array of JSON objects containing:
// {"correlation_id": "id of url", "original_url": "url to be shortened"}
// returns array of following jsons:
// {"correlation_id": "id of url", "short_url": "url that was shortened"}
func BatchShorten(dbController storage.DatabaseORM, config *cfg.ConfigT) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		type reqElement struct {
			LineID string `json:"correlation_id"`
			URL    string `json:"original_url"`
		}
		type resElement struct {
			LineID string `json:"correlation_id"`
			URL    string `json:"short_url"`
		}
		var newReqBody []*reqElement
		var newResBody []*resElement
		owner, ok := ctx.Get("owner")
		if !ok {
			ctx.Status(http.StatusBadRequest)
			return
		}

		if err := ctx.BindJSON(&newReqBody); err != nil {
			return
		}

		for _, element := range newReqBody {
			result, _, err := utils.Shorten(ctx, dbController, element.URL, owner.(string), config)
			if err != nil {
				ctx.Status(http.StatusBadRequest)
				return
			}
			newResBody = append(newResBody, &resElement{LineID: element.LineID, URL: result})
		}

		ctx.IndentedJSON(http.StatusCreated, newResBody)
	}
}
