// Package benchmarkdebug связывает benchmark-тесты и пошаговую трассировку.
package benchmarkdebug

import (
	"fmt"
	"strings"

	"github.com/rinatkh/classroom7/internal/quicksort"
)

// FormatTrace превращает снимки Quick Sort в таблицу для терминала.
// scan=- означает, что цикл сравнения ещё не начался.
func FormatTrace(steps []quicksort.Step) string {
	var builder strings.Builder
	for index, step := range steps {
		scan := "-"
		if step.Scan >= 0 {
			scan = fmt.Sprint(step.Scan)
		}
		fmt.Fprintf(&builder,
			"шаг %02d | глубина=%d | границы=%d..%d | scan=%s | store=%d | pivot=%d@%d | %v | %s\n",
			index+1, step.Depth, step.Left, step.Right, scan, step.Store,
			step.PivotValue, step.PivotIndex, step.Values, step.Action)
	}
	return strings.TrimSuffix(builder.String(), "\n")
}

// Example строит полную трассировку маленького входа.
//
// Как не нужно: эту функцию нельзя использовать для измерения скорости Sort.
// Создание строк и копий намеренно замедляет выполнение. Для скорости служат
// BenchmarkSearch и BenchmarkSort в benchmarkdebug_test.go.
func Example() string {
	values := []int{7, 2, 9, 3, 5}
	stats, steps := quicksort.SortWithTrace(values)
	return fmt.Sprintf("До: [7 2 9 3 5]\n%s\nПосле: %v; сравнений=%d; перестановок=%d\nИзмерения: make bench-search / make bench-sort",
		FormatTrace(steps), values, stats.Comparisons, stats.Swaps)
}
