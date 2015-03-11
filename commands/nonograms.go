package main

import (
	"log"
	"github.com/Lavos/nonograms"
)

func main(){
	g := nonograms.New()
	log.Printf("%#v", g)

	g.Run()
}
