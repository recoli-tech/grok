package grok

import (
	"github.com/gin-gonic/gin"
)

// SetAccessTokenInContext ...
func SetAccessTokenInContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		c.Set("access_token", token)
		c.Next()
		return
	}
}
