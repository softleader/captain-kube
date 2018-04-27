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

func Release(workdir, playbooks string, ctx iris.Context) {
	book := playbook.NewRelease()
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

	if slice.Contains(book.Tags, "retag") {
		tmp, err := ioutil.TempDir("/tmp", "")
		if err != nil {
			ctx.StreamWriter(pipe.Println(err.Error()))
			return
		}
		defer os.RemoveAll(tmp) // clean up
		images, err := docker.Retag(&opts, book.ChartPath, book.SourceRegistry, tmp)
		if err != nil {
			ctx.StreamWriter(pipe.Println(err.Error()))
			return
		}
		for _, i := range images {
			book.Images = append(book.Images, i...)
		}
		book.Tags = append(book.Tags, "push")
	}
	_, _, err = ansible.Play(&opts, *book)
	if err != nil {
		ctx.StreamWriter(pipe.Println(err.Error()))
		return
	}
}
