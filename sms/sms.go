package sms

type VRAM struct {
	characters [448]Tile  // sprite/tile patterns
	nameTable  Tilemap    // screen map
	sat        [256]uint8 // Sprite Attribute Table (x,y coords, tile ID)...not needed in the conversion?
}

type CRAM struct {
	palette1 [16]Colour // background palette
	palette2 [16]Colour // sprite and background palette
}
