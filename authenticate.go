package grok

import (
	"net/http"
	"strings"
	"time"

	"github.com/auth0-community/go-auth0"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"

	"gopkg.in/square/go-jose.v2"
)

const (
	// AuthClaimNamespace ...
	AuthClaimNamespace = "https://api.recoli.com.br/"
)

// Authenticate ...
type Authenticate interface {
	Middleware() gin.HandlerFunc
}

// Auth0Authenticate ...
type Auth0Authenticate struct {
	memoryCache    *cache.Cache
	auth           *APIAuth
	auth0Validator *auth0.JWTValidator
}

// CreateAuthenticate ...
func CreateAuthenticate(auth *APIAuth, cache *cache.Cache) Authenticate {
	if auth.Fake {
		return NewFakeAuthenticate(
			auth.FakeConfig.Authenticated,
			auth.FakeConfig.Claims,
		)
	}

	return NewAuthenticate(auth, cache)
}

// NewAuthenticate ...
func NewAuthenticate(auth *APIAuth, cache *cache.Cache) Authenticate {
	a := &Auth0Authenticate{auth: auth, memoryCache: cache}

	a.auth0Validator = auth0.NewValidator(
		auth0.NewConfiguration(
			auth0.NewJWKClient(
				auth0.JWKClientOptions{
					URI: a.auth.JWKS,
				},
				nil),
			a.auth.Audience,
			a.auth.Tenant,
			jose.RS256,
		),
		nil,
	)

	return a
}

// Middleware ...
func (a *Auth0Authenticate) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwt := c.Request.Header.Get("authorization")

		if claims, found := a.memoryCache.Get(jwt); found {
			a.setKeys(c, claims.(map[string]interface{}))
			c.Next()
			return
		}

		token, err := a.auth0Validator.ValidateRequest(c.Request)

		if err != nil {
			c.Error(err)
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		claims := make(map[string]interface{})
		if err := a.auth0Validator.Claims(c.Request, token, &claims); err != nil {
			c.Error(err)
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		a.setKeys(c, claims)

		if exp, ok := claims["exp"]; ok {
			float := exp.(float64)
			a.memoryCache.Set(jwt, claims, time.Second*time.Duration(int64(float)))
		}

		c.Next()
	}
}

func (a *Auth0Authenticate) setKeys(ctx *gin.Context, claims map[string]interface{}) {
	for key, value := range claims {
		if strings.Index(key, AuthClaimNamespace) >= 0 {
			key = strings.Replace(key, AuthClaimNamespace, "", -1)
		}

		ctx.Set(key, value)
	}
}

// GetSubFromContext removes auth0| prefix from sub
func GetSubFromContext(ctx gin.Context) string {
	sub := ctx.GetString("sub")
	splited := strings.Split(sub, "|")

	if len(splited) < 2 {
		return sub
	}

	return splited[1]
}

// RequiredClaims ...
func RequiredClaims(claims ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, c := range claims {
			if _, exists := ctx.Get(c); !exists {
				ctx.AbortWithStatus(http.StatusForbidden)
				return
			}
		}

		ctx.Next()
	}
}
