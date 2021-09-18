// Package tile represents an unique SMS tile (8x8 pixels)
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
	"fmt"
	"image"
	"image/color"

	"github.com/disintegration/imaging"
)

const Size = 8 // tile size in pixels

type Tile struct {
	// SMS tile data in planar format
	planarData [32]uint8

	// original image (tile) data; location and colour data
	info         Info
	orientations map[Orientation]image.Image

	// duplicate tiles located in the image (based on their RGBA colours);
	// exact match, along with vertically and horizontally flipped
	duplicates []Info
}

func NewNormalOrientation(info Info, tileImage image.Image) *Tile {
	t := Tile{
		info:         info,
		orientations: make(map[Orientation]image.Image, 4),
	}
	t.orientations[OrientationNormal] = tileImage
	return &t
}

func NewWithOrientations(info Info, tileImage image.Image) *Tile {
	t := NewNormalOrientation(info, tileImage)
	t.generateFlippedOrientations()
	return t
}

func (t Tile) OrientationAt(y, x int, orientation Orientation) (color.Color, error) {
	o, ok := t.orientations[orientation]
	if !ok {
		return color.NRGBA{}, fmt.Errorf("invalid orientation: %016b", orientation)
	}
	return o.At(x, y), nil
}

// AddDuplicate tile to the duplicates slice.
func (t *Tile) AddDuplicate(info Info) {
	t.duplicates = append(t.duplicates, info)
}

// Info returns the location/orientation info for the current tile.
func (t Tile) Info() Info {
	return t.info
}

// Duplicates returns a slice of the duplicates.
func (t Tile) Duplicates() []Info {
	return t.duplicates
}

// IsDuplicate tests the tile image for matching colours.
// If no match is found, then the image is flipped vertically, horizontally,
// and in both planes, and tested again after each.
func (t Tile) IsDuplicate(tile *Tile) (Orientation, bool) {
	// TODO: use goroutines?
	if t.matchingColours(tile, OrientationNormal) {
		return OrientationNormal, true
	} else if t.matchingColours(tile, OrientationFlippedH) {
		return OrientationFlippedH, true
	} else if t.matchingColours(tile, OrientationFlippedV) {
		return OrientationFlippedV, true
	} else if t.matchingColours(tile, OrientationFlippedVH) {
		return OrientationFlippedVH, true
	}
	return OrientationNormal, false
}

// tests if the pixel colours in two tiles are an exact match
func (t Tile) matchingColours(testTile *Tile, orientation Orientation) bool {
	base := t.orientations[orientation]
	tileX, tileY := base.Bounds().Dx(), base.Bounds().Dy()

	tile := testTile.orientations[OrientationNormal]
	if tile.Bounds().Dx() != tileX || tile.Bounds().Dy() != tileY {
		return false
	}

	for y := 0; y < tileY; y++ {
		for x := 0; x < tileX; x++ {
			tr, tg, tb, ta := tile.At(x, y).RGBA()
			r, g, b, a := base.At(x, y).RGBA()
			if tr != r || tg != g || tb != b || ta != a {
				return false
			}
		}
	}

	return true
}

// Generate image data for each orientation, from the existing image data.
// This increases the tile data size but saves a great deal of processing time.
func (t *Tile) generateFlippedOrientations() {
	img := t.orientations[OrientationNormal]

	// TODO: use goroutines
	t.orientations[OrientationFlippedH] = imaging.FlipH(img)

	flippedV := imaging.FlipV(img)
	t.orientations[OrientationFlippedV] = flippedV
	t.orientations[OrientationFlippedVH] = imaging.FlipH(flippedV)
}
