package sms_test

import (
	"testing"

	"github.com/mrcook/smstilemap/sms"
)

func TestTilemap_Width(t *testing.T) {
	tm := sms.Tilemap{}
	if tm.Width() != 32 {
		t.Errorf("expected tilemap width to be 32, got %d", tm.Width())
	}
}

func TestTilemap_Height(t *testing.T) {
	tm := sms.Tilemap{}
	if tm.Height() != 24 {
		t.Errorf("expected tilemap height to be 24, got %d", tm.Height())
	}
}

func TestTilemap_ExtendedHeight(t *testing.T) {
	tm := sms.Tilemap{}
	if tm.ExtendedHeight() != 28 {
		t.Errorf("expected tilemap extended height to be 28, got %d", tm.ExtendedHeight())
	}
}

func TestTilemap_Get(t *testing.T) {
	tm := sms.Tilemap{}
	word := sms.Word{TileNumber: 447}
	_ = tm.Set(23, 0, word)

	t.Run("with valid data", func(t *testing.T) {
		got, err := tm.Get(23, 0)
		if err != nil {
			t.Fatalf("unexpected error: %q", err)
		}
		if got.TileNumber != word.TileNumber {
			t.Errorf("expected to get correct data, got tile id %d", got.TileNumber)
		}
	})

	t.Run("with out of bounds row indexing", func(t *testing.T) {
		_, err := tm.Get(24, 31)
		if err == nil {
			t.Fatal("expected an error")
		}
		if err.Error() != "get tilemap out of bounds indexing, max is (23,31), requested (24,31)" {
			t.Errorf("unexpected error message, requested '%s'", err)
		}
	})

	t.Run("with out of bounds col indexing", func(t *testing.T) {
		_, err := tm.Get(23, 32)
		if err == nil {
			t.Fatal("expected an error")
		}
		if err.Error() != "get tilemap out of bounds indexing, max is (23,31), requested (23,32)" {
			t.Errorf("unexpected error message, requested '%s'", err)
		}
	})
}

func TestTilemap_Set(t *testing.T) {
	tm := sms.Tilemap{}
	word := sms.Word{TileNumber: 1}

	t.Run("with valid data", func(t *testing.T) {
		_ = tm.Set(0, 0, word)
		got, err := tm.Get(0, 0)
		if err != nil {
			t.Fatalf("unexpected error: %q", err)
		}
		if got.TileNumber != word.TileNumber {
			t.Errorf("expected correct data to have been set, got tile id %d", got.TileNumber)
		}
	})

	t.Run("with out of bounds row indexing", func(t *testing.T) {
		err := tm.Set(24, 31, word)
		if err == nil {
			t.Fatal("expected an error")
		}
		if err.Error() != "set tilemap out of bounds indexing, max is (23,31), requested (24,31)" {
			t.Errorf("unexpected error message, got '%s'", err)
		}
	})

	t.Run("with out of bounds col indexing", func(t *testing.T) {
		err := tm.Set(23, 32, word)
		if err == nil {
			t.Fatal("expected an error")
		}
		if err.Error() != "set tilemap out of bounds indexing, max is (23,31), requested (23,32)" {
			t.Errorf("unexpected error message, got '%s'", err)
		}
	})
}

func TestTilemap_Words(t *testing.T) {
	tm := sms.Tilemap{}
	_ = tm.Set(0, 0, sms.Word{Priority: true, TileNumber: 1})
	_ = tm.Set(0, 31, sms.Word{VerticalFlip: true, TileNumber: 511})
	_ = tm.Set(23, 31, sms.Word{HorizontalFlip: true, TileNumber: 257})

	words := tm.Words()

	if len(words) != 896 {
		t.Fatalf("expected tilemap data to contain the full 896 entries, got %d", len(words))
	}

	if words[0] != 0b0001000000000001 {
		t.Errorf("expected tilemap to include tile #1, got %016b", words[0])
	}
	if words[31] != 0b0000010111111111 {
		t.Errorf("expected tilemap to include tile #2, got %016b", words[31])
	}
	if words[767] != 0b0000001100000001 {
		t.Errorf("expected tilemap to include tile #3, got %016b", words[767])
	}

	t.Run("unset entries in the tilemap return zero values", func(t *testing.T) {
		if words[1] != 0b0000000000000000 {
			t.Errorf("expected unset entry to be a zero value, got %016b", words[1])
		}
	})
}
