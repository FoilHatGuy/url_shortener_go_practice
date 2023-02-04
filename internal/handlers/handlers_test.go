package handlers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
		//{
		//	name:    "Get req",
		//	method:  "GET",
		//	body:    "",
		//	handler: ReceiveURL,
		//	target:  "/XVlBzgbaiC",
		//	want: want{
		//		code:        200,
		//		response:    "",
		//		contentType: "",
		//	},
		//},
		{
			name:    "no such url",
			method:  "GET",
			body:    "",
			handler: ReceiveURL,
			target:  "/nosuchurl_",
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
			target:  "/urltoolongtobevalid",
			want: want{
				code:        400,
				response:    "",
				contentType: "",
			},
		},
	}
	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {
			body := bytes.NewReader([]byte(tt.body))
			request := httptest.NewRequest(tt.method, tt.target, body)

			// создаём новый Recorder
			w := httptest.NewRecorder()
			// определяем хендлер
			h := http.HandlerFunc(tt.handler)
			// запускаем сервер
			h.ServeHTTP(w, request)
			res := w.Result()

			// проверяем код ответа
			if res.StatusCode != tt.want.code {
				t.Errorf("Expected status code %d, got %d", tt.want.code, w.Code)
			}

			// получаем и проверяем тело запроса
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			if len(string(resBody)) != len(tt.want.response) {
				t.Errorf("Expected body %s, got %s", tt.want.response, w.Body.String())
			}

			// заголовок ответа
			if res.Header.Get("Content-Type") != tt.want.contentType {
				t.Errorf("Expected Content-Type %s, got %s", tt.want.contentType, res.Header.Get("Content-Type"))
			}
			defer res.Body.Close()
		})
	}
}
