package grok_test

import (
	"testing"

	"github.com/recoli-tech/grok"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestObjectIDFromHex(t *testing.T) {
	objectID := primitive.NewObjectID().Hex()

	assert.NotEqual(t, primitive.NilObjectID, grok.ObjectIDFromHex(objectID))
	assert.Equal(t, primitive.NilObjectID, grok.ObjectIDFromHex("wrongobjectid"))
}
