package main

import (
	"image/png"
	"os"

	"github.com/mrcook/smstilemap/image"
)

func main() {
	f, err := os.Open("/home/michael/code-go/jetpac-sms-8bit.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	img, err := png.Decode(f)
	if err != nil {
		panic(err)
	}

	t, err := image.FromImage(img)
	if err != nil {
		panic(err)
	}

	t.FindAllUniqueTiles()

	// newImg := t.ToImage() // convert back to a RGBA image
	newImg := t.UniqueTilesToImage() // convert back to a RGBA image

	// save to new png
	nf, err := os.Create("/home/michael/code-go/test.png")
	if err != nil {
		panic(err)
	}
	defer nf.Close()
	err = png.Encode(nf, newImg)
	if err != nil {
		panic(err)
	}
}
