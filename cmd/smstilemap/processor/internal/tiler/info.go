package tiler

import "github.com/mrcook/smstilemap/sms/orientation"

type info struct {
	col, row    int // of the tile as used in the tilemap
	orientation orientation.Orientation
}

func (i *info) Row() int {
	return i.row
}

func (i *info) Col() int {
	return i.col
}

func (i *info) Orientation() orientation.Orientation {
	return i.orientation
}
