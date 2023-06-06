package server

import (
	"bytes"
	rand "crypto/rand"
	"encoding/json"
	"io"
	"math/big"
	"net/http"
	"testing"
	"time"

	"shortener/internal/auth"
	"shortener/internal/cfg"
	"shortener/internal/storage"
)

func BenchmarkServer(b *testing.B) {
	b.ReportAllocs()
	// initializing server
	config := cfg.New(cfg.FromDefaults())
	auth.New(config)
	storage.New(config)
	go Run(config)
	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	time.Sleep(1 * time.Second)

	var strings []string
	urls := map[string]string{}
	for i := 0; i < b.N; i++ {
		str := generateString(30)
		b := bytes.NewBuffer([]byte(str))

		resp, _ := client.Post(config.Server.BaseURL+"/", "text/plain", b)
		body, _ := io.ReadAll(resp.Body)
		strings = append(strings, string(body))
		urls[string(body)] = str
		_ = resp.Body.Close()
	}
	for i := 0; i < b.N*5; i++ {
		bigN, _ := rand.Int(rand.Reader, big.NewInt(int64(len(strings))))
		d := bigN.Int64()
		short := strings[d]

		resp, err := client.Get(short)
		if http.StatusTemporaryRedirect != resp.StatusCode {
			b.Errorf("urls don't match because %t", err)
		}
		err = resp.Body.Close()
		if err != nil {
			return
		}

		num, _ := rand.Int(rand.Reader, big.NewInt(100))
		ch := num.Int64()
		if ch < 15 {
			reqBody, _ := json.Marshal([]string{short})
			srcReader := bytes.NewBuffer(reqBody)
			req, _ := http.NewRequest("DELETE", config.Server.BaseURL+"/api/user/urls",
				srcReader)
			respD, _ := client.Do(req)
			err = req.Body.Close()
			if err != nil {
				b.Errorf("couldn't delete %t, status: %d", err, respD.StatusCode)
			}
			err = respD.Body.Close()
			if err != nil || http.StatusAccepted != respD.StatusCode {
				b.Errorf("couldn't delete %t, status: %d", err, respD.StatusCode)
			}
		}
	}
}

var (
	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	length  = big.NewInt(int64(len(letters)))
)

func generateString(n int) string {
	b := make([]rune, n)
	for i := range b {
		num, _ := rand.Int(rand.Reader, length)
		b[i] = letters[num.Int64()]
	}
	return string(b)
}
