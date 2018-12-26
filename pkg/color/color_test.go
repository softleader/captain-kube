package color

import (
	"fmt"
	"github.com/mgutz/ansi"
	"strconv"
	"testing"
)

func TestPick(t *testing.T) {
	n := 999
	picked := Pick(n)
	if l := len(picked); l != n {
		t.Errorf("expected pick %v color but got %v", n, l)
	}
	expectedUnique := picked[0:len(colors)]
	if c := countUnique(expectedUnique); c != len(colors) {
		t.Errorf("expected first %v of picked color should be UNIQUE", len(colors))
	}
	for i := 0; i < len(expectedUnique); i++ {
		fmt.Println(string(expectedUnique[i]([]byte(strconv.Itoa(i)))))
	}
}

func TestColor(t *testing.T) {
	phosphorize := ansi.ColorFunc("green+h:black")
	fmt.Println(phosphorize("Bring back the 80s!"))

	lime := ansi.ColorCode("green+h:black")
	reset := ansi.ColorCode("reset")
	fmt.Println(lime, "Bring back the 80s!", reset)
}

func countUnique(s []func([]byte) []byte) (count int) {
	unique := make(map[interface{}]bool)
	for v := range s {
		if _, found := unique[v]; !found {
			count++
		}
	}
	return
}
