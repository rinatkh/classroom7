package search

import (
	"reflect"
	"testing"
)

func TestLinearSearch(t *testing.T) {
	values := []int{4, 8, 15, 16, 23, 42}
	for _, test := range []struct {
		name   string
		target int
		want   Result
	}{
		{name: "first", target: 4, want: Result{Index: 0, Found: true, Comparisons: 1}},
		{name: "last", target: 42, want: Result{Index: 5, Found: true, Comparisons: 6}},
		{name: "missing", target: 7, want: Result{Index: -1, Found: false, Comparisons: 6}},
	} {
		t.Run(test.name, func(t *testing.T) {
			if got := LinearSearch(values, test.target); got != test.want {
				t.Fatalf("LinearSearch(%d) = %+v; want %+v", test.target, got, test.want)
			}
		})
	}
}

func TestBinarySearchAndTrace(t *testing.T) {
	values := []int{3, 7, 11, 15, 19, 23, 27}
	got, steps := BinarySearchTrace(values, 23)
	if want := (Result{Index: 5, Found: true, Comparisons: 2}); got != want {
		t.Fatalf("BinarySearchTrace() result = %+v; want %+v", got, want)
	}
	wantSteps := []Step{
		{Left: 0, Right: 7, Middle: 3, Value: 15, Decision: "ищем справа"},
		{Left: 4, Right: 7, Middle: 5, Value: 23, Decision: "найдено"},
	}
	if !reflect.DeepEqual(steps, wantSteps) {
		t.Fatalf("steps = %#v; want %#v", steps, wantSteps)
	}
	if direct := BinarySearch(values, 23); direct != got {
		t.Fatalf("BinarySearch() = %+v; want %+v", direct, got)
	}
}

func TestBinarySearchMissingAndEmpty(t *testing.T) {
	for _, values := range [][]int{nil, {1, 3, 5, 7}} {
		got, steps := BinarySearchTrace(values, 4)
		if got.Found || got.Index != -1 || got.Comparisons != len(steps) {
			t.Fatalf("BinarySearchTrace(%v, 4) = %+v, %#v", values, got, steps)
		}
		direct := BinarySearch(values, 4)
		if direct.Found || direct.Index != -1 {
			t.Fatalf("BinarySearch(%v, 4) = %+v; want missing", values, direct)
		}
	}
}

func TestSearchInsert(t *testing.T) {
	values := []int{1, 3, 5, 6}
	for target, want := range map[int]int{0: 0, 1: 0, 2: 1, 5: 2, 7: 4} {
		if got := SearchInsert(values, target); got != want {
			t.Errorf("SearchInsert(%d) = %d; want %d", target, got, want)
		}
	}
}

func TestExample(t *testing.T) {
	const want = "target=23: linear index=5 comparisons=6; binary index=5 comparisons=2; insert(20)=5"
	if got := Example(); got != want {
		t.Fatalf("Example() = %q; want %q", got, want)
	}
}
