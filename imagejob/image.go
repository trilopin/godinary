package imagejob

import (
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/disintegration/imaging"
)

// Image contains image attributes
type Image struct {
	Width       int
	Height      int
	AspectRatio float32
	Content     image.Image
	URL         string
	Format      string
}

// Download retrieves url into io.Reader
func (img Image) Download() (io.Reader, error) {

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

	if img.URL == "" {
		return nil, errors.New("SourceURL not found in image")
	}

	resp, err := c.Get(img.URL)
	if err != nil {
		return nil, errors.New("Cannot download image")
	}
	return resp.Body, nil
}

// Decode takes body reader and loads image into struct
func (img *Image) Decode(body io.Reader) error {
	image, err := imaging.Decode(body)
	if err != nil {
		return errors.New("Cannot decode image")
	}
	img.Content = image
	img.extractInfo()
	return nil
}

// extractInfo calculates image dimensions and store them
func (img *Image) extractInfo() error {
	bounds := img.Content.Bounds()
	img.Height = bounds.Max.Y
	img.Width = bounds.Max.X
	img.AspectRatio = float32(img.Width) / float32(img.Height)
	return nil
}

//
func Encode(img image.Image, w io.Writer, format string, quality int) error {
	var err error

	if quality <= 0 || quality > 100 {
		quality = 75
	}

	switch format {
	case "jpg", "jpeg":
		var rgba *image.RGBA
		if nrgba, ok := img.(*image.NRGBA); ok {
			if nrgba.Opaque() {
				rgba = &image.RGBA{
					Pix:    nrgba.Pix,
					Stride: nrgba.Stride,
					Rect:   nrgba.Rect,
				}
			}
		}
		if rgba != nil {
			err = jpeg.Encode(w, rgba, &jpeg.Options{Quality: quality})
		} else {
			err = jpeg.Encode(w, img, &jpeg.Options{Quality: quality})
		}

	case "png":
		err = png.Encode(w, img)
	case "gif":
		err = gif.Encode(w, img, &gif.Options{NumColors: 256})
	default:
		err = errors.New("Unsupported format")
	}
	return err
}
