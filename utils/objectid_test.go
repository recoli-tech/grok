package utils_test

import "github.com/raafvargas/grok/utils"
import "testing"

import "go.mongodb.org/mongo-driver/bson/primitive"

import "github.com/stretchr/testify/assert"

func TestObjectIDFromHex(t *testing.T) {
	objectID := primitive.NewObjectID().Hex()

	assert.NotEqual(t, primitive.NilObjectID, utils.ObjectIDFromHex(objectID))
	assert.Equal(t, primitive.NilObjectID, utils.ObjectIDFromHex("wrongobjectid"))
}
