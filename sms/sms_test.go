package sms_test

import (
	"testing"

	"github.com/mrcook/smstilemap/sms"
)

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
