package route

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/softleader/captain-kube/docker"
	"github.com/softleader/captain-kube/pipe"
	"github.com/softleader/captain-kube/sh"
	"github.com/softleader/captain-kube/tmpl"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strconv"
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

const save = `#!/usr/bin/env bash
{{ $registry := index . "registry" }}
{{- range $source, $images := index . "images" }}
##---
# Source: {{ $source }}
{{- range $key, $image := $images }}
docker save -o ./{{ $image.Name }}.tar {{ $image.Registry }}/{{ $image.Name }}
{{- end }}
{{- end }}

exit 0`

const load = `#!/usr/bin/env bash
{{ $registry := index . "registry" }}
{{- range $source, $images := index . "images" }}
##---
# Source: {{ $source }}
{{- range $key, $image := $images }}
docker load -i ./{{ $image.Name }}.tar
{{- end }}
{{- end }}

exit 0`

func Pull(playbooks string, ctx iris.Context) {
	tmp, err := ioutil.TempDir(os.TempDir(), "")
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

	var sourceRegistry, registry string
	changeRegistry := ctx.FormValue("c")
	if change, _ := strconv.ParseBool(changeRegistry); change {
		sourceRegistry = ctx.FormValue("sr")
		registry = ctx.FormValue("r")
	}

	data := make(map[string]interface{})
	data["images"], err = docker.PullAndChangeRegistry(&opts, chartPath, sourceRegistry, registry, tmp)
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
	tmp, err := ioutil.TempDir(os.TempDir(), "")
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

	sourceRegistry := ctx.FormValue("sr")
	registry := ctx.FormValue("r")

	if registry == "" || sourceRegistry == "" {
		ctx.StreamWriter(pipe.Printf("Both registry: '%s' and sourceRegistry: %s are required", registry, sourceRegistry))
		return
	}

	data := make(map[string]interface{})
	data["registry"] = registry
	data["images"], err = docker.Retag(&opts, chartPath, sourceRegistry, tmp)
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

func Save(playbooks string, ctx iris.Context) {
	tmp, err := ioutil.TempDir(os.TempDir(), "")
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

	script, err := ioutil.TempFile(tmp, "save-")
	if err != nil {
		ctx.StreamWriter(pipe.Println(err.Error()))
		return
	}
	err = tmpl.CompileTo(save, data, script.Name())
	if err != nil {
		ctx.StreamWriter(pipe.Println(err.Error()))
		return
	}
	ctx.StreamWriter(pipe.Println("generated " + script.Name()))
}

func Load(playbooks string, ctx iris.Context) {
	tmp, err := ioutil.TempDir(os.TempDir(), "")
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

	script, err := ioutil.TempFile(tmp, "load-")
	if err != nil {
		ctx.StreamWriter(pipe.Println(err.Error()))
		return
	}
	err = tmpl.CompileTo(load, data, script.Name())
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
