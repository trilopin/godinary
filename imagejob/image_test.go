package imagejob

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDownload(t *testing.T) {
	img := Image{URL: testURL}
	err := img.Download(nil)
	assert.Nil(t, err)
}

func TestDownloadFailBecauseNoURL(t *testing.T) {
	img := Image{}
	err := img.Download(nil)
	assert.Equal(t, err, errors.New("SourceURL not found in image"))
}

func TestDownloadFailBecauseBadURL(t *testing.T) {
	img := Image{URL: "fake"}
	err := img.Download(nil)
	assert.Equal(t, err, errors.New("Cannot download image"))
}
