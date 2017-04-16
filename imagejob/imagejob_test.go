package imagejob

import (
	"errors"
	"testing"

	bimg "gopkg.in/h2non/bimg.v1"

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
				URL:  testURL,
				Hash: "9c2eb35928a2ee6ab8221c393fd306348b1235e282ecb32f0e41ca1bba6e90a9",
			},
			Target: Image{
				Format: bimg.JPEG,
				Hash:   "9c2eb35928a2ee6ab8221c393fd306348b1235e282ecb32f0e41ca1bba6e90a9",
			},
			Filters: map[string]string{"crop": "scale"},
		},
		"without filters",
	},
	{
		"w_400/" + testURL,
		ImageJob{
			Source: Image{
				URL:  testURL,
				Hash: "9c2eb35928a2ee6ab8221c393fd306348b1235e282ecb32f0e41ca1bba6e90a9",
			},
			Target: Image{
				Width:  400,
				Format: bimg.JPEG,
				Hash:   "30e1d866a540d350b083cd5989c65ae0ce7b74b6cd4a9a3bf578868160c19ab5",
			},
			Filters: map[string]string{"crop": "scale"},
		},
		"with one filter",
	},
	{
		"w_400,c_limit,h_600,f_jpg/" + testURL,
		ImageJob{
			Source: Image{
				URL:  testURL,
				Hash: "9c2eb35928a2ee6ab8221c393fd306348b1235e282ecb32f0e41ca1bba6e90a9",
			},
			Target: Image{
				Width:  400,
				Height: 600,
				Format: bimg.JPEG,
				Hash:   "3d5bd40f726f97e8ef21b079953d437204a10f1e42d4066025f87c6c1914f195",
			},
			Filters: map[string]string{"crop": "limit"},
		},
		"with multiple filter jpg",
	},
	{
		"w_400,c_limit,h_600,f_webp,q_50/" + testURL,
		ImageJob{
			Source: Image{
				URL:  testURL,
				Hash: "9c2eb35928a2ee6ab8221c393fd306348b1235e282ecb32f0e41ca1bba6e90a9",
			},
			Target: Image{
				Width:   400,
				Height:  600,
				Quality: 50,
				Format:  bimg.WEBP,
				Hash:    "e03aa68f05f0aef45860d96878e64facbbc6b48f2f04c22ce44fff3021daa5bb",
			},
			Filters: map[string]string{"crop": "limit"},
		},
		"with multiple filter webp",
	},
	{
		"w_400,c_limit,h_600,f_auto,q_65/" + testURL,
		ImageJob{
			Source: Image{
				URL:  testURL,
				Hash: "9c2eb35928a2ee6ab8221c393fd306348b1235e282ecb32f0e41ca1bba6e90a9",
			},
			Target: Image{
				Width:   400,
				Height:  600,
				Quality: 65,
				Format:  bimg.JPEG,
				Hash:    "ccc39c9c087b451f95a41c12f2f67ace0c4a337d74dfcc0d28faed8336ff4df8",
			},
			Filters: map[string]string{"crop": "limit"},
		},
		"with multiple filter auto jpeg",
	},
	{
		"w_400,c_limit,h_600,f_png/" + testURL,
		ImageJob{
			Source: Image{
				URL:  testURL,
				Hash: "9c2eb35928a2ee6ab8221c393fd306348b1235e282ecb32f0e41ca1bba6e90a9",
			},
			Target: Image{
				Width:  400,
				Height: 600,
				Format: bimg.PNG,
				Hash:   "e976de427ede66ef4c1af9216b59c76b2c030fec0537882dfe507c99d5c542fb",
			},
			Filters: map[string]string{"crop": "limit"},
		},
		"with multiple filter png",
	},
	{
		"w_400,c_limit,h_600,f_gif/" + testURL,
		ImageJob{
			Source: Image{
				URL:  testURL,
				Hash: "9c2eb35928a2ee6ab8221c393fd306348b1235e282ecb32f0e41ca1bba6e90a9",
			},
			Target: Image{
				Width:  400,
				Height: 600,
				Format: bimg.GIF,
				Hash:   "25be0d36afcd9cf64e785b5cf52f13c332d1b9c6d544c69a02c6a51cc1c40743",
			},
			Filters: map[string]string{"crop": "limit"},
		},
		"with multiple filter gif",
	},
	{
		"w_400,c_limit,h_600,f_jpeg/" + testURL,
		ImageJob{
			Source: Image{
				URL:  testURL,
				Hash: "9c2eb35928a2ee6ab8221c393fd306348b1235e282ecb32f0e41ca1bba6e90a9",
			},
			Target: Image{
				Width:  400,
				Height: 600,
				Format: bimg.JPEG,
				Hash:   "7bccb1bc86df66b2e348de55c100634cdf407a815b628b8dab0d8a8ac9519b7f",
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
	{
		"w_100,c_limit,h_500,q_fake/" + testURL,
		errors.New("Quality is not integer"),
		"Quality is not an integer",
	},
}

func TestParseFail(t *testing.T) {
	for _, test := range parserErrorCases {
		img := NewImageJob()
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
		job := NewImageJob()
		job.Source.Width = test.sourceWidth
		job.Source.Height = test.sourceHeight
		job.Source.AspectRatio = float32(test.sourceWidth) / float32(test.sourceHeight)
		job.Target.Width = test.targetWidth
		job.Target.Height = test.targetHeight
		job.Filters["crop"] = test.crop

		err := job.crop()
		assert.Nil(t, err)
		assert.Equal(t, test.expectedHeight, job.Target.Height, test.message)
		assert.Equal(t, test.expectedWidth, job.Target.Width, test.message)
	}
}
