package main

import (
	"flag"
	"image"
	"image/png"
	"log"
	"os"

	"github.com/mrcook/smstilemap/sms"
)

func main() {
	srcFilename := parseCliForFilename()
	dstFilename := srcFilename + "-new.png" // simple but works

	pngImage, err := openPNG(srcFilename)
	if err != nil {
		log.Fatal(err)
	}

	vdp := sms.SMS{}

	// convert PNG image to a tiled representation
	if err := vdp.FromImage(pngImage); err != nil {
		log.Fatal(err)
	}

	// convert the tiles back to a normal image
	dstImage, err := vdp.ToImage()
	if err != nil {
		log.Fatal(err)
	}

	// save to new png
	if err := saveImage(dstImage, dstFilename); err != nil {
		log.Fatal(err)
	}
}

func openPNG(filename string) (image.Image, error) {
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
