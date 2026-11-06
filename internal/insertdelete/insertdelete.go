// Package insertdelete показывает, какие элементы сдвигаются при изменении среза.
package insertdelete

import "fmt"

// InsertAt вставляет value перед index и возвращает новый заголовок среза.
// Допустимы 0 <= index <= len(values); неверный индекс не меняет вход.
// Сдвиг — O(n); рост массива может потребовать O(n) новой памяти.
func InsertAt[T any](values []T, index int, value T) ([]T, bool) {
	if index < 0 || index > len(values) {
		return values, false
	}

	values = append(values, value)
	copy(values[index+1:], values[index:len(values)-1])
	values[index] = value
	return values, true
}

// DeleteAt удаляет элемент с сохранением порядка и очищает старый хвост.
// Очистка освобождает ссылки из невидимой части массива.
// Неверный индекс не меняет вход; время O(n-index), память O(1).
func DeleteAt[T any](values []T, index int) ([]T, bool) {
	if index < 0 || index >= len(values) {
		return values, false
	}

	copy(values[index:], values[index+1:])
	var zero T
	values[len(values)-1] = zero
	return values[:len(values)-1], true
}

// DeleteUnordered заменяет удалённый элемент последним и очищает хвост.
// Порядок может измениться. Время и дополнительная память O(1).
func DeleteUnordered[T any](values []T, index int) ([]T, bool) {
	if index < 0 || index >= len(values) {
		return values, false
	}

	last := len(values) - 1
	values[index] = values[last]
	var zero T
	values[last] = zero
	return values[:last], true
}

// Put вставляет или обновляет ключ. Запись в nil-map вызывает panic.
func Put(values map[string]int, key string, value int) {
	values[key] = value
}

// Remove удаляет ключ. delete безопасен и для отсутствующего ключа, и для nil-map.
func Remove(values map[string]int, key string) {
	delete(values, key)
}

// Example возвращает предсказуемое состояние после четырёх операций.
func Example() string {
	values := []int{10, 20, 30}
	values, _ = InsertAt(values, 1, 15)
	values, _ = DeleteAt(values, 2)
	values, _ = DeleteUnordered(values, 0)

	ages := make(map[string]int)
	Put(ages, "Аня", 20)
	Put(ages, "Аня", 21)
	Remove(ages, "нет")
	return fmt.Sprintf("slice=%v; map[Аня]=%d", values, ages["Аня"])
}
