package tiler

// Orientation is used to indicate what flipped/rotated orientation a duplicate tile has.
type Orientation uint16

const (
	OrientationNormal Orientation = iota
	OrientationFlippedV
	OrientationFlippedH
	OrientationFlippedVH
)
