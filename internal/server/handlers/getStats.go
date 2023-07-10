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
	return func(ctx *gin.Context) {
		k := ctx.GetHeader("X-Real-IP")
		if k == "" {
			ctx.Status(http.StatusInternalServerError)
			return
		}
		ip := net.ParseIP(k)
		_, cidr, err := net.ParseCIDR(config.Server.TrustedSubnet)
		if err != nil {
			ctx.Status(http.StatusInternalServerError)
			return
		}
		valid := cidr.Contains(ip)
		if !valid {
			ctx.Status(http.StatusForbidden)
			return
		}

		data, err := dbController.GetStats(ctx)
		if err != nil {
			ctx.Status(http.StatusInternalServerError)
			return
		}
		ctx.JSON(http.StatusOK, data)
	}
}
