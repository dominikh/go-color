package color

import (
	"slices"
	"testing"
)

func TestStep(t *testing.T) {
	t.Run("smooth", func(t *testing.T) {
		c1 := Make(LinearSRGB, 0, 0, 0, 1)
		c2 := Make(LinearSRGB, 1, 0, 0, 1)
		const N = 101
		want := make([]float64, N)
		for i := range want {
			want[i] = float64(i) / 100
		}
		got := slices.Collect(Step(&c1, &c2, LinearSRGB, LinearSRGB, N))
		if len(got) != N {
			t.Fatalf("got %d steps, want %d", len(got), N)
		}
		for i := range N {
			g := got[i].Values[0]
			w := want[i]
			if g != w {
				t.Fatalf("step %d: got value %g, want %g", i, g, w)
			}
		}
	})

	t.Run("anchored", func(t *testing.T) {
		c1 := Make(LinearSRGB, 0, 0, 0, 1)
		c2 := Make(LinearSRGB, 1, 0, 0, 1)

		for i := range 1000 {
			got := slices.Collect(Step(&c1, &c2, LinearSRGB, LinearSRGB, i+2))
			if got[0] != c1 {
				t.Fatalf("got first step %v, want %v", got[0], c1)
			}
			if got[len(got)-1] != c2 {
				t.Fatalf("got last step %v, want %v", got[len(got)-1], c2)
			}
		}
	})

}
