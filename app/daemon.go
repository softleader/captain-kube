package app

import (
	"path"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

const daemonYaml = "daemon.yaml"

type Daemon struct {
	DefaultValue `yaml:"defaultValue"`
}

type DefaultValue struct {
	Inventory      string   `yaml:"inventory"`
	Tags           []string `yaml:"tags"`
	Namespace      string   `yaml:"namespace"`
	Version        string   `yaml:"version"`
	Verbose        bool     `yaml:"verbose"`
	SourceRegistry string   `yaml:"sourceRegistry"`
	Registry       string   `yaml:"registry"`
}

func GetDaemon(workdir string) (d *Daemon) {
	raw, err := ioutil.ReadFile(path.Join(workdir, daemonYaml))
	if err != nil {
		return
	}
	yaml.Unmarshal(raw, &d)
	return
}

//func (d DefaultValue) DeepCopyTo(i interface{}) {
//	b, _ := json.Marshal(d)
//	json.Unmarshal(b, i)
//}
