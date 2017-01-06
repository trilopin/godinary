package imagejob

import (
	"net/http"
	"os"
	"strconv"
	"strings"
)

// globalSemaphore controls concurrent http client requests
var globalSemaphore = make(chan struct{}, func() int {
	maxRequests, err := strconv.Atoi(os.Getenv("GODINARY_MAX_REQUEST"))
	if maxRequests == 0 || err != nil {
		return 20
	}
	return maxRequests
}())

// Fetch takes url + params in url to download iamge from url and apply filters
func Fetch(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	urlInfo := strings.Replace(r.URL.Path, "/v0.1/fetch/", "", 1)
	job := NewImageJob()
	if err := job.Parse(urlInfo); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	globalSemaphore <- struct{}{}
	body, err := job.Source.Download()
	<-globalSemaphore

	if err != nil {
		http.Error(w, "Cannot download image", http.StatusInternalServerError)
		return
	}

	job.Source.Decode(body)

	w.Header().Set("Content-Type", "image/"+job.Target.Format)
	if err := job.Process(w); err != nil {
		http.Error(w, "Cannot process Image", http.StatusInternalServerError)
		return
	}

}
