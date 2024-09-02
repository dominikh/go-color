package color

import (
	"fmt"
	"testing"
)

func BenchmarkParse(b *testing.B) {
	for range b.N {
		Parse(`color(xyz 0.1 0.2 0.3 / 0.4)`)
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
		Parse(s)
	})
}

func ExampleParse() {
	c, ok := Parse("color(lab 0.4 30% 0.2 / 1)")
	fmt.Println(c, ok)
	// Output:
	// color(--lab 0.400000 -50.000000 0.200000) true
}
