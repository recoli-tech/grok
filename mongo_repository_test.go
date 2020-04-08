package grok_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/raafvargas/grok"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RepositoryTestSuite struct {
	suite.Suite
	assert     *assert.Assertions
	settings   *grok.Settings
	client     *mongo.Client
	collection *mongo.Collection
	repository *grok.MongoRepository
}

type MongoTestDocument struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Value string             `bson:"value"`
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

func (s *RepositoryTestSuite) SetupTest() {
	s.assert = assert.New(s.T())
	s.settings = new(grok.Settings)

	err := grok.FromYAML("tests/config.yaml", s.settings)
	s.assert.NoError(err)

	s.client = grok.NewMongoConnection(s.settings.Mongo.ConnectionString)
	s.collection = s.client.Database(s.settings.Mongo.Database).Collection("grok")
	s.repository = grok.NewMongoRepository("ID", reflect.TypeOf(MongoTestDocument{}), s.collection)
}

func (s *RepositoryTestSuite) TestRepository() {
	doc := &MongoTestDocument{
		Value: uuid.New().String(),
	}

	result, err := s.repository.Insert(context.Background(), doc)

	s.assert.NoError(err)
	s.assert.IsType(&MongoTestDocument{}, result)
	s.assert.False(result.(*MongoTestDocument).ID.IsZero())

	oldValue := doc.Value
	result.(*MongoTestDocument).Value = uuid.New().String()

	err = s.repository.Update(context.Background(), result.(*MongoTestDocument).ID, doc)
	s.assert.NoError(err)

	result2, err := s.repository.FindByID(context.Background(), result.(*MongoTestDocument).ID)
	s.assert.NoError(err)
	s.assert.NotEqual(oldValue, result2.(*MongoTestDocument).Value)
}
