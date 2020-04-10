package grok

import (
	"reflect"

	"go.mongodb.org/mongo-driver/mongo"
)

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
	if reflect.TypeOf(err).Kind() == reflect.Struct {
		mongoError := mapping.mappingMongoError(err)

		if mongoError != nil {
			return true
		}

		return false
	}

	_, has := mapping[err]

	return has
}

// Get ...
func (mapping ErrorMapping) Get(err error) error {
	mongoError := mapping.mappingMongoError(err)

	if mongoError != nil {
		return mongoError
	}

	return mapping[err]
}

func (mapping ErrorMapping) mappingMongoError(err error) error {
	if exp, ok := err.(mongo.WriteException); ok {
		if len(exp.WriteErrors) > 0 {
			if result, has := mapping[mongo.WriteError{Code: exp.WriteErrors[0].Code}]; has {
				return result
			}
			return nil
		}
	}

	return nil
}
