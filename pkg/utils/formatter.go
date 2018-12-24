package utils

import (
	"bytes"
	"github.com/Sirupsen/logrus"
)

// 什麼都不 format 的 formatter
type PlainFormatter struct {
}

func (f *PlainFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString(entry.Message)
	buf.WriteByte('\n')
	return buf.Bytes(), nil
}
