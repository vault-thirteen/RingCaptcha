package c

import (
	"errors"
	"image"
	"image/color"
	"math"

	"github.com/vault-thirteen/RingCaptcha/models"
	"github.com/vault-thirteen/auxie/random"
)

func CreateCaptchaImage(w, h uint, useSample bool, blend bool) (canvas *image.NRGBA, ringCount uint, err error) {
	if (w < m.CaptchaImageMinWidth) || (h < m.CaptchaImageMinHeight) {
		return nil, 0, errors.New(m.Err_Dimensions)
	}

	var bg *image.NRGBA
	bg, err = createBG(w, h)
	if err != nil {
		return nil, 0, err
	}

	ringCount, err = random.Uint(m.RingMinCount, selectMaxRingCount(w, h))
	if err != nil {
		return nil, 0, err
	}

	layers := make([]*image.NRGBA, 0, ringCount)
	var layer *image.NRGBA

	for i := uint(0); i < ringCount; i++ {
		layer, err = createRingLayer(w, h, useSample, blend)
		if err != nil {
			return nil, 0, err
		}

		layers = append(layers, layer)
	}

	// Merge layers.
	canvas = image.NewNRGBA(image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: int(w), Y: int(h)}})
	FillCanvasWithColour(canvas, color.Transparent)
	for _, l := range layers {
		canvas, err = BlendImages(canvas, l)
		if err != nil {
			return nil, 0, err
		}
	}

	canvas, err = BlendImages(bg, canvas)
	if err != nil {
		return nil, 0, err
	}

	return canvas, ringCount, nil
}

func selectMaxRingCount(w, h uint) (maxRingCount uint) {
	minDim := math.Min(float64(w), float64(h))

	if minDim <= 128 {
		return MinUint(m.RingMinCount+1, m.RingMaxCount)
	}
	if minDim <= 160 {
		return MinUint(m.RingMinCount+2, m.RingMaxCount)
	}

	return m.RingMaxCount
}

func createRingLayer(w, h uint, useSample bool, blend bool) (canvas *image.NRGBA, err error) {
	canvas = image.NewNRGBA(
		image.Rectangle{
			Min: image.Point{X: 0, Y: 0},
			Max: image.Point{X: int(w), Y: int(h)},
		},
	)

	var br *SimpleBrush
	br, err = createBrush(w, h)
	if err != nil {
		return nil, err
	}

	// Select a center point.
	var center Point2D
	center, err = createPoint(w, h)
	if err != nil {
		return nil, err
	}

	// Select a ring radius.
	var ringRadius float64
	ringRadius, err = createRadius(uint(math.Round(math.Min(float64(w), float64(h)) * m.KR)))
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
	DrawRingWithSimpleBrush(canvas, br, center, ringRadius, da, useSample, blend)

	return canvas, nil
}

func createBrush(canvasW, canvasH uint) (br *SimpleBrush, err error) {
	var col *BrushColour
	col, err = createColour()
	if (err != nil) || (col == nil) {
		return nil, err
	}

	br = &SimpleBrush{Colour: *col}

	minDimension := math.Min(float64(canvasW), float64(canvasH))

	maxRadius := uint(math.Round(math.Max(float64(m.BrushOuterRadiusMin), math.Min(minDimension/16.0, float64(m.BrushOuterRadiusMax)/1.0))))

	var brushOuterRadius uint
	brushOuterRadius, err = random.Uint(m.BrushOuterRadiusMin, maxRadius)
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
		return nil, errors.New(m.Err_BrushRadiusRatio)
	}

	return br, nil
}

func createColour() (c *BrushColour, err error) {
	var Cr, Cg, Cb uint

	Cr, err = random.Uint(0, m.ColourComponentMaxValue)
	if err != nil {
		return nil, err
	}

	Cg, err = random.Uint(0, m.ColourComponentMaxValue)
	if err != nil {
		return nil, err
	}
	Cb, err = random.Uint(0, m.ColourComponentMaxValue)
	if err != nil {
		return nil, err
	}

	return &BrushColour{
		R: float64(Cr) / m.ColourComponentMaxValue,
		G: float64(Cg) / m.ColourComponentMaxValue,
		B: float64(Cb) / m.ColourComponentMaxValue,
		A: 1.0,
	}, nil
}

func createPoint(xMax, yMax uint) (p Point2D, err error) {
	var x, y uint

	x, err = random.Uint(0, xMax)
	if err != nil {
		return Point2D{}, err
	}

	y, err = random.Uint(0, yMax)
	if err != nil {
		return Point2D{}, err
	}

	p = Point2D{
		X: float64(x),
		Y: float64(y),
	}

	return p, nil
}

func createRadius(rMax uint) (r float64, err error) {
	var t uint
	t, err = random.Uint(m.RingMinRadius, rMax)
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
	if (kd < m.KDMin) || (kd > m.KDMax) {
		return 0, errors.New(m.Err_DensityCoefficient)
	}

	return RadianToDegree(kd * 2 * rB / rR), nil
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
		err = FillCanvasWithHGradient(bg, ColourWhite, ColourBlack)
	case 2:
		err = FillCanvasWithHGradient(bg, ColourBlack, ColourWhite)
	case 3:
		err = FillCanvasWithVGradient(bg, ColourWhite, ColourBlack)
	case 4:
		err = FillCanvasWithVGradient(bg, ColourBlack, ColourWhite)
	default:
		return nil, errors.New(m.Err_Anomaly)
	}

	if err != nil {
		return nil, err
	}

	return bg, nil
}
