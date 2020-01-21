package repositories_test

import (
	"testing"

	"github.com/raafvargas/grok/repositories"
	"github.com/raafvargas/grok/settings"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MongoTestSuite struct {
	suite.Suite
	assert   *assert.Assertions
	settings *settings.Settings
}

func TestMongoTestSuite(t *testing.T) {
	suite.Run(t, new(MongoTestSuite))
}

func (s *MongoTestSuite) SetupSuite() {
	s.assert = assert.New(s.T())
	s.settings = &settings.Settings{}
	settings.FromYAML("../tests/config.yaml", s.settings)
}

func (s *MongoTestSuite) TestConnect() {
	s.assert.NotPanics(func() {
		repositories.NewMongoConnection(s.settings.Mongo.ConnectionString)
	})
}

func (s *MongoTestSuite) TestConnectFail() {
	s.assert.Panics(func() {
		repositories.NewMongoConnection("nohost")
	})
}
