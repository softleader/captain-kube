package helm

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"os"
	"os/exec"
)

func Template(log *logrus.Logger, chart, outputDir string) (err error) {
	if err = ensureDirExist(outputDir); err != nil {
		return
	}
	cmd := exec.Command("helm", "template", "--output-dir", outputDir, chart)
	if log.IsLevelEnabled(logrus.DebugLevel) {
		cmd.Stdout = log.Writer()
		cmd.Stderr = log.Writer()
	}
	err = cmd.Run()
	return
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
