package app

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/softleader/captain-kube/app/route"
	"github.com/softleader/captain-kube/slice"
	"github.com/kataras/iris/core/router"
)

func NewApplication(args *Args) *iris.Application {
	app := iris.New()
	d := GetDaemon(args.Workdir)

	tmpl := iris.HTML("templates", ".html")
	tmpl.Reload(true)
	tmpl.AddFunc("Contains", slice.Contains)
	tmpl.AddFunc("NotContains", func(vs []string, s string) bool {
		return !slice.Contains(vs, s)
	})
	app.RegisterView(tmpl)
	app.UseGlobal(func(ctx iris.Context) {
		ctx.ViewData("daemon", d)
		ctx.ViewData("app", args)
		ctx.Next() // execute the next handler, in this case the main one.
	})

	app.StaticWeb(args.ContextPath, "./static")

	root := rootParty(app, args)
	{
		root.Get("/", func(ctx context.Context) {
			ctx.View("index.html")
		})

		testing := root.Party("/testing")
		{
			testing.Get("/", func(ctx context.Context) {
				ctx.View("testing.html")
			})

			testing.Post("/", func(ctx context.Context) {
				route.Testing(args.Workdir, args.Playbooks, ctx)
			})
		}

		staging := root.Party("/staging")
		{
			staging.Get("/", func(ctx context.Context) {
				ctx.View("staging.html")
			})

			staging.Post("/", func(ctx context.Context) {
				route.Staging(args.Workdir, args.Playbooks, ctx)
			})
		}

		production := root.Party("/production")
		{
			production.Get("/", func(ctx context.Context) {
				ctx.View("production.html")
			})

			production.Post("/", func(ctx context.Context) {
				route.Production(args.Workdir, args.Playbooks, ctx)
			})
		}

		script := root.Party("/script")
		{
			pull := script.Party("/image-pull")
			{
				pull.Get("/", func(ctx context.Context) {
					ctx.View("image-pull.html")
				})

				pull.Post("/", func(ctx context.Context) {
					route.Pull(args.Playbooks, ctx)
				})
			}

			retag := script.Party("/image-retag")
			{
				retag.Get("/", func(ctx context.Context) {
					ctx.View("image-retag.html")
				})

				retag.Post("/{source_registry:string}/{registry:string}", func(ctx context.Context) {
					route.Retag(args.Playbooks, ctx)
				})
			}

			save := script.Party("/image-save")
			{
				save.Get("/", func(ctx context.Context) {
					ctx.View("image-save.html")
				})

				save.Post("/", func(ctx context.Context) {
					route.Save(args.Playbooks, ctx)
				})
			}

			load := script.Party("/image-load")
			{
				load.Get("/", func(ctx context.Context) {
					ctx.View("image-load.html")
				})

				load.Post("/", func(ctx context.Context) {
					route.Load(args.Playbooks, ctx)
				})
			}

			script.Get("/", route.DownloadScript)
		}
	}

	return app
}

func rootParty(app *iris.Application, args *Args) (root router.Party) {
	relativePath := "/"
	if args.ContextPath != "" {
		app.Any(relativePath, func(ctx context.Context) {
			ctx.Redirect(args.ContextPath)
		})
		relativePath = args.ContextPath
	}
	root = app.Party(relativePath)
	return
}
