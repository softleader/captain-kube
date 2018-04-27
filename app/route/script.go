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

const pullScript = `
#!/usr/bin/env bash

{{- range $key, $value := index . "images" }}
docker pull {{ $value.Name }}
{{- end }}

exit 0
`

const retagAndPushScript = `
#!/usr/bin/env bash

{{- $registry := index . "registry" }}
{{- range $key, $value := index . "images" }}
docker tag {{ $value.Name }} {{ $registry }}/{{ $value.Name }} && docker push {{ $registry }}/{{ $value.Name }}
{{- end }}

exit 0
`

func Pull() {

}

func Retag(workdir, playbooks string, ctx iris.Context) {
	var chart, chartPath string
	ctx.UploadFormFiles(workdir, func(context context.Context, file *multipart.FileHeader) {
		chart = file.Filename
		chartPath = path.Join(workdir, file.Filename)
	})
	opts := sh.Options{
		Ctx:     &ctx,
		Pwd:     playbooks,
		Verbose: true,
	}
	tmp, err := ioutil.TempDir("/tmp", "script.")
	if err != nil {
		ctx.StreamWriter(pipe.Println(err.Error()))
		return
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
	if script != "" {
		defer os.Remove(filepath.Base(script))
		ctx.SendFile(filepath.Base(script)+".sh", script)
	}
}
