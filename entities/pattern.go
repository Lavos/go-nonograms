package entities

import (
	"crypto/rand"
)

const (
	ByteEmpty = iota
	ByteFilled
	ByteCrossedOut
)

type Matrix [][]byte

func NewMatrix(rows, columns int) Matrix {
	m := make(Matrix, rows)

	for y := 0; y < rows; y++ {
		m[y] = make([]byte, columns)
	}

	return m
}

func (m Matrix) Randomize() {
	for _, row := range m {
		rand.Read(row)

		for x, b := range row {
			row[x] = b % 2
		}
	}
}
