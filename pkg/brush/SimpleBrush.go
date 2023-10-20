package brush

import (
	"image"
	"math"

	"github.com/vault-thirteen/RingCaptcha/pkg/colour"
	g "github.com/vault-thirteen/RingCaptcha/pkg/geometry"
)

// SimpleBrush is a basic round brush with soft edges.
type SimpleBrush struct {
	InnerRadius float64
	OuterRadius float64
	Colour      BrushColour
}

// UseSimpleBrush creates a point using a simple brush.
func UseSimpleBrush(canvas *image.NRGBA, brush *SimpleBrush, blend bool, u g.Point2D) {
	var rB, gB, bB, aB float64
	var rA, gA, bA, aA float64

	p0 := g.Point2D{X: math.Round(u.X - brush.OuterRadius), Y: math.Round(u.Y - brush.OuterRadius)}
	p1 := g.Point2D{X: math.Round(u.X + brush.OuterRadius), Y: math.Round(u.Y + brush.OuterRadius)}
	rA, gA, bA = brush.Colour.R, brush.Colour.G, brush.Colour.B
	dOIR := brush.OuterRadius - brush.InnerRadius

	var dtc float64 // Distance from current pixel to brush center point.
	var bcS, bcIdealS BrushColour
	switch blend {
	case true: // Use blending.
		for x := p0.X; x <= p1.X; x++ {
			for y := p0.Y; y <= p1.Y; y++ {
				bcS = brush.Colour
				rB, gB, bB, aB = colour.GetColorComponentsS(canvas.At(int(x), int(y)))
				rA, gA, bA, _ = colour.GetColorComponentsS(bcS)

				dtc = g.Distance2D(g.Point2D{X: x, Y: y}, u)
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
				rB, gB, bB, aB = colour.GetColorComponentsS(canvas.At(int(x), int(y)))
				rA, gA, bA, _ = colour.GetColorComponentsS(bcS)

				dtc = g.Distance2D(g.Point2D{X: x, Y: y}, u)
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
