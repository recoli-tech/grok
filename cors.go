package grok

import (
	"github.com/gin-gonic/gin"
)

// CORS ...
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Accept, *")
		c.Next()
	}
}
