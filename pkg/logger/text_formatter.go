package logger

import (
	"bytes"
	"fmt"
	"time"
)

type TextFormatter struct {
	Timestamp       bool
	TimestampLayout string
	Prefix          string
	Buffer          *bytes.Buffer
}

func NewTextFormatter() *TextFormatter {
	return &TextFormatter{
		Timestamp:       true,
		TimestampLayout: time.RFC3339,
	}
}

func (f *TextFormatter) WithTimestamp(timestamp bool) *TextFormatter {
	f.Timestamp = timestamp
	return f
}

func (f *TextFormatter) WithTimestampLayout(layout string) *TextFormatter {
	f.TimestampLayout = layout
	return f
}

func (f *TextFormatter) WithPrefix(prefix string) *TextFormatter {
	f.Prefix = prefix
	return f
}

func (f *TextFormatter) WithBuffer(buffer *bytes.Buffer) *TextFormatter {
	f.Buffer = buffer
	return f
}

func (f *TextFormatter) Format(lv Level, b []byte) ([]byte, error) {
	var buf *bytes.Buffer
	if f.Buffer != nil {
		buf = f.Buffer
	} else {
		buf = &bytes.Buffer{}
	}
	if f.Timestamp {
		if _, err := buf.WriteString(fmt.Sprintf("%s ", time.Now().Format(f.TimestampLayout))); err != nil {
			return b, err
		}
	}
	if p := f.Prefix; len(p) > 0 {
		if _, err := buf.WriteString(fmt.Sprintf("%s ", p)); err != nil {
			return b, err
		}
	}
	if _, err := buf.WriteString(fmt.Sprintf("%s: ", lv.String())); err != nil {
		return b, err
	}
	if _, err := buf.Write(b); err != nil {
		return b, err
	}
	return buf.Bytes(), nil
}
