package image

import (
	"fmt"
	"os"
	"testing"

	bimg "gopkg.in/h2non/bimg.v1"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	img := Image{}
	r, _ := os.Open("testdata/fiveyears.jpg")
	img.Load(r)
	assert.NotNil(t, img.Content)
}

func TestExtractInfo(t *testing.T) {
	img := Image{}
	r, _ := os.Open("testdata/fiveyears.jpg")
	img.Load(r)
	err := img.ExtractInfo()
	assert.Nil(t, err)
	assert.Equal(t, img.Height, 733)
	assert.Equal(t, img.Width, 1262)
	assert.Equal(t, img.AspectRatio, float32(1.7216917))
}

func TestExtractInfoFail(t *testing.T) {
	img := Image{}
	r, _ := os.Open("testdata/corrupted-image.jpg.jpg")
	img.Load(r)
	err := img.ExtractInfo()
	assert.Equal(t, err, fmt.Errorf("can't extract dimensions: Unsupported image format"))

}

func TestDownload(t *testing.T) {
	img := Image{URL: testURL}
	err := img.Download(nil)
	assert.Nil(t, err)
}

func TestDownloadFailBecauseNoURL(t *testing.T) {
	img := Image{}
	err := img.Download(nil)
	assert.Equal(t, err, fmt.Errorf("sourceURL not found in image"))
}

func TestDownloadFailBecauseBadURL(t *testing.T) {
	img := Image{URL: "fake"}
	err := img.Download(nil)
	assert.Equal(t, err, fmt.Errorf("cannot download image fake: Get fake: unsupported protocol scheme \"\""))
}

func TestProcess(t *testing.T) {
	source := Image{}
	img := Image{Width: 300, Height: 400, Format: bimg.WEBP, Quality: 75}
	r, _ := os.Open("testdata/fiveyears.jpg")
	source.Load(r)

	img.Process(source, nil)
	newImage := bimg.NewImage(img.RawContent)
	size, _ := newImage.Size()
	assert.Equal(t, size.Height, 400)
	assert.Equal(t, size.Width, 300)
}
