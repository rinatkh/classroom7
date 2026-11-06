// Package maps показывает, когда индекс на основе map ускоряет повторный поиск.
package maps

import "fmt"

// BuildFirstIndex сохраняет первый индекс каждого числа.
// Построение — ожидаемо O(n) времени и O(n) дополнительной памяти.
func BuildFirstIndex(values []int) map[int]int {
	index := make(map[int]int, len(values))
	for position, value := range values {
		if _, exists := index[value]; !exists {
			index[value] = position
		}
	}
	return index
}

// Lookup читает готовый индекс; отсутствующий ключ даёт -1, false.
// Второй результат отличает отсутствие от корректного индекса 0.
func Lookup(index map[int]int, target int) (int, bool) {
	position, exists := index[target]
	if !exists {
		return -1, false
	}
	return position, true
}

// Frequencies считает появления значений за ожидаемое O(n) время.
// Память O(k), где k — число различных значений.
func Frequencies(values []int) map[int]int {
	frequencies := make(map[int]int)
	for _, value := range values {
		frequencies[value]++
	}
	return frequencies
}

// ContainsDuplicate завершает обход на первом повторе.
// make заранее резервирует место для n ключей; полный вызов не обещает O(1).
// Ожидаемое время O(n), дополнительная память O(n).
func ContainsDuplicate(values []int) bool {
	seen := make(map[int]struct{}, len(values))
	for _, value := range values {
		if _, exists := seen[value]; exists {
			return true
		}
		seen[value] = struct{}{}
	}
	return false
}

// Example показывает цену построения один раз и два быстрых чтения.
func Example() string {
	values := []int{8, 3, 8, 5}
	index := BuildFirstIndex(values)
	position, found := Lookup(index, 8)
	missing, missingFound := Lookup(index, 9)
	return fmt.Sprintf("index(8)=%d found=%t; index(9)=%d found=%t; count(8)=%d; duplicate=%t",
		position, found, missing, missingFound, Frequencies(values)[8], ContainsDuplicate(values))
}
