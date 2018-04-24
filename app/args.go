package app

import (
	"flag"
	"fmt"
	"encoding/json"
)

type Args struct {
	Playbooks string
	Workdir   string
	Addr      string
	Port      int
}

func NewArgs() *Args {
	a := Args{}
	flag.StringVar(&a.Playbooks, "playbooks", "", "Docker 中 ansible playbooks 的目錄")
	flag.StringVar(&a.Workdir, "workdir", "", "Docker 中 mount 出去的目錄, 通常會放客戶端的 hosts")
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
