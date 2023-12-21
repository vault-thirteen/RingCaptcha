package rc

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"math"
	"os"

	ae "github.com/vault-thirteen/auxie/errors"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/ccitt"
	_ "golang.org/x/image/colornames"
	_ "golang.org/x/image/riff"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/vector"
	_ "golang.org/x/image/vp8"
	_ "golang.org/x/image/vp8l"
	_ "golang.org/x/image/webp"
)

const (
	ErrCanvasIsTooSmall              = "canvas is too small"
	ErrImagesHaveDifferentDimensions = "images have different dimensions"
)

func GetImageFromFilePath(filePath string) (img image.Image, err error) {
	var f *os.File
	f, err = os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		derr := f.Close()
		if derr != nil {
			err = ae.Combine(err, derr)
		}
	}()

	img, _, err = image.Decode(f)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func ConvertImageToNRGBA(in image.Image) (out *image.NRGBA) {
	var ok bool
	out, ok = in.(*image.NRGBA)
	if ok {
		return out
	}

	b := in.Bounds()
	out = image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(out, out.Bounds(), in, b.Min, draw.Src)

	return out
}

func ConvertImageToRGBA(in image.Image) (out *image.RGBA) {
	var ok bool
	out, ok = in.(*image.RGBA)
	if ok {
		return out
	}

	b := in.Bounds()
	out = image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(out, out.Bounds(), in, b.Min, draw.Src)

	return out
}

// DrawImageToCanvas draws an image on another canvas.
func DrawImageToCanvas(in *image.RGBA, canvas image.Rectangle) (out *image.NRGBA) {
	bnds := in.Bounds()
	out = image.NewNRGBA(image.Rect(0, 0, canvas.Dx(), canvas.Dy()))
	draw.Draw(out, out.Bounds(), in, bnds.Min, draw.Src)

	return out
}

// BlendImages blends A and B images into O image.
// B is the base layer, i.e. the bottom layer.
// A is the applied layer, i.e. the top layer.
func BlendImages(iB *image.NRGBA, iA *image.NRGBA) (iO *image.NRGBA, err error) {
	if (iB.Rect.Dx() != iA.Rect.Dx()) || (iB.Rect.Dy() != iA.Rect.Dy()) {
		return nil, errors.New(ErrImagesHaveDifferentDimensions)
	}

	bnds := iB.Bounds()
	iO = image.NewNRGBA(image.Rect(0, 0, bnds.Dx(), bnds.Dy()))
	xMax := iB.Rect.Dx()
	yMax := iB.Rect.Dy()
	for x := 0; x < xMax; x++ {
		for y := 0; y < yMax; y++ {
			iO.Set(x, y, BlendColourOverlay(iB.At(x, y), iA.At(x, y)))
		}
	}

	return iO, nil
}

// BlendColourOverlay blends A and B pixels into O pixel.
// B is the base layer, i.e. the bottom layer.
// A is the applied layer, i.e. the top layer.
func BlendColourOverlay(cB, cA color.Color) (cO color.Color) {
	rB, gB, bB, aB := GetColorComponentsPM(cB)
	rA, gA, bA, aA := GetColorComponentsPM(cA)

	rO := BlendChannelOverlay(rB, rA, aB, aA)
	gO := BlendChannelOverlay(gB, gA, aB, aA)
	bO := BlendChannelOverlay(bB, bA, aB, aA)
	aO := BlendAlphaOverlay(aB, aA)

	return MakeRGBA64(rO, gO, bO, aO)
}

// MakeRGBA64 creates an RGBA64 colour with the provided channels.
func MakeRGBA64(r, g, b, a float64) color.RGBA64 {
	return color.RGBA64{
		R: uint16(math.Round(r * C1)),
		G: uint16(math.Round(g * C1)),
		B: uint16(math.Round(b * C1)),
		A: uint16(math.Round(a * C1)),
	}
}

// BlendChannelOverlay blends A and B channels into O channel.
// B is the base layer, i.e. the bottom layer.
// A is the applied layer, i.e. the top layer.
func BlendChannelOverlay(cB, cA, aB, aA float64) (cO float64) {
	if 2*cB < aB {
		return 2*cB*cA + cA*(1.0-aB) + cB*(1.0-aA)
	} else {
		return cB*(1.0+aA) + cA*(1.0+aB) - 2.0*cB*cA - aB*aA
	}
}

// BlendAlphaOverlay blends A and B alpha channels into O channel.
// B is the base layer, i.e. the bottom layer.
// A is the applied layer, i.e. the top layer.
func BlendAlphaOverlay(cB, cA float64) (cO float64) {
	return cB + cA - cB*cA
}

// FillCanvasWithColour fills the canvas with colour.
func FillCanvasWithColour(canvas *image.NRGBA, col color.Color) {
	xMax := canvas.Rect.Dx()
	yMax := canvas.Rect.Dy()

	for x := 0; x < xMax; x++ {
		for y := 0; y < yMax; y++ {
			canvas.Set(x, y, col)
		}
	}
}

// FillCanvasWithHGradient fills the canvas with a horizontal gradient.
// cL specifies the left side colour and cR is the right side colour.
func FillCanvasWithHGradient(canvas *image.NRGBA, cL, cR BrushColour) (err error) {
	xMax := canvas.Rect.Dx() // dX.
	yMax := canvas.Rect.Dy() // dY.

	minDim := MinInt(xMax, yMax)
	if minDim < 3 {
		return errors.New(ErrCanvasIsTooSmall)
	}

	// Colour deltas.
	var dCr = (cR.R - cL.R) / float64(xMax)
	var dCg = (cR.G - cL.G) / float64(xMax)
	var dCb = (cR.B - cL.B) / float64(xMax)
	var dCa = (cR.A - cL.A) / float64(xMax)

	var c = cL
	for x := 0; x < xMax; x++ {
		// Colour change.
		c.R += dCr
		c.G += dCg
		c.B += dCb
		c.A += dCa

		for y := 0; y < yMax; y++ {
			canvas.Set(x, y, c)
		}
	}

	return nil
}

// FillCanvasWithVGradient fills the canvas with a vertical gradient.
// cT specifies the top side colour and cB is the bottom side colour.
func FillCanvasWithVGradient(canvas *image.NRGBA, cT, cB BrushColour) (err error) {
	xMax := canvas.Rect.Dx() // dX.
	yMax := canvas.Rect.Dy() // dY.

	minDim := MinInt(xMax, yMax)
	if minDim < 3 {
		return errors.New(ErrCanvasIsTooSmall)
	}

	// Colour deltas.
	var dCr = (cB.R - cT.R) / float64(yMax)
	var dCg = (cB.G - cT.G) / float64(yMax)
	var dCb = (cB.B - cT.B) / float64(yMax)
	var dCa = (cB.A - cT.A) / float64(yMax)

	var c = cT
	for y := 0; y < yMax; y++ {
		// Colour change.
		c.R += dCr
		c.G += dCg
		c.B += dCb
		c.A += dCa

		for x := 0; x < xMax; x++ {
			canvas.Set(x, y, c)
		}
	}

	return nil
}
