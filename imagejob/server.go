package imagejob

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// globalSemaphore controls concurrent http client requests
var specificThrotling = make(map[string]chan struct{}, 20)
var globalThrotling = make(chan struct{}, func() int {
	maxRequests, err := strconv.Atoi(os.Getenv("GODINARY_MAX_REQUEST"))
	if maxRequests == 0 || err != nil {
		maxRequests = 20
	}
	return maxRequests
}())

// Concurrency is a handler for testing concurrency levels
func Concurrency(w http.ResponseWriter, r *http.Request) {
	domainThrotle, ok := specificThrotling["fake"]
	if !ok {
		domainThrotle = make(chan struct{}, 10)
		specificThrotling["fake"] = domainThrotle
	}

	globalThrotling <- struct{}{}
	fmt.Println("global acquired")
	domainThrotle <- struct{}{}
	fmt.Println("domain acquired")
	time.Sleep(time.Millisecond * 100)
	<-domainThrotle
	<-globalThrotling
	fmt.Println("finished")
}

// Fetch takes url + params in url to download image from url and apply filters
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

	domain, err := topDomain(job.Source.URL)
	if err != nil {
		http.Error(w, "Cannot parse hostname", http.StatusInternalServerError)
		return
	}
	domainThrotle, ok := specificThrotling[domain]
	if !ok {
		domainThrotle = make(chan struct{}, 1)
		specificThrotling[domain] = domainThrotle
	}

	globalThrotling <- struct{}{}
	domainThrotle <- struct{}{}
	body, err := job.Source.Download()
	<-domainThrotle
	<-globalThrotling

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

func topDomain(URL string) (string, error) {
	info, err := url.Parse(URL)
	if err != nil {
		return "", errors.New("Cannot parse hostname")
	}
	parts := strings.Split(info.Host, ".")
	if len(parts) <= 1 {
		return "", errors.New("Cannot parse hostname")
	}
	return strings.Join(parts[len(parts)-2:], "."), nil
}
