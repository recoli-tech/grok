package grok

import (
	"context"
	"encoding/json"
	"reflect"

	"cloud.google.com/go/pubsub"

	"github.com/sirupsen/logrus"
)

// PubSubSubscriber ...
type PubSubSubscriber struct {
	client       *pubsub.Client
	handler      func(interface{}) error
	subscriberID string
	topicID      string
	handleType   reflect.Type
}

// PubSubSubscriberOption ...
type PubSubSubscriberOption func(*PubSubSubscriber)

// NewPubSubSubscriber ...
func NewPubSubSubscriber(opts ...PubSubSubscriberOption) *PubSubSubscriber {
	subscriber := new(PubSubSubscriber)

	for _, opt := range opts {
		opt(subscriber)
	}

	return subscriber
}

// WithClient ...
func WithClient(c *pubsub.Client) PubSubSubscriberOption {
	return func(s *PubSubSubscriber) {
		s.client = c
	}
}

// WithHandler ...
func WithHandler(h func(interface{}) error) PubSubSubscriberOption {
	return func(s *PubSubSubscriber) {
		s.handler = h
	}
}

// WithPubSubSubscriberID ...
func WithPubSubSubscriberID(id string) PubSubSubscriberOption {
	return func(s *PubSubSubscriber) {
		s.subscriberID = id
	}
}

// WithTopicID ...
func WithTopicID(t string) PubSubSubscriberOption {
	return func(s *PubSubSubscriber) {
		s.topicID = t
	}
}

// WithType ...
func WithType(t reflect.Type) PubSubSubscriberOption {
	return func(s *PubSubSubscriber) {
		s.handleType = t
	}
}

// Run ...
func (s *PubSubSubscriber) Run(ctx context.Context) error {
	subscriber, err := createSubscriptionIfNotExists(s.client, s.subscriberID, s.topicID)

	if err != nil {
		logrus.WithError(err).
			Errorf("error starting %s", s.subscriberID)
		return err
	}

	logrus.Infof("starting consumer %s with topic %s", s.subscriberID, s.topicID)
	return subscriber.Receive(ctx, func(c context.Context, message *pubsub.Message) {
		body := reflect.New(s.handleType).Interface()
		err := json.Unmarshal(message.Data, body)

		if err != nil {
			logrus.WithError(err).WithField("content", string(message.Data)).
				Errorf("cannot unmarshal message %s", message.ID)
			message.Nack()
			return
		}

		err = s.handler(body)

		if err != nil {
			logrus.WithError(err).
				Errorf("error processing message %s", message.ID)

			message.Nack()
		}

		message.Ack()
	})
}

func createSubscriptionIfNotExists(client *pubsub.Client, subscriberID, topicID string) (*pubsub.Subscription, error) {
	subscriber := client.Subscription(subscriberID)

	exists, err := subscriber.Exists(context.Background())

	if err != nil || exists {
		return subscriber, err
	}

	topic, err := createTopicIfNotExists(client, topicID)

	if err != nil {
		logrus.WithError(err).
			Errorf("error creating topic %s", topicID)
		return nil, err
	}

	subscriber, err = client.CreateSubscription(context.Background(), subscriberID, pubsub.SubscriptionConfig{
		Topic: topic,
	})

	if err != nil {
		logrus.WithError(err).
			Errorf("error creating subscription %s", subscriberID)
		return nil, err
	}
	return subscriber, nil
}
