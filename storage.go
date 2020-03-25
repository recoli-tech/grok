package grok

import (
	"context"
	"crypto/tls"
	"net/http"

	"cloud.google.com/go/storage"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
)

// CreateStorageClient ...
func CreateStorageClient(settings *GCPSettings) *storage.Client {
	switch {
	case settings.Storage.Fake:
		return FakeStorageClient(settings)
	default:
		client, err := storage.NewClient(context.Background())
		if err != nil {
			logrus.WithError(err).Fatal("error creating storage client")
		}
		return client
	}
}

// FakeStorageClient ...
func FakeStorageClient(settings *GCPSettings) *storage.Client {
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client, _ := storage.NewClient(
		context.Background(),
		option.WithEndpoint(settings.Storage.Endpoint),
		option.WithHTTPClient(&http.Client{Transport: transCfg}),
	)

	return client
}
