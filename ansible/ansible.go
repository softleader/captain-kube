package ansible

import (
	"github.com/softleader/captain-kube/sh"
	"strings"
)

type Book interface {
	Yaml() []string
	I() string
	T() []string
	E() string
	V() string
}

func Play(opts *sh.Options, book Book) (arg string, out string, err error) {
	return sh.C(opts, commandOf(book))
}

func commandOf(book Book) (command string) {
	s := []string{"ansible-playbook", "-i", book.I()}
	if t := book.T(); len(t) > 0 {
		s = append(s, "-t", "\""+strings.Join(t, ",")+"\"")
	}
	if e := book.E(); e != "" {
		s = append(s, "-e", "'"+e+"'")
	}
	s = append(s, book.V())
	s = append(s, book.Yaml()...)
	command = strings.Join(s, " ")
	return
}
