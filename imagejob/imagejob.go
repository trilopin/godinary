package imagejob

import (
	"errors"
	"image"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/disintegration/imaging"
)

// ImageJob stores info for image processing
type ImageJob struct {
	Image             image.Image
	SourceURL         string
	SourceWidth       int
	SourceHeight      int
	SourceAspectRatio float32
	TargetWidth       int
	TargetHeight      int
	TargetFormat      string
	Filters           map[string]string
}

// NewImageJob constructs a default empty struct and return a pointer to it
func NewImageJob() *ImageJob {
	var job ImageJob
	job.Filters = make(map[string]string)
	job.Filters["crop"] = "scale"
	return &job
}

// Parse creates a Job struct from string
func (img *ImageJob) Parse(fetchData string) error {
	var offset int
	var err error

	parts := strings.SplitN(fetchData, "/", 2)
	if parts[0] != "http:" {
		filters := strings.Split(parts[0], ",")
		for _, v := range filters {
			filter := strings.Split(v, "_")
			switch filter[0] {
			case "h":
				img.TargetHeight, err = strconv.Atoi(filter[1])
				if err != nil {
					return errors.New("TargetHeight is not integer")
				}
			case "w":
				img.TargetWidth, err = strconv.Atoi(filter[1])
				if err != nil {
					return errors.New("TargetWidth is not integer")
				}
			case "f":
				allowed := map[string]bool{
					"jpg":  true,
					"jpeg": true,
					"png":  true,
					"gif":  true,
				}
				if !allowed[filter[1]] {
					return errors.New("Format not allowed")
				}
				img.TargetFormat = filter[1]
			case "c":
				allowed := map[string]bool{
					"limit": true,
					"fit":   true,
					"scale": true,
				}
				if !allowed[filter[1]] {
					return errors.New("Crop not allowed")
				}
				img.Filters["crop"] = filter[1]
			}
		}
		offset = len(parts[0]) + 1
	}
	img.SourceURL, _ = url.QueryUnescape(fetchData[offset:])
	return nil
}

// Download retrieves url into io.Reader
func Download(img *ImageJob) (io.Reader, error) {

	c := &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   2 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   2 * time.Second,
			ResponseHeaderTimeout: 2 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}

	if img.SourceURL == "" {
		return nil, errors.New("SourceURL not found in image")
	}

	resp, err := c.Get(img.SourceURL)
	if err != nil {
		return nil, errors.New("Cannot download image")
	}
	return resp.Body, nil
}

// Decode takes body reader and loads image into struct
func (img *ImageJob) Decode(body io.Reader) error {
	image, err := imaging.Decode(body)
	if err != nil {
		return errors.New("Cannot decode image")
	}
	img.Image = image
	img.extractInfo()
	return nil
}

func (img *ImageJob) extractInfo() error {
	bounds := img.Image.Bounds()
	img.SourceHeight = bounds.Max.Y
	img.SourceWidth = bounds.Max.X
	img.SourceAspectRatio = float32(img.SourceWidth) / float32(img.SourceHeight)
	return nil
}

// crop calculates the best strategy to crop the image
func (img *ImageJob) crop() error {
	switch img.Filters["crop"] {
	// Preserve aspect ratio, bigger dimension is selected
	case "fit":
		if img.TargetHeight > img.TargetWidth {
			img.TargetWidth = 0
		} else {
			img.TargetHeight = 0
		}
		// Same as Fit but limiting size to original image
	case "limit":
		if img.TargetHeight > img.SourceHeight || img.TargetWidth > img.SourceWidth {
			img.TargetWidth = img.SourceWidth
			img.TargetHeight = img.SourceHeight
		} else {
			if img.TargetHeight > img.TargetWidth {
				img.TargetWidth = 0
			} else {
				img.TargetHeight = 0
			}
		}
	}

	return nil
}

// Process transforms image
func (img *ImageJob) Process(writer io.Writer) error {
	if img.Image == nil {
		return errors.New("Image not found")
	}
	// Tweaks height and width parameters (Resize will launch it)
	img.crop()

	log.Printf(
		"%s from %dx%d to %dx%d",
		img.SourceURL,
		img.SourceWidth,
		img.SourceHeight,
		img.TargetWidth,
		img.TargetHeight,
	)
	transformedImg := imaging.Resize(img.Image, img.TargetWidth, img.TargetHeight, imaging.Lanczos)
	imaging.Encode(writer, transformedImg, 0)
	return nil
}
