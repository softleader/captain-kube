package arc

import (
	"fmt"
	"github.com/mholt/archiver"
	"github.com/sirupsen/logrus"
	"os"
)

// Extract 解壓縮
func Extract(log *logrus.Logger, source, destination string) (err error) {
	if err = ensureDirEmpty(destination); err != nil {
		return
	}

	log.Debugf("extracting archive to %q", destination)

	if err = archiver.Unarchive(source, destination); err != nil { // find Unarchiver by header
		var arc interface{}
		if arc, err = archiver.ByExtension(source); err != nil { // try again to find by extension
			return
		}
		return arc.(archiver.Unarchiver).Unarchive(source, destination)
	}
	return
}

func ensureDirEmpty(path string) error {
	if fi, err := os.Stat(path); err != nil {
		if err = os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("could not create %s: %s", path, err)
		}
		return nil
	} else if !fi.IsDir() {
		return fmt.Errorf("%s must be a directory", path)
	}
	// if goes here, dir already exist, so let's delete it
	return os.RemoveAll(path)
}
