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

	for i := 0; i < len(b.tiles); i++ {
		b.drawTileAt(img, i, yOffset, xOffset, b.tiles[i].Info().Orientation)
		xOffset += tile.Size
		if xOffset >= cols {
			xOffset = 0
			yOffset += tile.Size
		}
	}

	return img
}

// ToTileMappedNRGBA converts the tiles to a new NRGBA image, with all tiles
// mapped to their correct position.
func (b Background) ToTileMappedNRGBA() *image.NRGBA {
	img := image.NewNRGBA(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: b.imageInfo.width, Y: b.imageInfo.height},
	})

	for i := 0; i < len(b.tiles); i++ {
		t := &b.tiles[i]

		row := t.Info().Row * tile.Size
		col := t.Info().Col * tile.Size
		b.drawTileAt(img, i, row, col, t.Info().Orientation)

		for _, info := range t.Duplicates() {
			row = info.Row * tile.Size
			col = info.Col * tile.Size
			b.drawTileAt(img, i, row, col, info.Orientation)
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
	for i := 0; i < len(b.tiles); i++ {
		if orientation, dupe := b.tiles[i].IsDuplicate(img); dupe {
			info.Orientation = orientation
			b.tiles[i].AddDuplicate(info)
			return
		}
	}

	b.tiles = append(b.tiles, tile.NewTile(info, img))
}

func (b Background) drawTileAt(img *image.NRGBA, tileIndex, pxOffsetY, pxOffsetX int, orientation tile.Orientation) {
	tt := b.tiles[tileIndex].ImageInOrientation(orientation)

	for y := 0; y < tile.Size; y++ {
		for x := 0; x < tile.Size; x++ {
			colour := tt.At(x, y)
			img.Set(pxOffsetX+x, pxOffsetY+y, colour)
		}
	}
}
