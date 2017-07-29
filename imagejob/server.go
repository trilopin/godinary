package imagejob

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/trilopin/godinary/storage"

	bimg "gopkg.in/h2non/bimg.v1"
)

var (
	// MaxRequest is global max concurrency for download
	MaxRequest int
	// MaxRequestPerDomain is max concurrency per out domain
	MaxRequestPerDomain int
	// SpecificThrotling is semaphore per domain
	SpecificThrotling map[string]chan struct{}
	// GlobalThrotling is global semaphore
	GlobalThrotling chan struct{}
)

// Fetch takes url + params in url to download image from url and apply filters
func Fetch(w http.ResponseWriter, r *http.Request) {
	var body io.Reader
	var dSem float64

	t1 := time.Now()
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	urlInfo := strings.Replace(r.URL.Path, "/image/fetch/", "", 1)
	job := NewImageJob()
	job.AcceptWebp = strings.Contains(r.Header["Accept"][0], "image/webp")

	if err := job.Parse(urlInfo); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	domain, err := topDomain(job.Source.URL)
	if err != nil || domain == "" {
		http.Error(w, "Cannot parse domain", http.StatusInternalServerError)
		return
	}
	domainThrotle := make(chan struct{}, MaxRequestPerDomain)
	SpecificThrotling[domain] = domainThrotle

	// derived image is already cached
	if body, err = storage.StorageDriver.Read(job.Target.Hash); err == nil {
		if cached, err2 := ioutil.ReadAll(body); err2 == nil {
			writeImage(w, cached, job.Target.Format)
			log.Printf("CACHED - TOTAL %0.5f", time.Since(t1).Seconds())
			return
		}
	}

	// Download if original image does not exists at storage, load otherwise
	body, err = storage.StorageDriver.Read(job.Source.Hash)
	if err == nil {
		job.Source.Load(body)
	} else {
		tSem := time.Now()
		log.Printf("SEM %s %d/%d %d/%d", domain, len(GlobalThrotling), cap(GlobalThrotling), len(domainThrotle), cap(domainThrotle))
		GlobalThrotling <- struct{}{}
		domainThrotle <- struct{}{}
		dSem = time.Since(tSem).Seconds()
		err = job.Source.Download(storage.StorageDriver)
		<-domainThrotle
		<-GlobalThrotling

		if err != nil {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
	}
	t2 := time.Now()

	job.Source.ExtractInfo()
	job.crop()

	// do the process thing
	if err := job.Target.Process(job.Source, storage.StorageDriver); err != nil {
		log.Println("Error processing image ", job.Source.URL, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	t3 := time.Now()

	writeImage(w, job.Target.RawContent, job.Target.Format)
	log.Printf(
		"NEW - TOTAL %0.5f => SEM %0.5f, DOWN %0.5f, PROC %0.5f",
		time.Since(t1).Seconds(), dSem,
		t2.Sub(t1).Seconds()-dSem, t3.Sub(t2).Seconds())

}

func topDomain(URL string) (string, error) {
	info, err := url.Parse(URL)
	if err != nil {
		return "", errors.New("Cannot parse hostname")
	}
	return info.Host, nil
}

func writeImage(w http.ResponseWriter, buffer []byte, format bimg.ImageType) {
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer)))
	w.Header().Set("Content-Type", fmt.Sprintf("image/%s", bimg.ImageTypes[format]))
	w.Write(buffer)
}
