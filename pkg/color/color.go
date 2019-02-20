package color

import (
	"bytes"
	"github.com/mgutz/ansi"
	"math/rand"
	"time"
)

var colors = []string{
	//ansi.Black,
	ansi.Red,
	ansi.Green,
	ansi.Yellow,
	ansi.Blue,
	ansi.Magenta,
	ansi.Cyan,
	ansi.White,
	//ansi.LightBlack,
	ansi.LightRed,
	ansi.LightGreen,
	ansi.LightYellow,
	ansi.LightBlue,
	ansi.LightMagenta,
	ansi.LightCyan,
	// ansi.LightWhite,
}

// Plain 代表什麼都不上色
var Plain = func(b []byte) []byte {
	return b
}

// Pick 隨機挑選 n 個顏色
func Pick(n int) (c []func([]byte) []byte) {
	shuffled := shuffle(colors)
	d := len(shuffled)
	for i := 0; i < n; i++ {
		style := shuffled[i%d]
		c = append(c, func(b []byte) []byte {
			if len(b) == 0 {
				return b
			}
			buf := bytes.NewBufferString(style)
			buf.Write(b)
			buf.WriteString(ansi.Reset)
			return buf.Bytes()
		})
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
