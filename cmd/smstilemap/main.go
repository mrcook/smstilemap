package main

import (
	"flag"
	"image"
	"image/png"
	"log"
	"os"

	img "github.com/mrcook/smstilemap/image"
)

func main() {
	srcFilename := parseCliForFilename()
	dstFilename := srcFilename + "-new.png" // simple but works

	srcImage, err := openImage(srcFilename)
	if err != nil {
		log.Fatal(err)
	}

	// process it
	dstImage, err := processImage(srcImage)
	if err != nil {
		log.Fatal(err)
	}

	// save to new png
	if err := saveImage(dstImage, dstFilename); err != nil {
		log.Fatal(err)
	}
}

func processImage(tiled image.Image) (image.Image, error) {
	t, err := img.FromImage(tiled)
	if err != nil {
		panic(err)
	}
	t.FindAllUniqueTiles()
	dstImage := t.UniqueTilesToImage() // convert back to a RGBA image

	return dstImage, nil
}

func openImage(filename string) (image.Image, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	decodedImage, err := png.Decode(f)
	if err != nil {
		return nil, err
	}

	return decodedImage, nil
}

func saveImage(i image.Image, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	err = png.Encode(f, i)
	if err != nil {
		return err
	}
	return nil
}

func parseCliForFilename() string {
	filename := flag.String("src", "", "Source PNG image filename")
	flag.Parse()
	if len(*filename) == 0 {
		log.Fatal("source filename is required")
	}
	return *filename
}
