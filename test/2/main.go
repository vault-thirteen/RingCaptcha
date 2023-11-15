package main

import (
	"image"
	"image/color"
	"log"
	"path/filepath"

	rc "github.com/vault-thirteen/RingCaptcha"
	"github.com/vault-thirteen/RingCaptcha/test/common"
)

func main() {
	tdf := common.GetTestDataFolder()

	err := processImages(filepath.Join(tdf, common.OutputFolderName, "test_2.png"))
	if err != nil {
		log.Fatalln(err)
	}
}

func processImages(outputFilePath string) (err error) {
	var img = image.NewNRGBA(image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 500, Y: 500}})
	rc.FillCanvasWithColour(img, color.Transparent)
	paintSomeShit(img)

	err = rc.SaveImageAsPngFile(img, outputFilePath)
	if err != nil {
		return err
	}

	return nil
}

func paintSomeShit(canvas *image.NRGBA) {
	br1 := &rc.SolidBrush{InnerRadius: 16.0, Colour: rc.ColourGreen}

	p := rc.Point2D{X: 100, Y: 200}
	rc.UseSolidBrush(canvas, br1, p)

	br2 := &rc.SimpleBrush{InnerRadius: 16.0, OuterRadius: 32.0, Colour: rc.ColourGreen}

	p = rc.Point2D{X: 100, Y: 100}
	rc.UseSimpleBrush(canvas, br2, true, p)

	// Burn emulator.
	// With blending disabled, burning does not happen !
	p = rc.Point2D{X: 200, Y: 100}
	rc.UseSimpleBrush(canvas, br2, false, p)
	p = rc.Point2D{X: 250, Y: 100}
	rc.UseSimpleBrush(canvas, br2, false, p)
	for i := 1; i <= 50; i++ {
		p = rc.Point2D{X: 225, Y: 100}
		rc.UseSimpleBrush(canvas, br2, false, p)
	}

	// Line emulator.
	// Everything is good without blendinrc.
	rc.DrawLineWithSimpleBrush(canvas, br2, rc.Point2D{X: 200, Y: 200}, rc.Point2D{X: 400, Y: 200}, false)

	// Line #2.
	br3 := &rc.SimpleBrush{InnerRadius: 16.0, OuterRadius: 32.0, Colour: rc.ColourYellow}
	rc.DrawLineWithSimpleBrush(canvas, br3, rc.Point2D{X: 100, Y: 300}, rc.Point2D{X: 400, Y: 300}, false)

	// Line #3.
	// No blendinrc.
	br4 := &rc.SimpleBrush{InnerRadius: 16.0, OuterRadius: 32.0, Colour: rc.ColourCyan}
	rc.DrawLineWithSimpleBrush(canvas, br4, rc.Point2D{X: 350, Y: 150}, rc.Point2D{X: 350, Y: 350}, false)

	// Ring #1.
	br5 := &rc.SimpleBrush{InnerRadius: 4.0, OuterRadius: 6.0, Colour: rc.ColourMagenta}
	rc.DrawRingWithSimpleBrush(canvas, br5, rc.Point2D{X: 100, Y: 400}, 50, 15.0, false, false)
}
