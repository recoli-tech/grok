package grok

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// EnsureStoreFromPath ...
func EnsureStoreFromPath(paramName string) gin.HandlerFunc {
	return ensureStoreByKind("path", paramName)
}

// EnsureStoreFromQuery ...
func EnsureStoreFromQuery(paramName string) gin.HandlerFunc {
	return ensureStoreByKind("query", paramName)
}

func ensureStoreByKind(kind string, paramName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		value := ""

		switch kind {
		case "path":
			value = c.Param(paramName)
		case "query":
			value = c.Query(paramName)
		}

		if err := EnsureStore(c, value); err != nil {
			c.Error(err)
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}

// EnsureStore ...
func EnsureStore(ctx *gin.Context, storeID string) error {
	value, exists := ctx.Get("stores")

	if !exists || value == nil {
		return errors.New("stores parameter not found in user claims")
	}

	slice, ok := value.([]interface{})

	if !ok {
		return errors.New("stores parameter not found in user claims")
	}

	for _, store := range slice {
		if s, ok := store.(string); ok && s == storeID {
			return nil
		}
	}

	return errors.New("user not allowed to store")
}
