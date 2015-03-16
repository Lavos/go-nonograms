package entities

import (
	sf "bitbucket.org/krepa098/gosfml2"
)

const (
	StateEmpty = iota
	StateFilled
	StateCrossedOut
)

var (
	TextureRects = []sf.IntRect{
		sf.IntRect{
			Top: 0,
			Left: 0,
			Width: GridSize,
			Height: GridSize,
		},

		sf.IntRect{
			Top: 0,
			Left: GridToPixelsi(1),
			Width: GridSize,
			Height: GridSize,
		},

		sf.IntRect{
			Top: 0,
			Left: GridToPixelsi(2),
			Width: GridSize,
			Height: GridSize,
		},
	}
)

type Tile struct {
	Sprite *sf.Sprite

	State byte
}

func NewTile(texture *sf.Texture) *Tile {
	sprite, _ := sf.NewSprite(texture)
	sprite.SetTextureRect(TextureRects[StateEmpty])

	return &Tile{
		Sprite: sprite,
	}
}

func (t *Tile) SetHighlight(enabled bool) {
	if enabled {
		t.Sprite.SetColor(sf.ColorYellow())
	} else {
		t.Sprite.SetColor(sf.ColorWhite())
	}
}

func (t *Tile) Draw(target sf.RenderTarget, renderStates sf.RenderStates) {
	return 
	t.Sprite.Draw(target, renderStates)
}

func (t *Tile) SetState(state byte) {
	t.State = state

	t.Sprite.SetTextureRect(TextureRects[state])
}

func (t *Tile) CrossOut() {
	switch t.State {
	case StateCrossedOut:
		t.SetState(StateEmpty)

	case StateFilled:
		t.SetState(StateCrossedOut)

	case StateEmpty:
		t.SetState(StateCrossedOut)
	}
}

func (t *Tile) Fill() {
	switch t.State {
	case StateEmpty:
		t.SetState(StateFilled)

	case StateFilled:
		t.SetState(StateEmpty)
	}
}
