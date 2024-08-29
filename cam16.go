package color

// The viewing environment. For advice on choosing values, see "Usage Guidelines
// for CIECAM97s" (2000) by Moroney.
type _environment struct {
	White *Chromaticity
	// The average luminance of the environment in cd/m² (a.k.a. nits). Under a
	// "gray world" assumption this is 20% of the luminance of a white
	// reference.
	AdaptingLuminance float64
	// The relative luminance of the nearby background (out to 10°), relative to
	// Y=1 of the white point.
	BackgroundLuminance float64
	// A description of the peripheral area's luminance compared to that of the
	// scene. 0 ("dark") denotes a fully dark room, such as a movie cinema. 1
	// ("dim") denotes a dim room, such as an average media room, and 2
	// ("average") denotes a normally lit room. Other values between 0 and 2
	// will interpolate between these options.
	Surround float64
	// If set to true, the eyes are assumed to be fully adapted to the
	// illuminant. By default, the degree of discounting will be set based on
	// the other fields.
	Discounting bool
}
