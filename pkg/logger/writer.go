package logger

import "io"

func write(out io.Writer, fmt Formatter, lv Level, b []byte) (int, error) {
	f, err := fmt.Format(lv, b)
	if err != nil {
		return 0, err
	}
	return out.Write(f)
}


