package sms_test

import (
	"testing"

	"github.com/mrcook/smstilemap/sms"
)

func TestSMS_ScreenDimensions(t *testing.T) {
	sega := sms.SMS{}

	if sega.WidthInPixels() != 256 {
		t.Errorf("expected screen width to be 256px, got %dpx", sega.WidthInPixels())
	}
	if sega.WidthInTiles() != 32 {
		t.Errorf("expected screen tile width to be 32, got %d", sega.WidthInTiles())
	}
	if sega.HeightInPixels() != 224 {
		t.Errorf("expected screen height to be 224px, got %dpx", sega.HeightInPixels())
	}
	if sega.HeightInTiles() != 28 {
		t.Errorf("expected screen tile height to be 28, got %d", sega.HeightInTiles())
	}
	if sega.VisibleHeightInPixels() != 192 {
		t.Errorf("expected visible screen height to be 192px, got %dpx", sega.VisibleHeightInPixels())
	}
	if sega.VisibleHeightInTiles() != 24 {
		t.Errorf("expected visible screen tile height to be 24, got %d", sega.VisibleHeightInTiles())
	}
}

func TestSMS_TileAt(t *testing.T) {
	sega := sms.SMS{}

	t.Run("get the saved tile correctly", func(t *testing.T) {
		// generate a tile that can be checked
		tile := sms.Tile{}
		_ = tile.SetPaletteIdAt(1, 1, 5)
		pos, _ := sega.AddTile(&tile)

		foundTile, err := sega.TileAt(pos)
		if err != nil {
			t.Fatalf("unexpected error: %q", err)
		}
		// get the pid and check it matches
		px, _ := foundTile.PaletteIdAt(1, 1)
		if px != 5 {
			t.Errorf("expect to find tile with correct palette ID, got %d", px)
		}
	})

	t.Run("when pid is greater than the length of the tile slice", func(t *testing.T) {
		_, err := sega.TileAt(sms.MaxTileCount)
		if err == nil {
			t.Fatal("expected an error")
		} else if err.Error() != "invalid tile ID" {
			t.Errorf("expected error message, got '%s'", err)
		}
	})
}

func TestSMS_AddTile(t *testing.T) {
	tile := sms.Tile{}
	sega := sms.SMS{}

	pos, err := sega.AddTile(&tile)
	if err != nil {
		t.Fatalf("unexpected error: %q", err)
	}
	if pos != 0 {
		t.Errorf("expected tile to be placed in first slot, tile id was %d", pos)
	}

	pos, err = sega.AddTile(&tile)
	if err != nil {
		t.Fatalf("unexpected error: %q", err)
	}
	if pos != 1 {
		t.Errorf("expected next tile to be placed in second slot, tile id was %d", pos)
	}
}

func TestSMS_TilemapEntryAt(t *testing.T) {
	vdp := sms.SMS{}
	word := sms.Word{TileNumber: 56}
	_ = vdp.AddTilemapEntryAt(20, 19, word)

	t.Run("returns an entry from the tilemap", func(t *testing.T) {
		got, err := vdp.TilemapEntryAt(20, 19)
		if err != nil {
			t.Fatalf("unexpected error: %q", err)
		}
		if got.TileNumber != 56 {
			t.Errorf("expected entry to have been set correctly, tile id was %d", got.TileNumber)
		}
	})

	t.Run("when tilemap is given bad inputs", func(t *testing.T) {
		_, err := vdp.TilemapEntryAt(28, 32)
		if err == nil {
			t.Fatal("expected an error")
		}
	})
}

func TestSMS_AddTilemapEntryAt(t *testing.T) {
	vdp := sms.SMS{}
	word := sms.Word{TileNumber: 199}

	t.Run("successfully adds an entry to the tilemap", func(t *testing.T) {
		err := vdp.AddTilemapEntryAt(27, 31, word)
		if err != nil {
			t.Fatalf("unexpected error: %q", err)
		}

		got, err := vdp.TilemapEntryAt(27, 31)
		if err != nil {
			t.Fatalf("unexpected error: %q", err)
		}
		if got.TileNumber != 199 {
			t.Errorf("expected entry to have been set correctly, tile id was %d", got.TileNumber)
		}
	})

	t.Run("when tilemap is given bad inputs", func(t *testing.T) {
		err := vdp.AddTilemapEntryAt(28, 32, word)
		if err == nil {
			t.Fatal("expected an error")
		}
	})
}

func TestSMS_PaletteColour(t *testing.T) {
	sega := sms.SMS{}
	_, _ = sega.AddPaletteColour(sms.Colour(0b00000011))
	colour := sms.Colour(0b00111111)
	pid, _ := sega.AddPaletteColour(colour)

	t.Run("return correct colour for given ID", func(t *testing.T) {
		gotColour, _ := sega.PaletteColour(pid)
		if gotColour != colour {
			t.Errorf("expected correct colour, got %08b", pid)
		}
	})

	t.Run("when no colour has been set for the requested palette ID", func(t *testing.T) {
		_, err := sega.PaletteColour(pid + 1)
		if err == nil {
			t.Errorf("expected an error")
		}
	})
}

func TestSMS_PaletteIdForColour(t *testing.T) {
	t.Run("return position of matching colour", func(t *testing.T) {
		sega := sms.SMS{}
		_, _ = sega.AddPaletteColour(sms.Colour(0b00000011))
		colour := sms.Colour(0b00111111)
		_, _ = sega.AddPaletteColour(colour)

		pos, _ := sega.PaletteIdForColour(colour)
		if pos != 1 {
			t.Errorf("expected existing position, got %d", pos)
		}
	})

	t.Run("when no colour match found", func(t *testing.T) {
		sega := sms.SMS{}
		_, err := sega.PaletteIdForColour(sms.Colour(0b00000011))
		if err == nil {
			t.Fatal("expected an error")
		} else if err.Error() != "colour not found" {
			t.Errorf("unexpected error message, got '%s'", err)
		}
	})
}

func TestSMS_AddPaletteColour(t *testing.T) {
	t.Run("return position", func(t *testing.T) {
		sega := sms.SMS{}
		_, _ = sega.AddPaletteColour(sms.Colour(0b00000011))

		pos, _ := sega.AddPaletteColour(sms.Colour(0b00010101))
		if pos != 1 {
			t.Errorf("expected colour to be added at first slot, got %d", pos)
		}
	})

	t.Run("when palette is full, return an error", func(t *testing.T) {
		sega := sms.SMS{}
		for i := 0; i < 32; i++ {
			_, _ = sega.AddPaletteColour(sms.Colour(i))
		}

		_, err := sega.AddPaletteColour(sms.Colour(0b00111111))
		if err == nil {
			t.Fatalf("expected error")
		} else if err.Error() != "palette full" {
			t.Errorf("expect a valid error message, got '%s'", err)
		}
	})
}
