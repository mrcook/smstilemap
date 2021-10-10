package imager

import (
	"fmt"
	"image"

	"github.com/mrcook/smstilemap/cmd/smstilemap/imager/internal/tiler"
	"github.com/mrcook/smstilemap/sms"
	"github.com/mrcook/smstilemap/sms/orientation"
)

const tileSize = 8

type Imager struct {
	sega  *sms.SMS
	tiled *tiler.Tiled
}

// FromImage converts the given image into SMS image data.
func (t *Imager) FromImage(img image.Image) error {
	tiled, err := imageToTiled(img, tileSize)
	if err != nil {
		return err
	}
	t.tiled = tiled

	sega, err := tiledToSMS(tiled)
	if err != nil {
		return err
	}
	t.sega = sega

	return nil
}

// TilemapToImage converts the tiled data to a new NRGBA image, with all tiles mapped
// back to their original positions.
func (t *Imager) TilemapToImage() (image.Image, error) {
	return convertScreenToImage(t.tiled)
}

func imageToTiled(img image.Image, tileSize int) (*tiler.Tiled, error) {
	// validate image is suitable for the SMS
	if img == nil {
		return nil, fmt.Errorf("source image is nil")
	} else if img.Bounds().Dx() > sms.MaxScreenWidth || img.Bounds().Dy() > sms.MaxScreenHeight {
		return nil, fmt.Errorf("image size too big for SMS screen (%d x %d)", sms.MaxScreenWidth, sms.MaxScreenHeight)
	}

	// convert incoming image to a tiled representation
	tiled := tiler.FromImage(img, tileSize)

	// now make sure there are not too many colours for the SMS
	if tiled.ColourCount() > sms.MaxColourCount {
		return nil, fmt.Errorf("too many unique colours for SMS (max: %d)", sms.MaxColourCount)
	}

	return tiled, nil
}

func tiledToSMS(tiled *tiler.Tiled) (*sms.SMS, error) {
	sega := sms.SMS{}

	for i := 0; i < tiled.TileCount(); i++ {
		tile, _ := tiled.GetTile(i)

		// add all colours from tile to SMS palette (if not already present)
		for _, c := range tile.Palette() {
			r, g, b, _ := c.RGBA()
			colour := sms.FromNearestMatchRGB(uint8(r), uint8(g), uint8(b))
			if _, err := sega.AddPaletteColour(colour); err != nil {
				return nil, err
			}
		}

		// convert to an SMS tile, matching colours to SMS palette colours
		// add SMS tile to tiles list
		// update tilemap with the tile location
	}

	return &sega, nil
}

func convertScreenToImage(bg *tiler.Tiled) (image.Image, error) {
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
		if err := drawTileAt(bgTile, img, y, x, orientation.Normal); err != nil {
			return nil, err
		}

		for did := 0; did < bgTile.DuplicateCount(); did++ {
			info, _ := bgTile.GetDuplicateInfo(did)
			y = info.Row() * tileSize
			x = info.Col() * tileSize
			if err := drawTileAt(bgTile, img, y, x, info.Orientation()); err != nil {
				return nil, err
			}
		}
	}

	return img, nil
}

func drawTileAt(tile *tiler.Tile, img *image.NRGBA, pxOffsetY, pxOffsetX int, orientation orientation.Orientation) error {
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
