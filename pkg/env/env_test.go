package env

import (
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestLookup(t *testing.T) {
	key, value := "hello", "world"
	if v := Lookup(key, value); v != value {
		t.Errorf("should be %q, but got %q", value, v)
	}
	os.Setenv(key, value)
	if v := Lookup(key, strings.Repeat(value, 2)); v != value {
		t.Errorf("should be %q, but got %q", value, v)
	}
}

func TestLookupBool(t *testing.T) {
	key := "hello"
	value := false
	if v := LookupBool(key, value); v != value {
		t.Errorf("should be %v, but got %v", value, v)
	}
	os.Setenv(key, strconv.FormatBool(value))
	if v := LookupBool(key, !value); v != value {
		t.Errorf("should be %v, but got %v", value, v)
	}
}

func TestLookupInt(t *testing.T) {
	key := "hello"
	value := 1
	if v := LookupInt(key, value); v != value {
		t.Errorf("should be %q, but got %q", value, v)
	}
	os.Setenv(key, strconv.Itoa(value))
	if v := LookupInt(key, value*10); v != value {
		t.Errorf("should be %v, but got %v", value, v)
	}
}
