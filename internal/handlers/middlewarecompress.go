package handlers

import "github.com/gin-gonic/gin"

func ArchiveData() gin.HandlerFunc {
	return func(c *gin.Context) {
		// if compressed, decompress
		c.Next()
		// if accepts compressed, compress
	}
}
