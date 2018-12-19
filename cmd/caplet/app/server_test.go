package app

import (
	"testing"
)

func TestRetrieveServer(t *testing.T) {
	if _, err := retrieveServer("non-exist"); err == nil {
		t.Error("error must exist")
	}
}
