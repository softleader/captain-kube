package chart

import "fmt"

type openShiftInstaller struct {
	endpoint, chart string
}

func (i *openShiftInstaller) Install() error {
	return fmt.Errorf("OpenShift is not supported yet")
}
