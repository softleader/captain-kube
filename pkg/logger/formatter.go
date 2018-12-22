package logger

type Formatter interface {
	Format(Level, []byte) ([]byte, error)
}
