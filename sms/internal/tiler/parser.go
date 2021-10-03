package tiler

import (
	"fmt"
	"image"
	"image/color"
)

// imageTile holds an 8x8 pixel tile from the original image
type imageTile struct {
	posX, posY int                    // tile location in rows, cols.
	palette    map[string]color.Color // using the HEX colour value as the key (e.g. "#FF0000")
	image      *image.NRGBA
}

// convertToTiles converts a pixel based image to a slice of tiles, with
// each tile containing its original location and colour data.
func convertToTiles(img image.Image, tileSize int) (tiles []imageTile) {
	tileBounds := image.Rectangle{Min: image.Point{}, Max: image.Point{X: tileSize, Y: tileSize}}

	// the offsets enable moving the 'cursor' to the next tile location
	for rowOffset := 0; rowOffset < img.Bounds().Dy(); rowOffset += tileSize {
		for colOffset := 0; colOffset < img.Bounds().Dx(); colOffset += tileSize {
			newTile := imageTile{
				posX:    rowOffset / tileSize,
				posY:    colOffset / tileSize,
				image:   image.NewNRGBA(tileBounds),
				palette: make(map[string]color.Color),
			}

			// fetch the 8x8 tile colour data
			for y := 0; y < tileSize; y++ {
				for x := 0; x < tileSize; x++ {
					colour := img.At(colOffset+x, rowOffset+y)

					// add the pixel colour to the tile image
					newTile.image.Set(x, y, colour)

					// add the colour to the tile palette
					hex := colourToHex(colour)
					if _, found := newTile.palette[hex]; !found {
						newTile.palette[hex] = colour
					}
				}
			}

			tiles = append(tiles, newTile)
		}
	}
	return
}

// Hex returns the hex "html" representation of the color, as in #ff0080.
// Stolen from github.com/lucasb-eyer/go-colorful
func colourToHex(c color.Color) string {
	r, g, b, _ := c.RGBA()

	// Add 0.5 for rounding
	return fmt.Sprintf("#%02x%02x%02x", uint8(float64(r)*255.0+0.5), uint8(float64(g)*255.0+0.5), uint8(float64(b)*255.0+0.5))
}
