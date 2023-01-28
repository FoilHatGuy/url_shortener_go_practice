package handlers

import (
	"fmt"
	"io"
	"net/http"
	"shortener/internal/storage"
	"strconv"
	"strings"
)

const ( //config
	urlLength = 10
	host      = "localhost"
	port      = 8080
)

func ReceiveURL(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "POST":
		urlBytes, _ := io.ReadAll(request.Body)
		url := string(urlBytes)

		shortURL := storage.Database.AddURL(url, urlLength)

		fmt.Printf("Input url: %s\n", url)
		fmt.Printf("Short url: %s\n\n", shortURL)

		writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
		result := "http://" + host + ":" + strconv.FormatInt(port, 10) + "/" + shortURL
		writer.WriteHeader(201)
		_, err := writer.Write([]byte(result))
		if err != nil {
			return
		}

	case "GET":
		urlArray := strings.Split(request.URL.Path[1:], "/")
		fmt.Printf("Input url: %s\n", urlArray[0])
		if len(urlArray) == 1 && len(urlArray[0]) == urlLength {
			url, err := storage.Database.GetURL(urlArray[0])
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
