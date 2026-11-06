package dataset

import (
	"fmt"
	"slices"
	"testing"
)

func TestGenerator(t *testing.T) {
	for _, n := range []int{-3, -1, 0, 1, 2, 3, 4, 17, 100, 1000, 10000} {
		t.Run(fmt.Sprint(n), func(t *testing.T) {
			a, b := GenerateRandomSlice(n), GenerateRandomSlice(n)
			if len(a) != max(0, n) || !slices.Equal(a, b) {
				t.Fatal("size/reproducibility")
			}
			for _, v := range a {
				if v < -1000000 || v > 1000000 {
					t.Fatal("range")
				}
			}
			if len(a) > 0 {
				saved := b[0]
				a[0]++
				if b[0] != saved {
					t.Fatal("alias")
				}
			}
			if !slices.Equal(b, GenerateRandomSliceWithSeed(n, Seed)) {
				t.Fatal("explicit seed")
			}
		})
	}
	if slices.Equal(GenerateRandomSliceWithSeed(100, 1), GenerateRandomSliceWithSeed(100, 2)) {
		t.Fatal("seed ignored")
	}
}
func TestScenarios(t *testing.T) {
	for _, name := range Names {
		for _, n := range []int{-1, 0, 1, 2, 3, 5, 10, 20, 21, 100, 1000} {
			t.Run(fmt.Sprintf("%s/%d", name, n), func(t *testing.T) {
				a, b := Scenario(n, name), Scenario(n, name)
				if len(a) != max(0, n) || !slices.Equal(a, b) {
					t.Fatal("scenario reproducibility")
				}
				switch name {
				case "sorted":
					if !slices.IsSorted(a) {
						t.Fatal("sorted")
					}
				case "reverse":
					for i := 1; i < len(a); i++ {
						if a[i-1] < a[i] {
							t.Fatal("reverse")
						}
					}
				case "equal":
					for _, v := range a {
						if v != 7 {
							t.Fatal("equal")
						}
					}
				case "duplicates":
					for _, v := range a {
						if v < -7 || v > 7 {
							t.Fatal("duplicates range")
						}
					}
				case "nearly":
					sorted := Scenario(n, "sorted")
					for i := 0; i+1 < len(a); i += 20 {
						a[i], a[i+1] = a[i+1], a[i]
					}
					if !slices.Equal(a, sorted) {
						t.Fatal("nearly swaps")
					}
				}
			})
		}
	}
}
func TestUnknownScenario(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("unknown scenario must panic")
		}
	}()
	Scenario(2, "unknown")
}
