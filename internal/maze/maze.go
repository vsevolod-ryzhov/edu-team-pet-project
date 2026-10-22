// Package maze описывает лабиринт, разбор его текстового представления
// и отрисовку результата.
package maze

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Клетки лабиринта. Значения совпадают с символами протокола (docs/PROTOCOL.md).
const (
	Wall    = '#'
	Free    = '.'
	StartCh = 'S'
	GoalCh  = 'G'
)

// Point — клетка лабиринта. X — столбец, Y — строка, нумерация с нуля.
type Point struct {
	X, Y int
}

func (p Point) String() string { return fmt.Sprintf("%d %d", p.X, p.Y) }

// Maze — прямоугольный лабиринт с одним стартом и одной целью.
type Maze struct {
	Width  int
	Height int
	Rows   []string
	Start  Point
	Goal   Point
}

// Parse читает лабиринт в формате протокола и проверяет его корректность.
func Parse(r io.Reader) (*Maze, error) {
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	var w, h int
	if !sc.Scan() {
		return nil, fmt.Errorf("пустой ввод: ожидалась строка с размерами")
	}
	if _, err := fmt.Sscanf(strings.TrimSpace(sc.Text()), "%d %d", &w, &h); err != nil {
		return nil, fmt.Errorf("первая строка должна быть %q: %w", "<ширина> <высота>", err)
	}
	if w <= 0 || h <= 0 {
		return nil, fmt.Errorf("размеры должны быть положительными, получено %dx%d", w, h)
	}

	m := &Maze{Width: w, Height: h, Rows: make([]string, 0, h)}
	starts, goals := 0, 0
	for y := 0; y < h; y++ {
		if !sc.Scan() {
			return nil, fmt.Errorf("заявлено %d строк, прочитано %d", h, y)
		}
		row := strings.TrimRight(sc.Text(), "\r\n")
		if len(row) != w {
			return nil, fmt.Errorf("строка %d: длина %d, ожидалось %d", y, len(row), w)
		}
		for x := 0; x < w; x++ {
			switch row[x] {
			case Wall, Free:
			case StartCh:
				starts++
				m.Start = Point{X: x, Y: y}
			case GoalCh:
				goals++
				m.Goal = Point{X: x, Y: y}
			default:
				return nil, fmt.Errorf("строка %d, столбец %d: неизвестный символ %q", y, x, row[x])
			}
		}
		m.Rows = append(m.Rows, row)
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("чтение ввода: %w", err)
	}
	if starts != 1 {
		return nil, fmt.Errorf("стартов должно быть ровно 1, найдено %d", starts)
	}
	if goals != 1 {
		return nil, fmt.Errorf("целей должно быть ровно 1, найдено %d", goals)
	}
	return m, nil
}

// At возвращает символ клетки. Для клеток за границами — Wall,
// чтобы вызывающему коду не приходилось проверять границы отдельно.
func (m *Maze) At(p Point) byte {
	if p.X < 0 || p.Y < 0 || p.X >= m.Width || p.Y >= m.Height {
		return Wall
	}
	return m.Rows[p.Y][p.X]
}

// Passable сообщает, можно ли войти в клетку.
func (m *Maze) Passable(p Point) bool { return m.At(p) != Wall }

// Encode возвращает лабиринт в формате протокола — именно это Go-часть
// отправляет решателю в stdin.
func (m *Maze) Encode() string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d %d\n", m.Width, m.Height)
	for _, row := range m.Rows {
		b.WriteString(row)
		b.WriteByte('\n')
	}
	return b.String()
}

// Render рисует лабиринт, отмечая путь символом '*'.
// Старт и цель остаются видны: по ним проще понять, что путь на месте.
func (m *Maze) Render(path []Point) string {
	grid := make([][]byte, m.Height)
	for y, row := range m.Rows {
		grid[y] = []byte(row)
	}
	for _, p := range path {
		if p.X < 0 || p.Y < 0 || p.X >= m.Width || p.Y >= m.Height {
			continue
		}
		if c := grid[p.Y][p.X]; c == StartCh || c == GoalCh {
			continue
		}
		grid[p.Y][p.X] = '*'
	}
	var b strings.Builder
	for _, row := range grid {
		b.Write(row)
		b.WriteByte('\n')
	}
	return b.String()
}

// ValidatePath проверяет, что путь действительно ведёт от старта к цели
// по свободным клеткам и без прыжков. Нужна, чтобы поймать ошибку в чужом
// решателе сразу, а не разбираться потом с кривой картинкой.
func (m *Maze) ValidatePath(path []Point) error {
	if len(path) == 0 {
		return fmt.Errorf("путь пуст")
	}
	if path[0] != m.Start {
		return fmt.Errorf("путь начинается в %v, а старт в %v", path[0], m.Start)
	}
	if path[len(path)-1] != m.Goal {
		return fmt.Errorf("путь заканчивается в %v, а цель в %v", path[len(path)-1], m.Goal)
	}
	for i, p := range path {
		if !m.Passable(p) {
			return fmt.Errorf("клетка %d (%v) — стена или вне лабиринта", i, p)
		}
		if i == 0 {
			continue
		}
		prev := path[i-1]
		dx, dy := abs(p.X-prev.X), abs(p.Y-prev.Y)
		if dx+dy != 1 {
			return fmt.Errorf("клетки %d (%v) и %d (%v) не соседние", i-1, prev, i, p)
		}
	}
	return nil
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
