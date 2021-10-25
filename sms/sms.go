// Package sms can be used to construct a set of objects for converting to SMS
// VDP data (palettes, tiles, tilemap, etc.)
//
// * VRAM has an area dedicated to tiles called the Character generator
//   (Sega calls tiles 'Characters'), along with the tilemap and SAT.
// * CRAM stores two palettes of 16 colours each. Accessed using a base address
//   of $C000.
//
// A typical memory map for the VRAM is as follows (noting that the screen
// display and sprite information tables can be moved to any desired point):
//
//   $C020 ---------------------------------------------------------------
//         Palette data: 2 palettes of 16 colours each.
//   $C000 ---------------------------------------------------------------
//
//   $4000 ---------------------------------------------------------------
//         Sprite Attribute Table: array of all the defined sprites
//   $3F00 ---------------------------------------------------------------
//         Tilemap (name table): 32x28 table of tile IDs & attributes
//   $3800 ---------------------------------------------------------------
//         Sprite/tile patterns, 256..447
//   $2000 ---------------------------------------------------------------
//         Sprite/tile patterns, 0..255
//   $0000 ---------------------------------------------------------------
package sms

import (
	"fmt"
)

const (
	ScreenWidth          = 256 // screen width in pixels
	ScreenHeight         = 192 // screen height in pixels
	ExtendedScreenHeight = 224 // extended 'mode 4' screen height in pixels on the SMS
	MaxColourCount       = 64  // maximum colours the SMS supports
	MaxTileCount         = 448 // maximum number of tiles the VDP can store
)

type SMS struct {
	// The Character generator (sprite/tile patterns) is 14 KB in size.
	// Each tile occupies 32 bytes, allowing up to 448 unique tiles to be stored.
	characters [MaxTileCount]*Tile

	// The Screen Map can hold the positions of the 786 tiles (896 in the
	// extended mode 4) and is 1792 bytes in size. Each entry is 2-bytes wide
	// and contains the address of the tile in the Character generator, along
	// with their attributes.
	nameTable Tilemap

	// The SAT (Sprite Attribute Table) is a 256-byte area in VideoRam that
	// contains an array of all the sprites defined, its entries are
	// similar to the background layer, except each sprite contain two
	// additional values representing the X/Y coordinates (x,y coords, tile ID).
	// For the majority of cases the table is stored at VRAM address $3F00.
	// NOTE: probably not needed in this library.
	sat [256]uint8

	// Palette of 32 colours (2x16) used for the background and sprite palettes.
	// Accessed using a base address of $C000.
	palette Palette
}

// WidthInPixels returns the maximum screen width in pixels.
func (s *SMS) WidthInPixels() int {
	return ScreenWidth
}

// WidthInTiles returns the maximum screen width calculated as 8x8 tiles.
func (s *SMS) WidthInTiles() int {
	return s.nameTable.Width()
}

// HeightInPixels returns the maximum screen height in pixels.
func (s *SMS) HeightInPixels() int {
	return ScreenHeight
}

// ExtendedHeightInPixels returns the extended mode 4 screen height in pixels.
func (s *SMS) ExtendedHeightInPixels() int {
	return ExtendedScreenHeight
}

// HeightInTiles returns the maximum screen height calculated as 8x8 tiles.
func (s *SMS) HeightInTiles() int {
	return s.nameTable.Height()
}

// ExtendedHeightInTiles returns the extended mode 4 screen height calculated as 8x8 tiles.
func (s *SMS) ExtendedHeightInTiles() int {
	return s.nameTable.ExtendedHeight()
}

// TileAt returns a reference to the character generator tile using the given ID.
func (s *SMS) TileAt(tileId uint16) (*Tile, error) {
	if int(tileId) >= len(s.characters) {
		return nil, fmt.Errorf("invalid tile ID")
	}
	return s.characters[tileId], nil
}

// AddTile adds a tile at the next available slot, returning its index position.
func (s *SMS) AddTile(t *Tile) (uint16, error) {
	for i, chr := range s.characters {
		if chr == nil {
			s.characters[i] = t
			return uint16(i), nil
		}
	}
	return 0, fmt.Errorf("tile memory full")
}

// TilemapEntryAt returns the tile info from the tilemap for the requested location.
func (s *SMS) TilemapEntryAt(row, col int) (*Word, error) {
	return s.nameTable.Get(row, col)
}

// AddTilemapEntryAt adds the tile info to the tilemap at the requested location.
func (s *SMS) AddTilemapEntryAt(row, col int, word Word) error {
	return s.nameTable.Set(row, col, word)
}

// PaletteColour returns the colour for the given palette ID.
func (s *SMS) PaletteColour(id PaletteId) (Colour, error) {
	return s.palette.ColourAt(id)
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

// TileData returns all tiles as a slice of bytes.
func (s *SMS) TileData() (data []uint8) {
	for _, tile := range s.characters {
		if tile == nil {
			continue
		}
		for _, b := range tile.Bytes() {
			data = append(data, b)
		}
	}
	return
}

// TilemapData returns the tilemap data as a slice of 16-bit words.
func (s *SMS) TilemapData() (data []uint16) {
	return s.nameTable.Words()
}

// PaletteData returns the palette data as a slice of bytes.
func (s *SMS) PaletteData() (data [32]uint8) {
	return s.palette.Bytes()
}
