package nonograms

import (
	"log"
	"runtime"
	"time"
	sf "bitbucket.org/krepa098/gosfml2"
)

type Game struct {
	Window *sf.RenderWindow
	Grid *Grid

	Subscribers map[sf.EventType][]*Subscription
}

func New() *Game {
	runtime.LockOSThread()
	renderWindow := sf.NewRenderWindow(sf.VideoMode{800, 600, 32}, "Events (GoSFML2)", sf.StyleDefault, sf.DefaultContextSettings())

	game := &Game{
		Window: renderWindow,
		Subscribers: make(map[sf.EventType][]*Subscription),
	}

	game.Grid = NewGrid(game)

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
			log.Printf("FPS: %d", fps)
			fps = 0

		default:
			fps++

			for event := g.Window.PollEvent(); event != nil; event = g.Window.PollEvent() {
				g.distribute(event)
			}

			g.Window.Clear(sf.Color{50, 200, 50, 0})
			g.Window.Draw(g.Grid.Shape, sf.DefaultRenderStates())
			g.Window.Display()
		}
	}
}

func (g *Game) Subscribe(sub *Subscription) {
	current_subscribers := g.Subscribers[sub.EventType]
	current_subscribers = append(current_subscribers, sub)
	g.Subscribers[sub.EventType] = current_subscribers
}

func (g *Game) Unsubscribe(sub *Subscription) {
	current_subscribers, ok := g.Subscribers[sub.EventType]

	if !ok {
		return
	}

	var index int
	for i, c := range current_subscribers {
		if c == sub {
			index = i
			break
		}
	}

	current_subscribers = append(current_subscribers[:index], current_subscribers[index+1:]...)
	g.Subscribers[sub.EventType] = current_subscribers
}

func (g *Game) distribute (event sf.Event) {
	event_type := event.Type()
	subscribers, ok := g.Subscribers[event_type]

	if !ok {
		return
	}

	for _, sub := range subscribers {
		sub.Callback(event)
	}
}
