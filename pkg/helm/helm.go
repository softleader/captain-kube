package helm

import (
	"fmt"
	"github.com/softleader/captain-kube/pkg/logger"
	"os"
	"os/exec"
)

func Template(log *logger.Logger, chart, outputDir string) (err error) {
	if err = ensureDirExist(outputDir); err != nil {
		return
	}
	cmd := exec.Command("helm", "template", "--output-dir", outputDir, chart)
	if log.IsLevelEnabled(logger.DebugLevel) {
		cmd.Stdout = log.GetOutput()
		cmd.Stderr = log.GetOutput()
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
