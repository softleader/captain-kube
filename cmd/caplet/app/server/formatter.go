package server

import (
	"bytes"
	"fmt"
	"github.com/Sirupsen/logrus"
)

var text = new(logrus.TextFormatter)

type HostnameFormatter struct {
	Hostname string
}

func (f *HostnameFormatter) Format(entry *logrus.Entry) (b []byte, err error) {
	if b, err = text.Format(entry); err != nil {
		return
	} else {
		var buf bytes.Buffer
		buf.WriteString(fmt.Sprintf("[%s] ", f.Hostname))
		buf.Write(b)
		return buf.Bytes(), nil
	}
}
