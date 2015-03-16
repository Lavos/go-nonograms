package entities

import (
	sf "bitbucket.org/krepa098/gosfml2"
)

var (
	HintTextureRects []sf.IntRect
)

func init(){
	HintTextureRects = make([]sf.IntRect, 10)

	for x := 0; x < 10; x++ {
		HintTextureRects[x] = sf.IntRect{
			Top: 0,
			Left: GridToPixelsi(x),
			Width: GridSize,
			Height: GridSize,
		}
	}
}

type Hint struct {
	Sprite *sf.Sprite

	State byte
}

func NewHint(texture *sf.Texture, num int) *Hint {
	sprite, _ := sf.NewSprite(texture)
	sprite.SetTextureRect(HintTextureRects[num])

	return &Hint{
		Sprite: sprite,
	}
}

func (h *Hint) SetHighlight(enabled bool) {
	if enabled {
		h.Sprite.SetColor(sf.ColorYellow())
	} else {
		h.Sprite.SetColor(sf.ColorWhite())
	}
}

func (h *Hint) Draw(target sf.RenderTarget, renderStates sf.RenderStates) {
	h.Sprite.Draw(target, renderStates)
}

func (h *Hint) SetState(state byte) {
	h.State = state
	h.Sprite.SetColor(sf.ColorRed())
}

func (h *Hint) CrossOut() {
	switch h.State {
	case StateCrossedOut:
		h.SetState(StateEmpty)

	case StateEmpty:
		h.SetState(StateCrossedOut)
	}
}
