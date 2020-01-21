package container

import "github.com/raafvargas/grok/controllers"

// Container ...
type Container interface {
	Close() error
	Controllers() []controllers.Controller
}
