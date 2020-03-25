package grok

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// EnsureStoreFromPath ...
func EnsureStoreFromPath(provider Provider, paramName string) gin.HandlerFunc {
	return ensureStoreByKind("path", provider, paramName)
}

// EnsureStoreFromQuery ...
func EnsureStoreFromQuery(provider Provider, paramName string) gin.HandlerFunc {
	return ensureStoreByKind("query", provider, paramName)
}

func ensureStoreByKind(kind string, provider Provider, paramName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("sub")
		value := ""

		switch kind {
		case "path":
			value = c.Param(paramName)
		case "query":
			value = c.Query(paramName)
		}

		if err := ensureStore(provider, userID, value); err != nil {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}

func ensureStore(provider Provider, userID, storeID string) error {
	user, err := provider.Fetch(userID)

	if err != nil {
		logrus.WithError(err).
			Error("error fetching user data")
		return err
	}

	for _, store := range user.Stores {
		if store == storeID {
			return nil
		}
	}

	return errors.New("user not allowed to store")
}
