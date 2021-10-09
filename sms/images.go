package sms

// This file contains a set of helper functions for converting to and from
// standard image.Image data.
// They're kept here to keep the main SMS file less cluttered.

import (
	"fmt"
	"image"

	"github.com/mrcook/smstilemap/sms/internal/tiler"
	"github.com/mrcook/smstilemap/sms/orientation"
)

// FromImage converts the given image into SMS image data.
func (s *SMS) FromImage(img image.Image) error {
	return s.readImageOntoSMS(img, tileSize)
}

// TilemapToImage converts the tiled data to a new NRGBA image, with all tiles mapped
// back to their original positions.
func (s *SMS) TilemapToImage() (image.Image, error) {
	return s.convertScreenToImage()
}

func (s *SMS) readImageOntoSMS(img image.Image, tileSize int) error {
	// validate image is suitable for the SMS
	if img == nil {
		return fmt.Errorf("source image is nil")
	} else if img.Bounds().Dx() > maxScreenWidth || img.Bounds().Dy() > maxScreenHeight {
		return fmt.Errorf("image size too big for SMS screen (%d x %d)", maxScreenWidth, maxScreenHeight)
	}

	// convert incoming image to a tiled representation
	s.tiledImg = tiler.FromImage(img, tileSize)

	// now make sure there are not too many colours for the SMS
	if s.tiledImg.ColourCount() > maxColourCount {
		return fmt.Errorf("too many unique colours for SMS (max: %d)", maxColourCount)
	}

	// convert all background tiles to planar data and add to tilemap
	for i := 0; i < s.tiledImg.TileCount(); i++ {
		if tile, err := s.tiledImg.GetTile(i); err != nil {
			return err
		} else {
			tileId, err := s.AddTile(&Tile{}) // TODO: convert to Tile
			if err != nil {
				return fmt.Errorf("out of tile space")
			}
			s.addTilemapEntries(tileId, tile)
		}
	}

	return nil
}

func (s *SMS) addTilemapEntries(tileID int, tile *tiler.Tile) {
	// add the normal orientation tile
	s.AddTilemapEntry(tileID, tile.Row(), tile.Col(), tile.Orientation())

	// add any duplicate (flipped) tiles
	for i := 0; i < tile.DuplicateCount(); i++ {
		inf, err := tile.GetDuplicateInfo(i)
		if err != nil {
			break // TODO: break?
		}
		s.AddTilemapEntry(tileID, inf.Row(), inf.Col(), inf.Orientation())
	}
}

func (s *SMS) convertScreenToImage() (image.Image, error) {
	bg := s.tiledImg

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
		if err := s.drawTileAt(bgTile, img, y, x, orientation.Normal); err != nil {
			return nil, err
		}

		for did := 0; did < bgTile.DuplicateCount(); did++ {
			info, _ := bgTile.GetDuplicateInfo(did)
			y = info.Row() * tileSize
			x = info.Col() * tileSize
			if err := s.drawTileAt(bgTile, img, y, x, info.Orientation()); err != nil {
				return nil, err
			}
		}
	}

	return img, nil
}

func (s SMS) drawTileAt(t *tiler.Tile, img *image.NRGBA, pxOffsetY, pxOffsetX int, orientation orientation.Orientation) error {
	for y := 0; y < tileSize; y++ {
		for x := 0; x < tileSize; x++ {
			colour, err := t.OrientationAt(y, x, orientation)
			if err != nil {
				return fmt.Errorf("draw tile error: %w", err)
			}
			img.Set(pxOffsetX+x, pxOffsetY+y, colour)
		}
	}
	return nil
}
