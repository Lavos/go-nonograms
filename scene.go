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

type Scene struct {
	Origin sf.Vector2f

	Rows, Columns int
	HintTopRows, HintLeftColumns int


	PlayTray *Tray
	HintTrayTop *Tray
	HintTrayLeft *Tray

	Solved bool

	GoalMatrix Matrix
	WorkingMatrix Matrix

	Drawers []sf.Drawer

	Mode byte
}

func NewScene (tm *TextureManager, origin sf.Vector2f, rows, columns int) *Scene {
	// trays
	playTray := NewTray(tm, sf.ColorWhite(), TrayTypePlay)
	hintTrayTop := NewTray(tm, sf.Color{251, 233, 194, 255}, TrayTypeHint)
	hintTrayLeft := NewTray(tm, sf.Color{251, 233, 194, 255}, TrayTypeHint)

	playTray.SetPosition(origin.Plus(sf.Vector2f{ 343, 183 }))
	playTray.SetSize(rows, columns)

	var hints_top_rows,
	hints_left_columns int

	hints_top_origin := origin.Plus(sf.Vector2f{ 343, 0 })
	hints_left_origin := origin.Plus(sf.Vector2f{ 0, 183 })

	switch rows {
	case 5, 10:
		hints_top_rows = 5
		hints_top_origin.Y = origin.Y + 97

	case 15, 20:
		hints_top_rows = 10
		hints_top_origin.Y = origin.Y + 15
	}

	switch columns {
	case 5, 10:
		hints_left_columns = 5
		hints_left_origin.X = origin.X + 257

	case 15, 20:
		hints_left_columns = 10
		hints_left_origin.X = origin.X + 175

	case 30:
		hints_left_columns = 20
		hints_left_origin.X = origin.X + 11
	}

	hintTrayTop.SetSize(hints_top_rows, columns)
	hintTrayLeft.SetSize(rows, hints_left_columns)
	hintTrayTop.SetPosition(hints_top_origin)
	hintTrayLeft.SetPosition(hints_left_origin)

	drawers := []sf.Drawer{
		hintTrayTop,
		hintTrayLeft,
		playTray,
	}

	s := &Scene{
		Origin: origin,

		Rows: rows,
		Columns: columns,

		HintTopRows: hints_top_rows,
		HintLeftColumns: hints_left_columns,

		PlayTray: playTray,
		HintTrayTop: hintTrayTop,
		HintTrayLeft: hintTrayLeft,

		Drawers: drawers,
	}

	log.Printf("%#v", s)
	return s
}

func (s *Scene) SetWorkingMatrix(matrix Matrix) {
	s.WorkingMatrix = matrix
}

func (s *Scene) SetGoalMatrix(matrix Matrix) {
	rows := len(matrix)
	columns := len(matrix[0])

	if rows != s.Rows || columns != s.Columns {
		log.Printf("matrix size does not match.")
		return
	}

	s.SetWorkingMatrix(NewMatrix(rows, columns))
	s.GoalMatrix = matrix
	s.Solved = false

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
			quad, ok := s.HintTrayTop.QuadFromIndexes(x, (s.HintTopRows) - len(column) + i)

			if !ok {
				continue
			}

			s.HintTrayTop.SetState(quad, h + 1)
		}
	}

	for y, row := range row_hints {
		for i, h := range row {
			quad, ok := s.HintTrayLeft.QuadFromIndexes((s.HintLeftColumns) - len(row) + i, y)

			if !ok {
				continue
			}

			s.HintTrayLeft.SetState(quad, h + 1)
		}
	}
}

func (s *Scene) Logic() {
	s.Solved = s.CheckIfSolved()
}

func (s *Scene) CheckIfSolved() bool {
	for y, row := range s.GoalMatrix {
		for x, b := range row {
			if b == ByteFilled && b != s.WorkingMatrix[y][x] {
				return false
			}

			if s.WorkingMatrix[y][x] == ByteFilled && b != ByteFilled {
				return false
			}
		}
	}

	return true
}

func (s *Scene) Draw(target sf.RenderTarget, renderStates sf.RenderStates) {
	for _, h := range s.Drawers {
		h.Draw(target, renderStates)
	}
}

func (s *Scene) HandleEvent(event sf.Event) {
	// defer TrackTime(time.Now(), "HandleEvent")

	if s.WorkingMatrix == nil {
		return
	}

	switch event.Type() {
	case sf.EventTypeMouseButtonPressed:
		if sf.IsMouseButtonPressed(sf.MouseLeft) {
			s.Mode = ModeFill
		}

		if sf.IsMouseButtonPressed(sf.MouseRight) {
			s.Mode = ModeCrossOut
		}
	}
}

func (s *Scene) HandleViewEvent(event sf.Event, coords sf.Vector2f) {
	if s.WorkingMatrix == nil {
		return
	}

	switch event.Type() {
	case sf.EventTypeMouseButtonReleased:
		var quad []sf.Vertex
		var ok bool
		var x, y int

		x, y = s.PlayTray.IndexesFromCoords(coords)
		quad, ok = s.PlayTray.QuadFromIndexes(x, y)

		if ok {
			switch s.WorkingMatrix[y][x] {
			case ByteFilled:
				s.WorkingMatrix[y][x] = ByteEmpty

			case ByteEmpty:
				s.WorkingMatrix[y][x] = s.Mode

			case ByteCrossedOut:
				s.WorkingMatrix[y][x] = ByteEmpty
			}

			s.PlayTray.SetState(quad, s.WorkingMatrix[y][x])
			break
		}

		x, y = s.HintTrayTop.IndexesFromCoords(coords)
		quad, ok = s.HintTrayTop.QuadFromIndexes(x, y)

		if ok {
			s.HintTrayTop.ToggleHighlight(quad)
			break
		}

		x, y = s.HintTrayLeft.IndexesFromCoords(coords)
		quad, ok = s.HintTrayLeft.QuadFromIndexes(x, y)

		if ok {
			s.HintTrayLeft.ToggleHighlight(quad)
			break
		}
	}
}
