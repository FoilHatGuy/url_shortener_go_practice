package server

import (
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

	result, err := storage.Controller.GetURL(inputURL)
	fmt.Printf("Output url: %s, %t\n", result, err == nil)

	if err != nil {
		c.Status(http.StatusBadRequest)
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

	result, err := shorten(inputURL, owner.(string))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	c.String(http.StatusCreated, "%v", result)
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

	result, err := shorten(newReqBody.URL, owner.(string))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	newResBody := struct {
		Result string `json:"result"`
	}{result}
	c.IndentedJSON(http.StatusCreated, newResBody)
}

func pingDatabase(c *gin.Context) {
	ping := storage.Controller.Ping()
	if ping {
		c.Status(http.StatusOK)
	} else {
		c.Status(http.StatusInternalServerError)
	}

}

func getAllOwnedURL(c *gin.Context) {
	owner, ok := c.Get("owner")
	if !ok {
		fmt.Println("NO OWNER CONTEXT")
		c.Status(http.StatusBadRequest)
		return
	}

	result, err := storage.Controller.GetURLByOwner(owner.(string))
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

func shorten(inputURL string, owner string) (string, error) {

	_, err := url.Parse(inputURL)
	if err != nil {
		return "", errors.New("bad URL")
	}
	shortURL, err := storage.Controller.AddURL(inputURL, owner)
	if err != nil {
		return "", err
	}

	fmt.Printf("Input url: %s\n", inputURL)
	fmt.Printf("Short url: %s\n\n", shortURL)

	result, _ := url.Parse(cfg.Server.BaseURL)
	result = result.JoinPath(shortURL)
	return result.String(), nil
}
