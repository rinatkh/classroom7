package maps

import (
	"reflect"
	"testing"
)

func TestBuildFirstIndexAndLookup(t *testing.T) {
	index := BuildFirstIndex([]int{8, 3, 8, 5})
	want := map[int]int{8: 0, 3: 1, 5: 3}
	if !reflect.DeepEqual(index, want) {
		t.Fatalf("BuildFirstIndex() = %v; want %v", index, want)
	}
	if got, ok := Lookup(index, 5); !ok || got != 3 {
		t.Fatalf("Lookup(5) = %d, %v; want 3, true", got, ok)
	}
	if got, ok := Lookup(index, 42); ok || got != -1 {
		t.Fatalf("Lookup(42) = %d, %v; want -1, false", got, ok)
	}
}

func TestFrequencies(t *testing.T) {
	got := Frequencies([]int{2, 1, 2, 3, 2})
	want := map[int]int{1: 1, 2: 3, 3: 1}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Frequencies() = %v; want %v", got, want)
	}
}

func TestContainsDuplicate(t *testing.T) {
	for _, test := range []struct {
		values []int
		want   bool
	}{
		{values: nil, want: false},
		{values: []int{1, 2, 3}, want: false},
		{values: []int{1, 2, 1}, want: true},
	} {
		if got := ContainsDuplicate(test.values); got != test.want {
			t.Errorf("ContainsDuplicate(%v) = %t; want %t", test.values, got, test.want)
		}
	}
}

func TestExample(t *testing.T) {
	const want = "index(8)=0 found=true; index(9)=-1 found=false; count(8)=2; duplicate=true"
	if got := Example(); got != want {
		t.Fatalf("Example() = %q; want %q", got, want)
	}
}
