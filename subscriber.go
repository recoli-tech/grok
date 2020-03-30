package grok

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"cloud.google.com/go/pubsub"

	"github.com/sirupsen/logrus"
)

// PubSubSubscriber ...
type PubSubSubscriber struct {
	client                 *pubsub.Client
	handler                func(interface{}) error
	subscriberID           string
	topicID                string
	handleType             reflect.Type
	maxRetries             int
	producer               *PubSubProducer
	maxRetriesAttribute    string
	maxOutstandingMessages int
}

// PubSubSubscriberOption ...
type PubSubSubscriberOption func(*PubSubSubscriber)

// NewPubSubSubscriber ...
func NewPubSubSubscriber(opts ...PubSubSubscriberOption) *PubSubSubscriber {
	subscriber := new(PubSubSubscriber)
	subscriber.maxRetries = 5
	subscriber.maxOutstandingMessages = pubsub.DefaultReceiveSettings.MaxOutstandingMessages

	for _, opt := range opts {
		opt(subscriber)
	}

	subscriber.maxRetriesAttribute = "retries"
	subscriber.producer = NewPubSubProducer(subscriber.client)

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

// WithMaxRetries - default 5
func WithMaxRetries(maxRetries int) PubSubSubscriberOption {
	return func(s *PubSubSubscriber) {
		s.maxRetries = maxRetries
	}
}

//WitMaxOutstandingMessages ...
func WitMaxOutstandingMessages(maxOutstandingMessages int) PubSubSubscriberOption {
	return func(s *PubSubSubscriber) {
		s.maxOutstandingMessages = maxOutstandingMessages
	}
}

// Run ...
func (s *PubSubSubscriber) Run(ctx context.Context) error {
	subscriber, err := createSubscriptionIfNotExists(s.client, s.subscriberID, s.topicID)
	subscriber.ReceiveSettings.MaxOutstandingMessages = s.maxOutstandingMessages

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
				Errorf("cannot unmarshal message %s - sending to dlq", message.ID)

			s.dlq(message, err)

			message.Ack()
			return
		}

		err = s.handler(body)

		if err != nil {
			logrus.WithError(err).
				Errorf("error processing message %s", message.ID)

			switch s.getRetries(message) >= s.maxRetries {
			case true:
				if err := s.dlq(message, err); err != nil {
					logrus.WithError(err).
						Errorf("error sending message %s to dlq", message.ID)
				}
				break
			case false:
				if err := s.retry(message, body); err != nil {
					logrus.WithError(err).
						Errorf("error retrying message %s", message.ID)
				}
				break
			}
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

func (s *PubSubSubscriber) retry(message *pubsub.Message, body interface{}) error {
	retries := s.getRetries(message)
	retries++

	message.Attributes[s.maxRetriesAttribute] = strconv.Itoa(retries)

	return s.producer.PublishWihAttribrutes(s.topicID, body, message.Attributes)
}

func (s *PubSubSubscriber) dlq(message *pubsub.Message, e error) error {
	dlq := fmt.Sprintf("%s_dlq", s.topicID)

	logrus.Infof("sending message %s to %s", message.ID, dlq)

	_, err := createTopicIfNotExists(s.client, dlq)

	if err != nil {
		return err
	}

	attributes := make(map[string]string)
	attributes["error"] = e.Error()

	return s.producer.PublishWihAttribrutes(dlq, message.Data, attributes)
}

func (s *PubSubSubscriber) getRetries(message *pubsub.Message) int {
	if message.Attributes == nil {
		message.Attributes = make(map[string]string)
	}

	retries := 0
	attribute, ok := message.Attributes[s.maxRetriesAttribute]

	if ok {
		retries, _ = strconv.Atoi(attribute)
	}

	return retries
}
