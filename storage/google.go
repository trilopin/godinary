package storage

import (
	"io"

	"github.com/spf13/viper"

	"golang.org/x/net/context"
	"google.golang.org/api/option"

	gs "cloud.google.com/go/storage"
)

// GoogleStorageDriver struct
type GoogleStorageDriver struct {
	bucketName string
	bucket     *gs.BucketHandle
}

// NewGoogleStorageDriver constructs new GoogleStorageDriver
func NewGoogleStorageDriver() *GoogleStorageDriver {
	var gsw GoogleStorageDriver
	var err error
	var client *gs.Client

	gceProject := viper.GetString("gce_project")
	if gceProject == "" {
		panic("GODINARY_GCE_PROJECT should be setted")
	}

	gsw.bucketName = viper.GetString("gs_bucket")
	if gsw.bucketName == "" {
		panic("GODINARY_GS_BUCKET should be setted")
	}

	serviceAccount := viper.GetString("gs_credentials")

	ctx := context.Background()
	if serviceAccount == "" {
		client, err = gs.NewClient(ctx)
	} else {
		client, err = gs.NewClient(ctx, option.WithServiceAccountFile(serviceAccount))
	}
	if err != nil {
		panic("error in gstorage")
	}
	gsw.bucket = client.Bucket(gsw.bucketName)
	return &gsw
}

// Write in filesystem a bytearray
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
		// raven.CaptureError(err, nil) // it's called in a goroutine
		return nil, err
	}
	return rc, nil
}
