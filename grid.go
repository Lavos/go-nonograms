package nonograms

import (
	// "log"
	sf "bitbucket.org/krepa098/gosfml2"
)

type Grid struct {
	Shape *sf.RectangleShape
	Active bool

	Eventer Eventer

	Tiles [][]*Tile
}

func NewGrid (e Eventer) *Grid {
	g := &Grid{
		Eventer: e,
		Tiles: make([][]*Tile, 0),
	}

	shape, _ := sf.NewRectangleShape();
	g.Shape = shape
	g.Shape.SetPosition(sf.Vector2f{100, 100})
	g.Shape.SetSize(sf.Vector2f{100, 100})

	g.Shape.SetFillColor(sf.ColorBlack())

	mm_sub := &Subscription{
		EventType: sf.EventTypeMouseMoved,
		Callback: g.HandleHighlight,
	}

	e.Subscribe(mm_sub)

	r_sub := &Subscription{
		EventType: sf.EventTypeMouseButtonReleased,
		Callback: g.HandleRelease,
	}

	e.Subscribe(r_sub)

	return g
}

func (g *Grid) HandleHighlight(e sf.Event) {
	mm := e.(sf.EventMouseMoved)
	bounds := g.Shape.GetGlobalBounds()

	g.Highlight(bounds.Contains(float32(mm.X), float32(mm.Y)))
}

func (g *Grid) HandleRelease(e sf.Event) {
	if g.Active {
		g.Shape.SetSize(sf.Vector2f{150, 150})
	} else {
		g.Shape.SetSize(sf.Vector2f{100, 100})
	}
}

func (g *Grid) Highlight(enabled bool) {
	if enabled {
		g.Shape.SetFillColor(sf.ColorRed())
	} else {
		g.Shape.SetFillColor(sf.ColorBlack())
	};

	g.Active = enabled;
}
