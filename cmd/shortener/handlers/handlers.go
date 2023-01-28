package handlers

import (
	"fmt"
	"github.com/FoilHatGuy/url_shortener_go_practice/storage"
	"io"
	"net/http"
	"strings"
)

const ( //config
	urlLength = 10
)

func ReceiveUrl(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "POST":
		urlBytes, _ := io.ReadAll(request.Body)
		url := string(urlBytes)

		shortUrl := storage.Database.AddUrl(url, urlLength)

		fmt.Printf("Input url: %s\n", url)
		fmt.Printf("Short url: %s\n\n", shortUrl)

		writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := writer.Write([]byte(shortUrl))
		if err != nil {
			return
		}

	case "GET":
		urlArray := strings.Split(request.URL.Path[1:], "/")
		fmt.Printf("Input url: %s\n", urlArray[0])
		if len(urlArray) == 1 && len(urlArray[0]) == urlLength {
			url, err := storage.Database.GetUrl(urlArray[0])
			fmt.Printf("Output url: %s, %t\n", url, err == nil)
			if err == nil {
				fmt.Printf("get complete\n\n")
				writer.Header().Set("Location", url)
				writer.WriteHeader(http.StatusTemporaryRedirect)
			} else {
				writer.WriteHeader(http.StatusBadRequest)
			}
		}
	}
}
