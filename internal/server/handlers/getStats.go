package handlers

import (
	"net"
	"net/http"

	"shortener/internal/cfg"

	"github.com/gin-gonic/gin"

	"shortener/internal/storage"
)

// GetStats
// Ping server+database activity
func GetStats(dbController storage.DatabaseORM, config *cfg.ConfigT) gin.HandlerFunc {
	return func(c *gin.Context) {
		k := c.GetHeader("X-Real-IP")
		if k == "" {
			c.Status(http.StatusInternalServerError)
			return
		}
		ip := net.ParseIP(k)
		_, cidr, err := net.ParseCIDR(config.Server.TrustedSubnet)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		valid := cidr.Contains(ip)
		if !valid {
			c.Status(http.StatusForbidden)
			return
		}

		data, err := dbController.GetStats(c)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, data)
	}
}
