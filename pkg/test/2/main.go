package main

import (
	"image"
	"image/color"
	"log"
	"path/filepath"

	br "github.com/vault-thirteen/RingCaptcha/pkg/brush"
	g "github.com/vault-thirteen/RingCaptcha/pkg/geometry"
	im "github.com/vault-thirteen/RingCaptcha/pkg/image"
	"github.com/vault-thirteen/RingCaptcha/pkg/os"
	"github.com/vault-thirteen/RingCaptcha/pkg/shape"
	"github.com/vault-thirteen/RingCaptcha/pkg/test/common"
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
	im.FillCanvasWithColour(img, color.Transparent)
	paintSomeShit(img)

	err = os.SaveImageAsPngFile(img, outputFilePath)
	if err != nil {
		return err
	}

	return nil
}

func paintSomeShit(canvas *image.NRGBA) {
	br1 := &br.SolidBrush{InnerRadius: 16.0, Colour: br.ColourGreen}

	p := g.Point2D{X: 100, Y: 200}
	br.UseSolidBrush(canvas, br1, p)

	br2 := &br.SimpleBrush{InnerRadius: 16.0, OuterRadius: 32.0, Colour: br.ColourGreen}

	p = g.Point2D{X: 100, Y: 100}
	br.UseSimpleBrush(canvas, br2, true, p)

	// Burn emulator.
	// With blending disabled, burning does not happen !
	p = g.Point2D{X: 200, Y: 100}
	br.UseSimpleBrush(canvas, br2, false, p)
	p = g.Point2D{X: 250, Y: 100}
	br.UseSimpleBrush(canvas, br2, false, p)
	for i := 1; i <= 50; i++ {
		p = g.Point2D{X: 225, Y: 100}
		br.UseSimpleBrush(canvas, br2, false, p)
	}

	// Line emulator.
	// Everything is good without blending.
	shape.DrawLineWithSimpleBrush(canvas, br2, g.Point2D{X: 200, Y: 200}, g.Point2D{X: 400, Y: 200}, false)

	// Line #2.
	br3 := &br.SimpleBrush{InnerRadius: 16.0, OuterRadius: 32.0, Colour: br.ColourYellow}
	shape.DrawLineWithSimpleBrush(canvas, br3, g.Point2D{X: 100, Y: 300}, g.Point2D{X: 400, Y: 300}, false)

	// Line #3.
	// No blending.
	br4 := &br.SimpleBrush{InnerRadius: 16.0, OuterRadius: 32.0, Colour: br.ColourCyan}
	shape.DrawLineWithSimpleBrush(canvas, br4, g.Point2D{X: 350, Y: 150}, g.Point2D{X: 350, Y: 350}, false)

	// Ring #1.
	br5 := &br.SimpleBrush{InnerRadius: 4.0, OuterRadius: 6.0, Colour: br.ColourMagenta}
	shape.DrawRingWithSimpleBrush(canvas, br5, g.Point2D{X: 100, Y: 400}, 50, 15.0, false, false)
}
