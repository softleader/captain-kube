package verbose

import (
	"fmt"
	"io"
	"log"
	"os"
)

var Enabled = false
var WithTime = true

type Printer struct {
	enable, withTime bool
}

func Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error) {
	return Printer{Enabled, WithTime}.Fprintf(w, format, a...)
}

func Printf(format string, a ...interface{}) (n int, err error) {
	return Printer{Enabled, WithTime}.Printf(format, a...)
}

func Fprint(w io.Writer, a ...interface{}) (n int, err error) {
	return Printer{Enabled, WithTime}.Fprint(w, a...)
}

func Print(a ...interface{}) (n int, err error) {
	return Printer{Enabled, WithTime}.Print(a...)
}

func Fprintln(w io.Writer, a ...interface{}) (n int, err error) {
	return Printer{Enabled, WithTime}.Fprintln(w, a...)
}

func Println(a ...interface{}) (n int, err error) {
	return Printer{Enabled, WithTime}.Println(a...)
}

func (p *Printer) Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error) {
	if Enabled {
		if WithTime {
			log.New(w, "", log.LstdFlags).Printf(format, a...)
		} else {
			return fmt.Fprintf(w, format, a...)
		}
	}
	return
}

func (p *Printer) Printf(format string, a ...interface{}) (n int, err error) {
	return p.Fprintf(os.Stdout, format, a...)
}

func (p *Printer) Fprint(w io.Writer, a ...interface{}) (n int, err error) {
	if Enabled {
		if WithTime {
			log.New(w, "", log.LstdFlags).Print(a...)
		} else {
			return fmt.Fprint(w, a...)
		}
	}
	return
}

func (p *Printer) Print(a ...interface{}) (n int, err error) {
	return p.Fprint(os.Stdout, a...)
}

func (p *Printer) Fprintln(w io.Writer, a ...interface{}) (n int, err error) {
	if Enabled {
		if WithTime {
			log.New(w, "", log.LstdFlags).Println(a...)
		} else {
			return fmt.Fprintln(w, a...)
		}
	}
	return
}

func (p *Printer) Println(a ...interface{}) (n int, err error) {
	return p.Fprintln(os.Stdout, a...)
}
