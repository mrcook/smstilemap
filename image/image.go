package image

import (
	"fmt"
	"image"

	"github.com/mrcook/smstilemap/tile"
)

type TiledMap struct {
	rows  int            // # of tiled rows
	cols  int            // # of tiles columns
	table [][]*tile.Tile // the image stored as a tiled map (rows/cols).

	// a list of uniquely coloured tiles (based on their RGBA colours), including
	// any Vertically and Horizontally flipped tiles.
	uniques []uniqueTile
}

// UniqueTile stores the location of the unique tile along with the locations
// of all duplicate tiles found in the tiled map table.
type uniqueTile struct {
	Location   location
	Duplicates []location
}

type location struct {
	X, Y int
}

// FromImage converts an image.Image to a tiled a map of 8x8 tiles.
func FromImage(img image.Image) (tiled *TiledMap, err error) {
	rows := img.Bounds().Dy()
	cols := img.Bounds().Dx()
	if cols%tile.Size != 0 || rows%tile.Size != 0 {
		err = fmt.Errorf("image must be multiples of 8px, got %dx%d pixels", cols, rows)
		return
	}

	tiled = &TiledMap{
		rows: rows / tile.Size,
		cols: cols / tile.Size,
	}

	for yOffset := 0; yOffset < rows; yOffset += tile.Size {
		var rowTiles []*tile.Tile
		for xOffset := 0; xOffset < cols; xOffset += tile.Size {
			colTile := tile.New()
			for y := 0; y < tile.Size; y++ {
				for x := 0; x < tile.Size; x++ {
					pixelColour := img.At(xOffset+x, yOffset+y)
					colTile.NRGBA.Set(x, y, pixelColour)
				}
			}
			rowTiles = append(rowTiles, colTile)
		}
		tiled.table = append(tiled.table, rowTiles)
	}
	return
}

// ToImage converts the tiled image to an RGBA image.
func (t *TiledMap) ToImage() *image.NRGBA {
	rows := t.rows * tile.Size
	cols := t.cols * tile.Size
	img := image.NewNRGBA(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: cols, Y: rows},
	})

	xOffset := 0
	yOffset := 0
	for _, rowTiles := range t.table {
		for _, colTile := range rowTiles {
			for y := 0; y < tile.Size; y++ {
				for x := 0; x < tile.Size; x++ {
					colour := colTile.NRGBA.At(x, y)
					img.Set(xOffset+x, yOffset+y, colour)
				}
			}
			xOffset += tile.Size
		}
		xOffset = 0
		yOffset += tile.Size
	}
	return img
}
