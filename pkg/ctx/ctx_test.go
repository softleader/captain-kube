package ctx

import (
	"bufio"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/captain"
	"github.com/softleader/captain-kube/pkg/utils"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestLoadContexts(t *testing.T) {
	tmp, err := ioutil.TempDir(os.TempDir(), "ctx")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(tmp)

	b := bufio.NewWriter(nil)
	log := logrus.New()
	log.SetFormatter(&utils.PlainFormatter{})
	logrus.SetOutput(b)

	ctxFile := filepath.Join(tmp, "test-ctx.yaml")

	ctx, err := LoadContexts(log, ctxFile)
	if err != nil {
		t.Error(err)
	}

	if err = ctx.Add("foo", []string{"-e", "localhost"}); err != nil {
		t.Error(err)
	}
	actual, err := LoadContexts(log, ctxFile)
	if err != nil {
		t.Error(err)
	}
	if len(actual.Contexts) != 1 {
		t.Errorf("suppose to have 1 context")
		t.FailNow()
	}
	if foo, found := actual.Contexts["foo"]; !found {
		t.Errorf("suppose to have foo key")
		t.FailNow()
	} else if foo.Endpoint == nil {
		t.Errorf("suppose to have non nil endpoint of foo")
		t.FailNow()
	} else if host := foo.Endpoint.Host; host != "localhost" {
		t.Errorf("suppose to have localhost endpoint.host of foo, but got %s", host)
		t.FailNow()
	}

	if err = ctx.Add("bar", []string{"--endpoint-port", "9876", "--endpoint", "192.168.1.93"}); err != nil {
		t.Error(err)
	}
	actual, err = LoadContexts(log, ctxFile)
	if err != nil {
		t.Error(err)
	}
	if len(actual.Contexts) != 2 {
		t.Errorf("suppose to have 2 context")
		t.FailNow()
	}
	if bar, found := actual.Contexts["bar"]; !found {
		t.Errorf("suppose to have bar key")
		t.FailNow()
	} else if bar.Endpoint == nil {
		t.Errorf("suppose to have non nil endpoint of bar")
		t.FailNow()
	} else if host := bar.Endpoint.Host; host != "192.168.1.93" {
		t.Errorf("suppose to have 192.168.1.93 endpoint.host of bar, but got %s", host)
		t.FailNow()
	} else if port := bar.Endpoint.Port; port != 9876 {
		t.Errorf("suppose to have 9876 endpoint.port of bar, but got %v", port)
		t.FailNow()
	}

	if err := ctx.Switch("bar"); err != nil {
		t.Error(err)
		t.FailNow()
	}
	actual, err = LoadContexts(log, ctxFile)
	if err != nil {
		t.Error(err)
	}
	if a := actual.Active; a != "bar" {
		t.Errorf("active context should be bar, bot get %s", a)
		t.FailNow()
	}
	if p := actual.Previous; p != "" {
		t.Errorf("previous context should be empty, bot get %s", p)
		t.FailNow()
	}

	if err := ctx.Switch("foo"); err != nil {
		t.Error(err)
		t.FailNow()
	}
	actual, err = LoadContexts(log, ctxFile)
	if err != nil {
		t.Error(err)
	}
	if a := actual.Active; a != "foo" {
		t.Errorf("active context should be foo, bot get %s", a)
		t.FailNow()
	}
	if p := actual.Previous; p != "bar" {
		t.Errorf("previous context should be bar, bot get %s", p)
		t.FailNow()
	}

	if err := ctx.Switch("-"); err != nil {
		t.Error(err)
		t.FailNow()
	}
	actual, err = LoadContexts(log, ctxFile)
	if err != nil {
		t.Error(err)
	}
	if a := actual.Active; a != "bar" {
		t.Errorf("active context should be bar, bot get %s", a)
		t.FailNow()
	}
	if p := actual.Previous; p != "foo" {
		t.Errorf("previous context should be foo, bot get %s", p)
		t.FailNow()
	}

	if err := ctx.Delete("bar"); err != nil {
		t.Error(err)
		t.FailNow()
	}
	actual, err = LoadContexts(log, ctxFile)
	if err != nil {
		t.Error(err)
	}
	if len(actual.Contexts) != 1 {
		t.Errorf("suppose to have 1 context")
		t.FailNow()
	}
	if a := actual.Active; a != "" {
		t.Errorf("active context should be empty, bot get %s", a)
		t.FailNow()
	}
	if p := actual.Previous; p != "foo" {
		t.Errorf("previous context should be foo, bot get %s", p)
		t.FailNow()
	}
}

func TestExpandEnv(t *testing.T) {
	ctx := newContext(false)
	ctx.HelmTiller.Endpoint = "192.168.1.93"
	ctx.HelmTiller.Account = "hello-tiller"
	ctx.HelmTiller.Password = "secret"
	os.Setenv(captain.DefaultTillerPassword, "surprised")
	os.Setenv(captain.EnvPort, strconv.Itoa(captain.DefaultPort))

	if err := ctx.expandEnv(); err != nil {
		t.Error(err)
	}
	if e := ctx.HelmTiller.Endpoint; e != "192.168.1.93" {
		t.Errorf("new endpoint should be 192.168.1.93, but got %s", e)
	}
	if a := ctx.HelmTiller.Account; a != "hello-tiller" {
		t.Errorf("new account should be hello-tiller, but got %s", a)
	}
	if p := ctx.HelmTiller.Password; p != "secret" {
		t.Errorf("new endpoint should be secret, but got %s", p)
	}
	if p := ctx.Endpoint.Port; p != captain.DefaultPort {
		t.Errorf("new endpoint port should be %v, but got %v", captain.DefaultPort, p)
	}
}
