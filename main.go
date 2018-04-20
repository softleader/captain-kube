package main

import (
	"github.com/softleader/captain-kube/app"
	"github.com/kataras/iris"
	"strconv"
	"mime/multipart"
	"github.com/kataras/iris/context"
	"github.com/softleader/captain-kube/sh"
	"github.com/softleader/captain-kube/ansible"
	"github.com/softleader/captain-kube/ansible/playbook"
	"encoding/json"
	"path"
)

func main() {
	args := app.NewArgs()

	// https://github.com/kataras/iris
	newApp(args).Run(
		iris.Addr(args.Addr+":"+strconv.Itoa(args.Port)),
		iris.WithoutVersionChecker,
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,         // enables faster json serialization and more
		iris.WithPostMaxMemory(32<<20), // with post limit at 32 MB.
	)
}

func newApp(args *app.Args) *iris.Application {
	app := iris.New()

	app.Post("/staging", func(ctx iris.Context) {

		book := playbook.NewStaging()
		ctx.UploadFormFiles(args.Workspace, func(context context.Context, file *multipart.FileHeader) {
			book.Chart = path.Join(args.HostWorkspace, file.Filename)
		})
		body := ctx.GetHeader("Captain-Kube")
		json.Unmarshal([]byte(body), &book)
		opts := sh.Options{
			Ctx:     &ctx,
			Pwd:     args.Ansible,
			Verbose: book.V(),
		}
		book.Inventory = path.Join(args.Workspace, book.Inventory)
		ansible.Play(&opts, *book)
	})

	return app
}
