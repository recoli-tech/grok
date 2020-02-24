package grok

import (
	"context"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Healthz ...
type Healthz struct {
	client *mongo.Client
	checks []func(*Healthz) error
}

// HealtzOption ...
type HealtzOption func(*Healthz)

// WithMongo ...
func WithMongo(client *mongo.Client) HealtzOption {
	return func(h *Healthz) {
		h.client = client
		h.checks = append(h.checks, func(healthz *Healthz) error {
			return healthz.client.Ping(context.Background(), readpref.Primary())
		})
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
