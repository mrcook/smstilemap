package sms

// All graphics on the Master System are built up from 8×8 pixel tiles.
// Each pixel is a palette index from 0 to 15, i.e. 4 bits.
//
// The tile data is in a planar format, split by tile row. That means that the
// first byte contains the least significant bit, bit 0, of each pixel in the
// top row of the tile. The second byte contains bit 1 of each pixel, the third
// bit 2, and the fourth bit 3. Thus the top eight pixels are represented by
// the first four bytes of data, split by "bitplane". The process is repeated
// for consecutive rows of the tile, producing 32 bytes total.
//
// In the most typical VRAM layout, 14KB of the total 16KB is available for
// tiles; that is enough space for 448 tiles. (With some tricks you can get
// space for a few more.)

// Tile is a type holding the colour data for an 8x8 pixel tile
type Tile struct {
	pixels [64]PaletteId // TODO: or make it an [8][8]PaletteId slice?

	palette *Palette // TODO: store here are just pass references in the funcs?
}

// ToPlanarData converts a tile to an SMS planar data slice.
func (t *Tile) ToPlanarData() [32]uint8 {
	return [32]uint8{} // TODO: implement
}
