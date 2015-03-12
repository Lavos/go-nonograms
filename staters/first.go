package staters

import (
	sf "bitbucket.org/krepa098/gosfml2"
	"github.com/Lavos/nonograms/entities"
)

type First struct {
	Drawers []sf.Drawer
	Eventers []Eventer
}

func NewFirst () *First {
	g := entities.NewGrid()
	drawers := []sf.Drawer{ g }
	eventers := []Eventer{ g }

	return &First{
		Drawers: drawers,
		Eventers: eventers,
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
