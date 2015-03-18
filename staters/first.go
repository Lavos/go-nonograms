package staters

import (
	"log"
	sf "bitbucket.org/krepa098/gosfml2"
	"github.com/Lavos/nonograms/entities"
)

type First struct {
	Drawers []sf.Drawer
	Eventers []Eventer
	Logicers []Logicer
}

func NewFirst () *First {
	background_texture, _ := sf.NewTextureFromFile("../assets/grid.png", nil)
	background_texture.SetRepeated(true)
	background, _ := sf.NewSprite(background_texture)
	background.SetTextureRect(sf.IntRect{
		Top: 0,
		Left: 0,
		Width: 960,
		Height: 544,
	})

	m := entities.NewMatrix(20, 20)
	m.Randomize()

	log.Printf("Matrix: %#v", m)

	g := entities.NewGrid()
	g.Render(m)

	drawers := []sf.Drawer{ background, g }
	eventers := []Eventer{ g }
	logicers := []Logicer { g }

	return &First{
		Drawers: drawers,
		Eventers: eventers,
		Logicers: logicers,
	}
}

func (f *First) Draw(target sf.RenderTarget, renderStates sf.RenderStates) {
	for _, d := range f.Drawers {
		d.Draw(target, renderStates)
	}
}

func (f *First) HandleEvent(event sf.Event) {
	for _, e := range f.Eventers {
		e.HandleEvent(event)
	}
}

func (f *First) Logic() {
	for _, e := range f.Logicers {
		e.Logic()
	}
}
