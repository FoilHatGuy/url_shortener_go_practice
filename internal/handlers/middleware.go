package handlers

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net"
	"net/http"
	"strings"
)

type customWriter struct {
	body   []byte
	writer gin.ResponseWriter
}

func ArchiveData() gin.HandlerFunc {
	return func(c *gin.Context) {
		var wb *ResponseBuffer
		if w, ok := c.Writer.(gin.ResponseWriter); ok {
			wb = NewResponseBuffer(w)
			c.Writer = wb
		}

		contentType := c.GetHeader("Content-Type")
		fmt.Println(contentType)
		if strings.Contains(contentType, "gzip") {
			body := &bytes.Buffer{}
			_, err := body.ReadFrom(c.Request.Body)
			if err != nil {
				return
			}
			defer c.Request.Body.Close()
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
			defer c.Request.Body.Close()
			fmt.Println(body.String())
			c.Set("Body", body.String())
		}

		c.Next()

		acceptsType := c.GetHeader("Accept-Encoding")
		if strings.Contains(acceptsType, "gzipp") {
			fmt.Println("dasdsdsadasdasdasdasdasdsa")
			c.Writer.Header().Set("Content-Encoding", "gzip")

			data := wb.Body
			fmt.Println("string before compression\t", data.String())
			wb.Body.Reset()

			buffer := bytes.Buffer{}

			gz := gzip.NewWriter(&buffer)
			_, _ = gz.Write(data.Bytes())
			err := gz.Close()

			fmt.Println("bytes after compression  \t", buffer.Bytes())
			if err != nil {
				return
			}
			defer gz.Close()

			_, err = wb.Write(buffer.Bytes())
			if err != nil {
				return
			}

			fmt.Println("string after compression\t", wb.Body.String())

		}
		wb.Flush()
	}
}

type ResponseBuffer struct {
	Response gin.ResponseWriter // the actual ResponseWriter to flush to
	status   int                // the HTTP response code from WriteHeader
	Body     *bytes.Buffer      // the response content body
	Flushed  bool
}

func (w *ResponseBuffer) Pusher() http.Pusher {
	return w.Response.Pusher()
}

func NewResponseBuffer(w gin.ResponseWriter) *ResponseBuffer {
	return &ResponseBuffer{
		Response: w, status: 200, Body: &bytes.Buffer{},
	}
}

func (w *ResponseBuffer) Header() http.Header {
	return w.Response.Header() // use the actual response header
}

func (w *ResponseBuffer) Write(buf []byte) (int, error) {

	w.Body.Write(buf)
	return len(buf), nil
}

func (w *ResponseBuffer) WriteString(s string) (n int, err error) {
	//w.WriteHeaderNow()
	//n, err = io.WriteString(w.ResponseWriter, s)
	//w.size += n
	n, err = w.Write([]byte(s))
	return
}

func (w *ResponseBuffer) Written() bool {
	return w.Body.Len() != -1
}

func (w *ResponseBuffer) WriteHeader(status int) {
	w.status = status
}

func (w *ResponseBuffer) WriteHeaderNow() {
	//if !w.Written() {
	//	w.size = 0
	//	w.ResponseWriter.WriteHeader(w.status)
	//}
}

func (w *ResponseBuffer) Status() int {
	return w.status
}

func (w *ResponseBuffer) Size() int {
	return w.Body.Len()
}

func (w *ResponseBuffer) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	//if w.size < 0 {
	//	w.size = 0
	//}
	return w.Response.(http.Hijacker).Hijack()
}

func (w *ResponseBuffer) CloseNotify() <-chan bool {
	return w.Response.(http.CloseNotifier).CloseNotify()
}

func (w *ResponseBuffer) Flush() {
	if w.Flushed {
		return
	}
	w.Response.WriteHeader(w.status)
	if w.Body.Len() > 0 {
		_, err := w.Response.Write(w.Body.Bytes())
		if err != nil {
			panic(err)
		}
		w.Body.Reset()
	}
	w.Flushed = true
}
