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

func Fetch(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	img := ImageJob{}

	if err := img.New(ps.ByName("info")[1:]); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if err := img.Download(); err != nil {
		http.Error(w, "Cannot download image", http.StatusInternalServerError)
		return
	}

	if err := img.Process(w); err != nil {
		http.Error(w, "Cannot process Image", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/"+img.TargetFormat)
}
