package gg_test

import (
	"fmt"
	"testing"

	"github.com/mrcook/smstilemap/gg"
)

func TestColour_FromRGB(t *testing.T) {
	table := []struct {
		colour  gg.Colour
		r, g, b uint8
	}{
		{0b0000000000000000, 0, 0, 0},
		{0b0000010101010101, 85, 85, 85},
		{0b0000101010101010, 170, 170, 170},
		{0b0000111111111111, 255, 255, 255},
		{0b0000000000001111, 255, 0, 0},
		{0b0000000011110000, 0, 255, 0},
		{0b0000111100000000, 0, 0, 255},
	}

	for index, testColour := range table {
		data := gg.ColourDataForRGB(testColour.r, testColour.g, testColour.b)

		t.Run(fmt.Sprintf("0b%16b should return correct Colour value", index), func(t *testing.T) {
			if data.Index != testColour.colour {
				t.Errorf("expected correct colour, got: %016b", data.Index)
			}
		})
	}
}

func TestColour_FromNearestMatchRGB(t *testing.T) {
	table := []struct {
		colour  gg.Colour
		r, g, b uint8
	}{
		{0b0000000000000000, 8, 0, 0},
		{0b0000000000000001, 9, 0, 0},
		{0b0000000000000001, 25, 0, 0},
		{0b0000000000000010, 26, 0, 0},
		{0b0000000000000010, 42, 0, 0},
		{0b0000000000000011, 43, 0, 0},
		{0b0000000000000011, 59, 0, 0},
		{0b0000000000000100, 60, 0, 0},
		{0b0000000000000100, 76, 0, 0},
		{0b0000000000000101, 77, 0, 0},
		{0b0000000000000101, 93, 0, 0},
		{0b0000000000000110, 94, 0, 0},
		{0b0000000000000110, 110, 0, 0},
		{0b0000000000000111, 111, 0, 0},
		{0b0000000000000111, 127, 0, 0},
		{0b0000000000001000, 128, 0, 0},
		{0b0000000000001000, 144, 0, 0},
		{0b0000000000001001, 145, 0, 0},
		{0b0000000000001001, 161, 0, 0},
		{0b0000000000001010, 162, 0, 0},
		{0b0000000000001010, 178, 0, 0},
		{0b0000000000001011, 179, 0, 0},
		{0b0000000000001011, 195, 0, 0},
		{0b0000000000001100, 196, 0, 0},
		{0b0000000000001100, 212, 0, 0},
		{0b0000000000001101, 213, 0, 0},
		{0b0000000000001101, 229, 0, 0},
		{0b0000000000001110, 230, 0, 0},
		{0b0000000000001110, 246, 0, 0},
		{0b0000000000001111, 247, 0, 0},
		{0b0000000000001111, 254, 0, 0},

		{0b0000000000000000, 0, 0, 8},
		{0b0000010101010101, 77, 85, 93},
		{0b0000101010101010, 162, 170, 178},
		{0b0000111111111111, 247, 255, 254},

		{0b0000000000001111, 247, 0, 0},
		{0b0000000011110000, 0, 247, 0},
		{0b0000111100000000, 0, 0, 247},
	}

	for index, testColour := range table {
		data := gg.ColourDataForNearestRGB(testColour.r, testColour.g, testColour.b)

		t.Run(fmt.Sprintf("it should match to the nearest colour: 0b%16b", index), func(t *testing.T) {
			if data.Index != testColour.colour {
				t.Errorf("expected %016b, got: %016b", testColour.colour, data.Index)
			}
		})
	}
}
