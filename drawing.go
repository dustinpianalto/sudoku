package main

import "github.com/veandco/go-sdl2/sdl"

func drawNumber(pos pos, color sdl.Color, size, num int, pixels []byte) {
	startX := int(pos.x) - (size * 3) / 2
	startY := int(pos.y) - (size * 5) / 2

	for i, v := range nums[num] {
		if v == 1 {
			for y := startY; y < startY + size; y++ {
				for x := startX; x < startX + size; x++ {
					setPixel(x, y, color, pixels)
				}
			}
		}
		startX += size
		if (i + 1) % 3 == 0 {
			startY += size
			startX -= size * 3
		}
	}
}
