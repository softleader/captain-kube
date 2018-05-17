package pipe

import (
	"io"
	"fmt"
)

const closeAndFlush = false

func Println(a ...interface{}) func(w io.Writer) bool {
	return func(w io.Writer) bool {
		fmt.Fprintln(w, a...)
		return closeAndFlush
	}
}

func Print(a ...interface{}) func(w io.Writer) bool {
	return func(w io.Writer) bool {
		fmt.Fprint(w, a...)
		return closeAndFlush
	}
}

func Printf(format string, a ...interface{}) func(w io.Writer) bool {
	return func(w io.Writer) bool {
		fmt.Fprintf(w, format, a...)
		return closeAndFlush
	}
}
