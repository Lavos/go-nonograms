package entities

import (
	"log"
)

var (
	ExamplePattern = []byte{10,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 0, 0, 0, 1, 0, 0,
		1, 1, 0, 1, 0, 0, 0, 0, 0, 0,
		1, 1, 0, 1, 0, 1, 0, 0, 0, 0,
		1, 1, 0, 1, 0, 0, 0, 1, 0, 0,
		1, 1, 0, 1, 0, 1, 0, 0, 0, 0,
		1, 0, 1, 1, 0, 0, 0, 0, 0, 0,
		1, 0, 1, 1, 0, 0, 0, 1, 0, 0,
		1, 1, 1, 1, 0, 0, 0, 0, 0, 1,
	}
)

type Pattern struct {
	Columns int
	Rows int
	Bytes []byte
	Matrix [][]byte
}

func PatternFromBytes (p []byte) *Pattern {
	columns := int(p[0])
	tiles := p[1:len(p)]
	rows := len(tiles) / columns

	log.Printf("len(tiles): %d, columns: %d, rows: %d", len(tiles), columns, rows)

	matrix := make([][]byte, rows)

	for x := 0; x < rows; x++ {
		matrix[x] = tiles[(x * rows):((x * rows) + columns)]
	}

	log.Printf("Matrix: %#v", matrix)

	return &Pattern{
		Columns: columns,
		Rows: rows,
		Bytes: p,
		Matrix: matrix,
	}
}
