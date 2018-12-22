package chart

import (
	"fmt"
	"github.com/softleader/captain-kube/pkg/logger"
)

type gcpInstaller struct {
	endpoint, chart string
}

// TODO: 因為只有 Google 可以將資料放在台灣, 客戶多半會優先選擇 GCP, 因此在未來將會需要優先整合
func (i *gcpInstaller) Install(log *logger.Logger) error {
	return fmt.Errorf("GCP is not supported yet")
}
