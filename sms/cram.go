package sms

// CRAM (Colour RAM) stores two palettes of 16 colours each. Each entry is
// 6-bit wide and each 2-bit set defines one colour from the RGB model.
// This means that there are 64 colours available to choose from.
type CRAM struct {
	palette1 [16]Colour // background palette
	palette2 [16]Colour // sprite and background palette
}
