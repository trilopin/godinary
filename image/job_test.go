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
		input       string
		isFetch     bool
		job         Job
		description string
	}{
		{
			"w_400,c_limit,h_600,f_auto,q_65/" + testURL,
			true,
			Job{
				Source: Image{
					URL:  testURL,
					Hash: "",
				},
				Target: Image{
					Width:   400,
					Height:  600,
					Quality: 65,
					Format:  bimg.JPEG,
					Hash:    "",
				},
				Filters: map[string]string{"crop": "limit"},
				Hasher:  fakeHasher,
			},
			"multiple filters",
		},
	}
	for _, test := range cases {
		job := NewJob()
		job.Hasher = fakeHasher
		err := job.Parse(test.input, test.isFetch)
		assert.Nil(t, err)
		assert.Equal(t, test.job, *job, test.description)
	}
}

func TestParseURL(t *testing.T) {
	cases := []struct {
		input       string
		isFetch     bool
		image       string
		filters     string
		description string
	}{
		// fetch cases
		{testURL, true, testURL, "", "fetch without filters"},
		{testSecureURL, true, testSecureURL, "", "fetch secure without filters"},
		{"w_400/" + testURL, true, testURL, "w_400", "fetch with filter"},
		{"w_400/" + testURL, true, testURL, "w_400", "fetch with multiple filter"},
		{"w_400,c_limit,h_600,f_jpg/" + testSecureURL, true, testSecureURL, "w_400,c_limit,h_600,f_jpg", "fetch secure with filter"},
		{"w_400,c_limit,h_600,f_jpg/" + testSecureURL, true, testSecureURL, "w_400,c_limit,h_600,f_jpg", "fetch secure with multiple filter"},
		// upload cases
		{"file.jpg", false, "file.jpg", "", "upload without filters"},
		{"folder/file.jpg", false, "file.jpg", "", "upload without filters and folder"},
		{"w_400/file.jpg", false, "file.jpg", "w_400", "upload with one filters"},
		{"w_400,c_limit,h_600,f_jpeg/file.jpg", false, "file.jpg", "w_400,c_limit,h_600,f_jpeg", "upload with multiple filters"},
		{"w_400,c_limit,h_600,f_jpeg/folder/file.jpg", false, "file.jpg", "w_400,c_limit,h_600,f_jpeg", "upload with multiple filters and folder"},
	}
	for _, test := range cases {
		filters, image, err := parseURL(test.input, test.isFetch)
		assert.Nil(t, err)
		assert.Equal(t, test.image, image, test.description)
		assert.Equal(t, test.filters, filters, test.description)
	}
}

func TestParseFilters(t *testing.T) {
	cases := []struct {
		input       string
		expected    *Job
		err         error
		description string
	}{
		{
			"w_400",
			&Job{Target: Image{Width: 400, Format: bimg.JPEG}, Filters: map[string]string{"crop": "scale"}},
			nil, "only width",
		},
		{
			"h_400,w_400",
			&Job{Target: Image{Height: 400, Width: 400, Format: bimg.JPEG}, Filters: map[string]string{"crop": "scale"}},
			nil, "width - height",
		},
		{
			"f_png",
			&Job{Target: Image{Format: bimg.PNG}, Filters: map[string]string{"crop": "scale"}},
			nil, "png format",
		},
		{
			"f_gif",
			&Job{Target: Image{Format: bimg.GIF}, Filters: map[string]string{"crop": "scale"}},
			nil, "gif format",
		},
		{
			"f_webp",
			&Job{Target: Image{Format: bimg.WEBP}, Filters: map[string]string{"crop": "scale"}},
			nil, "gif format",
		},
		{
			"h_400",
			&Job{Target: Image{Height: 400, Format: bimg.JPEG}, Filters: map[string]string{"crop": "scale"}},
			nil, "only height",
		},
		{
			"c_limit",
			&Job{Target: Image{Format: bimg.JPEG}, Filters: map[string]string{"crop": "limit"}},
			nil, "only crop",
		},
		{
			"h_400,w_400,f_png,q_55,c_limit",
			&Job{Target: Image{Height: 400, Width: 400, Format: bimg.PNG, Quality: 55}, Filters: map[string]string{"crop": "limit"}},
			nil, "all filters",
		},
		{"c_fake", nil, errors.New("crop \"fake\" not allowed"), "Crop  not accepted"},
	}
	for _, test := range cases {
		job := NewJob()
		err := job.parseFilters(test.input)
		assert.Equal(t, test.err, err, test.description)
		if err == nil {
			assert.Equal(t, test.expected.Target.Height, job.Target.Height, test.description)
			assert.Equal(t, test.expected.Target.Width, job.Target.Width, test.description)
			assert.Equal(t, test.expected.Target.Format, job.Target.Format, test.description)
			assert.Equal(t, test.expected.Target.Quality, job.Target.Quality, test.description)
			assert.Equal(t, test.expected.Filters, job.Filters, test.description)
		} else {
			t.Log(err)
		}
	}
}

func TestParseCrop(t *testing.T) {
	cases := []struct {
		crop        string
		expected    string
		err         error
		description string
	}{
		{"scale", "scale", nil, "Scale"},
		{"limit", "limit", nil, "Limit"},
		{"fit", "fit", nil, "Fit"},
		{"fake", "scale", errors.New("crop \"fake\" not allowed"), "Error case not accepted"},
	}
	for _, test := range cases {
		crop, err := parseCrop(test.crop)
		assert.Equal(t, test.err, err, test.description)
		assert.Equal(t, test.expected, crop, test.description)
	}
}

func TestParseFormat(t *testing.T) {
	cases := []struct {
		format      string
		webp        bool
		expected    bimg.ImageType
		err         error
		description string
	}{
		{"jpg", false, bimg.JPEG, nil, "jpg format"},
		{"jpeg", false, bimg.JPEG, nil, "jpeg format"},
		{"png", false, bimg.PNG, nil, "png format"},
		{"gif", false, bimg.GIF, nil, "gif format"},
		{"auto", false, bimg.JPEG, nil, "auto without webp"},
		// {"auto", true, bimg.WEBP, nil, "auto with webp"},
		{"fake", false, bimg.JPEG, errors.New("format \"fake\" not allowed"), "Error case not accepted"},
	}
	for _, test := range cases {
		format, err := parseFormat(test.format, test.webp)
		assert.Equal(t, test.err, err, test.description)
		assert.Equal(t, test.expected, format, test.description)
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
		err := img.Parse(test.url, true)
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
