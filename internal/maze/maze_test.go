package maze

import (
	"strings"
	"testing"
)

func TestParseValid(t *testing.T) {
	m, err := Parse(strings.NewReader("5 3\n#####\n#S.G#\n#####\n"))
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if m.Width != 5 || m.Height != 3 {
		t.Errorf("размеры = %dx%d, ожидалось 5x3", m.Width, m.Height)
	}
	if want := (Point{X: 1, Y: 1}); m.Start != want {
		t.Errorf("старт = %v, ожидалось %v", m.Start, want)
	}
	if want := (Point{X: 3, Y: 1}); m.Goal != want {
		t.Errorf("цель = %v, ожидалось %v", m.Goal, want)
	}
}

func TestParseErrors(t *testing.T) {
	cases := map[string]string{
		"пустой ввод":        "",
		"нет размеров":       "#####\n",
		"строка короче":      "5 2\n####\n#S.G#\n",
		"строк меньше":       "5 3\n#####\n#S.G#\n",
		"нет старта":         "5 3\n#####\n#..G#\n#####\n",
		"два старта":         "5 3\n#####\n#SSG#\n#####\n",
		"нет цели":           "5 3\n#####\n#S..#\n#####\n",
		"неизвестный символ": "5 3\n#####\n#SxG#\n#####\n",
		"нулевая ширина":     "0 1\n\n",
	}
	for name, in := range cases {
		t.Run(name, func(t *testing.T) {
			if _, err := Parse(strings.NewReader(in)); err == nil {
				t.Fatal("ожидалась ошибка, получено nil")
			}
		})
	}
}

func TestEncodeRoundTrip(t *testing.T) {
	const src = "7 5\n#######\n#S#...#\n#.#.#.#\n#...#G#\n#######\n"
	m, err := Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("разбор: %v", err)
	}
	if got := m.Encode(); got != src {
		t.Errorf("Encode() = %q, ожидалось %q", got, src)
	}
}

func TestRenderKeepsStartAndGoal(t *testing.T) {
	m, err := Parse(strings.NewReader("5 3\n#####\n#S.G#\n#####\n"))
	if err != nil {
		t.Fatal(err)
	}
	got := m.Render([]Point{{X: 1, Y: 1}, {X: 2, Y: 1}, {X: 3, Y: 1}})
	if want := "#####\n#S*G#\n#####\n"; got != want {
		t.Errorf("Render() = %q, ожидалось %q", got, want)
	}
}

func TestValidatePath(t *testing.T) {
	m, err := Parse(strings.NewReader("5 3\n#####\n#S.G#\n#####\n"))
	if err != nil {
		t.Fatal(err)
	}
	good := []Point{{X: 1, Y: 1}, {X: 2, Y: 1}, {X: 3, Y: 1}}
	if err := m.ValidatePath(good); err != nil {
		t.Errorf("корректный путь отклонён: %v", err)
	}

	bad := map[string][]Point{
		"пустой":        nil,
		"не от старта":  {{X: 2, Y: 1}, {X: 3, Y: 1}},
		"не до цели":    {{X: 1, Y: 1}, {X: 2, Y: 1}},
		"прыжок":        {{X: 1, Y: 1}, {X: 3, Y: 1}},
		"через стену":   {{X: 1, Y: 1}, {X: 1, Y: 0}, {X: 2, Y: 0}, {X: 3, Y: 0}, {X: 3, Y: 1}},
		"вне лабиринта": {{X: 1, Y: 1}, {X: 1, Y: -1}, {X: 3, Y: 1}},
	}
	for name, path := range bad {
		t.Run(name, func(t *testing.T) {
			if err := m.ValidatePath(path); err == nil {
				t.Fatal("ожидалась ошибка, получено nil")
			}
		})
	}
}
