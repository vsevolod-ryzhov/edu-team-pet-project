// Команда maze читает лабиринт, просит решатель найти путь и рисует результат.
//
//	go run ./cmd/maze -input testdata/simple.txt
//	go run ./cmd/maze -input testdata/simple.txt -solver ./solver/solver
package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"flag"

	"edu-team-pet-project/internal/maze"
	"edu-team-pet-project/internal/solver"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "ошибка:", err)
		os.Exit(1)
	}
}

func run() error {
	input := flag.String("input", "-", "файл с лабиринтом; \"-\" — читать stdin")
	solverPath := flag.String("solver", envOr("MAZE_SOLVER", "reference"),
		"путь до решателя; \"reference\" — встроенный эталон на Go")
	timeout := flag.Duration("timeout", 5*time.Second, "сколько ждать внешний решатель")
	check := flag.Bool("check", true, "проверять корректность полученного пути")
	printRaw := flag.Bool("print-raw", false,
		"печатать ответ в формате протокола вместо картинки (для сверки через diff)")
	flag.Parse()

	src := io.Reader(os.Stdin)
	if *input != "-" {
		f, err := os.Open(*input)
		if err != nil {
			return err
		}
		defer f.Close()
		src = f
	}

	m, err := maze.Parse(src)
	if err != nil {
		return fmt.Errorf("разбор лабиринта: %w", err)
	}

	var s solver.Solver = solver.Reference{}
	if *solverPath != "reference" {
		s = solver.External{Path: *solverPath, Timeout: *timeout}
	}

	path, err := s.Solve(m)
	switch {
	case errors.Is(err, solver.ErrNoPath) && *printRaw:
		fmt.Print(solver.EncodePath(nil))
		return nil
	case errors.Is(err, solver.ErrNoPath):
		fmt.Printf("решатель: %s\n", s.Name())
		fmt.Println("пути нет")
		fmt.Print(m.Render(nil))
		return nil
	case err != nil:
		return err
	}

	if *check {
		if err := m.ValidatePath(path); err != nil {
			return fmt.Errorf("решатель %s вернул некорректный путь: %w", s.Name(), err)
		}
	}

	// В сыром режиме — только stdout в формате протокола, без заголовков:
	// иначе diff с выводом C++-решателя покажет лишние строки.
	if *printRaw {
		fmt.Print(solver.EncodePath(path))
		return nil
	}

	fmt.Printf("решатель: %s\n", s.Name())
	fmt.Printf("длина пути: %d клеток\n", len(path))
	fmt.Print(m.Render(path))
	return nil
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
