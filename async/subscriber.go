package async

import (
	"context"
	"encoding/json"
	"reflect"

	"cloud.google.com/go/pubsub"

	"github.com/sirupsen/logrus"
)

// Subscriber ...
type Subscriber struct {
	client       *pubsub.Client
	handler      func(interface{}) error
	subscriberID string
	topicID      string
	handleType   reflect.Type
}

// SubscriberOption ...
type SubscriberOption func(*Subscriber)

// NewSubscriber ...
func NewSubscriber(opts ...SubscriberOption) *Subscriber {
	subscriber := new(Subscriber)

	for _, opt := range opts {
		opt(subscriber)
	}

	return subscriber
}

// WithClient ...
func WithClient(c *pubsub.Client) SubscriberOption {
	return func(s *Subscriber) {
		s.client = c
	}
}

// WithHandler ...
func WithHandler(h func(interface{}) error) SubscriberOption {
	return func(s *Subscriber) {
		s.handler = h
	}
}

// WithSubscriberID ...
func WithSubscriberID(id string) SubscriberOption {
	return func(s *Subscriber) {
		s.subscriberID = id
	}
}

// WithTopicID ...
func WithTopicID(t string) SubscriberOption {
	return func(s *Subscriber) {
		s.topicID = t
	}
}

// WithType ...
func WithType(t reflect.Type) SubscriberOption {
	return func(s *Subscriber) {
		s.handleType = t
	}
}

// Run ...
func (s *Subscriber) Run(ctx context.Context) error {
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
