package grok

import (
	"context"

	"cloud.google.com/go/pubsub"
	"github.com/sirupsen/logrus"
)

// CreatePubSubClient ...
func CreatePubSubClient(settings *GCPSettings) *pubsub.Client {
	switch {
	case settings.PubSub.Fake:
		return FakePubSubClient(settings.PubSub.Endpoint)
	default:
		pubsub, err := pubsub.NewClient(context.Background(), settings.ProjectID)
		if err != nil {
			logrus.WithError(err).Fatal("error creating pubsub client")
		}
		return pubsub
	}
}
