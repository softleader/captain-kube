package logger

var (
	plain = new(PlainFormatter)
)

type PlainFormatter struct {
}

func (f *PlainFormatter) Format(lv Level, b []byte) ([]byte, error) {
	return b, nil
}
