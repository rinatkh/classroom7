package benchmarkdebug

import (
	"fmt"
	"runtime"
	"slices"
	"testing"

	"github.com/rinatkh/classroom7/internal/dataset"
	"github.com/rinatkh/classroom7/internal/insertdelete"
	algomaps "github.com/rinatkh/classroom7/internal/maps"
	"github.com/rinatkh/classroom7/internal/mergesort"
	"github.com/rinatkh/classroom7/internal/quicksort"
	"github.com/rinatkh/classroom7/internal/search"
	"github.com/rinatkh/classroom7/internal/simplesort"
)

// Глобальные приёмники сохраняют наблюдаемый результат измеряемой работы.
var (
	searchSink bool
	indexSink  int
	sortSink   int
	sliceSink  []int
	mapSink    map[int]int
)

// BenchmarkSearch измеряет готовый поиск на прогретых, повторяемых целях.
// Это не модель всех реальных запросов; смешанная серия есть ниже отдельно.
// При повторах сравниваем наличие: BinarySearch может вернуть любой индекс,
// map хранит первый исходный индекс, slices.BinarySearch — нижнюю границу.
func BenchmarkSearch(b *testing.B) {
	for _, n := range []int{100, 1000, 10000, 100000, 1000000} {
		unique := make([]int, n)
		for i := range unique {
			unique[i] = i * 2
		}
		for _, scenario := range []string{"first", "middle", "last", "missing", "duplicates"} {
			values := unique
			target := unique[n/2]
			switch scenario {
			case "first":
				target = unique[0]
			case "last":
				target = unique[n-1]
			case "missing":
				target = -1
			case "duplicates":
				values = slices.Clone(unique)
				for i := range values {
					values[i] = (i / 4) * 2
				}
				target = values[n/2]
			}
			index := algomaps.BuildFirstIndex(values)
			for _, algorithm := range []string{"linear", "binary", "map", "slices-binary"} {
				b.Run(fmt.Sprintf("%s/%s/n=%d", algorithm, scenario, n), func(b *testing.B) {
					b.ReportAllocs()
					b.ResetTimer()
					switch algorithm {
					case "linear":
						for i := 0; i < b.N; i++ {
							r := search.LinearSearch(values, target)
							indexSink, searchSink = r.Index, r.Found
						}
					case "binary":
						for i := 0; i < b.N; i++ {
							r := search.BinarySearch(values, target)
							indexSink, searchSink = r.Index, r.Found
						}
					case "map":
						for i := 0; i < b.N; i++ {
							indexSink, searchSink = algomaps.Lookup(index, target)
						}
					case "slices-binary":
						for i := 0; i < b.N; i++ {
							indexSink, searchSink = slices.BinarySearch(values, target)
						}
					}
					b.StopTimer()
					if searchSink != (scenario != "missing") {
						b.Fatal("неверный результат поиска")
					}
					runtime.KeepAlive(values)
					runtime.KeepAlive(index)
				})
			}
		}
	}
}

// BenchmarkSearchBuild включает выделение памяти и всю подготовку индекса.
// Возвращаемый индекс после сортировки относится к копии, не к исходному порядку.
func BenchmarkSearchBuild(b *testing.B) {
	for _, n := range []int{100, 1000, 10000, 100000, 1000000} {
		source := dataset.GenerateRandomSlice(n)
		b.Run(fmt.Sprintf("map-build/n=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				mapSink = algomaps.BuildFirstIndex(source)
			}
			b.StopTimer()
			if len(mapSink) == 0 {
				b.Fatal("пустой индекс")
			}
		})
		b.Run(fmt.Sprintf("clone-sort/n=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				sliceSink = slices.Clone(source)
				slices.Sort(sliceSink)
			}
			b.StopTimer()
			if !slices.IsSorted(sliceSink) {
				b.Fatal("не отсортировано")
			}
		})
	}
}

// BenchmarkSearchQueries измеряет ВСЮ задачу: подготовка + q смешанных запросов.
// Сравниваем только наличие, поскольку сортировка меняет исходные индексы.
func BenchmarkSearchQueries(b *testing.B) {
	for _, n := range []int{100, 10000} {
		source := dataset.GenerateRandomSlice(n)
		targets := []int{source[0], source[n/2], source[n-1], 2000001}
		for _, q := range []int{1, 100, 1000} {
			for _, algorithm := range []string{"linear", "sort-binary", "map"} {
				b.Run(fmt.Sprintf("%s/n=%d/q=%d", algorithm, n, q), func(b *testing.B) {
					b.ReportAllocs()
					b.ResetTimer()
					for i := 0; i < b.N; i++ {
						count := 0
						switch algorithm {
						case "linear":
							for j := 0; j < q; j++ {
								if search.LinearSearch(source, targets[j%4]).Found {
									count++
								}
							}
						case "sort-binary":
							work := slices.Clone(source)
							slices.Sort(work)
							for j := 0; j < q; j++ {
								if search.BinarySearch(work, targets[j%4]).Found {
									count++
								}
							}
						case "map":
							index := algomaps.BuildFirstIndex(source)
							for j := 0; j < q; j++ {
								if _, ok := index[targets[j%4]]; ok {
									count++
								}
							}
						}
						indexSink = count
					}
					b.StopTimer()
					if indexSink != q-q/4 {
						b.Fatal("разошлись результаты запросов")
					}
				})
			}
		}
	}
}

type sortAlgorithm struct {
	name string
	run  func([]int)
	maxN int
}

var sortAlgorithms = []sortAlgorithm{
	{"bubble", func(v []int) { simplesort.BubbleSort(v) }, 512},
	{"selection", func(v []int) { simplesort.SelectionSort(v) }, 512},
	{"insertion", func(v []int) { simplesort.InsertionSort(v) }, 512},
	{"merge", mergesort.Sort, 10000},
	// Последний pivot деградирует также на equal: поэтому лимит как у O(n²).
	{"quick", func(v []int) { quicksort.Sort(v) }, 512},
	{"slices", slices.Sort[[]int, int], 10000},
}

// BenchmarkSort: copy-sort — copy в готовый буфер входит в таймер (не allocation).
// sort-only — batch до 32 независимых буферов, copy между StopTimer/StartTimer.
// Обе группы включают внутренние allocation алгоритма (например, Merge).
func BenchmarkSort(b *testing.B) {
	for _, n := range []int{100, 512, 10000} {
		for _, scenario := range dataset.Names {
			source := dataset.Scenario(n, scenario)
			for _, algorithm := range sortAlgorithms {
				if n > algorithm.maxN {
					continue
				}
				for _, mode := range []string{"copy-sort", "sort-only"} {
					b.Run(fmt.Sprintf("%s/%s/%s/n=%d", mode, algorithm.name, scenario, n), func(b *testing.B) {
						b.ReportAllocs()
						if mode == "copy-sort" {
							work := make([]int, n)
							b.ResetTimer()
							for i := 0; i < b.N; i++ {
								copy(work, source)
								algorithm.run(work)
								sortSink = work[n/2]
							}
							b.StopTimer()
							if !slices.IsSorted(work) {
								b.Fatal("сортировка не завершена")
							}
							runtime.KeepAlive(work)
						} else {
							const batch = 32
							buffers := make([][]int, batch)
							for i := range buffers {
								buffers[i] = make([]int, n)
							}
							b.ResetTimer()
							b.StopTimer()
							for done := 0; done < b.N; {
								count := min(batch, b.N-done)
								for j := 0; j < count; j++ {
									copy(buffers[j], source)
								}
								b.StartTimer()
								for j := 0; j < count; j++ {
									algorithm.run(buffers[j])
									sortSink = buffers[j][n/2]
								}
								b.StopTimer()
								for j := 0; j < count; j++ {
									if !slices.IsSorted(buffers[j]) {
										b.Fatal("сортировка не завершена")
									}
								}
								done += count
							}
							runtime.KeepAlive(buffers)
						}
					})
				}
			}
		}
	}
}

// resetBatch исключает восстановление изменяемых структур из измеряемого участка.
// На одной итерации ровно одна mutation; размер не растёт вместе с b.N.
func resetBatch(b *testing.B, reset, operation func(int), check func(int) bool) {
	b.Helper()
	b.ReportAllocs()
	b.ResetTimer()
	b.StopTimer()
	for done := 0; done < b.N; {
		count := min(256, b.N-done)
		for j := 0; j < count; j++ {
			reset(j)
		}
		b.StartTimer()
		for j := 0; j < count; j++ {
			operation(j)
		}
		b.StopTimer()
		for j := 0; j < count; j++ {
			if !check(j) {
				b.Fatal("неверный результат mutation")
			}
		}
		done += count
	}
}

// BenchmarkStructures явно различает индекс среза, поиск значения, порядок и ключ.
// Вставка в середину среза сохраняет порядок; map вставляет новую пару по ключу.
func BenchmarkStructures(b *testing.B) {
	for _, n := range []int{100, 1000} {
		source := make([]int, n)
		for i := range source {
			source[i] = i
		}
		index := algomaps.BuildFirstIndex(source)
		b.Run(fmt.Sprintf("value-lookup/slice-linear/n=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				r := search.LinearSearch(source, n-1)
				indexSink, searchSink = r.Index, r.Found
			}
			b.StopTimer()
			if !searchSink || indexSink != n-1 {
				b.Fatal("lookup")
			}
		})
		b.Run(fmt.Sprintf("value-lookup/map-key/n=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				indexSink, searchSink = index[n-1]
			}
			b.StopTimer()
			if !searchSink || indexSink != n-1 {
				b.Fatal("lookup")
			}
		})
		// Контроль: известная позиция — другой контракт, не поиск по значению.
		b.Run(fmt.Sprintf("known-index/slice-read/n=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				indexSink = source[n-1]
			}
			b.StopTimer()
			if indexSink != n-1 {
				b.Fatal("read")
			}
		})
		for _, name := range []string{"slice-insert-ordered", "slice-delete-ordered", "slice-delete-unordered", "slice-update-known-index", "slice-update-by-key", "map-insert-key", "map-delete-key", "map-update-key"} {
			b.Run(fmt.Sprintf("%s/n=%d", name, n), func(b *testing.B) {
				type entry struct{ key, value int }
				records := make([][]entry, 256)
				buffers := make([][]int, 256)
				tables := make([]map[int]int, 256)
				for j := range buffers {
					buffers[j] = make([]int, n, n+1)
					if name == "slice-update-by-key" {
						records[j] = make([]entry, n)
					}
					tables[j] = make(map[int]int, n+1)
					for k := range source {
						tables[j][k] = k
					}
				}
				reset := func(j int) {
					buffers[j] = buffers[j][:n]
					copy(buffers[j], source)
					for k := range records[j] {
						records[j][k] = entry{k, k}
					}
					// Сохраняем длину map и восстанавливаем только затронутые ключи.
					delete(tables[j], n)
					tables[j][n/2] = n / 2
					tables[j][n-1] = n - 1
				}
				operation := func(j int) {
					switch name {
					case "slice-insert-ordered":
						buffers[j], _ = insertdelete.InsertAt(buffers[j], n/2, n)
					case "slice-delete-ordered":
						buffers[j], _ = insertdelete.DeleteAt(buffers[j], n/2)
					case "slice-delete-unordered":
						buffers[j], _ = insertdelete.DeleteUnordered(buffers[j], n/2)
					case "slice-update-known-index":
						buffers[j][n-1] = -1
					case "slice-update-by-key":
						// Ищем запись по key и меняем только value: как map-update-key.
						for k := range records[j] {
							if records[j][k].key == n-1 {
								records[j][k].value = -1
								break
							}
						}
					case "map-insert-key":
						tables[j][n] = n
					case "map-delete-key":
						delete(tables[j], n/2)
					case "map-update-key":
						tables[j][n-1] = -1
					}
				}
				check := func(j int) bool {
					switch name {
					case "slice-insert-ordered":
						return len(buffers[j]) == n+1 && buffers[j][n/2] == n
					case "slice-delete-ordered":
						return len(buffers[j]) == n-1 && buffers[j][n/2] == n/2+1
					case "slice-delete-unordered":
						return len(buffers[j]) == n-1 && buffers[j][n/2] == n-1
					case "slice-update-known-index":
						return buffers[j][n-1] == -1
					case "slice-update-by-key":
						return records[j][n-1].key == n-1 && records[j][n-1].value == -1
					case "map-insert-key":
						return len(tables[j]) == n+1 && tables[j][n] == n
					case "map-delete-key":
						_, ok := tables[j][n/2]
						return !ok && len(tables[j]) == n-1
					default:
						return tables[j][n-1] == -1
					}
				}
				resetBatch(b, reset, operation, check)
				runtime.KeepAlive(buffers)
				runtime.KeepAlive(tables)
				runtime.KeepAlive(records)
			})
		}
	}
}

// BenchmarkMemory считает allocation при построении, НЕ полный retained heap.
// Буфер структуры уходит в глобальный sink; GC может освободить предыдущий.
func BenchmarkMemory(b *testing.B) {
	for _, n := range []int{100, 1000, 10000, 100000} {
		source := dataset.GenerateRandomSlice(n)
		b.Run(fmt.Sprintf("slice-build/n=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				sliceSink = slices.Clone(source)
			}
			b.StopTimer()
			if !slices.Equal(sliceSink, source) {
				b.Fatal("clone")
			}
		})
		b.Run(fmt.Sprintf("map-index-build/n=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				mapSink = algomaps.BuildFirstIndex(source)
			}
			b.StopTimer()
			if len(mapSink) == 0 {
				b.Fatal("build")
			}
		})
	}
}
