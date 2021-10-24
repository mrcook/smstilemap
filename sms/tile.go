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
	pixels [tileSize][tileSize]PaletteId
}

// Size returns the size of the tile (8x8).
func (t *Tile) Size() int {
	return tileSize
}

// PaletteIdAt returns the SMS palette ID from the tile for the requested pixel.
func (t *Tile) PaletteIdAt(row, col int) (PaletteId, error) {
	if row >= t.Size() || col >= t.Size() {
		return 0, fmt.Errorf("tile indexing out of bounds, requested (%d,%d), tile size is %d", row, col, t.Size())
	}
	return t.pixels[row][col], nil
}

// SetPaletteIdAt sets the SMS palette ID for the pixel at row/col.
func (t *Tile) SetPaletteIdAt(row, col int, pid PaletteId) error {
	if row >= t.Size() || col >= t.Size() {
		return fmt.Errorf("tile indexing out of bounds, requested (%d,%d), tile size is %d", row, col, t.Size())
	}
	t.pixels[row][col] = pid
	return nil
}

// AsTilemap returns a copy of the tile with any tilemap vertical/horizontal
// flipped states applied.
func (t *Tile) AsTilemap(word *Word) *Tile {
	tile := *t // make a copy

	if word.HorizontalFlip {
		for row := 0; row < t.Size(); row++ {
			for col := 0; col < t.Size(); col++ {
				a := col
				b := t.Size() - 1 - col
				if a >= b {
					break
				}
				tile.pixels[row][a], tile.pixels[row][b] = tile.pixels[row][b], tile.pixels[row][a]
			}
		}
	}
	if word.VerticalFlip {
		for row := 0; row < t.Size(); row++ {
			a := row
			b := t.Size() - 1 - row
			if a >= b {
				break
			}
			for col := 0; col < t.Size(); col++ {
				tile.pixels[a][col], tile.pixels[b][col] = tile.pixels[b][col], tile.pixels[a][col]
			}
		}
	}

	return &tile
}

// Bytes converts the tiles to planar data, returning the result as a slice of bytes.
// Each row of 8 pixels requires 4 bit planes, totalling 4-bytes per row, with 32 in total.
func (t *Tile) Bytes() (data []uint8) {
	for _, rowPixels := range t.pixels {
		// 4 bit planes for each row of 8 pixels
		var planes [4]uint8

		for i, pid := range rowPixels {
			// Assign the bit of each palette ID nibble to the correct bit plane
			// Note: 4 planes = 4-bits of a palette ID number (0-15), assigned as:
			//   bit-0 > plane 0, bit-1 > plane 1, bit-2 > plane 2, bit-3 > plane 3
			// Each tile pixel aligns with the same bit number on each plane.
			for plane := 0; plane < 4; plane++ {
				pixel := uint8(pid >> plane) // shift each bit of the nibble to bit-0 position
				pixel &= 0b00000001          // make sure this is the only set bit
				pixel <<= 7 - i              // shift the bit to the correct pixel location
				planes[plane] |= pixel       // set the bit (pixel location) on the plane
			}
		}

		// add the 4 planes to the slice
		for _, plane := range planes {
			data = append(data, plane)
		}
	}
	return
}
