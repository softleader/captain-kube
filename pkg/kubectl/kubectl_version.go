package kubectl

import (
	"gopkg.in/yaml.v2"
	"os/exec"
	"strings"
)

// KubeVersion 代表 kubectl version 的結果
type KubeVersion struct {
	ClientVersion ClientVersion `yaml:"clientVersion"`
	ServerVersion ServerVersion `yaml:"serverVersion"`
}

// ClientVersion 代表 kubectl version 的 client 版本
type ClientVersion struct {
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

// ServerVersion 代表 kubectl version 的 server 版本
type ServerVersion struct {
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

func (sv *ServerVersion) IsICP() bool {
	return strings.Contains(sv.GitVersion, "icp")
}

func (sv *ServerVersion) IsGCP() bool {
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
