package async

import (
	"context"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
)

// FakeClient ...
func FakeClient(endpoint string) *pubsub.Client {
	conn, _ := grpc.Dial(endpoint, grpc.WithInsecure())

	client, _ := pubsub.NewClient(
		context.Background(),
		"fake_client",
		option.WithGRPCConn(conn))

	return client
}
