package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"path"
	"strings"

	"github.com/mrcook/smstilemap/cmd/smstilemap/processor"
	"github.com/mrcook/smstilemap/sms"
)

var (
	inputFilename  *string
	outputFormat   *string
	outputFilename *string
)

func init() {
	inputFilename = flag.String("in", "", "Input PNG filename")
	outputFormat = flag.String("format", "png", "Format to output: png")
	outputFilename = flag.String("out", "", "Output PNG filename (optional)")
	v := flag.Bool("v", false, "Display version number")

	flag.Parse()

	if *v {
		fmt.Printf("%s v%s\n", os.Args[0], sms.Version)
		os.Exit(0)
	}

	if len(*inputFilename) == 0 {
		fmt.Println("ERROR: 'in' filename is required!")
		fmt.Println()
		flag.Usage()
		os.Exit(2)
	}

	// TODO: support SMS output
	*outputFormat = "png"

	if len(*outputFilename) == 0 {
		*outputFilename = setOutputFilename()
	}

	outputExtension := ".png"
	if path.Ext(*outputFilename) != outputExtension {
		*outputFilename += outputExtension
	}
}

func main() {
	srcImage, err := openImage(*inputFilename)
	if err != nil {
		log.Fatal(err)
	}

	// convert the PNG image to an SMS representation
	sega, err := processor.ImageToSms(srcImage)
	if err != nil {
		log.Fatal(err)
	}

	// convert SMS tilemap data back to a normal image
	dstImage, err := processor.SmsToImage(sega)
	if err != nil {
		log.Fatal(err)
	}

	// save to new png
	if err := saveImage(dstImage, *outputFilename); err != nil {
		log.Fatal(err)
	}
}

func setOutputFilename() string {
	dir := path.Dir(*inputFilename)
	file := path.Base(*inputFilename)

	file = strings.ReplaceAll(file, path.Ext(file), "")
	file += "-new.png"

	return path.Join(dir, file)
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
