package handlers

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"shortener/internal/cfg"
	"shortener/internal/storage"
	"testing"
	"time"
)

func Run() {
	r := gin.Default()
	baseRouter := r.Group("")
	{
		baseRouter.Use(ArchiveData())
		baseRouter.GET("/:shortURL", GetShortURL)
		baseRouter.POST("/", PostURL)
		api := baseRouter.Group("/api")
		{
			api.POST("/shorten", PostAPIURL)
		}
	}

	t := time.NewTicker(time.Duration(cfg.Storage.AutosaveInterval) * time.Second)
	storage.Database.LoadData()
	go func() {
		for range t.C {
			//fmt.Print("AUTOSAVE\n")
			storage.Database.SaveData()
		}
	}()
	fmt.Println("SERVER LISTENING ON", cfg.Server.Address)
	log.Fatal(r.Run(cfg.Server.Address))
}

func TestReceiveURL(t *testing.T) {
	cfg.Initialize()
	type want struct {
		acceptType  string
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name   string
		method string
		body   string
		target string
		want   want
	}{
		{
			name:   "Post req",
			method: "POST",
			body:   "http://a30ac6lti.biz/fc6pql9n/duut2ohnkaja",
			target: "http://localhost:8080/",
			want: want{
				acceptType:  "text/plain",
				code:        201,
				response:    "http://localhost:8080/XVlBzgbaiC",
				contentType: "text/plain; charset=utf-8",
			},
		},
		// TODO: to complete this autotoest data should be stored on drive
		{
			name:   "Get req",
			method: "GET",
			body:   "",
			target: "http://localhost:8080/XVlBzgbaiC",
			want: want{
				acceptType:  "text/plain",
				code:        200,
				response:    "",
				contentType: "",
			},
		},
		{
			name:   "no such url",
			method: "GET",
			body:   "",
			target: "http://localhost:8080/nosuchurl_",
			want: want{
				acceptType:  "text/plain",
				code:        400,
				response:    "",
				contentType: "",
			},
		},
		{
			name:   "url too long to be valid",
			method: "GET",
			body:   "",
			target: "http://localhost:8080/urltoolongtobevalid",
			want: want{
				acceptType:  "text/plain",
				code:        400,
				response:    "",
				contentType: "",
			},
		},
	}
	go Run()
	client := &http.Client{}
	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {
			var res *http.Response
			if tt.method == "GET" {
				var err error
				//res, err = http.Get(tt.target)
				body := bytes.NewReader([]byte(tt.body))
				r, _ := http.NewRequest("GET", tt.target, body)
				r.Header.Add("Accept-Encoding", tt.want.acceptType)
				res, err = client.Do(r)
				if err != nil {
					return
				}
				defer res.Body.Close()
			} else if tt.method == "POST" {
				var err error
				body := bytes.NewReader([]byte(tt.body))
				r, _ := http.NewRequest("POST", tt.target, body)
				r.Header.Add("Accept-Encoding", tt.want.acceptType)
				//res, err = http.Post(tt.target, "text/plain; charset=utf-8", body)
				res, err = client.Do(r)
				if err != nil {
					panic(err)
				}
				if err != nil {
					return
				}
				defer res.Body.Close()
			}
			if res.StatusCode != tt.want.code {
				t.Errorf("Expected status code %d, got %d", tt.want.code, res.StatusCode)
			}

			// получаем и проверяем тело запроса
			resBody, err := io.ReadAll(res.Body)
			fmt.Print(string(resBody))
			if err != nil {
				t.Fatal(err)
			}
			if string(resBody) != tt.want.response {
				t.Errorf("Expected body %s, got %s", tt.want.response, string(resBody))
			}

			// заголовок ответа
			if res.Header.Get("Content-Type") != tt.want.contentType {
				t.Errorf("Expected Content-Type %s, got %s", tt.want.contentType, res.Header.Get("Content-Type"))
			}
		})
	}
}
