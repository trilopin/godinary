package parser

import (
	"errors"
	"net/url"
	"strconv"
	"strings"
)

// Job stores info for image processing
type Job struct {
	SourceURL     string
	SourceWidth   int
	SourceHeight  int
	DesiredWidth  int
	DesiredHeight int
	DesiredFormat string
	FixedRatio    bool
	Filters       map[string]string
}

// Parse takes a string from URL and returns structured info as Job
func (job *Job) New(fetchData string) error {
	var offset int
	var err error

	job.Filters = map[string]string{}
	job.FixedRatio = false

	parts := strings.SplitN(fetchData, "/", 2)
	if parts[0] != "http:" {
		filters := strings.Split(parts[0], ",")
		for _, v := range filters {
			filter := strings.Split(v, "_")
			switch filter[0] {
			case "h":
				job.DesiredHeight, err = strconv.Atoi(filter[1])
				if err != nil {
					return errors.New("DesiredHeight is not integer")
				}
			case "w":
				job.DesiredWidth, err = strconv.Atoi(filter[1])
				if err != nil {
					return errors.New("DesiredWidth is not integer")
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
				job.DesiredFormat = filter[1]
			case "c":
				job.FixedRatio = filter[1] == "limit"
			}
		}
		offset = len(parts[0]) + 1
	}
	job.SourceURL, _ = url.QueryUnescape(fetchData[offset:])
	return nil
}
