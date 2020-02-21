package grok

import (
	"net/http"
	"time"

	"github.com/auth0-community/go-auth0"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"

	"gopkg.in/square/go-jose.v2"
)

// Authenticate ...
func Authenticate(auth *APIAuth, cache *cache.Cache) gin.HandlerFunc {
	return func(c *gin.Context) {
		jwt := c.Request.Header.Get("authorization")

		if value, found := cache.Get(jwt); found {
			for key, value := range value.(map[string]interface{}) {
				c.Set(key, value)
			}
			c.Next()
			return
		}

		client := auth0.NewJWKClient(auth0.JWKClientOptions{URI: auth.JWKS}, nil)
		configuration := auth0.NewConfiguration(client, auth.Audience, auth.Tenant, jose.RS256)
		validator := auth0.NewValidator(configuration, nil)

		token, err := validator.ValidateRequest(c.Request)

		if err != nil {
			c.Error(err)
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		claims := make(map[string]interface{})
		if err := validator.Claims(c.Request, token, &claims); err != nil {
			c.Error(err)
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		for key, value := range claims {
			c.Set(key, value)
		}

		if exp, ok := claims["exp"]; ok {
			float := exp.(float64)
			cache.Set(jwt, claims, time.Second*time.Duration(int64(float)))
		}

		c.Next()
	}
}
