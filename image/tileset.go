package image

import (
	"fmt"
	"image"
)

// Background represents a tiled version of an image, consisting of unique 8x8 tiles.
type Background struct {
	metadata metadata // original image details
	tiles    []Tile   // a set of unique tiles
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
	}
	tiles := imageToTiles(img)
	bg.generateUniqueTileList(tiles)

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
	for _, t := range tiles {
		b.addTile(t.row, t.col, t.image)
	}
}

// add a tile to the current background tiles, either as a new unique tile
// or as a duplicate of an existing tile, when flipped in one of the supported
// vertical/horizontal orientations.
func (b *Background) addTile(row, col int, img image.Image) {
	// iterate over existing tiles and add as duplicate if a match is found
	for i := 0; i < len(b.tiles); i++ {
		t := New(row, col, img)
		if orientation, dupe := b.tiles[i].IsDuplicate(t); dupe {
			b.tiles[i].AddDuplicateInfo(row, col, orientation)
			return
		}
	}

	b.tiles = append(b.tiles, *NewWithOrientations(row, col, img))
}
