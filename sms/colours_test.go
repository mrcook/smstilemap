package sms_test

import (
	"fmt"
	"testing"

	"github.com/mrcook/smstilemap/sms"
)

func TestColour_FromRGB(t *testing.T) {
	table := []struct {
		colour  sms.Colour
		r, g, b uint8
	}{
		{0b00000000, 0, 0, 0},
		{0b00010101, 85, 85, 85},
		{0b00101010, 170, 170, 170},
		{0b00111111, 255, 255, 255},
		{0b00000011, 255, 0, 0},
		{0b00001100, 0, 255, 0},
		{0b00110000, 0, 0, 255},
	}

	for index, testData := range table {
		data := sms.ColourDataForRGB(testData.r, testData.g, testData.b)

		t.Run(fmt.Sprintf("%08b should return correct Colour value", index), func(t *testing.T) {
			if data.Index != testData.colour {
				t.Errorf("expected correct colour, got: %08b", data.Index)
			}
		})
	}
}

func TestColour_ColourDataForNearestRGB(t *testing.T) {
	table := []struct {
		colour  sms.Colour
		r, g, b uint8
	}{
		{0b00000000, 0, 0, 52},
		{0b00010101, 53, 85, 127},
		{0b00101010, 128, 170, 202},
		{0b00111111, 203, 255, 255},
		{0b00000011, 203, 0, 0},
		{0b00001100, 0, 203, 0},
		{0b00110000, 0, 0, 203},
	}

	for index, testData := range table {
		data := sms.ColourDataForNearestRGB(testData.r, testData.g, testData.b)

		t.Run(fmt.Sprintf("it should match to the nearest colour: %08b", index), func(t *testing.T) {
			if data.Index != testData.colour {
				t.Errorf("expected correct match, got: %08b", data.Index)
			}
		})
	}
}
