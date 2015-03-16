package entities

import (
	"log"
	"io"
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/base64"
)

var (
	Patterns5x5 = [1][5][5]byte{
		{
			{ 1, 0, 0, 0, 0 },
			{ 1, 0, 0, 0, 0 },
			{ 1, 0, 0, 0, 0 },
			{ 1, 0, 0, 0, 0 },
			{ 1, 0, 0, 0, 0 },
		},
	}
)

type Patterner interface {
	Matrix() Pattern
	String() string
}

type Pattern [][]byte

type Pattern5x5 [5][5]byte

func (p *Pattern5x5) Matrix() Pattern {
	m := make([][]byte, 5)

	for y, row := range p {
		row_slice := make([]byte, 5)
		copy(row_slice, row[:])
		m[y] = row_slice
	}

	return m
}

func (p *Pattern5x5) EncodeTo(w io.Writer) {
	binary.Write(w, binary.LittleEndian, p)
}

func (p *Pattern5x5) DecodeFrom(r io.Reader) {
	binary.Read(r, binary.LittleEndian, p)
}

func (p *Pattern5x5) String() string {
	var b bytes.Buffer
	p.EncodeTo(&b)
	return base64.StdEncoding.EncodeToString(b.Bytes())
}

func (p *Pattern5x5) Randomize() {
	log.Printf("Original: %#v", p)

	for y, row := range p {
		rand.Read(row[:])

		for x, b := range row {
			row[x] = b % 2
		}

		p[y] = row
	}

	log.Printf("Randomized: %#v", p)
}
