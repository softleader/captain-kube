package ansible

import (
	"github.com/softleader/captain-kube/sh"
	"strings"
	"path"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

type Book interface {
	Yaml() []string
	I() string
	T() []string
	E() string
	V() bool
}

const daemonYaml = "daemon.yaml"

func ExtendsDefaultValues(workdir string, book interface{}) (err error) {
	raw, err := ioutil.ReadFile(path.Join(workdir, daemonYaml))
	if err != nil {
		return
	}
	yaml.Unmarshal(raw, book)
	return
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
	if book.V() {
		s = append(s, "-v")
	}
	s = append(s, book.Yaml()...)
	command = strings.Join(s, " ")
	return
}
