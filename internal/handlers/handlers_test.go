package handlers

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"testing"
)

func Run() {
	r := gin.Default()

	r.GET("/:shortURL", GetShortURL)
	r.POST("/", PostURL)
	api := r.Group("/api")
	{
		api.POST("/shorten", PostApiURL)
	}
	log.Fatal(r.Run())
}

func TestReceiveURL(t *testing.T) {
	type want struct {
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
				code:        400,
				response:    "",
				contentType: "",
			},
		},
	}
	go Run()
	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {
			var res *http.Response
			if tt.method == "GET" {
				var err error
				res, err = http.Get(tt.target)
				if err != nil {
					return
				}
			} else if tt.method == "POST" {
				var err error
				body := bytes.NewReader([]byte(tt.body))
				res, err = http.Post(tt.target, "text/plain; charset=utf-8", body)
				if err != nil {
					return
				}

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
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {

				}
			}(res.Body)
		})
	}
}
