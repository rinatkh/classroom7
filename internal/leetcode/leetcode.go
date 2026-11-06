// Package leetcode содержит совместный разбор одной внешней задачи.
package leetcode

import (
	"fmt"

	"github.com/rinatkh/classroom7/internal/search"
)

// BinarySearch704 возвращает индекс target или -1, не меняя вход.
// Вход отсортирован и содержит различные int. Время до O(log n), память O(1).
// Оригинал: https://leetcode.com/problems/binary-search/
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
