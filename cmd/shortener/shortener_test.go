package main

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
	"testing"
	"time"
)

type ServerTestSuite struct {
	suite.Suite
	client http.Client
}

func (suite *ServerTestSuite) SetupSuite() {
	go main()
	suite.client = http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	cfg.Storage.StorageType = "none"
	time.Sleep(1 * time.Second)
}
func (suite *ServerTestSuite) TestPing() {

	respG, err := suite.client.Get(cfg.Server.BaseURL + "/ping")
	suite.Assert().NoError(err)
	suite.Assert().Equal(http.StatusOK, respG.StatusCode)

	err = respG.Body.Close()
	suite.Assert().NoError(err)
}

func (suite *ServerTestSuite) TestGetPostRequest() {
	srcURL := "https://www.TestGetPostRequest.com"
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

func (suite *ServerTestSuite) TestBatchRequest() {
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

	respP, err := suite.client.Post(cfg.Server.BaseURL+"/api/shorten/batch", "application/json", srcReader)
	fmt.Println("POST response:\t\t", respP)
	fmt.Println("POST error   :\t\t", err)
	fmt.Println("current suite:\t\t", suite)
	suite.Assert().NoError(err)
	suite.Assert().Equal(http.StatusCreated, respP.StatusCode)

	suite.Assert().NoError(err)
	var resBody []resElement
	bodyR, err := io.ReadAll(respP.Body)
	suite.Assert().NoError(err)
	err = json.Unmarshal(bodyR, &resBody)
	suite.Assert().NoError(err)

	fmt.Println("response body: ", resBody[0].URL)
	respG, err := suite.client.Get(resBody[0].URL)
	suite.Assert().NoError(err)
	suite.Assert().Equal(http.StatusTemporaryRedirect, respG.StatusCode)

	bodyG := respG.Header.Get("Location")
	fmt.Println("RESULT URL:\t", bodyG)
	suite.Equal(srcURL, string(bodyG))

	err = respP.Body.Close()
	suite.Assert().NoError(err)
	err = respG.Body.Close()
	suite.Assert().NoError(err)
}

func (suite *ServerTestSuite) TestGzipRequest() {
	if true {
		return
	} else {
		srcURL := "https://www.TestGzipRequest.com"
		srcReader := bytes.NewBuffer([]byte(srcURL))
		fmt.Println("INPUT URL:\t", srcURL)

		var b bytes.Buffer
		gz := gzip.NewWriter(&b)

		_, err := gz.Write(srcReader.Bytes())
		suite.Assert().NoError(err)

		respP, err := suite.client.Post(cfg.Server.BaseURL+"/", "text/plain", &b)

		fmt.Println("POST response:\t\t", respP)
		fmt.Println("POST error   :\t\t", err)
		fmt.Println("current suite:\t\t", suite)
		suite.Assert().NoError(err)
		suite.Assert().Equal(http.StatusCreated, respP.StatusCode)

		var bodyP []byte
		bodyP, err = io.ReadAll(respP.Body)
		suite.Assert().NoError(err)

		fmt.Println("SHORT URL:\t", string(bodyP))

		//respG, err := suite.client.Get(string(bodyP))
		req, _ := http.NewRequest("GET", string(bodyP), &b)
		req.Header.Set("Accept-Encoding", "application/gzip")
		respG, _ := suite.client.Do(req)
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
}

func (suite *ServerTestSuite) TestDeleteRequest() {
	srcURL := "https://www.TestDeleteRequest.com"
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

	b := bytes.NewBuffer(bodyP)
	req, err := http.NewRequest("DELETE", cfg.Server.BaseURL+"/user/urls", b)
	suite.Assert().NoError(err)
	_, err = suite.client.Do(req)
	suite.Assert().NoError(err)

	respG, err := suite.client.Get(string(bodyP))
	suite.Assert().NoError(err)
	suite.Assert().Equal(http.StatusTemporaryRedirect, respG.StatusCode)

	bodyG := respG.Header.Get("Location")
	fmt.Println("RESULT URL:\t", bodyG)

	suite.Equal(srcURL, bodyG)

	err = req.Body.Close()
	suite.Assert().NoError(err)
	err = respP.Body.Close()
	suite.Assert().NoError(err)
	err = respG.Body.Close()
	suite.Assert().NoError(err)
}

func (suite *ServerTestSuite) TestGetUserRequest() {
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}
