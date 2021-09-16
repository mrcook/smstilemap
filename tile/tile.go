// Package tile represents a Sega Master System tile.
//
// All graphics on the Master System are built up from 8×8 pixel tiles.
// Each pixel is a palette index from 0 to 15, i.e. 4 bits.
//
// The tile data is in a planar format, split by tile row. That means that the
// first byte contains the least significant bit, bit 0, of each pixel in the
// top row of the tile. The second byte contains bit 1 of each pixel, the third
// bit 2, and the fourth bit 3. Thus the top eight pixels are represented by
// the first four bytes of data, split by "bitplane". The process is repeated
// for consecutive rows of the tile, producing 32 bytes total.
//
// In the most typical VRAM layout, 14KB of the total 16KB is available for
// tiles; that is enough space for 448 tiles. (With some tricks you can get
// space for a few more.)
package tile

import (
	"image"

	"github.com/disintegration/imaging"
)

const Size = 8

// Tile is an 8x8 pixel image.
type Tile struct {
	NRGBA *image.NRGBA

	data [32]uint8 // tile data in planar format
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
