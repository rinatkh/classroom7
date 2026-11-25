SHELL := /bin/bash
GO ?= go
GO_PACKAGES := ./...
UNIT_PACKAGES := ./internal/... ./cmd/benchreport
INTEGRATION_PACKAGES := ./test/integration/...
BIN_DIR := bin
COVERAGE_FILE ?= coverage.out
COVERAGE_THRESHOLD ?= 90.0
PACKAGE_FILE ?= $(BIN_DIR)/classroom7-linux-amd64.tar.gz
CMDS := 01_complexity 02_search 03_maps 04_insert_delete 05_simple_sorts 06_quicksort 07_benchmark_debug 08_leetcode 09_mergesort

.PHONY: help deps-check mod-check fmt fmt-check vet test test-unit test-integration test-race coverage coverage-check build package clean run-all ci bench bench-all bench-search bench-sort bench-structures bench-memory test-leetcode trace-mergesort trace-quicksort $(addprefix run-,$(CMDS))

BENCHTIME ?= 100ms
BENCHCOUNT ?= 3
BENCHCPU ?= 1
BENCHVIEW ?= compact
BENCHN ?= 0
BENCHSCENARIO ?=
BENCHDIR ?= $(BIN_DIR)/benchmarks
BENCH_GROUPS := search search-build search-queries sort structures memory
BENCH_TARGETS := $(addprefix bench-,$(BENCH_GROUPS))
BUILD_CMDS := $(CMDS) benchreport
.PHONY: $(BENCH_TARGETS) $(addsuffix -raw,$(BENCH_TARGETS)) bench-all-raw

help: ## СПРАВКА | Все публичные команды и параметры
	@awk 'BEGIN { FS=":.*## " } /^[a-zA-Z0-9_-]+:.*## / { split($$2, p, "[|]"); g=p[1]; sub(/ +$$/, "", g); sub(/^ +/, "", p[2]); if (!(g in known)) { order[++n]=g; known[g]=1 } text[g]=text[g] sprintf("  make %-25s %s\n", $$1, p[2]) } END { for (i=1;i<=n;i++) printf "\n%s\n%s", order[i], text[order[i]] }' $(MAKEFILE_LIST)
	@printf '\nПараметры: BENCHTIME=100ms BENCHCOUNT=3 BENCHCPU=1 BENCHVIEW=compact\n'
	@printf 'BENCHN=0 — все размеры; BENCHSCENARIO= — все сценарии; BENCHDIR=bin/benchmarks.\n'
	@printf 'Быстрый пример: make bench-search BENCHN=100 BENCHSCENARIO=last BENCHTIME=20ms\n'
	@printf 'Полный вид: make bench-search BENCHVIEW=full\n'
	@printf 'Тщательнее: make bench-search BENCHTIME=500ms BENCHCOUNT=5\n'
	@printf 'Raw: make bench-search-raw; все шесть групп: make bench-all (последовательно).\n'

deps-check: ## ЗАВИСИМОСТИ | Скачать и проверить зависимости
	$(GO) mod download
	$(GO) mod verify

mod-check: ## ЗАВИСИМОСТИ | go mod tidy и проверка diff go.mod/go.sum
	$(GO) mod tidy
	@if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then \
		git diff --exit-code -- go.mod go.sum; \
	else \
		echo "Проверка git diff пропущена: каталог не является Git-репозиторием"; \
	fi

fmt: ## КАЧЕСТВО | Применить gofmt ко всем Go-файлам
	gofmt -w $$(find . -name '*.go' -not -path './$(BIN_DIR)/*')

fmt-check: ## КАЧЕСТВО | Проверить gofmt без изменения файлов
	@files="$$(gofmt -l $$(find . -name '*.go' -not -path './$(BIN_DIR)/*'))"; \
	if [ -n "$$files" ]; then echo "Найдены неотформатированные Go-файлы:"; echo "$$files"; exit 1; fi

vet: ## КАЧЕСТВО | Проверить статический анализ
	$(GO) vet $(GO_PACKAGES)

test: test-unit test-integration ## ТЕСТЫ | Unit + integration

test-unit: ## ТЕСТЫ | Модульные тесты internal
	$(GO) test $(UNIT_PACKAGES)

test-integration: ## ТЕСТЫ | Запуск девяти примеров
	$(GO) test $(INTEGRATION_PACKAGES)

test-race: ## ТЕСТЫ | Все тесты с race detector
	$(GO) test -race $(GO_PACKAGES)

coverage: ## КАЧЕСТВО | Сформировать отчёт покрытия internal
	$(GO) test $(UNIT_PACKAGES) -coverpkg=./internal/... -covermode=atomic -coverprofile=$(COVERAGE_FILE)
	$(GO) tool cover -func=$(COVERAGE_FILE)

coverage-check: coverage ## КАЧЕСТВО | Проверить порог покрытия (COVERAGE_THRESHOLD=90)
	@coverage="$$($(GO) tool cover -func=$(COVERAGE_FILE) | awk '/^total:/ {gsub("%", "", $$3); print $$3}')"; \
	awk -v coverage="$$coverage" -v threshold="$(COVERAGE_THRESHOLD)" 'BEGIN { \
		if (coverage + 0 < threshold + 0) { printf "Покрытие %.1f%% ниже порога %.1f%%\n", coverage, threshold; exit 1 } \
		printf "Покрытие %.1f%% не ниже порога %.1f%%\n", coverage, threshold; \
	}'

run-all: ## ЗАПУСК ПРИМЕРОВ | Все девять примеров последовательно
	@set -e; for cmd in $(CMDS); do \
		echo "== $$cmd =="; \
		$(GO) run ./cmd/$$cmd; \
	done

run-01_complexity: ## ЗАПУСК ПРИМЕРОВ | Подсчёт операций и Big O
	$(GO) run ./cmd/01_complexity
run-02_search: ## ЗАПУСК ПРИМЕРОВ | Линейный и бинарный поиск
	$(GO) run ./cmd/02_search
run-03_maps: ## ЗАПУСК ПРИМЕРОВ | Частоты и индекс map
	$(GO) run ./cmd/03_maps
run-04_insert_delete: ## ЗАПУСК ПРИМЕРОВ | Вставка/удаление slice и map
	$(GO) run ./cmd/04_insert_delete
run-05_simple_sorts: ## ЗАПУСК ПРИМЕРОВ | Bubble, Selection, Insertion
	$(GO) run ./cmd/05_simple_sorts
run-06_quicksort: ## ЗАПУСК ПРИМЕРОВ | Quick Sort
	$(GO) run ./cmd/06_quicksort
run-07_benchmark_debug: ## ЗАПУСК ПРИМЕРОВ | Пошаговая трассировка Quick Sort
	$(GO) run ./cmd/07_benchmark_debug
run-08_leetcode: ## ЗАПУСК ПРИМЕРОВ | Binary Search 704
	$(GO) run ./cmd/08_leetcode

run-09_mergesort: ## ЗАПУСК ПРИМЕРОВ | Merge Sort и шаги слияния
	$(GO) run ./cmd/09_mergesort

trace-mergesort: run-09_mergesort ## ОТЛАДКА | Псевдоним run-09_mergesort

test-leetcode: ## ТЕСТЫ | Тесты адаптации LeetCode 704
	$(GO) test ./internal/leetcode

trace-quicksort: run-07_benchmark_debug ## ОТЛАДКА | Псевдоним run-07_benchmark_debug

bench-search: ## БЕНЧМАРКИ | Готовый поиск; compact показывает last, CSV — все сценарии
bench-search-build: ## БЕНЧМАРКИ | Отдельно построение map и новая копия + сортировка
bench-search-queries: ## БЕНЧМАРКИ | Полный сценарий: подготовка + q запросов
bench-sort: ## БЕНЧМАРКИ | Сортировки; compact показывает random, CSV — все сценарии
bench-structures: ## БЕНЧМАРКИ | Slice/map с разделением контрактов операций
bench-memory: ## БЕНЧМАРКИ | Allocation при построении структур, не retained heap
bench-search-raw: ## СЫРОЙ ВЫВОД | Исходный testing.B: готовый поиск
bench-search-build-raw: ## СЫРОЙ ВЫВОД | Исходный testing.B: подготовка
bench-search-queries-raw: ## СЫРОЙ ВЫВОД | Исходный testing.B: подготовка + q запросов
bench-sort-raw: ## СЫРОЙ ВЫВОД | Исходный testing.B: сортировки
bench-structures-raw: ## СЫРОЙ ВЫВОД | Исходный testing.B: slice/map
bench-memory-raw: ## СЫРОЙ ВЫВОД | Исходный testing.B: allocation

# Имя target определяет группу; описания публичных targets — реестр help.
$(BENCH_TARGETS) $(addsuffix -raw,$(BENCH_TARGETS)):
	@$(GO) run ./cmd/benchreport -go "$(GO)" -kind $(patsubst bench-%,%,$(patsubst %-raw,%,$@)) -benchtime "$(BENCHTIME)" -count $(BENCHCOUNT) -cpu "$(BENCHCPU)" -view "$(BENCHVIEW)" -n $(BENCHN) -scenario "$(BENCHSCENARIO)" -dir "$(BENCHDIR)" $(if $(filter %-raw,$@),-raw,)

# Подзадачи выполняются последовательно даже при make -j.
bench-all: ## БЕНЧМАРКИ | Все шесть групп, последовательно
	@set -e; for group in $(BENCH_GROUPS); do $(MAKE) --no-print-directory bench-$$group; done

bench-all-raw: ## СЫРОЙ ВЫВОД | Все шесть групп testing.B, последовательно
	@set -e; for group in $(BENCH_GROUPS); do $(MAKE) --no-print-directory bench-$$group-raw; done

bench: bench-all ## БЕНЧМАРКИ | Псевдоним bench-all

build: ## СБОРКА | Собрать девять примеров и benchreport в bin
	@mkdir -p $(BIN_DIR)
	@set -e; for cmd in $(BUILD_CMDS); do $(GO) build -buildvcs=false -o $(BIN_DIR)/$$cmd ./cmd/$$cmd; done

package: build ## СБОРКА | Упаковать собранные исполняемые файлы
	@mkdir -p $(BIN_DIR)
	tar -czf $(PACKAGE_FILE) -C $(BIN_DIR) $(BUILD_CMDS)

ci: deps-check mod-check fmt-check vet test-unit test-integration test-race coverage-check build package ## КАЧЕСТВО | Полная локальная проверка, сборка и упаковка

clean: ## СБОРКА | Удалить bin (включая benchmark-отчёты) и coverage
	rm -rf $(BIN_DIR) $(COVERAGE_FILE)
