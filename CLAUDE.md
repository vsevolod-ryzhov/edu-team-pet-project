# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Что это

Учебный проект для совместной работы ментора (Go) и стажёра (C++).
Цель — практика git, docker и командной работы, поэтому **простота важнее
функциональности**: не предлагать cgo, gRPC, protobuf, CMake и прочее,
что повышает порог входа. Взаимодействие Go и C++ — только через
stdin/stdout, формат описан в `docs/PROTOCOL.md`.

## Границы ответственности

`cmd/` и `internal/` — зона ментора (Go). `solver/` — зона стажёра (C++).
`docs/PROTOCOL.md` и `testdata/` общие. Не менять `solver/` и протокол
в одностороннем порядке: контракт между частями — предмет договорённости,
а не деталь реализации.

## Команды

```
go run ./cmd/maze -input testdata/corner.txt
go run ./cmd/maze -input testdata/corner.txt -solver ./solver/solver
go test ./...
go test ./internal/solver/ -run TestDecodePath
```

## Что держать в голове при правках

- `internal/solver.DecodePath` — Go-половина контракта. Любое изменение
  формата должно одновременно отражаться в `docs/PROTOCOL.md` и в
  `solver/README.md`, иначе стажёр будет писать не то.
- `solver.Reference` (BFS на Go) существует не как «настоящая» реализация,
  а как эталон для сверки и как способ не блокировать Go-часть, пока C++
  не готов. Его нельзя удалять как дубль.
- `maze.ValidatePath` вызывается CLI по умолчанию: смысл в том, чтобы
  ошибка в чужом решателе всплывала сразу и с внятным сообщением.
- Комментарии в коде объясняют *почему*, а не *что* — код читает новичок.
