// Package sms can be used to construct a set of objects for converting to SMS
// VDP data (palettes, tiles, tilemap, etc.)
//
// VRAM has an area dedicated to tiles called the Character generator (Sega
// calls tiles 'Characters'), along with the tilemap and SAT.
// CRAM stores two palettes of 16 colours each.
package sms

import (
	"fmt"
)

const (
	MaxScreenWidth  = 256 // screen width in pixels
	MaxScreenHeight = 224 // screen height in pixels, only 192px are visible on the SMS
	MaxColourCount  = 64  // maximum colours the SMS supports

	maxTileCount = 448 // maximum number of tiles the VDP can store
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

// AddTilemapEntryAt adds the tile info to the tilemap at the requested location.
func (s *SMS) AddTilemapEntryAt(row, col int, word Word) error {
	return s.nameTable.Set(row, col, word)
}

// TilemapEntryAt returns the tile info from the tilemap for the requested location.
func (s *SMS) TilemapEntryAt(row, col int) (*Word, error) {
	return s.nameTable.Get(row, col)
}

func (s *SMS) AddPaletteColour(colour Colour) (int, error) {
	return s.palette.AddColour(colour)
}
