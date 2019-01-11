package ctx

import (
	"bufio"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/utils"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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

	if err = ctx.Add("foo", []string{"-e", "localhost"}, false); err != nil {
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
	if args, found := actual.Contexts["foo"]; !found {
		t.Errorf("suppose to have foo key")
		t.FailNow()
	} else {
		foo, err := NewContext(args...)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		if host := foo.Endpoint.Host; host != "localhost" {
			t.Errorf("suppose to have localhost endpoint.host of foo, but got %s", host)
			t.FailNow()
		}
	}
	if err = ctx.Add("foo", []string{"-e", "localhost"}, false); err == nil {
		t.Error("should have error")
	} else if !strings.HasSuffix(err.Error(), "already exists") {
		t.Error("should have error suffix: already exists")
	}
	if err = ctx.Add("foo", []string{"-e", "localhost"}, true); err != nil {
		t.Error(err)
	}
	actual, err = LoadContexts(log, ctxFile)
	if err != nil {
		t.Error(err)
	}
	if len(actual.Contexts) != 1 {
		t.Errorf("suppose to have 1 context")
		t.FailNow()
	}

	if err = ctx.Add("bar", []string{"--endpoint-port", "9876", "--endpoint", "192.168.1.93"}, false); err != nil {
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
	if args, found := actual.Contexts["bar"]; !found {
		t.Errorf("suppose to have bar key")
		t.FailNow()
	} else {
		bar, err := NewContext(args...)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		if host := bar.Endpoint.Host; host != "192.168.1.93" {
			t.Errorf("suppose to have 192.168.1.93 endpoint.host of bar, but got %s", host)
			t.FailNow()
		} else if port := bar.Endpoint.Port; port != 9876 {
			t.Errorf("suppose to have 9876 endpoint.port of bar, but got %v", port)
			t.FailNow()
		}
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
	if len(actual.Contexts) != 2 {
		t.Errorf("suppose to have 2 context")
		t.FailNow()
	}

	if err := ctx.Rename("bar", "bar-bee-que"); err != nil {
		t.Error(err)
		t.FailNow()
	}
	actual, err = LoadContexts(log, ctxFile)
	if err != nil {
		t.Error(err)
	}
	if a := actual.Active; a != "bar-bee-que" {
		t.Errorf("active context should be bar-bee-que, bot get %s", a)
		t.FailNow()
	}
	if p := actual.Previous; p != "foo" {
		t.Errorf("previous context should be foo, bot get %s", p)
		t.FailNow()
	}

	if err := ctx.Delete("bar-bee-que"); err != nil {
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
