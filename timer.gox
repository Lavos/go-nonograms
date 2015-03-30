package nonograms

import (
	"time"
	"log"
	sf "bitbucket.org/krepa098/gosfml2"
)

var (
	TimerTextureRects []sf.IntRect
)

func init(){
	s := make([]sf.IntRect, 12)

	for x := 0; x < 12; x++ {
		s[x] = sf.IntRect{
			Left: x * GridSize,
			Top: 0,
			Width: GridSize,
			Height: GridSize,
		}
	}

	TimerTextureRects = s
}

type Timer struct {
	Drawers []sf.Drawer
	Texture *sf.Texture

	MinutesTens *sf.Sprite
	MinutesOnes *sf.Sprite
	SecondsTens *sf.Sprite
	SecondsOnes *sf.Sprite

	Seconds int
	StopChan chan bool
}

func NewTimer () *Timer {
	drawers := make([]sf.Drawer, 6)
	texture, _ := sf.NewTextureFromFile("../assets/timer.png", nil)

	seconds_divider, _ := sf.NewSprite(texture)
	seconds_divider.SetTextureRect(TimerTextureRects[11])
	seconds_divider.SetPosition(sf.Vector2f{ GridToPixelsf(9), GridToPixelsf(9) })
	drawers[0] = seconds_divider

	minutes_divider, _ := sf.NewSprite(texture)
	minutes_divider.SetTextureRect(TimerTextureRects[10])
	minutes_divider.SetPosition(sf.Vector2f{ GridToPixelsf(6), GridToPixelsf(9) })
	drawers[1] = minutes_divider

	seconds_ones, _ := sf.NewSprite(texture)
	seconds_ones.SetTextureRect(TimerTextureRects[0])
	seconds_ones.SetPosition(sf.Vector2f{ GridToPixelsf(8), GridToPixelsf(9) })
	drawers[2] = seconds_ones

	seconds_tens, _ := sf.NewSprite(texture)
	seconds_tens.SetTextureRect(TimerTextureRects[0])
	seconds_tens.SetPosition(sf.Vector2f{ GridToPixelsf(7), GridToPixelsf(9) })
	drawers[3] = seconds_tens

	minutes_ones, _ := sf.NewSprite(texture)
	minutes_ones.SetTextureRect(TimerTextureRects[0])
	minutes_ones.SetPosition(sf.Vector2f{ GridToPixelsf(5), GridToPixelsf(9) })
	drawers[4] = minutes_ones

	minutes_tens, _ := sf.NewSprite(texture)
	minutes_tens.SetTextureRect(TimerTextureRects[0])
	minutes_tens.SetPosition(sf.Vector2f{ GridToPixelsf(4), GridToPixelsf(9) })
	drawers[5] = minutes_tens

	return &Timer{
		Drawers: drawers,
		Texture: texture,

		MinutesTens: minutes_tens,
		MinutesOnes: minutes_ones,
		SecondsTens: seconds_tens,
		SecondsOnes: seconds_ones,

		StopChan: make(chan bool),
	}
}

func (t *Timer) Update(){
	minutes := t.Seconds / 60
	seconds := t.Seconds % 60

	minutes_tens := minutes / 10
	minutes_ones := minutes % 10

	seconds_tens := seconds / 10
	seconds_ones := seconds % 10

	t.MinutesTens.SetTextureRect(TimerTextureRects[minutes_tens])
	t.MinutesOnes.SetTextureRect(TimerTextureRects[minutes_ones])
	t.SecondsTens.SetTextureRect(TimerTextureRects[seconds_tens])
	t.SecondsOnes.SetTextureRect(TimerTextureRects[seconds_ones])
}

func (t *Timer) Start() {
	t.Seconds = 0
	ticker := time.NewTicker(time.Second)

	go func(){
TimerLoop:
		for {
			select {
			case <-ticker.C:
				t.Seconds++
				t.Update()

			case <-t.StopChan:
				break TimerLoop
			}
		}

		ticker.Stop()
		log.Printf("Ticker Stopped.")
	}()

	t.Update()
}

func (t *Timer) Stop() {
	t.StopChan <- true
}

func (t *Timer) Draw(target sf.RenderTarget, renderStates sf.RenderStates) {
	for _, d := range t.Drawers {
		d.Draw(target, renderStates)
	}
}
