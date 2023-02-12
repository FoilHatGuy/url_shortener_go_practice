package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"shortener/internal/cfg"
	"shortener/internal/storage"
)

func GetShortURL(ctx *gin.Context) {

	//fmt.Printf("--------------data: %v\n", storage.Database.GetData())
	inputURL := ctx.Params.ByName("shortURL")
	fmt.Printf("Input url: %q\n\n", inputURL)
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

func PostURL(ctx *gin.Context) {
	buf := make([]byte, 1024)
	num, _ := ctx.Request.Body.Read(buf)
	inputURL := string(buf[0:num])

	_, err := url.Parse(inputURL)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	shortURL := storage.Database.AddURL(inputURL)

	fmt.Printf("Input url: %s\n", inputURL)
	host := cfg.Router.BaseURL
	result := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   shortURL,
	}
	fmt.Printf("Short url: %s\n\n", result.String())

	ctx.String(http.StatusCreated, "%v", result.String())
}

func PostAPIURL(ctx *gin.Context) {
	var newReqBody struct {
		URL string `json:"url"`
	}

	if err := ctx.BindJSON(&newReqBody); err != nil {
		return
	}

	_, err := url.Parse(newReqBody.URL)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	shortURL := storage.Database.AddURL(newReqBody.URL)

	//fmt.Printf("Input url: %s\n", newReqBody.URL)
	//fmt.Printf("Short url: %s\n\n", shortURL)

	result := url.URL{
		Scheme: "http",
		Host:   cfg.Server.Host + ":" + cfg.Server.Port,
		Path:   shortURL,
	}

	newResBody := struct {
		Result string `json:"result"`
	}{result.String()}
	ctx.IndentedJSON(http.StatusCreated, newResBody)
}
