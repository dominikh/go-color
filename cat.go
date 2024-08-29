package color

import "fmt"

var (
	Bradford = &CAT{
		ToCone: [3][3]float64{
			{+0.8951, +0.2664, -0.1614},
			{-0.7502, +1.7135, +0.0367},
			{+0.0389, -0.0685, +1.0296},
		},
		FromCone: [3][3]float64{
			{0.9869929054667121, -0.14705425642099013, 0.15996265166373122},
			{0.4323052697233945, 0.5183602715367774, 0.049291228212855594},
			{-0.00852866457517732, 0.04004282165408486, 0.96848669578755},
		},
	}

	CAT16 = &CAT{
		ToCone: [3][3]float64{
			{0.401288, 0.650173, -0.051461},
			{-0.250268, 1.204414, 0.045854},
			{-0.002079, 0.048952, 0.953127},
		},
		FromCone: [3][3]float64{
			{1.862067855087233, -1.0112546305316845, 0.14918677544445172},
			{0.3875265432361372, 0.6214474419314753, -0.008973985167612521},
			{-0.01584149884933386, -0.03412293802851557, 1.0499644368778496},
		},
	}
)

var (
	// Standard illuminants for the CIE 1931 standard observer, from tables T.3,
	// T.8, T.8.2, and T.9 in CIE 15:2004.
	WhitesCIE2004TwoDegA      = &Chromaticity{0.44758, 040745}
	WhitesCIE2004TwoDegC      = &Chromaticity{0.31006, 0.31616}
	WhitesCIE2004TwoDegD50    = &Chromaticity{0.34567, 0.35851}
	WhitesCIE2004TwoDegD55    = &Chromaticity{0.33243, 0.34744}
	WhitesCIE2004TwoDegD65    = &Chromaticity{0.31272, 0.32903}
	WhitesCIE2004TwoDegD75    = &Chromaticity{0.29903, 0.31488}
	WhitesCIE2004TwoDegFL1    = &Chromaticity{0.3131, 0.3371}
	WhitesCIE2004TwoDegFL2    = &Chromaticity{0.3721, 0.3751}
	WhitesCIE2004TwoDegFL3    = &Chromaticity{0.4091, 0.3941}
	WhitesCIE2004TwoDegFL3_1  = &Chromaticity{0.4407, 0.4033}
	WhitesCIE2004TwoDegFL3_2  = &Chromaticity{0.3808, 0.3734}
	WhitesCIE2004TwoDegFL3_3  = &Chromaticity{0.3153, 0.3439}
	WhitesCIE2004TwoDegFL3_4  = &Chromaticity{0.4429, 0.4043}
	WhitesCIE2004TwoDegFL3_5  = &Chromaticity{0.3749, 0.3672}
	WhitesCIE2004TwoDegFL3_6  = &Chromaticity{0.3488, 0.36}
	WhitesCIE2004TwoDegFL3_7  = &Chromaticity{0.4384, 0.4045}
	WhitesCIE2004TwoDegFL3_8  = &Chromaticity{0.382, 0.3832}
	WhitesCIE2004TwoDegFL3_9  = &Chromaticity{0.3499, 0.3591}
	WhitesCIE2004TwoDegFL3_10 = &Chromaticity{0.3455, 0.356}
	WhitesCIE2004TwoDegFL3_11 = &Chromaticity{0.3245, 0.3434}
	WhitesCIE2004TwoDegFL3_12 = &Chromaticity{0.4377, 0.4037}
	WhitesCIE2004TwoDegFL3_13 = &Chromaticity{0.383, 0.3724}
	WhitesCIE2004TwoDegFL3_14 = &Chromaticity{0.3447, 0.3609}
	WhitesCIE2004TwoDegFL3_15 = &Chromaticity{0.3127, 0.3288}
	WhitesCIE2004TwoDegFL4    = &Chromaticity{0.4402, 0.4031}
	WhitesCIE2004TwoDegFL5    = &Chromaticity{0.3138, 0.3452}
	WhitesCIE2004TwoDegFL6    = &Chromaticity{0.3779, 0.3882}
	WhitesCIE2004TwoDegFL7    = &Chromaticity{0.3129, 0.3292}
	WhitesCIE2004TwoDegFL8    = &Chromaticity{0.3458, 0.3586}
	WhitesCIE2004TwoDegFL9    = &Chromaticity{0.3741, 0.3727}
	WhitesCIE2004TwoDegFL10   = &Chromaticity{0.3458, 0.3588}
	WhitesCIE2004TwoDegFL11   = &Chromaticity{0.3805, 0.3769}
	WhitesCIE2004TwoDegFL12   = &Chromaticity{0.4370, 0.4042}
	WhitesCIE2004TwoDegHP1    = &Chromaticity{0.533, 0.415}
	WhitesCIE2004TwoDegHP2    = &Chromaticity{0.4778, 0.4158}
	WhitesCIE2004TwoDegHP3    = &Chromaticity{0.4302, 0.4075}
	WhitesCIE2004TwoDegHP4    = &Chromaticity{0.3812, 0.3797}
	WhitesCIE2004TwoDegHP5    = &Chromaticity{0.3776, 0.3713}

	// Standard illuminants for the CIE 1964 standard observer, from table T.3
	// in CIE 15:2004.
	WhitesCIE2004TenDegA   = &Chromaticity{0.45117, 0.40594}
	WhitesCIE2004TenDegC   = &Chromaticity{0.31039, 0.31905}
	WhitesCIE2004TenDegD50 = &Chromaticity{0.34773, 0.35952}
	WhitesCIE2004TenDegD55 = &Chromaticity{0.33412, 0.34877}
	WhitesCIE2004TenDegD65 = &Chromaticity{0.31381, 0.33098}
	WhitesCIE2004TenDegD75 = &Chromaticity{0.29968, 0.31740}

	// The D50 white point as defined in [CSS Color Module Level 4]. This
	// corresponds to [WhitesCIE2004TwoDegD50] but rounded to 4 digits.
	//
	// [CSS Color Module Level 4]: https://www.w3.org/TR/css-color-4/
	WhitesCSSD50 = &Chromaticity{0.3457, 0.3585}

	// The D65 white point as specified by sRGB. This corresponds to
	// [WhitesCIE2004TwoDegD65] but rounded to 4 digits.
	WhitesSRGBD65 = &Chromaticity{0.3127, 0.3290}
)

// MakeCIEDaylightIlluminant computes a daylight illuminant at a nominal
// correlated color temperature. The illuminant's correlated color temperature
// will be approximately equal to the nominal value, but not exactly so.
//
// Note that due to pecularities in rounding in the CIE standards, values
// returned by this function will not exactly match predefined daylight
// illuminants.
//
// The color temperature, specified in Kelvin, must be between 4000 K and
// 25,000 K.
func MakeCIEDaylightIlluminant(temp float64) Chromaticity {
	switch {
	case temp < 4000, temp > 25_000:
		fallthrough
	default:
		// function not defined for this range
		panic(fmt.Sprintf("color temperature %v is not in range [4000, 25000]", temp))
	case temp <= 7000:
		// Formula taken from CIE 15:2004, page 3, equations 3.2 and 3.3.
		x := (-4.6070e9)/(temp*temp*temp) + 2.9678e6/(temp*temp) + 0.09911e3/temp + 0.244063
		y := -3*x*x + 2.870*x - 0.275
		return Chromaticity{x, y}
	case temp <= 25_000:
		// Formula taken from CIE 15:2004, pages 3-4, equations 3.2 and 3.4.
		x := (-2.0064e9)/(temp*temp*temp) + 1.9018e6/(temp*temp) + 0.24748e3/temp + 0.237040
		y := -3*x*x + 2.870*x - 0.275
		return Chromaticity{x, y}
	}
}

// CAT represents a chromatic adaptation transform. It consists of two matrices,
// one for converting from XYZ to cone responses and one for converting from
// cone responses back to XYZ.
//
// Given a CAT, colors can be adapted between any two white points, either by
// using [CAT.Adapt] for one-offs, or by combining [CAT.Matrix] and [Adapt],
// which allows reusing matrices computed for pairs of white points.
type CAT struct {
	ToCone   [3][3]float64
	FromCone [3][3]float64
}

func (cat *CAT) Adapt(xyz *[3]float64, src, dst *Chromaticity) [3]float64 {
	m := cat.Matrix(src, dst)
	return Adapt(xyz, &m)
}

func (cat *CAT) Matrix(src, dst *Chromaticity) [3][3]float64 {
	ws := src.XYZ()
	wd := dst.XYZ()

	coneS := mulVecMat(&ws, &cat.ToCone)
	coneD := mulVecMat(&wd, &cat.ToCone)

	ρS := coneS[0]
	γS := coneS[1]
	βS := coneS[2]

	ρD := coneD[0]
	γD := coneD[1]
	βD := coneD[2]

	a := cat.FromCone[0][0]
	b := cat.FromCone[0][1]
	c := cat.FromCone[0][2]
	d := cat.FromCone[1][0]
	e := cat.FromCone[1][1]
	f := cat.FromCone[1][2]
	g := cat.FromCone[2][0]
	h := cat.FromCone[2][1]
	i := cat.FromCone[2][2]

	rρ := ρD / ρS
	rγ := γD / γS
	rβ := βD / βS
	m_ := [3][3]float64{
		{a * rρ, b * rγ, c * rβ},
		{d * rρ, e * rγ, f * rβ},
		{g * rρ, h * rγ, i * rβ},
	}
	return mulMatMat(&m_, &cat.ToCone)
}

func Adapt(xyz *[3]float64, m *[3][3]float64) [3]float64 {
	return mulVecMat(xyz, m)
}
