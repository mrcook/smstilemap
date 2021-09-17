package image

import (
	"image"
)

type orientation int

const (
	normal   orientation = 1
	flippedV orientation = 2
	flippedH orientation = 3
)

type info struct {
	X, Y        int
	orientation orientation
}

// UniqueTile stores the location of the unique tile along with the locations
// of all duplicate tiles found in the tiled map table.
type uniqueTile struct {
	info       info
	Duplicates []info
}

// FindAllUniqueTiles processes the tile list, recording only unique tiles.
func (t *TiledImage) FindAllUniqueTiles() {
	for y, rowTiles := range t.tiles {
		for x, testTile := range rowTiles {
			t.updateUniques(testTile, x, y)
		}
	}
}

func (t *TiledImage) updateUniques(testTile *Tile, x, y int) {
	inf := info{X: x, Y: y, orientation: normal}

	for i, u := range t.uniques {
		baseTile := t.tiles[u.info.Y][u.info.X]
		if o, matched := t.tileMatches(baseTile, testTile); matched {
			inf.orientation = o
			t.uniques[i].Duplicates = append(t.uniques[i].Duplicates, inf)
			return
		}
	}
	t.uniques = append(t.uniques, uniqueTile{info: inf})
}

// TODO: run these in goroutines
func (t *TiledImage) tileMatches(baseTile, testTile *Tile) (orientation, bool) {
	if t.tilePixelsMatch(baseTile, testTile) {
		return normal, true
	} else if t.tilePixelsMatch(baseTile, testTile.FlipV()) {
		return flippedV, true
	} else if t.tilePixelsMatch(baseTile, testTile.FlipH()) {
		return flippedH, true
	}
	return normal, false
}

// tests if the pixels (colours) in two tiles are an exact match
func (t *TiledImage) tilePixelsMatch(baseTile, testTile *Tile) bool {
	if len(baseTile.NRGBA.Pix) != len(testTile.NRGBA.Pix) {
		return false
	}
	for i, attr := range baseTile.NRGBA.Pix {
		if attr != testTile.NRGBA.Pix[i] {
			return false
		}
	}
	return true
}

// UniqueTilesToImage converts the unique tiles to an RGBA image.
func (t *TiledImage) UniqueTilesToImage() *image.NRGBA {
	rows := len(t.uniques) / t.cols
	if len(t.uniques)%t.cols > 0 {
		rows += 1 // make up missing row
	}
	rows *= Size
	cols := t.cols * Size

	img := image.NewNRGBA(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: cols, Y: rows},
	})

	xOffset := 0
	yOffset := 0

	for _, uniq := range t.uniques {
		curTile := t.tiles[uniq.info.Y][uniq.info.X]

		for y := 0; y < Size; y++ {
			for x := 0; x < Size; x++ {
				colour := curTile.NRGBA.At(x, y)
				img.Set(xOffset+x, yOffset+y, colour)
			}
		}

		xOffset += Size
		if xOffset >= cols {
			xOffset = 0
			yOffset += Size
		}
	}

	return img
}

func (t *TiledImage) locationToCoords(loc int) (x, y int) {
	x = loc % t.cols
	y = loc / t.cols
	return
}
