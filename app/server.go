package app

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/softleader/captain-kube/app/route"
	"github.com/softleader/captain-kube/ansible/playbook"
	"github.com/softleader/captain-kube/ansible"
	"github.com/softleader/captain-kube/slice"
)

func NewApplication(args *Args) *iris.Application {
	app := iris.New()

	tmpl := iris.HTML("templates", ".html")
	tmpl.Reload(true)
	tmpl.AddFunc("Contains", slice.Contains)
	app.RegisterView(tmpl)

	app.StaticWeb("/", "./static")

	app.Get("/", func(ctx context.Context) {
		ctx.View("index.html")
	})

	testing := app.Party("/testing")
	{
		testing.Get("/", func(ctx context.Context) {
			dft := playbook.NewTesting()
			ansible.ExtendsDefaultValues(args.Workdir, dft)
			ctx.ViewData("dft", dft)
			ctx.View("testing.html")
		})

		testing.Post("/", func(ctx context.Context) {
			route.Testing(args.Workdir, args.Playbooks, ctx)
		})
	}

	staging := app.Party("/staging")
	{
		staging.Get("/", func(ctx context.Context) {
			dft := playbook.NewStaging()
			ansible.ExtendsDefaultValues(args.Workdir, dft)
			ctx.ViewData("dft", dft)
			ctx.View("staging.html")
		})

		staging.Post("/", func(ctx context.Context) {
			route.Staging(args.Workdir, args.Playbooks, ctx)
		})
	}

	production := app.Party("/production")
	{
		production.Get("/", func(ctx context.Context) {
			dft := playbook.NewProduction()
			ansible.ExtendsDefaultValues(args.Workdir, dft)
			ctx.ViewData("dft", dft)
			ctx.View("production.html")
		})

		production.Post("/", func(ctx context.Context) {
			route.Production(args.Workdir, args.Playbooks, ctx)
		})
	}

	script := app.Party("/script")
	{
		pull := script.Party("/pull")
		{
			pull.Get("/", func(ctx context.Context) {
				ctx.View("pull.html")
			})

			pull.Post("/", func(ctx context.Context) {
				route.Pull(args.Workdir, args.Playbooks, ctx)
			})
		}

		retag := script.Party("/retag")
		{
			retag.Get("/", func(ctx context.Context) {
				ctx.View("retag.html")
			})

			retag.Post("/{source_registry:string}/{registry:string}", func(ctx context.Context) {
				route.Retag(args.Workdir, args.Playbooks, ctx)
			})
		}

		script.Get("/", route.DownloadScript)
	}

	return app
}
