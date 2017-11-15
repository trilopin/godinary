package storage

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

var gsd *GoogleStorageDriver

func setup() {
	// create one file in GS
	GCEProject := "godinary"
	GSBucket := "godinary-tests"
	GSCredentials := "../godinary-e7f1383f7309.json"
	if os.Getenv("GODKEY") == "" {
		return
	}

	gsd := &GoogleStorageDriver{
		BucketName:  GSBucket,
		ProjectName: GCEProject,
		Credentials: GSCredentials,
	}
	err := gsd.Init()
	if err != nil {
		panic(err)
	}

	//write one file on aabbccddeeff
	ctx := context.Background()
	wc := gsd.bucket.Object("aa/bb/cc/dd/ee/aabbccddeeff").NewWriter(ctx)
	err = wc.Write([]byte("Content"))
	if err != nil {
		panic(err)
	}

}

func shutdown() {
	// drop files created in GS
}

func TestGoogleNewReaderWhenDoesNotExist(t *testing.T) {
	if gsd == nil {
		t.Skip()
	}
	r, err := gsd.NewReader("aabbccddeeff", "prefix/")
	assert.Nil(t, r)
	assert.NotNil(t, err)
}

func TestGoogleWrite(t *testing.T) {
	if gsd == nil {
		t.Skip()
	}
	content := []byte("content")
	err := gsd.Write(content, "aabbccddeeff", "prefix/")
	assert.Nil(t, err)
}

func TestGoogleNewReaderWhenExist(t *testing.T) {
	if gsd == nil {
		t.Skip()
	}
	r, err := gsd.NewReader("aabbccddeeff", "prefix/")
	assert.Nil(t, r)
	assert.NotNil(t, err)
}

func TestMain(m *testing.M) {
	setup()
	exitCode := m.Run()
	shutdown()
	os.Exit(exitCode)
}
