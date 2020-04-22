package grok_test

import (
	"context"
	"reflect"
	"testing"

	"cloud.google.com/go/pubsub"

	"github.com/recoli-tech/grok"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PubSubSubscriberTestSuite struct {
	suite.Suite
	assert   *assert.Assertions
	settings *grok.Settings
	client   *pubsub.Client
	producer *grok.PubSubProducer
}

func TestPubSubSubscriberTestSuite(t *testing.T) {
	suite.Run(t, new(PubSubSubscriberTestSuite))
}

func (s *PubSubSubscriberTestSuite) SetupTest() {
	s.assert = assert.New(s.T())
	s.settings = &grok.Settings{}
	grok.FromYAML("tests/config.yaml", s.settings)
	s.client = grok.FakePubSubClient(s.settings.GCP.PubSub.Endpoint)
	s.producer = grok.NewPubSubProducer(s.client)

}

func (s *PubSubSubscriberTestSuite) TestSubscribe() {
	ctx := context.Background()

	received := make(chan bool, 1)

	subscriberID := "subs"
	topicID := "topic"

	message := map[string]interface{}{"ping": "pong"}

	go func() {
		grok.NewPubSubSubscriber(
			grok.WithClient(s.client),
			grok.WithTopicID(topicID),
			grok.WithPubSubSubscriberID(subscriberID),
			grok.WithType(reflect.TypeOf(message)),
			grok.WithHandler(func(data interface{}) error {
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
