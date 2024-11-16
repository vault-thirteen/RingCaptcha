package main

import (
	"fmt"
	"image"
	"log"
	"path/filepath"

	"github.com/vault-thirteen/RingCaptcha/creator"
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
		img, err = c.GetImageFromFilePath(imgPath)
		if err != nil {
			return fmt.Errorf("getImageFromFilePath error for file '%s': %w", imgPath, err)
		}

		rgbaImg = c.ConvertImageToRGBA(img)

		imgs = append(imgs, rgbaImg)
	}

	ru := imgs[0].Rect.Union(imgs[1].Rect)

	imgB := c.DrawImageToCanvas(imgs[0], ru)
	imgA := c.DrawImageToCanvas(imgs[1], ru)

	var imgO *image.NRGBA
	imgO, err = c.BlendImages(imgB, imgA)
	if err != nil {
		return err
	}

	paintSomeShit(imgO)

	err = c.SaveImageAsPngFile(imgO, outputFilePath)
	if err != nil {
		return err
	}

	return nil
}

func paintSomeShit(canvas *image.NRGBA) {
	br1 := &c.SolidBrush{InnerRadius: 32.0, Colour: c.ColourRed}
	p := c.Point2D{X: 570, Y: 338}
	c.UseSolidBrush(canvas, br1, p)

	br2 := &c.SimpleBrush{InnerRadius: 16.0, OuterRadius: 32.0, Colour: c.ColourGreen}
	p = c.Point2D{X: 472, Y: 345}
	c.UseSimpleBrush(canvas, br2, true, p)

	br3 := &c.SimpleBrush{InnerRadius: 16.0, OuterRadius: 32.0, Colour: c.ColourBlue}
	p = c.Point2D{X: 650, Y: 250}
	c.UseSimpleBrush(canvas, br3, true, p)

	br4 := &c.SimpleBrush{InnerRadius: 16.0, OuterRadius: 32.0, Colour: c.ColourBlue}
	p = c.Point2D{X: 750, Y: 500}
	c.UseSimpleBrush(canvas, br4, true, p)
}
