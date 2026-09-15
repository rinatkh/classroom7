package simplesort

import (
	"reflect"
	"testing"
)

func TestSorts(t *testing.T) {
	sorts := map[string]func([]int) Report{
		"bubble":    BubbleSort,
		"selection": SelectionSort,
		"insertion": InsertionSort,
	}

	for name, sortFn := range sorts {
		t.Run(name, func(t *testing.T) {
			for _, input := range [][]int{nil, {}, {1}, {3, 1, 2}, {2, 2, 1}, {-1, 4, 0}} {
				values := append([]int(nil), input...)
				sortFn(values)
				want := append([]int(nil), input...)
				standardInsertion(want)
				if !reflect.DeepEqual(values, want) {
					t.Fatalf("sort(%v) = %v; want %v", input, values, want)
				}
			}
		})
	}
}

func standardInsertion(values []int) {
	for i := 1; i < len(values); i++ {
		for j := i; j > 0 && values[j] < values[j-1]; j-- {
			values[j], values[j-1] = values[j-1], values[j]
		}
	}
}

func TestBubbleAlreadySortedStopsEarly(t *testing.T) {
	values := []int{1, 2, 3, 4}
	report := BubbleSort(values)
	if report != (Report{Comparisons: 3, Writes: 0}) {
		t.Fatalf("BubbleSort(sorted) = %+v", report)
	}
}

func TestBubbleTrace(t *testing.T) {
	got := BubbleTrace([]int{3, 2, 1})
	want := [][]int{{3, 2, 1}, {2, 1, 3}, {1, 2, 3}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("BubbleTrace() = %v; want %v", got, want)
	}
}

func TestExample(t *testing.T) {
	const want = "sorted=[1 2 4 5]; bubble=6/8; selection=6/4; insertion=6/7"
	if got := Example(); got != want {
		t.Fatalf("Example() = %q; want %q", got, want)
	}
}
