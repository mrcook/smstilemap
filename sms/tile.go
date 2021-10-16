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

// ToPlanarData converts a tile to an SMS planar data slice.
func (t *Tile) ToPlanarData() [planarDataSize]uint8 {
	return [planarDataSize]uint8{} // TODO: implement
}

// Flipped copies the source tile, applying the flipped states.
func (t *Tile) Flipped(word *Word) *Tile {
	if !word.VerticalFlip && !word.HorizontalFlip {
		return t
	}
	flipped := *t // make a copy

	if word.HorizontalFlip {
		for row := 0; row < t.Size(); row++ {
			for col := 0; col < t.Size(); col++ {
				a := col
				b := t.Size() - 1 - col
				if a >= b {
					break
				}
				flipped.pixels[row][a], flipped.pixels[row][b] = flipped.pixels[row][b], flipped.pixels[row][a]
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
				flipped.pixels[a][col], flipped.pixels[b][col] = flipped.pixels[b][col], flipped.pixels[a][col]
			}
		}
	}

	return &flipped
}
