package container

import "github.com/raafvargas/grok/http"

// Container ...
type Container interface {
	Close() error
	Controllers() []http.Controller
}
