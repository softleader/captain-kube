package chart

import (
	"io"
	"os/exec"
)

type icpInstaller struct {
	endpoint          string
	chart             string
	username          string
	password          string
	account           string
	skipSslValidation bool
}

func (i *icpInstaller) Install(out io.Writer) error {
	args := []string{"pr", "login", "-a", i.endpoint, "-u", i.username, "-p", i.password, "-c", i.account}
	if i.skipSslValidation {
		args = append(args, "--skip-ssl-validation")
	}
	cmd := exec.Command("bx", args...)
	cmd.Stdout = out
	cmd.Stderr = out
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("bx", "pr", "load-helm-char", "--archive", i.chart)
	cmd.Stdout = out
	cmd.Stderr = out
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
