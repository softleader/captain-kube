package charts

import (
	"testing"
)

func TestBefore(t *testing.T) {
	expected := "hub.softleader.com.tw"
	actual := before("hub.softleader.com.tw/captain-kube:v10.0.0", "/")
	if expected != actual {
		t.Errorf("Expected '%s', but got '%v'", expected, actual)
	}
}

func TestAfter(t *testing.T) {
	expected := "captain-kube:v10.0.0"
	actual := after("hub.softleader.com.tw/captain-kube:v10.0.0", "/")
	if expected != actual {
		t.Errorf("Expected '%s', but got '%v'", expected, actual)
	}
}
