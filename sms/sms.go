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
	MaxScreenWidth         = 256 // screen width in pixels
	MaxScreenHeight        = 224 // screen height in pixels
	MaxVisibleScreenHeight = 192 // visible screen height in pixels on the SMS
	MaxColourCount         = 64  // maximum colours the SMS supports
	MaxTileCount           = 448 // maximum number of tiles the VDP can store
)

type SMS struct {
	// The Character generator (sprite/tile patterns) is 14 KB in size.
	// Each tile occupies 32 bytes, allowing up to 448 unique tiles to be stored.
	characters [MaxTileCount]*Tile

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

func (s *SMS) WidthInPixels() int {
	return MaxScreenWidth
}

func (s *SMS) WidthInTiles() int {
	return s.nameTable.Width()
}

func (s *SMS) HeightInPixels() int {
	return MaxScreenHeight
}

func (s *SMS) VisibleHeightInPixels() int {
	return MaxVisibleScreenHeight
}

func (s *SMS) HeightInTiles() int {
	return s.nameTable.Height()
}

func (s *SMS) VisibleHeightInTiles() int {
	return s.nameTable.VisibleHeight()
}

func (s *SMS) TileAt(id uint16) (*Tile, error) {
	if int(id) >= len(s.characters) {
		return nil, fmt.Errorf("invalid tile ID")
	}
	return s.characters[id], nil
}

// AddTile adds a tile at the next available slot, returning its index position.
func (s *SMS) AddTile(t *Tile) (uint16, error) {
	for i, chr := range s.characters {
		if chr == nil {
			s.characters[i] = t
			return uint16(i), nil
		}
	}
	return 0, fmt.Errorf("no space available")
}

// TilemapEntryAt returns the tile info from the tilemap for the requested location.
func (s *SMS) TilemapEntryAt(row, col int) (*Word, error) {
	return s.nameTable.Get(row, col)
}

// AddTilemapEntryAt adds the tile info to the tilemap at the requested location.
func (s *SMS) AddTilemapEntryAt(row, col int, word Word) error {
	return s.nameTable.Set(row, col, word)
}

// PaletteIdForColour returns the palette index position for the requested colour.
// If the colour is not found, an error is returned.
func (s *SMS) PaletteIdForColour(colour Colour) (PaletteId, error) {
	return s.palette.PaletteIdFor(colour)
}

// AddPaletteColour in the first available palette slot and return its index position.
// When the palette already contains the colour, its position is returned.
// An error is returned when the palette is full.
func (s *SMS) AddPaletteColour(colour Colour) (PaletteId, error) {
	return s.palette.AddColour(colour)
}

// PaletteColour returns the colour for the given palette ID.
func (s *SMS) PaletteColour(id PaletteId) (Colour, error) {
	return s.palette.ColourAt(id)
}
