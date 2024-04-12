package sms

// Colour represents a single RGB colour on the SMS.
//
// An SMS colour consists of a single byte, giving a total 64 possible colours.
// Each colour uses two bits from the byte, representing a value between 0 and 3:
//
//	Bit:   7 6  |  5 4 |  3 2  | 1 0
//	  %: Unused | Blue | Green | Red
type Colour uint8

// SMS returns the SMS `byte` for the colour.
func (c Colour) SMS() uint8 {
	return uint8(c)
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
