package captcha

import (
	"errors"
	"image"
	"image/color"
	"math"

	"github.com/vault-thirteen/RingCaptcha/pkg/brush"
	"github.com/vault-thirteen/RingCaptcha/pkg/geometry"
	img "github.com/vault-thirteen/RingCaptcha/pkg/image"
	math2 "github.com/vault-thirteen/RingCaptcha/pkg/math"
	rad "github.com/vault-thirteen/RingCaptcha/pkg/radian"
	"github.com/vault-thirteen/RingCaptcha/pkg/shape"
	"github.com/vault-thirteen/auxie/random"
)

const (
	CaptchaImageMinWidth    = 128
	CaptchaImageMinHeight   = 128
	RingMinCount            = 3
	RingMaxCount            = 6
	RingMinRadius           = 24
	BrushOuterRadiusMin     = 2
	BrushOuterRadiusMax     = 32
	ColourComponentMaxValue = 65535
	KR                      = 0.5
	KDMin                   = 1.0
	KDMax                   = 1.5
)

const (
	ErrDimensions         = "dimensions error"
	ErrBrushRadiusRatio   = "brush radius ratio error"
	ErrDensityCoefficient = "density coefficient error"
	ErrAnomaly            = "anomaly"
)

func CreateCaptchaImage(w, h uint) (canvas *image.NRGBA, ringCount uint, err error) {
	if (w < CaptchaImageMinWidth) || (h < CaptchaImageMinHeight) {
		return nil, 0, errors.New(ErrDimensions)
	}

	var bg *image.NRGBA
	bg, err = createBG(w, h)
	if err != nil {
		return nil, 0, err
	}

	ringCount, err = random.Uint(RingMinCount, selectMaxRingCount(w, h))
	if err != nil {
		return nil, 0, err
	}

	layers := make([]*image.NRGBA, 0, ringCount)
	var layer *image.NRGBA

	for i := uint(0); i < ringCount; i++ {
		layer, err = createRingLayer(w, h)
		if err != nil {
			return nil, 0, err
		}

		layers = append(layers, layer)
	}

	// Merge layers.
	canvas = image.NewNRGBA(image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: int(w), Y: int(h)}})
	img.FillCanvasWithColour(canvas, color.Transparent)
	for _, l := range layers {
		canvas, err = img.BlendImages(canvas, l)
		if err != nil {
			return nil, 0, err
		}
	}

	canvas, err = img.BlendImages(bg, canvas)
	if err != nil {
		return nil, 0, err
	}

	return canvas, ringCount, nil
}

func selectMaxRingCount(w, h uint) (maxRingCount uint) {
	minDim := math.Min(float64(w), float64(h))

	if minDim <= 128 {
		return math2.MinUint(RingMinCount+1, RingMaxCount)
	}
	if minDim <= 160 {
		return math2.MinUint(RingMinCount+2, RingMaxCount)
	}

	return RingMaxCount
}

func createRingLayer(w, h uint) (canvas *image.NRGBA, err error) {
	canvas = image.NewNRGBA(
		image.Rectangle{
			Min: image.Point{X: 0, Y: 0},
			Max: image.Point{X: int(w), Y: int(h)}},
	)

	var br *brush.SimpleBrush
	br, err = createBrush(w, h)
	if err != nil {
		return nil, err
	}

	// Select a center point.
	var center geometry.Point2D
	center, err = createPoint(w, h)
	if err != nil {
		return nil, err
	}

	// Select a ring radius.
	var ringRadius float64
	ringRadius, err = createRadius(uint(math.Round(math.Min(float64(w), float64(h)) * KR)))
	if err != nil {
		return nil, err
	}

	// Select a density.
	var da float64
	da, err = selectDensity(br.OuterRadius, ringRadius)
	if err != nil {
		return nil, err
	}

	// Draw a ring.
	shape.DrawRingWithSimpleBrush(canvas, br, center, ringRadius, da)

	return canvas, nil
}

func createBrush(canvasW, canvasH uint) (br *brush.SimpleBrush, err error) {
	var col *brush.BrushColour
	col, err = createColour()
	if (err != nil) || (col == nil) {
		return nil, err
	}

	br = &brush.SimpleBrush{
		Colour: *col,
	}

	minDimension := math.Min(float64(canvasW), float64(canvasH))

	maxRadius := uint(math.Round(math.Max(
		float64(BrushOuterRadiusMin),
		math.Min(
			minDimension/16.0,
			float64(BrushOuterRadiusMax)/1.0,
		),
	)))

	var brushOuterRadius uint
	brushOuterRadius, err = random.Uint(BrushOuterRadiusMin, maxRadius)
	if err != nil {
		return nil, err
	}

	br.OuterRadius = float64(brushOuterRadius)

	// Select a random ratio between inner and outer radius,
	// i.e. Rr = dRio / Ri, where dRio = Ro - Ri.
	// 1 => 1.0, i.e. dRio =   1 × Ri => Ro = Ri × 2	=> Ri = Ro / 2.0;
	// 2 => 0.5, i.e. dRio = 1/2 × Ri => Ro = Ri × 1.5	=> Ri = Ro / 1.5;
	// 3 => 2.0, i.e. dRio =   2 × Ri => Ro = Ri × 3	=> Ri = Ro / 3.0.
	var brushRadiusRatioType uint
	brushRadiusRatioType, err = random.Uint(1, 3)
	if err != nil {
		return nil, err
	}

	switch brushRadiusRatioType {
	case 1:
		br.InnerRadius = float64(brushOuterRadius) / 2.0
	case 2:
		br.InnerRadius = float64(brushOuterRadius) / 1.5
	case 3:
		br.InnerRadius = float64(brushOuterRadius) / 3.0
	default:
		return nil, errors.New(ErrBrushRadiusRatio)
	}

	return br, nil
}

func createColour() (c *brush.BrushColour, err error) {
	var Cr, Cg, Cb uint

	Cr, err = random.Uint(0, ColourComponentMaxValue)
	if err != nil {
		return nil, err
	}

	Cg, err = random.Uint(0, ColourComponentMaxValue)
	if err != nil {
		return nil, err
	}
	Cb, err = random.Uint(0, ColourComponentMaxValue)
	if err != nil {
		return nil, err
	}

	return &brush.BrushColour{
		R: float64(Cr) / ColourComponentMaxValue,
		G: float64(Cg) / ColourComponentMaxValue,
		B: float64(Cb) / ColourComponentMaxValue,
		A: 1.0,
	}, nil
}

func createPoint(xMax, yMax uint) (p geometry.Point2D, err error) {
	var x, y uint

	x, err = random.Uint(0, xMax)
	if err != nil {
		return geometry.Point2D{}, err
	}

	y, err = random.Uint(0, yMax)
	if err != nil {
		return geometry.Point2D{}, err
	}

	p = geometry.Point2D{
		X: float64(x),
		Y: float64(y),
	}

	return p, nil
}

func createRadius(rMax uint) (r float64, err error) {
	var t uint
	t, err = random.Uint(RingMinRadius, rMax)
	if err != nil {
		return 0.0, err
	}

	return float64(t), nil
}

func selectDensity(rB, rR float64) (da float64, err error) {
	// Select a random density coefficient.
	// 0 => 1.0 i.e. KDMin;
	// 5 => 1.5 i.e. KDMax.
	var kdSwitch uint
	kdSwitch, err = random.Uint(0, 5)
	if err != nil {
		return 0, err
	}

	// [0; 5] -> [1.0; 1.5]
	var kd = 1.0 + float64(kdSwitch)/10
	if (kd < KDMin) || (kd > KDMax) {
		return 0, errors.New(ErrDensityCoefficient)
	}

	return rad.RadianToDegree(kd * 2 * rB / rR), nil
}

func createBG(w, h uint) (bg *image.NRGBA, err error) {
	bg = image.NewNRGBA(image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: int(w), Y: int(h)}})

	var bgType uint

	bgType, err = random.Uint(1, 4)
	if err != nil {
		return nil, err
	}

	switch bgType {
	case 1:
		err = img.FillCanvasWithHGradient(bg, brush.ColourWhite, brush.ColourBlack)
	case 2:
		err = img.FillCanvasWithHGradient(bg, brush.ColourBlack, brush.ColourWhite)
	case 3:
		err = img.FillCanvasWithVGradient(bg, brush.ColourWhite, brush.ColourBlack)
	case 4:
		err = img.FillCanvasWithVGradient(bg, brush.ColourBlack, brush.ColourWhite)
	default:
		return nil, errors.New(ErrAnomaly)
	}

	if err != nil {
		return nil, err
	}

	return bg, nil
}
