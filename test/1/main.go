package main

import (
	"fmt"
	"image"
	"log"
	"path/filepath"

	rc "github.com/vault-thirteen/RingCaptcha"
	"github.com/vault-thirteen/RingCaptcha/test/common"
)

func main() {
	tdf := common.GetTestDataFolder()

	err := processImages(
		[]string{
			filepath.Join(tdf, "Cat.jpg"),
			filepath.Join(tdf, "Test.png"),
		},
		filepath.Join(tdf, common.OutputFolderName, "test_1.png"),
	)
	if err != nil {
		log.Fatalln(err)
	}
}

func processImages(imgPaths []string, outputFilePath string) (err error) {
	var imgs = make([]*image.RGBA, 0, len(imgPaths))

	var img image.Image
	var rgbaImg *image.RGBA
	for _, imgPath := range imgPaths {
		img, err = rc.GetImageFromFilePath(imgPath)
		if err != nil {
			return fmt.Errorf("getImageFromFilePath error for file '%s': %w", imgPath, err)
		}

		rgbaImg = rc.ConvertImageToRGBA(img)
		if err != nil {
			return fmt.Errorf("convertImageToRGBA error for file '%s': %w", imgPath, err)
		}

		imgs = append(imgs, rgbaImg)
	}

	ru := imgs[0].Rect.Union(imgs[1].Rect)

	imgB := rc.DrawImageToCanvas(imgs[0], ru)
	imgA := rc.DrawImageToCanvas(imgs[1], ru)

	var imgO *image.NRGBA
	imgO, err = rc.BlendImages(imgB, imgA)
	if err != nil {
		return err
	}

	paintSomeShit(imgO)

	err = rc.SaveImageAsPngFile(imgO, outputFilePath)
	if err != nil {
		return err
	}

	return nil
}

func paintSomeShit(canvas *image.NRGBA) {
	br1 := &rc.SolidBrush{InnerRadius: 32.0, Colour: rc.ColourRed}
	p := rc.Point2D{X: 570, Y: 338}
	rc.UseSolidBrush(canvas, br1, p)

	br2 := &rc.SimpleBrush{InnerRadius: 16.0, OuterRadius: 32.0, Colour: rc.ColourGreen}
	p = rc.Point2D{X: 472, Y: 345}
	rc.UseSimpleBrush(canvas, br2, true, p)

	br3 := &rc.SimpleBrush{InnerRadius: 16.0, OuterRadius: 32.0, Colour: rc.ColourBlue}
	p = rc.Point2D{X: 650, Y: 250}
	rc.UseSimpleBrush(canvas, br3, true, p)

	br4 := &rc.SimpleBrush{InnerRadius: 16.0, OuterRadius: 32.0, Colour: rc.ColourBlue}
	p = rc.Point2D{X: 750, Y: 500}
	rc.UseSimpleBrush(canvas, br4, true, p)
}
