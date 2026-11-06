package search

import (
	"github.com/rinatkh/classroom7/internal/testutil"
	"slices"
	"testing"
)

func TestSearchContracts(t *testing.T) {
	for _, tc := range testutil.SearchCases() {
		t.Run(tc.Name, func(t *testing.T) {
			saved := slices.Clone(tc.Values)
			first := slices.Index(tc.Values, tc.Target)
			t.Run("linear", func(t *testing.T) {
				got := LinearSearch(tc.Values, tc.Target)
				wantSteps := len(tc.Values)
				if first >= 0 {
					wantSteps = first + 1
				}
				if got.Index != first || got.Found != (first >= 0) || got.Comparisons != wantSteps {
					t.Fatalf("linear %+v, first=%d steps=%d", got, first, wantSteps)
				}
			})
			t.Run("binary", func(t *testing.T) { checkSearch(t, tc.Values, tc.Target, BinarySearch(tc.Values, tc.Target)) })
			t.Run("binary-trace", func(t *testing.T) {
				got, steps := BinarySearchTrace(tc.Values, tc.Target)
				checkSearch(t, tc.Values, tc.Target, got)
				if got != BinarySearch(tc.Values, tc.Target) || got.Comparisons != len(steps) {
					t.Fatal("trace diverged")
				}
				previous := len(tc.Values) + 1
				for _, step := range steps {
					if step.Left < 0 || step.Right > len(tc.Values) || step.Middle < step.Left || step.Middle >= step.Right || step.Right-step.Left >= previous {
						t.Fatalf("bad interval %+v", step)
					}
					if step.Value != tc.Values[step.Middle] {
						t.Fatal("wrong middle")
					}
					previous = step.Right - step.Left
				}
			})
			t.Run("lower-bound", func(t *testing.T) {
				got := SearchInsert(tc.Values, tc.Target)
				want, _ := slices.BinarySearch(tc.Values, tc.Target)
				if got != want {
					t.Fatalf("lower=%d want=%d", got, want)
				}
				for _, v := range tc.Values[:got] {
					if v >= tc.Target {
						t.Fatal("left invariant")
					}
				}
				for _, v := range tc.Values[got:] {
					if v < tc.Target {
						t.Fatal("right invariant")
					}
				}
			})
			if !slices.Equal(saved, tc.Values) {
				t.Fatal("search changed input")
			}
		})
	}
}
func checkSearch(t *testing.T, v []int, target int, got Result) {
	t.Helper()
	first := slices.Index(v, target)
	if got.Found != (first >= 0) || (!got.Found && got.Index != -1) || (got.Found && (got.Index < 0 || got.Index >= len(v) || v[got.Index] != target)) {
		t.Fatalf("bad result %+v", got)
	}
	limit := 0
	for n := len(v); n > 0; n /= 2 {
		limit++
	}
	if got.Comparisons > limit {
		t.Fatalf("%d > logarithmic bound %d", got.Comparisons, limit)
	}
}
