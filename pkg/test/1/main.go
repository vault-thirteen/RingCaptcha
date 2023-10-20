package main

import (
	"fmt"
	"image"
	"log"
	"path/filepath"

	"github.com/vault-thirteen/RingCaptcha/pkg/brush"
	g "github.com/vault-thirteen/RingCaptcha/pkg/geometry"
	im "github.com/vault-thirteen/RingCaptcha/pkg/image"
	"github.com/vault-thirteen/RingCaptcha/pkg/test/common"
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
		img, err = im.GetImageFromFilePath(imgPath)
		if err != nil {
			return fmt.Errorf("getImageFromFilePath error for file '%s': %w", imgPath, err)
		}

		rgbaImg = im.ConvertImageToRGBA(img)
		if err != nil {
			return fmt.Errorf("convertImageToRGBA error for file '%s': %w", imgPath, err)
		}

		imgs = append(imgs, rgbaImg)
	}

	//TODO
	ru := imgs[0].Rect.Union(imgs[1].Rect)

	imgB := im.DrawImageToCanvas(imgs[0], ru)
	imgA := im.DrawImageToCanvas(imgs[1], ru)

	var imgO *image.NRGBA
	imgO, err = im.BlendImages(imgB, imgA)
	if err != nil {
		return err
	}

	paintSomeShit(imgO)

	err = im.SaveImageAsPngFile(imgO, outputFilePath)
	if err != nil {
		return err
	}

	return nil
}

func paintSomeShit(canvas *image.NRGBA) {
	br1 := &brush.SolidBrush{InnerRadius: 32.0, Colour: brush.ColourRed}
	p := g.Point2D{X: 570, Y: 338}
	brush.UseSolidBrush(canvas, br1, p)

	br2 := &brush.SimpleBrush{InnerRadius: 16.0, OuterRadius: 32.0, Colour: brush.ColourGreen}
	p = g.Point2D{X: 472, Y: 345}
	brush.UseSimpleBrush(canvas, br2, true, p)

	br3 := &brush.SimpleBrush{InnerRadius: 16.0, OuterRadius: 32.0, Colour: brush.ColourBlue}
	p = g.Point2D{X: 650, Y: 250}
	brush.UseSimpleBrush(canvas, br3, true, p)

	br4 := &brush.SimpleBrush{InnerRadius: 16.0, OuterRadius: 32.0, Colour: brush.ColourBlue}
	p = g.Point2D{X: 750, Y: 500}
	brush.UseSimpleBrush(canvas, br4, true, p)
}
