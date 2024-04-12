package sms

// ColourDataForColour returns the colour data for the requested colour.
func ColourDataForColour(c Colour) ColourData {
	for _, colour := range AllColours {
		if c.Equal(colour.Index) {
			return colour
		}
	}
	return AllColours[0] // defaults to black
}

// ColourDataForRGB returns the colour data for the requested RGB values.
func ColourDataForRGB(r, g, b uint8) ColourData {
	for _, colour := range AllColours {
		if r == colour.R && g == colour.G && b == colour.B {
			return colour
		}
	}
	return AllColours[0] // defaults to black
}

// ColourDataForNearestRGB tries to match the given 8-bit RGB values to the
// closest SMS colour data.
func ColourDataForNearestRGB(r, g, b uint8) ColourData {
	return ColourDataForRGB(matchToNearest(r), matchToNearest(g), matchToNearest(b))
}

// Converts an 8-bit colour (e.g. R from RGB) to its nearest SMS equivalent.
//
// A more exact conversion might be:
//
//	0-42 = 0, 43-127 = 85, 128-211 = 170, 212-255 = 255
//
// TODO: perhaps a more advanced conversion is needed.
func matchToNearest(colour uint8) uint8 {
	if colour < 53 {
		return 0
	} else if colour < 128 {
		return 85
	} else if colour < 203 {
		return 170
	} else {
		return 255
	}
}

type ColourData struct {
	Index   Colour
	R, G, B uint8
	HTML    string
}

// AllColours is a list of all 64 Master System colours.
var AllColours = []ColourData{
	{0b00000000, 0, 0, 0, "#000000"},
	{0b00000001, 85, 0, 0, "#550000"},
	{0b00000010, 170, 0, 0, "#AA0000"},
	{0b00000011, 255, 0, 0, "#FF0000"},
	{0b00000100, 0, 85, 0, "#005500"},
	{0b00000101, 85, 85, 0, "#555500"},
	{0b00000110, 170, 85, 0, "#AA5500"},
	{0b00000111, 255, 85, 0, "#FF5500"},
	{0b00001000, 0, 170, 0, "#00AA00"},
	{0b00001001, 85, 170, 0, "#55AA00"},
	{0b00001010, 170, 170, 0, "#AAAA00"},
	{0b00001011, 255, 170, 0, "#FFAA00"},
	{0b00001100, 0, 255, 0, "#00FF00"},
	{0b00001101, 85, 255, 0, "#55FF00"},
	{0b00001110, 170, 255, 0, "#AAFF00"},
	{0b00001111, 255, 255, 0, "#FFFF00"},
	{0b00010000, 0, 0, 85, "#000055"},
	{0b00010001, 85, 0, 85, "#550055"},
	{0b00010010, 170, 0, 85, "#AA0055"},
	{0b00010011, 255, 0, 85, "#FF0055"},
	{0b00010100, 0, 85, 85, "#005555"},
	{0b00010101, 85, 85, 85, "#555555"},
	{0b00010110, 170, 85, 85, "#AA5555"},
	{0b00010111, 255, 85, 85, "#FF5555"},
	{0b00011000, 0, 170, 85, "#00AA55"},
	{0b00011001, 85, 170, 85, "#55AA55"},
	{0b00011010, 170, 170, 85, "#AAAA55"},
	{0b00011011, 255, 170, 85, "#FFAA55"},
	{0b00011100, 0, 255, 85, "#00FF55"},
	{0b00011101, 85, 255, 85, "#55FF55"},
	{0b00011110, 170, 255, 85, "#AAFF55"},
	{0b00011111, 255, 255, 85, "#FFFF55"},
	{0b00100000, 0, 0, 170, "#0000AA"},
	{0b00100001, 85, 0, 170, "#5500AA"},
	{0b00100010, 170, 0, 170, "#AA00AA"},
	{0b00100011, 255, 0, 170, "#FF00AA"},
	{0b00100100, 0, 85, 170, "#0055AA"},
	{0b00100101, 85, 85, 170, "#5555AA"},
	{0b00100110, 170, 85, 170, "#AA55AA"},
	{0b00100111, 255, 85, 170, "#FF55AA"},
	{0b00101000, 0, 170, 170, "#00AAAA"},
	{0b00101001, 85, 170, 170, "#55AAAA"},
	{0b00101010, 170, 170, 170, "#AAAAAA"},
	{0b00101011, 255, 170, 170, "#FFAAAA"},
	{0b00101100, 0, 255, 170, "#00FFAA"},
	{0b00101101, 85, 255, 170, "#55FFAA"},
	{0b00101110, 170, 255, 170, "#AAFFAA"},
	{0b00101111, 255, 255, 170, "#FFFFAA"},
	{0b00110000, 0, 0, 255, "#0000FF"},
	{0b00110001, 85, 0, 255, "#5500FF"},
	{0b00110010, 170, 0, 255, "#AA00FF"},
	{0b00110011, 255, 0, 255, "#FF00FF"},
	{0b00110100, 0, 85, 255, "#0055FF"},
	{0b00110101, 85, 85, 255, "#5555FF"},
	{0b00110110, 170, 85, 255, "#AA55FF"},
	{0b00110111, 255, 85, 255, "#FF55FF"},
	{0b00111000, 0, 170, 255, "#00AAFF"},
	{0b00111001, 85, 170, 255, "#55AAFF"},
	{0b00111010, 170, 170, 255, "#AAAAFF"},
	{0b00111011, 255, 170, 255, "#FFAAFF"},
	{0b00111100, 0, 255, 255, "#00FFFF"},
	{0b00111101, 85, 255, 255, "#55FFFF"},
	{0b00111110, 170, 255, 255, "#AAFFFF"},
	{0b00111111, 255, 255, 255, "#FFFFFF"},
}
