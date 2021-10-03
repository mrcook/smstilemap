package sms

import (
	"image"

	"github.com/mrcook/smstilemap/sms/internal/tiler"
)

const (
	maxWidth       = 256
	maxHeight      = 224
	maxColourCount = 64
)

// SMS represents an image used for generating SMS character/palette data.
type SMS struct {
	videoRAM  VRAM
	colourRAM CRAM

	tiledImg *tiler.Tiled // tiled version of the source image
}

// FromImage converts the given image into SMS image data.
func (s *SMS) FromImage(img image.Image) error {
	return s.readImageOntoSMS(img)
}

// ToImage converts the tiled data to a new NRGBA image, with all tiles mapped
// back to their original positions.
func (s *SMS) ToImage() (image.Image, error) {
	return s.convertToImage()
}
