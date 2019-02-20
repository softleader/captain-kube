package helm

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strings"
)

// Template 執行 helm template
func Template(log *logrus.Logger, chart, outputDir string, set ...string) (err error) {
	if err = ensureDirExist(outputDir); err != nil {
		return
	}
	args := []string{"template", "--output-dir", outputDir, chart}
	if len(set) > 0 {
		args = append(args, "--set", strings.Join(set, ","))
	}
	cmd := exec.Command("helm", args...)
	if log.IsLevelEnabled(logrus.DebugLevel) {
		log.Out.Write([]byte(fmt.Sprintln(strings.Join(cmd.Args, " "))))
		cmd.Stdout = log.Out
	}
	cmd.Stderr = log.Out
	if err := cmd.Start(); err != nil {
		return err
	}
	return cmd.Wait()
}

func ensureDirExist(path string) error {
	if fi, err := os.Stat(path); err != nil {
		if err = os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("could not create %s: %s", path, err)
		}
		return nil
	} else if !fi.IsDir() {
		return fmt.Errorf("%s must be a directory", path)
	}
	return nil
}
