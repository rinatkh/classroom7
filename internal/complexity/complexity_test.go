package complexity

import "testing"

func TestAt(t *testing.T) {
	values := []int{10, 20, 30}
	if got, ok := At(values, 1); !ok || got != 20 {
		t.Fatalf("At(valid) = %d, %v; want 20, true", got, ok)
	}
	for _, index := range []int{-1, len(values)} {
		if got, ok := At(values, index); ok || got != 0 {
			t.Fatalf("At(%d) = %d, %v; want 0, false", index, got, ok)
		}
	}
}

func TestSum(t *testing.T) {
	for _, test := range []struct {
		name   string
		values []int
		want   Measurement
	}{
		{name: "empty", values: nil, want: Measurement{}},
		{name: "numbers", values: []int{4, -1, 7}, want: Measurement{Value: 10, Steps: 3}},
	} {
		t.Run(test.name, func(t *testing.T) {
			if got := Sum(test.values); got != test.want {
				t.Fatalf("Sum(%v) = %+v; want %+v", test.values, got, test.want)
			}
		})
	}
}

func TestHalvingSteps(t *testing.T) {
	for n, want := range map[int]int{0: 0, 1: 0, 2: 1, 3: 2, 16: 4, 17: 5} {
		if got := HalvingSteps(n); got != want {
			t.Errorf("HalvingSteps(%d) = %d; want %d", n, got, want)
		}
	}
}

func TestPairComparisons(t *testing.T) {
	for n, want := range map[int]int{0: 0, 1: 0, 2: 1, 4: 6, 5: 10} {
		if got := PairComparisons(make([]int, n)); got != want {
			t.Errorf("PairComparisons(n=%d) = %d; want %d", n, got, want)
		}
	}
}

func TestExample(t *testing.T) {
	const want = "O(1): value=40; O(n): sum=100 steps=4; O(log n): steps=4; O(n²): pairs=6"
	if got := Example(); got != want {
		t.Fatalf("Example() = %q; want %q", got, want)
	}
}
