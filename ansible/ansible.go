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
	V() bool
}

func Play(opts *sh.Options, book Book) (arg string, out string, err error) {
	return sh.C(opts, commandOf(book))
}

func commandOf(book Book) (command string) {
	s := []string{"ansible-playbook", "-i", book.I()}
	if len(book.T()) > 0 {
		s = append(s, "-t", "\""+strings.Join(book.T(), ",")+"\"")
	}
	s = append(s, "-e", book.E())
	if book.V() {
		s = append(s, "-v")
	}
	s = append(s, book.Yaml()...)
	command = strings.Join(s, " ")
	return
}
