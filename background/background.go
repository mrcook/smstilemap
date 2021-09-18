package background

import (
	"fmt"
	"image"

	"github.com/mrcook/smstilemap/tile"
)

// Background is an SMS tile set consisting of unique 8x8 tiles.
type Background struct {
	metadata metadata
	tiles    []tile.Tile
}

// FromImage returns a new Background tile set from the given image data.
func FromImage(img image.Image) *Background {
	bg := Background{
		metadata: metadata{
			Rows:   img.Bounds().Dy() / tile.Size,
			Cols:   img.Bounds().Dx() / tile.Size,
			Width:  img.Bounds().Dx(),
			Height: img.Bounds().Dy(),
		},
	}
	tiles := imageToTiles(img)
	bg.generateUniqueTileList(tiles)

	return &bg
}

// GetTile returns the tile for the given index number.
func (b Background) GetTile(id int) (*tile.Tile, error) {
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
	info := tile.Info{Row: row, Col: col, Orientation: tile.OrientationNormal}

	// iterate over existing tiles and add as duplicate if a match is found
	for i := 0; i < len(b.tiles); i++ {
		inf := tile.Info{Col: col, Row: row, Orientation: tile.OrientationNormal}
		t := tile.NewNormalOrientation(inf, img)
		if orientation, dupe := b.tiles[i].IsDuplicate(t); dupe {
			info.Orientation = orientation
			b.tiles[i].AddDuplicate(info)
			return
		}
	}

	b.tiles = append(b.tiles, *tile.NewWithOrientations(info, img))
}
