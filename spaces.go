package color

import (
	"fmt"
	"math"
	"slices"
	"strings"
	"sync"
)

func init() {
	RegisterColorSpace(XYZ_D50)
	RegisterColorSpace(XYZ_D65)
	RegisterColorSpace(LinearDisplayP3)
	RegisterColorSpace(DisplayP3)
	RegisterColorSpace(LinearSRGB)
	RegisterColorSpace(SRGB)
	RegisterColorSpace(Oklab)
	RegisterColorSpace(Oklch)
	RegisterColorSpace(ProPhoto)
	RegisterColorSpace(LinearProPhoto)
	RegisterColorSpace(Lab)
	RegisterColorSpace(LCh)
}

var (
	colorSpacesMu sync.RWMutex
	colorSpaces   = map[string]*ColorSpace{}

	infty = [2]float64{math.Inf(-1), math.Inf(1)}
	norm  = [2]float64{0, 1}
)

// LookupColorSpace looks up a registered (see [RegisterColorSpace]) color space by ID.
func LookupColorSpace(id string) (*ColorSpace, bool) {
	id = strings.TrimPrefix(id, "--")
	colorSpacesMu.RLock()
	defer colorSpacesMu.RUnlock()
	cs, ok := colorSpaces[id]
	return cs, ok
}

// RegisterColorSpace registers a color space. This allows it to be referenced
// by ID in 'color()' expressions as parsed by [ParseColor] and looked up by
// [LookupColorSpace].
//
// All color spaces provided by this package are automatically registered.
func RegisterColorSpace(cs *ColorSpace) {
	colorSpacesMu.Lock()
	defer colorSpacesMu.Unlock()
	registerColorSpace(cs)
}

func registerColorSpace(cs *ColorSpace) {
	if _, ok := colorSpaces[cs.ID]; ok {
		// Trying to register the same color space ID more than once might point
		// to a mistake, but it might also be the result of us registering base
		// spaces, so we can't panic here.
		return
	}
	colorSpaces[cs.ID] = cs
	if cs.Base != nil {
		if _, ok := colorSpaces[cs.Base.ID]; !ok {
			registerColorSpace(cs.Base)
		}
	}
}

// ColorSpace describes a color space, such as sRGB or HSV.
//
// Color spaces form a tree. Every space, except for [XYZ_D65], has a base space
// and can be converted to and from it. Every space can be converted to any
// other space by finding a common "connection space". Often, the common space
// is XYZ D65.
//
// Color spaces are bit-width-agnostic and values are often stored in a
// normalized form. For example, the R, G, and B coordinates in sRGB will be in
// the range [0, 1], not [0, 255].
//
// The white point of a color space is described by the White field. This field
// only serves as documentation and does not affect how conversions are carried
// out. That is, simply changing the white point of an existing color space will
// not have the intended effect.
//
// When creating new color spaces, you must call [ColorSpace.Init] once you're
// done.
type ColorSpace struct {
	ID       string
	Name     string
	White    *Chromaticity
	Base     *ColorSpace
	Coords   [3]Coordinate
	FromBase func(c *[3]float64) [3]float64
	ToBase   func(c *[3]float64) [3]float64

	path []*ColorSpace
}

func (cs *ColorSpace) Init() *ColorSpace {
	if cs.Coords == ([3]Coordinate{}) {
		cs.Coords = cs.Base.Coords
	}
	if cs.White == nil && cs.Base != nil {
		cs.White = cs.Base.White
	}

	for i := range cs.Coords {
		coord := &cs.Coords[i]
		if coord.RefRange == ([2]float64{}) {
			coord.RefRange = coord.Range
		}
	}

	// if cs.GamutSpace == nil {
	// 	var isPolar bool
	// 	for _, coord := range cs.Coords {
	// 		if coord.IsAngle {
	// 			isPolar = true
	// 			break
	// 		}
	// 	}
	// 	if isPolar {
	// 		cs.GamutSpace = cs.Base
	// 	} else {
	// 		cs.GamutSpace = cs
	// 	}
	// }

	orig := cs
	var out []*ColorSpace
	for cs != nil {
		out = append(out, cs)
		cs = cs.Base
	}
	slices.Reverse(out)
	orig.path = out
	return orig
}

func (cs *ColorSpace) InGamut(values [3]float64) bool {
	const ϵ = 0.000075
	// if cs.GamutSpace != cs {
	// 	values = cs.Convert(cs.GamutSpace, values)
	// 	return cs.GamutSpace.InGamut(values)
	// }

	for i, v := range values {
		meta := cs.Coords[i]
		if !meta.IsAngle {
			min := meta.Range[0]
			max := meta.Range[1]
			if !(v >= min-ϵ && v <= max+ϵ) {
				return false
			}
		}
	}
	return true
}

func (cs *ColorSpace) Convert(to *ColorSpace, coords [3]float64) [3]float64 {
	ourPath := cs.path
	theirPath := to.path

	// Determine the connection space by finding the lowest common ancestor of
	// the source and destination spaces in the color space tree.
	connIdx := -1
	for i := range min(len(ourPath), len(theirPath)) {
		if ourPath[i] == theirPath[i] {
			connIdx = i
		} else {
			break
		}
	}
	if connIdx == -1 {
		// Every space should be connectable through XYZ.
		panic(fmt.Sprintf("internal error: couldn't find connection space for %s and %s",
			cs.Name, to.Name))
	}

	// Convert from our space to the connection space
	for i := len(ourPath) - 1; i > connIdx; i-- {
		coords = ourPath[i].ToBase(&coords)
	}
	// Convert from connection space to destination space
	for i := connIdx + 1; i < len(theirPath); i++ {
		coords = theirPath[i].FromBase(&coords)
	}

	return coords
}

// NewXYZSpace returns a new CIE XYZ color space with the specified name, ID, and
// white point.
func NewXYZSpace(name, id string, white *Chromaticity) *ColorSpace {
	// OPT(dh): because all white point conversions go through D65, converting
	// between two non-D65 white points uses two instead of one matrix. For
	// example, we'd do D50->D65->D75, instead of the more direct D50->D75. This
	// is slower, and introduces more floating point error.
	//
	// In practice, most color spaces use D65 or D50, anyway.
	toD65 := Bradford.Matrix(white, XYZ_D65.White)
	fromD65 := Bradford.Matrix(XYZ_D65.White, white)
	return (&ColorSpace{
		ID:    id,
		Name:  name,
		White: white,
		Base:  XYZ_D65,
		FromBase: func(c *[3]float64) [3]float64 {
			return Adapt(c, &fromD65)
		},
		ToBase: func(c *[3]float64) [3]float64 {
			return Adapt(c, &toD65)
		},
	}).Init()
}

var XYZ_D50 = NewXYZSpace("XYZ D50", "xyz-d50", WhitesCSSD50)

var XYZ_D65 = (&ColorSpace{
	ID:   "xyz-d65",
	Name: "XYZ D65",
	Coords: [3]Coordinate{
		{Name: "X", Range: infty, RefRange: norm},
		{Name: "Y", Range: infty, RefRange: norm},
		{Name: "Z", Range: infty, RefRange: norm},
	},
	White: WhitesSRGBD65,
}).Init()

var LinearDisplayP3 = newRGBColorSpace(
	&rgbColorSpace{
		ID:   "display-p3-linear",
		Name: "Linear Display P3",
		Base: XYZ_D65,
		ToBase: [3][3]float64{
			{0.4865709486482162, 0.26566769316909306, 0.1982172852343625},
			{0.2289745640697488, 0.6917385218365064, 0.079286914093745},
			{0.0000000000000000, 0.04511338185890264, 1.043944368900976},
		},
		FromBase: [3][3]float64{
			{2.493496911941425, -0.9313836179191239, -0.40271078445071684},
			{-0.8294889695615747, 1.7626640603183463, 0.023624685841943577},
			{0.03584583024378447, -0.07617238926804182, 0.9568845240076872},
		},
	},
)

var DisplayP3 = (&ColorSpace{
	ID:   "display-p3",
	Name: "Display P3",
	Base: LinearDisplayP3,
	// Gamma encoding is the same as sRGB
	ToBase:   SRGB.ToBase,
	FromBase: SRGB.FromBase,
}).Init()

var LinearSRGB = newRGBColorSpace(
	&rgbColorSpace{
		ID:   "srgb-linear",
		Name: "Linear sRGB",
		Base: XYZ_D65,
		// sRGB Matrices taken from
		// https://github.com/w3c/csswg-drafts/pull/7320/commits/e835926c83e10342c0c43fd2f1ccbff1b35c3f07
		ToBase: [3][3]float64{
			{506752.0 / 1228815.0, 87881.0 / 245763.0, 12673.0 / 70218.0},
			{87098.0 / 409605.0, 175762.0 / 245763.0, 12673.0 / 175545.0},
			{7918.0 / 409605.0, 87881.0 / 737289.0, 1001167.0 / 1053270.0},
		},
		FromBase: [3][3]float64{
			{12831.0 / 3959.0, -329.0 / 214.0, -1974.0 / 3959.0},
			{-851781.0 / 878810.0, 1648619.0 / 878810.0, 36519.0 / 878810.0},
			{705.0 / 12673.0, -2585.0 / 12673.0, 705.0 / 667.0},
		},
	},
)

type rgbColorSpace struct {
	ID       string
	Name     string
	Base     *ColorSpace
	ToBase   [3][3]float64
	FromBase [3][3]float64
}

func newRGBColorSpace(space *rgbColorSpace) *ColorSpace {
	return (&ColorSpace{
		ID:     space.ID,
		Name:   space.Name,
		Coords: RGBCoordinates,
		Base:   space.Base,
		ToBase: func(c *[3]float64) [3]float64 {
			return mulVecMat(c, &space.ToBase)
		},
		FromBase: func(c *[3]float64) [3]float64 {
			return mulVecMat(c, &space.FromBase)
		},
	}).Init()
}

var SRGB = (&ColorSpace{
	ID:   "srgb",
	Name: "sRGB",
	Base: LinearSRGB,
	FromBase: func(c *[3]float64) [3]float64 {
		// TODO(dh): should this use the piecewise function, or a flat 2.2
		// gamma? See discussion in
		// https://gitlab.freedesktop.org/pq/color-and-hdr/-/issues/12

		f := func(ch float64) float64 {
			var sign float64
			if ch < 0 {
				sign = -1.0
			} else {
				sign = 1.0
			}
			abs := ch * sign

			if abs > 0.0031308 {
				return sign * (1.055*(math.Pow(abs, 1.0/2.4)) - 0.055)
			} else {
				return 12.92 * ch
			}
		}
		return [3]float64{f(c[0]), f(c[1]), f(c[2])}
	},
	ToBase: func(c *[3]float64) [3]float64 {
		// TODO(dh): same concern as FromBase
		f := func(ch float64) float64 {
			var sign float64
			if ch < 0 {
				sign = -1
			} else {
				sign = 1
			}
			abs := ch * sign
			if abs <= 0.04045 {
				return ch / 12.92
			} else {
				return sign * math.Pow((abs+0.055)/1.055, 2.4)
			}
		}
		return [3]float64{f(c[0]), f(c[1]), f(c[2])}
	},
}).Init()

// Matrices have been recalculated for consistent reference white;
// see https://github.com/w3c/csswg-drafts/issues/6642#issuecomment-943521484
var (
	oklabXyzToLms = [3][3]float64{
		{0.8190224379967030, 0.3619062600528904, -0.1288737815209879},
		{0.0329836539323885, 0.9292868615863434, 0.0361446663506424},
		{0.0481771893596242, 0.2642395317527308, 0.6335478284694309},
	}

	oklabLmsToLab = [3][3]float64{
		{0.2104542683093140, 0.7936177747023054, -0.0040720430116193},
		{1.9779985324311684, -2.4285922420485799, 0.4505937096174110},
		{0.0259040424655478, 0.7827717124575296, -0.8086757549230774},
	}

	oklabLabToLms = [3][3]float64{
		{1.0000000000000000, 0.3963377773761749, 0.2158037573099136},
		{1.0000000000000000, -0.1055613458156586, -0.0638541728258133},
		{1.0000000000000000, -0.0894841775298119, -1.2914855480194092},
	}

	oklabLmsToXyz = [3][3]float64{
		{1.2268798758459243, -0.5578149944602171, 0.2813910456659647},
		{-0.0405757452148008, 1.1122868032803170, -0.0717110580655164},
		{-0.0763729366746601, -0.4214933324022432, 1.5869240198367816},
	}
)

var Oklab = (&ColorSpace{
	ID:   "oklab",
	Name: "Oklab",
	Coords: [3]Coordinate{
		{Name: "Lightness", Range: infty, RefRange: norm},
		{Name: "a", Range: infty, RefRange: [2]float64{-0.4, 0.4}},
		{Name: "b", Range: infty, RefRange: [2]float64{-0.4, 0.4}},
	},
	Base: XYZ_D65,
	FromBase: func(c *[3]float64) [3]float64 {
		lms := mulVecMat(c, &oklabXyzToLms)

		lms_ := [3]float64{
			math.Cbrt(lms[0]),
			math.Cbrt(lms[1]),
			math.Cbrt(lms[2]),
		}
		lab := mulVecMat(&lms_, &oklabLmsToLab)
		return lab
	},
	ToBase: func(c *[3]float64) [3]float64 {
		lms := mulVecMat(c, &oklabLabToLms)
		lms_ := [3]float64{
			lms[0] * lms[0] * lms[0],
			lms[1] * lms[1] * lms[1],
			lms[2] * lms[2] * lms[2],
		}

		xyz := mulVecMat(&lms_, &oklabLmsToXyz)
		return xyz
	},
}).Init()

var Oklch = (&ColorSpace{
	ID:   "oklch",
	Name: "Oklch",
	Coords: [3]Coordinate{
		{Name: "Lightness", Range: infty, RefRange: norm},
		{Name: "Chroma", Range: infty, RefRange: [2]float64{0, 0.4}},
		{Name: "Hue", Range: infty, IsAngle: true, RefRange: [2]float64{0, 360}},
	},
	Base: Oklab,
	FromBase: func(c *[3]float64) [3]float64 {
		return labToLCH(c, 0.8/1e5)
	},
	ToBase: LCh.ToBase,
}).Init()

var Lab = (&ColorSpace{
	ID:   "lab",
	Name: "Lab",
	Coords: [3]Coordinate{
		{Name: "Lightness", Range: infty, RefRange: [2]float64{0, 100}},
		{Name: "a", Range: infty, RefRange: [2]float64{-125, 125}},
		{Name: "b", Range: infty, RefRange: [2]float64{-125, 125}},
	},
	Base: XYZ_D50,
	FromBase: func(c *[3]float64) [3]float64 {
		const (
			ϵ  = 216.0 / 24389.0
			ϵ3 = 24.0 / 116.0
			κ  = 24389.0 / 27.0
		)

		white := WhitesCSSD50.XYZ()
		xyz := *c
		xyz[0] /= white[0]
		xyz[1] /= white[1]
		xyz[2] /= white[2]

		f := func(x float64) float64 {
			if x > ϵ {
				return math.Cbrt(x)
			} else {
				return (κ*x + 16) / 116.0
			}
		}
		x_ := f(xyz[0])
		y_ := f(xyz[1])
		z_ := f(xyz[2])

		l := 116.0*y_ - 16
		a := 500.0 * (x_ - y_)
		b := 200.0 * (y_ - z_)

		return [3]float64{l, a, b}
	},
	ToBase: func(c *[3]float64) [3]float64 {
		const (
			ϵ  = 216.0 / 24389.0
			ϵ3 = 24.0 / 116.0
			κ  = 24389.0 / 27.0
		)

		l, a, b := c[0], c[1], c[2]
		f1 := (l + 16.0) / 116.0
		f0 := a/500.0 + f1
		f2 := f1 - b/200.0

		var x, y, z float64
		if f0 > ϵ3 {
			x = f0 * f0 * f0
		} else {
			x = (116*f0 - 16) / κ
		}
		if l > 8 {
			n := (l + 16.0) / 116.0
			y = n * n * n
		} else {
			y = l / κ
		}
		if f2 > ϵ3 {
			z = f2 * f2 * f2
		} else {
			z = (116.0*f2 - 16) / κ
		}

		white := WhitesCSSD50.XYZ()
		x /= white[0]
		y /= white[1]
		z /= white[2]

		return [3]float64{x, y, z}
	},
}).Init()

var LCh = (&ColorSpace{
	ID:   "lch",
	Name: "LCh",
	Coords: [3]Coordinate{
		{Name: "Lightness", Range: infty, RefRange: [2]float64{0, 100}},
		{Name: "Chroma", Range: infty, RefRange: [2]float64{0, 150}},
		{Name: "Hue", Range: infty, IsAngle: true, RefRange: [2]float64{0, 360}},
	},
	Base: Lab,
	FromBase: func(c *[3]float64) [3]float64 {
		return labToLCH(c, 250.0/1e5)
	},
	ToBase: func(cl *[3]float64) [3]float64 {
		// XXX handle achromatic h
		l, c, h := cl[0], cl[1], cl[2]
		if c < 0 {
			c = 0
		}
		a := c * math.Cos(h*math.Pi/180.0)
		b := c * math.Sin(h*math.Pi/180)
		return [3]float64{l, a, b}
	},
}).Init()

func labToLCH(lab *[3]float64, ϵ float64) [3]float64 {
	l, a, b := lab[0], lab[1], lab[2]
	achromatic := math.Abs(a) < ϵ && math.Abs(b) < ϵ
	var c, h float64
	if achromatic {
		c = 0
		// XXX color.js uses null for achromatic
		h = 0
	} else {
		c = math.Sqrt(a*a + b*b)
		h_ := math.Atan2(b, a) * 180 / math.Pi
		h = math.Mod(h_+360, 360)
	}
	return [3]float64{l, c, h}
}

var LinearProPhoto = newRGBColorSpace(
	&rgbColorSpace{
		ID:   "prophoto-rgb-linear",
		Name: "Linear ProPhoto",
		Base: XYZ_D50,
		ToBase: [3][3]float64{
			{0.79776664490064230, 0.13518129740053308, 0.03134773412839220},
			{0.28807482881940130, 0.71183523424187300, 0.00008993693872564},
			{0.00000000000000000, 0.00000000000000000, 0.82510460251046020},
		},
		FromBase: [3][3]float64{
			{1.34578688164715830, -0.25557208737979464, -0.05110186497554526},
			{-0.54463070512490190, 1.50824774284514680, 0.02052744743642139},
			{0.00000000000000000, 0.00000000000000000, 1.21196754563894520},
		},
	},
)

var ProPhoto = (&ColorSpace{
	ID:     "prophoto-rgb",
	Name:   "ProPhoto",
	Base:   LinearProPhoto,
	Coords: RGBCoordinates,
	ToBase: func(c *[3]float64) [3]float64 {
		f := func(v float64) float64 {
			if v < 16.0/512.0 {
				return v / 16.0
			} else {
				return math.Pow(v, 1.8)
			}
		}
		return [3]float64{
			f(c[0]),
			f(c[1]),
			f(c[2]),
		}
	},
	FromBase: func(c *[3]float64) [3]float64 {
		f := func(v float64) float64 {
			if v >= 1.0/512.0 {
				return math.Pow(v, (1.0 / 1.8))
			} else {
				return 16 * v
			}
		}
		return [3]float64{
			f(c[0]),
			f(c[1]),
			f(c[2]),
		}
	},
}).Init()

func mulVecMat(vec *[3]float64, m *[3][3]float64) [3]float64 {
	return [3]float64{
		m[0][0]*vec[0] + m[0][1]*vec[1] + m[0][2]*vec[2],
		m[1][0]*vec[0] + m[1][1]*vec[1] + m[1][2]*vec[2],
		m[2][0]*vec[0] + m[2][1]*vec[1] + m[2][2]*vec[2],
	}
}

func mulMatMat(m1, m2 *[3][3]float64) [3][3]float64 {
	return [3][3]float64{
		{
			m1[0][0]*m2[0][0] + m1[0][1]*m2[1][0] + m1[0][2]*m2[2][0],
			m1[0][0]*m2[0][1] + m1[0][1]*m2[1][1] + m1[0][2]*m2[2][1],
			m1[0][0]*m2[0][2] + m1[0][1]*m2[1][2] + m1[0][2]*m2[2][2],
		},
		{
			m1[1][0]*m2[0][0] + m1[1][1]*m2[1][0] + m1[1][2]*m2[2][0],
			m1[1][0]*m2[0][1] + m1[1][1]*m2[1][1] + m1[1][2]*m2[2][1],
			m1[1][0]*m2[0][2] + m1[1][1]*m2[1][2] + m1[1][2]*m2[2][2],
		},
		{
			m1[2][0]*m2[0][0] + m1[2][1]*m2[1][0] + m1[2][2]*m2[2][0],
			m1[2][0]*m2[0][1] + m1[2][1]*m2[1][1] + m1[2][2]*m2[2][1],
			m1[2][0]*m2[0][2] + m1[2][1]*m2[1][2] + m1[2][2]*m2[2][2],
		},
	}
}
