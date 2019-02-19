package server

import (
	"encoding/hex"
	"fmt"
	"github.com/softleader/captain-kube/pkg/proto"
	"io/ioutil"
)

func saveChart(chart *captainkube_v2.Chart, path string) error {
	body := chart.GetContent()
	if len(body) == 0 {
		hexadecimal := chart.GetContentHex()
		if len(hexadecimal) == 0 {
			return fmt.Errorf("chart is required, but got %+v", chart)
		}
		decode, err := hex.DecodeString(hexadecimal)
		if err != nil {
			return err
		}
		body = decode
	}
	return ioutil.WriteFile(path, body, 0644)
}
