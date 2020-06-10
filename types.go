package main

type pos struct {
	x, y int32
}

type gameState int
const (
	start gameState = iota
	play
	pause
	end
)

type mouseState struct {
	Left 	bool
	Right 	bool
	x 		int32
	y 		int32
}

type cell struct {
	x    int
	y    int
	size int
}