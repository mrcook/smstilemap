package tiler

type info struct {
	col, row    int // of the tile as used in the tilemap
	orientation Orientation
}

func (i *info) Row() int {
	return i.row
}

func (i *info) Col() int {
	return i.col
}

func (i *info) Orientation() Orientation {
	return i.orientation
}
