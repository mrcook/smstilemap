package sms

import "fmt"

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
// In the most typical VideoRam layout, 14KB of the total 16KB is available for
// tiles; that is enough space for 448 tiles. (With some tricks you can get
// space for a few more.)

const (
	tileSize       = 8 // SMS tiles are 8x8 pixels
	planarDataSize = 32
)

// Tile is a type holding the colour data for an 8x8 pixel tile
type Tile struct {
	pixels [tileSize][tileSize]PaletteColourId
}

// SetPixelAt sets a pixel in the tile at row/col with an ID from the colour palette.
func (t *Tile) SetPixelAt(row, col int, colour PaletteColourId) error {
	if row >= tileSize || col >= tileSize {
		return fmt.Errorf("tile indexing out of bounds, requested (%d,%d), tile size is %d", row, col, tileSize)
	}
	t.pixels[row][col] = colour
	return nil
}

// PixelAt gets the palette colour from the tile at row/col.
func (t *Tile) PixelAt(row, col int) (PaletteColourId, error) {
	if row >= tileSize || col >= tileSize {
		return 0, fmt.Errorf("tile indexing out of bounds, requested (%d,%d), tile size is %d", row, col, tileSize)
	}
	return t.pixels[row][col], nil
}

// ToPlanarData converts a tile to an SMS planar data slice.
func (t *Tile) ToPlanarData() [planarDataSize]uint8 {
	return [planarDataSize]uint8{} // TODO: implement
}
