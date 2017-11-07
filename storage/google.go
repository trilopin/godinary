package storage

import (
	"io"

	"golang.org/x/net/context"
	"google.golang.org/api/option"

	gs "cloud.google.com/go/storage"
)

// GoogleStorageDriver struct
type GoogleStorageDriver struct {
	BucketName  string
	ProjectName string
	Credentials string
	bucket      *gs.BucketHandle
}

// Init initialises bucket connection
func (gsw *GoogleStorageDriver) Init() error {
	var client *gs.Client
	var err error
	ctx := context.Background()
	if gsw.Credentials == "" {
		client, err = gs.NewClient(ctx)
	} else {
		client, err = gs.NewClient(ctx, option.WithServiceAccountFile(gsw.Credentials))
	}
	if err != nil {
		return err
	}
	gsw.bucket = client.Bucket(gsw.BucketName)
	return nil
}

// Write in Google storage a bytearray
func (gsw *GoogleStorageDriver) Write(buf []byte, hash string, prefix string) error {
	ctx := context.Background()
	_, newHash := makeFoldersFromHash(hash, prefix, 5)
	wc := gsw.bucket.Object(newHash).NewWriter(ctx)
	defer wc.Close()
	if _, err := wc.Write(buf); err != nil {
		return err
	}
	return nil
}

// NewReader produces a handler for file in google storage
func (gsw *GoogleStorageDriver) NewReader(hash string, prefix string) (io.ReadCloser, error) {
	ctx := context.Background()
	_, newHash := makeFoldersFromHash(hash, prefix, 5)
	rc, err := gsw.bucket.Object(newHash).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	return rc, nil
}
