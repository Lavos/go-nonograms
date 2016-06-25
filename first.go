package nonograms

import (
	"log"
	"time"
	"math"
	sf "bitbucket.org/krepa098/gosfml2"
	"golang.org/x/mobile/sprite/clock"
)

type First struct {
	// Timer *Timer

	PlayAreaView *sf.View
	RenderTarget sf.RenderTarget

	Drawers []sf.Drawer
	Eventers []Eventer
	ViewEventers []ViewEventer
	Logicers []Logicer
}

func NewFirst (tm *TextureManager, rt sf.RenderTarget) *First {
	m := NewMatrix(5, 5)
	m.Randomize()

	view := sf.NewViewFromRect(sf.FloatRect{0, 0, ViewWidth, ViewHeight})

	drawers := make([]sf.Drawer, 0)
	eventers := make([]Eventer, 0)
	view_eventers := make([]ViewEventer, 0)
	logicers := make([]Logicer, 0)

	log.Printf("Matrix: %#v", m)

	for y, _ := range Modes {
		for x, size := range Sizes {
			s := NewScene(tm, sf.Vector2f{ float32(x * ViewWidth), float32(y * ViewHeight) }, size[0], size[1])
			s.SetGoalMatrix(m)
			drawers = append(drawers, s)
			eventers = append(eventers, s)
			view_eventers = append(view_eventers, s)
			logicers = append(logicers, s)
		}
	}

	// timer := NewTimer()
	// timer.Start()

	return &First{
		// Timer: timer,
		PlayAreaView: view,
		RenderTarget: rt,

		Drawers: drawers,
		Eventers: eventers,
		ViewEventers: view_eventers,
		Logicers: logicers,
	}
}

func (f *First) MoveView(targetOrigin sf.Vector2f) {
	baseOrigin := f.PlayAreaView.GetCenter()
	diffOrigin := targetOrigin.Minus(baseOrigin)

	log.Printf("base %#v", baseOrigin)
	log.Printf("target %#v", targetOrigin)
	log.Printf("diff %#v", diffOrigin)

	delta_x := float32(math.Abs(float64(diffOrigin.X)))
	delta_y := float32(math.Abs(float64(diffOrigin.Y)))

	log.Printf("deltas: x %f y %f", delta_x, delta_y)

	go func(){
		var counter_x, counter_y clock.Time

		if diffOrigin.X > 0 {
			counter_x = clock.Time(baseOrigin.X)
		} else {
			counter_x = clock.Time(targetOrigin.X)
		}

		if diffOrigin.Y > 0 {
			counter_y = clock.Time(baseOrigin.Y)
		} else {
			counter_y = clock.Time(targetOrigin.Y)
		}

		ticker := time.NewTicker(time.Second / 60)

		var x, y float32

loop:
		for {
			select {
			case <-ticker.C:
				counter_x += clock.Time(delta_x / 60)
				counter_y += clock.Time(delta_y / 60)

				if diffOrigin.X > 0 {
					x = (clock.EaseInOut(clock.Time(baseOrigin.X), clock.Time(targetOrigin.X), counter_x) * delta_x) + baseOrigin.X
				} else {
					x = baseOrigin.X - (clock.EaseInOut(clock.Time(targetOrigin.X), clock.Time(baseOrigin.X), counter_x) * delta_x)
				}

				if diffOrigin.Y > 0 {
					y = (clock.EaseInOut(clock.Time(baseOrigin.Y), clock.Time(targetOrigin.Y), counter_y) * delta_y) + baseOrigin.Y
				} else {
					y = baseOrigin.Y - (clock.EaseInOut(clock.Time(targetOrigin.Y), clock.Time(baseOrigin.Y), counter_y) * delta_y)
				}

				// log.Printf("x %f y %f", x, y)

				f.PlayAreaView.SetCenter(sf.Vector2f{x, y})

				if x == targetOrigin.X && y == targetOrigin.Y {
					break loop
				}
			}
		}

		ticker.Stop()
	}()
}

func (f *First) Draw(target sf.RenderTarget, renderStates sf.RenderStates) {
	target.SetView(f.PlayAreaView)

	for _, d := range f.Drawers {
		d.Draw(target, renderStates)
	}

	target.SetView(target.GetDefaultView())
}

func (f *First) HandleEvent(event sf.Event) {
	for _, e := range f.Eventers {
		e.HandleEvent(event)
	}

	coords := f.RenderTarget.MapPixelToCoords(sf.Vector2i{ CurrentMousePosition.X, CurrentMousePosition.Y }, f.PlayAreaView)

	for _, e := range f.ViewEventers {
		e.HandleViewEvent(event, coords)
	}
}

func (f *First) Logic() {
	for _, e := range f.Logicers {
		e.Logic()
	}

	/* if f.Scene.Solved {
		// f.Timer.Stop()
		m := NewMatrix(5, 5)
		m.Randomize()

		f.Scene.Render(m)
		// f.Timer.Start()
	} */
}
