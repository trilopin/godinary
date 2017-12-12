package image

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	bimg "gopkg.in/h2non/bimg.v1"
)

// Hasher interface
type Hasher interface {
	Hash(s string) string
}

// Sha256 hasher
type Sha256 struct {
}

// Hash return a sha256 hash os input string
func (hasher *Sha256) Hash(s string) string {
	ht := sha256.New()
	ht.Write([]byte(s))
	return hex.EncodeToString(ht.Sum(nil))
}

// Job manages image transformation
type Job struct {
	Source     Image
	Target     Image
	Filters    map[string]string
	AcceptWebp bool
	Hasher     Hasher
}

// NewJob constructs a default empty struct and return a pointer to it
func NewJob() *Job {
	var job Job
	job.Filters = make(map[string]string)
	job.Filters["crop"] = "scale"
	job.Target.Format = bimg.JPEG
	job.Hasher = &Sha256{}
	return &job
}

func parseFormat(format string, acceptWebp bool) (bimg.ImageType, error) {
	switch format {
	case "jpg", "jpeg":
		return bimg.JPEG, nil
	case "webp":
		return bimg.WEBP, nil
	case "auto":
		if acceptWebp {
			return bimg.JPEG, nil // WEBP disabled
		}
		return bimg.JPEG, nil
	case "png":
		return bimg.PNG, nil
	case "gif":
		return bimg.GIF, nil
	default:
		return bimg.JPEG, fmt.Errorf("format \"%s\" not allowed", format)
	}
}

func parseCrop(crop string) (string, error) {
	allowed := map[string]bool{
		"limit": true,
		"fit":   true,
		"scale": true,
	}
	if !allowed[crop] {
		return "scale", fmt.Errorf("crop \"%s\" not allowed", crop)
	}
	return crop, nil
}

func (job *Job) parseFilters(s string) error {
	var err error
	filters := strings.Split(s, ",")
	for _, v := range filters {
		filter := strings.Split(v, "_")
		switch filter[0] {
		case "h":
			if job.Target.Height, err = strconv.Atoi(filter[1]); err != nil {
				return fmt.Errorf("targetHeight is not integer: %v", err)
			}
		case "w":
			if job.Target.Width, err = strconv.Atoi(filter[1]); err != nil {
				return fmt.Errorf("targetWidth is not integer: %v", err)
			}
		case "q":
			if job.Target.Quality, err = strconv.Atoi(filter[1]); err != nil {
				return fmt.Errorf("quality is not integer: %v", err)
			}
		case "f":
			if job.Target.Format, err = parseFormat(filter[1], job.AcceptWebp); err != nil {
				return err
			}
		case "c":
			if job.Filters["crop"], err = parseCrop(filter[1]); err != nil {
				return err
			}
		}
	}
	return nil
}

// Parse creates a Job struct from string
func (job *Job) Parse(fetchData string) error {
	var offset int

	parts := strings.SplitN(fetchData, "/", 2)
	if len(parts) > 1 && parts[0] != "http:" && parts[0] != "https:" {
		if err := job.parseFilters(parts[0]); err != nil {
			return err
		}
		offset = len(parts[0]) + 1
	}
	job.Source.URL, _ = url.QueryUnescape(fetchData[offset:])
	// Temporary hack until complete rework of parsing
	if strings.Count(job.Source.URL, "/") == 1 {
		parts = strings.Split(job.Source.URL, "/")
		job.Source.URL = parts[1]
	}
	job.Target.Hash = job.Hasher.Hash(fetchData + string(job.Target.Format))
	job.Source.Hash = job.Hasher.Hash(job.Source.URL)

	return nil
}

// Crop calculates the best strategy to crop the image
func (job *Job) Crop() error {

	// reset dimensions
	switch job.Filters["crop"] {
	// Preserve aspect ratio, bigger dimension is selected
	case "fit":
		if job.Target.Height > job.Target.Width {
			job.Target.Width = int(float32(job.Target.Height) * job.Source.AspectRatio)
		} else {
			job.Target.Height = int(float32(job.Target.Width) / job.Source.AspectRatio)
		}
	// Same as Fit but limiting size to original image
	case "limit":
		if job.Target.Height > job.Source.Height || job.Target.Width > job.Source.Width {
			job.Target.Width = job.Source.Width
			job.Target.Height = job.Source.Height
		} else {
			if job.Target.Height > job.Target.Width {
				job.Target.Width = int(float32(job.Target.Height) * job.Source.AspectRatio)
			} else {
				job.Target.Height = int(float32(job.Target.Width) / job.Source.AspectRatio)
			}
		}
	// do not preserve nothing, respect callers decision
	case "scale":
		if job.Target.Width == 0 {
			job.Target.Width = job.Target.Height
		}
		if job.Target.Height == 0 {
			job.Target.Height = job.Target.Width
		}
	}
	return nil
}
