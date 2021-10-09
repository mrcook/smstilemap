package sms

// Colour RAM stores two palettes of 16 colours each.
//
// The first sixteen colours are the background palette and the second sixteen
// are the sprite palette, which can also be used by the background tiles. Each
// entry is 6-bits wide and each 2-bit pair defines one colour from the RGB
// model, allowing for a possible 64 colours.
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

// Palette defines two palettes, each with 16 colours.
type Palette struct {
	// TODO: a single slice or two separate? If single, maybe change Palette from a struct to a slice?
	// palette1 [16]Colour // background palette
	// palette2 [16]Colour // sprite and background palette
	colours [32]Colour
}

// PaletteId references one of the possible 32 palette colours.
type PaletteId uint8
