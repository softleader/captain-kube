package tmpl

import (
	"html/template"
	"bytes"
	"os"
	"io/ioutil"
)

func CompileTo(text string, data interface{}, dest string) error {
	buf, err := compile(text, data)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(dest, buf.Bytes(), os.ModePerm)
}

func compile(text string, data interface{}) (buf bytes.Buffer, err error) {
	t := template.Must(template.New("").Parse(text))
	err = t.Execute(&buf, data)
	return
}
