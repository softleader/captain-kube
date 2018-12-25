package color

import (
	"github.com/mgutz/ansi"
	"math/rand"
	"strconv"
	"time"
)

func Pick(n int) (colors []func(string) string) {
	picked := make(map[int]func(string) string)

	for {
		color := random(0, 255)
		if _, found := picked[n]; found {
			continue
		}
		picked[n] = ansi.ColorFunc(strconv.Itoa(color))
		if len(picked)+1 == n {
			break
		}
	}
	for _, color := range picked {
		colors = append(colors, color)
	}
	return
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}
