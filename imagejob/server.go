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

	raven "github.com/getsentry/raven-go"
	"github.com/spf13/viper"
	"github.com/trilopin/godinary/storage"
	bimg "gopkg.in/h2non/bimg.v1"
)

var (
	// SpecificThrotling is semaphore per domain
	SpecificThrotling map[string]chan struct{}
	// GlobalThrotling is global semaphore
	GlobalThrotling chan struct{}
)

type ServerOpts struct {
	MaxRequest          int
	MaxRequestPerDomain int
	Port                string
	Domain              string
	AllowedReferers     []string
}

// RobotsTXT return robots.txt valid for complete disallow
func RobotsTXT(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "User-Agent: *")
	fmt.Fprintln(w, "Disallow: /")
}

// Up is the health check for application
func Up(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "up")
}

// Fetch takes url + params in url to download image from url and apply filters
func Fetch(opts *ServerOpts) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var reader io.ReadCloser
		var dSem float64

		t1 := time.Now()
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		urlInfo := strings.Replace(r.URL.Path, "/image/fetch/", "", 1)
		job := NewImageJob()
		acceptHeader, ok := r.Header["Accept"]
		job.AcceptWebp = ok && strings.Contains(acceptHeader[0], "image/webp")

		if err := job.Parse(urlInfo); err != nil {
			raven.CaptureErrorAndWait(err, nil)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		domain, err := topDomain(job.Source.URL)
		if err != nil || domain == "" {
			if err != nil {
				raven.CaptureErrorAndWait(err, nil)
			}
			http.Error(w, "Cannot parse domain", http.StatusInternalServerError)
			return
		}
		if _, ok := SpecificThrotling[domain]; !ok {
			SpecificThrotling[domain] = make(chan struct{}, opts.MaxRequestPerDomain)
		}

		// derived image is already cached
		if reader, err = storage.StorageDriver.NewReader(job.Target.Hash, "derived/"); err == nil {
			defer reader.Close()
			if cached, err2 := ioutil.ReadAll(reader); err2 == nil {
				if err = writeImage(w, cached, job.Target.Format); err == nil {
					log.Printf("CACHED - TOTAL %0.5f", time.Since(t1).Seconds())
				} else {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
				return
			}
		}

		// Download if original image does not exists at storage, load otherwise
		reader, err = storage.StorageDriver.NewReader(job.Source.Hash, "source/")
		if err == nil {
			defer reader.Close()
			job.Source.Load(reader)
		} else {
			tSem := time.Now()
			log.Printf("SEM %s %d/%d %d/%d", domain, len(GlobalThrotling), cap(GlobalThrotling), len(SpecificThrotling[domain]), cap(SpecificThrotling[domain]))
			GlobalThrotling <- struct{}{}
			SpecificThrotling[domain] <- struct{}{}
			dSem = time.Since(tSem).Seconds()
			err = job.Source.Download(storage.StorageDriver)
			<-SpecificThrotling[domain]
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

		if err = writeImage(w, job.Target.RawContent, job.Target.Format); err == nil {
			log.Printf(
				"NEW - TOTAL %0.5f => SEM %0.5f, DOWN %0.5f, PROC %0.5f",
				time.Since(t1).Seconds(), dSem,
				t2.Sub(t1).Seconds()-dSem, t3.Sub(t2).Seconds())
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}

}

func topDomain(URL string) (string, error) {
	info, err := url.Parse(URL)
	if err != nil {
		return "", errors.New("Cannot parse hostname")
	}
	return info.Host, nil
}

func writeImage(w http.ResponseWriter, buffer []byte, format bimg.ImageType) error {
	w.Header().Set("Cache-Control", "public, max-age="+viper.GetString("cdn_ttl"))
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer)))
	w.Header().Set("Content-Type", fmt.Sprintf("image/%s", bimg.ImageTypes[format]))
	_, err := w.Write(buffer)

	if err != nil {
		log.Println("Error writing response ", err)
		raven.CaptureErrorAndWait(err, nil)
	}
	return err
}
