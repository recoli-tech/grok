package grok

import (
	"context"
	"reflect"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

// MongoRepository ...
type MongoRepository struct {
	idProperty   string
	documentType reflect.Type
	collection   *mongo.Collection
}

// NewMongoRepository ...
func NewMongoRepository(idProperty string, documentType reflect.Type, collection *mongo.Collection) *MongoRepository {
	return &MongoRepository{
		idProperty:   idProperty,
		documentType: documentType,
		collection:   collection,
	}
}

// Insert ...
func (r *MongoRepository) Insert(ctx context.Context, document interface{}) (interface{}, error) {
	result, err := r.collection.InsertOne(ctx, document)

	if err != nil {
		return nil, err
	}

	field := reflect.ValueOf(document).Elem().FieldByName(r.idProperty)

	if !field.IsValid() || !field.CanSet() {
		logrus.WithField("document", document).
			Panicf("property %s is invalid or cannot be set", r.idProperty)
	}

	field.Set(reflect.ValueOf(result.InsertedID))

	return document, err
}

// Update ...
func (r *MongoRepository) Update(ctx context.Context, id interface{}, document interface{}) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{
		"_id": id,
	}, bson.M{"$set": document})

	return err
}

// FindByID ...
func (r *MongoRepository) FindByID(ctx context.Context, id interface{}) (interface{}, error) {
	doc := reflect.New(r.documentType).Interface()

	err := r.collection.FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(doc)

	return doc, err
}
