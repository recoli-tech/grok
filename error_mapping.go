package grok

// ErrorMapping ...
type ErrorMapping map[error]error

var (
	// DefaultErrorMapping ...
	DefaultErrorMapping = ErrorMapping{}
)

// Register ...
func (mapping ErrorMapping) Register(k error, v error) {
	mapping[k] = v
}

// Exists ...
func (mapping ErrorMapping) Exists(err error) bool {
	_, has := mapping[err]

	return has
}

// Get ...
func (mapping ErrorMapping) Get(err error) error {
	return mapping[err]
}
