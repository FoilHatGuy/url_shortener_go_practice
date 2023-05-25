package server

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/suite"
	"io"
	"net/http"
	_ "net/http/pprof"
	"shortener/internal/cfg"
	"shortener/internal/storage"
	"testing"
	"time"
)

type ServerTestSuite struct {
	suite.Suite
	client http.Client
	config *cfg.ConfigT
}

func (s *ServerTestSuite) SetupSuite() {
	config := cfg.Initialize()
	s.config = config
	s.config.Storage.StorageType = "none"
	storage.Initialize(config)
	Run(config)
	s.client = http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	time.Sleep(1 * time.Second)
}
func (s *ServerTestSuite) TestPing() {

	respG, err := s.client.Get(s.config.Server.BaseURL + "/ping")
	s.Assert().NoError(err)
	s.Assert().Equal(http.StatusOK, respG.StatusCode)

	err = respG.Body.Close()
	s.Assert().NoError(err)
}

func (s *ServerTestSuite) TestGetPostRequest() {
	srcURL := "https://www.TestGetPostRequest.com"
	srcReader := bytes.NewBuffer([]byte(srcURL))
	fmt.Println("INPUT URL:\t", srcURL)

	respP, err := s.client.Post(s.config.Server.BaseURL+"/", "text/plain", srcReader)
	fmt.Println("POST response:\t\t", respP)
	fmt.Println("POST error   :\t\t", err)
	s.Assert().NoError(err)
	s.Assert().Equal(http.StatusCreated, respP.StatusCode)

	var bodyP []byte
	bodyP, err = io.ReadAll(respP.Body)
	s.Assert().NoError(err)

	fmt.Println("SHORT URL:\t", string(bodyP))

	respG, err := s.client.Get(string(bodyP))
	s.Assert().NoError(err)
	s.Assert().Equal(http.StatusTemporaryRedirect, respG.StatusCode)

	bodyG := respG.Header.Get("Location")
	fmt.Println("RESULT URL:\t", bodyG)

	s.Equal(srcURL, bodyG)

	err = respP.Body.Close()
	s.Assert().NoError(err)
	err = respG.Body.Close()
	s.Assert().NoError(err)
}

func (s *ServerTestSuite) TestBatchRequest() {
	srcURL := "https://www.TestBatchRequest.com"
	type reqElement struct {
		LineID string `json:"correlation_id"`
		URL    string `json:"original_url"`
	}
	type resElement struct {
		LineID string `json:"correlation_id"`
		URL    string `json:"short_url"`
	}
	reqBody, _ := json.Marshal([]reqElement{{"TestBatchRequest", srcURL}})
	srcReader := bytes.NewBuffer(reqBody)
	fmt.Println("INPUT URL:\t", srcReader.String())

	respP, err := s.client.Post(s.config.Server.BaseURL+"/api/shorten/batch", "application/json", srcReader)
	fmt.Println("POST response:\t\t", respP)
	fmt.Println("POST error   :\t\t", err)
	s.Assert().NoError(err)
	s.Assert().Equal(http.StatusCreated, respP.StatusCode)

	s.Assert().NoError(err)
	var resBody []resElement
	bodyR, err := io.ReadAll(respP.Body)
	s.Assert().NoError(err)
	err = json.Unmarshal(bodyR, &resBody)
	s.Assert().NoError(err)

	fmt.Println("response body: ", resBody[0].URL)
	respG, err := s.client.Get(resBody[0].URL)
	s.Assert().NoError(err)
	s.Assert().Equal(http.StatusTemporaryRedirect, respG.StatusCode)

	bodyG := respG.Header.Get("Location")
	fmt.Println("RESULT URL:\t", bodyG)
	s.Equal(srcURL, string(bodyG))

	err = respP.Body.Close()
	s.Assert().NoError(err)
	err = respG.Body.Close()
	s.Assert().NoError(err)
}

func (s *ServerTestSuite) TestGzipRequest() {
	if true {
		return
	} else {
		srcURL := "https://www.TestGzipRequest.com"
		srcReader := bytes.NewBuffer([]byte(srcURL))
		fmt.Println("INPUT URL:\t", srcURL)

		var b bytes.Buffer
		gz := gzip.NewWriter(&b)

		_, err := gz.Write(srcReader.Bytes())
		s.Assert().NoError(err)

		respP, err := s.client.Post(s.config.Server.BaseURL+"/", "text/plain", &b)

		fmt.Println("POST response:\t\t", respP)
		fmt.Println("POST error   :\t\t", err)
		fmt.Println("current suite:\t\t", s)
		s.Assert().NoError(err)
		s.Assert().Equal(http.StatusCreated, respP.StatusCode)

		var bodyP []byte
		bodyP, err = io.ReadAll(respP.Body)
		s.Assert().NoError(err)

		fmt.Println("SHORT URL:\t", string(bodyP))

		//respG, err := suite.client.Get(string(bodyP))
		req, _ := http.NewRequest("GET", string(bodyP), &b)
		req.Header.Set("Accept-Encoding", "application/gzip")
		respG, _ := s.client.Do(req)
		s.Assert().NoError(err)
		s.Assert().Equal(http.StatusTemporaryRedirect, respG.StatusCode)

		bodyG := respG.Header.Get("Location")

		fmt.Println("RESULT URL:\t", bodyG)

		s.Equal(srcURL, bodyG)

		err = respP.Body.Close()
		s.Assert().NoError(err)
		err = respG.Body.Close()
		s.Assert().NoError(err)
	}
}

func (s *ServerTestSuite) TestDeleteRequest() {
	srcURL := "https://www.TestDeleteRequest.com"
	srcReader := bytes.NewBuffer([]byte(srcURL))
	fmt.Println("INPUT URL:\t", srcURL)

	respP, err := s.client.Post(s.config.Server.BaseURL+"/", "text/plain", srcReader)
	fmt.Println("POST response:\t\t", respP)
	fmt.Println("POST error   :\t\t", err)
	s.Assert().NoError(err)
	s.Assert().Equal(http.StatusCreated, respP.StatusCode)

	var bodyP []byte
	bodyP, err = io.ReadAll(respP.Body)
	s.Assert().NoError(err)

	fmt.Println("SHORT URL:\t", string(bodyP))

	b := bytes.NewBuffer(bodyP)
	req, err := http.NewRequest("DELETE", s.config.Server.BaseURL+"/user/urls", b)
	s.Assert().NoError(err)
	resp, err := s.client.Do(req)
	s.Assert().NoError(err)
	err = resp.Body.Close()
	s.Assert().NoError(err)

	respG, err := s.client.Get(string(bodyP))
	s.Assert().NoError(err)
	s.Assert().Equal(http.StatusTemporaryRedirect, respG.StatusCode)

	bodyG := respG.Header.Get("Location")
	fmt.Println("RESULT URL:\t", bodyG)

	s.Equal(srcURL, bodyG)

	err = req.Body.Close()
	s.Assert().NoError(err)
	err = respP.Body.Close()
	s.Assert().NoError(err)
	err = respG.Body.Close()
	s.Assert().NoError(err)
}

func (s *ServerTestSuite) TestGetUserRequest() {
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}
