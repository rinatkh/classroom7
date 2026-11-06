package quicksort

import (
	"fmt"
	"github.com/rinatkh/classroom7/internal/testutil"
	"slices"
	"testing"
)

func TestQuickContracts(t *testing.T) {
	for _, tc := range testutil.Cases() {
		t.Run(tc.Name, func(t *testing.T) {
			want := slices.Clone(tc.Values)
			slices.Sort(want)
			work := slices.Clone(tc.Values)
			stats := Sort(work)
			if !slices.Equal(work, want) {
				t.Fatalf("sort=%v want=%v", work, want)
			}
			t.Run("trace", func(t *testing.T) {
				work = slices.Clone(tc.Values)
				traced, steps := SortWithTrace(work)
				if !slices.Equal(work, want) || traced != stats {
					t.Fatal("trace and sort disagree")
				}
				if len(steps) > 0 && len(work) > 0 {
					first := slices.Clone(steps[0].Values)
					work[0]++
					if !slices.Equal(steps[0].Values, first) {
						t.Fatal("snapshot aliases input")
					}
				}
			})
		})
	}
}
func TestPartitionInvariant(t *testing.T) {
	for _, tc := range testutil.Cases() {
		if len(tc.Values) == 0 {
			continue
		}
		t.Run(tc.Name, func(t *testing.T) {
			// Защитные элементы вне диапазона не должны меняться.
			values := append([]int{12345}, tc.Values...)
			values = append(values, 67890)
			left, right := 1, len(values)-2
			pivot := values[right]
			want := slices.Clone(values[left : right+1])
			slices.Sort(want)
			stats := Stats{}
			p := partition(values, left, right, 0, &stats, nil, false)
			if p < left || p > right || values[p] != pivot {
				t.Fatalf("p=%d values=%v", p, values)
			}
			for _, v := range values[left:p] {
				if v >= pivot {
					t.Fatal("left not < pivot")
				}
			}
			for _, v := range values[p+1 : right+1] {
				if v < pivot {
					t.Fatal("right not >= pivot")
				}
			}
			got := slices.Clone(values[left : right+1])
			slices.Sort(got)
			if !slices.Equal(got, want) || values[0] != 12345 || values[len(values)-1] != 67890 {
				t.Fatal("permutation or bounds")
			}
			if stats.Comparisons != right-left {
				t.Fatal("comparison count")
			}
		})
	}
}
func TestLastPivotDegeneration(t *testing.T) {
	for _, scenario := range []string{"sorted", "reverse", "equal"} {
		for _, n := range []int{2, 3, 4, 8, 16, 32, 64, 128, 256, 512} {
			t.Run(fmt.Sprintf("%s/n=%d", scenario, n), func(t *testing.T) {
				v := make([]int, n)
				for i := range v {
					switch scenario {
					case "sorted":
						v[i] = i
					case "reverse":
						v[i] = n - i
					default:
						v[i] = 7
					}
				}
				got := Sort(v)
				if got.Comparisons != n*(n-1)/2 || got.MaxDepth != n-1 || !slices.IsSorted(v) {
					t.Fatalf("stats=%+v", got)
				}
			})
		}
	}
}
