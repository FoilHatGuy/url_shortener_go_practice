package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/url"
	"shortener/internal/cfg"
	"shortener/internal/storage"
)

func getShortURL(c *gin.Context) {

	//fmt.Printf("--------------data: %v\n", storage.Controller.GetData())
	inputURL := c.Params.ByName("shortURL")
	fmt.Printf("Input url: %q\n", inputURL)
	if len(inputURL) != cfg.Shortener.URLLength {
		c.Status(http.StatusBadRequest)
		return
	}

	result, ok, err := storage.Controller.GetURL(c, inputURL)
	fmt.Printf("Output url: %s, %t\n", result, err == nil)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	if !ok {
		c.Status(http.StatusGone)
		return
	}
	//fmt.Printf("get complete\n\n")
	c.Redirect(307, result)

}

func postURL(c *gin.Context) {
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

	result, added, err := shorten(inputURL, owner.(string), c)
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

func postAPIURL(c *gin.Context) {
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

	result, added, err := shorten(newReqBody.URL, owner.(string), c)
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

func pingDatabase(c *gin.Context) {
	ping := storage.Controller.Ping(c)
	if ping {
		c.Status(http.StatusOK)
	} else {
		c.Status(http.StatusInternalServerError)
	}

}
func delete(c *gin.Context) {
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

	go func() {
		err := storage.Controller.Delete(c, urls, owner.(string))
		if err != nil {
			return
		}
	}()
	c.Status(http.StatusAccepted)

}

func getAllOwnedURL(c *gin.Context) {
	owner, ok := c.Get("owner")
	if !ok {
		fmt.Println("NO OWNER CONTEXT")
		c.Status(http.StatusBadRequest)
		return
	}

	result, err := storage.Controller.GetURLByOwner(c, owner.(string))
	if err != nil {
		fmt.Println("ERROR WHILE GETTING DATA FROM DB")
		c.Status(http.StatusBadRequest)
		return
	}
	if result != nil {
		c.IndentedJSON(http.StatusOK, result)
	} else {
		c.Status(http.StatusNoContent)
	}
}

func shorten(inputURL string, owner string, ctx context.Context) (string, bool, error) {

	_, err := url.Parse(inputURL)
	if err != nil {
		return "", false, errors.New("bad URL")
	}
	shortURL, added, err := storage.Controller.AddURL(ctx, inputURL, owner)
	if err != nil {
		return "", added, err
	}

	fmt.Printf("Input url: %s\n", inputURL)
	fmt.Printf("Short url: %s\n\n", shortURL)

	result, _ := url.Parse(cfg.Server.BaseURL)
	result = result.JoinPath(shortURL)
	return result.String(), added, nil
}

func batchShorten(c *gin.Context) {
	type reqElement struct {
		LineID string `json:"correlation_id"`
		URL    string `json:"original_url"`
	}
	type resElement struct {
		LineID string `json:"correlation_id"`
		URL    string `json:"short_url"`
	}
	var newReqBody []reqElement
	var newResBody []resElement
	owner, ok := c.Get("owner")
	if !ok {
		c.Status(http.StatusBadRequest)
		return
	}

	if err := c.BindJSON(&newReqBody); err != nil {
		return
	}

	for _, element := range newReqBody {
		result, _, err := shorten(element.URL, owner.(string), c)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		newResBody = append(newResBody, resElement{element.LineID, result})
	}

	c.IndentedJSON(http.StatusCreated, newResBody)
}
