package benchreport

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"go/types"
	"io"
	"strings"
	"testing"
)

func TestAggregateArithmeticMean(t *testing.T) {
	samples := []Sample{sample(t, "BenchmarkSearch/linear/last/n=100", 1), sample(t, "BenchmarkSearch/linear/last/n=100", 2), sample(t, "BenchmarkSearch/linear/last/n=100", 9)}
	for i := range samples {
		samples[i].Iterations = uint64(i*i*1000 + 1)
		samples[i].Bytes = float64(i * 100)
		samples[i].Allocs = float64(i)
	}
	rows, err := Aggregate(samples, 3)
	if err != nil || len(rows) != 1 {
		t.Fatalf("%+v, %v", rows, err)
	}
	r := rows[0]
	// Среднее 4, медиана 2. b.N умышленно разный и не является весом.
	if r.MeanNS != 4 || r.MedianNS != 2 || r.MinNS != 1 || r.MaxNS != 9 || r.MeanBytes != 100 || r.MeanAllocs != 1 || r.Repeats != 3 {
		t.Fatal(r)
	}
	even, err := Aggregate(samples[:2], 0)
	if err != nil || even[0].MedianNS != 1.5 {
		t.Fatalf("%+v, %v", even, err)
	}
	for _, test := range []struct {
		values []Sample
		count  int
	}{{nil, 3}, {samples, -1}, {samples[:2], 3}} {
		if _, err := Aggregate(test.values, test.count); err == nil {
			t.Fatal("accepted missing repetitions")
		}
	}
}

func TestAggregateDoesNotMixDimensions(t *testing.T) {
	base := sample(t, "BenchmarkSearch/linear/last/n=100", 1)
	samples := []Sample{base}
	changes := []func(*Sample){
		func(s *Sample) { s.Name.N = 1000 }, func(s *Sample) { s.Name.N = 512 },
		func(s *Sample) { s.Name.Scenario = "first" }, func(s *Sample) { s.Name.Procs = 8 },
		func(s *Sample) { s.Name.Q = 100 }, func(s *Sample) { s.Name.Mode = "other" },
		func(s *Sample) { s.Name.Root = "BenchmarkMemory" }, func(s *Sample) { s.Name.Algorithm = "binary" },
		func(s *Sample) { s.Environment.CPU = "Other CPU" }, func(s *Sample) { s.Environment.GOOS = "darwin" },
		func(s *Sample) { s.Environment.GOARCH = "arm64" }, func(s *Sample) { s.Environment.Package = "other/pkg" },
	}
	for _, change := range changes {
		copy := base
		change(&copy)
		samples = append(samples, copy)
	}
	rows, err := Aggregate(samples, 1)
	if err != nil || len(rows) != len(samples) {
		t.Fatalf("groups mixed: %+v, %v", rows, err)
	}
	var ns []int
	for _, r := range rows {
		if r.Environment == base.Environment && r.Name.Procs == 1 && r.Name.Root == base.Name.Root && r.Name.Scenario == "last" && r.Name.Mode == "" && r.Name.Q == 0 && r.Name.Algorithm == "linear" {
			ns = append(ns, r.Name.N)
		}
	}
	if fmt.Sprint(ns) != "[100 512 1000]" {
		t.Fatalf("not numeric order: %v", ns)
	}
}

type errorWriter struct{}

func (errorWriter) Write([]byte) (int, error) { return 0, errors.New("write error") }

func TestCSV(t *testing.T) {
	s := sample(t, "BenchmarkSearch/map/last/n=100-8", 1.136)
	s.Environment.CPU = "CPU, \"quoted\""
	rows, _ := Aggregate([]Sample{s}, 1)
	var out bytes.Buffer
	if err := WriteCSV(&out, rows); err != nil {
		t.Fatal(err)
	}
	records, err := csv.NewReader(&out).ReadAll()
	if err != nil || len(records) != 2 || len(records[0]) != 18 {
		t.Fatalf("%v %v", records, err)
	}
	if records[1][3] != s.Environment.CPU || records[1][4] != "8" || records[1][9] != "100" || records[1][12] != "1.136" {
		t.Fatal(records[1])
	}
	if err := WriteCSV(errorWriter{}, rows); err == nil {
		t.Fatal("lost flush error")
	}
	// Переполнение буфера csv.Writer обнаруживает ошибку и во время Write.
	rows[0].Environment.CPU = strings.Repeat("x", 10000)
	if err := WriteCSV(errorWriter{}, rows); err == nil {
		t.Fatal("lost record error")
	}
}

func TestUnits(t *testing.T) {
	for _, test := range []struct {
		value       float64
		time, bytes string
	}{
		{0, "0 ns", "0 B"}, {1.136, "1.136 ns", "1.136 B"}, {999, "999 ns", "999 B"},
		{1000, "1 µs", "1000 B"}, {1024, "1.024 µs", "1 KiB"},
		{1000000, "1 ms", "976.6 KiB"}, {1048576, "1.049 ms", "1 MiB"},
		{1000000000, "1 s", "953.7 MiB"}, {1073741824, "1.074 s", "1 GiB"},
		{0.00001, "1e-05 ns", "1e-05 B"},
	} {
		t.Run(fmt.Sprint(test.value), func(t *testing.T) {
			if got := FormatTime(test.value); got != test.time {
				t.Fatalf("%s want %s", got, test.time)
			}
			if got := FormatBytes(test.value); got != test.bytes {
				t.Fatalf("%s want %s", got, test.bytes)
			}
		})
	}
	for _, test := range []struct {
		arch  string
		bytes int
	}{{"amd64", 8}, {"arm64", 8}, {"386", 4}, {"arm", 4}, {"wasm", 8}, {"unknown", 0}} {
		if IntBytes(test.arch) != test.bytes {
			t.Fatal(test)
		}
	}
}

func TestIntBytesMatchesGoTypes(t *testing.T) {
	for _, arch := range []string{"386", "amd64", "amd64p32", "arm", "arm64", "loong64", "mips", "mipsle", "mips64", "mips64le", "ppc64", "ppc64le", "riscv64", "s390x", "sparc64", "wasm"} {
		t.Run(arch, func(t *testing.T) {
			want := types.SizesFor("gc", arch).Sizeof(types.Typ[types.Int])
			if int64(IntBytes(arch)) != want {
				t.Fatalf("%s: got %d, want %d", arch, IntBytes(arch), want)
			}
		})
	}
}

func TestPattern(t *testing.T) {
	for _, test := range []struct {
		kind, scenario string
		n              int
		want           string
	}{
		{"search", "", 0, "^BenchmarkSearch$/.*/.*$/n=[0-9]+$"},
		{"search", "last", 100, "^BenchmarkSearch$/.*/last$/n=100$"},
		{"sort", "random", 512, "^BenchmarkSort$/.*/.*/random$/n=512$"},
		{"sort", "", 0, "^BenchmarkSort$/.*/.*/.*$/n=[0-9]+$"},
		{"search-build", "", 100, "^BenchmarkSearchBuild$/.*/n=100$"},
		{"search-queries", "", 100, "^BenchmarkSearchQueries$/.*/n=100$/q=[0-9]+$"},
		{"structures", "", 0, "^BenchmarkStructures$/.*"},
		{"memory", "", 100, "^BenchmarkMemory$/.*/n=100$"},
	} {
		got, err := Pattern(test.kind, test.n, test.scenario)
		if err != nil || got != test.want {
			t.Fatalf("%+v => %q %v", test, got, err)
		}
	}
	for _, test := range []struct {
		kind, scenario string
		n              int
	}{{"other", "", 0}, {"search", "", -1}, {"search", "average", 0}, {"sort", "first", 0}, {"memory", "random", 0}} {
		if _, err := Pattern(test.kind, test.n, test.scenario); err == nil {
			t.Fatal(test)
		}
	}
	if SizeFilter(0) != "" || SizeFilter(100) != "100" {
		t.Fatal("size filter")
	}
}

func reportSamples(t *testing.T) []Sample {
	var samples []Sample
	names := []string{
		"BenchmarkSearch/linear/last/n=100", "BenchmarkSearch/binary/last/n=100", "BenchmarkSearch/map/last/n=100", "BenchmarkSearch/slices-binary/last/n=100",
		"BenchmarkSearch/linear/first/n=100", "BenchmarkSearch/linear/middle/n=100", "BenchmarkSearch/linear/missing/n=100", "BenchmarkSearch/linear/duplicates/n=100",
		"BenchmarkSearchBuild/map-build/n=100", "BenchmarkSearchBuild/clone-sort/n=100",
		"BenchmarkSearchQueries/linear/n=100/q=100", "BenchmarkSearchQueries/map/n=100/q=100", "BenchmarkSearchQueries/sort-binary/n=100/q=100",
		"BenchmarkMemory/slice-build/n=100", "BenchmarkMemory/map-index-build/n=100",
		"BenchmarkStructures/value-lookup/slice-linear/n=100", "BenchmarkStructures/value-lookup/map-key/n=100",
		"BenchmarkStructures/known-index/slice-read/n=100",
	}
	for _, operation := range []string{"slice-insert-ordered", "slice-delete-ordered", "slice-delete-unordered", "slice-update-known-index", "slice-update-by-key", "map-insert-key", "map-delete-key", "map-update-key"} {
		names = append(names, "BenchmarkStructures/"+operation+"/n=100")
	}
	for _, algorithm := range []string{"bubble", "selection", "insertion", "merge", "quick", "slices"} {
		names = append(names, "BenchmarkSort/sort-only/"+algorithm+"/random/n=100")
	}
	names = append(names, "BenchmarkSort/copy-sort/merge/random/n=10000", "BenchmarkSort/sort-only/merge/sorted/n=100")
	for i, name := range names {
		samples = append(samples, sample(t, name, float64(i+1)))
	}
	unknown := sample(t, "BenchmarkSearch/linear/last/n=1000-2", 0)
	unknown.Environment.GOARCH = "unknown"
	unknown.Environment.CPU = ""
	samples = append(samples, unknown)
	other := sample(t, "BenchmarkSearch/binary/last/n=1000-2", 0)
	other.Environment = unknown.Environment
	samples = append(samples, other)
	return samples
}

func TestRender(t *testing.T) {
	rows, err := Aggregate(reportSamples(t), 1)
	if err != nil {
		t.Fatal(err)
	}
	compact, err := Render(rows, "compact")
	if err != nil {
		t.Fatal(err)
	}
	full, err := Render(rows, "full")
	if err != nil {
		t.Fatal(err)
	}
	for _, part := range []string{"Среднее арифметическое", "b.N", "на n его НЕ умножают", "не retained heap", "100 × 8 B", "не измерялось", "ограничение размера", "стек: O(log n)..O(n)", "O(n) + стек O(log n)", "время всей серии", "GOMAXPROCS=2", "неизвестен", "скрыто 5", "max/min ="} {
		if !strings.Contains(compact, part) {
			t.Errorf("missing %q", part)
		}
	}
	if strings.Contains(compact, "## ПОИСК В ГОТОВЫХ ДАННЫХ · элемент в начале") || !strings.Contains(full, "## ПОИСК В ГОТОВЫХ ДАННЫХ · элемент в начале") || strings.Contains(full, "скрыто ") {
		t.Fatal("compact/full filtering")
	}
	if strings.Count(full, "| Binary Search") != 2 {
		t.Fatal("merged or missing algorithm rows")
	}
	// Теории не переносят время первого или среднего target на весь класс.
	for _, scenario := range []string{"first", "middle", "last", "missing", "duplicates"} {
		if !strings.Contains(full, scenarioLabel(scenario)) {
			t.Fatal(scenario)
		}
	}
	for _, args := range []struct {
		rows []Row
		view string
	}{{nil, "full"}, {rows, "bad"}} {
		if _, err := Render(args.rows, args.view); err == nil {
			t.Fatal("accepted invalid report")
		}
	}
}

func TestRenderNoInventedMeasurements(t *testing.T) {
	samples := []Sample{sample(t, "BenchmarkSort/sort-only/merge/sorted/n=100", 1)}
	rows, _ := Aggregate(samples, 1)
	report, err := Render(rows, "compact")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(report, "| Bubble Sort") || strings.Contains(report, "| slices.Sort") || !strings.Contains(report, "sorted") {
		t.Fatal(report)
	}
	// Если импортирован сторонний опыт за пределом учебного лимита,
	// его реальная строка не дублируется пометкой «не измерялось».
	rows, _ = Aggregate([]Sample{sample(t, "BenchmarkSort/sort-only/bubble/random/n=1000", 1)}, 1)
	report, err = Render(rows, "full")
	if err != nil || strings.Count(report, "| Bubble Sort") != 1 {
		t.Fatalf("%v %s", err, report)
	}
	// Без last выбранный first не исчезает из компактного вывода.
	rows, _ = Aggregate([]Sample{sample(t, "BenchmarkSearch/linear/first/n=100", 1)}, 1)
	report, err = Render(rows, "compact")
	if err != nil || !strings.Contains(report, "элемент в начале") {
		t.Fatalf("%v %s", err, report)
	}
}

func TestCSVLargeHeaderFailure(t *testing.T) {
	// io.Discard подтверждает, что публичный writer не требует файла.
	rows, _ := Aggregate(reportSamples(t), 1)
	if err := WriteCSV(io.Discard, rows); err != nil {
		t.Fatal(err)
	}
}

func TestTableMarkdownSeparator(t *testing.T) {
	got := table([][]string{{"k", "Время"}, {"3", "1 ns"}})
	separator := strings.Split(got, "\n")[1]
	for _, cell := range strings.Split(strings.Trim(separator, " |"), "|") {
		if len(strings.TrimSpace(cell)) < 3 {
			t.Fatalf("short Markdown separator: %q", separator)
		}
	}
}
