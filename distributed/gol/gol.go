package gol

// Params provides the details of how to run the Game of Life and which image to load.
type Params struct {
	Turns       int
	Threads     int
	ImageWidth  int
	ImageHeight int
}

// Run starts the processing of Game of Life. It should initialise channels and goroutines.
func Run(p Params, events chan<- Event, keyPresses <-chan rune) {

	//	TODO: Put the missing channels in here.
	ioCommand := make(chan ioCommand)
	ioIdle := make(chan bool)
	filenameChannel := make(chan string)
	inChannel := make(chan uint8)
	outChannel := make(chan uint8)

	ioChannels := ioChannels{
		command:  ioCommand,
		idle:     ioIdle,
		filename: filenameChannel,
		output:   outChannel,
		input:    inChannel,
	}
	go startIo(p, ioChannels)

	distributorChannels := distributorChannels{
		events:     events,
		ioCommand:  ioCommand,
		ioIdle:     ioIdle,
		ioFilename: filenameChannel,
		ioOutput:   outChannel,
		ioInput:    inChannel,
	}
	distributor(p, distributorChannels)
}
