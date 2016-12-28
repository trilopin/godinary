package parser

import (
	"errors"
	"image"
	"io"
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
	Image        image.Image
	SourceURL    string
	SourceWidth  int
	SourceHeight int
	TargetWidth  int
	TargetHeight int
	TargetFormat string
	FixedRatio   bool
	Filters      map[string]string
}

// New creates a Job struct from string
func (img *ImageJob) New(fetchData string) error {
	var offset int
	var err error

	img.Filters = map[string]string{}
	img.FixedRatio = false

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
				img.FixedRatio = filter[1] == "limit"
			}
		}
		offset = len(parts[0]) + 1
	}
	img.SourceURL, _ = url.QueryUnescape(fetchData[offset:])
	return nil
}

// Download retrieves and decodes remote image
func (img *ImageJob) Download() error {

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
		return errors.New("SourceURL not found in image")
	}

	resp, err := c.Get(img.SourceURL)
	if err != nil {
		return errors.New("Cannot download image")
	}

	image, err := imaging.Decode(resp.Body)
	if err != nil {
		return errors.New("Cannot decode image")
	}
	img.Image = image
	bounds := image.Bounds()
	img.SourceHeight = bounds.Max.Y
	img.SourceWidth = bounds.Max.X
	return nil
}

// Process transforms image
func (img *ImageJob) Process(writer io.Writer) error {
	transformedImg := imaging.Resize(img.Image, img.TargetWidth, img.TargetHeight, imaging.Lanczos)
	imaging.Encode(writer, transformedImg, 0)
	return nil
}
