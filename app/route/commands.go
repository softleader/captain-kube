package route

import (
	"github.com/kataras/iris"
	"github.com/softleader/captain-kube/ansible/playbook"
	"encoding/json"
	"github.com/softleader/captain-kube/pipe"
	"github.com/softleader/captain-kube/sh"
	"path"
	"github.com/softleader/captain-kube/ansible"
)

func Commands(workdir, playbooks string, ctx iris.Context) {
	book := playbook.MewCommands()
	body := ctx.GetHeader("Captain-Kube")
	err := json.Unmarshal([]byte(body), &book)
	if err != nil {
		ctx.StreamWriter(pipe.Println(err.Error()))
		return
	}
	opts := sh.Options{
		Ctx:     &ctx,
		Pwd:     playbooks,
		Verbose: book.V() != "",
	}
	book.Inventory = path.Join(workdir, book.Inventory)
	_, _, err = ansible.Play(&opts, *book)
	if err != nil {
		ctx.StreamWriter(pipe.Println(err.Error()))
		return
	}
}
