package chart

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
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
	if err := loginBxPr(out, i.endpoint, i.username, i.password, i.account, i.skipSslValidation); err != nil {
		return err
	}

	if err := loadHelmChart(out, i.chart); err != nil {
		return err
	}

	return nil
}

func format(endpoint string) string {
	if strings.HasPrefix(endpoint, "http") {
		return endpoint
	}
	return fmt.Sprintf("https://%s:8443", endpoint)
}

func loginBxPr(out io.Writer, endpoint, username, password, account string, skipSslValidation bool) error {
	args := []string{"pr", "login", "-a", format(endpoint), "-u", username, "-p", password, "-c", account}
	if skipSslValidation {
		args = append(args, "--skip-ssl-validation")
	}
	cmd := exec.Command("bx", args...)
	cmd.Stdout = out
	cmd.Stderr = out
	return cmd.Run()
}

func loadHelmChart(out io.Writer, chart string) error {
	cmd := exec.Command("bx", "pr", "load-helm-chart", "--archive", chart)
	cmd.Stdout = out
	cmd.Stderr = out
	return cmd.Run()
}
