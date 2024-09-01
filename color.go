// See https://github.com/w3c/csswg-drafts/issues/6618 for a discussion on the
// specific white points used in the CSS Color Module Level 4 spec. tl;dr: there
// are a dozen different definitions of D65, we have to pick one. Same for D50.

// Useful links:
// - https://facelessuser.github.io/coloraide/
// - https://colorjs.io/
// - https://www.w3.org/TR/css-color-4

// Package color provides types and functions for working with colors.
//
// # Color
//
// The core type is [Color], a color space-aware type for representing colors in
// any 3-axis coordinate system. For example, one color may be represented in a
// cartesian sRGB color space, while another is represented in the cylindrical
// Oklch. Colors can be instantiated directly or be created with the [Make]
// helper. Additionally, the CSS 'color()' syntax is supported and can be parsed
// by [Parse]. The package provides numerous color spaces, as variables of
// type [Space].
//
// Colors can be freely converted between any two color spaces, with optional
// gamut mapping. To convert a Color without gamut mapping, use [Color.Convert].
// If the destination color space is narrower than the source space (for
// example, when converting from Display P3 to sRGB), the new Color might not be
// in gamut and have values outside the expected range. Package color does not
// apply gamut mapping automatically for two reasons:
//
//  1. No gamut mapping approach is perfect and we leave the choice of which
//     algorithm to use up to the user.
//
//  2. Automatic gamut mapping would prevent colors from being roundtrippable.
//     Without gamut mapping, it is possible to go from a wider to a narrower and
//     back to a wider color space without losing information (except for the
//     introduction of rounding errors).
//
// A decent gamut mapping operation with a relative colorimetric intent is
// implemented by [GamutMapCSS]. This is the same algorithm that is used by
// browsers that implement the [CSS Color Module Level 4] standard.
//
// Whether a color is in gamut for its color space can be checked with
// [Color.InGamut]. There is also [Color.InGamutOf] that checks if a color, when
// converted to a target color space, will be in gamut.
//
// The following example shows the creation of a very saturated pink (more
// saturated than can be displayed by non-wide-gamut displays) using Oklch and
// its conversion to sRGB with and without gamut mapping.
//
//	veryPink := Make(Oklch, 0.65, 0.29, 0, 1)
//	veryPinkSRGB := veryPink.Convert(SRGB)
//	pinkSRGB := GamutMapCSS(&veryPink, SRGB)
//
//	fmt.Println(veryPink, veryPink.InGamut())
//	fmt.Println(veryPinkSRGB, veryPinkSRGB.InGamut())
//	fmt.Println(pinkSRGB, pinkSRGB.InGamut())
//
// Output:
//
//	color(--oklch 0.650000 0.290000 0.000000) true
//	color(srgb 1.040595 -0.191616 0.533106) false
//	color(srgb 1.000000 0.000000 0.533824) true
//
// This is everything that is needed to create and use colors. The following
// sections describe more advanced but optional functionality.
//
// # Color spaces
//
// Without an associated color space, a triple like (0.6, 0.2, 0) conveys no
// meaning. It might represent a brown color in sRGB (when representing values
// as floating point values in the range [0, 1]), or a soft pink in Oklch.
//
// Package color supports numerous color spaces and allows users to define their
// own. They are represented by the [Space] type. Color spaces include meta
// data such as an ID (for the use with [Parse]) and name as well as a
// [white point]. However, their main functionality is provided by having a
// "base" color space to and from which a color space can be converted. The base
// space might be closely related, such as sRGB being based on linear sRGB (as
// the former is just the latter with a transfer function applied). It might
// also be a more fundamental color space, often XYZ with the appropriate white
// point. Together, all color spaces form a tree, with XYZ D65 as the root. This
// tree allows converting between any two color spaces by finding a common
// ancestor (the conversion space) and applying a series of conversions from the
// source space to the conversion space and from the conversion space to the
// target space.
//
// # Chromatic adaptation
//
// Different color spaces may have different white points, requiring [chromatic
// adaptation] when converting between them. This is automatically handled
// during color space conversion. Every color space's chain of base spaces will
// eventually lead to an XYZ color space with the appropriate white point. When
// converting between two XYZ color spaces with different white points, we use
// the Bradford method to adapt between them.
//
// It is also possible to manually apply chromatic adaptation using the [CAT]
// type. This package provides the [Bradford] and [CAT16] transformations, but
// any transformation that is based on using a 3x3 matrix to map to a cone
// response and another 3x3 matrix to map from a cone response can be
// implemented using [CAT].
//
// Furthermore, we provide definitions for numerous white points, expressed as
// xy chromaticities using the [Chromaticity] type. Additionally, CIE daylight
// illuminants can be created with [MakeCIEDaylightIlluminant].
//
// New XYZ spaces with custom white points can be created with [NewXYZSpace].
//
// # Other functionality
//
// The perceptual difference between two colors can be computed by functions
// whose names start with "Delta", such as [DeltaEOK2]. Different functions have
// different tradeoffs.
//
// The contrast of two colors can be computed by functions whose names start
// with "Contrast", such as [ContrastWeber]. Different functions have different
// tradeoffs.
//
// [Step] creates color gradients by linearly interpolating between two colors
// in a color space of your choice.
//
// [CSS Color Module Level 4]: https://www.w3.org/TR/css-color-4/
// [white point]: https://en.wikipedia.org/wiki/White_point
// [chromatic adaptation]: https://en.wikipedia.org/wiki/Chromatic_adaptation
package color

// TODO(dh): https://github.com/WICG/color-api/issues/30

// a98rgb-linear.js
// a98rgb.js
// acescc.js
// acescg.js
// cam16.js
// hct.js
// hpluv.js
// hsl.js
// hsluv.js
// hsv.js
// hwb.js
// ictcp.js
// jzazbz.js
// jzczhz.js
// lab-d65.js
// lchuv.js
// luv.js
// okhsl.js
// okhsv.js
// oklrab.js
// oklrch.js
// rec2020-linear.js
// rec2020.js
// rec2100-hlg.js
// rec2100-pq.js
// xyz-abs-d65.js

import (
	"fmt"
	"iter"
)

// Make is a convenience function for initializing colors.
func Make(space *Space, p1, p2, p3, alpha float64) Color {
	if alpha < 0 {
		alpha = 0
	} else if alpha > 1 {
		alpha = 1
	}
	return Color{
		Values: [3]float64{p1, p2, p3},
		Space:  space,
		Alpha:  alpha,
	}
}

func lerp(x, y float64, a float64) float64 {
	return x*(1.0-a) + y*a
}

// Step computes num colors that lie between c1 and c2, interpolating in the in
// color space and returning them in the out color space, without applying any
// gamut mapping.
func Step(c1, c2 *Color, in, out *Space, num int) iter.Seq[Color] {
	return func(yield func(Color) bool) {
		c1in := c1.Convert(in)
		c2in := c2.Convert(in)

		for i := range num {
			t := float64(i+1) / float64(num)
			c := Make(
				in,
				lerp(c1in.Values[0], c2in.Values[0], t),
				lerp(c1in.Values[1], c2in.Values[1], t),
				lerp(c1in.Values[2], c2in.Values[2], t),
				lerp(c1in.Alpha, c2in.Alpha, t),
			)
			cout := c.Convert(out)
			if !yield(cout) {
				return
			}
		}
	}
}

// Chromaticity describes a color's chromaticity in the CIE 1931 xy color space.
type Chromaticity struct {
	X float64
	Y float64
}

// XYZ converts the xy chromaticity to the X, Y, and Z tristimulus values, with
// Y = 1.
func (chr *Chromaticity) XYZ() [3]float64 {
	return [3]float64{
		chr.X / chr.Y,
		1,
		(1 - chr.X - chr.Y) / chr.Y,
	}
}

// Color represents a color with 3 coordinates in some color space. The meaning
// of the values depends on the color space.
//
// The values of a color may be out of gamut for the color space. This is
// allowed so that conversions between color spaces do not lose any information,
// even if the destination space is smaller than the source space. The package
// provides functions for explicit gamut mapping.
//
// For convenience, colors include an alpha channel, commonly used for opacity
// or coverage. The alpha value doesn't affect operations such as color space
// conversions, gamut mapping, or distance metrics and will simply be preserved.
// [Step], however, will interpolate between the start and end alpha values.
type Color struct {
	Values [3]float64
	Space  *Space
	Alpha  float64
}

func (c Color) String() string {
	var isCSS bool
	switch c.Space.ID {
	case "srgb", "srgb-linear", "display-p3", "a98-rgb", "prophoto-rgb",
		"rec2020", "xyz-d50", "xyz-d65":
		isCSS = true
	}

	id := c.Space.ID
	if !isCSS {
		id = "--" + id
	}

	if c.Alpha != 1 {
		return fmt.Sprintf("color(%s %f %f %f / %f)",
			id, c.Values[0], c.Values[1], c.Values[2], c.Alpha)
	} else {
		return fmt.Sprintf("color(%s %f %f %f)",
			id, c.Values[0], c.Values[1], c.Values[2])
	}
}

// Convert converts c from its current color space to a different color space.
// It does not apply any gamut mapping.
func (c *Color) Convert(space *Space) Color {
	if c.Space == space {
		return *c
	}

	return Color{
		Values: c.Space.Convert(space, c.Values),
		Space:  space,
		Alpha:  c.Alpha,
	}
}

// InGamut reports whether c's values are in gamut of its color space.
func (c *Color) InGamut() bool {
	return c.Space.InGamut(c.Values)
}

// InGamutOf reports whether c, when converted to space, is in gamut.
func (c *Color) InGamutOf(space *Space) bool {
	cc := c.Convert(space)
	return cc.InGamut()
}

// GamutMapCSS uses the [CSS gamut mapping algorithm] to map individual colors
// to a destination color space. It implements a relative colorimetric intent.
// That is, colors that are already inside the target gamut are unchanged. This
// is intended for mapping individual colors, not for mapping images.
//
// For some limitations of this algorithm, see [1] and [2].
//
// [CSS gamut mapping algorithm]: https://www.w3.org/TR/css-color-4/#css-gamut-mapping
// [1]: https://github.com/w3c/csswg-drafts/issues/7071
// [2]: https://github.com/w3c/csswg-drafts/issues/9449
func GamutMapCSS(c *Color, to *Space) Color {
	// 1. if destination has no gamut limits (XYZ-D65, XYZ-D50, Lab, LCH,
	// Oklab, Oklch) convert origin to destination and return it as the
	// gamut mapped color
	if to.Coords[0].Range == infty &&
		to.Coords[1].Range == infty &&
		to.Coords[2].Range == infty {
		return c.Convert(to)
	}

	cOklch := c.Convert(Oklch)
	if cOklch.Values[0] >= 1 {
		out := Make(Oklab, 1, 0, 0, c.Alpha)
		return out.Convert(to)
	}
	if cOklch.Values[0] <= 0 {
		out := Make(Oklab, 0, 0, 0, c.Alpha)
		return out.Convert(to)
	}

	if out := cOklch.Convert(to); out.InGamut() {
		return out
	}

	// The just noticeable difference between two colors in Oklch
	const jnd = 0.02
	const ϵ = 0.0001

	clip := func(cc *Color) Color {
		clamp := func(f, low, high float64) float64 {
			if f < low {
				return low
			}
			if f > high {
				return high
			}
			return f
		}
		ccc := cc.Convert(to)
		ccc.Values[0] = clamp(ccc.Values[0], ccc.Space.Coords[0].Range[0], ccc.Space.Coords[0].Range[1])
		ccc.Values[1] = clamp(ccc.Values[1], ccc.Space.Coords[1].Range[0], ccc.Space.Coords[1].Range[1])
		ccc.Values[2] = clamp(ccc.Values[2], ccc.Space.Coords[2].Range[0], ccc.Space.Coords[2].Range[1])
		return ccc
	}

	current := cOklch
	clipped := clip(&current)
	e := DeltaEOK(&clipped, &current)
	if e < jnd {
		return clipped
	}
	min := 0.0
	max := cOklch.Values[1]
	minInGamut := true
	for max-min > ϵ {
		chroma := (min + max) / 2
		current.Values[1] = chroma
		if minInGamut && current.InGamutOf(to) {
			min = chroma
			continue
		} else if !current.InGamutOf(to) {
			clipped = clip(&current)
			e = DeltaEOK(&clipped, &current)
			if e < jnd {
				if jnd-e < ϵ {
					return clipped
				} else {
					minInGamut = false
					min = chroma
				}
			} else {
				max = chroma
				continue
			}
		}
	}
	return clipped
}

// Coordinate is metadata describing a coordinate of a color space.
type Coordinate struct {
	// Name is the human readable name of the coordinate.
	Name string
	// Range describes the range of values that are in gamut. For some
	// coordinates in some color spaces, this will be [-∞, ∞].
	Range [2]float64
	// Range describes the values that map to 0% and 100%. If not set, defaults
	// to Range.
	RefRange [2]float64
	// IsAngle is true for coordinates that represent angles, such as color hue.
	IsAngle bool
}

var RGBCoordinates = [3]Coordinate{
	{Name: "Red", Range: [2]float64{0, 1}},
	{Name: "Green", Range: [2]float64{0, 1}},
	{Name: "Blue", Range: [2]float64{0, 1}},
}
