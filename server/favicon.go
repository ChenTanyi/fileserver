package server

import (
	_ "embed"

	"github.com/gin-gonic/gin"
)

//go:embed assets/favicon.svg
var favicon []byte

func FaviconMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/favicon.ico" {
			c.Data(200, "image/svg+xml", favicon)
			c.Abort()
			return
		}
		c.Next()
	}
}
