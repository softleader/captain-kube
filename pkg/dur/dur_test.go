package dur

import (
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	if p := Parse("5"); p != time.Duration(5)*time.Second {
		t.Errorf("should get 5 seconds duration, but got %s", p)
	}
}
