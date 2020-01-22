package grok

import "github.com/gin-gonic/gin"

import "net/http"

// Authorize ...
func Authorize(scope string) gin.HandlerFunc {
	return func(c *gin.Context) {
		permissions := c.GetStringSlice("permissions")

		for _, permission := range permissions {
			if permission == scope {
				c.Next()
				return
			}
		}

		c.AbortWithStatus(http.StatusForbidden)
	}
}
