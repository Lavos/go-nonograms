package main

import (
	"log"
)

var (
	a = [5]byte{0, 0, 0, 1, 1}
	b = [5]byte{0, 0, 0, 1, 1}
)

func main () {
	log.Printf("%#v", make([][]byte, 5))
}
