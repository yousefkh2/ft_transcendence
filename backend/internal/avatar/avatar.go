// Package avatar checks uploaded profile pictures and re-encodes them before they are stored.
package avatar

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
)

// MaxSide is the largest width or height accepted. The frontend sends 256 x 256 pictures.
const MaxSide = 1024

var (
	ErrFormat     = errors.New("Please upload a PNG or JPEG image.")
	ErrDimensions = errors.New("The picture can be at most 1024 x 1024 pixels.")
)

// Normalize decodes a PNG or JPEG and encodes it again in the same format. That proves the file is
// a real image and drops metadata such as the GPS location a phone camera writes into photos.
func Normalize(data []byte) (out []byte, contentType string, err error) {
	// The size check runs on the header only, before a huge image is decoded into memory.
	cfg, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil || (format != "png" && format != "jpeg") {
		return nil, "", ErrFormat
	}
	if cfg.Width > MaxSide || cfg.Height > MaxSide {
		return nil, "", ErrDimensions
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, "", ErrFormat
	}
	var buf bytes.Buffer
	if format == "png" {
		err = png.Encode(&buf, img)
	} else {
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 90})
	}
	if err != nil {
		return nil, "", err
	}
	return buf.Bytes(), "image/" + format, nil
}
