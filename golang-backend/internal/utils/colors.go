package utils

import (
	"fmt"
	"math/rand"
)

// HSL -> RGB helper
func hslToRgb(h, s, l float64) (int, int, int) {
	var r, g, b float64

	c := (1 - abs(2*l-1)) * s
	x := c * (1 - abs(float64(int(h/60)%2)-1))
	m := l - c/2

	switch {
	case h < 60:
		r, g, b = c, x, 0
	case h < 120:
		r, g, b = x, c, 0
	case h < 180:
		r, g, b = 0, c, x
	case h < 240:
		r, g, b = 0, x, c
	case h < 300:
		r, g, b = x, 0, c
	default:
		r, g, b = c, 0, x
	}

	return int((r + m) * 255), int((g + m) * 255), int((b + m) * 255)
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func RandomSoftHexColor() string {
	h := float64(rand.Intn(360))
	s := 0.45 + rand.Float64()*0.25 // 45% - 70%
	l := 0.55 + rand.Float64()*0.20 // 55% - 75%
	r, g, b := hslToRgb(h, s, l)
	return fmt.Sprintf("#%02X%02X%02X", r, g, b)
}
