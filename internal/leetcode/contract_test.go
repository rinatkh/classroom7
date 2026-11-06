package leetcode

import (
	"github.com/rinatkh/classroom7/internal/testutil"
	"slices"
	"testing"
)

func Test704Contracts(t *testing.T) {
	for _, tc := range testutil.SearchCases() {
		t.Run(tc.Name, func(t *testing.T) {
			// Внешняя задача требует уникальных чисел; Compact сохраняет sorted порядок.
			values := slices.Compact(slices.Clone(tc.Values))
			saved := slices.Clone(values)
			want := slices.Index(values, tc.Target)
			if got := BinarySearch704(values, tc.Target); got != want {
				t.Fatalf("%d want %d", got, want)
			}
			if !slices.Equal(values, saved) {
				t.Fatal("input mutated")
			}
		})
	}
}
