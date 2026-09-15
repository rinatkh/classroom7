package quicksort

import (
	"reflect"
	"slices"
	"testing"
)

func TestSort(t *testing.T) {
	for _, input := range [][]int{nil, {}, {1}, {3, 1, 2}, {5, 1, 1, 0, -2}, {1, 2, 3, 4}} {
		values := append([]int(nil), input...)
		stats := Sort(values)
		want := append([]int(nil), input...)
		slices.Sort(want)
		if !reflect.DeepEqual(values, want) {
			t.Fatalf("Sort(%v) = %v; want %v", input, values, want)
		}
		if len(values) > 1 && stats.Calls == 0 {
			t.Fatalf("Sort(%v) returned empty stats: %+v", input, stats)
		}
	}
}

func TestSortWithTraceCopiesStates(t *testing.T) {
	values := []int{3, 1, 2}
	stats, steps := SortWithTrace(values)
	if !slices.Equal(values, []int{1, 2, 3}) || stats.Comparisons == 0 || len(steps) == 0 {
		t.Fatalf("SortWithTrace() = %v, %+v, %d steps", values, stats, len(steps))
	}
	first := append([]int(nil), steps[0].Values...)
	values[0] = 99
	if !reflect.DeepEqual(steps[0].Values, first) {
		t.Fatal("первый снимок изменился вместе с исходным срезом")
	}
	foundBase := false
	for _, step := range steps {
		if step.Action == "один элемент уже отсортирован" {
			foundBase = true
		}
	}
	if !foundBase {
		t.Fatal("трассировка не содержит базовый случай")
	}
}

func TestExample(t *testing.T) {
	const want = "Quick Sort: [2 3 5 7 9]; comparisons=6 swaps=3 maxDepth=2"
	if got := Example(); got != want {
		t.Fatalf("Example() = %q; want %q", got, want)
	}
}
