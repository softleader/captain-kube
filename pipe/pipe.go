package pipe

import (
	"io"
	"fmt"
)

func Println(a ...interface{}) func(w io.Writer) bool {
	return func(w io.Writer) bool {
		fmt.Fprintln(w, a...)
		return false
	}
}

func Print(a ...interface{}) func(w io.Writer) bool {
	return func(w io.Writer) bool {
		fmt.Fprint(w, a...)
		return false
	}
}

func Printf(format string, a ...interface{}) func(w io.Writer) bool {
	return func(w io.Writer) bool {
		fmt.Fprintf(w, format, a...)
		return false
	}
}
