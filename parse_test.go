package color

import "testing"

func BenchmarkParse(b *testing.B) {
	for range b.N {
		ParseColor(`color(xyz 0.1 0.2 0.3 / 0.4)`)
	}
}

func FuzzParse(f *testing.F) {
	f.Add(`color(xyz 0.1 0.2 0.3 / 0.4)`)
	f.Add(`color(xyz-d65 0.1 0.2 0.3 / 0.4)`)
	f.Add(`color(--oklab 0.1 0.2 0.3 / 0.4)`)
	f.Add(`color(oklab 0.1 0.2 0.3 / 0.4)`)
	f.Add(`color(oklab 0.1 0.2 0.3 / 40%)`)
	f.Add(`color(oklab 0.1 0.2 0.3)`)
	f.Add(`color(oklab 10% 0.2 0.3)`)

	f.Fuzz(func(t *testing.T, s string) {
		ParseColor(s)
	})
}
