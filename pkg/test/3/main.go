package main

import (
	"fmt"
	"image"
	"log"
	"path/filepath"

	"github.com/vault-thirteen/RingCaptcha/pkg/captcha"
	im "github.com/vault-thirteen/RingCaptcha/pkg/image"
	"github.com/vault-thirteen/RingCaptcha/pkg/test/common"
	"github.com/vault-thirteen/auxie/random"
)

func main() {
	tdf := common.GetTestDataFolder()

	err := processImages(filepath.Join(tdf, common.OutputFolderName))
	if err != nil {
		log.Fatalln(err)
	}
}

func processImages(outputFolderPath string) (err error) {
	for i := 1; i <= 10; i++ {
		err = processImage(outputFolderPath, i)
		if err != nil {
			return err
		}
	}

	return nil
}

func processImage(outputFolderPath string, n int) (err error) {
	outputFilePath := filepath.Join(outputFolderPath, fmt.Sprintf("%v.png", n))

	var dim uint
	dim, err = random.Uint(128, 320)
	if err != nil {
		return err
	}

	var canvas *image.NRGBA
	var ringCount uint
	canvas, ringCount, err = captcha.CreateCaptchaImage(dim, dim)
	if err != nil {
		return err
	}

	err = im.SaveImageAsPngFile(canvas, outputFilePath)
	if err != nil {
		return err
	}

	log.Println(fmt.Sprintf("w=%v, h=%v, n=%v", dim, dim, ringCount))

	return nil
}
