package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"math"
)

func setPixel(x, y int, c sdl.Color, pixels []byte) {
	index := (y*winWidth + x) * 4
	if index < len(pixels) - 4 && index >= 0 {
		pixels[index] = c.R
		pixels[index+1] = c.G
		pixels[index+2] = c.B
		//pixels[index+3] = c.A
	}
}

func clearScreen(pixels []byte) {
	for i := range pixels {
		pixels[i] = 0
	}
}

func getCenter() pos {
	return pos{int32(winWidth / 2), int32(winHeight / 2)}
}

func float32Lerp(a, b, pct float32) float32 {
	return a + pct * (b - a)
}

func byteLerp(b1, b2 byte, pct float32) byte {
	return byte(float32(b1) + pct * (float32(b2) - float32(b1)))
}

func colorLerp(c1, c2 sdl.Color, pct float32) sdl.Color {
	return sdl.Color{
		R: byteLerp(c1.R, c2.R, pct),
		G: byteLerp(c1.G, c2.G, pct),
		B: byteLerp(c1.B, c2.B, pct),
		A: byteLerp(c1.A, c2.A, pct),
	}
}

func clamp(min, max, v int) int {
	if v < min {
		return min
	} else if v > max {
		return max
	}
	return v
}

func getGradient(c1, c2 sdl.Color) []sdl.Color {
	var result = make([]sdl.Color, 256)

	for i := range result {
		pct := float32(i) / float32(255)
		result[i] = colorLerp(c1, c2, pct)
	}

	return result
}

func getDualGradient(c1, c2, c3, c4 sdl.Color) []sdl.Color {
	var result = make([]sdl.Color, 256)

	for i := range result {
		pct := float32(i) / float32(255)
		if pct < 0.5 {
			result[i] = colorLerp(c1, c2, pct * 2)
		} else {
			result[i] = colorLerp(c3, c4, pct * 1.5 - 0.5)
		}
	}

	return result
}

func getMouseState() mouseState {
	mouseX, mouseY, s := sdl.GetMouseState()
	state := mouseState{
		Left:  s & sdl.BUTTON_LEFT != 0,
		Right: s & sdl.BUTTON_RIGHT != 0,
		x:     mouseX,
		y:     mouseY,
	}
	return state
}

func getSelectedCell(x, y int32, size int) (int, int) {
	cellX := int(math.Floor(float64(x / int32(size))))
	cellY := int(math.Floor(float64(y / int32(size))))
	return cellX, cellY
}