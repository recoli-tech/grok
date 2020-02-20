package grok

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

// MockedResponseResult ...
type MockedResponseResult struct {
	Body    string
	Status  int
	Headers http.Header
}

// HTTPClientMock ...
type HTTPClientMock struct {
	client    *http.Client
	history   []string
	result    map[string]*MockedResponseResult
	roundTrip RoundTripFunc
}

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

// Client ...
func (h *HTTPClientMock) Client() *http.Client {
	return h.client
}

// AddMock ...
func (h *HTTPClientMock) AddMock(uri string, result *MockedResponseResult) {
	WithMock(uri, result)(h)
}

// History ...
func (h *HTTPClientMock) History() []string {
	return h.history
}

// NewHTTPClientMock ...
func NewHTTPClientMock(options ...HTTPMockOption) *HTTPClientMock {
	mock := &HTTPClientMock{}
	mock.result = make(map[string]*MockedResponseResult)
	mock.roundTrip = mock.defaultRoundTrip

	for _, opt := range options {
		opt(mock)
	}

	mock.client = &http.Client{
		Transport: mock.roundTrip}

	return mock
}

// HTTPMockOption ...
type HTTPMockOption func(*HTTPClientMock)

// WithMock ...
func WithMock(url string, result *MockedResponseResult) HTTPMockOption {
	return func(mock *HTTPClientMock) {
		if result.Headers == nil {
			result.Headers = make(http.Header)
		}

		mock.result[url] = result
	}
}

// WithRoundTrip ...
func WithRoundTrip(roundTrip RoundTripFunc) HTTPMockOption {
	return func(mock *HTTPClientMock) {
		mock.roundTrip = roundTrip
	}
}

func (h *HTTPClientMock) defaultRoundTrip(req *http.Request) *http.Response {
	uri := req.URL.RequestURI()

	h.history = append(h.history, uri)

	result, ok := h.result[uri]

	if !ok {
		result = &MockedResponseResult{}
		result.Body = "OK"
		result.Status = http.StatusOK
		result.Headers = make(http.Header)
	}

	return &http.Response{
		StatusCode: result.Status,
		Body:       ioutil.NopCloser(bytes.NewBufferString(result.Body)),
		Header:     result.Headers,
	}
}
