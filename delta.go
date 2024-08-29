package color

import "math"

// TODO:
// 2000
// CMC
// HCT
// ITP
// EJz

// DeltaDistance computes the Euclidean distance in the provided color space.
func DeltaDistance(reference, sample *Color, space *ColorSpace) float64 {
	ref := reference.Convert(space)
	s := sample.Convert(space)
	Δ0 := ref.Values[0] - s.Values[0]
	Δ1 := ref.Values[1] - s.Values[1]
	Δ2 := ref.Values[2] - s.Values[2]
	return math.Hypot(math.Hypot(Δ0, Δ1), Δ2)
}

// DeltaE76 computes the CIE 1976 color difference using the Euclidean distance
// in the [Lab] color space.
func DeltaE76(reference, sample *Color) float64 {
	return DeltaDistance(reference, sample, Lab)
}

// DeltaEOK computes the color difference using the Euclidean distance in the
// [Oklab] color space.
func DeltaEOK(reference, sample *Color) float64 {
	return DeltaDistance(reference, sample, Oklab)
}

// DeltaEOK2 computes the color difference using the Euclidean distance in the
// [Oklab] color space, with the a and b axes scaled by a factor of 2, for
// better uniformity.
func DeltaEOK2(reference, sample *Color) float64 {
	// See
	// https://github.com/w3c/csswg-drafts/issues/6642#issuecomment-945714988
	// and
	// https://github.com/color-js/color.js/blob/40e7a059c639bafde14504627e62791588c63100/src/deltaE/deltaEOK2.js
	// for background on the scaling.

	ref := reference.Convert(Oklab)
	s := sample.Convert(Oklab)

	Δ0 := ref.Values[0] - s.Values[0]
	Δ1 := 2 * (ref.Values[1] - s.Values[1])
	Δ2 := 2 * (ref.Values[2] - s.Values[2])
	return math.Hypot(math.Hypot(Δ0, Δ1), Δ2)
}
