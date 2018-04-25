package docker

import (
	"github.com/softleader/captain-kube/sh"
	"fmt"
	"path"
	"github.com/softleader/captain-kube/tgz"
	"github.com/softleader/captain-kube/helm"
	"github.com/softleader/captain-kube/tmpl"
	"github.com/softleader/captain-kube/charts"
)

const dockerPullScript = `
#!/usr/bin/env bash
{{ range $key, $value := . }}
docker pull {{ $value.Name }}
{{ end }}
exit 0
`

const retagAndPushScript = `
#!/usr/bin/env bash
{{ $registry := .Registry }}
{{ range $key, $value := .Images }}
docker tag {{ $value.Name }} {{ $registry }}/{{ $value.RemoteName }} && docker push {{ $registry }}/{{ $value.RemoteName }}
{{ end }}
exit 0
`

func Pull(opts *sh.Options, tar, tmp string) (string, string, error) {
	err := tgz.Extract(opts, tar, tmp)
	// 不確定為啥 tar 的輸出都在 err 中..
	//if err != nil {
	//	fmt.Println(err)
	//	return err
	//}

	// Render chart templates locally
	rendered, err := helm.Template(opts, tmp)
	if err != nil {
		fmt.Println(err)
		return "", "", err
	}

	images, err := charts.CollectImages(rendered)
	if err != nil {
		fmt.Println(err)
		return "", "", err
	}

	script := "docker-pull.sh"
	scriptPath := path.Join(tmp, script)
	err = tmpl.CompileTo(dockerPullScript, images, scriptPath)
	if err != nil {
		fmt.Println(err)
		return "", "", err
	}

	return script, scriptPath, nil
}

func RetagAndPush(opts *sh.Options, tar, registry, tmp string) (string, string, error) {
	err := tgz.Extract(opts, tar, tmp)
	// 不確定為啥 tar 的輸出都在 err 中..
	//if err != nil {
	//	fmt.Println(err)
	//	return err
	//}

	// Render chart templates locally
	rendered, err := helm.Template(opts, tmp)
	if err != nil {
		fmt.Println(err)
		return "", "", err
	}

	retag := Retag{
		Registry: registry,
	}
	retag.Images, err =
		charts.CollectImages(rendered)
	if err != nil {
		fmt.Println(err)
		return "", "", err
	}

	script := "retag-and-push.sh"
	scriptPath := path.Join(tmp, script)
	err = tmpl.CompileTo(retagAndPushScript, retag, scriptPath)
	if err != nil {
		fmt.Println(err)
		return "", "", err
	}

	return script, scriptPath, nil
}

type Retag struct {
	Registry string
	Images   []charts.Image
}
