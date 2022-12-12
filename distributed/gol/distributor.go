package gol

import (
	"fmt"
	"log"
	"net/rpc"
	"uk.ac.bris.cs/gameoflife/stubs"
)

type distributorChannels struct {
	events     chan<- Event
	ioCommand  chan<- ioCommand
	ioIdle     <-chan bool
	ioFilename chan<- string
	ioOutput   chan<- uint8
	ioInput    <-chan uint8
}

// distributor divides the work between workers and interacts with other goroutines.
func distributor(p Params, c distributorChannels) {
	// TODO: Create a 2D slice to store the world.
	filename := fmt.Sprintf("%vx%v", p.ImageWidth, p.ImageHeight)
	c.ioCommand <- ioInput
	c.ioFilename <- filename
	world := make([][]uint8, p.ImageWidth)
	for i := 0; i < p.ImageWidth; i++ {
		world[i] = make([]uint8, p.ImageHeight)
	}
	//fmt.Println("read c.ioInput")
	for x := 0; x < p.ImageWidth; x++ {
		for y := 0; y < p.ImageHeight; y++ {
			world[x][y] = <-c.ioInput
		}
	}

	// TODO: Execute all turns of the Game of Life.
	//server := flag.String("server", "127.0.0.1:8030", "IP:port string to connect to as server")
	server := "3.86.88.141:8030"
	//flag.Parse()
	client, err := rpc.Dial("tcp", server)
	if err != nil {
		log.Println("rpc Dial error ", err)
	}
	defer client.Close()
	response := makeCall(client, world, p)
	alive := response.Cell
	turn := response.Turn
	for _, cell := range alive {
		c.events <- CellFlipped{
			CompletedTurns: turn,
			Cell:           cell,
		}
	}
	// complete one turn, render gui frame
	c.events <- TurnComplete{CompletedTurns: turn}

	// TODO: Report the final state using FinalTurnCompleteEvent.
	c.events <- FinalTurnComplete{
		CompletedTurns: turn,
		Alive:          alive,
	}

	// Make sure that the Io has finished any output before exiting.
	c.ioCommand <- ioCheckIdle
	<-c.ioIdle

	c.events <- StateChange{turn, Quitting}

	// Close the channel to stop the SDL goroutine gracefully. Removing may cause deadlock.
	close(c.events)
}

func makeCall(client *rpc.Client, world [][]uint8, params Params) *stubs.Response {
	fmt.Println("send call")
	request := stubs.Request{World: world, Turns: params.Turns, Height: params.ImageHeight, Width: params.ImageWidth}
	response := new(stubs.Response)
	err := client.Call(stubs.GameOfLifeHandler, request, response)
	if err != nil {
		log.Println("error: ", err)
	}
	return response
}
