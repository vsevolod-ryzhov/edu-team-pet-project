# Многостадийная сборка: собираем в одном образе, запускаем в другом.
# Так в финальный образ не попадают ни компилятор, ни исходники —
# сравните размер образа с размером golang:1.27 и всё станет понятно.

# --- стадия 1: сборка Go-части -----------------------------------------------
FROM golang:1.27-bookworm AS go-build
WORKDIR /src
COPY go.mod ./
COPY cmd/ ./cmd/
COPY internal/ ./internal/
# CGO_ENABLED=0 даёт статический бинарь: его можно положить в образ,
# в котором нет вообще никаких библиотек.
RUN CGO_ENABLED=0 go build -trimpath -o /out/maze ./cmd/maze

# --- стадия 2: сборка решателя на C++ ----------------------------------------
# TODO(вместе, PR №5): добавить стадию
#   FROM gcc:14 AS cpp-build
# которая собирает solver/ и отдаёт готовый бинарь в финальный образ.
# Пока решатель на C++ не готов, образ работает с эталонным решателем на Go.

# --- стадия 3: то, что запускается -------------------------------------------
FROM debian:bookworm-slim
COPY --from=go-build /out/maze /usr/local/bin/maze
COPY testdata/ /testdata/
ENTRYPOINT ["maze"]
CMD ["-input", "/testdata/corner.txt"]
