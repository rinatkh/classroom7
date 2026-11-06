package simplesort

import (
	"fmt"
	"github.com/rinatkh/classroom7/internal/testutil"
	"slices"
	"testing"
)

func TestSortContracts(t *testing.T) {
	for _, a := range []struct {
		name  string
		run   func([]int) Report
		trace func([]int) [][]int
	}{
		{"bubble", BubbleSort, BubbleTrace}, {"selection", SelectionSort, SelectionTrace}, {"insertion", InsertionSort, InsertionTrace},
	} {
		t.Run(a.name, func(t *testing.T) {
			for _, tc := range testutil.Cases() {
				t.Run(tc.Name, func(t *testing.T) {
					input := slices.Clone(tc.Values)
					saved := slices.Clone(tc.Values)
					want := slices.Clone(tc.Values)
					slices.Sort(want)
					a.run(input)
					if !slices.Equal(input, want) {
						t.Fatalf("got %v want %v", input, want)
					}
					t.Run("trace", func(t *testing.T) {
						input = slices.Clone(tc.Values)
						states := a.trace(input)
						if len(states) == 0 || !slices.Equal(states[0], saved) || !slices.Equal(states[len(states)-1], want) {
							t.Fatalf("bad trace %v", states)
						}
						if len(input) > 0 {
							first := slices.Clone(states[0])
							input[0]++
							if !slices.Equal(states[0], first) {
								t.Fatal("trace aliases input")
							}
							if len(states) > 1 {
								last := slices.Clone(states[len(states)-1])
								states[0][0]++
								if !slices.Equal(states[len(states)-1], last) {
									t.Fatal("snapshots alias")
								}
							}
						}
					})
					if !slices.Equal(tc.Values, saved) {
						t.Fatal("fixture changed")
					}
				})
			}
		})
	}
}
func TestSelectionComparisonsSorted(t *testing.T) {
	for n := 2; n <= 12; n++ {
		t.Run(fmt.Sprintf("n=%d", n), func(t *testing.T) {
			v := make([]int, n)
			for i := range v {
				v[i] = i
			}
			got := SelectionSort(v)
			if got.Comparisons != n*(n-1)/2 || got.Writes != 0 {
				t.Fatalf("%+v", got)
			}
		})
	}
}
