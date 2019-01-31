package processor

import (
	"bense4ger/image-resizer/files"
	"bense4ger/image-resizer/images"
	"fmt"
	"image/jpeg"
	"log"
	"math"
	"os"
	"path"

	"github.com/nfnt/resize"
)

// ImageProcessor encapsulates image processing functionality
type ImageProcessor struct {
	InputDir   string
	OutputDir  string
	ResizePC   uint
	FileHelper files.FileHelper
}

// Must ensures that an ImageProcessor has been created correctly
func Must(ip *ImageProcessor) (*ImageProcessor, error) {
	if len(ip.InputDir) == 0 {
		return nil, fmt.Errorf("No input directory")
	}

	if len(ip.OutputDir) == 0 {
		return nil, fmt.Errorf("No output directory")
	}

	return ip, nil
}

// Do performs in the image processing
func (ip *ImageProcessor) Do() error {
	files, err := ip.FileHelper.FileList()
	if err != nil {
		return fmt.Errorf("Do: Error getting files: %s", err.Error())
	}

	if len(files) == 0 {
		log.Println("No files to process")
		return nil
	}

	ip.processImages(files)
	return nil
}

func (ip *ImageProcessor) processImages(files []os.FileInfo) {
	errors := make(chan error)
	origImg := make(chan *images.ImageContainer)
	done := make(chan bool)

	go ip.worker(len(files), origImg, errors, done)

	for _, f := range files {
		go ip.FileHelper.ReadFile(f.Name(), origImg, errors)
	}

	<-done
}

func (ip *ImageProcessor) worker(total int, inputImg chan *images.ImageContainer, errors chan error, done chan bool) {
	count := 0
	errCount := 0
	rsOutput := make(chan *images.ImageContainer)

	fwOk := make(chan string)

	for {
		select {
		case oI := <-inputImg:
			go ip.resizeImage(oI, rsOutput)
		case rI := <-rsOutput:
			go ip.writeImage(rI, fwOk, errors)
		case fn := <-fwOk:
			count++
			log.Printf("Output File (%d): %s\n", count, fn)
		case e := <-errors:
			errCount++
			log.Printf("Error (%d): %s\n", errCount, e.Error())
		}

		if (count + errCount) == total {
			done <- true
		}
	}
}

func (ip *ImageProcessor) writeImage(img *images.ImageContainer, output chan string, errors chan error) {
	fName := path.Join(ip.OutputDir, img.Name)
	f, err := os.Create(fName)
	if err != nil {
		errors <- fmt.Errorf("writeImage: %s", err.Error())
		return
	}
	defer f.Close()
	jpeg.Encode(f, img.Image, nil)

	output <- fName
}

func (ip *ImageProcessor) resizeImage(input *images.ImageContainer, output chan *images.ImageContainer) {
	oW := input.Image.Bounds().Max.X
	oH := input.Image.Bounds().Max.Y
	rs := float64(ip.ResizePC) / float64(100)
	nW := math.Floor(float64(oW) * rs)
	nH := math.Floor(float64(oH) * rs)

	rImg := resize.Resize(uint(nW), uint(nH), input.Image, resize.Lanczos3)

	input.Image = rImg
	output <- input
}
