// Package simplesort содержит учебные реализации трёх сортировок O(n²).
package simplesort

import "fmt"

// Report показывает, сколько сравнений и записей в срез выполнил алгоритм.
// Это учебные счётчики, а не точное число процессорных инструкций.
type Report struct {
	Comparisons int
	Writes      int
}

// BubbleSort сортирует на месте обменами соседей; равные не переставляет.
// Ранний выход даёт O(n) на sorted; среднее/худшее O(n²), память O(1).
func BubbleSort(values []int) Report { return bubbleSort(values, nil) }

func bubbleSort(values []int, observe func()) Report {
	report := Report{}
	for end := len(values) - 1; end > 0; end-- {
		swapped := false
		for index := 0; index < end; index++ {
			report.Comparisons++
			if values[index] > values[index+1] {
				values[index], values[index+1] = values[index+1], values[index]
				report.Writes += 2
				swapped = true
			}
		}
		if observe != nil {
			observe()
		}
		if !swapped {
			break
		}
	}
	return report
}

// BubbleTrace возвращает копию среза после каждого полного прохода.
// Копии намеренно требуют O(n²) памяти в худшем случае: это цена визуализации,
// а не память самого BubbleSort.
func BubbleTrace(values []int) [][]int { return trace(values, bubbleSort) }

// SelectionTrace сохраняет начальное состояние и результат каждого выбора минимума.
func SelectionTrace(values []int) [][]int { return trace(values, selectionSort) }

// InsertionTrace сохраняет начальное состояние и результат каждой вставки в префикс.
func InsertionTrace(values []int) [][]int { return trace(values, insertionSort) }

func trace(values []int, algorithm func([]int, func()) Report) [][]int {
	states := [][]int{append([]int(nil), values...)}
	algorithm(values, func() { states = append(states, append([]int(nil), values...)) })
	return states
}

// SelectionSort ставит минимум оставшегося хвоста на очередную позицию.
// Всегда n*(n-1)/2 сравнений, память O(1). Вариант с обменом неустойчив.
func SelectionSort(values []int) Report { return selectionSort(values, nil) }

func selectionSort(values []int, observe func()) Report {
	report := Report{}
	for position := 0; position < len(values)-1; position++ {
		minimum := position
		for candidate := position + 1; candidate < len(values); candidate++ {
			report.Comparisons++
			if values[candidate] < values[minimum] {
				minimum = candidate
			}
		}
		if minimum != position {
			values[position], values[minimum] = values[minimum], values[position]
			report.Writes += 2
		}
		if observe != nil {
			observe()
		}
	}
	return report
}

// InsertionSort вставляет current в отсортированный левый префикс.
// Равные не сдвигаются через current: сортировка устойчива.
// Время O(n+I), где I — число инверсий; память O(1).
func InsertionSort(values []int) Report { return insertionSort(values, nil) }

func insertionSort(values []int, observe func()) Report {
	report := Report{}
	for position := 1; position < len(values); position++ {
		current := values[position]
		insert := position

		for insert > 0 {
			report.Comparisons++
			if values[insert-1] <= current {
				break
			}
			values[insert] = values[insert-1]
			report.Writes++
			insert--
		}

		values[insert] = current
		report.Writes++
		if observe != nil {
			observe()
		}
	}
	return report
}

// Example сортирует одинаковый вход тремя способами и показывает счётчики.
func Example() string {
	input := []int{5, 1, 4, 2}
	bubbleValues := append([]int(nil), input...)
	selectionValues := append([]int(nil), input...)
	insertionValues := append([]int(nil), input...)
	bubble := BubbleSort(bubbleValues)
	selection := SelectionSort(selectionValues)
	insertion := InsertionSort(insertionValues)
	return fmt.Sprintf("sorted=%v; bubble=%d/%d; selection=%d/%d; insertion=%d/%d",
		bubbleValues, bubble.Comparisons, bubble.Writes, selection.Comparisons, selection.Writes,
		insertion.Comparisons, insertion.Writes)
}
