package middleware

import (
	"compress/gzip"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Gunzip
// Performs the data decompression if the contentType is gzip.
// If no errors met during unpacking, passes the request to next handler.
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

// Gzip
// Performs the data compression if the acceptsType is application/gzip.
// Adds layer to gin.ResponseWriter that performs compression
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

// gzipWriter is a custom writer user for compressing responses
type gzipWriter struct {
	gin.ResponseWriter
	writer *gzip.Writer
}

// WriteString replaces original method so that the output is compressed using gzip
func (g *gzipWriter) WriteString(s string) (int, error) {
	res, err := g.writer.Write([]byte(s))
	if err != nil {
		return 0, fmt.Errorf("while writing with gzip:\n %w", err)
	}
	return res, nil
}

// Write replaces original method so that the output is compressed using gzip
func (g *gzipWriter) Write(data []byte) (int, error) {
	res, err := g.writer.Write(data)
	if err != nil {
		return 0, fmt.Errorf("while writing with gzip:\n %w", err)
	}
	return res, nil
}
