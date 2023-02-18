package server

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"shortener/internal/cfg"
	"strings"
	"testing"
	"time"
)

func TestReceiveURL(t *testing.T) {
	cfg.Initialize()
	type want struct {
		acceptType  string
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		encoding string
		name     string
		method   string
		body     string
		target   string
		want     want
	}{
		{
			name:     "Post req",
			method:   "POST",
			body:     "http://a30ac6lti.biz/fc6pql9n/duut2ohnkaja",
			target:   "http://localhost:8080/",
			encoding: "none",
			want: want{
				acceptType:  "text/plain",
				code:        201,
				response:    "http://localhost:8080/XVlBzgbaiC",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:     "Post API req",
			method:   "POST",
			body:     "{\"url\":\"http://google.com\"}",
			target:   "http://localhost:8080/api/shorten",
			encoding: "none",
			want: want{
				acceptType:  "text/plain",
				code:        201,
				response:    "{\n    \"result\": \"http://localhost:8080/XVlBzgbaiC\"\n}",
				contentType: "application/json; charset=utf-8",
			},
		},
		{
			name:     "Get from API",
			method:   "GET",
			body:     "",
			target:   "http://localhost:8080/XVlBzgbaiC",
			encoding: "none",
			want: want{
				acceptType:  "text/plain",
				code:        200,
				response:    "",
				contentType: "",
			},
		},
		{
			name:     "Get req",
			method:   "GET",
			body:     "",
			target:   "http://localhost:8080/XVlBzgbaiC",
			encoding: "none",
			want: want{
				acceptType:  "text/plain",
				code:        200,
				response:    "",
				contentType: "",
			},
		},
		{
			name:     "no such url",
			method:   "GET",
			body:     "",
			target:   "http://localhost:8080/nosuchurl_",
			encoding: "none",
			want: want{
				acceptType:  "text/plain",
				code:        400,
				response:    "",
				contentType: "",
			},
		},
		{
			name:     "url too long to be valid",
			method:   "GET",
			body:     "",
			target:   "http://localhost:8080/urltoolongtobevalid",
			encoding: "none",
			want: want{
				acceptType:  "text/plain",
				code:        400,
				response:    "",
				contentType: "",
			},
		},
	}
	go Run()
	client := &http.Client{
		CheckRedirect: noRedirect,
	}
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
				if tt.encoding == "gzip" {
					body := bytes.NewBuffer([]byte{})
					gzipR := gzip.NewWriter(body)
					fmt.Printf("%x\n", []byte(tt.body))
					_, err = gzipR.Write([]byte(tt.body))
					if err != nil {
						return
					}
					fmt.Printf("%x\n", body.Bytes())
					defer gzipR.Close()
				}
				body := bytes.NewReader([]byte(tt.body))
				r, _ := http.NewRequest("POST", tt.target, body)
				if tt.encoding == "gzip" {
					r.Header.Add("Content-Encoding", tt.encoding)
				}
				r.Header.Add("Accept-Encoding", tt.want.acceptType)
				//res, err = http.Post(tt.target, "text/plain; charset=utf-8", body)
				res, err = client.Do(r)
				if err != nil {
					return
				}
				defer res.Body.Close()
			}
			if res.StatusCode != tt.want.code {
				t.Errorf("Expected status code %d, got %d", tt.want.code, res.StatusCode)
			}

			// получаем и проверяем тело запроса
			var resBody []byte
			contentType := res.Header.Get("Content-Encoding")
			if !strings.Contains(contentType, "gzip") {
				fmt.Println("Reading body")
				resBody, _ = io.ReadAll(res.Body)
			} else {
				fmt.Println("Unpacking body")
				gzipR, err := gzip.NewReader(res.Body)
				if err != nil {
					return
				}
				defer gzipR.Close()
				resBody, _ = io.ReadAll(gzipR)
			}

			fmt.Printf("RECEIVED\nBODY:\t%s\nSTATUS:\t%v\n", string(resBody), res.StatusCode)

			if string(resBody) != tt.want.response {
				t.Errorf("Expected body %s, got %s", tt.want.response, string(resBody))
			}

			// заголовок ответа
			if res.Header.Get("Content-Type") != tt.want.contentType {
				t.Errorf("Expected Content-Type %s, got %s", tt.want.contentType, res.Header.Get("Content-Type"))
			}

		})
		time.Sleep(1 * time.Second)
	}
}

func noRedirect(req *http.Request, via []*http.Request) error {
	return errors.New("Don't redirect!")
}