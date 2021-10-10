package sms_test

import (
	"fmt"
	"testing"

	"github.com/mrcook/smstilemap/sms"
)

func TestColour_RGB(t *testing.T) {
	table := map[string]struct {
		colour  uint8
		r, g, b uint8
		html    string
	}{
		"black":     {0b00000000, 0, 0, 0, "#000000"},
		"darkgrey":  {0b00010101, 85, 85, 85, "#555555"},
		"lightgrey": {0b00101010, 170, 170, 170, "#AAAAAA"},
		"white":     {0b00111111, 255, 255, 255, "#FFFFFF"},
		"red":       {0b00000011, 255, 0, 0, "#FF0000"},
		"green":     {0b00001100, 0, 255, 0, "#00FF00"},
		"blue":      {0b00110000, 0, 0, 255, "#0000FF"},
	}

	for label, data := range table {
		col := sms.Colour(data.colour)

		t.Run(fmt.Sprintf("%s should return correct SMS value", label), func(t *testing.T) {
			b := col.SMS()
			if b != data.colour {
				t.Errorf("expected correct SMS value, got: %d", b)
			}
		})

		t.Run(fmt.Sprintf("%s should return correct RGB values", label), func(t *testing.T) {
			r, b, g := col.RGB()
			if r != data.r && g != data.g && b != data.b {
				t.Errorf("expected correct RGB values, got: %d, %d, %d", r, g, b)
			}
		})

		t.Run(fmt.Sprintf("%s should return correct HTML string", label), func(t *testing.T) {
			html := col.HTML()
			if html != data.html {
				t.Errorf("expected correct HTML string, got: '%s'", html)
			}
		})
	}
}

func TestColour_FromRGB(t *testing.T) {
	table := map[string]struct {
		colour  sms.Colour
		r, g, b uint8
	}{
		"black":     {0b00000000, 0, 0, 0},
		"darkgrey":  {0b00010101, 85, 85, 85},
		"lightgrey": {0b00101010, 170, 170, 170},
		"white":     {0b00111111, 255, 255, 255},
		"red":       {0b00000011, 255, 0, 0},
		"green":     {0b00001100, 0, 255, 0},
		"blue":      {0b00110000, 0, 0, 255},
	}

	for label, data := range table {
		col := sms.FromRGB(data.r, data.g, data.b)

		t.Run(fmt.Sprintf("%s should return correct Colour value", label), func(t *testing.T) {
			if col != data.colour {
				t.Errorf("expected correct colour, got: %08b", col)
			}
		})
	}
}

func TestColour_FromNearestMatchRGB(t *testing.T) {
	table := map[string]struct {
		colour  sms.Colour
		r, g, b uint8
	}{
		"black":     {0b00000000, 0, 0, 52},
		"darkgrey":  {0b00010101, 53, 85, 127},
		"lightgrey": {0b00101010, 128, 170, 202},
		"white":     {0b00111111, 203, 255, 255},
		"red":       {0b00000011, 203, 0, 0},
		"green":     {0b00001100, 0, 203, 0},
		"blue":      {0b00110000, 0, 0, 203},
	}

	for label, data := range table {
		col := sms.FromNearestMatchRGB(data.r, data.g, data.b)

		t.Run(fmt.Sprintf("it should match to the nearest colour: %s", label), func(t *testing.T) {
			if col != data.colour {
				t.Errorf("expected correct match, got: %08b", col)
			}
		})
	}
}

func TestColour_Equal(t *testing.T) {
	baseColour := 0b00101101
	colour := sms.Colour(baseColour)

	t.Run("returns true when the colours match", func(t *testing.T) {
		testColour := sms.Colour(baseColour)
		if !colour.Equal(testColour) {
			t.Errorf("expected the colours to match")
		}
	})

	t.Run("returns false when do not match", func(t *testing.T) {
		testColour := sms.Colour(0b00000001)
		if colour.Equal(testColour) {
			t.Errorf("expected 0b%08b to not match 0b%08b", testColour.SMS(), colour.SMS())
		}
	})
}
