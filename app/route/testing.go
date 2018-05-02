package route

import (
	"path"

	"github.com/kataras/iris"
	"github.com/softleader/captain-kube/ansible"
	"github.com/softleader/captain-kube/ansible/playbook"
	"github.com/softleader/captain-kube/pipe"
	"github.com/softleader/captain-kube/sh"
)

func Testing(workdir, playbooks string, ctx iris.Context) {
	book := playbook.NewTesting()
	opts := sh.Options{
		Ctx:     &ctx,
		Pwd:     playbooks,
		Verbose: book.V(),
	}
	book.Inventory = path.Join(workdir, book.Inventory)
	_, _, err := ansible.Play(&opts, *book)
	if err != nil {
		ctx.StreamWriter(pipe.Println(err.Error()))
		return
	}
}