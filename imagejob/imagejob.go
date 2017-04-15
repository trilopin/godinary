package imagejob

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/url"
	"strconv"
	"strings"

	bimg "gopkg.in/h2non/bimg.v1"
)

// ImageJob manages image transformation
type ImageJob struct {
	Source  Image
	Target  Image
	Filters map[string]string
}

// NewImageJob constructs a default empty struct and return a pointer to it
func NewImageJob() *ImageJob {
	var job ImageJob
	job.Filters = make(map[string]string)
	job.Filters["crop"] = "scale"
	job.Target.Format = bimg.WEBP
	return &job
}

// Parse creates a Job struct from string
func (job *ImageJob) Parse(fetchData string) error {
	var offset int
	var err error

	ht := sha256.New()
	ht.Write([]byte(fetchData))
	job.Target.Hash = hex.EncodeToString(ht.Sum(nil))

	parts := strings.SplitN(fetchData, "/", 2)
	if parts[0] != "http:" {
		filters := strings.Split(parts[0], ",")
		for _, v := range filters {
			filter := strings.Split(v, "_")
			switch filter[0] {
			case "h":
				job.Target.Height, err = strconv.Atoi(filter[1])
				if err != nil {
					return errors.New("TargetHeight is not integer")
				}
			case "w":
				job.Target.Width, err = strconv.Atoi(filter[1])
				if err != nil {
					return errors.New("TargetWidth is not integer")
				}
			case "q":
				job.Target.Quality, err = strconv.Atoi(filter[1])
				if err != nil {
					return errors.New("Quality is not integer")
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
				job.Target.Format = bimg.WEBP
			case "c":
				allowed := map[string]bool{
					"limit": true,
					"fit":   true,
					"scale": true,
				}
				if !allowed[filter[1]] {
					return errors.New("Crop not allowed")
				}
				job.Filters["crop"] = filter[1]
			}
		}
		offset = len(parts[0]) + 1
	}
	job.Source.URL, _ = url.QueryUnescape(fetchData[offset:])

	hs := sha256.New()
	hs.Write([]byte(job.Source.URL))
	job.Source.Hash = hex.EncodeToString(hs.Sum(nil))

	return nil
}

// crop calculates the best strategy to crop the image
func (job *ImageJob) crop() error {

	// reset dimensions
	if job.Target.Width == 0 {
		job.Target.Width = int(float32(job.Target.Height) / job.Source.AspectRatio)
	}
	if job.Target.Height == 0 {
		job.Target.Height = int(float32(job.Target.Width) * job.Source.AspectRatio)
	}

	switch job.Filters["crop"] {
	// Preserve aspect ratio, bigger dimension is selected
	case "fit":
		if job.Target.Height > job.Target.Width {
			job.Target.Width = int(float32(job.Target.Height) / job.Source.AspectRatio)
		} else {
			job.Target.Height = int(float32(job.Target.Width) * job.Source.AspectRatio)
		}
		// Same as Fit but limiting size to original image
	case "limit":
		if job.Target.Height > job.Source.Height || job.Target.Width > job.Source.Width {
			job.Target.Width = job.Source.Width
			job.Target.Height = job.Source.Height
		} else {
			if job.Target.Height > job.Target.Width {
				job.Target.Width = int(float32(job.Target.Height) / job.Source.AspectRatio)
			} else {
				job.Target.Height = int(float32(job.Target.Width) * job.Source.AspectRatio)
			}
		}
	}

	return nil
}
