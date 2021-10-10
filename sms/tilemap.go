package sms

import "fmt"

// Tilemap represents the background graphics on the Master System screen,
// which is 256x224 pixels (32x28 8x8 tiles). This "virtual screen" is slightly
// larger than the viewport (256x224 px/32x24 tiles), allowing the viewport to
// scroll smoothly, with updates to the tilemap happening in the off-screen
// parts. Each entry in the tilemap represents one tile on the virtual screen.
//
// The VideoRam can hold 448 unique tiles (14kb) - under normal usage.
//
// It would take 768 unique tiles to fill the viewport area, or 896 to
// completely fill the tilemap, therefore the tilemap has to be built using a
// repetition of available VideoRam tiles.
//
// Data format:
//
// The data for a tilemap entry uses a total of 13 bits, stored in two bytes:
//
// Bit  |15 14 13|    12    |    11   |      10       |        9        | 8 7 6 5 4 3 2 1 0
// Data | Unused | Priority | Palette | Vertical flip | Horizontal flip |    Tile number
//
// The data is stored in VideoRam (usually at location $3800), in little-endian
// format, and takes up 1792 bytes (32x28x2 bytes).
//
// Flags:
//
// Flipping:
// Vertical and horizontal flipping allows symmetric objects to be created
// with fewer tiles, thereby allowing greater variety in the graphics.
//
// Priority:
// When a tile has its priority bit set, all pixels with index greater than 0
// will be drawn on top of sprites. You must therefore choose a single colour
// in palette position 0 to be the background colour for such tiles, and they
// will have a "blank" background. Careful use of tile priority can make the
// graphics seem more multi-layered.
// https://www.smspower.org/maxim/HowToProgram/Tilemap

const (
	tilemapRows = 28
	tilemapCols = 32
)

// Tilemap represents the background graphics on the Master System screen,
type Tilemap struct {
	table [tilemapRows][tilemapCols]Word
}

func (t *Tilemap) Get(row, col int) (*Word, error) {
	if row >= tilemapRows || col >= tilemapCols {
		return nil, fmt.Errorf("get tilemap out of bounds indexing, max is (%d,%d), requested (%d,%d)", tilemapRows-1, tilemapCols-1, row, col)
	}

	return &t.table[row][col], nil
}

func (t *Tilemap) Set(row, col int, word Word) error {
	if row >= tilemapRows || col >= tilemapCols {
		return fmt.Errorf("set tilemap out of bounds indexing, max is (%d,%d), requested (%d,%d)", tilemapRows-1, tilemapCols-1, row, col)
	}

	t.table[row][col] = word
	return nil
}
