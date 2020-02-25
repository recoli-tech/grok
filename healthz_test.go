package grok_test

import (
	"testing"

	"github.com/raafvargas/grok"
	"github.com/stretchr/testify/assert"
)

func TestHealthz(t *testing.T) {
	settings := &grok.Settings{}
	grok.FromYAML("tests/config.yaml", settings)

	t.Run("Mongo Success", func(t *testing.T) {
		healthz := grok.NewHealthz(
			grok.WithMongo(),
			grok.WithHealthzSettings(settings))

		err := healthz.Healthz()

		assert.NoError(t, err)
	})
}
