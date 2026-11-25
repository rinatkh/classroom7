# Проверочный замер после патча — 25 ноября 2026

Это реальные короткие измерения, а не перенос чисел из презентации. Среднее — **арифметическое по трём прогонам**; исходные прогоны остаются в отдельных raw.

Окружение: go version go1.22.12 linux/amd64, CPU **AMD EPYC 9V74 80-Core Processor**, GOMAXPROCS=1. Общий контейнер, частота CPU и чужая нагрузка не зафиксированы. Разброс виден в CSV (min/max), статистическая значимость не оценивалась. На Mac результаты будут другими.

## Команды

```bash
make bench-search BENCHSCENARIO=last BENCHTIME=100ms BENCHCOUNT=3
make bench-sort BENCHSCENARIO=random BENCHTIME=20ms BENCHCOUNT=3
make bench-search-build BENCHTIME=20ms BENCHCOUNT=3
make bench-search-queries BENCHTIME=20ms BENCHCOUNT=3
make bench-structures BENCHTIME=5ms BENCHCOUNT=3
make bench-memory BENCHTIME=20ms BENCHCOUNT=3
```

Для серьёзного сравнения увеличьте время до 500ms–1s и число прогонов до 5–10. Особенно шумен 5ms-пример Structures: он проверяет полный путь сбора отчёта, не даёт основания ранжировать наносекундные операции.

## Готовый поиск последнего элемента

Каждая ячейка — время одного полного поиска. Вход и map-индекс подготовлены; их построение исключено. B/op и allocs/op здесь равны нулю в напечатанных Go метриках. Это не нулевая память массива или индекса.

| n | Linear | Binary | map | slices.BinarySearch |
| ---: | ---: | ---: | ---: | ---: |
| 100 | 32.16 ns | 7.141 ns | 10.06 ns | 8.298 ns |
| 1000 | 317.6 ns | 11.82 ns | 8.115 ns | 9.838 ns |
| 10000 | 3.167 µs | 18.56 ns | 8.521 ns | 13.32 ns |
| 100000 | 30.15 µs | 27.36 ns | 9.455 ns | 16.34 ns |
| 1000000 | 307.3 µs | 30.38 ns | 8.516 ns | 20.12 ns |

## Сортировки: random, n=512

| Режим | Алгоритм | Среднее / op | Средний B/op | Средний allocs/op |
| --- | --- | ---: | ---: | ---: |
| copy-sort | bubble | 92.35 µs | 0 | 0 |
| copy-sort | insertion | 30.43 µs | 0 | 0 |
| copy-sort | merge | 21.56 µs | 4096 | 1 |
| copy-sort | quick | 9.17 µs | 0 | 0 |
| copy-sort | selection | 160.8 µs | 0 | 0 |
| copy-sort | slices | 5.816 µs | 0 | 0 |
| sort-only | bubble | 105.5 µs | 0 | 0 |
| sort-only | insertion | 29.77 µs | 0 | 0 |
| sort-only | merge | 20.57 µs | 4096 | 1 |
| sort-only | quick | 9.917 µs | 0 | 0 |
| sort-only | selection | 161.9 µs | 0 | 0 |
| sort-only | slices | 6.038 µs | 0 | 0 |

`copy-sort`: copy в готовый буфер + сортировка; allocation этого рабочего буфера исключена. `sort-only`: копирование вне таймера, выделение внутреннего Merge-буфера внутри. У Quick остаётся стек; нулевой B/op его не отменяет.

## Полные данные

- [summary.csv](summary.csv) — 106 агрегированных строк; все средние, медиана/min/max времени. Разные группы не усреднялись вместе.
- [environment.json](environment.json) — шесть запусков, точные параметры и версия Go.
- [search-raw.txt](search-raw.txt) — 20 сценариев × 3 прогона.
- [sort-raw.txt](sort-raw.txt) — 28 сценариев × 3 прогона.
- [search-build-raw.txt](search-build-raw.txt) — 10 сценариев × 3 прогона.
- [search-queries-raw.txt](search-queries-raw.txt) — 18 сценариев × 3 прогона.
- [structures-raw.txt](structures-raw.txt) — 22 сценария × 3 прогона.
- [memory-raw.txt](memory-raw.txt) — 8 сценариев × 3 прогона.

Всего 318 строк реальных измерений. Остальные формы входа сохранены в лаборатории и проверены отдельно полным smoke-прогоном: 326 сценариев × 3, benchtime=1x. Smoke-значения не используются в этих сравнительных таблицах.

CSV объединяет строки шести опытов для удобства чтения, **не** смешивает прогоны разных команд в одном среднем. Чтобы повторно форматировать один raw: `go run ./cmd/benchreport -input docs/benchmarks/fix-2026-11-25/search-raw.txt -count 3 -view full`. Остальные raw импортируются каждый отдельно.

При экспорте этих `*-raw.txt` убрано только пробельное дополнение справа
в строке `cpu:`. Имена, b.N и все измеренные метрики сохранены без изменения.
Сам `cmd/benchreport` записывает исходный raw в bin без такой нормализации.
