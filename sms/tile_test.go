package sms_test

import (
	"testing"

	"github.com/mrcook/smstilemap/sms"
)

func TestTile_PixelAt(t *testing.T) {
	tile := sms.Tile{}

	colourId := sms.PaletteId(2)
	_ = tile.SetPixelAt(7, 7, colourId)

	t.Run("retrieving the pixel", func(t *testing.T) {
		pid, _ := tile.PixelAt(7, 7)
		if pid != colourId {
			t.Fatalf("expected colour ID of %d, got %d", colourId, pid)
		}
	})

	t.Run("when row is bigger than tile size", func(t *testing.T) {
		_, err := tile.PixelAt(8, 7)
		if err == nil {
			t.Fatal("expected an error")
		}
		if err.Error() != "tile indexing out of bounds, requested (8,7), tile size is 8" {
			t.Errorf("expected correct error message, got '%s'", err.Error())
		}
	})

	t.Run("when col is bigger than tile size", func(t *testing.T) {
		_, err := tile.PixelAt(7, 8)
		if err == nil {
			t.Fatal("expected an error")
		}
		if err.Error() != "tile indexing out of bounds, requested (7,8), tile size is 8" {
			t.Errorf("expected correct error message, got '%s'", err.Error())
		}
	})
}

func TestTile_SetPixelAt(t *testing.T) {
	tile := sms.Tile{}
	colourId := sms.PaletteId(31)

	t.Run("setting the pixel", func(t *testing.T) {
		_ = tile.SetPixelAt(0, 0, colourId)
		pid, _ := tile.PixelAt(0, 0)
		if pid != colourId {
			t.Fatalf("expected colour ID of %d, got %d", colourId, pid)
		}
	})

	t.Run("when row is bigger than tile size", func(t *testing.T) {
		err := tile.SetPixelAt(8, 0, colourId)
		if err == nil {
			t.Fatal("expected an error")
		}
		if err.Error() != "tile indexing out of bounds, requested (8,0), tile size is 8" {
			t.Errorf("expected correct error message, got '%s'", err.Error())
		}
	})

	t.Run("when col is bigger than tile size", func(t *testing.T) {
		err := tile.SetPixelAt(0, 8, colourId)
		if err == nil {
			t.Fatal("expected an error")
		}
		if err.Error() != "tile indexing out of bounds, requested (0,8), tile size is 8" {
			t.Errorf("expected correct error message, got '%s'", err.Error())
		}
	})
}
