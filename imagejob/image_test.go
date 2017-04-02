package imagejob

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"

	"github.com/disintegration/imaging"
	"github.com/stretchr/testify/assert"
)

func TestDownload(t *testing.T) {
	img := Image{URL: testURL}
	body, err := img.Download(nil)
	assert.Nil(t, err)
	assert.NotNil(t, body)
}

func TestDownloadFailBecauseNoURL(t *testing.T) {
	img := Image{}
	body, err := img.Download(nil)
	assert.Nil(t, body)
	assert.Equal(t, err, errors.New("SourceURL not found in image"))
}

func TestDownloadFailBecauseBadURL(t *testing.T) {
	img := Image{URL: "fake"}
	body, err := img.Download(nil)
	assert.Nil(t, body)
	assert.Equal(t, err, errors.New("Cannot download image"))
}

func TestDecodeFail(t *testing.T) {
	img := Image{URL: "https://github.com"}
	body, _ := img.Download(nil)
	err := img.Decode(body)
	assert.NotNil(t, body)
	assert.Equal(t, err, errors.New("Cannot decode image"))
}

func TestDecode(t *testing.T) {
	img := Image{URL: testURL}
	body, _ := img.Download(nil)
	err := img.Decode(body)
	assert.Nil(t, err)
	assert.NotNil(t, body)
	assert.NotNil(t, img.Content, "Downloaded image should be not nil")
	assert.Equal(t, img.Height, 800)
	assert.Equal(t, img.Width, 566)
}

var EncodeCases = []struct {
	file    string
	format  string
	quality int
	err     error
}{
	{"testdata/fiveyears.jpg", "jpg", 120, nil},
	{"testdata/fiveyears.jpg", "jpg", 80, nil},
	{"testdata/fiveyears.jpg", "jpg", -1, nil},
	{"testdata/fiveyears.jpg", "jpeg", 80, nil},
	{"testdata/fiveyears.jpg", "gif", 80, nil},
	{"testdata/fiveyears.jpg", "png", 80, nil},
	{"testdata/fiveyears.jpg", "fake", 80, errors.New("Unsupported format")},
}

func TestEncode(t *testing.T) {
	for _, test := range EncodeCases {
		out, _ := ioutil.TempFile("/tmp/", "godinary")
		defer os.Remove(out.Name())

		img, _ := imaging.Open(test.file)

		err := Encode(img, out, test.format, test.quality)
		if test.err == nil {
			assert.Nil(t, err)
			_, err = os.Stat(out.Name())
			assert.Nil(t, err)
		} else {
			assert.Equal(t, test.err, err)
		}
	}

}
