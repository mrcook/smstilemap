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
	if tm.Height() != 28 {
		t.Errorf("expected tilemap height to be 28, got %d", tm.Height())
	}
}

func TestTilemap_VisibleHeight(t *testing.T) {
	tm := sms.Tilemap{}
	if tm.VisibleHeight() != 24 {
		t.Errorf("expected tilemap visible height to be 24, got %d", tm.VisibleHeight())
	}
}

func TestTilemap_Get(t *testing.T) {
	tm := sms.Tilemap{}
	word := sms.Word{TileNumber: 447}
	_ = tm.Set(27, 0, word)

	t.Run("with valid data", func(t *testing.T) {
		got, err := tm.Get(27, 0)
		if err != nil {
			t.Fatalf("unexpected error: %q", err)
		}
		if got.TileNumber != word.TileNumber {
			t.Errorf("expected to get correct data, got tile id %d", got.TileNumber)
		}
	})

	t.Run("with out of bounds row indexing", func(t *testing.T) {
		_, err := tm.Get(28, 31)
		if err == nil {
			t.Fatal("expected an error")
		}
		if err.Error() != "get tilemap out of bounds indexing, max is (27,31), requested (28,31)" {
			t.Errorf("unexpected error message, requested '%s'", err)
		}
	})

	t.Run("with out of bounds col indexing", func(t *testing.T) {
		_, err := tm.Get(27, 32)
		if err == nil {
			t.Fatal("expected an error")
		}
		if err.Error() != "get tilemap out of bounds indexing, max is (27,31), requested (27,32)" {
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
		err := tm.Set(28, 31, word)
		if err == nil {
			t.Fatal("expected an error")
		}
		if err.Error() != "set tilemap out of bounds indexing, max is (27,31), requested (28,31)" {
			t.Errorf("unexpected error message, got '%s'", err)
		}
	})

	t.Run("with out of bounds col indexing", func(t *testing.T) {
		err := tm.Set(27, 32, word)
		if err == nil {
			t.Fatal("expected an error")
		}
		if err.Error() != "set tilemap out of bounds indexing, max is (27,31), requested (27,32)" {
			t.Errorf("unexpected error message, got '%s'", err)
		}
	})
}
