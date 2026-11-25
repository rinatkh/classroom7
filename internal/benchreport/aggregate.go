package benchreport

import (
	"encoding/csv"
	"fmt"
	"io"
	"sort"
	"strconv"
)

// Row содержит арифметические средние по прогонам с одним ключом.
// MedianNS, MinNS и MaxNS дополняют среднее, но не заменяют его.
type Row struct {
	Environment                                           Environment
	Name                                                  Name
	Repeats                                               int
	MeanNS, MeanBytes, MeanAllocs, MedianNS, MinNS, MaxNS float64
}

type key struct {
	Environment
	Name
}

// Aggregate не смешивает n, q, сценарий, режим, GOMAXPROCS и окружение.
// expected=0 допускает неизвестное число повторов при импорте старого лога.
func Aggregate(samples []Sample, expected int) ([]Row, error) {
	if len(samples) == 0 || expected < 0 {
		return nil, fmt.Errorf("нет результатов или неверное число повторов")
	}
	groups := make(map[key][]Sample)
	for _, sample := range samples {
		k := key{sample.Environment, sample.Name}
		groups[k] = append(groups[k], sample)
	}
	rows := make([]Row, 0, len(groups))
	for k, group := range groups {
		if expected > 0 && len(group) != expected {
			return nil, fmt.Errorf("%s/%s n=%d: получено %d повторов вместо %d", k.Root, k.Algorithm, k.N, len(group), expected)
		}
		row := Row{Environment: k.Environment, Name: k.Name, Repeats: len(group)}
		times := make([]float64, len(group))
		for i, sample := range group {
			// Онлайн-среднее не суммирует b.N и не взвешивает им прогоны.
			divisor := float64(i + 1)
			row.MeanNS += (sample.NS - row.MeanNS) / divisor
			row.MeanBytes += (sample.Bytes - row.MeanBytes) / divisor
			row.MeanAllocs += (sample.Allocs - row.MeanAllocs) / divisor
			times[i] = sample.NS
		}
		sort.Float64s(times)
		row.MinNS, row.MaxNS = times[0], times[len(times)-1]
		row.MedianNS = times[len(times)/2]
		if len(times)%2 == 0 {
			row.MedianNS = times[len(times)/2-1]/2 + row.MedianNS/2
		}
		rows = append(rows, row)
	}
	sort.Slice(rows, func(i, j int) bool {
		a, b := rows[i], rows[j]
		// Размер сравнивается численно: 100, 512, 1000, 10000.
		if a.Environment != b.Environment {
			return environmentKey(a.Environment) < environmentKey(b.Environment)
		}
		if a.Name.Procs != b.Name.Procs {
			return a.Name.Procs < b.Name.Procs
		}
		if a.Name.Root != b.Name.Root {
			return a.Name.Root < b.Name.Root
		}
		if a.Name.N != b.Name.N {
			return a.Name.N < b.Name.N
		}
		if a.Name.Mode != b.Name.Mode {
			return a.Name.Mode < b.Name.Mode
		}
		if a.Name.Scenario != b.Name.Scenario {
			return a.Name.Scenario < b.Name.Scenario
		}
		if a.Name.Q != b.Name.Q {
			return a.Name.Q < b.Name.Q
		}
		return a.Name.Algorithm < b.Name.Algorithm
	})
	return rows, nil
}

func environmentKey(env Environment) string {
	return env.GOOS + "\x00" + env.GOARCH + "\x00" + env.Package + "\x00" + env.CPU
}

// WriteCSV сохраняет все сценарии, включая скрытые в компактном представлении.
// В raw.txt остаются исходные b.N и точность печати testing.B.
func WriteCSV(writer io.Writer, rows []Row) error {
	output := csv.NewWriter(writer)
	if err := output.Write([]string{"goos", "goarch", "pkg", "cpu", "gomaxprocs", "benchmark", "algorithm", "mode", "scenario", "n", "q", "repeats", "mean_ns_per_op", "median_ns_per_op", "min_ns_per_op", "max_ns_per_op", "mean_B_per_op", "mean_allocs_per_op"}); err != nil {
		return err
	}
	for _, row := range rows {
		n, env := row.Name, row.Environment
		record := []string{env.GOOS, env.GOARCH, env.Package, env.CPU, strconv.Itoa(n.Procs), n.Root, n.Algorithm, n.Mode, n.Scenario, strconv.Itoa(n.N), strconv.Itoa(n.Q), strconv.Itoa(row.Repeats)}
		for _, value := range []float64{row.MeanNS, row.MedianNS, row.MinNS, row.MaxNS, row.MeanBytes, row.MeanAllocs} {
			record = append(record, strconv.FormatFloat(value, 'g', -1, 64))
		}
		if err := output.Write(record); err != nil {
			return err
		}
	}
	output.Flush()
	return output.Error()
}
