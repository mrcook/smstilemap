package sms

type Word struct {
	Priority       bool   // bit 12: tile is displayed in front of sprites when set
	PaletteSelect  bool   // bit 11: use tile palette or sprite palette (when set)
	VerticalFlip   bool   // bit 10: flip vertically
	HorizontalFlip bool   // bit 09: flip horizontally
	TileNumber     uint16 // tile definition number to use (0..511)
}

func (w Word) ToUint() uint16 {
	var value uint16

	if w.Priority {
		value |= 0b0001000000000000
	}
	if w.PaletteSelect {
		value |= 0b0000100000000000
	}
	if w.VerticalFlip {
		value |= 0b0000010000000000
	}
	if w.HorizontalFlip {
		value |= 0b0000001000000000
	}
	if w.TileNumber <= 511 {
		value |= w.TileNumber
	}

	return value
}