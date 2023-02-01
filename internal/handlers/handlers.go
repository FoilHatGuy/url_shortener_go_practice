package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"shortener/internal/storage"
	"strconv"
)

const ( //config
	urlLength = 10
	host      = "localhost"
	port      = 8080
)

func SendURL(writer http.ResponseWriter, request *http.Request) {
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

}

func ReceiveURL(writer http.ResponseWriter, request *http.Request) {
	url := chi.URLParam(request, "shortURL")
	fmt.Printf("Input url: %s\n", url)
	if len(url) == urlLength {
		url, err := storage.Database.GetURL(url)
		fmt.Printf("Output url: %s, %t\n", url, err == nil)
		if err == nil {
			fmt.Printf("get complete\n\n")
			writer.Header().Set("Location", url)
			writer.WriteHeader(307)
			_, _ = writer.Write([]byte(url))
		} else {
			writer.WriteHeader(http.StatusBadRequest)
		}
	} else {
		writer.WriteHeader(400)
	}

}
