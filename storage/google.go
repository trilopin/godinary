package storage

import (
	"io"
	"os"

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

	gceProject := os.Getenv("GODINARY_GCE_PROJECT")
	if gceProject == "" {
		panic("GODINARY_GCE_PROJECT should be setted")
	}

	gsw.bucketName = os.Getenv("GODINARY_GS_BUCKET")
	if gsw.bucketName == "" {
		panic("GODINARY_GS_BUCKET should be setted")
	}

	serviceAccount := os.Getenv("GODINARY_GS_CREDENTIALS")
	if serviceAccount == "" {
		panic("GODINARY_GS_CREDENTIALS should be setted")
	}

	ctx := context.Background()
	client, err := gs.NewClient(ctx, option.WithServiceAccountFile(serviceAccount))
	if err != nil {
		panic("error in gstorage")
	}
	gsw.bucket = client.Bucket(gsw.bucketName)
	return &gsw
}

// Write in filesystem a bytearray
func (gsw *GoogleStorageDriver) Write(buf []byte, hash string) error {
	ctx := context.Background()
	_, newHash := makeFoldersFromHash(hash, "", 5)
	wc := gsw.bucket.Object(newHash).NewWriter(ctx)
	defer wc.Close()
	if _, err := wc.Write(buf); err != nil {
		return err
	}
	return nil
}

func (gsw *GoogleStorageDriver) NewReader(hash string) (io.ReadCloser, error) {
	ctx := context.Background()
	_, newHash := makeFoldersFromHash(hash, "", 5)
	rc, err := gsw.bucket.Object(newHash).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	return rc, nil
}
