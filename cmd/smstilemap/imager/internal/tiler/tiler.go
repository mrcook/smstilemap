// Package tiler converts a standard image.Image to a tiled representation.
// Tiles are read as 8x8 pixel images starting a the top-left of the image (0,0).
//
// This package is currently marked as 'internal' as it has very rough and
// untested code. Once this has been improved, and its API has been better
// defined, it may be moved to a public package.
package tiler

import (
	"fmt"
	"image"
	"image/color"
)

// Tiled represents a tiled version of a the original image, consisting of unique 8x8 tiles.
type Tiled struct {
	tileSize int // width & height of tile (usually 8x8 pixels)
	width    int // image width in pixels
	height   int // image height in pixels
	rows     int // image row count (in 8x8 tiles)
	cols     int // image column count (in 8x8 tiles)

	tiles   []Tile  // a set of unique tiles making up the image
	palette Palette // map of unique colours found in the image
}

// FromImage returns a new tile set from the given image data.
// The tile size is the width/height of a tile in pixels, and must in be multiples of 8px.
func FromImage(img image.Image, tileSize int) *Tiled {
	if tileSize == 0 || tileSize%8 != 0 {
		tileSize = 8
	}

	bg := Tiled{
		tileSize: tileSize,
		rows:     img.Bounds().Dy() / tileSize,
		cols:     img.Bounds().Dx() / tileSize,
		width:    img.Bounds().Dx(),
		height:   img.Bounds().Dy(),
		palette:  make(map[string]color.Color),
	}

	tiles := convertToTiles(img, tileSize)
	bg.generateUniqueTileList(tiles)

	return &bg
}

func (b *Tiled) Width() int {
	return b.width
}

func (b *Tiled) Height() int {
	return b.height
}

// GetTile returns the tile for the given index number.
func (b *Tiled) GetTile(id int) (*Tile, error) {
	if id >= b.TileCount() {
		return nil, fmt.Errorf("background tile index out of range: %d", id)
	}
	return &b.tiles[id], nil
}

// TileCount is the total number of unique tiles in the background image.
func (b *Tiled) TileCount() int {
	return len(b.tiles)
}

// ColourCount is the total number of unique colours in the image.
func (b *Tiled) ColourCount() int {
	return len(b.palette)
}

// processes the tile list, recording all unique tiles, and adding duplicate
// info if the tile is already present.
func (b *Tiled) generateUniqueTileList(tiles []imageTile) {
	for _, tile := range tiles {
		b.addTile(&tile)
	}
}

// add a tile to the current background tiles, either as a new unique tile
// or as a duplicate of an existing tile, when flipped in one of the supported
// vertical/horizontal orientations.
func (b *Tiled) addTile(tile *imageTile) {
	// add as a duplicate if an existing tile match is found
	for i := 0; i < len(b.tiles); i++ {
		if b.addIfDuplicate(i, tile) {
			return
		}
	}

	// if not duplicate found, add as a new tile
	b.addTileColoursToPalette(tile)
	b.tiles = append(b.tiles, *NewWithOrientations(tile.posX, tile.posY, b.tileSize, tile.palette, tile.image))
}

// if a tile is a duplicate, add it to the duplicates list
func (b *Tiled) addIfDuplicate(tileID int, tile *imageTile) bool {
	t := New(tile.posX, tile.posY, b.tileSize, tile.palette, tile.image)
	if orientation, dupe := b.tiles[tileID].IsDuplicate(t); dupe {
		b.tiles[tileID].AddDuplicateInfo(tile.posX, tile.posY, orientation)
		return true
	}
	return false
}

// adds the tile palette to the global palette data
func (b *Tiled) addTileColoursToPalette(tile *imageTile) {
	for k, v := range tile.palette {
		if _, found := b.palette[k]; !found {
			b.palette[k] = v
		}
	}
}
