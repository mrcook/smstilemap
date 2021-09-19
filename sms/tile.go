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

type Tile struct {
	planarData [32]uint8 // TODO: or use [8][4]uint8 ?
}
