package imagejob

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

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
		"/image/fetch/w_500,c_limit/http://upload.wikimedia.org/wikipedia/commons/0/0c/Scarlett_Johansson_Césars_2014.jpg",
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
		"/image/fetch/w_500,c_limit/",
		"GET",
		500,
		"Non existent URI",
	},
}

func setupModule() {

	os.Setenv("GODINARY_FS_BASE", "/tmp/.godinary/")
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
			assert.NotEqual(t, "", rr.Body.String())
		}
	}
}
