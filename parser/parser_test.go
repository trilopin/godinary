package parser

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testURL = "http://upload.wikimedia.org/wikipedia/commons/0/0c/Scarlett_Johansson_Césars_2014.jpg"

var parserCases = []struct {
	url         string
	expected    Job
	description string
}{
	{
		testURL,
		Job{
			SourceURL:  testURL,
			Filters:    map[string]string{},
			FixedRatio: false,
		},
		"without filters",
	},
	{
		"w_400/" + testURL,
		Job{
			SourceURL:    testURL,
			Filters:      map[string]string{},
			DesiredWidth: 400,
			FixedRatio:   false,
		},
		"with one filter",
	},
	{
		"w_400,c_limit,h_600,f_jpg/" + testURL,
		Job{
			SourceURL:     testURL,
			Filters:       map[string]string{},
			DesiredWidth:  400,
			DesiredHeight: 600,
			DesiredFormat: "jpg",
			FixedRatio:    true,
		},
		"with multiple filter jpg",
	},
	{
		"w_400,c_limit,h_600,f_png/" + testURL,
		Job{
			SourceURL:     testURL,
			Filters:       map[string]string{},
			DesiredWidth:  400,
			DesiredHeight: 600,
			DesiredFormat: "png",
			FixedRatio:    true,
		},
		"with multiple filter png",
	},
	{
		"w_400,c_limit,h_600,f_gif/" + testURL,
		Job{
			SourceURL:     testURL,
			Filters:       map[string]string{},
			DesiredWidth:  400,
			DesiredHeight: 600,
			DesiredFormat: "gif",
			FixedRatio:    true,
		},
		"with multiple filter gif",
	},
	{
		"w_400,c_limit,h_600,f_jpeg/" + testURL,
		Job{
			SourceURL:     testURL,
			Filters:       map[string]string{},
			DesiredWidth:  400,
			DesiredHeight: 600,
			DesiredFormat: "jpeg",
			FixedRatio:    true,
		},
		"with multiple filter jpeg",
	},
}

func TestNew(t *testing.T) {
	for _, test := range parserCases {
		job := Job{}
		err := job.New(test.url)
		assert.Nil(t, err)
		assert.Equal(t, test.expected, job, test.description)
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
		errors.New("DesiredHeight is not integer"),
		"DesiredHeight is not integer",
	},
	{
		"w_OOO,c_limit,h_500/" + testURL,
		errors.New("DesiredWidth is not integer"),
		"DesiredWidth is not integer",
	},
}

func TestParseWithError(t *testing.T) {
	for _, test := range parserErrorCases {
		job := Job{}
		err := job.New(test.url)
		assert.Equal(t, test.err, err, test.description)
	}
}
