package app

import (
	"flag"
	"fmt"
	"encoding/json"
)

type Args struct {
	Ansible       string
	Workspace     string
	HostWorkspace string
	Addr          string
	Port          int
}

func NewArgs() *Args {
	a := Args{}
	flag.StringVar(&a.Ansible, "ansible", "", "Docker 中 ansible playbook 的位置")
	flag.StringVar(&a.Workspace, "workspace", "", "Docker 中 mount 出去的位置")
	flag.StringVar(&a.HostWorkspace, "host-workspace", "", "Host 中 mount 給 docker 的位置")
	flag.StringVar(&a.Addr, "addr", "", " Determine application address (default blank)")
	flag.IntVar(&a.Port, "port", 10080, "Determine application port")
	flag.Parse()

	marshaled, err := json.MarshalIndent(a, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Println("Parsed flags:", string(marshaled))

	return &a
}
