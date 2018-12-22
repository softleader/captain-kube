package logger

const (
	ErrorLevel Level = iota
	WarnLevel
	InfoLevel
	DebugLevel
)

type Level uint32

var AllLevels = []Level{
	ErrorLevel,
	WarnLevel,
	InfoLevel,
	DebugLevel,
}

func (level Level) String() string {
	switch level {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	}
	return "unknown"
}
