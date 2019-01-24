package kubectl

import (
	"gopkg.in/yaml.v2"
	"os/exec"
	"strings"
)

func NewKubeVersion(vendor string) *KubeVersion {
	return &KubeVersion{
		Server: Info{
			GitCommit: vendor,
		},
	}
}

// KubeVersion 代表 kubectl version 的內容
type KubeVersion struct {
	Client Info `yaml:"clientVersion"`
	Server Info `yaml:"serverVersion"`
}

// Info 代表了跟版本有關的欄位
type Info struct {
	BuildDate struct {
	} `yaml:"buildDate"`
	Compiler     string `yaml:"compiler"`
	GitCommit    string `yaml:"gitCommit"`
	GitTreeState string `yaml:"gitTreeState"`
	GitVersion   string `yaml:"gitVersion"`
	GoVersion    string `yaml:"goVersion"`
	Major        string `yaml:"major"`
	Minor        string `yaml:"minor"`
	Platform     string `yaml:"platform"`
}

func (sv *Info) IsICP() bool {
	return strings.Contains(sv.GitVersion, "icp")
}

func (sv *Info) IsGCP() bool {
	// TODO 還不知道怎麼判斷
	return false
}

// Version returns the version of kubernetes server
func Version() (*KubeVersion, error) {
	args := []string{"--kubeconfig", kubeconfig, "version", "-o", "yaml"}
	cmd := exec.Command(kubectl, args...)
	b, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	kv := &KubeVersion{}
	if err = yaml.Unmarshal(b, kv); err != nil {
		return nil, err
	}
	return kv, nil
}
