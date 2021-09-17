package image

import (
	"fmt"
	"image"
)

type TiledImage struct {
	rows  int       // # of tiled rows
	cols  int       // # of tiles columns
	tiles [][]*Tile // the image stored as a tiled map (rows/cols).

	// a list of unique tiles (based on their RGBA colours);
	// exact match, along with vertically and horizontally flipped.
	uniques []uniqueTile
}

// FromImage converts an image.Image to a tiled a map of 8x8 tiles.
func FromImage(img image.Image) (tiled *TiledImage, err error) {
	rows := img.Bounds().Dy()
	cols := img.Bounds().Dx()
	if cols%Size != 0 || rows%Size != 0 {
		err = fmt.Errorf("image must be multiples of 8px, got %dx%d pixels", cols, rows)
		return
	}

	tiled = &TiledImage{
		rows: rows / Size,
		cols: cols / Size,
	}

	for yOffset := 0; yOffset < rows; yOffset += Size {
		var rowTiles []*Tile
		for xOffset := 0; xOffset < cols; xOffset += Size {
			colTile := New()
			for y := 0; y < Size; y++ {
				for x := 0; x < Size; x++ {
					pixelColour := img.At(xOffset+x, yOffset+y)
					colTile.NRGBA.Set(x, y, pixelColour)
				}
			}
			rowTiles = append(rowTiles, colTile)
		}
		tiled.tiles = append(tiled.tiles, rowTiles)
	}
	return
}

// ToImage converts the tiled image to an RGBA image.
func (t *TiledImage) ToImage() *image.NRGBA {
	rows := t.rows * Size
	cols := t.cols * Size
	img := image.NewNRGBA(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: cols, Y: rows},
	})

	xOffset := 0
	yOffset := 0
	for _, rowTiles := range t.tiles {
		for _, colTile := range rowTiles {
			for y := 0; y < Size; y++ {
				for x := 0; x < Size; x++ {
					colour := colTile.NRGBA.At(x, y)
					img.Set(xOffset+x, yOffset+y, colour)
				}
			}
			xOffset += Size
		}
		xOffset = 0
		yOffset += Size
	}
	return img
}
