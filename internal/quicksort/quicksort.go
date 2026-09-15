// Package quicksort содержит прозрачную учебную реализацию Quick Sort.
package quicksort

import "fmt"

// Stats хранит число сравнений, перестановок, вызовов и максимальную глубину.
type Stats struct {
	Comparisons int
	Swaps       int
	Calls       int
	MaxDepth    int
}

// Step — снимок состояния во время разбиения среза.
// Values является копией, поэтому следующий шаг не изменит предыдущий снимок.
type Step struct {
	Depth      int
	Left       int
	Right      int
	Scan       int
	Store      int
	PivotIndex int
	PivotValue int
	Values     []int
	Action     string
}

// Sort сортирует values на месте и возвращает учебную статистику.
//
// Среднее время Quick Sort — O(n log n), худшее — O(n²), если разбиения каждый
// раз получаются очень неравными. Дополнительная память расходуется стеком
// рекурсивных вызовов: обычно O(log n), в худшем случае O(n).
//
// Это учебная реализация с последним элементом как pivot. В рабочем коде обычно
// выбирают slices.Sort: библиотечная реализация учитывает больше сложных случаев.
//
// Как не нужно: нельзя утверждать, что Quick Sort гарантированно O(n log n).
// Нельзя также забывать базовый случай, иначе рекурсия не остановится.
func Sort(values []int) Stats {
	stats := Stats{}
	quickSort(values, 0, len(values)-1, 0, &stats, nil, false)
	return stats
}

// SortWithTrace сортирует values и сохраняет состояния для урока и отладчика.
// Сами снимки требуют дополнительную память; это не часть обычного Sort.
func SortWithTrace(values []int) (Stats, []Step) {
	stats := Stats{}
	steps := make([]Step, 0)
	quickSort(values, 0, len(values)-1, 0, &stats, &steps, true)
	return stats, steps
}

func quickSort(values []int, left, right, depth int, stats *Stats, steps *[]Step, trace bool) {
	stats.Calls++
	if depth > stats.MaxDepth {
		stats.MaxDepth = depth
	}

	if left >= right {
		if trace && left == right {
			record(steps, values, depth, left, right, -1, left, left, values[left], "один элемент уже отсортирован")
		}
		return
	}

	pivotIndex := partition(values, left, right, depth, stats, steps, trace)
	quickSort(values, left, pivotIndex-1, depth+1, stats, steps, trace)
	quickSort(values, pivotIndex+1, right, depth+1, stats, steps, trace)
}

func partition(values []int, left, right, depth int, stats *Stats, steps *[]Step, trace bool) int {
	pivot := values[right]
	store := left
	if trace {
		record(steps, values, depth, left, right, -1, store, right, pivot, "выбрали опорное значение")
	}

	for scan := left; scan < right; scan++ {
		stats.Comparisons++
		if values[scan] < pivot {
			compared := values[scan]
			if scan != store {
				values[scan], values[store] = values[store], values[scan]
				stats.Swaps++
			}
			store++
			if trace {
				action := fmt.Sprintf("%d < %d: переносим в левую часть", compared, pivot)
				record(steps, values, depth, left, right, scan, store, right, pivot, action)
			}
		} else if trace {
			action := fmt.Sprintf("%d >= %d: оставляем справа", values[scan], pivot)
			record(steps, values, depth, left, right, scan, store, right, pivot, action)
		}
	}

	if store != right {
		values[store], values[right] = values[right], values[store]
		stats.Swaps++
	}
	if trace {
		record(steps, values, depth, left, right, right, store, store, pivot, "поставили опорное значение на окончательную позицию")
	}
	return store
}

func record(steps *[]Step, values []int, depth, left, right, scan, store, pivotIndex, pivotValue int, action string) {
	*steps = append(*steps, Step{
		Depth:      depth,
		Left:       left,
		Right:      right,
		Scan:       scan,
		Store:      store,
		PivotIndex: pivotIndex,
		PivotValue: pivotValue,
		Values:     append([]int(nil), values...),
		Action:     action,
	})
}

// Example сортирует маленький срез и показывает статистику текущей реализации.
func Example() string {
	values := []int{7, 2, 9, 3, 5}
	stats := Sort(values)
	return fmt.Sprintf("Quick Sort: %v; comparisons=%d swaps=%d maxDepth=%d", values, stats.Comparisons, stats.Swaps, stats.MaxDepth)
}
