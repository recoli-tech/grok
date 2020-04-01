package grok

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
)

// UploadKind ...
type UploadKind string

const (
	// GoogleStorageURL ...
	GoogleStorageURL = "https://storage.googleapis.com/%s/%s"
)

const (
	// PublicUpload ...
	PublicUpload = "public"
	// PrivateUpload ...
	PrivateUpload = "private"
)

// UploadProvider ...
type UploadProvider interface {
	New(ctx context.Context, filename string, kind UploadKind) io.WriteCloser
	FromRequest(ctx context.Context, file *multipart.FileHeader, kind UploadKind) (url string, err error)
	FromBase64(ctx context.Context, base64 string, kind UploadKind) (url string, err error)
}

// GCSUploadProvider ...
type GCSUploadProvider struct {
	bucketName string
	bucket     *storage.BucketHandle
}

// NewGCSUploadProvider ...
func NewGCSUploadProvider(client *storage.Client, projectID, bucket string) *GCSUploadProvider {
	gcs := new(GCSUploadProvider)
	gcs.ensureBucket(client, projectID, bucket)

	return gcs
}

// New ...
func (gcs *GCSUploadProvider) New(ctx context.Context, filename string, kind UploadKind) io.WriteCloser {
	o := gcs.bucket.Object(filename)
	w := o.NewWriter(ctx)

	switch kind {
	case PublicUpload:
		w.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}
	case PrivateUpload:
		w.ACL = []storage.ACLRule{{Entity: storage.AllAuthenticatedUsers, Role: storage.RoleReader}}
	}

	return w
}

// FromBase64 ...
func (gcs *GCSUploadProvider) FromBase64(ctx context.Context, value string, kind UploadKind) (url string, err error) {
	fileName := uuid.New().String()
	writer := gcs.New(ctx, fileName, kind)

	data, err := base64.StdEncoding.DecodeString(value)

	if err != nil {
		return "", err
	}

	if _, err := writer.Write(data); err != nil {
		return "", err
	}

	if err := writer.Close(); err != nil {
		return "", err
	}

	return gcs.locationURL(fileName), nil
}

// FromRequest ...
func (gcs *GCSUploadProvider) FromRequest(ctx context.Context, file *multipart.FileHeader, kind UploadKind) (string, error) {
	writer := gcs.New(ctx, file.Filename, kind)

	f, err := file.Open()

	if err != nil {
		return "", err
	}

	if _, err := io.Copy(writer, f); err != nil {
		return "", err
	}

	if err := writer.Close(); err != nil {
		return "", err
	}

	return gcs.locationURL(file.Filename), nil
}

func (gcs *GCSUploadProvider) locationURL(filename string) string {
	return fmt.Sprintf(GoogleStorageURL, gcs.bucketName, filename)
}

func (gcs *GCSUploadProvider) ensureBucket(client *storage.Client, projectID, bucket string) {
	b := client.Bucket(bucket)

	if _, err := b.Attrs(context.Background()); err != nil {
		if err != storage.ErrBucketNotExist {
			panic(err)
		}

		if err := b.Create(context.Background(), projectID, &storage.BucketAttrs{
			Name: bucket,
		}); err != nil {
			panic(err)
		}
	}

	gcs.bucket = b
	gcs.bucketName = bucket
}
