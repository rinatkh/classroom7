# GitHub Actions — Classroom 7

Workflow `.github/workflows/ci.yml` продолжает проект после
`classroom7_fix.patch`. Все команды выполняются из корня отдельного
репозитория classroom7. Makefile и готовые алгоритмы сохранены.

## Автоматический запуск

Проверки запускаются при каждом push и при открытии, обновлении или повторном
открытии pull request. Повторный push в ту же ветку отменяет устаревший запуск
этого workflow. Имя ветки не ограничено main/master.

Используется Go из `go.mod`: для текущего `go 1.22` устанавливается patch-версия
ветки 1.22. Автоматическая смена toolchain отключена. На Linux включён CGO для
race detector. Секреты для этого workflow не требуются.

| Шаг | Команда / результат |
| --- | --- |
| Полная проверка | `make ci`: зависимости, чистота go.mod/go.sum, gofmt, vet, unit/integration tests, race, coverage ≥90%, build/package |
| Benchmark smoke | `make bench-all BENCHTIME=1x BENCHCOUNT=1 BENCHCPU=1`: все 326 подзадач, без оценки реального быстродействия |
| Выдача сборки | Архив Linux/amd64, coverage и raw/CSV/Markdown benchmark-отчёты в Artifacts текущего запуска, хранение 7 дней |

Ошибка make или теста делает CI красным. Шаг сохранения доступных диагностических
файлов выполняется также при ошибке, сохраняя статус неуспешной проверки.
Числа smoke нельзя сравнивать с полноценными локальными измерениями.

## Как запустить на GitHub

После применения дополнительного патча сохраните `.github` в Git и отправьте
коммит в свой репозиторий:

```bash
git add .github
git commit -m "Add GitHub Actions for lesson 7"
git push
```

Откройте **Actions → Classroom 7 CI → последний запуск**. Проверки появятся
также в **Checks** pull request. В конце страницы запуска находятся Artifacts.
В скачанном артефакте архив `classroom7-linux-amd64.tar.gz` содержит десять
исполняемых файлов: девять примеров и benchreport. Для запуска на Mac собирайте
проект локально: эта выдаваемая сборка предназначена для Linux/amd64.

Ручной запуск: **Actions → Classroom 7 CI → Run workflow → ветка → Run workflow**.
Кнопка появляется после попадания workflow в default branch репозитория;
доступны ветки, содержащие этот файл. Если Actions отключены в fork, сначала
включите их на вкладке Actions. Для первого PR из чужого fork GitHub может
потребовать одобрение запуска владельцем репозитория.

Локальная проверка тех же команд:

```bash
make ci
make bench-all BENCHTIME=1x BENCHCOUNT=1 BENCHCPU=1
```

Изменения публикуются дополнительным патчем; удалённый GitHub run появится
после commit/push или ручного запуска. Сборка доставляется через Artifacts;
настройка сервера для развёртывания этому алгоритмическому проекту не требуется.

Документация: [Go в GitHub Actions](https://docs.github.com/en/actions/tutorials/build-and-test-code/go),
[события запуска](https://docs.github.com/en/actions/how-tos/write-workflows/choose-when-workflows-run/trigger-a-workflow),
[setup-go](https://github.com/actions/setup-go),
[upload-artifact](https://github.com/actions/upload-artifact).
