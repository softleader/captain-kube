package server

import (
	"encoding/hex"
	"github.com/softleader/captain-kube/pkg/proto"
	"io/ioutil"
)

func saveChart(chart *captainkube_v2.Chart, path string) error {
	body := chart.GetContent()
	if len(body) == 0 {
		decode, err := hex.DecodeString(chart.GetContentHex())
		if err != nil {
			return err
		}
		body = decode
	}
	return ioutil.WriteFile(path, body, 0644)
}
