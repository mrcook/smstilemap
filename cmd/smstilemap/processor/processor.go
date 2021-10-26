package processor

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path"
	"strings"

	"github.com/mrcook/smstilemap/assembly"
	"github.com/mrcook/smstilemap/cmd/smstilemap/processor/internal/tiler"
	"github.com/mrcook/smstilemap/sms"
)

type Processor struct {
	pngInputFilename string
	outputDirectory  string
	baseFilename     string

	image image.Image
	sega  sms.SMS
}

func New(srcFilename, outputDir string) *Processor {
	return &Processor{
		pngInputFilename: srcFilename,
		outputDirectory:  outputDirectory(outputDir, srcFilename),
		baseFilename:     baseFilename(srcFilename),
	}
}

func (p *Processor) CreateOutputDirectory() error {
	if err := os.MkdirAll(p.outputDirectory, 0755); err != nil {
		return fmt.Errorf("error creating output directory: %w", err)
	}
	return nil
}

func (p *Processor) PngToSMS() error {
	if err := p.readPNG(p.pngInputFilename); err != nil {
		return fmt.Errorf("PNG input file error: %w", err)
	}
	if err := p.imageToSMS(); err != nil {
		return fmt.Errorf("PNG to SMS data error: %w", err)
	}
	return nil
}

func (p *Processor) ToAssembly() error {
	var sb strings.Builder

	sb.WriteString(assembly.Tilemap(p.sega.TilemapData()).String())
	sb.WriteString("\n")
	sb.WriteString(assembly.Palettes(p.sega.PaletteData()).String())
	sb.WriteString("\n")
	sb.WriteString(assembly.Tiles(p.sega.TileData()).String())

	f, err := os.Create(path.Join(p.outputDirectory, p.baseFilename+".asm"))
	if err != nil {
		return fmt.Errorf("error creating ASM file: %w", err)
	}
	defer f.Close()

	if _, err := f.WriteString(sb.String()); err != nil {
		return fmt.Errorf("error writing SMS assembly to file: %w", err)
	}
	return nil
}

// SaveTilesToImage converts the SMS tiles to an image
func (p *Processor) SaveTilesToImage() error {
	dstImage, err := p.smsTilesToImage()
	if err != nil {
		return err
	}
	return p.saveImageToFilename(dstImage, p.pngTilesFilename())
}

// SaveTilemapToImage converts the SMS tilemap data back to a normal image
func (p *Processor) SaveTilemapToImage() error {
	dstImage, err := p.smsToImage()
	if err != nil {
		return err
	}
	return p.saveImageToFilename(dstImage, p.pngFilename())
}

func (p *Processor) readPNG(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	p.image, err = png.Decode(f)
	if err != nil {
		return err
	}

	return nil
}

// convert the PNG image to an SMS representation
func (p *Processor) imageToSMS() error {
	// validate image is suitable for conversion to the SMS
	if p.image == nil {
		return fmt.Errorf("source image is nil")
	} else if p.image.Bounds().Dx() > sms.ScreenWidth || p.image.Bounds().Dy() > sms.ScreenHeight {
		return fmt.Errorf("image size too big for SMS screen (%d x %d)", sms.ScreenWidth, sms.ScreenHeight)
	}
	tiled := tiler.FromImage(p.image, 8)

	// check there are too many colours for the SMS
	if tiled.ColourCount() > sms.MaxColourCount {
		return fmt.Errorf("too many unique colours for SMS (max: %d)", sms.MaxColourCount)
	}

	// add the image tiles to the SMS
	for i := 0; i < tiled.TileCount(); i++ {
		tile, _ := tiled.GetTile(i)
		if err := p.convertAndAddTileToSms(tile); err != nil {
			return err
		}
	}
	return nil
}

func (p *Processor) convertAndAddTileToSms(tile *tiler.Tile) error {
	if err := p.addTileColoursToSmsPalette(tile); err != nil {
		return fmt.Errorf("error adding colours to SMS palette: %w", err)
	}

	smsTile, err := p.convertToSmsTile(tile)
	if err != nil {
		return fmt.Errorf("error converting image tile to SMS tile: %w", err)
	}

	tid, err := p.sega.AddTile(smsTile)
	if err != nil {
		return err
	}

	if err := p.addTileToTilemap(tile, tid); err != nil {
		return fmt.Errorf("error adding tile to SMS tilemap: %w", err)
	}

	return nil
}

// make sure all tile colours are added to the SMS palette
func (p *Processor) addTileColoursToSmsPalette(tile *tiler.Tile) error {
	for _, c := range tile.Palette() {
		r, g, b, _ := c.RGBA()
		colour := sms.FromNearestMatchRGB(uint8(r), uint8(g), uint8(b))
		if _, err := p.sega.AddPaletteColour(colour); err != nil {
			return err
		}
	}
	return nil
}

// convert to an SMS tile, matching colours to SMS palette colours
func (p *Processor) convertToSmsTile(tile *tiler.Tile) (*sms.Tile, error) {
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
			pid, err := p.sega.PaletteIdForColour(colour)
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
func (p *Processor) addTileToTilemap(tile *tiler.Tile, tileId uint16) error {
	word := sms.Word{TileNumber: tileId}

	// the tile
	word.SetFlippedStateFromOrientation(p.smsOrientation(tile.Orientation()))
	if err := p.sega.AddTilemapEntryAt(tile.Row(), tile.Col(), word); err != nil {
		return err
	}

	// and its duplicates
	for did := 0; did < tile.DuplicateCount(); did++ {
		inf, err := tile.GetDuplicateInfo(did)
		if err != nil {
			return err
		}

		word.SetFlippedStateFromOrientation(p.smsOrientation(inf.Orientation()))
		if err := p.sega.AddTilemapEntryAt(inf.Row(), inf.Col(), word); err != nil {
			return err
		}
	}
	return nil
}

// converts a tiler orientation to an SMS orientation.
func (p *Processor) smsOrientation(or tiler.Orientation) sms.Orientation {
	switch or {
	case tiler.OrientationFlippedV:
		return sms.OrientationFlippedV
	case tiler.OrientationFlippedH:
		return sms.OrientationFlippedH
	case tiler.OrientationFlippedVH:
		return sms.OrientationFlippedVH
	default:
		return sms.OrientationNormal
	}
}

// smsTilesToImage converts the SMS tile data to a new NRGBA image.
func (p *Processor) smsTilesToImage() (image.Image, error) {
	tileSize := 8
	height, width := p.tileSheetSizeInPixels(tileSize)

	img := image.NewNRGBA(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: width, Y: height},
	})

	errorMessage := "drawing tile to image"
	rowOffset := 0
	colOffset := 0

	for i := uint16(0); i < sms.MaxTileCount; i++ {
		tile, err := p.sega.TileAt(i)
		if err != nil {
			return nil, err
		} else if tile == nil {
			break
		}

		// draw the tile to the image
		for y := 0; y < tileSize; y++ {
			for x := 0; x < tileSize; x++ {
				paletteId, err := tile.PaletteIdAt(y, x)
				if err != nil {
					return nil, fmt.Errorf("%s: %w", errorMessage, err)
				}
				colour, err := p.sega.PaletteColour(paletteId)
				if err != nil {
					return nil, fmt.Errorf("%s: %w", errorMessage, err)
				}
				img.Set(colOffset+x, rowOffset+y, colour)
			}
		}

		// next column
		colOffset += tileSize

		// next row?
		if colOffset >= width {
			colOffset = 0
			rowOffset += tileSize
		}
	}

	return img, nil
}

func (p *Processor) tileSheetSizeInPixels(tileSize int) (int, int) {
	height := sms.MaxTileCount / p.sega.WidthInTiles()
	if sms.MaxTileCount%p.sega.WidthInTiles() > 0 {
		height += 1 // make up missing row
	}
	height *= tileSize // set height in pixels
	width := p.sega.WidthInPixels()

	return height, width
}

// smsToImage converts the SMS data to a new NRGBA image, with the tile layout
// as defined in the tilemap name table.
func (p *Processor) smsToImage() (image.Image, error) {
	img := image.NewNRGBA(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: p.sega.WidthInPixels(), Y: p.sega.HeightInPixels()},
	})

	for row := 0; row < p.sega.HeightInTiles(); row++ {
		for col := 0; col < p.sega.WidthInTiles(); col++ {
			if err := p.drawTilemapEntry(img, row, col); err != nil {
				return nil, err
			}
		}
	}

	return img, nil
}

// draws a tile to the image using the tilemap entry data
func (p *Processor) drawTilemapEntry(img *image.NRGBA, row, col int) error {
	tile, err := p.smsTileForTilemapEntryAt(row, col)
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
			colour, err := p.sega.PaletteColour(paletteId)
			if err != nil {
				return fmt.Errorf("%s: %w", errorMessage, err)
			}
			img.Set(pxOffsetX+x, pxOffsetY+y, colour)
		}
	}
	return nil
}

// returns the mapped tile for the given row/col.
func (p *Processor) smsTileForTilemapEntryAt(row, col int) (*sms.Tile, error) {
	processingErrorMessage := "converting tilemap tile to correctly flipped tile"

	mapEntry, err := p.sega.TilemapEntryAt(row, col)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", processingErrorMessage, err)
	}
	tile, err := p.sega.TileAt(mapEntry.TileNumber)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", processingErrorMessage, err)
	}

	// set the correct orientation based on tilemap entry.
	return tile.AsTilemap(mapEntry), nil
}

func (p *Processor) saveImageToFilename(i image.Image, filename string) error {
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

func (p *Processor) pngFilename() string {
	return path.Join(p.outputDirectory, p.baseFilename+"-generated.png")
}

func (p *Processor) pngTilesFilename() string {
	return path.Join(p.outputDirectory, p.baseFilename+"-tiles.png")
}

func outputDirectory(dir, inFilename string) string {
	if len(dir) > 0 {
		return dir
	}
	return path.Dir(inFilename)
}

func baseFilename(input string) string {
	file := path.Base(input)
	return strings.ReplaceAll(file, path.Ext(file), "")
}
