package chart

import "fmt"

type pksInstaller struct {
	endpoint, chart string
}

func (i *pksInstaller) Install() error {
	return fmt.Errorf("PKS is not supported yet")
}
