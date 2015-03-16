package entities

import (
	"log"
	sf "bitbucket.org/krepa098/gosfml2"
)

const (
	ModeFill = iota
	ModeCrossOut
)

func GridToPixelsf (grid_index int) float32 {
	return float32(GridToPixelsi(grid_index))
}

func GridToPixelsi (grid_index int) int {
	return GridSize * grid_index
}

type Grid struct {
	TileMap *sf.VertexArray

	Solved bool
	Patterner Patterner

	Drawers []sf.Drawer

	SuccessMessage *sf.Text

	Tiles [][]*Tile

	MousePosition sf.EventMouseMoved
	Mode byte
	HighlightedTile *Tile
}

func NewGrid (p Patterner) *Grid {
	matrix := p.Matrix()
	log.Printf("Matrix: %#v", matrix)
	rows := len(matrix)
	columns := len(matrix[0])

	font, _ := sf.NewFontFromFile("TerminalVector.ttf")
	texture, err := sf.NewTextureFromFile("../assets/tile.png", nil)

	log.Printf("Texture Error: %s", err)
	hint_texture, err := sf.NewTextureFromFile("../assets/hints.png", nil)

	log.Printf("Texture Error: %s", err)

	drawers := make([]sf.Drawer, 0)
	tiles := make([][]*Tile, rows)

	tm := NewTileMap()
	tm.SetSize(rows, columns)

	drawers = append(drawers, tm)

	success_message, _ := sf.NewText(font)
	success_message.SetCharacterSize(12)
	success_message.SetString("Completed!")
	success_message.SetPosition(sf.Vector2f{10, 120})


	row_hints := make([][]int, rows)
	column_hints := make([][]int, columns)
	column_consecutive := make([]int, columns)

	// set positions, set states, and determine hints
	for y, row := range matrix {
		var row_consecutive int

		row_slice := make([]*Tile, columns)

		for x, b := range row {
			t := NewTile(texture)

			if b == StateFilled {
				row_consecutive++
				column_consecutive[x]++

				t.SetState(StateFilled)
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

			t.Sprite.SetPosition(sf.Vector2f{GridToPixelsf(x +11), GridToPixelsf(y +11)});

			row_slice[x] = t
			drawers = append(drawers, t)

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

		tiles[y] = row_slice
	}

	for y, row := range row_hints {
		for x, n := range row {
			h := NewHint(hint_texture, n)
			h.Sprite.SetPosition(sf.Vector2f{ GridToPixelsf(11 + x - len(row)), GridToPixelsf(y + 11) })
			drawers = append(drawers, h)
		}
	}

	for x, column := range column_hints {
		for y, n := range column {
			h := NewHint(hint_texture, n)
			h.Sprite.SetPosition(sf.Vector2f{ GridToPixelsf(x + 11), GridToPixelsf(11 + y - len(column)) })
			drawers = append(drawers, h)
		}
	}

	g := &Grid{
		Patterner: p,
		Drawers: drawers,
		Tiles: tiles,
		SuccessMessage: success_message,
	}

	return g
}

func (g *Grid) Logic() {
	g.Solved = g.CheckIfSolved()
}

func (g *Grid) CheckIfSolved() bool {
	for y, row := range g.Patterner.Matrix() {
		for x, b := range row {
			if b == StateFilled && g.Tiles[y][x].State != b {
				return false
			}

			if b == StateEmpty && g.Tiles[y][x].State == StateFilled {
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

	for _, d := range g.Drawers {
		d.Draw(target, renderStates)
	}
}

func (g *Grid) HandleEvent(event sf.Event) {
	switch event.Type() {
	case sf.EventTypeMouseButtonPressed:
		if sf.IsMouseButtonPressed(sf.MouseLeft) {
			g.Mode = ModeFill
		}

		if sf.IsMouseButtonPressed(sf.MouseRight) {
			g.Mode = ModeCrossOut
		}

	case sf.EventTypeMouseMoved:
		var b sf.FloatRect
		g.MousePosition = event.(sf.EventMouseMoved)
		var target_tile *Tile

highlight:
		for _, row := range g.Tiles {
			for _, tile := range row {
				b = tile.Sprite.GetGlobalBounds()

				if b.Contains(float32(g.MousePosition.X), float32(g.MousePosition.Y)) {
					target_tile = tile
					break highlight
				}
			}
		}

		if g.HighlightedTile != nil {
			g.HighlightedTile.SetHighlight(false)
		}

		if target_tile != nil {
			g.HighlightedTile = target_tile
			target_tile.SetHighlight(true)
		}

	case sf.EventTypeMouseButtonReleased:
		var b sf.FloatRect

		// TODO check if the event is in range for iteratoring
button:
		for _, row := range g.Tiles {
			for _, tile := range row {
				b = tile.Sprite.GetGlobalBounds()

				if b.Contains(float32(g.MousePosition.X), float32(g.MousePosition.Y)) {
					switch g.Mode {
					case ModeFill:
						tile.Fill()

					case ModeCrossOut:
						tile.CrossOut()
					}

					break button
				}
			}
		}
	}
}
