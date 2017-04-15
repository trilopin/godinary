package imagejob

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	bimg "gopkg.in/h2non/bimg.v1"

	"github.com/trilopin/godinary/storage"
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
	var body io.Reader

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

	// derived image is already cached
	if body, err = storage.StorageDriver.Read(job.Target.Hash); err == nil {
		if cached, err2 := ioutil.ReadAll(body); err2 != nil {
			writeImage(w, cached, job.Target.Format)
			return
		}
	}

	// Download if does not exists at storage, load otherwise
	body, err = storage.StorageDriver.Read(job.Source.Hash)
	if err == nil {
		job.Source.Load(body)
	} else {
		globalThrotling <- struct{}{}
		domainThrotle <- struct{}{}
		err = job.Source.Download(storage.StorageDriver)
		<-domainThrotle
		<-globalThrotling

		if err != nil {
			http.Error(w, "Cannot download image", http.StatusInternalServerError)
			return
		}
	}

	job.Source.ExtractInfo()
	job.crop()

	// do the process thing
	if err := job.Target.Process(job.Source, storage.StorageDriver); err != nil {
		log.Println(err)
		http.Error(w, "Cannot process Image", http.StatusInternalServerError)
	}

	writeImage(w, job.Target.RawContent, job.Target.Format)
}

func writeImage(w http.ResponseWriter, buffer []byte, format bimg.ImageType) {
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer)))
	w.Header().Set("Content-Type", fmt.Sprintf("image/%s", bimg.ImageTypes[format]))
	w.Write(buffer)
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
