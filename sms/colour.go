package sms

// Colour represents a single RGB colour on the SMS.
//
// An SMS colour consists of a single byte, giving a total 64 possible colours.
// Each colour uses two bits from the byte, representing a value between 0 and 3:
//
//   Bit:   7 6  |  5 4 |  3 2  | 1 0
//     %: Unused | Blue | Green | Red
type Colour uint8

// FromRGB converts 8-bit RGB values to one of the 64 SMS colours.
// Conversion: 0-52 = 0, 53-127 = 85, 128-202 = 170, 203-255 = 255
func FromRGB(r, g, b uint8) Colour {
	r = matchToNearestSmsColour(r)
	g = matchToNearestSmsColour(g)
	b = matchToNearestSmsColour(b)

	for value, colour := range smsColours {
		if r == colour.r && g == colour.g && b == colour.b {
			return value
		}
	}
	return 0
}

func (c Colour) RGB() (r, g, b uint8) {
	if col, ok := smsColours[c]; ok {
		return col.r, col.g, col.b
	}
	return 0, 0, 0
}

const htmlBlack = "#000000"

func (c Colour) HTML() string {
	if col, ok := smsColours[c]; ok {
		return col.html
	}
	return htmlBlack
}

// Converts an 8-bit colour (e.g. R from RGB) to its nearest SMS equivalent.
//
// A more exact conversion might be:
// 	 0-42 = 0, 43-127 = 85, 128-211 = 170, 212-255 = 255
// TODO: perhaps a more advanced conversion is needed.
func matchToNearestSmsColour(colour uint8) uint8 {
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

type smsPaletteColour struct {
	r, g, b uint8
	html    string
}

// Complete SMS colour palette.
var smsColours = map[Colour]smsPaletteColour{
	0b00000000: {r: 0, g: 0, b: 0, html: "#000000"},
	0b00000001: {r: 85, g: 0, b: 0, html: "#550000"},
	0b00000010: {r: 170, g: 0, b: 0, html: "#AA0000"},
	0b00000011: {r: 255, g: 0, b: 0, html: "#FF0000"},
	0b00000100: {r: 0, g: 85, b: 0, html: "#005500"},
	0b00000101: {r: 85, g: 85, b: 0, html: "#555500"},
	0b00000110: {r: 170, g: 85, b: 0, html: "#AA5500"},
	0b00000111: {r: 255, g: 85, b: 0, html: "#FF5500"},
	0b00001000: {r: 0, g: 170, b: 0, html: "#00AA00"},
	0b00001001: {r: 85, g: 170, b: 0, html: "#55AA00"},
	0b00001010: {r: 170, g: 170, b: 0, html: "#AAAA00"},
	0b00001011: {r: 255, g: 170, b: 0, html: "#FFAA00"},
	0b00001100: {r: 0, g: 255, b: 0, html: "#00FF00"},
	0b00001101: {r: 85, g: 255, b: 0, html: "#55FF00"},
	0b00001110: {r: 170, g: 255, b: 0, html: "#AAFF00"},
	0b00001111: {r: 255, g: 255, b: 0, html: "#FFFF00"},
	0b00010000: {r: 0, g: 0, b: 85, html: "#000055"},
	0b00010001: {r: 85, g: 0, b: 85, html: "#550055"},
	0b00010010: {r: 170, g: 0, b: 85, html: "#AA0055"},
	0b00010011: {r: 255, g: 0, b: 85, html: "#FF0055"},
	0b00010100: {r: 0, g: 85, b: 85, html: "#005555"},
	0b00010101: {r: 85, g: 85, b: 85, html: "#555555"},
	0b00010110: {r: 170, g: 85, b: 85, html: "#AA5555"},
	0b00010111: {r: 255, g: 85, b: 85, html: "#FF5555"},
	0b00011000: {r: 0, g: 170, b: 85, html: "#00AA55"},
	0b00011001: {r: 85, g: 170, b: 85, html: "#55AA55"},
	0b00011010: {r: 170, g: 170, b: 85, html: "#AAAA55"},
	0b00011011: {r: 255, g: 170, b: 85, html: "#FFAA55"},
	0b00011100: {r: 0, g: 255, b: 85, html: "#00FF55"},
	0b00011101: {r: 85, g: 255, b: 85, html: "#55FF55"},
	0b00011110: {r: 170, g: 255, b: 85, html: "#AAFF55"},
	0b00011111: {r: 255, g: 255, b: 85, html: "#FFFF55"},
	0b00100000: {r: 0, g: 0, b: 170, html: "#0000AA"},
	0b00100001: {r: 85, g: 0, b: 170, html: "#5500AA"},
	0b00100010: {r: 170, g: 0, b: 170, html: "#AA00AA"},
	0b00100011: {r: 255, g: 0, b: 170, html: "#FF00AA"},
	0b00100100: {r: 0, g: 85, b: 170, html: "#0055AA"},
	0b00100101: {r: 85, g: 85, b: 170, html: "#5555AA"},
	0b00100110: {r: 170, g: 85, b: 170, html: "#AA55AA"},
	0b00100111: {r: 255, g: 85, b: 170, html: "#FF55AA"},
	0b00101000: {r: 0, g: 170, b: 170, html: "#00AAAA"},
	0b00101001: {r: 85, g: 170, b: 170, html: "#55AAAA"},
	0b00101010: {r: 170, g: 170, b: 170, html: "#AAAAAA"},
	0b00101011: {r: 255, g: 170, b: 170, html: "#FFAAAA"},
	0b00101100: {r: 0, g: 255, b: 170, html: "#00FFAA"},
	0b00101101: {r: 85, g: 255, b: 170, html: "#55FFAA"},
	0b00101110: {r: 170, g: 255, b: 170, html: "#AAFFAA"},
	0b00101111: {r: 255, g: 255, b: 170, html: "#FFFFAA"},
	0b00110000: {r: 0, g: 0, b: 255, html: "#0000FF"},
	0b00110001: {r: 85, g: 0, b: 255, html: "#5500FF"},
	0b00110010: {r: 170, g: 0, b: 255, html: "#AA00FF"},
	0b00110011: {r: 255, g: 0, b: 255, html: "#FF00FF"},
	0b00110100: {r: 0, g: 85, b: 255, html: "#0055FF"},
	0b00110101: {r: 85, g: 85, b: 255, html: "#5555FF"},
	0b00110110: {r: 170, g: 85, b: 255, html: "#AA55FF"},
	0b00110111: {r: 255, g: 85, b: 255, html: "#FF55FF"},
	0b00111000: {r: 0, g: 170, b: 255, html: "#00AAFF"},
	0b00111001: {r: 85, g: 170, b: 255, html: "#55AAFF"},
	0b00111010: {r: 170, g: 170, b: 255, html: "#AAAAFF"},
	0b00111011: {r: 255, g: 170, b: 255, html: "#FFAAFF"},
	0b00111100: {r: 0, g: 255, b: 255, html: "#00FFFF"},
	0b00111101: {r: 85, g: 255, b: 255, html: "#55FFFF"},
	0b00111110: {r: 170, g: 255, b: 255, html: "#AAFFFF"},
	0b00111111: {r: 255, g: 255, b: 255, html: "#FFFFFF"},
}
