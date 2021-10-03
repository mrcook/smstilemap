package sms

import (
	"github.com/mrcook/smstilemap/sms/internal/tiler"
)

// VRAM (Video RAM) has an area dedicated to tiles called the Character
// generator (Sega calls tiles 'Characters'), along with the tilemap and SAT.
type VRAM struct {
	// The Character generator (sprite/tile patterns) is 14 KB in size.
	// Each tile occupies 32 bytes, allowing up to 448 unique tiles to be stored.
	characters [448]Tile

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

func (v *VRAM) addCharacter(i int, tile *tiler.Tile) {
	// add planar data
	v.characters[i] = Tile{} // TODO: generate the data
}

func (v *VRAM) addTilemapEntry(tile *tiler.Tile) {
	word := Word{} // TODO: generate the word
	inf := tile.Info()
	v.nameTable.Set(inf.Row(), inf.Col(), word)

	// add any duplicates
	for i := 0; i < tile.DuplicateCount(); i++ {
		inf, err := tile.GetDuplicateInfo(i)
		if err != nil {
			break // TODO: break?
		}
		word := Word{} // TODO: generate the word
		v.nameTable.Set(inf.Row(), inf.Col(), word)
	}
}
