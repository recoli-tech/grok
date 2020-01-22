package grok_test

import "github.com/raafvargas/grok"
import "testing"

import "go.mongodb.org/mongo-driver/bson/primitive"

import "github.com/stretchr/testify/assert"

func TestObjectIDFromHex(t *testing.T) {
	objectID := primitive.NewObjectID().Hex()

	assert.NotEqual(t, primitive.NilObjectID, grok.ObjectIDFromHex(objectID))
	assert.Equal(t, primitive.NilObjectID, grok.ObjectIDFromHex("wrongobjectid"))
}
