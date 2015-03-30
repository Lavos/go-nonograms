package main

import (
	"log"
	"time"
	"golang.org/x/mobile/sprite/clock"
)

func main () {
	var val float32
	var counter clock.Time = 2000
	var t1 float32 = 2000
	var t2 float32 = 100

	ticker := time.NewTicker(time.Second /60)

	for {
		select {
		case <-ticker.C:
			counter -= 100

			if t1 < t2 {
				val = (clock.EaseInOut(clock.Time(t1), clock.Time(t2), counter) * float32(t2)) + t1
			} else {
				val = (clock.EaseInOut(clock.Time(t2), clock.Time(t1), counter) * float32(t2)) - t1
			}


			log.Printf("[%d] val: %#v", counter, val)
		}
	}
}
