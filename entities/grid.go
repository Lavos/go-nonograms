package entities

import (
	"log"
	"time"
	sf "bitbucket.org/krepa098/gosfml2"
)

const (
	ModeFill = iota
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
	TileMap *TileMap

	Solved bool

	GoalMatrix Matrix
	WorkingMatrix Matrix

	Hints []sf.Drawer

	SuccessMessage *sf.Text
	HintTexture *sf.Texture

	MousePosition sf.EventMouseMoved
	Mode byte
}

func NewGrid () *Grid {
	tm := NewTileMap()

	font, _ := sf.NewFontFromFile("TerminalVector.ttf")
	success_message, _ := sf.NewText(font)
	success_message.SetCharacterSize(12)
	success_message.SetString("Completed!")
	success_message.SetPosition(sf.Vector2f{10, 120})

	hint_texture, _ := sf.NewTextureFromFile("../assets/hints.png", nil)

	return &Grid{
		TileMap: tm,

		SuccessMessage: success_message,
		HintTexture: hint_texture,
	}
}

func (g *Grid) Render(matrix Matrix) {
	rows := len(matrix)
	columns := len(matrix[0])

	g.WorkingMatrix = NewMatrix(rows, columns)
	g.GoalMatrix = matrix

	g.TileMap.SetSize(columns, rows)

	g.Hints = make([]sf.Drawer, 0)
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

	for y, row := range row_hints {
		for x, n := range row {
			h := NewHint(g.HintTexture, n)
			h.Sprite.SetPosition(sf.Vector2f{ GridToPixelsf(11 + x - len(row)), GridToPixelsf(y + 11) })
			g.Hints = append(g.Hints, h)
		}
	}

	for x, column := range column_hints {
		for y, n := range column {
			h := NewHint(g.HintTexture, n)
			h.Sprite.SetPosition(sf.Vector2f{ GridToPixelsf(x + 11), GridToPixelsf(11 + y - len(column)) })
			g.Hints = append(g.Hints, h)
		}
	}
}

func (g *Grid) Logic() {
	g.Solved = g.CheckIfSolved()
}

func (g *Grid) CheckIfSolved() bool {

	for y, row := range g.GoalMatrix {
		for x, b := range row {
			if b != g.WorkingMatrix[y][x] {
				return false
			}
		}
	}

	return true
}

func (g *Grid) Draw(target sf.RenderTarget, renderStates sf.RenderStates) {
	if g.Solved {
		g.SuccessMessage.Draw(target, renderStates)
	}

	for _, h := range g.Hints {
		h.Draw(target, renderStates)
	}

	g.TileMap.Draw(target, renderStates)
}

func (g *Grid) HandleEvent(event sf.Event) {
	defer TrackTime(time.Now(), "HandleEvent")

	switch event.Type() {
	case sf.EventTypeMouseButtonPressed:
		if sf.IsMouseButtonPressed(sf.MouseLeft) {
			g.Mode = ModeFill
		}

		if sf.IsMouseButtonPressed(sf.MouseRight) {
			g.Mode = ModeCrossOut
		}

	case sf.EventTypeMouseMoved:
		g.MousePosition = event.(sf.EventMouseMoved)

	case sf.EventTypeMouseButtonReleased:
		x, y := g.TileMap.CoordsFromPosition(g.MousePosition.X, g.MousePosition.Y)
		quad, ok := g.TileMap.QuadFromCoords(x, y)

		log.Printf("x: %d, y: %d", x, y)
		log.Printf("quad: %#v, ok: %t", quad, ok)

		if !ok {
			return
		}

		log.Printf("Working: %d", g.WorkingMatrix[y][x])

		switch g.WorkingMatrix[y][x] {
		case ByteFilled:
			g.WorkingMatrix[y][x] = ByteEmpty

		case ByteEmpty:
			g.WorkingMatrix[y][x] = ByteFilled
		}

		g.TileMap.SetState(quad, g.WorkingMatrix[y][x])
	}
}
