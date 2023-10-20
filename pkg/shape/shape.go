package shape

import (
	"image"
	"math"

	"github.com/vault-thirteen/RingCaptcha/pkg/brush"
	g "github.com/vault-thirteen/RingCaptcha/pkg/geometry"
	rad "github.com/vault-thirteen/RingCaptcha/pkg/radian"
)

// DrawLineWithSimpleBrush draws a straight line using a simple brush.
func DrawLineWithSimpleBrush(canvas *image.NRGBA, br *brush.SimpleBrush, p1, p2 g.Point2D) {
	dx := p2.X - p1.X
	dy := p2.Y - p1.Y
	var p g.Point2D
	if dx > dy {
		dj := dy / dx
		for i := float64(0); i <= dx; i++ {
			p.X = p1.X + i
			p.Y = p1.Y + dj
			brush.UseSimpleBrush(canvas, br, false, p)
		}
	} else {
		dj := dx / dy
		for i := float64(0); i <= dy; i++ {
			p.Y = p1.Y + i
			p.X = p1.X + dj
			brush.UseSimpleBrush(canvas, br, false, p)
		}
	}

}

func DrawRingWithSimpleBrush(canvas *image.NRGBA, br *brush.SimpleBrush, center g.Point2D, radius float64, daDegree float64) {
	da := rad.DegreeToRadian(daDegree)                            // Delta alpha.
	aMax := 2*math.Pi - 1.2*calculatePhiB(br.OuterRadius, radius) // Alpha maximum.
	n := int(math.Floor(aMax / da))                               // Number of brush uses.
	daNormalized := 2 * math.Pi / float64(n)                      // Delta alpha after aligning.

	a := 0.0
	var p g.Point2D
	for i := 0; i < n; i++ {
		p.X = math.Round(center.X + radius*math.Cos(a))
		p.Y = math.Round(center.Y + radius*math.Sin(a))
		brush.UseSimpleBrush(canvas, br, false, p)

		// Next.
		a = a + daNormalized
	}
}

// calculatePhiB calculates an angle of the ring circle taken by a brush's
// circle.
func calculatePhiB(rB, rR float64) (phiB float64) {
	cosPhi1 := (2*rR*rR - rB*rB) / (2 * rR * rR)
	phi1 := math.Acos(cosPhi1)

	return 2 * phi1
}
