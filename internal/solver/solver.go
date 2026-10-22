// Package solver содержит две реализации поиска пути: эталонную на Go
// и обёртку вокруг внешнего процесса (решателя на C++).
//
// Обе реализуют Solver, поэтому CLI не знает, кто именно считает путь.
package solver

import (
	"errors"
	"fmt"
	"strings"

	"edu-team-pet-project/internal/maze"
)

// ErrNoPath означает, что лабиринт корректен, но цель недостижима.
// Это не ошибка выполнения, а один из двух нормальных ответов решателя.
var ErrNoPath = errors.New("путь не найден")

// Solver ищет путь от старта к цели.
type Solver interface {
	Solve(m *maze.Maze) ([]maze.Point, error)
	Name() string
}

// DecodePath разбирает ответ решателя в формате протокола (docs/PROTOCOL.md).
// Вынесено отдельно, потому что это и есть та половина контракта, которую
// держит Go-часть: по ней удобно писать тесты, не запуская C++.
func DecodePath(out string) ([]maze.Point, error) {
	lines := strings.Split(strings.ReplaceAll(strings.TrimSpace(out), "\r\n", "\n"), "\n")
	if len(lines) == 0 || lines[0] == "" {
		return nil, fmt.Errorf("решатель ничего не вывел")
	}

	head := strings.TrimSpace(lines[0])
	if head == "NO_PATH" {
		return nil, ErrNoPath
	}

	var n int
	if _, err := fmt.Sscanf(head, "OK %d", &n); err != nil {
		return nil, fmt.Errorf("первая строка должна быть %q или %q, получено %q", "OK <N>", "NO_PATH", head)
	}
	if n < 0 {
		return nil, fmt.Errorf("длина пути отрицательна: %d", n)
	}
	if got := len(lines) - 1; got != n {
		return nil, fmt.Errorf("заявлено %d клеток, в выводе %d", n, got)
	}

	path := make([]maze.Point, 0, n)
	for i, line := range lines[1:] {
		var p maze.Point
		if _, err := fmt.Sscanf(strings.TrimSpace(line), "%d %d", &p.X, &p.Y); err != nil {
			return nil, fmt.Errorf("клетка %d: ожидалось %q, получено %q", i, "<x> <y>", line)
		}
		path = append(path, p)
	}
	return path, nil
}
