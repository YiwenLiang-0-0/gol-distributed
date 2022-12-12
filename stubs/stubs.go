package stubs

import (
	"uk.ac.bris.cs/gameoflife/util"
)

var GameOfLifeHandler = "GolOperations.GameOfLife"

type Response struct {
	NewWorld [][]uint8
	Cell     []util.Cell
	Turn     int
}

type Request struct {
	World  [][]uint8
	Turns  int
	Height int
	Width  int
}
