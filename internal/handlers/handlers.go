package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	url "net/url"
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
	inputURL := string(urlBytes)

	_, err := url.Parse(inputURL)
	if err == nil {
		shortURL := storage.Database.AddURL(inputURL, urlLength)

		fmt.Printf("Input url: %s\n", inputURL)
		fmt.Printf("Short url: %s\n\n", shortURL)

		writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
		result := url.URL{
			Scheme: "http",
			Host:   host + ":" + strconv.FormatInt(port, 10),
			Path:   shortURL,
		}
		writer.WriteHeader(http.StatusCreated)
		_, err := writer.Write([]byte(result.String()))
		if err != nil {
			return
		}
	}
}

func ReceiveURL(writer http.ResponseWriter, request *http.Request) {
	inputURL := chi.URLParam(request, "shortURL")
	fmt.Printf("Input url: %s\n", inputURL)
	if len(inputURL) == urlLength {
		result, err := storage.Database.GetURL(inputURL)
		fmt.Printf("Output url: %s, %t\n", result, err == nil)
		if err == nil {
			fmt.Printf("get complete\n\n")
			writer.Header().Set("Location", result)
			writer.WriteHeader(307)
			_, _ = writer.Write([]byte(result))
		} else {
			writer.WriteHeader(http.StatusBadRequest)
		}
	} else {
		writer.WriteHeader(http.StatusBadRequest)
	}

}
