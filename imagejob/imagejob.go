package imagejob

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"net/url"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
)

// ImageJob manages image transformation
type ImageJob struct {
	Source  Image
	Target  Image
	Hash    string
	Filters map[string]string
}

// NewImageJob constructs a default empty struct and return a pointer to it
func NewImageJob() *ImageJob {
	var job ImageJob
	job.Filters = make(map[string]string)
	job.Filters["crop"] = "scale"
	return &job
}

// Parse creates a Job struct from string
func (job *ImageJob) Parse(fetchData string) error {
	var offset int
	var err error

	h := sha256.New()
	h.Write([]byte(fetchData))
	job.Hash = hex.EncodeToString(h.Sum(nil))

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
				job.Target.Format = filter[1]
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
	return nil
}

// crop calculates the best strategy to crop the image
func (job *ImageJob) crop() error {
	switch job.Filters["crop"] {
	// Preserve aspect ratio, bigger dimension is selected
	case "fit":
		if job.Target.Height > job.Target.Width {
			job.Target.Width = 0
		} else {
			job.Target.Height = 0
		}
		// Same as Fit but limiting size to original image
	case "limit":
		if job.Target.Height > job.Source.Height || job.Target.Width > job.Source.Width {
			job.Target.Width = job.Source.Width
			job.Target.Height = job.Source.Height
		} else {
			if job.Target.Height > job.Target.Width {
				job.Target.Width = 0
			} else {
				job.Target.Height = 0
			}
		}
	}

	return nil
}

// Process transforms image
func (job *ImageJob) Process(writer io.Writer) error {
	if job.Source.Content == nil {
		return errors.New("Image not found")
	}
	// Tweaks height and width parameters (Resize will launch it)
	job.crop()

	log.Printf(
		"%s from %dx%d to %dx%d",
		job.Source.URL,
		job.Source.Width,
		job.Source.Height,
		job.Target.Width,
		job.Target.Height,
	)
	transformedImg := imaging.Resize(
		job.Source.Content,
		job.Target.Width,
		job.Target.Height,
		imaging.Lanczos)
	imaging.Encode(writer, transformedImg, 0)
	return nil
}
