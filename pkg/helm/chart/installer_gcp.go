package chart

import "fmt"

type gcpInstaller struct {
	endpoint, chart string
}

// TODO: 因為只有 Google 可以將資料放在台灣, 客戶多半會優先選擇 GCP, 因此在未來將會需要優先整合
func (i *gcpInstaller) Install() error {
	return fmt.Errorf("GCP is not supported yet")
}
