package main

import (
	"fmt"
	"image"

	"github.com/mrcook/smstilemap/cmd/smstilemap/internal/tiler"
	"github.com/mrcook/smstilemap/sms"
	"github.com/mrcook/smstilemap/sms/orientation"
)

const (
	tileSize = 8
)

type tiledImage struct {
	tiled *tiler.Tiled
}

// fromImage converts the given image into SMS image data.
func (t *tiledImage) fromImage(img image.Image) error {
	if err := t.processImage(img, tileSize); err != nil {
		return err
	}
	return nil
}

// tilemapToImage converts the tiled data to a new NRGBA image, with all tiles mapped
// back to their original positions.
func (t *tiledImage) tilemapToImage() (image.Image, error) {
	return t.convertScreenToImage()
}

func (t *tiledImage) processImage(img image.Image, tileSize int) error {
	// validate image is suitable for the SMS
	if img == nil {
		return fmt.Errorf("source image is nil")
	} else if img.Bounds().Dx() > sms.MaxScreenWidth || img.Bounds().Dy() > sms.MaxScreenHeight {
		return fmt.Errorf("image size too big for SMS screen (%d x %d)", sms.MaxScreenWidth, sms.MaxScreenHeight)
	}

	// convert incoming image to a tiled representation
	t.tiled = tiler.FromImage(img, tileSize)

	// now make sure there are not too many colours for the SMS
	if t.tiled.ColourCount() > sms.MaxColourCount {
		return fmt.Errorf("too many unique colours for SMS (max: %d)", sms.MaxColourCount)
	}

	return nil
}

func (t *tiledImage) convertScreenToImage() (image.Image, error) {
	bg := t.tiled

	if bg == nil {
		return nil, fmt.Errorf("no image data available to convert")
	}

	img := image.NewNRGBA(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: bg.Width(), Y: bg.Height()},
	})

	for i := 0; i < bg.TileCount(); i++ {
		bgTile, _ := bg.GetTile(i)

		y := bgTile.RowInPixels()
		x := bgTile.ColInPixels()
		if err := t.drawTileAt(bgTile, img, y, x, orientation.Normal); err != nil {
			return nil, err
		}

		for did := 0; did < bgTile.DuplicateCount(); did++ {
			info, _ := bgTile.GetDuplicateInfo(did)
			y = info.Row() * tileSize
			x = info.Col() * tileSize
			if err := t.drawTileAt(bgTile, img, y, x, info.Orientation()); err != nil {
				return nil, err
			}
		}
	}

	return img, nil
}

func (t *tiledImage) drawTileAt(tile *tiler.Tile, img *image.NRGBA, pxOffsetY, pxOffsetX int, orientation orientation.Orientation) error {
	for y := 0; y < tileSize; y++ {
		for x := 0; x < tileSize; x++ {
			colour, err := tile.OrientationAt(y, x, orientation)
			if err != nil {
				return fmt.Errorf("draw tile error: %w", err)
			}
			img.Set(pxOffsetX+x, pxOffsetY+y, colour)
		}
	}
	return nil
}
