// Package complexity содержит маленькие функции, на которых удобно считать шаги.
package complexity

import "fmt"

// Measurement хранит результат и число основных операций.
type Measurement struct {
	Value int
	Steps int
}

// At возвращает элемент по индексу или 0, false при неверной границе.
// Время и дополнительная память — O(1).
func At(values []int, index int) (int, bool) {
	if index < 0 || index >= len(values) {
		return 0, false
	}
	return values[index], true
}

// Sum возвращает сумму и число сложений. Время O(n), память O(1).
func Sum(values []int) Measurement {
	measurement := Measurement{}
	for _, value := range values {
		measurement.Value += value
		measurement.Steps++
	}
	return measurement
}

// HalvingSteps считает деления с округлением вверх до n <= 1.
// Для n <= 1 результат 0; время O(log n), память O(1).
func HalvingSteps(n int) int {
	steps := 0
	for n > 1 {
		n = n/2 + n%2 // Округление вверх без переполнения n+1 при MaxInt.
		steps++
	}
	return steps
}

// PairComparisons считает пары индексов first < second.
// Результат n*(n-1)/2; время O(n²), память O(1).
func PairComparisons(values []int) int {
	steps := 0
	for first := 0; first < len(values); first++ {
		for second := first + 1; second < len(values); second++ {
			steps++
		}
	}
	return steps
}

// Example возвращает короткий предсказуемый результат для cmd и integration-теста.
func Example() string {
	values := []int{10, 20, 30, 40}
	last, _ := At(values, 3)
	sum := Sum(values)
	return fmt.Sprintf("O(1): value=%d; O(n): sum=%d steps=%d; O(log n): steps=%d; O(n²): pairs=%d",
		last, sum.Value, sum.Steps, HalvingSteps(16), PairComparisons(values))
}
