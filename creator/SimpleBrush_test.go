package c

import (
	"image"
	"testing"

	"github.com/vault-thirteen/auxie/tester"
)

// go test -bench=. -benchmem

func Benchmark_UseSimpleBrush(b *testing.B) {
	w := 640
	h := 480
	canvas := image.NewNRGBA(image.Rect(0, 0, w, h))
	br := &SimpleBrush{InnerRadius: 16.0, OuterRadius: 32.0, Colour: ColourGreen}
	p := Point2D{X: 0, Y: 0}

	for i := 0; i < b.N; i++ {
		UseSimpleBrush(canvas, br, true, p)
	}
}

func Benchmark_UseSimpleBrushS(b *testing.B) {
	w := 640
	h := 480
	canvas := image.NewNRGBA(image.Rect(0, 0, w, h))
	br := &SimpleBrush{InnerRadius: 16.0, OuterRadius: 32.0, Colour: ColourGreen}
	p := Point2D{X: 0, Y: 0}

	for i := 0; i < b.N; i++ {
		UseSimpleBrushS(canvas, br, true, p)
	}
}

func Test_GetSample(t *testing.T) {
	aTest := tester.New(t)
	br := &SimpleBrush{InnerRadius: 16.0, OuterRadius: 32.0, Colour: ColourGreen}
	aTest.MustBeEqual(br.sample, (*image.NRGBA)(nil))
	s := br.GetSample()
	aTest.MustBeDifferent(br.sample, (*image.NRGBA)(nil))
	aTest.MustBeEqual(br.sample, s)
}
