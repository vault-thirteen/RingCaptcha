package rc

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"

	"github.com/vault-thirteen/errorz"
)

const ErrFFileExtensionMismatch = "file extension mismatch: png vs %s"

const FileExtPng = `.png`

func SaveImageAsPngFile(img image.Image, filePath string) (err error) {
	if filepath.Ext(filePath) != FileExtPng {
		return fmt.Errorf(ErrFFileExtensionMismatch, filepath.Ext(filePath))
	}

	var f *os.File
	f, err = os.Create(filePath)
	if err != nil {
		return err
	}

	defer func() {
		derr := f.Close()
		if derr != nil {
			err = errorz.Combine(err, derr)
		}
	}()

	err = png.Encode(f, img)
	if err != nil {
		return err
	}

	return nil
}

func makeRecordFilePath(imagesFolder string, id string) (path string) {
	return filepath.Join(imagesFolder, id+"."+ImageFileExt)
}
