package background

import (
	"image"

	"github.com/mrcook/smstilemap/tile"
)

// ImageInfo contains data about the original image
type imageInfo struct {
	width  int // in pixels
	height int // in pixels
}

// imageTile holds an 8x8 pixel tile from the original image
type imageTile struct {
	row, col int // tile location
	image    *image.NRGBA
}

// imageToTiles converts a pixel based image to a slice of tiles, with
// each tile containing its original location and colour data.
func imageToTiles(img image.Image) (tiles []imageTile) {
	tileBounds := image.Rectangle{Min: image.Point{}, Max: image.Point{X: tile.Size, Y: tile.Size}}

	// the offsets enable moving the 'cursor' to the next tile location
	for rowOffset := 0; rowOffset < img.Bounds().Dy(); rowOffset += tile.Size {
		for colOffset := 0; colOffset < img.Bounds().Dx(); colOffset += tile.Size {
			newTile := imageTile{
				row:   rowOffset / tile.Size,
				col:   colOffset / tile.Size,
				image: image.NewNRGBA(tileBounds),
			}

			// fetch the 8x8 tile colour data
			for y := 0; y < tile.Size; y++ {
				for x := 0; x < tile.Size; x++ {
					colour := img.At(colOffset+x, rowOffset+y)
					newTile.image.Set(x, y, colour)
				}
			}

			tiles = append(tiles, newTile)
		}
	}
	return
}
