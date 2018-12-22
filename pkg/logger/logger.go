package logger

import (
	"fmt"
	"io"
	"sync/atomic"
)

func New(output io.Writer) *Logger {
	return &Logger{
		out: output,
		lv:  InfoLevel,
		fmt: plain,
	}
}

type Logger struct {
	out io.Writer
	lv  Level
	fmt Formatter
}

func (l *Logger) Debugf(format string, args ...interface{}) (int, error) {
	if l.IsLevelEnabled(DebugLevel) {
		return write(l.out, l.fmt, l.lv, []byte(fmt.Sprintf(format, args...)))
	}
	return 0, nil
}

func (l *Logger) Infof(format string, args ...interface{}) (int, error) {
	if l.IsLevelEnabled(InfoLevel) {
		return write(l.out, l.fmt, l.lv, []byte(fmt.Sprintf(format, args...)))
	}
	return 0, nil
}

func (l *Logger) Printf(format string, args ...interface{}) (int, error) {
	return write(l.out, plain, l.lv, []byte(fmt.Sprintf(format, args...)))
}

func (l *Logger) Warnf(format string, args ...interface{}) (int, error) {
	if l.IsLevelEnabled(WarnLevel) {
		return write(l.out, l.fmt, l.lv, []byte(fmt.Sprintf(format, args...)))
	}
	return 0, nil
}

func (l *Logger) Errorf(format string, args ...interface{}) (int, error) {
	if l.IsLevelEnabled(ErrorLevel) {
		return write(l.out, l.fmt, l.lv, []byte(fmt.Sprintf(format, args...)))
	}
	return 0, nil
}

func (l *Logger) Debug(args ...interface{}) (int, error) {
	if l.IsLevelEnabled(DebugLevel) {
		return write(l.out, l.fmt, l.lv, []byte(fmt.Sprint(args...)))
	}
	return 0, nil
}

func (l *Logger) Info(args ...interface{}) (int, error) {
	if l.IsLevelEnabled(InfoLevel) {
		return write(l.out, l.fmt, l.lv, []byte(fmt.Sprint(args...)))
	}
	return 0, nil
}

func (l *Logger) Print(args ...interface{}) (int, error) {
	return write(l.out, plain, l.lv, []byte(fmt.Sprint(args...)))
}

func (l *Logger) Warn(args ...interface{}) (int, error) {
	if l.IsLevelEnabled(WarnLevel) {
		return write(l.out, l.fmt, l.lv, []byte(fmt.Sprint(args...)))
	}
	return 0, nil
}

func (l *Logger) Error(args ...interface{}) (int, error) {
	if l.IsLevelEnabled(ErrorLevel) {
		return write(l.out, l.fmt, l.lv, []byte(fmt.Sprint(args...)))
	}
	return 0, nil
}

func (l *Logger) Debugln(args ...interface{}) (int, error) {
	if l.IsLevelEnabled(DebugLevel) {
		return write(l.out, l.fmt, l.lv, []byte(fmt.Sprintln(args...)))
	}
	return 0, nil
}

func (l *Logger) Infoln(args ...interface{}) (int, error) {
	if l.IsLevelEnabled(InfoLevel) {
		return write(l.out, l.fmt, l.lv, []byte(fmt.Sprintln(args...)))
	}
	return 0, nil
}

func (l *Logger) Println(args ...interface{}) (int, error) {
	return write(l.out, plain, l.lv, []byte(fmt.Sprintln(args...)))
}

func (l *Logger) Warnln(args ...interface{}) (int, error) {
	if l.IsLevelEnabled(WarnLevel) {
		return write(l.out, l.fmt, l.lv, []byte(fmt.Sprintln(args...)))
	}
	return 0, nil
}

func (l *Logger) Errorln(args ...interface{}) (int, error) {
	if l.IsLevelEnabled(ErrorLevel) {
		return write(l.out, l.fmt, l.lv, []byte(fmt.Sprintln(args...)))
	}
	return 0, nil
}

func (l *Logger) GetLevel() Level {
	return Level(atomic.LoadUint32((*uint32)(&l.lv)))
}

func (l *Logger) GetOutput() io.Writer {
	return l.out
}

func (l *Logger) WithLevel(level Level) *Logger {
	atomic.StoreUint32((*uint32)(&l.lv), uint32(level))
	return l
}

func (l *Logger) WithVerbose(verbose bool) *Logger {
	if verbose {
		l.WithLevel(DebugLevel)
	}
	return l
}

func (l *Logger) WithFormatter(formatter Formatter) *Logger {
	l.fmt = formatter
	return l
}

func (l *Logger) IsLevelEnabled(level Level) bool {
	return l.GetLevel() >= level
}
