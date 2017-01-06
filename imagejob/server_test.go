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

func TestFetch(t *testing.T) {
	req, _ := http.NewRequest("GET", goodUrl, nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Fetch)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, 200, rr.Code)
}

func TestBadMethod(t *testing.T) {
	req, _ := http.NewRequest("POST", goodUrl, nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Fetch)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, 405, rr.Code)
}

func TestBadFilter(t *testing.T) {
	req, _ := http.NewRequest("GET", badFilter, nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Fetch)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, 400, rr.Code)
}

func TestBadURL(t *testing.T) {
	req, _ := http.NewRequest("GET", badURL, nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Fetch)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, 500, rr.Code)
}
