package image

import (
	"io"
	"log"

	"github.com/disintegration/imaging"
	"github.com/trilopin/godinary/parser"
)

// Process applies all jobs to image contained into buffer reader an
// produces a new image
func Process(buf io.Reader, img parser.ImageJob, writer io.Writer) error {
	image, err := imaging.Decode(buf)
	if err != nil {
		log.Fatal("Cannot decode image")
		return err
	}
	image = imaging.Resize(image, img.TargetWidth, img.TargetHeight, imaging.Lanczos)
	imaging.Encode(writer, image, 0)
	return nil
}
