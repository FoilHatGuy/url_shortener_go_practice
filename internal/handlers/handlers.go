package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"shortener/internal/storage"
	"strconv"
)

const ( //config
	urlLength = 10
	host      = "localhost"
	port      = 8080
)

func GetShortURL(ctx *gin.Context) {

	//fmt.Printf("--------------data: %v\n", storage.Database.GetData())
	inputURL := ctx.Params.ByName("shortURL")
	fmt.Printf("Input url: %q\n\n", inputURL)
	if len(inputURL) != urlLength {
		ctx.Status(http.StatusBadRequest)
		return
	}

	result, err := storage.Database.GetURL(inputURL)
	fmt.Printf("Output url: %s, %t\n", result, err == nil)

	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	fmt.Printf("get complete\n\n")
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
	fmt.Printf("Short url: %s\n\n", shortURL)

	result := url.URL{
		Scheme: "http",
		Host:   host + ":" + strconv.FormatInt(port, 10),
		Path:   shortURL,
	}
	ctx.String(http.StatusCreated, "%v", result.String())
}

func PostApiURL(ctx *gin.Context) {
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

	fmt.Printf("Input url: %s\n", newReqBody.URL)
	fmt.Printf("Short url: %s\n\n", shortURL)

	result := url.URL{
		Scheme: "http",
		Host:   host + ":" + strconv.FormatInt(port, 10),
		Path:   shortURL,
	}

	newResBody := struct {
		Result string `json:"result"`
	}{result.String()}
	ctx.IndentedJSON(http.StatusCreated, newResBody)
}
