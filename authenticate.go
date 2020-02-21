package grok

import (
	"net/http"

	"github.com/auth0-community/go-auth0"
	"github.com/gin-gonic/gin"

	"gopkg.in/square/go-jose.v2"
)

// Authenticate ...
func Authenticate(tenant, jwks string, audience []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		client := auth0.NewJWKClient(auth0.JWKClientOptions{URI: jwks}, nil)
		configuration := auth0.NewConfiguration(client, audience, tenant, jose.RS256)
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
	}
}

// AuthenticateWithConfig ...
func AuthenticateWithConfig(auth *APIAuth) gin.HandlerFunc {
	return Authenticate(auth.Tenant, auth.JWKS, auth.Audience)
}
