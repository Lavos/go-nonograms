package nonograms

import (
	sf "bitbucket.org/krepa098/gosfml2"
)

const (

)

type Eventer interface {
	Subscribe(*Subscription)
	Unsubscribe(*Subscription)
}

type Subscription struct {
	EventType sf.EventType
	Callback func(sf.Event)
}
