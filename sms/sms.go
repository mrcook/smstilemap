package sms

import (
	"image"

	"github.com/mrcook/smstilemap/sms/internal/tiler"
)

const (
	tileSize       = 8   // 8x8 pixel tiles
	maxWidth       = 256 // screen width in pixels
	maxHeight      = 224 // screen height in pixels, only 192px are visible
	maxColourCount = 64  // maximum system colours
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
