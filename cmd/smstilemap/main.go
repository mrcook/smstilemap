package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"

	background "github.com/mrcook/smstilemap/image"
)

func main() {
	srcFilename := parseCliForFilename()
	dstFilename := srcFilename + "-new.png" // simple but works

	srcImage, err := openImage(srcFilename)
	if err != nil {
		log.Fatal(err)
	}

	// process it
	bg := background.FromImage(srcImage)
	// dstImage, err := toImage(bg) // only unique tiles
	dstImage, err := toTileMappedImage(bg)
	if err != nil {
		log.Fatal(err)
	}

	// save to new png
	if err := saveImage(dstImage, dstFilename); err != nil {
		log.Fatal(err)
	}
}

// ToImage converts the tiles to an NRGBA image.
func toImage(bg *background.Background) (image.Image, error) {
	rows := bg.TileCount() / bg.Info().Cols
	if bg.TileCount()%bg.Info().Cols > 0 {
		rows += 1 // make up missing row
	}
	rows *= background.Size
	cols := bg.Info().Width

	img := image.NewNRGBA(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: cols, Y: rows},
	})

	xOffset := 0
	yOffset := 0

	for i := 0; i < bg.TileCount(); i++ {
		bgTile, _ := bg.GetTile(i)
		if err := drawTileAt(bgTile, img, i, yOffset, xOffset, background.OrientationNormal); err != nil {
			return nil, err
		}
		xOffset += background.Size
		if xOffset >= cols {
			xOffset = 0
			yOffset += background.Size
		}
	}

	return img, nil
}

// Converts the tiles to a new NRGBA image, with all tiles mapped back to
// their original positions.
func toTileMappedImage(bg *background.Background) (image.Image, error) {
	img := image.NewNRGBA(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: bg.Info().Width, Y: bg.Info().Height},
	})

	for i := 0; i < bg.TileCount(); i++ {
		bgTile, _ := bg.GetTile(i)

		y := bgTile.RowInPixels()
		x := bgTile.ColInPixels()
		if err := drawTileAt(bgTile, img, i, y, x, background.OrientationNormal); err != nil {
			return nil, err
		}

		for did := 0; did < bgTile.DuplicateCount(); did++ {
			info, _ := bgTile.GetDuplicateInfo(did)
			y = info.Row() * background.Size
			x = info.Col() * background.Size
			if err := drawTileAt(bgTile, img, i, y, x, info.Orientation()); err != nil {
				return nil, err
			}
		}
	}

	return img, nil
}

func drawTileAt(t *background.Tile, img *image.NRGBA, tileIndex, pxOffsetY, pxOffsetX int, orientation background.Orientation) error {
	for y := 0; y < background.Size; y++ {
		for x := 0; x < background.Size; x++ {
			colour, err := t.OrientationAt(y, x, orientation)
			if err != nil {
				return fmt.Errorf("draw tile error: %w", err)
			}
			img.Set(pxOffsetX+x, pxOffsetY+y, colour)
		}
	}
	return nil
}

func openImage(filename string) (image.Image, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	decodedImage, err := png.Decode(f)
	if err != nil {
		return nil, err
	}

	return decodedImage, nil
}

func saveImage(i image.Image, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	err = png.Encode(f, i)
	if err != nil {
		return err
	}
	return nil
}

func parseCliForFilename() string {
	filename := flag.String("src", "", "Source PNG image filename")
	flag.Parse()
	if len(*filename) == 0 {
		log.Fatal("source filename is required")
	}
	return *filename
}
