// Package search показывает линейный и бинарный поиск на одном контракте.
package search

import (
	"fmt"
	"slices"
)

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

// LinearSearch возвращает первое совпадение или Index=-1, Found=false.
// Время от O(1) до O(n), дополнительная память O(1).
func LinearSearch(values []int, target int) Result {
	for index, value := range values {
		if value == target {
			return Result{Index: index, Found: true, Comparisons: index + 1}
		}
	}
	return Result{Index: -1, Found: false, Comparisons: len(values)}
}

// BinarySearch ищет в неубывающем срезе; при повторах допустим любой индекс.
// Активный диапазон [left, right): правая граница не входит.
// Comparisons считает проверенные середины. Время до O(log n), память O(1).
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

// SearchInsert возвращает первую позицию со значением >= target, включая len(values).
// Вход неубывающий; время O(log n) на непустом входе, память O(1).
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
	standardIndex, standardFound := slices.BinarySearch(values, 23)
	return fmt.Sprintf("target=23: linear index=%d comparisons=%d; binary index=%d comparisons=%d; insert(20)=%d; slices.BinarySearch=%d/%t",
		linear.Index, linear.Comparisons, binary.Index, binary.Comparisons, SearchInsert(values, 20), standardIndex, standardFound)
}
