package sms_test

import (
	"fmt"
	"testing"

	"github.com/mrcook/smstilemap/sms"
)

func TestWord_ToUint(t *testing.T) {
	table := map[uint16]sms.Word{
		0b0000000000000001: {false, false, false, false, 1},
		0b0000000111111111: {false, false, false, false, 511},
		0b0001000000000001: {true, false, false, false, 1},
		0b0000100000000001: {false, true, false, false, 1},
		0b0000010000000001: {false, false, true, false, 1},
		0b0000001000000001: {false, false, false, true, 1},
		0b0001111111111111: {true, true, true, true, 511},
	}

	for expected, word := range table {
		t.Run(fmt.Sprintf("converting tilemap word to %016b", expected), func(t *testing.T) {
			result := word.ToUint()
			if result != expected {
				t.Errorf("invalid tilemap entry, got %016b", result)
			}
		})
	}

	t.Run("when tile number exceeds 511", func(t *testing.T) {
		word := sms.Word{TileNumber: 512}
		result := word.ToUint()

		if result != 0 {
			t.Errorf("expected tile number to default to zero, got %016b", result)
		}
	})
}

func TestWord_SetFlippedStateFromOrientation(t *testing.T) {
	t.Run("when flipped vertically", func(t *testing.T) {
		word := sms.Word{
			HorizontalFlip: true, // to make sure this gets reset
		}
		word.SetFlippedStateFromOrientation(sms.OrientationFlippedV)
		if !word.VerticalFlip {
			t.Fatal("expected vertical flip to be set")
		}
		if word.HorizontalFlip {
			t.Fatal("expected horizontal flip to be reset")
		}
	})

	t.Run("when flipped horizontally", func(t *testing.T) {
		word := sms.Word{
			VerticalFlip: true, // to make sure this gets reset
		}
		word.SetFlippedStateFromOrientation(sms.OrientationFlippedH)
		if !word.HorizontalFlip {
			t.Fatal("expected horizontal flip to be set")
		}
		if word.VerticalFlip {
			t.Fatal("expected vertical flip to be reset")
		}
	})

	t.Run("when flipped vertically and horizontally", func(t *testing.T) {
		word := sms.Word{}
		word.SetFlippedStateFromOrientation(sms.OrientationFlippedVH)
		if !word.VerticalFlip {
			t.Fatal("expected vertical flip to be set")
		}
		if !word.HorizontalFlip {
			t.Fatal("expected horizontal flip to be set")
		}
	})

	t.Run("when normal orientation (not flipped)", func(t *testing.T) {
		word := sms.Word{
			VerticalFlip:   true, // to make sure this gets reset
			HorizontalFlip: true, // to make sure this gets reset
		}
		word.SetFlippedStateFromOrientation(sms.OrientationNormal)
		if word.VerticalFlip {
			t.Fatal("expected vertical flip to be reset")
		}
		if word.HorizontalFlip {
			t.Fatal("expected horizontal flip to be reset")
		}
	})
}
