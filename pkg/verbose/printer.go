package verbose

import (
	"fmt"
	"io"
	"log"
	"os"
)

var Enabled = false
var WithTime = true

func Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error) {
	if Enabled {
		if WithTime {
			log.New(w, "", log.LstdFlags).Printf(format, a...)
		} else {
			return fmt.Fprintf(w, format, a...)
		}
	}
	return
}

func Printf(format string, a ...interface{}) (n int, err error) {
	return Fprintf(os.Stdout, format, a...)
}

func Fprint(w io.Writer, a ...interface{}) (n int, err error) {
	if Enabled {
		if WithTime {
			log.New(w, "", log.LstdFlags).Print(a...)
		} else {
			return fmt.Fprint(w, a...)
		}
	}
	return
}

func Print(a ...interface{}) (n int, err error) {
	return Fprint(os.Stdout, a...)
}

func Fprintln(w io.Writer, a ...interface{}) (n int, err error) {
	if Enabled {
		if WithTime {
			log.New(w, "", log.LstdFlags).Println(a...)
		} else {
			return fmt.Fprintln(w, a...)
		}
	}
	return
}

func Println(a ...interface{}) (n int, err error) {
	return Fprintln(os.Stdout, a...)
}
