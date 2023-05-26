package server

import (
	"bytes"
	rand "crypto/rand"
	"encoding/json"
	"io"
	"math/big"
	mrng "math/rand"
	"net/http"
	_ "net/http/pprof"
	"shortener/internal/cfg"
	"shortener/internal/security"
	"shortener/internal/storage"
	"testing"
	"time"
)

func BenchmarkServer(b *testing.B) {
	b.ReportAllocs()
	// initialising server
	config := cfg.Initialize()
	security.Init(config)
	config.Storage.StorageType = "none"
	storage.Initialize(config)
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
		d := mrng.Intn(len(strings))
		short := strings[d]

		resp, err := client.Get(short)
		if http.StatusTemporaryRedirect != resp.StatusCode {
			b.Errorf("urls don't match because %t", err)
		}
		err = resp.Body.Close()
		if err != nil {
			return
		}

		ch := mrng.Intn(100)
		if ch < 15 {
			reqBody, _ := json.Marshal([]string{short})
			srcReader := bytes.NewBuffer(reqBody)
			req, _ := http.NewRequest("DELETE", config.Server.BaseURL+"/api/user/urls",
				srcReader)
			respD, _ := client.Do(req)
			err = req.Body.Close()
			if err != nil {
				return
			}
			err = respD.Body.Close()
			if err != nil {
				return
			}
			if http.StatusAccepted != respD.StatusCode {
				b.Errorf("couldn't delete %t, status: %d", err, respD.StatusCode)
			}

		}

	}
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var length = big.NewInt(int64(len(letters)))

func generateString(n int) string {
	b := make([]rune, n)
	for i := range b {
		num, _ := rand.Int(rand.Reader, length)
		b[i] = letters[num.Int64()]
	}
	return string(b)
}
