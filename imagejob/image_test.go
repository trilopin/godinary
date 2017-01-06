package imagejob

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDownload(t *testing.T) {
	img := Image{URL: testURL}
	body, err := img.Download()
	assert.Nil(t, err)
	assert.NotNil(t, body)
}

func TestDownloadFailBecauseNoURL(t *testing.T) {
	img := Image{}
	body, err := img.Download()
	assert.Nil(t, body)
	assert.Equal(t, err, errors.New("SourceURL not found in image"))
}

func TestDownloadFailBecauseBadURL(t *testing.T) {
	img := Image{URL: "fake"}
	body, err := img.Download()
	assert.Nil(t, body)
	assert.Equal(t, err, errors.New("Cannot download image"))
}

func TestDecodeFail(t *testing.T) {
	img := Image{URL: "https://github.com"}
	body, _ := img.Download()
	err := img.Decode(body)
	assert.NotNil(t, body)
	assert.Equal(t, err, errors.New("Cannot decode image"))
}

func TestDecode(t *testing.T) {
	img := Image{URL: testURL}
	body, _ := img.Download()
	err := img.Decode(body)
	assert.Nil(t, err)
	assert.NotNil(t, body)
	assert.NotNil(t, img.Content, "Downloaded image should be not nil")
	assert.Equal(t, img.Height, 800)
	assert.Equal(t, img.Width, 566)
}
