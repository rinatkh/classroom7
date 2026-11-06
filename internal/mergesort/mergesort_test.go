package mergesort

import (
	"cmp"
	"github.com/rinatkh/classroom7/internal/testutil"
	"slices"
	"strings"
	"testing"
)

func TestMergeContracts(t *testing.T) {
	for _, tc := range testutil.Cases() {
		t.Run(tc.Name, func(t *testing.T) {
			want := slices.Clone(tc.Values)
			slices.Sort(want)
			for _, mode := range []string{"int", "func", "trace"} {
				t.Run(mode, func(t *testing.T) {
					work := slices.Clone(tc.Values)
					switch mode {
					case "int":
						Sort(work)
					case "func":
						SortFunc(work, cmp.Compare[int])
					default:
						steps := SortWithTrace(work)
						if len(work) > 1 && len(steps) != len(work)-1 {
							t.Fatal("each merge expected")
						}
						for _, s := range steps {
							if !(0 <= s.Left && s.Left < s.Middle && s.Middle < s.Right && s.Right <= len(work)) || !slices.IsSorted(s.Values[s.Left:s.Right]) {
								t.Fatalf("bad step %+v", s)
							}
						}
						if len(steps) > 0 {
							saved := slices.Clone(steps[0].Values)
							work[0]++
							if !slices.Equal(saved, steps[0].Values) {
								t.Fatal("alias")
							}
							work[0]--
						}
					}
					if !slices.Equal(work, want) {
						t.Fatalf("got %v want %v", work, want)
					}
				})
			}
		})
	}
}
func TestMergeStability(t *testing.T) {
	type record struct{ Key, Order int }
	for _, tc := range testutil.Cases() {
		t.Run(tc.Name, func(t *testing.T) {
			values := make([]record, len(tc.Values))
			for i, k := range tc.Values {
				values[i] = record{k, i}
			}
			// Order обозначает исходную позицию, сравниваем ТОЛЬКО Key.
			SortFunc(values, func(a, b record) int { return cmp.Compare(a.Key, b.Key) })
			for i := 1; i < len(values); i++ {
				if values[i-1].Key > values[i].Key || (values[i-1].Key == values[i].Key && values[i-1].Order > values[i].Order) {
					t.Fatalf("unstable %v", values)
				}
			}
		})
	}
}
func TestExample(t *testing.T) {
	if !strings.Contains(Example(), "Merge Sort: [1 2 4 5]") {
		t.Fatal(Example())
	}
}
