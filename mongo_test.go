package grok_test

import (
	"testing"

	"github.com/recoli-tech/grok"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MongoTestSuite struct {
	suite.Suite
	assert   *assert.Assertions
	settings *grok.Settings
}

func TestMongoTestSuite(t *testing.T) {
	suite.Run(t, new(MongoTestSuite))
}

func (s *MongoTestSuite) SetupSuite() {
	s.assert = assert.New(s.T())
	s.settings = &grok.Settings{}
	grok.FromYAML("tests/config.yaml", s.settings)
}

func (s *MongoTestSuite) TestConnect() {
	s.assert.NotPanics(func() {
		grok.NewMongoConnection(s.settings.Mongo.ConnectionString)
	})
}

func (s *MongoTestSuite) TestConnectFail() {
	s.assert.Panics(func() {
		grok.NewMongoConnection("nohost")
	})
}
