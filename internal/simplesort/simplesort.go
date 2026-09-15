// Package simplesort содержит учебные реализации трёх сортировок O(n²).
package simplesort

import "fmt"

// Report показывает, сколько сравнений и записей в срез выполнил алгоритм.
// Это учебные счётчики, а не точное число процессорных инструкций.
type Report struct {
	Comparisons int
	Writes      int
}

// BubbleSort несколько раз проходит слева направо и меняет соседей местами,
// если они стоят в неправильном порядке. После каждого прохода один из самых
// больших элементов оказывается справа на окончательной позиции.
//
// Худшее и среднее время — O(n²), дополнительная память — O(1).
// С флагом swapped уже отсортированный вход завершается за O(n).
//
// Как не нужно: не используйте пузырьковую сортировку для больших рабочих
// данных только потому, что её код короткий. Здесь она нужна, чтобы увидеть шаги.
func BubbleSort(values []int) Report {
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
		if !swapped {
			break
		}
	}
	return report
}

// BubbleTrace возвращает копию среза после каждого полного прохода.
// Копии намеренно требуют O(n²) памяти в худшем случае: это цена визуализации,
// а не память самого BubbleSort.
func BubbleTrace(values []int) [][]int {
	states := [][]int{append([]int(nil), values...)}
	for end := len(values) - 1; end > 0; end-- {
		swapped := false
		for index := 0; index < end; index++ {
			if values[index] > values[index+1] {
				values[index], values[index+1] = values[index+1], values[index]
				swapped = true
			}
		}
		states = append(states, append([]int(nil), values...))
		if !swapped {
			break
		}
	}
	return states
}

// SelectionSort на каждой позиции ищет минимальное значение в правой части
// и ставит его на текущую позицию. Сравнений всегда порядка n², даже если вход
// уже отсортирован. Дополнительная память — O(1).
//
// Как не нужно: выбор минимального элемента не делает алгоритм O(n), потому что
// для каждой позиции оставшаяся часть просматривается заново.
func SelectionSort(values []int) Report {
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
	}
	return report
}

// InsertionSort поддерживает отсортированную левую часть и вставляет в неё
// следующий элемент. Для почти отсортированных маленьких данных работает хорошо.
// Худшее время — O(n²), лучший случай — O(n), дополнительная память — O(1).
//
// Как не нужно: после сдвигов обязательно записать current в найденную позицию.
// Иначе значение будет потеряно, а один элемент продублируется.
func InsertionSort(values []int) Report {
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
