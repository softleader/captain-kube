package utils

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"strings"
	"fmt"
)

var ln = fmt.Sprintln()

// 什麼都不 format 的 formatter
type PlainFormatter struct {
}

func (f *PlainFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var buf *bytes.Buffer
	if entry.Buffer != nil {
		buf = entry.Buffer
	} else {
		buf = &bytes.Buffer{}
	}
	if entry.Message != "" {
		buf.WriteString(entry.Message)
	}
	if !strings.HasSuffix(entry.Message, ln) {
		buf.WriteString(ln)
	}
	return buf.Bytes(), nil
}
