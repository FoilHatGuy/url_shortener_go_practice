package main

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/suite"
	"io"
	"net/http"
	"shortener/internal/cfg"
	"shortener/internal/server"
	"shortener/internal/storage"
	"testing"
	"time"
)

//
//func TestReceiveURL(t *testing.T) {
//	type want struct {
//		acceptType  string
//		code        int
//		response    string
//		contentType string
//	}
//	tests := []struct {
//		encoding string
//		name     string
//		method   string
//		body     string
//		target   string
//		want     want
//	}{
//		{
//			name:     "Post req",
//			method:   "POST",
//			body:     "http://a30ac6lti.biz/fc6pql9n/duut2ohnkaja",
//			target:   "http://localhost:8080/",
//			encoding: "none",
//			want: want{
//				acceptType:  "text/plain",
//				code:        201,
//				response:    "http://localhost:8080/XVlBzgbaiC",
//				contentType: "text/plain; charset=utf-8",
//			},
//		},
//		{
//			name:     "Post API req",
//			method:   "POST",
//			body:     "{\"url\":\"http://google.com\"}",
//			target:   "http://localhost:8080/api/shorten",
//			encoding: "none",
//			want: want{
//				acceptType:  "text/plain",
//				code:        201,
//				response:    "{\n    \"result\": \"http://localhost:8080/XVlBzgbaiC\"\n}",
//				contentType: "application/json; charset=utf-8",
//			},
//		},
//		{
//			name:     "Get from API",
//			method:   "GET",
//			body:     "",
//			target:   "http://localhost:8080/XVlBzgbaiC",
//			encoding: "none",
//			want: want{
//				acceptType:  "text/plain",
//				code:        200,
//				response:    "",
//				contentType: "",
//			},
//		},
//		//{
//		//	name:     "Get req",
//		//	method:   "GET",
//		//	body:     "",
//		//	target:   "http://localhost:8080/MRAjWwhTHc",
//		//	encoding: "none",
//		//	want: want{
//		//		acceptType:  "text/plain",
//		//		code:        200,
//		//		response:    "",
//		//		contentType: "",
//		//	},
//		//},
//		{
//			name:     "no such url",
//			method:   "GET",
//			body:     "",
//			target:   "http://localhost:8080/nosuchurl_",
//			encoding: "none",
//			want: want{
//				acceptType:  "text/plain",
//				code:        400,
//				response:    "",
//				contentType: "",
//			},
//		},
//		{
//			name:     "url too long to be valid",
//			method:   "GET",
//			body:     "",
//			target:   "http://localhost:8080/urltoolongtobevalid",
//			encoding: "none",
//			want: want{
//				acceptType:  "text/plain",
//				code:        400,
//				response:    "",
//				contentType: "",
//			},
//		},
//	}
//	cfg.Initialize()
//	go Run()
//	client := &http.Client{
//		CheckRedirect: noRedirect,
//	}
//	for _, tt := range tests {
//		// запускаем каждый тест
//		t.Run(tt.name, func(t *testing.T) {
//			var res *http.Response
//			if tt.method == "GET" {
//				var err error
//				//res, err = http.Get(tt.target)
//				body := bytes.NewReader([]byte(tt.body))
//				r, _ := http.NewRequest("GET", tt.target, body)
//				r.Header.Add("Accept-Encoding", tt.want.acceptType)
//				res, err = client.Do(r)
//				if err != nil {
//					return
//				}
//				defer res.Body.Close()
//
//			} else if tt.method == "POST" {
//				var err error
//				if tt.encoding == "gzip" {
//					body := bytes.NewBuffer([]byte{})
//					gzipR := gzip.NewWriter(body)
//					fmt.Printf("%x\n", []byte(tt.body))
//					_, err = gzipR.Write([]byte(tt.body))
//					if err != nil {
//						return
//					}
//					fmt.Printf("%x\n", body.Bytes())
//					defer gzipR.Close()
//				}
//				body := bytes.NewReader([]byte(tt.body))
//				r, _ := http.NewRequest("POST", tt.target, body)
//				if tt.encoding == "gzip" {
//					r.Header.Add("Content-Encoding", tt.encoding)
//				}
//				r.Header.Add("Accept-Encoding", tt.want.acceptType)
//				//res, err = http.Post(tt.target, "text/plain; charset=utf-8", body)
//				res, err = client.Do(r)
//				if err != nil {
//					return
//				}
//				defer res.Body.Close()
//			}
//			if res.StatusCode != tt.want.code {
//				t.Errorf("Expected status code %d, got %d", tt.want.code, res.StatusCode)
//			}
//
//			// получаем и проверяем тело запроса
//			var resBody []byte
//			contentType := res.Header.Get("Content-Encoding")
//			if !strings.Contains(contentType, "gzip") {
//				fmt.Println("Reading body")
//				resBody, _ = io.ReadAll(res.Body)
//			} else {
//				fmt.Println("Unpacking body")
//				gzipR, err := gzip.NewReader(res.Body)
//				if err != nil {
//					return
//				}
//				defer gzipR.Close()
//				resBody, _ = io.ReadAll(gzipR)
//			}
//
//			fmt.Printf("RECEIVED\nBODY:\t%s\nSTATUS:\t%v\n", string(resBody), res.StatusCode)
//
//			if len(string(resBody)) != len(tt.want.response) {
//				t.Errorf("Expected body %s, got %s", tt.want.response, string(resBody))
//			}
//
//			// заголовок ответа
//			if res.Header.Get("Content-Type") != tt.want.contentType {
//				t.Errorf("Expected Content-Type %s, got %s", tt.want.contentType, res.Header.Get("Content-Type"))
//			}
//
//		})
//		time.Sleep(1 * time.Second)
//	}
//}

//	type testInput struct {
//		encoding string
//		name     string
//		method   string
//		body     string
//		target   string
//	}
//
//	type testOutput struct {
//		acceptType  string
//		code        int
//		response    string
//		contentType string
//	}
type ServerTestSuite struct {
	suite.Suite
	client http.Client
	//input  testInput
	//output testOutput
}

func (suite *ServerTestSuite) SetupTest() {
	cfg.Initialize()
	storage.Initialize()
	go server.Run()
	suite.client = http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	cfg.Storage.StorageType = "none"
	time.Sleep(1 * time.Second)
}

func (suite *ServerTestSuite) TestGetPostRequest() {
	srcURL := "https://www.google.com"
	srcReader := bytes.NewBuffer([]byte(srcURL))
	fmt.Println("INPUT URL:\t", srcURL)

	respP, err := suite.client.Post(cfg.Server.BaseURL+"/", "text/plain", srcReader)
	fmt.Println("POST response:\t\t", respP)
	fmt.Println("POST error   :\t\t", err)
	fmt.Println("current suite:\t\t", suite)
	suite.Assert().NoError(err)
	suite.Assert().Equal(http.StatusCreated, respP.StatusCode)

	var bodyP []byte
	bodyP, err = io.ReadAll(respP.Body)
	suite.Assert().NoError(err)

	fmt.Println("SHORT URL:\t", string(bodyP))

	respG, err := suite.client.Get(string(bodyP))
	suite.Assert().NoError(err)
	suite.Assert().Equal(http.StatusTemporaryRedirect, respG.StatusCode)

	bodyG := respG.Header.Get("Location")
	fmt.Println("RESULT URL:\t", bodyG)

	suite.Equal(srcURL, bodyG)

	err = respP.Body.Close()
	suite.Assert().NoError(err)
	err = respG.Body.Close()
	suite.Assert().NoError(err)
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}
