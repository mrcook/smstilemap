package processor

import (
	"fmt"
	"image"

	"github.com/mrcook/smstilemap/cmd/smstilemap/processor/internal/tiler"
	"github.com/mrcook/smstilemap/sms"
)

// ImageToSms converts an image into SMS tile data.
func ImageToSms(img image.Image) (*sms.SMS, error) {
	// validate image is suitable for conversion to the SMS
	if img == nil {
		return nil, fmt.Errorf("source image is nil")
	} else if img.Bounds().Dx() > sms.MaxScreenWidth || img.Bounds().Dy() > sms.MaxScreenHeight {
		return nil, fmt.Errorf("image size too big for SMS screen (%d x %d)", sms.MaxScreenWidth, sms.MaxScreenHeight)
	}
	tiled := tiler.FromImage(img, 8)

	// check there are too many colours for the SMS
	if tiled.ColourCount() > sms.MaxColourCount {
		return nil, fmt.Errorf("too many unique colours for SMS (max: %d)", sms.MaxColourCount)
	}

	// add the image tiles to the SMS
	sega := sms.SMS{}
	for i := 0; i < tiled.TileCount(); i++ {
		tile, _ := tiled.GetTile(i)
		if err := convertAndAddTileToSms(&sega, tile); err != nil {
			return nil, err
		}
	}
	return &sega, nil
}

// SmsToImage converts the SMS data to a new NRGBA image, with the tile layout
// as defined in the tilemap name table.
func SmsToImage(sega *sms.SMS) (image.Image, error) {
	if sega == nil {
		return nil, fmt.Errorf("no SMS data available to convert")
	}

	img := image.NewNRGBA(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: sega.WidthInPixels(), Y: sega.VisibleHeightInPixels()},
	})

	for row := 0; row < sega.VisibleHeightInTiles(); row++ {
		for col := 0; col < sega.WidthInTiles(); col++ {
			if err := drawTilemapEntry(sega, img, row, col); err != nil {
				return nil, err
			}
		}
	}

	return img, nil
}

func convertAndAddTileToSms(sega *sms.SMS, tile *tiler.Tile) error {
	if err := addTileColoursToSmsPalette(sega, tile); err != nil {
		return fmt.Errorf("error adding colours to SMS palette: %w", err)
	}

	smsTile, err := convertToSmsTile(sega, tile)
	if err != nil {
		return fmt.Errorf("error converting image tile to SMS tile: %w", err)
	}

	tid, err := sega.AddTile(smsTile)
	if err != nil {
		return err
	}

	if err := addTileToTilemap(sega, tile, tid); err != nil {
		return fmt.Errorf("error adding tile to SMS tilemap: %w", err)
	}

	return nil
}

// make sure all tile colours are added to the SMS palette
func addTileColoursToSmsPalette(sega *sms.SMS, tile *tiler.Tile) error {
	for _, c := range tile.Palette() {
		r, g, b, _ := c.RGBA()
		colour := sms.FromNearestMatchRGB(uint8(r), uint8(g), uint8(b))
		if _, err := sega.AddPaletteColour(colour); err != nil {
			return err
		}
	}
	return nil
}

// convert to an SMS tile, matching colours to SMS palette colours
func convertToSmsTile(sega *sms.SMS, tile *tiler.Tile) (*sms.Tile, error) {
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
func addTileToTilemap(sega *sms.SMS, tile *tiler.Tile, tileId uint16) error {
	word := sms.Word{TileNumber: tileId}

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
func drawTilemapEntry(sega *sms.SMS, img *image.NRGBA, row, col int) error {
	tile, err := smsTileForTilemapEntryAt(sega, row, col)
	if err != nil {
		return err
	}

	errorMessage := "drawing tile to image"
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
func smsTileForTilemapEntryAt(sega *sms.SMS, row, col int) (*sms.Tile, error) {
	processingErrorMessage := "converting tilemap tile to correctly flipped tile"

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
