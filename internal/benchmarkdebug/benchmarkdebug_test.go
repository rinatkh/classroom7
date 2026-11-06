package benchmarkdebug

import (
	"strings"
	"testing"

	"github.com/rinatkh/classroom7/internal/quicksort"
)

func TestFormatTrace(t *testing.T) {
	steps := []quicksort.Step{{
		Depth: 0, Left: 0, Right: 2, Scan: -1, Store: 0,
		PivotIndex: 2, PivotValue: 2, Values: []int{3, 1, 2}, Action: "выбрали",
	}}
	const want = "шаг 01 | глубина=0 | границы=0..2 | scan=- | store=0 | pivot=2@2 | [3 1 2] | выбрали"
	if got := FormatTrace(steps); got != want {
		t.Fatalf("FormatTrace() = %q; want %q", got, want)
	}
	if got := FormatTrace(nil); got != "" {
		t.Fatalf("FormatTrace(nil) = %q; want empty", got)
	}
}

func TestExample(t *testing.T) {
	got := Example()
	for _, part := range []string{"До: [7 2 9 3 5]", "шаг 01", "После: [2 3 5 7 9]", "make bench-search"} {
		if !strings.Contains(got, part) {
			t.Fatalf("Example() does not contain %q:\n%s", part, got)
		}
	}
}
