package solver

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"edu-team-pet-project/internal/maze"
)

// External запускает внешний решатель (бинарь на C++) как отдельный процесс:
// лабиринт уходит ему в stdin, путь приходит из stdout.
//
// Такой способ выбран намеренно вместо cgo и сетевых вызовов: со стороны C++
// это ровно std::cin и std::cout, а граница между частями остаётся наглядной —
// её видно в диспетчере задач.
type External struct {
	Path    string        // путь до исполняемого файла
	Timeout time.Duration // 0 — ждать сколько угодно
}

func (e External) Name() string { return "external (" + e.Path + ")" }

func (e External) Solve(m *maze.Maze) ([]maze.Point, error) {
	ctx := context.Background()
	if e.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, e.Timeout)
		defer cancel()
	}

	cmd := exec.CommandContext(ctx, e.Path)
	cmd.Stdin = strings.NewReader(m.Encode())

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return nil, fmt.Errorf("решатель %s не ответил за %s", e.Path, e.Timeout)
	}
	if err != nil {
		// Непустой stderr — это сообщение решателя человеку, его и показываем:
		// по протоколу диагностика идёт туда, а не в stdout.
		if msg := strings.TrimSpace(stderr.String()); msg != "" {
			return nil, fmt.Errorf("решатель %s завершился с ошибкой (%v): %s", e.Path, err, msg)
		}
		return nil, fmt.Errorf("решатель %s завершился с ошибкой: %w", e.Path, err)
	}

	path, err := DecodePath(stdout.String())
	if err != nil {
		return nil, fmt.Errorf("не разобрать ответ решателя %s: %w", e.Path, err)
	}
	return path, nil
}
