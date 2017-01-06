package imagejob

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const goodUrl = "/v0.1/fetch/w_500,c_limit/http://upload.wikimedia.org/wikipedia/commons/0/0c/Scarlett_Johansson_Césars_2014.jpg"
const badFilter = "/v0.1/fetch/w_pp,c_limit/http://upload.wikimedia.org/wikipedia/commons/0/0c/Scarlett_Johansson_Césars_2014.jpg"
const badURL = "/v0.1/fetch/w_500,c_limit/http://fakedomain.com/wikiped.jpg"

var fetchCases = []struct {
	url     string
	method  string
	status  int
	message string
}{
	{
		"/v0.1/fetch/w_500,c_limit/http://upload.wikimedia.org/wikipedia/commons/0/0c/Scarlett_Johansson_Césars_2014.jpg",
		"GET",
		200,
		"Regular use",
	},
	{
		"/v0.1/fetch/w_500,c_limit/http://upload.wikimedia.org/wikipedia/commons/0/0c/Scarlett_Johansson_Césars_2014.jpg",
		"POST",
		405,
		"Bad Method POST",
	},
	{
		"/v0.1/fetch/w_500,c_limit/http://upload.wikimedia.org/wikipedia/commons/0/0c/Scarlett_Johansson_Césars_2014.jpg",
		"PUT",
		405,
		"Bad Method PUT",
	},
	{
		"/v0.1/fetch/w_pp,c_limit/http://upload.wikimedia.org/wikipedia/commons/0/0c/Scarlett_Johansson_Césars_2014.jpg",
		"GET",
		400,
		"Wrong filter",
	},
	{
		"/v0.1/fetch/w_500,c_limit/http://upload.com/wikihansson/Césars_2014.jpg",
		"GET",
		500,
		"Non existent URI",
	},
}

func TestFetch(t *testing.T) {
	for _, test := range fetchCases {
		req, _ := http.NewRequest(test.method, test.url, nil)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(Fetch)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, test.status, rr.Code)
	}
}
