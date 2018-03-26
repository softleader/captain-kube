package main

import (
	"testing"
	"github.com/spf13/cobra/doc"
	"github.com/softleader/captain-kube/app"
)

func TestGenerateMarkdownForCommandTree(t *testing.T) {
	ckube := app.NewCaptainKubeCommand()
	err := doc.GenMarkdownTree(ckube, "./docs/man")
	if err != nil {
		t.Error(err)
	}
}
