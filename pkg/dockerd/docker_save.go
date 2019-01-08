package dockerd

import (
	"bufio"
	"context"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/helm/chart"
	"os"
	"path/filepath"
)

func Save(log *logrus.Logger, images []*chart.Image, output string, force bool) error {
	expanded, err := homedir.Expand(output)
	if err != nil {
		expanded = output
	}
	abs, err := filepath.Abs(expanded)
	if err != nil {
		return err
	}
	if _, err := os.Stat(abs); err == nil {
		if !force {
			return fmt.Errorf("could not save images to %q, because it's alreay exist", abs)
		}
		log.Debugf("removing exist file: %s\n", abs)
		if err := os.RemoveAll(abs); err != nil {
			return err
		}
	}
	if len(images) == 0 {
		log.Debugln("no images to save")
		return nil
	}
	f, err := os.Create(abs)
	if err != nil {
		return err
	}
	defer f.Close()
	ctx := context.Background()
	cli, err := docker.NewClientFromEnv()
	if err != nil {
		return err
	}
	var exports []string
	for _, image := range images {
		exports = append(exports, image.String())
	}
	w := bufio.NewWriter(f)
	defer w.Flush()
	opts := docker.ExportImagesOptions{
		Context:      ctx,
		Names:        exports,
		OutputStream: w,
	}
	log.Println("Exporting images..")
	if err := cli.ExportImages(opts); err != nil {
		return err
	}
	log.Printf("Successfully exported images and saved it to: %s\n", abs)
	return nil
}
