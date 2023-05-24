package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	utils "shortener/internal/server/handlers/utils"
)

// PostAPIURL
// Handler for batch shortening of urs.
// Takes the field "url" from request body and returns the result in "url" field for accessing original url.
func PostAPIURL(c *gin.Context) {
	var newReqBody struct {
		URL string `json:"url"`
	}
	owner, ok := c.Get("owner")
	if !ok {
		c.Status(http.StatusBadRequest)
		return
	}

	if err := c.BindJSON(&newReqBody); err != nil {
		return
	}

	result, added, err := utils.Shorten(newReqBody.URL, owner.(string), c)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	newResBody := struct {
		Result string `json:"result"`
	}{result}
	if added {
		c.IndentedJSON(http.StatusCreated, newResBody)
	} else {
		c.IndentedJSON(http.StatusConflict, newResBody)
	}
}
