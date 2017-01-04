package imagejob

import (
	"net/http"
	"os"
	"strconv"

	"github.com/julienschmidt/httprouter"
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
func Fetch(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	job := NewImageJob()

	if err := job.Parse(ps.ByName("info")[1:]); err != nil {
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

	if err := job.Process(w); err != nil {
		http.Error(w, "Cannot process Image", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/"+job.Target.Format)
}
