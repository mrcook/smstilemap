package sms_test

import (
	"testing"

	"github.com/mrcook/smstilemap/sms"
)

func TestTile_Size(t *testing.T) {
	tile := sms.Tile{}
	if tile.Size() != 8 {
		t.Errorf("expected a tile size of 8, got %d", tile.Size())
	}
}

func TestTile_PaletteIdAt(t *testing.T) {
	tile := sms.Tile{}

	colourId := sms.PaletteId(2)
	_ = tile.SetPaletteIdAt(7, 7, colourId)

	t.Run("retrieving the pixel", func(t *testing.T) {
		pid, _ := tile.PaletteIdAt(7, 7)
		if pid != colourId {
			t.Fatalf("expected colour ID of %d, got %d", colourId, pid)
		}
	})

	t.Run("when row is bigger than tile size", func(t *testing.T) {
		_, err := tile.PaletteIdAt(8, 7)
		if err == nil {
			t.Fatal("expected an error")
		}
		if err.Error() != "tile indexing out of bounds, requested (8,7), tile size is 8" {
			t.Errorf("expected correct error message, got '%s'", err.Error())
		}
	})

	t.Run("when col is bigger than tile size", func(t *testing.T) {
		_, err := tile.PaletteIdAt(7, 8)
		if err == nil {
			t.Fatal("expected an error")
		}
		if err.Error() != "tile indexing out of bounds, requested (7,8), tile size is 8" {
			t.Errorf("expected correct error message, got '%s'", err.Error())
		}
	})
}

func TestTile_SetPaletteIdAt(t *testing.T) {
	tile := sms.Tile{}
	colourId := sms.PaletteId(31)

	t.Run("setting the pixel", func(t *testing.T) {
		_ = tile.SetPaletteIdAt(0, 0, colourId)
		pid, _ := tile.PaletteIdAt(0, 0)
		if pid != colourId {
			t.Fatalf("expected colour ID of %d, got %d", colourId, pid)
		}
	})

	t.Run("when row is bigger than tile size", func(t *testing.T) {
		err := tile.SetPaletteIdAt(8, 0, colourId)
		if err == nil {
			t.Fatal("expected an error")
		}
		if err.Error() != "tile indexing out of bounds, requested (8,0), tile size is 8" {
			t.Errorf("expected correct error message, got '%s'", err.Error())
		}
	})

	t.Run("when col is bigger than tile size", func(t *testing.T) {
		err := tile.SetPaletteIdAt(0, 8, colourId)
		if err == nil {
			t.Fatal("expected an error")
		}
		if err.Error() != "tile indexing out of bounds, requested (0,8), tile size is 8" {
			t.Errorf("expected correct error message, got '%s'", err.Error())
		}
	})
}

func TestTile_AsTilemap(t *testing.T) {
	tile := sms.Tile{}
	counter := sms.PaletteId(0)
	for row := 0; row < tile.Size(); row++ {
		for col := 0; col < tile.Size(); col++ {
			_ = tile.SetPaletteIdAt(row, col, counter)
			counter++
		}
	}

	t.Run("when flipped vertically", func(t *testing.T) {
		word := sms.Word{VerticalFlip: true}
		copied := tile.AsTilemap(&word)

		for col := 0; col < 8; col++ {
			fid, _ := copied.PaletteIdAt(0, col)
			tid, _ := tile.PaletteIdAt(7, col)
			if fid != tid {
				t.Errorf("vertical flip error, expected %dx%d to be %d, got %d", 0, col, tid, fid)
			}
			fid, _ = copied.PaletteIdAt(3, col)
			tid, _ = tile.PaletteIdAt(4, col)
			if fid != tid {
				t.Errorf("vertical flip error, expected pixel %dx%d to be %d, got %d", 3, col, tid, fid)
			}
			fid, _ = copied.PaletteIdAt(7, col)
			tid, _ = tile.PaletteIdAt(0, col)
			if fid != tid {
				t.Errorf("vertical flip error, expected pixel %dx%d to be %d, got %d", 7, col, tid, fid)
			}
		}
	})

	t.Run("when flipped horizontally", func(t *testing.T) {
		word := sms.Word{HorizontalFlip: true}
		copied := tile.AsTilemap(&word)

		for row := 0; row < 8; row++ {
			fid, _ := copied.PaletteIdAt(row, 0)
			tid, _ := tile.PaletteIdAt(row, 7)
			if fid != tid {
				t.Errorf("horizontal flip error, expected %dx%d to be %d, got %d", row, 0, tid, fid)
			}
			fid, _ = copied.PaletteIdAt(row, 3)
			tid, _ = tile.PaletteIdAt(row, 4)
			if fid != tid {
				t.Errorf("horizontal flip error, expected pixel %dx%d to be %d, got %d", row, 3, tid, fid)
			}
			fid, _ = copied.PaletteIdAt(row, 7)
			tid, _ = tile.PaletteIdAt(row, 0)
			if fid != tid {
				t.Errorf("horizontal flip error, expected pixel %dx%d to be %d, got %d", row, 7, tid, fid)
			}
		}
	})
}

func TestTile_Bytes(t *testing.T) {
	tile := sms.Tile{}
	_ = tile.SetPaletteIdAt(0, 0, sms.PaletteId(0b00000001))
	_ = tile.SetPaletteIdAt(0, 1, sms.PaletteId(0b00000010))
	_ = tile.SetPaletteIdAt(0, 2, sms.PaletteId(0b00000100))
	_ = tile.SetPaletteIdAt(0, 3, sms.PaletteId(0b00001000))
	_ = tile.SetPaletteIdAt(0, 4, sms.PaletteId(0b00001000))
	_ = tile.SetPaletteIdAt(0, 5, sms.PaletteId(0b00000100))
	_ = tile.SetPaletteIdAt(0, 6, sms.PaletteId(0b00000010))
	_ = tile.SetPaletteIdAt(0, 7, sms.PaletteId(0b00000001))
	data := tile.Bytes()

	t.Run("convert tile to planar data", func(t *testing.T) {
		if data[0] != 0b10000001 {
			t.Errorf("expected plane 0 to contain correct nibble values, got %08b", data[0])
		}
		if data[1] != 0b01000010 {
			t.Errorf("expected plane 1 to contain correct nibble values, got %08b", data[17])
		}
		if data[2] != 0b00100100 {
			t.Errorf("expected plane 2 to contain correct nibble values, got %08b", data[31])
		}
		if data[3] != 0b00011000 {
			t.Errorf("expected plane 3 to contain correct nibble values, got %08b", data[31])
		}
	})

	t.Run("when no pixel was set", func(t *testing.T) {
		if data[4] != 0b00000000 {
			t.Errorf("expected a default value of 0, got %08b", data[1])
		}
	})
}
