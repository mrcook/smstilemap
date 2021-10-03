// Package tiler converts a standard image.Image to a tiled representation.
// Tiles are read as 8x8 pixel images starting a the top-left of the image (0,0).
//
// Currently this package is marked as 'internal', however, this may change in
// the future once its API has been better defined.
package tiler

import (
	"fmt"
	"image"
	"image/color"
)

// Background represents a tiled version of an image, consisting of unique 8x8 tiles.
// TODO: rename this to something better!!!
type Background struct {
	metadata metadata               // original image details
	tiles    []Tile                 // a set of unique tiles
	palette  map[string]color.Color // map of unique colours found in the image
}

// FromImage returns a new Background tile set from the given image data.
func FromImage(img image.Image) *Background {
	bg := Background{
		metadata: metadata{
			Rows:   img.Bounds().Dy() / Size,
			Cols:   img.Bounds().Dx() / Size,
			Width:  img.Bounds().Dx(),
			Height: img.Bounds().Dy(),
		},
		palette: make(map[string]color.Color, 64),
	}
	tiles := imageToTiles(img)
	bg.generateUniqueTileList(tiles)
	bg.metadata.UniqueColourCount = len(bg.palette)

	return &bg
}

// GetTile returns the tile for the given index number.
func (b Background) GetTile(id int) (*Tile, error) {
	if id >= b.TileCount() {
		return nil, fmt.Errorf("background tile index out of range: %d", id)
	}
	return &b.tiles[id], nil
}

// TileCount is the total number of unique tiles in the background image.
func (b Background) TileCount() int {
	return len(b.tiles)
}

// Info is the background image metadata (width, height, etc.).
func (b Background) Info() *metadata {
	return &b.metadata
}

// processes the tile list, recording all unique tiles, and adding duplicate
// info if the tile is already present.
func (b *Background) generateUniqueTileList(tiles []imageTile) {
	for _, tile := range tiles {
		b.addTile(&tile)
	}
}

// add a tile to the current background tiles, either as a new unique tile
// or as a duplicate of an existing tile, when flipped in one of the supported
// vertical/horizontal orientations.
func (b *Background) addTile(tile *imageTile) {
	// add as a duplicate if an existing tile match is found
	for i := 0; i < len(b.tiles); i++ {
		if b.addIfDuplicate(i, tile) {
			return
		}
	}

	// if not duplicate found, add as a new tile
	b.addTileColoursToPalette(tile)
	b.tiles = append(b.tiles, *newWithOrientations(tile.row, tile.col, tile.image))
}

// if a tile is a duplicate, add it to the duplicates list
func (b *Background) addIfDuplicate(tileID int, tile *imageTile) bool {
	t := New(tile.row, tile.col, tile.image)
	if orientation, dupe := b.tiles[tileID].IsDuplicate(t); dupe {
		b.tiles[tileID].AddDuplicateInfo(tile.row, tile.col, orientation)
		return true
	}
	return false
}

// adds the tile palette to the global palette data
func (b *Background) addTileColoursToPalette(tile *imageTile) {
	for k, v := range tile.palette {
		if _, found := b.palette[k]; !found {
			b.palette[k] = v
		}
	}
}
