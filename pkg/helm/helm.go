package helm

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

func Template(out io.Writer, chart, outputDir string) (err error) {
	if err = ensureDirExist(outputDir); err != nil {
		return
	}
	cmd := exec.Command("helm", "template", "--output-dir", outputDir, chart)
	cmd.Stdout = out
	cmd.Stderr = out
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
