package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/raafvargas/grok/controllers"
	"github.com/raafvargas/grok/server"
	"github.com/raafvargas/grok/settings"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ControllerTestSuite struct {
	suite.Suite
	assert   *assert.Assertions
	settings *settings.Settings
	server   *server.Server
}

type testContainer struct{}

func (c *testContainer) Controllers() []controllers.Controller {
	return nil
}

func (c *testContainer) Close() error {
	return nil
}

func TestControllerTestSuite(t *testing.T) {
	suite.Run(t, new(ControllerTestSuite))
}

func (s *ControllerTestSuite) SetupTest() {
	container := &testContainer{}
	s.assert = assert.New(s.T())
	s.settings = settings.FromYAML("../tests/config.yaml")
	s.server = server.New(
		server.WithSettings(s.settings),
		server.WithContainer(container))
}

func (s *ControllerTestSuite) TestNotFound() {
	req := httptest.NewRequest("GET", "/notfound", nil)
	response := httptest.NewRecorder()

	s.server.Engine.ServeHTTP(response, req)

	s.assert.Equal(http.StatusNotFound, response.Code)
}

func (s *ControllerTestSuite) TestSwagger() {
	req := httptest.NewRequest("GET", "/swagger", nil)
	response := httptest.NewRecorder()

	s.server.Engine.ServeHTTP(response, req)

	s.assert.Equal(http.StatusOK, response.Code)
}
