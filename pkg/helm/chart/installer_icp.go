package chart

import (
	"fmt"
	"github.com/Sirupsen/logrus"
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

func (i *icpInstaller) Install(log *logrus.Logger) error {
	if err := loginBxPr(log, i.endpoint, i.username, i.password, i.account, i.skipSslValidation); err != nil {
		return err
	}

	if err := loadHelmChart(log, i.chart); err != nil {
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

func loginBxPr(log *logrus.Logger, endpoint, username, password, account string, skipSslValidation bool) error {
	args := []string{"pr", "login", "-a", format(endpoint), "-u", username, "-p", password, "-c", account}
	if skipSslValidation {
		args = append(args, "--skip-ssl-validation")
	}
	cmd := exec.Command("bx", args...)
	if log.IsLevelEnabled(logrus.DebugLevel) {
		cmd.Stdout = log.Writer()
		cmd.Stderr = log.Writer()
	}
	return cmd.Run()
}

func loadHelmChart(log *logrus.Logger, chart string) error {
	cmd := exec.Command("bx", "pr", "load-helm-chart", "--archive", chart)
	if log.IsLevelEnabled(logrus.DebugLevel) {
		cmd.Stdout = log.Writer()
		cmd.Stderr = log.Writer()
	}
	return cmd.Run()
}
