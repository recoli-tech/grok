package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/raafvargas/grok/container"
	"github.com/raafvargas/grok/middlewares"
	"github.com/raafvargas/grok/settings"

	"github.com/gin-gonic/gin"

	"github.com/sirupsen/logrus"
)

// Server wraps API configurations.
type Server struct {
	Engine *gin.Engine
	router *gin.RouterGroup

	cors      bool
	settings  *settings.Settings
	Container container.Container
}

// Option wrapps all server configurations
type Option func(server *Server)

func init() {
	gin.SetMode("release")
}

// WithContainer adds a container to the server
func WithContainer(c container.Container) Option {
	return func(server *Server) {
		server.Container = c
	}
}

// WithSettings sets server configurations
func WithSettings(settings *settings.Settings) Option {
	return func(server *Server) {
		server.settings = settings
	}
}

// WithCORS enables CORS
func WithCORS() Option {
	return func(server *Server) {
		server.cors = true
	}
}

// New creates a new API server
func New(opts ...Option) *Server {
	server := &Server{}

	for _, opt := range opts {
		opt(server)
	}

	server.Engine = gin.New()
	server.Engine.Use(gin.Recovery())
	server.Engine.Use(middlewares.Logging())

	if server.cors {
		server.Engine.Use(middlewares.CORS())
	}

	server.Engine.NoRoute(func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNotFound)
	})

	server.router = server.Engine.Group("")

	server.router.GET("/swagger", middlewares.Swagger(server.settings.API.Swagger))

	for _, ctrl := range server.Container.Controllers() {
		ctrl.RegisterRoutes(server.router)
	}

	return server
}

// Run starts the server.
func (server *Server) Run() {
	defer server.Container.Close()

	srv := http.Server{
		Addr:    server.settings.API.Host,
		Handler: server.Engine,
	}

	sigs := make(chan os.Signal)
	signal.Notify(sigs, os.Interrupt)

	go func() {
		sig := <-sigs

		logrus.Infof("caught sig: %+v", sig)
		logrus.Info("waiting 5 seconds to finish processing")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			logrus.WithField("error", err).Error("shotdown error")
		}
	}()

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logrus.WithField("error", err).Info("startup error")
	}
}
