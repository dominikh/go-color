package color

import (
	"regexp"
	"strconv"
)

var reColor = regexp.MustCompile(`^color\(` +
	`([a-zA-Z0-9-]+) ` +
	`((?:[+-]?\d+|[+-]?\d*\.\d+(?:[eE][+-]?\d+)?)%?) ` +
	`((?:[+-]?\d+|[+-]?\d*\.\d+(?:[eE][+-]?\d+)?)%?) ` +
	`((?:[+-]?\d+|[+-]?\d*\.\d+(?:[eE][+-]?\d+)?)%?)` +
	`(?: / ((?:[+-]?\d+|[+-]?\d*\.\d+(?:[eE][+-]?\d+)?)%?))?\);?$`)

// Parse parses colors in the CSS 'color()' format. The double dash for
// non-standard color spaces is optional.
func Parse(s string) (Color, bool) {
	m := reColor.FindStringSubmatch(s)
	if m == nil {
		return Color{}, false
	}

	space := m[1]
	x := m[2]
	y := m[3]
	z := m[4]
	a := m[5]

	if space == "xyz" {
		space = "xyz-d65"
	}
	cs, ok := LookupSpace(space)
	if !ok {
		return Make(SRGB, 0, 0, 0, 1), false
	}

	var values [4]float64
	parseValue := func(idx int, s string) bool {
		if idx == 3 && len(s) == 0 {
			values[3] = 1
			return true
		}

		if s[len(s)-1] == '%' {
			f, err := strconv.ParseFloat(s[:len(s)-1], 64)
			if err != nil {
				// Even inputs that pass the regex can get here, e.g. because of
				// absurdly large values.
				return false
			}
			if f < 0 {
				f = 0
			}
			if f > 100 {
				f = 100
			}
			f /= 100
			if idx == 3 {
				values[3] = f
			} else {
				rng := cs.Coords[idx].RefRange
				values[idx] = lerp(rng[0], rng[1], f)
			}
		} else {
			f, err := strconv.ParseFloat(s, 64)
			if err != nil {
				// Even inputs that pass the regex can get here, e.g. because of
				// absurdly large values.
				return false
			}
			if idx == 3 {
				if f < 0 {
					f = 0
				}
				if f > 1 {
					f = 1
				}
			}
			values[idx] = f
		}
		return true
	}

	parseValue(0, x)
	parseValue(1, y)
	parseValue(2, z)
	parseValue(3, a)

	return Make(cs, values[0], values[1], values[2], values[3]), true
}
