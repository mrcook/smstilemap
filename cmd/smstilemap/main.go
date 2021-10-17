package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mrcook/smstilemap/cmd/smstilemap/processor"
	"github.com/mrcook/smstilemap/sms"
)

var (
	inputFilename   *string
	outputDirectory *string
	testLibrary     *bool
)

func init() {
	inputFilename = flag.String("in", "", "Input PNG filename")
	outputDirectory = flag.String("dir", "", "Output directory for generated files (default: input filename directory)")
	testLibrary = flag.Bool("test", true, "Test SMS library by generating a new PNG file")
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
}

func main() {
	pro := processor.New(*inputFilename, *outputDirectory)

	if err := pro.CreateOutputDirectory(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := pro.PngToSMS(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if *testLibrary {
		if err := pro.ExportSmsToPngImage(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
