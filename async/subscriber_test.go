package async_test

import (
	"context"
	"reflect"
	"testing"

	"cloud.google.com/go/pubsub"
	"github.com/raafvargas/grok/async"
	"github.com/raafvargas/grok/settings"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SubscriberTestSuite struct {
	suite.Suite
	assert   *assert.Assertions
	settings *settings.Settings
	client   *pubsub.Client
	producer *async.PubSubProducer
}

func TestSubscriberTestSuite(t *testing.T) {
	suite.Run(t, new(SubscriberTestSuite))
}

func (s *SubscriberTestSuite) SetupTest() {
	s.assert = assert.New(s.T())
	s.settings = &settings.Settings{}
	settings.FromYAML("../tests/config.yaml", s.settings)
	s.client = async.FakeClient(s.settings.GCP.PubSub.Endpoint)
	s.producer = async.NewPubSubProducer(s.client)
}

func (s *SubscriberTestSuite) TestSubscribe() {
	ctx := context.Background()

	received := make(chan bool, 1)

	subscriberID := "subs"
	topicID := "topic"

	message := map[string]interface{}{"ping": "pong"}

	go func() {
		async.NewSubscriber(
			async.WithClient(s.client),
			async.WithTopicID(topicID),
			async.WithSubscriberID(subscriberID),
			async.WithType(reflect.TypeOf(message)),
			async.WithHandler(func(data interface{}) error {
				defer func() { received <- true }()

				value, ok := data.(*map[string]interface{})
				s.assert.True(ok)
				s.assert.Equal("pong", (*value)["ping"])

				return nil
			}),
		).
			Run(ctx)
	}()

	err := s.producer.Publish(topicID, message)

	s.assert.NoError(err)

	<-received
}
