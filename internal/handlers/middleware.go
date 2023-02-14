package handlers

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
)

func ArchiveData() gin.HandlerFunc {
	return func(c *gin.Context) {
		contentType := c.GetHeader("Content-Type")
		c.Set("responseType", "")
		c.Set("responseStatus", 200)
		c.Set("responseBody", bytes.NewBuffer([]byte{}))
		fmt.Println(contentType)
		if strings.Contains(contentType, "gzip") {
			body := &bytes.Buffer{}
			_, err := body.ReadFrom(c.Request.Body)
			if err != nil {
				return
			}
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {

				}
			}(c.Request.Body)
			fmt.Println(body.String())

			reader := bytes.NewReader(body.Bytes())
			gzreader, e1 := gzip.NewReader(reader)
			if e1 != nil {
				fmt.Println(e1) // Maybe panic here, depends on your error handling.
			}

			output, e2 := io.ReadAll(gzreader)
			if e2 != nil {
				fmt.Println(e2)
			}

			result := string(output)
			fmt.Printf("ungzipped:%v", result)
			c.Set("Body", result)
		} else {
			body := &bytes.Buffer{}
			_, err := body.ReadFrom(c.Request.Body)
			if err != nil {
				return
			}
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {

				}
			}(c.Request.Body)
			fmt.Println(body.String())
			c.Set("Body", body.String())
		}

		c.Next()

		acceptsType := c.GetHeader("Accept-Encoding")
		respT, okT := c.Get("responseType")
		respType := respT.(string)
		respS, okS := c.Get("responseStatus")
		respStatus := respS.(int)
		respB, okB := c.Get("responseBody")
		respBody := respB.(*bytes.Buffer)
		if !(okT && okB && okS) {
			c.Status(http.StatusBadRequest)
			return
		}
		if strings.Contains(acceptsType, "gzip") {
			c.Writer.Header().Set("Content-Encoding", "gzip")

			fmt.Println("string before compression\t", respBody.String())

			buffer := bytes.Buffer{}
			gz := gzip.NewWriter(&buffer)
			_, _ = gz.Write(respBody.Bytes())
			err := gz.Close()
			if err != nil {
				return
			}
			fmt.Println("bytes after compression  \t", buffer.Bytes())
			fmt.Println("String after compression  \t", buffer.String())

			c.Data(respStatus, "application/gzip", buffer.Bytes())
		} else {
			switch respType {
			case "json":
				//FIXME only known structure so no problems, need to be dynamic
				newResBody := struct {
					Result string `json:"result"`
				}{}
				err := json.Unmarshal(respBody.Bytes(), &newResBody)
				if err != nil {
					return
				}
				c.IndentedJSON(respStatus, newResBody)
			case "text":
				c.String(respStatus, respBody.String())
			case "none":
				c.Status(respStatus)
			case "redirect":
				c.Redirect(http.StatusTemporaryRedirect, respBody.String())

			}
		}
	}
}
