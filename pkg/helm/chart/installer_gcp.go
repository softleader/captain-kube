package chart

import "fmt"

type gcpInstaller struct {
	endpoint, chart string
}

func (i *gcpInstaller) Install() error {
	return fmt.Errorf("GCP is not supported yet")
}
