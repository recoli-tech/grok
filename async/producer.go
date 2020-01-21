package async

import (
	"context"
	"encoding/json"
	"time"

	"cloud.google.com/go/pubsub"
)

// PubSubProducer ...
type PubSubProducer struct {
	client *pubsub.Client
}

// NewPubSubProducer ...
func NewPubSubProducer(client *pubsub.Client) *PubSubProducer {
	return &PubSubProducer{client: client}
}

// Publish ...
func (p *PubSubProducer) Publish(topicID string, data interface{}) error {
	body, err := json.Marshal(data)

	if err != nil {
		return err
	}

	topic, err := createTopicIfNotExists(p.client, topicID)

	if err != nil {
		return err
	}

	_, err = topic.
		Publish(context.Background(), &pubsub.Message{
			Data:        body,
			PublishTime: time.Now(),
		}).
		Get(context.Background())

	return err
}

func createTopicIfNotExists(client *pubsub.Client, id string) (*pubsub.Topic, error) {
	topic := client.Topic(id)
	exists, _ := topic.Exists(context.Background())

	if exists {
		return topic, nil
	}

	return client.CreateTopic(context.Background(), id)
}
