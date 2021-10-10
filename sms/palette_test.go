package sms_test

import (
	"testing"

	"github.com/mrcook/smstilemap/sms"
)

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
		} else if err.Error() != "palette index out of bounds, got 32, max value is 31" {
			t.Errorf("expected correct error message, got '%s", err.Error())
		}
	})

	t.Run("when colour at position is not enabled", func(t *testing.T) {
		_, err := pal.ColourAt(0)
		if err == nil {
			t.Fatal("expected an error")
		} else if err.Error() != "uninitialised colour" {
			t.Errorf("expected correct error message, got '%s", err.Error())
		}
	})
}

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
		} else if err.Error() != "palette index out of bounds, got 32, max value is 31" {
			t.Errorf("expected correct error message, got '%s", err.Error())
		}
	})
}

func TestPalette_AddColour(t *testing.T) {
	t.Run("and return first position", func(t *testing.T) {
		pal := sms.Palette{}
		pos, err := pal.AddColour(sms.Colour(0b00010101))
		if err != nil {
			t.Fatalf("unexpected error, got '%s'", err)
		}
		if pos != 0 {
			t.Errorf("expected colour to be added at first slot, got %d", pos)
		}
	})

	t.Run("add colour to first available slot", func(t *testing.T) {
		pal := sms.Palette{}
		_ = pal.SetColourAt(0, sms.Colour(0b00111111))
		_ = pal.SetColourAt(2, sms.Colour(0b00111111))

		pos, err := pal.AddColour(sms.Colour(0b00000001))
		if err != nil {
			t.Fatalf("unexpected error, got '%s'", err)
		}
		if pos != 1 {
			t.Errorf("expected colour to be added at second slot, got %d", pos)
		}
	})

	t.Run("with an existing colour, return its position", func(t *testing.T) {
		pal := sms.Palette{}
		colour := sms.Colour(0b00111111)
		_ = pal.SetColourAt(2, colour)

		pos, err := pal.AddColour(colour)
		if err != nil {
			t.Fatalf("unexpected error, got '%s'", err)
		}
		if pos != 2 {
			t.Errorf("expected existing position, got %d", pos)
		}
	})

	t.Run("when two identical colours are present, return position of first", func(t *testing.T) {
		pal := sms.Palette{}
		colour := sms.Colour(0b00111111)
		_ = pal.SetColourAt(5, colour)
		_ = pal.SetColourAt(20, colour)

		pos, err := pal.AddColour(colour)
		if err != nil {
			t.Fatalf("unexpected error, got '%s'", err)
		}
		if pos != 5 {
			t.Errorf("expected colour to be added at second slot, got %d", pos)
		}
	})

	t.Run("when palette is full, return an error", func(t *testing.T) {
		pal := sms.Palette{}
		colour := sms.Colour(0b00111111)
		for i := sms.PaletteId(0); i < 32; i++ {
			_ = pal.SetColourAt(i, colour)
		}

		_, err := pal.AddColour(sms.Colour(0b00000011))
		if err == nil {
			t.Fatalf("expected error")
		} else if err.Error() != "palette full" {
			t.Errorf("expect a valid error message, got '%s'", err)
		}
	})
}

func TestPalette_PaletteIdFor(t *testing.T) {
	t.Run("return position of matching colour", func(t *testing.T) {
		pal := sms.Palette{}
		colour := sms.Colour(0b00111111)
		_ = pal.SetColourAt(2, colour)

		pos, err := pal.PaletteIdFor(colour)
		if err != nil {
			t.Fatalf("unexpected error, got '%s'", err)
		}
		if pos != 2 {
			t.Errorf("expected existing position, got %d", pos)
		}
	})

	t.Run("when no colour match found", func(t *testing.T) {
		pal := sms.Palette{}
		_ = pal.SetColourAt(2, sms.Colour(0b00111111))

		_, err := pal.PaletteIdFor(sms.Colour(0b00000011))
		if err == nil {
			t.Fatal("expected an error")
		} else if err.Error() != "colour not found" {
			t.Errorf("unexpected error message, got '%s'", err)
		}
	})
}
