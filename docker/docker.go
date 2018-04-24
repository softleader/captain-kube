package docker

import (
	"github.com/softleader/captain-kube/sh"
	"io/ioutil"
	"path/filepath"
	"fmt"
	"os"
	"gopkg.in/yaml.v2"
	"github.com/softleader/captain-kube/chart"
	"path"
	"html/template"
	"bytes"
)

type PullScript struct {
	Images []string
}

const dockerPullScript = `
#!/usr/bin/env bash
{{ range $key, $value := .Images }}
docker pull {{ $value }}
{{ end }}
exit 0
`

func (s PullScript) String() string {
	t := template.Must(template.New("docker-pull").Parse(dockerPullScript))
	var buf bytes.Buffer
	t.Execute(&buf, s)
	return buf.String()
}

func PullImage(opts *sh.Options, dir string, chart string) (string, string, error) {
	// untar
	_, _, err := sh.C(opts, "tar zxvf", chart, "-C", dir, "--strip 1")
	// 不確定為啥 tar 的輸出都在 err 中..
	//if err != nil {
	//	fmt.Println(err)
	//	return err
	//}

	// Render chart templates locally
	rendered := filepath.Join(dir, "rendered")
	_, _, err = sh.C(opts, "mkdir -p", rendered, "&& helm template --output-dir", rendered, dir)
	if err != nil {
		fmt.Println(err)
		return "", "", err
	}

	s := PullScript{}

	err = filepath.Walk(rendered, func(path string, info os.FileInfo, err error) error {
		return pull(&s, path, info, err)
	})
	if err != nil {
		fmt.Println(err)
		return "", "", err
	}

	script := "docker-pull.sh"
	scriptPath := path.Join(dir, script)
	err = ioutil.WriteFile(scriptPath, []byte(s.String()), os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return "", "", err
	}

	return script, scriptPath, nil
}

func pull(script *PullScript, path string, f os.FileInfo, err error) error {
	if !f.IsDir() && filepath.Ext(path) == ".yaml" {
		fmt.Printf("pull: %s\n", path)
		in, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		fmt.Println(string(in))
		t := chart.Template{}
		yaml.Unmarshal(in, &t)
		for _, c := range t.Spec.Template.Spec.Containers {
			script.Images = append(script.Images, c.Image)
		}
	}
	return nil
}
