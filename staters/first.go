package staters

import (
	sf "bitbucket.org/krepa098/gosfml2"
	"github.com/Lavos/nonograms/entities"
)

type First struct {
	Drawers []sf.Drawer
	Eventers []Eventer
	Logicers []Logicer
}

func NewFirst () *First {
	g := entities.NewGrid(entities.PatternFromBytes(entities.ExamplePattern))
	drawers := []sf.Drawer{ g }
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
