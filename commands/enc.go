package main

import (
	"log"
	// "encoding/binary"
	// "encoding/base64"
)

type Pattern [][]byte

type Patterner interface {
	Matrix() Pattern
	// String() string
}

type Pattern5x5 [5][5]byte

func (p Pattern5x5) Matrix() Pattern {
	log.Printf("%#v", p[:])

	m := make([][]byte, 5)

	for y, row := range p {
		m[y] = row[:]
	}

	return m
}


func getMatrix(p Patterner) {
	m := p.Matrix()
	log.Printf("%#v", m)
}

func main () {
	p := Pattern5x5{
		{ 1, 0, 0, 0, 0 },
		{ 1, 0, 0, 0, 0 },
		{ 1, 0, 0, 0, 0 },
		{ 1, 0, 0, 0, 0 },
		{ 1, 0, 0, 0, 0 },
	}

	getMatrix(p)
}
