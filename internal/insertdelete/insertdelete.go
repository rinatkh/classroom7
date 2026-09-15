// Package insertdelete показывает, какие элементы сдвигаются при изменении среза.
package insertdelete

import "fmt"

// InsertAt вставляет value перед index и возвращает новый заголовок среза.
// Допустимы индексы от 0 до len(values) включительно.
//
// Что происходит:
//  1. append увеличивает длину на одну ячейку; иногда создаётся новый массив.
//  2. copy сдвигает хвост вправо.
//  3. В освободившееся место записывается value.
//
// В худшем случае сдвигается n элементов: время O(n). Если массив приходится
// увеличить, временно нужна новая область O(n); при свободной capacity — O(1).
//
// Как не нужно: после append нужно использовать возвращённый срез. Старый
// заголовок может по-прежнему указывать на массив прежней длины.
func InsertAt[T any](values []T, index int, value T) ([]T, bool) {
	if index < 0 || index > len(values) {
		return values, false
	}

	values = append(values, value)
	copy(values[index+1:], values[index:len(values)-1])
	values[index] = value
	return values, true
}

// DeleteAt удаляет элемент с сохранением порядка остальных значений.
//
// Хвост сдвигается влево, поэтому худшее время O(n). Последнюю ячейку старой
// длины обнуляем: если T содержит указатель, удалённый объект не удерживается
// невидимой частью базового массива. Дополнительная память — O(1).
//
// Как не нужно: выражение append(values[:i], values[i+1:]...) работает для
// удаления, но без обнуления хвоста может дольше удерживать ссылочный объект.
func DeleteAt[T any](values []T, index int) ([]T, bool) {
	if index < 0 || index >= len(values) {
		return values, false
	}

	copy(values[index:], values[index+1:])
	var zero T
	values[len(values)-1] = zero
	return values[:len(values)-1], true
}

// DeleteUnordered удаляет элемент за O(1), если порядок не важен.
// На место удаляемого элемента ставится последний, затем длина уменьшается.
//
// Как не нужно: нельзя выбирать этот вариант, когда порядок значений является
// частью контракта функции. Быстрое удаление намеренно меняет порядок.
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

// Put записывает новое значение по ключу или заменяет старое.
// Для готовой ненулевой map операция обычно рассматривается как O(1) в среднем.
//
// Как не нужно: запись в nil-map вызывает panic. Её нужно создать через make
// или литерал map до первой записи.
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
