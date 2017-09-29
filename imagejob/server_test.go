package imagejob

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	bimg "gopkg.in/h2non/bimg.v1"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/trilopin/godinary/storage"
)

var fetchCases = []struct {
	url     string
	method  string
	status  int
	message string
}{
	{
		"/image/fetch/w_100,h_100,c_limit/http://upload.wikimedia.org/wikipedia/commons/0/0c/Scarlett_Johansson_Césars_2014.jpg",
		"GET",
		200,
		"Regular use",
	},
	{
		"/image/fetch/w_500,c_limit/http://upload.wikimedia.org/wikipedia/commons/0/0c/Scarlett_Johansson_Césars_2014.jpg",
		"POST",
		405,
		"Bad Method POST",
	},
	{
		"/image/fetch/w_500,c_limit/http://upload.wikimedia.org/wikipedia/commons/0/0c/Scarlett_Johansson_Césars_2014.jpg",
		"PUT",
		405,
		"Bad Method PUT",
	},
	{
		"/image/fetch/w_pp,c_limit/http://upload.wikimedia.org/wikipedia/commons/0/0c/Scarlett_Johansson_Césars_2014.jpg",
		"GET",
		400,
		"Wrong filter",
	},
	{
		"/image/fetch/w_500,c_limit/http://fake.dot.org/wiksdafadsfasdfadsfipedi",
		"GET",
		404,
		"Non existent URI",
	},
}

func setupModule() {

	viper.Set("fs_base", "/tmp/.godinary/")
	storage.StorageDriver = storage.NewFileDriver()
	MaxRequest = 2
	MaxRequestPerDomain = 1
	SpecificThrotling = make(map[string]chan struct{}, MaxRequestPerDomain)
	GlobalThrotling = make(chan struct{}, MaxRequest)
}

func TestFetch(t *testing.T) {
	setupModule()
	defer os.RemoveAll("/tmp/.godinary")

	for _, test := range fetchCases {
		req, _ := http.NewRequest(test.method, test.url, nil)
		req.Header.Set("Accept", "image/webp")
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(Fetch)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, test.status, rr.Code, test.message)

		if test.status == 200 {
			buffer, err := ioutil.ReadAll(rr.Body)
			assert.Nil(t, err)
			img := bimg.NewImage(buffer)
			size, _ := img.Size()
			assert.Equal(t, 141, size.Height, "height")
			assert.Equal(t, 100, size.Width, "width")
		}
	}
}
