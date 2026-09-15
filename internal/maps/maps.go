// Package maps показывает, когда индекс на основе map ускоряет повторный поиск.
package maps

import "fmt"

// BuildFirstIndex один раз обходит срез и запоминает первый индекс каждого числа.
// Построение занимает O(n) ожидаемого времени и O(n) дополнительной памяти.
//
// Чтение готовой map обычно рассматривают как O(1) в среднем. Спецификация Go
// описывает поведение map, но не обещает конкретную асимптотическую сложность,
// поэтому в учебном материале мы не называем это гарантированным худшим случаем.
//
// Как не нужно: если нужен всего один поиск, строить большую map может быть
// дороже линейного прохода. Индекс особенно полезен для многих запросов.
func BuildFirstIndex(values []int) map[int]int {
	index := make(map[int]int, len(values))
	for position, value := range values {
		if _, exists := index[value]; !exists {
			index[value] = position
		}
	}
	return index
}

// Lookup читает индекс из готовой map и возвращает -1, если ключа нет.
//
// Как не нужно: выражение index[target] без проверки exists не отличает
// отсутствующий ключ от ключа, которому действительно соответствует индекс 0.
func Lookup(index map[int]int, target int) (int, bool) {
	position, exists := index[target]
	if !exists {
		return -1, false
	}
	return position, true
}

// Frequencies считает количество появлений каждого значения за один проход.
// Время — O(n) в среднем для map-операций, дополнительная память — O(k),
// где k — число различных значений и k <= n.
func Frequencies(values []int) map[int]int {
	frequencies := make(map[int]int)
	for _, value := range values {
		frequencies[value]++
	}
	return frequencies
}

// ContainsDuplicate завершает работу при первом повторе.
// Лучший случай — O(1), худший ожидаемый для полного прохода — O(n), память O(n).
//
// Как не нужно: вложенные циклы тоже найдут повтор, но сделают до O(n²)
// сравнений. Это полезный исходный вариант для обсуждения, но не лучший итог.
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
