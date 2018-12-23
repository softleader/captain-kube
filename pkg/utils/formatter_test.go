package utils

import (
	"bytes"
	"github.com/Sirupsen/logrus"
	"testing"
)

func TestPlainFormatter_Format(t *testing.T) {
	log := logrus.New()
	log.SetFormatter(&PlainFormatter{})
	b := bytes.NewBuffer(nil)
	log.SetOutput(b)
	log.Println("123456789")
	if b.String() != "123456789\n" {
		t.Error("out should be 123456789\n")
	}
}