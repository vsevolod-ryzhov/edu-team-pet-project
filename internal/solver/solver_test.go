package solver

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"edu-team-pet-project/internal/maze"
)

func parse(t *testing.T, src string) *maze.Maze {
	t.Helper()
	m, err := maze.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("разбор лабиринта: %v", err)
	}
	return m
}

func TestDecodePathOK(t *testing.T) {
	path, err := DecodePath("OK 3\n1 1\n2 1\n3 1\n")
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	want := []maze.Point{{X: 1, Y: 1}, {X: 2, Y: 1}, {X: 3, Y: 1}}
	if len(path) != len(want) {
		t.Fatalf("длина пути = %d, ожидалось %d", len(path), len(want))
	}
	for i := range want {
		if path[i] != want[i] {
			t.Errorf("клетка %d = %v, ожидалось %v", i, path[i], want[i])
		}
	}
}

func TestDecodePathNoPath(t *testing.T) {
	if _, err := DecodePath("NO_PATH\n"); !errors.Is(err, ErrNoPath) {
		t.Fatalf("ошибка = %v, ожидалось ErrNoPath", err)
	}
}

// Это тесты на контракт: именно так Go-часть поймает решатель,
// который отвечает «почти по протоколу».
func TestDecodePathErrors(t *testing.T) {
	cases := map[string]string{
		"пустой вывод":           "",
		"мусор вместо заголовка": "готово\n1 1\n",
		"счётчик больше фактa":   "OK 3\n1 1\n",
		"счётчик меньше факта":   "OK 1\n1 1\n2 1\n",
		"координаты не числа":    "OK 1\nгде-то там\n",
		"одна координата":        "OK 1\n1\n",
		"отрицательный счётчик":  "OK -1\n",
	}
	for name, out := range cases {
		t.Run(name, func(t *testing.T) {
			if _, err := DecodePath(out); err == nil || errors.Is(err, ErrNoPath) {
				t.Fatalf("ожидалась ошибка разбора, получено %v", err)
			}
		})
	}
}

func TestReferenceFindsShortestPath(t *testing.T) {
	m := parse(t, "9 5\n#########\n#S......#\n#.......#\n#......G#\n#########\n")
	path, err := Reference{}.Solve(m)
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if err := m.ValidatePath(path); err != nil {
		t.Fatalf("путь некорректен: %v", err)
	}
	// Манхэттенское расстояние 6+2, клеток на один больше, стен на пути нет.
	if want := 9; len(path) != want {
		t.Errorf("длина пути = %d, ожидалось %d", len(path), want)
	}
}

func TestReferenceNoPath(t *testing.T) {
	m := parse(t, "5 5\n#####\n#S#G#\n#.#.#\n#.#.#\n#####\n")
	if _, err := (Reference{}).Solve(m); !errors.Is(err, ErrNoPath) {
		t.Fatalf("ошибка = %v, ожидалось ErrNoPath", err)
	}
}

// Внешний решатель подменяется shell-скриптом: так проверяется сам механизм
// запуска процесса, без зависимости от собранного C++.
func TestExternalRunsProcess(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("скрипт-заглушка рассчитан на POSIX-оболочку")
	}
	m := parse(t, "5 3\n#####\n#S.G#\n#####\n")

	fake := filepath.Join(t.TempDir(), "fake-solver")
	script := "#!/bin/sh\ncat > /dev/null\nprintf 'OK 3\\n1 1\\n2 1\\n3 1\\n'\n"
	if err := os.WriteFile(fake, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}

	path, err := External{Path: fake}.Solve(m)
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if err := m.ValidatePath(path); err != nil {
		t.Fatalf("путь некорректен: %v", err)
	}
}

func TestExternalReportsStderr(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("скрипт-заглушка рассчитан на POSIX-оболочку")
	}
	m := parse(t, "5 3\n#####\n#S.G#\n#####\n")

	fake := filepath.Join(t.TempDir(), "failing-solver")
	script := "#!/bin/sh\ncat > /dev/null\necho 'сломался на строке 2' >&2\nexit 2\n"
	if err := os.WriteFile(fake, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}

	_, err := External{Path: fake}.Solve(m)
	if err == nil {
		t.Fatal("ожидалась ошибка")
	}
	if !strings.Contains(err.Error(), "сломался на строке 2") {
		t.Errorf("сообщение из stderr потерялось: %v", err)
	}
}
