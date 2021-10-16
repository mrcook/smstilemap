package sms

// Orientation is a 16-bit number for use in the named table entry: ---pcvhnnnnnnnnn
// This binary number contains only the Vertical and Horizontal bits flipped so
// that it can be OR-ed directly when generating the tilemap data.
type Orientation uint16

const (
	OrientationNormal    Orientation = 0b0000000000000000
	OrientationFlippedV  Orientation = 0b0000010000000000
	OrientationFlippedH  Orientation = 0b0000001000000000
	OrientationFlippedVH Orientation = 0b0000011000000000
)
