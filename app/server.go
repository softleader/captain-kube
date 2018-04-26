package app

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/softleader/captain-kube/app/route"
)

func NewApplication(args *Args) *iris.Application {
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

		staging.Post("/", func(ctx context.Context) {
			route.Staging(args.Workdir, args.Playbooks, ctx)
		})
	}

	release := app.Party("/release")
	{
		release.Get("/", func(ctx context.Context) {
			ctx.View("release.html")
		})

		release.Post("/", func(ctx context.Context) {
			route.Release(args.Workdir, args.Playbooks, ctx)
		})
	}

	pull := app.Party("/pull")
	{
		pull.Get("/", func(ctx context.Context) {
			ctx.View("pull.html")
		})

		pull.Post("/", route.Pull)
	}

	retag := app.Party("/retag")
	{
		retag.Get("/", func(ctx context.Context) {
			ctx.View("retag.html")
		})

		retag.Post("/", route.Retag)
	}

	return app
}
