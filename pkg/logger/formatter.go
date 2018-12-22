package logger

// formatter 會是在輸出之前, 再次的對輸出內容 format
type Formatter interface {
	Format(Level, []byte) ([]byte, error)
}
