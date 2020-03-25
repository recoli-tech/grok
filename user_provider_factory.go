package grok

import (
	"os"
	"time"

	"github.com/patrickmn/go-cache"
	"gopkg.in/auth0.v3/management"
)

// UserProviderFactory ...
func UserProviderFactory(settings *Settings) Provider {
	switch settings.UserProvider.Kind {
	case "fake":
		return createMockProvider(settings)
	default:
		return createAuth0Provider(settings)
	}
}

func createAuth0Provider(settings *Settings) Provider {
	cache := cache.New(
		time.Duration(settings.UserProvider.Auth0.CacheTTL)*time.Minute,
		10*time.Second)

	clientID := settings.UserProvider.Auth0.ClientID
	clientSecret := settings.UserProvider.Auth0.ClientSecret

	if settings.UserProvider.Auth0.ClientFrom == "environment" {
		clientID = os.Getenv(settings.UserProvider.Auth0.ClientIDEnv)
		clientSecret = os.Getenv(settings.UserProvider.Auth0.ClientSecretEnv)
	}

	management, _ := management.New(
		settings.UserProvider.Auth0.Domain,
		clientID,
		clientSecret,
	)

	return NewAuth0Provider(cache, management)
}

func createMockProvider(settings *Settings) Provider {
	users := []*User{}

	for _, user := range settings.UserProvider.Mock {
		users = append(users, &User{
			ID:     user["id"].(string),
			Stores: parseStores(user["stores"].([]interface{})),
			Email:  user["email"].(string),
		})
	}

	return NewMockProvider(users...)
}

func parseStores(stores []interface{}) []string {
	ss := []string{}
	for _, s := range stores {
		ss = append(ss, s.(string))
	}

	return ss
}
