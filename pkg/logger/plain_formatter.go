package logger

var (
	plain = new(PlainFormatter)
)

type PlainFormatter struct {
}

// 什麼都不做的 formatter
func (f *PlainFormatter) Format(lv Level, b []byte) ([]byte, error) {
	return b, nil
}
