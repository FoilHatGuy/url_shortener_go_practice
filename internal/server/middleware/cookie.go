package middleware

import (
	"github.com/gin-gonic/gin"
	"shortener/internal/cfg"
	. "shortener/internal/security"
)

func Cooker() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("user")
		var key string
		if err == nil {
			//fmt.Println("UID COOKIE PRESENT:\n", cookie)

			key, err = AuthEngine.Validate(cookie)
			//fmt.Println("VALIDATION RESULT:\n", key, err)
			if err == nil {
				c.SetCookie("user", cookie, cfg.Server.CookieLifetime, "/", cfg.Server.BaseURL, false, true)
				c.Set("owner", key)
				//fmt.Println("UID KEY:\n", key)
				c.Next()
				return
			}
		}
		//fmt.Println("UID COOKIE MET ERROR:\n", err)
		cookie, key = AuthEngine.Generate()
		//fmt.Println("NEW COOKIE GENERATED:\n", cookie)
		//fmt.Println("NEW UID KEY:\n", key)
		c.SetCookie("user", cookie, cfg.Server.CookieLifetime, "/", cfg.Server.BaseURL, false, true)
		c.Set("owner", key)
		c.Next()
	}
}
