// Package mergesort показывает устойчивое слияние с одним общим буфером.
package mergesort

import (
	"cmp"
	"fmt"
	"slices"
)

// Step сохраняет состояние после слияния [Left,Middle) и [Middle,Right).
// Values — независимый снимок всего среза; трассировка не участвует в benchmark.
type Step struct {
	Left, Middle, Right int
	Values              []int
}

// Sort сортирует вход на месте, но использует дополнительный массив O(n).
// Без пропуска готовых слияний лучший/средний/худший случай — O(n log n).
// Стек — O(log n), общий буфер — O(n). Пустой вход уже отсортирован.
func Sort(values []int) { SortFunc(values, cmp.Compare[int]) }

// SortFunc позволяет увидеть устойчивость на записях с одинаковым ключом.
// compare должен задавать согласованный порядок и возвращать 0 для равных ключей.
func SortFunc[T any](values []T, compare func(T, T) int) { sortWith(values, compare, nil) }

func sortWith[T any](values []T, compare func(T, T) int, observe func(int, int, int)) {
	if len(values) < 2 {
		return
	}
	buffer := make([]T, len(values)) // Один буфер на весь запуск, не на каждый вызов.
	var sortRange func(int, int)
	sortRange = func(left, right int) {
		if right-left < 2 {
			return
		} // Условие завершения рекурсии.
		middle := left + (right-left)/2
		sortRange(left, middle)
		sortRange(middle, right)
		// Перед перезаписью сохраняем обе отсортированные половины.
		copy(buffer[left:right], values[left:right])
		i, j := left, middle
		for write := left; write < right; write++ {
			// При равенстве берём ЛЕВЫЙ элемент: так сохраняется устойчивость.
			if i < middle && (j == right || compare(buffer[i], buffer[j]) <= 0) {
				values[write] = buffer[i]
				i++
			} else {
				values[write] = buffer[j]
				j++
			}
		}
		if observe != nil {
			observe(left, middle, right)
		}
	}
	sortRange(0, len(values))
}

// SortWithTrace использует тот же алгоритм, добавляя снимок после каждого merge.
// Не оценивайте память обычного Sort по памяти истории снимков.
func SortWithTrace(values []int) []Step {
	steps := make([]Step, 0)
	sortWith(values, cmp.Compare[int], func(left, middle, right int) {
		steps = append(steps, Step{left, middle, right, slices.Clone(values)})
	})
	return steps
}

// Example показывает малый вход, границы слияний и итог.
func Example() string {
	values := []int{5, 1, 4, 2}
	steps := SortWithTrace(values)
	return fmt.Sprintf("Merge Sort: %v; merge steps=%v", values, steps)
}
