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
	"github.com/softleader/captain-kube/docker"
	"github.com/softleader/captain-kube/pipe"
	"io/ioutil"
	"os"
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
			ctx.UploadFormFiles(args.Workdir, func(context context.Context, file *multipart.FileHeader) {
				book.Chart = file.Filename
				book.ChartPath = path.Join(args.Workdir, file.Filename)
			})
			body := ctx.GetHeader("Captain-Kube")
			err := json.Unmarshal([]byte(body), &book)
			if err != nil {
				ctx.StreamWriter(pipe.Println(err.Error()))
				return
			}
			opts := sh.Options{
				Ctx:     &ctx,
				Pwd:     args.Playbooks,
				Verbose: book.V(),
			}
			book.Inventory = path.Join(args.Workdir, book.Inventory)
			if book.DockerPull {
				tmp, err := ioutil.TempDir("/tmp", ".staging")
				if err != nil {
					ctx.StreamWriter(pipe.Println(err.Error()))
					return
				}
				defer os.RemoveAll(tmp) // clean up
				book.Script, book.ScriptPath, err = docker.Pull(&opts, path.Join(args.Workdir, book.Chart), tmp)
				if err != nil {
					ctx.StreamWriter(pipe.Println(err.Error()))
					return
				}
				book.Tags = append(book.Tags, "script")
			}
			_, _, err = ansible.Play(&opts, *book)
			if err != nil {
				ctx.StreamWriter(pipe.Println(err.Error()))
				return
			}
		})
	}

	release := app.Party("/release")
	{
		release.Get("/", func(ctx context.Context) {
			ctx.View("release.html")
		})

		release.Post("/", func(ctx iris.Context) {
			book := playbook.NewRelease()
			ctx.UploadFormFiles(args.Workdir, func(context context.Context, file *multipart.FileHeader) {
				book.Chart = file.Filename
				book.ChartPath = path.Join(args.Workdir, file.Filename)
			})
			body := ctx.GetHeader("Captain-Kube")
			err := json.Unmarshal([]byte(body), &book)
			if err != nil {
				ctx.StreamWriter(pipe.Println(err.Error()))
				return
			}
			opts := sh.Options{
				Ctx:     &ctx,
				Pwd:     args.Playbooks,
				Verbose: book.V(),
			}
			book.Inventory = path.Join(args.Workdir, book.Inventory)

			tmp, err := ioutil.TempDir("/tmp", ".release")
			if err != nil {
				ctx.StreamWriter(pipe.Println(err.Error()))
				return
			}
			defer os.RemoveAll(tmp) // clean up
			book.Script, book.ScriptPath, err = docker.RetagAndPush(&opts, path.Join(args.Workdir, book.Chart), book.Registry, tmp)
			if err != nil {
				ctx.StreamWriter(pipe.Println(err.Error()))
				return
			}
			book.Tags = append(book.Tags, "script")

			_, _, err = ansible.Play(&opts, *book)
			if err != nil {
				ctx.StreamWriter(pipe.Println(err.Error()))
				return
			}
		})
	}

	return app
}
