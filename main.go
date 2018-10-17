package main

import (
	"github.com/kataras/iris"
	"github.com/softleader/captain-kube/app"
	"strconv"
)

func main() {
	args := app.NewArgs()

	// https://github.com/kataras/iris
	app.NewApplication(args).Run(
		iris.Addr(args.Addr+":"+strconv.Itoa(args.Port)),
		iris.WithoutVersionChecker,
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,         // enables faster json serialization and more
		iris.WithPostMaxMemory(32<<20), // with post limit at 32 MB.
	)
}
