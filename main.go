package main

import (
	"github.com/softleader/captain-kube/app"
	"github.com/kataras/iris"
	"strconv"
	"mime/multipart"
	"github.com/kataras/iris/context"
	"github.com/softleader/captain-kube/ansible/playbook"
	"path"
	"github.com/softleader/captain-kube/sh"
	"github.com/softleader/captain-kube/ansible"
	"encoding/json"
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

	tmpl := iris.HTML("templates", ".html")
	tmpl.Reload(true)
	app.RegisterView(tmpl)

	app.StaticWeb("/", "./static")


	app.Get("/", func(ctx context.Context) {
		ctx.Redirect("/staging")
	})

	staging := app.Party("/staging")
	{
		staging.Get("/", func(ctx context.Context) {
			ctx.View("staging.html")
		})

		staging.Post("/", func(ctx iris.Context) {

			book := playbook.NewStaging()
			ctx.UploadFormFiles(args.Workspace, func(context context.Context, file *multipart.FileHeader) {
				book.Chart = path.Join(args.HostWorkspace, file.Filename)
			})
			body := ctx.GetHeader("Captain-Kube")
			json.Unmarshal([]byte(body), &book)
			opts := sh.Options{
				Ctx:     &ctx,
				Pwd:     args.Playbooks,
				Verbose: book.V(),
			}
			book.Inventory = path.Join(args.Workspace, book.Inventory)
			ansible.Play(&opts, *book)
		})
	}

	release := app.Party("/release")
	{
		release.Get("/", func(ctx context.Context) {
			ctx.View("release.html")
		})

		release.Post("/", func(ctx iris.Context) {

			book := playbook.NewStaging()
			ctx.UploadFormFiles(args.Workspace, func(context context.Context, file *multipart.FileHeader) {
				book.Chart = path.Join(args.HostWorkspace, file.Filename)
			})
			body := ctx.GetHeader("Captain-Kube")
			json.Unmarshal([]byte(body), &book)
			opts := sh.Options{
				Ctx:     &ctx,
				Pwd:     args.Playbooks,
				Verbose: book.V(),
			}
			book.Inventory = path.Join(args.Workspace, book.Inventory)
			ansible.Play(&opts, *book)
		})
	}

	return app
}
