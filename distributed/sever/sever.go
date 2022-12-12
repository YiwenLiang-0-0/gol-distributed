package main

import (
	"flag"
	"math/rand"
	"net"
	"net/rpc"
	"time"
	"uk.ac.bris.cs/gameoflife/stubs"
	"uk.ac.bris.cs/gameoflife/util"
)

func calculateNextState(Height int, Width int, world [][]uint8) [][]uint8 {
	temp := make([][]uint8, Width)

	for i := 0; i < Width; i++ {
		temp[i] = make([]uint8, Height)
	}
	// direction arrays
	dx := [8]int{-1, 1, 0, 0, -1, 1, -1, 1}
	dy := [8]int{0, 0, -1, 1, -1, 1, 1, -1}

	for baseX := 0; baseX < len(world); baseX++ {
		for baseY := 0; baseY < len(world[0]); baseY++ {
			n := 0 // alive cell
			for k := 0; k < 8; k++ {
				x := baseX + dx[k]
				y := baseY + dy[k]
				if x < 0 {
					x += len(world)
				}
				if x >= len(world) {
					x -= len(world)
				}
				if y < 0 {
					y += len(world)
				}
				if y >= len(world) {
					y -= len(world)
				}
				if world[x][y] == byte(255) {
					n += 1
				}
			}
			//fmt.Println(n)
			x, y := baseX, baseY
			switch {
			case n < 2:
				temp[x][y] = 0
			case n == 3:
				temp[x][y] = byte(255)
				//fmt.Println("alive", x, y)
			case n > 3:
				temp[x][y] = 0
			default:
				temp[x][y] = world[x][y]
			}
		}
	}
	return temp
}

func calculateAliveCells(world [][]uint8) []util.Cell {
	var res []util.Cell
	for i := 0; i < len(world); i++ {
		for j := 0; j < len(world[0]); j++ {
			if world[i][j] == uint8(255) {
				res = append(res, util.Cell{
					X: j,
					Y: i,
				})
			}
		}
	}
	return res
}

type GolOperations struct{}

func (s *GolOperations) GameOfLife(req stubs.Request, res *stubs.Response) (err error) {
	world := req.World
	for Turn := 0; Turn < req.Turns; Turn++ {
		newWorld := calculateNextState(req.Height, req.Width, world)
		world = newWorld
		// send alive cell for gui update
	}
	res.Cell = calculateAliveCells(world)
	res.Turn = req.Turns
	return
}

func main() {
	pAddr := flag.String("port", "8030", "Port to listen on")
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	err := rpc.Register(&GolOperations{})
	if err != nil {
	}
	listener, err := net.Listen("tcp", ":"+*pAddr)
	if err != nil {
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {

		}
	}(listener)
	rpc.Accept(listener)
}
