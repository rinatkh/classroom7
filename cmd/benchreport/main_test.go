package main

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const validLog = "goos: linux\ngoarch: amd64\npkg: example/benchmarkdebug\ncpu: Test CPU\nBenchmarkSearch/linear/last/n=100 10 1.136 ns/op 0 B/op 0 allocs/op\nPASS\nok example/benchmarkdebug 0.01s\n"

func TestValidateFlags(t *testing.T) {
	for _, value := range []string{"1x", "100x", "20ms", "1s", "0.5s"} {
		if err := validateTime(value); err != nil {
			t.Fatal(value, err)
		}
	}
	for _, value := range []string{"", "0", "0x", "-1x", "1.2x", "-1ms", "0s", "100"} {
		if err := validateTime(value); err == nil {
			t.Fatal(value)
		}
	}
	for _, value := range []string{"1", "1,2", "8"} {
		if err := validateCPU(value); err != nil {
			t.Fatal(value, err)
		}
	}
	for _, value := range []string{"", "0", "-1", "1,1", "1,", "1.5", "a", "1, 2"} {
		if err := validateCPU(value); err == nil {
			t.Fatal(value)
		}
	}
	for _, args := range [][]string{{"extra"}, {"-view", "unknown"}, {"-kind", "bad"}, {"-count", "0"}, {"-count", "-1"}, {"-benchtime", "wrong"}, {"-cpu", "0"}, {"-n", "-1"}, {"-unknown"}} {
		var out bytes.Buffer
		if err := run(context.Background(), args, &out, &out); err == nil {
			t.Fatal(args)
		}
	}
}

func TestImport(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "input.txt")
	if err := os.WriteFile(input, []byte(validLog), 0644); err != nil {
		t.Fatal(err)
	}
	for _, raw := range []bool{false, true} {
		t.Run(map[bool]string{false: "summary", true: "raw"}[raw], func(t *testing.T) {
			output := filepath.Join(dir, t.Name())
			args := []string{"-input", input, "-count", "1", "-dir", output}
			if raw {
				args = append(args, "-raw")
			}
			var out bytes.Buffer
			if err := run(context.Background(), args, &out, &out); err != nil {
				t.Fatal(err)
			}
			if raw != strings.Contains(out.String(), "BenchmarkSearch/linear/last/n=100 10 1.136") {
				t.Fatal(out.String())
			}
			if !raw && !strings.Contains(out.String(), "1.136 ns") {
				t.Fatal(out.String())
			}
			entries, err := os.ReadDir(output)
			if err != nil || len(entries) != 1 {
				t.Fatalf("%v %v", entries, err)
			}
			for _, name := range []string{"raw.txt", "summary.csv", "report.md", "environment.json"} {
				if _, err := os.Stat(filepath.Join(output, entries[0].Name(), name)); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}

func TestNoPartialSuccess(t *testing.T) {
	for name, log := range map[string]string{"interrupted": strings.TrimSuffix(validLog, "PASS\nok example/benchmarkdebug 0.01s\n") + "FAIL\n", "missing repetitions": validLog} {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			input := filepath.Join(dir, "input.txt")
			if err := os.WriteFile(input, []byte(log), 0644); err != nil {
				t.Fatal(err)
			}
			var out bytes.Buffer
			err := run(context.Background(), []string{"-input", input, "-dir", filepath.Join(dir, "reports")}, &out, &out)
			if err == nil || strings.Contains(out.String(), "Среднее арифметическое") {
				t.Fatalf("partial success: %v %s", err, out.String())
			}
			files, _ := filepath.Glob(filepath.Join(dir, "reports", "*", "summary.csv"))
			if len(files) != 0 {
				t.Fatal(files)
			}
		})
	}
}

func TestRunGoProcess(t *testing.T) {
	for _, test := range []struct {
		name, exit string
		raw        bool
	}{{"summary", "0", false}, {"raw", "0", true}, {"exit failure", "1", false}} {
		t.Run(test.name, func(t *testing.T) {
			dir := t.TempDir()
			program := filepath.Join(dir, "fake-go")
			// Команда с корректными строками, но ненулевым exit, не имеет права
			// публиковать таблицы: одних слов PASS в stdout недостаточно.
			script := "#!/bin/sh\nif [ \"$1\" = version ]; then echo 'go version test linux/amd64'; exit 0; fi\ncat <<'DATA'\n" + validLog + "DATA\nexit " + test.exit + "\n"
			if err := os.WriteFile(program, []byte(script), 0755); err != nil {
				t.Fatal(err)
			}
			args := []string{"-go", program, "-dir", filepath.Join(dir, "reports"), "-count", "1", "-benchtime", "1x"}
			if test.raw {
				args = append(args, "-raw")
			}
			var out bytes.Buffer
			err := run(context.Background(), args, &out, &out)
			if (err == nil) != (test.exit == "0") {
				t.Fatalf("%v %s", err, out.String())
			}
			files, _ := filepath.Glob(filepath.Join(dir, "reports", "*", "summary.csv"))
			if (len(files) > 0) != (test.exit == "0") {
				t.Fatal(files)
			}
			metadataFiles, _ := filepath.Glob(filepath.Join(dir, "reports", "*", "environment.json"))
			if len(metadataFiles) != 1 {
				t.Fatal(metadataFiles)
			}
			data, err := os.ReadFile(metadataFiles[0])
			if err != nil {
				t.Fatal(err)
			}
			var metadata map[string]any
			if err := json.Unmarshal(data, &metadata); err != nil {
				t.Fatal(err)
			}
			want := "success"
			if test.exit != "0" {
				want = "failed"
			}
			if metadata["state"] != want {
				t.Fatal(metadata)
			}
			if test.exit != "0" && strings.Contains(out.String(), "Среднее арифметическое") {
				t.Fatal(out.String())
			}
		})
	}
}

func TestFileAndCommandErrors(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "file")
	if err := os.WriteFile(file, nil, 0644); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{{"-dir", file}, {"-dir", filepath.Join(file, "nested")}, {"-dir", filepath.Join(dir, "new"), "-input", filepath.Join(dir, "missing")}, {"-dir", filepath.Join(dir, "new"), "-go", filepath.Join(dir, "missing")}} {
		var out bytes.Buffer
		if err := run(context.Background(), args, &out, &out); err == nil {
			t.Fatal(args)
		}
	}
	if err := writeMetadata(filepath.Join(dir, "new-missing", "file"), map[string]any{}); err == nil {
		t.Fatal("lost write error")
	}
	if err := writeMetadata(file, map[string]any{"bad": make(chan int)}); err == nil {
		t.Fatal("lost marshal error")
	}
}
