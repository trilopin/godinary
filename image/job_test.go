package image

import (
	"errors"
	"fmt"
	"testing"

	bimg "gopkg.in/h2non/bimg.v1"

	"github.com/stretchr/testify/assert"
)

const testURL = "http://upload.wikimedia.org/wikipedia/commons/0/0c/Scarlett_Johansson_Césars_2014.jpg"
const testSecureURL = "https://upload.wikimedia.org/wikipedia/commons/0/0c/Scarlett_Johansson_Césars_2014.jpg"

var testFiles = map[string]string{
	"jpg":  "testdata/fiveyears.jpg",
	"jpeg": "testdata/highres_315125112.jpeg",
	"gif":  "testdata/Jake_ma_kanapke.gif",
	"png":  "testdata/baboon.png",
}

type FakeHasher struct{}

func (fh *FakeHasher) Hash(s string) string {
	return ""
}

func TestParse(t *testing.T) {
	fakeHasher := &FakeHasher{}
	cases := []struct {
		url         string
		expected    Job
		description string
	}{
		{
			testURL,
			Job{
				Source:  Image{URL: testURL, Hash: ""},
				Target:  Image{Format: bimg.JPEG, Hash: ""},
				Filters: map[string]string{"crop": "scale"},
				Hasher:  fakeHasher,
			},
			"without filters",
		},
		{
			testSecureURL,
			Job{
				Source:  Image{URL: testSecureURL, Hash: ""},
				Target:  Image{Format: bimg.JPEG, Hash: ""},
				Filters: map[string]string{"crop": "scale"},
				Hasher:  fakeHasher,
			},
			"without filters secure",
		},
		{
			"w_400/" + testURL,
			Job{
				Source:  Image{URL: testURL, Hash: ""},
				Target:  Image{Width: 400, Format: bimg.JPEG, Hash: ""},
				Filters: map[string]string{"crop": "scale"},
				Hasher:  fakeHasher,
			},
			"with one filter",
		},
		{
			"w_400/" + testSecureURL,
			Job{
				Source:  Image{URL: testSecureURL, Hash: ""},
				Target:  Image{Width: 400, Format: bimg.JPEG, Hash: ""},
				Filters: map[string]string{"crop": "scale"},
				Hasher:  fakeHasher,
			},
			"with one filter secure",
		},
		{
			"w_400,c_limit,h_600,f_jpg/" + testURL,
			Job{
				Source:  Image{URL: testURL, Hash: ""},
				Target:  Image{Width: 400, Height: 600, Format: bimg.JPEG, Hash: ""},
				Filters: map[string]string{"crop": "limit"},
				Hasher:  fakeHasher,
			},
			"with multiple filter jpg",
		},
		{
			"w_400,c_limit,h_600,f_webp,q_50/" + testURL,
			Job{
				Source:  Image{URL: testURL, Hash: ""},
				Target:  Image{Width: 400, Height: 600, Quality: 50, Format: bimg.WEBP, Hash: ""},
				Filters: map[string]string{"crop": "limit"},
				Hasher:  fakeHasher,
			},
			"with multiple filter webp",
		},
		{
			"w_400,c_limit,h_600,f_auto,q_65/" + testURL,
			Job{
				Source:  Image{URL: testURL, Hash: ""},
				Target:  Image{Width: 400, Height: 600, Quality: 65, Format: bimg.JPEG, Hash: ""},
				Filters: map[string]string{"crop": "limit"},
				Hasher:  fakeHasher,
			},
			"with multiple filter auto jpeg",
		},
		{
			"w_400,c_limit,h_600,f_png/" + testURL,
			Job{
				Source:  Image{URL: testURL, Hash: ""},
				Target:  Image{Width: 400, Height: 600, Format: bimg.PNG, Hash: ""},
				Filters: map[string]string{"crop": "limit"},
				Hasher:  fakeHasher,
			},
			"with multiple filter png",
		},
		{
			"w_400,c_limit,h_600,f_gif/" + testURL,
			Job{
				Source:  Image{URL: testURL, Hash: ""},
				Target:  Image{Width: 400, Height: 600, Format: bimg.GIF, Hash: ""},
				Filters: map[string]string{"crop": "limit"},
				Hasher:  fakeHasher,
			},
			"with multiple filter gif",
		},
		{
			"w_400,c_limit,h_600,f_jpeg/" + testURL,
			Job{
				Source:  Image{URL: testURL, Hash: ""},
				Target:  Image{Width: 400, Height: 600, Format: bimg.JPEG, Hash: ""},
				Filters: map[string]string{"crop": "limit"},
				Hasher:  fakeHasher,
			},
			"with multiple filter jpeg",
		},
		{
			"w_400,c_limit,h_600,f_jpeg/file.jpg",
			Job{
				Source:  Image{URL: "file.jpg", Hash: ""},
				Target:  Image{Width: 400, Height: 600, Format: bimg.JPEG, Hash: ""},
				Filters: map[string]string{"crop": "limit"},
				Hasher:  fakeHasher,
			},
			"plain uploaded file with multiple filter jpeg",
		},
		{
			"file.jpg",
			Job{
				Source:  Image{URL: "file.jpg", Hash: ""},
				Target:  Image{Format: bimg.JPEG, Hash: ""},
				Filters: map[string]string{"crop": "scale"},
				Hasher:  fakeHasher,
			},
			"plain uploaded file without filters",
		},
		{
			"w_400,c_limit,h_600,f_jpeg/ffolder/file.jpg",
			Job{
				Source:  Image{URL: "file.jpg", Hash: ""},
				Target:  Image{Width: 400, Height: 600, Format: bimg.JPEG, Hash: ""},
				Filters: map[string]string{"crop": "limit"},
				Hasher:  fakeHasher,
			},
			"plain uploaded foldered file without filters",
		},
		{
			"folder/file.jpg",
			Job{
				Source:  Image{URL: "file.jpg", Hash: ""},
				Target:  Image{Format: bimg.JPEG, Hash: ""},
				Filters: map[string]string{"crop": "scale"},
				Hasher:  fakeHasher,
			},
			"plain uploaded foldered file with filters",
		},
	}
	for _, test := range cases {
		job := NewJob()
		job.Hasher = fakeHasher
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
		fmt.Errorf("format \"fake\" not allowed"),
		"Bad image format",
	},
	{
		"w_400,c_limit,h_pp/" + testURL,
		fmt.Errorf("targetHeight is not integer: strconv.Atoi: parsing \"pp\": invalid syntax"),
		"TargetHeight is not integer",
	},
	{
		"w_OOO,c_limit,h_500/" + testURL,
		fmt.Errorf("targetWidth is not integer: strconv.Atoi: parsing \"OOO\": invalid syntax"),
		"TargetWidth is not integer",
	},
	{
		"w_100,c_fake,h_500/" + testURL,
		errors.New("crop \"fake\" not allowed"),
		"Crop is not allowed",
	},
	{
		"w_100,c_limit,h_500,q_fake/" + testURL,
		fmt.Errorf("quality is not integer: strconv.Atoi: parsing \"fake\": invalid syntax"),
		"Quality is not an integer",
	},
}

func TestParseFail(t *testing.T) {
	for _, test := range parserErrorCases {
		img := NewJob()
		err := img.Parse(test.url)
		assert.Equal(t, test.err, err, test.description)
	}
}

var cropCases = []struct {
	crop           string
	sourceWidth    int
	sourceHeight   int
	targetWidth    int
	targetHeight   int
	expectedWidth  int
	expectedHeight int
	message        string
}{
	{"scale", 1000, 500, 100, 50, 100, 50, "scale hor"},
	{"scale", 500, 1000, 50, 100, 50, 100, "scale vert"},
	{"fit", 1000, 500, 100, 50, 100, 50, "fit hor-hor"},
	{"fit", 500, 1000, 100, 50, 100, 200, "fit ver-ver"},
	{"fit", 1000, 500, 50, 100, 200, 100, "fit hor-ver"},
	{"fit", 500, 1000, 50, 100, 50, 100, "fit ver-hor"},
	{"fit", 50, 100, 500, 1000, 500, 1000, "fit bigger"},
	{"fit", 2000, 1000, 0, 2000, 4000, 2000, "fit bigger without w"},
	{"limit", 50, 100, 500, 1000, 50, 100, "limit bigger"},
	{"limit", 1000, 500, 100, 50, 100, 50, "limit hor-hor"},
	{"limit", 500, 1000, 100, 50, 100, 200, "limit ver-ver"},
	{"limit", 1000, 500, 50, 100, 200, 100, "limit hor-ver"},
	{"limit", 500, 1000, 50, 100, 50, 100, "limit ver-hor"},
}

func TestCrop(t *testing.T) {
	for _, test := range cropCases {
		job := NewJob()
		job.Source.Width = test.sourceWidth
		job.Source.Height = test.sourceHeight
		job.Source.AspectRatio = float32(test.sourceWidth) / float32(test.sourceHeight)
		job.Target.Width = test.targetWidth
		job.Target.Height = test.targetHeight
		job.Filters["crop"] = test.crop

		err := job.Crop()
		assert.Nil(t, err)
		assert.Equal(t, test.expectedHeight, job.Target.Height, test.message)
		assert.Equal(t, test.expectedWidth, job.Target.Width, test.message)
	}
}
