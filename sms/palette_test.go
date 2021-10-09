package sms_test

import (
	"testing"

	"github.com/mrcook/smstilemap/sms"
)

func TestPalette_SetColourAt(t *testing.T) {
	t.Run("with valid position", func(t *testing.T) {
		pal := sms.Palette{}
		colour := sms.FromRGB(85, 85, 85)
		if err := pal.SetColourAt(5, colour); err != nil {
			t.Fatalf("unexpected error, got '%s", err)
		}

		want := colour.SMS()
		got, _ := pal.ColourAt(5)
		if got.SMS() != want {
			t.Errorf("expected correct colour to be returned, got '0b%08b', expected '0b%08b'", got.SMS(), want)
		}
	})

	t.Run("when position is greater than palette size", func(t *testing.T) {
		pal := sms.Palette{}
		err := pal.SetColourAt(32, sms.Colour(0))
		if err == nil {
			t.Fatal("expected an error")
		}
		if err.Error() != "palette index out of bounds, got 32, max value is 31" {
			t.Errorf("expected correct error message, got '%s", err.Error())
		}
	})
}

func TestPalette_ColourAt(t *testing.T) {
	pal := sms.Palette{}

	t.Run("with valid position", func(t *testing.T) {
		colour := sms.FromRGB(170, 170, 170)
		_ = pal.SetColourAt(31, colour)

		want := colour.SMS()
		got, err := pal.ColourAt(31)

		if err != nil {
			t.Fatal("unexpected error")
		}
		if got.SMS() != want {
			t.Errorf("expected correct colour to be returned, got '0b%08b', expected '0b%08b'", got.SMS(), want)
		}
	})

	t.Run("when position is greater than palette size", func(t *testing.T) {
		_, err := pal.ColourAt(32)
		if err == nil {
			t.Fatal("expected an error")
		}
		if err.Error() != "palette index out of bounds, got 32, max value is 31" {
			t.Errorf("expected correct error message, got '%s", err.Error())
		}
	})
}
