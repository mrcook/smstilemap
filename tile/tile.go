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
	"image"
	"image/color"

	"github.com/disintegration/imaging"
)

const Size = 8 // tile size in pixels

type Tile struct {
	// SMS tile data in planar format
	planarData [32]uint8

	// original image (tile) data; location and colour data
	info      Info
	imageData image.Image // Image interface or NRGBA?

	// duplicate tiles located in the image (based on their RGBA colours);
	// exact match, along with vertically and horizontally flipped
	duplicates []Info
}

func NewTile(info Info, img image.Image) Tile {
	return Tile{info: info, imageData: img}
}

// AddDuplicate tile to the duplicates slice.
func (t *Tile) AddDuplicate(info Info) {
	t.duplicates = append(t.duplicates, info)
}

// Info returns the location/orientation info for the current tile.
func (t Tile) Info() Info {
	return t.info
}

// OriginalColorAt : TODO: only for testing? do we really want to keep this?
func (t Tile) OriginalColorAt(x, y int) color.Color {
	return t.imageData.At(x, y)
}

// IsDuplicate tests the tile image for matching colours.
// If no match is found, then the image is flipped vertically, horizontally,
// and in both planes, and tested again after each.
func (t Tile) IsDuplicate(img image.Image) (Orientation, bool) {
	// TODO: run these in goroutines?

	if t.matchingColours(img) {
		return OrientationNormal, true
	} else if t.matchingColours(imaging.FlipH(img)) {
		return OrientationFlippedH, true
	}

	flippedV := imaging.FlipV(img)
	if t.matchingColours(flippedV) {
		return OrientationFlippedV, true
	} else if t.matchingColours(imaging.FlipH(flippedV)) {
		return OrientationFlippedVH, true
	}

	return OrientationNormal, false
}

// tests if the pixel colours in two tiles are an exact match
func (t Tile) matchingColours(testImg image.Image) bool {
	tile := t.imageData
	tileX, tileY := tile.Bounds().Dx(), tile.Bounds().Dy()

	if testImg.Bounds().Dx() != tileX || testImg.Bounds().Dy() != tileY {
		return false
	}

	for y := 0; y < tileY; y++ {
		for x := 0; x < tileX; x++ {
			tr, tg, tb, ta := testImg.At(x, y).RGBA()
			r, g, b, a := tile.At(x, y).RGBA()
			if tr != r || tg != g || tb != b || ta != a {
				return false
			}
		}
	}

	return true
}
