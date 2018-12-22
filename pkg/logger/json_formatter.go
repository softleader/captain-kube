package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

type JsonFormatter struct {
	Timestamp       bool
	TimestampLayout string
	Prefix          string
	PrettyPrint     bool
	Buffer          *bytes.Buffer
}

func (f *JsonFormatter) Format(lv Level, b []byte) ([]byte, error) {
	data := make(map[string]string)

	if f.Timestamp {
		data["time"] = time.Now().Format(f.TimestampLayout)
	}
	if p := f.Prefix; len(p) > 0 {
		data["prefix"] = f.Prefix
	}
	data["level"] = lv.String()
	data["msg"] = string(b)

	var buf *bytes.Buffer
	if f.Buffer != nil {
		buf = f.Buffer
	} else {
		buf = &bytes.Buffer{}
	}

	encoder := json.NewEncoder(buf)
	if f.PrettyPrint {
		encoder.SetIndent("", "  ")
	}
	if err := encoder.Encode(data); err != nil {
		return nil, fmt.Errorf("failed to marshal log to JSON, %v", err)
	}

	return buf.Bytes(), nil
}
