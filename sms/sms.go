// Package sms can be used to construct a set of objects for converting to SMS
// VDP data (palettes, tiles, tilemap, etc.)
//
// VRAM has an area dedicated to tiles called the Character generator (Sega
// calls tiles 'Characters'), along with the tilemap and SAT.
// CRAM stores two palettes of 16 colours each.
package sms

import (
	"fmt"

	"github.com/mrcook/smstilemap/sms/orientation"

	"github.com/mrcook/smstilemap/sms/internal/tiler"
)

const (
	tileSize        = 8   // SMS tiles are 8x8 pixels
	maxTileCount    = 448 // maximum number of tiles the VDP can store
	maxColourCount  = 64  // maximum colours the SMS supports
	maxScreenWidth  = 256 // screen width in pixels
	maxScreenHeight = 224 // screen height in pixels, only 192px are visible on the SMS
)

type SMS struct {
	// The Character generator (sprite/tile patterns) is 14 KB in size.
	// Each tile occupies 32 bytes, allowing up to 448 unique tiles to be stored.
	characters [maxTileCount]*Tile

	// The Screen Map can hold the positions of the 896 tiles (768 visible) and
	// is 1792 bytes in size. Each entry is 2-bytes wide and contains the address
	// of the tile in the Character generator, along with their attributes.
	nameTable Tilemap

	// Palette of 32 colours (2x16) used for the background and sprite palettes.
	palette Palette

	// The SAT (Sprite Attribute Table) is a 256-byte area in VideoRam that
	// contains an array of all the sprites defined, its entries are
	// similar to the background layer, except each sprite contain two
	// additional values representing the X/Y coordinates (x,y coords, tile ID).
	// NOTE: probably not needed in this library.
	sat [256]uint8

	tiledImg *tiler.Tiled // TODO: should not import tiler
}

// AddTile adds a tile at the next available slot, returning its index position.
func (s *SMS) AddTile(t *Tile) (int, error) {
	for i, chr := range s.characters {
		if chr == nil {
			s.characters[i] = t
			return i, nil
		}
	}
	return 0, fmt.Errorf("no space available")
}

func (s *SMS) TileAt(id int) (*Tile, error) {
	if id > len(s.characters) {
		return nil, fmt.Errorf("invalid tile ID")
	}
	return s.characters[id], nil
}

func (s *SMS) AddTilemapEntry(tileID, row, col int, or orientation.Orientation) {
	word := Word{
		Priority:      false, // set as a background tile
		PaletteSelect: false, // use the background tile palette
		TileNumber:    uint16(tileID),
	}
	word.SetFlippedStateFromOrientation(or)
	s.nameTable.Set(row, col, word)
}
