package imagejob

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testURL = "http://upload.wikimedia.org/wikipedia/commons/0/0c/Scarlett_Johansson_Césars_2014.jpg"

var parserCases = []struct {
	url         string
	expected    ImageJob
	description string
}{
	{
		testURL,
		ImageJob{
			SourceURL:  testURL,
			Filters:    map[string]string{},
			FixedRatio: false,
		},
		"without filters",
	},
	{
		"w_400/" + testURL,
		ImageJob{
			SourceURL:   testURL,
			Filters:     map[string]string{},
			TargetWidth: 400,
			FixedRatio:  false,
		},
		"with one filter",
	},
	{
		"w_400,c_limit,h_600,f_jpg/" + testURL,
		ImageJob{
			SourceURL:    testURL,
			Filters:      map[string]string{},
			TargetWidth:  400,
			TargetHeight: 600,
			TargetFormat: "jpg",
			FixedRatio:   true,
		},
		"with multiple filter jpg",
	},
	{
		"w_400,c_limit,h_600,f_png/" + testURL,
		ImageJob{
			SourceURL:    testURL,
			Filters:      map[string]string{},
			TargetWidth:  400,
			TargetHeight: 600,
			TargetFormat: "png",
			FixedRatio:   true,
		},
		"with multiple filter png",
	},
	{
		"w_400,c_limit,h_600,f_gif/" + testURL,
		ImageJob{
			SourceURL:    testURL,
			Filters:      map[string]string{},
			TargetWidth:  400,
			TargetHeight: 600,
			TargetFormat: "gif",
			FixedRatio:   true,
		},
		"with multiple filter gif",
	},
	{
		"w_400,c_limit,h_600,f_jpeg/" + testURL,
		ImageJob{
			SourceURL:    testURL,
			Filters:      map[string]string{},
			TargetWidth:  400,
			TargetHeight: 600,
			TargetFormat: "jpeg",
			FixedRatio:   true,
		},
		"with multiple filter jpeg",
	},
}

func TestNew(t *testing.T) {
	for _, test := range parserCases {
		img := ImageJob{}
		err := img.New(test.url)
		assert.Nil(t, err)
		assert.Equal(t, test.expected, img, test.description)
	}
}

var parserErrorCases = []struct {
	url         string
	err         error
	description string
}{
	{
		"w_400,c_limit,h_600,f_fake/" + testURL,
		errors.New("Format not allowed"),
		"Bad image format",
	},
	{
		"w_400,c_limit,h_pp/" + testURL,
		errors.New("TargetHeight is not integer"),
		"TargetHeight is not integer",
	},
	{
		"w_OOO,c_limit,h_500/" + testURL,
		errors.New("TargetWidth is not integer"),
		"TargetWidth is not integer",
	},
}

func TestNewWithError(t *testing.T) {
	for _, test := range parserErrorCases {
		img := ImageJob{}
		err := img.New(test.url)
		assert.Equal(t, test.err, err, test.description)
	}
}

func TestDownload(t *testing.T) {
	img := ImageJob{SourceURL: testURL}
	err := img.Download()
	assert.Nil(t, err)
	assert.NotNil(t, img.Image, "Downloaded image should be not nil")
	assert.Equal(t, img.SourceHeight, 800)
	assert.Equal(t, img.SourceWidth, 566)
}

func TestDownloadFailBecauseNoURL(t *testing.T) {
	img := ImageJob{}
	err := img.Download()
	assert.Nil(t, img.Image)
	assert.Equal(t, err, errors.New("SourceURL not found in image"))
	assert.NotNil(t, err)
}

func TestDownloadFailBecauseBadURL(t *testing.T) {
	img := ImageJob{SourceURL: "fake"}
	err := img.Download()
	assert.Nil(t, img.Image)
	assert.Equal(t, err, errors.New("Cannot download image"))
	assert.NotNil(t, err)
}

func TestDownloadFailBecauseNoImage(t *testing.T) {
	img := ImageJob{SourceURL: "https://github.com"}
	err := img.Download()
	assert.Nil(t, img.Image)
	assert.Equal(t, err, errors.New("Cannot decode image"))
	assert.NotNil(t, err)
}
