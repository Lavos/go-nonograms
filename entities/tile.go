package entities

import (
	sf "bitbucket.org/krepa098/gosfml2"
)

const (
	StateEmpty = iota
	StateFilled
	StateCrossedOut
)


type Tile struct {
	Shape *sf.RectangleShape

	State byte
}

func NewTile() *Tile {
	shape, _ := sf.NewRectangleShape()

	shape.SetSize(sf.Vector2f{20, 20})

	return &Tile{
		Shape: shape,
	}
}

func (t *Tile) SetHighlight(enabled bool) {

}

func (t *Tile) SetState(state byte) {
	t.State = state

	switch state {
	case StateEmpty:
		t.Shape.SetFillColor(sf.ColorWhite())

	case StateFilled:
		t.Shape.SetFillColor(sf.ColorBlack())

	case StateCrossedOut:
		t.Shape.SetFillColor(sf.ColorRed())
	}
}

func (t *Tile) Activate() {
	switch t.State {
	case StateEmpty:
		t.SetState(StateFilled)

	case StateFilled:
		t.SetState(StateEmpty)
	}
}
