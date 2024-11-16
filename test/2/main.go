package main

import (
	"image"
	"image/color"
	"log"
	"path/filepath"

	"github.com/vault-thirteen/RingCaptcha/creator"
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
	c.FillCanvasWithColour(img, color.Transparent)
	paintSomeShit(img)

	err = c.SaveImageAsPngFile(img, outputFilePath)
	if err != nil {
		return err
	}

	return nil
}

func paintSomeShit(canvas *image.NRGBA) {
	br1 := &c.SolidBrush{InnerRadius: 16.0, Colour: c.ColourGreen}

	p := c.Point2D{X: 100, Y: 200}
	c.UseSolidBrush(canvas, br1, p)

	br2 := &c.SimpleBrush{InnerRadius: 16.0, OuterRadius: 32.0, Colour: c.ColourGreen}

	p = c.Point2D{X: 100, Y: 100}
	c.UseSimpleBrush(canvas, br2, true, p)

	// Burn emulator.
	// With blending disabled, burning does not happen !
	p = c.Point2D{X: 200, Y: 100}
	c.UseSimpleBrush(canvas, br2, false, p)
	p = c.Point2D{X: 250, Y: 100}
	c.UseSimpleBrush(canvas, br2, false, p)
	for i := 1; i <= 50; i++ {
		p = c.Point2D{X: 225, Y: 100}
		c.UseSimpleBrush(canvas, br2, false, p)
	}

	// Line emulator.
	// Everything is good without blending.
	c.DrawLineWithSimpleBrush(canvas, br2, c.Point2D{X: 200, Y: 200}, c.Point2D{X: 400, Y: 200}, false)

	// Line #2.
	br3 := &c.SimpleBrush{InnerRadius: 16.0, OuterRadius: 32.0, Colour: c.ColourYellow}
	c.DrawLineWithSimpleBrush(canvas, br3, c.Point2D{X: 100, Y: 300}, c.Point2D{X: 400, Y: 300}, false)

	// Line #3.
	// No blending.
	br4 := &c.SimpleBrush{InnerRadius: 16.0, OuterRadius: 32.0, Colour: c.ColourCyan}
	c.DrawLineWithSimpleBrush(canvas, br4, c.Point2D{X: 350, Y: 150}, c.Point2D{X: 350, Y: 350}, false)

	// Ring #1.
	br5 := &c.SimpleBrush{InnerRadius: 4.0, OuterRadius: 6.0, Colour: c.ColourMagenta}
	c.DrawRingWithSimpleBrush(canvas, br5, c.Point2D{X: 100, Y: 400}, 50, 15.0, false, false)
}
