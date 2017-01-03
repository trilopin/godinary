package imagejob

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/disintegration/imaging"
	"github.com/stretchr/testify/assert"
)

const testURL = "http://upload.wikimedia.org/wikipedia/commons/0/0c/Scarlett_Johansson_Césars_2014.jpg"

var testFiles = map[string]string{
	"jpg":  "testdata/fiveyears.jpg",
	"jpeg": "testdata/highres_315125112.jpeg",
	"gif":  "testdata/Jake_ma_kanapke.gif",
	"png":  "testdata/baboon.png",
}

var parserCases = []struct {
	url         string
	expected    ImageJob
	description string
}{
	{
		testURL,
		ImageJob{
			SourceURL: testURL,
			Filters:   map[string]string{"crop": "scale"},
		},
		"without filters",
	},
	{
		"w_400/" + testURL,
		ImageJob{
			SourceURL:   testURL,
			Filters:     map[string]string{"crop": "scale"},
			TargetWidth: 400,
		},
		"with one filter",
	},
	{
		"w_400,c_limit,h_600,f_jpg/" + testURL,
		ImageJob{
			SourceURL:    testURL,
			Filters:      map[string]string{"crop": "limit"},
			TargetWidth:  400,
			TargetHeight: 600,
			TargetFormat: "jpg",
		},
		"with multiple filter jpg",
	},
	{
		"w_400,c_limit,h_600,f_png/" + testURL,
		ImageJob{
			SourceURL:    testURL,
			Filters:      map[string]string{"crop": "limit"},
			TargetWidth:  400,
			TargetHeight: 600,
			TargetFormat: "png",
		},
		"with multiple filter png",
	},
	{
		"w_400,c_limit,h_600,f_gif/" + testURL,
		ImageJob{
			SourceURL:    testURL,
			Filters:      map[string]string{"crop": "limit"},
			TargetWidth:  400,
			TargetHeight: 600,
			TargetFormat: "gif",
		},
		"with multiple filter gif",
	},
	{
		"w_400,c_limit,h_600,f_jpeg/" + testURL,
		ImageJob{
			SourceURL:    testURL,
			Filters:      map[string]string{"crop": "limit"},
			TargetWidth:  400,
			TargetHeight: 600,
			TargetFormat: "jpeg",
		},
		"with multiple filter jpeg",
	},
}

func TestParse(t *testing.T) {
	for _, test := range parserCases {
		img := NewImageJob()
		err := img.Parse(test.url)
		assert.Nil(t, err)
		assert.Equal(t, test.expected, *img, test.description)
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
	{
		"w_100,c_fake,h_500/" + testURL,
		errors.New("Crop not allowed"),
		"Crop is not allowed",
	},
}

func TestParseFail(t *testing.T) {
	for _, test := range parserErrorCases {
		img := NewImageJob()
		err := img.Parse(test.url)
		assert.Equal(t, test.err, err, test.description)
	}
}

func TestDownload(t *testing.T) {
	img := NewImageJob()
	img.SourceURL = testURL
	body, err := Download(img)
	assert.Nil(t, err)
	assert.NotNil(t, body)
}

func TestDownloadFailBecauseNoURL(t *testing.T) {
	img := NewImageJob()
	body, err := Download(img)
	assert.Nil(t, body)
	assert.Equal(t, err, errors.New("SourceURL not found in image"))
}

func TestDownloadFailBecauseBadURL(t *testing.T) {
	img := NewImageJob()
	img.SourceURL = "fake"
	body, err := Download(img)
	assert.Nil(t, body)
	assert.Equal(t, err, errors.New("Cannot download image"))
}

func TestDecodeFail(t *testing.T) {
	img := NewImageJob()
	img.SourceURL = "https://github.com"
	body, _ := Download(img)
	err := img.Decode(body)
	assert.Nil(t, img.Image)
	assert.NotNil(t, body)
	assert.Equal(t, err, errors.New("Cannot decode image"))
}

func TestDecode(t *testing.T) {
	img := NewImageJob()
	img.SourceURL = testURL
	body, _ := Download(img)
	err := img.Decode(body)
	assert.Nil(t, err)
	assert.NotNil(t, body)
	assert.NotNil(t, img.Image, "Downloaded image should be not nil")
	assert.Equal(t, img.SourceHeight, 800)
	assert.Equal(t, img.SourceWidth, 566)
}

func TestProcess(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	for _, test := range testFiles {
		img, _ := imaging.Open(test)
		out, _ := ioutil.TempFile("/tmp/", "godinary")

		image := NewImageJob()
		image.Image = img
		image.TargetWidth = 40
		image.TargetHeight = 60
		image.TargetFormat = "jpg"
		image.extractInfo()

		err := image.Process(out)
		assert.Nil(t, err)

		resImg, _ := imaging.Open(out.Name())
		bounds := resImg.Bounds()
		assert.Equal(t, bounds.Max.Y, 60)
		assert.Equal(t, bounds.Max.X, 40)
		os.Remove(out.Name())
	}
}

func TestProcessFitHorizontal(t *testing.T) {
	img, _ := imaging.Open(testFiles["jpg"])
	out, _ := ioutil.TempFile("/tmp/", "godinary")

	image := NewImageJob()
	image.Image = img
	image.TargetWidth = 60
	image.TargetHeight = 40
	image.TargetFormat = "jpg"
	image.Filters["crop"] = "fit"
	image.extractInfo()

	err := image.Process(out)
	assert.Nil(t, err)

	resImg, _ := imaging.Open(out.Name())
	bounds := resImg.Bounds()
	assert.Equal(t, 60, bounds.Max.X)
	assert.Equal(t, 35, bounds.Max.Y)
	os.Remove(out.Name())
}

func TestProcessLimitHorizontal(t *testing.T) {
	img, _ := imaging.Open(testFiles["jpg"])
	out, _ := ioutil.TempFile("/tmp/", "godinary")

	image := NewImageJob()
	image.Image = img
	image.TargetWidth = 6000
	image.TargetHeight = 2000
	image.TargetFormat = "jpg"
	image.Filters["crop"] = "limit"
	image.extractInfo()

	err := image.Process(out)
	assert.Nil(t, err)

	resImg, _ := imaging.Open(out.Name())
	bounds := resImg.Bounds()
	assert.Equal(t, 1262, bounds.Max.X)
	assert.Equal(t, 733, bounds.Max.Y)
	os.Remove(out.Name())
}

func TestProcessFail(t *testing.T) {
	out, _ := ioutil.TempFile("/tmp/", "godinary")

	image := NewImageJob()
	image.TargetWidth = 40
	image.TargetHeight = 60
	image.TargetFormat = "jpg"

	err := image.Process(out)
	assert.Nil(t, image.Image)
	assert.Equal(t, err, errors.New("Image not found"))
}
