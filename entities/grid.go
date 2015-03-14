package entities

import (
	"log"
	"strings"
	"strconv"
	sf "bitbucket.org/krepa098/gosfml2"
)

type Grid struct {
	Pattern *Pattern

	Texts []*sf.Text
	Tiles [][]*Tile

	MousePosition sf.EventMouseMoved
	HighlightedTile *Tile
}

func NewGrid (pattern *Pattern) *Grid {
	font, _ := sf.NewFontFromFile("TerminalVector.ttf")
	tiles := make([][]*Tile, pattern.Rows)
	texts := make([]*sf.Text, pattern.Rows)

	row_hints := make([][]int, pattern.Rows)
	column_hints := make([][]int, pattern.Columns)

	column_consecutive := make([]int, len(pattern.Matrix[0]))

	// set positions, states, and hints
	for y, row := range pattern.Matrix {
		var row_consecutive int

		row_slice := make([]*Tile, len(pattern.Matrix))

		for x, b := range row {
			t := NewTile()

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

				t.SetState(StateEmpty)
			}

			t.Shape.SetPosition(sf.Vector2f{float32((x * 25) + 100), float32((y * 25) + 100)});

			row_slice[x] = t

			if y == len(pattern.Matrix)-1 {
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

	for i, h := range row_hints {
		s := make([]string, len(h))

		for ih, hn := range h {
			s[ih] = strconv.Itoa(hn)
		}

		text, _ := sf.NewText(font)
		text.SetCharacterSize(12)
		text.SetString(strings.Join(s, " "))
		text.SetPosition(sf.Vector2f{float32(50), float32(103 + (i * 25))})

		texts[i] = text
	}

	for i, h := range column_hints {
		s := make([]string, len(h))

		for ih, hn := range h {
			s[ih] = strconv.Itoa(hn)
		}

		log.Printf("Column %d: %s", i, strings.Join(s, " "))
	}

	g := &Grid{
		Pattern: pattern,
		Tiles: tiles,
		Texts: texts,
	}

	return g
}

func (g *Grid) Logic() {
	// probably don't need to do this 60/sec
	log.Printf("Solved: %t", g.CheckIfSolved())
}

func (g *Grid) CheckIfSolved() bool {
	for y, row := range g.Pattern.Matrix {
		for x, b := range row {
			if g.Tiles[y][x].State != b {
				return false
			}
		}
	}

	return true
}

func (g *Grid) Draw(target sf.RenderTarget, renderStates sf.RenderStates) {
	for _, text := range g.Texts {
		text.Draw(target, renderStates)
	}

	for _, row := range g.Tiles {
		for _, cell := range row {
			cell.Shape.Draw(target, renderStates)
		}
	}
}

func (g *Grid) HandleEvent(event sf.Event) {
	switch event.Type() {
	case sf.EventTypeMouseMoved:
		var b sf.FloatRect
		g.MousePosition = event.(sf.EventMouseMoved)
		var target_tile *Tile

highlight:
		for _, row := range g.Tiles {
			for _, tile := range row {
				b = tile.Shape.GetGlobalBounds()

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
				b = tile.Shape.GetGlobalBounds()

				if b.Contains(float32(g.MousePosition.X), float32(g.MousePosition.Y)) {
					tile.Activate()
					break button
				}
			}
		}
	}
}
