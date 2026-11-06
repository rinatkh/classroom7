// Package dataset создаёт одинаковые воспроизводимые входы для тестов и измерений.
package dataset

import (
	"math/rand"
	"slices"
)

// Seed фиксирует последовательность, чтобы повторные запуски сравнивали один вход.
const Seed int64 = 7

// GenerateRandomSlice возвращает n псевдослучайных int из [-1000000, 1000000].
// Для n <= 0 возвращается nil. Каждый вызов независим; глобального rand.Seed нет.
func GenerateRandomSlice(n int) []int { return GenerateRandomSliceWithSeed(n, Seed) }

// GenerateRandomSliceWithSeed позволяет отдельно исследовать влияние другого входа.
func GenerateRandomSliceWithSeed(n int, seed int64) []int {
	if n <= 0 {
		return nil
	}
	random := rand.New(rand.NewSource(seed))
	values := make([]int, n)
	for i := range values {
		values[i] = random.Intn(2000001) - 1000000
	}
	return values
}

// Names задаёт порядок шести сценариев в таблицах benchmark.
var Names = []string{"random", "sorted", "reverse", "nearly", "equal", "duplicates"}

// Scenario возвращает новый независимый срез. Неизвестное имя — ошибка вызывающего.
// nearly: отсортированный вход с одним соседним обменом на каждые 20 позиций.
func Scenario(n int, name string) []int {
	values := GenerateRandomSlice(n)
	switch name {
	case "random":
	case "sorted":
		slices.Sort(values)
	case "reverse":
		slices.Sort(values)
		slices.Reverse(values)
	case "nearly":
		slices.Sort(values)
		for i := 0; i+1 < len(values); i += 20 {
			values[i], values[i+1] = values[i+1], values[i]
		}
	case "equal":
		for i := range values {
			values[i] = 7
		}
	case "duplicates":
		for i := range values {
			values[i] %= 8
		}
	default:
		panic("неизвестный сценарий: " + name)
	}
	return values
}
