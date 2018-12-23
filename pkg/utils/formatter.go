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

// 固定在 logrus.TextFormatter 之前加入一段 prefix 的 formatter
type PrefixFormatter struct {
	Prefix string
	logrus.TextFormatter
}

func (f *PrefixFormatter) Format(entry *logrus.Entry) (b []byte, err error) {
	if b, err = f.Format(entry); err != nil {
		return
	} else {
		var buf bytes.Buffer
		buf.WriteString(f.Prefix)
		buf.Write(b)
		return buf.Bytes(), nil
	}
}
