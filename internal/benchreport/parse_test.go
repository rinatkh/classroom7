package benchreport

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"
)

const header = "goos: linux\ngoarch: amd64\npkg: example/benchmarkdebug\ncpu: Test CPU\n"
const resultLine = "BenchmarkSearch/linear/last/n=100 98816059 1.136 ns/op 0 B/op 0 allocs/op"
const footer = "\nPASS\nok  example/benchmarkdebug 0.123s\n"

func TestParseName(t *testing.T) {
	for _, test := range []struct {
		raw  string
		want Name
	}{
		{"BenchmarkSearch/linear/first/n=100", Name{Root: "BenchmarkSearch", Algorithm: "linear", Scenario: "first", N: 100, Procs: 1}},
		{"BenchmarkSearch/binary/middle/n=1000-8", Name{Root: "BenchmarkSearch", Algorithm: "binary", Scenario: "middle", N: 1000, Procs: 8}},
		{"BenchmarkSearch/map/last/n=10000-1", Name{Root: "BenchmarkSearch", Algorithm: "map", Scenario: "last", N: 10000, Procs: 1}},
		{"BenchmarkSearch/slices-binary/duplicates/n=100", Name{Root: "BenchmarkSearch", Algorithm: "slices-binary", Scenario: "duplicates", N: 100, Procs: 1}},
		{"BenchmarkSearchBuild/clone-sort/n=100", Name{Root: "BenchmarkSearchBuild", Algorithm: "clone-sort", N: 100, Procs: 1}},
		{"BenchmarkSearchBuild/map-build/n=100", Name{Root: "BenchmarkSearchBuild", Algorithm: "map-build", N: 100, Procs: 1}},
		{"BenchmarkSearchQueries/sort-binary/n=100/q=1000-4", Name{Root: "BenchmarkSearchQueries", Algorithm: "sort-binary", N: 100, Q: 1000, Procs: 4}},
		{"BenchmarkSort/copy-sort/merge/random/n=512", Name{Root: "BenchmarkSort", Mode: "copy-sort", Algorithm: "merge", Scenario: "random", N: 512, Procs: 1}},
		{"BenchmarkSort/sort-only/quick/equal/n=100", Name{Root: "BenchmarkSort", Mode: "sort-only", Algorithm: "quick", Scenario: "equal", N: 100, Procs: 1}},
		{"BenchmarkMemory/map-index-build/n=100", Name{Root: "BenchmarkMemory", Algorithm: "map-index-build", N: 100, Procs: 1}},
		{"BenchmarkMemory/slice-build/n=100", Name{Root: "BenchmarkMemory", Algorithm: "slice-build", N: 100, Procs: 1}},
		{"BenchmarkStructures/value-lookup/slice-linear/n=100", Name{Root: "BenchmarkStructures", Mode: "value-lookup", Algorithm: "slice-linear", N: 100, Procs: 1}},
		{"BenchmarkStructures/value-lookup/map-key/n=100", Name{Root: "BenchmarkStructures", Mode: "value-lookup", Algorithm: "map-key", N: 100, Procs: 1}},
		{"BenchmarkStructures/known-index/slice-read/n=100", Name{Root: "BenchmarkStructures", Mode: "known-index", Algorithm: "slice-read", N: 100, Procs: 1}},
	} {
		t.Run(test.raw, func(t *testing.T) {
			got, err := ParseName(test.raw)
			if err != nil || got != test.want {
				t.Fatalf("ParseName() = %+v, %v; want %+v", got, err, test.want)
			}
		})
	}
	for operation, mode := range map[string]string{
		"slice-insert-ordered": "ordered-insert", "slice-delete-ordered": "ordered-delete",
		"slice-delete-unordered": "unordered-delete", "slice-update-known-index": "known-index-update",
		"slice-update-by-key": "update-by-key", "map-update-key": "update-by-key",
		"map-insert-key": "key-insert", "map-delete-key": "key-delete",
	} {
		t.Run(operation, func(t *testing.T) {
			got, err := ParseName("BenchmarkStructures/" + operation + "/n=1000")
			if err != nil || got.Algorithm != operation || got.Mode != mode || got.N != 1000 {
				t.Fatalf("%+v, %v", got, err)
			}
		})
	}
}

func TestParseNameRejectsMalformed(t *testing.T) {
	for _, raw := range []string{
		"", "BenchmarkSearch", "BenchmarkSearch/linear", "BenchmarkOther/a/n=1",
		"BenchmarkSearch/unknown/last/n=100", "BenchmarkSearch/linear/average/n=100",
		"BenchmarkSearch/linear/last/size=100", "BenchmarkSearch/linear/last/n=0",
		"BenchmarkSearch/linear/last/n=-1", "BenchmarkSearch/linear/last/n=x",
		"BenchmarkSearch/linear/last/n=100-0", "BenchmarkSearch/linear/last/n=100-99999999999999999999999",
		"BenchmarkSearch/linear/last/n=100/extra", "BenchmarkSearchQueries/map/n=100/count=1",
		"BenchmarkSearchQueries/map/n=100/q=0", "BenchmarkSearchQueries/map/n=100/q=x",
		"BenchmarkSearchQueries/no/n=100/q=1", "BenchmarkSearchBuild/no/n=100",
		"BenchmarkSort/unknown/merge/random/n=100", "BenchmarkSort/sort-only/no/random/n=100",
		"BenchmarkSort/sort-only/merge/no/n=100", "BenchmarkMemory/no/n=100",
		"BenchmarkStructures/no/n=100", "BenchmarkStructures/value-lookup/slice-read/n=100",
		"BenchmarkStructures/known-index/map-key/n=100",
	} {
		t.Run(raw, func(t *testing.T) {
			if _, err := ParseName(raw); err == nil {
				t.Fatalf("accepted %q", raw)
			}
		})
	}
}

func TestParseMetricsAndCompletion(t *testing.T) {
	text := header + resultLine + "\n" + strings.ReplaceAll(resultLine, "1.136 ns/op 0 B/op 0 allocs/op", "2 allocs/op 2.5e3 ns/op 1024 B/op") + footer
	got, err := Parse(strings.NewReader(text))
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0].Iterations != 98816059 || got[0].NS != 1.136 || got[1].NS != 2500 || got[1].Bytes != 1024 || got[1].Allocs != 2 {
		t.Fatalf("incorrect metrics: %+v", got)
	}
	if !reflect.DeepEqual(got[0].Environment, Environment{"linux", "amd64", "example/benchmarkdebug", "Test CPU"}) {
		t.Fatal(got[0].Environment)
	}
	// Заголовок CPU необязателен; CRLF и пустые строки допустимы.
	_, err = Parse(strings.NewReader(strings.ReplaceAll(strings.ReplaceAll(text, "cpu: Test CPU\n", ""), "\n", "\r\n\r\n")))
	if err != nil {
		t.Fatal(err)
	}
}

func TestParseRejectsBadLogs(t *testing.T) {
	bad := map[string]string{
		"пустой": "", "только заголовок": header, "нет результатов": header + "PASS\nok example/benchmarkdebug 1s\n",
		"нет PASS": header + resultLine + "\nok example/benchmarkdebug 1s\n",
		"нет ok":   header + resultLine + "\nPASS\n", "обрыв": header + resultLine,
		"нет окружения": resultLine + footer, "нет goarch": strings.ReplaceAll(header, "goarch: amd64\n", "") + resultLine + footer,
		"FAIL":                   header + resultLine + "\nFAIL\nFAIL example/benchmarkdebug 1s\n",
		"прерывание":             header + resultLine + "\nBenchmarkSort/sort-only/merge/sorted/n=10000 ^Csignal: interrupt\nFAIL\n",
		"повтор PASS":            header + resultLine + "\nPASS\n" + footer,
		"результат после PASS":   header + resultLine + "\nPASS\n" + resultLine + footer,
		"строки после ok":        header + resultLine + footer + "FAIL\n",
		"другое окружение":       header + resultLine + "\ngoos: darwin\n" + resultLine + footer,
		"повтор заголовка":       header + "goos: linux\n" + resultLine + footer,
		"другой пакет":           header + resultLine + "\nPASS\nok other 1s\n",
		"кэшированный результат": header + resultLine + "\nPASS\nok example/benchmarkdebug (cached)\n",
		"плохое время ok":        header + resultLine + "\nPASS\nok example/benchmarkdebug NaNs\n",
		"нулевой b.N":            header + strings.ReplaceAll(resultLine, "98816059", "0") + footer,
		"отрицательный b.N":      header + strings.ReplaceAll(resultLine, "98816059", "-1") + footer,
		"дробный b.N":            header + strings.ReplaceAll(resultLine, "98816059", "1.5") + footer,
		"NaN":                    header + strings.ReplaceAll(resultLine, "1.136", "NaN") + footer,
		"Inf":                    header + strings.ReplaceAll(resultLine, "1.136", "+Inf") + footer,
		"отрицательная метрика":  header + strings.ReplaceAll(resultLine, "1.136", "-1") + footer,
		"не число":               header + strings.ReplaceAll(resultLine, "1.136", "hello") + footer,
		"неизвестная единица":    header + strings.ReplaceAll(resultLine, "ns/op", "ms/op") + footer,
		"повтор единицы":         header + strings.ReplaceAll(resultLine, "B/op", "ns/op") + footer,
		"нет benchmem":           header + strings.ReplaceAll(resultLine, " 0 B/op 0 allocs/op", "") + footer,
		"неизвестное имя":        header + strings.ReplaceAll(resultLine, "BenchmarkSearch", "BenchmarkNo") + footer,
		"посторонний текст":      header + "incorrect\n" + resultLine + footer,
		"пустое поле окружения":  "goos:\n" + resultLine + footer,
	}
	for name, text := range bad {
		t.Run(name, func(t *testing.T) {
			if rows, err := Parse(strings.NewReader(text)); err == nil || rows != nil {
				t.Fatalf("expected no partial success, got %+v, %v", rows, err)
			}
		})
	}
}

type brokenReader struct{}

func (brokenReader) Read([]byte) (int, error) { return 0, errors.New("read error") }

func TestParseReadError(t *testing.T) {
	for _, reader := range []io.Reader{brokenReader{}, strings.NewReader(strings.Repeat("x", 1024*1024+1))} {
		if _, err := Parse(reader); err == nil || !strings.Contains(err.Error(), "чтение") {
			t.Fatalf("expected read error: %v", err)
		}
	}
}

func sample(t *testing.T, name string, ns float64) Sample {
	t.Helper()
	parsed, err := ParseName(name)
	if err != nil {
		t.Fatal(err)
	}
	return Sample{Environment: Environment{"linux", "amd64", "example/benchmarkdebug", "Test CPU"}, Name: parsed, Iterations: 10, NS: ns, Bytes: 8, Allocs: 1}
}

func TestParseAllLabShapes(t *testing.T) {
	// Все шесть форм сценария сортировки и все пять целей поиска доступны
	// formatter без переименования исходных benchmark.
	for _, scenario := range []string{"first", "middle", "last", "missing", "duplicates"} {
		for _, algorithm := range []string{"linear", "binary", "map", "slices-binary"} {
			sample(t, fmt.Sprintf("BenchmarkSearch/%s/%s/n=100", algorithm, scenario), 1)
		}
	}
	for _, scenario := range []string{"random", "sorted", "reverse", "nearly", "equal", "duplicates"} {
		for _, algorithm := range []string{"bubble", "selection", "insertion", "merge", "quick", "slices"} {
			sample(t, fmt.Sprintf("BenchmarkSort/sort-only/%s/%s/n=100", algorithm, scenario), 1)
		}
	}
	for _, algorithm := range []string{"linear", "sort-binary", "map"} {
		sample(t, "BenchmarkSearchQueries/"+algorithm+"/n=100/q=100", 1)
	}
}
