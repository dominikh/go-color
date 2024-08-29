package color

// TODO:
// APCA
// Lstar
// Michelson
// WCAG21
// Weber
// DeltaPhiStar

func luminance(c *Color) float64 {
	return c.Convert(XYZ_D65).Values[1]
}

// ContrastWeber computes the Weber luminance contrast.
func ContrastWeber(c1, c2 *Color) float64 {
	y1 := max(luminance(c1), 0)
	y2 := max(luminance(c2), 0)

	if y2 > y1 {
		y1, y2 = y2, y1
	}

	if y2 == 0 {
		// the darkest sRGB color above black is #000001 and this produces a
		// plain Weber contrast of ~45647. So, setting the divide-by-zero result
		// at 50000 is a reasonable max clamp for the plain Weber
		return 50_000
	} else {
		return (y1 - y2) / y2
	}
}

// ContrastMichelson computes the Michelson contrast.
func ContrastMichelson(c1, c2 *Color) float64 {
	y1 := max(luminance(c1), 0)
	y2 := max(luminance(c2), 0)

	if y2 > y1 {
		y1, y2 = y2, y1
	}

	if y1+y2 == 0 {
		return 0
	}
	return (y1 - y2) / (y1 + y2)
}
