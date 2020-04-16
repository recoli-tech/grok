package grok

import (
	"errors"
	"fmt"

	"github.com/patrickmn/go-cache"
	"gopkg.in/auth0.v3/management"
)

// User ...
type User struct {
	ID     string
	Stores []string
	Email  string
}

// Provider ...
type Provider interface {
	Fetch(id string) (*User, error)
}

type auth0Provider struct {
	cache           *cache.Cache
	auth0Management *management.Management
}

// NewAuth0Provider ...
func NewAuth0Provider(
	cache *cache.Cache,
	auth0Management *management.Management) Provider {
	return &auth0Provider{
		cache:           cache,
		auth0Management: auth0Management,
	}
}

func (p *auth0Provider) Fetch(id string) (*User, error) {
	if user, ok := p.cache.Get(id); ok {
		return user.(*User), nil
	}

	list, err := p.auth0Management.User.Search(
		management.Parameter("q", fmt.Sprintf("user_id:*%s*", id)))

	if err != nil {
		return nil, err
	}

	if list.Length <= 0 {
		return nil, errors.New("user not found")
	}

	auth0User := list.Users[0]
	user := new(User)
	user.ID = id
	user.Email = auth0User.GetEmail()
	value, ok := auth0User.UserMetadata["stores"]

	if !ok {
		return nil, errors.New("no store found")
	}

	if stores, ok := value.([]interface{}); ok {
		for _, s := range stores {
			if store, ok := s.(string); ok {
				user.Stores = append(user.Stores, store)
			}
		}
	}

	p.cache.SetDefault(id, user)

	return user, nil
}

type mockProvider struct {
	users []*User
}

// NewMockProvider ...
func NewMockProvider(users ...*User) Provider {
	return &mockProvider{
		users: users,
	}
}

func (p *mockProvider) Fetch(id string) (*User, error) {
	for _, u := range p.users {
		if id == u.ID {
			return u, nil
		}
	}

	return nil, errors.New("user not found")
}
