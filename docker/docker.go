package docker

import (
	"github.com/softleader/captain-kube/sh"
	"fmt"
	"github.com/softleader/captain-kube/tgz"
	"github.com/softleader/captain-kube/helm"
	"github.com/softleader/captain-kube/charts"
)

func Pull(opts *sh.Options, tar, tmp string) (images map[string][]charts.Image, err error) {
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

	images, err = charts.CollectImages(rendered, func(registry string) bool {
		return true
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	return
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
		charts.CollectImages(rendered, func(registry string) bool {
			return registry == sourceRegistry
		})
	if err != nil {
		fmt.Println(err)
		return
	}

	return
}
