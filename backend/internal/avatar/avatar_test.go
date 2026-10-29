package avatar_test

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"testing"

	"transcendence/backend/internal/avatar"
)

func testImage(w, h int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	img.Set(0, 0, color.RGBA{R: 255, A: 255})
	return img
}

func encode(t *testing.T, format string, img image.Image) []byte {
	t.Helper()
	var buf bytes.Buffer
	var err error
	switch format {
	case "png":
		err = png.Encode(&buf, img)
	case "jpeg":
		err = jpeg.Encode(&buf, img, nil)
	case "gif":
		err = gif.Encode(&buf, img, nil)
	}
	if err != nil {
		t.Fatalf("encode %s: %v", format, err)
	}
	return buf.Bytes()
}

func TestNormalize_AcceptsPNGAndJPEG(t *testing.T) {
	for _, format := range []string{"png", "jpeg"} {
		out, contentType, err := avatar.Normalize(encode(t, format, testImage(256, 256)))
		if err != nil {
			t.Fatalf("%s: unexpected error %v", format, err)
		}
		if contentType != "image/"+format {
			t.Errorf("%s: content type %q", format, contentType)
		}
		if _, got, err := image.DecodeConfig(bytes.NewReader(out)); err != nil || got != format {
			t.Errorf("%s: output is not a valid %s (%v)", format, format, err)
		}
	}
}

// A JPEG with an EXIF block (where cameras put the GPS location) must come back without it.
func TestNormalize_DropsMetadata(t *testing.T) {
	plain := encode(t, "jpeg", testImage(64, 64))
	exif := append([]byte{0xFF, 0xE1, 0x00, 0x0E}, []byte("Exif\x00\x00GPS-DATA")...)
	withExif := append(append(append([]byte{}, plain[:2]...), exif...), plain[2:]...)

	out, _, err := avatar.Normalize(withExif)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	if bytes.Contains(out, []byte("GPS-DATA")) {
		t.Error("metadata survived re-encoding")
	}
}

func TestNormalize_Rejects(t *testing.T) {
	cases := []struct {
		name string
		data []byte
		want error
	}{
		{"not an image", []byte("<svg onload=alert(1)>"), avatar.ErrFormat},
		{"empty", nil, avatar.ErrFormat},
		{"gif", encode(t, "gif", testImage(32, 32)), avatar.ErrFormat},
		{"too wide", encode(t, "png", testImage(avatar.MaxSide+1, 10)), avatar.ErrDimensions},
		{"too tall", encode(t, "png", testImage(10, avatar.MaxSide+1)), avatar.ErrDimensions},
	}
	for _, c := range cases {
		if _, _, err := avatar.Normalize(c.data); !errors.Is(err, c.want) {
			t.Errorf("%s: got %v, want %v", c.name, err, c.want)
		}
	}
}
