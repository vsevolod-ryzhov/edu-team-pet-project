package solver

import "edu-team-pet-project/internal/maze"

// Reference — эталонный решатель на Go: поиск в ширину (BFS).
//
// Он нужен не как «настоящая» реализация, а как опора:
//   - Go-часть можно разрабатывать и тестировать до того, как готов C++;
//   - с его ответом можно сверить ответ C++ на тех же лабиринтах.
//
// BFS подходит, потому что все переходы «стоят» одинаково: первая волна,
// дошедшая до цели, и есть кратчайший путь.
type Reference struct{}

func (Reference) Name() string { return "reference (Go, BFS)" }

func (Reference) Solve(m *maze.Maze) ([]maze.Point, error) {
	// cameFrom служит сразу двум целям: отмечает посещённые клетки
	// и помнит, откуда мы в них пришли, чтобы развернуть путь в конце.
	cameFrom := make(map[maze.Point]maze.Point, m.Width*m.Height)
	cameFrom[m.Start] = m.Start

	queue := []maze.Point{m.Start}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		if cur == m.Goal {
			return restore(cameFrom, m.Start, m.Goal), nil
		}
		for _, next := range neighbours(cur) {
			if !m.Passable(next) {
				continue
			}
			if _, seen := cameFrom[next]; seen {
				continue
			}
			cameFrom[next] = cur
			queue = append(queue, next)
		}
	}
	return nil, ErrNoPath
}

func neighbours(p maze.Point) [4]maze.Point {
	return [4]maze.Point{
		{X: p.X, Y: p.Y - 1},
		{X: p.X, Y: p.Y + 1},
		{X: p.X - 1, Y: p.Y},
		{X: p.X + 1, Y: p.Y},
	}
}

func restore(cameFrom map[maze.Point]maze.Point, start, goal maze.Point) []maze.Point {
	var reversed []maze.Point
	for cur := goal; ; cur = cameFrom[cur] {
		reversed = append(reversed, cur)
		if cur == start {
			break
		}
	}
	path := make([]maze.Point, len(reversed))
	for i, p := range reversed {
		path[len(reversed)-1-i] = p
	}
	return path
}
