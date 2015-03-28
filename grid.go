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
	playTray := NewTray(tm, sf.ColorWhite(), TrayTypePlay)
	hintTrayTop := NewTray(tm, sf.Color{ 251, 233, 194, 255}, TrayTypeHint)
	hintTrayLeft := NewTray(tm, sf.Color{ 251, 233, 194, 255}, TrayTypeHint)

	playTray.SetPosition(sf.Vector2f{ 343, 183 })

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

	var hints_top_rows,
	hints_left_columns int

	hints_top_origin := sf.Vector2f{ 343, 0 }
	hints_left_origin := sf.Vector2f{ 0, 183 }

	switch rows {
	case 5, 10:
		hints_top_rows = 5
		hints_top_origin.Y = 97

	case 15, 20:
		hints_top_rows = 10
		hints_top_origin.Y = 15
	}

	switch columns {
	case 5, 10:
		hints_left_columns = 5
		hints_left_origin.X = 257

	case 15, 20:
		hints_left_columns = 10
		hints_left_origin.X = 175

	case 30:
		hints_left_columns = 20
		hints_left_origin.X = 11
	}

	g.WorkingMatrix = NewMatrix(rows, columns)
	g.GoalMatrix = matrix

	g.PlayTray.SetSize(rows, columns)

	g.HintTrayTop.SetSize(hints_top_rows, columns)
	g.HintTrayLeft.SetSize(rows, hints_left_columns)
	g.HintTrayTop.SetPosition(hints_top_origin)
	g.HintTrayLeft.SetPosition(hints_left_origin)

	g.Drawers = []sf.Drawer{
		g.HintTrayTop,
		g.HintTrayLeft,
		g.PlayTray,
	}

	g.Eventers = make([]Eventer, 0)
	row_hints := make([][]byte, rows)
	column_hints := make([][]byte, columns)
	column_consecutive := make([]byte, columns)

	// determine hints
	for y, row := range matrix {
		var row_consecutive byte

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

	// set hints into trays
	for x, column := range column_hints {
		for i, h := range column {
			log.Printf("Column Hint %d %d: %d", x, i, h)

			quad, ok := g.HintTrayTop.QuadFromCoords(x, (hints_top_rows) - len(column) + i)

			if !ok {
				log.Printf("Not OK!")
				continue
			}

			g.HintTrayTop.SetState(quad, h + 1)
		}
	}

	for y, row := range row_hints {
		for i, h := range row {
			quad, ok := g.HintTrayLeft.QuadFromCoords((hints_left_columns) - len(row) + i, y)

			if !ok {
				log.Printf("Not OK!")
				continue
			}

			g.HintTrayLeft.SetState(quad, h + 1)
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
		var quad []sf.Vertex
		var ok bool
		var x, y int

		x, y = g.PlayTray.CoordsFromPosition(CurrentMousePosition.X, CurrentMousePosition.Y)
		quad, ok = g.PlayTray.QuadFromCoords(x, y)

		if ok {
			switch g.WorkingMatrix[y][x] {
			case ByteFilled:
				g.WorkingMatrix[y][x] = ByteEmpty

			case ByteEmpty:
				g.WorkingMatrix[y][x] = g.Mode

			case ByteCrossedOut:
				g.WorkingMatrix[y][x] = ByteEmpty
			}

			g.PlayTray.SetState(quad, g.WorkingMatrix[y][x])
			break
		}

		x, y = g.HintTrayTop.CoordsFromPosition(CurrentMousePosition.X, CurrentMousePosition.Y)
		quad, ok = g.HintTrayTop.QuadFromCoords(x, y)

		if ok {
			g.HintTrayTop.ToggleHighlight(quad)
			break
		}

		x, y = g.HintTrayLeft.CoordsFromPosition(CurrentMousePosition.X, CurrentMousePosition.Y)
		quad, ok = g.HintTrayLeft.QuadFromCoords(x, y)

		if ok {
			g.HintTrayLeft.ToggleHighlight(quad)
			break
		}
	}

	for _, e := range g.Eventers {
		e.HandleEvent(event)
	}
}
