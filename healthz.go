package grok

import (
	"context"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Healthz ...
type Healthz struct {
	settings *Settings
	checks   []func(*Healthz) error
}

// HealtzOption ...
type HealtzOption func(*Healthz)

// WithMongo ...
func WithMongo() HealtzOption {
	return func(h *Healthz) {
		h.checks = append(h.checks, func(healthz *Healthz) error {
			client := NewMongoConnection(h.settings.Mongo.ConnectionString)
			defer client.Disconnect(context.Background())

			return client.Ping(context.Background(), readpref.Primary())
		})
	}
}

// WithHealthzSettings ...
func WithHealthzSettings(s *Settings) HealtzOption {
	return func(h *Healthz) {
		h.settings = s
	}
}

// NewHealthz ...
func NewHealthz(options ...HealtzOption) *Healthz {
	h := new(Healthz)
	h.checks = []func(*Healthz) error{}

	for _, o := range options {
		o(h)
	}

	return h
}

// HTTPHealthz ...
func HTTPHealthz(options ...HealtzOption) gin.HandlerFunc {
	h := NewHealthz(options...)
	return h.HTTP()
}

// ConsumerHealthz ...
func ConsumerHealthz(settingsFlag string, options ...HealtzOption) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		settings := &struct {
			Grok *Settings `yaml:"grok"`
		}{}
		err := FromYAML(cmd.Flag(settingsFlag).Value.String(), settings)

		if err != nil {
			logrus.WithError(err).
				Panic("error loading settings")
		}

		options = append(options, WithHealthzSettings(settings.Grok))

		h := NewHealthz(options...)

		if err := h.Healthz(); err != nil {
			logrus.WithError(err).
				Panic("health checke failed")
		}
	}
}

// Healthz ...
func (h *Healthz) Healthz() error {
	wg := new(sync.WaitGroup)

	errCh := make(chan error, len(h.checks))
	doneCh := make(chan bool, len(h.checks))

	for _, check := range h.checks {
		wg.Add(1)
		go func(c func(*Healthz) error) {
			defer wg.Done()
			if err := c(h); err != nil {
				errCh <- err
			}
		}(check)
	}

	go func() {
		wg.Wait()
		doneCh <- true
	}()

	<-doneCh

	close(errCh)
	close(doneCh)

	if len(errCh) > 0 {
		return <-errCh
	}

	return nil
}

// HTTP ...
func (h *Healthz) HTTP() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := h.Healthz(); err != nil {
			ctx.Error(err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ctx.Status(http.StatusOK)
	}
}
