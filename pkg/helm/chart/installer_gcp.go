package chart

import "fmt"

type gcpInstaller struct {
	chart string
}

func (i *gcpInstaller) Install() error {
	return fmt.Errorf("GCP is not supported yet")
}
