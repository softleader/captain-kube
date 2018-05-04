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
	tmp, err := ioutil.TempDir("/tmp", "")
	if err != nil {
		ctx.StreamWriter(pipe.Println(err.Error()))
		return
	}
	defer os.RemoveAll(tmp) // clean up

	book := playbook.NewStaging()
	ctx.UploadFormFiles(tmp, func(context context.Context, file *multipart.FileHeader) {
		book.Chart = file.Filename
		book.ChartPath = path.Join(tmp, file.Filename)
	})
	body := ctx.GetHeader("Captain-Kube")
	err = json.Unmarshal([]byte(body), &book)
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
		images, err := docker.Pull(&opts, book.ChartPath, tmp)
		if err != nil {
			ctx.StreamWriter(pipe.Println(err.Error()))
			return
		}
		for _, i := range images {
			book.Images = append(book.Images, i...)
		}
	}
	_, _, err = ansible.Play(&opts, *book)
	if err != nil {
		ctx.StreamWriter(pipe.Println(err.Error()))
		return
	}
}
