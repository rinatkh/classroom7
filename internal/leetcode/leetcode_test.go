package leetcode

import "testing"

func TestBinarySearch704(t *testing.T) {
	values := []int{-1, 0, 3, 5, 9, 12}
	for target, want := range map[int]int{-1: 0, 9: 4, 12: 5, 2: -1, 20: -1} {
		if got := BinarySearch704(values, target); got != want {
			t.Errorf("BinarySearch704(%d) = %d; want %d", target, got, want)
		}
	}
	if got := BinarySearch704(nil, 1); got != -1 {
		t.Fatalf("BinarySearch704(nil, 1) = %d; want -1", got)
	}
}

func TestExample(t *testing.T) {
	const want = "LeetCode 704: target=9 -> 4; target=2 -> -1; time=O(log n); memory=O(1)"
	if got := Example(); got != want {
		t.Fatalf("Example() = %q; want %q", got, want)
	}
}
