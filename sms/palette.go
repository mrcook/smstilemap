package sms

// Palette defines which colours we can use. There are actually two palettes;
// one for the background, and one for the sprites. (The sprite palette can be
// used by the background too.) Each palette contains 16 entries.
//
// The first sixteen colours are the background palette, the second sixteen are
// the sprite palette.
// Each pixel of each tile is represented by four bits, giving a number between
// 0 and 15. This number is used to select which colour to use.
//
// Each palette entry is one of the 64 possible colours on the Master System.
// To pick a colour, you must choose a number between 0 and 3 for each of the
// red, green and blue colour channels. Then combine them in a byte:
//
//   Bit:   7 6  |  5 4 |  3 2  | 1 0
//     %: Unused | Blue | Green | Red
//
// So, for example, if there was a little blue, no green and a lot of red, the
// colour would be %00010011.
type Palette struct {
	// TODO: or two palettes? And if so, make Palette a slice not a struct?
	// Background [16]Colour
	// Sprites    [16]Colour
	Colours [32]Colour
}

// PaletteId references one of the possible 32 palette colours.
type PaletteId uint8
