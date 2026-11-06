package maps

import (
	"github.com/rinatkh/classroom7/internal/testutil"
	"maps"
	"slices"
	"testing"
)

func TestMapContracts(t *testing.T) {
	for _, tc := range testutil.Cases() {
		t.Run(tc.Name, func(t *testing.T) {
			saved := slices.Clone(tc.Values)
			wantIndex, wantCounts := map[int]int{}, map[int]int{}
			duplicate := false
			for i, v := range tc.Values {
				if wantCounts[v] > 0 {
					duplicate = true
				} else {
					wantIndex[v] = i
				}
				wantCounts[v]++
			}
			t.Run("build", func(t *testing.T) {
				if !maps.Equal(BuildFirstIndex(tc.Values), wantIndex) {
					t.Fatal("first index")
				}
			})
			t.Run("frequencies", func(t *testing.T) {
				if !maps.Equal(Frequencies(tc.Values), wantCounts) {
					t.Fatal("counts")
				}
			})
			t.Run("duplicate", func(t *testing.T) {
				if ContainsDuplicate(tc.Values) != duplicate {
					t.Fatal("duplicate")
				}
			})
			t.Run("lookup", func(t *testing.T) {
				index := BuildFirstIndex(tc.Values)
				for v, want := range wantIndex {
					got, ok := Lookup(index, v)
					if !ok || got != want {
						t.Fatal("present key")
					}
				}
				missing := 42
				for {
					if _, ok := wantIndex[missing]; !ok {
						break
					}
					missing++
				}
				if got, ok := Lookup(index, missing); got != -1 || ok {
					t.Fatal("missing key")
				}
			})
			if !slices.Equal(saved, tc.Values) {
				t.Fatal("input changed")
			}
		})
	}
}
