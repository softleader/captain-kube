package chart

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os/exec"
	"strings"
)

func loginBxPr(log *logrus.Logger, endpoint, username, password, account string, skipSslValidation bool) error {
	args := []string{"pr", "login", "-a", formatAccess(endpoint), "-u", username, "-p", password, "-c", account}
	if skipSslValidation {
		args = append(args, "--skip-ssl-validation")
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

func formatAccess(endpoint string) string {
	if strings.HasPrefix(endpoint, "http") {
		return endpoint
	}
	return fmt.Sprintf("https://%s:8443", endpoint)
}
