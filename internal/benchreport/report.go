package benchreport

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

// FormatTime меняет только единицу, не умножая ns/op на размер входа.
func FormatTime(ns float64) string {
	units := []string{"ns", "µs", "ms", "s"}
	i := 0
	for ns >= 1000 && i < len(units)-1 {
		ns /= 1000
		i++
	}
	return number(ns) + " " + units[i]
}

// FormatBytes использует двоичные единицы. Это allocation на op, если значение
// взято из B/op; тот же formatter не превращает его в retained heap.
func FormatBytes(bytes float64) string {
	units := []string{"B", "KiB", "MiB", "GiB"}
	i := 0
	for bytes >= 1024 && i < len(units)-1 {
		bytes /= 1024
		i++
	}
	return number(bytes) + " " + units[i]
}

func number(value float64) string { return strconv.FormatFloat(value, 'g', 4, 64) }

// IntBytes выводится из архитектуры самого лога, а не машины, читающей файл.
// Неизвестная архитектура не получает выдуманной оценки размера данных.
func IntBytes(arch string) int {
	switch arch {
	case "386", "arm", "mips", "mipsle", "amd64p32":
		return 4
	case "amd64", "arm64", "ppc64", "ppc64le", "mips64", "mips64le", "s390x", "riscv64", "loong64", "sparc64", "wasm":
		return 8
	default:
		return 0
	}
}

// Pattern ограничивает измерения до запуска testing.B, а не отбрасывает
// полученные результаты. n=0 означает все предусмотренные лабораторией размеры.
func Pattern(kind string, n int, scenario string) (string, error) {
	if n < 0 {
		return "", fmt.Errorf("n должен быть >= 0")
	}
	roots := map[string]string{"search": "BenchmarkSearch", "search-build": "BenchmarkSearchBuild", "search-queries": "BenchmarkSearchQueries", "sort": "BenchmarkSort", "structures": "BenchmarkStructures", "memory": "BenchmarkMemory"}
	root, ok := roots[kind]
	if !ok {
		return "", fmt.Errorf("неизвестная группа %q", kind)
	}
	if scenario != "" && !(kind == "search" && oneOf(scenario, "first", "middle", "last", "missing", "duplicates") || kind == "sort" && oneOf(scenario, "random", "sorted", "reverse", "nearly", "equal", "duplicates")) {
		return "", fmt.Errorf("сценарий %q не подходит группе %s", scenario, kind)
	}
	size := "n=[0-9]+$"
	if n > 0 {
		size = fmt.Sprintf("n=%d$", n)
	}
	shape := ""
	switch kind {
	case "search":
		if scenario == "" {
			scenario = ".*"
		}
		shape = ".*/" + scenario + "$/" + size
	case "sort":
		if scenario == "" {
			scenario = ".*"
		}
		shape = ".*/.*/" + scenario + "$/" + size
	case "search-queries":
		shape = ".*/" + size + "/q=[0-9]+$"
	case "structures":
		// У операций разная глубина имени; n фильтруется внутри Go benchreport
		// через -bench для каждого из двух вариантов отдельно (см. SizeFilter).
		shape = ".*"
	default:
		shape = ".*/" + size
	}
	return "^" + root + "$/" + shape, nil
}

// SizeFilter передаётся самим benchmark через окружение только для структуры
// с неодинаковой глубиной b.Run. Пустое значение сохраняет исходную матрицу.
func SizeFilter(n int) string {
	if n == 0 {
		return ""
	}
	return strconv.Itoa(n)
}

type reportGroup struct {
	Environment
	Root, Mode, Scenario string
	N, Q, Procs          int
}

func groupOf(row Row) reportGroup {
	n := row.Name
	return reportGroup{row.Environment, n.Root, n.Mode, n.Scenario, n.N, n.Q, n.Procs}
}

func title(g reportGroup) string {
	switch g.Root {
	case "BenchmarkSearch":
		return "ПОИСК В ГОТОВЫХ ДАННЫХ · " + scenarioLabel(g.Scenario)
	case "BenchmarkSearchBuild":
		return "ПОДГОТОВКА СТРУКТУР · построение / копия + сортировка"
	case "BenchmarkSearchQueries":
		return fmt.Sprintf("ПОДГОТОВКА + %d ЗАПРОСОВ · время всей серии", g.Q)
	case "BenchmarkSort":
		mode := "только сортировка"
		if g.Mode == "copy-sort" {
			mode = "копирование + сортировка (готовый буфер)"
		}
		return "СОРТИРОВКИ · " + g.Scenario + " · " + mode
	case "BenchmarkMemory":
		return "ALLOCATION ПРИ ПОСТРОЕНИИ · не retained heap"
	default:
		labels := map[string]string{"value-lookup": "поиск по значению / ключу", "known-index": "чтение по известному индексу", "ordered-insert": "вставка с сохранением порядка", "ordered-delete": "удаление с сохранением порядка", "unordered-delete": "удаление с изменением порядка", "known-index-update": "обновление по известному индексу", "update-by-key": "обновление значения по ключу", "key-insert": "вставка по ключу без порядка", "key-delete": "удаление по ключу без порядка"}
		return "SLICE / MAP · " + labels[g.Mode]
	}
}

func scenarioLabel(scenario string) string {
	return map[string]string{"first": "элемент в начале", "middle": "конкретный элемент из середины", "last": "последний элемент", "missing": "элемент отсутствует", "duplicates": "повторяющиеся значения"}[scenario]
}

func algorithmLabel(algorithm string) string {
	labels := map[string]string{"linear": "Linear Search", "binary": "Binary Search", "map": "Map Lookup", "slices-binary": "slices.BinarySearch", "bubble": "Bubble Sort", "selection": "Selection Sort", "insertion": "Insertion Sort", "merge": "Merge Sort", "quick": "Quick (last pivot)", "slices": "slices.Sort", "sort-binary": "Копия + Sort + Binary", "map-build": "Построение map", "clone-sort": "Новая копия + Sort", "slice-build": "Построение slice", "map-index-build": "Построение map-индекса"}
	if label, ok := labels[algorithm]; ok {
		return label
	}
	return algorithm
}

// theory показывает верхнюю границу времени, а не число реально сделанных
// сравнений. Для map ожидаемая оценка подписывается отдельно от гарантии.
func theory(n Name) (time, memory string) {
	if n.Root == "BenchmarkSearchQueries" {
		switch n.Algorithm {
		case "linear":
			return "O(q·n)", "O(1)"
		case "map":
			return "O(n+q) ожид.", "O(n) индекс"
		default:
			return "O(n log n+q log n)", "O(n) копия"
		}
	}
	if n.Root == "BenchmarkSearchBuild" || n.Root == "BenchmarkMemory" {
		if n.Algorithm == "clone-sort" {
			return "O(n log n)", "O(n) копия"
		}
		if n.Algorithm == "slice-build" {
			return "O(n)", "O(n) результат"
		}
		return "O(n) ожид.", "O(n) индекс"
	}
	if n.Root == "BenchmarkSort" {
		switch n.Algorithm {
		case "merge":
			return "O(n log n)", "O(n) + стек O(log n)"
		case "quick":
			return "O(n²) худш.", "стек: O(log n)..O(n)"
		case "slices":
			return "O(n log n) худш.¹", "стек O(log n)¹"
		default:
			return "O(n²) худш.", "O(1)"
		}
	}
	if strings.HasPrefix(n.Algorithm, "map") || n.Algorithm == "map" {
		return "O(1) ожид.", "O(1), индекс готов"
	}
	if n.Algorithm == "binary" || n.Algorithm == "slices-binary" {
		return "O(log n) худш.", "O(1)"
	}
	if n.Algorithm == "linear" || n.Algorithm == "slice-linear" || oneOf(n.Algorithm, "slice-insert-ordered", "slice-delete-ordered", "slice-update-by-key") {
		return "O(n) худш.", "O(1), буфер готов"
	}
	return "O(1)", "O(1), буфер готов"
}

// Render — Markdown-таблицы, читаемые также в терминале. compact скрывает
// только дополнительные формы входа; CSV всегда содержит все измеренные строки.
func Render(rows []Row, view string) (string, error) {
	if len(rows) == 0 || !oneOf(view, "compact", "full") {
		return "", fmt.Errorf("нет строк или неизвестный вид отчёта %q", view)
	}
	var out strings.Builder
	fmt.Fprintln(&out, "# Лаборатория benchmark — реальные измерения")
	fmt.Fprintln(&out, "Среднее арифметическое по прогонам; k — число прогонов (-count), не b.N.")
	fmt.Fprintln(&out, "n — размер входа. Первая числовая колонка raw — b.N: вызовы внутри одного прогона.")
	fmt.Fprintln(&out, "Время / op — один полный вызов (либо вся серия q запросов); на n его НЕ умножают.")
	fmt.Fprintln(&out, "Heap / op — выделенные байты B/op; allocs/op — число выделений. Это НЕ вся память.")
	fmt.Fprintln(&out, "0 B/op не учитывает готовый вход, буферы вне таймера и стек; Go округляет allocation/op вниз.")
	fmt.Fprintln(&out, "Теория объясняет рост, benchmark не доказывает Big O. middle не означает средний случай.")
	preferred := map[string]bool{}
	for _, row := range rows {
		if row.Name.Root == "BenchmarkSearch" && row.Name.Scenario == "last" || row.Name.Root == "BenchmarkSort" && row.Name.Scenario == "random" {
			preferred[row.Name.Root] = true
		}
	}
	var groups []reportGroup
	groupRows := map[reportGroup][]Row{}
	hidden := 0
	for _, row := range rows {
		n := row.Name
		if view == "compact" && preferred[n.Root] && (n.Root == "BenchmarkSearch" && n.Scenario != "last" || n.Root == "BenchmarkSort" && n.Scenario != "random") {
			hidden++
			continue
		}
		g := groupOf(row)
		if _, exists := groupRows[g]; !exists {
			groups = append(groups, g)
		}
		groupRows[g] = append(groupRows[g], row)
	}
	if hidden > 0 {
		fmt.Fprintf(&out, "Компактный вид: скрыто %d агрегированных строк других сценариев. Все — в summary.csv; BENCHVIEW=full покажет их в таблицах.\n", hidden)
	}
	// Aggregate уже сортирует числовые размеры. Стабильный порядок групп также
	// сохраняется при передаче заранее агрегированных строк из импортера.
	sort.SliceStable(groups, func(i, j int) bool {
		a, b := groups[i], groups[j]
		if a.Environment != b.Environment {
			return environmentKey(a.Environment) < environmentKey(b.Environment)
		}
		if a.Procs != b.Procs {
			return a.Procs < b.Procs
		}
		if a.Root != b.Root {
			return a.Root < b.Root
		}
		if a.N != b.N {
			return a.N < b.N
		}
		if a.Mode != b.Mode {
			return a.Mode < b.Mode
		}
		if a.Scenario != b.Scenario {
			return a.Scenario < b.Scenario
		}
		return a.Q < b.Q
	})
	for _, g := range groups {
		fmt.Fprintf(&out, "\n## %s — n = %d\n", title(g), g.N)
		cpu := g.CPU
		if cpu == "" {
			cpu = "модель не указана"
		}
		fmt.Fprintf(&out, "%s/%s · %s · GOMAXPROCS=%d\n", g.GOOS, g.GOARCH, cpu, g.Procs)
		if width := IntBytes(g.GOARCH); width != 0 {
			fmt.Fprintf(&out, "Вход []int: ≈ %s (%d × %d B), без header, запаса cap и map-индекса.\n", FormatBytes(float64(g.N)*float64(width)), g.N, width)
		} else {
			fmt.Fprintln(&out, "Размер int для архитектуры лога неизвестен; оценка памяти входа не рассчитана.")
		}
		cells := [][]string{{"Алгоритм / операция", "Среднее / op", "Heap / op", "allocs/op", "k", "Время: теория", "Доп. память"}}
		measured := map[string]bool{}
		for _, row := range groupRows[g] {
			measured[row.Name.Algorithm] = true
			time, memory := theory(row.Name)
			cells = append(cells, []string{algorithmLabel(row.Name.Algorithm), FormatTime(row.MeanNS), FormatBytes(row.MeanBytes), number(row.MeanAllocs), strconv.Itoa(row.Repeats), time, memory})
		}
		if g.Root == "BenchmarkSort" && g.N > 512 {
			for _, algorithm := range []string{"bubble", "selection", "insertion", "quick"} {
				if measured[algorithm] {
					continue
				}
				cells = append(cells, []string{algorithmLabel(algorithm), "не измерялось", "—", "—", "—", "ограничение размера", "n ≤ 512"})
			}
		}
		out.WriteString(table(cells))
		// Сравнение только одинакового контракта. Вставка в упорядоченный slice
		// и вставка ключа в map намеренно никогда не выбирают общего победителя.
		comparable := g.Root == "BenchmarkSearch" || g.Root == "BenchmarkSearchQueries" || g.Root == "BenchmarkSort" || g.Root == "BenchmarkStructures" && oneOf(g.Mode, "value-lookup", "update-by-key")
		list := groupRows[g]
		if comparable && len(list) > 1 {
			best, worst := list[0], list[0]
			for _, row := range list[1:] {
				if row.MeanNS < best.MeanNS {
					best = row
				}
				if row.MeanNS > worst.MeanNS {
					worst = row
				}
			}
			if best.MeanNS > 0 {
				fmt.Fprintf(&out, "Минимальное среднее здесь: %s; отношение max/min = %s×. Это не оценка статистической значимости.\n", algorithmLabel(best.Name.Algorithm), number(worst.MeanNS/best.MeanNS))
			}
		}
	}
	fmt.Fprintln(&out, "\nMap lookup требует готового O(n)-индекса; ожидаемое O(1) не является гарантией худшего случая.")
	fmt.Fprintln(&out, "У Binary лучший случай O(1); slices.BinarySearch ищет нижнюю границу и продолжает деление диапазона.")
	fmt.Fprintln(&out, "Bubble с ранним выходом и Insertion: лучший O(n); Selection — O(n²) сравнений даже на sorted.")
	fmt.Fprintln(&out, "Quick с последним pivot: обычно O(n log n), sorted/reverse/equal дают O(n²) и стек O(n).")
	fmt.Fprintln(&out, "¹ Оценки slices.Sort относятся к реализации Go 1.22; версию фактического запуска см. в environment.json.")
	return out.String(), nil
}

func table(cells [][]string) string {
	widths := make([]int, len(cells[0]))
	for i := range widths {
		widths[i] = 3 // Не меньше трёх дефисов в разделителе Markdown-таблицы.
	}
	for _, row := range cells {
		for i, cell := range row {
			if n := utf8.RuneCountInString(cell); n > widths[i] {
				widths[i] = n
			}
		}
	}
	var result strings.Builder
	for j, row := range cells {
		result.WriteString("| ")
		for i, cell := range row {
			result.WriteString(cell)
			result.WriteString(strings.Repeat(" ", widths[i]-utf8.RuneCountInString(cell)))
			result.WriteString(" | ")
		}
		result.WriteByte('\n')
		if j == 0 {
			result.WriteString("| ")
			for _, width := range widths {
				result.WriteString(strings.Repeat("-", width))
				result.WriteString(" | ")
			}
			result.WriteByte('\n')
		}
	}
	return result.String()
}
