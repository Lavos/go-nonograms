package staters

import (
	sf "bitbucket.org/krepa098/gosfml2"
)

type Eventer interface {
	HandleEvent(sf.Event)
}

type Logicer interface {
	Logic()
}

type Stater interface {
	sf.Drawer
	Eventer
	Logicer
}
