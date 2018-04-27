package route

import (
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"

	"github.com/kataras/iris/context"

	"github.com/kataras/iris"
	"github.com/softleader/captain-kube/ansible"
	"github.com/softleader/captain-kube/ansible/playbook"
	"github.com/softleader/captain-kube/docker"
	"github.com/softleader/captain-kube/pipe"
	"github.com/softleader/captain-kube/sh"
	"github.com/softleader/captain-kube/slice"
)

func Staging(workdir, playbooks string, ctx iris.Context) {
	book := playbook.NewStaging()
	ctx.UploadFormFiles(workdir, func(context context.Context, file *multipart.FileHeader) {
		book.Chart = file.Filename
		book.ChartPath = path.Join(workdir, file.Filename)
	})
	body := ctx.GetHeader("Captain-Kube")
	err := json.Unmarshal([]byte(body), &book)
	if err != nil {
		ctx.StreamWriter(pipe.Println(err.Error()))
		return
	}
	opts := sh.Options{
		Ctx:     &ctx,
		Pwd:     playbooks,
		Verbose: book.V(),
	}
	book.Inventory = path.Join(workdir, book.Inventory)
	if slice.Contains(book.Tags, "pull") {
		tmp, err := ioutil.TempDir("/tmp", "staging.")
		if err != nil {
			ctx.StreamWriter(pipe.Println(err.Error()))
			return
		}
		defer os.RemoveAll(tmp) // clean up
		book.Images, err = docker.Pull(&opts, book.ChartPath, tmp)
		if err != nil {
			ctx.StreamWriter(pipe.Println(err.Error()))
			return
		}
	}
	_, _, err = ansible.Play(&opts, *book)
	if err != nil {
		ctx.StreamWriter(pipe.Println(err.Error()))
		return
	}
}
