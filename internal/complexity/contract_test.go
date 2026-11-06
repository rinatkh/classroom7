package complexity

import (
	"fmt"
	"github.com/rinatkh/classroom7/internal/testutil"
	"math"
	"math/bits"
	"testing"
)

func TestComplexityContracts(t *testing.T) {
	for _, tc := range testutil.Cases() {
		t.Run(tc.Name, func(t *testing.T) {
			// Суммирование int намеренно оставляет обычную семантику переполнения Go.
			want := 0
			for _, v := range tc.Values {
				want += v
			}
			got := Sum(tc.Values)
			if got.Value != want || got.Steps != len(tc.Values) {
				t.Fatal("sum")
			}
			if PairComparisons(tc.Values) != len(tc.Values)*(len(tc.Values)-1)/2 {
				t.Fatal("pairs")
			}
			for _, i := range []int{-1, 0, len(tc.Values) / 2, len(tc.Values) - 1, len(tc.Values)} {
				got, ok := At(tc.Values, i)
				valid := i >= 0 && i < len(tc.Values)
				if ok != valid || (ok && got != tc.Values[i]) {
					t.Fatal("at bounds")
				}
			}
		})
	}
	for _, n := range []int{-2, -1, 0, 1, 2, 3, 4, 7, 8, 16, 17, 100, math.MaxInt} {
		t.Run(fmt.Sprintf("halving/%d", n), func(t *testing.T) {
			want := 0
			if n > 1 {
				want = bits.Len(uint(n - 1))
			}
			if got := HalvingSteps(n); got != want {
				t.Fatalf("%d want %d", got, want)
			}
		})
	}
}
