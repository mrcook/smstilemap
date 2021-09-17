package background

import (
	"image"

	"github.com/mrcook/smstilemap/tile"
)

// Background represents an image consisting of unique tiles.
type Background struct {
	rows, cols int
	imageInfo  imageInfo
	tiles      []tile.Tile
}

func FromImage(img image.Image) *Background {
	bg := Background{
		rows: img.Bounds().Dy() / tile.Size,
		cols: img.Bounds().Dx() / tile.Size,
		imageInfo: imageInfo{
			width:  img.Bounds().Dx(),
			height: img.Bounds().Dy(),
		},
	}
	tiles := imageToTiles(img)
	bg.generateUniqueTileList(tiles)

	return &bg
}

// ToNRGBA converts the tiles to an NRGBA image.
func (b Background) ToNRGBA() *image.NRGBA {
	rows := len(b.tiles) / b.cols
	if len(b.tiles)%b.cols > 0 {
		rows += 1 // make up missing row
	}
	rows *= tile.Size
	cols := b.imageInfo.width

	img := image.NewNRGBA(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: cols, Y: rows},
	})

	xOffset := 0
	yOffset := 0

	for _, uniq := range b.tiles {
		for y := 0; y < tile.Size; y++ {
			for x := 0; x < tile.Size; x++ {
				colour := uniq.OriginalColorAt(x, y)
				img.Set(xOffset+x, yOffset+y, colour)
			}
		}

		xOffset += tile.Size
		if xOffset >= cols {
			xOffset = 0
			yOffset += tile.Size
		}
	}

	return img
}

// processes the tile list, recording all unique tiles, and adding duplicate
// info if the tile is already present.
func (b *Background) generateUniqueTileList(tiles []imageTile) {
	for _, t := range tiles {
		b.addTile(t.row, t.col, t.image)
	}
}

func (b *Background) addTile(row, col int, img image.Image) {
	info := tile.Info{Row: row, Col: col, Orientation: tile.OrientationNormal}

	// iterate over existing tiles and add as duplicate if a match is found
	for _, t := range b.tiles {
		if orientation, dupe := t.IsDuplicate(img); dupe {
			info.Orientation = orientation
			t.AddDuplicate(info)
			return
		}
	}

	b.tiles = append(b.tiles, tile.NewTile(info, img))
}
