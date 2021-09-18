package tilemap_test

import (
	"fmt"
	"testing"

	"github.com/mrcook/smstilemap/tilemap"
)

func TestWord_ToUint(t *testing.T) {
	table := map[uint16]tilemap.Word{
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
		word := tilemap.Word{TileNumber: 512}
		result := word.ToUint()

		if result != 0 {
			t.Errorf("expected tile number to default to zero, got %016b", result)
		}
	})
}
