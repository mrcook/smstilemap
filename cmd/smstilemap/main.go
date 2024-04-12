package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mrcook/smstilemap/cmd/smstilemap/processor"
)

const version = "0.1.1"

var (
	inputFilename   *string
	outputDirectory *string
	outputFormat    *string
	testLibrary     *bool
)

func init() {
	inputFilename = flag.String("in", "", "Input PNG filename")
	outputDirectory = flag.String("out", "", "Output directory for generated files (default: input filename directory)")
	outputFormat = flag.String("fmt", "asm", "Output format: asm, tiles")
	testLibrary = flag.Bool("test", false, "Test SMS library by generating a new PNG file")
	v := flag.Bool("v", false, "Display version number")

	flag.Parse()

	if *v {
		fmt.Printf("%s v%s\n", os.Args[0], version)
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
		if err := pro.SaveTilemapToImage(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	var err error

	switch *outputFormat {
	case "asm":
		err = pro.ToAssembly()
	case "tiles":
		err = pro.SaveTilesToImage()
	default:
		fmt.Println("ERROR: 'fmt' unknown output format!")
		fmt.Println()
		flag.Usage()
		os.Exit(2)
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
