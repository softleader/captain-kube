package chart

import (
	"fmt"
	"github.com/softleader/captain-kube/pkg/logger"
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

func (i *icpInstaller) Install(log *logger.Logger) error {
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

func loginBxPr(log *logger.Logger, endpoint, username, password, account string, skipSslValidation bool) error {
	args := []string{"pr", "login", "-a", format(endpoint), "-u", username, "-p", password, "-c", account}
	if skipSslValidation {
		args = append(args, "--skip-ssl-validation")
	}
	cmd := exec.Command("bx", args...)
	if log.IsLevelEnabled(logger.DebugLevel) {
		cmd.Stdout = log.GetOutput()
		cmd.Stderr = log.GetOutput()
	}
	return cmd.Run()
}

func loadHelmChart(log *logger.Logger, chart string) error {
	cmd := exec.Command("bx", "pr", "load-helm-chart", "--archive", chart)
	if log.IsLevelEnabled(logger.DebugLevel) {
		cmd.Stdout = log.GetOutput()
		cmd.Stderr = log.GetOutput()
	}
	return cmd.Run()
}
