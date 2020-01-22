package grok

// Container ...
type Container interface {
	Close() error
	Controllers() []APIController
}
