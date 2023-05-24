package middleware

import (
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func Gunzip() gin.HandlerFunc {
	return func(c *gin.Context) {
		contentType := c.GetHeader("Content-Encoding")
		if !strings.Contains(contentType, "gzip") {
			c.Next()
			return
		}
		gzipR, err := gzip.NewReader(c.Request.Body)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		defer func(gzipR *gzip.Reader) {
			err := gzipR.Close()
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return
			}
		}(gzipR)
		c.Request.Body = gzipR
		c.Next()
	}
}

func Gzip() gin.HandlerFunc {
	return func(c *gin.Context) {
		acceptsType := c.GetHeader("Accept-Encoding")
		if !strings.Contains(acceptsType, "gzip") {
			return
		}
		gzipW, err := gzip.NewWriterLevel(c.Writer, gzip.BestSpeed)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		c.Writer.Header().Set("Content-Encoding", "gzip")
		c.Writer = &gzipWriter{c.Writer, gzipW}
		defer func(gzipW *gzip.Writer) {
			err := gzipW.Close()
			if err != nil {
				c.Status(http.StatusBadRequest)
				return
			}
		}(gzipW)
		c.Next()
	}
}

type gzipWriter struct {
	gin.ResponseWriter
	writer *gzip.Writer
}

func (g *gzipWriter) WriteString(s string) (int, error) {
	return g.writer.Write([]byte(s))
}

func (g *gzipWriter) Write(data []byte) (int, error) {
	return g.writer.Write(data)
}
