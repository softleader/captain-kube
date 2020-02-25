package ctx

import (
	"bytes"
	"github.com/imdario/mergo"
	"github.com/spf13/cobra"
)

// Context 封裝了很多常用的 flags
type Context struct {
	Global       *Global
	Endpoint     *Endpoint
	HelmTiller   *HelmTiller
	HelmChart    *HelmChart
	RegistryAuth *RegistryAuth
	ReTag        *ReTag
}

// NewContextFromEnv 建立 Context 建議並試著從 OS Env 中取得預設值
func NewContextFromEnv() (c *Context) {
	c = &Context{
		Global:       newGlobalFromEnv(),
		Endpoint:     newEndpointFromEnv(),
		HelmTiller:   newHelmTillerFromEnv(),
		HelmChart:    newHelmChartFromEnv(),
		RegistryAuth: newRegistryAuthFromEnv(),
		ReTag:        newReTagFromEnv(),
	}
	return
}

// FlagsString 返回所有 Context 中支援的 flag 說明文字
func FlagsString() (string, error) {
	c := &Context{
		Global:       &Global{},
		Endpoint:     &Endpoint{},
		HelmTiller:   &HelmTiller{},
		HelmChart:    &HelmChart{},
		RegistryAuth: &RegistryAuth{},
		ReTag:        &ReTag{},
	}
	cmd := &cobra.Command{}
	cmd.SetUsageTemplate(`{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}`)
	f := cmd.Flags()
	addFlags(c, f)
	bb := new(bytes.Buffer)
	cmd.SetOutput(bb)
	if err := cmd.Usage(); err != nil {
		return "", err
	}
	return bb.String(), nil
}

// NewContext 依照傳入的 args 建立 Context 物件
func NewContext(args ...string) (*Context, error) {
	c := &Context{
		Global:       &Global{},
		Endpoint:     &Endpoint{},
		HelmTiller:   &HelmTiller{},
		HelmChart:    &HelmChart{},
		RegistryAuth: &RegistryAuth{},
		ReTag:        &ReTag{},
	}
	if len(args) == 0 {
		return c, nil
	}
	cmd := &cobra.Command{}
	f := cmd.Flags()
	addFlags(c, f)
	return c, cmd.ParseFlags(args)
}

// ExpandEnv Merge OS Env 到當前的 Context
func (ctx *Context) ExpandEnv() error {
	defaultCtx := NewContextFromEnv()
	return mergo.Merge(ctx, defaultCtx)
}
