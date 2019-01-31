package main

import (
	"bense4ger/image-resizer/files"
	"bense4ger/image-resizer/processor"
	"flag"
	"log"
	"os"
)

var (
	inDir  string
	outDir string
	pc     uint
)

func main() {
	fh, err := files.Must(&files.FSHelper{
		WorkingDir: inDir,
		Extension:  ".jpg",
	})
	if err != nil {
		log.Fatalf("Failed to created FileHelper: %s", err.Error())
	}

	p, err := processor.Must(&processor.ImageProcessor{
		InputDir:   inDir,
		OutputDir:  outDir,
		FileHelper: fh,
		ResizePC:   pc,
	})
	if err != nil {
		log.Fatalf("Failed to create ImageProcessor: %s", err.Error())
	}

	err = p.Do()
	if err != nil {
		log.Fatalf("Error executing: %s", err.Error())
	}

	log.Println("Complete")
	os.Exit(0)
}

func init() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting working directory: %s", err.Error())
	}

	flag.StringVar(&inDir, "indir", wd, "The image input directory.  Defaults to the directory this program is executed in")
	flag.StringVar(&outDir, "outdir", "", "The image output directory")
	flag.UintVar(&pc, "pc", 50, "The resize percent, defaults to 50")

	flag.Parse()
}
