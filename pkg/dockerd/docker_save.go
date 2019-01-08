package dockerd

import (
	"bufio"
	"context"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/helm/chart"
	"os"
)

func Save(log *logrus.Logger, images []*chart.Image, output string) error {
	if _, err := os.Stat(output); err == nil {
		return fmt.Errorf("could not save images to %q, because it's alreay exist", output)
	}
	if len(images) == 0 {
		log.Debugln("no images to save")
		return nil
	}
	f, err := os.Create(output)
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
	log.Printf("saving images to %q:", output)
	for _, image := range images {
		n := image.String()
		log.Println(n)
		exports = append(exports, n)
	}
	w := bufio.NewWriter(f)
	defer w.Flush()
	opts := docker.ExportImagesOptions{
		Context:      ctx,
		Names:        exports,
		OutputStream: w,
	}
	return cli.ExportImages(opts)
}
