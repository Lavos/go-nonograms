package nonograms

import (
	sf "bitbucket.org/krepa098/gosfml2"
)

type Eventer interface {
	HandleEvent(sf.Event)
}

type ViewEventer interface {
	Eventer

	HandleViewEvent(sf.Event, sf.Vector2f)
}

type Logicer interface {
	Logic()
}

type Stater interface {
	sf.Drawer
	Eventer
	Logicer
}
