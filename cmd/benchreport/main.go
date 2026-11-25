// benchreport запускает testing.B и публикует сводку только после успешного
// завершения go test и проверки всех строк. Raw сохраняется и при ошибке.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/rinatkh/classroom7/internal/benchreport"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := run(ctx, os.Args[1:], os.Stdout, os.Stderr); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return
		}
		fmt.Fprintln(os.Stderr, "benchreport:", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("benchreport", flag.ContinueOnError)
	flags.SetOutput(stderr)
	kind := flags.String("kind", "search", "search, search-build, search-queries, sort, structures, memory")
	n := flags.Int("n", 0, "один размер входа; 0 — вся матрица")
	scenario := flags.String("scenario", "", "один сценарий поиска или сортировки")
	benchTime := flags.String("benchtime", "100ms", "время подзадачи или число вызовов, например 1x")
	count := flags.Int("count", 3, "число прогонов; при -input 0 отключает проверку количества")
	cpu := flags.String("cpu", "1", "список GOMAXPROCS через запятую, например 1,2")
	view := flags.String("view", "compact", "compact или full")
	goCommand := flags.String("go", "go", "путь к Go-компилятору")
	dir := flags.String("dir", "bin/benchmarks", "каталог результатов; внутри создаётся новый подкаталог")
	input := flags.String("input", "", "готовый чистый raw.txt вместо нового измерения")
	raw := flags.Bool("raw", false, "показывать исходный вывод go test вместо таблиц")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("лишние аргументы: %v", flags.Args())
	}
	if *view != "compact" && *view != "full" {
		return fmt.Errorf("view: нужен compact или full")
	}
	pattern, err := benchreport.Pattern(*kind, *n, *scenario)
	if err != nil {
		return err
	}
	if *count < 0 || *count == 0 && *input == "" {
		return fmt.Errorf("count должен быть > 0 при измерении")
	}
	if err := validateTime(*benchTime); err != nil {
		return err
	}
	if err := validateCPU(*cpu); err != nil {
		return err
	}
	if err := os.MkdirAll(*dir, 0755); err != nil {
		return err
	}
	outputDir, err := os.MkdirTemp(*dir, *kind+"-")
	if err != nil {
		return err
	}
	metadata := map[string]any{"started_at": time.Now().UTC().Format(time.RFC3339), "state": "running", "kind": *kind, "n": *n, "scenario": *scenario, "view": *view, "count": *count, "cpu": *cpu, "benchtime": *benchTime}
	// Записываем состояние и при сбое. Ошибка сохранения не скрывает причину
	// неуспешного измерения, но и не позволяет объявить результат успешным.
	metadataPath := filepath.Join(outputDir, "environment.json")
	if err := writeMetadata(metadataPath, metadata); err != nil {
		return err
	}
	rawPath := filepath.Join(outputDir, "raw.txt")
	var buffer bytes.Buffer
	if *input != "" {
		data, err := os.ReadFile(*input)
		if err != nil {
			return err
		}
		buffer.Write(data)
		metadata["source"] = *input
		metadata["go_version"] = "не указана в raw; импорт"
		if err := os.WriteFile(rawPath, data, 0644); err != nil {
			return err
		}
		if *raw {
			if _, err := stdout.Write(data); err != nil {
				return err
			}
		}
	} else {
		version, err := exec.CommandContext(ctx, *goCommand, "version").CombinedOutput()
		if err != nil {
			return fmt.Errorf("go version: %w: %s", err, version)
		}
		metadata["go_version"] = strings.TrimSpace(string(version))
		commandArgs := []string{"test", "./internal/benchmarkdebug", "-run", "^$", "-bench", pattern, "-benchmem", "-benchtime=" + *benchTime, "-count=" + strconv.Itoa(*count), "-cpu=" + *cpu}
		metadata["command"] = append([]string{*goCommand}, commandArgs...)
		file, err := os.Create(rawPath)
		if err != nil {
			return err
		}
		var destination io.Writer = io.MultiWriter(file, &buffer)
		if *raw {
			destination = io.MultiWriter(file, &buffer, stdout)
		} else {
			fmt.Fprintf(stdout, "Измерение %s: benchtime=%s, count=%d, cpu=%s. Результаты: %s\n", *kind, *benchTime, *count, *cpu, outputDir)
		}
		command := exec.CommandContext(ctx, *goCommand, commandArgs...)
		command.Stdout, command.Stderr = destination, destination
		command.WaitDelay = 5 * time.Second
		// Не наследуем случайный фильтр от прошлого запуска.
		command.Env = append(os.Environ(), "LESSON7_BENCH_N="+benchreport.SizeFilter(*n))
		runErr := command.Run()
		closeErr := file.Close()
		if runErr != nil || closeErr != nil {
			metadata["state"] = "failed"
			metadata["error"] = fmt.Sprint(runErr, closeErr)
			if err := writeMetadata(metadataPath, metadata); err != nil {
				return err
			}
			return fmt.Errorf("benchmark не завершился успешно (%v, close=%v); raw: %s", runErr, closeErr, rawPath)
		}
	}
	samples, parseErr := benchreport.Parse(bytes.NewReader(buffer.Bytes()))
	var rows []benchreport.Row
	if parseErr == nil {
		rows, parseErr = benchreport.Aggregate(samples, *count)
	}
	if parseErr != nil {
		metadata["state"], metadata["error"] = "failed", parseErr.Error()
		if err := writeMetadata(metadataPath, metadata); err != nil {
			return err
		}
		return fmt.Errorf("%w; raw: %s", parseErr, rawPath)
	}
	report, err := benchreport.Render(rows, *view)
	if err != nil {
		return err
	}
	settings := fmt.Sprintf("Go: %v; benchtime=%s; -count=%d; вид=%s.\n", metadata["go_version"], *benchTime, *count, *view)
	if *input != "" {
		settings = "Импорт raw: параметры времени и версия Go неизвестны; count проверяет число строк каждой группы.\n"
	} else if strings.HasSuffix(*benchTime, "x") {
		settings += "Фиксированное число вызовов: проверка веток, а не устойчивый замер скорости.\n"
	}
	report = settings + report
	var csv bytes.Buffer
	if err := benchreport.WriteCSV(&csv, rows); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(outputDir, "summary.csv"), csv.Bytes(), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(outputDir, "report.md"), []byte(report), 0644); err != nil {
		return err
	}
	metadata["state"] = "success"
	metadata["finished_at"] = time.Now().UTC().Format(time.RFC3339)
	metadata["environment"] = rows[0].Environment
	metadata["int_bytes"] = benchreport.IntBytes(rows[0].Environment.GOARCH)
	if err := writeMetadata(metadataPath, metadata); err != nil {
		return err
	}
	if !*raw {
		if _, err := io.WriteString(stdout, report); err != nil {
			return err
		}
	}
	fmt.Fprintf(stdout, "\nСохранено: %s/{raw.txt,summary.csv,report.md,environment.json}\n", outputDir)
	return nil
}

func validateTime(value string) error {
	if strings.HasSuffix(value, "x") {
		n, err := strconv.ParseUint(strings.TrimSuffix(value, "x"), 10, 64)
		if err == nil && n > 0 {
			return nil
		}
	} else if duration, err := time.ParseDuration(value); err == nil && duration > 0 {
		return nil
	}
	return fmt.Errorf("benchtime: нужна положительная длительность или Nx, получено %q", value)
}

func validateCPU(value string) error {
	seen := map[int]bool{}
	for _, part := range strings.Split(value, ",") {
		n, err := strconv.Atoi(part)
		if err != nil || n < 1 || seen[n] {
			return fmt.Errorf("cpu: нужны разные положительные целые, получено %q", value)
		}
		seen[n] = true
	}
	return nil
}

func writeMetadata(path string, value map[string]any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0644)
}
