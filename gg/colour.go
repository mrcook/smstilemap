package gg

// Colour represents a single RGB colour on the Game Gear.
//
// The Game Gear palette is exactly the same as the Master System except:
// each colour channel consists of 2 bytes with four bits per colour, giving
// 4096 possible colours. This requires there to be 64 bytes of colour RAM
// instead of 32.
//
//	Byte:   1                     | 0
//	 Bit:   7 6 | 5 4 | 3 2 | 1 0 | 7 6 | 5 4 | 3 2 | 1 0
//	    :    Unused   |   Blue    |   Green   |    Red
type Colour uint16

// GG returns the Game Gear `word` for the colour.
func (c Colour) GG() uint16 {
	return uint16(c)
}

// RGB returns the RGB values for the colour.
func (c Colour) RGB() (r, g, b uint8) {
	if data := ColourDataForColour(c); data.Index.Equal(c) {
		return data.R, data.G, data.B
	}
	return 0, 0, 0
}

// RGBA implements the Go `color.Color` interface.
func (c Colour) RGBA() (r, g, b, a uint32) {
	cR, cG, cB := c.RGB()

	r = uint32(cR)
	r |= r << 8
	g = uint32(cG)
	g |= g << 8
	b = uint32(cB)
	b |= b << 8
	a = uint32(255)
	a |= a << 8
	return
}

// HTML returns a HTML compatible hex value for the colour.
func (c Colour) HTML() string {
	if data := ColourDataForColour(c); data.Index.Equal(c) {
		return data.HTML
	}
	return AllColours[0].HTML // default to black
}

// Equal compares the given colour and returns true if it matches.
func (c Colour) Equal(colour Colour) bool {
	return c == colour
}
