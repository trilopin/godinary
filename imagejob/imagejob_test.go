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
			Source: Image{
				URL: testURL,
			},
			Filters: map[string]string{"crop": "scale"},
		},
		"without filters",
	},
	{
		"w_400/" + testURL,
		ImageJob{
			Source: Image{
				URL: testURL,
			},
			Target: Image{
				Width: 400,
			},
			Filters: map[string]string{"crop": "scale"},
		},
		"with one filter",
	},
	{
		"w_400,c_limit,h_600,f_jpg/" + testURL,
		ImageJob{
			Source: Image{
				URL: testURL,
			},
			Target: Image{
				Width:  400,
				Height: 600,
				Format: "jpg",
			},
			Filters: map[string]string{"crop": "limit"},
		},
		"with multiple filter jpg",
	},
	{
		"w_400,c_limit,h_600,f_png/" + testURL,
		ImageJob{
			Source: Image{
				URL: testURL,
			},
			Target: Image{
				Width:  400,
				Height: 600,
				Format: "png",
			},
			Filters: map[string]string{"crop": "limit"},
		},
		"with multiple filter png",
	},
	{
		"w_400,c_limit,h_600,f_gif/" + testURL,
		ImageJob{
			Source: Image{
				URL: testURL,
			},
			Target: Image{
				Width:  400,
				Height: 600,
				Format: "gif",
			},
			Filters: map[string]string{"crop": "limit"},
		},
		"with multiple filter gif",
	},
	{
		"w_400,c_limit,h_600,f_jpeg/" + testURL,
		ImageJob{
			Source: Image{
				URL: testURL,
			},
			Target: Image{
				Width:  400,
				Height: 600,
				Format: "jpeg",
			},
			Filters: map[string]string{"crop": "limit"},
		},
		"with multiple filter jpeg",
	},
}

func TestParse(t *testing.T) {
	for _, test := range parserCases {
		job := NewImageJob()
		err := job.Parse(test.url)
		assert.Nil(t, err)
		assert.Equal(t, test.expected, *job, test.description)
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
	job := NewImageJob()
	job.Source.URL = testURL
	body, err := job.Source.Download()
	assert.Nil(t, err)
	assert.NotNil(t, body)
}

func TestDownloadFailBecauseNoURL(t *testing.T) {
	job := NewImageJob()
	body, err := job.Source.Download()
	assert.Nil(t, body)
	assert.Equal(t, err, errors.New("SourceURL not found in image"))
}

func TestDownloadFailBecauseBadURL(t *testing.T) {
	job := NewImageJob()
	job.Source.URL = "fake"
	body, err := job.Source.Download()
	assert.Nil(t, body)
	assert.Equal(t, err, errors.New("Cannot download image"))
}

func TestDecodeFail(t *testing.T) {
	job := NewImageJob()
	job.Source.URL = "https://github.com"
	body, _ := job.Source.Download()
	err := job.Source.Decode(body)
	assert.NotNil(t, body)
	assert.Equal(t, err, errors.New("Cannot decode image"))
}

func TestDecode(t *testing.T) {
	job := NewImageJob()
	job.Source.URL = testURL
	body, _ := job.Source.Download()
	err := job.Source.Decode(body)
	assert.Nil(t, err)
	assert.NotNil(t, body)
	assert.NotNil(t, job.Source.Content, "Downloaded image should be not nil")
	assert.Equal(t, job.Source.Height, 800)
	assert.Equal(t, job.Source.Width, 566)
}

func TestProcess(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	for _, test := range testFiles {
		img, _ := imaging.Open(test)
		out, _ := ioutil.TempFile("/tmp/", "godinary")

		job := NewImageJob()
		job.Source.Content = img
		job.Target.Width = 40
		job.Target.Height = 60
		job.Target.Format = "jpg"
		job.Source.extractInfo()

		err := job.Process(out)
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

	job := NewImageJob()
	job.Source.Content = img
	job.Target.Width = 60
	job.Target.Height = 40
	job.Target.Format = "jpg"
	job.Filters["crop"] = "fit"
	job.Source.extractInfo()

	err := job.Process(out)
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

	job := NewImageJob()
	job.Source.Content = img
	job.Target.Width = 6000
	job.Target.Height = 2000
	job.Target.Format = "jpg"
	job.Filters["crop"] = "limit"
	job.Source.extractInfo()

	err := job.Process(out)
	assert.Nil(t, err)

	resImg, _ := imaging.Open(out.Name())
	bounds := resImg.Bounds()
	assert.Equal(t, 1262, bounds.Max.X)
	assert.Equal(t, 733, bounds.Max.Y)
	os.Remove(out.Name())
}

func TestProcessFail(t *testing.T) {
	out, _ := ioutil.TempFile("/tmp/", "godinary")

	job := NewImageJob()
	job.Target.Width = 40
	job.Target.Height = 60
	job.Target.Format = "jpg"

	err := job.Process(out)
	assert.Nil(t, job.Source.Content)
	assert.Equal(t, err, errors.New("Image not found"))
}
