package grok_test

import (
	"testing"

	"github.com/raafvargas/grok"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ProducerTestSuite struct {
	suite.Suite
	assert   *assert.Assertions
	settings *grok.Settings
}

func TestProducerTestSuite(t *testing.T) {
	suite.Run(t, new(ProducerTestSuite))
}

func (s *ProducerTestSuite) SetupTest() {
	s.assert = assert.New(s.T())
	s.settings = &grok.Settings{}
	grok.FromYAML("tests/config.yaml", s.settings)
}

func (s *ProducerTestSuite) TestPublish() {
	producer := grok.NewPubSubProducer(
		grok.FakePubSubClient(s.settings.GCP.PubSub.Endpoint))

	err := producer.Publish("test-topic", map[string]interface{}{"ping": "pong"})

	s.assert.NoError(err)
}
