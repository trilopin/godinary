package image

import (
	"errors"
	"fmt"
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
	expected    Job
	description string
}{
	{
		testURL,
		Job{
			Source: Image{
				URL:  testURL,
				Hash: "9c2eb35928a2ee6ab8221c393fd306348b1235e282ecb32f0e41ca1bba6e90a9",
			},
			Target: Image{
				Format: bimg.JPEG,
				Hash:   "7568485b3a4ec0d74e9f0220d7c83fa4a6e4f1e4399f6903eda9483fdc8cdbc0",
			},
			Filters: map[string]string{"crop": "scale"},
		},
		"without filters",
	},
	{
		"w_400/" + testURL,
		Job{
			Source: Image{
				URL:  testURL,
				Hash: "9c2eb35928a2ee6ab8221c393fd306348b1235e282ecb32f0e41ca1bba6e90a9",
			},
			Target: Image{
				Width:  400,
				Format: bimg.JPEG,
				Hash:   "7499a3a704204fcc4811d8af480ba2f27916f37132c295687d0e3b9ae2ca992f",
			},
			Filters: map[string]string{"crop": "scale"},
		},
		"with one filter",
	},
	{
		"w_400,c_limit,h_600,f_jpg/" + testURL,
		Job{
			Source: Image{
				URL:  testURL,
				Hash: "9c2eb35928a2ee6ab8221c393fd306348b1235e282ecb32f0e41ca1bba6e90a9",
			},
			Target: Image{
				Width:  400,
				Height: 600,
				Format: bimg.JPEG,
				Hash:   "00f5cefda635ef7e5a657e8363ea3fd15ef90028b9ea2f8372fa828631d506ca",
			},
			Filters: map[string]string{"crop": "limit"},
		},
		"with multiple filter jpg",
	},
	{
		"w_400,c_limit,h_600,f_webp,q_50/" + testURL,
		Job{
			Source: Image{
				URL:  testURL,
				Hash: "9c2eb35928a2ee6ab8221c393fd306348b1235e282ecb32f0e41ca1bba6e90a9",
			},
			Target: Image{
				Width:   400,
				Height:  600,
				Quality: 50,
				Format:  bimg.WEBP,
				Hash:    "ac799d0068df5144351ebed0f03fc846c6263c44cc543d31b8e6b33cc8f10103",
			},
			Filters: map[string]string{"crop": "limit"},
		},
		"with multiple filter webp",
	},
	{
		"w_400,c_limit,h_600,f_auto,q_65/" + testURL,
		Job{
			Source: Image{
				URL:  testURL,
				Hash: "9c2eb35928a2ee6ab8221c393fd306348b1235e282ecb32f0e41ca1bba6e90a9",
			},
			Target: Image{
				Width:   400,
				Height:  600,
				Quality: 65,
				Format:  bimg.JPEG,
				Hash:    "511d78de767e504b5b8c805011ed6e4d9c26ce76f9f36d50ff33e29a96d0c3a4",
			},
			Filters: map[string]string{"crop": "limit"},
		},
		"with multiple filter auto jpeg",
	},
	{
		"w_400,c_limit,h_600,f_png/" + testURL,
		Job{
			Source: Image{
				URL:  testURL,
				Hash: "9c2eb35928a2ee6ab8221c393fd306348b1235e282ecb32f0e41ca1bba6e90a9",
			},
			Target: Image{
				Width:  400,
				Height: 600,
				Format: bimg.PNG,
				Hash:   "506c667899dc573c359a15d6063bd96f67f006b86f83b207294c51d8fc45affc",
			},
			Filters: map[string]string{"crop": "limit"},
		},
		"with multiple filter png",
	},
	{
		"w_400,c_limit,h_600,f_gif/" + testURL,
		Job{
			Source: Image{
				URL:  testURL,
				Hash: "9c2eb35928a2ee6ab8221c393fd306348b1235e282ecb32f0e41ca1bba6e90a9",
			},
			Target: Image{
				Width:  400,
				Height: 600,
				Format: bimg.GIF,
				Hash:   "9ed8868e8d61f26e717129d0bac06787041e6048d1dcf5a906c23ac0efa3d4af",
			},
			Filters: map[string]string{"crop": "limit"},
		},
		"with multiple filter gif",
	},
	{
		"w_400,c_limit,h_600,f_jpeg/" + testURL,
		Job{
			Source: Image{
				URL:  testURL,
				Hash: "9c2eb35928a2ee6ab8221c393fd306348b1235e282ecb32f0e41ca1bba6e90a9",
			},
			Target: Image{
				Width:  400,
				Height: 600,
				Format: bimg.JPEG,
				Hash:   "962aa5597fbb18caabfacfee9a1b612ad478341d4e8a573682788b773d874a37",
			},
			Filters: map[string]string{"crop": "limit"},
		},
		"with multiple filter jpeg",
	},
	{
		"w_400,c_limit,h_600,f_jpeg/file.jpg",
		Job{
			Source: Image{
				URL:  "file.jpg",
				Hash: "91a721b7244245b40e368346edc97c311439a0260e4d51f60023eb1bc86d7238",
			},
			Target: Image{
				Width:  400,
				Height: 600,
				Format: bimg.JPEG,
				Hash:   "4d083201587ff3b5b4d268253c5070b650d64051b646a29bc352ea57878eb41b",
			},
			Filters: map[string]string{"crop": "limit"},
		},
		"plain uploaded file with multiple filter jpeg",
	},
	{
		"file.jpg",
		Job{
			Source: Image{
				URL:  "file.jpg",
				Hash: "91a721b7244245b40e368346edc97c311439a0260e4d51f60023eb1bc86d7238",
			},
			Target: Image{
				Format: bimg.JPEG,
				Hash:   "79fba77d9b83b2a93091d8a50c51e98337ff9f500ada410aef0ed39ceb3b6aad",
			},
			Filters: map[string]string{"crop": "limit"},
		},
		"plain uploaded file without filters",
	},
}

func TestParse(t *testing.T) {
	for _, test := range parserCases {
		job := NewJob()
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
