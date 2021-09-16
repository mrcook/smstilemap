package image

import (
	"image"

	"github.com/mrcook/smstilemap/tile"
)

// FindAllUniqueTiles processes the tile list, recording only unique tiles.
func (t *TiledMap) FindAllUniqueTiles() {
	for y, rowTiles := range t.table {
		for x, testTile := range rowTiles {
			t.updateUniques(testTile, x, y)
		}
	}
}

func (t *TiledMap) updateUniques(testTile *tile.Tile, x, y int) {
	loc := location{X: x, Y: y}

	for i, u := range t.uniques {
		baseTile := t.table[u.Location.Y][u.Location.X]
		if t.tileMatches(baseTile, testTile) {
			t.uniques[i].Duplicates = append(t.uniques[i].Duplicates, loc)
			return
		}
	}
	t.uniques = append(t.uniques, uniqueTile{Location: loc})
}

// TODO: run these in goroutines
func (t *TiledMap) tileMatches(baseTile, testTile *tile.Tile) (matched bool) {
	if t.tilePixelsMatch(baseTile, testTile) {
		matched = true
	} else if t.tilePixelsMatch(baseTile, testTile.FlipV()) {
		matched = true
	} else if t.tilePixelsMatch(baseTile, testTile.FlipH()) {
		matched = true
	}
	return
}

// tests if the pixels (colours) in two tiles are an exact match
func (t *TiledMap) tilePixelsMatch(baseTile, testTile *tile.Tile) bool {
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
func (t *TiledMap) UniqueTilesToImage() *image.NRGBA {
	rows := len(t.uniques) / t.cols
	if len(t.uniques)%t.cols > 0 {
		rows += 1 // make up missing row
	}
	rows *= tile.Size
	cols := t.cols * tile.Size

	img := image.NewNRGBA(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: cols, Y: rows},
	})

	xOffset := 0
	yOffset := 0

	for _, uniq := range t.uniques {
		curTile := t.table[uniq.Location.Y][uniq.Location.X]

		for y := 0; y < tile.Size; y++ {
			for x := 0; x < tile.Size; x++ {
				colour := curTile.NRGBA.At(x, y)
				img.Set(xOffset+x, yOffset+y, colour)
			}
		}

		xOffset += tile.Size
		if xOffset >= cols {
			xOffset = 0
			yOffset += tile.Size
		}
	}

	return img
}

func (t *TiledMap) locationToCoords(loc int) (x, y int) {
	x = loc % t.cols
	y = loc / t.cols
	return
}
