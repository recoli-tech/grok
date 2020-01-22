package grok

import (
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	// Validator ...
	Validator = NewValidator()
)

// NewValidator ...
func NewValidator() *validator.Validate {
	validate := validator.New()

	validate.RegisterValidation("objectid", func(fl validator.FieldLevel) bool {
		if _, err := primitive.ObjectIDFromHex(fl.Field().String()); err != nil {
			return false
		}

		return true
	})

	return validate
}
