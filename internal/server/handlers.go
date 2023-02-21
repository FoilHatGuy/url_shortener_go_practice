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

func getShortURL(ctx *gin.Context) {

	//fmt.Printf("--------------data: %v\n", storage.Database.GetData())
	inputURL := ctx.Params.ByName("shortURL")
	fmt.Printf("Input url: %q\n", inputURL)
	if len(inputURL) != cfg.Shortener.URLLength {
		ctx.Status(http.StatusBadRequest)
		return
	}

	result, err := storage.Database.GetURL(inputURL)
	fmt.Printf("Output url: %s, %t\n", result, err == nil)

	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	//fmt.Printf("get complete\n\n")
	ctx.Redirect(307, result)

}

func postURL(ctx *gin.Context) {
	buf, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	inputURL := string(buf)

	result, err := shorten(inputURL)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	ctx.String(http.StatusCreated, "%v", result)
}

func postAPIURL(ctx *gin.Context) {
	var newReqBody struct {
		URL string `json:"url"`
	}

	if err := ctx.BindJSON(&newReqBody); err != nil {
		return
	}

	result, err := shorten(newReqBody.URL)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
	}

	newResBody := struct {
		Result string `json:"result"`
	}{result}
	ctx.IndentedJSON(http.StatusCreated, newResBody)
}

func shorten(inputURL string) (string, error) {

	_, err := url.Parse(inputURL)
	if err != nil {
		return "", errors.New("bad URL")
	}
	shortURL := storage.Database.AddURL(inputURL)

	fmt.Printf("Input url: %s\n", inputURL)
	fmt.Printf("Short url: %s\n\n", shortURL)

	result, _ := url.Parse(cfg.Server.BaseURL)
	result = result.JoinPath(shortURL)
	return result.String(), nil
}
