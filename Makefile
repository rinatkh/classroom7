SHELL := /bin/bash
GO ?= go
GO_PACKAGES := ./...
UNIT_PACKAGES := ./internal/...
INTEGRATION_PACKAGES := ./test/integration/...
BIN_DIR := bin
COVERAGE_FILE ?= coverage.out
COVERAGE_THRESHOLD ?= 90.0
PACKAGE_FILE ?= $(BIN_DIR)/classroom7-linux-amd64.tar.gz
CMDS := 01_complexity 02_search 03_maps 04_insert_delete 05_simple_sorts 06_quicksort 07_benchmark_debug 08_leetcode

.PHONY: help deps-check mod-check fmt fmt-check vet test test-unit test-integration test-race coverage coverage-check build package clean run-all ci bench bench-search bench-sort trace-quicksort $(addprefix run-,$(CMDS))

help:
	@echo "Доступные команды:"
	@echo "  make run-all          - запустить восемь примеров урока"
	@echo "  make trace-quicksort  - показать состояния Quick Sort по шагам"
	@echo "  make bench-search     - сравнить линейный, бинарный поиск и map"
	@echo "  make bench-sort       - сравнить учебные и библиотечную сортировки"
	@echo "  make bench            - выполнить все benchmark-тесты"
	@echo "  make test-unit        - запустить модульные тесты"
	@echo "  make test-integration - запустить интеграционные тесты команд"
	@echo "  make ci               - выполнить полную локальную проверку"

deps-check:
	$(GO) mod download
	$(GO) mod verify

mod-check:
	$(GO) mod tidy
	@if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then \
		git diff --exit-code -- go.mod; \
		if [ -f go.sum ]; then git diff --exit-code -- go.sum; fi; \
	else \
		echo "Проверка git diff пропущена: каталог не является Git-репозиторием"; \
	fi

fmt:
	gofmt -w $$(find . -name '*.go' -not -path './$(BIN_DIR)/*')

fmt-check:
	@files="$$(gofmt -l $$(find . -name '*.go' -not -path './$(BIN_DIR)/*'))"; \
	if [ -n "$$files" ]; then echo "Найдены неотформатированные Go-файлы:"; echo "$$files"; exit 1; fi

vet:
	$(GO) vet $(GO_PACKAGES)

test: test-unit test-integration

test-unit:
	$(GO) test $(UNIT_PACKAGES)

test-integration:
	$(GO) test $(INTEGRATION_PACKAGES)

test-race:
	$(GO) test -race $(UNIT_PACKAGES)

coverage:
	$(GO) test $(UNIT_PACKAGES) -covermode=atomic -coverprofile=$(COVERAGE_FILE)
	$(GO) tool cover -func=$(COVERAGE_FILE)

coverage-check: coverage
	@coverage="$$($(GO) tool cover -func=$(COVERAGE_FILE) | awk '/^total:/ {gsub("%", "", $$3); print $$3}')"; \
	awk -v coverage="$$coverage" -v threshold="$(COVERAGE_THRESHOLD)" 'BEGIN { \
		if (coverage + 0 < threshold + 0) { printf "Покрытие %.1f%% ниже порога %.1f%%\n", coverage, threshold; exit 1 } \
		printf "Покрытие %.1f%% не ниже порога %.1f%%\n", coverage, threshold; \
	}'

run-all:
	@for cmd in $(CMDS); do \
		echo "== $$cmd =="; \
		$(GO) run ./cmd/$$cmd; \
	done

run-01_complexity:
	$(GO) run ./cmd/01_complexity
run-02_search:
	$(GO) run ./cmd/02_search
run-03_maps:
	$(GO) run ./cmd/03_maps
run-04_insert_delete:
	$(GO) run ./cmd/04_insert_delete
run-05_simple_sorts:
	$(GO) run ./cmd/05_simple_sorts
run-06_quicksort:
	$(GO) run ./cmd/06_quicksort
run-07_benchmark_debug:
	$(GO) run ./cmd/07_benchmark_debug
run-08_leetcode:
	$(GO) run ./cmd/08_leetcode

trace-quicksort: run-07_benchmark_debug

bench-search:
	$(GO) test ./internal/benchmarkdebug -run '^$$' -bench '^BenchmarkSearch' -benchmem -count=3

bench-sort:
	$(GO) test ./internal/benchmarkdebug -run '^$$' -bench '^BenchmarkSort' -benchmem -count=3

bench:
	$(GO) test ./internal/benchmarkdebug -run '^$$' -bench . -benchmem -count=3

build:
	@mkdir -p $(BIN_DIR)
	@for cmd in $(CMDS); do $(GO) build -o $(BIN_DIR)/$$cmd ./cmd/$$cmd; done

package: build
	@mkdir -p $(BIN_DIR)
	tar -czf $(PACKAGE_FILE) -C $(BIN_DIR) $(CMDS)

ci: deps-check mod-check fmt-check vet test-unit test-integration test-race coverage-check build package

clean:
	rm -rf $(BIN_DIR) $(COVERAGE_FILE)
