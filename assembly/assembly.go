package assembly

import (
	"fmt"
	"strings"
)

func Tiles(data []uint8) *strings.Builder {
	var sb strings.Builder
	lines := tileToHexStrings(data[:])

	sb.WriteString("; Tile data (characters)\n")
	sb.WriteString("; An 8x8 pixel tile is represented by 4x8 bytes. Each horizontal byte specifies\n")
	sb.WriteString("; 2 pixels, with each nibble of the byte referencing a palette ID.\n")
	sb.WriteString("TileData:\n")
	tileNumber := 0
	for i, line := range lines {
		if i == 0 || i%2 == 0 {
			sb.WriteString(fmt.Sprintf("; tile %03d:\n", tileNumber))
			tileNumber++
		}
		sb.WriteString(fmt.Sprintf(".db %s\n", line))
	}
	sb.WriteString("TileDataEnd:\n")
	return &sb
}

func Tilemap(data []uint16) *strings.Builder {
	var sb strings.Builder
	lines := tilemapToBinaryStrings(data[:])

	sb.WriteString("; Tilemap data (the name table)\n")
	sb.WriteString("; A matrix of 28 rows and 32 columns consisting of 16-bit [WORD] values:\n")
	sb.WriteString(";   Bit  |15 14 13|    12    |    11     |      10       |        9        | 8 7 6 5 4 3 2 1 0\n")
	sb.WriteString(";   Data | Unused | Priority | Palette # | Vertical flip | Horizontal flip |    Tile number\n")
	sb.WriteString("Tilemap:\n")
	row := 0
	for i, line := range lines {
		if i == 0 || i%8 == 0 {
			if row == 24 {
				break // don't show the unused last 4 rows of the 28 row tilemap
			}
			sb.WriteString(fmt.Sprintf("; row %02d\n", row))
			row++
		}
		sb.WriteString(fmt.Sprintf(".dw %s\n", line))
	}
	sb.WriteString("TilemapEnd:\n")
	return &sb
}

func Palettes(data [32]uint8) *strings.Builder {
	var sb strings.Builder
	lines := paletteToBinaryStrings(data[:])

	sb.WriteString("; Palette data; two 16 colour palettes\n")
	sb.WriteString(";   Bit:   7 6  |  5 4 |  3 2  | 1 0\n")
	sb.WriteString(";     %: Unused | Blue | Green | Red\n")
	sb.WriteString("PaletteData:\n")
	sb.WriteString("; palette 1\n")
	sb.WriteString(fmt.Sprintf(".db %s\n", lines[0]))
	sb.WriteString(fmt.Sprintf(".db %s\n", lines[1]))
	sb.WriteString("; palette 2\n")
	sb.WriteString(fmt.Sprintf(".db %s\n", lines[2]))
	sb.WriteString(fmt.Sprintf(".db %s\n", lines[3]))
	sb.WriteString("PaletteDataEnd:\n")
	return &sb
}

func tileToHexStrings(data []uint8) []string {
	bytesPerLine := 16
	var lines []string
	var pixels []string
	var pid = 0
	for _, b := range data {
		pixels = append(pixels, fmt.Sprintf("$%02X", b))
		if pid == bytesPerLine-1 {
			lines = append(lines, strings.Join(pixels[:], ", "))
			pid = 0
			pixels = []string{}
		} else {
			pid++
		}
	}
	return lines
}

func paletteToBinaryStrings(data []uint8) []string {
	bytesPerLine := 8
	var lines []string
	var pixels []string
	var pid = 0
	for _, b := range data {
		pixels = append(pixels, fmt.Sprintf("%%%08b", b))
		if pid == bytesPerLine-1 {
			lines = append(lines, strings.Join(pixels[:], ", "))
			pid = 0
			pixels = []string{}
		} else {
			pid++
		}
	}
	return lines
}

func tilemapToBinaryStrings(data []uint16) []string {
	wordsPerLine := 4
	var lines []string
	var pixels []string
	var pid = 0
	for _, b := range data {
		pixels = append(pixels, fmt.Sprintf("%%%016b", b))
		if pid == wordsPerLine-1 {
			lines = append(lines, strings.Join(pixels[:], ", "))
			pid = 0
			pixels = []string{}
		} else {
			pid++
		}
	}
	return lines
}
