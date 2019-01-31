package files

import (
	"bense4ger/image-resizer/images"
	"fmt"
	"image"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// FileHelper defines the shape of a struct that is used by the processor
type FileHelper interface {
	FileList() ([]os.FileInfo, error)
	ReadFile(filename string, fileChan chan *images.ImageContainer, errChan chan error)
}

// FSHelper is a type of FileHelper
type FSHelper struct {
	WorkingDir string
	Extension  string
}

// Must ensures that a FSHelper is configured correctly
func Must(h *FSHelper) (*FSHelper, error) {
	if len(h.WorkingDir) == 0 {
		return nil, fmt.Errorf("FSHelper: No working directory")
	}

	if len(h.Extension) == 0 {
		return nil, fmt.Errorf("FSHelper: No file extension")
	}

	return h, nil
}

// FileList gets the number of files in the reciever's working directory
func (h *FSHelper) FileList() ([]os.FileInfo, error) {
	files, err := ioutil.ReadDir(h.WorkingDir)
	if err != nil {
		return nil, fmt.Errorf("FileList: %s", err.Error())
	}

	list := make([]os.FileInfo, 0)
	for _, f := range files {
		if strings.HasSuffix(f.Name(), h.Extension) {
			list = append(list, f)
		}
	}

	return list, nil
}

// ReadFile reads a file from the reciever's working directory and outputs an image to the fileChan channel
func (h *FSHelper) ReadFile(fileName string, fileChan chan *images.ImageContainer, errChan chan error) {
	ip := path.Join(h.WorkingDir, fileName)
	rdr, err := os.Open(ip)
	if err != nil {
		errChan <- fmt.Errorf("ReadFile: %s", err.Error())
		return
	}
	defer rdr.Close()

	im, _, err := image.Decode(rdr)
	if err != nil {
		errChan <- fmt.Errorf("ReadFile: %s", err.Error())
		return
	}

	fileChan <- &images.ImageContainer{
		Name:  fileName,
		Image: im,
	}
}
