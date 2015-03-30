package nonograms

import (
	"log"
	"runtime"
	"time"
	sf "bitbucket.org/krepa098/gosfml2"
)

const (
	ViewWidth = 960
	ViewHeight = 544
)


var (
	CurrentMousePosition sf.EventMouseMoved
	Modes = []string{ "Random" }
	Sizes = [][]int{
		{  5,  5 },
		{ 10, 10 },
		{ 15, 15 },
		{ 20, 20 },
		{ 20, 30 },
	}
)

type Game struct {
	Window *sf.RenderWindow
	TextureManager *TextureManager

	CurrentState Stater
	States []Stater
}

func New(root string) *Game {
	start := time.Now()

	runtime.LockOSThread()
	renderWindow := sf.NewRenderWindow(sf.VideoMode{ViewWidth, ViewHeight, 32}, "Nonograms", sf.StyleClose, sf.DefaultContextSettings())

	tm := NewTextureManger(root)

	f := NewFirst(tm)

	a := time.NewTicker(time.Second * 5)
	toggle := false

	go func(){
		for {
			<-a.C
			toggle = !toggle

			if toggle {
				f.MoveView(sf.Vector2f{ 1440, 272 })
			} else {
				f.MoveView(sf.Vector2f{ 480, 272 })
			}
		}
	}()

	s := []Stater{ f }

	game := &Game{
		Window: renderWindow,
		TextureManager: tm,

		States: s,
		CurrentState: s[0],
	}

	log.Printf("New game in: %s", time.Now().Sub(start))
	return game
}

func (g *Game) Run () {
	log.Print("Running...")
	// g.Window.SetFramerateLimit(60)

	t := time.NewTicker(time.Second)
	limit := time.NewTicker(time.Second / 60)

	var frames uint64
	var fps int


	for g.Window.IsOpen(){
		select {
		case <-t.C:
			log.Printf("FPS: %d", fps)
			log.Printf("Frames: %d", frames)
			fps = 0

		case <-limit.C:
			frames++
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
