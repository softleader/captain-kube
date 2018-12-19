package verbose

import (
	"fmt"
	"io"
	"os"
)

var Verbose = false

func Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error) {
	if Verbose {
		return fmt.Fprintf(w, format, a...)
	}
	return
}

func Printf(format string, a ...interface{}) (n int, err error) {
	return Fprintf(os.Stdout, format, a...)
}

func Fprint(w io.Writer, a ...interface{}) (n int, err error) {
	if Verbose {
		return fmt.Fprint(w, a...)
	}
	return
}

func Print(a ...interface{}) (n int, err error) {
	return Fprint(os.Stdout, a...)
}

func Fprintln(w io.Writer, a ...interface{}) (n int, err error) {
	if Verbose {
		return fmt.Fprintln(w, a...)
	}
	return
}

func Println(a ...interface{}) (n int, err error) {
	return Fprintln(os.Stdout, a...)
}
