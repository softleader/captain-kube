package route

import (
	"mime/multipart"
	"path"
	"github.com/softleader/captain-kube/sh"
	"io/ioutil"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/softleader/captain-kube/pipe"
	"github.com/softleader/captain-kube/docker"
	"path/filepath"
	"github.com/softleader/captain-kube/tmpl"
	"os"
)

const pullScript = `#!/usr/bin/env bash
{{ range $source, $images := index . "images" }}
##---
# Source: {{ $source }}
{{- range $key, $image := $images }}
docker pull {{ $image.Registry }}/{{ $image.Name }}
{{- end }}
{{- end }}

exit 0`

const retagAndPushScript = `#!/usr/bin/env bash
{{ $registry := index . "registry" }}
{{- range $source, $images := index . "images" }}
##---
# Source: {{ $source }}
{{- range $key, $image := $images }}
docker tag {{ $image.Registry }}/{{ $image.Name }} {{ $registry }}/{{ $image.Name }} && docker push {{ $registry }}/{{ $image.Name }}
{{- end }}
{{- end }}

exit 0`

func Pull(playbooks string, ctx iris.Context) {
	tmp, err := ioutil.TempDir("/tmp", "")
	if err != nil {
		ctx.StreamWriter(pipe.Println(err.Error()))
		return
	}

	var chartPath string
	ctx.UploadFormFiles(tmp, func(context context.Context, file *multipart.FileHeader) {
		chartPath = path.Join(tmp, file.Filename)
	})
	opts := sh.Options{
		Ctx:     &ctx,
		Pwd:     playbooks,
		Verbose: true,
	}

	data := make(map[string]interface{})
	data["images"], err = docker.Pull(&opts, chartPath, tmp)
	if err != nil {
		ctx.StreamWriter(pipe.Println(err.Error()))
		return
	}

	script, err := ioutil.TempFile(tmp, "pull-")
	if err != nil {
		ctx.StreamWriter(pipe.Println(err.Error()))
		return
	}
	err = tmpl.CompileTo(pullScript, data, script.Name())
	if err != nil {
		ctx.StreamWriter(pipe.Println(err.Error()))
		return
	}
	ctx.StreamWriter(pipe.Println("generated " + script.Name()))
}

func Retag(playbooks string, ctx iris.Context) {
	tmp, err := ioutil.TempDir("/tmp", "")
	if err != nil {
		ctx.StreamWriter(pipe.Println(err.Error()))
		return
	}

	var chartPath string
	ctx.UploadFormFiles(tmp, func(context context.Context, file *multipart.FileHeader) {
		chartPath = path.Join(tmp, file.Filename)
	})
	opts := sh.Options{
		Ctx:     &ctx,
		Pwd:     playbooks,
		Verbose: true,
	}

	data := make(map[string]interface{})
	data["registry"] = ctx.Params().Get("registry")
	data["images"], err = docker.Retag(&opts, chartPath, ctx.Params().Get("source_registry"), tmp)
	if err != nil {
		ctx.StreamWriter(pipe.Println(err.Error()))
		return
	}

	script, err := ioutil.TempFile(tmp, "retag-")
	if err != nil {
		ctx.StreamWriter(pipe.Println(err.Error()))
		return
	}
	err = tmpl.CompileTo(retagAndPushScript, data, script.Name())
	if err != nil {
		ctx.StreamWriter(pipe.Println(err.Error()))
		return
	}
	ctx.StreamWriter(pipe.Println("generated " + script.Name()))
}

func DownloadScript(ctx iris.Context) {
	script := ctx.FormValue("file")
	if scriptExists(script) {
		defer os.RemoveAll(filepath.Dir(script))
		ctx.SendFile(script, filepath.Base(script))
	}
}

func scriptExists(path string) bool {
	if path == "" {
		return false
	}
	f, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	if f.IsDir() {
		return false
	}
	return true
}
