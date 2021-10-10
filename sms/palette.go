package sms

import "fmt"

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

const paletteSize = 32

// PaletteId references one of the possible 32 palette colours.
type PaletteId uint8

// Palette defines two palettes, each with 16 colours.
type Palette struct {
	colours [paletteSize]entry
}

type entry struct {
	colour  Colour
	enabled bool // required as colours are initialised to black
}

// ColourAt returns the colour stored at the given index position.
func (p *Palette) ColourAt(pos PaletteId) (Colour, error) {
	if pos >= paletteSize {
		return 0, fmt.Errorf("palette index out of bounds, got %d, max value is %d", pos, paletteSize-1)
	}
	if !p.colours[pos].enabled {
		return 0, fmt.Errorf("uninitialised colour")
	}
	return p.colours[pos].colour, nil
}

// SetColourAt sets the palette colour at the given index position.
func (p *Palette) SetColourAt(pos PaletteId, colour Colour) error {
	if pos >= paletteSize {
		return fmt.Errorf("palette index out of bounds, got %d, max value is %d", pos, paletteSize-1)
	}
	p.colours[pos].colour = colour
	p.colours[pos].enabled = true
	return nil
}

// AddColour in the first available slot and return its index position.
// When the palette already contains the colour its position is returned,
// or an error is the palette is full.
func (p *Palette) AddColour(colour Colour) (PaletteId, error) {
	if pos, err := p.PaletteIdFor(colour); err == nil {
		return pos, nil
	}

	for i, _ := range p.colours {
		if !p.colours[i].enabled {
			p.colours[i].colour = colour
			p.colours[i].enabled = true
			return PaletteId(i), nil
		}
	}

	return 0, fmt.Errorf("palette full")
}

// PaletteIdFor returns the position ID for a matching colour.
// If the colour is not found, an error is returned.
func (p *Palette) PaletteIdFor(colour Colour) (PaletteId, error) {
	for i, _ := range p.colours {
		if p.colours[i].enabled && p.colours[i].colour.Equal(colour) {
			return PaletteId(i), nil
		}
	}
	return 0, fmt.Errorf("colour not found")
}
