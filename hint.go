package nonograms

import (
	"log"
	sf "bitbucket.org/krepa098/gosfml2"
)

var (
	HintTextureRects []sf.IntRect
)

func init(){
	HintTextureRects = make([]sf.IntRect, 20)

	for x := 0; x < 20; x++ {
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

func (h *Hint) HandleEvent(event sf.Event) {
	switch event.Type() {
	case sf.EventTypeMouseButtonReleased:
		bounds := h.Sprite.GetGlobalBounds()

		if bounds.Contains(float32(CurrentMousePosition.X), float32(CurrentMousePosition.Y)) {
			h.ToggleState()
		}
	}
}

func (h *Hint) Draw(target sf.RenderTarget, renderStates sf.RenderStates) {
	h.Sprite.Draw(target, renderStates)
}

func (h *Hint) SetState(state byte) {
	h.State = state

	switch h.State {
	case ByteEmpty:
		h.Sprite.SetColor(sf.ColorWhite())

	case ByteCrossedOut:
		h.Sprite.SetColor(sf.Color{ 255, 255, 255, 125 })
	}
}

func (h *Hint) ToggleState() {
	log.Printf("toggle state")

	switch h.State {
	case ByteCrossedOut:
		h.SetState(ByteEmpty)

	case ByteEmpty:
		h.SetState(ByteCrossedOut)
	}
}
