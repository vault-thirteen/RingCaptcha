package c

import (
	"image"
	"math"
)

// SimpleBrush is a basic round brush with soft edges.
type SimpleBrush struct {
	InnerRadius float64
	OuterRadius float64
	Colour      BrushColour

	hasSample bool
	sample    *image.NRGBA
}

func (sb *SimpleBrush) GetSample() (sample *image.NRGBA) {
	if sb.hasSample {
		return sb.sample
	}

	sb.sample = sb.createSample()
	sb.hasSample = true

	return sb.sample
}

func (sb *SimpleBrush) createSample() (sample *image.NRGBA) {
	sample = image.NewNRGBA(image.Rect(0, 0, int(2*sb.OuterRadius), int(2*sb.OuterRadius)))
	br := &SimpleBrush{InnerRadius: sb.InnerRadius, OuterRadius: sb.OuterRadius, Colour: sb.Colour}
	p := Point2D{X: sb.OuterRadius, Y: sb.OuterRadius}
	UseSimpleBrush(sample, br, true, p)

	return sample
}

// UseSimpleBrush creates a point using a simple brush.
func UseSimpleBrush(canvas *image.NRGBA, brush *SimpleBrush, blend bool, u Point2D) {
	var rB, gB, bB, aB float64
	var rA, gA, bA, aA float64

	p0 := Point2D{X: math.Round(u.X - brush.OuterRadius), Y: math.Round(u.Y - brush.OuterRadius)}
	p1 := Point2D{X: math.Round(u.X + brush.OuterRadius), Y: math.Round(u.Y + brush.OuterRadius)}
	rA, gA, bA = brush.Colour.R, brush.Colour.G, brush.Colour.B
	dOIR := brush.OuterRadius - brush.InnerRadius

	var dtc float64 // Distance from current pixel to brush center point.
	var bcS, bcIdealS BrushColour
	switch blend {
	case true: // Use blending.
		for x := p0.X; x <= p1.X; x++ {
			for y := p0.Y; y <= p1.Y; y++ {
				bcS = brush.Colour
				rB, gB, bB, aB = GetColorComponentsS(canvas.At(int(x), int(y)))
				rA, gA, bA, _ = GetColorComponentsS(bcS)

				dtc = Distance2D(Point2D{X: x, Y: y}, u)
				if dtc <= brush.OuterRadius {
					if dtc > brush.InnerRadius {
						// Linear fade.
						aA = 1.0 - ((dtc - brush.InnerRadius) / dOIR)
						bcS.A = aA + aB*(1-aA)
						if bcS.A != 0 {
							bcS.R = (rA*aA + rB*aB*(1.0-aA)) / bcS.A
							bcS.G = (gA*aA + gB*aB*(1.0-aA)) / bcS.A
							bcS.B = (bA*aA + bB*aB*(1.0-aA)) / bcS.A
						} else {
							bcS.R = rA
							bcS.G = gA
							bcS.B = bA
						}
					}

					// Apply the bcS to a pixel.
					canvas.Set(int(x), int(y), bcS.StraightToPremultiplied())
				}
			}
		}

	case false: // Do not use blending.
		for x := p0.X; x <= p1.X; x++ {
			for y := p0.Y; y <= p1.Y; y++ {
				bcS = brush.Colour
				rB, gB, bB, aB = GetColorComponentsS(canvas.At(int(x), int(y)))
				rA, gA, bA, _ = GetColorComponentsS(bcS)

				dtc = Distance2D(Point2D{X: x, Y: y}, u)
				if dtc <= brush.OuterRadius {
					if dtc > brush.InnerRadius {
						// Linear fade.
						aA = 1.0 - ((dtc - brush.InnerRadius) / dOIR)
						bcS.A = aA + aB*(1-aA)
						if bcS.A != 0 {
							bcS.R = (rA*aA + rB*aB*(1.0-aA)) / bcS.A
							bcS.G = (gA*aA + gB*aB*(1.0-aA)) / bcS.A
							bcS.B = (bA*aA + bB*aB*(1.0-aA)) / bcS.A

							// Calculate the second (ideal BC).
							bcIdealS.A = aA
							bcIdealS.R = (rA * aA) / bcIdealS.A
							bcIdealS.G = (gA * aA) / bcIdealS.A
							bcIdealS.B = (bA * aA) / bcIdealS.A
						} else {
							bcS.R = rA
							bcS.G = gA
							bcS.B = bA

							// Calculate the second (ideal BC).
							bcIdealS = bcS
						}

						if bcIdealS.A < aB {
							continue
						}
						if bcS.A > bcIdealS.A {
							bcS = bcIdealS
						}

					}

					// Apply the bcS to a pixel.
					canvas.Set(int(x), int(y), bcS.StraightToPremultiplied())
				}
			}
		}
	}
}

// UseSimpleBrushS creates a point using a simple brush with a sample.
func UseSimpleBrushS(canvas *image.NRGBA, brush *SimpleBrush, blend bool, u Point2D) {
	var rB, gB, bB, aB float64
	var rA, gA, bA, aA float64

	p0 := Point2D{X: math.Round(u.X - brush.OuterRadius), Y: math.Round(u.Y - brush.OuterRadius)}
	p1 := Point2D{X: math.Round(u.X + brush.OuterRadius), Y: math.Round(u.Y + brush.OuterRadius)}
	rA, gA, bA = brush.Colour.R, brush.Colour.G, brush.Colour.B

	var bcS, bcIdealS BrushColour // Brush colour, straight.
	switch blend {
	case true: // Use blending.
		x0 := int(p0.X)
		y0 := int(p0.Y)
		xMax := int(p1.X)
		yMax := int(p1.Y)
		for x := x0; x <= xMax; x++ {
			for y := y0; y <= yMax; y++ {
				bcS = brush.Colour
				rB, gB, bB, aB = GetColorComponentsS(canvas.At(x, y))
				rA, gA, bA, aA = GetColorComponentsS(brush.GetSample().At(x-x0, y-y0))

				bcS.A = aA + aB*(1-aA)
				if bcS.A != 0 {
					bcS.R = (rA*aA + rB*aB*(1.0-aA)) / bcS.A
					bcS.G = (gA*aA + gB*aB*(1.0-aA)) / bcS.A
					bcS.B = (bA*aA + bB*aB*(1.0-aA)) / bcS.A
				} else {
					bcS.R = rA
					bcS.G = gA
					bcS.B = bA
				}

				canvas.Set(x, y, bcS.StraightToPremultiplied())
			}
		}

	case false: // Do not use blending.
		x0 := int(p0.X)
		y0 := int(p0.Y)
		xMax := int(p1.X)
		yMax := int(p1.Y)
		for x := x0; x <= xMax; x++ {
			for y := y0; y <= yMax; y++ {
				bcS = brush.Colour
				rB, gB, bB, aB = GetColorComponentsS(canvas.At(x, y))
				rA, gA, bA, aA = GetColorComponentsS(brush.GetSample().At(x-x0, y-y0))

				bcS.A = aA + aB*(1-aA)
				if bcS.A != 0 {
					bcS.R = (rA*aA + rB*aB*(1.0-aA)) / bcS.A
					bcS.G = (gA*aA + gB*aB*(1.0-aA)) / bcS.A
					bcS.B = (bA*aA + bB*aB*(1.0-aA)) / bcS.A

					// Calculate the second (ideal BC).
					bcIdealS.A = aA
					bcIdealS.R = (rA * aA) / bcIdealS.A
					bcIdealS.G = (gA * aA) / bcIdealS.A
					bcIdealS.B = (bA * aA) / bcIdealS.A
				} else {
					bcS.R = rA
					bcS.G = gA
					bcS.B = bA

					// Calculate the second (ideal BC).
					bcIdealS = bcS
				}

				if bcIdealS.A < aB {
					continue
				}
				if bcS.A > bcIdealS.A {
					bcS = bcIdealS
				}

				canvas.Set(x, y, bcS.StraightToPremultiplied())
			}
		}
	}
}
