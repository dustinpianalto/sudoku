package main

import (
	"context"
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"log"
	"sudokuSolver"
)

var (
	winWidth = 541
	winHeight = 541
)

func main() {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		log.Fatal(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow(
		"SUDOKU",
		sdl.WINDOWPOS_CENTERED,
		sdl.WINDOWPOS_CENTERED,
		int32(winWidth),
		int32(winHeight),
		sdl.WINDOW_SHOWN)
	if err != nil {
		log.Fatal(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		log.Fatal(err)
	}
	defer renderer.Destroy()

	tex, err := renderer.CreateTexture(
		sdl.PIXELFORMAT_ABGR8888,
		sdl.TEXTUREACCESS_STREAMING,
		int32(winWidth),
		int32(winHeight))
	if err != nil {
		log.Fatal(err)
	}
	defer tex.Destroy()


	grid := [][]int{
		{4,6,0,0,0,0,1,0,7},
		{2,0,0,4,7,0,9,0,8},
		{0,0,0,8,0,2,0,0,0},
		{0,0,0,9,4,0,0,8,0},
		{5,0,1,0,0,0,6,0,4},
		{0,8,0,0,2,5,0,0,0},
		{0,0,0,3,0,7,0,0,0},
		{7,0,6,0,8,4,0,0,2},
		{8,0,3,0,0,0,0,7,1},
	}
	origGrid := make([][]int, len(grid))
	for i := range grid {
		origGrid[i] = make([]int, len(grid[i]))
		copy(origGrid[i], grid[i])
	}
	fmt.Println(origGrid)



	pixels := make([]byte, winHeight*winWidth*4)

	keyState := sdl.GetKeyboardState()
	selectedCell := cell{
		x:    0,
		y:    0,
		size: 60,
	}
	mouseState := mouseState{
		Left:  false,
		Right: false,
		x:     0,
		y:     0,
	}
	selectedTexture, err := renderer.CreateTexture(
		sdl.PIXELFORMAT_ABGR8888,
		sdl.TEXTUREACCESS_STREAMING,
		int32(selectedCell.size),
		int32(selectedCell.size),
		)
	if err != nil {
		log.Fatal(err)
	}
	defer selectedTexture.Destroy()
	selectedBox := make([]byte, selectedCell.size*selectedCell.size*4)
	for i := 0; i < len(selectedBox); i = i + 4 {
		selectedBox[i] = 0
		selectedBox[i+1] = 0
		selectedBox[i+2] = 255
		selectedBox[i+3] = 100
	}
	err = selectedTexture.SetBlendMode(sdl.BLENDMODE_BLEND)
	if err != nil {
		log.Println(err)
	}
	selectedTexture.Update(nil,selectedBox, selectedCell.size*4)
	gameState := start
	solverCtx, solverCancel := context.WithCancel(context.Background())

	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}
		mouseState = getMouseState()
		if mouseState.Left {
			selectedCell.x, selectedCell.y = getSelectedCell(mouseState.x, mouseState.y, selectedCell.size)
		}
		if keyState[sdl.SCANCODE_0] != 0 {
			grid[selectedCell.y][selectedCell.x] = 0
		} else if keyState[sdl.SCANCODE_1] != 0 {
			grid[selectedCell.y][selectedCell.x] = 1
		} else if keyState[sdl.SCANCODE_2] != 0 {
			grid[selectedCell.y][selectedCell.x] = 2
		} else if keyState[sdl.SCANCODE_3] != 0 {
			grid[selectedCell.y][selectedCell.x] = 3
		} else if keyState[sdl.SCANCODE_4] != 0 {
			grid[selectedCell.y][selectedCell.x] = 4
		} else if keyState[sdl.SCANCODE_5] != 0 {
			grid[selectedCell.y][selectedCell.x] = 5
		} else if keyState[sdl.SCANCODE_6] != 0 {
			grid[selectedCell.y][selectedCell.x] = 6
		} else if keyState[sdl.SCANCODE_7] != 0 {
			grid[selectedCell.y][selectedCell.x] = 7
		} else if keyState[sdl.SCANCODE_8] != 0 {
			grid[selectedCell.y][selectedCell.x] = 8
		} else if keyState[sdl.SCANCODE_9] != 0 {
			grid[selectedCell.y][selectedCell.x] = 9
		}

		if keyState[sdl.SCANCODE_S] != 0 && gameState == start {
			gameState = play
			for i := range grid {
				copy(origGrid[i], grid[i])
			}
			fmt.Println(origGrid)
			go sudokuSolver.SolveSudoku(grid, 0, 0, solverCtx)
		}
		if keyState[sdl.SCANCODE_P] != 0 {
			gameState = end
			solverCancel()
		}
		if keyState[sdl.SCANCODE_R] != 0 {
			gameState = start
			for i := range origGrid {
				copy(grid[i], origGrid[i])
			}
			solverCtx, solverCancel = context.WithCancel(context.Background())
		}

		clearScreen(pixels)
		for i, row := range grid {
			for j, col := range row {
				if col == 0 {
					continue
				}
				p := pos {
					x: int32((60 * (j + 1)) - 30),
					y: int32((60 * (i + 1)) - 30),
				}
				var c sdl.Color
				if origGrid[i][j] == 0 {
					c = sdl.Color{
						R: 255,
						G: 255,
						B: 0,
						A: 255,
					}
				} else {
					c = sdl.Color{
						R: 255,
						G: 100,
						B: 0,
						A: 200,
					}
				}
				drawNumber(p, c, 8, col, pixels)
			}
		}

		tex.Update(nil, pixels, winWidth * 4)
		renderer.Copy(tex, nil, nil)
		for x := 0; x <= winWidth; x = x + 60 {
			if x % 180 == 0 {
				renderer.SetDrawColor(255,255,255,255)
			} else {
				renderer.SetDrawColor(100,100,100,255)
			}
			renderer.DrawLine(int32(x), int32(0), int32(x), int32(winHeight))
		}
		for y := 0; y <= winWidth; y = y + 60 {
			if y % 180 == 0 {
				renderer.SetDrawColor(255,255,255,255)
			} else {
				renderer.SetDrawColor(100,100,100,255)
			}
			renderer.DrawLine(int32(0), int32(y), int32(winWidth), int32(y))
		}

		renderer.Copy(selectedTexture, nil, &sdl.Rect{
			X: int32(selectedCell.x * selectedCell.size),
			Y: int32(selectedCell.y * selectedCell.size),
			W: int32(selectedCell.size),
			H: int32(selectedCell.size),
		})

		renderer.Present()

		sdl.Delay(16)
	}
}
