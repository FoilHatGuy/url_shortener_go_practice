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
	return func(c *gin.Context) {
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
		owner, ok := c.Get("owner")
		if !ok {
			c.Status(http.StatusBadRequest)
			return
		}

		if err := c.BindJSON(&newReqBody); err != nil {
			return
		}

		for _, element := range newReqBody {
			result, _, err := utils.Shorten(c, dbController, element.URL, owner.(string), config)
			if err != nil {
				c.Status(http.StatusBadRequest)
				return
			}
			newResBody = append(newResBody, &resElement{LineID: element.LineID, URL: result})
		}

		c.IndentedJSON(http.StatusCreated, newResBody)
	}
}
