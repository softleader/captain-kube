package dockerd

import (
	"context"
	"github.com/fsouza/go-dockerclient"
	"github.com/gosuri/uitable"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/helm/chart"
	"strings"
)

func Images(log *logrus.Logger, image string) ([]*chart.Image, error) {
	log.Debugf("listing docker images of %q", image)

	ctx := context.Background()

	cli, err := docker.NewClientFromEnv()
	if err != nil {
		return nil, err
	}

	// 參數準備
	options := docker.ListImagesOptions{
		Context: ctx,
		All:     true,
		Filter:  image,
	}

	i, err := cli.ListImages(options)
	if err != nil {
		return nil, err
	}

	var images []*chart.Image
	for _, img := range i {
		if len(img.RepoTags) > 0 {
			images = append(images, chart.NewImage(img.RepoTags[len(img.RepoTags)-1]))
		}
	}
	return images, nil
}

func ImagesWithTagConstraint(log *logrus.Logger, image, constraint string) ([]*chart.Image, error) {
	list, err := Images(log, image)
	if err != nil {
		return nil, err
	}
	log.Debugf("%v image(s) found", len(list))
	if len(list) == 0 {
		return list, nil
	}
	log.Debugf("checking the tag constraint: %q", constraint)
	var filtered []*chart.Image
	table := uitable.New()
	for _, i := range list {
		if ok, err := i.CheckTag(constraint); err != nil && i.Tag != constraint { // 如果 tag 就等於是 constraint, 就算沒有符合 semver2 也當成就是要刪除的目標吧
			log.Errorf("skip %q due to the tag constraint check: %s", i.String(), err)
			continue
		} else if !ok {
			table.AddRow(i.String(), "NOT match")
			continue
		}
		table.AddRow(i.String(), "match")
		filtered = append(filtered, i)
	}
	// 為了讓 caplet 在呼叫這個 function 時, 從 grpc streaming 出去時可以正常被斷行
	// 因此這邊我們切開斷行, 一行一行印出
	for _, row := range strings.Split(table.String(), "\n") {
		log.Debugf(row)
	}
	return filtered, nil
}
