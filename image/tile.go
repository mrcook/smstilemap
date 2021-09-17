package image

import (
	"image"

	"github.com/disintegration/imaging"
)

const Size = 8

// Tile is an 8x8 pixel image.
type Tile struct {
	NRGBA *image.NRGBA
}

// New returns a new tile with of size 8x8
func New() *Tile {
	return &Tile{
		NRGBA: image.NewNRGBA(image.Rectangle{
			Min: image.Point{X: 0, Y: 0},
			Max: image.Point{X: Size, Y: Size},
		}),
	}
}

// FlipV returns a new tile, of the current tile flipped vertically
func (t *Tile) FlipV() *Tile {
	return &Tile{NRGBA: imaging.FlipV(t.NRGBA)}
}

// FlipH returns a new tile, of the current tile flipped horizontally
func (t *Tile) FlipH() *Tile {
	return &Tile{NRGBA: imaging.FlipH(t.NRGBA)}
}
