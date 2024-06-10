package ad_image

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
	"os"
)

type Image struct {
	img image.Image
}

// Load loads the image file from the given path.
func Load(fp string) (*Image, error) {
	f, e := os.Open(fp)
	if e != nil {
		return nil, e
	}
	defer f.Close()

	img, e := png.Decode(f)
	return &Image{img: img}, e
}

// GetImage returns the image.
func (i *Image) GetImage() image.Image {
	return i.img
}

func (i *Image) Base64() string {
	var dst bytes.Buffer
	_ = png.Encode(&dst, i.img)
	return base64.StdEncoding.EncodeToString(dst.Bytes())
}
