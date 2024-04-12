package gg_test

import (
	"fmt"
	"testing"

	"github.com/mrcook/smstilemap/gg"
)

func TestColour_RGB(t *testing.T) {
	table := []struct {
		colour  uint16
		r, g, b uint8
		html    string
	}{
		{0b0000000000000000, 0, 0, 0, "#000000"},
		{0b0000010101010101, 85, 85, 85, "#555555"},
		{0b0000101010101010, 170, 170, 170, "#AAAAAA"},
		{0b0000111111111111, 255, 255, 255, "#FFFFFF"},
		{0b0000000000001111, 255, 0, 0, "#FF0000"},
		{0b0000000011110000, 0, 255, 0, "#00FF00"},
		{0b0000111100000000, 0, 0, 255, "#0000FF"},
	}

	for index, data := range table {
		col := gg.Colour(data.colour)

		t.Run(fmt.Sprintf("0b%16b should return correct GG value", index), func(t *testing.T) {
			b := col.GG()
			if b != data.colour {
				t.Errorf("expected %d, got: %d", data.colour, b)
			}
		})

		t.Run(fmt.Sprintf("0b%16b should return correct RGB values", index), func(t *testing.T) {
			r, b, g := col.RGB()
			if r != data.r && g != data.g && b != data.b {
				t.Errorf("expected correct RGB values, got: %d, %d, %d", r, g, b)
			}
		})

		t.Run(fmt.Sprintf("0b%16b should return correct RGBA values", index), func(t *testing.T) {
			r, b, g, a := col.RGBA()
			if r != uint32(data.r)*257 && g != uint32(data.g)*257 && b != uint32(data.b)*257 && a != 65535 {
				t.Errorf("expected correct RGB values, got: %d, %d, %d", r, g, b)
			}
		})

		t.Run(fmt.Sprintf("0b%16b should return correct HTML string", index), func(t *testing.T) {
			html := col.HTML()
			if html != data.html {
				t.Errorf("expected correct HTML string, got: '%s'", html)
			}
		})
	}
}

func TestColour_Equal(t *testing.T) {
	baseColour := 0b00101101
	colour := gg.Colour(baseColour)

	t.Run("returns true when the colours match", func(t *testing.T) {
		testColour := gg.Colour(baseColour)
		if !colour.Equal(testColour) {
			t.Errorf("expected the colours to match")
		}
	})

	t.Run("returns false when do not match", func(t *testing.T) {
		testColour := gg.Colour(0b00000001)
		if colour.Equal(testColour) {
			t.Errorf("expected 0b%08b to not match 0b%08b", testColour.GG(), colour.GG())
		}
	})
}
