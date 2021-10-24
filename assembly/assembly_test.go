package assembly_test

import (
	"strings"
	"testing"

	"github.com/mrcook/smstilemap/assembly"
)

func TestAssembly_TileData(t *testing.T) {
	var testTile = [8][4]uint8{
		{0b00010001, 0b00010001, 0b00010001, 0b00010001},
		{0b00100010, 0b00100010, 0b00100010, 0b00100010},
		{0b00110011, 0b00110011, 0b00110011, 0b00110011},
		{0b01000100, 0b01000100, 0b01000100, 0b01000100},
		{0b01010101, 0b01010101, 0b01010101, 0b01010101},
		{0b01100110, 0b01100110, 0b01100110, 0b01100110},
		{0b01110111, 0b01110111, 0b01110111, 0b01110111},
		{0b10001000, 0b10001000, 0b10001000, 0b10001000},
	}
	var data []uint8
	for i := 0; i < 2; i++ {
		for _, cols := range testTile {
			for _, b := range cols {
				data = append(data, b)
			}
		}
	}

	got := assembly.Tiles(data).String()
	want := `TileData:
; tile 000:
.db $11, $11, $11, $11, $22, $22, $22, $22, $33, $33, $33, $33, $44, $44, $44, $44
.db $55, $55, $55, $55, $66, $66, $66, $66, $77, $77, $77, $77, $88, $88, $88, $88
; tile 001:
.db $11, $11, $11, $11, $22, $22, $22, $22, $33, $33, $33, $33, $44, $44, $44, $44
.db $55, $55, $55, $55, $66, $66, $66, $66, $77, $77, $77, $77, $88, $88, $88, $88
TileDataEnd:
`
	lines := strings.Split(got, "\n")
	got = strings.Join(lines[3:], "\n")
	if got != want {
		t.Errorf("unexpected output, got:\n%s", got)
	}
}

func TestAssembly_TilemapData(t *testing.T) {
	tilemapData := make([]uint16, 896)
	tilemapData[0] = 0b0000000000000001
	tilemapData[1] = 0b0000001000000001
	tilemapData[8] = 0b0000000000000110
	tilemapData[9] = 0b0000010000000110

	fullTilemap := assembly.Tilemap(tilemapData)
	want := `Tilemap:
; row 00
.dw %0000000000000001, %0000001000000001, %0000000000000000, %0000000000000000
.dw %0000000000000000, %0000000000000000, %0000000000000000, %0000000000000000
.dw %0000000000000110, %0000010000000110, %0000000000000000, %0000000000000000
.dw %0000000000000000, %0000000000000000, %0000000000000000, %0000000000000000`

	lines := strings.Split(fullTilemap.String(), "\n")
	got := strings.Join(lines[4:10], "\n")

	if got != want {
		t.Errorf("unexpected output, got:\n%s", got)
	}
}

func TestAssembly_PaletteData(t *testing.T) {
	var paletteData [32]uint8
	paletteData[0] = 0b00000011
	paletteData[15] = 0b00001100
	paletteData[16] = 0b00111111
	paletteData[31] = 0b00110000

	got := assembly.Palettes(paletteData).String()
	want := `PaletteData:
; palette 1
.db %00000011, %00000000, %00000000, %00000000, %00000000, %00000000, %00000000, %00000000
.db %00000000, %00000000, %00000000, %00000000, %00000000, %00000000, %00000000, %00001100
; palette 2
.db %00111111, %00000000, %00000000, %00000000, %00000000, %00000000, %00000000, %00000000
.db %00000000, %00000000, %00000000, %00000000, %00000000, %00000000, %00000000, %00110000
PaletteDataEnd:
`
	lines := strings.Split(got, "\n")
	got = strings.Join(lines[3:], "\n")
	if got != want {
		t.Errorf("unexpected output, got:\n%s", got)
	}
}
