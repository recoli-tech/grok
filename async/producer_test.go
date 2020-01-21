package async_test

import (
	"testing"

	"github.com/raafvargas/grok/async"
	"github.com/raafvargas/grok/settings"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ProducerTestSuite struct {
	suite.Suite
	assert   *assert.Assertions
	settings *settings.Settings
}

func TestProducerTestSuite(t *testing.T) {
	suite.Run(t, new(ProducerTestSuite))
}

func (s *ProducerTestSuite) SetupTest() {
	s.assert = assert.New(s.T())
	s.settings = settings.FromYAML("../tests/config.yaml")
}

func (s *ProducerTestSuite) TestPublish() {
	producer := async.NewPubSubProducer(
		async.FakeClient(s.settings.GCP.PubSub.Endpoint))

	err := producer.Publish("test-topic", map[string]interface{}{"ping": "pong"})

	s.assert.NoError(err)
}
