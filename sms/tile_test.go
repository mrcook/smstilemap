package sms_test

import (
	"testing"

	"github.com/mrcook/smstilemap/sms"
)

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

func TestTile_Flipped(t *testing.T) {
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
		flipped := tile.Flipped(&word)

		for col := 0; col < 8; col++ {
			fid, _ := flipped.PaletteIdAt(0, col)
			tid, _ := tile.PaletteIdAt(7, col)
			if fid != tid {
				t.Errorf("vertical flip error, expected %dx%d to be %d, got %d", 0, col, tid, fid)
			}
			fid, _ = flipped.PaletteIdAt(3, col)
			tid, _ = tile.PaletteIdAt(4, col)
			if fid != tid {
				t.Errorf("vertical flip error, expected pixel %dx%d to be %d, got %d", 3, col, tid, fid)
			}
			fid, _ = flipped.PaletteIdAt(7, col)
			tid, _ = tile.PaletteIdAt(0, col)
			if fid != tid {
				t.Errorf("vertical flip error, expected pixel %dx%d to be %d, got %d", 7, col, tid, fid)
			}
		}
	})

	t.Run("when flipped horizontally", func(t *testing.T) {
		word := sms.Word{HorizontalFlip: true}
		flipped := tile.Flipped(&word)

		for row := 0; row < 8; row++ {
			fid, _ := flipped.PaletteIdAt(row, 0)
			tid, _ := tile.PaletteIdAt(row, 7)
			if fid != tid {
				t.Errorf("horizontal flip error, expected %dx%d to be %d, got %d", row, 0, tid, fid)
			}
			fid, _ = flipped.PaletteIdAt(row, 3)
			tid, _ = tile.PaletteIdAt(row, 4)
			if fid != tid {
				t.Errorf("horizontal flip error, expected pixel %dx%d to be %d, got %d", row, 3, tid, fid)
			}
			fid, _ = flipped.PaletteIdAt(row, 7)
			tid, _ = tile.PaletteIdAt(row, 0)
			if fid != tid {
				t.Errorf("horizontal flip error, expected pixel %dx%d to be %d, got %d", row, 7, tid, fid)
			}
		}
	})
}
