package handlers

import (
	"bytes"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"strconv"
	"testing"
)

func serve() {
	r := chi.NewRouter()

	r.Post("/", SendURL)
	r.Get("/{shortURL:[a-zA-Z]{"+strconv.FormatInt(urlLength, 10)+"}}", ReceiveURL)

	http.ListenAndServe("localhost:8080", r)
}

func TestReceiveURL(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name    string
		method  string
		body    string
		handler func(http.ResponseWriter, *http.Request)
		target  string
		want    want
	}{
		{
			name:    "Post req",
			method:  "POST",
			body:    "http://google.com",
			handler: SendURL,
			target:  "/",
			want: want{
				code:        201,
				response:    "http://localhost:8080/XVlBzgbaiC",
				contentType: "text/plain; charset=utf-8",
			},
		},
		// TODO: to complete this autotoest data should be stored on drive
		{
			name:    "Get req",
			method:  "GET",
			body:    "",
			handler: ReceiveURL,
			target:  "http://localhost:8080/XVlBzgbaiC",
			want: want{
				code:        200,
				response:    "",
				contentType: "",
			},
		},
		{
			name:    "no such url",
			method:  "GET",
			body:    "",
			handler: ReceiveURL,
			target:  "http://localhost:8080/nosuchurl_",
			want: want{
				code:        400,
				response:    "",
				contentType: "",
			},
		},
		{
			name:    "url too long to be valid",
			method:  "GET",
			body:    "",
			handler: ReceiveURL,
			target:  "http://localhost:8080/urltoolongtobevalid",
			want: want{
				code:        400,
				response:    "",
				contentType: "",
			},
		},
	}
	go serve()
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
				res, err = http.Post(tt.target, "text/plain", body)
				if err != nil {
					return
				}

			}
			//request := httptest.NewRequest(tt.method, tt.target, body)

			//w := httptest.NewRecorder()
			//res := w.Result()
			//res, err := client.Do(request)
			// проверяем код ответа
			if res.StatusCode != tt.want.code {
				t.Errorf("Expected status code %d, got %d", tt.want.code, res.StatusCode)
			}

			// получаем и проверяем тело запроса
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			if len(string(resBody)) != len(tt.want.response) {
				t.Errorf("Expected body %s, got %s", tt.want.response, string(resBody))
			}

			// заголовок ответа
			if res.Header.Get("Content-Type") != tt.want.contentType {
				t.Errorf("Expected Content-Type %s, got %s", tt.want.contentType, res.Header.Get("Content-Type"))
			}
			defer res.Body.Close()
		})
	}
}
