package tiler

import (
	"fmt"
	"image"
	"image/color"

	"github.com/disintegration/imaging"
)

// Tile is an 8x8 pixel image tile.
type Tile struct {
	tileSize int // normally 8x8 pixels
	*info        // basic data; row, col, and orientation

	// location/orientation data for all duplicate tiles located in the image,
	// based on their RGBA colours; exact match, vertically and horizontally flipped
	duplicates []info

	orientations map[Orientation]image.Image // the tile image data in all its orientations
}

func New(row, col, tileSize int, tileImage image.Image) *Tile {
	t := Tile{
		tileSize:     tileSize,
		info:         &info{row: row, col: col, orientation: OrientationNormal},
		orientations: make(map[Orientation]image.Image, 4),
	}
	t.orientations[OrientationNormal] = tileImage
	return &t
}

// A new tile, with all its different flipped orientations generated
func newWithOrientations(row, col, tileSize int, tileImage image.Image) *Tile {
	t := New(row, col, tileSize, tileImage)
	t.generateFlippedOrientations()
	return t
}

// RowInPixels is the tile row in pixels, as located in the source image.
func (t *Tile) RowInPixels() int {
	return t.row * t.tileSize
}

// ColInPixels is the tile column in pixels, as located in the source image.
func (t *Tile) ColInPixels() int {
	return t.col * t.tileSize
}

func (t *Tile) OrientationAt(y, x int, orientation Orientation) (color.Color, error) {
	o, ok := t.orientations[orientation]
	if !ok {
		return color.NRGBA{}, fmt.Errorf("invalid orientation: %016b", orientation)
	}
	return o.At(x, y), nil
}

// AddDuplicateInfo tile to the duplicates slice.
func (t *Tile) AddDuplicateInfo(row, col int, orientation Orientation) {
	inf := info{row: row, col: col, orientation: orientation}
	t.duplicates = append(t.duplicates, inf)
}

// DuplicateCount returns number of duplicates for the tile.
func (t *Tile) DuplicateCount() int {
	return len(t.duplicates)
}

// GetDuplicateInfo returns the duplicate at the given index number.
func (t *Tile) GetDuplicateInfo(id int) (*info, error) {
	if id >= len(t.duplicates) {
		return nil, fmt.Errorf("tile duplicate index out of range: %d", id)
	}
	return &t.duplicates[id], nil
}

// IsDuplicate tests the tile image for matching colours.
// If no match is found, then the image is flipped vertically, horizontally,
// and in both planes, and tested again after each.
func (t *Tile) IsDuplicate(tile *Tile) (Orientation, bool) {
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
func (t *Tile) matchingColours(testTile *Tile, orientation Orientation) bool {
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
