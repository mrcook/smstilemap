package sms

import (
	"github.com/mrcook/smstilemap/sms/orientation"
)

// VRAM (Video RAM) has an area dedicated to tiles called the Character
// generator (Sega calls tiles 'Characters'), along with the tilemap and SAT.
type VRAM struct {
	// The Character generator (sprite/tile patterns) is 14 KB in size.
	// Each tile occupies 32 bytes, allowing up to 448 unique tiles to be stored.
	characters [maxTileCount]Tile

	// The Screen Map can hold the positions of the 896 tiles (768 visible) and
	// is 1792 bytes in size. Each entry is 2-bytes wide and contains the address
	// of the tile in the Character generator, along with their attributes.
	nameTable Tilemap

	// The SAT (Sprite Attribute Table) is a 256-byte area in VRAM that
	// contains an array of all the sprites defined, its entries are
	// similar to the background layer, except each sprite contain two
	// additional values representing the X/Y coordinates (x,y coords, tile ID).
	// NOTE: probably not needed in this library.
	sat [256]uint8
}

// adds the planar data for each tile
func (v *VRAM) addTile(tileID int) {
	t := Tile{}
	v.characters[tileID] = t
}

func (v *VRAM) addTilemapEntry(tileID, row, col int, or orientation.Orientation) {
	word := Word{
		Priority:      false, // set as a background tile
		PaletteSelect: false, // use the background tile palette
		TileNumber:    uint16(tileID),
	}
	word.SetFlippedStateFromOrientation(or)

	v.nameTable.Set(row, col, word)
}
