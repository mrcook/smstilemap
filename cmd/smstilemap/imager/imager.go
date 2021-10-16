package imager

import (
	"fmt"
	"image"

	"github.com/mrcook/smstilemap/cmd/smstilemap/imager/internal/tiler"
	"github.com/mrcook/smstilemap/sms"
)

const tileSize = 8

type Imager struct {
	sega  *sms.SMS
	tiled *tiler.Tiled
}

// FromImage converts the given image into SMS image data.
func (t *Imager) FromImage(img image.Image) error {
	tiled, err := t.imageToTiled(img, tileSize)
	if err != nil {
		return err
	}
	t.tiled = tiled

	sega, err := t.tiledToSMS(tiled)
	if err != nil {
		return err
	}
	t.sega = sega

	return nil
}

// SmsToImage converts the SMS data to a new NRGBA image, with the tile layout
// as defined in the tilemap name table.
func (t *Imager) SmsToImage() (image.Image, error) {
	if t.sega == nil {
		return nil, fmt.Errorf("no image data available to convert")
	}

	img := image.NewNRGBA(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: t.sega.WidthInPixels(), Y: t.sega.VisibleHeightInPixels()},
	})

	for row := 0; row < t.sega.VisibleHeightInTiles(); row++ {
		for col := 0; col < t.sega.WidthInTiles(); col++ {
			if err := t.drawTilemapEntryFor(row, col, t.sega, img); err != nil {
				return nil, err
			}
		}
	}

	return img, nil
}

func (t *Imager) imageToTiled(img image.Image, tileSize int) (*tiler.Tiled, error) {
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

func (t *Imager) tiledToSMS(tiled *tiler.Tiled) (*sms.SMS, error) {
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

		smsTile, err := t.convertToSmsTile(&sega, tile)
		if err != nil {
			return nil, err
		}
		tid, err := sega.AddTile(smsTile)
		if err != nil {
			return nil, err
		}

		if err := t.updateTilemap(tid, &sega, tile); err != nil {
			return nil, err
		}
	}

	return &sega, nil
}

// convert to an SMS tile, matching colours to SMS palette colours
func (t *Imager) convertToSmsTile(sega *sms.SMS, tile *tiler.Tile) (*sms.Tile, error) {
	smsTile := sms.Tile{}

	for row := 0; row < tile.Size(); row++ {
		for col := 0; col < tile.Size(); col++ {
			// get the colour for the current pixel
			c, err := tile.OrientationAt(row, col, tile.Orientation())
			if err != nil {
				return nil, err
			}
			r, g, b, _ := c.RGBA()
			colour := sms.FromNearestMatchRGB(uint8(r), uint8(g), uint8(b))

			// find the palette ID for the colour
			pid, err := sega.PaletteIdForColour(colour)
			if err != nil {
				return nil, err
			}

			// sets the pixel colour using the palette ID
			if err := smsTile.SetPaletteIdAt(row, col, pid); err != nil {
				return nil, err
			}
		}
	}

	return &smsTile, nil
}

// update tilemap with the tile+duplicate locations
func (t *Imager) updateTilemap(tid uint16, sega *sms.SMS, tile *tiler.Tile) error {
	word := sms.Word{TileNumber: tid}

	// the tile
	word.SetFlippedStateFromOrientation(tile.Orientation())
	if err := sega.AddTilemapEntryAt(tile.Row(), tile.Col(), word); err != nil {
		return err
	}

	// and its duplicates
	for did := 0; did < tile.DuplicateCount(); did++ {
		inf, err := tile.GetDuplicateInfo(did)
		if err != nil {
			return err
		}

		word.SetFlippedStateFromOrientation(inf.Orientation())
		if err := sega.AddTilemapEntryAt(inf.Row(), inf.Col(), word); err != nil {
			return err
		}
	}
	return nil
}

// draws a tile to the image using the tilemap entry data
func (t *Imager) drawTilemapEntryFor(row, col int, sega *sms.SMS, img *image.NRGBA) error {
	tile, err := t.smsTileForTilemapEntryAt(sega, row, col)
	if err != nil {
		return err
	}

	errorMessage := "draw tile error"
	pxOffsetY := row * tile.Size()
	pxOffsetX := col * tile.Size()

	for y := 0; y < tile.Size(); y++ {
		for x := 0; x < tile.Size(); x++ {
			paletteId, err := tile.PaletteIdAt(y, x)
			if err != nil {
				return fmt.Errorf("%s: %w", errorMessage, err)
			}
			colour, err := sega.PaletteColour(paletteId)
			if err != nil {
				return fmt.Errorf("%s: %w", errorMessage, err)
			}
			img.Set(pxOffsetX+x, pxOffsetY+y, colour)
		}
	}
	return nil
}

// returns the mapped tile for the given row/col.
func (t *Imager) smsTileForTilemapEntryAt(sega *sms.SMS, row, col int) (*sms.Tile, error) {
	processingErrorMessage := "tilemap to orientated tile error"

	mapEntry, err := sega.TilemapEntryAt(row, col)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", processingErrorMessage, err)
	}
	tile, err := sega.TileAt(mapEntry.TileNumber)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", processingErrorMessage, err)
	}

	// set the correct orientation based on tilemap entry.
	return tile.AsTilemap(mapEntry), nil
}
