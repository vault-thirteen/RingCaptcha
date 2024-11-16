package c

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/vault-thirteen/RingCaptcha/models"
	ae "github.com/vault-thirteen/auxie/errors"
)

func SaveImageAsPngFile(img image.Image, filePath string) (err error) {
	if filepath.Ext(filePath) != m.FileExtFullPng {
		return fmt.Errorf(m.ErrF_FileExtensionMismatch, filepath.Ext(filePath))
	}

	var f *os.File
	f, err = os.Create(filePath)
	if err != nil {
		return err
	}

	defer func() {
		derr := f.Close()
		if derr != nil {
			err = ae.Combine(err, derr)
		}
	}()

	err = png.Encode(f, img)
	if err != nil {
		return err
	}

	return nil
}

func DeleteImageFile(filePath string) (err error) {
	if filepath.Ext(filePath) != m.FileExtFullPng {
		return fmt.Errorf(m.ErrF_FileExtensionMismatch, filepath.Ext(filePath))
	}

	err = os.Remove(filePath)
	if err != nil {
		return err
	}

	return nil
}

func MakeRecordFilePath(imagesFolder string, id string) (path string) {
	return filepath.Join(imagesFolder, MakeFileName(id))
}

func MakeFileName(id string) (fileName string) {
	return id + "." + m.ImageFileExt
}

func FileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}
