package chart

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os/exec"
	"strings"
)

type icpDeleter struct {
	endpoint          string
	username          string
	password          string
	account           string
	skipSslValidation bool
	chartName         string
	chartVersion      string
}

// Delete 執行 ICP 的 helm chart delete
func (i *icpDeleter) Delete(log *logrus.Logger) error {
	if err := loginBxPr(log, i.endpoint, i.username, i.password, i.account, i.skipSslValidation); err != nil {
		return err
	}
	if err := deleteHelmChart(log, i.endpoint, i.chartName, i.chartVersion); err != nil {
		return err
	}
	return nil
}

func deleteHelmChart(log *logrus.Logger, endpoint, chartName, chartVersion string) error {
	args := []string{"pr", "delete-helm-chart", "--clustername", endpoint, "--name", chartName}
	if len(chartVersion) > 0 {
		args = append(args, "--version", chartVersion)
	}
	cmd := exec.Command("bx", args...)
	if log.IsLevelEnabled(logrus.DebugLevel) {
		log.Out.Write([]byte(fmt.Sprintln(strings.Join(cmd.Args, " "))))
		cmd.Stdout = log.Out
		cmd.Stderr = log.Out
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	return cmd.Wait()
}
