package middleware

import (
	"github.com/gin-gonic/gin"
	"shortener/internal/cfg"
	sec "shortener/internal/security"
	"strings"
)

var config *cfg.ConfigT

func init() {
	config = cfg.Initialize()
}
func Cooker() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("user")
		var key string
		if err == nil {
			//fmt.Println("UID COOKIE PRESENT:\n", cookie)

			key, err = sec.AuthEngine.Validate(cookie)
			//fmt.Println("VALIDATION RESULT:\n", key, err)
			if err == nil {
				c.SetCookie("user", cookie, config.Server.CookieLifetime, "/",
					strings.Split(config.Server.Address, ":")[0], false, true)
				c.Set("owner", key)
				//fmt.Println("UID KEY:\n", key)
				c.Next()
				return
			}
		}
		//fmt.Println("UID COOKIE MET ERROR:\n", err)
		cookie, key = sec.AuthEngine.Generate()
		//fmt.Println("NEW COOKIE GENERATED:\n", cookie)
		//fmt.Println("NEW UID KEY:\n", key)
		c.SetCookie("user", cookie, config.Server.CookieLifetime, "/", config.Server.BaseURL, false, true)
		c.Set("owner", key)
		c.Next()
	}
}
