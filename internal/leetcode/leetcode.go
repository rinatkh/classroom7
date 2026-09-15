// Package leetcode содержит совместный разбор одной внешней задачи.
package leetcode

import (
	"fmt"

	"github.com/rinatkh/classroom7/internal/search"
)

// BinarySearch704 решает адаптированную задачу LeetCode 704 Binary Search.
// Оригинал условия: https://leetcode.com/problems/binary-search/
//
// Контракт своими словами:
//  1. На входе отсортированный по возрастанию срез различных int.
//  2. Если target существует, возвращаем его индекс.
//  3. Если target отсутствует, возвращаем -1.
//
// Как рассуждать:
//  1. Слово «отсортированный» разрешает бинарный поиск.
//  2. Держим полуинтервал [left, right), в котором target ещё возможен.
//  3. Сравнение с серединой исключает половину диапазона.
//  4. Цикл заканчивается находкой или пустым диапазоном.
//  5. Время O(log n), дополнительная память O(1).
//
// Как не нужно: сортировать вход внутри функции нельзя. Это изменит индексы
// исходного среза и добавит работу O(n log n), хотя условие уже даёт сортировку.
func BinarySearch704(values []int, target int) int {
	result := search.BinarySearch(values, target)
	if !result.Found {
		return -1
	}
	return result.Index
}

// Example показывает найденный и отсутствующий target.
func Example() string {
	values := []int{-1, 0, 3, 5, 9, 12}
	return fmt.Sprintf("LeetCode 704: target=9 -> %d; target=2 -> %d; time=O(log n); memory=O(1)",
		BinarySearch704(values, 9), BinarySearch704(values, 2))
}
