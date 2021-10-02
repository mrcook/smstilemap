package sms

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

// CRAM (Colour RAM) stores two palettes of 16 colours each. Each entry is
// 6-bit wide and each 2-bit set defines one colour from the RGB model.
// This means that there are 64 colours available to choose from.
type CRAM struct {
	palette1 [16]Colour // background palette
	palette2 [16]Colour // sprite and background palette
}
