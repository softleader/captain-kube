package color

import "testing"

func TestPick(t *testing.T) {
	n := 10
	picked := Pick(n)
	if l := len(picked); l != n {
		t.Errorf("expected pick %v color but got %v", n, l)
	}

	if unique := countUnique(picked); unique != n {
		t.Errorf("expected pick %v UNIQUE color but got %v", n, unique)
	}
}

func countUnique(s []func(string) string) (count int) {
	unique := make(map[interface{}]bool)
	for v := range s {
		if _, found := unique[v]; !found {
			count++
		}
	}
	return
}
