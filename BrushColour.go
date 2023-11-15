package rc

import (
	"math"
)

type BrushColour struct {
	R float64
	G float64
	B float64
	A float64
}

var (
	ColourWhite   = BrushColour{R: 1.0, G: 1.0, B: 1.0, A: 1.0}
	ColourRed     = BrushColour{R: 1.0, G: 0.0, B: 0.0, A: 1.0}
	ColourGreen   = BrushColour{R: 0.0, G: 1.0, B: 0.0, A: 1.0}
	ColourBlue    = BrushColour{R: 0.0, G: 0.0, B: 1.0, A: 1.0}
	ColourCyan    = BrushColour{R: 0.0, G: 1.0, B: 1.0, A: 1.0}
	ColourMagenta = BrushColour{R: 1.0, G: 0.0, B: 1.0, A: 1.0}
	ColourYellow  = BrushColour{R: 1.0, G: 1.0, B: 0.0, A: 1.0}
	ColourBlack   = BrushColour{R: 0.0, G: 0.0, B: 0.0, A: 1.0}
)

// RGBA method of the color.Color interface.
func (bc BrushColour) RGBA() (r, g, b, a uint32) {
	return uint32(math.Round(bc.R * C1)),
		uint32(math.Round(bc.G * C1)),
		uint32(math.Round(bc.B * C1)),
		uint32(math.Round(bc.A * C1))
}

func (bc BrushColour) StraightToPremultiplied() (out BrushColour) {
	return BrushColour{
		R: bc.R * bc.A,
		G: bc.G * bc.A,
		B: bc.B * bc.A,
		A: bc.A,
	}
}

func (bc BrushColour) PremultipliedToStraight() (out BrushColour) {
	return BrushColour{
		R: bc.R / bc.A,
		G: bc.G / bc.A,
		B: bc.B / bc.A,
		A: bc.A,
	}
}
