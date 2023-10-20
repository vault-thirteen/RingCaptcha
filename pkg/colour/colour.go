package colour

import "image/color"

const (
	// C1 is a maximum colour channel value of a pixel in an 'image.RGBA' of
	// the built-in library.
	C1 = float64(65_535)
)

// GetColorComponentsPM extracts colour's components as usable floating point
// numbers. The returned values are pre-multiplied with alpha channel.
func GetColorComponentsPM(c color.Color) (r, g, b, a float64) {
	// Golang's RGBA Colour model returns channels as 16-bit values inside a
	// 32-bit variable. Holy shit !
	rCu32, gCu32, bCu32, aCu32 := c.RGBA()

	return float64(rCu32) / C1, float64(gCu32) / C1,
		float64(bCu32) / C1, float64(aCu32) / C1
}

// GetColorComponentsS extracts colour's components as usable floating point
// numbers. The returned values use straight alpha channel.
func GetColorComponentsS(c color.Color) (r, g, b, a float64) {
	// Golang's RGBA Colour model returns channels as 16-bit values inside a
	// 32-bit variable. Holy shit !
	rCu32, gCu32, bCu32, aCu32 := c.RGBA()

	if aCu32 != 0 {
		r = float64(rCu32) / float64(aCu32)
		g = float64(gCu32) / float64(aCu32)
		b = float64(bCu32) / float64(aCu32)
		a = float64(aCu32) / C1
		return r, g, b, a
	} else {
		return 0.0, 0.0, 0.0, 0.0
	}
}
