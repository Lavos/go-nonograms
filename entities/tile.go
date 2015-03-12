package entities

import (
	sf "bitbucket.org/krepa098/gosfml2"
)

type Tile struct {
	Shape *sf.RectangleShape

	Clicked bool
}

func NewTile() *Tile {
	shape, _ := sf.NewRectangleShape()

	shape.SetFillColor(sf.ColorBlue())
	shape.SetSize(sf.Vector2f{20, 20})

	return &Tile{
		Shape: shape,
	}
}

func (t *Tile) SetHighlight(enabled bool) {
	if t.Clicked {
		return
	}

	if enabled {
		t.Shape.SetFillColor(sf.ColorRed())
	} else {
		t.Shape.SetFillColor(sf.ColorBlack())
	}
}

func (t *Tile) Activate() {
	t.Clicked = true
	t.Shape.SetFillColor(sf.ColorYellow())
}
