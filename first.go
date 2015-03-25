package nonograms

import (
	"log"
	sf "bitbucket.org/krepa098/gosfml2"
)

type First struct {
	Grid *Grid
	Timer *Timer

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

	m := NewMatrix(5, 5)
	m.Randomize()

	log.Printf("Matrix: %#v", m)

	g := NewGrid()
	g.Render(m)

	timer := NewTimer()
	timer.Start()

	drawers := []sf.Drawer{ background, g, timer }
	eventers := []Eventer{ g }
	logicers := []Logicer { g }

	return &First{
		Grid: g,
		Timer: timer,
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

	if f.Grid.Solved {
		f.Timer.Stop()
		m := NewMatrix(5, 5)
		m.Randomize()

		f.Grid.Render(m)
		f.Timer.Start()
	}
}
