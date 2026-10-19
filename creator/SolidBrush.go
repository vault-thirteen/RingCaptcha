package c

import (
	"image"
	"math"
)

// SolidBrush is a basic round brush without any soft edges.
type SolidBrush struct {
	InnerRadius float64
	Colour      BrushColour
}

// UseSolidBrush creates a point using a solid brush.
func UseSolidBrush(canvas *image.NRGBA, brush *SolidBrush, u Point2D) {
	var rB, gB, bB, aB float64
	var rA, gA, bA, aA float64

	p0 := Point2D{X: math.Round(u.X - brush.InnerRadius), Y: math.Round(u.Y - brush.InnerRadius)}
	p1 := Point2D{X: math.Round(u.X + brush.InnerRadius), Y: math.Round(u.Y + brush.InnerRadius)}
	rA, gA, bA = brush.Colour.R, brush.Colour.G, brush.Colour.B

	var dtc float64 // Distance from current pixel to brush center point.
	var bc BrushColour
	for x := p0.X; x <= p1.X; x++ {
		for y := p0.Y; y <= p1.Y; y++ {
			bc = brush.Colour
			rB, gB, bB, aB = GetColorComponentsS(canvas.At(int(x), int(y)))

			dtc = Distance2D(Point2D{X: x, Y: y}, u)
			if dtc <= brush.InnerRadius {
				if math.Abs(dtc-brush.InnerRadius) < 0.1 {
					// Make the rough edge a bit softer.
					aA = 0.5
					bc.A = aA + aB*(1-aA)
					if bc.A != 0 {
						bc.R = (rA*aA + rB*aB*(1-aA)) / bc.A
						bc.G = (gA*aA + gB*aB*(1-aA)) / bc.A
						bc.B = (bA*aA + bB*aB*(1-aA)) / bc.A
					} else {
						bc.R = rA
						bc.G = gA
						bc.B = bA
					}
				}

				// Apply the colour to a pixel.
				canvas.Set(int(x), int(y), bc)
			}
		}
	}
}
