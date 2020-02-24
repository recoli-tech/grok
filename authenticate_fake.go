package grok

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// FakeAuthenticate ...
type FakeAuthenticate struct {
	claims        map[string]interface{}
	authenticated bool
}

// NewFakeAuthenticate ...
func NewFakeAuthenticate(authenticated bool, claims map[string]interface{}) Authenticate {
	return &FakeAuthenticate{
		authenticated: authenticated,
		claims:        claims,
	}
}

// Middleware ...
func (a *FakeAuthenticate) Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !a.authenticated {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}

		for k, v := range a.claims {
			ctx.Set(k, v)
		}
	}
}
