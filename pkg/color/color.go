package color

import (
	"github.com/mgutz/ansi"
	"math/rand"
	"time"
)

var colors = []string{
	"black",
	"red",
	"green",
	"yellow",
	"blue",
	"magenta",
	"cyan",
	"white",
	//"black+h",
	"red+h",
	"green+h",
	"yellow+h",
	"blue+h",
	"magenta+h",
	"cyan+h",
	//"white+h",
}

func Pick(n int) (c []func(string) string) {
	shuffled := shuffle(colors)
	d := len(shuffled)
	for i := 0; i < n; i++ {
		code := shuffled[i%d]
		c = append(c, ansi.ColorFunc(code))
	}
	return
}

func shuffle(a []string) (b []string) {
	rand.Seed(time.Now().UTC().UnixNano())
	b = append(a[:0:0], a...)
	rand.Shuffle(len(b), func(i, j int) {
		b[i], b[j] = b[j], b[i]
	})
	return
}
