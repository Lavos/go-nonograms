package nonograms

import (
	"log"
	"time"
	sf "bitbucket.org/krepa098/gosfml2"
)

const (
	ModeEmpty = iota
	ModeFill
	ModeCrossOut
)

func TrackTime (start time.Time, name string) {
	log.Printf("%s: %s", name, time.Since(start))
}

func GridToPixelsf (grid_index int) float32 {
	return float32(GridToPixelsi(grid_index))
}

func GridToPixelsi (grid_index int) int {
	return GridSize * grid_index
}

type Grid struct {
	PlayTray *Tray
	HintTrayTop *Tray
	HintTrayLeft *Tray

	Solved bool

	GoalMatrix Matrix
	WorkingMatrix Matrix

	Drawers []sf.Drawer
	Eventers []Eventer

	Mode byte
}

func NewGrid (tm *TextureManager) *Grid {
	playTray := NewTray(tm, sf.Vector2f{ 100, 100 }, sf.ColorRed(), TrayTypePlay)
	hintTrayTop := NewTray(tm, sf.Vector2f{ 100, 0 }, sf.ColorBlue(), TrayTypeHint)
	hintTrayLeft := NewTray(tm, sf.Vector2f{ 0, 100 }, sf.ColorBlue(), TrayTypeHint)

	return &Grid{
		PlayTray: playTray,
		HintTrayTop: hintTrayTop,
		HintTrayLeft: hintTrayLeft,
	}
}

func (g *Grid) Render(matrix Matrix) {
	g.Solved = false
	rows := len(matrix)
	columns := len(matrix[0])

	log.Printf("Rows: %d Columns: %d", rows, columns)

	g.WorkingMatrix = NewMatrix(rows, columns)
	g.GoalMatrix = matrix

	g.PlayTray.SetSize(rows, columns)
	g.HintTrayTop.SetSize(5, 5)
	g.HintTrayLeft.SetSize(5, 5)

	g.Drawers = []sf.Drawer{
		g.HintTrayTop,
		g.HintTrayLeft,
		g.PlayTray,
	}

	g.Eventers = make([]Eventer, 0)
	row_hints := make([][]int, rows)
	column_hints := make([][]int, columns)
	column_consecutive := make([]int, columns)

	// determine hints
	for y, row := range matrix {
		var row_consecutive int

		for x, b := range row {
			if b == ByteFilled {
				row_consecutive++
				column_consecutive[x]++
			} else {
				if row_consecutive != 0 {
					row_hints[y] = append(row_hints[y], row_consecutive)
				}

				if column_consecutive[x] != 0 {
					column_hints[x] = append(column_hints[x], column_consecutive[x])
				}

				row_consecutive = 0
				column_consecutive[x] = 0
			}

			if y == rows-1 {
				if column_consecutive[x] != 0 {
					column_hints[x] = append(column_hints[x], column_consecutive[x])
					continue
				}

				if len(column_hints[x]) == 0 {
					column_hints[x] = append(column_hints[x], 0)
				}
			}
		}

		if len(row_hints[y]) == 0 || row_consecutive != 0 {
			row_hints[y] = append(row_hints[y], row_consecutive)
		}
	}
}

func (g *Grid) Logic() {
	g.Solved = g.CheckIfSolved()
}

func (g *Grid) CheckIfSolved() bool {
	for y, row := range g.GoalMatrix {
		for x, b := range row {
			if b == ByteFilled && b != g.WorkingMatrix[y][x] {
				return false
			}

			if g.WorkingMatrix[y][x] == ByteFilled && b != ByteFilled {
				return false
			}
		}
	}

	return true
}

func (g *Grid) Draw(target sf.RenderTarget, renderStates sf.RenderStates) {
	for _, h := range g.Drawers {
		h.Draw(target, renderStates)
	}
}

func (g *Grid) HandleEvent(event sf.Event) {
	// defer TrackTime(time.Now(), "HandleEvent")

	switch event.Type() {
	case sf.EventTypeMouseButtonPressed:
		if sf.IsMouseButtonPressed(sf.MouseLeft) {
			g.Mode = ModeFill
		}

		if sf.IsMouseButtonPressed(sf.MouseRight) {
			g.Mode = ModeCrossOut
		}

	case sf.EventTypeMouseButtonReleased:
		x, y := g.PlayTray.CoordsFromPosition(CurrentMousePosition.X, CurrentMousePosition.Y)
		quad, ok := g.PlayTray.QuadFromCoords(x, y)

		if !ok {
			break
		}

		switch g.WorkingMatrix[y][x] {
		case ByteFilled:
			g.WorkingMatrix[y][x] = ByteEmpty

		case ByteEmpty:
			g.WorkingMatrix[y][x] = g.Mode

		case ByteCrossedOut:
			g.WorkingMatrix[y][x] = ByteEmpty
		}

		g.PlayTray.SetState(quad, g.WorkingMatrix[y][x])
	}

	for _, e := range g.Eventers {
		e.HandleEvent(event)
	}
}
