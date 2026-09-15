// Package search показывает линейный и бинарный поиск на одном контракте.
package search

import "fmt"

// Result описывает итог поиска и позволяет увидеть реальное число сравнений.
type Result struct {
	Index       int
	Found       bool
	Comparisons int
}

// Step хранит одно состояние бинарного поиска перед изменением границ.
type Step struct {
	Left     int
	Right    int
	Middle   int
	Value    int
	Decision string
}

// LinearSearch проверяет элементы слева направо.
//
// Лучший случай — нужное значение стоит первым: O(1).
// Худший случай — значение последнее или отсутствует: O(n).
// Дополнительная память — O(1).
//
// Как не нужно: нельзя написать O(1) только потому, что внутри цикла одно
// сравнение. Это сравнение может повториться n раз.
func LinearSearch(values []int, target int) Result {
	for index, value := range values {
		if value == target {
			return Result{Index: index, Found: true, Comparisons: index + 1}
		}
	}
	return Result{Index: -1, Found: false, Comparisons: len(values)}
}

// BinarySearch ищет target в срезе, отсортированном по возрастанию.
//
// На каждом шаге остаётся примерно половина диапазона, поэтому худшее время —
// O(log n), а дополнительная память итеративной версии — O(1).
//
// Как не нужно: нельзя запускать бинарный поиск на неотсортированном срезе.
// Код может завершиться без panic, но вернуть неверный ответ.
func BinarySearch(values []int, target int) Result {
	left, right := 0, len(values)
	comparisons := 0
	for left < right {
		middle := left + (right-left)/2
		comparisons++
		switch {
		case values[middle] == target:
			return Result{Index: middle, Found: true, Comparisons: comparisons}
		case values[middle] < target:
			left = middle + 1
		default:
			right = middle
		}
	}
	return Result{Index: -1, Found: false, Comparisons: comparisons}
}

// BinarySearchTrace выполняет тот же бинарный поиск и сохраняет каждый шаг.
// Правая граница right не входит в диапазон: активная часть — [left, right).
func BinarySearchTrace(values []int, target int) (Result, []Step) {
	left, right := 0, len(values)
	steps := make([]Step, 0)

	for left < right {
		middle := left + (right-left)/2
		step := Step{Left: left, Right: right, Middle: middle, Value: values[middle]}

		switch {
		case values[middle] == target:
			step.Decision = "найдено"
			steps = append(steps, step)
			return Result{Index: middle, Found: true, Comparisons: len(steps)}, steps
		case values[middle] < target:
			step.Decision = "ищем справа"
			left = middle + 1
		default:
			step.Decision = "ищем слева"
			right = middle
		}

		steps = append(steps, step)
	}

	return Result{Index: -1, Found: false, Comparisons: len(steps)}, steps
}

// SearchInsert возвращает индекс найденного target или место, куда его нужно
// вставить, чтобы сохранить сортировку. Это поиск левой границы за O(log n).
//
// Как не нужно: после неудачного BinarySearch нельзя без дополнительной логики
// считать последний middle правильной позицией вставки.
func SearchInsert(values []int, target int) int {
	left, right := 0, len(values)
	for left < right {
		middle := left + (right-left)/2
		if values[middle] < target {
			left = middle + 1
		} else {
			right = middle
		}
	}
	return left
}

// Example возвращает сравнимые результаты поиска одного значения.
func Example() string {
	values := []int{3, 7, 11, 15, 19, 23, 27}
	linear := LinearSearch(values, 23)
	binary := BinarySearch(values, 23)
	return fmt.Sprintf("target=23: linear index=%d comparisons=%d; binary index=%d comparisons=%d; insert(20)=%d",
		linear.Index, linear.Comparisons, binary.Index, binary.Comparisons, SearchInsert(values, 20))
}
