package docker

import (
	"fmt"
	"github.com/softleader/captain-kube/charts"
	"github.com/softleader/captain-kube/helm"
	"github.com/softleader/captain-kube/sh"
	"github.com/softleader/captain-kube/tgz"
)

func PullAndChangeRegistry(opts *sh.Options, tar, sourceRegistry, registry, tmp string) (images map[string][]charts.Image, err error) {
	err = tgz.Extract(opts, tar, tmp)
	// 不確定為啥 tar 的輸出都在 err 中..
	//if err != nil {
	//	fmt.Println(err)
	//	return err
	//}

	// Render chart templates locally
	rendered, err := helm.Template(opts, tmp)
	if err != nil {
		fmt.Println(err)
		return
	}

	images, err = charts.CollectImages(rendered, func(image charts.Image) bool {
		return true
	}, func(image charts.Image) charts.Image {
		if registry != "" && sourceRegistry != "" && image.Registry == sourceRegistry {
			image.Registry = registry
		}
		return image
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	return
}

func Pull(opts *sh.Options, tar, tmp string) (images map[string][]charts.Image, err error) {
	return PullAndChangeRegistry(opts, tar, "", "", tmp)
}

func Retag(opts *sh.Options, tar, sourceRegistry, tmp string) (images map[string][]charts.Image, err error) {
	err = tgz.Extract(opts, tar, tmp)
	// 不確定為啥 tar 的輸出都在 err 中..
	//if err != nil {
	//	fmt.Println(err)
	//	return err
	//}

	// Render chart templates locally
	rendered, err := helm.Template(opts, tmp)
	if err != nil {
		fmt.Println(err)
		return
	}

	images, err =
		charts.CollectImages(rendered, func(image charts.Image) bool {
			return image.Registry == sourceRegistry
		}, func(image charts.Image) charts.Image {
			return image
		})
	if err != nil {
		fmt.Println(err)
		return
	}

	return
}
