package sh

import (
	"strings"
	"os/exec"
	"fmt"
	"github.com/kataras/iris"
	"github.com/softleader/captain-kube/pipe"
	"bytes"
	"errors"
)

type Options struct {
	Ctx     *iris.Context
	Pwd     string
	Verbose bool
}

type output struct {
	ctx     *iris.Context
	buf     bytes.Buffer
	Verbose bool
}

func C(opts *Options, commands ...string) (arg string, out string, err error) {
	arg = strings.Join(commands, " ")
	cmd := exec.Command("sh", "-c", arg)
	if opts.Pwd != "" {
		cmd.Dir = opts.Pwd
	}

	if opts.Verbose {
		if opts.Ctx != nil {
			(*opts.Ctx).StreamWriter(pipe.Printf("$ %v\n", arg))
		} else {
			fmt.Printf("$ %v\n", arg)
		}
	}

	stdout := output{ctx: opts.Ctx, Verbose: opts.Verbose}
	stderr := output{ctx: opts.Ctx, Verbose: opts.Verbose}
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return "", "", errors.New(fmt.Sprint(err) + ": " + stderr.buf.String())
	}
	if stderr.buf.Len() > 0 {
		return "", "", errors.New(stderr.buf.String())
	}

	return arg, stdout.buf.String(), nil
}

func (o *output) Write(b []byte) (n int, err error) {
	o.buf.Write(b)
	if o.ctx != nil {
		s := string(b)
		if o.Verbose {
			fmt.Print(s)
		}
		(*o.ctx).StreamWriter(pipe.Print(s))
	}
	return len(b), nil
}
