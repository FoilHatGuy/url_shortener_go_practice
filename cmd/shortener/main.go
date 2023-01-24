package main

import (
	"fmt"
	"github.com/FoilHatGuy/url_shortener_go_practice/storage"
	"io"
	"log"
	"net/http"
)

const ( //config
	urlLength = 10
	host      = "localhost"
	port      = 8080
)

var database = new(storage.Storage)

func main() {
	server := http.Server{
		Addr: fmt.Sprintf("%s:%d", host, port),
	}

	http.HandleFunc("/", receiveUrl)
	//http.HandleFunc("/:id", provideUrl)
	//http.HandleFunc("/*", provideUrl)

	log.Fatal(server.ListenAndServe())
}

func receiveUrl(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "POST":
		urlBytes, _ := io.ReadAll(request.Body)
		url := string(urlBytes)

		shortUrl := database.AddUrl(url, urlLength)

		fmt.Printf("Input url: %s\n", url)
		fmt.Printf("SHort url: %s\n\n", shortUrl)

		writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := writer.Write([]byte(shortUrl))
		if err != nil {
			return
		}

	case "GET":
		fmt.Printf("Input url: %s\n", request.URL.Fragment)

	}
}
