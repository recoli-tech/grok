package grok_test

import (
	"context"
	"testing"

	"github.com/raafvargas/grok"
	"github.com/stretchr/testify/assert"
)

func TestHealthz(t *testing.T) {
	settings := &grok.Settings{}
	grok.FromYAML("tests/config.yaml", settings)

	t.Run("Mongo Success", func(t *testing.T) {
		client := grok.NewMongoConnection(settings.Mongo.ConnectionString)

		healthz := grok.NewHealthz(
			grok.WithMongo(client))

		err := healthz.Healthz()

		assert.NoError(t, err)
	})

	t.Run("Mongo Error", func(t *testing.T) {
		client := grok.NewMongoConnection(settings.Mongo.ConnectionString)

		healthz := grok.NewHealthz(
			grok.WithMongo(client))

		err := client.Disconnect(context.Background())

		assert.NoError(t, err)

		err = healthz.Healthz()

		assert.Error(t, err)
	})
}
