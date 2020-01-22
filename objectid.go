package grok

import "go.mongodb.org/mongo-driver/bson/primitive"

// ObjectIDFromHex ...
func ObjectIDFromHex(hex string) primitive.ObjectID {
	id, _ := primitive.ObjectIDFromHex(hex)
	return id
}
