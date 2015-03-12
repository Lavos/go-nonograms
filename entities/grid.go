package entities

import (
	sf "bitbucket.org/krepa098/gosfml2"
)

type Grid struct {
	Shape *sf.RectangleShape

	Tiles [][]*Tile

	MousePosition sf.EventMouseMoved
	HighlightedTile *Tile
}

func NewGrid () *Grid {
	g := &Grid{
		Tiles: make([][]*Tile, 10),
	}

	// set positions
	for x := 0; x < 10; x++ {
		row := make([]*Tile, 10)

		for y := 0; y < 10; y++ {
			t := NewTile()
			t.Shape.SetPosition(sf.Vector2f{float32((x * 25) + 100), float32((y * 25) + 100)});

			row[y] = t
		}

		g.Tiles[x] = row
	}

	shape, _ := sf.NewRectangleShape()
	g.Shape = shape
	g.Shape.SetPosition(sf.Vector2f{100, 100})
	g.Shape.SetSize(sf.Vector2f{100, 100})

	g.Shape.SetFillColor(sf.ColorBlack())

	return g
}

func (g *Grid) Draw(target sf.RenderTarget, renderStates sf.RenderStates) {
	g.Shape.Draw(target, renderStates)

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
