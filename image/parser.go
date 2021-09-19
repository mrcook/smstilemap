package image

import (
	"image"
)

// Info contains data about the original image
type metadata struct {
	Rows   int // in 8x8 tiles
	Cols   int // in 8x8 tiles
	Width  int // in pixels
	Height int // in pixels
}

// imageTile holds an 8x8 pixel tile from the original image
type imageTile struct {
	row, col int // tile location
	image    *image.NRGBA
}

// imageToTiles converts a pixel based image to a slice of tiles, with
// each tile containing its original location and colour data.
func imageToTiles(img image.Image) (tiles []imageTile) {
	tileBounds := image.Rectangle{Min: image.Point{}, Max: image.Point{X: Size, Y: Size}}

	// the offsets enable moving the 'cursor' to the next tile location
	for rowOffset := 0; rowOffset < img.Bounds().Dy(); rowOffset += Size {
		for colOffset := 0; colOffset < img.Bounds().Dx(); colOffset += Size {
			newTile := imageTile{
				row:   rowOffset / Size,
				col:   colOffset / Size,
				image: image.NewNRGBA(tileBounds),
			}

			// fetch the 8x8 tile colour data
			for y := 0; y < Size; y++ {
				for x := 0; x < Size; x++ {
					colour := img.At(colOffset+x, rowOffset+y)
					newTile.image.Set(x, y, colour)
				}
			}

			tiles = append(tiles, newTile)
		}
	}
	return
}
