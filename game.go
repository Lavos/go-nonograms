package nonograms

import (
	"log"
	"runtime"
	"time"
	sf "bitbucket.org/krepa098/gosfml2"
)

var (
	CurrentMousePosition sf.EventMouseMoved
)

type Game struct {
	Window *sf.RenderWindow

	CurrentState Stater
	States []Stater
}

func New() *Game {
	start := time.Now()

	runtime.LockOSThread()
	renderWindow := sf.NewRenderWindow(sf.VideoMode{960, 544, 32}, "Nonograms", sf.StyleClose, sf.DefaultContextSettings())

	s := make([]Stater, 1)
	s[0] = NewFirst()

	game := &Game{
		Window: renderWindow,
		States: s,
		CurrentState: s[0],
	}

	log.Printf("New game in: %s", time.Now().Sub(start))
	return game
}

func (g *Game) Run () {
	log.Print("Running...")
	g.Window.SetFramerateLimit(60)

	t := time.NewTicker(time.Second)
	var fps int

	for g.Window.IsOpen(){
		select {
		case <-t.C:
			// log.Printf("FPS: %d", fps)
			fps = 0

		default:
			fps++

			for event := g.Window.PollEvent(); event != nil; event = g.Window.PollEvent() {
				switch event.Type() {
				case sf.EventTypeClosed:
					g.Window.Close()

				case sf.EventTypeMouseMoved:
					CurrentMousePosition = event.(sf.EventMouseMoved)

				default:
					g.CurrentState.HandleEvent(event)
				}
			}

			g.CurrentState.Logic()

			g.Window.Clear(sf.Color{50, 200, 50, 0})
			g.Window.Draw(g.CurrentState, sf.DefaultRenderStates())
			g.Window.Display()
		}
	}
}
