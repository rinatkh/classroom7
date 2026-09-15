package benchmarkdebug

import (
	"fmt"
	"math/rand"
	"slices"
	"strings"
	"testing"

	algomaps "github.com/rinatkh/classroom7/internal/maps"
	"github.com/rinatkh/classroom7/internal/quicksort"
	"github.com/rinatkh/classroom7/internal/search"
	"github.com/rinatkh/classroom7/internal/simplesort"
)

var (
	searchSink bool
	indexSink  int
	sortSink   int
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

// BenchmarkSearch сравнивает три разных сценария на одинаковых данных.
// Подготовка отсортированного среза и готовой map находится вне таймера поиска.
// Построение map измеряется отдельной веткой map-build.
func BenchmarkSearch(b *testing.B) {
	for _, size := range []int{64, 4096, 262144} {
		values := make([]int, size)
		for index := range values {
			values[index] = index * 2
		}
		target := values[len(values)-1]
		index := algomaps.BuildFirstIndex(values)

		b.Run(fmt.Sprintf("linear/n=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for iteration := 0; iteration < b.N; iteration++ {
				result := search.LinearSearch(values, target)
				searchSink = result.Found
			}
		})

		b.Run(fmt.Sprintf("binary/n=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for iteration := 0; iteration < b.N; iteration++ {
				result := search.BinarySearch(values, target)
				searchSink = result.Found
			}
		})

		b.Run(fmt.Sprintf("map-read/n=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for iteration := 0; iteration < b.N; iteration++ {
				position, found := algomaps.Lookup(index, target)
				indexSink, searchSink = position, found
			}
		})

		b.Run(fmt.Sprintf("map-build/n=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for iteration := 0; iteration < b.N; iteration++ {
				built := algomaps.BuildFirstIndex(values)
				indexSink = len(built)
			}
		})
	}
}

// BenchmarkSort включает копирование одинакового входа в измеряемый участок
// для каждого алгоритма. Это честно сравнивает одинаковую полезную работу.
// Bubble и Selection запускаются только на небольших n, чтобы урок не зависал.
func BenchmarkSort(b *testing.B) {
	for _, size := range []int{64, 512} {
		source := randomValues(size)
		// Срез, а не map, фиксирует порядок запуска подзадач benchmark.
		// Это не устраняет шум компьютера, но делает повторные запуски понятнее.
		algorithms := []struct {
			name   string
			sortFn func([]int)
		}{
			{name: "bubble", sortFn: func(values []int) { simplesort.BubbleSort(values) }},
			{name: "selection", sortFn: func(values []int) { simplesort.SelectionSort(values) }},
			{name: "insertion", sortFn: func(values []int) { simplesort.InsertionSort(values) }},
			{name: "quick", sortFn: func(values []int) { quicksort.Sort(values) }},
			{name: "slices", sortFn: slices.Sort[[]int, int]},
		}

		for _, algorithm := range algorithms {
			b.Run(fmt.Sprintf("%s/n=%d", algorithm.name, size), func(b *testing.B) {
				work := make([]int, size)
				b.ReportAllocs()
				b.ResetTimer()
				for iteration := 0; iteration < b.N; iteration++ {
					copy(work, source)
					algorithm.sortFn(work)
					sortSink = work[size/2]
				}
			})
		}
	}
}

func randomValues(size int) []int {
	random := rand.New(rand.NewSource(7))
	values := make([]int, size)
	for index := range values {
		values[index] = random.Intn(size * 4)
	}
	return values
}
